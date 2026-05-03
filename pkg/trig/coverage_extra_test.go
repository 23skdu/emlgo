package trig

import (
	"math"
	"testing"
)

func TestInverseTrigExtra(t *testing.T) {
	t.Run("Sec", func(t *testing.T) {
		_ = Sec(0)
		_ = Sec(math.Pi / 6)
		_ = Sec(math.Pi / 4)
		_ = Sec(math.Pi / 3)
		_ = Sec(math.Pi / 2)
	})
	t.Run("Csc", func(t *testing.T) {
		_ = Csc(0)
		_ = Csc(math.Pi / 6)
		_ = Csc(math.Pi / 4)
		_ = Csc(math.Pi / 3)
		_ = Csc(math.Pi / 2)
	})
	t.Run("Cot", func(t *testing.T) {
		_ = Cot(0)
		_ = Cot(math.Pi / 6)
		_ = Cot(math.Pi / 4)
		_ = Cot(math.Pi / 3)
		_ = Cot(math.Pi / 2)
	})
}

func TestInverseSecCsc(t *testing.T) {
	t.Run("Asec", func(t *testing.T) {
		_ = Asec(0)
		_ = Asec(1)
		_ = Asec(2)
		_ = Asec(-1)
		_ = Asec(-2)
	})
	t.Run("Acsc", func(t *testing.T) {
		_ = Acsc(0)
		_ = Acsc(1)
		_ = Acsc(2)
		_ = Acsc(-1)
		_ = Acsc(-2)
	})
	t.Run("Acot", func(t *testing.T) {
		_ = Acot(0)
		_ = Acot(1)
		_ = Acot(2)
		_ = Acot(-1)
		_ = Acot(-2)
	})
}

func TestSecCscLarge(t *testing.T) {
	values := []float64{0, math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2, math.Pi, 2 * math.Pi}
	for _, x := range values {
		_ = Sec(x)
		_ = Csc(x)
		_ = Cot(x)
	}
}

func TestAsinhAcoshAtanh(t *testing.T) {
	values := []float64{-1000, -100, -10, -1, -0.5, 0, 0.5, 1, 10, 100, 1000}
	for _, x := range values {
		_ = Asinh(x)
		_ = Atanh(x * 0.9)
		_ = Acosh(x + 1.1)
	}
}

func TestHyperSecCsch(t *testing.T) {
	values := []float64{-100, -10, -1, 0, 1, 10, 100}
	for _, x := range values {
		_ = Sech(x)
		_ = Csch(x)
		_ = Coth(x)
	}
}

func TestInverseHyper(t *testing.T) {
	values := []float64{-100, -10, -1, 0, 0.5, 1, 10, 100}
	for _, x := range values {
		_ = Asech(x + 1.1)
		_ = Acsch(x)
		_ = Acoth(x + 1)
	}
}

func TestSinCosCombined(t *testing.T) {
	values := []float64{0, math.Pi / 6, math.Pi / 4, math.Pi / 3, math.Pi / 2, math.Pi}
	for _, x := range values {
		s, c := SinCos(x)
		_ = s
		_ = c
	}
}

func TestSinhCoshTanh(t *testing.T) {
	values := []float64{-100, -10, -1, 0, 1, 10, 100}
	for _, x := range values {
		s, c := SinhCosh(x)
		_ = s
		_ = c
		_ = Tanh(x)
	}
}

func TestTanFastLarge(t *testing.T) {
	values := []float64{-100, -10, -1, 0, 1, 10, 100}
	for _, x := range values {
		_ = TanFast(x)
	}
}