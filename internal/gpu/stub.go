//go:build !cuda && !(darwin && arm64)

package gpu

import (
	"fmt"
	"runtime"
)

// Override in tests to exercise specific paths.
var cudaSupportedOS = func() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "windows"
}

var getDevicesFn = getDevicesImpl

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

func (d *Device) ExpBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) LogBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) SinBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) CosBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) TanBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) SinhBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) CoshBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) TanhBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) SqrtBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) EmlBatch(x, y []float64) ([]float64, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) AbsBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) NegBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) PowBatch(x []float64, exponent float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) InvBatch(x []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) FmaBatch(a, b, c []float64) ([]float64, error) {
	_ = d
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) ExpBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) LogBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) SinBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) CosBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) TanBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) SinhBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) CoshBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) TanhBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) SqrtBatchTo(dst, x []float64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) EmlBatchTo(dst, x, y []float64) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

func (d *Device) DotBatch(a, b []float64) (float64, error) {
	_ = d
	_ = a
	_ = b
	return 0, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// Float32 operations
func (d *Device) ExpBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) ExpBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) LogBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) LogBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CosBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CosBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinhBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinhBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CoshBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CoshBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanhBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanhBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SqrtBatchF32(x []float32) ([]float32, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SqrtBatchToF32(dst, x []float32) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) EmlBatchF32(x, y []float32) ([]float32, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) EmlBatchToF32(dst, x, y []float32) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) DotBatchF32(a, b []float32) (float32, error) {
	_ = d
	_ = a
	_ = b
	return 0, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// Complex64 operations
func (d *Device) ExpBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) ExpBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) LogBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) LogBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CosBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CosBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinhBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinhBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CoshBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CoshBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanhBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanhBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SqrtBatchC64(x []complex64) ([]complex64, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SqrtBatchToC64(dst, x []complex64) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) EmlBatchC64(x, y []complex64) ([]complex64, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) EmlBatchToC64(dst, x, y []complex64) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) DotBatchC64(a, b []complex64) (complex64, error) {
	_ = d
	_ = a
	_ = b
	return 0, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}

// Complex128 operations
func (d *Device) ExpBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) ExpBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) LogBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) LogBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CosBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CosBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinhBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SinhBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CoshBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) CoshBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanhBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) TanhBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SqrtBatchC128(x []complex128) ([]complex128, error) {
	_ = d
	_ = x
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) SqrtBatchToC128(dst, x []complex128) error {
	_ = d
	_ = dst
	_ = x
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) EmlBatchC128(x, y []complex128) ([]complex128, error) {
	_ = d
	_ = x
	_ = y
	return nil, fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
func (d *Device) EmlBatchToC128(dst, x, y []complex128) error {
	_ = d
	_ = dst
	_ = x
	_ = y
	return fmt.Errorf("GPU execution not available (build with -tags cuda)")
}
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

func AllocatePinnedF32(size int) ([]float32, error) {
	return make([]float32, size), nil
}

func FreePinnedF32(p []float32) {
	_ = p
}

func AllocatePinnedC64(size int) ([]complex64, error) {
	return make([]complex64, size), nil
}

func FreePinnedC64(p []complex64) {
	_ = p
}

func AllocatePinnedC128(size int) ([]complex128, error) {
	return make([]complex128, size), nil
}

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
