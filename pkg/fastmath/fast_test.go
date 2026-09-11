package fastmath

import (
	"math"
	"testing"
)

func TestFastExpMinimaxAccuracy(t *testing.T) {
	testPoints := []float64{
		-100.0, -10.0, -5.0, -1.0, -0.5, -0.1, -0.01, 0.0,
		0.01, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 50.0, 100.0, 500.0,
	}
	for _, x := range testPoints {
		got := FastExpMinimax(x)
		want := math.Exp(x)
		relErr := math.Abs(got-want) / want
		if relErr > 1.5e-6 {
			t.Errorf("FastExpMinimax(%v): got %v, want %v, relErr %v > 1.5e-6", x, got, want, relErr)
		}
	}
}

func TestFastLogMinimaxAccuracy(t *testing.T) {
	testPoints := []float64{
		1e-10, 1e-5, 0.01, 0.1, 0.5, 0.707, 1.0, 1.414, 2.0, math.E,
		5.0, 10.0, 100.0, 1000.0, 1e6, 1e15,
	}
	for _, y := range testPoints {
		got := FastLogMinimax(y)
		want := math.Log(y)
		relErr := math.Abs(got-want) / math.Abs(want)
		if math.Abs(want) < 1e-6 {
			relErr = math.Abs(got - want)
		}
		if relErr > 5e-5 {
			t.Errorf("FastLogMinimax(%v): got %v, want %v, relErr %v > 5e-5", y, got, want, relErr)
		}
	}
}

func TestFastExpBitCast(t *testing.T) {
	for _, x := range []float64{-2.0, -1.0, 0.0, 0.5, 1.0, 2.0} {
		got := FastExpBitCast(x)
		want := math.Exp(x)
		relErr := math.Abs(got-want) / want
		if relErr > 0.05 { // Schraudolph relative error is ~1-3%
			t.Errorf("FastExpBitCast(%v): got %v, want %v, relErr %v > 5%%", x, got, want, relErr)
		}
	}
}

func TestGuardrails(t *testing.T) {
	// LnRegularized at zero should be finite and equal to ln(eps)
	eps := 1e-6
	valAtZero := LnRegularized(0, eps)
	wantAtZero := math.Log(eps)
	if math.Abs(valAtZero-wantAtZero) > 1e-4 {
		t.Errorf("LnRegularized(0) = %v, want %v", valAtZero, wantAtZero)
	}

	// Negative values should be smooth and positive
	negVal := LnRegularized(-5.0, eps)
	posVal := LnRegularized(5.0, eps)
	if math.Abs(negVal-posVal) > 1e-6 {
		t.Errorf("LnRegularized should be symmetric: f(-5)=%v, f(5)=%v", negVal, posVal)
	}

	// Clamped Exp should not overflow to +Inf
	hugeExp := ExpClamped(1000.0)
	if math.IsInf(hugeExp, 1) || math.IsNaN(hugeExp) {
		t.Errorf("ExpClamped(1000) returned Inf or NaN: %v", hugeExp)
	}

	// FastEmlRegularized should never produce NaN across negative domain
	for _, y := range []float64{-100.0, -1.0, -0.001, 0.0, 0.001, 1.0, 100.0} {
		res := FastEmlRegularized(1.0, y, DefaultEpsilon)
		if math.IsNaN(res) || math.IsInf(res, 0) {
			t.Errorf("FastEmlRegularized(1.0, %v) produced non-finite: %v", y, res)
		}
	}
}

func BenchmarkFastExp(b *testing.B) {
	x := 1.234
	for b.Loop() {
		_ = FastExp(x)
	}
}

func BenchmarkStdExp(b *testing.B) {
	x := 1.234
	for b.Loop() {
		_ = math.Exp(x)
	}
}

func BenchmarkFastLog(b *testing.B) {
	y := 2.345
	for b.Loop() {
		_ = FastLog(y)
	}
}

func BenchmarkStdLog(b *testing.B) {
	y := 2.345
	for b.Loop() {
		_ = math.Log(y)
	}
}
