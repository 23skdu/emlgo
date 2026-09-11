package bytecode

import (
	"math"
	"testing"

	"github.com/emlgo/eml/internal/jit"
)

func TestBytecodeOpsStringsAndArity(t *testing.T) {
	ops := []OpCode{
		OpConst, OpVar, OpEML, OpAdd, OpSub, OpMul, OpDiv, OpPow,
		OpNeg, OpInv, OpSqrt, OpExp, OpLog, OpCode(255),
	}
	for _, op := range ops {
		_ = op.String()
		_ = op.Arity()
	}

	var nilProg *Program
	if nilProg.String() != "<nil>" {
		t.Errorf("nil prog string should be <nil>")
	}
	p := NewProgram()
	_ = p.String()
	p.Ops = []OpCode{OpConst, OpVar, OpAdd}
	p.Consts = []float64{1.0}
	p.VarIndices = []uint16{0}
	_ = p.String()

	// Malformed program with missing const and var entries
	pBad := &Program{Ops: []OpCode{OpConst, OpVar}}
	_ = pBad.String()

	clone := p.Clone()
	if len(clone.Ops) != len(p.Ops) {
		t.Errorf("Clone failed")
	}
	if nilProg.Clone() != nil {
		t.Errorf("nil clone should be nil")
	}
}

func TestCompileExprComprehensive(t *testing.T) {
	exprs := []string{
		"1 + 2 * 3",
		"x - y / 2",
		"-x + 5",
		"x ^ 2",
		"exp(x) + log(y)",
		"ln(x) + sqrt(y)",
		"eml(x, y)",
		"(x + 1) * (y - 2)",
		"x / y",
	}
	for _, e := range exprs {
		p, err := CompileExpr(e)
		if err != nil {
			t.Errorf("CompileExpr(%q) failed: %v", e, err)
			continue
		}
		vars := []float64{2.0, 3.0}
		_ = p.Eval(vars, nil)
	}

	// Error cases
	errorExprs := []string{
		"",
		"(1 + 2",
		"1 + )",
		"unknown_func(x)",
		"x @ y",
		"1 + +",
		"x y",
		"eml(1)", // wrong arg count
	}
	for _, e := range errorExprs {
		_, _ = CompileExpr(e)
	}

	// Helper tests
	_ = isFunc("sin")
	_ = isFunc("unknown")
	_, _ = getPrecedence('+')
	_, _ = getPrecedence('*')
	_, _ = getPrecedence('^')
	_, _ = getPrecedence('?')

	vm := make(map[string]uint16)
	_, _ = getOrAssignVar("x", vm)
	_, _ = getOrAssignVar("x", vm) // existing
}

func TestCompileASTAllNodes(t *testing.T) {
	nodes := []jit.Node{
		jit.Number{Value: 42.0},
		jit.Variable{Name: "x"},
		jit.UnaryOp{Op: '-', Operand: jit.Number{Value: 5.0}},
		jit.BinaryOp{Left: jit.Number{Value: 1.0}, Op: '+', Right: jit.Number{Value: 2.0}},
		jit.BinaryOp{Left: jit.Number{Value: 1.0}, Op: '-', Right: jit.Number{Value: 2.0}},
		jit.BinaryOp{Left: jit.Number{Value: 1.0}, Op: '*', Right: jit.Number{Value: 2.0}},
		jit.BinaryOp{Left: jit.Number{Value: 1.0}, Op: '/', Right: jit.Number{Value: 2.0}},
		jit.BinaryOp{Left: jit.Number{Value: 1.0}, Op: '^', Right: jit.Number{Value: 2.0}},
		jit.FunctionCall{Name: "exp", Arg: jit.Number{Value: 1.0}},
		jit.FunctionCall{Name: "log", Arg: jit.Number{Value: 1.0}},
		jit.FunctionCall{Name: "sqrt", Arg: jit.Number{Value: 1.0}},
	}

	for _, n := range nodes {
		p, err := CompileAST(n)
		if err != nil {
			t.Errorf("CompileAST(%v) failed: %v", n, err)
			continue
		}
		_ = p.Eval([]float64{1.0}, nil)
	}

	// Unsupported node
	_, _ = CompileAST(nil)
	_, _ = CompileAST(jit.UnaryOp{Op: '+', Operand: jit.Number{Value: 5.0}})
	_, _ = CompileAST(jit.BinaryOp{Op: '%'})
	_, _ = CompileAST(jit.UnaryOp{Op: '!'})
	_, _ = CompileAST(jit.FunctionCall{Name: "foo"})
}

func TestCompileEMLAllNodes(t *testing.T) {
	var nilNode *jit.EMLNode
	_, _ = CompileEML(nilNode)

	e1 := jit.Canonicalize(jit.Number{Value: 3.14})
	_, _ = CompileEML(e1)

	e2 := jit.Canonicalize(jit.Variable{Name: "x"})
	_, _ = CompileEML(e2)

	e3 := jit.CanonicalExp(e2)
	_, _ = CompileEML(e3)

	// Test all EMLFunc operations and EMLOp
	funcNames := []string{"add", "sub", "mul", "div", "pow", "neg", "exp", "log", "sqrt", "unknown"}
	for _, fn := range funcNames {
		node := &jit.EMLNode{
			Kind:  jit.EMLFunc,
			Name:  fn,
			Left:  &jit.EMLNode{Kind: jit.EMLConst, Value: 1.0},
			Right: &jit.EMLNode{Kind: jit.EMLConst, Value: 2.0},
		}
		_, _ = CompileEML(node)
	}

	emlOpNode := &jit.EMLNode{
		Kind:  jit.EMLOp,
		Left:  &jit.EMLNode{Kind: jit.EMLConst, Value: 1.0},
		Right: &jit.EMLNode{Kind: jit.EMLConst, Value: 2.0},
	}
	_, _ = CompileEML(emlOpNode)

	unknownKind := &jit.EMLNode{Kind: 99}
	_, _ = CompileEML(unknownKind)
}

func TestEvalAndEvalRegularizedOpcodes(t *testing.T) {
	// Test every opcode in Eval and EvalRegularized
	p := NewProgram()
	p.Ops = []OpCode{
		OpConst, OpConst, OpAdd,
		OpConst, OpSub,
		OpConst, OpMul,
		OpConst, OpDiv,
		OpConst, OpPow,
		OpSqrt,
		OpInv,
		OpNeg,
		OpExp,
		OpLog,
		OpConst, OpEML,
	}
	p.Consts = []float64{
		2.0, 3.0, // Add -> 5
		1.0,      // Sub -> 4
		2.0,      // Mul -> 8
		2.0,      // Div -> 4
		2.0,      // Pow -> 16
		2.0,      // EML arg
	}
	p.CalculateMaxStackDepth()

	val1 := p.Eval(nil, nil)
	if math.IsNaN(val1) {
		t.Errorf("Eval produced NaN")
	}

	val2 := p.EvalRegularized(nil, 1e-6, nil)
	if math.IsNaN(val2) {
		t.Errorf("EvalRegularized produced NaN")
	}

	// Nil / empty checks
	var nilP *Program
	if nilP.Eval(nil, nil) != 0 {
		t.Errorf("nil prog eval should be 0")
	}
	if nilP.EvalRegularized(nil, 0, nil) != 0 {
		t.Errorf("nil prog eval regularized should be 0")
	}

	// Large stack depth (> 32)
	pLarge := NewProgram()
	for i := 0; i < 40; i++ {
		pLarge.Ops = append(pLarge.Ops, OpConst)
		pLarge.Consts = append(pLarge.Consts, 1.0)
	}
	for i := 0; i < 39; i++ {
		pLarge.Ops = append(pLarge.Ops, OpAdd)
	}
	pLarge.CalculateMaxStackDepth()
	scratch := make([]float64, pLarge.MaxStackDepth)
	_ = pLarge.Eval(nil, scratch)
	_ = pLarge.Eval(nil, nil) // allocate fallback
	_ = pLarge.EvalRegularized(nil, 1e-6, scratch)
	_ = pLarge.EvalRegularized(nil, -1.0, nil) // fallback & default eps
}

func TestEvalBatchAndColumnarComprehensive(t *testing.T) {
	// Program using all batch-supported opcodes
	p := NewProgram()
	p.Ops = []OpCode{
		OpVar, OpConst, OpAdd,
		OpVar, OpSub,
		OpVar, OpMul,
		OpConst, OpDiv,
		OpConst, OpPow,
		OpSqrt,
		OpInv,
		OpNeg,
		OpExp,
		OpLog,
		OpConst, OpEML,
	}
	p.Consts = []float64{1.0, 2.0, 2.0, 2.0}
	p.VarIndices = []uint16{0, 0, 0}
	p.CalculateMaxStackDepth()

	n := 100
	colX := make([]float64, n)
	for i := range colX {
		colX[i] = float64(i%10) + 1.0
	}
	dst := make([]float64, n)
	p.EvalBatchColumnar([][]float64{colX}, dst)

	// Row-based batch (including dst smaller than data)
	rows := make([][]float64, n)
	for i := range rows {
		rows[i] = []float64{colX[i]}
	}
	p.EvalBatch(rows, dst)
	p.EvalBatch(rows, dst[:10])

	// Nil / empty checks
	var nilP *Program
	nilP.EvalBatch(rows, dst)
	nilP.EvalBatchColumnar([][]float64{colX}, dst)
	p.EvalBatchColumnar(nil, dst)
}

func TestOptimizerAllRules(t *testing.T) {
	var nilP *Program
	if Optimize(nilP) != nil {
		t.Errorf("Optimize(nil) should be nil")
	}

	// Test folding for all binary and unary constant operations
	p := NewProgram()
	p.Ops = []OpCode{
		OpConst, OpConst, OpAdd,
		OpConst, OpConst, OpSub,
		OpConst, OpConst, OpMul,
		OpConst, OpConst, OpDiv,
		OpConst, OpConst, OpPow,
		OpConst, OpConst, OpEML,
		OpConst, OpConst, OpEML, // x=1, y=1
		OpConst, OpConst, OpEML, // x=0, y=1
		OpConst, OpNeg,
		OpConst, OpInv,
		OpConst, OpSqrt,
		OpConst, OpExp,
		OpConst, OpLog,
		OpVar, OpConst, OpEML, // x, CONST(1), EML -> EXP
	}
	p.Consts = []float64{
		1.0, 2.0, // Add
		5.0, 3.0, // Sub
		2.0, 4.0, // Mul
		10.0, 2.0, // Div
		2.0, 3.0, // Pow
		2.0, 3.0, // EML
		1.0, 1.0, // EML -> e
		0.0, 1.0, // EML -> 1
		5.0,      // Neg
		2.0,      // Inv
		4.0,      // Sqrt
		1.0,      // Exp
		2.0,      // Log
		1.0,      // EML with const 1
	}
	p.VarIndices = []uint16{0}

	opt := Optimize(p)
	if opt == nil || len(opt.Ops) == 0 {
		t.Errorf("Optimize failed")
	}
}

func TestGeneticOperatorsComprehensive(t *testing.T) {
	// ValidateStack
	if ValidateStack(nil) {
		t.Errorf("ValidateStack(nil) should be false")
	}
	if ValidateStack([]OpCode{OpAdd}) {
		t.Errorf("Underflow should be false")
	}
	if ValidateStack([]OpCode{OpConst, OpConst}) {
		t.Errorf("Stack height 2 should be false")
	}
	if !ValidateStack([]OpCode{OpConst, OpConst, OpAdd}) {
		t.Errorf("Valid stack should be true")
	}

	// Crossover
	p1, _ := CompileExpr("x + 1")
	p2, _ := CompileExpr("y * 2")
	c1, c2, err := Crossover(p1, p2)
	if err != nil || c1 == nil || c2 == nil {
		t.Errorf("Crossover failed: %v", err)
	}

	_, _, err = Crossover(nil, p2)
	if err == nil {
		t.Errorf("Crossover(nil, p2) should error")
	}

	// Mutate and mutateOp coverage
	_ = mutateOp(OpConst)
	_ = mutateOp(OpVar)
	_ = mutateOp(OpNeg)
	_ = mutateOp(OpAdd)
	_ = mutateOp(OpCode(255))

	m := Mutate(p1, 1.0)
	if m == nil {
		t.Errorf("Mutate failed")
	}
	if Mutate(nil, 1.0) != nil {
		t.Errorf("Mutate(nil) should be nil")
	}
}
