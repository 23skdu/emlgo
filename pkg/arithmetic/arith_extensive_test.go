package arithmetic

import (
	"math"
	"testing"
)

func TestAddComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"zero", 0, 0, 0},
		{"pos_pos", 3, 5, 8},
		{"neg_neg", -3, -5, -8},
		{"pos_neg", 5, -3, 2},
		{"neg_pos", -3, 5, 2},
		{"float", 1.5, 2.5, 4},
		{"inf_pos", math.Inf(1), 1, math.Inf(1)},
		{"neg_inf", math.Inf(-1), 1, math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.x, tt.y)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestSubComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"zero", 0, 0, 0},
		{"pos", 5, 3, 2},
		{"neg", 3, 5, -2},
		{"float", 5.5, 2.5, 3},
		{"inf", math.Inf(1), math.Inf(1), math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sub(tt.x, tt.y)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Sub(%v, %v) = %v, want NaN", tt.x, tt.y, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Sub(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestMulComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"zero_zero", 0, 0, 0},
		{"zero_nonzero", 0, 5, 0},
		{"pos", 3, 5, 15},
		{"neg", -3, 5, -15},
		{"both_neg", -3, -5, 15},
		{"float", 1.5, 2, 3},
		{"inf_pos", math.Inf(1), 1, math.Inf(1)},
		{"inf_neg", math.Inf(1), -1, math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mul(tt.x, tt.y)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Mul(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestDivComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"pos", 6, 3, 2},
		{"neg", 6, -3, -2},
		{"zero_nonzero", 0, 5, 0},
		{"pos_zero", 1, 0, math.Inf(1)},
		{"neg_zero", -1, 0, math.Inf(-1)},
		{"zero_zero", 0, 0, math.NaN()},
		{"float", 7.5, 2.5, 3},
		{"inf_div", math.Inf(1), 1, math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Div(tt.x, tt.y)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Div(%v, %v) = %v, want NaN", tt.x, tt.y, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Div(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestModComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"pos", 10, 3, 1},
		{"neg", -10, 3, -1},
		{"zero", 0, 3, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mod(tt.x, tt.y)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Mod(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestPowComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"pos_int", 2, 3, 8},
		{"zero_exp", 0, 2, 0},
		{"one_exp", 5, 1, 5},
		{"any_zero_exp", 3, 0, 1},
		{"sqrt", 4, 0.5, 2},
		{"neg_even", -2, 2, 4},
		{"neg_odd", -2, 3, -8},
		{"neg_nonint", -2, 0.5, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pow(tt.x, tt.y)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Pow(%v, %v) = %v, want NaN", tt.x, tt.y, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Pow(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestPowInt(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		n        int
		expected float64
	}{
		{"pos", 2, 3, 8},
		{"neg_exp", 2, -2, 0.25},
		{"zero_exp", 5, 0, 1},
		{"zero_base", 0, 2, 0},
		{"one_exp", 10, 1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PowInt(tt.x, tt.n)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("PowInt(%v, %d) = %v, want %v", tt.x, tt.n, got, tt.expected)
			}
		})
	}
}

func TestLogBaseComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, base  float64
		expected float64
	}{
		{"log2_4", 4, 2, 2},
		{"log10_100", 100, 10, 2},
		{"log_e_e", math.E, math.E, 1},
		{"zero", 0, 2, math.NaN()},
		{"neg", -1, 2, math.NaN()},
		{"base_one", 4, 1, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LogBase(tt.x, tt.base)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("LogBase(%v, %v) = %v, want NaN", tt.x, tt.base, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("LogBase(%v, %v) = %v, want %v", tt.x, tt.base, got, tt.expected)
			}
		})
	}
}

func TestSqrtComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pos", 4, 2},
		{"one", 1, 1},
		{"float", 2, math.Sqrt(2)},
		{"neg", -1, math.NaN()},
		{"inf", math.Inf(1), math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Sqrt(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCbrtComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pos", 8, 2},
		{"neg", -8, -2},
		{"one", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cbrt(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Cbrt(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestHypotComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"pos", 3, 4, 5},
		{"zero", 0, 0, 0},
		{"one", 1, 1, math.Sqrt2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Hypot(tt.x, tt.y)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Hypot(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestMaxComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"pos", 3, 5, 5},
		{"neg", -3, -5, -3},
		{"nan", math.NaN(), 5, 5},
		{"equal", 5, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Max(tt.x, tt.y)
			if got != tt.expected {
				t.Errorf("Max(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestMinComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"pos", 3, 5, 3},
		{"neg", -3, -5, -5},
		{"nan", math.NaN(), 5, 5},
		{"equal", 5, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Min(tt.x, tt.y)
			if got != tt.expected {
				t.Errorf("Min(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestFloorComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 3.7, 3},
		{"neg", -3.7, -4},
		{"int", 5, 5},
		{"zero", 0, 0},
		{"inf", math.Inf(1), math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Floor(tt.x)
			if got != tt.expected {
				t.Errorf("Floor(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCeilComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 3.2, 4},
		{"neg", -3.2, -3},
		{"int", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ceil(tt.x)
			if got != tt.expected {
				t.Errorf("Ceil(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestTruncComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 3.7, 3},
		{"neg", -3.7, -3},
		{"int", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Trunc(tt.x)
			if got != tt.expected {
				t.Errorf("Trunc(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestRoundComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 3.5, 4},
		{"neg", -3.5, -4},
		{"int", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Round(tt.x)
			if got != tt.expected {
				t.Errorf("Round(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAbsComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 5, 5},
		{"neg", -5, 5},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Abs(tt.x)
			if got != tt.expected {
				t.Errorf("Abs(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestNegComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 5, -5},
		{"neg", -5, 5},
		{"zero", 0, -0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Neg(tt.x)
			if got != tt.expected && !(math.Signbit(got) && math.Signbit(tt.expected)) {
				t.Errorf("Neg(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestInvComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 2, 0.5},
		{"neg", -2, -0.5},
		{"one", 1, 1},
		{"zero", 0, math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Inv(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Inv(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestSquareComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 3, 9},
		{"neg", -3, 9},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Square(tt.x)
			if got != tt.expected {
				t.Errorf("Square(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCubeComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pos", 2, 8},
		{"neg", -2, -8},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cube(tt.x)
			if got != tt.expected {
				t.Errorf("Cube(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestLogComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"e", math.E, 1},
		{"ten", 10, math.Ln10},
		{"zero", 0, math.Inf(-1)},
		{"neg", -1, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Log(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Log(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Log(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestLog1pComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.Ln2},
		{"neg_one", -0.5, math.Log(0.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Log1p(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Log1p(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestExpM1Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.E - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpM1(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("ExpM1(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestFMAComprehensive(t *testing.T) {
	got := FMA(2, 3, 1)
	expected := 7
	if got != float64(expected) {
		t.Errorf("FMA(2, 3, 1) = %v, want 7", got)
	}
}

func TestGCDComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected int64
	}{
		{"pos", 48, 18, 6},
		{"both_pos", 100, 25, 25},
		{"neg", -48, 18, 6},
		{"zero", 0, 5, 5},
		{"one", 1, 100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GCD(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("GCD(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestLCMComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected int64
	}{
		{"pos", 4, 6, 12},
		{"both", 5, 5, 5},
		{"zero", 0, 5, 0},
		{"one", 1, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LCM(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("LCM(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func close(a, b, tol float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return true
	}
	if math.IsInf(a, -1) && math.IsInf(b, -1) {
		return true
	}
	return math.Abs(a-b) <= tol
}

func TestExpComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"one", 1, math.E},
		{"neg", -1, 1 / math.E},
		{"pos", 2, math.E * math.E},
		{"inf_pos", math.Inf(1), math.Inf(1)},
		{"inf_neg", math.Inf(-1), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exp(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Exp(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestExpm1Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.E - 1},
		{"neg", -1, 1/math.E - 1},
		{"inf_pos", math.Inf(1), math.Inf(1)},
		{"inf_neg", math.Inf(-1), -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Expm1(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Expm1(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Expm1(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestRemainderComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"pos", 7, 3, 1},
		{"neg", -7, 3, -1},
		{"zero", 0, 3, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Remainder(tt.x, tt.y)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Remainder(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestLogBase2Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"two", 2, 1},
		{"half", 0.5, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LogBase2(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("LogBase2(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("LogBase2(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestLogBase10Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"ten", 10, 1},
		{"hundred", 100, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LogBase10(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("LogBase10(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("LogBase10(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestIntFunctionsComprehensive(t *testing.T) {
	t.Run("IntAdd", func(t *testing.T) {
		if got := IntAdd(2, 3); got != 5 {
			t.Errorf("IntAdd(2, 3) = %v, want 5", got)
		}
	})
	t.Run("IntSub", func(t *testing.T) {
		if got := IntSub(5, 3); got != 2 {
			t.Errorf("IntSub(5, 3) = %v, want 2", got)
		}
	})
	t.Run("IntMul", func(t *testing.T) {
		if got := IntMul(3, 4); got != 12 {
			t.Errorf("IntMul(3, 4) = %v, want 12", got)
		}
	})
	t.Run("IntDiv", func(t *testing.T) {
		if got := IntDiv(10, 3); got != 3 {
			t.Errorf("IntDiv(10, 3) = %v, want 3", got)
		}
	})
	t.Run("IntDiv_by_zero", func(t *testing.T) {
		if got := IntDiv(10, 0); got != 0 {
			t.Errorf("IntDiv(10, 0) = %v, want 0", got)
		}
	})
	t.Run("IntMod", func(t *testing.T) {
		if got := IntMod(10, 3); got != 1 {
			t.Errorf("IntMod(10, 3) = %v, want 1", got)
		}
	})
	t.Run("IntMod_by_zero", func(t *testing.T) {
		if got := IntMod(10, 0); got != 0 {
			t.Errorf("IntMod(10, 0) = %v, want 0", got)
		}
	})
	t.Run("IntAbs_pos", func(t *testing.T) {
		if got := IntAbs(5); got != 5 {
			t.Errorf("IntAbs(5) = %v, want 5", got)
		}
	})
	t.Run("IntAbs_neg", func(t *testing.T) {
		if got := IntAbs(-5); got != 5 {
			t.Errorf("IntAbs(-5) = %v, want 5", got)
		}
	})
	t.Run("IntMax", func(t *testing.T) {
		if got := IntMax(3, 7); got != 7 {
			t.Errorf("IntMax(3, 7) = %v, want 7", got)
		}
	})
	t.Run("IntMin", func(t *testing.T) {
		if got := IntMin(3, 7); got != 3 {
			t.Errorf("IntMin(3, 7) = %v, want 3", got)
		}
	})
}

func TestUintFunctionsComprehensive(t *testing.T) {
	t.Run("UintAdd", func(t *testing.T) {
		if got := UintAdd(2, 3); got != 5 {
			t.Errorf("UintAdd(2, 3) = %v, want 5", got)
		}
	})
	t.Run("UintSub", func(t *testing.T) {
		if got := UintSub(5, 3); got != 2 {
			t.Errorf("UintSub(5, 3) = %v, want 2", got)
		}
	})
	t.Run("UintMul", func(t *testing.T) {
		if got := UintMul(3, 4); got != 12 {
			t.Errorf("UintMul(3, 4) = %v, want 12", got)
		}
	})
	t.Run("UintDiv", func(t *testing.T) {
		if got := UintDiv(10, 3); got != 3 {
			t.Errorf("UintDiv(10, 3) = %v, want 3", got)
		}
	})
	t.Run("UintDiv_by_zero", func(t *testing.T) {
		if got := UintDiv(10, 0); got != 0 {
			t.Errorf("UintDiv(10, 0) = %v, want 0", got)
		}
	})
	t.Run("UintMod", func(t *testing.T) {
		if got := UintMod(10, 3); got != 1 {
			t.Errorf("UintMod(10, 3) = %v, want 1", got)
		}
	})
	t.Run("UintMod_by_zero", func(t *testing.T) {
		if got := UintMod(10, 0); got != 0 {
			t.Errorf("UintMod(10, 0) = %v, want 0", got)
		}
	})
	t.Run("UintMax", func(t *testing.T) {
		if got := UintMax(3, 7); got != 7 {
			t.Errorf("UintMax(3, 7) = %v, want 7", got)
		}
	})
	t.Run("UintMin", func(t *testing.T) {
		if got := UintMin(3, 7); got != 3 {
			t.Errorf("UintMin(3, 7) = %v, want 3", got)
		}
	})
}

func TestBatchFunctionsComprehensive(t *testing.T) {
	testData := []float64{1, 2, 3, 4, 5}
	testDataNeg := []float64{-1, -2, -3, -4, -5}
	testDataSmall := []float64{1, 2}

	t.Run("SqrtBatch", func(t *testing.T) {
		result := SqrtBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("SqrtBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("SqrtBatch_empty", func(t *testing.T) {
		result := SqrtBatch([]float64{})
		if result != nil && len(result) != 0 {
			t.Errorf("SqrtBatch empty = %v", result)
		}
	})
	t.Run("AbsBatch", func(t *testing.T) {
		result := AbsBatch(testDataNeg)
		for i, v := range result {
			if v < 0 {
				t.Errorf("AbsBatch[%d] = %v, want >= 0", i, v)
			}
		}
	})
	t.Run("AbsBatch_small", func(t *testing.T) {
		result := AbsBatch(testDataSmall)
		if len(result) != len(testDataSmall) {
			t.Errorf("AbsBatch length = %v, want %v", len(result), len(testDataSmall))
		}
	})
	t.Run("FloorBatch", func(t *testing.T) {
		result := FloorBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("FloorBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("CeilBatch", func(t *testing.T) {
		result := CeilBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("CeilBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("TruncBatch", func(t *testing.T) {
		result := TruncBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("TruncBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("Log1pBatch", func(t *testing.T) {
		result := Log1pBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("Log1pBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("Expm1Batch", func(t *testing.T) {
		result := Expm1Batch(testData)
		if len(result) != len(testData) {
			t.Errorf("Expm1Batch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("PowBatch", func(t *testing.T) {
		result := PowBatch(testData, 2)
		if len(result) != len(testData) {
			t.Errorf("PowBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("CbrtBatch", func(t *testing.T) {
		result := CbrtBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("CbrtBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("HypotBatch", func(t *testing.T) {
		result := HypotBatch(testData, testDataNeg)
		if len(result) != len(testData) {
			t.Errorf("HypotBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("HypotBatch_empty", func(t *testing.T) {
		result := HypotBatch([]float64{}, []float64{})
		if result != nil && len(result) != 0 {
			t.Errorf("HypotBatch empty = %v", result)
		}
	})
	t.Run("MaxBatch", func(t *testing.T) {
		result := MaxBatch(testData, 3)
		if len(result) != len(testData) {
			t.Errorf("MaxBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("MinBatch", func(t *testing.T) {
		result := MinBatch(testData, 3)
		if len(result) != len(testData) {
			t.Errorf("MinBatch length = %v, want %v", len(result), len(testData))
		}
	})
}

func TestFusedBatchFunctionsComprehensive(t *testing.T) {
	testData := []float64{1, 2, 3, 4, 5}
	testData2 := []float64{2, 3, 4, 5, 6}

	t.Run("AddBatch", func(t *testing.T) {
		result := AddBatch(testData, testData2)
		if len(result) != len(testData) {
			t.Errorf("AddBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("SubBatch", func(t *testing.T) {
		result := SubBatch(testData, testData2)
		if len(result) != len(testData) {
			t.Errorf("SubBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("MulBatch", func(t *testing.T) {
		result := MulBatch(testData, testData2)
		if len(result) != len(testData) {
			t.Errorf("MulBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("DivBatch", func(t *testing.T) {
		result := DivBatch(testData, testData2)
		if len(result) != len(testData) {
			t.Errorf("DivBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("AddScalarBatch", func(t *testing.T) {
		result := AddScalarBatch(testData, 10)
		if len(result) != len(testData) {
			t.Errorf("AddScalarBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("MulScalarBatch", func(t *testing.T) {
		result := MulScalarBatch(testData, 2)
		if len(result) != len(testData) {
			t.Errorf("MulScalarBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("NegBatch", func(t *testing.T) {
		result := NegBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("NegBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("InvBatch", func(t *testing.T) {
		result := InvBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("InvBatch length = %v, want %v", len(result), len(testData))
		}
	})
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(1.5, 2.5)
	}
}

func BenchmarkMul(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Mul(1.5, 2.5)
	}
}

func BenchmarkDiv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Div(1.5, 2.5)
	}
}

func BenchmarkPow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pow(1.5, 2.5)
	}
}

func BenchmarkSqrt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sqrt(1.5)
	}
}