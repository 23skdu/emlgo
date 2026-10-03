package jit

import (
	"math"
	"strings"
	"testing"
)

func TestStep34_DiffIsZeroDerivative(t *testing.T) {
	for _, name := range []string{"ceil", "floor", "trunc", "round"} {
		if !diffIsZeroDerivative(name) {
			t.Errorf("diffIsZeroDerivative(%q) = false, want true", name)
		}
		node := &EMLNode{Kind: EMLFunc, Name: name, Left: varNode()}
		d := Diff(node)
		if d == nil || d.Kind != EMLConst || d.Value != 0 {
			t.Errorf("Diff(%s(x)) = %v, want 0", name, d)
		}
	}
	if diffIsZeroDerivative("unknown") {
		t.Error("diffIsZeroDerivative(unknown) = true, want false")
	}
}

func TestStep34_GammaDerivativeNaN(t *testing.T) {
	gammaNode := &EMLNode{Kind: EMLFunc, Name: "gamma", Left: varNode()}
	d := Diff(gammaNode)
	if d == nil || d.Kind != EMLConst || !math.IsNaN(d.Value) {
		t.Errorf("Diff(gamma(x)) = %v, want NaN", d)
	}
	val := DiffEval(gammaNode, 2.5)
	if !math.IsNaN(val) {
		t.Errorf("DiffEval(gamma, 2.5) = %v, want NaN", val)
	}
}

func TestStep34_SimplifySubZero(t *testing.T) {
	// 0 - x -> neg(x)
	expr := emlFuncBinary("sub", constNode(0), varNode())
	simplified := Simplify(expr)
	if simplified.Kind != EMLFunc || simplified.Name != "neg" || simplified.Left.Kind != EMLVar {
		t.Errorf("Simplify(0 - x) = %v, want neg(x)", simplified)
	}
}

func TestStep34_SecondDerivativeEval(t *testing.T) {
	// f(x) = x^3 => f'(x) = 3x^2 => f''(x) = 6x. At x = 2, f''(2) = 12.
	cube := emlFuncBinary("pow", varNode(), constNode(3))
	d1 := DiffEval(cube, 2.0)
	if math.Abs(d1-12.0) > 1e-9 {
		t.Errorf("DiffEval(x^3, 2) = %v, want 12", d1)
	}
	d2 := SecondDerivativeEval(cube, 2.0)
	if math.Abs(d2-12.0) > 1e-9 {
		t.Errorf("SecondDerivativeEval(x^3, 2) = %v, want 12", d2)
	}

	// f(x) = sin(x) => f'(x) = cos(x) => f''(x) = -sin(x). At x = 0.5, f''(0.5) = -sin(0.5).
	sinExpr := CanonicalSin(varNode())
	d2Sin := SecondDerivativeEval(sinExpr, 0.5)
	wantSin2 := -math.Sin(0.5)
	if math.Abs(d2Sin-wantSin2) > 1e-9 {
		t.Errorf("SecondDerivativeEval(sin(x), 0.5) = %v, want %v", d2Sin, wantSin2)
	}

	// Compare with central second finite difference
	h := 1e-5
	fPlus := EMLEval(sinExpr, 0.5+h)
	fMid := EMLEval(sinExpr, 0.5)
	fMinus := EMLEval(sinExpr, 0.5-h)
	num2 := (fPlus - 2*fMid + fMinus) / (h * h)
	if math.Abs(d2Sin-num2) > 1e-4 {
		t.Errorf("SecondDerivativeEval vs central difference: got %v, want %v", d2Sin, num2)
	}
}

func TestStep34_All25Functions(t *testing.T) {
	funcs := []struct {
		name         string
		testPoint    float64
		smoothDiff   bool
		expectedDiff float64
	}{
		{"sin", 0.5, true, math.Cos(0.5)},
		{"cos", 0.5, true, -math.Sin(0.5)},
		{"exp", 0.5, true, math.Exp(0.5)},
		{"log", 2.0, true, 1.0 / 2.0},
		{"sqrt", 2.0, true, 0.5 / math.Sqrt(2.0)},
		{"tan", 0.5, true, 1.0 / (math.Cos(0.5) * math.Cos(0.5))},
		{"asin", 0.5, true, 1.0 / math.Sqrt(1.0-0.25)},
		{"acos", 0.5, true, -1.0 / math.Sqrt(1.0-0.25)},
		{"atan", 0.5, true, 1.0 / (1.0 + 0.25)},
		{"abs", 0.5, true, 1.0},
		{"cbrt", 2.0, true, 1.0 / (3.0 * math.Pow(2.0, 2.0/3.0))},
		{"log2", 2.0, true, 1.0 / (2.0 * math.Ln2)},
		{"log10", 2.0, true, 1.0 / (2.0 * math.Ln10)},
		{"ceil", 2.3, false, 0.0},
		{"floor", 2.3, false, 0.0},
		{"trunc", 2.3, false, 0.0},
		{"round", 2.3, false, 0.0},
		{"sinh", 0.5, true, math.Cosh(0.5)},
		{"cosh", 0.5, true, math.Sinh(0.5)},
		{"tanh", 0.5, true, 1.0 - math.Pow(math.Tanh(0.5), 2)},
		{"asinh", 0.5, true, 1.0 / math.Sqrt(0.25+1.0)},
		{"acosh", 2.0, true, 1.0 / math.Sqrt(4.0-1.0)},
		{"atanh", 0.5, true, 1.0 / (1.0 - 0.25)},
		{"erf", 0.5, true, (2.0 / math.Sqrt(math.Pi)) * math.Exp(-0.25)},
		{"gamma", 2.5, false, math.NaN()},
	}

	for _, fn := range funcs {
		t.Run(fn.name, func(t *testing.T) {
			astNode := FunctionCall{Name: fn.name, Arg: Variable{Name: "x"}}
			emlAst := &EMLNode{Kind: EMLFunc, Name: fn.name, Left: varNode()}

			// 1. EMLEval vs Eval
			emlVal := EMLEval(emlAst, fn.testPoint)
			evalVal := Eval(astNode, fn.testPoint)
			if math.IsNaN(emlVal) && math.IsNaN(evalVal) {
				// OK
			} else if math.Abs(emlVal-evalVal) > 1e-12 {
				t.Errorf("%s EMLEval(%v) = %v, Eval = %v", fn.name, fn.testPoint, emlVal, evalVal)
			}

			// 2. DiffEval
			diffVal := DiffEval(emlAst, fn.testPoint)
			if !fn.smoothDiff {
				if math.IsNaN(fn.expectedDiff) {
					if !math.IsNaN(diffVal) {
						t.Errorf("%s DiffEval = %v, want NaN", fn.name, diffVal)
					}
				} else {
					if diffVal != fn.expectedDiff {
						t.Errorf("%s DiffEval = %v, want %v", fn.name, diffVal, fn.expectedDiff)
					}
				}
			} else {
				if math.Abs(diffVal-fn.expectedDiff) > 1e-7 {
					t.Errorf("%s DiffEval = %v, want %v", fn.name, diffVal, fn.expectedDiff)
				}
				// Check central difference approximation
				h := 1e-6
				fPlus := EMLEval(emlAst, fn.testPoint+h)
				fMinus := EMLEval(emlAst, fn.testPoint-h)
				numDiff := (fPlus - fMinus) / (2 * h)
				if math.Abs(diffVal-numDiff) > 1e-4 {
					t.Errorf("%s DiffEval = %v vs numerical %v", fn.name, diffVal, numDiff)
				}

				// Check SecondDerivativeEval vs numerical 2nd difference
				d2Val := SecondDerivativeEval(emlAst, fn.testPoint)
				fMid := EMLEval(emlAst, fn.testPoint)
				num2Diff := (fPlus - 2*fMid + fMinus) / (h * h)
				if math.Abs(d2Val-num2Diff) > 1e-3 {
					t.Errorf("%s SecondDerivativeEval = %v vs numerical %v", fn.name, d2Val, num2Diff)
				}
			}
		})
	}
}

func TestStep34_DeepRecursionLimits(t *testing.T) {
	// Build an EMLNode tree with depth 1005 (exceeding MaxASTDepth = 1000)
	var deepEML *EMLNode = varNode()
	for i := 0; i < 1005; i++ {
		deepEML = &EMLNode{Kind: EMLFunc, Name: "sin", Left: deepEML}
	}

	if sz := emlSizeDepth(deepEML, MaxASTDepth+1); sz != 0 {
		t.Errorf("emlSizeDepth(deep, MaxASTDepth+1) = %d, want 0", sz)
	}
	if d := depthWithLimit(deepEML, MaxASTDepth+1); d != 0 {
		t.Errorf("depthWithLimit(deep, MaxASTDepth+1) = %d, want 0", d)
	}
	if v := EMLEval(deepEML, 0.5); v != 0 {
		t.Errorf("EMLEval(deep) = %v, want 0", v)
	}
	if v := EMLEvalRegularized(deepEML, 0.5, 1e-6); v != 0 {
		t.Errorf("EMLEvalRegularized(deep) = %v, want 0", v)
	}
	if s := Simplify(deepEML); s == nil {
		t.Errorf("Simplify(deep) returned nil")
	}
	if diff := Diff(deepEML); diff == nil || diff.Kind != EMLConst || diff.Value != 0 {
		t.Errorf("Diff(deep) = %v, want constNode(0)", diff)
	}
	if eq := Equiv(deepEML, deepEML); eq {
		t.Errorf("Equiv(deep, deep) = true, want false")
	}
	if s := decompileDepth(deepEML, MaxASTDepth+1); s != "" {
		t.Errorf("decompileDepth(deep, MaxASTDepth+1) = %q, want empty", s)
	}
	if s := decompileLaTeXDepth(deepEML, MaxASTDepth+1); s != "" {
		t.Errorf("decompileLaTeXDepth(deep, MaxASTDepth+1) = %q, want empty", s)
	}

	// Test wrapLatexPow directly
	if got := wrapLatexPow(nil); got != "" {
		t.Errorf("wrapLatexPow(nil) = %q, want empty", got)
	}
	powNode := &EMLNode{Kind: EMLFunc, Name: "add", Left: varNode(), Right: constNode(1)}
	if got := wrapLatexPow(powNode); got != "(x + 1.0)" {
		t.Errorf("wrapLatexPow(add) = %q, want (x + 1.0)", got)
	}

	// Build a Node tree with depth 1005
	var deepNode Node = Variable{Name: "x"}
	for i := 0; i < 1005; i++ {
		deepNode = FunctionCall{Name: "sin", Arg: deepNode}
	}

	if c := canonicalizeDepth(deepNode, MaxASTDepth+1); c != nil {
		t.Errorf("canonicalizeDepth(deepNode, MaxASTDepth+1) = %v, want nil", c)
	}
	if v := Eval(deepNode, 0.5); v != 0 {
		t.Errorf("Eval(deepNode) = %v, want 0", v)
	}
	if err := validateVars(deepNode, map[string]bool{"x": true}); err == nil {
		t.Errorf("validateVars(deepNode) expected error, got nil")
	}
	if s := FormatExpr(deepNode); !strings.Contains(s, "...") {
		t.Errorf("FormatExpr(deepNode) = %q, want containing '...'", s)
	}

	// Test parser nesting depth limit
	deepInput := strings.Repeat("(", 1005) + "x" + strings.Repeat(")", 1005)
	_, err := Parse(deepInput)
	if err == nil {
		t.Errorf("Parse(1005 parens) expected error, got nil")
	} else if !strings.Contains(err.Error(), "maximum depth") {
		t.Errorf("Parse(1005 parens) error = %v, want maximum depth exceeded", err)
	}

	// Test Arena with depth limit
	arena := NewArena(2048)
	var deepArena *ArenaNode = arena.Alloc()
	deepArena.Kind = kindVariable
	for i := 0; i < 1005; i++ {
		n := arena.Alloc()
		n.Kind = kindUnaryOp
		n.Op = '-'
		n.Left = deepArena
		deepArena = n
	}
	if iface := arena.toInterfaceDepth(deepArena, MaxASTDepth+1); iface != nil {
		t.Errorf("arena.toInterfaceDepth(deep, MaxASTDepth+1) = %v, want nil", iface)
	}
	if aNode := arena.fromInterfaceDepth(deepNode, MaxASTDepth+1); aNode != nil {
		t.Errorf("arena.fromInterfaceDepth(deep, MaxASTDepth+1) = %v, want nil", aNode)
	}
	if v := EvalArena(deepArena, 0.5); v != 0 {
		t.Errorf("EvalArena(deep) = %v, want 0", v)
	}
}

func TestStep34_FormattingWrappers(t *testing.T) {
	if wrapDecomp(nil, "+", true) != "" {
		t.Error("wrapDecomp(nil) should be empty")
	}
	powNode := &EMLNode{Kind: EMLFunc, Name: "pow", Left: varNode(), Right: constNode(2.0)}
	if got := wrapDecomp(powNode, "pow", true); got != "(x^2.0)" {
		t.Errorf("wrapDecomp(pow) = %q, want (x^2.0)", got)
	}

	x := Variable{Name: "x"}
	two := Number{Value: 2}
	bin := BinaryOp{Op: '+', Left: x, Right: two}
	if got := wrapParen(bin); got != "(x + 2)" {
		t.Errorf("wrapParen(bin) = %q, want (x + 2)", got)
	}
	if got := wrapBinOp(bin, '*', true); got != "(x + 2)" {
		t.Errorf("wrapBinOp(bin) = %q, want (x + 2)", got)
	}
	if got := wrapPower(bin); got != "(x + 2)" {
		t.Errorf("wrapPower(bin) = %q, want (x + 2)", got)
	}
}
