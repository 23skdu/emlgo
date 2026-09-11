//go:build amd64
// +build amd64

package jit

import (
	"fmt"
	"math"
	"testing"
)

func TestJITSimple(t *testing.T) {
	tests := []struct {
		expr string
		x    float64
		want float64
	}{
		{"x", 5, 5},
		{"x + 1", 5, 6},
		{"x - 1", 5, 4},
		{"2 * x", 5, 10},
		{"x / 2", 10, 5},
		{"-x", 5, -5},
		{"x^2", 4, 16},
		{"x^3", 3, 27},
		{"(x + 1)^2", 2, 9},
		{"0", 42, 0},
	}
	for _, tc := range tests {
		c := NewCompiler()
		f, err := c.Compile(tc.expr)
		if err != nil {
			t.Fatalf("compile %q: %v", tc.expr, err)
		}
		got := f(tc.x)
		if math.Abs(got-tc.want) > 1e-14 {
			t.Errorf("%s at x=%v: f = %v, want %v", tc.expr, tc.x, got, tc.want)
		}
	}
}

func TestJITPolynomial(t *testing.T) {
	expr := "x^2 + 2*x + 1"
	c := NewCompiler()
	f, err := c.Compile(expr)
	if err != nil {
		t.Fatal(err)
	}
	for x := -5.0; x <= 5; x++ {
		got := f(x)
		want := x*x + 2*x + 1
		if math.Abs(got-want) > 1e-14 {
			t.Errorf("f(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestJITLargePolynomial(t *testing.T) {
	expr := "x^5 - 3*x^4 + 2*x^3 - x^2 + 5*x - 7"
	c := NewCompiler()
	f, err := c.Compile(expr)
	if err != nil {
		t.Fatal(err)
	}
	for x := -3.0; x <= 3; x += 0.5 {
		got := f(x)
		want := math.Pow(x, 5) - 3*math.Pow(x, 4) + 2*math.Pow(x, 3) - x*x + 5*x - 7
		if math.Abs(got-want) > 1e-13 {
			t.Errorf("f(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestJITFunctionCalls(t *testing.T) {
	tests := []struct {
		expr string
		x    float64
		want float64
	}{
		{"sin(x)", 0, 0},
		{"sin(x)", 1.5707963267948966, 1},
		{"cos(x)", 0, 1},
		{"cos(x)", 3.141592653589793, -1},
		{"exp(x)", 0, 1},
		{"exp(x)", 1, math.E},
		{"log(x)", 1, 0},
		{"log(x)", math.E, 1},
		{"sqrt(x)", 4, 2},
		{"sqrt(x)", 9, 3},
		{"tan(x)", 0, 0},
		{"abs(x)", -5, 5},
		{"abs(x)", 5, 5},
		{"floor(x)", 3.7, 3},
		{"ceil(x)", 3.2, 4},
		{"trunc(x)", -3.7, -3},
		{"log2(x)", 8, 3},
		{"log10(x)", 100, 2},
		{"cbrt(x)", 27, 3},
		{"asin(sin(x))", 0.5, 0.5},
		{"atan(x)", 1, 0.7853981633974483},
		{"sqrt(x^2 + 1)", 3, math.Sqrt(10)},
		{"sin(x)^2 + cos(x)^2", 1.5, 1},
	}
	for _, tc := range tests {
		c := NewCompiler()
		f, err := c.Compile(tc.expr)
		if err != nil {
			t.Fatalf("compile %q: %v", tc.expr, err)
		}
		got := f(tc.x)
		if math.Abs(got-tc.want) > 1e-10 {
			t.Errorf("%s at x=%v: f = %v, want %v (diff=%v)", tc.expr, tc.x, got, tc.want, math.Abs(got-tc.want))
		}
	}
}

func TestJITNegativeExponent(t *testing.T) {
	tests := []struct {
		expr string
		x    float64
		want float64
	}{
		{"x^-1", 2, 0.5},
		{"x^-1", 4, 0.25},
		{"x^-2", 3, 1.0 / 9.0},
		{"x^-3", 2, 0.125},
	}
	for _, tc := range tests {
		c := NewCompiler()
		f, err := c.Compile(tc.expr)
		if err != nil {
			t.Fatalf("compile %q: %v", tc.expr, err)
		}
		got := f(tc.x)
		if math.Abs(got-tc.want) > 1e-12 {
			t.Errorf("%s at x=%v: f = %v, want %v", tc.expr, tc.x, got, tc.want)
		}
	}
}

func TestJITNonIntegerExponent(t *testing.T) {
	tests := []struct {
		expr string
		x    float64
		want float64
	}{
		{"x^0.5", 4, 2},
		{"x^0.5", 9, 3},
		{"x^0.25", 16, 2},
		{"x^1.5", 4, 8},
		{"x^-0.5", 4, 0.5},
		{"x^0.5", 0.25, 0.5},
		{"x^0.1", 1, 1},
		{"x^3.7", 2, math.Pow(2, 3.7)},
	}
	for _, tc := range tests {
		c := NewCompiler()
		f, err := c.Compile(tc.expr)
		if err != nil {
			t.Fatalf("compile %q: %v", tc.expr, err)
		}
		got := f(tc.x)
		if math.Abs(got-tc.want) > 1e-10 {
			t.Errorf("%s at x=%v: f = %v, want %v (diff=%v)", tc.expr, tc.x, got, tc.want, math.Abs(got-tc.want))
		}
	}
}

func TestJITVariableExponent(t *testing.T) {
	tests := []struct {
		expr string
		x    float64
		want float64
	}{
		{"x^x", 2, 4},
		{"x^x", 3, 27},
		{"x^(x-1)", 3, 9},
		{"x^(2*x)", 2, 16},
		{"x^(0+1)", 5, 5},
	}
	for _, tc := range tests {
		c := NewCompiler()
		f, err := c.Compile(tc.expr)
		if err != nil {
			t.Fatalf("compile %q: %v", tc.expr, err)
		}
		got := f(tc.x)
		ast, _ := Parse(tc.expr)
		want := Eval(ast, tc.x)
		if math.Abs(got-want) > 1e-10 {
			t.Errorf("%s at x=%v: f = %v, want %v (diff=%v)", tc.expr, tc.x, got, want, math.Abs(got-want))
		}
	}
}

func TestJITAccuracyVsEval(t *testing.T) {
	c := NewCompiler()
	expr := "(x^4 + 2*x^3 - 3*x^2 + 4*x - 5) / (x^2 + 1)"
	f, err := c.Compile(expr)
	if err != nil {
		t.Fatal(err)
	}
	ast, _ := Parse(expr)
	for x := -5.0; x <= 5; x += 0.25 {
		got := f(x)
		want := Eval(ast, x)
		diff := math.Abs(got - want)
		if diff > 1e-13 {
			t.Errorf("f(%v) = %v, eval = %v, diff = %v", x, got, want, diff)
		}
	}
}

func TestJITRegisterAllocationDepth(t *testing.T) {
	// A nested expression demanding multiple scratch registers
	expr := "(((x + 1) * (x + 2)) + ((x + 3) * (x + 4))) * (((x + 5) * (x + 6)) + ((x + 7) * (x + 8)))"
	c := NewCompiler()
	f, err := c.Compile(expr)
	if err != nil {
		t.Fatalf("Failed to compile deep expression: %v", err)
	}
	x := 1.5
	got := f(x)
	want := (((x + 1) * (x + 2)) + ((x + 3) * (x + 4))) * (((x + 5) * (x + 6)) + ((x + 7) * (x + 8)))
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("Deep expression got %v, want %v", got, want)
	}
}

func TestJITBinaryExponentiationBySquaring(t *testing.T) {
	exponents := []int{0, 1, 2, 3, 4, 5, 8, 9, 15, 16, 31, 32}
	c := NewCompiler()
	for _, n := range exponents {
		expr := fmt.Sprintf("x^%d", n)
		f, err := c.Compile(expr)
		if err != nil {
			t.Fatalf("Failed to compile %s: %v", expr, err)
		}
		for _, x := range []float64{0.0, 1.0, 2.0, -1.5, 3.14} {
			got := f(x)
			want := math.Pow(x, float64(n))
			if math.Abs(got-want) > 1e-12 {
				t.Errorf("x=%v, exponent=%v: got %v, want %v", x, n, got, want)
			}
		}
	}
}

func TestJITRegisterSpillError(t *testing.T) {
	// Build an extremely deeply nested tree of binary operations to run out of registers.
	// We have 14 scratch registers, so a depth > 15 should spill.
	expr := "x"
	for i := 0; i < 20; i++ {
		expr = fmt.Sprintf("(%s) * (%s)", expr, expr)
	}
	c := NewCompiler()
	_, err := c.Compile(expr)
	if err == nil {
		t.Fatal("Expected compiler error due to register exhaustion, but got nil")
	}
}

func TestJITHyperbolicAndSpecial(t *testing.T) {
	funcs := []struct {
		name string
		eval func(float64) float64
	}{
		{"sinh", math.Sinh},
		{"cosh", math.Cosh},
		{"tanh", math.Tanh},
		{"asinh", math.Asinh},
		{"erf", math.Erf},
	}
	c := NewCompiler()
	for _, fn := range funcs {
		expr := fmt.Sprintf("%s(x)", fn.name)
		f, err := c.Compile(expr)
		if err != nil {
			t.Fatalf("Failed to compile %s: %v", expr, err)
		}
		for _, x := range []float64{0.0, 0.5, 1.0, -1.0} {
			got := f(x)
			want := fn.eval(x)
			if math.Abs(got-want) > 1e-12 {
				t.Errorf("%s(%v): got %v, want %v", fn.name, x, got, want)
			}
		}
	}
}

