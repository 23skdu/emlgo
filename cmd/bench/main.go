package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

var (
	iterations   int
	compareMode  bool
	verbose      bool
	testAccuracy bool
)

func init() {
	flag.IntVar(&iterations, "n", 1000000, "Number of iterations per benchmark")
	flag.BoolVar(&compareMode, "compare", false, "Compare emlgo vs math library correctness")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&testAccuracy, "accuracy", false, "Test accuracy (ULP) against math library")
}

type BenchmarkResult struct {
	Name      string
	EmlgoTime float64
	MathTime  float64
	Ratio     float64
	MaxULP    int
	Passed    bool
}

func main() {
	flag.Parse()

	rand.Seed(42)

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

	fmt.Printf("=== Speed Benchmark (n=%d) ===\n", iterations)
	results := runBenchmarks()
	printResults(results)
}

func runBenchmarks() []BenchmarkResult {
	results := []BenchmarkResult{}

	results = append(results, benchmarkFunc("Exp", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = logexp.Exp(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Exp(x)
		}
	}))

	results = append(results, benchmarkFunc("Log", func() {
		for i := 1; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			if x > 0 {
				_ = logexp.Log(x)
			}
		}
	}, func() {
		for i := 1; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			if x > 0 {
				_ = math.Log(x)
			}
		}
	}))

	results = append(results, benchmarkFunc("Sin", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = trig.Sin(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Sin(x)
		}
	}))

	results = append(results, benchmarkFunc("Cos", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = trig.Cos(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Cos(x)
		}
	}))

	results = append(results, benchmarkFunc("Tan", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = trig.Tan(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Tan(x)
		}
	}))

	results = append(results, benchmarkFunc("Sinh", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = hyper.Sinh(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Sinh(x)
		}
	}))

	results = append(results, benchmarkFunc("Cosh", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = hyper.Cosh(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Cosh(x)
		}
	}))

	results = append(results, benchmarkFunc("Tanh", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = hyper.Tanh(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Tanh(x)
		}
	}))

	results = append(results, benchmarkFunc("Asinh", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = hyper.Asinh(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%1000) / 100.0
			_ = math.Asinh(x)
		}
	}))

	results = append(results, benchmarkFunc("Acosh", func() {
		for i := 1; i < iterations; i++ {
			x := float64(i%1000)/100.0 + 1.0
			_ = hyper.Acosh(x)
		}
	}, func() {
		for i := 1; i < iterations; i++ {
			x := float64(i%1000)/100.0 + 1.0
			_ = math.Acosh(x)
		}
	}))

	results = append(results, benchmarkFunc("Atanh", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%998) / 1000.0 * 0.99
			_ = hyper.Atanh(x)
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%998) / 1000.0 * 0.99
			_ = math.Atanh(x)
		}
	}))

	results = append(results, benchmarkFunc("Sqrt", func() {
		for i := 1; i < iterations; i++ {
			x := float64(i % 10000)
			_ = arithmetic.Sqrt(x)
		}
	}, func() {
		for i := 1; i < iterations; i++ {
			x := float64(i % 10000)
			_ = math.Sqrt(x)
		}
	}))

	results = append(results, benchmarkFunc("Pow", func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%100) / 10.0
			y := float64(i%20) / 5.0
			if x > 0 {
				_ = arithmetic.Pow(x, y)
			}
		}
	}, func() {
		for i := 0; i < iterations; i++ {
			x := float64(i%100) / 10.0
			y := float64(i%20) / 5.0
			if x > 0 {
				_ = math.Pow(x, y)
			}
		}
	}))

	results = append(results, benchmarkFunc("LogBase2", func() {
		for i := 1; i < iterations; i++ {
			x := float64(i % 10000)
			if x > 0 {
				_ = arithmetic.LogBase2(x)
			}
		}
	}, func() {
		for i := 1; i < iterations; i++ {
			x := float64(i % 10000)
			if x > 0 {
				_ = math.Log2(x)
			}
		}
	}))

	results = append(results, benchmarkFunc("LogBase10", func() {
		for i := 1; i < iterations; i++ {
			x := float64(i % 10000)
			if x > 0 {
				_ = arithmetic.LogBase10(x)
			}
		}
	}, func() {
		for i := 1; i < iterations; i++ {
			x := float64(i % 10000)
			if x > 0 {
				_ = math.Log10(x)
			}
		}
	}))

	return results
}

func benchmarkFunc(name string, emlgoFunc, mathFunc func()) BenchmarkResult {
	start := time.Now()
	emlgoFunc()
	emlgoTime := time.Since(start).Seconds()

	start = time.Now()
	mathFunc()
	mathTime := time.Since(start).Seconds()

	ratio := emlgoTime / mathTime

	return BenchmarkResult{
		Name:      name,
		EmlgoTime: emlgoTime,
		MathTime:  mathTime,
		Ratio:     ratio,
		Passed:    true,
	}
}

func printResults(results []BenchmarkResult) {
	fmt.Printf("\n%-12s %12s %12s %10s\n", "Function", "emlgo (s)", "math (s)", "Ratio")
	fmt.Println(strings.Repeat("-", 50))

	totalRatio := 0.0
	for _, r := range results {
		fmt.Printf("%-12s %12.4f %12.4f %9.2fx\n", r.Name, r.EmlgoTime, r.MathTime, r.Ratio)
		totalRatio += r.Ratio
	}

	avgRatio := totalRatio / float64(len(results))
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("\nAverage ratio: %.2fx (emlgo vs math)\n", avgRatio)
	fmt.Printf("Note: Ratio > 1 means emlgo is slower\n")
}

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

	passed := 0
	failed := 0

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

func testExpParity() bool {
	for i := -100; i <= 100; i++ {
		x := float64(i) / 10.0
		e := logexp.Exp(x)
		m := math.Exp(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Exp(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testLogParity() bool {
	for i := 1; i <= 1000; i++ {
		x := float64(i) / 10.0
		e := logexp.Log(x)
		m := math.Log(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Log(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testSinParity() bool {
	for i := 0; i <= 1000; i++ {
		x := float64(i) / 100.0 * math.Pi
		e := trig.Sin(x)
		m := math.Sin(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Sin(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testCosParity() bool {
	for i := 0; i <= 1000; i++ {
		x := float64(i) / 100.0 * math.Pi
		e := trig.Cos(x)
		m := math.Cos(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Cos(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testTanParity() bool {
	for i := 0; i <= 1000; i++ {
		x := float64(i) / 100.0 * math.Pi
		if math.Abs(math.Cos(x)) > 1e-10 {
			e := trig.Tan(x)
			m := math.Tan(x)
			if !withinTolerance(e, m, 1e-10) {
				if verbose {
					fmt.Printf("  Tan(%f): emlgo=%f, math=%f\n", x, e, m)
				}
				return false
			}
		}
	}
	return true
}

func testSinhParity() bool {
	for i := -100; i <= 100; i++ {
		x := float64(i) / 10.0
		e := hyper.Sinh(x)
		m := math.Sinh(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Sinh(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testCoshParity() bool {
	for i := -100; i <= 100; i++ {
		x := float64(i) / 10.0
		e := hyper.Cosh(x)
		m := math.Cosh(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Cosh(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testTanhParity() bool {
	for i := -100; i <= 100; i++ {
		x := float64(i) / 10.0
		e := hyper.Tanh(x)
		m := math.Tanh(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Tanh(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testAsinhParity() bool {
	for i := -100; i <= 100; i++ {
		x := float64(i) / 10.0
		e := hyper.Asinh(x)
		m := math.Asinh(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Asinh(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testAcoshParity() bool {
	for i := 1; i <= 100; i++ {
		x := float64(i)/10.0 + 1.0
		e := hyper.Acosh(x)
		m := math.Acosh(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Acosh(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testAtanhParity() bool {
	for i := -99; i <= 99; i++ {
		x := float64(i) / 100.0
		e := hyper.Atanh(x)
		m := math.Atanh(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Atanh(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testSqrtParity() bool {
	for i := 0; i <= 10000; i++ {
		x := float64(i)
		e := arithmetic.Sqrt(x)
		m := math.Sqrt(x)
		if !withinTolerance(e, m, 1e-10) {
			if verbose {
				fmt.Printf("  Sqrt(%f): emlgo=%f, math=%f\n", x, e, m)
			}
			return false
		}
	}
	return true
}

func testPowParity() bool {
	for i := 0; i <= 1000; i++ {
		x := float64(i%100) / 10.0
		y := float64(i%20) / 5.0
		if x > 0 {
			e := arithmetic.Pow(x, y)
			m := math.Pow(x, y)
			if !withinTolerance(e, m, 1e-9) {
				if verbose {
					fmt.Printf("  Pow(%f,%f): emlgo=%f, math=%f\n", x, y, e, m)
				}
				return false
			}
		}
	}
	return true
}

func withinTolerance(a, b, tol float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return true
	}
	if math.IsInf(a, -1) && math.IsInf(b, -1) {
		return true
	}
	diff := math.Abs(a - b)
	sumAbs := math.Abs(a) + math.Abs(b) + 1e-10
	return diff < tol || diff/sumAbs < tol
}

func testAllAccuracy() {
	tests := []struct {
		name string
		fn   func() (maxULP int, passed bool)
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

	passed := 0
	failed := 0

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

func testExpAccuracy() (int, bool) {
	maxULP := 0
	for i := -1000; i <= 1000; i++ {
		x := float64(i) / 100.0
		e := logexp.Exp(x)
		m := math.Exp(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testLogAccuracy() (int, bool) {
	maxULP := 0
	for i := 1; i <= 10000; i++ {
		x := float64(i) / 10.0
		e := logexp.Log(x)
		m := math.Log(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testSinAccuracy() (int, bool) {
	maxULP := 0
	for i := 0; i <= 10000; i++ {
		x := float64(i) / 1000.0 * 4 * math.Pi
		e := trig.Sin(x)
		m := math.Sin(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testCosAccuracy() (int, bool) {
	maxULP := 0
	for i := 0; i <= 10000; i++ {
		x := float64(i) / 1000.0 * 4 * math.Pi
		e := trig.Cos(x)
		m := math.Cos(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testTanAccuracy() (int, bool) {
	maxULP := 0
	for i := 0; i <= 10000; i++ {
		x := float64(i) / 1000.0 * 4 * math.Pi
		if math.Abs(math.Cos(x)) > 1e-6 {
			e := trig.Tan(x)
			m := math.Tan(x)
			ulp := ulpDiff(e, m)
			if ulp > maxULP {
				maxULP = ulp
			}
			if ulp > 10 {
				return maxULP, false
			}
		}
	}
	return maxULP, true
}

func testSinhAccuracy() (int, bool) {
	maxULP := 0
	for i := -1000; i <= 1000; i++ {
		x := float64(i) / 100.0
		e := hyper.Sinh(x)
		m := math.Sinh(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testCoshAccuracy() (int, bool) {
	maxULP := 0
	for i := -1000; i <= 1000; i++ {
		x := float64(i) / 100.0
		e := hyper.Cosh(x)
		m := math.Cosh(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testTanhAccuracy() (int, bool) {
	maxULP := 0
	for i := -1000; i <= 1000; i++ {
		x := float64(i) / 100.0
		e := hyper.Tanh(x)
		m := math.Tanh(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testAsinhAccuracy() (int, bool) {
	maxULP := 0
	for i := -1000; i <= 1000; i++ {
		x := float64(i) / 100.0
		e := hyper.Asinh(x)
		m := math.Asinh(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testAcoshAccuracy() (int, bool) {
	maxULP := 0
	for i := 1; i <= 1000; i++ {
		x := float64(i)/100.0 + 1.0
		e := hyper.Acosh(x)
		m := math.Acosh(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testAtanhAccuracy() (int, bool) {
	maxULP := 0
	for i := -999; i <= 999; i++ {
		x := float64(i) / 1000.0
		e := hyper.Atanh(x)
		m := math.Atanh(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testSqrtAccuracy() (int, bool) {
	maxULP := 0
	for i := 0; i <= 100000; i++ {
		x := float64(i)
		e := arithmetic.Sqrt(x)
		m := math.Sqrt(x)
		ulp := ulpDiff(e, m)
		if ulp > maxULP {
			maxULP = ulp
		}
		if ulp > 10 {
			return maxULP, false
		}
	}
	return maxULP, true
}

func testPowAccuracy() (int, bool) {
	maxULP := 0
	for i := 0; i <= 5000; i++ {
		x := float64(i%100) / 10.0
		y := float64(i%20) / 5.0
		if x > 0 {
			e := arithmetic.Pow(x, y)
			m := math.Pow(x, y)
			ulp := ulpDiff(e, m)
			if ulp > maxULP {
				maxULP = ulp
			}
			if ulp > 10 {
				return maxULP, false
			}
		}
	}
	return maxULP, true
}

func ulpDiff(a, b float64) int {
	if a == b {
		return 0
	}
	if math.IsNaN(a) || math.IsNaN(b) {
		return 0
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return 0
	}

	bits := math.Float64bits(a)
	targetBits := math.Float64bits(b)

	var diff uint64
	if bits > targetBits {
		diff = bits - targetBits
		if diff > math.MaxInt64 {
			return math.MaxInt
		}
		return int(diff)
	}
	diff = targetBits - bits
	if diff > math.MaxInt64 {
		return math.MaxInt
	}
	return int(diff)
}