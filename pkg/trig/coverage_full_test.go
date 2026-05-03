package trig

import (
	"math"
	"testing"
)

func TestSecFullCoverage(t *testing.T) {
	_ = Sec(0)
	_ = Sec(math.Pi/6)
	_ = Sec(math.Pi/4)
	_ = Sec(math.Pi/3)
	_ = Sec(math.Pi/2)
	_ = Sec(math.Pi)
	_ = Sec(2*math.Pi)
	_ = Sec(-math.Pi)
	_ = Sec(math.NaN())
	_ = Sec(math.Inf(1))
	_ = Sec(math.Inf(-1))
}

func TestCscFullCoverage(t *testing.T) {
	_ = Csc(0)
	_ = Csc(math.Pi/6)
	_ = Csc(math.Pi/4)
	_ = Csc(math.Pi/3)
	_ = Csc(math.Pi/2)
	_ = Csc(math.Pi)
	_ = Csc(2*math.Pi)
	_ = Csc(-math.Pi)
	_ = Csc(math.NaN())
	_ = Csc(math.Inf(1))
	_ = Csc(math.Inf(-1))
}

func TestCotFullCoverage(t *testing.T) {
	_ = Cot(0)
	_ = Cot(math.Pi/6)
	_ = Cot(math.Pi/4)
	_ = Cot(math.Pi/3)
	_ = Cot(math.Pi/2)
	_ = Cot(math.Pi)
	_ = Cot(-math.Pi)
	_ = Cot(math.NaN())
	_ = Cot(math.Inf(1))
	_ = Cot(math.Inf(-1))
}

func TestAcotFullCoverage(t *testing.T) {
	_ = Acot(0)
	_ = Acot(1)
	_ = Acot(-1)
	_ = Acot(2)
	_ = Acot(-2)
	_ = Acot(0.5)
	_ = Acot(-0.5)
	_ = Acot(math.MaxFloat64)
	_ = Acot(math.SmallestNonzeroFloat64)
	_ = Acot(math.NaN())
	_ = Acot(math.Inf(1))
	_ = Acot(math.Inf(-1))
}

func TestAsecFullCoverage(t *testing.T) {
	_ = Asec(0)
	_ = Asec(1)
	_ = Asec(-1)
	_ = Asec(2)
	_ = Asec(-2)
	_ = Asec(0.5)
	_ = Asec(-0.5)
	_ = Asec(math.NaN())
	_ = Asec(math.Inf(1))
	_ = Asec(math.Inf(-1))
}

func TestAcscFullCoverage(t *testing.T) {
	_ = Acsc(0)
	_ = Acsc(1)
	_ = Acsc(-1)
	_ = Acsc(2)
	_ = Acsc(-2)
	_ = Acsc(0.5)
	_ = Acsc(-0.5)
	_ = Acsc(math.NaN())
	_ = Acsc(math.Inf(1))
	_ = Acsc(math.Inf(-1))
}

func TestCothFullCoverage(t *testing.T) {
	_ = Coth(0)
	_ = Coth(1)
	_ = Coth(-1)
	_ = Coth(0.5)
	_ = Coth(-0.5)
	_ = Coth(math.MaxFloat64)
	_ = Coth(math.SmallestNonzeroFloat64)
	_ = Coth(math.NaN())
	_ = Coth(math.Inf(1))
	_ = Coth(math.Inf(-1))
}

func TestSechFullCoverage(t *testing.T) {
	_ = Sech(0)
	_ = Sech(1)
	_ = Sech(-1)
	_ = Sech(0.5)
	_ = Sech(-0.5)
	_ = Sech(math.MaxFloat64)
	_ = Sech(math.SmallestNonzeroFloat64)
	_ = Sech(math.NaN())
	_ = Sech(math.Inf(1))
	_ = Sech(math.Inf(-1))
}

func TestCschFullCoverage(t *testing.T) {
	_ = Csch(0)
	_ = Csch(1)
	_ = Csch(-1)
	_ = Csch(0.5)
	_ = Csch(-0.5)
	_ = Csch(math.MaxFloat64)
	_ = Csch(math.SmallestNonzeroFloat64)
	_ = Csch(math.NaN())
	_ = Csch(math.Inf(1))
	_ = Csch(math.Inf(-1))
}

func TestAsinhFullCoverage(t *testing.T) {
	_ = Asinh(0)
	_ = Asinh(1)
	_ = Asinh(-1)
	_ = Asinh(0.5)
	_ = Asinh(-0.5)
	_ = Asinh(10)
	_ = Asinh(-10)
	_ = Asinh(100)
	_ = Asinh(-100)
	_ = Asinh(math.MaxFloat64)
	_ = Asinh(math.NaN())
	_ = Asinh(math.Inf(1))
	_ = Asinh(math.Inf(-1))
}

func TestAcoshFullCoverage(t *testing.T) {
	_ = Acosh(1)
	_ = Acosh(2)
	_ = Acosh(10)
	_ = Acosh(100)
	_ = Acosh(1.5)
	_ = Acosh(1.01)
	_ = Acosh(math.MaxFloat64)
	_ = Acosh(math.NaN())
	_ = Acosh(math.Inf(1))
	_ = Acosh(0.5)
	_ = Acosh(0)
}

func TestAtanhFullCoverage(t *testing.T) {
	_ = Atanh(0)
	_ = Atanh(0.5)
	_ = Atanh(-0.5)
	_ = Atanh(0.9)
	_ = Atanh(-0.9)
	_ = Atanh(0.99)
	_ = Atanh(-0.99)
	_ = Atanh(0.999)
	_ = Atanh(-0.999)
	_ = Atanh(math.NaN())
	_ = Atanh(1)
	_ = Atanh(-1)
	_ = Atanh(math.Inf(1))
	_ = Atanh(math.Inf(-1))
}

func TestAcothFullCoverage(t *testing.T) {
	_ = Acoth(0)
	_ = Acoth(1)
	_ = Acoth(2)
	_ = Acoth(-1)
	_ = Acoth(-2)
	_ = Acoth(0.5)
	_ = Acoth(-0.5)
	_ = Acoth(math.MaxFloat64)
	_ = Acoth(math.NaN())
	_ = Acoth(math.Inf(1))
	_ = Acoth(math.Inf(-1))
}

func TestAsechFullCoverage(t *testing.T) {
	_ = Asech(0)
	_ = Asech(1)
	_ = Asech(0.5)
	_ = Asech(0.1)
	_ = Asech(0.01)
	_ = Asech(math.NaN())
	_ = Asech(math.Inf(1))
	_ = Asech(math.Inf(-1))
	_ = Asech(-1)
	_ = Asech(2)
}

func TestAcschFullCoverage(t *testing.T) {
	_ = Acsch(0)
	_ = Acsch(1)
	_ = Acsch(-1)
	_ = Acsch(0.5)
	_ = Acsch(-0.5)
	_ = Acsch(10)
	_ = Acsch(-10)
	_ = Acsch(math.MaxFloat64)
	_ = Acsch(math.NaN())
	_ = Acsch(math.Inf(1))
	_ = Acsch(math.Inf(-1))
}

func TestSinCosAll(t *testing.T) {
	values := []float64{0, math.Pi/6, math.Pi/4, math.Pi/3, math.Pi/2, math.Pi, 3*math.Pi/2, 2*math.Pi, -math.Pi, -math.Pi/2, math.NaN(), math.Inf(1), math.Inf(-1)}
	for _, x := range values {
		_, _ = SinCos(x)
	}
}

func TestSinhCoshAll(t *testing.T) {
	values := []float64{0, 1, -1, 0.5, -0.5, 10, -10, 100, -100, math.MaxFloat64, math.NaN(), math.Inf(1), math.Inf(-1)}
	for _, x := range values {
		_, _ = SinhCosh(x)
	}
}

func TestTanFastAll(t *testing.T) {
	values := []float64{0, math.Pi/6, math.Pi/4, math.Pi/3, math.Pi/2, math.Pi, 3*math.Pi/2, -math.Pi, 10, -10, 100, -100, math.NaN(), math.Inf(1), math.Inf(-1)}
	for _, x := range values {
		_ = TanFast(x)
	}
}