package arithmetic_test

import (
	"fmt"
	"math"

	"github.com/emlgo/eml/pkg/arithmetic"
)

func ExampleAdd() {
	result := arithmetic.Add(3.0, 4.0)
	fmt.Printf("Add(3, 4) = %.1f\n", result)
	// Output: Add(3, 4) = 7.0
}

func ExamplePow() {
	result := arithmetic.Pow(2.0, 10.0)
	fmt.Printf("Pow(2, 10) = %.0f\n", result)
	// Output: Pow(2, 10) = 1024
}

func ExamplePow_numericalStability() {
	// Pow near x=1 uses Log1p for numerical stability
	result := arithmetic.Pow(1.0+1e-15, 1e6)
	fmt.Printf("Pow(1+1e-15, 1e6) = %.6f\n", result)
	_ = result // approximately 1.001
}

func ExampleSqrt() {
	result := arithmetic.Sqrt(144.0)
	fmt.Printf("Sqrt(144) = %.0f\n", result)
	// Output: Sqrt(144) = 12
}

func ExampleGCD() {
	result := arithmetic.GCD(12, 18)
	fmt.Printf("GCD(12, 18) = %d\n", result)
	// Output: GCD(12, 18) = 6
}

func ExampleLCM() {
	result := arithmetic.LCM(4, 6)
	fmt.Printf("LCM(4, 6) = %d\n", result)
	// Output: LCM(4, 6) = 12
}

func ExampleAddBatch() {
	x := []float64{1.0, 2.0, 3.0, 4.0}
	y := []float64{10.0, 20.0, 30.0, 40.0}
	result := arithmetic.AddBatch(x, y)
	fmt.Printf("AddBatch: %v\n", result)
	// Output: AddBatch: [11 22 33 44]
}

func ExampleFMA() {
	// Fused multiply-add: a*b + c computed with a single rounding
	result := arithmetic.FMA(3.0, 4.0, 5.0)
	fmt.Printf("FMA(3, 4, 5) = %.0f\n", result)
	// Output: FMA(3, 4, 5) = 17
}

func ExampleRound() {
	values := []float64{-1.5, -1.4, -0.5, 0.5, 1.4, 1.5}
	for _, v := range values {
		fmt.Printf("Round(%+.1f) = %+.0f\n", v, arithmetic.Round(v))
	}
	// Output:
	// Round(-1.5) = -2
	// Round(-1.4) = -1
	// Round(-0.5) = -1
	// Round(+0.5) = +1
	// Round(+1.4) = +1
	// Round(+1.5) = +2
}

func ExampleHypot() {
	result := arithmetic.Hypot(3.0, 4.0)
	fmt.Printf("Hypot(3, 4) = %.1f\n", result)
	_ = math.Sqrt(25) // 5.0
	// Output: Hypot(3, 4) = 5.0
}
