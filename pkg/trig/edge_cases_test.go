package trig

import (
	"math"
	"testing"
)

func TestTrigEdgeCases(t *testing.T) {
	t.Run("TanFast_ZeroCos", func(t *testing.T) {
		TanFast(math.Pi / 2)
	})
	t.Run("IEEE754_Poles_And_Asymptotics", func(t *testing.T) {
		if !math.IsInf(Atanh(1.0), 1) {
			t.Fatalf("Atanh(1.0) = %v, want +Inf", Atanh(1.0))
		}
		if !math.IsInf(Atanh(-1.0), -1) {
			t.Fatalf("Atanh(-1.0) = %v, want -Inf", Atanh(-1.0))
		}
		negZero := math.Copysign(0.0, -1.0)
		if !math.IsInf(Acsch(negZero), -1) {
			t.Fatalf("Acsch(-0.0) = %v, want -Inf", Acsch(negZero))
		}
		if !math.IsInf(Cot(0.0), 1) || !math.IsInf(Cot(negZero), -1) {
			t.Fatalf("Cot(±0.0) = (%v, %v), want (+Inf, -Inf)", Cot(0.0), Cot(negZero))
		}
		if got, want := Asinh(-1e160), math.Asinh(-1e160); math.Abs(got-want) > 1e-12*math.Abs(want) {
			t.Fatalf("Asinh(-1e160) = %v, want %v", got, want)
		}
		if got, want := Acosh(1e160), math.Acosh(1e160); math.Abs(got-want) > 1e-12*math.Abs(want) {
			t.Fatalf("Acosh(1e160) = %v, want %v", got, want)
		}
	})
}

func FuzzTrigPythagoreanIdentity(f *testing.F) {
	f.Add(0.0)
	f.Add(0.5)
	f.Add(-1.25)
	f.Add(math.Pi / 4)
	f.Fuzz(func(t *testing.T, x float64) {
		if math.IsNaN(x) || math.IsInf(x, 0) || math.Abs(x) > 1e6 {
			return
		}
		s, c := Sin(x), Cos(x)
		sum := s*s + c*c
		if math.Abs(sum-1.0) > 1e-9 {
			t.Fatalf("Sin(%g)^2 + Cos(%g)^2 = %g, want 1.0", x, x, sum)
		}
	})
}
