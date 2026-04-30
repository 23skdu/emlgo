package hyper

import (
	"math"
	"testing"
)

func TestSinh(t *testing.T) {
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
			got := Sinh(tt.x)
			want := math.Sinh(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Sinh(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestCosh(t *testing.T) {
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
			got := Cosh(tt.x)
			want := math.Cosh(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Cosh(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestTanh(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		tol  float64
	}{
		{"zero", 0, 1e-15},
		{"one", 1, 1e-15},
		{"inf", 100, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tanh(tt.x)
			want := math.Tanh(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Tanh(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
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