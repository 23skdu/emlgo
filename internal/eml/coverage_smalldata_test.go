package eml

import (
	"testing"
)

func TestAllScalarMath(t *testing.T) {
	values := []float64{-2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2}
	for _, x := range values {
		_ = nativeTan(x)
		_ = nativeInv(x)
		_ = nativeNeg(x)
		_ = nativeAbs(x)
		_ = nativeLog10(x)
	}
}

func TestAllNativeOps(t *testing.T) {
	pairs := [][2]float64{{1, 2}, {2, 1}, {0.5, 0.3}, {-1, 2}, {1, -2}}
	for _, p := range pairs {
		_ = copysign(p[0], p[1])
		_ = nativeMod(p[0], p[1])
		_ = nativeRemainder(p[0], p[1])
		_ = nativeHypot(p[0], p[1])
		_ = nativeMax(p[0], p[1])
		_ = nativeMin(p[0], p[1])
	}
}

func TestSmallSIMD(t *testing.T) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7} {
		x := make([]float64, n)
		y := make([]float64, n)
		for i := range x {
			x[i] = float64(i+1) * 0.5
			y[i] = float64(i+1) * 0.3
		}
		_ = AddSIMD(x, y)
		_ = SubSIMD(x, y)
		_ = MulSIMD(x, y)
		_ = DivSIMD(x, y)
		_ = AddScalarSIMD(x, 1.5)
		_ = MulScalarSIMD(x, 1.5)
	}
}

func TestSmallTrigSIMD(t *testing.T) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7} {
		x := make([]float64, n)
		r := make([]float64, n)
		for i := range x {
			x[i] = float64(i+1) * 0.2
		}
		ExpSIMDTo(x, r)
		LogSIMDTo(x, r)
		SinSIMDTo(x, r)
		CosSIMDTo(x, r)
		TanSIMDTo(x, r)
		SqrtSIMDTo(x, r)
		_ = SqrtSIMD(x)
	}
}

func TestSmallSinCos(t *testing.T) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7} {
		x := make([]float64, n)
		s := make([]float64, n)
		c := make([]float64, n)
		for i := range x {
			x[i] = float64(i+1) * 0.2
		}
		SinCosSIMDTo(x, s, c)
	}
}

func TestSmallFused(t *testing.T) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7} {
		a := make([]float64, n)
		b := make([]float64, n)
		r := make([]float64, n)
		for i := range a {
			a[i] = float64(i+1) * 0.1
			b[i] = float64(i+1) * 0.2
		}
		ExpMulTo(a, b, r)
		ExpAddTo(a, b, r)
		LogDivTo(a, b, r)
		LogSubTo(a, b, r)
	}
}