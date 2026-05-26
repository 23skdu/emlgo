package jit

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func FuzzParse(f *testing.F) {
	seeds := []string{
		"x", "42", "3.14",
		"x+1", "2*x", "x/2", "x-1",
		"x^2", "x^3+2*x+1",
		"-x", "-(x+1)",
		"sin(x)", "cos(x)", "exp(x)", "log(x)", "sqrt(x)",
		"tan(x)", "asin(x)", "acos(x)", "atan(x)", "abs(x)",
		"x^2+2*x+1", "(x+1)*(x-1)",
		"x^5-3*x^4+2*x^3-x^2+5*x-7",
		"(((x+2)*x+3)*x+4)*x+5",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if strings.ContainsAny(input, "\x00\n\r\t") {
			t.Skip()
		}

		node, err := Parse(input)
		if err != nil {
			return
		}

		if node == nil {
			t.Errorf("Parse returned nil node with nil error for: %q", input)
			return
		}

		s := node.String()
		if s == "" {
			t.Errorf("String() returned empty for: %q", input)
		}

		fs := FormatExpr(node)
		if fs == "" {
			t.Errorf("FormatExpr returned empty for: %q", input)
		}

		fa := FormatAST(node)
		if fa == "" {
			t.Errorf("FormatAST returned empty for: %q", input)
		}

		for _, x := range []float64{-100, -1, 0, 1, 2, 100} {
			result := Eval(node, x)
			if math.IsNaN(result) || math.IsInf(result, 0) {
				continue
			}
			if _, err := strconv.ParseFloat(strconv.FormatFloat(result, 'g', -1, 64), 64); err != nil {
				t.Errorf("Eval(%q, %v) = %v (unrepresentable)", input, x, result)
			}
		}

		reparsed, err := Parse(fs)
		if err != nil {
			t.Skip()
		}
		for _, x := range []float64{-5, 0, 1, 5} {
			orig := Eval(node, x)
			round := Eval(reparsed, x)
			if math.IsNaN(orig) && math.IsNaN(round) {
				continue
			}
			if math.IsInf(orig, 0) && math.IsInf(round, 0) {
				continue
			}
			if math.Abs(orig-round) > 1e-12 {
				t.Errorf("round-trip mismatch for %q (FormatExpr=%q) at x=%v: %v vs %v",
					input, fs, x, orig, round)
			}
		}
	})
}

func FuzzEval(f *testing.F) {
	exprs := []string{
		"x", "42", "x+1", "2*x", "x^2", "x^3+2*x+1",
		"sin(x)", "cos(x)", "exp(x)", "log(x+2)", "sqrt(x+2)",
		"-x", "x/(x+1)", "(x+1)*(x-1)",
	}
	for _, e := range exprs {
		f.Add(e)
	}

	f.Fuzz(func(t *testing.T, expr string) {
		if strings.ContainsAny(expr, "\x00\n\r\t") {
			t.Skip()
		}

		node, err := Parse(expr)
		if err != nil {
			t.Skip()
		}

		for _, x := range []float64{-10, -1, -0.5, 0, 0.5, 1, 10} {
			result := Eval(node, x)
			if math.IsNaN(result) || math.IsInf(result, 0) {
				continue
			}
			if result < -1e308 || result > 1e308 {
				continue
			}
		}
	})
}
