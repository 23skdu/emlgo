//go:build cuda

package gpu

/*
#cgo LDFLAGS: -L${SRCDIR}/../../cuda -lemcl_capi -lcuda -lcudart
#cgo CFLAGS: -I${SRCDIR}/../../cuda
#include "eml_capi.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

var initialized bool

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
		computeMajor, computeMinor                                    C.int
		memoryBytes                                                    C.longlong
		maxThreads, warpSize, clockRate                                C.int
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

// launchBatch is the generic implementation for all batch kernel launches.
func launchBatch(x []float64, fn func(int, int) C.int) ([]float64, error) {
	if !initialized {
		return nil, fmt.Errorf("CUDA not initialized")
	}
	if len(x) == 0 {
		return nil, nil
	}

	n := len(x)
	var d_x, d_result unsafe.Pointer

	size := C.longlong(n * 8) // float64 = 8 bytes
	if err := C.eml_allocate(&d_x, size); err != 0 {
		return nil, fmt.Errorf("device memory allocation failed: error %d", int(err))
	}
	defer C.eml_free(d_x)

	if err := C.eml_allocate(&d_result, size); err != 0 {
		return nil, fmt.Errorf("device memory allocation failed: error %d", int(err))
	}
	defer C.eml_free(d_result)

	if err := C.eml_copy_to_device(d_x, unsafe.Pointer(&x[0]), size); err != 0 {
		return nil, fmt.Errorf("copy to device failed: error %d", int(err))
	}

	if err := fn(int(uintptr(d_result)), int(uintptr(d_x))); err != 0 {
		return nil, fmt.Errorf("kernel launch failed: error %d", int(err))
	}

	if err := C.eml_sync_device(); err != 0 {
		return nil, fmt.Errorf("device sync failed: error %d", int(err))
	}

	result := make([]float64, n)
	if err := C.eml_copy_to_host(unsafe.Pointer(&result[0]), d_result, size); err != 0 {
		return nil, fmt.Errorf("copy to host failed: error %d", int(err))
	}

	return result, nil
}

func (d *Device) ExpBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_exp(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) LogBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_log(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) SinBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_sin(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) CosBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_cos(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) TanBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_tan(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) SinhBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_sinh(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) CoshBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_cosh(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) TanhBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_tanh(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) SqrtBatch(x []float64) ([]float64, error) {
	return launchBatch(x, func(dResult, dX int) C.int {
		return C.eml_launch_sqrt(
			unsafe.Pointer(uintptr(dX)),
			unsafe.Pointer(uintptr(dResult)),
			C.int(len(x)), C.int(DefaultBlockSize),
		)
	})
}

func (d *Device) EmlBatch(x, y []float64) ([]float64, error) {
	if !initialized {
		return nil, fmt.Errorf("CUDA not initialized")
	}
	if len(x) != len(y) {
		return nil, fmt.Errorf("x and y must have same length: %d vs %d", len(x), len(y))
	}
	if len(x) == 0 {
		return nil, nil
	}

	n := len(x)
	var d_x, d_y, d_result unsafe.Pointer

	size := C.longlong(n * 8)
	for _, ptr := range []*unsafe.Pointer{&d_x, &d_y, &d_result} {
		if err := C.eml_allocate(ptr, size); err != 0 {
			return nil, fmt.Errorf("device memory allocation failed: error %d", int(err))
		}
	}
	defer C.eml_free(d_x)
	defer C.eml_free(d_y)
	defer C.eml_free(d_result)

	if err := C.eml_copy_to_device(d_x, unsafe.Pointer(&x[0]), size); err != 0 {
		return nil, fmt.Errorf("copy x to device failed: error %d", int(err))
	}
	if err := C.eml_copy_to_device(d_y, unsafe.Pointer(&y[0]), size); err != 0 {
		return nil, fmt.Errorf("copy y to device failed: error %d", int(err))
	}

	if err := C.eml_launch_eml(d_x, d_y, d_result, C.int(n), C.int(DefaultBlockSize)); err != 0 {
		return nil, fmt.Errorf("kernel launch failed: error %d", int(err))
	}

	if err := C.eml_sync_device(); err != 0 {
		return nil, fmt.Errorf("device sync failed: error %d", int(err))
	}

	result := make([]float64, n)
	if err := C.eml_copy_to_host(unsafe.Pointer(&result[0]), d_result, size); err != 0 {
		return nil, fmt.Errorf("copy to host failed: error %d", int(err))
	}

	return result, nil
}

// ---------- Async Streams ----------

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

// launchBatchStream is the streamed (async) version of launchBatch.
func launchBatchStream(x []float64, stream *Stream, fn func(int, int, uintptr) C.int) ([]float64, error) {
	if !initialized {
		return nil, fmt.Errorf("CUDA not initialized")
	}
	if len(x) == 0 {
		return nil, nil
	}

	n := len(x)
	var d_x, d_result unsafe.Pointer

	size := C.longlong(n * 8)
	for _, ptr := range []*unsafe.Pointer{&d_x, &d_result} {
		if err := C.eml_allocate(ptr, size); err != 0 {
			return nil, fmt.Errorf("device memory allocation failed: error %d", int(err))
		}
	}
	defer C.eml_free(d_x)
	defer C.eml_free(d_result)

	if err := C.eml_copy_to_device(d_x, unsafe.Pointer(&x[0]), size); err != 0 {
		return nil, fmt.Errorf("copy to device failed: error %d", int(err))
	}

	var streamHandle C.longlong
	if stream != nil {
		streamHandle = C.longlong(stream.handle)
	}

	if err := fn(int(uintptr(d_result)), int(uintptr(d_x)), uintptr(streamHandle)); err != 0 {
		return nil, fmt.Errorf("kernel launch failed: error %d", int(err))
	}

	if err := C.eml_sync_stream(streamHandle); err != 0 {
		return nil, fmt.Errorf("stream sync failed: error %d", int(err))
	}

	result := make([]float64, n)
	if err := C.eml_copy_to_host(unsafe.Pointer(&result[0]), d_result, size); err != 0 {
		return nil, fmt.Errorf("copy to host failed: error %d", int(err))
	}

	return result, nil
}

// ---------- Pinned Memory ----------

// AllocatePinned allocates page-locked (pinned) host memory for zero-copy transfers.
func AllocatePinned(size int) ([]float64, error) {
	if !initialized || size <= 0 {
		return make([]float64, size), nil
	}
	var ptr unsafe.Pointer
	byteSize := C.longlong(size * 8)
	if err := C.eml_allocate_pinned(&ptr, byteSize); err != 0 {
		return nil, fmt.Errorf("pinned memory allocation failed: error %d", int(err))
	}
	// Wrap the pinned pointer in a Go slice.
	data := (*[1 << 30]float64)(ptr)[:size:size]
	return data, nil
}

// FreePinned frees page-locked host memory allocated with AllocatePinned.
func FreePinned(p []float64) {
	if len(p) == 0 {
		return
	}
	// Get the base pointer of the slice
	ptr := unsafe.Pointer(&p[0])
	C.eml_free_pinned(ptr)
}
