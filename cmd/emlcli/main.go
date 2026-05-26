package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/internal/gpu"
	"github.com/emlgo/eml/internal/jit"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "demo":
		runDemo()
	case "gpu-status":
		runGpuStatus()
	case "gpu-bench":
		runGpuBench()
	case "gpu-verify":
		runGpuVerify()
	case "jit-test":
		runJitTest()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("EML (Exp-Minus-Log) CLI Tool")
	fmt.Println("Usage: eml [command]")
	fmt.Println("\nCommands:")
	fmt.Println("  demo        Run library demo")
	fmt.Println("  gpu-status  Check GPU availability and status")
	fmt.Println("  gpu-bench   Run GPU vs CPU performance benchmarks")
	fmt.Println("  gpu-verify  Verify GPU results against math library (ULP)")
	fmt.Println("  jit-test    Test JIT polynomial compilation")
}

func runDemo() {
	fmt.Println("EML Library Demo")
	fmt.Println("==================")

	fmt.Println("\n--- Core EML Operator ---")
	fmt.Printf("eml(1, 1) = %v (should be e - 0 = e)\n", eml.Eml(1, 1))
	fmt.Printf("eml(2, 1) = %v (should be e^2)\n", eml.Eml(2, 1))
	fmt.Printf("eml(1, eml(1, 1)) = %v (should be 1/e)\n", eml.Eml(1, eml.Eml(1, 1)))

	fmt.Println("\n--- Exponential & Logarithmic ---")
	fmt.Printf("Exp(1) = %v (should be e)\n", logexp.Exp(1))
	fmt.Printf("Log(e) = %v (should be 1)\n", logexp.Log(2.718281828459045))

	fmt.Println("\n--- Trigonometric ---")
	fmt.Printf("Sin(0) = %v\n", trig.Sin(0))
	fmt.Printf("Cos(0) = %v\n", trig.Cos(0))
	fmt.Printf("Tan(0) = %v\n", trig.Tan(0))

	fmt.Println("\n--- Hyperbolic ---")
	fmt.Printf("Sinh(0) = %v\n", hyper.Sinh(0))
	fmt.Printf("Cosh(0) = %v\n", hyper.Cosh(0))
	fmt.Printf("Tanh(0) = %v\n", hyper.Tanh(0))

	fmt.Println("\n--- Arithmetic ---")
	fmt.Printf("Add(2, 3) = %v\n", arithmetic.Add(2, 3))
	fmt.Printf("Mul(2, 3) = %v\n", arithmetic.Mul(2, 3))
	fmt.Printf("Pow(2, 3) = %v\n", arithmetic.Pow(2, 3))
	fmt.Printf("Sqrt(4) = %v\n", arithmetic.Sqrt(4))

	fmt.Println("\n--- SIMD Support ---")
	fmt.Printf("Has AVX2: %v\n", eml.HasAVX2())
	fmt.Printf("Has AVX512: %v\n", eml.HasAVX512())
	fmt.Printf("Has Neon: %v\n", eml.HasNeon())
	fmt.Printf("Has WASM SIMD: %v\n", eml.HasWasmSIMD())

	fmt.Println("\nDone!")
}

func runGpuStatus() {
	fmt.Println(gpu.Status())

	devices, err := gpu.GetDevices()
	if err != nil {
		fmt.Printf("Error querying devices: %v\n", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("\nNo CUDA-capable GPUs detected.")
		fmt.Println("Tip: Build with -tags cuda and ensure libeml_capi.so is in your library path.")
		return
	}

	fmt.Println("\n--- GPU Details ---")
	for _, d := range devices {
		fmt.Printf("  Device %d: %s\n", d.ID, d.Name)
		fmt.Printf("    Compute Capability: %d.%d\n", d.ComputeMajor, d.ComputeMinor)
		fmt.Printf("    Memory: %.1f GB\n", float64(d.MemoryBytes)/1e9)
		fmt.Printf("    Max Threads/Block: %d\n", d.MaxThreadsPerBlock)
		fmt.Printf("    Warp Size: %d\n", d.WarpSize)
		fmt.Printf("    Clock Rate: %.2f GHz\n", float64(d.ClockRateKHz)/1e6)
	}
}

// GPU batch operations available for benchmarking/verification.
type gpuOp struct {
	Name string
	// RunGPU executes the operation on the GPU for the given device and input.
	RunGPU func(d *gpu.Device, x []float64) ([]float64, error)
	// RunCPU is the CPU reference (single-element) for verification.
	CPU func(float64) float64
	// CPU batch function for performance comparison.
	CPUBatch func(x []float64) []float64
}

var gpuOps = []gpuOp{
	{"Exp", (*gpu.Device).ExpBatch, math.Exp, logexp.ExpBatch},
	{"Log", (*gpu.Device).LogBatch, math.Log, func(x []float64) []float64 {
		r := make([]float64, len(x))
		for i := range x {
			r[i] = math.Log(x[i])
		}
		return r
	}},
	{"Sin", (*gpu.Device).SinBatch, math.Sin, trig.SinBatch},
	{"Cos", (*gpu.Device).CosBatch, math.Cos, trig.CosBatch},
	{"Tan", (*gpu.Device).TanBatch, math.Tan, trig.TanBatch},
	{"Sinh", (*gpu.Device).SinhBatch, math.Sinh, hyper.SinhBatch},
	{"Cosh", (*gpu.Device).CoshBatch, math.Cosh, hyper.CoshBatch},
	{"Tanh", (*gpu.Device).TanhBatch, math.Tanh, hyper.TanhBatch},
	{"Sqrt", (*gpu.Device).SqrtBatch, math.Sqrt, arithmetic.SqrtBatch},
}

func runGpuBench() {
	devices, err := gpu.GetDevices()
	if err != nil {
		fmt.Printf("GPU error: %v\n", err)
		return
	}
	if len(devices) == 0 {
		fmt.Println("No GPU devices available.")
		return
	}

	device := &devices[0]

	sizes := []int{1024, 4096, 16384, 65536, 262144}
	runs := 5

	fmt.Printf("Device: %s (SM %d.%d, %d MB)\n",
		device.Name, device.ComputeMajor, device.ComputeMinor,
		device.MemoryBytes/1024/1024)
	fmt.Printf("Runs per size: %d\n\n", runs)

	for _, op := range gpuOps {
		fmt.Printf("--- %s ---\n", op.Name)
		fmt.Printf("%-10s %14s %14s %10s\n", "Size", "GPU (s)", "CPU (s)", "Speedup")
		fmt.Println(strings.Repeat("-", 54))

		for _, n := range sizes {
			data := make([]float64, n)
			for i := range data {
				data[i] = float64(i%100+1) / 100.0
			}

			// GPU timing
			var gpuTotal time.Duration
			for range runs {
				start := time.Now()
				_, err := op.RunGPU(device, data)
				if err != nil {
					gpuTotal = 0
					break
				}
				gpuTotal += time.Since(start)
			}
			if gpuTotal == 0 {
				fmt.Printf("%-10d %14s %14s %10s\n", n, "N/A", "N/A", "N/A")
				continue
			}
			avgGPU := gpuTotal.Seconds() / float64(runs)

			// CPU timing
			start := time.Now()
			for range runs {
				_ = op.CPUBatch(data)
			}
			avgCPU := time.Since(start).Seconds() / float64(runs)

			speedup := avgCPU / avgGPU
			fmt.Printf("%-10d %14.6f %14.6f %9.2fx\n", n, avgGPU, avgCPU, speedup)
		}
		fmt.Println()
	}
}

func runGpuVerify() {
	devices, err := gpu.GetDevices()
	if err != nil {
		fmt.Printf("GPU error: %v\n", err)
		return
	}
	if len(devices) == 0 {
		fmt.Println("No GPU devices available.")
		return
	}

	device := &devices[0]
	verifier := gpu.DefaultVerifier()

	// Test points: mix of common ranges and edge cases.
	testPoints := []float64{
		0.0, 0.5, 1.0, 2.0, 10.0, 100.0,
		-0.5, -1.0, -2.0, -10.0,
		0.1, 0.01, 1e-10, 1e10,
		math.Pi, math.E, math.Ln2, math.Ln10,
		math.SmallestNonzeroFloat64,
		math.MaxFloat64,
		math.Inf(1), math.Inf(-1),
		math.NaN(),
	}

	fmt.Printf("Verifying GPU results on %s against Go math library...\n\n", device.Name)
	fmt.Println("ULP tolerance: 1")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-8s %-16s %-8s %s\n", "Op", "Value", "MaxULP", "Status")
	fmt.Println(strings.Repeat("-", 80))

	allPassed := true
	for _, op := range gpuOps {
		var maxULP uint64
		var anyFailure bool

		for _, x := range testPoints {
			input := []float64{x}
			result, err := op.RunGPU(device, input)
			if err != nil {
				fmt.Printf("%-8s %-16v %-8s FAIL (launch: %v)\n", op.Name, x, "-", err)
				allPassed = false
				continue
			}
			ulp, _, _ := verifier.VerifyOp(op.Name, input, result, op.CPU)
			if ulp > maxULP {
				maxULP = ulp
			}
			if ulp > 1 {
				anyFailure = true
				allPassed = false
			}
		}

		status := "PASS"
		if anyFailure {
			status = "FAIL"
		}
		fmt.Printf("%-8s %-16s %-8d %s\n", op.Name, "(all points)", maxULP, status)
	}

	fmt.Println(strings.Repeat("-", 80))
	if allPassed {
		fmt.Println("✓ All GPU results within 1 ULP of math library")
	} else {
		fmt.Println("✗ Some GPU results exceed 1 ULP tolerance")
	}
}

func runJitTest() {
	fmt.Println("JIT Polynomial Compilation Test")
	c := jit.NewCompiler()
	f, err := c.Compile("x^2 + 2x + 1")
	if err != nil {
		fmt.Printf("JIT Error: %v\n", err)
		return
	}
	fmt.Printf("JIT Compiled f(2) = %v\n", f(2))
}
