package trig_test

import (
	"fmt"
	"math"

	"github.com/emlgo/eml/pkg/trig"
)

func ExampleSin() {
	fmt.Printf("sin(π/4) = %.6f\n", trig.Sin(math.Pi/4))
	// Output: sin(π/4) = 0.707107
}

func ExampleCos() {
	fmt.Printf("cos(0) = %.1f\n", trig.Cos(0))
	// Output: cos(0) = 1.0
}

func ExampleTan() {
	fmt.Printf("tan(π/4) = %.6f\n", trig.Tan(math.Pi/4))
	// Output: tan(π/4) = 1.000000
}

func ExampleAsin() {
	fmt.Printf("asin(1) = %.6f\n", trig.Asin(1.0))
	// Output: asin(1) = 1.570796
}

func ExampleAtan2() {
	fmt.Printf("atan2(1, 1) = %.6f\n", trig.Atan2(1.0, 1.0))
	// Output: atan2(1, 1) = 0.785398
}

func ExampleSinBatch() {
	x := []float64{0, math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2}
	result := trig.SinBatch(x)
	for i, v := range result {
		fmt.Printf("sin(%d*π/%.0f) = %.4f\n", i+1, 12.0/float64(i+1), v)
	}
	// Output:
	// sin(1*π/12) = 0.0000
	// sin(2*π/6) = 0.5000
	// sin(3*π/4) = 0.7071
	// sin(4*π/3) = 0.8660
	// sin(5*π/2) = 1.0000
}
