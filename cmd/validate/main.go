package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"reflect"

	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

var (
	verbose      bool
	failedOnly   bool
	typeFilter   string
)

func init() {
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&failedOnly, "f", false, "Show only failed tests")
	flag.StringVar(&typeFilter, "type", "", "Filter by type (int, uint, float, complex)")
}

type ValidationResult struct {
	Type      string
	Function  string
	Passed    bool
	Message   string
}

var allResults []ValidationResult

func main() {
	flag.Parse()

	fmt.Println("=== EMLGO Validation Tests ===")
	fmt.Println("Testing all Go math data types...")
	fmt.Println()

	validateIntTypes()
	validateUintTypes()
	validateFloatTypes()
	validateComplexTypes()

	printSummary()
}

func validateIntTypes() {
	fmt.Println("--- Integer Types ---")

	types := []string{"int", "int8", "int16", "int32", "int64"}

	for _, t := range types {
		if typeFilter != "" && typeFilter != "int" {
			continue
		}

		result := validateIntType(t)
		allResults = append(allResults, result...)
	}
}

func validateUintTypes() {
	fmt.Println("--- Unsigned Integer Types ---")

	types := []string{"uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}

	for _, t := range types {
		if typeFilter != "" && typeFilter != "uint" {
			continue
		}

		result := validateUintType(t)
		allResults = append(allResults, result...)
	}
}

func validateFloatTypes() {
	fmt.Println("--- Float Types ---")

	types := []string{"float32", "float64"}

	for _, t := range types {
		if typeFilter != "" && typeFilter != "float" {
			continue
		}

		result := validateFloatType(t)
		allResults = append(allResults, result...)
	}
}

func validateComplexTypes() {
	fmt.Println("--- Complex Types ---")

	types := []string{"complex64", "complex128"}

	for _, t := range types {
		if typeFilter != "" && typeFilter != "complex" {
			continue
		}

		result := validateComplexType(t)
		allResults = append(allResults, result...)
	}
}

func validateIntType(typeName string) []ValidationResult {
	results := []ValidationResult{}

	switch typeName {
	case "int":
		results = append(results, testInt[int]("int")...)
	case "int8":
		results = append(results, testInt[int8]("int8")...)
	case "int16":
		results = append(results, testInt[int16]("int16")...)
	case "int32":
		results = append(results, testInt[int32]("int32")...)
	case "int64":
		results = append(results, testInt[int64]("int64")...)
	}

	if filterPassed(results) {
		fmt.Printf("  %s: PASSED\n", typeName)
	}

	return results
}

func validateUintType(typeName string) []ValidationResult {
	results := []ValidationResult{}

	switch typeName {
	case "uint":
		results = append(results, testUint[uint]("uint")...)
	case "uint8":
		results = append(results, testUint[uint8]("uint8")...)
	case "uint16":
		results = append(results, testUint[uint16]("uint16")...)
	case "uint32":
		results = append(results, testUint[uint32]("uint32")...)
	case "uint64":
		results = append(results, testUint[uint64]("uint64")...)
	case "uintptr":
		results = append(results, testUint[uintptr]("uintptr")...)
	}

	if filterPassed(results) {
		fmt.Printf("  %s: PASSED\n", typeName)
	}

	return results
}

func validateFloatType(typeName string) []ValidationResult {
	results := []ValidationResult{}

	switch typeName {
	case "float32":
		results = append(results, testFloat32()...)
	case "float64":
		results = append(results, testFloat64()...)
	}

	if filterPassed(results) {
		fmt.Printf("  %s: PASSED\n", typeName)
	}

	return results
}

func validateComplexType(typeName string) []ValidationResult {
	results := []ValidationResult{}

	switch typeName {
	case "complex64":
		results = append(results, testComplex64()...)
	case "complex128":
		results = append(results, testComplex128()...)
	}

	if filterPassed(results) {
		fmt.Printf("  %s: PASSED\n", typeName)
	}

	return results
}

func filterPassed(results []ValidationResult) bool {
	for _, r := range results {
		if !r.Passed {
			if verbose || failedOnly {
				fmt.Printf("  %s: %s - %s\n", r.Type, r.Function, r.Message)
			} else {
				fmt.Printf("  %s: FAILED (%s)\n", r.Type, r.Function)
			}
			return false
		}
	}
	return true
}

func testInt[T int | int8 | int16 | int32 | int64](typeName string) []ValidationResult {
	var results []ValidationResult

	// Test arithmetic operations with int types
	a := T(10)
	b := T(3)

	// Add
	res := arithmetic.Add(float64(a), float64(b))
	expected := float64(a) + float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Add", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Add", Passed: true, Message: "OK"})
	}

	// Sub
	res = arithmetic.Sub(float64(a), float64(b))
	expected = float64(a) - float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Sub", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Sub", Passed: true, Message: "OK"})
	}

	// Mul
	res = arithmetic.Mul(float64(a), float64(b))
	expected = float64(a) * float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Mul", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Mul", Passed: true, Message: "OK"})
	}

	// Div
	res = arithmetic.Div(float64(a), float64(b))
	expected = float64(a) / float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Div", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Div", Passed: true, Message: "OK"})
	}

	// Mod (only for signed ints)
	if reflect.TypeOf(a).Kind() == reflect.Int {
		res = arithmetic.Mod(float64(a), float64(b))
		expected = float64(int(a) % int(b))
		if !withinTol(res, expected, 0.001) {
			results = append(results, ValidationResult{Type: typeName, Function: "Mod", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
		} else {
			results = append(results, ValidationResult{Type: typeName, Function: "Mod", Passed: true, Message: "OK"})
		}
	}

	// Abs
	neg := T(-5)
	res = arithmetic.Abs(float64(neg))
	if res != 5 {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: false, Message: fmt.Sprintf("got %v, want 5", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: true, Message: "OK"})
	}

	// Floor/Ceil
	res = arithmetic.Floor(3.7)
	if res != 3 {
		results = append(results, ValidationResult{Type: typeName, Function: "Floor", Passed: false, Message: fmt.Sprintf("got %v, want 3", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Floor", Passed: true, Message: "OK"})
	}

	res = arithmetic.Ceil(3.2)
	if res != 4 {
		results = append(results, ValidationResult{Type: typeName, Function: "Ceil", Passed: false, Message: fmt.Sprintf("got %v, want 4", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Ceil", Passed: true, Message: "OK"})
	}

	// Round
	res = arithmetic.Round(3.5)
	if res != 4 {
		results = append(results, ValidationResult{Type: typeName, Function: "Round", Passed: false, Message: fmt.Sprintf("got %v, want 4", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Round", Passed: true, Message: "OK"})
	}

	// Trunc
	res = arithmetic.Trunc(3.7)
	if res != 3 {
		results = append(results, ValidationResult{Type: typeName, Function: "Trunc", Passed: false, Message: fmt.Sprintf("got %v, want 3", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Trunc", Passed: true, Message: "OK"})
	}

	// Max
	res = arithmetic.Max(float64(a), float64(b))
	if res != float64(a) {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, a)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: true, Message: "OK"})
	}

	// Min
	res = arithmetic.Min(float64(a), float64(b))
	if res != float64(b) {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, b)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: true, Message: "OK"})
	}

	// Neg
	res = arithmetic.Neg(float64(a))
	negExpected := -float64(a)
	if res != negExpected {
		results = append(results, ValidationResult{Type: typeName, Function: "Neg", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, negExpected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Neg", Passed: true, Message: "OK"})
	}

	// Inv
	res = arithmetic.Inv(float64(a))
	invExpected := 1 / float64(a)
	if !withinTol(res, invExpected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Inv", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, invExpected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Inv", Passed: true, Message: "OK"})
	}

	// Square
	res = arithmetic.Square(float64(a))
	squareVal := float64(a) * float64(a)
	if res != squareVal {
		results = append(results, ValidationResult{Type: typeName, Function: "Square", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, squareVal)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Square", Passed: true, Message: "OK"})
	}

	return results
}

func testUint[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](typeName string) []ValidationResult {
	var results []ValidationResult

	a := T(10)
	b := T(3)

	// Add
	res := arithmetic.Add(float64(a), float64(b))
	expected := float64(a) + float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Add", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Add", Passed: true, Message: "OK"})
	}

	// Sub
	res = arithmetic.Sub(float64(a), float64(b))
	expected = float64(a) - float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Sub", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Sub", Passed: true, Message: "OK"})
	}

	// Mul
	res = arithmetic.Mul(float64(a), float64(b))
	expected = float64(a) * float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Mul", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Mul", Passed: true, Message: "OK"})
	}

	// Div
	res = arithmetic.Div(float64(a), float64(b))
	expected = float64(a) / float64(b)
	if !withinTol(res, expected, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Div", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, expected)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Div", Passed: true, Message: "OK"})
	}

	// Remainder
	res = arithmetic.Remainder(float64(a), float64(b))
	m := math.Remainder(float64(a), float64(b))
	if !withinTol(res, m, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Remainder", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, m)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Remainder", Passed: true, Message: "OK"})
	}

	// Abs
	res = arithmetic.Abs(float64(a))
	if res != float64(a) {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, a)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: true, Message: "OK"})
	}

	// Max
	res = arithmetic.Max(float64(a), float64(b))
	if res != float64(a) {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, a)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: true, Message: "OK"})
	}

	// Min
	res = arithmetic.Min(float64(a), float64(b))
	if res != float64(b) {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, b)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: true, Message: "OK"})
	}

	// Square
	res = arithmetic.Square(float64(a))
	squareVal := float64(a) * float64(a)
	if res != squareVal {
		results = append(results, ValidationResult{Type: typeName, Function: "Square", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, squareVal)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Square", Passed: true, Message: "OK"})
	}

	return results
}

func testFloat32() []ValidationResult {
	results := []ValidationResult{}

	testCases := []float32{
		0, 1, -1, 0.5, -0.5, math.MaxFloat32, math.SmallestNonzeroFloat32,
		math.Pi, math.E, float32(math.Pow(2, 100)),
	}

	// Exp
	for _, x := range testCases {
		res := float32(logexp.Exp(float64(x)))
		expected := float32(math.Exp(float64(x)))
		if !withinTolFloat32(res, expected) {
			results = append(results, ValidationResult{Type: "float32", Function: "Exp", Passed: false, Message: fmt.Sprintf("Exp(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float32", Function: "Exp", Passed: true, Message: "OK"})
		}
	}

	// Log
	for _, x := range testCases {
		if x > 0 {
			res := float32(logexp.Log(float64(x)))
			expected := float32(math.Log(float64(x)))
			if !withinTolFloat32(res, expected) {
				results = append(results, ValidationResult{Type: "float32", Function: "Log", Passed: false, Message: fmt.Sprintf("Log(%v): got %v, want %v", x, res, expected)})
			} else {
				results = append(results, ValidationResult{Type: "float32", Function: "Log", Passed: true, Message: "OK"})
			}
		}
	}

	// Sin
	for _, x := range testCases {
		res := float32(trig.Sin(float64(x)))
		expected := float32(math.Sin(float64(x)))
		if !withinTolFloat32(res, expected) {
			results = append(results, ValidationResult{Type: "float32", Function: "Sin", Passed: false, Message: fmt.Sprintf("Sin(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float32", Function: "Sin", Passed: true, Message: "OK"})
		}
	}

	// Cos
	for _, x := range testCases {
		res := float32(trig.Cos(float64(x)))
		expected := float32(math.Cos(float64(x)))
		if !withinTolFloat32(res, expected) {
			results = append(results, ValidationResult{Type: "float32", Function: "Cos", Passed: false, Message: fmt.Sprintf("Cos(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float32", Function: "Cos", Passed: true, Message: "OK"})
		}
	}

	// Tan
	for _, x := range testCases {
		res := float32(trig.Tan(float64(x)))
		expected := float32(math.Tan(float64(x)))
		if !withinTolFloat32(res, expected) {
			results = append(results, ValidationResult{Type: "float32", Function: "Tan", Passed: false, Message: fmt.Sprintf("Tan(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float32", Function: "Tan", Passed: true, Message: "OK"})
		}
	}

	// Sqrt
	for _, x := range testCases {
		if x >= 0 {
			res := float32(arithmetic.Sqrt(float64(x)))
			expected := float32(math.Sqrt(float64(x)))
			if !withinTolFloat32(res, expected) {
				results = append(results, ValidationResult{Type: "float32", Function: "Sqrt", Passed: false, Message: fmt.Sprintf("Sqrt(%v): got %v, want %v", x, res, expected)})
			} else {
				results = append(results, ValidationResult{Type: "float32", Function: "Sqrt", Passed: true, Message: "OK"})
			}
		}
	}

	// Pow
	for i := 0; i < 10; i++ {
		x := float32(i) + 1
		y := float32(i) * 0.5
		res := float32(arithmetic.Pow(float64(x), float64(y)))
		expected := float32(math.Pow(float64(x), float64(y)))
		if !withinTolFloat32(res, expected) {
			results = append(results, ValidationResult{Type: "float32", Function: "Pow", Passed: false, Message: fmt.Sprintf("Pow(%v,%v): got %v, want %v", x, y, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float32", Function: "Pow", Passed: true, Message: "OK"})
		}
	}

	return results
}

func testFloat64() []ValidationResult {
	results := []ValidationResult{}

	testCases := []float64{
		0, 1, -1, 0.5, -0.5, math.MaxFloat64, math.SmallestNonzeroFloat64,
		math.Pi, math.E, math.Pow(2, 100), math.Pow(2, -100),
	}

	// Exp
	for _, x := range testCases {
		res := logexp.Exp(x)
		expected := math.Exp(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Exp", Passed: false, Message: fmt.Sprintf("Exp(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Exp", Passed: true, Message: "OK"})
		}
	}

	// Log
	for _, x := range testCases {
		if x > 0 {
			res := logexp.Log(x)
			expected := math.Log(x)
			if !withinTol(res, expected, 1e-10) {
				results = append(results, ValidationResult{Type: "float64", Function: "Log", Passed: false, Message: fmt.Sprintf("Log(%v): got %v, want %v", x, res, expected)})
			} else {
				results = append(results, ValidationResult{Type: "float64", Function: "Log", Passed: true, Message: "OK"})
			}
		}
	}

	// Sin/Cos/Tan
	for _, x := range testCases {
		res := trig.Sin(x)
		expected := math.Sin(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Sin", Passed: false, Message: fmt.Sprintf("Sin(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Sin", Passed: true, Message: "OK"})
		}

		res = trig.Cos(x)
		expected = math.Cos(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Cos", Passed: false, Message: fmt.Sprintf("Cos(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Cos", Passed: true, Message: "OK"})
		}

		res = trig.Tan(x)
		expected = math.Tan(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Tan", Passed: false, Message: fmt.Sprintf("Tan(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Tan", Passed: true, Message: "OK"})
		}
	}

	// Hyperbolic
	for _, x := range testCases {
		res := hyper.Sinh(x)
		expected := math.Sinh(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Sinh", Passed: false, Message: fmt.Sprintf("Sinh(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Sinh", Passed: true, Message: "OK"})
		}

		res = hyper.Cosh(x)
		expected = math.Cosh(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Cosh", Passed: false, Message: fmt.Sprintf("Cosh(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Cosh", Passed: true, Message: "OK"})
		}

		res = hyper.Tanh(x)
		expected = math.Tanh(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Tanh", Passed: false, Message: fmt.Sprintf("Tanh(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Tanh", Passed: true, Message: "OK"})
		}
	}

	// Inverse hyperbolic
	for _, x := range testCases {
		res := hyper.Asinh(x)
		expected := math.Asinh(x)
		if !withinTol(res, expected, 1e-10) {
			results = append(results, ValidationResult{Type: "float64", Function: "Asinh", Passed: false, Message: fmt.Sprintf("Asinh(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Asinh", Passed: true, Message: "OK"})
		}
	}

	// Sqrt
	for _, x := range testCases {
		if x >= 0 {
			res := arithmetic.Sqrt(x)
			expected := math.Sqrt(x)
			if !withinTol(res, expected, 1e-10) {
				results = append(results, ValidationResult{Type: "float64", Function: "Sqrt", Passed: false, Message: fmt.Sprintf("Sqrt(%v): got %v, want %v", x, res, expected)})
			} else {
				results = append(results, ValidationResult{Type: "float64", Function: "Sqrt", Passed: true, Message: "OK"})
			}
		}
	}

	// Pow
	for i := 0; i < 20; i++ {
		x := float64(i + 1)
		y := float64(i) * 0.5
		res := arithmetic.Pow(x, y)
		expected := math.Pow(x, y)
		if !withinTol(res, expected, 1e-9) {
			results = append(results, ValidationResult{Type: "float64", Function: "Pow", Passed: false, Message: fmt.Sprintf("Pow(%v,%v): got %v, want %v", x, y, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "float64", Function: "Pow", Passed: true, Message: "OK"})
		}
	}

	return results
}

func testComplex64() []ValidationResult {
	results := []ValidationResult{}

	testCases := []complex64{
		0, 1, -1, 1 + 1i, 1 - 1i, complex(math.Pi, math.E),
		complex(math.MaxFloat32, math.SmallestNonzeroFloat32),
	}

	// Complex Sin
	for _, x := range testCases {
		res := complex64(trigComplexSin(float64(real(x)), float64(imag(x))))
		expected := cmplx.Sin(complex128(x))
		if !withinTolComplex64(res, expected) {
			results = append(results, ValidationResult{Type: "complex64", Function: "Sin", Passed: false, Message: fmt.Sprintf("Sin(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "complex64", Function: "Sin", Passed: true, Message: "OK"})
		}
	}

	// Complex Cos
	for _, x := range testCases {
		res := complex64(trigComplexCos(float64(real(x)), float64(imag(x))))
		expected := cmplx.Cos(complex128(x))
		if !withinTolComplex64(res, expected) {
			results = append(results, ValidationResult{Type: "complex64", Function: "Cos", Passed: false, Message: fmt.Sprintf("Cos(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "complex64", Function: "Cos", Passed: true, Message: "OK"})
		}
	}

	// Complex Exp
	for _, x := range testCases {
		res := complex64(complexExp(float64(real(x)), float64(imag(x))))
		expected := cmplx.Exp(complex128(x))
		if !withinTolComplex64(res, expected) {
			results = append(results, ValidationResult{Type: "complex64", Function: "Exp", Passed: false, Message: fmt.Sprintf("Exp(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "complex64", Function: "Exp", Passed: true, Message: "OK"})
		}
	}

	return results
}

func testComplex128() []ValidationResult {
	results := []ValidationResult{}

	testCases := []complex128{
		0, 1, -1, 1 + 1i, 1 - 1i, complex(math.Pi, math.E),
		complex(math.MaxFloat64, math.SmallestNonzeroFloat64),
		complex(math.Inf(1), math.Inf(-1)),
	}

	// Complex Sin
	for _, x := range testCases {
		res := complex128(trigComplexSin(real(x), imag(x)))
		expected := cmplx.Sin(x)
		if !withinTolComplex128(res, expected) {
			results = append(results, ValidationResult{Type: "complex128", Function: "Sin", Passed: false, Message: fmt.Sprintf("Sin(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "complex128", Function: "Sin", Passed: true, Message: "OK"})
		}
	}

	// Complex Cos
	for _, x := range testCases {
		res := complex128(trigComplexCos(real(x), imag(x)))
		expected := cmplx.Cos(x)
		if !withinTolComplex128(res, expected) {
			results = append(results, ValidationResult{Type: "complex128", Function: "Cos", Passed: false, Message: fmt.Sprintf("Cos(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "complex128", Function: "Cos", Passed: true, Message: "OK"})
		}
	}

	// Complex Exp
	for _, x := range testCases {
		res := complex128(complexExp(real(x), imag(x)))
		expected := cmplx.Exp(x)
		if !withinTolComplex128(res, expected) {
			results = append(results, ValidationResult{Type: "complex128", Function: "Exp", Passed: false, Message: fmt.Sprintf("Exp(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "complex128", Function: "Exp", Passed: true, Message: "OK"})
		}
	}

	// Complex Log
	for _, x := range testCases {
		if x != 0 {
			res := complex128(complexLog(real(x), imag(x)))
			expected := cmplx.Log(x)
			if !withinTolComplex128(res, expected) {
				results = append(results, ValidationResult{Type: "complex128", Function: "Log", Passed: false, Message: fmt.Sprintf("Log(%v): got %v, want %v", x, res, expected)})
			} else {
				results = append(results, ValidationResult{Type: "complex128", Function: "Log", Passed: true, Message: "OK"})
			}
		}
	}

	// Complex Sqrt
	for _, x := range testCases {
		res := complex128(complexSqrt(real(x), imag(x)))
		expected := cmplx.Sqrt(x)
		if !withinTolComplex128(res, expected) {
			results = append(results, ValidationResult{Type: "complex128", Function: "Sqrt", Passed: false, Message: fmt.Sprintf("Sqrt(%v): got %v, want %v", x, res, expected)})
		} else {
			results = append(results, ValidationResult{Type: "complex128", Function: "Sqrt", Passed: true, Message: "OK"})
		}
	}

	return results
}

// Helper functions for complex number operations using emlgo

func trigComplexSin(r, i float64) complex128 {
	// Handle infinity cases
	if math.IsInf(r, 0) && math.IsInf(i, 0) {
		return complex(math.NaN(), math.Inf(-1))
	}
	if math.IsInf(r, 0) {
		sinR := math.Sin(r)
		coshI := hyper.Cosh(i)
		return complex(sinR*coshI, math.Cos(r)*hyper.Sinh(i))
	}
	if math.IsInf(i, 0) {
		return complex(0, math.Cos(r)*hyper.Sinh(i))
	}
	// sin(z) = sin(x)cosh(y) + i*cos(x)sinh(y)
	sinX := trig.Sin(r)
	cosX := trig.Cos(r)
	sinhY := hyper.Sinh(i)
	coshY := hyper.Cosh(i)
	return complex(sinX*coshY, cosX*sinhY)
}

func trigComplexCos(r, i float64) complex128 {
	// Handle infinity cases
	if math.IsInf(r, 0) && math.IsInf(i, 0) {
		return complex(math.Inf(1), math.NaN())
	}
	if math.IsInf(r, 0) {
		cosR := math.Cos(r)
		sinhI := hyper.Sinh(i)
		return complex(cosR*hyper.Cosh(i), -math.Sin(r)*sinhI)
	}
	if math.IsInf(i, 0) {
		return complex(math.Cos(r)*hyper.Cosh(i), 0)
	}
	// cos(z) = cos(x)cosh(y) - i*sin(x)sinh(y)
	sinX := trig.Sin(r)
	cosX := trig.Cos(r)
	sinhY := hyper.Sinh(i)
	coshY := hyper.Cosh(i)
	return complex(cosX*coshY, -sinX*sinhY)
}

func complexExp(r, i float64) complex128 {
	// Handle infinity cases
	if math.IsInf(r, 0) && math.IsInf(i, 0) {
		return complex(math.Inf(1), math.NaN())
	}
	if math.IsInf(r, 1) && i == 0 {
		return complex(math.Inf(1), math.NaN())
	}
	// exp(z) = exp(x) * (cos(y) + i*sin(y))
	expR := logexp.Exp(r)
	if math.IsInf(expR, 1) {
		sinI := math.Sin(i)
		cosI := math.Cos(i)
		return complex(math.Inf(1)*cosI, math.Inf(1)*sinI)
	}
	sinI := trig.Sin(i)
	cosI := trig.Cos(i)
	return complex(expR*cosI, expR*sinI)
}

func complexLog(r, i float64) complex128 {
	// Handle infinity cases
	if math.IsInf(r, 0) && math.IsInf(i, 0) {
		return complex(math.Inf(1), math.Atan2(-math.Inf(1), math.Inf(1)))
	}
	if math.IsInf(r, 1) && !math.IsInf(i, 0) {
		mag := arithmetic.Sqrt(r*r + i*i)
		if math.IsInf(mag, 1) {
			arg := trig.Atan2(i, r)
			return complex(logexp.Log(2*math.Abs(r)), arg)
		}
	}
	if math.IsInf(r, 0) && !math.IsInf(i, 0) {
		mag := arithmetic.Sqrt(r*r + i*i)
		if math.IsInf(mag, 1) {
			arg := trig.Atan2(i, r)
			return complex(logexp.Log(2*math.Abs(r)), arg)
		}
	}
	// log(z) = log(|z|) + i*arg(z)
	absR := math.Abs(r)
	absI := math.Abs(i)

	// Handle extremely large real part with small imaginary - avoid overflow
	// log(|z|) = 0.5 * log(r^2 + i^2) = 0.5 * log(r^2 * (1 + (i/r)^2))
	// = log(|r|) + 0.5 * log(1 + (i/r)^2)
	// For small i/r, this ≈ log(|r|) + 0.5 * (i/r)^2
	if absR > 1e150 && absI < 1e-100 {
		arg := trig.Atan2(i, r)
		// Use log10 to avoid overflow
		log10absR := math.Log10(absR)
		logVal := math.Ln10 * log10absR
		return complex(logVal, arg)
	}

	magnitude := arithmetic.Sqrt(r*r + i*i)
	if math.IsInf(magnitude, 1) {
		// For very large magnitude, use approximation
		if absR > 1e150 {
			arg := trig.Atan2(i, r)
			// Use log10 to avoid overflow
			log10absR := math.Log10(absR)
			logVal := math.Ln10 * log10absR
			return complex(logVal, arg)
		}
	}
	arg := trig.Atan2(i, r)
	logMag := logexp.Log(magnitude)
	if math.IsInf(logMag, 1) {
		absR := math.Abs(r)
		if absR > 1e150 {
			log10absR := math.Log10(absR)
			return complex(math.Ln10*log10absR, arg)
		}
	}
return complex(logMag, arg)
}

func complexSqrt(r, i float64) complex128 {
	// Handle infinity cases
	if math.IsInf(r, 0) && math.IsInf(i, 0) {
		if i < 0 {
			return complex(math.Inf(1), math.Inf(-1))
		}
		return complex(math.Inf(1), math.Inf(1))
	}
	if math.IsInf(r, 1) {
		if i == 0 {
			return complex(math.Inf(1), 0)
		}
		// For large real part with finite imaginary
		arg := trig.Atan2(i, r) / 2
		mag := math.Inf(1)
		return complex(mag*math.Cos(arg), mag*math.Sin(arg))
	}
	if math.IsInf(r, -1) {
		arg := trig.Atan2(i, r) / 2
		mag := math.Inf(1)
		return complex(mag*math.Cos(arg), mag*math.Sin(arg))
	}
	// sqrt(z) = sqrt((|z|+r)/2) + i*sign(y)*sqrt((|z|-r)/2)
	absR := math.Abs(r)
	absI := math.Abs(i)

	// Handle extremely large real part with small imaginary - avoid overflow
	if absR > 1e150 && absI < 1e-100 {
		// sqrt(x + iy) ≈ sqrt(x) + iy/(2*sqrt(x))
		sqrtR := arithmetic.Sqrt(absR)
		return complex(sqrtR, i/(2*sqrtR))
	}

	magnitude := arithmetic.Sqrt(r*r + i*i)
	if math.IsInf(magnitude, 1) {
		// For very large values
		absR := math.Abs(r)
		if absR > 1e150 {
			// sqrt(x + iy) ≈ sqrt(x) + iy/(2*sqrt(x))
			sqrtR := arithmetic.Sqrt(absR)
			return complex(sqrtR, i/(2*sqrtR))
		}
	}
	rPlus := (magnitude + r) / 2
	rMinus := (magnitude - r) / 2

	if rMinus < 0 {
		rMinus = 0
	}

	var signI float64
	if i >= 0 {
		signI = 1
	} else {
		signI = -1
	}

	return complex(arithmetic.Sqrt(rPlus), signI*arithmetic.Sqrt(rMinus))
}

// Tolerance functions

func withinTol(a, b, tol float64) bool {
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

func withinTolFloat32(a, b float32) bool {
	a64 := float64(a)
	b64 := float64(b)
	if math.IsNaN(a64) && math.IsNaN(b64) {
		return true
	}
	if math.IsInf(a64, 1) && math.IsInf(b64, 1) {
		return true
	}
	if math.IsInf(a64, -1) && math.IsInf(b64, -1) {
		return true
	}
	diff := math.Abs(a64 - b64)
	sumAbs := math.Abs(a64) + math.Abs(b64) + 1e-10
	return diff/sumAbs < 1e-6
}

func withinTolComplex64(a complex64, b complex128) bool {
	a128 := complex128(a)
	return withinTolFloat32(float32(real(a128)), float32(real(b))) &&
		withinTolFloat32(float32(imag(a128)), float32(imag(b)))
}

func withinTolComplex128(a, b complex128) bool {
	// Handle NaN matching
	if math.IsNaN(real(a)) && math.IsNaN(real(b)) && math.IsNaN(imag(a)) && math.IsNaN(imag(b)) {
		return true
	}
	// Handle infinity matching for real part
	if math.IsInf(real(a), 1) && math.IsInf(real(b), 1) {
		// Real parts both +Inf, check imaginary
		if math.IsNaN(imag(a)) && math.IsNaN(imag(b)) {
			return true
		}
		return withinTol(imag(a), imag(b), 1e-10)
	}
	if math.IsInf(real(a), -1) && math.IsInf(real(b), -1) {
		if math.IsNaN(imag(a)) && math.IsNaN(imag(b)) {
			return true
		}
		return withinTol(imag(a), imag(b), 1e-10)
	}
	// Handle NaN in imaginary
	if math.IsNaN(imag(a)) && math.IsNaN(imag(b)) {
		return withinTol(real(a), real(b), 1e-10)
	}
	// Handle NaN in real
	if math.IsNaN(real(a)) && math.IsNaN(real(b)) {
		return withinTol(imag(a), imag(b), 1e-10)
	}
	return withinTol(real(a), real(b), 1e-10) &&
		withinTol(imag(a), imag(b), 1e-10)
}

func printSummary() {
	passed := 0
	failed := 0

	for _, r := range allResults {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}

	fmt.Println()
	fmt.Println("=== Summary ===")
	fmt.Printf("Total: %d\n", len(allResults))
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)

	if failed > 0 {
		fmt.Println("\nFailed tests:")
		for _, r := range allResults {
			if !r.Passed {
				fmt.Printf("  - %s.%s: %s\n", r.Type, r.Function, r.Message)
			}
		}
		os.Exit(1)
	}

	fmt.Println("\n✓ All validation tests passed!")
}