package eml

import (
	"math"
)

// Standard mathematical constants.
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


// IsNaN reports whether f is "not-a-number".
func IsNaN(f float64) bool {
	return math.IsNaN(f)
}


// IsInf reports whether f is an infinity, according to sign.
func IsInf(f float64, sign int) bool {
	return math.IsInf(f, sign)
}


// IsFinite reports whether f is neither NaN nor Inf.
func IsFinite(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}


// Signbit reports whether f is negative or negative zero.
func Signbit(f float64) bool {
	return signbit(f)
}


// NaN returns an IEEE 754 "not-a-number" value.
func NaN() float64 {
	return math.NaN()
}


// Inf returns an infinity value, according to sign.
func Inf(sign int) float64 {
	return math.Inf(sign)
}


// Floor returns the greatest integer value less than or equal to x.
func Floor(x float64) float64 {
	return math.Floor(x)
}


// Ceil returns the least integer value greater than or equal to x.
func Ceil(x float64) float64 {
	return math.Ceil(x)
}


// Trunc returns the integer value of x.
func Trunc(x float64) float64 {
	return math.Trunc(x)
}


// Round returns the nearest integer, rounding half away from zero.
func Round(x float64) float64 {
	return math.Round(x)
}


// Abs returns the absolute value of x.
func Abs(x float64) float64 {
	return math.Abs(x)
}


// Neg returns the negation of x.
func Neg(x float64) float64 {
	return -x
}


// Inv returns the inverse of x (1/x).
func Inv(x float64) float64 {
	return 1 / x
}


// Exp returns e**x, the base-e exponential of x.
func Exp(x float64) float64 {
	return math.Exp(x)
}


// Log returns the natural logarithm of x.
func Log(x float64) float64 {
	return math.Log(x)
}


// Log1p returns the natural logarithm of 1 plus its argument x.
func Log1p(x float64) float64 {
	return math.Log1p(x)
}


// Expm1 returns e**x - 1.
func Expm1(x float64) float64 {
	return math.Expm1(x)
}


// Sqrt returns the square root of x.
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}


// Cbrt returns the cube root of x.
func Cbrt(x float64) float64 {
	return math.Cbrt(x)
}


// Pow returns x**y.
func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}


// PowInt returns x**n.
func PowInt(x float64, n int) float64 {
	return math.Pow(x, float64(n))
}


// Sin returns the sine of the radian argument x.
func Sin(x float64) float64 {
	return math.Sin(x)
}


// Cos returns the cosine of the radian argument x.
func Cos(x float64) float64 {
	return math.Cos(x)
}


// Tan returns the tangent of the radian argument x.
func Tan(x float64) float64 {
	return math.Tan(x)
}


// Sincos returns Sin(x), Cos(x).
func Sincos(x float64) (sin, cos float64) {
	return math.Sincos(x)
}


// Asin returns the arcsine of x in radians.
func Asin(x float64) float64 {
	return math.Asin(x)
}


// Acos returns the arccosine of x in radians.
func Acos(x float64) float64 {
	return math.Acos(x)
}


// Atan returns the arctangent of x in radians.
func Atan(x float64) float64 {
	return math.Atan(x)
}


// Atan2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
func Atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}


// Sinh returns the hyperbolic sine of x.
func Sinh(x float64) float64 {
	return math.Sinh(x)
}


// Cosh returns the hyperbolic cosine of x.
func Cosh(x float64) float64 {
	return math.Cosh(x)
}


// Tanh returns the hyperbolic tangent of x.
func Tanh(x float64) float64 {
	return math.Tanh(x)
}


// Asinh returns the inverse hyperbolic sine of x.
func Asinh(x float64) float64 {
	return math.Asinh(x)
}


// Acosh returns the inverse hyperbolic cosine of x.
func Acosh(x float64) float64 {
	return math.Acosh(x)
}


// Atanh returns the inverse hyperbolic tangent of x.
func Atanh(x float64) float64 {
	return math.Atanh(x)
}


// Hypot returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow.
func Hypot(x, y float64) float64 {
	return math.Hypot(x, y)
}


// Max returns the larger of x or y.
func Max(x, y float64) float64 {
	return nativeMax(x, y)
}


// Min returns the smaller of x or y.
func Min(x, y float64) float64 {
	return nativeMin(x, y)
}


// Mod returns the floating-point remainder of x/y. The magnitude of the result is less than y and its sign agrees with that of x.
func Mod(x, y float64) float64 {
	return math.Mod(x, y)
}


// Remainder returns the IEEE 754 floating-point remainder of x/y.
func Remainder(x, y float64) float64 {
	return math.Remainder(x, y)
}


// Log10 returns the decimal logarithm of x.
func Log10(x float64) float64 {
	return math.Log10(x)
}


// Copysign returns a value with the magnitude of x and the sign of y.
func Copysign(x, y float64) float64 {
	return math.Copysign(x, y)
}


// Modf returns integer and fractional floating-point parts of x.
func Modf(x float64) (intPart, fracPart float64) {
	return math.Modf(x)
}

