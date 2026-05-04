package eml

import (
	"math"
)

// Bitwise helpers for IEEE 754 floating point
func f64bits(f float64) uint64     { return math.Float64bits(f) }
func f64frombits(b uint64) float64 { return math.Float64frombits(b) }

func signbit(x float64) bool {
	return f64bits(x)&(1<<63) != 0
}

//go:inline
func nan() float64 {
	return math.NaN()
}

//go:inline
func inf(sign int) float64 {
	return math.Inf(sign)
}

//go:inline
func isNaN(f float64) bool {
	return math.IsNaN(f)
}

//go:inline
func isInf(f float64, sign int) bool {
	return math.IsInf(f, sign)
}

func nativeSqrt(x float64) float64 {
	return math.Sqrt(x)
}

func nativeExp(x float64) float64 {
	return math.Exp(x)
}

func nativeLog(x float64) float64 {
	return math.Log(x)
}

func nativeSin(x float64) float64 {
	return math.Sin(x)
}

func nativeCos(x float64) float64 {
	return math.Cos(x)
}

func nativeLog2(x float64) float64 {
	return math.Log2(x)
}

func nativeLog10(x float64) float64 {
	return math.Log10(x)
}

func nativeSincos(x float64) (sin, cos float64) {
	return math.Sincos(x)
}

func nativeExpm1(x float64) float64 {
	return math.Expm1(x)
}

func nativeTan(x float64) float64 {
	return math.Tan(x)
}

func nativeAtan(x float64) float64 {
	return math.Atan(x)
}

func nativeAtan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

func nativeAsin(x float64) float64 {
	return math.Asin(x)
}

func nativeAcos(x float64) float64 {
	return math.Acos(x)
}

func nativeAsinh(x float64) float64 {
	return math.Asinh(x)
}

func nativeAcosh(x float64) float64 {
	return math.Acosh(x)
}

func nativeAtanh(x float64) float64 {
	return math.Atanh(x)
}

func nativeLog1p(x float64) float64 {
	return math.Log1p(x)
}

func nativePow(x, y float64) float64 {
	return math.Pow(x, y)
}

func nativeInv(x float64) float64 {
	return 1 / x
}

func nativeNeg(x float64) float64 {
	return -x
}

func nativeAbs(x float64) float64 {
	return math.Abs(x)
}

func nativeMod(x, y float64) float64 {
	return math.Mod(x, y)
}

func nativeRemainder(x, y float64) float64 {
	return math.Remainder(x, y)
}

func nativeHypot(x, y float64) float64 {
	return math.Hypot(x, y)
}

func nativeCbrt(x float64) float64 {
	return math.Cbrt(x)
}

func nativeMax(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x > y {
		return x
	}
	return y
}

func nativeMin(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x < y {
		return x
	}
	return y
}


func copysign(x, y float64) float64 {
	return math.Copysign(x, y)
}

func floor(x float64) float64 {
	return math.Floor(x)
}

func ceil(x float64) float64 {
	return math.Ceil(x)
}

func trunc(x float64) float64 {
	return math.Trunc(x)
}

func round(x float64) float64 {
	return math.Round(x)
}
