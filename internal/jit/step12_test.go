package jit

import (
	"fmt"
	"math"
	"testing"
)

func TestEMLEvalAndDiffSupportAllFunctions(t *testing.T) {
	funcs := []string{
		"sin", "cos", "tan", "asin", "acos", "atan", "abs", "cbrt",
		"log2", "log10", "sinh", "cosh", "tanh", "asinh", "acosh", "atanh",
		"erf", "gamma", "sqrt", "exp", "log",
	}

	for _, fn := range funcs {
		t.Run(fn, func(t *testing.T) {
			expr := fmt.Sprintf("%s(x)", fn)
			ast, err := Parse(expr)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", expr, err)
			}
			emlTree := Canonicalize(ast)

			x := 0.5
			if fn == "acosh" {
				x = 1.5
			}

			// EMLEval must return the correct non-zero mathematical value
			val := EMLEval(emlTree, x)
			want := Eval(ast, x)
			if math.Abs(val-want) > 1e-9 {
				t.Fatalf("EMLEval %s(%v) = %v, want %v", fn, x, val, want)
			}

			// Diff must return a non-zero derivative for functions with analytical derivatives
			if fn != "gamma" {
				diffNode := Diff(emlTree)
				dVal := EMLEval(diffNode, x)
				if math.IsNaN(dVal) || math.Abs(dVal) < 1e-9 {
					t.Fatalf("Diff(%s) evaluated to %v at x=%v", fn, dVal, x)
				}
			}
		})
	}
}
