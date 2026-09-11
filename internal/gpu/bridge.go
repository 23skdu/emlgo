//go:build cuda && cgo

package gpu

/*
#cgo LDFLAGS: -L${SRCDIR}/../../cuda -Wl,-rpath,${SRCDIR}/../../cuda -leml_capi -lcuda -lcudart
#cgo CFLAGS: -I${SRCDIR}/../../cuda
#include "eml_capi.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

var initialized bool

type deviceBuffer struct {
	ptr  unsafe.Pointer
	size int64
}

var (
	poolMu sync.Mutex
	pool   []deviceBuffer
)

func allocDeviceBuffer(size int64) (unsafe.Pointer, error) {
	poolMu.Lock()
	for i, b := range pool {
		if b.size >= size {
			ptr := b.ptr
			pool[i] = pool[len(pool)-1]
			pool = pool[:len(pool)-1]
			poolMu.Unlock()
			return ptr, nil
		}
	}
	poolMu.Unlock()

	var ptr unsafe.Pointer
	if err := C.eml_allocate(&ptr, C.longlong(size)); err != 0 {
		return nil, fmt.Errorf("device memory allocation failed: error %d", int(err))
	}
	return ptr, nil
}

func freeDeviceBuffer(ptr unsafe.Pointer, size int64) {
	poolMu.Lock()
	pool = append(pool, deviceBuffer{ptr: ptr, size: size})
	poolMu.Unlock()
}

func init() {
	if C.eml_init() == 0 {
		initialized = true
	}
}

// Init initializes the GPU subsystem.
func Init() error {
	if C.eml_init() == 0 {
		initialized = true
		return nil
	}
	return fmt.Errorf("failed to initialize CUDA")
}

// Shutdown cleans up the GPU subsystem.
func Shutdown() {
	if initialized {
		poolMu.Lock()
		for _, b := range pool {
			C.eml_free(b.ptr)
		}
		pool = nil
		poolMu.Unlock()
		C.eml_cleanup()
		initialized = false
	}
}

// GetDevices returns a list of available GPU devices.
func GetDevices() ([]Device, error) {
	var count C.int
	err := C.eml_get_device_count(&count)
	if err != 0 {
		return nil, fmt.Errorf("cudaGetDeviceCount failed: error %d", int(err))
	}
	if count == 0 {
		return nil, nil
	}

	devices := make([]Device, int(count))
	for i := 0; i < int(count); i++ {
		d, err := getDeviceProps(i)
		if err != nil {
			return nil, err
		}
		devices[i] = d
	}
	return devices, nil
}

func getDeviceProps(id int) (Device, error) {
	var (
		computeMajor, computeMinor      C.int
		memoryBytes                     C.longlong
		maxThreads, warpSize, clockRate C.int
	)
	err := C.eml_get_device_props(
		C.int(id),
		&computeMajor, &computeMinor,
		&memoryBytes,
		&maxThreads, &warpSize, &clockRate,
	)
	if err != 0 {
		return Device{}, fmt.Errorf("failed to get device %d properties: error %d", id, int(err))
	}

	nameBuf := make([]C.char, 256)
	err = C.eml_get_device_name(C.int(id), &nameBuf[0], 256)
	if err != 0 {
		return Device{}, fmt.Errorf("failed to get device %d name: error %d", id, int(err))
	}

	return Device{
		ID:                 id,
		Name:               C.GoString(&nameBuf[0]),
		ComputeMajor:       int(computeMajor),
		ComputeMinor:       int(computeMinor),
		MemoryBytes:        int64(memoryBytes),
		MaxThreadsPerBlock: int(maxThreads),
		WarpSize:           int(warpSize),
		ClockRateKHz:       int(clockRate),
	}, nil
}

// Status returns a human-readable status of the GPU subsystem.
func Status() string {
	if !initialized {
		if C.eml_init() == 0 {
			initialized = true
		} else {
			return "CUDA initialization failed"
		}
	}

	devices, err := GetDevices()
	if err != nil {
		return fmt.Sprintf("GPU Error: %v", err)
	}
	if len(devices) == 0 {
		return "No GPU devices found"
	}

	s := fmt.Sprintf("Found %d GPU device(s):\n", len(devices))
	for _, d := range devices {
		s += fmt.Sprintf("  [%d] %s (SM %d.%d, %d MB, %d threads/block, %d MHz)\n",
			d.ID, d.Name, d.ComputeMajor, d.ComputeMinor,
			d.MemoryBytes/1024/1024, d.MaxThreadsPerBlock, d.ClockRateKHz/1000)
	}
	return s
}

// ---------- Generic Batch Launchers ----------

func launchUnary[T any](x []T, elemSize int, fn func(dResult, dX unsafe.Pointer) C.int) ([]T, error) {
	if len(x) == 0 {
		return nil, nil
	}
	result := make([]T, len(x))
	if err := launchUnaryTo(result, x, elemSize, fn); err != nil {
		return nil, err
	}
	return result, nil
}

func launchUnaryTo[T any](dst, x []T, elemSize int, fn func(dResult, dX unsafe.Pointer) C.int) error {
	if !initialized {
		return fmt.Errorf("CUDA not initialized")
	}
	if len(x) == 0 {
		return nil
	}
	if len(dst) != len(x) {
		return fmt.Errorf("slice length mismatch: dst %d != x %d", len(dst), len(x))
	}

	n := len(x)
	size := C.longlong(n * elemSize)

	// Optimization: if in-place, use a single device buffer
	if &dst[0] == &x[0] {
		d_buf, allocErr := allocDeviceBuffer(int64(size))
		if allocErr != nil {
			return allocErr
		}
		defer freeDeviceBuffer(d_buf, int64(size))

		if err := C.eml_copy_to_device(d_buf, unsafe.Pointer(&x[0]), size); err != 0 {
			return fmt.Errorf("copy to device failed: error %d", int(err))
		}
		if err := fn(d_buf, d_buf); err != 0 {
			return fmt.Errorf("kernel launch failed: error %d", int(err))
		}
		if err := C.eml_sync_device(); err != 0 {
			return fmt.Errorf("device sync failed: error %d", int(err))
		}
		if err := C.eml_copy_to_host(unsafe.Pointer(&dst[0]), d_buf, size); err != 0 {
			return fmt.Errorf("copy to host failed: error %d", int(err))
		}
		return nil
	}

	d_x, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return allocErr
	}
	defer freeDeviceBuffer(d_x, int64(size))

	d_result, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return allocErr
	}
	defer freeDeviceBuffer(d_result, int64(size))

	if err := C.eml_copy_to_device(d_x, unsafe.Pointer(&x[0]), size); err != 0 {
		return fmt.Errorf("copy to device failed: error %d", int(err))
	}
	if err := fn(d_result, d_x); err != 0 {
		return fmt.Errorf("kernel launch failed: error %d", int(err))
	}
	if err := C.eml_sync_device(); err != 0 {
		return fmt.Errorf("device sync failed: error %d", int(err))
	}
	if err := C.eml_copy_to_host(unsafe.Pointer(&dst[0]), d_result, size); err != 0 {
		return fmt.Errorf("copy to host failed: error %d", int(err))
	}

	return nil
}

func launchBinary[T any](x, y []T, elemSize int, fn func(dResult, dX, dY unsafe.Pointer) C.int) ([]T, error) {
	if len(x) != len(y) {
		return nil, fmt.Errorf("x and y must have same length: %d vs %d", len(x), len(y))
	}
	if len(x) == 0 {
		return nil, nil
	}
	result := make([]T, len(x))
	if err := launchBinaryTo(result, x, y, elemSize, fn); err != nil {
		return nil, err
	}
	return result, nil
}

func launchBinaryTo[T any](dst, x, y []T, elemSize int, fn func(dResult, dX, dY unsafe.Pointer) C.int) error {
	if !initialized {
		return fmt.Errorf("CUDA not initialized")
	}
	if len(x) != len(y) || len(dst) != len(x) {
		return fmt.Errorf("slice length mismatch: dst=%d, x=%d, y=%d", len(dst), len(x), len(y))
	}
	if len(x) == 0 {
		return nil
	}

	n := len(x)
	size := C.longlong(n * elemSize)

	d_x, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return allocErr
	}
	defer freeDeviceBuffer(d_x, int64(size))

	d_y, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return allocErr
	}
	defer freeDeviceBuffer(d_y, int64(size))

	d_result, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return allocErr
	}
	defer freeDeviceBuffer(d_result, int64(size))

	if err := C.eml_copy_to_device(d_x, unsafe.Pointer(&x[0]), size); err != 0 {
		return fmt.Errorf("copy x to device failed: error %d", int(err))
	}
	if err := C.eml_copy_to_device(d_y, unsafe.Pointer(&y[0]), size); err != 0 {
		return fmt.Errorf("copy y to device failed: error %d", int(err))
	}
	if err := fn(d_result, d_x, d_y); err != 0 {
		return fmt.Errorf("kernel launch failed: error %d", int(err))
	}
	if err := C.eml_sync_device(); err != 0 {
		return fmt.Errorf("device sync failed: error %d", int(err))
	}
	if err := C.eml_copy_to_host(unsafe.Pointer(&dst[0]), d_result, size); err != 0 {
		return fmt.Errorf("copy to host failed: error %d", int(err))
	}

	return nil
}

func launchDot[T any](a, b []T, elemSize int, fn func(dA, dB, dResult unsafe.Pointer) C.int) (T, error) {
	var zero T
	if !initialized {
		return zero, fmt.Errorf("CUDA not initialized")
	}
	if len(a) != len(b) {
		return zero, fmt.Errorf("slice length mismatch: %d vs %d", len(a), len(b))
	}
	if len(a) == 0 {
		return zero, nil
	}

	n := len(a)
	size := C.longlong(n * elemSize)

	d_a, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return zero, allocErr
	}
	defer freeDeviceBuffer(d_a, int64(size))

	d_b, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return zero, allocErr
	}
	defer freeDeviceBuffer(d_b, int64(size))

	d_res, allocErr := allocDeviceBuffer(int64(elemSize))
	if allocErr != nil {
		return zero, allocErr
	}
	defer freeDeviceBuffer(d_res, int64(elemSize))

	if err := C.eml_copy_to_device(d_a, unsafe.Pointer(&a[0]), size); err != 0 {
		return zero, fmt.Errorf("copy a failed: error %d", int(err))
	}
	if err := C.eml_copy_to_device(d_b, unsafe.Pointer(&b[0]), size); err != 0 {
		return zero, fmt.Errorf("copy b failed: error %d", int(err))
	}

	if err := fn(d_a, d_b, d_res); err != 0 {
		return zero, fmt.Errorf("dot kernel launch failed: error %d", int(err))
	}
	if err := C.eml_sync_device(); err != 0 {
		return zero, fmt.Errorf("device sync failed: error %d", int(err))
	}

	var res T
	if err := C.eml_copy_to_host(unsafe.Pointer(&res), d_res, C.longlong(elemSize)); err != 0 {
		return zero, fmt.Errorf("copy result failed: error %d", int(err))
	}
	return res, nil
}

// ============================================================================
// Float64 Device Methods
// ============================================================================

func (d *Device) ExpBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) ExpBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) LogBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) LogBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CosBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CosBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinhBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinhBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CoshBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CoshBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanhBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanhBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SqrtBatch(x []float64) ([]float64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SqrtBatchTo(dst, x []float64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) EmlBatch(x, y []float64) ([]float64, error) {
	return launchBinary(x, y, 8, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) EmlBatchTo(dst, x, y []float64) error {
	return launchBinaryTo(dst, x, y, 8, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) DotBatch(a, b []float64) (float64, error) {
	return launchDot(a, b, 8, func(dA, dB, dRes unsafe.Pointer) C.int {
		return C.eml_launch_dot(dA, dB, dRes, C.int(len(a)), C.int(DefaultBlockSize))
	})
}

// ============================================================================
// Float32 Device Methods
// ============================================================================

func (d *Device) ExpBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) ExpBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) LogBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) LogBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CosBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CosBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinhBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinhBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CoshBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CoshBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanhBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanhBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SqrtBatchF32(x []float32) ([]float32, error) {
	return launchUnary(x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SqrtBatchToF32(dst, x []float32) error {
	return launchUnaryTo(dst, x, 4, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) EmlBatchF32(x, y []float32) ([]float32, error) {
	return launchBinary(x, y, 4, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml_f32(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) EmlBatchToF32(dst, x, y []float32) error {
	return launchBinaryTo(dst, x, y, 4, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml_f32(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) DotBatchF32(a, b []float32) (float32, error) {
	return launchDot(a, b, 4, func(dA, dB, dRes unsafe.Pointer) C.int {
		return C.eml_launch_dot_f32(dA, dB, dRes, C.int(len(a)), C.int(DefaultBlockSize))
	})
}

// ============================================================================
// Complex64 Device Methods
// ============================================================================

func (d *Device) ExpBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) ExpBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) LogBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) LogBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CosBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CosBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinhBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinhBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CoshBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CoshBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanhBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanhBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SqrtBatchC64(x []complex64) ([]complex64, error) {
	return launchUnary(x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SqrtBatchToC64(dst, x []complex64) error {
	return launchUnaryTo(dst, x, 8, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) EmlBatchC64(x, y []complex64) ([]complex64, error) {
	return launchBinary(x, y, 8, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml_c64(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) EmlBatchToC64(dst, x, y []complex64) error {
	return launchBinaryTo(dst, x, y, 8, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml_c64(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) DotBatchC64(a, b []complex64) (complex64, error) {
	return launchDot(a, b, 8, func(dA, dB, dRes unsafe.Pointer) C.int {
		return C.eml_launch_dot_c64(dA, dB, dRes, C.int(len(a)), C.int(DefaultBlockSize))
	})
}

// ============================================================================
// Complex128 Device Methods
// ============================================================================

func (d *Device) ExpBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) ExpBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_exp_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) LogBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) LogBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_log_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sin_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CosBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CosBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cos_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tan_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SinhBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SinhBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sinh_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) CoshBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) CoshBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_cosh_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) TanhBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) TanhBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_tanh_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) SqrtBatchC128(x []complex128) ([]complex128, error) {
	return launchUnary(x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) SqrtBatchToC128(dst, x []complex128) error {
	return launchUnaryTo(dst, x, 16, func(dRes, dX unsafe.Pointer) C.int {
		return C.eml_launch_sqrt_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) EmlBatchC128(x, y []complex128) ([]complex128, error) {
	return launchBinary(x, y, 16, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml_c128(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}
func (d *Device) EmlBatchToC128(dst, x, y []complex128) error {
	return launchBinaryTo(dst, x, y, 16, func(dRes, dX, dY unsafe.Pointer) C.int {
		return C.eml_launch_eml_c128(dX, dY, dRes, C.int(len(x)), C.int(DefaultBlockSize))
	})
}

func (d *Device) DotBatchC128(a, b []complex128) (complex128, error) {
	return launchDot(a, b, 16, func(dA, dB, dRes unsafe.Pointer) C.int {
		return C.eml_launch_dot_c128(dA, dB, dRes, C.int(len(a)), C.int(DefaultBlockSize))
	})
}

// ============================================================================
// Async Streams
// ============================================================================

// NewStream creates a new asynchronous stream.
func NewStream() (*Stream, error) {
	h := C.eml_create_stream()
	if h == 0 {
		return nil, fmt.Errorf("failed to create CUDA stream")
	}
	return &Stream{handle: uintptr(h)}, nil
}

// Sync synchronizes the stream, blocking until all operations complete.
func (s *Stream) Sync() error {
	if s == nil || s.handle == 0 {
		return nil
	}
	if err := C.eml_sync_stream(C.longlong(s.handle)); err != 0 {
		return fmt.Errorf("stream sync failed: error %d", int(err))
	}
	return nil
}

// Destroy releases the stream resources.
func (s *Stream) Destroy() error {
	if s == nil || s.handle == 0 {
		return nil
	}
	err := C.eml_destroy_stream(C.longlong(s.handle))
	s.handle = 0
	if err != 0 {
		return fmt.Errorf("stream destroy failed: error %d", int(err))
	}
	return nil
}

// launchUnaryStream launches a unary kernel on an async stream.
func launchUnaryStream[T any](x []T, elemSize int, stream *Stream, fn func(dRes, dX unsafe.Pointer, s C.longlong) C.int) ([]T, error) {
	if !initialized {
		return nil, fmt.Errorf("CUDA not initialized")
	}
	if len(x) == 0 {
		return nil, nil
	}

	n := len(x)
	size := C.longlong(n * elemSize)

	d_x, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return nil, allocErr
	}
	defer freeDeviceBuffer(d_x, int64(size))

	d_result, allocErr := allocDeviceBuffer(int64(size))
	if allocErr != nil {
		return nil, allocErr
	}
	defer freeDeviceBuffer(d_result, int64(size))

	if err := C.eml_copy_to_device(d_x, unsafe.Pointer(&x[0]), size); err != 0 {
		return nil, fmt.Errorf("copy to device failed: error %d", int(err))
	}

	var streamHandle C.longlong
	if stream != nil {
		streamHandle = C.longlong(stream.handle)
	}

	if err := fn(d_result, d_x, streamHandle); err != 0 {
		return nil, fmt.Errorf("kernel launch failed: error %d", int(err))
	}

	if err := C.eml_sync_stream(streamHandle); err != 0 {
		return nil, fmt.Errorf("stream sync failed: error %d", int(err))
	}

	result := make([]T, n)
	if err := C.eml_copy_to_host(unsafe.Pointer(&result[0]), d_result, size); err != 0 {
		return nil, fmt.Errorf("copy to host failed: error %d", int(err))
	}

	return result, nil
}

func (d *Device) ExpBatchStream(x []float64, s *Stream) ([]float64, error) {
	return launchUnaryStream(x, 8, s, func(dRes, dX unsafe.Pointer, sh C.longlong) C.int {
		return C.eml_launch_exp_stream(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize), sh)
	})
}

func (d *Device) ExpBatchStreamF32(x []float32, s *Stream) ([]float32, error) {
	return launchUnaryStream(x, 4, s, func(dRes, dX unsafe.Pointer, sh C.longlong) C.int {
		return C.eml_launch_exp_stream_f32(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize), sh)
	})
}

func (d *Device) ExpBatchStreamC64(x []complex64, s *Stream) ([]complex64, error) {
	return launchUnaryStream(x, 8, s, func(dRes, dX unsafe.Pointer, sh C.longlong) C.int {
		return C.eml_launch_exp_stream_c64(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize), sh)
	})
}

func (d *Device) ExpBatchStreamC128(x []complex128, s *Stream) ([]complex128, error) {
	return launchUnaryStream(x, 16, s, func(dRes, dX unsafe.Pointer, sh C.longlong) C.int {
		return C.eml_launch_exp_stream_c128(dX, dRes, C.int(len(x)), C.int(DefaultBlockSize), sh)
	})
}

// ============================================================================
// Pinned Host Memory
// ============================================================================

func AllocatePinned(size int) ([]float64, error) {
	if !initialized || size <= 0 {
		return make([]float64, size), nil
	}
	var ptr unsafe.Pointer
	byteSize := C.longlong(size * 8)
	if err := C.eml_allocate_pinned(&ptr, byteSize); err != 0 {
		return nil, fmt.Errorf("pinned memory allocation failed: error %d", int(err))
	}
	data := (*[1 << 30]float64)(ptr)[:size:size]
	return data, nil
}

func FreePinned(p []float64) {
	if len(p) == 0 {
		return
	}
	ptr := unsafe.Pointer(&p[0])
	C.eml_free_pinned(ptr)
}

func AllocatePinnedF32(size int) ([]float32, error) {
	if !initialized || size <= 0 {
		return make([]float32, size), nil
	}
	var ptr unsafe.Pointer
	byteSize := C.longlong(size * 4)
	if err := C.eml_allocate_pinned(&ptr, byteSize); err != 0 {
		return nil, fmt.Errorf("pinned memory allocation failed: error %d", int(err))
	}
	data := (*[1 << 30]float32)(ptr)[:size:size]
	return data, nil
}

func FreePinnedF32(p []float32) {
	if len(p) == 0 {
		return
	}
	ptr := unsafe.Pointer(&p[0])
	C.eml_free_pinned(ptr)
}

func AllocatePinnedC64(size int) ([]complex64, error) {
	if !initialized || size <= 0 {
		return make([]complex64, size), nil
	}
	var ptr unsafe.Pointer
	byteSize := C.longlong(size * 8)
	if err := C.eml_allocate_pinned(&ptr, byteSize); err != 0 {
		return nil, fmt.Errorf("pinned memory allocation failed: error %d", int(err))
	}
	data := (*[1 << 30]complex64)(ptr)[:size:size]
	return data, nil
}

func FreePinnedC64(p []complex64) {
	if len(p) == 0 {
		return
	}
	ptr := unsafe.Pointer(&p[0])
	C.eml_free_pinned(ptr)
}

func AllocatePinnedC128(size int) ([]complex128, error) {
	if !initialized || size <= 0 {
		return make([]complex128, size), nil
	}
	var ptr unsafe.Pointer
	byteSize := C.longlong(size * 16)
	if err := C.eml_allocate_pinned(&ptr, byteSize); err != 0 {
		return nil, fmt.Errorf("pinned memory allocation failed: error %d", int(err))
	}
	data := (*[1 << 30]complex128)(ptr)[:size:size]
	return data, nil
}

func FreePinnedC128(p []complex128) {
	if len(p) == 0 {
		return
	}
	ptr := unsafe.Pointer(&p[0])
	C.eml_free_pinned(ptr)
}

func (d *Device) AbsBatch(x []float64) ([]float64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not implemented for AbsBatch")
}

func (d *Device) NegBatch(x []float64) ([]float64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not implemented for NegBatch")
}

func (d *Device) PowBatch(x []float64, exponent float64) ([]float64, error) {
	_ = d
	_ = x
	_ = exponent
	return nil, fmt.Errorf("GPU execution not implemented for PowBatch")
}

func (d *Device) InvBatch(x []float64) ([]float64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not implemented for InvBatch")
}

func (d *Device) FmaBatch(a, b, c []float64) ([]float64, error) {
	_ = d
	_ = a
	_ = b
	_ = c
	return nil, fmt.Errorf("GPU execution not implemented for FmaBatch")
}
