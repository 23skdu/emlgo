package trig

import (
	"math"
	"testing"
)

func TestSin(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		tol  float64
	}{
		{"zero", 0, 1e-15},
		{"pi", math.Pi, 1e-15},
		{"pi/2", math.Pi / 2, 1e-15},
		{"pi/4", math.Pi / 4, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sin(tt.x)
			want := math.Sin(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Sin(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestCos(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		tol  float64
	}{
		{"zero", 0, 1e-15},
		{"pi", math.Pi, 1e-15},
		{"pi/2", math.Pi / 2, 1e-15},
		{"pi/4", math.Pi / 4, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cos(tt.x)
			want := math.Cos(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Cos(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestTan(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		tol  float64
	}{
		{"zero", 0, 1e-15},
		{"pi/4", math.Pi / 4, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tan(tt.x)
			want := math.Tan(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Tan(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestAsin(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		tol  float64
	}{
		{"zero", 0, 1e-15},
		{"half", 0.5, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asin(tt.x)
			want := math.Asin(tt.x)
			if math.Abs(got-want) > tt.tol {
				t.Errorf("Asin(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestAsinNaN(t *testing.T) {
	if !math.IsNaN(Asin(2)) {
		t.Error("Asin(2) should be NaN")
	}
}

func BenchmarkSin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sin(1.5)
	}
}

func BenchmarkCos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Cos(1.5)
	}
}