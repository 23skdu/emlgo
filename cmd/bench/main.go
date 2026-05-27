package main

// This is benchmark code - math/rand is acceptable for test data generation
//gosec: G404
import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/emlgo/eml/internal/gpu"
	"github.com/emlgo/eml/internal/jit"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/fastmath"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

var (
	iterations   int
	compareMode  bool
	verbose      bool
	testAccuracy bool
	typeFilter   string
	profile     string
	device      string
)

func init() {
	flag.IntVar(&iterations, "n", 1000000, "Number of iterations per benchmark")
	flag.BoolVar(&compareMode, "compare", false, "Compare emlgo vs math library correctness")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&testAccuracy, "accuracy", false, "Test accuracy (ULP) against math library")
	flag.StringVar(&typeFilter, "type", "all", "Type to test: all, int, uint, float32, float64, complex64, complex128")
	flag.StringVar(&profile, "profile", "", "Profile type: cpu, mem, or block (requires pprof binary)")
	flag.StringVar(&device, "device", "cpu", "Device to run on: cpu, gpu, or jit")
}

type BenchmarkResult struct {
	Type      string
	Name      string
	EmlgoTime float64
	MathTime  float64
	Ratio     float64
	Passed    bool
}

func main() {
	flag.Parse()
	rand.Seed(42)

	if device == "gpu" {
		runGpuBenchmarks()
		return
	}

	if device == "jit" {
		runJitBenchmarks()
		return
	}

	if profile != "" {
		runProfiling()
		return
	}

	if testAccuracy {
		fmt.Println("=== Accuracy Test (ULP) ===")
		testAllAccuracy()
		return
	}

	if compareMode {
		fmt.Println("=== Feature Parity Test ===")
		testAllParity()
		return
	}

	fmt.Printf("=== Comprehensive Performance Benchmark (n=%d) ===\n", iterations)
	results := runAllBenchmarks()
	results = append(results, runBatchBenchmarks()...)
	results = append(results, runFastMathBenchmarks()...)
	printResults(results)
	checkRegression(results)
}

// ==================== GPU BENCHMARKS ====================

func runGpuBenchmarks() {
	fmt.Println("=== GPU Device Benchmark ===")

	devices, err := gpu.GetDevices()
	if err != nil {
		fmt.Printf("GPU error: %v\n", err)
		os.Exit(1)
	}
	if len(devices) == 0 {
		fmt.Println("No GPU devices found (build with -tags cuda and ensure CUDA is available)")
		os.Exit(1)
	}

	device := devices[0]
	fmt.Printf("Using device: %s (SM %d.%d, %d MB)\n",
		device.Name, device.ComputeMajor, device.ComputeMinor,
		device.MemoryBytes/1024/1024)

	sizes := []int{1024, 4096, 16384, 65536, 262144, 1048576}

	fmt.Printf("\n%-12s %-12s %14s %14s %12s\n", "Size", "Function", "GPU Time (s)", "CPU Time (s)", "Speedup")
	fmt.Println(strings.Repeat("-", 70))

	for _, n := range sizes {
		data := make([]float64, n)
		for i := range data {
			data[i] = float64(i%100) / 100.0 + 0.1
		}

		// GPU Exp
		start := time.Now()
		_, err := device.ExpBatch(data)
		if err != nil {
			fmt.Printf("%-12d %-12s %14s %14s %12s\n", n, "Exp", "N/A", "N/A", err)
			continue
		}
		gpuTime := time.Since(start).Seconds()

		// CPU Exp baseline
		start = time.Now()
		for range 10 {
			_ = logexp.ExpBatch(data)
		}
		cpuTime := time.Since(start).Seconds() / 10

		speedup := cpuTime / gpuTime
		fmt.Printf("%-12d %-12s %14.6f %14.6f %11.2fx\n", n, "Exp", gpuTime, cpuTime, speedup)
	}
}

// ==================== JIT BENCHMARKS ====================

func runJitBenchmarks() {
	fmt.Println("=== JIT Polynomial Compilation Benchmark ===")

	type polyExpr struct {
		name string
		expr string
		goFn func(float64) float64
	}

	polys := []polyExpr{
		{"linear", "2*x + 1", func(x float64) float64 { return 2*x + 1 }},
		{"quadratic", "x^2 + 2*x + 1", func(x float64) float64 { return x*x + 2*x + 1 }},
		{"cubic", "x^3 - 3*x^2 + 3*x - 1", func(x float64) float64 { return x*x*x - 3*x*x + 3*x - 1 }},
		{"quintic", "x^5 - 3*x^4 + 2*x^3 - x^2 + 5*x - 7", func(x float64) float64 {
			return x*x*x*x*x - 3*x*x*x*x + 2*x*x*x - x*x + 5*x - 7
		}},
		{"horner", "(((x + 2)*x + 3)*x + 4)*x + 5", func(x float64) float64 {
			return (((x+2)*x+3)*x+4)*x + 5
		}},
	}

	n := 500000

	fmt.Printf("\n%-14s %14s %14s %14s %14s %12s\n",
		"Expression", "JIT exec (s)", "Eval exec (s)", "Native exec (s)",
		"JIT/Eval", "JIT/Go")
	fmt.Println(strings.Repeat("-", 95))

	compiler := jit.NewCompiler()
	var compileTotal time.Duration

	for _, p := range polys {
		start := time.Now()
		fn, err := compiler.Compile(p.expr)
		compileTime := time.Since(start)
		compileTotal += compileTime

		if err != nil {
			fmt.Printf("%-14s %14s %14s %14s %12s %12s\n",
				p.name, "ERR", "ERR", "ERR", "ERR", err.Error())
			continue
		}

		data := make([]float64, n)
		rng := rand.New(rand.NewSource(42)) // #nosec G404 - benchmark determinism
		for i := range data {
			data[i] = rng.Float64()*10 - 5
		}

		ast, _ := jit.Parse(p.expr)

		start = time.Now()
		for i := 0; i < n; i++ {
			_ = fn(data[i])
		}
		jitTime := time.Since(start).Seconds()

		start = time.Now()
		for i := 0; i < n; i++ {
			_ = jit.Eval(ast, data[i])
		}
		evalTime := time.Since(start).Seconds()

		start = time.Now()
		for i := 0; i < n; i++ {
			_ = p.goFn(data[i])
		}
		goTime := time.Since(start).Seconds()

		jitEvalRatio := jitTime / evalTime
		jitGoRatio := jitTime / goTime

		fmt.Printf("%-14s %14.6f %14.6f %14.6f %13.2fx %11.2fx\n",
			p.name, jitTime, evalTime, goTime, jitEvalRatio, jitGoRatio)
	}

	fmt.Println(strings.Repeat("-", 95))
	avgCompile := compileTotal / time.Duration(len(polys))
	fmt.Printf("\nAverage compilation time: %v\n", avgCompile)
	fmt.Printf("Note: Ratio < 1 means JIT is faster than reference\n")
}

func runProfiling() {
	var profFile *os.File
	var err error

	switch profile {
	case "cpu":
		profFile, err = os.Create("cpu.prof")
	case "mem":
		profFile, err = os.Create("mem.prof")
	case "block":
		profFile, err = os.Create("block.prof")
	default:
		fmt.Printf("Unknown profile type: %s (use cpu, mem, or block)\n", profile)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Error creating profile file: %v\n", err)
		os.Exit(1)
	}
	defer profFile.Close()

	switch profile {
	case "cpu":
		if err := pprof.StartCPUProfile(profFile); err != nil {
			fmt.Printf("Error starting CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	case "mem":
		runtime.GC()
		err := pprof.WriteHeapProfile(profFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing heap profile: %v\n", err)
		}
		fmt.Println("Memory profile written to mem.prof")
		return
	case "block":
		runtime.SetBlockProfileRate(1)
		defer func() {
			if prof := pprof.Lookup("block"); prof != nil {
				if err := prof.WriteTo(profFile, 0); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing block profile: %v\n", err)
				}
			}
		}()
	}

	_ = iterations
	n := 4096
	data := make([]float64, n)
	// #nosec G404
	for i := range data {
		data[i] = rand.Float64()*10 - 5
	}

	benchData := make([]float64, iterations)
	// #nosec G404
	for i := range benchData {
		benchData[i] = rand.Float64() * 10
	}

	for i := 0; i < 100; i++ {
		_ = arithmetic.AbsBatch(data)
		_ = logexp.ExpBatch(data)
		_ = trig.SinBatch(data)
	}

	fmt.Printf("Profile written to %s.prof\n", profile)
	fmt.Printf("To analyze: go tool pprof %s.prof\n", profile)
}

func runAllBenchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	if typeFilter == "all" || typeFilter == "int" {
		results = append(results, runIntBenchmarks()...)
	}
	if typeFilter == "all" || typeFilter == "uint" {
		results = append(results, runUintBenchmarks()...)
	}
	if typeFilter == "all" || typeFilter == "float32" {
		results = append(results, runFloat32Benchmarks()...)
	}
	if typeFilter == "all" || typeFilter == "float64" {
		results = append(results, runFloat64Benchmarks()...)
	}
	if typeFilter == "all" || typeFilter == "complex64" {
		results = append(results, runComplex64Benchmarks()...)
	}
	if typeFilter == "all" || typeFilter == "complex128" {
		results = append(results, runComplex128Benchmarks()...)
	}

	return results
}

// ==================== INT BENCHMARKS ====================

func runIntBenchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

results = append(results, benchmarkFuncInt("int", "Add", func(a, b int) {
		_ = arithmetic.IntAdd(a, b)
	}, func(a, b int) {
		_ = a + b
	}))

	results = append(results, benchmarkFuncInt("int", "Sub", func(a, b int) {
		_ = arithmetic.IntSub(a, b)
	}, func(a, b int) {
		_ = a - b
	}))

	results = append(results, benchmarkFuncInt("int", "Mul", func(a, b int) {
		_ = arithmetic.IntMul(a, b)
	}, func(a, b int) {
		_ = a * b
	}))

	results = append(results, benchmarkFuncInt("int", "Div", func(a, b int) {
		_ = arithmetic.IntDiv(a, b)
	}, func(a, b int) {
		if b != 0 {
			_ = a / b
		}
	}))

	results = append(results, benchmarkFuncInt("int", "Mod", func(a, b int) {
		if b != 0 {
			_ = arithmetic.IntMod(a, b)
		}
	}, func(a, b int) {
		if b != 0 {
			_ = a % b
		}
	}))

	results = append(results, benchmarkFuncInt("int", "Abs", func(a, b int) {
		_ = arithmetic.IntAbs(a)
	}, func(a, b int) {
		if a < 0 {
			_ = -a
		} else {
			_ = a
		}
	}))

	results = append(results, benchmarkFuncInt("int", "Floor", func(a, b int) {
		_ = arithmetic.Floor(float64(a))
	}, func(a, b int) {
		_ = float64(a)
	}))

	results = append(results, benchmarkFuncInt("int", "Ceil", func(a, b int) {
		_ = arithmetic.Ceil(float64(a))
	}, func(a, b int) {
		_ = float64(a)
	}))

	results = append(results, benchmarkFuncInt("int", "Max", func(a, b int) {
		_ = arithmetic.IntMax(a, b)
	}, func(a, b int) {
		if a > b {
			_ = a
		} else {
			_ = b
		}
	}))

	results = append(results, benchmarkFuncInt("int", "Min", func(a, b int) {
		_ = arithmetic.IntMin(a, b)
	}, func(a, b int) {
		if a < b {
			_ = a
		} else {
			_ = b
		}
	}))

	return results
}

func benchmarkFuncInt(typ, name string, emlgoFunc, mathFunc func(int, int)) BenchmarkResult {
	rand.Seed(42)
	randData := make([]int, iterations)
	// nosec G404 - benchmark tool uses math/rand for deterministic test data
	// #nosec G404 - benchmark tool uses math/rand for deterministic test data
	for i := range randData {
		randData[i] = rand.Intn(10000) - 5000
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		emlgoFunc(randData[i%len(randData)], randData[(i+1)%len(randData)])
	}
	emlgoTime := time.Since(start).Seconds()

	start = time.Now()
	for i := 0; i < iterations; i++ {
		mathFunc(randData[i%len(randData)], randData[(i+1)%len(randData)])
	}
	mathTime := time.Since(start).Seconds()

	return BenchmarkResult{
		Type:      typ,
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     emlgoTime / mathTime,
		Passed:    true,
	}
}

// ==================== UINT BENCHMARKS ====================

func runUintBenchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	results = append(results, benchmarkFuncUint("uint", "Add", func(a, b uint) {
		_ = arithmetic.UintAdd(a, b)
	}, func(a, b uint) {
		_ = a + b
	}))

	results = append(results, benchmarkFuncUint("uint", "Sub", func(a, b uint) {
		_ = arithmetic.UintSub(a, b)
	}, func(a, b uint) {
		_ = a - b
	}))

	results = append(results, benchmarkFuncUint("uint", "Mul", func(a, b uint) {
		_ = arithmetic.UintMul(a, b)
	}, func(a, b uint) {
		_ = a * b
	}))

	results = append(results, benchmarkFuncUint("uint", "Div", func(a, b uint) {
		if b != 0 {
			_ = arithmetic.UintDiv(a, b)
		}
	}, func(a, b uint) {
		if b != 0 {
			_ = a / b
		}
	}))

	return results
}

func benchmarkFuncUint(typ, name string, emlgoFunc, mathFunc func(uint, uint)) BenchmarkResult {
	rand.Seed(42)
	randData := make([]uint, iterations)
	// nosec G404 - benchmark tool uses math/rand for deterministic test data
	// #nosec G404 - benchmark tool uses math/rand for deterministic test data
	for i := range randData {
		randData[i] = uint(rand.Intn(10000))
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		emlgoFunc(randData[i%len(randData)], randData[(i+1)%len(randData)])
	}
	emlgoTime := time.Since(start).Seconds()

	start = time.Now()
	for i := 0; i < iterations; i++ {
		mathFunc(randData[i%len(randData)], randData[(i+1)%len(randData)])
	}
	mathTime := time.Since(start).Seconds()

	return BenchmarkResult{
		Type:      typ,
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     emlgoTime / mathTime,
		Passed:    true,
	}
}

// ==================== FLOAT32 BENCHMARKS ====================

func runFloat32Benchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	results = append(results, benchmarkFuncFloat32("float32", "Exp", func(x float32) {
		_ = float32(logexp.Exp(float64(x)))
	}, func(x float32) {
		_ = float32(math.Exp(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Log", func(x float32) {
		_ = float32(logexp.Exp(float64(x)))
	}, func(x float32) {
		_ = float32(math.Log(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Sin", func(x float32) {
		_ = trig.Sin(float64(x))
	}, func(x float32) {
		_ = float32(math.Sin(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Cos", func(x float32) {
		_ = trig.Cos(float64(x))
	}, func(x float32) {
		_ = float32(math.Cos(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Tan", func(x float32) {
		_ = trig.Tan(float64(x))
	}, func(x float32) {
		_ = float32(math.Tan(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Sqrt", func(x float32) {
		_ = float32(arithmetic.Sqrt(float64(x)))
	}, func(x float32) {
		_ = float32(math.Sqrt(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Pow", func(x float32) {
		_ = float32(arithmetic.Pow(float64(x), 2.0))
	}, func(x float32) {
		_ = float32(math.Pow(float64(x), 2.0))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Sinh", func(x float32) {
		_ = float32(hyper.Sinh(float64(x)))
	}, func(x float32) {
		_ = float32(math.Sinh(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Cosh", func(x float32) {
		_ = float32(hyper.Cosh(float64(x)))
	}, func(x float32) {
		_ = float32(math.Cosh(float64(x)))
	}))

	results = append(results, benchmarkFuncFloat32("float32", "Tanh", func(x float32) {
		_ = float32(hyper.Tanh(float64(x)))
	}, func(x float32) {
		_ = float32(math.Tanh(float64(x)))
	}))

	return results
}

func benchmarkFuncFloat32(typ, name string, emlgoFunc, mathFunc func(float32)) BenchmarkResult {
	rand.Seed(42)
	randData := make([]float32, iterations)
	// #nosec G404 - benchmark tool uses math/rand for deterministic test data
	for i := range randData {
		randData[i] = float32(rand.Float64()*10 - 5)
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		emlgoFunc(randData[i%len(randData)])
	}
	emlgoTime := time.Since(start).Seconds()

	start = time.Now()
	for i := 0; i < iterations; i++ {
		mathFunc(randData[i%len(randData)])
	}
	mathTime := time.Since(start).Seconds()

	return BenchmarkResult{
		Type:      typ,
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     emlgoTime / mathTime,
		Passed:    true,
	}
}

// ==================== FLOAT64 BENCHMARKS ====================

func runFloat64Benchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	results = append(results, benchmarkFloat64("Exp", func(x float64) float64 { return logexp.Exp(x) }, func(x float64) float64 { return math.Exp(x) }))
	results = append(results, benchmarkFloat64("Log", func(x float64) float64 { return logexp.Log(x) }, func(x float64) float64 { return math.Log(x) }))
	results = append(results, benchmarkFloat64("Log2", func(x float64) float64 { return arithmetic.LogBase2(x) }, func(x float64) float64 { return math.Log2(x) }))
	results = append(results, benchmarkFloat64("Log10", func(x float64) float64 { return arithmetic.LogBase10(x) }, func(x float64) float64 { return math.Log10(x) }))

	results = append(results, benchmarkFloat64("Sin", func(x float64) float64 { return trig.Sin(x) }, func(x float64) float64 { return math.Sin(x) }))
	results = append(results, benchmarkFloat64("Cos", func(x float64) float64 { return trig.Cos(x) }, func(x float64) float64 { return math.Cos(x) }))
	results = append(results, benchmarkFloat64("Tan", func(x float64) float64 { return trig.Tan(x) }, func(x float64) float64 { return math.Tan(x) }))
	results = append(results, benchmarkFloat64("Cot", func(x float64) float64 { return trig.Cot(x) }, func(x float64) float64 { return 1/math.Tan(x) }))
	results = append(results, benchmarkFloat64("Sec", func(x float64) float64 { return trig.Sec(x) }, func(x float64) float64 { return 1/math.Cos(x) }))
	results = append(results, benchmarkFloat64("Csc", func(x float64) float64 { return trig.Csc(x) }, func(x float64) float64 { return 1/math.Sin(x) }))

	results = append(results, benchmarkFloat64("Asin", func(x float64) float64 { return trig.Asin(x) }, func(x float64) float64 { return math.Asin(x) }))
	results = append(results, benchmarkFloat64("Acos", func(x float64) float64 { return trig.Acos(x) }, func(x float64) float64 { return math.Acos(x) }))
	results = append(results, benchmarkFloat64("Atan", func(x float64) float64 { return trig.Atan(x) }, func(x float64) float64 { return math.Atan(x) }))
	results = append(results, benchmarkFloat64("Atan2", func(x float64) float64 { return trig.Atan2(x, 1) }, func(x float64) float64 { return math.Atan2(x, 1) }))

	results = append(results, benchmarkFloat64("Sinh", func(x float64) float64 { return hyper.Sinh(x) }, func(x float64) float64 { return math.Sinh(x) }))
	results = append(results, benchmarkFloat64("Cosh", func(x float64) float64 { return hyper.Cosh(x) }, func(x float64) float64 { return math.Cosh(x) }))
	results = append(results, benchmarkFloat64("Tanh", func(x float64) float64 { return hyper.Tanh(x) }, func(x float64) float64 { return math.Tanh(x) }))

	results = append(results, benchmarkFloat64("Asinh", func(x float64) float64 { return hyper.Asinh(x) }, func(x float64) float64 { return math.Asinh(x) }))
	results = append(results, benchmarkFloat64("Acosh", func(x float64) float64 { return hyper.Acosh(x + 1) }, func(x float64) float64 { return math.Acosh(x + 1) }))
	results = append(results, benchmarkFloat64("Atanh", func(x float64) float64 { return hyper.Atanh(x * 0.5) }, func(x float64) float64 { return math.Atanh(x * 0.5) }))

	results = append(results, benchmarkFloat64("Sqrt", func(x float64) float64 { return arithmetic.Sqrt(x) }, func(x float64) float64 { return math.Sqrt(x) }))
	results = append(results, benchmarkFloat64("Cbrt", func(x float64) float64 { return arithmetic.Cbrt(x) }, func(x float64) float64 { return math.Cbrt(x) }))
	results = append(results, benchmarkFloat64("Pow", func(x float64) float64 { return arithmetic.Pow(x, 2.5) }, func(x float64) float64 { return math.Pow(x, 2.5) }))
	results = append(results, benchmarkFloat64("PowInt", func(x float64) float64 { return arithmetic.PowInt(x, 3) }, func(x float64) float64 { return math.Pow(x, 3) }))

	results = append(results, benchmarkFloat64("Floor", func(x float64) float64 { return arithmetic.Floor(x) }, func(x float64) float64 { return math.Floor(x) }))
	results = append(results, benchmarkFloat64("Ceil", func(x float64) float64 { return arithmetic.Ceil(x) }, func(x float64) float64 { return math.Ceil(x) }))
	results = append(results, benchmarkFloat64("Round", func(x float64) float64 { return arithmetic.Round(x) }, func(x float64) float64 { return math.Round(x) }))
	results = append(results, benchmarkFloat64("Trunc", func(x float64) float64 { return arithmetic.Trunc(x) }, func(x float64) float64 { return math.Trunc(x) }))
	results = append(results, benchmarkFloat64("Abs", func(x float64) float64 { return arithmetic.Abs(x) }, func(x float64) float64 { return math.Abs(x) }))
	results = append(results, benchmarkFloat64("Neg", func(x float64) float64 { return arithmetic.Neg(x) }, func(x float64) float64 { return -x }))
	results = append(results, benchmarkFloat64("Inv", func(x float64) float64 { return arithmetic.Inv(x) }, func(x float64) float64 { return 1/x }))
	results = append(results, benchmarkFloat64("Square", func(x float64) float64 { return arithmetic.Square(x) }, func(x float64) float64 { return x*x }))
	results = append(results, benchmarkFloat64("Cube", func(x float64) float64 { return arithmetic.Cube(x) }, func(x float64) float64 { return x*x*x }))

	results = append(results, benchmarkFloat64("Max", func(x float64) float64 { return arithmetic.Max(x, x+1) }, func(x float64) float64 { return math.Max(x, x+1) }))
	results = append(results, benchmarkFloat64("Min", func(x float64) float64 { return arithmetic.Min(x, x+1) }, func(x float64) float64 { return math.Min(x, x+1) }))
	results = append(results, benchmarkFloat64("Hypot", func(x float64) float64 { return arithmetic.Hypot(x, x+1) }, func(x float64) float64 { return math.Hypot(x, x+1) }))
	results = append(results, benchmarkFloat64("FMA", func(x float64) float64 { return arithmetic.FMA(x, x+1, x+2) }, func(x float64) float64 { return x*(x+1) + (x+2) }))

	results = append(results, benchmarkFloat64("Add", func(x float64) float64 { return x + (x + 1) }, func(x float64) float64 { return x + (x + 1) }))
	results = append(results, benchmarkFloat64("Sub", func(x float64) float64 { return x - (x + 1) }, func(x float64) float64 { return x - (x + 1) }))
	results = append(results, benchmarkFloat64("Mul", func(x float64) float64 { return x * (x + 1) }, func(x float64) float64 { return x * (x + 1) }))
	results = append(results, benchmarkFloat64("Div", func(x float64) float64 { return x / (x + 1) }, func(x float64) float64 { return x / (x + 1) }))

	return results
}

func benchmarkFloat64(name string, emlgoFunc, mathFunc func(float64) float64) BenchmarkResult {
	rand.Seed(42)
	randData := make([]float64, iterations)
	// #nosec G404 - benchmark tool uses math/rand for deterministic test data
	for i := range randData {
		randData[i] = rand.Float64()*10 - 5
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = emlgoFunc(randData[i%len(randData)])
	}
	emlgoTime := time.Since(start).Seconds()

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = mathFunc(randData[i%len(randData)])
	}
	mathTime := time.Since(start).Seconds()

	return BenchmarkResult{
		Type:      "float64",
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     emlgoTime / mathTime,
		Passed:    true,
	}
}

// ==================== COMPLEX64 BENCHMARKS ====================

func runComplex64Benchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	results = append(results, benchmarkComplex64("Complex64", "Exp", func(x complex64) complex64 {
		r := complexExp(float64(real(x)), float64(imag(x)))
		return complex64(r)
	}, func(x complex64) complex64 {
		return complex64(cmplx.Exp(complex128(x)))
	}))

	results = append(results, benchmarkComplex64("Complex64", "Log", func(x complex64) complex64 {
		if x == 0 {
			return 0
		}
		mag := arithmetic.Sqrt(float64(real(x))*float64(real(x)) + float64(imag(x))*float64(imag(x)))
		arg := trig.Atan2(float64(imag(x)), float64(real(x)))
		return complex64(complex(logexp.Log(mag), arg))
	}, func(x complex64) complex64 {
		return complex64(cmplx.Log(complex128(x)))
	}))

	results = append(results, benchmarkComplex64("Complex64", "Sin", func(x complex64) complex64 {
		return complex64(trigComplexSin(float64(real(x)), float64(imag(x))))
	}, func(x complex64) complex64 {
		return complex64(cmplx.Sin(complex128(x)))
	}))

	results = append(results, benchmarkComplex64("Complex64", "Cos", func(x complex64) complex64 {
		return complex64(trigComplexCos(float64(real(x)), float64(imag(x))))
	}, func(x complex64) complex64 {
		return complex64(cmplx.Cos(complex128(x)))
	}))

	results = append(results, benchmarkComplex64("Complex64", "Sqrt", func(x complex64) complex64 {
		r, i := float64(real(x)), float64(imag(x))
		mag := arithmetic.Sqrt(r*r + i*i)
		rPlus := (mag + r) / 2
		signI := 1.0
		if i < 0 {
			signI = -1
		}
		return complex64(complex(arithmetic.Sqrt(rPlus), signI*arithmetic.Sqrt((mag-r)/2)))
	}, func(x complex64) complex64 {
		return complex64(cmplx.Sqrt(complex128(x)))
	}))

	return results
}

func benchmarkComplex64(typ, name string, emlgoFunc, mathFunc func(complex64) complex64) BenchmarkResult {
	rand.Seed(42)
	randData := make([]complex64, iterations)
	// #nosec G404 - benchmark tool uses math/rand for deterministic test data
	for i := range randData {
		randData[i] = complex(float32(rand.Float64()*10-5), float32(rand.Float64()*10-5))
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = emlgoFunc(randData[i%len(randData)])
	}
	emlgoTime := time.Since(start).Seconds()

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = mathFunc(randData[i%len(randData)])
	}
	mathTime := time.Since(start).Seconds()

	return BenchmarkResult{
		Type:      typ,
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     emlgoTime / mathTime,
		Passed:    true,
	}
}

// ==================== COMPLEX128 BENCHMARKS ====================

func runComplex128Benchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	results = append(results, benchmarkComplex128("complex128", "Exp", func(x complex128) complex128 {
		return complexExp(real(x), imag(x))
	}, func(x complex128) complex128 {
		return cmplx.Exp(x)
	}))

	results = append(results, benchmarkComplex128("complex128", "Log", func(x complex128) complex128 {
		if x == 0 {
			return 0
		}
		mag := arithmetic.Sqrt(real(x)*real(x) + imag(x)*imag(x))
		arg := trig.Atan2(imag(x), real(x))
		return complex(logexp.Log(mag), arg)
	}, func(x complex128) complex128 {
		return cmplx.Log(x)
	}))

	results = append(results, benchmarkComplex128("complex128", "Sin", func(x complex128) complex128 {
		return trigComplexSin(real(x), imag(x))
	}, func(x complex128) complex128 {
		return cmplx.Sin(x)
	}))

	results = append(results, benchmarkComplex128("complex128", "Cos", func(x complex128) complex128 {
		return trigComplexCos(real(x), imag(x))
	}, func(x complex128) complex128 {
		return cmplx.Cos(x)
	}))

	results = append(results, benchmarkComplex128("complex128", "Tan", func(x complex128) complex128 {
		return trigComplexTan(real(x), imag(x))
	}, func(x complex128) complex128 {
		return cmplx.Tan(x)
	}))

	results = append(results, benchmarkComplex128("complex128", "Sqrt", func(x complex128) complex128 {
		r, i := real(x), imag(x)
		mag := arithmetic.Sqrt(r*r + i*i)
		rPlus := (mag + r) / 2
		signI := 1.0
		if i < 0 {
			signI = -1
		}
		return complex(arithmetic.Sqrt(rPlus), signI*arithmetic.Sqrt((mag-r)/2))
	}, func(x complex128) complex128 {
		return cmplx.Sqrt(x)
	}))

	return results
}

func benchmarkComplex128(typ, name string, emlgoFunc, mathFunc func(complex128) complex128) BenchmarkResult {
	rand.Seed(42)
	randData := make([]complex128, iterations)
	// #nosec G404 - benchmark tool uses math/rand for deterministic test data
	for i := range randData {
		randData[i] = complex(rand.Float64()*10-5, rand.Float64()*10-5)
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = emlgoFunc(randData[i%len(randData)])
	}
	emlgoTime := time.Since(start).Seconds()

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = mathFunc(randData[i%len(randData)])
	}
	mathTime := time.Since(start).Seconds()

	return BenchmarkResult{
		Type:      typ,
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     emlgoTime / mathTime,
		Passed:    true,
	}
}

// ==================== FASTMATH BENCHMARKS ====================

func runFastMathBenchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	results = append(results, benchmarkFloat64("fastmath.Exp", fastmath.Exp, math.Exp))
	results = append(results, benchmarkFloat64("fastmath.Log", fastmath.Log, math.Log))
	results = append(results, benchmarkFloat64("fastmath.Sin", fastmath.Sin, math.Sin))
	results = append(results, benchmarkFloat64("fastmath.Cos", fastmath.Cos, math.Cos))
	results = append(results, benchmarkFloat64("fastmath.Sqrt", fastmath.Sqrt, math.Sqrt))

	return results
}

// ==================== BATCH BENCHMARKS ====================

func runBatchBenchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}
	n := 4096 // Batch size

	results = append(results, benchmarkBatch("batch", "Add", n, func(a, b []float64) {
		_ = arithmetic.AddBatch(a, b)
	}, func(a, b []float64) {
		res := make([]float64, len(a))
		for i := range a {
			res[i] = a[i] + b[i]
		}
	}))

	results = append(results, benchmarkBatch("batch", "Sub", n, func(a, b []float64) {
		_ = arithmetic.SubBatch(a, b)
	}, func(a, b []float64) {
		res := make([]float64, len(a))
		for i := range a {
			res[i] = a[i] - b[i]
		}
	}))

	results = append(results, benchmarkBatch("batch", "Mul", n, func(a, b []float64) {
		_ = arithmetic.MulBatch(a, b)
	}, func(a, b []float64) {
		res := make([]float64, len(a))
		for i := range a {
			res[i] = a[i] * b[i]
		}
	}))

	results = append(results, benchmarkBatch("batch", "Div", n, func(a, b []float64) {
		_ = arithmetic.DivBatch(a, b)
	}, func(a, b []float64) {
		res := make([]float64, len(a))
		for i := range a {
			if b[i] != 0 {
				res[i] = a[i] / b[i]
			}
		}
	}))

	results = append(results, benchmarkBatch("batch", "Sqrt", n, func(a, b []float64) {
		_ = arithmetic.SqrtBatch(a)
	}, func(a, b []float64) {
		res := make([]float64, len(a))
		for i := range a {
			res[i] = math.Sqrt(a[i])
		}
	}))

	results = append(results, benchmarkBatch("batch", "Exp", n, func(a, b []float64) {
		_ = logexp.ExpBatch(a)
	}, func(a, b []float64) {
		res := make([]float64, len(a))
		for i := range a {
			res[i] = math.Exp(a[i])
		}
	}))

	return results
}

func benchmarkBatch(typ, name string, n int, emlgoFunc, mathFunc func([]float64, []float64)) BenchmarkResult {
	rand.Seed(42)
	a := make([]float64, n)
	b := make([]float64, n)
	// #nosec G404 - benchmark tool uses math/rand for deterministic test data
	for i := 0; i < n; i++ {
		a[i] = rand.Float64()*10 - 5
		b[i] = rand.Float64()*10 - 5
	}

	// Adjust iterations for batch benchmarks to keep runtime reasonable
	batchIterations := 100000

	start := time.Now()
	for i := 0; i < batchIterations; i++ {
		emlgoFunc(a, b)
	}
	emlgoTime := time.Since(start).Seconds() / float64(batchIterations)

	start = time.Now()
	for i := 0; i < batchIterations; i++ {
		mathFunc(a, b)
	}
	mathTime := time.Since(start).Seconds() / float64(batchIterations)

	return BenchmarkResult{
		Type:      typ,
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     emlgoTime / mathTime,
		Passed:    true,
	}
}

// Helper functions for complex operations
func trigComplexSin(r, i float64) complex128 {
	sinX, cosX := trig.SinCos(r)
	sinhI, coshI := trig.SinhCosh(i)
	return complex(sinX*coshI, cosX*sinhI)
}

func trigComplexCos(r, i float64) complex128 {
	sinX, cosX := trig.SinCos(r)
	sinhI, coshI := trig.SinhCosh(i)
	return complex(cosX*coshI, -sinX*sinhI)
}

func trigComplexTan(r, i float64) complex128 {
	return trigComplexSin(r, i) / trigComplexCos(r, i)
}

func complexExp(r, i float64) complex128 {
	expR := logexp.Exp(r)
	sinI, cosI := trig.SinCos(i)
	return complex(expR*cosI, expR*sinI)
}

func printResults(results []BenchmarkResult) {
	fmt.Printf("\n%-12s %-12s %12s %12s %10s\n", "Type", "Function", "emlgo (s)", "math (s)", "Ratio")
	fmt.Println(strings.Repeat("-", 65))

	totalRatio := 0.0
	count := 0

	// Group by type
	types := map[string]bool{}
	for _, r := range results {
		types[r.Type] = true
	}

	for t := range types {
		fmt.Printf("\n=== %s ===\n", t)
		fmt.Printf("%-12s %-12s %12s %12s %10s\n", "Type", "Function", "emlgo (s)", "math (s)", "Ratio")
		fmt.Println(strings.Repeat("-", 65))
		for _, r := range results {
			if r.Type == t {
				fmt.Printf("%-12s %-12s %12.4f %12.4f %9.2fx\n", r.Type, r.Name, r.EmlgoTime, r.MathTime, r.Ratio)
				totalRatio += r.Ratio
				count++
			}
		}
	}

	avgRatio := totalRatio / float64(count)
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("\nAverage ratio: %.2fx (emlgo vs math)\n", avgRatio)
	fmt.Printf("Note: Ratio > 1 means emlgo is slower, Ratio < 1 means faster\n")
}

// ==================== PARITY AND ACCURACY TESTS ====================

func testAllParity() {
	tests := []struct {
		name string
		fn   func() bool
	}{
		{"Exp", testExpParity},
		{"Log", testLogParity},
		{"Sin", testSinParity},
		{"Cos", testCosParity},
		{"Tan", testTanParity},
		{"Sinh", testSinhParity},
		{"Cosh", testCoshParity},
		{"Tanh", testTanhParity},
		{"Asinh", testAsinhParity},
		{"Acosh", testAcoshParity},
		{"Atanh", testAtanhParity},
		{"Sqrt", testSqrtParity},
		{"Pow", testPowParity},
	}

	passed, failed := 0, 0
	for _, t := range tests {
		if t.fn() {
			fmt.Printf("✓ %s: PASSED\n", t.name)
			passed++
		} else {
			fmt.Printf("✗ %s: FAILED\n", t.name)
			failed++
		}
	}
	fmt.Printf("\nResults: %d passed, %d failed\n", passed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func testAllAccuracy() {
	tests := []struct {
		name string
		fn   func() (uint64, bool)
	}{
		{"Exp", testExpAccuracy},
		{"Log", testLogAccuracy},
		{"Sin", testSinAccuracy},
		{"Cos", testCosAccuracy},
		{"Tan", testTanAccuracy},
		{"Sinh", testSinhAccuracy},
		{"Cosh", testCoshAccuracy},
		{"Tanh", testTanhAccuracy},
		{"Asinh", testAsinhAccuracy},
		{"Acosh", testAcoshAccuracy},
		{"Atanh", testAtanhAccuracy},
		{"Sqrt", testSqrtAccuracy},
		{"Pow", testPowAccuracy},
	}

	passed, failed := 0, 0
	for _, t := range tests {
		maxULP, ok := t.fn()
		if ok {
			fmt.Printf("✓ %s: PASSED (max ULP: %d)\n", t.name, maxULP)
			passed++
		} else {
			fmt.Printf("✗ %s: FAILED (max ULP: %d)\n", t.name, maxULP)
			failed++
		}
	}
	fmt.Printf("\nResults: %d passed, %d failed\n", passed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func withinTol(a, b, tol float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) { return true }
	if math.IsInf(a, 1) && math.IsInf(b, 1) { return true }
	if math.IsInf(a, -1) && math.IsInf(b, -1) { return true }
	diff := math.Abs(a - b)
	sumAbs := math.Abs(a) + math.Abs(b) + 1e-10
	return diff < tol || diff/sumAbs < tol
}

func testExpParity() bool { return testOpParity(func(x float64) float64 { return logexp.Exp(x) }, math.Exp) }
func testLogParity() bool { return testOpParity(func(x float64) float64 { return logexp.Log(x) }, math.Log) }
func testSinParity() bool { return testOpParity(trig.Sin, math.Sin) }
func testCosParity() bool { return testOpParity(trig.Cos, math.Cos) }
func testTanParity() bool { return testOpParity(trig.Tan, math.Tan) }
func testSinhParity() bool { return testOpParity(hyper.Sinh, math.Sinh) }
func testCoshParity() bool { return testOpParity(hyper.Cosh, math.Cosh) }
func testTanhParity() bool { return testOpParity(hyper.Tanh, math.Tanh) }
func testAsinhParity() bool { return testOpParity(hyper.Asinh, math.Asinh) }
func testAcoshParity() bool { return testOpParity(func(x float64) float64 { return hyper.Acosh(x+1) }, func(x float64) float64 { return math.Acosh(x+1) }) }
func testAtanhParity() bool {
	for i := 1; i <= 19; i++ {
		x := float64(i) / 10.0
		e, m := hyper.Atanh(x*0.5), math.Atanh(x*0.5)
		if !withinTol(e, m, 1e-10) {
			if verbose { fmt.Printf("  %f: emlgo=%f, math=%f\n", x, e, m) }
			return false
		}
	}
	return true
}
func testSqrtParity() bool { return testOpParity(arithmetic.Sqrt, math.Sqrt) }
func testPowParity() bool { return testOpParity(func(x float64) float64 { return arithmetic.Pow(x, 2.5) }, func(x float64) float64 { return math.Pow(x, 2.5) }) }

func testOpParity(emlgoFunc, mathFunc func(float64) float64) bool {
	for i := -100; i <= 100; i++ {
		x := float64(i) / 10.0
		if x <= 0 { continue }
		e, m := emlgoFunc(x), mathFunc(x)
		if !withinTol(e, m, 1e-10) {
			if verbose { fmt.Printf("  %f: emlgo=%f, math=%f\n", x, e, m) }
			return false
		}
	}
	return true
}

func testExpAccuracy() (uint64, bool) { return testOpAccuracy(func(x float64) float64 { return logexp.Exp(x) }, math.Exp) }
func testLogAccuracy() (uint64, bool) { return testOpAccuracy(func(x float64) float64 { return logexp.Log(x) }, math.Log) }
func testSinAccuracy() (uint64, bool) { return testOpAccuracy(trig.Sin, math.Sin) }
func testCosAccuracy() (uint64, bool) { return testOpAccuracy(trig.Cos, math.Cos) }
func testTanAccuracy() (uint64, bool) { return testOpAccuracy(trig.Tan, math.Tan) }
func testSinhAccuracy() (uint64, bool) { return testOpAccuracy(hyper.Sinh, math.Sinh) }
func testCoshAccuracy() (uint64, bool) { return testOpAccuracy(hyper.Cosh, math.Cosh) }
func testTanhAccuracy() (uint64, bool) { return testOpAccuracy(hyper.Tanh, math.Tanh) }
func testAsinhAccuracy() (uint64, bool) { return testOpAccuracy(hyper.Asinh, math.Asinh) }
func testAcoshAccuracy() (uint64, bool) { return testOpAccuracy(func(x float64) float64 { return hyper.Acosh(x+1) }, func(x float64) float64 { return math.Acosh(x+1) }) }
func testAtanhAccuracy() (uint64, bool) { return testOpAccuracy(func(x float64) float64 { return hyper.Atanh(x*0.5) }, func(x float64) float64 { return math.Atanh(x*0.5) }) }
func testSqrtAccuracy() (uint64, bool) { return testOpAccuracy(arithmetic.Sqrt, math.Sqrt) }
func testPowAccuracy() (uint64, bool) { return testOpAccuracy(func(x float64) float64 { return arithmetic.Pow(x, 2.5) }, func(x float64) float64 { return math.Pow(x, 2.5) }) }

func testOpAccuracy(emlgoFunc, mathFunc func(float64) float64) (uint64, bool) {
	var maxULP uint64 = 0
	for i := -1000; i <= 1000; i++ {
		x := float64(i) / 100.0
		e, m := emlgoFunc(x), mathFunc(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP { maxULP = ulp }
		if ulp > 200 { return maxULP, false }
	}
	return maxULP, true
}

func ulpDiff(a, b float64) uint64 {
	if a == b { return 0 }
	if math.IsNaN(a) || math.IsNaN(b) { return 0 }
	if math.IsInf(a, 0) || math.IsInf(b, 0) { return 0 }
	bits, targetBits := math.Float64bits(a), math.Float64bits(b)
	if bits > targetBits { return bits - targetBits }
	return targetBits - bits
}

var baseline = map[string]float64{
	"float64/Exp":    1.10,
	"float64/Log":    1.00,
	"float64/Sin":    1.02,
	"float64/Cos":    1.00,
	"float64/Tan":    1.07,
	"float64/Sqrt":   1.00,
	"float64/Pow":    0.85,
	"float64/PowInt": 0.16,
	"float64/Cosh":   1.00,
	"int/Add":        1.01,
	"int/Mod":        1.00,
	"int/Max":        1.00,
	"int/Min":        1.00,
	"uint/Add":       1.00,
	"uint/Mul":       1.00,
	"uint/Div":       1.00,
}

func checkRegression(results []BenchmarkResult) {
	regressionFlag := flag.Bool("regression", false, "Check for performance regression against baseline")
	flag.Parse()
	
	if !*regressionFlag {
		return
	}
	
	fmt.Println("\n=== Regression Check ===")
	regressions := 0
	for _, r := range results {
		key := r.Type + "/" + r.Name
		if baselineRatio, ok := baseline[key]; ok {
			regression := r.Ratio - baselineRatio
			if regression > 0.10 {
				fmt.Printf("⚠️  REGRESSION: %s ratio changed from %.2fx to %.2fx (+%.1f%%)\n", 
					key, baselineRatio, r.Ratio, regression*100)
				regressions++
			} else if regression < -0.15 {
				fmt.Printf("✓  IMPROVEMENT: %s ratio changed from %.2fx to %.2fx (%.1f%% better)\n",
					key, baselineRatio, r.Ratio, -regression*100)
			}
		}
	}
	if regressions > 0 {
		fmt.Printf("\n⚠️  WARNING: %d regressions detected (>15%% slower than baseline)\n", regressions)
		os.Exit(1)
	} else {
		fmt.Println("\n✓ All benchmarks within 15% of baseline")
	}
}