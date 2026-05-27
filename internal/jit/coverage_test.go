//go:build amd64

package jit

import (
	"math"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

// -- ast.go: nodeSigil --

func TestNodeSigil(t *testing.T) {
	(Number{}).nodeSigil()
	(Variable{}).nodeSigil()
	(UnaryOp{}).nodeSigil()
	(BinaryOp{}).nodeSigil()
	(FunctionCall{}).nodeSigil()
}

// -- codegen.go: sqrtsd, addPool dedup, REX paths, encoder methods --

func TestEncoderSqrtsd(t *testing.T) {
	var e encoder
	e.sqrtsd(0, 1)
	if len(e.code) != 4 {
		t.Fatalf("sqrtsd emitted %d bytes, want 4", len(e.code))
	}
}

func TestEncoderAddPoolDedup(t *testing.T) {
	var e encoder
	idx0 := e.addPool(3.14)
	idx1 := e.addPool(3.14)
	if idx0 != idx1 {
		t.Fatalf("addPool dedup failed: %d vs %d", idx0, idx1)
	}
}

func TestEncoderLoadConstantREX(t *testing.T) {
	var e encoder
	e.addPool(1.0)
	e.loadConstant(8, 0) // dst >= 8 triggers REX 0x44
	if len(e.code) < 9 {
		t.Fatalf("loadConstant REX emitted %d bytes, want >= 9", len(e.code))
	}
	if e.code[0] != 0x44 {
		t.Fatalf("expected REX.W prefix 0x44, got %02x", e.code[0])
	}
}

func TestEncoderMovsdStoreREX(t *testing.T) {
	var e encoder
	e.movsdStore(8) // reg >= 8 triggers REX 0x41
	if len(e.code) < 6 {
		t.Fatalf("movsdStore REX emitted %d bytes, want >= 6", len(e.code))
	}
	if e.code[0] != 0x41 {
		t.Fatalf("expected REX.B prefix 0x41, got %02x", e.code[0])
	}
}

func TestEncoderMovsdLoadREX(t *testing.T) {
	var e encoder
	e.movsdLoad(8) // reg >= 8 triggers REX 0x44
	if len(e.code) < 6 {
		t.Fatalf("movsdLoad REX emitted %d bytes, want >= 6", len(e.code))
	}
	if e.code[0] != 0x44 {
		t.Fatalf("expected REX.W prefix 0x44, got %02x", e.code[0])
	}
}

func TestEncoderHelperMethods(t *testing.T) {
	var e encoder
	e.emit(0x90)
	if len(e.code) != 1 || e.code[0] != 0x90 {
		t.Fatal("emit failed")
	}
	e.emit32(0x12345678)
	if len(e.code) != 5 {
		t.Fatal("emit32 failed")
	}
	e.emit64(0xDEADBEEF)
	if len(e.code) != 13 {
		t.Fatal("emit64 failed")
	}
	if e.rex(0, 0, 0, 0) != 0x40 {
		t.Fatal("rex empty failed")
	}
	if e.rex(1, 0, 0, 0) != 0x48 {
		t.Fatal("rex W failed")
	}
	if e.modrm(3, 0, 1) != 0xC1 {
		t.Fatalf("modrm got %02x", e.modrm(3, 0, 1))
	}
	if e.sib(0, 4, 4) != 0x24 {
		t.Fatalf("sib got %02x", e.sib(0, 4, 4))
	}
}

func TestEncoderSse2REX(t *testing.T) {
	var e encoder
	e.sse2(0xF2, 0x58, 8, 0) // dst >= 8
	if len(e.code) < 5 {
		t.Fatalf("sse2 REX emitted %d bytes, want >= 5", len(e.code))
	}
	if e.code[0] != 0x44 && e.code[0] != 0x41 {
		t.Fatalf("expected REX prefix, got %02x", e.code[0])
	}
}

func TestEncoderPushPop(t *testing.T) {
	var e encoder
	e.push()
	if len(e.code) != 9 {
		t.Fatalf("push emitted %d bytes, want 9", len(e.code))
	}
	e.popTo(0)
	if len(e.code) != 18 {
		t.Fatalf("push+pop emitted %d bytes, want 18", len(e.code))
	}
}

// -- codegen.go: gen FunctionCall error (operand in UnaryOp) --

func TestCodegenUnaryFunctionCall(t *testing.T) {
	n := UnaryOp{Op: '-', Operand: FunctionCall{Name: "sin", Arg: Variable{}}}
	_, err := compileToCode(n)
	if err == nil {
		t.Fatal("expected error for function call operand in unary op")
	}
}

// -- codegen.go: gen BinaryOp left/right error paths --

func TestCodegenBinaryLeftError(t *testing.T) {
	n := BinaryOp{
		Left:  FunctionCall{Name: "sin", Arg: Variable{}},
		Op:    '+',
		Right: Number{1},
	}
	_, err := compileToCode(n)
	if err == nil || !strings.Contains(err.Error(), "function calls not supported") {
		t.Fatalf("expected function call error, got: %v", err)
	}
}

func TestCodegenBinaryRightError(t *testing.T) {
	n := BinaryOp{
		Left:  Number{1},
		Op:    '+',
		Right: FunctionCall{Name: "cos", Arg: Number{0}},
	}
	_, err := compileToCode(n)
	if err == nil || !strings.Contains(err.Error(), "function calls not supported") {
		t.Fatalf("expected function call error, got: %v", err)
	}
}

// -- codegen.go: gen BinaryOp unsupported operator --

func TestCodegenUnsupportedOp(t *testing.T) {
	n := BinaryOp{
		Left:  Number{1},
		Op:    '%',
		Right: Number{2},
	}
	_, err := compileToCode(n)
	if err == nil || !strings.Contains(err.Error(), "unsupported operator") {
		t.Fatalf("expected unsupported operator error, got: %v", err)
	}
}

// -- codegen.go: genPow n==0 (x^0) --

func TestCodegenPowZero(t *testing.T) {
	f, err := NewCompiler().Compile("x^0")
	if err != nil {
		t.Fatal(err)
	}
	got := f(42)
	if got != 1 {
		t.Fatalf("x^0 at 42 = %v, want 1", got)
	}
}

// -- codegen.go: genPow gen(base) error (sin(x)^2) --

func TestCodegenPowBaseError(t *testing.T) {
	n := BinaryOp{
		Left:  FunctionCall{Name: "sin", Arg: Variable{}},
		Op:    '^',
		Right: Number{2},
	}
	_, err := compileToCode(n)
	if err == nil || !strings.Contains(err.Error(), "function calls not supported") {
		t.Fatalf("expected function call error, got: %v", err)
	}
}

// -- eval.go: asin, acos, atan, fallthrough return 0 --

func TestEvalAsinAcosAtan(t *testing.T) {
	tests := []struct {
		expr string
		want float64
	}{
		{"asin(0)", 0},
		{"acos(1)", 0},
		{"atan(0)", 0},
	}
	for _, tc := range tests {
		n, err := Parse(tc.expr)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.expr, err)
		}
		got := Eval(n, 0)
		if math.Abs(got-tc.want) > 1e-14 {
			t.Errorf("%s: Eval = %v, want %v", tc.expr, got, tc.want)
		}
	}
}

func TestEvalUnknownBinaryOpFallthrough(t *testing.T) {
	n := BinaryOp{Left: Number{1}, Op: '%', Right: Number{2}}
	got := Eval(n, 0)
	if got != 0 {
		t.Fatalf("Eval unknown op = %v, want 0", got)
	}
}

func TestEvalUnknownFunctionFallthrough(t *testing.T) {
	n := FunctionCall{Name: "unknown", Arg: Number{1}}
	got := Eval(n, 0)
	if got != 0 {
		t.Fatalf("Eval unknown function = %v, want 0", got)
	}
}

// -- jit.go: Compile error paths --

func TestCompileParseError(t *testing.T) {
	_, err := NewCompiler().Compile("2x")
	if err == nil || !strings.Contains(err.Error(), "parse error") {
		t.Fatalf("expected parse error, got: %v", err)
	}
}

func TestCompileCodegenError(t *testing.T) {
	_, err := NewCompiler().Compile("sin(x)")
	if err == nil || !strings.Contains(err.Error(), "codegen error") {
		t.Fatalf("expected codegen error, got: %v", err)
	}
}

// -- parser.go: lexer EOF fallthrough (unknown char) --

func TestParseUnknownChar(t *testing.T) {
	_, err := Parse("@")
	if err == nil {
		t.Fatal("expected error for unknown char")
	}
}

// -- parser.go: tokTypeFromRune default --

func TestTokTypeFromRuneDefault(t *testing.T) {
	got := tokTypeFromRune('@')
	if got != tokEOF {
		t.Fatalf("tokTypeFromRune('@') = %d, want tokEOF(%d)", got, tokEOF)
	}
}

// -- parser.go: right-side errors in parseMulDiv, parsePower, parseUnary, parseBase --

func TestParseMulDivRightError(t *testing.T) {
	_, err := Parse("1 / *")
	if err == nil {
		t.Fatal("expected error from right-side parseMulDiv failure")
	}
}

func TestParsePowerRightError(t *testing.T) {
	_, err := Parse("x ^ *")
	if err == nil {
		t.Fatal("expected error from right-side parsePower failure")
	}
}

func TestParseUnaryError(t *testing.T) {
	_, err := Parse("- *")
	if err == nil {
		t.Fatal("expected error from parseUnary operand failure")
	}
}

func TestParseInvalidNumber(t *testing.T) {
	_, err := Parse("12.34.56")
	if err == nil || !strings.Contains(err.Error(), "invalid number") {
		t.Fatalf("expected invalid number error, got: %v", err)
	}
}

func TestParseParenExprError(t *testing.T) {
	_, err := Parse("( + 1)")
	if err == nil {
		t.Fatal("expected error for empty parenthesized expression")
	}
}

func TestParseFunctionArgError(t *testing.T) {
	_, err := Parse("sin( + 1)")
	if err == nil {
		t.Fatal("expected error for invalid function argument")
	}
}

func TestParseFunctionMissingRParen(t *testing.T) {
	_, err := Parse("sin(1")
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected missing ) error, got: %v", err)
	}
}

// -- parser.go: FormatAST, FormatExpr, wrap helpers --

func TestFormatAST(t *testing.T) {
	n, err := Parse("x + 1")
	if err != nil {
		t.Fatal(err)
	}
	got := FormatAST(n)
	if got != "(x + 1)" {
		t.Fatalf("FormatAST = %q, want (x + 1)", got)
	}
}

func TestFormatExprAllTypes(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"42", "42"},
		{"x", "x"},
		{"-x", "-x"},
		{"x + x", "x + x"},
		{"x * 2", "x * 2"},
		{"x / 2", "x / 2"},
		{"x - 1", "x - 1"},
		{"x + 1", "x + 1"},
		{"x * 2", "x * 2"},
		{"x / 2", "x / 2"},
		{"x - 1", "x - 1"},
		{"x^2", "x^2"},
		{"sin(x)", "sin(x)"},
		{"x + x", "x + x"},
		{"x - x", "x - x"},
		{"x * x", "x * x"},
		{"x / x", "x / x"},
	}
	for _, tc := range tests {
		n, err := Parse(tc.expr)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.expr, err)
		}
		got := FormatExpr(n)
		if got != tc.want {
			t.Errorf("FormatExpr(%q) = %q, want %q", tc.expr, got, tc.want)
		}
	}
}

func TestFormatExprPrecedence(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"x + x * x", "x + x * x"},
		{"x * x + x", "x * x + x"},
		{"(x + 1) * 2", "(x + 1) * 2"},
		{"2 * (x + 1)", "2 * (x + 1)"},
		{"x - x - x", "x - x - x"},
		{"-(x + 1)", "-(x + 1)"},
		{"-x^2", "-(x^2)"},
	}
	for _, tc := range tests {
		n, err := Parse(tc.expr)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.expr, err)
		}
		got := FormatExpr(n)
		if got != tc.want {
			t.Errorf("FormatExpr(%q) = %q, want %q", tc.expr, got, tc.want)
		}
	}
}

func TestFormatExprUnknownType(t *testing.T) {
	got := FormatExpr(struct{ Node }{Number{42}})
	if got != "42" {
		t.Fatalf("FormatExpr unknown node = %q, want 42", got)
	}
}

func TestWrapParen(t *testing.T) {
	tests := []struct {
		name string
		n    Node
		want string
	}{
		{"binary", BinaryOp{Left: Number{1}, Op: '+', Right: Number{2}}, "(1 + 2)"},
		{"number", Number{42}, "42"},
		{"variable", Variable{}, "x"},
	}
	for _, tc := range tests {
		got := wrapParen(tc.n)
		if got != tc.want {
			t.Errorf("wrapParen(%s) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestWrapBinOp(t *testing.T) {
	n := BinaryOp{Left: Number{1}, Op: '+', Right: Number{2}}
	gotLeft := wrapBinOp(n.Left, '*', true)
	if gotLeft != "1" {
		t.Fatalf("wrapBinOp number = %q, want 1", gotLeft)
	}
	binLeft := BinaryOp{Left: Number{1}, Op: '+', Right: Number{2}}
	got := wrapBinOp(binLeft, '*', true)
	if got != "(1 + 2)" {
		t.Fatalf("wrapBinOp lower prec = %q, want (1 + 2)", got)
	}
	binPower := BinaryOp{Left: Number{2}, Op: '^', Right: Number{3}}
	got = wrapBinOp(binPower, '^', true)
	if got != "(2^3)" {
		t.Fatalf("wrapBinOp same prec power = %q, want (2^3)", got)
	}
	got = wrapBinOp(binPower, '+', true)
	if got != "2^3" {
		t.Fatalf("wrapBinOp higher prec = %q, want 2^3", got)
	}
	sub := BinaryOp{Left: Number{1}, Op: '-', Right: Number{2}}
	got = wrapBinOp(sub, '-', false)
	if got != "(1 - 2)" {
		t.Fatalf("wrapBinOp right-assoc minus = %q, want (1 - 2)", got)
	}
	got = wrapBinOp(sub, '-', true)
	if got != "1 - 2" {
		t.Fatalf("wrapBinOp left-assoc minus = %q, want 1 - 2", got)
	}
}

func TestWrapPower(t *testing.T) {
	tests := []struct {
		name string
		n    Node
		want string
	}{
		{"binary", BinaryOp{Left: Number{1}, Op: '+', Right: Number{2}}, "(1 + 2)"},
		{"unary", UnaryOp{Op: '-', Operand: Number{5}}, "(-5)"},
		{"number", Number{42}, "42"},
	}
	for _, tc := range tests {
		got := wrapPower(tc.n)
		if got != tc.want {
			t.Errorf("wrapPower(%s) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// -- UnaryOp.String, FunctionCall.String --

func TestUnaryOpString(t *testing.T) {
	n := UnaryOp{Op: '-', Operand: Number{5}}
	got := n.String()
	if got != "(-5)" {
		t.Fatalf("UnaryOp.String() = %q, want (-5)", got)
	}
}

func TestFunctionCallString(t *testing.T) {
	n := FunctionCall{Name: "sin", Arg: Number{0}}
	got := n.String()
	if got != "sin(0)" {
		t.Fatalf("FunctionCall.String() = %q, want sin(0)", got)
	}
}

// -- edge case: format expr with empty node --
type emptyNode struct{}

func (emptyNode) nodeSigil()         {}
func (emptyNode) String() string     { return "?" }
func (emptyNode) GoString() string   { return "?" }
func TestFormatExprEdgeCases(t *testing.T) {
	n := BinaryOp{Left: Variable{}, Op: '+', Right: Variable{}}
	got := FormatExpr(n)
	if got != "x + x" {
		t.Fatalf("FormatExpr(x + x) = %q, want x + x", got)
	}
}

func TestFormatExprFuncCall(t *testing.T) {
	n, err := Parse("sin(x)")
	if err != nil {
		t.Fatal(err)
	}
	got := FormatExpr(n)
	if got != "sin(x)" {
		t.Fatalf("FormatExpr(sin(x)) = %q, want sin(x)", got)
	}
}

func TestFormatASTNil(t *testing.T) {
	// The default case inside FormatExpr for unknown Node types falls back to n.String()
	n := struct{ Node }{Number{99}}
	got := FormatAST(n)
	if got != "99" {
		t.Fatalf("FormatAST with embedded = %q, want 99", got)
	}
}

// -- wrapBinOp: prec default (parser.go:312) --

func TestWrapBinOpPrecDefault(t *testing.T) {
	n := BinaryOp{Left: Number{1}, Op: '%', Right: Number{2}}
	got := wrapBinOp(n, '+', true)
	if got != "(1 % 2)" {
		t.Fatalf("wrapBinOp prec default = %q, want (1 %% 2)", got)
	}
}

// -- AllocateExecutableMemory error paths via RLIMIT_AS --

func TestCompileMemoryAllocationError(t *testing.T) {
	runtime.LockOSThread()
	var oldRlim unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_AS, &oldRlim); err != nil {
		runtime.UnlockOSThread()
		t.Skip("cannot get RLIMIT_AS:", err)
	}
	zero := unix.Rlimit{Cur: 0, Max: oldRlim.Max}
	if err := unix.Setrlimit(unix.RLIMIT_AS, &zero); err != nil {
		runtime.UnlockOSThread()
		t.Skip("cannot set RLIMIT_AS:", err)
	}
	_, err := NewCompiler().Compile("x")
	if err == nil || !strings.Contains(err.Error(), "memory allocation") {
		unix.Setrlimit(unix.RLIMIT_AS, &oldRlim)
		runtime.UnlockOSThread()
		t.Fatalf("expected memory allocation error, got: %v", err)
	}
	unix.Setrlimit(unix.RLIMIT_AS, &oldRlim)
	runtime.UnlockOSThread()
}

// Test AllocateExecutableMemory happy path
func TestAllocateExecMemory(t *testing.T) {
	code := []byte{0xC3} // RET
	ptr, err := AllocateExecutableMemory(code)
	if err != nil {
		t.Fatal(err)
	}
	_ = ptr
}
func TestFormatExprComplexPrecedence(t *testing.T) {
	n := BinaryOp{
		Left:  BinaryOp{Left: Number{1}, Op: '+', Right: Number{2}},
		Op:    '*',
		Right: Number{3},
	}
	got := FormatExpr(n)
	if got != "(1 + 2) * 3" {
		t.Fatalf("FormatExpr = %q, want (1 + 2) * 3", got)
	}

	n2 := BinaryOp{
		Left:  Number{1},
		Op:    '+',
		Right: BinaryOp{Left: Number{2}, Op: '*', Right: Number{3}},
	}
	got2 := FormatExpr(n2)
	if got2 != "1 + 2 * 3" {
		t.Fatalf("FormatExpr = %q, want 1 + 2 * 3", got2)
	}
}
