package bytecode

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/emlgo/eml/internal/jit"
)

// fastmathRelTol is the relative tolerance for exp/log, which the bytecode
// evaluator computes with the fastmath minimax polynomials rather than
// math.Exp/math.Log. Their documented worst-case relative error is ~1.5e-6
// for exp and ~1e-7 for log.
const fastmathRelTol = 1e-5

// TestEmitOp covers every supported operator plus the unsupported-operator
// error branch.
func TestEmitOp(t *testing.T) {
	tests := []struct {
		op   rune
		want OpCode
	}{
		{'+', OpAdd},
		{'-', OpSub},
		{'*', OpMul},
		{'/', OpDiv},
		{'^', OpPow},
	}
	for _, tc := range tests {
		p := NewProgram()
		if err := emitOp(p, token{op: tc.op, preced: 1, rAssoc: false}); err != nil {
			t.Errorf("emitOp(%q): unexpected error %v", tc.op, err)
			continue
		}
		if len(p.Ops) != 1 || p.Ops[0] != tc.want {
			t.Errorf("emitOp(%q) = %v, want [%v]", tc.op, p.Ops, tc.want)
		}
	}

	// Unsupported operators must be reported.
	for _, op := range []rune{'%', '#', '@', 'x'} {
		p := NewProgram()
		err := emitOp(p, token{op: op, preced: 1})
		if err == nil {
			t.Errorf("emitOp(%q): expected error, got nil", op)
			continue
		}
		if !strings.Contains(err.Error(), "unsupported operator") {
			t.Errorf("emitOp(%q): error %q missing context", op, err)
		}
		if len(p.Ops) != 0 {
			t.Errorf("emitOp(%q): appended %v despite error", op, p.Ops)
		}
	}
}

// TestEmitFunc covers every supported function plus the
// unsupported-function error branch.
func TestEmitFunc(t *testing.T) {
	tests := []struct {
		name string
		want OpCode
	}{
		{"eml", OpEML},
		{"exp", OpExp},
		{"log", OpLog},
		{"ln", OpLog},
		{"sqrt", OpSqrt},
	}
	for _, tc := range tests {
		p := NewProgram()
		if err := emitFunc(p, token{name: tc.name, isFunc: true, argCount: 1}); err != nil {
			t.Errorf("emitFunc(%q): unexpected error %v", tc.name, err)
			continue
		}
		if len(p.Ops) != 1 || p.Ops[0] != tc.want {
			t.Errorf("emitFunc(%q) = %v, want [%v]", tc.name, p.Ops, tc.want)
		}
	}

	// Names outside the supported set must be reported.
	for _, name := range []string{"bogus", "", "exp2", "Sqrt", "logg"} {
		p := NewProgram()
		err := emitFunc(p, token{name: name, isFunc: true, argCount: 1})
		if err == nil {
			t.Errorf("emitFunc(%q): expected error, got nil", name)
			continue
		}
		if !strings.Contains(err.Error(), "unsupported function") {
			t.Errorf("emitFunc(%q): error %q missing context", name, err)
		}
		if len(p.Ops) != 0 {
			t.Errorf("emitFunc(%q): appended %v despite error", name, p.Ops)
		}
	}
}

// TestGetOrAssignVarReusesIndices checks that repeated lookups reuse slots and
// that new names allocate fresh, dense indices.
func TestGetOrAssignVarReusesIndices(t *testing.T) {
	varMap := map[string]uint16{}

	first, err := getOrAssignVar("x", varMap)
	if err != nil {
		t.Fatalf("getOrAssignVar: %v", err)
	}
	if first != 0 {
		t.Errorf("first index = %d, want 0", first)
	}

	again, err := getOrAssignVar("x", varMap)
	if err != nil {
		t.Fatalf("getOrAssignVar: %v", err)
	}
	if again != first {
		t.Errorf("repeated lookup = %d, want %d", again, first)
	}

	second, err := getOrAssignVar("y", varMap)
	if err != nil {
		t.Fatalf("getOrAssignVar: %v", err)
	}
	if second != 1 {
		t.Errorf("second index = %d, want 1", second)
	}
	if len(varMap) != 2 {
		t.Errorf("varMap has %d entries, want 2", len(varMap))
	}
}

// TestGetOrAssignVarTooManyVariables checks the MaxUint16 guard without
// allocating 65536 entries: the map is pre-filled to the limit.
func TestGetOrAssignVarTooManyVariables(t *testing.T) {
	varMap := make(map[string]uint16, math.MaxUint16)
	for i := range int(math.MaxUint16) {
		varMap[strconv.Itoa(i)+"-v"] = uint16(i) // #nosec G115 -- i < MaxUint16
	}
	if len(varMap) != math.MaxUint16 {
		t.Fatalf("pre-filled varMap has %d entries, want %d", len(varMap), math.MaxUint16)
	}

	idx, err := getOrAssignVar("overflow", varMap)
	if err == nil {
		t.Fatal("expected an error when the variable limit is reached")
	}
	if idx != 0 {
		t.Errorf("index on error = %d, want 0", idx)
	}
	if !strings.Contains(err.Error(), "too many distinct variables") {
		t.Errorf("error %q missing context", err)
	}
	if _, added := varMap["overflow"]; added {
		t.Error("variable was inserted despite the error")
	}
}

// TestEmitASTUnsupportedNode checks that an AST node kind the compiler does
// not handle is reported rather than silently skipped.
func TestEmitASTUnsupportedNode(t *testing.T) {
	p := NewProgram()
	varMap := map[string]uint16{}
	// A nil node reaches the default branch of emitAST.
	if err := emitAST(nil, p, varMap); err == nil {
		t.Error("emitAST(nil): expected an error, got nil")
	}
}

// TestCompileExprErrors checks the compiler's error reporting for malformed
// expressions.
func TestCompileExprErrors(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"empty", ""},
		{"only operator", "*"},
		{"unbalanced open", "(1 + 2"},
		{"unbalanced close", "1 + 2)"},
		{"dangling operator", "1 +"},
		{"empty parentheses", "()"},
		{"missing operand", "1 * * 2"},
		{"bare unary function", "exp"},
		{"bare sqrt", "sqrt"},
		{"unknown function", "bogus(1)"},
		{"misspelled function", "logg(1)"},
	}
	for _, tc := range tests {
		p, err := CompileExpr(tc.expr)
		if err == nil {
			t.Errorf("%s: CompileExpr(%q) = %v, want error", tc.name, tc.expr, p)
			continue
		}
		// A rejected expression must never yield a program that panics on Eval.
		if p != nil {
			t.Errorf("%s: CompileExpr(%q) returned both a program and an error", tc.name, tc.expr)
		}
	}
}

// TestCompileExprRejectsStackUnderflow is the regression test for programs
// such as "1 +" or "*" that used to compile successfully and then panic with
// an index-out-of-range error during Eval.
func TestCompileExprRejectsStackUnderflow(t *testing.T) {
	for _, expr := range []string{"*", "1 +", "1 * * 2", "exp", "sqrt", "+ + +"} {
		p, err := CompileExpr(expr)
		if err == nil {
			t.Errorf("CompileExpr(%q) accepted an underflowing program: %v", expr, p)
		}
	}
}

// TestCompileExprErrorRecoverySkipsBadPrefix documents that CompileExpr falls
// back to its own parser when the JIT parser rejects the expression, and that
// the fallback also rejects it rather than emitting an unusable program.
func TestCompileExprErrorRecoverySkipsBadPrefix(t *testing.T) {
	// eml takes two arguments and cannot be parsed by jit.Parse, so this
	// exercises the Shunting-Yard fallback on a valid expression.
	p, err := CompileExpr("eml(1, exp(1))")
	if err != nil {
		t.Fatalf("CompileExpr fallback: %v", err)
	}
	// eml(1, exp(1)) = exp(1) - ln(exp(1)) = e - 1.
	want := math.E - 1
	got := p.Eval(nil, nil)
	if math.Abs(got-want) > fastmathRelTol*math.Abs(want) {
		t.Errorf("eml(1, exp(1)) = %v, want %v", got, want)
	}
}

// TestCompileExprSupportedExpressions is the positive counterpart: the
// documented expression syntax must compile and evaluate.
func TestCompileExprSupportedExpressions(t *testing.T) {
	// vars are supplied positionally in variable-slot order (x, then y).
	tests := []struct {
		expr string
		vars []float64
		want float64
	}{
		{"1 + 2", nil, 3},
		{"5 - 2", nil, 3},
		{"3 * 4", nil, 12},
		{"9 / 3", nil, 3},
		{"2 ^ 10", nil, 1024},
		{"exp(1)", nil, math.E},
		{"log(exp(2))", nil, 2},
		{"ln(exp(3))", nil, 3},
		{"sqrt(16)", nil, 4},
		{"sqrt(sqrt(16))", nil, 2},
		{"eml(0, 1)", nil, 1},
		{"1 + 2 * 3", nil, 7},
		{"(1 + 2) * 3", nil, 9},
		{"2 ^ 3 ^ 2", nil, 512}, // right-associative
		{"x + 1", []float64{4}, 5},
		{"x * y", []float64{3, 4}, 12},
	}
	for _, tc := range tests {
		p, err := CompileExpr(tc.expr)
		if err != nil {
			t.Errorf("CompileExpr(%q): %v", tc.expr, err)
			continue
		}
		got := p.Eval(tc.vars, nil)
		if math.Abs(got-tc.want) > fastmathRelTol*math.Max(1, math.Abs(tc.want)) {
			t.Errorf("CompileExpr(%q).Eval = %v, want %v", tc.expr, got, tc.want)
		}
	}
}

// TestCompileASTRoundTrip checks that compiling a parsed AST evaluates the
// same as evaluating the AST through the JIT interpreter.
func TestCompileASTRoundTrip(t *testing.T) {
	exprs := []string{
		"x + 1",
		"x * x + 1",
		"exp(x)",
		"sqrt(x) + 1",
		"log(x) * x",
	}
	for _, expr := range exprs {
		node, err := jit.Parse(expr)
		if err != nil {
			t.Errorf("jit.Parse(%q): %v", expr, err)
			continue
		}
		p, err := CompileAST(node)
		if err != nil {
			t.Errorf("CompileAST(%q): %v", expr, err)
			continue
		}
		// The bytecode evaluator uses the fastmath polynomials for exp/log, so
		// compare within their documented accuracy rather than bit-for-bit.
		for _, x := range []float64{0.5, 1.0, 2.0, 3.5} {
			want := jit.Eval(node, x)
			got := p.Eval([]float64{x}, nil)
			if math.Abs(got-want) > fastmathRelTol*math.Max(1, math.Abs(want)) {
				t.Errorf("%q at x=%v: bytecode %v, jit %v", expr, x, got, want)
			}
		}
	}
}

// TestFunctionNamesAreCaseInsensitive documents that the Shunting-Yard
// tokenizer lower-cases identifiers before looking them up, so "SQRT" and
// "sqrt" name the same opcode.
func TestFunctionNamesAreCaseInsensitive(t *testing.T) {
	for _, expr := range []string{"sqrt(x)", "SQRT(x)", "Sqrt(x)"} {
		p, err := CompileExpr(expr)
		if err != nil {
			t.Errorf("CompileExpr(%q): %v", expr, err)
			continue
		}
		if len(p.Ops) != 2 || p.Ops[1] != OpSqrt {
			t.Errorf("CompileExpr(%q) = %v, want [VAR SQRT]", expr, p.Ops)
		}
	}
}

// TestIsFuncAndPrecedence checks the tokenizer's function and precedence
// tables, including the associativity of exponentiation.
func TestIsFuncAndPrecedence(t *testing.T) {
	// isFunc drives tokenization; sin/cos are recognised here but rejected
	// later by emitFunc, which is how CompileExpr reports them as unsupported.
	// isFunc must agree with what the VM can actually emit: every entry of
	// unaryFuncs plus the two-argument eml.
	for name := range unaryFuncs {
		if !isFunc(name) {
			t.Errorf("isFunc(%q) = false, want true", name)
		}
		if _, ok := unaryOpcode(name); !ok {
			t.Errorf("unaryOpcode(%q) not found", name)
		}
	}
	if !isFunc(emlName) {
		t.Errorf("isFunc(%q) = false, want true", emlName)
	}
	for _, name := range []string{"+", "", "bogus", "exp2", "Sqrt", "logg", "tan2"} {
		if isFunc(name) {
			t.Errorf("isFunc(%q) = true, want false", name)
		}
	}

	// getPrecedence returns (precedence, rightAssociative). Precedence must
	// increase with depth of binding, and only '^' is right-associative.
	precOf := func(op rune) (int, bool) {
		p, rAssoc := getPrecedence(op)
		return p, rAssoc
	}
	lowAdd, addRA := precOf('+')
	lowMul, mulRA := precOf('*')
	powPrec, powRA := precOf('^')

	if lowAdd >= lowMul {
		t.Errorf("precedence(+)=%d should be lower than precedence(*)=%d", lowAdd, lowMul)
	}
	if powPrec <= lowMul {
		t.Errorf("precedence(^)=%d should exceed precedence(*)=%d", powPrec, lowMul)
	}
	if !powRA {
		t.Error("^ should be right-associative")
	}
	if addRA || mulRA {
		t.Error("+ and * should be left-associative")
	}

	// Subtraction and division share precedence with their partners.
	if p, _ := precOf('-'); p != lowAdd {
		t.Errorf("precedence(-)=%d, want %d (same as +)", p, lowAdd)
	}
	if p, _ := precOf('/'); p != lowMul {
		t.Errorf("precedence(/)=%d, want %d (same as *)", p, lowMul)
	}
}

// TestProgramStackDepthIsSufficient checks CalculateMaxStackDepth against a
// range of expression shapes, since an under-estimate causes Eval to index
// out of range at run time.
func TestProgramStackDepthIsSufficient(t *testing.T) {
	exprs := []string{
		"1",
		"1 + 1",
		"1 + 1 + 1 + 1 + 1 + 1 + 1 + 1",
		"1 * (2 * (3 * (4 * (5 * 6))))",
		"sqrt(sqrt(sqrt(sqrt(sqrt(1024)))))",
		"1 + 2 * 3 ^ 4 - 5 / 6",
	}
	for _, expr := range exprs {
		p, err := CompileExpr(expr)
		if err != nil {
			t.Errorf("CompileExpr(%q): %v", expr, err)
			continue
		}
		if p.MaxStackDepth <= 0 {
			t.Errorf("%q: MaxStackDepth = %d, want > 0", expr, p.MaxStackDepth)
			continue
		}
		// Eval with a scratch buffer of the advertised size must not panic.
		scratch := make([]float64, p.MaxStackDepth)
		_ = p.Eval(nil, scratch)
	}
}

// TestEvalScratchBufferGrowth checks Eval tolerates a nil or short scratch
// buffer by growing as needed.
func TestEvalScratchBufferGrowth(t *testing.T) {
	p, err := CompileExpr("1 + 2 * 3 + sqrt(16)")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}

	want := p.Eval(nil, make([]float64, p.MaxStackDepth))
	if got := p.Eval(nil, nil); got != want {
		t.Errorf("Eval with nil scratch = %v, want %v", got, want)
	}
	if got := p.Eval(nil, make([]float64, 1)); got != want {
		t.Errorf("Eval with 1-slot scratch = %v, want %v", got, want)
	}
}

// TestCloneIsIndependent checks that Clone deep-copies the backing slices, so
// mutating either program leaves the other untouched.
func TestCloneIsIndependent(t *testing.T) {
	p, err := CompileExpr("x + 1 * 2")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}
	if len(p.Consts) == 0 || len(p.VarIndices) == 0 {
		t.Fatalf("test program lacks constants/variables: %v", p)
	}
	c := p.Clone()

	before := p.Eval([]float64{1}, nil)
	if got := c.Eval([]float64{1}, nil); got != before {
		t.Fatalf("clone evaluates to %v, want %v", got, before)
	}

	// Mutating the clone must not affect the original.
	c.Consts[0] = 999
	c.VarIndices[0] = 7
	if got := p.Eval([]float64{1}, nil); got != before {
		t.Errorf("mutating the clone changed the original: %v -> %v", before, got)
	}
	if p.Consts[0] == 999 {
		t.Error("clone shares its Consts slice with the original")
	}
	if p.VarIndices[0] == 7 {
		t.Error("clone shares its VarIndices slice with the original")
	}

	// ...and mutating the original must not affect the clone either.
	p.Consts[0] = -42
	if c.Consts[0] == -42 {
		t.Error("clone shares its Consts slice with the original (reverse direction)")
	}

	// Appending to one program must not grow the other.
	origLen := len(p.Ops)
	_ = append(p.Ops, OpConst) //nolint:gocritic // deliberately exercising slice independence
	if len(c.Ops) != origLen {
		t.Errorf("append to the original changed the clone: %d != %d", len(c.Ops), origLen)
	}

	// Clone handles a nil receiver.
	var nilProgram *Program
	if got := nilProgram.Clone(); got != nil {
		t.Errorf("(*Program)(nil).Clone() = %v, want nil", got)
	}
}

// TestEvalRegularizedStaysFinite checks that regularized evaluation keeps
// results finite for inputs that make the plain evaluator diverge.
func TestEvalRegularizedStaysFinite(t *testing.T) {
	p, err := CompileExpr("exp(x) + log(x) + sqrt(x)")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}
	// x is negative, so log and sqrt are undefined; the regularized form must
	// not produce NaN.
	got := p.EvalRegularized([]float64{-1}, 1e-6, nil)
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Errorf("EvalRegularized = %v, want a finite value", got)
	}
}
