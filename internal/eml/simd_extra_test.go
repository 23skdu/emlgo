package eml

import (
	"testing"
)

func TestExpMulAddBatch(t *testing.T) {
	for _, n := range []int{1, 8, 16, 64} {
		a := make([]float64, n)
		b := make([]float64, n)
		r := make([]float64, n)
		for i := range a {
			a[i] = float64(i) * 0.1
			b[i] = float64(i) * 0.2
		}
		ExpMulTo(a, b, r)
		ExpAddTo(a, b, r)
		LogDivTo(a, b, r)
		LogSubTo(a, b, r)
	}
}

func TestHyperbolic(t *testing.T) {
	for _, x := range []float64{-1, 0, 0.5, 1} {
		_ = Sinh(x)
		_ = Cosh(x)
		_ = Tanh(x)
	}
}

func TestNativePair(t *testing.T) {
	for _, p := range [][2]float64{{1, 2}, {-1, 1}} {
		_ = copysign(p[0], p[1])
		_ = nativeMod(p[0], p[1])
		_ = nativeRemainder(p[0], p[1])
		_ = nativeHypot(p[0], p[1])
	}
}

func TestNativeUn(t *testing.T) {
	for _, x := range []float64{0, 1, -1} {
		_ = nativeTan(x)
		_ = nativeInv(x)
		_ = nativeNeg(x)
		_ = nativeAbs(x)
		_ = nativeLog10(x)
	}
}

func TestModf2(t *testing.T) {
	for _, x := range []float64{0, 1.5, -1.5} {
		i, f := Modf(x)
		_ = i
		_ = f
	}
}

func TestSIMDops2(t *testing.T) {
	for _, n := range []int{1, 8, 16, 64} {
		x := make([]float64, n)
		y := make([]float64, n)
		for i := range x {
			x[i] = float64(i) * 0.1
			y[i] = float64(i) * 0.2
		}
		_ = AddSIMD(x, y)
		_ = SubSIMD(x, y)
		_ = MulSIMD(x, y)
		_ = DivSIMD(x, y)
		_ = AddScalarSIMD(x, 3.14)
		_ = MulScalarSIMD(x, 3.14)
	}
}

func TestSqrtScalar2(t *testing.T) {
	_ = SqrtScalar(0)
	_ = SqrtScalar(1)
	_ = SqrtScalar(100)
}

func TestDetect(t *testing.T) {
	detectSIMD()
}

func TestSinhCoshBatch2(t *testing.T) {
	for _, n := range []int{1, 8, 16, 64} {
		x := make([]float64, n)
		for i := range x {
			x[i] = float64(i) * 0.01
		}
		_ = SinhBatch(x)
		_ = CoshBatch(x)
	}
}

func TestSIMD2(t *testing.T) {
	for _, n := range []int{1, 8, 16, 64} {
		x := make([]float64, n)
		y := make([]float64, n)
		out := make([]float64, n)
		SIMD(x, y, out)
	}
}


func TestVariousSIMD2(t *testing.T) {
	for _, n := range []int{1, 8, 16, 64} {
		x := make([]float64, n)
		r := make([]float64, n)
		for i := range x {
			x[i] = float64(i) * 0.1
		}
		ExpSIMDTo(x, r)
		LogSIMDTo(x, r)
		SinSIMDTo(x, r)
		CosSIMDTo(x, r)
		TanSIMDTo(x, r)
		SqrtSIMDTo(x, r)
		s := make([]float64, n)
		c := make([]float64, n)
		SinCosSIMDTo(x, s, c)
		_ = SqrtSIMD(x)
	}
}