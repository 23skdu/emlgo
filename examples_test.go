package eml_test

import (
	"fmt"
	"math"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

func ExampleEml() {
	result := eml.Eml(1, 1)
	fmt.Printf("eml(1, 1) = %v (should be e - 0 = e)\n", result)
	fmt.Printf("math.E  = %v\n", math.E)
}

func ExampleEmlOne() {
	result := eml.EmlOne(2)
	fmt.Printf("eml(2, 1) = e^2 = %v\n", result)
	fmt.Printf("math.Exp(2) = %v\n", math.Exp(2))
}

func ExampleOneEml() {
	result := eml.OneEml(2)
	fmt.Printf("eml(1, 2) = -ln(2) = %v\n", result)
	fmt.Printf("-math.Log(2) = %v\n", -math.Log(2))
}

func ExampleExp() {
	result := logexp.Exp(1)
	fmt.Printf("Exp(1) = %v (should be e)\n", result)
}

func ExampleLog() {
	result := logexp.Log(math.E)
	fmt.Printf("Log(e) = %v (should be 1)\n", result)
}

func ExampleSin() {
	result := trig.Sin(math.Pi / 2)
	fmt.Printf("sin(π/2) = %v (should be 1)\n", result)
}

func ExampleCos() {
	result := trig.Cos(0)
	fmt.Printf("cos(0) = %v (should be 1)\n", result)
}

func ExampleTanh() {
	result := hyper.Tanh(0)
	fmt.Printf("tanh(0) = %v (should be 0)\n", result)
}

func ExampleAdd() {
	result := arithmetic.Add(2, 3)
	fmt.Printf("Add(2, 3) = %v\n", result)
}

func ExampleMul() {
	result := arithmetic.Mul(2, 3)
	fmt.Printf("Mul(2, 3) = %v\n", result)
}

func ExampleSqrt() {
	result := arithmetic.Sqrt(16)
	fmt.Printf("Sqrt(16) = %v\n", result)
}

func ExampleHasAVX2() {
	fmt.Printf("Has AVX2: %v\n", eml.HasAVX2())
	fmt.Printf("Has AVX512: %v\n", eml.HasAVX512())
}