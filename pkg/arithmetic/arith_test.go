package arithmetic

import (
	"math"
	"testing"
)

func TestPowStabilityNearOne(t *testing.T) {
	tests := []struct {
		x, y float64
		tol  float64
	}{
		{1.0 + 1e-15, 1000, 1e-12},
		{1.0 - 1e-15, 1000, 1e-12},
		{1.0 + 1e-10, 1e6, 1e-4},
		{0.999999, 1000000, 1e-5},
		{1.000001, 1000000, 1e-5},
	}
	for _, tc := range tests {
		got := Pow(tc.x, tc.y)
		want := math.Pow(tc.x, tc.y)
		diff := math.Abs(got - want)
		if diff > tc.tol {
			t.Errorf("Pow(%v, %v) = %v, want %v (diff=%v, tol=%v)", tc.x, tc.y, got, want, diff, tc.tol)
		}
	}
}

func TestLogStabilityNearOne(t *testing.T) {
	tests := []struct {
		x   float64
		tol float64
	}{
		{1.0 + 1e-15, 1e-27},
		{1.0 - 1e-15, 1e-27},
		{1.0 + 1e-10, 1e-22},
		{1.0 + 1e-8, 1e-20},
	}
	for _, tc := range tests {
		got := Log(tc.x)
		want := math.Log1p(tc.x - 1)
		diff := math.Abs(got - want)
		if diff > tc.tol {
			t.Errorf("Log(%v) = %v, want %v (diff=%v, tol=%v)", tc.x, got, want, diff, tc.tol)
		}
	}
}
