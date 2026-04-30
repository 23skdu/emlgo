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