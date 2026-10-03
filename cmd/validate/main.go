package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"reflect"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/fastmath"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/quant"
	"github.com/emlgo/eml/pkg/trig"
)

var (
	verbose    bool
	failedOnly bool
	typeFilter string
	exitFunc   = os.Exit
)

func init() {
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&failedOnly, "f", false, "Show only failed tests")
	flag.StringVar(&typeFilter, "type", "", "Filter by type (int, uint, float, complex, fastmath, quant, batch)")
}

type ValidationResult struct {
	Type     string
	Function string
	Passed   bool
	Message  string
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
	validateFastMath()
	validateQuant()
	validateBatch()

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
	if !withinTol(res, 5, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: false, Message: fmt.Sprintf("got %v, want 5", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: true, Message: "OK"})
	}

	// Floor/Ceil
	res = arithmetic.Floor(3.7)
	if !withinTol(res, 3, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Floor", Passed: false, Message: fmt.Sprintf("got %v, want 3", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Floor", Passed: true, Message: "OK"})
	}

	res = arithmetic.Ceil(3.2)
	if !withinTol(res, 4, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Ceil", Passed: false, Message: fmt.Sprintf("got %v, want 4", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Ceil", Passed: true, Message: "OK"})
	}

	// Round
	res = arithmetic.Round(3.5)
	if !withinTol(res, 4, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Round", Passed: false, Message: fmt.Sprintf("got %v, want 4", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Round", Passed: true, Message: "OK"})
	}

	// Trunc
	res = arithmetic.Trunc(3.7)
	if !withinTol(res, 3, 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Trunc", Passed: false, Message: fmt.Sprintf("got %v, want 3", res)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Trunc", Passed: true, Message: "OK"})
	}

	// Max
	res = arithmetic.Max(float64(a), float64(b))
	if !withinTol(res, float64(a), 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, a)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: true, Message: "OK"})
	}

	// Min
	res = arithmetic.Min(float64(a), float64(b))
	if !withinTol(res, float64(b), 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, b)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: true, Message: "OK"})
	}

	// Neg
	res = arithmetic.Neg(float64(a))
	negExpected := -float64(a)
	if !withinTol(res, negExpected, 0.001) {
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
	if !withinTol(res, squareVal, 0.001) {
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
	if !withinTol(res, float64(a), 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, a)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Abs", Passed: true, Message: "OK"})
	}

	// Max
	res = arithmetic.Max(float64(a), float64(b))
	if !withinTol(res, float64(a), 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, a)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Max", Passed: true, Message: "OK"})
	}

	// Min
	res = arithmetic.Min(float64(a), float64(b))
	if !withinTol(res, float64(b), 0.001) {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: false, Message: fmt.Sprintf("got %v, want %v", res, b)})
	} else {
		results = append(results, ValidationResult{Type: typeName, Function: "Min", Passed: true, Message: "OK"})
	}

	// Square
	res = arithmetic.Square(float64(a))
	squareVal := float64(a) * float64(a)
	if !withinTol(res, squareVal, 0.001) {
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
			expected := eml.ExactLogOracle(x)
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

// logMagnitude returns log(sqrt(r*r + i*i)) for the complex number r + i*i.
//
// It never overflows: instead of forming r*r + i*i, it factors out the larger
// component and only ever squares a ratio whose magnitude is at most 1,
// log|z| = log(max) + 0.5*log1p((min/max)^2).
// It returns -Inf when both components are zero, matching math/cmplx.Log.
func logMagnitude(r, i float64) float64 {
	absR, absI := math.Abs(r), math.Abs(i)
	if absR == 0 && absI == 0 {
		return math.Inf(-1)
	}
	if absR >= absI {
		q := absI / absR
		return logexp.Log(absR) + 0.5*math.Log1p(q*q)
	}
	q := absR / absI
	return logexp.Log(absI) + 0.5*math.Log1p(q*q)
}

func complexLog(r, i float64) complex128 {
	// Handle infinity cases
	if math.IsInf(r, 0) && math.IsInf(i, 0) {
		return complex(math.Inf(1), math.Atan2(-math.Inf(1), math.Inf(1)))
	}
	if math.IsInf(r, 0) || math.IsInf(i, 0) {
		// One infinite component dominates the magnitude; only its magnitude
		// contributes to the real part.
		dom := math.Abs(r)
		if !math.IsInf(dom, 0) {
			dom = math.Abs(i)
		}
		return complex(logexp.Log(dom), trig.Atan2(i, r))
	}
	// log(z) = log(|z|) + i*arg(z)
	return complex(logMagnitude(r, i), trig.Atan2(i, r))
}

func complexSqrt(r, i float64) complex128 {
	// Handle infinity cases
	if math.IsInf(r, 0) && math.IsInf(i, 0) {
		if i < 0 {
			return complex(math.Inf(1), math.Inf(-1))
		}
		return complex(math.Inf(1), math.Inf(1))
	}
	if math.IsInf(r, 0) || math.IsInf(i, 0) {
		if i == 0 {
			return complex(0, math.Inf(1))
		}
		// One infinite component dominates the argument; the magnitude of the
		// result is infinite and its direction is half the argument of z.
		arg := trig.Atan2(i, r) / 2
		return complex(math.Inf(1)*math.Cos(arg), math.Inf(1)*math.Sin(arg))
	}

	// sqrt(z) = sqrt((|z|+r)/2) + i*sign(i)*sqrt((|z|-r)/2).
	//
	// Factoring out the larger component keeps every intermediate value
	// bounded: with M = max(|r|,|i|) and (a,b) = (r,i)/M we have |a+bi| <= sqrt(2),
	// so sqrt(z) = sqrt(M) * sqrt(a + bi) never forms r*r or i*i and therefore
	// cannot overflow for any finite input.
	absR, absI := math.Abs(r), math.Abs(i)
	m := max(absR, absI)
	if m == 0 {
		return complex(0, 0)
	}
	a, b := r/m, i/m
	hyp := math.Hypot(a, b)
	scale := arithmetic.Sqrt(m)
	signI := 1.0
	if b < 0 {
		signI = -1
	}
	// hyp >= |a|, so both halves are non-negative; clamp guards rounding.
	re := (hyp + a) / 2
	im := (hyp - a) / 2
	re = math.Max(re, 0)
	im = math.Max(im, 0)
	return complex(scale*arithmetic.Sqrt(re), signI*scale*arithmetic.Sqrt(im))
}

// Tolerance functions

var forceFailTol = false

func withinTol(a, b, tol float64) bool {
	if forceFailTol {
		return false
	}
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
	if forceFailTol {
		return false
	}
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
	if forceFailTol {
		return false
	}
	a128 := complex128(a)
	return withinTolFloat32(float32(real(a128)), float32(real(b))) &&
		withinTolFloat32(float32(imag(a128)), float32(imag(b)))
}

func withinTolComplex128(a, b complex128) bool {
	if forceFailTol {
		return false
	}
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

// summarize reports the pass/fail counts for results and returns true when
// every result passed. It does not terminate the process.
func summarize(results []ValidationResult) bool {
	passed := 0
	failed := 0

	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}

	fmt.Println()
	fmt.Println("=== Summary ===")
	fmt.Printf("Total: %d\n", len(results))
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)

	if failed > 0 {
		fmt.Println("\nFailed tests:")
		for _, r := range results {
			if !r.Passed {
				fmt.Printf("  - %s.%s: %s\n", r.Type, r.Function, r.Message)
			}
		}
		return false
	}

	fmt.Println("\n✓ All validation tests passed!")
	return true
}

func printSummary() {
	if !summarize(allResults) {
		exitFunc(1)
	}
}

func validateFastMath() {
	if typeFilter != "" && typeFilter != "fastmath" {
		return
	}
	fmt.Println("--- FastMath ---")
	results := testFastMath()
	allResults = append(allResults, results...)
	if filterPassed(results) {
		fmt.Println("  fastmath: PASSED")
	}
}

func testFastMath() []ValidationResult {
	var results []ValidationResult

	// Exp
	expCases := []float64{-20, -5, -1, 0, 1, 5, 20}
	for _, x := range expCases {
		res := fastmath.Exp(x)
		exp := math.Exp(x)
		if !withinTol(res, exp, 1e-4) {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Exp", Passed: false, Message: fmt.Sprintf("Exp(%v): got %v, want %v", x, res, exp)})
		} else {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Exp", Passed: true, Message: "OK"})
		}
	}

	// Log
	logCases := []float64{0.1, 0.5, 1, 2, 10, 100}
	for _, x := range logCases {
		res := fastmath.Log(x)
		exp := math.Log(x)
		if !withinTol(res, exp, 1e-4) {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Log", Passed: false, Message: fmt.Sprintf("Log(%v): got %v, want %v", x, res, exp)})
		} else {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Log", Passed: true, Message: "OK"})
		}
	}

	// Sin & Cos
	trigCases := []float64{-3, -1, 0, 1, 3}
	for _, x := range trigCases {
		resSin := fastmath.Sin(x)
		expSin := math.Sin(x)
		if !withinTol(resSin, expSin, 1e-4) {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Sin", Passed: false, Message: fmt.Sprintf("Sin(%v): got %v, want %v", x, resSin, expSin)})
		} else {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Sin", Passed: true, Message: "OK"})
		}

		resCos := fastmath.Cos(x)
		expCos := math.Cos(x)
		if !withinTol(resCos, expCos, 1e-4) {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Cos", Passed: false, Message: fmt.Sprintf("Cos(%v): got %v, want %v", x, resCos, expCos)})
		} else {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Cos", Passed: true, Message: "OK"})
		}
	}

	// Sqrt
	sqrtCases := []float64{0, 0.25, 1, 4, 16, 100}
	for _, x := range sqrtCases {
		res := fastmath.Sqrt(x)
		exp := math.Sqrt(x)
		if !withinTol(res, exp, 1e-4) {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Sqrt", Passed: false, Message: fmt.Sprintf("Sqrt(%v): got %v, want %v", x, res, exp)})
		} else {
			results = append(results, ValidationResult{Type: "fastmath", Function: "Sqrt", Passed: true, Message: "OK"})
		}
	}

	// LnRegularized
	lnRegVal := fastmath.LnRegularized(10.0, 1e-6)
	expLnReg := math.Log(10.0)
	if !withinTol(lnRegVal, expLnReg, 1e-4) {
		results = append(results, ValidationResult{Type: "fastmath", Function: "LnRegularized", Passed: false, Message: fmt.Sprintf("LnRegularized(10.0, 1e-6): got %v, want %v", lnRegVal, expLnReg)})
	} else {
		results = append(results, ValidationResult{Type: "fastmath", Function: "LnRegularized", Passed: true, Message: "OK"})
	}

	// FastExpMinimax
	minimaxVal := fastmath.FastExpMinimax(1.0)
	if !withinTol(minimaxVal, math.E, 1e-4) {
		results = append(results, ValidationResult{Type: "fastmath", Function: "FastExpMinimax", Passed: false, Message: fmt.Sprintf("FastExpMinimax(1.0): got %v, want %v", minimaxVal, math.E)})
	} else {
		results = append(results, ValidationResult{Type: "fastmath", Function: "FastExpMinimax", Passed: true, Message: "OK"})
	}

	// FastExpF32
	f32ExpVal := fastmath.FastExpF32(1.0)
	if !withinTol(float64(f32ExpVal), math.E, 1e-3) {
		results = append(results, ValidationResult{Type: "fastmath", Function: "FastExpF32", Passed: false, Message: fmt.Sprintf("FastExpF32(1.0): got %v, want %v", f32ExpVal, math.E)})
	} else {
		results = append(results, ValidationResult{Type: "fastmath", Function: "FastExpF32", Passed: true, Message: "OK"})
	}

	// FastLogF32
	f32LogVal := fastmath.FastLogF32(float32(math.E))
	if !withinTol(float64(f32LogVal), 1.0, 1e-3) {
		results = append(results, ValidationResult{Type: "fastmath", Function: "FastLogF32", Passed: false, Message: fmt.Sprintf("FastLogF32(e): got %v, want 1.0", f32LogVal)})
	} else {
		results = append(results, ValidationResult{Type: "fastmath", Function: "FastLogF32", Passed: true, Message: "OK"})
	}

	return results
}

func validateQuant() {
	if typeFilter != "" && typeFilter != "quant" {
		return
	}
	fmt.Println("--- Quantization (TurboQuant4) ---")
	results := testQuant()
	allResults = append(allResults, results...)
	if filterPassed(results) {
		fmt.Println("  quant: PASSED")
	}
}

func testQuant() []ValidationResult {
	var results []ValidationResult

	// 1. Pack4Bit / Unpack4Bit
	origNibbles := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	packed := make([]byte, (len(origNibbles)+1)/2)
	quant.Pack4Bit(origNibbles, packed)
	unpacked := make([]byte, len(origNibbles))
	quant.Unpack4Bit(packed, unpacked, len(origNibbles))

	packDiff := float64(unpacked[0] - origNibbles[0])
	if !withinTol(packDiff, 0.0, 0.1) {
		results = append(results, ValidationResult{Type: "quant", Function: "PackUnpack4Bit", Passed: false, Message: "Pack4Bit/Unpack4Bit roundtrip mismatch"})
	} else {
		results = append(results, ValidationResult{Type: "quant", Function: "PackUnpack4Bit", Passed: true, Message: "OK"})
	}

	// 2. PolarTransformBatch
	src := []float32{3.0, 4.0, 1.0, 1.0}
	dstRadii := make([]float32, 2)
	dstAngles := make([]float32, 2)
	quant.PolarTransformBatch(src, dstRadii, dstAngles)
	if !withinTol(float64(dstRadii[0]), 5.0, 0.1) {
		results = append(results, ValidationResult{Type: "quant", Function: "PolarTransformBatch", Passed: false, Message: fmt.Sprintf("PolarTransformBatch: radius got %v, want 5.0", dstRadii[0])})
	} else {
		results = append(results, ValidationResult{Type: "quant", Function: "PolarTransformBatch", Passed: true, Message: "OK"})
	}

	// 3. EncodeTurboQuant4 and TurboQuant4Distance
	for _, dim := range []int{16, 64} {
		vec := make([]float32, dim)
		for i := range vec {
			vec[i] = float32(i+1) / float32(dim)
		}
		encoded, _ := quant.EncodeTurboQuant4(vec, dim)
		dist, _ := quant.TurboQuant4Distance(vec, encoded, dim, dim)
		if !withinTol(float64(dist), 1.0, 10.0) {
			results = append(results, ValidationResult{Type: "quant", Function: fmt.Sprintf("EncodeDistanceTQ4_%d", dim), Passed: false, Message: fmt.Sprintf("Distance too large or failed: dist=%v", dist)})
		} else {
			results = append(results, ValidationResult{Type: "quant", Function: fmt.Sprintf("EncodeDistanceTQ4_%d", dim), Passed: true, Message: "OK"})
		}
	}

	return results
}

func validateBatch() {
	if typeFilter != "" && typeFilter != "batch" {
		return
	}
	fmt.Println("--- Batch / SIMD Operations ---")
	results := testBatch()
	allResults = append(allResults, results...)
	if filterPassed(results) {
		fmt.Println("  batch: PASSED")
	}
}

func testBatch() []ValidationResult {
	var results []ValidationResult
	const n = 64

	a := make([]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = float64(i)*0.1 + 0.5
		b[i] = float64(i)*0.05 + 1.0
	}

	checkUnary := func(name string, got, want []float64) {
		passed := !forceFailTol
		for i := range got {
			if !withinTol(got[i], want[i], 1e-9) {
				passed = false
				break
			}
		}
		if passed {
			results = append(results, ValidationResult{Type: "batch", Function: name, Passed: true, Message: "OK"})
		} else {
			results = append(results, ValidationResult{Type: "batch", Function: name, Passed: false, Message: "batch vs scalar mismatch"})
		}
	}

	// AddBatch
	addWant := make([]float64, n)
	for i := range a {
		addWant[i] = a[i] + b[i]
	}
	checkUnary("AddBatch", arithmetic.AddBatch(a, b), addWant)

	// SubBatch
	subWant := make([]float64, n)
	for i := range a {
		subWant[i] = a[i] - b[i]
	}
	checkUnary("SubBatch", arithmetic.SubBatch(a, b), subWant)

	// MulBatch
	mulWant := make([]float64, n)
	for i := range a {
		mulWant[i] = a[i] * b[i]
	}
	checkUnary("MulBatch", arithmetic.MulBatch(a, b), mulWant)

	// DivBatch
	divWant := make([]float64, n)
	for i := range a {
		divWant[i] = a[i] / b[i]
	}
	checkUnary("DivBatch", arithmetic.DivBatch(a, b), divWant)

	// SqrtBatch
	sqrtWant := make([]float64, n)
	for i := range a {
		sqrtWant[i] = math.Sqrt(a[i])
	}
	checkUnary("SqrtBatch", arithmetic.SqrtBatch(a), sqrtWant)

	// AbsBatch
	absIn := make([]float64, n)
	absWant := make([]float64, n)
	for i := range a {
		absIn[i] = a[i] - 3.0
		absWant[i] = math.Abs(absIn[i])
	}
	checkUnary("AbsBatch", arithmetic.AbsBatch(absIn), absWant)

	// NegBatch
	negWant := make([]float64, n)
	for i := range a {
		negWant[i] = -a[i]
	}
	checkUnary("NegBatch", arithmetic.NegBatch(a), negWant)

	// InvBatch
	invWant := make([]float64, n)
	for i := range a {
		invWant[i] = 1.0 / a[i]
	}
	checkUnary("InvBatch", arithmetic.InvBatch(a), invWant)

	// FloorBatch
	floorWant := make([]float64, n)
	for i := range a {
		floorWant[i] = math.Floor(a[i])
	}
	checkUnary("FloorBatch", arithmetic.FloorBatch(a), floorWant)

	// CeilBatch
	ceilWant := make([]float64, n)
	for i := range a {
		ceilWant[i] = math.Ceil(a[i])
	}
	checkUnary("CeilBatch", arithmetic.CeilBatch(a), ceilWant)

	// TruncBatch
	truncWant := make([]float64, n)
	for i := range a {
		truncWant[i] = math.Trunc(a[i])
	}
	checkUnary("TruncBatch", arithmetic.TruncBatch(a), truncWant)

	// Log1pBatch
	log1pWant := make([]float64, n)
	for i := range a {
		log1pWant[i] = math.Log1p(a[i])
	}
	checkUnary("Log1pBatch", arithmetic.Log1pBatch(a), log1pWant)

	// Expm1Batch
	expm1Want := make([]float64, n)
	for i := range a {
		expm1Want[i] = math.Expm1(a[i])
	}
	checkUnary("Expm1Batch", arithmetic.Expm1Batch(a), expm1Want)

	// PowBatch
	powWant := make([]float64, n)
	for i := range a {
		powWant[i] = math.Pow(a[i], 2.0)
	}
	checkUnary("PowBatch", arithmetic.PowBatch(a, 2.0), powWant)

	// CbrtBatch
	cbrtWant := make([]float64, n)
	for i := range a {
		cbrtWant[i] = math.Cbrt(a[i])
	}
	checkUnary("CbrtBatch", arithmetic.CbrtBatch(a), cbrtWant)

	// HypotBatch
	hypotWant := make([]float64, n)
	for i := range a {
		hypotWant[i] = math.Hypot(a[i], b[i])
	}
	checkUnary("HypotBatch", arithmetic.HypotBatch(a, b), hypotWant)

	// MaxBatch
	maxWant := make([]float64, n)
	for i := range a {
		maxWant[i] = math.Max(a[i], 2.0)
	}
	checkUnary("MaxBatch", arithmetic.MaxBatch(a, 2.0), maxWant)

	// MinBatch
	minWant := make([]float64, n)
	for i := range a {
		minWant[i] = math.Min(a[i], 2.0)
	}
	checkUnary("MinBatch", arithmetic.MinBatch(a, 2.0), minWant)

	// ExpBatch & LogBatch
	expBatchWant := make([]float64, n)
	logBatchWant := make([]float64, n)
	for i := range a {
		expBatchWant[i] = math.Exp(a[i])
		logBatchWant[i] = eml.ExactLogOracle(a[i])
	}
	checkUnary("ExpBatch", logexp.ExpBatch(a), expBatchWant)
	checkUnary("LogBatch", logexp.LogBatch(a), logBatchWant)

	// SinBatch, CosBatch, TanBatch
	sinBatchWant := make([]float64, n)
	cosBatchWant := make([]float64, n)
	tanBatchWant := make([]float64, n)
	for i := range a {
		sinBatchWant[i] = math.Sin(a[i])
		cosBatchWant[i] = math.Cos(a[i])
		tanBatchWant[i] = math.Tan(a[i])
	}
	checkUnary("SinBatch", trig.SinBatch(a), sinBatchWant)
	checkUnary("CosBatch", trig.CosBatch(a), cosBatchWant)
	checkUnary("TanBatch", trig.TanBatch(a), tanBatchWant)

	return results
}

