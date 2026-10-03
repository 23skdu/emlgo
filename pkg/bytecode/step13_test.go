package bytecode

import (
	"math/rand"
	"strings"
	"testing"
)

func TestStep13FinalizeValidation(t *testing.T) {
	// Unresolvable opcode
	pUnknown := &Program{
		Ops:    []OpCode{OpCode(200)},
		Consts: nil,
	}
	// Must pass validateStackFn to reach opcode loop
	oldValidate := validateStackFn
	validateStackFn = func([]OpCode) bool { return true }
	defer func() { validateStackFn = oldValidate }()

	err := finalize(pUnknown)
	if err == nil || !strings.Contains(err.Error(), "unresolvable opcode") {
		t.Fatalf("expected unresolvable opcode error, got %v", err)
	}

	// Constant count mismatch
	pConstMismatch := &Program{
		Ops:    []OpCode{OpConst},
		Consts: []float64{1.0, 2.0}, // 2 consts, 1 OpConst
	}
	err = finalize(pConstMismatch)
	if err == nil || !strings.Contains(err.Error(), "constant count mismatch") {
		t.Fatalf("expected constant count mismatch error, got %v", err)
	}

	// Var index count mismatch
	pVarMismatch := &Program{
		Ops:        []OpCode{OpVar},
		VarIndices: []uint16{}, // 0 var indices, 1 OpVar
	}
	err = finalize(pVarMismatch)
	if err == nil || !strings.Contains(err.Error(), "var index count mismatch") {
		t.Fatalf("expected var index count mismatch error, got %v", err)
	}
}

func TestStep13EvalStackUnderflows(t *testing.T) {
	ops := []OpCode{
		OpEML, OpAdd, OpSub, OpMul, OpDiv, OpPow,
		OpNeg, OpInv, OpSqrt, OpExp, OpLog, OpSin,
	}

	for _, op := range ops {
		t.Run("Eval_"+op.String(), func(t *testing.T) {
			p := &Program{
				Ops:           []OpCode{op},
				MaxStackDepth: 2,
			}
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("expected panic on %s underflow", op.String())
				}
			}()
			p.Eval([]float64{1.0}, nil)
		})

		t.Run("EvalRegularized_"+op.String(), func(t *testing.T) {
			p := &Program{
				Ops:           []OpCode{op},
				MaxStackDepth: 2,
			}
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("expected panic on %s underflow in EvalRegularized", op.String())
				}
			}()
			p.EvalRegularized([]float64{1.0}, 1e-6, nil)
		})
	}
}

func TestStep13EvalBatchColumnarScratchSafety(t *testing.T) {
	// Unknown opcode
	t.Run("UnknownOpcode", func(t *testing.T) {
		p := &Program{
			Ops:           []OpCode{OpCode(222)},
			MaxStackDepth: 2,
		}
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("expected panic on unknown opcode")
			}
		}()
		s := NewBatchScratch(p, 4)
		p.EvalBatchColumnarScratch(nil, make([]float64, 4), &s)
	})

	// Stack underflows
	ops := []OpCode{
		OpEML, OpAdd, OpSub, OpMul, OpDiv, OpPow,
		OpNeg, OpInv, OpSqrt, OpExp, OpLog, OpSin,
	}
	for _, op := range ops {
		t.Run("Underflow_"+op.String(), func(t *testing.T) {
			p := &Program{
				Ops:           []OpCode{op},
				MaxStackDepth: 2,
			}
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("expected panic on %s underflow in columnar", op.String())
				}
			}()
			s := NewBatchScratch(p, 4)
			p.EvalBatchColumnarScratch(nil, make([]float64, 4), &s)
		})
	}

	// Short column padding and offset past column
	t.Run("ShortColumnPadding", func(t *testing.T) {
		// p evaluates x (Var 0)
		p := &Program{
			Ops:           []OpCode{OpVar},
			VarIndices:    []uint16{0},
			MaxStackDepth: 2,
		}
		// data has 3 elements, dst has 2048 elements -> elements 3..2047 must be padded with 0
		data := [][]float64{{10.0, 20.0, 30.0}}
		dst := make([]float64, 2048)
		s := NewBatchScratch(p, 0)
		p.EvalBatchColumnarScratch(data, dst, &s)

		if dst[0] != 10.0 || dst[1] != 20.0 || dst[2] != 30.0 {
			t.Fatalf("unexpected data: %v", dst[:3])
		}
		for i := 3; i < 2048; i++ {
			if dst[i] != 0 {
				t.Fatalf("expected 0 padding at %d, got %v", i, dst[i])
			}
		}
	})
}

func TestStep19GAScalingAndParsimony(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	p1, _ := CompileExpr("x + 1")
	p2, _ := CompileExpr("x * 2")

	// 1. CrossoverOne edge cases
	if _, err := CrossoverOne(nil, p2, rng); err == nil {
		t.Errorf("CrossoverOne(nil, p2) should fail")
	}
	if _, err := CrossoverOne(p1, nil, rng); err == nil {
		t.Errorf("CrossoverOne(p1, nil) should fail")
	}
	if _, err := CrossoverOne(p1, p2, nil); err == nil {
		t.Errorf("CrossoverOne with nil rng should fail")
	}
	if _, err := CrossoverOne(&Program{}, p2, rng); err == nil {
		t.Errorf("CrossoverOne with empty ops should fail")
	}

	// CrossoverOne fallback with single-op program (empty split points)
	pLeaf := &Program{Ops: []OpCode{OpVar}, VarIndices: []uint16{0}, MaxStackDepth: 1}
	childLeaf, err := CrossoverOne(pLeaf, p2, rng)
	if err != nil || len(childLeaf.Ops) != len(pLeaf.Ops) {
		t.Fatalf("CrossoverOne leaf failed: %v", err)
	}

	// CrossoverOne stack validation failure fallback
	oldValidate := validateStackFn
	validateStackFn = func([]OpCode) bool { return false }
	cInvalid, err := CrossoverOne(p1, p2, rng)
	validateStackFn = oldValidate
	if err != nil || len(cInvalid.Ops) != len(p1.Ops) {
		t.Fatalf("CrossoverOne invalid stack fallback failed: %v", err)
	}

	// Normal CrossoverOne
	cValid, err := CrossoverOne(p1, p2, rng)
	if err != nil || cValid == nil || !ValidateStack(cValid.Ops) {
		t.Fatalf("CrossoverOne valid failed: %v", err)
	}

	// 2. Parsimony in GA Search
	d := polyDataset(16)
	resParsimony, err := Search(Config{
		PopulationSize:   32,
		Generations:      10,
		Workers:          2,
		Elitism:          2,
		LocalSearchSteps: 5,
		Parsimony:        0.05,
		Seed:             100,
	}, d, LeastSquares(1e-6))
	if err != nil {
		t.Fatalf("Search with Parsimony: %v", err)
	}
	if resParsimony.Best == nil {
		t.Errorf("expected best program from Parsimony search")
	}

	// 3. Early termination on PlateauGenerations and AdaptiveMutation
	resPlateau, err := Search(Config{
		PopulationSize:     32,
		Generations:        100,
		Workers:            2,
		Elitism:            2,
		LocalSearchSteps:   0,
		PlateauGenerations: 5,
		AdaptiveMutation:   true,
		Seed:               123,
	}, d, LeastSquares(1e-6))
	if err != nil {
		t.Fatalf("Search with Plateau: %v", err)
	}
	if len(resPlateau.History) >= 100 {
		t.Errorf("expected early plateau termination, ran all %d gens", len(resPlateau.History))
	}

	// 4. Early termination on StopTolerance and Optimize
	resTolerance, err := Search(Config{
		PopulationSize:   64,
		Generations:      50,
		Workers:          2,
		Elitism:          4,
		LocalSearchSteps: 20,
		StopTolerance:    5.0, // broad tolerance to trigger immediately
		Optimize:         true,
		Seed:             777,
	}, d, LeastSquares(1e-6))
	if err != nil {
		t.Fatalf("Search with StopTolerance: %v", err)
	}
	if len(resTolerance.History) >= 50 {
		t.Errorf("expected early stop on tolerance, ran all %d gens", len(resTolerance.History))
	}

	// 5. Serial Optimize branch (Workers: 1)
	resSerialOpt, err := Search(Config{
		PopulationSize:   16,
		Generations:      1,
		Workers:          1,
		Elitism:          2,
		LocalSearchSteps: 1,
		Optimize:         true,
		Seed:             42,
	}, d, LeastSquares(1e-6))
	if err != nil || resSerialOpt.Best == nil {
		t.Fatalf("Search with serial Optimize: %v", err)
	}
}
