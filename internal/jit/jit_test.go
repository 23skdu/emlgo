package jit

import (
	"math"
	"testing"
)

func TestParseNumber(t *testing.T) {
	n, err := Parse("42")
	if err != nil {
		t.Fatal(err)
	}
	num, ok := n.(Number)
	if !ok || num.Value != 42 {
		t.Fatalf("expected Number(42), got %T(%v)", n, n)
	}
}

func TestParseFloat(t *testing.T) {
	n, err := Parse("3.14")
	if err != nil {
		t.Fatal(err)
	}
	num, ok := n.(Number)
	if !ok || math.Abs(num.Value-3.14) > 1e-15 {
		t.Fatalf("expected Number(3.14), got %T(%v)", n, n)
	}
}

func TestParseVariable(t *testing.T) {
	n, err := Parse("x")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := n.(Variable); !ok {
		t.Fatalf("expected Variable, got %T", n)
	}
}

func TestParseAdd(t *testing.T) {
	n, err := Parse("x + 1")
	if err != nil {
		t.Fatal(err)
	}
	b, ok := n.(BinaryOp)
	if !ok || b.Op != '+' {
		t.Fatalf("expected BinaryOp(+), got %T(%v)", n, n)
	}
	if _, ok := b.Left.(Variable); !ok {
		t.Fatalf("expected Variable left, got %T", b.Left)
	}
	r, ok := b.Right.(Number)
	if !ok || r.Value != 1 {
		t.Fatalf("expected Number(1) right, got %T(%v)", b.Right, b.Right)
	}
}

func TestParseMul(t *testing.T) {
	n, err := Parse("2*x")
	if err != nil {
		t.Fatal(err)
	}
	b, ok := n.(BinaryOp)
	if !ok || b.Op != '*' {
		t.Fatalf("expected BinaryOp(*), got %T(%v)", n, n)
	}
}

func TestParsePower(t *testing.T) {
	n, err := Parse("x^3")
	if err != nil {
		t.Fatal(err)
	}
	b, ok := n.(BinaryOp)
	if !ok || b.Op != '^' {
		t.Fatalf("expected BinaryOp(^), got %T(%v)", n, n)
	}
	r, ok := b.Right.(Number)
	if !ok || r.Value != 3 {
		t.Fatalf("expected Number(3) right, got %T(%v)", b.Right, b.Right)
	}
}

func TestParseParens(t *testing.T) {
	n, err := Parse("(x + 1) * 2")
	if err != nil {
		t.Fatal(err)
	}
	b, ok := n.(BinaryOp)
	if !ok || b.Op != '*' {
		t.Fatalf("expected BinaryOp(*), got %T(%v)", n, n)
	}
	if _, ok := b.Left.(BinaryOp); !ok {
		t.Fatalf("expected BinaryOp left, got %T", b.Left)
	}
}

func TestParseUnaryMinus(t *testing.T) {
	n, err := Parse("-x")
	if err != nil {
		t.Fatal(err)
	}
	u, ok := n.(UnaryOp)
	if !ok || u.Op != '-' {
		t.Fatalf("expected UnaryOp(-), got %T(%v)", n, n)
	}
}

func TestParseComplex(t *testing.T) {
	expr := "x^2 + 2*x + 1"
	n, err := Parse(expr)
	if err != nil {
		t.Fatal(err)
	}
	if n.String() != "(((x ^ 2) + (2 * x)) + 1)" {
		t.Fatalf("unexpected AST: %s", n)
	}
}

func TestParseFunction(t *testing.T) {
	n, err := Parse("sin(x)")
	if err != nil {
		t.Fatal(err)
	}
	fc, ok := n.(FunctionCall)
	if !ok || fc.Name != "sin" {
		t.Fatalf("expected FunctionCall(sin), got %T(%v)", n, n)
	}
}

func TestParseNestedFunction(t *testing.T) {
	n, err := Parse("sqrt(x^2 + 1)")
	if err != nil {
		t.Fatal(err)
	}
	fc, ok := n.(FunctionCall)
	if !ok || fc.Name != "sqrt" {
		t.Fatalf("expected FunctionCall(sqrt), got %T(%v)", n, n)
	}
}

func TestParseErrors(t *testing.T) {
	cases := []string{
		"2x",
		"x + ",
		"* x",
		"x )",
		"(x + 1",
		"sin x",
		"unknown(1)",
	}
	for _, expr := range cases {
		_, err := Parse(expr)
		if err == nil {
			t.Errorf("expected error for %q", expr)
		}
	}
}

func TestEvalNumber(t *testing.T) {
	n := Number{Value: 3.14}
	if got := Eval(n, 0); got != 3.14 {
		t.Fatalf("Eval = %v, want 3.14", got)
	}
}

func TestEvalVariable(t *testing.T) {
	n := Variable{}
	if got := Eval(n, 42); got != 42 {
		t.Fatalf("Eval = %v, want 42", got)
	}
}

func TestEvalBinaryOp(t *testing.T) {
	tests := []struct {
		expr string
		x    float64
		want float64
	}{
		{"x + 1", 5, 6},
		{"x - 1", 5, 4},
		{"2 * x", 5, 10},
		{"x / 2", 10, 5},
		{"x^2", 4, 16},
		{"x^3", 3, 27},
		{"(x + 1)^2", 2, 9},
		{"x^2 + 2*x + 1", 3, 16},
	}
	for _, tc := range tests {
		n, err := Parse(tc.expr)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.expr, err)
		}
		got := Eval(n, tc.x)
		if math.Abs(got-tc.want) > 1e-14 {
			t.Errorf("%s at x=%v: Eval = %v, want %v", tc.expr, tc.x, got, tc.want)
		}
	}
}

func TestEvalUnaryMinus(t *testing.T) {
	n := UnaryOp{Op: '-', Operand: Variable{}}
	if got := Eval(n, 5); got != -5 {
		t.Fatalf("Eval(-x, 5) = %v, want -5", got)
	}
}

func TestEvalFunction(t *testing.T) {
	tests := []struct {
		expr string
		x    float64
		want float64
	}{
		{"sin(0)", 0, 0},
		{"cos(0)", 0, 1},
		{"exp(0)", 0, 1},
		{"log(1)", 0, 0},
		{"sqrt(4)", 0, 2},
		{"abs(-3)", 0, 3},
		{"tan(0)", 0, 0},
	}
	for _, tc := range tests {
		n, err := Parse(tc.expr)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.expr, err)
		}
		got := Eval(n, tc.x)
		if math.Abs(got-tc.want) > 1e-14 {
			t.Errorf("%s: Eval = %v, want %v", tc.expr, got, tc.want)
		}
	}
}
