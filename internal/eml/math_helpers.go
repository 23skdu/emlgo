package eml

import (
	"math"
)

const (
	Pi              = 3.141592653589793238462643383279
	PiOver2         = 1.570796326794896619231321691640
	PiOver4         = 0.785398163397448309615660845820
	E               = 2.718281828459045235360287471353
	Sqrt2           = 1.414213562373095048801688724210
	Ln2             = 0.693147180559945309417232121458
	Ln10            = 2.302585092994045684017991454684
	MaxFloat64      = 1.797693134862315708145274237317043544980e+308
	SmallestNonZero = 4.940656458412465441765687928082213877e-324
)

func IsNaN(f float64) bool {
	return math.IsNaN(f)
}

func IsInf(f float64, sign int) bool {
	return math.IsInf(f, sign)
}

func IsFinite(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}

func Signbit(f float64) bool {
	return signbit(f)
}

func NaN() float64 {
	return math.NaN()
}

func Inf(sign int) float64 {
	return math.Inf(sign)
}

func Floor(x float64) float64 {
	return math.Floor(x)
}

func Ceil(x float64) float64 {
	return math.Ceil(x)
}

func Trunc(x float64) float64 {
	return math.Trunc(x)
}

func Round(x float64) float64 {
	return math.Round(x)
}

func Abs(x float64) float64 {
	return math.Abs(x)
}

func Neg(x float64) float64 {
	return -x
}

func Inv(x float64) float64 {
	return 1 / x
}

func Exp(x float64) float64 {
	return math.Exp(x)
}

func Log(x float64) float64 {
	return math.Log(x)
}

func Log1p(x float64) float64 {
	return math.Log1p(x)
}

func Expm1(x float64) float64 {
	return math.Expm1(x)
}

func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

func Cbrt(x float64) float64 {
	return math.Cbrt(x)
}

func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}

func PowInt(x float64, n int) float64 {
	return math.Pow(x, float64(n))
}

func Sin(x float64) float64 {
	return math.Sin(x)
}

func Cos(x float64) float64 {
	return math.Cos(x)
}

func Tan(x float64) float64 {
	return math.Tan(x)
}

func Sincos(x float64) (sin, cos float64) {
	return math.Sincos(x)
}

func Asin(x float64) float64 {
	return math.Asin(x)
}

func Acos(x float64) float64 {
	return math.Acos(x)
}

func Atan(x float64) float64 {
	return math.Atan(x)
}

func Atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

func Sinh(x float64) float64 {
	return math.Sinh(x)
}

func Cosh(x float64) float64 {
	return math.Cosh(x)
}

func Tanh(x float64) float64 {
	return math.Tanh(x)
}

func Asinh(x float64) float64 {
	return math.Asinh(x)
}

func Acosh(x float64) float64 {
	return math.Acosh(x)
}

func Atanh(x float64) float64 {
	return math.Atanh(x)
}

func Hypot(x, y float64) float64 {
	return math.Hypot(x, y)
}

func Max(x, y float64) float64 {
	return nativeMax(x, y)
}

func Min(x, y float64) float64 {
	return nativeMin(x, y)
}

func Mod(x, y float64) float64 {
	return math.Mod(x, y)
}

func Remainder(x, y float64) float64 {
	return math.Remainder(x, y)
}

func Log10(x float64) float64 {
	return math.Log10(x)
}

func Copysign(x, y float64) float64 {
	return math.Copysign(x, y)
}

func Modf(x float64) (intPart, fracPart float64) {
	return math.Modf(x)
}
