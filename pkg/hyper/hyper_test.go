package hyper

import (
	"math"
	"testing"
)

func TestSinhAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pos_one", 1, math.Sinh(1)},
		{"neg_one", -1, math.Sinh(-1)},
		{"pos_two", 2, math.Sinh(2)},
		{"neg_two", -2, math.Sinh(-2)},
		{"large_pos", 10, math.Sinh(10)},
		{"large_neg", -10, math.Sinh(-10)},
		{"nan", math.NaN(), math.NaN()},
		{"inf_pos", math.Inf(1), math.Inf(1)},
		{"inf_neg", math.Inf(-1), math.Inf(-1)},
		{"half", 0.5, math.Sinh(0.5)},
		{"pi", math.Pi, math.Sinh(math.Pi)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sinh(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Sinh(%v) = %v, want NaN", tt.x, got)
				}
			} else if math.IsInf(tt.expected, 0) {
				if !math.IsInf(got, 0) {
					t.Errorf("Sinh(%v) = %v, want %v", tt.x, got, tt.expected)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Sinh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCoshAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"pos_one", 1, math.Cosh(1)},
		{"neg_one", -1, math.Cosh(-1)},
		{"pos_two", 2, math.Cosh(2)},
		{"neg_two", -2, math.Cosh(-2)},
		{"large_pos", 10, math.Cosh(10)},
		{"large_neg", -10, math.Cosh(-10)},
		{"nan", math.NaN(), math.NaN()},
		{"inf_pos", math.Inf(1), math.Inf(1)},
		{"inf_neg", math.Inf(-1), math.Inf(1)},
		{"half", 0.5, math.Cosh(0.5)},
		{"pi", math.Pi, math.Cosh(math.Pi)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cosh(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Cosh(%v) = %v, want NaN", tt.x, got)
				}
			} else if math.IsInf(tt.expected, 0) {
				if !math.IsInf(got, 1) {
					t.Errorf("Cosh(%v) = %v, want %v", tt.x, got, tt.expected)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Cosh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestTanhAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pos_one", 1, math.Tanh(1)},
		{"neg_one", -1, math.Tanh(-1)},
		{"pos_two", 2, math.Tanh(2)},
		{"neg_two", -2, math.Tanh(-2)},
		{"large_pos", 100, 1},
		{"large_neg", -100, -1},
		{"nan", math.NaN(), math.NaN()},
		{"inf_pos", math.Inf(1), 1},
		{"inf_neg", math.Inf(-1), -1},
		{"half", 0.5, math.Tanh(0.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tanh(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Tanh(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Tanh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAsinhAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pos_one", 1, math.Asinh(1)},
		{"neg_one", -1, math.Asinh(-1)},
		{"pos_two", 2, math.Asinh(2)},
		{"neg_two", -2, math.Asinh(-2)},
		{"large_pos", 100, math.Asinh(100)},
		{"large_neg", -100, math.Asinh(-100)},
		{"nan", math.NaN(), math.NaN()},
		{"inf_pos", math.Inf(1), math.Inf(1)},
		{"inf_neg", math.Inf(-1), math.Inf(-1)},
		{"half", 0.5, math.Asinh(0.5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asinh(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Asinh(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Asinh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAcoshAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"pos_two", 2, math.Acosh(2)},
		{"neg", 0.5, math.NaN()},
		{"zero", 0, math.NaN()},
		{"large", 100, math.Acosh(100)},
		{"nan", math.NaN(), math.NaN()},
		{"inf", math.Inf(1), math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acosh(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Acosh(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Acosh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAtanhAll(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"half", 0.5, math.Atanh(0.5)},
		{"neg_half", -0.5, math.Atanh(-0.5)},
		{"pos_one", 1, math.NaN()},
		{"neg_one", -1, math.NaN()},
		{"greater_one", 2, math.NaN()},
		{"less_neg_one", -2, math.NaN()},
		{"nan", math.NaN(), math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atanh(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Atanh(%v) = %v, want NaN", tt.x, got)
			} else if !math.IsNaN(tt.expected) && !close(got, tt.expected, 1e-10) {
				t.Errorf("Atanh(%v) = %v, want %v", tt.x, got, tt.expected)
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

func BenchmarkSinh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sinh(1.5)
	}
}

func BenchmarkCosh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Cosh(1.5)
	}
}

func BenchmarkTanh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tanh(1.5)
	}
}

func BenchmarkAsinh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Asinh(1.5)
	}
}

func BenchmarkAcosh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Acosh(1.5)
	}
}

func BenchmarkAtanh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Atanh(0.5)
	}
}