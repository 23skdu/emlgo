package hyper

import (
	"math"
	"testing"
)

func TestTanhExhaustive(t *testing.T) {
	values := []float64{-1e10, -1e5, -10, -1, -0.1, 0, 0.1, 1, 10, 1e5, 1e10, math.MaxFloat64, math.SmallestNonzeroFloat64}
	for _, x := range values {
		_ = Tanh(x)
	}
}

func TestAsinhExhaustive(t *testing.T) {
	values := []float64{-1e10, -1e5, -10, -1, -0.1, -1e-10, 0, 1e-10, 0.1, 1, 10, 1e5, 1e10, math.MaxFloat64, math.SmallestNonzeroFloat64}
	for _, x := range values {
		_ = Asinh(x)
	}
}

func TestAtanhExhaustive(t *testing.T) {
	values := []float64{-1 + 1e-10, -0.9, -0.5, -0.1, -1e-10, 0, 1e-10, 0.1, 0.5, 0.9, 1 - 1e-10}
	for _, x := range values {
		_ = Atanh(x)
	}
}

func TestAcoshLarge(t *testing.T) {
	values := []float64{1, 1.1, 2, 10, 100, 1e5, 1e10, math.MaxFloat64}
	for _, x := range values {
		_ = Acosh(x)
	}
}