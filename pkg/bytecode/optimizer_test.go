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
