package logexp_test

import (
	"fmt"
	"math"

	"github.com/emlgo/eml/pkg/logexp"
)

func ExampleExp() {
	fmt.Printf("exp(1) = %.6f\n", logexp.Exp(1))
	// Output: exp(1) = 2.718282
}

func ExampleLog() {
	fmt.Printf("log(e) = %.6f\n", logexp.Log(math.E))
	// Output: log(e) = 1.000000
}

func ExampleExpBatch() {
	x := []float64{0, 0.5, 1.0, 1.5, 2.0}
	result := logexp.ExpBatch(x)
	for i, v := range result {
		fmt.Printf("exp(%.1f) = %.4f\n", float64(i)*0.5, v)
	}
	// Output:
	// exp(0.0) = 1.0000
	// exp(0.5) = 1.6487
	// exp(1.0) = 2.7183
	// exp(1.5) = 4.4817
	// exp(2.0) = 7.3891
}

func ExampleExpFast() {
	fmt.Printf("exp(0) = %.1f\n", logexp.ExpFast(0))
	// Output: exp(0) = 1.0
}

func ExampleLogFast() {
	fmt.Printf("log(1) = %.1f\n", logexp.LogFast(1))
	// Output: log(1) = 0.0
}
