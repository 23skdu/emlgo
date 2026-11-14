//go:build !cuda || !cgo
// +build !cuda !cgo

package gpu

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Override in tests to exercise specific paths.
var cudaSupportedOS = func() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "windows"
}

var getDevicesFn = getDevicesImpl

// SetGetDevicesForTesting overrides the device lookup for testing purposes.
// It returns a restore function that reverts to the previous implementation.
func SetGetDevicesForTesting(fn func() ([]Device, error)) func() {
	prev := getDevicesFn
	getDevicesFn = fn
	return func() { getDevicesFn = prev }
}

// GetDevices returns a list of available GPU devices.
// Without the cuda build tag, always returns empty.
func GetDevices() ([]Device, error) {
	return getDevicesFn()
}

func getDevicesImpl() ([]Device, error) {
	if !cudaSupportedOS() {
		return nil, fmt.Errorf("CUDA only supported on Linux and Windows")
	}
	return nil, nil
}

// Status returns a human-readable status string for the GPU subsystem.
func Status() string {
	devices, err := getDevicesFn()
	if err != nil {
		return fmt.Sprintf("GPU Error: %v", err)
	}
	if len(devices) == 0 {
		return "No GPU devices found (build with -tags cuda to enable)"
	}
	return fmt.Sprintf("Found %d GPU device(s)", len(devices))
}

// Init initializes the GPU subsystem. Stub: always succeeds.
func Init() error {
	return nil
}

// Shutdown cleans up the GPU subsystem. Stub: no-op.
func Shutdown() {
	_ = 1
}

var expBatchHook func(x []float64) ([]float64, error)

// SetExpBatchForTesting allows testing ExpBatch behavior.
func SetExpBatchForTesting(fn func(x []float64) ([]float64, error)) func() {
	prev := expBatchHook
	expBatchHook = fn
	return func() { expBatchHook = prev }
}

// ExpBatch computes the element-wise exponential of x on the device, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in, unless overridden by
// SetExpBatchForTesting.
func (d *Device) ExpBatch(x []float64) ([]float64, error) {
	_ = d
	if expBatchHook != nil {
		return expBatchHook(x)
	}
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatch computes the element-wise natural logarithm of x on the device, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatch computes the element-wise sine of x on the device, returning a newly allocated result.
// Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatch computes the element-wise cosine of x on the device, returning a newly allocated result.
// Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatch computes the element-wise tangent of x on the device, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatch computes the element-wise hyperbolic sine of x on the device, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatch computes the element-wise hyperbolic cosine of x on the device, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatch computes the element-wise hyperbolic tangent of x on the device, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatch computes the element-wise square root of x on the device, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatch computes the element-wise EML operator exp(x) - ln(y) over corresponding elements of x
// and y, returning a newly allocated result. Stub: returns an error because CUDA support was not
// compiled in.
func (d *Device) EmlBatch(x, y []float64) ([]float64, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// AbsBatch computes the element-wise absolute value of x on the device, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) AbsBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// NegBatch computes the element-wise negation of x on the device, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) NegBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// PowBatch raises each element of x to the power exponent on the device, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) PowBatch(x []float64, exponent float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// InvBatch computes the element-wise reciprocal of x on the device, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) InvBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// FmaBatch computes the element-wise fused multiply-add a*b + c over a, b, and c on the device,
// returning a newly allocated result. Stub: returns an error because CUDA support was not compiled
// in.
func (d *Device) FmaBatch(a, b, c []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// ExpBatchTo computes the element-wise exponential of x on the device and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) ExpBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatchTo computes the element-wise natural logarithm of x on the device and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatchTo computes the element-wise sine of x on the device and writes the result into dst.
// Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatchTo computes the element-wise cosine of x on the device and writes the result into dst.
// Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatchTo computes the element-wise tangent of x on the device and writes the result into dst.
// Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatchTo computes the element-wise hyperbolic sine of x on the device and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatchTo computes the element-wise hyperbolic cosine of x on the device and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatchTo computes the element-wise hyperbolic tangent of x on the device and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatchTo computes the element-wise square root of x on the device and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatchTo computes the element-wise EML operator exp(x) - ln(y) over corresponding elements of x
// and y and writes the result into dst. Stub: returns an error because CUDA support was not
// compiled in.
func (d *Device) EmlBatchTo(dst, x, y []float64) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// DotBatch computes the dot product of a and b on the device. Stub: returns an error because CUDA
// support was not compiled in.
func (d *Device) DotBatch(a, b []float64) (float64, error) {
	_ = d
	_ = a
	_ = b
	return 0, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// Float32 operations
// ExpBatchF32 computes the element-wise exponential of each float32 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) ExpBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// ExpBatchToF32 computes the element-wise exponential of each float32 in x and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) ExpBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatchF32 computes the element-wise natural logarithm of each float32 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatchToF32 computes the element-wise natural logarithm of each float32 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatchF32 computes the element-wise sine of each float32 in x, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatchToF32 computes the element-wise sine of each float32 in x and writes the result into dst.
// Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatchF32 computes the element-wise cosine of each float32 in x, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatchToF32 computes the element-wise cosine of each float32 in x and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatchF32 computes the element-wise tangent of each float32 in x, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatchToF32 computes the element-wise tangent of each float32 in x and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatchF32 computes the element-wise hyperbolic sine of each float32 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatchToF32 computes the element-wise hyperbolic sine of each float32 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatchF32 computes the element-wise hyperbolic cosine of each float32 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatchToF32 computes the element-wise hyperbolic cosine of each float32 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatchF32 computes the element-wise hyperbolic tangent of each float32 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatchToF32 computes the element-wise hyperbolic tangent of each float32 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatchF32 computes the element-wise square root of each float32 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatchToF32 computes the element-wise square root of each float32 in x and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatchF32 computes the element-wise EML operator exp(x) - ln(y) over corresponding float32
// elements of x and y, returning a newly allocated result. Stub: returns an error because CUDA
// support was not compiled in.
func (d *Device) EmlBatchF32(x, y []float32) ([]float32, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatchToF32 computes the element-wise EML operator exp(x) - ln(y) over corresponding float32
// elements of x and y and writes the result into dst. Stub: returns an error because CUDA support
// was not compiled in.
func (d *Device) EmlBatchToF32(dst, x, y []float32) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// DotBatchF32 computes the dot product of the float32 elements of a and b on the device. Stub:
// returns an error because CUDA support was not compiled in.
func (d *Device) DotBatchF32(a, b []float32) (float32, error) {
	_ = d
	_ = a
	_ = b
	return 0, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// Complex64 operations
// ExpBatchC64 computes the element-wise exponential of each complex64 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) ExpBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// ExpBatchToC64 computes the element-wise exponential of each complex64 in x and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) ExpBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatchC64 computes the element-wise natural logarithm of each complex64 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatchToC64 computes the element-wise natural logarithm of each complex64 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatchC64 computes the element-wise sine of each complex64 in x, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatchToC64 computes the element-wise sine of each complex64 in x and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatchC64 computes the element-wise cosine of each complex64 in x, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatchToC64 computes the element-wise cosine of each complex64 in x and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatchC64 computes the element-wise tangent of each complex64 in x, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatchToC64 computes the element-wise tangent of each complex64 in x and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatchC64 computes the element-wise hyperbolic sine of each complex64 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatchToC64 computes the element-wise hyperbolic sine of each complex64 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatchC64 computes the element-wise hyperbolic cosine of each complex64 in x, returning a
// newly allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatchToC64 computes the element-wise hyperbolic cosine of each complex64 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatchC64 computes the element-wise hyperbolic tangent of each complex64 in x, returning a
// newly allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatchToC64 computes the element-wise hyperbolic tangent of each complex64 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatchC64 computes the element-wise square root of each complex64 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatchToC64 computes the element-wise square root of each complex64 in x and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatchC64 computes the element-wise EML operator exp(x) - ln(y) over corresponding complex64
// elements of x and y, returning a newly allocated result. Stub: returns an error because CUDA
// support was not compiled in.
func (d *Device) EmlBatchC64(x, y []complex64) ([]complex64, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatchToC64 computes the element-wise EML operator exp(x) - ln(y) over corresponding complex64
// elements of x and y and writes the result into dst. Stub: returns an error because CUDA support
// was not compiled in.
func (d *Device) EmlBatchToC64(dst, x, y []complex64) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// DotBatchC64 computes the dot product of the complex64 elements of a and b on the device. Stub:
// returns an error because CUDA support was not compiled in.
func (d *Device) DotBatchC64(a, b []complex64) (complex64, error) {
	_ = d
	_ = a
	_ = b
	return 0, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// Complex128 operations
// ExpBatchC128 computes the element-wise exponential of each complex128 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) ExpBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// ExpBatchToC128 computes the element-wise exponential of each complex128 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) ExpBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatchC128 computes the element-wise natural logarithm of each complex128 in x, returning a
// newly allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// LogBatchToC128 computes the element-wise natural logarithm of each complex128 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) LogBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatchC128 computes the element-wise sine of each complex128 in x, returning a newly allocated
// result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinBatchToC128 computes the element-wise sine of each complex128 in x and writes the result into
// dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatchC128 computes the element-wise cosine of each complex128 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CosBatchToC128 computes the element-wise cosine of each complex128 in x and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CosBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatchC128 computes the element-wise tangent of each complex128 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanBatchToC128 computes the element-wise tangent of each complex128 in x and writes the result
// into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatchC128 computes the element-wise hyperbolic sine of each complex128 in x, returning a
// newly allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SinhBatchToC128 computes the element-wise hyperbolic sine of each complex128 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SinhBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatchC128 computes the element-wise hyperbolic cosine of each complex128 in x, returning a
// newly allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// CoshBatchToC128 computes the element-wise hyperbolic cosine of each complex128 in x and writes
// the result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) CoshBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatchC128 computes the element-wise hyperbolic tangent of each complex128 in x, returning a
// newly allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// TanhBatchToC128 computes the element-wise hyperbolic tangent of each complex128 in x and writes
// the result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) TanhBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatchC128 computes the element-wise square root of each complex128 in x, returning a newly
// allocated result. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// SqrtBatchToC128 computes the element-wise square root of each complex128 in x and writes the
// result into dst. Stub: returns an error because CUDA support was not compiled in.
func (d *Device) SqrtBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatchC128 computes the element-wise EML operator exp(x) - ln(y) over corresponding complex128
// elements of x and y, returning a newly allocated result. Stub: returns an error because CUDA
// support was not compiled in.
func (d *Device) EmlBatchC128(x, y []complex128) ([]complex128, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// EmlBatchToC128 computes the element-wise EML operator exp(x) - ln(y) over corresponding
// complex128 elements of x and y and writes the result into dst. Stub: returns an error because
// CUDA support was not compiled in.
func (d *Device) EmlBatchToC128(dst, x, y []complex128) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// DotBatchC128 computes the dot product of the complex128 elements of a and b on the device. Stub:
// returns an error because CUDA support was not compiled in.
func (d *Device) DotBatchC128(a, b []complex128) (complex128, error) {
	_ = d
	_ = a
	_ = b
	return 0, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// NewStream creates an async stream. Stub: returns nil.
func NewStream() (*Stream, error) {
	return nil, fmt.Errorf("streams not available without cuda tag")
}

// AllocatePinned allocates page-locked host memory. Stub: returns regular memory.
func AllocatePinned(size int) ([]float64, error) {
	return make([]float64, size), nil
}

// FreePinned frees page-locked host memory. Stub: no-op.
func FreePinned(p []float64) {
	_ = p
}

// AllocatePinnedF32 allocates page-locked float32 host memory. Stub: returns regular memory.
func AllocatePinnedF32(size int) ([]float32, error) {
	return make([]float32, size), nil
}

// FreePinnedF32 releases page-locked float32 host memory. Stub: no-op.
func FreePinnedF32(p []float32) {
	_ = p
}

// AllocatePinnedC64 allocates page-locked complex64 host memory. Stub: returns regular memory.
func AllocatePinnedC64(size int) ([]complex64, error) {
	return make([]complex64, size), nil
}

// FreePinnedC64 releases page-locked complex64 host memory. Stub: no-op.
func FreePinnedC64(p []complex64) {
	_ = p
}

// AllocatePinnedC128 allocates page-locked complex128 host memory. Stub: returns regular memory.
func AllocatePinnedC128(size int) ([]complex128, error) {
	return make([]complex128, size), nil
}

// FreePinnedC128 releases page-locked complex128 host memory. Stub: no-op.
func FreePinnedC128(p []complex128) {
	_ = p
}

// Sync is a no-op in stub mode.
func (s *Stream) Sync() error {
	_ = s
	return nil
}

// Destroy is a no-op in stub mode.
func (s *Stream) Destroy() error {
	_ = s
	return nil
}

// CopyToDeviceAsync is a no-op stub in non-CUDA builds.
func CopyToDeviceAsync(dst, src unsafe.Pointer, size int, stream *Stream) error {
	_ = dst
	_ = src
	_ = size
	_ = stream
	return fmt.Errorf("CUDA not available (build with -tags cuda)")
}

// CopyToHostAsync is a no-op stub in non-CUDA builds.
func CopyToHostAsync(dst, src unsafe.Pointer, size int, stream *Stream) error {
	_ = dst
	_ = src
	_ = size
	_ = stream
	return fmt.Errorf("CUDA not available (build with -tags cuda)")
}
