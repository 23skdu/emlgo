package bytecode

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/emlgo/eml/internal/jit"
)

func TestFunctabPanics(t *testing.T) {
	defer func() {
		initFuncTables(defaultUnaryCodes, unaryFuncs)
	}()

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("expected panic for unknown function")
			}
		}()
		initFuncTables(map[string]OpCode{"unknown_func": OpSin}, unaryFuncs)
	}()

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("expected panic for drifted tables")
			}
		}()
		initFuncTables(map[string]OpCode{"sin": OpSin}, unaryFuncs)
	}()
}

func TestOptimizerEdgeCases(t *testing.T) {
	if got := Optimize(nil); got != nil {
		t.Fatalf("Optimize(nil) = %v, want nil", got)
	}
	if got := Optimize(&Program{}); got == nil || len(got.Ops) != 0 {
		t.Fatalf("Optimize(empty) failed")
	}

	// Const underflow
	pConst := &Program{Ops: []OpCode{OpConst}}
	if got := Optimize(pConst); got == nil || len(got.Ops) != 1 {
		t.Fatalf("Optimize missing const failed")
	}

	// Var underflow
	pVar := &Program{Ops: []OpCode{OpVar}}
	if got := Optimize(pVar); got == nil || len(got.Ops) != 1 {
		t.Fatalf("Optimize missing var failed")
	}

	// Unary stack underflow
	for _, op := range []OpCode{OpNeg, OpInv, OpExp, OpLog, OpSin} {
		p := &Program{Ops: []OpCode{op}}
		if got := Optimize(p); got == nil || len(got.Ops) != 1 {
			t.Fatalf("Optimize stack underflow for %v failed", op)
		}
	}

	// Binary stack underflow
	for _, op := range []OpCode{OpAdd, OpSub, OpMul, OpDiv, OpPow, OpEML} {
		p := &Program{Ops: []OpCode{OpConst, op}, Consts: []float64{1.0}}
		if got := Optimize(p); got == nil || len(got.Ops) != 2 {
			t.Fatalf("Optimize stack underflow for %v failed", op)
		}
	}

	// Non-const child unary
	pExp := &Program{Ops: []OpCode{OpVar, OpExp}, VarIndices: []uint16{0}}
	if got := Optimize(pExp); got == nil || len(got.Ops) != 2 {
		t.Fatalf("Optimize exp(var) failed")
	}

	pLog := &Program{Ops: []OpCode{OpVar, OpLog}, VarIndices: []uint16{0}}
	if got := Optimize(pLog); got == nil || len(got.Ops) != 2 {
		t.Fatalf("Optimize log(var) failed")
	}

	pSin := &Program{Ops: []OpCode{OpVar, OpSin}, VarIndices: []uint16{0}}
	if got := Optimize(pSin); got == nil || len(got.Ops) != 2 {
		t.Fatalf("Optimize sin(var) failed")
	}

	// OpDiv left is 0: 0 / x -> 0
	pDiv0 := &Program{Ops: []OpCode{OpConst, OpVar, OpDiv}, Consts: []float64{0.0}, VarIndices: []uint16{0}}
	optDiv0 := Optimize(pDiv0)
	if len(optDiv0.Ops) != 1 || optDiv0.Ops[0] != OpConst || optDiv0.Consts[0] != 0.0 {
		t.Fatalf("Optimize 0 / x failed: %+v", optDiv0)
	}

	// Unknown opcode
	pUnknown := &Program{Ops: []OpCode{OpCode(255)}}
	if got := Optimize(pUnknown); got == nil || len(got.Ops) != 1 {
		t.Fatalf("Optimize unknown op failed")
	}

	// Unused consts / var indices
	pUnused := &Program{Ops: []OpCode{OpConst}, Consts: []float64{1.0, 2.0}}
	if got := Optimize(pUnused); got == nil || len(got.Consts) != 2 {
		t.Fatalf("Optimize unused consts failed")
	}
}

func TestGeneticEdgeCases(t *testing.T) {
	pools := buildMutationPools()
	if len(pools) == 0 {
		t.Fatal("empty mutation pools")
	}

	s := []int{1, 2}
	if got := dropAt(s, 5); len(got) != 2 {
		t.Fatalf("dropAt out of range failed: %v", got)
	}

	if got := splitPoints(nil); got != nil {
		t.Fatalf("splitPoints(nil) = %v, want nil", got)
	}

	rng := rand.New(rand.NewSource(42))
	p1 := &Program{Ops: []OpCode{OpConst}, Consts: []float64{1.0}}
	p2 := &Program{Ops: []OpCode{OpConst}, Consts: []float64{2.0}}
	c1, c2, err := Crossover(p1, p2, rng)
	if err != nil || c1 == nil || c2 == nil {
		t.Fatalf("Crossover single op failed: %v", err)
	}

	progA, err := CompileExpr("x + 1")
	if err != nil {
		t.Fatal(err)
	}
	progB, err := CompileExpr("x + 2")
	if err != nil {
		t.Fatal(err)
	}

	oldValidate := validateStackFn
	validateStackFn = func([]OpCode) bool { return false }
	c1, c2, err = Crossover(progA, progB, rng)
	validateStackFn = oldValidate
	if err != nil || c1 == nil || c2 == nil {
		t.Fatalf("Crossover with invalid stack check failed: %v", err)
	}

	rProg1 := RandomProgram(0, 0, 0, nil)
	if rProg1 == nil {
		t.Fatal("RandomProgram with zeros failed")
	}
	rProg2 := RandomProgram(1, 10, 5, rng)
	if rProg2 == nil {
		t.Fatal("RandomProgram min > max failed")
	}
}

func TestGAEdgeCases(t *testing.T) {
	s := timeSeed()
	if s == 0 {
		t.Fatal("timeSeed returned 0")
	}

	dX := Dataset{X: [][]float64{{1, 2}}}
	if dX.Samples() != 2 {
		t.Fatalf("Samples with X only = %d, want 2", dX.Samples())
	}

	dY := Dataset{X: nil, Y: []float64{1, 2}}
	if dY.Samples() != 2 {
		t.Fatalf("Samples with Y only = %d, want 2", dY.Samples())
	}

	dMismatch := Dataset{X: [][]float64{{1, 2}}, Y: []float64{1}}
	if err := dMismatch.validate(); err == nil {
		t.Fatal("expected validate error for mismatched X and Y")
	}

	fit := LeastSquares(-1)
	if fit == nil {
		t.Fatal("LeastSquares(-1) returned nil")
	}

	cfgStruct := Config{LocalSearchSteps: -1}
	cfg := cfgStruct.withDefaults()
	if cfg.LocalSearchSteps != 0 {
		t.Fatalf("withDefaults LocalSearchSteps = %d, want 0", cfg.LocalSearchSteps)
	}

	res := Result{BestFitness: 0.0}
	if !res.Converged(1e-4) {
		t.Fatal("expected Converged to be true")
	}

	// Search with Seed: 0 and X: nil
	resSearch, err := Search(Config{
		PopulationSize: 10,
		Generations:    2,
		Seed:           0,
		Workers:        1,
	}, Dataset{X: nil, Y: []float64{1, 2}}, LeastSquares(1e-6))
	if err != nil || resSearch.Best == nil {
		t.Fatalf("Search with Seed 0 failed: %v", err)
	}

	// scorePopulation branch lo >= len(pop)
	pop := make([]*Program, 97)
	for i := range pop {
		pop[i] = &Program{Ops: []OpCode{OpConst}, Consts: []float64{1.0}, MaxStackDepth: 1}
	}
	scores := make([]float64, 97)
	scorePopulation(Config{Workers: 12}, pop, Dataset{X: [][]float64{{1}}, Y: []float64{1}}, LeastSquares(1e-6), scores)
}

func TestCompilerEdgeCases(t *testing.T) {
	// Variable name empty
	pVar, err := CompileAST(jit.Variable{Name: ""})
	if err != nil || pVar == nil {
		t.Fatalf("CompileAST empty var name failed: %v", err)
	}

	// Max variables overflow
	varMap := make(map[string]uint16, math.MaxUint16)
	for i := 0; i < math.MaxUint16; i++ {
		varMap[fmt.Sprintf("v%d", i)] = uint16(i)
	}
	if _, err := getOrAssignVar("overflow", varMap); err == nil {
		t.Fatal("expected getOrAssignVar overflow error")
	}
	if err := emitAST(jit.Variable{Name: "overflow"}, NewProgram(), varMap); err == nil {
		t.Fatal("expected emitAST var overflow error")
	}
	if err := emitEML(&jit.EMLNode{Kind: jit.EMLVar, Name: "overflow"}, NewProgram(), varMap); err == nil {
		t.Fatal("expected emitEML var overflow error")
	}
	// Variable overflow in parseShuntingYard
	// Build an expression that jit.Parse rejects (using %) so it goes to parseShuntingYard
	var sb strings.Builder
	for i := 0; i < math.MaxUint16; i++ {
		sb.WriteString(fmt.Sprintf("v%d + ", i))
	}
	sb.WriteString("vOverflow")
	if _, err := parseShuntingYard(sb.String()); err == nil {
		t.Fatal("expected parseShuntingYard var overflow error")
	}

	// BinaryOp Right error
	if _, err := CompileAST(jit.BinaryOp{Left: jit.Number{Value: 1}, Right: jit.UnaryOp{Op: '+'}}); err == nil {
		t.Fatal("expected BinaryOp Right error")
	}

	// BinaryOp unsupported op
	if _, err := CompileAST(jit.BinaryOp{Op: '%', Left: jit.Number{Value: 1}, Right: jit.Number{Value: 2}}); err == nil {
		t.Fatal("expected BinaryOp unsupported op error")
	}

	// finalize validateStackFn error
	oldValidate := validateStackFn
	validateStackFn = func([]OpCode) bool { return false }
	if _, err := CompileAST(jit.Number{Value: 1}); err == nil {
		t.Fatal("expected finalize stack check error")
	}
	validateStackFn = oldValidate

	// emitEML edge cases
	if err := emitEML(nil, NewProgram(), make(map[string]uint16)); err != nil {
		t.Fatalf("emitEML(nil) failed: %v", err)
	}
	if err := emitEML(&jit.EMLNode{Kind: jit.EMLVar, Name: ""}, NewProgram(), make(map[string]uint16)); err != nil {
		t.Fatalf("emitEML empty var name failed: %v", err)
	}
	badNode := &jit.EMLNode{Kind: jit.EMLNodeKind(255)}
	if err := emitEML(&jit.EMLNode{Kind: jit.EMLOp, Left: badNode}, NewProgram(), make(map[string]uint16)); err == nil {
		t.Fatal("expected emitEML Left error")
	}
	if err := emitEML(&jit.EMLNode{Kind: jit.EMLOp, Left: &jit.EMLNode{Kind: jit.EMLConst}, Right: badNode}, NewProgram(), make(map[string]uint16)); err == nil {
		t.Fatal("expected emitEML Right error")
	}
	if err := emitEML(&jit.EMLNode{Kind: jit.EMLFunc, Left: badNode}, NewProgram(), make(map[string]uint16)); err == nil {
		t.Fatal("expected emitEML func Left error")
	}
	if err := emitEML(&jit.EMLNode{Kind: jit.EMLFunc, Left: &jit.EMLNode{Kind: jit.EMLConst}, Right: badNode}, NewProgram(), make(map[string]uint16)); err == nil {
		t.Fatal("expected emitEML func Right error")
	}

	// parseShuntingYard syntax and error branches
	if _, err := CompileExpr("1.2.3"); err == nil {
		t.Fatal("expected ParseFloat error for 1.2.3")
	}
	if _, err := CompileExpr("eml(x + 1, 2)"); err != nil {
		t.Fatalf("CompileExpr eml(x + 1, 2) failed: %v", err)
	}
	if _, err := CompileExpr("eml(x % 1, 2)"); err == nil {
		t.Fatal("expected unsupported operator in eml args")
	}
	if _, err := CompileExpr("(x % 1)"); err == nil {
		t.Fatal("expected unsupported operator in parens")
	}

	// Delete unaryCodes["sin"] temporarily
	delete(unaryCodes, "sin")
	if _, err := parseShuntingYard("sin(x)"); err == nil {
		t.Fatal("expected emitFunc error for missing unaryCode")
	}
	if _, err := parseShuntingYard("x + sin"); err == nil {
		t.Fatal("expected emitFunc error for missing unaryCode in opStack")
	}
	unaryCodes["sin"] = OpSin

	if _, err := parseShuntingYard("x % y % z"); err == nil {
		t.Fatal("expected emitOp error in operator loop")
	}
	if _, err := parseShuntingYard("x % y + z"); err == nil {
		t.Fatal("expected error in x % y + z")
	}
	if _, err := parseShuntingYard("x %"); err == nil {
		t.Fatal("expected emitOp error at end")
	}
}

func TestEvalEdgeCases(t *testing.T) {
	scratch := make([]float64, 4)
	pMissingVar := &Program{Ops: []OpCode{OpVar}, VarIndices: []uint16{0}, MaxStackDepth: 1}
	if got := pMissingVar.Eval([]float64{}, scratch); got != 0 {
		t.Fatalf("Eval missing var = %f, want 0", got)
	}
	if got := pMissingVar.EvalRegularized([]float64{}, 1e-6, scratch); got != 0 {
		t.Fatalf("EvalRegularized missing var = %f, want 0", got)
	}

	pUnknown := &Program{Ops: []OpCode{OpCode(255)}, MaxStackDepth: 1}
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("expected Eval to panic on unknown opcode")
			}
		}()
		pUnknown.Eval([]float64{}, scratch)
	}()
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("expected EvalRegularized to panic on unknown opcode")
			}
		}()
		pUnknown.EvalRegularized([]float64{}, 1e-6, scratch)
	}()

	pBatch := &Program{Ops: []OpCode{OpConst}, Consts: []float64{1.0}, MaxStackDepth: 0}
	cols := [][]float64{{1, 2}}
	dst := make([]float64, 2)
	batchScratch := NewBatchScratch(pBatch, 0)
	pBatch.EvalBatchColumnarScratch(cols, dst, &batchScratch)
	if dst[0] != 1.0 {
		t.Fatalf("EvalBatchColumnarScratch depth 0 failed")
	}

	pDeep := &Program{Ops: []OpCode{OpConst}, Consts: []float64{1.0}, MaxStackDepth: 35}
	pDeep.EvalBatch([][]float64{{1}}, dst[:1])
	if dst[0] != 1.0 {
		t.Fatalf("EvalBatch depth > 32 failed")
	}

	pEmpty := &Program{Ops: []OpCode{}}
	pEmpty.CalculateMaxStackDepth()
	if pEmpty.MaxStackDepth != 1 {
		t.Fatalf("empty Ops stack depth = %d, want 1", pEmpty.MaxStackDepth)
	}
}
