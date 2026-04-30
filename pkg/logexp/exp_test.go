package logexp

import (
	"math"
	"testing"
)

func TestExp(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		tol  float64
	}{
		{"zero", 0, 1e-15},
		{"one", 1, 1e-15},
		{"negone", -1, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exp(tt.x)
			want := math.Exp(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Exp(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		tol  float64
	}{
		{"one", 1, 1e-15},
		{"e", math.E, 1e-15},
		{"two", 2, 1e-15},
		{"ten", 10, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Log(tt.x)
			want := math.Log(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Log(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestLogNaN(t *testing.T) {
	if !math.IsNaN(Log(0)) {
		t.Error("Log(0) should be NaN")
	}
	if !math.IsNaN(Log(-1)) {
		t.Error("Log(-1) should be NaN")
	}
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