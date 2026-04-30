package eml

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestEmlComplexAll(t *testing.T) {
	tests := []struct {
		name     string
		x        complex128
		y        complex128
		expected complex128
	}{
		{"zero", 0, 1, 1 - 0i},
		{"one_one", 1, 1, cmplx.Exp(1) - cmplx.Log(1)},
		{"real_imag", complex(1, 1), complex(2, 0), cmplx.Exp(complex(1, 1)) - cmplx.Log(complex(2, 0))},
		{"negative_real", complex(-1, 0), complex(1, 0), cmplx.Exp(-1) - 0},
		{"pure_imag", complex(0, 1), complex(1, 0), cmplx.Exp(complex(0, 1)) - 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmlComplex(tt.x, tt.y)
			want := tt.expected
			if !closeComplex(got, want, 1e-10) {
				t.Errorf("EmlComplex(%v, %v) = %v, want %v", tt.x, tt.y, got, want)
			}
		})
	}
}

func TestEmlComplexOneAll(t *testing.T) {
	tests := []struct {
		name     string
		x        complex128
		expected complex128
	}{
		{"zero", 0, 1},
		{"one", 1, cmplx.Exp(1)},
		{"real", 2, cmplx.Exp(2)},
		{"imag", complex(0, 1), cmplx.Exp(complex(0, 1))},
		{"complex", complex(1, 1), cmplx.Exp(complex(1, 1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmlComplexOne(tt.x)
			want := tt.expected
			if !closeComplex(got, want, 1e-10) {
				t.Errorf("EmlComplexOne(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func TestOneEmlComplexAll(t *testing.T) {
	tests := []struct {
		name     string
		x        complex128
		expected complex128
	}{
		{"one", 1, cmplx.Exp(1) - cmplx.Log(1)},
		{"two", 2, cmplx.Exp(1) - cmplx.Log(2)},
		{"e", complex(math.E, 0), cmplx.Exp(1) - cmplx.Log(complex(math.E, 0))},
		{"complex", complex(1, 1), cmplx.Exp(1) - cmplx.Log(complex(1, 1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OneEmlComplex(tt.x)
			want := tt.expected
			if !closeComplex(got, want, 1e-10) {
				t.Errorf("OneEmlComplex(%v) = %v, want %v", tt.x, got, want)
			}
		})
	}
}

func closeComplex(a, b complex128, tol float64) bool {
	if cmplx.IsNaN(a) && cmplx.IsNaN(b) {
		return true
	}
	if cmplx.IsInf(a) && cmplx.IsInf(b) {
		return true
	}
	return math.Abs(real(a)-real(b)) <= tol && math.Abs(imag(a)-imag(b)) <= tol
}

func BenchmarkEmlComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EmlComplex(complex(1.5, 0.5), complex(2.5, 0))
	}
}

func BenchmarkEmlComplexOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EmlComplexOne(complex(1.5, 0.5))
	}
}