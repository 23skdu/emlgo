package eml

import (
	"math"
	"testing"
)

func TestEmlAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		y        float64
		expected float64
	}{
		{"x0_y1", 0, 1, 1},
		{"x1_y1", 1, 1, math.E},
		{"x2_y1", 2, 1, math.E * math.E},
		{"xneg1_y1", -1, 1, 1 / math.E},
		{"x0_y2", 0, 2, -math.Ln2 + 1},
		{"x1_y2", 1, 2, math.E - math.Ln2},
		{"xln2_y1", math.Ln2, 1, 2},
		{"xneg2_y1", -2, 1, 1 / (math.E * math.E)},
		{"x0_y0.5", 0, 0.5, 1 + math.Ln2},
		{"inf_x", math.Inf(1), 1, math.Inf(1)},
		{"inf_neg_x", math.Inf(-1), 1, 0},
		{"nan_x", math.NaN(), 1, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Eml(tt.x, tt.y)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Eml(%v, %v) = %v, want NaN", tt.x, tt.y, got)
				}
			} else if math.IsInf(tt.expected, 0) {
				if !math.IsInf(got, 0) {
					t.Errorf("Eml(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Eml(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestOneAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"one", 1, math.E},
		{"two", 2, math.E * math.E},
		{"neg_one", -1, 1 / math.E},
		{"neg_two", -2, 1 / (math.E * math.E)},
		{"half", 0.5, math.Sqrt(math.E)},
		{"neg_half", -0.5, 1 / math.Sqrt(math.E)},
		{"inf_pos", math.Inf(1), math.Inf(1)},
		{"inf_neg", math.Inf(-1), 0},
		{"nan", math.NaN(), math.NaN()},
		{"large_pos", 700, math.Exp(700)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := One(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("One(%v) = %v, want NaN", tt.x, got)
				}
			} else if math.IsInf(tt.expected, 0) {
				if !math.IsInf(got, 0) {
					t.Errorf("One(%v) = %v, want %v", tt.x, got, tt.expected)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("One(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestOneEmlAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, math.Exp(1) - math.Log(1)},
		{"e", math.E, math.Exp(1) - math.Log(math.E)},
		{"two", 2, math.Exp(1) - math.Log(2)},
		{"ten", 10, math.Exp(1) - math.Log(10)},
		{"half", 0.5, math.Exp(1) - math.Log(0.5)},
		{"pos_inf", math.Inf(1), math.Inf(1)},
		{"nan", math.NaN(), math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OneEml(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("OneEml(%v) = %v, want NaN", tt.x, got)
				}
			} else if math.IsInf(tt.expected, 0) {
				if !math.IsInf(got, 0) {
					t.Errorf("OneEml(%v) = %v, want %v", tt.x, got, tt.expected)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("OneEml(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestOneConst(t *testing.T) {
	if one != 1.0 {
		t.Errorf("one = %v, want 1.0", one)
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

func BenchmarkEml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Eml(1.5, 2.5)
	}
}

func BenchmarkOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		One(1.5)
	}
}

func BenchmarkOneEml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OneEml(1.5)
	}
}