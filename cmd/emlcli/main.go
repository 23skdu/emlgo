package main

import (
	"fmt"
	"os"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

func main() {
	fmt.Println("EML (Exp-Minus-Log) Library Demo")
	fmt.Println("==================================")

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

	fmt.Println("\nDone!")
	os.Exit(0)
}