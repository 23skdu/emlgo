package bytecode

import (
	"math"
	"testing"
)

// vmFunctions is every function name the VM supports, with a probe point and a
// reference value. The set is derived from the shared table so that a function
// added to it must also be covered here.
var vmFunctions = []struct {
	name  string
	probe float64
	want  float64
}{
	{"sin", 0.5, math.Sin(0.5)},
	{"cos", 0.5, math.Cos(0.5)},
	{"tan", 0.5, math.Tan(0.5)},
	{"asin", 0.5, math.Asin(0.5)},
	{"acos", 0.5, math.Acos(0.5)},
	{"atan", 0.5, math.Atan(0.5)},
	{"abs", -0.5, math.Abs(-0.5)},
	{"cbrt", 8, math.Cbrt(8)},
	{"log2", 8, math.Log2(8)},
	{"log10", 1000, math.Log10(1000)},
	{"ceil", 0.2, math.Ceil(0.2)},
	{"floor", 0.8, math.Floor(0.8)},
	{"trunc", 0.8, math.Trunc(0.8)},
	{"round", 3.7, math.Round(3.7)},
	{"sinh", 0.5, math.Sinh(0.5)},
	{"cosh", 0.5, math.Cosh(0.5)},
	{"tanh", 0.5, math.Tanh(0.5)},
	{"asinh", 0.5, math.Asinh(0.5)},
	{"acosh", 1.5, math.Acosh(1.5)},
	{"atanh", 0.5, math.Atanh(0.5)},
	{"erf", 0.5, math.Erf(0.5)},
	{"gamma", 5, math.Gamma(5)},
	{"sqrt", 4, math.Sqrt(4)},
	// exp and log go through the fastmath approximations inside the VM, so
	// they get a relative tolerance rather than exact equality.
	{"exp", 1, math.Exp(1)},
	{"ln", math.E, math.Log(math.E)},
}

// TestVMSupportsEveryFunction is the Step 5 gate: every function the JIT
// supports must compile and evaluate in the VM too.
func TestVMSupportsEveryFunction(t *testing.T) {
	// unaryFuncs holds 26 entries because "ln" is an alias of "log"; the test
	// covers 25 distinct function names.
	if len(unaryFuncs) != len(vmFunctions)+1 {
		t.Fatalf("test covers %d functions but the table has %d", len(vmFunctions), len(unaryFuncs))
	}

	for _, tc := range vmFunctions {
		t.Run(tc.name, func(t *testing.T) {
			expr := tc.name + "(x)"
			p, err := CompileExpr(expr)
			if err != nil {
				t.Fatalf("CompileExpr(%q): %v", expr, err)
			}

			got := p.Eval([]float64{tc.probe}, nil)
			// exp/ln use the fastmath polynomials; everything else is exact.
			tol := 0.0
			switch tc.name {
			case "exp", "ln":
				tol = 1e-5
			}
			if math.Abs(got-tc.want) > tol*math.Max(1, math.Abs(tc.want)) {
				t.Errorf("Eval(%q, %v) = %v, want %v", expr, tc.probe, got, tc.want)
			}
		})
	}
}

// TestVMBatchEvaluationAgreesWithScalar checks every function through both
// batch entry points, since they have their own dispatch chain.
func TestVMBatchEvaluationAgreesWithScalar(t *testing.T) {
	const n = 33 // deliberately not a multiple of any chunk size

	for _, tc := range vmFunctions {
		t.Run(tc.name, func(t *testing.T) {
			p, err := CompileExpr(tc.name + "(x)")
			if err != nil {
				t.Fatalf("CompileExpr: %v", err)
			}

			// Columnar layout: data[i] holds the values of variable i, so a
			// single-variable program is a single column.
			col := make([]float64, n)
			for i := range col {
				col[i] = tc.probe
			}
			// Row layout: data[i] is the variable vector for sample i.
			rows := make([][]float64, n)
			for i := range rows {
				rows[i] = []float64{tc.probe}
			}
			rowOut := make([]float64, n)

			scalar := p.Eval([]float64{tc.probe}, nil)

			p.EvalBatchColumnar([][]float64{col[:n]}, col[:n])
			p.EvalBatch(rows, rowOut)

			for i := 0; i < n; i++ {
				if col[i] != scalar {
					t.Errorf("EvalBatchColumnar[%d] = %v, want scalar %v", i, col[i], scalar)
					break
				}
				if rowOut[i] != scalar {
					t.Errorf("EvalBatch[%d] = %v, want scalar %v", i, rowOut[i], scalar)
					break
				}
			}
		})
	}
}

// TestVMFunctionOpArityAndLabel checks the per-opcode metadata stays consistent
// with the shared table.
func TestVMFunctionOpArityAndLabel(t *testing.T) {
	seen := map[OpCode][]string{}
	for name := range unaryFuncs {
		op, ok := unaryOpcode(name)
		if !ok {
			t.Errorf("unaryOpcode(%q) not found", name)
			continue
		}
		if op.Arity() != 1 {
			t.Errorf("%s (op %d) has arity %d, want 1", name, op, op.Arity())
		}
		if _, ok := funcName(op); !ok {
			t.Errorf("funcName(%d) not found for %q", op, name)
		}
		if op.String() == "" {
			t.Errorf("%s has an empty String()", name)
		}
		seen[op] = append(seen[op], name)
	}

	// "ln" is a deliberate alias of "log", so 26 names share 25 opcodes and the
	// only permitted collision is on OpLog.
	if len(seen) != len(unaryFuncs)-1 {
		t.Errorf("%d distinct opcodes for %d names; expected exactly one alias",
			len(seen), len(unaryFuncs))
	}
	for op, names := range seen {
		if len(names) == 1 {
			continue
		}
		sorted := append([]string(nil), names...)
		for i := range sorted {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[j] < sorted[i] {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		if op != OpLog || len(sorted) != 2 || sorted[0] != "ln" || sorted[1] != "log" {
			t.Errorf("op %d is shared by %v, want only the ln/log alias", op, sorted)
		}
	}
}

// TestVMRegularizedStaysFinite checks that no function can poison a genetic
// fitness with NaN, even when called far outside its domain.
func TestVMRegularizedStaysFinite(t *testing.T) {
	// Values chosen so that every function in the table is out of domain for
	// at least one of them.
	probes := []float64{-8, -1, -1e-9, 0, 1e-9, 1, 8, 1e8}

	for _, tc := range vmFunctions {
		t.Run(tc.name, func(t *testing.T) {
			p, err := CompileExpr(tc.name + "(x)")
			if err != nil {
				t.Fatalf("CompileExpr: %v", err)
			}
			for _, x := range probes {
				got := p.EvalRegularized([]float64{x}, 1e-6, nil)
				if math.IsNaN(got) {
					t.Errorf("EvalRegularized(%s, %v) = NaN", tc.name, x)
				}
			}
		})
	}
}

// TestVMUnknownFunctionStillRejected guards the error path after the table
// refactor.
func TestVMUnknownFunctionStillRejected(t *testing.T) {
	for _, expr := range []string{"bogus(x)", "logg(x)", "exp2(x)", "sinx(x)"} {
		if _, err := CompileExpr(expr); err == nil {
			t.Errorf("CompileExpr(%q) succeeded, want an error", expr)
		}
	}
}

// TestVMSqrtFoldingInOptimizer checks the optimizer's new sqrt case, which is
// reachable now that OpSqrt is in the table.
func TestVMSqrtFoldingInOptimizer(t *testing.T) {
	p, err := CompileExpr("sqrt(16)")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}
	got := p.Eval(nil, nil)
	if got != 4 {
		t.Errorf("Eval = %v, want 4", got)
	}
	if len(p.Ops) != 2 || p.Ops[0] != OpConst || p.Ops[1] != OpSqrt {
		t.Errorf("Ops = %v, want [CONST SQRT]", p.Ops)
	}
}

// TestOptimizerFoldsEveryFunction checks that Optimize collapses a constant
// unary-function subtree for every opcode, not just the handful that had
// explicit cases before the table refactor.
func TestOptimizerFoldsEveryFunction(t *testing.T) {
	for _, tc := range vmFunctions {
		t.Run(tc.name, func(t *testing.T) {
			p, err := CompileExpr(tc.name + "(16)")
			if err != nil {
				t.Fatalf("CompileExpr: %v", err)
			}
			before := p.Eval(nil, nil)
			opt := Optimize(p)

			if len(opt.Ops) != 1 || opt.Ops[0] != OpConst {
				t.Errorf("Optimize(%s(16)) = %v, want a single CONST", tc.name, opt.Ops)
				return
			}
			after := opt.Eval(nil, nil)
			if math.Abs(after-before) > 1e-9*math.Max(1, math.Abs(before)) {
				t.Errorf("Optimize(%s(16)) = %v, want %v", tc.name, after, before)
			}
		})
	}
}
