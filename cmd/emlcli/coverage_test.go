package main

import (
	"fmt"
	"testing"

	"github.com/emlgo/eml/internal/gpu"
	"github.com/emlgo/eml/internal/jit"
)

func TestEmlCli100PercentCoverage(t *testing.T) {
	oldExit := exitFunc
	exitCode := 0
	exitFunc = func(code int) {
		exitCode = code
	}
	defer func() { exitFunc = oldExit }()

	// 1. main() with < 2 args
	withArgs(t, []string{"eml"}, func() {
		exitCode = 0
		main()
		if exitCode != 1 {
			t.Errorf("expected exitCode 1, got %d", exitCode)
		}
	})

	// 2. main() with unknown command
	withArgs(t, []string{"eml", "unknown-cmd"}, func() {
		exitCode = 0
		main()
		if exitCode != 1 {
			t.Errorf("expected exitCode 1, got %d", exitCode)
		}
	})

	// 3. main() with each valid command
	cmds := []string{"demo", "gpu-status", "gpu-bench", "gpu-verify", "jit-test", "decompile"}
	for _, cmd := range cmds {
		withArgs(t, []string{"eml", cmd}, func() {
			exitCode = 0
			main()
		})
	}

	// 4. runGpuStatus with error and with devices
	restore := gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return nil, fmt.Errorf("mock error")
	})
	runGpuStatus()
	restore()

	mockDev := gpu.Device{
		ID:                 0,
		Name:               "Test Device",
		ComputeMajor:       8,
		ComputeMinor:       0,
		MemoryBytes:        8 * 1024 * 1024 * 1024,
		MaxThreadsPerBlock: 1024,
		WarpSize:           32,
		ClockRateKHz:       1500000,
	}

	restore = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return []gpu.Device{mockDev}, nil
	})
	runGpuStatus()
	restore()

	// 5. runGpuBench with error, empty, and devices
	restore = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return nil, fmt.Errorf("bench device error")
	})
	runGpuBench()
	restore()

	restore = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return nil, nil
	})
	runGpuBench()
	restore()

	restore = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return []gpu.Device{mockDev}, nil
	})
	// Temporarily replace gpuOps with a fast, small op
	oldGpuOps := gpuOps
	gpuOps = []gpuOp{
		{
			Name: "MockOpPass",
			RunGPU: func(d *gpu.Device, x []float64) ([]float64, error) {
				res := make([]float64, len(x))
				copy(res, x)
				return res, nil
			},
			CPU: func(x float64) float64 { return x },
			CPUBatch: func(x []float64) []float64 {
				r := make([]float64, len(x))
				copy(r, x)
				return r
			},
		},
		{
			Name: "MockOpFail",
			RunGPU: func(d *gpu.Device, x []float64) ([]float64, error) {
				return nil, fmt.Errorf("gpu launch fail")
			},
			CPU: func(x float64) float64 { return x },
			CPUBatch: func(x []float64) []float64 {
				return x
			},
		},
	}
	runGpuBench()
	gpuOps = oldGpuOps
	restore()

	// 6. runGpuVerify with error, empty, and devices (passing and failing)
	restore = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return nil, fmt.Errorf("verify device error")
	})
	runGpuVerify()
	restore()

	restore = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return nil, nil
	})
	runGpuVerify()
	restore()

	restore = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return []gpu.Device{mockDev}, nil
	})
	// Run default stub verification (errors out on launch)
	runGpuVerify()

	// Run with all passing operations (covering allPassed == true branch)
	gpuOps = []gpuOp{
		{
			Name: "PassOp",
			RunGPU: func(d *gpu.Device, x []float64) ([]float64, error) {
				res := make([]float64, len(x))
				copy(res, x)
				return res, nil
			},
			CPU: func(x float64) float64 { return x },
			CPUBatch: func(x []float64) []float64 { return x },
		},
	}
	runGpuVerify()

	// Run with failing operations (covering anyFailure branch)
	gpuOps = []gpuOp{
		{
			Name: "FailULPOp",
			RunGPU: func(d *gpu.Device, x []float64) ([]float64, error) {
				res := make([]float64, len(x))
				for i, v := range x {
					res[i] = v + 10.0 // force ULP > 1
				}
				return res, nil
			},
			CPU: func(x float64) float64 { return x },
			CPUBatch: func(x []float64) []float64 { return x },
		},
	}
	runGpuVerify()
	gpuOps = oldGpuOps
	restore()

	// 7. runJitTest with compile error
	oldJitExpr := jitTestExpr
	jitTestExpr = "1 +"
	runJitTest()
	jitTestExpr = oldJitExpr
}

func TestJitTestErrorBranch(t *testing.T) {
	// Exercise jit.Compile error
	c := jit.NewCompiler()
	_, err := c.Compile("1 +")
	if err == nil {
		t.Fatal("expected error")
	}
}
