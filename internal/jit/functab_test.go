package jit

import (
	"math"
	"testing"
)

// mathFuncNames is the authoritative list of built-in math functions. Tests
// iterate over it so that adding a function to the dispatch table without
// updating the parser or the arena name buffer is caught.
var mathFuncNames = []string{
	"sin", "cos", "exp", "log", "sqrt", "tan", "asin", "acos", "atan", "abs",
	"cbrt", "log2", "log10", "ceil", "floor", "trunc", "round",
	"sinh", "cosh", "tanh", "asinh", "acosh", "atanh", "erf", "gamma",
}

// TestMathFuncsCoversParser ensures every function the parser accepts has a
// dispatch entry, and vice versa.
func TestMathFuncsCoversParser(t *testing.T) {
	for _, name := range mathFuncNames {
		if _, ok := mathFuncs[name]; !ok {
			t.Errorf("mathFuncs missing %q accepted by the parser", name)
		}
		if !isFuncName(name) {
			t.Errorf("mathFuncs has %q but the parser rejects it", name)
		}
	}
	if len(mathFuncs) != len(mathFuncNames) {
		t.Errorf("mathFuncs has %d entries, want %d (table and test drifted apart)",
			len(mathFuncs), len(mathFuncNames))
	}
}

// TestCallMathFuncByNameUnknown checks unknown names are reported as unknown
// rather than silently evaluating to zero.
func TestCallMathFuncByNameUnknown(t *testing.T) {
	for _, name := range []string{"", "bogus", "sin2", "SIN"} {
		if _, ok := callMathFuncByName(name, 1.0); ok {
			t.Errorf("callMathFuncByName(%q) reported ok, want not-found", name)
		}
	}
}

// TestEvalSupportsAllFuncNames evaluates every built-in function through the
// tree interpreter and the arena interpreter and requires both to agree with
// the reference result from mathFuncs.
func TestEvalSupportsAllFuncNames(t *testing.T) {
	// 0.5 is in the domain of every function in the table.
	const arg = 0.5

	for _, name := range mathFuncNames {
		want, ok := callMathFuncByName(name, arg)
		if !ok {
			t.Fatalf("callMathFuncByName(%q) not found", name)
		}

		node := FunctionCall{Name: name, Arg: Variable{Name: "x"}}
		if got := Eval(node, arg); math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
			t.Errorf("Eval(%s(x)) = %v, want %v", name, got, want)
		}

		arena := NewArena(4)
		arenaNode := arena.FromInterface(node)
		if arenaNode == nil {
			t.Fatalf("FromInterface returned nil for %q", name)
		}
		if got := EvalArena(arenaNode, arg); math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
			t.Errorf("EvalArena(%s(x)) = %v, want %v", name, got, want)
		}
	}
}

// TestEvalArenaFunctionNameBufferBoundary checks that the 16-byte fixed name
// field in ArenaNode is wide enough for every built-in function name.
func TestEvalArenaFunctionNameBufferBoundary(t *testing.T) {
	longest := ""
	for _, name := range mathFuncNames {
		if len(name) > len(longest) {
			longest = name
		}
	}
	// ArenaNode.Name is [16]byte and FromInterface copies at most 15 bytes
	// plus a NUL terminator.
	if len(longest) > 15 {
		t.Fatalf("longest function name %q (%d bytes) does not fit the arena name buffer", longest, len(longest))
	}
}

// TestEvalArenaTruncatedNameIsNotEvaluated documents that a name longer than
// the arena buffer is truncated and therefore does not silently evaluate to
// the correct result. Guards against regressing into a silent-zero bug.
func TestEvalArenaTruncatedNameIsNotEvaluated(t *testing.T) {
	arena := NewArena(4)
	node := arena.Alloc()
	node.Kind = kindFuncCall
	copy(node.Name[:], "definitelynotafunc")
	node.Left = arena.Alloc()
	node.Left.Kind = kindVariable
	// Unknown names evaluate to 0 rather than panicking.
	if got := EvalArena(node, 1.0); got != 0 {
		t.Errorf("EvalArena with unknown name = %v, want 0", got)
	}
}

// TestSimplifyEmlPatternsPreserveValue extends the property test to the raw
// eml(u,v) patterns that the identity reductions in Simplify operate on.
// Reductions legitimately reassociate arithmetic, so results are compared with
// a tight relative tolerance rather than bit-for-bit.
func TestSimplifyEmlPatternsPreserveValue(t *testing.T) {
	x := varNode()
	one := constNode(1)
	two := constNode(2)

	cases := []struct {
		name string
		node *EMLNode
	}{
		{"eml(x,1)", emlNode(x, one)},
		{"eml(1,x)", emlNode(one, x)},
		{"eml(x,x)", emlNode(x, x)},
		{"eml(1,1)", emlNode(one, one)},
		{"eml(0,1)", emlNode(constNode(0), one)},
		{"eml(0,x)", emlNode(constNode(0), x)},
		{"eml(1,eml(eml(1,x),1))", emlNode(one, emlNode(emlNode(one, x), one))},
		{"eml(eml(1,eml(x,1)),eml(1,1))", emlNode(emlNode(one, emlNode(x, one)), emlNode(one, one))},
		{"eml(eml(1,x),eml(x,1))", emlNode(emlNode(one, x), emlNode(x, one))},
		{"eml(two,two)", emlNode(two, two)},
		{"neg(const)", emlFuncUnary("neg", two)},
		{"neg(neg(x))", emlFuncUnary("neg", emlFuncUnary("neg", x))},
		{"add(x,0)", emlFuncBinary("add", x, constNode(0))},
		{"add(0,x)", emlFuncBinary("add", constNode(0), x)},
		{"sub(x,0)", emlFuncBinary("sub", x, constNode(0))},
		{"mul(x,0)", emlFuncBinary("mul", x, constNode(0))},
		{"mul(x,1)", emlFuncBinary("mul", x, one)},
		{"div(x,1)", emlFuncBinary("div", x, one)},
		{"pow(x,0)", emlFuncBinary("pow", x, constNode(0))},
		{"pow(x,1)", emlFuncBinary("pow", x, one)},
		{"sqrt(const)", emlFuncUnary("sqrt", two)},
		{"sqrt(4)", emlFuncUnary("sqrt", constNode(4))},
	}

	for _, c := range cases {
		simplified := Simplify(c.node)
		if simplified == nil {
			t.Errorf("%s: Simplify returned nil", c.name)
			continue
		}
		// Values well inside the real domain keep both evaluations finite.
		for _, xv := range []float64{0.5, 1.0, 1.5, 2.0} {
			want := EMLEval(c.node, xv)
			got := EMLEval(simplified, xv)
			if math.IsNaN(want) && math.IsNaN(got) {
				continue
			}
			if math.Abs(got-want) > 1e-13*math.Max(1, math.Abs(want)) {
				t.Errorf("%s at x=%v: Simplify changed value %.17g -> %.17g (%s)",
					c.name, xv, want, got, Decompile(simplified))
			}
		}
	}
}

// TestSimplifyNilAndTrivial covers the trivial inputs.
func TestSimplifyNilAndTrivial(t *testing.T) {
	if got := Simplify(nil); got != nil {
		t.Errorf("Simplify(nil) = %v, want nil", got)
	}
	if got := Simplify(constNode(3.5)); got == nil || got.Value != 3.5 {
		t.Errorf("Simplify(const) = %v, want const 3.5", got)
	}
	if got := Simplify(varNode()); got == nil || got.Kind != EMLVar {
		t.Errorf("Simplify(var) = %v, want var", got)
	}
}

// TestSimplifyDoesNotInventNegationOrReciprocal locks in that the removed
// (and non-identity) reductions stay removed. Each pattern must evaluate to the
// same value before and after simplification, which fails if a bogus
// negation or reciprocal rewrite is reintroduced.
func TestSimplifyDoesNotInventNegationOrReciprocal(t *testing.T) {
	x := varNode()
	one := constNode(1)
	patterns := map[string]*EMLNode{
		"eml(eml(1,eml(x,1)),eml(1,1))": emlNode(
			emlNode(one, emlNode(x, one)),
			emlNode(one, one),
		),
		"eml(eml(1,x),eml(x,1))": emlNode(
			emlNode(one, x),
			emlNode(x, one),
		),
	}
	for name, node := range patterns {
		simplified := Simplify(node)
		for _, xv := range []float64{0.5, 0.75, 1.25} {
			want := EMLEval(node, xv)
			got := EMLEval(simplified, xv)
			if math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
				t.Errorf("%s at x=%v: got %v, want %v (simplified to %s)",
					name, xv, got, want, Decompile(simplified))
			}
		}
		// The naive -x and 1/x rewrites would produce these; neither may appear.
		if got := Decompile(simplified); got == "-x" || got == "1 / x" {
			t.Errorf("%s simplified to %q, which is not an identity for eml(u,v)=exp(u)-ln(v)",
				name, got)
		}
	}
}
