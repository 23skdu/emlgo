package hyper

import (
	"math"
	"testing"
)

func TestTanhEdgeCases(t *testing.T) {
	_ = Tanh(0)
	_ = Tanh(1)
	_ = Tanh(-1)
	_ = Tanh(0.5)
	_ = Tanh(-0.5)
	_ = Tanh(10)
	_ = Tanh(-10)
	_ = Tanh(100)
	_ = Tanh(-100)
	_ = Tanh(math.Inf(1))
	_ = Tanh(math.Inf(-1))
	_ = Tanh(math.NaN())
}

func TestAsinhEdgeCases(t *testing.T) {
	_ = Asinh(0)
	_ = Asinh(1)
	_ = Asinh(-1)
	_ = Asinh(0.5)
	_ = Asinh(-0.5)
	_ = Asinh(10)
	_ = Asinh(-10)
	_ = Asinh(100)
	_ = Asinh(-100)
	_ = Asinh(math.Inf(1))
	_ = Asinh(math.Inf(-1))
	_ = Asinh(math.NaN())
	_ = Asinh(1e-10)
	_ = Asinh(-1e-10)
	_ = Asinh(1e10)
	_ = Asinh(-1e10)
}

func TestAcoshEdgeCases(t *testing.T) {
	_ = Acosh(1)
	_ = Acosh(2)
	_ = Acosh(10)
	_ = Acosh(1.5)
	_ = Acosh(1.01)
	_ = Acosh(1.001)
	_ = Acosh(math.MaxFloat64)
	_ = Acosh(math.Inf(1))
	_ = Acosh(math.NaN())
	_ = Acosh(0.5)
	_ = Acosh(0)
}

func TestAtanhEdgeCases(t *testing.T) {
	_ = Atanh(0)
	_ = Atanh(0.5)
	_ = Atanh(-0.5)
	_ = Atanh(0.9)
	_ = Atanh(-0.9)
	_ = Atanh(0.99)
	_ = Atanh(-0.99)
	_ = Atanh(0.999)
	_ = Atanh(-0.999)
	_ = Atanh(1)
	_ = Atanh(-1)
	_ = Atanh(math.NaN())
	_ = Atanh(math.Inf(1))
	_ = Atanh(math.Inf(-1))
	_ = Atanh(0.9999)
	_ = Atanh(-0.9999)
}