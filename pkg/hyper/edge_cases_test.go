package hyper

import (
	"math"
	"testing"
)

func TestHyperEdgeCases(t *testing.T) {
	t.Run("Tanh_Edge", func(t *testing.T) {
		Tanh(800)
		Tanh(-800)
		Tanh(math.Inf(1))
		Tanh(math.Inf(-1))
		Tanh(math.NaN())
	})

	t.Run("Asinh_Edge", func(t *testing.T) {
		Asinh(1e200)
		Asinh(-1e200)
		Asinh(math.MaxFloat64)
		Asinh(-math.MaxFloat64)
		Asinh(math.Inf(1))
		Asinh(math.Inf(-1))
		Asinh(math.NaN())
		Asinh(0)
	})

	t.Run("Atanh_Edge", func(t *testing.T) {
		Atanh(1)
		Atanh(-1)
		Atanh(2)
		Atanh(-2)
		Atanh(0.99999999999999999)
		Atanh(-0.99999999999999999)
		Atanh(math.NaN())
		Atanh(0)
	})
}
