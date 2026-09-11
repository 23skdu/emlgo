//go:build !js || !wasm

package jit

import (
	"testing"
)

func TestCacheAll(t *testing.T) {
	ClearJITCache()

	fn1, err := CompileCached("x + 1")
	if err != nil {
		t.Fatalf("CompileCached failed: %v", err)
	}
	if fn1(5.0) != 6.0 {
		t.Errorf("fn1(5) = %v, want 6", fn1(5.0))
	}

	// Cache hit
	fn2, err := CompileCached("x + 1")
	if err != nil {
		t.Fatalf("CompileCached second call failed: %v", err)
	}
	if fn2(5.0) != 6.0 {
		t.Errorf("fn2(5) = %v, want 6", fn2(5.0))
	}

	// Cache error
	_, err = CompileCached("(")
	if err == nil {
		t.Errorf("CompileCached(\"(\") should error")
	}

	// LRU eviction
	c := newJITCache(2)
	c.put("a", fn1)
	c.put("b", fn2)
	c.put("a", fn1) // move to front
	c.put("c", fn1) // evicts b
	if _, ok := c.get("b"); ok {
		t.Errorf("b should have been evicted")
	}
	if _, ok := c.get("a"); !ok {
		t.Errorf("a should still be in cache")
	}
	c.clear()
	if _, ok := c.get("a"); ok {
		t.Errorf("cache should be empty after clear")
	}
}

func TestDecompileAll(t *testing.T) {
	if Decompile(nil) != "" {
		t.Errorf("Decompile(nil) should be empty")
	}
	if DecompileLaTeX(nil) != "" {
		t.Errorf("DecompileLaTeX(nil) should be empty")
	}
	if DecompileNodeToExpr(nil) != "" {
		t.Errorf("DecompileNodeToExpr(nil) should be empty")
	}

	funcs := []string{"exp", "log", "sin", "cos", "tan", "sqrt", "neg", "add", "sub", "mul", "div", "pow", "custom"}
	for _, fn := range funcs {
		node := &EMLNode{
			Kind:  EMLFunc,
			Name:  fn,
			Left:  &EMLNode{Kind: EMLConst, Value: 3.0},
			Right: &EMLNode{Kind: EMLConst, Value: 4.0},
		}
		_ = Decompile(node)
		_ = DecompileLaTeX(node)
		_ = DecompileNodeToExpr(node)
	}

	// Variables and EMLOp
	emlOp := &EMLNode{
		Kind:  EMLOp,
		Left:  &EMLNode{Kind: EMLVar},
		Right: &EMLNode{Kind: EMLConst, Value: 2.5},
	}
	_ = Decompile(emlOp)
	_ = DecompileLaTeX(emlOp)
	_ = DecompileNodeToExpr(emlOp)

	// Sub-expression with parentheses
	addNode := &EMLNode{
		Kind:  EMLFunc,
		Name:  "add",
		Left:  &EMLNode{Kind: EMLVar},
		Right: &EMLNode{Kind: EMLConst, Value: 1.0},
	}
	negAdd := &EMLNode{
		Kind: EMLFunc,
		Name: "neg",
		Left: addNode,
	}
	_ = Decompile(negAdd)
	_ = DecompileLaTeX(negAdd)

	powNode := &EMLNode{
		Kind:  EMLFunc,
		Name:  "pow",
		Left:  addNode,
		Right: &EMLNode{Kind: EMLConst, Value: 2.0},
	}
	_ = Decompile(powNode)
	_ = DecompileLaTeX(powNode)
}

func TestEvalVarsAllFunctions(t *testing.T) {
	funcs := []string{
		"sin", "cos", "exp", "log", "sqrt", "tan", "asin", "acos", "atan", "abs",
		"cbrt", "log2", "log10", "ceil", "floor", "trunc", "round",
		"sinh", "cosh", "tanh", "asinh", "acosh", "atanh", "erf", "gamma", "unknown",
	}
	for _, fn := range funcs {
		n := FunctionCall{Name: fn, Arg: Number{Value: 0.5}}
		_ = EvalVars(n, nil)
	}

	// Binary ops
	for _, op := range []rune{'+', '-', '*', '/', '^', '%'} {
		n := BinaryOp{Left: Number{Value: 2.0}, Op: op, Right: Number{Value: 3.0}}
		_ = EvalVars(n, nil)
	}
	// Variable lookup
	_ = EvalVars(Variable{Name: "y"}, map[string]float64{"y": 10.0})
	_ = EvalVars(Variable{Name: "unknown"}, nil)
	_ = EvalVars(UnaryOp{Op: '-', Operand: Number{Value: 5.0}}, nil)
}

func TestArenaCallMathFuncAll(t *testing.T) {
	arena := NewArena(1024)
	funcs := []string{
		"sin", "cos", "exp", "log", "sqrt", "tan", "asin", "acos", "atan", "abs",
		"cbrt", "log2", "log10", "ceil", "floor", "trunc", "unknown",
	}
	for _, fn := range funcs {
		node := arena.FromInterface(FunctionCall{Name: fn, Arg: Number{Value: 0.5}})
		_ = EvalArena(node, 0.5)
	}
}

func TestCanonicalComprehensive(t *testing.T) {
	x := Canonicalize(Variable{Name: "x"})
	_ = CanonicalSin(x)
	_ = CanonicalCos(x)
	_ = CanonicalSqrt(x)

	// EMLEvalRegularized branches
	funcs := []string{"sin", "cos", "exp", "log", "neg", "add", "sub", "mul", "div", "pow", "sqrt", "unknown"}
	for _, fn := range funcs {
		n := &EMLNode{
			Kind:  EMLFunc,
			Name:  fn,
			Left:  &EMLNode{Kind: EMLConst, Value: -2.0},
			Right: &EMLNode{Kind: EMLConst, Value: 3.0},
		}
		_ = EMLEvalRegularized(n, 1.0, 1e-6)
		_ = EMLEvalRegularized(n, 1.0, -1.0) // default eps
	}

	// EMLOp clamp branches in EMLEvalRegularized
	hugeLeft := &EMLNode{Kind: EMLOp, Left: constNode(800.0), Right: constNode(1.0)}
	_ = EMLEvalRegularized(hugeLeft, 1.0, 1e-6)
	tinyLeft := &EMLNode{Kind: EMLOp, Left: constNode(-800.0), Right: constNode(1.0)}
	_ = EMLEvalRegularized(tinyLeft, 1.0, 1e-6)

	// Diff on all functions
	for _, fn := range []string{"exp", "log", "sin", "cos", "sqrt", "neg", "add", "sub", "mul", "div", "pow", "unknown"} {
		n := &EMLNode{
			Kind:  EMLFunc,
			Name:  fn,
			Left:  x,
			Right: constNode(2.0),
		}
		d := Diff(n)
		_ = DiffEval(n, 1.5)
		_ = d
	}
	// General power diff
	powGen := &EMLNode{Kind: EMLFunc, Name: "pow", Left: x, Right: x}
	_ = Diff(powGen)
	// EMLOp diff
	emlOp := emlNode(x, x)
	_ = Diff(emlOp)
	_ = Diff(nil)
	_ = Diff(constNode(5.0))

	// Equiv checks
	if !Equiv(nil, nil) {
		t.Errorf("Equiv(nil, nil) should be true")
	}
	if Equiv(nil, x) || Equiv(x, nil) {
		t.Errorf("Equiv with one nil should be false")
	}
	if Equiv(constNode(1.0), constNode(2.0)) {
		t.Errorf("Equiv with different constants should be false")
	}
	if !Equiv(constNode(1.0), constNode(1.0)) {
		t.Errorf("Equiv with same constants should be true")
	}
	if Equiv(CanonicalSin(x), CanonicalCos(x)) {
		t.Errorf("Equiv with different func names should be false")
	}
	if Equiv(x, constNode(1.0)) {
		t.Errorf("Equiv with different kinds should be false")
	}
}
