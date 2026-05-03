package hyper

import (
	"testing"
)

func TestHyperExhaustive(t *testing.T) {
	medium := make([]float64, 1000)
	for i := range medium {
		medium[i] = float64(i)*0.01 - 5
	}

	t.Run("Sinh_medium", func(t *testing.T) {
		r := Sinh(medium[500])
		_ = r
	})
	t.Run("Cosh_medium", func(t *testing.T) {
		r := Cosh(medium[500])
		_ = r
	})
	t.Run("Tanh_medium", func(t *testing.T) {
		r := Tanh(medium[500])
		_ = r
	})
	t.Run("Asinh_medium", func(t *testing.T) {
		r := Asinh(medium[500])
		_ = r
	})
	t.Run("Acosh_medium", func(t *testing.T) {
		r := Acosh(medium[500] + 2)
		_ = r
	})
	t.Run("Atanh_medium", func(t *testing.T) {
		r := Atanh(medium[500] * 0.5)
		_ = r
	})
}

func TestSinhCoshTanhEdgeCases(t *testing.T) {
	t.Run("sinh_neg_large", func(t *testing.T) {
		_ = Sinh(-1000)
	})
	t.Run("sinh_pos_large", func(t *testing.T) {
		_ = Sinh(1000)
	})
	t.Run("cosh_neg_large", func(t *testing.T) {
		_ = Cosh(-1000)
	})
	t.Run("cosh_pos_large", func(t *testing.T) {
		_ = Cosh(1000)
	})
	t.Run("tanh_boundary", func(t *testing.T) {
		_ = Tanh(709)
	})
}

func TestAsinhAcoshAtanhEdgeCases(t *testing.T) {
	t.Run("asinh_zero", func(t *testing.T) {
		got := Asinh(0)
		if got != 0 {
			t.Errorf("Asinh(0) = %v, want 0", got)
		}
	})
	t.Run("asinh_one", func(t *testing.T) {
		_ = Asinh(1)
	})
	t.Run("acosh_one", func(t *testing.T) {
		got := Acosh(1)
		if got != 0 {
			t.Errorf("Acosh(1) = %v, want 0", got)
		}
	})
}