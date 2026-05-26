package main

import (
	"fmt"
	"os"

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
	fmt.Println("  gpu-bench   Run a quick GPU batch performance test")
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

	device := devices[0]
	fmt.Printf("Running GPU benchmark on %s...\n", device.Name)

	sizes := []int{1024, 16384, 262144, 1048576}
	for _, n := range sizes {
		data := make([]float64, n)
		for i := range data {
			data[i] = float64(i%100) / 100.0
		}

		result, err := device.ExpBatch(data)
		if err != nil {
			fmt.Printf("  n=%d: ExpBatch failed: %v\n", n, err)
			continue
		}
		// Verify a few results
		ok := true
		for i := 0; i < 10 && i < n; i++ {
			if result[i] == 0 {
				ok = false
				break
			}
		}
		if ok {
			fmt.Printf("  n=%d: ExpBatch OK (first=%.6f, last=%.6f)\n", n, result[0], result[n-1])
		} else {
			fmt.Printf("  n=%d: ExpBatch produced zeros - possible kernel issue\n", n)
		}
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
