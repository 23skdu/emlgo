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
