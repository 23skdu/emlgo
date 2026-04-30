package arithmetic

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		want float64
	}{
		{"basic", 2, 3, 5},
		{"zero", 0, 5, 5},
		{"neg", 2, -3, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.x, tt.y)
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		want float64
	}{
		{"basic", 2, 3, 6},
		{"zero", 0, 5, 0},
		{"one", 1, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mul(tt.x, tt.y)
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("Mul(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	if Div(6, 3) != 2 {
		t.Error("Div(6, 3) should be 2")
	}
}

func TestDivZero(t *testing.T) {
	if !math.IsInf(Div(1, 0), 1) {
		t.Error("Div(1, 0) should be +Inf")
	}
}

func TestPow(t *testing.T) {
	if Pow(2, 3) != 8 {
		t.Error("Pow(2, 3) should be 8")
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"four", 4, 2},
		{"one", 1, 1},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.x)
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestSqrtNaN(t *testing.T) {
	if !math.IsNaN(Sqrt(-1)) {
		t.Error("Sqrt(-1) should be NaN")
	}
}

func TestNeg(t *testing.T) {
	if Neg(5) != -5 {
		t.Error("Neg(5) should be -5")
	}
}

func TestInv(t *testing.T) {
	if Inv(2) != 0.5 {
		t.Error("Inv(2) should be 0.5")
	}
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