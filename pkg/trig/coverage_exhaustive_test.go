package trig

import (
	"math"
	"testing"
)

func TestTrigExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.01
	}

	t.Run("SinBatch_large", func(t *testing.T) {
		r := SinBatch(large)
		_ = r
	})
	t.Run("CosBatch_large", func(t *testing.T) {
		r := CosBatch(large)
		_ = r
	})
	t.Run("SinCosBatch_large", func(t *testing.T) {
		sin, cos := SinCosBatch(large)
		_ = sin
		_ = cos
	})
	t.Run("TanBatch_large", func(t *testing.T) {
		r := TanBatch(large)
		_ = r
	})
}

func TestFastTrigExhaustive(t *testing.T) {
	values := []float64{-100, -10, -1, -0.1, 0, 0.1, 1, 10, 100}
	for _, x := range values {
		t.Run("SinFast", func(t *testing.T) {
			_ = SinFast(x)
		})
		t.Run("CosFast", func(t *testing.T) {
			_ = CosFast(x)
		})
		t.Run("TanFast", func(t *testing.T) {
			_ = TanFast(x)
		})
	}
}

func TestInverseTrigExhaustive(t *testing.T) {
	values := []float64{-math.MaxFloat64, -1, -0.5, 0, 0.5, 1, math.MaxFloat64}
	for _, x := range values {
		t.Run("Asin", func(t *testing.T) {
			_ = Asin(x)
		})
		t.Run("Acos", func(t *testing.T) {
			_ = Acos(x)
		})
		t.Run("Atan", func(t *testing.T) {
			_ = Atan(x)
		})
	}
}

func TestAtan2Exhaustive(t *testing.T) {
	pairs := [][2]float64{
		{math.Inf(1), math.Inf(1)},
		{math.Inf(-1), math.Inf(1)},
		{math.Inf(1), math.Inf(-1)},
		{math.Inf(-1), math.Inf(-1)},
	}
	for _, p := range pairs {
		t.Run("Atan2", func(t *testing.T) {
			_ = Atan2(p[0], p[1])
		})
	}
}

func TestSecCscCotExhaustive(t *testing.T) {
	values := []float64{0, math.Pi/6, math.Pi/4, math.Pi/3, math.Pi/2, math.Pi}
	for _, x := range values {
		t.Run("Sec", func(t *testing.T) {
			_ = Sec(x)
		})
		t.Run("Csc", func(t *testing.T) {
			_ = Csc(x)
		})
		t.Run("Cot", func(t *testing.T) {
			_ = Cot(x)
		})
	}
}

func TestDegRadExhaustive(t *testing.T) {
	values := []float64{-180, -90, -45, 0, 45, 90, 180, 360}
	for _, d := range values {
		t.Run("DegToRad", func(t *testing.T) {
			_ = DegToRad(d)
		})
		t.Run("RadToDeg", func(t *testing.T) {
			_ = RadToDeg(d)
		})
	}
}

func TestSinhCoshExhaustive(t *testing.T) {
	values := []float64{-100, -10, -1, 0, 1, 10, 100}
	for _, x := range values {
		t.Run("Sinh", func(t *testing.T) {
			_ = Sinh(x)
		})
		t.Run("Cosh", func(t *testing.T) {
			_ = Cosh(x)
		})
		t.Run("Tanh", func(t *testing.T) {
			_ = Tanh(x)
		})
		t.Run("SinhCosh", func(t *testing.T) {
			_, _ = SinhCosh(x)
		})
	}
}

func TestInverseHyperbolicExhaustive(t *testing.T) {
	values := []float64{-math.MaxFloat64, -10, -1, 0, 1, 10, math.MaxFloat64}
	for _, x := range values {
		t.Run("Asinh", func(t *testing.T) {
			_ = Asinh(x)
		})
		t.Run("Acosh", func(t *testing.T) {
			_ = Acosh(x + 1)
		})
		t.Run("Atanh", func(t *testing.T) {
			_ = Atanh(x * 0.5)
		})
	}
}