package bytecode

import (
	"math"
	"math/rand"
	"testing"
)

func TestOptimizerSubtractionAndDivisionNotCorrupted(t *testing.T) {
	// x - 2
	p := NewProgram()
	p.Ops = []OpCode{OpVar, OpConst, OpSub}
	p.Consts = []float64{2.0}
	p.VarIndices = []uint16{0}
	p.CalculateMaxStackDepth()

	opt := Optimize(p)
	valOrig := p.Eval([]float64{5.0}, nil)
	valOpt := opt.Eval([]float64{5.0}, nil)
	if valOrig != 3.0 || valOpt != 3.0 {
		t.Fatalf("x - 2 failed: orig=%v, opt=%v, want 3.0", valOrig, valOpt)
	}

	// x / 2
	pDiv := NewProgram()
	pDiv.Ops = []OpCode{OpVar, OpConst, OpDiv}
	pDiv.Consts = []float64{2.0}
	pDiv.VarIndices = []uint16{0}
	pDiv.CalculateMaxStackDepth()

	optDiv := Optimize(pDiv)
	valOrigDiv := pDiv.Eval([]float64{10.0}, nil)
	valOptDiv := optDiv.Eval([]float64{10.0}, nil)
	if valOrigDiv != 5.0 || valOptDiv != 5.0 {
		t.Fatalf("x / 2 failed: orig=%v, opt=%v, want 5.0", valOrigDiv, valOptDiv)
	}

	// (x + 1) * (x + 2)
	pComplex := NewProgram()
	pComplex.Ops = []OpCode{OpVar, OpConst, OpAdd, OpVar, OpConst, OpAdd, OpMul}
	pComplex.Consts = []float64{1.0, 2.0}
	pComplex.VarIndices = []uint16{0, 0}
	pComplex.CalculateMaxStackDepth()

	optComplex := Optimize(pComplex)
	valOrigComp := pComplex.Eval([]float64{3.0}, nil)
	valOptComp := optComplex.Eval([]float64{3.0}, nil)
	if valOrigComp != 20.0 || valOptComp != 20.0 {
		t.Fatalf("(x + 1) * (x + 2) at x=3: orig=%v, opt=%v, want 20.0", valOrigComp, valOptComp)
	}
}

func TestOptimizerAlgebraicIdentities(t *testing.T) {
	tests := []struct {
		name string
		expr string
		vars []float64
		want float64
	}{
		{"x + 0", "x + 0", []float64{7.5}, 7.5},
		{"0 + x", "0 + x", []float64{7.5}, 7.5},
		{"x - 0", "x - 0", []float64{7.5}, 7.5},
		{"0 - x", "0 - x", []float64{7.5}, -7.5},
		{"x * 1", "x * 1", []float64{7.5}, 7.5},
		{"1 * x", "1 * x", []float64{7.5}, 7.5},
		{"x * 0", "x * 0", []float64{7.5}, 0.0},
		{"0 * x", "0 * x", []float64{7.5}, 0.0},
		{"x / 1", "x / 1", []float64{7.5}, 7.5},
		{"x ^ 0", "x ^ 0", []float64{7.5}, 1.0},
		{"x ^ 1", "x ^ 1", []float64{7.5}, 7.5},
		{"-(-x)", "-(-x)", []float64{7.5}, 7.5},
		{"eml(x, 1)", "eml(x, 1)", []float64{2.0}, math.Exp(2.0)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := CompileExpr(tc.expr)
			if err != nil {
				t.Fatalf("compile error: %v", err)
			}
			opt := Optimize(p)
			got := opt.Eval(tc.vars, nil)
			if math.Abs(got-tc.want) > 1e-7 {
				t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestOptimizerDifferentialRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	vars := []float64{2.5}

	for i := 0; i < 2000; i++ {
		p := RandomProgram(1, 3, 15, rng)
		opt := Optimize(p)

		origVal := p.Eval(vars, nil)
		optVal := opt.Eval(vars, nil)

		if math.IsNaN(origVal) && math.IsNaN(optVal) {
			continue
		}
		if math.IsInf(origVal, 0) && math.IsInf(optVal, 0) {
			if math.IsInf(origVal, 1) == math.IsInf(optVal, 1) {
				continue
			}
		}

		diff := math.Abs(origVal - optVal)
		maxVal := math.Max(math.Abs(origVal), math.Abs(optVal))
		relErr := 0.0
		if maxVal > 1e-9 {
			relErr = diff / maxVal
		}

		if diff > 1e-5 && relErr > 1e-5 {
			t.Fatalf("Iteration %d: orig=%v, opt=%v, diff=%v, p=%s, opt=%s",
				i, origVal, optVal, diff, p, opt)
		}
	}
}

func FuzzOptimizerDifferential(f *testing.F) {
	f.Add(float64(5.0), float64(2.0), uint8(0))
	f.Add(float64(-3.5), float64(1.5), uint8(1))
	f.Add(float64(10.0), float64(0.0), uint8(2))

	f.Fuzz(func(t *testing.T, x, c float64, opChoice uint8) {
		if math.IsNaN(x) || math.IsNaN(c) || math.IsInf(x, 0) || math.IsInf(c, 0) {
			return
		}
		var op OpCode
		switch opChoice % 6 {
		case 0:
			op = OpAdd
		case 1:
			op = OpSub
		case 2:
			op = OpMul
		case 3:
			op = OpDiv
		case 4:
			op = OpPow
		case 5:
			op = OpEML
		}

		p := NewProgram()
		p.Ops = []OpCode{OpVar, OpConst, op}
		p.Consts = []float64{c}
		p.VarIndices = []uint16{0}
		p.CalculateMaxStackDepth()

		opt := Optimize(p)
		vars := []float64{x}
		valOrig := p.Eval(vars, nil)
		valOpt := opt.Eval(vars, nil)

		if math.IsNaN(valOrig) && math.IsNaN(valOpt) {
			return
		}
		if math.IsInf(valOrig, 0) && math.IsInf(valOpt, 0) {
			return
		}

		diff := math.Abs(valOrig - valOpt)
		denom := math.Max(math.Abs(valOrig), math.Abs(valOpt))
		if denom > 1e-9 {
			diff /= denom
		}
		if diff > 1e-5 {
			t.Fatalf("x=%v, c=%v, op=%v: orig=%v, opt=%v", x, c, op, valOrig, valOpt)
		}
	})
}

func FuzzEvalBatchColumnarDifferential(f *testing.F) {
	f.Add(int64(42), float64(0.5), float64(1.25))
	f.Add(int64(99), float64(-0.75), float64(2.5))
	f.Fuzz(func(t *testing.T, seed int64, x0, x1 float64) {
		if math.IsNaN(x0) || math.IsNaN(x1) || math.IsInf(x0, 0) || math.IsInf(x1, 0) {
			return
		}
		if math.Abs(x0) > 10 || math.Abs(x1) > 10 {
			return
		}
		rng := rand.New(rand.NewSource(seed))
		p := RandomProgram(2, 3, 11, rng)
		colData := [][]float64{{x0, x1}, {x1, x0}}
		rowData := [][]float64{{x0, x1}, {x1, x0}}
		colDst := make([]float64, 2)
		rowDst := make([]float64, 2)
		scratch := NewBatchScratch(p, 8)
		p.EvalBatchColumnarScratch(colData, colDst, &scratch)
		p.EvalBatch(rowData, rowDst)
		for i := range colDst {
			if math.IsNaN(colDst[i]) && math.IsNaN(rowDst[i]) {
				continue
			}
			if math.IsInf(colDst[i], 0) && math.IsInf(rowDst[i], 0) && math.Signbit(colDst[i]) == math.Signbit(rowDst[i]) {
				continue
			}
			scale := math.Max(1.0, math.Abs(rowDst[i]))
			if math.Abs(colDst[i]-rowDst[i]) > 1e-6*scale {
				t.Fatalf("EvalBatchColumnarScratch[%d]=%v != EvalBatch[%d]=%v for %s", i, colDst[i], i, rowDst[i], p)
			}
		}
	})
}

func TestOracleBytecodeVMAndOptimizer(t *testing.T) {
	p := NewProgram()
	p.Ops = []OpCode{OpVar, OpConst, OpEML, OpConst, OpAdd}
	p.Consts = []float64{1.0, 0.0}
	p.VarIndices = []uint16{0}
	p.CalculateMaxStackDepth()

	opt := Optimize(p)
	xs := []float64{-1.0, 0.0, 0.5, 1.0, 2.0}
	colOut := make([]float64, len(xs))
	opt.EvalBatchColumnar([][]float64{xs}, colOut)
	for i, x := range xs {
		want := p.Eval([]float64{x}, nil)
		if math.Abs(colOut[i]-want) > 1e-6*math.Max(1.0, math.Abs(want)) {
			t.Fatalf("TestOracleBytecodeVMAndOptimizer[%d]: got %v, want %v", i, colOut[i], want)
		}
	}
}
