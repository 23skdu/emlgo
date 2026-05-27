//go:build darwin && arm64 && !cuda

package gpu

/*
#cgo LDFLAGS: -framework Metal -framework Foundation

#include <stdint.h>
#include <stdlib.h>

typedef struct {
    int id;
    char name[256];
    int64_t memoryBytes;
    int maxThreadsPerThreadgroup;
    int clockRateMHz;
} MetalDeviceProps;

int eml_metal_init(void);
void eml_metal_cleanup(void);
int eml_metal_get_device_count(void);
int eml_metal_get_device_props(int id, MetalDeviceProps *props);
int eml_metal_launch_unary(const char *name, const double *x,
                           double *result, int n);
int eml_metal_launch_binary(const char *name, const double *a,
                            const double *b, double *result, int n);
int eml_metal_launch_eml(const double *x, const double *y,
                         double *result, int n);
void eml_metal_sync(void);
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

var metalInitialized bool

func Init() error {
	if C.eml_metal_init() == 0 {
		metalInitialized = true
		return nil
	}
	return fmt.Errorf("failed to initialize Metal")
}

func Shutdown() {
	if metalInitialized {
		C.eml_metal_cleanup()
		metalInitialized = false
	}
}

var getDevicesFn = getMetalDevices

func GetDevices() ([]Device, error) {
	return getDevicesFn()
}

func getMetalDevices() ([]Device, error) {
	if !metalInitialized {
		if err := Init(); err != nil {
			return nil, err
		}
	}

	count := int(C.eml_metal_get_device_count())
	if count == 0 {
		return nil, nil
	}

	devices := make([]Device, count)
	for i := 0; i < count; i++ {
		if i > 1<<30 {
			continue
		}
		var props C.MetalDeviceProps
		if err := C.eml_metal_get_device_props(C.int(i), &props); err != 0 {
			return nil, fmt.Errorf("failed to get device %d properties", i)
		}
		devices[i] = Device{
			ID:                 int(props.id),
			Name:               C.GoString(&props.name[0]),
			MemoryBytes:        int64(props.memoryBytes),
			MaxThreadsPerBlock: int(props.maxThreadsPerThreadgroup),
			WarpSize:           1,
			ClockRateKHz:       int(props.clockRateMHz) * 1000,
		}
	}
	return devices, nil
}

func Status() string {
	if !metalInitialized {
		if err := Init(); err != nil {
			return fmt.Sprintf("Metal Error: %v", err)
		}
	}

	devices, err := getDevicesFn()
	if err != nil {
		return fmt.Sprintf("Metal Error: %v", err)
	}
	if len(devices) == 0 {
		return "No Metal-capable GPU devices found"
	}

	s := fmt.Sprintf("Found %d Metal GPU device(s):\n", len(devices))
	for _, d := range devices {
		s += fmt.Sprintf("  [%d] %s (%s, %d MB)\n",
			d.ID, d.Name, runtime.GOARCH, d.MemoryBytes/1024/1024)
	}
	return s
}

func launchBatch(x []float64, kernelName string) ([]float64, error) {
	if !metalInitialized {
		return nil, fmt.Errorf("Metal not initialized")
	}
	if len(x) == 0 {
		return nil, nil
	}

	n := len(x)
	if n > 1<<30 {
		return nil, fmt.Errorf("batch size %d exceeds max int32", n)
	}
	cName := C.CString(kernelName)
	defer C.free(unsafe.Pointer(cName))

	result := make([]float64, n)
	err := C.eml_metal_launch_unary(cName,
		(*C.double)(unsafe.Pointer(&x[0])),
		(*C.double)(unsafe.Pointer(&result[0])),
		C.int(n),
	)
	if err != 0 {
		return nil, fmt.Errorf("Metal kernel %s failed: error %d", kernelName, int(err))
	}
	return result, nil
}

func (d *Device) ExpBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_exp")
}

func (d *Device) LogBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_log")
}

func (d *Device) SinBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_sin")
}

func (d *Device) CosBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_cos")
}

func (d *Device) TanBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_tan")
}

func (d *Device) SinhBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_sinh")
}

func (d *Device) CoshBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_cosh")
}

func (d *Device) TanhBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_tanh")
}

func (d *Device) SqrtBatch(x []float64) ([]float64, error) {
	return launchBatch(x, "kernel_sqrt")
}

func (d *Device) EmlBatch(x, y []float64) ([]float64, error) {
	if !metalInitialized {
		return nil, fmt.Errorf("Metal not initialized")
	}
	if len(x) != len(y) {
		return nil, fmt.Errorf("x and y must have same length: %d vs %d", len(x), len(y))
	}
	if len(x) == 0 {
		return nil, nil
	}

	n := len(x)
	if n > 1<<30 {
		return nil, fmt.Errorf("batch size %d exceeds max int32", n)
	}
	result := make([]float64, n)
	err := C.eml_metal_launch_eml(
		(*C.double)(unsafe.Pointer(&x[0])),
		(*C.double)(unsafe.Pointer(&y[0])),
		(*C.double)(unsafe.Pointer(&result[0])),
		C.int(n),
	)
	if err != 0 {
		return nil, fmt.Errorf("Metal EML kernel failed: error %d", int(err))
	}
	return result, nil
}

func NewStream() (*Stream, error) {
	return nil, fmt.Errorf("streams not available with Metal backend")
}

func AllocatePinned(size int) ([]float64, error) {
	return make([]float64, size), nil
}

func FreePinned(p []float64) {}

func (s *Stream) Sync() error { return nil }

func (s *Stream) Destroy() error { return nil }
