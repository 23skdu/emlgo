package hyper

import (
	"math"
	"testing"
)

func TestHyperEdgeCasesMore(t *testing.T) {
	Tanh(1000)
	Tanh(-1000)
	Tanh(math.NaN())
	Atanh(0.5)
	Atanh(1)
	Atanh(-1)
	Atanh(2)
	Atanh(math.NaN())
}

func TestHyperFinal(t *testing.T) {
	Tanh(709.7)
	Asinh(1e200)
	Asinh(-1e200)
	Acosh(1e308)
	Atanh(0.9999999999999999)
	Atanh(-0.9999999999999999)
}
