package gpu

import (
	"fmt"
	"runtime"
)

// Device represents a GPU device.
type Device struct {
	ID   int
	Name string
}

// GetDevices returns a list of available GPU devices.
func GetDevices() ([]Device, error) {
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		return nil, fmt.Errorf("CUDA only supported on Linux and Windows currently")
	}
	// TODO: Use cgo to call cudaGetDeviceCount
	return nil, nil
}

// ExpBatch executes a batch exponential operation on the GPU.
func (d *Device) ExpBatch(x []float64) ([]float64, error) {
	// 1. Allocate device memory
	// 2. Copy data to device
	// 3. Launch kernel
	// 4. Copy results back
	return nil, fmt.Errorf("GPU execution not yet implemented")
}

// Status returns the current status of the GPU subsystem.
func Status() string {
	devices, err := GetDevices()
	if err != nil {
		return fmt.Sprintf("GPU Error: %v", err)
	}
	if len(devices) == 0 {
		return "No GPU devices found"
	}
	return fmt.Sprintf("Found %d GPU devices", len(devices))
}
