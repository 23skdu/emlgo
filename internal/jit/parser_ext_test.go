package jit

import (
	"errors"
	"math"
	"testing"
)

func TestRoundFunction(t *testing.T) {
	cases := []struct {
		in, want float64
	}{
		{2.4, 2}, {2.5, 3}, {-2.5, -3}, {0.9, 1}, {-0.4, 0},
	}
	for _, tc := range cases {
		n, err := Parse("round(x)")
		if err != nil {
			t.Fatalf("Parse round(x): %v", err)
		}
		got := Eval(n, tc.in)
		if got != tc.want {
			t.Errorf("round(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseWithVars(t *testing.T) {
	// Parse a multi-variable expression; y and z are treated as variables.
	n, err := ParseWithVars("x + y", []string{"x", "y"})
	if err != nil {
		t.Fatalf("ParseWithVars: %v", err)
	}
	got := EvalVars(n, map[string]float64{"x": 3, "y": 4})
	if got != 7 {
		t.Errorf("x+y with x=3,y=4 = %v, want 7", got)
	}
}

func TestEvalVars(t *testing.T) {
	// x * y
	n, err := ParseWithVars("x * y", []string{"x", "y"})
	if err != nil {
		t.Fatalf("ParseWithVars: %v", err)
	}
	cases := [][3]float64{{2, 3, 6}, {-1, 5, -5}, {0, 100, 0}, {1.5, 2, 3}}
	for _, tc := range cases {
		got := EvalVars(n, map[string]float64{"x": tc[0], "y": tc[1]})
		if math.Abs(got-tc[2]) > 1e-12 {
			t.Errorf("x*y: x=%v, y=%v → %v, want %v", tc[0], tc[1], got, tc[2])
		}
	}
}

func TestEvalVarsSinZ(t *testing.T) {
	n, err := ParseWithVars("sin(z)", []string{"z"})
	if err != nil {
		t.Fatalf("ParseWithVars sin(z): %v", err)
	}
	for _, z := range []float64{0, math.Pi / 2, math.Pi} {
		got := EvalVars(n, map[string]float64{"z": z})
		want := math.Sin(z)
		if math.Abs(got-want) > 1e-12 {
			t.Errorf("sin(z) z=%v: got %v, want %v", z, got, want)
		}
	}
}

func TestParseErrorStructured(t *testing.T) {
	_, err := Parse("x + * 3")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if pe.Msg == "" {
		t.Error("ParseError.Msg should not be empty")
	}
}

func TestParseErrorEOF(t *testing.T) {
	_, err := Parse("x +")
	if err == nil {
		t.Fatal("expected parse error for trailing operator")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
}

func TestParseErrorUnexpectedToken(t *testing.T) {
	_, err := Parse("x y")
	if err == nil {
		t.Fatal("expected parse error for 'x y'")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
}

func TestVariableNamePreserved(t *testing.T) {
	n, err := Parse("myvar")
	if err != nil {
		t.Fatalf("Parse myvar: %v", err)
	}
	v, ok := n.(Variable)
	if !ok {
		t.Fatalf("expected Variable, got %T", n)
	}
	if v.Name != "myvar" {
		t.Errorf("Variable.Name = %q, want %q", v.Name, "myvar")
	}
}

// FuzzParseNoPanic ensures the parser never panics on arbitrary input.
func FuzzParseNoPanic(f *testing.F) {
	seeds := []string{
		"x", "sin(x)", "x + 1", "x * y", "sqrt(x^2)", "",
		"((((", "1/0", "x^^2", "round(x)", "log(x + sin(y))",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parse(%q) panicked: %v", input, r)
			}
		}()
		_, _ = Parse(input)
	})
}
