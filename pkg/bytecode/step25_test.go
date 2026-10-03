package bytecode

import (
	"fmt"
	"math"
	"testing"

	"github.com/emlgo/eml/internal/jit"
)

func TestCompileEMLAllFunctions(t *testing.T) {
	funcs := []string{
		"sin", "cos", "tan", "asin", "acos", "atan", "abs", "cbrt",
		"log2", "log10", "ceil", "floor", "trunc", "round",
		"sinh", "cosh", "tanh", "asinh", "acosh", "atanh",
		"erf", "gamma", "sqrt", "exp", "log",
	}

	for _, fn := range funcs {
		t.Run(fn, func(t *testing.T) {
			expr := fmt.Sprintf("%s(x)", fn)
			ast, err := jit.Parse(expr)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", expr, err)
			}
			emlTree := jit.Canonicalize(ast)
			prog, err := CompileEML(emlTree)
			if err != nil {
				t.Fatalf("CompileEML for %s failed: %v", fn, err)
			}

			// Evaluate at a safe domain point
			x := 0.5
			if fn == "acosh" {
				x = 1.5
			}
			gotVM := prog.Eval([]float64{x}, nil)
			wantJIT := jit.Eval(ast, x)

			if math.IsNaN(gotVM) && math.IsNaN(wantJIT) {
				return
			}
			diff := math.Abs(gotVM - wantJIT)
			if diff > 1e-6 {
				t.Errorf("%s(%v): got VM %v, want JIT %v (diff %v)", fn, x, gotVM, wantJIT, diff)
			}
		})
	}
}

func FuzzCompileEMLRoundTrip(f *testing.F) {
	f.Add(float64(0.5), uint8(0))
	f.Add(float64(2.0), uint8(1))
	f.Add(float64(0.1), uint8(2))

	f.Fuzz(func(t *testing.T, x float64, fnChoice uint8) {
		funcs := []string{
			"sin", "cos", "tan", "asin", "acos", "atan", "abs", "cbrt",
			"sinh", "cosh", "tanh", "asinh", "sqrt", "exp", "log",
		}
		fn := funcs[int(fnChoice)%len(funcs)]
		switch fn {
		case "asin", "acos", "atanh":
			// domain (-1, 1)
			x = 0.5
		case "sqrt", "log":
			x = math.Abs(x) + 0.1
		}

		expr := fmt.Sprintf("%s(x)", fn)
		ast, err := jit.Parse(expr)
		if err != nil {
			return
		}
		emlTree := jit.Canonicalize(ast)
		prog, err := CompileEML(emlTree)
		if err != nil {
			t.Fatalf("CompileEML(%q) failed: %v", expr, err)
		}

		vmVal := prog.Eval([]float64{x}, nil)
		jitVal := jit.Eval(ast, x)

		if math.IsNaN(vmVal) && math.IsNaN(jitVal) {
			return
		}
		diff := math.Abs(vmVal - jitVal)
		denom := math.Max(math.Abs(vmVal), math.Abs(jitVal))
		if denom > 1e-9 {
			diff /= denom
		}
		if diff > 1e-4 {
			t.Fatalf("%s(%v): vm=%v, jit=%v", fn, x, vmVal, jitVal)
		}
	})
}
