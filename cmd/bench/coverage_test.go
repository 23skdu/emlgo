package main

import (
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/emlgo/eml/internal/gpu"
)

func TestBenchComprehensiveCoverage(t *testing.T) {
	oldExit := exitFunc
	exitCode := 0
	exitFunc = func(code int) {
		exitCode = code
	}
	defer func() {
		exitFunc = oldExit
		device = "cpu"
		profile = ""
		testAccuracy = false
		compareMode = false
		checkRegressionFlag = false
		forceFailTol = false
		forceFailAccuracy = false
		iterations = 1000000
		gpuBenchSizes = []int{1024, 4096, 16384, 65536, 262144, 1048576}
		jitPolys = defaultPolys
		jitBenchN = 500000
		typeFilter = "all"
		osCreateProfile = os.Create
		startCPUProfile = pprof.StartCPUProfile
		writeHeapProfile = pprof.WriteHeapProfile
		writeBlockProfile = func(p *pprof.Profile, w io.Writer) error { return p.WriteTo(w, 0) }
	}()

	// 1. Test GPU benchmarks branches
	device = "gpu"
	gpuBenchSizes = []int{16}

	// 1a. GPU error
	exitCode = 0
	restoreDev := gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return nil, fmt.Errorf("simulated GPU failure")
	})
	main()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on GPU error, got %d", exitCode)
	}

	// 1b. GPU no devices
	exitCode = 0
	_ = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return nil, nil
	})
	main()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on empty GPU list, got %d", exitCode)
	}

	// 1c. GPU device found, ExpBatch error
	exitCode = 0
	mockDevs := []gpu.Device{{ID: 0, Name: "Mock GPU", ComputeMajor: 8, ComputeMinor: 0, MemoryBytes: 1024 * 1024 * 1024}}
	_ = gpu.SetGetDevicesForTesting(func() ([]gpu.Device, error) {
		return mockDevs, nil
	})
	restoreExp := gpu.SetExpBatchForTesting(func(x []float64) ([]float64, error) {
		return nil, fmt.Errorf("mock exp failure")
	})
	main()

	// 1d. GPU device found, ExpBatch success
	_ = gpu.SetExpBatchForTesting(func(x []float64) ([]float64, error) {
		return make([]float64, len(x)), nil
	})
	main()

	restoreExp()
	restoreDev()
	device = "cpu"

	// 2. Test JIT benchmarks branches
	device = "jit"
	jitBenchN = 10

	// 2a. Normal JIT
	main()

	// 2b. JIT compile error branch
	jitPolys = []polyExpr{
		{"invalid", "unknown_func(x)", nil},
	}
	main()
	jitPolys = defaultPolys
	device = "cpu"

	// 3. Test Profiling branches
	// 3a. Unknown profile
	profile = "invalid_prof"
	exitCode = 0
	main()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on invalid profile, got %d", exitCode)
	}

	// 3b. Create profile error
	profile = "cpu"
	exitCode = 0
	osCreateProfile = func(name string) (*os.File, error) {
		return nil, fmt.Errorf("simulated create error")
	}
	main()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on create profile error, got %d", exitCode)
	}
	osCreateProfile = os.Create

	// 3c. StartCPUProfile error
	exitCode = 0
	startCPUProfile = func(w io.Writer) error {
		return fmt.Errorf("simulated cpu profile start error")
	}
	main()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on start cpu profile error, got %d", exitCode)
	}
	startCPUProfile = pprof.StartCPUProfile
	_ = os.Remove("cpu.prof")

	// 3d. Normal CPU profile
	iterations = 2
	main()
	_ = os.Remove("cpu.prof")

	// 3e. Mem profile with write error
	profile = "mem"
	writeHeapProfile = func(w io.Writer) error {
		return fmt.Errorf("simulated heap write error")
	}
	main()
	_ = os.Remove("mem.prof")

	// 3f. Normal mem profile
	writeHeapProfile = pprof.WriteHeapProfile
	main()
	_ = os.Remove("mem.prof")

	// 3g. Block profile with write error
	profile = "block"
	writeBlockProfile = func(p *pprof.Profile, w io.Writer) error {
		return fmt.Errorf("simulated block write error")
	}
	main()
	_ = os.Remove("block.prof")

	// 3h. Normal block profile
	writeBlockProfile = defaultWriteBlockProfile
	main()
	_ = os.Remove("block.prof")
	_ = defaultWriteBlockProfile(nil, nil)
	_ = defaultWriteBlockProfile(pprof.Lookup("block"), io.Discard)
	profile = ""

	// 4. Test Accuracy branches
	testAccuracy = true
	// 4a. Normal accuracy
	main()

	// 4b. Failed accuracy
	exitCode = 0
	forceFailAccuracy = true
	main()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on accuracy failure, got %d", exitCode)
	}
	forceFailAccuracy = false
	testAccuracy = false

	// 5. Test Parity branches
	compareMode = true
	verbose = true
	// 5a. Normal parity
	main()

	// 5b. Failed parity
	exitCode = 0
	forceFailTol = true
	main()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on parity failure, got %d", exitCode)
	}
	forceFailTol = false
	compareMode = false
	verbose = false

	// 6. Test Default benchmark run (all benchmarks, printResults)
	iterations = 2
	typeFilter = "all"
	checkRegressionFlag = false
	main()

	// 6a. checkRegression when flag is false
	checkRegressionFlag = false
	checkRegression(nil)

	// 6b. Regression check branches: regression detected
	checkRegressionFlag = true
	exitCode = 0
	checkRegression([]BenchmarkResult{
		{Type: "float64", Name: "Exp", Ratio: 99.0},
	})
	if exitCode != 1 {
		t.Errorf("expected exitCode 1 on benchmark regression, got %d", exitCode)
	}

	// 6c. Regression check branches: improvement detected
	exitCode = 0
	checkRegression([]BenchmarkResult{
		{Type: "float64", Name: "Exp", Ratio: 0.1},
	})
	if exitCode != 0 {
		t.Errorf("expected exitCode 0 on benchmark improvement, got %d", exitCode)
	}

	// 6d. Regression check branches: all benchmarks pass
	exitCode = 0
	checkRegression([]BenchmarkResult{
		{Type: "float64", Name: "Exp", Ratio: 1.10},
	})
	if exitCode != 0 {
		t.Errorf("expected exitCode 0 on all benchmarks passing, got %d", exitCode)
	}

	// 7. Individual typeFilter runs to cover typeFilter branches in runAllBenchmarks
	for _, tf := range []string{"int", "uint", "float32", "float64", "complex64", "complex128"} {
		typeFilter = tf
		_ = runAllBenchmarks()
	}

	// 8. Cover benchmarkBatch when iterations >= 100000 and batchIterations < 2000
	iterations = 100000
	_ = benchmarkBatch("batch", "Test", 4096, func(a, b []float64) {}, func(a, b []float64) {})
	iterations = 2
}

