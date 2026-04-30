package logexp

import (
	"math"
	"testing"
)

func TestExpAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"pos_one", 1, math.E},
		{"neg_one", 1, math.E},
		{"pos_two", 2, math.E * math.E},
		{"neg_two", -2, 1 / (math.E * math.E)},
		{"large_pos", 100, math.Exp(100)},
		{"large_neg", -100, 0},
		{"nan", math.NaN(), math.NaN()},
		{"inf_pos", math.Inf(1), math.Inf(1)},
		{"inf_neg", math.Inf(-1), 0},
		{"half", 0.5, math.Sqrt(math.E)},
		{"neg_half", -0.5, 1 / math.Sqrt(math.E)},
		{"pi", math.Pi, math.Exp(math.Pi)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exp(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Exp(%v) = %v, want NaN", tt.x, got)
			} else if math.IsInf(tt.expected, 0) {
				if !math.IsInf(got, 0) {
					t.Errorf("Exp(%v) = %v, want %v", tt.x, got, tt.expected)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Exp(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestLogAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"e", math.E, 1},
		{"e_squared", math.E * math.E, 2},
		{"ten", 10, math.Ln10},
		{"half", 0.5, -math.Ln2},
		{"pos_small", 0.001, math.Log(0.001)},
		{"zero", 0, math.NaN()},
		{"neg", -1, math.NaN()},
		{"nan", math.NaN(), math.NaN()},
		{"inf", math.Inf(1), math.Inf(1)},
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

func BenchmarkExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Exp(1.5)
	}
}

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log(1.5)
	}
}