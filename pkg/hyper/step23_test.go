package hyper

import (
	"math"
	"testing"

	"github.com/emlgo/eml/pkg/logexp"
)

func TestAsinhNegativeLargeValuesNoMinusInf(t *testing.T) {
	// Previously Asinh(-1e8) through Asinh(-1e150) returned -Inf
	values := []float64{-1e8, -1e10, -1e20, -1e50, -1e100, -1e150}
	for _, x := range values {
		got := Asinh(x)
		want := math.Asinh(x)
		if math.IsInf(got, -1) {
			t.Fatalf("Asinh(%e) returned -Inf", x)
		}
		relErr := math.Abs((got - want) / want)
		if relErr > 1e-10 {
			t.Errorf("Asinh(%e) = %v, want %v (relErr %e)", x, got, want, relErr)
		}
	}
}

func TestSinhSmallValuesNoCancellation(t *testing.T) {
	// Previously Sinh(1e-16) suffered 44.5% error
	smallValues := []float64{1e-16, 5e-16, 1e-15, 1e-12, 1e-8, -1e-16, -1e-12}
	for _, x := range smallValues {
		got := Sinh(x)
		want := math.Sinh(x)
		relErr := math.Abs((got - want) / want)
		if relErr > 1e-10 {
			t.Errorf("Sinh(%e) = %e, want %e (relErr %e)", x, got, want, relErr)
		}
	}
}

func TestLogZeroReturnsMinusInf(t *testing.T) {
	got := logexp.Log(0)
	if !math.IsInf(got, -1) {
		t.Fatalf("logexp.Log(0) = %v, want -Inf", got)
	}

	gotFast := logexp.LogFast(0)
	if !math.IsInf(gotFast, -1) {
		t.Fatalf("logexp.LogFast(0) = %v, want -Inf", gotFast)
	}

	gotNeg := logexp.Log(-1.0)
	if !math.IsNaN(gotNeg) {
		t.Fatalf("logexp.Log(-1.0) = %v, want NaN", gotNeg)
	}
}

func FuzzHyperbolicAccuracy(f *testing.F) {
	f.Add(float64(0.0))
	f.Add(float64(1.0))
	f.Add(float64(-1.0))
	f.Add(float64(1e-16))
	f.Add(float64(-1e8))
	f.Add(float64(50.0))

	f.Fuzz(func(t *testing.T, x float64) {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return
		}

		// Asinh
		gotAsinh := Asinh(x)
		wantAsinh := math.Asinh(x)
		diffAsinh := math.Abs(gotAsinh - wantAsinh)
		denomAsinh := math.Max(math.Abs(gotAsinh), math.Abs(wantAsinh))
		if denomAsinh > 1e-9 {
			diffAsinh /= denomAsinh
		}
		if diffAsinh > 1e-7 {
			t.Fatalf("Asinh(%e): got %v, want %v", x, gotAsinh, wantAsinh)
		}

		// Sinh (if within overflow range)
		if math.Abs(x) < 700.0 {
			gotSinh := Sinh(x)
			wantSinh := math.Sinh(x)
			diffSinh := math.Abs(gotSinh - wantSinh)
			denomSinh := math.Max(math.Abs(gotSinh), math.Abs(wantSinh))
			if denomSinh > 1e-9 {
				diffSinh /= denomSinh
			}
			if diffSinh > 1e-7 {
				t.Fatalf("Sinh(%e): got %v, want %v", x, gotSinh, wantSinh)
			}
		}

		// Tanh
		gotTanh := Tanh(x)
		wantTanh := math.Tanh(x)
		diffTanh := math.Abs(gotTanh - wantTanh)
		if diffTanh > 1e-7 {
			t.Fatalf("Tanh(%e): got %v, want %v", x, gotTanh, wantTanh)
		}
	})
}
