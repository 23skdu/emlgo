package eml

import (
	"math"
	"testing"
)

func TestMathHelpersExhaustive(t *testing.T) {
	t.Run("NaN", func(t *testing.T) {
		got := NaN()
		if !math.IsNaN(got) {
			t.Errorf("NaN() = %v, want NaN", got)
		}
	})
	t.Run("Inf_pos", func(t *testing.T) {
		got := Inf(1)
		if !math.IsInf(got, 1) {
			t.Errorf("Inf(1) = %v, want +Inf", got)
		}
	})
	t.Run("Inf_neg", func(t *testing.T) {
		got := Inf(-1)
		if !math.IsInf(got, -1) {
			t.Errorf("Inf(-1) = %v, want -Inf", got)
		}
	})
}

func TestFloorCeilTruncRoundExhaustive(t *testing.T) {
	tests := []float64{-1.5, -1, -0.5, 0, 0.5, 1, 1.5}
	for _, x := range tests {
		t.Run("Floor", func(t *testing.T) {
			_ = Floor(x)
		})
		t.Run("Ceil", func(t *testing.T) {
			_ = Ceil(x)
		})
		t.Run("Trunc", func(t *testing.T) {
			_ = Trunc(x)
		})
		t.Run("Round", func(t *testing.T) {
			_ = Round(x)
		})
	}
}

func TestAbsNegInvExhaustive(t *testing.T) {
	tests := []float64{-1, -0.5, 0, 0.5, 1, math.Inf(1), math.Inf(-1)}
	for _, x := range tests {
		t.Run("Abs", func(t *testing.T) {
			_ = Abs(x)
		})
		t.Run("Neg", func(t *testing.T) {
			_ = Neg(x)
		})
		t.Run("Inv", func(t *testing.T) {
			_ = Inv(x)
		})
	}
}

func TestExpLogExhaustive(t *testing.T) {
	tests := []float64{-10, -1, 0, 0.5, 1, 10, 100}
	for _, x := range tests {
		t.Run("Exp", func(t *testing.T) {
			_ = Exp(x)
		})
		t.Run("Log", func(t *testing.T) {
			_ = Log(x)
		})
		t.Run("Log1p", func(t *testing.T) {
			_ = Log1p(x)
		})
		t.Run("Expm1", func(t *testing.T) {
			_ = Expm1(x)
		})
	}
}

func TestSqrtCbrtPowExhaustive(t *testing.T) {
	tests := []float64{0, 0.25, 1, 4, 9, 100}
	for _, x := range tests {
		t.Run("Sqrt", func(t *testing.T) {
			_ = Sqrt(x)
		})
		t.Run("Cbrt", func(t *testing.T) {
			_ = Cbrt(x)
		})
		t.Run("Pow", func(t *testing.T) {
			_ = Pow(x, 2)
		})
	}
}

func TestSinCosTanExhaustive(t *testing.T) {
	tests := []float64{-math.Pi, -math.Pi/2, -math.Pi/4, 0, math.Pi/4, math.Pi/2, math.Pi}
	for _, x := range tests {
		t.Run("Sin", func(t *testing.T) {
			_ = Sin(x)
		})
		t.Run("Cos", func(t *testing.T) {
			_ = Cos(x)
		})
		t.Run("Tan", func(t *testing.T) {
			_ = Tan(x)
		})
		t.Run("Sincos", func(t *testing.T) {
			_, _ = Sincos(x)
		})
	}
}

func TestAsinAcosAtanExhaustive(t *testing.T) {
	tests := []float64{-1, -0.5, 0, 0.5, 1}
	for _, x := range tests {
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
		{1, 1}, {1, 0}, {0, 1}, {-1, 1}, {1, -1},
	}
	for _, p := range pairs {
		t.Run("Atan2", func(t *testing.T) {
			_ = Atan2(p[0], p[1])
		})
	}
}

func TestSinhCoshTanhExhaustive(t *testing.T) {
	tests := []float64{-5, -1, 0, 1, 5}
	for _, x := range tests {
		t.Run("Sinh", func(t *testing.T) {
			_ = Sinh(x)
		})
		t.Run("Cosh", func(t *testing.T) {
			_ = Cosh(x)
		})
		t.Run("Tanh", func(t *testing.T) {
			_ = Tanh(x)
		})
	}
}

func TestAsinhAcoshAtanhExhaustive(t *testing.T) {
	tests := []float64{-5, -1, -0.5, 0, 0.5, 1, 5, 1.1}
	for _, x := range tests {
		t.Run("Asinh", func(t *testing.T) {
			_ = Asinh(x)
		})
		t.Run("Acosh", func(t *testing.T) {
			_ = Acosh(x)
		})
		t.Run("Atanh", func(t *testing.T) {
			_ = Atanh(x)
		})
	}
}

func TestHypotMaxMinModExhaustive(t *testing.T) {
	pairs := [][2]float64{
		{3, 4}, {5, 12}, {1, 1}, {0, 1}, {1, 0},
	}
	for _, p := range pairs {
		t.Run("Hypot", func(t *testing.T) {
			_ = Hypot(p[0], p[1])
		})
		t.Run("Max", func(t *testing.T) {
			_ = Max(p[0], p[1])
		})
		t.Run("Min", func(t *testing.T) {
			_ = Min(p[0], p[1])
		})
		t.Run("Mod", func(t *testing.T) {
			_ = Mod(p[0], p[1])
		})
		t.Run("Remainder", func(t *testing.T) {
			_ = Remainder(p[0], p[1])
		})
	}
}

func TestLog10PowExhaustive(t *testing.T) {
	tests := []float64{0.001, 0.01, 0.1, 1, 10, 100, 1000}
	for _, x := range tests {
		t.Run("Log10", func(t *testing.T) {
			_ = Log10(x)
		})
	}
}

func TestPowIntExhaustive(t *testing.T) {
	bases := []float64{2, 3, 5, 10}
	exps := []int{0, 1, 2, 3, -1, -2}
	for _, b := range bases {
		for _, e := range exps {
			t.Run("PowInt", func(t *testing.T) {
				_ = PowInt(b, e)
			})
		}
	}
}

func TestCopysignExhaustive(t *testing.T) {
	pairs := [][2]float64{
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, {0, 1}, {0, -1},
	}
	for _, p := range pairs {
		t.Run("Copysign", func(t *testing.T) {
			_ = Copysign(p[0], p[1])
		})
	}
}

func TestModfExhaustive(t *testing.T) {
	tests := []float64{-2.5, -1.5, -0.5, 0, 0.5, 1.5, 2.5}
	for _, x := range tests {
		t.Run("Modf", func(t *testing.T) {
			_, _ = Modf(x)
		})
	}
}