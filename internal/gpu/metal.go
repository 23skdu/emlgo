//go:build darwin && arm64 && metal && cgo && !cuda

package gpu

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Metal -framework Foundation

#include <stdint.h>
#include <stdlib.h>

typedef struct {
    int id;
    char name[256];
    int64_t memoryBytes;
    int maxThreadsPerThreadgroup;
    int warpSize;
    int clockRateMHz;
} MetalDeviceProps;

int eml_metal_init(const char *msl_source);
void eml_metal_cleanup(void);
int eml_metal_get_device_count(void);
int eml_metal_get_device_props(int device_id, MetalDeviceProps *props);
int eml_metal_launch_unary(const char *kernel_name, const void *x, void *result, int n, int elem_size);
int eml_metal_launch_binary(const char *kernel_name, const void *x, const void *y, void *result, int n, int elem_size);
int eml_metal_launch_dot_f32(const float *a, const float *b, float *out, int n);
int eml_metal_launch_dot_c64(const float *a, const float *b, float *out_re, float *out_im, int n);
*/
import "C"

import (
	_ "embed"
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

//go:embed eml_metal_shaders.metal
var metalShaderSource string

// ErrMetalFloat64Unsupported is returned when a 64-bit floating-point batch method
// (float64 or complex128) is invoked on the Apple Metal backend. Apple Metal Shading
// Language provides native 32-bit IEEE 754 float (and float2) compute kernels and
// does not have hardware 64-bit double support; rather than silently truncating
// precision, Metal returns this explicit error.
var ErrMetalFloat64Unsupported = errors.New("metal GPU does not support native 64-bit floating-point (float64/complex128); use float32/complex64 batch methods or CPU SIMD")

var (
	metalMu          sync.Mutex
	metalInitialized bool
)

func init() {
	_ = Init()
}

// Init initializes the Apple Metal GPU subsystem and compiles the embedded
// Metal Shading Language kernels with strict IEEE-754 math mode.
func Init() error {
	metalMu.Lock()
	defer metalMu.Unlock()
	if metalInitialized {
		return nil
	}
	csrc := C.CString(metalShaderSource)
	defer C.free(unsafe.Pointer(csrc))
	if rc := C.eml_metal_init(csrc); rc != 0 {
		return fmt.Errorf("failed to initialize Metal: error %d", int(rc))
	}
	metalInitialized = true
	return nil
}

// Shutdown releases all Apple Metal GPU resources and cached compute pipelines.
func Shutdown() {
	metalMu.Lock()
	defer metalMu.Unlock()
	if metalInitialized {
		C.eml_metal_cleanup()
		metalInitialized = false
	}
}

var getDevicesFn = getMetalDevices

// SetGetDevicesForTesting overrides the device lookup for testing purposes.
// It returns a restore function that reverts to the previous implementation.
func SetGetDevicesForTesting(fn func() ([]Device, error)) func() {
	prev := getDevicesFn
	getDevicesFn = fn
	return func() { getDevicesFn = prev }
}

var expBatchHook func(x []float64) ([]float64, error)

// SetExpBatchForTesting allows testing ExpBatch behavior.
func SetExpBatchForTesting(fn func(x []float64) ([]float64, error)) func() {
	prev := expBatchHook
	expBatchHook = fn
	return func() { expBatchHook = prev }
}

// GetDevices returns the available Apple Metal GPU devices.
func GetDevices() ([]Device, error) {
	return getDevicesFn()
}

func getMetalDevices() ([]Device, error) {
	count := int(C.eml_metal_get_device_count())
	if count <= 0 {
		return nil, nil
	}
	devices := make([]Device, 0, count)
	for i := 0; i < count; i++ {
		var props C.MetalDeviceProps
		if rc := C.eml_metal_get_device_props(C.int(i), &props); rc != 0 {
			return nil, fmt.Errorf("failed to get Metal device %d properties: error %d", i, int(rc))
		}
		devices = append(devices, Device{
			ID:                 int(props.id),
			Name:               C.GoString(&props.name[0]),
			MemoryBytes:        int64(props.memoryBytes),
			MaxThreadsPerBlock: int(props.maxThreadsPerThreadgroup),
			WarpSize:           int(props.warpSize),
			ClockRateKHz:       int(props.clockRateMHz) * 1000,
		})
	}
	return devices, nil
}

// Status returns a human-readable status string for the Apple Metal GPU subsystem.
func Status() string {
	devices, err := GetDevices()
	if err != nil {
		return fmt.Sprintf("GPU Error: %v", err)
	}
	if len(devices) == 0 {
		return "No Metal GPU devices found"
	}
	s := fmt.Sprintf("Found %d Metal GPU device(s):\n", len(devices))
	for _, d := range devices {
		s += fmt.Sprintf("  [%d] %s (%d MB, %d threads/threadgroup, SIMD width %d)\n",
			d.ID, d.Name, d.MemoryBytes/1024/1024, d.MaxThreadsPerBlock, d.WarpSize)
	}
	return s
}

// ---------- Generic Metal Launchers ----------

func metalLaunchUnary[T any](kernelName string, x []T, elemSize int) ([]T, error) {
	if len(x) == 0 {
		return nil, nil
	}
	dst := make([]T, len(x))
	if err := metalLaunchUnaryTo(kernelName, dst, x, elemSize); err != nil {
		return nil, err
	}
	return dst, nil
}

func metalLaunchUnaryTo[T any](kernelName string, dst, x []T, elemSize int) error {
	if len(dst) != len(x) {
		return fmt.Errorf("slice length mismatch: dst %d != x %d", len(dst), len(x))
	}
	if len(x) == 0 {
		return nil
	}
	if err := Init(); err != nil {
		return err
	}
	ckernel := C.CString(kernelName)
	defer C.free(unsafe.Pointer(ckernel))
	rc := C.eml_metal_launch_unary(
		ckernel,
		unsafe.Pointer(&x[0]),
		unsafe.Pointer(&dst[0]),
		C.int(len(x)),
		C.int(elemSize),
	)
	if rc != 0 {
		return fmt.Errorf("metal kernel %s failed: error %d", kernelName, int(rc))
	}
	return nil
}

func metalLaunchBinary[T any](kernelName string, x, y []T, elemSize int) ([]T, error) {
	if len(x) != len(y) {
		return nil, fmt.Errorf("x and y must have same length: %d vs %d", len(x), len(y))
	}
	if len(x) == 0 {
		return nil, nil
	}
	dst := make([]T, len(x))
	if err := metalLaunchBinaryTo(kernelName, dst, x, y, elemSize); err != nil {
		return nil, err
	}
	return dst, nil
}

func metalLaunchBinaryTo[T any](kernelName string, dst, x, y []T, elemSize int) error {
	if len(x) != len(y) || len(dst) != len(x) {
		return fmt.Errorf("slice length mismatch: dst=%d, x=%d, y=%d", len(dst), len(x), len(y))
	}
	if len(x) == 0 {
		return nil
	}
	if err := Init(); err != nil {
		return err
	}
	ckernel := C.CString(kernelName)
	defer C.free(unsafe.Pointer(ckernel))
	rc := C.eml_metal_launch_binary(
		ckernel,
		unsafe.Pointer(&x[0]),
		unsafe.Pointer(&y[0]),
		unsafe.Pointer(&dst[0]),
		C.int(len(x)),
		C.int(elemSize),
	)
	if rc != 0 {
		return fmt.Errorf("metal kernel %s failed: error %d", kernelName, int(rc))
	}
	return nil
}

// ============================================================================
// Float32 Device Methods (native Apple Silicon Metal GPU execution)
// ============================================================================

// ExpBatchF32 computes exp(x) for each float32 element on the Metal GPU.
func (d *Device) ExpBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_exp_f32", x, 4)
}

// ExpBatchToF32 computes exp(x) on the Metal GPU and writes into dst.
func (d *Device) ExpBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_exp_f32", dst, x, 4)
}

// LogBatchF32 computes ln(x) for each float32 element on the Metal GPU.
func (d *Device) LogBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_log_f32", x, 4)
}

// LogBatchToF32 computes ln(x) on the Metal GPU and writes into dst.
func (d *Device) LogBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_log_f32", dst, x, 4)
}

// SinBatchF32 computes sin(x) for each float32 element on the Metal GPU.
func (d *Device) SinBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_sin_f32", x, 4)
}

// SinBatchToF32 computes sin(x) on the Metal GPU and writes into dst.
func (d *Device) SinBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_sin_f32", dst, x, 4)
}

// CosBatchF32 computes cos(x) for each float32 element on the Metal GPU.
func (d *Device) CosBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_cos_f32", x, 4)
}

// CosBatchToF32 computes cos(x) on the Metal GPU and writes into dst.
func (d *Device) CosBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_cos_f32", dst, x, 4)
}

// TanBatchF32 computes tan(x) for each float32 element on the Metal GPU.
func (d *Device) TanBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_tan_f32", x, 4)
}

// TanBatchToF32 computes tan(x) on the Metal GPU and writes into dst.
func (d *Device) TanBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_tan_f32", dst, x, 4)
}

// SinhBatchF32 computes sinh(x) for each float32 element on the Metal GPU.
func (d *Device) SinhBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_sinh_f32", x, 4)
}

// SinhBatchToF32 computes sinh(x) on the Metal GPU and writes into dst.
func (d *Device) SinhBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_sinh_f32", dst, x, 4)
}

// CoshBatchF32 computes cosh(x) for each float32 element on the Metal GPU.
func (d *Device) CoshBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_cosh_f32", x, 4)
}

// CoshBatchToF32 computes cosh(x) on the Metal GPU and writes into dst.
func (d *Device) CoshBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_cosh_f32", dst, x, 4)
}

// TanhBatchF32 computes tanh(x) for each float32 element on the Metal GPU.
func (d *Device) TanhBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_tanh_f32", x, 4)
}

// TanhBatchToF32 computes tanh(x) on the Metal GPU and writes into dst.
func (d *Device) TanhBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_tanh_f32", dst, x, 4)
}

// SqrtBatchF32 computes sqrt(x) for each float32 element on the Metal GPU.
func (d *Device) SqrtBatchF32(x []float32) ([]float32, error) {
	return metalLaunchUnary("kernel_sqrt_f32", x, 4)
}

// SqrtBatchToF32 computes sqrt(x) on the Metal GPU and writes into dst.
func (d *Device) SqrtBatchToF32(dst, x []float32) error {
	return metalLaunchUnaryTo("kernel_sqrt_f32", dst, x, 4)
}

// EmlBatchF32 computes exp(x) - ln(y) for float32 slices on the Metal GPU.
func (d *Device) EmlBatchF32(x, y []float32) ([]float32, error) {
	return metalLaunchBinary("kernel_eml_f32", x, y, 4)
}

// EmlBatchToF32 computes exp(x) - ln(y) for float32 slices on the Metal GPU and writes into dst.
func (d *Device) EmlBatchToF32(dst, x, y []float32) error {
	return metalLaunchBinaryTo("kernel_eml_f32", dst, x, y, 4)
}

// DotBatchF32 computes the float32 dot product of a and b on the Metal GPU.
func (d *Device) DotBatchF32(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("slice length mismatch: %d vs %d", len(a), len(b))
	}
	if len(a) == 0 {
		return 0, nil
	}
	if err := Init(); err != nil {
		return 0, err
	}
	var out C.float
	rc := C.eml_metal_launch_dot_f32(
		(*C.float)(unsafe.Pointer(&a[0])),
		(*C.float)(unsafe.Pointer(&b[0])),
		&out,
		C.int(len(a)),
	)
	if rc != 0 {
		return 0, fmt.Errorf("metal dot_f32 kernel failed: error %d", int(rc))
	}
	return float32(out), nil
}

// ============================================================================
// Complex64 Device Methods (native Apple Silicon Metal GPU execution)
// ============================================================================

// ExpBatchC64 computes exp(z) for each complex64 element on the Metal GPU.
func (d *Device) ExpBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_exp_c64", x, 8)
}

// ExpBatchToC64 computes exp(z) on the Metal GPU and writes into dst.
func (d *Device) ExpBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_exp_c64", dst, x, 8)
}

// LogBatchC64 computes ln(z) for each complex64 element on the Metal GPU.
func (d *Device) LogBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_log_c64", x, 8)
}

// LogBatchToC64 computes ln(z) on the Metal GPU and writes into dst.
func (d *Device) LogBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_log_c64", dst, x, 8)
}

// SinBatchC64 computes sin(z) for each complex64 element on the Metal GPU.
func (d *Device) SinBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_sin_c64", x, 8)
}

// SinBatchToC64 computes sin(z) on the Metal GPU and writes into dst.
func (d *Device) SinBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_sin_c64", dst, x, 8)
}

// CosBatchC64 computes cos(z) for each complex64 element on the Metal GPU.
func (d *Device) CosBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_cos_c64", x, 8)
}

// CosBatchToC64 computes cos(z) on the Metal GPU and writes into dst.
func (d *Device) CosBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_cos_c64", dst, x, 8)
}

// TanBatchC64 computes tan(z) for each complex64 element on the Metal GPU.
func (d *Device) TanBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_tan_c64", x, 8)
}

// TanBatchToC64 computes tan(z) on the Metal GPU and writes into dst.
func (d *Device) TanBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_tan_c64", dst, x, 8)
}

// SinhBatchC64 computes sinh(z) for each complex64 element on the Metal GPU.
func (d *Device) SinhBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_sinh_c64", x, 8)
}

// SinhBatchToC64 computes sinh(z) on the Metal GPU and writes into dst.
func (d *Device) SinhBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_sinh_c64", dst, x, 8)
}

// CoshBatchC64 computes cosh(z) for each complex64 element on the Metal GPU.
func (d *Device) CoshBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_cosh_c64", x, 8)
}

// CoshBatchToC64 computes cosh(z) on the Metal GPU and writes into dst.
func (d *Device) CoshBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_cosh_c64", dst, x, 8)
}

// TanhBatchC64 computes tanh(z) for each complex64 element on the Metal GPU.
func (d *Device) TanhBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_tanh_c64", x, 8)
}

// TanhBatchToC64 computes tanh(z) on the Metal GPU and writes into dst.
func (d *Device) TanhBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_tanh_c64", dst, x, 8)
}

// SqrtBatchC64 computes sqrt(z) for each complex64 element on the Metal GPU.
func (d *Device) SqrtBatchC64(x []complex64) ([]complex64, error) {
	return metalLaunchUnary("kernel_sqrt_c64", x, 8)
}

// SqrtBatchToC64 computes sqrt(z) on the Metal GPU and writes into dst.
func (d *Device) SqrtBatchToC64(dst, x []complex64) error {
	return metalLaunchUnaryTo("kernel_sqrt_c64", dst, x, 8)
}

// EmlBatchC64 computes exp(x) - ln(y) for complex64 slices on the Metal GPU.
func (d *Device) EmlBatchC64(x, y []complex64) ([]complex64, error) {
	return metalLaunchBinary("kernel_eml_c64", x, y, 8)
}

// EmlBatchToC64 computes exp(x) - ln(y) for complex64 slices on the Metal GPU and writes into dst.
func (d *Device) EmlBatchToC64(dst, x, y []complex64) error {
	return metalLaunchBinaryTo("kernel_eml_c64", dst, x, y, 8)
}

// DotBatchC64 computes the Hermitian inner product sum(a[i] * conj(b[i])) on the Metal GPU.
func (d *Device) DotBatchC64(a, b []complex64) (complex64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("slice length mismatch: %d vs %d", len(a), len(b))
	}
	if len(a) == 0 {
		return 0, nil
	}
	if err := Init(); err != nil {
		return 0, err
	}
	var outRe, outIm C.float
	rc := C.eml_metal_launch_dot_c64(
		(*C.float)(unsafe.Pointer(&a[0])),
		(*C.float)(unsafe.Pointer(&b[0])),
		&outRe,
		&outIm,
		C.int(len(a)),
	)
	if rc != 0 {
		return 0, fmt.Errorf("metal dot_c64 kernel failed: error %d", int(rc))
	}
	return complex(float32(outRe), float32(outIm)), nil
}

// ============================================================================
// Float64 & Complex128 Device Methods (explicitly unsupported on Metal)
// ============================================================================

// ExpBatch returns ErrMetalFloat64Unsupported unless overridden by SetExpBatchForTesting.
func (d *Device) ExpBatch(x []float64) ([]float64, error) {
	if expBatchHook != nil {
		return expBatchHook(x)
	}
	return nil, ErrMetalFloat64Unsupported
}

// ExpBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) ExpBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// LogBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) LogBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// LogBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) LogBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// SinBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) SinBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// SinBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) SinBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// CosBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) CosBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// CosBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) CosBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// TanBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) TanBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// TanBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) TanBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// SinhBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) SinhBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// SinhBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) SinhBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// CoshBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) CoshBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// CoshBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) CoshBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// TanhBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) TanhBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// TanhBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) TanhBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// SqrtBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) SqrtBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// SqrtBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) SqrtBatchTo(dst, x []float64) error { return ErrMetalFloat64Unsupported }

// EmlBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) EmlBatch(x, y []float64) ([]float64, error) {
	return nil, ErrMetalFloat64Unsupported
}

// EmlBatchTo returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) EmlBatchTo(dst, x, y []float64) error { return ErrMetalFloat64Unsupported }

// DotBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) DotBatch(a, b []float64) (float64, error) { return 0, ErrMetalFloat64Unsupported }

// AbsBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) AbsBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// NegBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) NegBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// PowBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) PowBatch(x []float64, exponent float64) ([]float64, error) {
	return nil, ErrMetalFloat64Unsupported
}

// InvBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) InvBatch(x []float64) ([]float64, error) { return nil, ErrMetalFloat64Unsupported }

// FmaBatch returns ErrMetalFloat64Unsupported because Metal has no hardware float64.
func (d *Device) FmaBatch(a, b, c []float64) ([]float64, error) {
	return nil, ErrMetalFloat64Unsupported
}

// ExpBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) ExpBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// ExpBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) ExpBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// LogBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) LogBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// LogBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) LogBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// SinBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) SinBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// SinBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) SinBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// CosBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) CosBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// CosBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) CosBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// TanBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) TanBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// TanBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) TanBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// SinhBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) SinhBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// SinhBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) SinhBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// CoshBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) CoshBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// CoshBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) CoshBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// TanhBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) TanhBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// TanhBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) TanhBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// SqrtBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) SqrtBatchC128(x []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// SqrtBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) SqrtBatchToC128(dst, x []complex128) error { return ErrMetalFloat64Unsupported }

// EmlBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) EmlBatchC128(x, y []complex128) ([]complex128, error) {
	return nil, ErrMetalFloat64Unsupported
}

// EmlBatchToC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) EmlBatchToC128(dst, x, y []complex128) error { return ErrMetalFloat64Unsupported }

// DotBatchC128 returns ErrMetalFloat64Unsupported because Metal has no hardware complex128.
func (d *Device) DotBatchC128(a, b []complex128) (complex128, error) {
	return 0, ErrMetalFloat64Unsupported
}

// ============================================================================
// Unified Memory & Stream Helpers
// ============================================================================

// NewStream returns a lightweight stream handle on Metal (commands are synchronized per batch).
func NewStream() (*Stream, error) {
	if err := Init(); err != nil {
		return nil, err
	}
	return &Stream{handle: 1}, nil
}

// Sync is a no-op because Metal batch dispatches synchronize before returning.
func (s *Stream) Sync() error { return nil }

// Destroy invalidates the Metal stream handle.
func (s *Stream) Destroy() error {
	if s != nil {
		s.handle = 0
	}
	return nil
}

// AllocatePinned allocates unified-memory host storage on Apple Silicon.
func AllocatePinned(size int) ([]float64, error) { return make([]float64, size), nil }

// FreePinned is a no-op on Apple Silicon unified memory (managed by Go GC).
func FreePinned(p []float64) { _ = p }

// AllocatePinnedF32 allocates unified-memory float32 storage on Apple Silicon.
func AllocatePinnedF32(size int) ([]float32, error) { return make([]float32, size), nil }

// FreePinnedF32 is a no-op on Apple Silicon unified memory (managed by Go GC).
func FreePinnedF32(p []float32) { _ = p }

// AllocatePinnedC64 allocates unified-memory complex64 storage on Apple Silicon.
func AllocatePinnedC64(size int) ([]complex64, error) { return make([]complex64, size), nil }

// FreePinnedC64 is a no-op on Apple Silicon unified memory (managed by Go GC).
func FreePinnedC64(p []complex64) { _ = p }

// AllocatePinnedC128 allocates unified-memory complex128 storage on Apple Silicon.
func AllocatePinnedC128(size int) ([]complex128, error) { return make([]complex128, size), nil }

// FreePinnedC128 is a no-op on Apple Silicon unified memory (managed by Go GC).
func FreePinnedC128(p []complex128) { _ = p }

// CopyToDeviceAsync returns an error on Metal because Apple Silicon uses unified shared memory.
func CopyToDeviceAsync(dst, src unsafe.Pointer, size int, stream *Stream) error {
	return fmt.Errorf("explicit async host-device copy is not required on Metal unified memory")
}

// CopyToHostAsync returns an error on Metal because Apple Silicon uses unified shared memory.
func CopyToHostAsync(dst, src unsafe.Pointer, size int, stream *Stream) error {
	return fmt.Errorf("explicit async device-host copy is not required on Metal unified memory")
}
