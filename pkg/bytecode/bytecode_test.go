package bytecode

import (
	"math"
	"testing"

	"github.com/emlgo/eml/internal/jit"
)

func TestBytecodeBasicEval(t *testing.T) {
	// Program for x + 2.5: VAR(0), CONST(2.5), ADD
	p := NewProgram()
	p.Ops = []OpCode{OpVar, OpConst, OpAdd}
	p.Consts = []float64{2.5}
	p.VarIndices = []uint16{0}
	p.CalculateMaxStackDepth()

	vars := []float64{3.5}
	got := p.Eval(vars, nil)
	want := 6.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Eval(3.5 + 2.5): got %v, want %v", got, want)
	}
}

func TestZeroAllocations(t *testing.T) {
	p := NewProgram()
	p.Ops = []OpCode{OpVar, OpConst, OpEML, OpConst, OpAdd}
	p.Consts = []float64{1.0, 5.0}
	p.VarIndices = []uint16{0}
	p.CalculateMaxStackDepth()

	vars := []float64{2.0}
	var scratch [32]float64

	// Verify that Eval performs zero heap allocations
	allocs := testing.AllocsPerRun(1000, func() {
		_ = p.Eval(vars, scratch[:])
	})
	if allocs > 0 {
		t.Errorf("Expected 0 allocations during Eval, got %v", allocs)
	}
}

func TestCompileExpr(t *testing.T) {
	expr := "x + 5"
	p, err := CompileExpr(expr)
	if err != nil {
		t.Fatalf("CompileExpr(%q) failed: %v", expr, err)
	}
	vars := []float64{10.0}
	got := p.Eval(vars, nil)
	want := 15.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("CompileExpr(%q): got %v, want %v", expr, got, want)
	}
}

func TestCompileEML(t *testing.T) {
	// Canonical exp(x) is eml(x, 1)
	expTree := jit.CanonicalExp(jit.Canonicalize(jit.Variable{Name: "x"}))
	p, err := CompileEML(expTree)
	if err != nil {
		t.Fatalf("CompileEML failed: %v", err)
	}

	for _, x := range []float64{-1.0, 0.0, 1.0, 2.0} {
		got := p.Eval([]float64{x}, nil)
		want := math.Exp(x)
		if math.Abs(got-want)/want > 1e-5 {
			t.Errorf("Eval(CompileEML(exp(%v))): got %v, want %v", x, got, want)
		}
	}
}

func TestEvalBatchColumnar(t *testing.T) {
	// Program: x * y + 3.0
	p := NewProgram()
	p.Ops = []OpCode{OpVar, OpVar, OpMul, OpConst, OpAdd}
	p.Consts = []float64{3.0}
	p.VarIndices = []uint16{0, 1}
	p.CalculateMaxStackDepth()

	n := 2048
	colX := make([]float64, n)
	colY := make([]float64, n)
	dst := make([]float64, n)

	for i := 0; i < n; i++ {
		colX[i] = float64(i) * 0.1
		colY[i] = float64(i%10) + 1.0
	}

	data := [][]float64{colX, colY}
	p.EvalBatchColumnar(data, dst)

	for i := 0; i < n; i++ {
		scalar := p.Eval([]float64{colX[i], colY[i]}, nil)
		if math.Abs(dst[i]-scalar) > 1e-9 {
			t.Fatalf("Batch mismatch at index %d: batch %v vs scalar %v", i, dst[i], scalar)
		}
	}
}

func TestOptimize(t *testing.T) {
	// eml(1, 1) -> e
	p := NewProgram()
	p.Ops = []OpCode{OpConst, OpConst, OpEML}
	p.Consts = []float64{1.0, 1.0}
	p.CalculateMaxStackDepth()

	opt := Optimize(p)
	if len(opt.Ops) != 1 || opt.Ops[0] != OpConst {
		t.Errorf("Optimize failed to fold eml(1, 1): got %s", opt)
	}
	if math.Abs(opt.Consts[0]-math.E) > 1e-9 {
		t.Errorf("Folded constant = %v, want e (%v)", opt.Consts[0], math.E)
	}
}

func TestGeneticStackValidation(t *testing.T) {
	valid := []OpCode{OpVar, OpConst, OpAdd}
	if !ValidateStack(valid) {
		t.Errorf("Expected %v to be valid RPN", valid)
	}

	invalidUnderflow := []OpCode{OpAdd, OpVar, OpConst}
	if ValidateStack(invalidUnderflow) {
		t.Errorf("Expected %v to be invalid (underflow)", invalidUnderflow)
	}

	invalidTooManyOutputs := []OpCode{OpVar, OpConst}
	if ValidateStack(invalidTooManyOutputs) {
		t.Errorf("Expected %v to be invalid (leaves 2 items on stack)", invalidTooManyOutputs)
	}
}

func BenchmarkBytecodeScalarEval(b *testing.B) {
	p := NewProgram()
	p.Ops = []OpCode{OpVar, OpConst, OpEML, OpConst, OpAdd}
	p.Consts = []float64{1.0, 2.0}
	p.VarIndices = []uint16{0}
	p.CalculateMaxStackDepth()

	vars := []float64{1.5}
	var scratch [32]float64

	for b.Loop() {
		_ = p.Eval(vars, scratch[:])
	}
}

func BenchmarkBytecodeBatchColumnar(b *testing.B) {
	p := NewProgram()
	p.Ops = []OpCode{OpVar, OpConst, OpEML, OpConst, OpAdd}
	p.Consts = []float64{1.0, 2.0}
	p.VarIndices = []uint16{0}
	p.CalculateMaxStackDepth()

	n := 10000
	col := make([]float64, n)
	for i := range col {
		col[i] = float64(i) * 0.001
	}
	dst := make([]float64, n)
	data := [][]float64{col}

	for b.Loop() {
		p.EvalBatchColumnar(data, dst)
	}
}
