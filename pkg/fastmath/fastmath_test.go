package fastmath

import (
	"math"
	"testing"
)

func TestSqrt(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, 1},
		{"pos", 4, 2},
		{"neg", -1, math.NaN()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Sqrt(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if got != tt.expected {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestFMA(t *testing.T) {
	tests := []struct {
		name     string
		x, y, z  float64
		expected float64
	}{
		{"pos", 2, 3, 1, 7},
		{"zero", 0, 5, 3, 3},
		{"neg", -2, 3, 1, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FMA(tt.x, tt.y, tt.z)
			if got != tt.expected {
				t.Errorf("FMA(%v, %v, %v) = %v, want %v", tt.x, tt.y, tt.z, got, tt.expected)
			}
		})
	}
}

func TestExp(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"one", 1, 2.718},
		{"neg_one", -1, 0.367},
		{"large_pos", 800, math.Inf(1)},
		{"large_neg", -800, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exp(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Exp(%v) = %v, want NaN", tt.x, got)
			}
			if !math.IsInf(tt.expected, 0) && !math.IsInf(got, 0) && math.Abs(got-tt.expected) > 1e-2 {
				t.Errorf("Exp(%v) = %v, want ~%v", tt.x, got, tt.expected)
			}
			if math.IsInf(tt.expected, 1) && !math.IsInf(got, 1) {
				t.Errorf("Exp(%v) = %v, want Inf", tt.x, got)
			}
			if tt.expected == 0 && got > 1e-10 {
				t.Errorf("Exp(%v) = %v, want 0", tt.x, got)
			}
		})
	}
}

func TestSin(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"small_pos", 0.1, 0.0998},
		{"small_neg", -0.1, -0.0998},
		{"large", 10, math.Sin(10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sin(tt.x)
			if math.Abs(got-tt.expected) > 1e-2 {
				t.Errorf("Sin(%v) = %v, want ~%v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCos(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"small_pos", 0.1, 0.995},
		{"large", 10, math.Cos(10)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cos(tt.x)
			if math.Abs(got-tt.expected) > 1e-2 {
				t.Errorf("Cos(%v) = %v, want ~%v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"e", math.E, 1},
		{"zero", 0, math.Inf(-1)},
		{"neg", -1, math.NaN()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Log(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Log(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if !close(got, tt.expected, 1e-2) {
				t.Errorf("Log(%v) = %v, want ~%v", tt.x, got, tt.expected)
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