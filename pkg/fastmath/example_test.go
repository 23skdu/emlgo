package fastmath_test

import (
	"fmt"
	"math"

	"github.com/emlgo/eml/pkg/fastmath"
)

func ExampleExp() {
	result := fastmath.Exp(1.0)
	fmt.Printf("fastmath.Exp(1) = %.6f\n", result)
	fmt.Printf("math.Exp(1)    = %.6f\n", math.Exp(1))
	// Output:
	// fastmath.Exp(1) = 2.718282
	// math.Exp(1)    = 2.718282
}

func ExampleLog() {
	result := fastmath.Log(math.E)
	fmt.Printf("fastmath.Log(e) = %.6f\n", result)
	fmt.Printf("math.Log(e)    = %.6f\n", math.Log(math.E))
	// Output:
	// fastmath.Log(e) = 1.000000
	// math.Log(e)    = 1.000000
}

func ExampleFastEml() {
	result := fastmath.FastEml(1.0, math.E)
	fmt.Printf("FastEml(1, e) = %.6f\n", result)
	// Output: FastEml(1, e) = 1.718282
}

func ExampleSin() {
	result := fastmath.Sin(math.Pi / 6)
	fmt.Printf("fastmath.Sin(π/6) = %.6f\n", result)
	// Output: fastmath.Sin(π/6) = 0.500000
}

func ExampleCos() {
	result := fastmath.Cos(0)
	fmt.Printf("fastmath.Cos(0) = %.1f\n", result)
	// Output: fastmath.Cos(0) = 1.0
}

func ExampleSqrt() {
	result := fastmath.Sqrt(144.0)
	fmt.Printf("fastmath.Sqrt(144) = %.0f\n", result)
	// Output: fastmath.Sqrt(144) = 12
}
