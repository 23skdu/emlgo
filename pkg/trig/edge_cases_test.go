package trig

import (
	"math"
	"testing"
)

func TestTrigEdgeCases(t *testing.T) {
	t.Run("TanFast_ZeroCos", func(t *testing.T) {
		TanFast(math.Pi / 2)
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
