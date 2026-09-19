package hyper_test

import (
	"fmt"

	"github.com/emlgo/eml/pkg/hyper"
)

func ExampleSinh() {
	fmt.Printf("sinh(1) = %.6f\n", hyper.Sinh(1))
	// Output: sinh(1) = 1.175201
}

func ExampleCosh() {
	fmt.Printf("cosh(0) = %.1f\n", hyper.Cosh(0))
	// Output: cosh(0) = 1.0
}

func ExampleTanh() {
	fmt.Printf("tanh(2) = %.6f\n", hyper.Tanh(2))
	// Output: tanh(2) = 0.964028
}

func ExampleAsinh() {
	fmt.Printf("asinh(1) = %.6f\n", hyper.Asinh(1))
	// Output: asinh(1) = 0.881374
}

func ExampleSinhBatch() {
	x := []float64{0, 0.5, 1.0}
	result := hyper.SinhBatch(x)
	fmt.Printf("SinhBatch: [%.4f, %.4f, %.4f]\n", result[0], result[1], result[2])
	// Output: SinhBatch: [0.0000, 0.5211, 1.1752]
}
