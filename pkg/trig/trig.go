package trig

import (
	"math"

	"github.com/23skdu/emlgo/internal/eml"
	"github.com/23skdu/emlgo/pkg/hyper"
)

var (
	pi          = eml.Pi
	piOver2     = eml.PiOver2
	inf         = eml.Inf
	nan         = eml.NaN
	isNaN       = eml.IsNaN
	isInf       = eml.IsInf
	nativeSin   = eml.Sin
	nativeCos   = eml.Cos
	nativeTan   = eml.Tan
	nativeAsin  = eml.Asin
	nativeAcos  = eml.Acos
	nativeAtan  = eml.Atan
	nativeAtan2 = eml.Atan2
)

// Sin returns the sine of x in radians.
// It returns NaN for NaN or infinite input and preserves the sign of zero.
func Sin(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	if x == 0 {
		return x // preserve signed zero
	}
	return nativeSin(x)
}

// Cos returns the cosine of x in radians.
// It returns NaN for NaN or infinite input.
func Cos(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return nativeCos(x)
}

// Tan returns the tangent of x in radians.
// It delegates straight to the standard math implementation with no NaN/Inf pre-checks.
func Tan(x float64) float64 {
	return nativeTan(x)
}

// Cot returns the cotangent of x, computed as Cos(x) / Sin(x).
// It returns NaN for NaN or infinite input.
func Cot(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return Cos(x) / Sin(x)
}

// Sec returns the secant of x, computed as 1 / cos(x).
// It returns NaN for NaN or infinite input.
func Sec(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return 1 / nativeCos(x)
}

// Csc returns the cosecant of x, computed as 1 / Sin(x).
// It returns NaN for NaN or infinite input.
func Csc(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return 1 / Sin(x)
}

// Asin returns the arcsine of x in radians via the standard math implementation.
// Inputs with |x| > 1 return NaN.
func Asin(x float64) float64 {
	return nativeAsin(x)
}

// Acos returns the arccosine of x in radians via the standard math implementation.
// Inputs with |x| > 1 return NaN.
func Acos(x float64) float64 {
	return nativeAcos(x)
}

// Atan returns the arctangent of x in radians via the standard math implementation.
func Atan(x float64) float64 {
	return nativeAtan(x)
}

// Atan2 returns the arc tangent of y/x, using the signs of the two to determine
// the quadrant of the result. It delegates to the standard math implementation.
func Atan2(y, x float64) float64 {
	return nativeAtan2(y, x)
}

// Acot returns the arccotangent of x, computed as pi/2 - Atan(x), so the result
// lies in (0, pi). It returns NaN for NaN or infinite input.
func Acot(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return piOver2 - Atan(x)
}

// Asec returns the arcsecant of x, computed as Acos(1/x).
// It returns NaN for NaN input and for |x| < 1.
func Asec(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x > -1 && x < 1 {
		return nan()
	}
	return Acos(1 / x)
}

// Acsc returns the arccosecant of x, computed as Asin(1/x).
// It returns NaN for NaN input and for |x| < 1.
func Acsc(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x > -1 && x < 1 {
		return nan()
	}
	return Asin(1 / x)
}

// Sinh returns the hyperbolic sine of x, delegating to hyper.Sinh.
func Sinh(x float64) float64 {
	return hyper.Sinh(x)
}

// Cosh returns the hyperbolic cosine of x, delegating to hyper.Cosh.
func Cosh(x float64) float64 {
	return hyper.Cosh(x)
}

// Tanh returns the hyperbolic tangent of x, delegating to hyper.Tanh.
func Tanh(x float64) float64 {
	return hyper.Tanh(x)
}

// Coth returns the hyperbolic cotangent of x, computed as 1 / hyper.Tanh(x).
// It returns NaN for NaN input and for x == 0.
func Coth(x float64) float64 {
	if isNaN(x) || x == 0 {
		return nan()
	}
	return 1 / hyper.Tanh(x)
}

// Sech returns the hyperbolic secant of x, computed as 1 / hyper.Cosh(x).
// It returns NaN for NaN input; infinite input yields 0.
func Sech(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	return 1 / hyper.Cosh(x)
}

// Csch returns the hyperbolic cosecant of x, computed as 1 / hyper.Sinh(x).
// It returns NaN for NaN input and for x == 0.
func Csch(x float64) float64 {
	if isNaN(x) || x == 0 {
		return nan()
	}
	return 1 / hyper.Sinh(x)
}

// Asinh returns the inverse hyperbolic sine of x, delegating to hyper.Asinh.
func Asinh(x float64) float64 {
	return hyper.Asinh(x)
}

// Acosh returns the inverse hyperbolic cosine of x, delegating to hyper.Acosh.
func Acosh(x float64) float64 {
	return hyper.Acosh(x)
}

// Atanh returns the inverse hyperbolic tangent of x, delegating to hyper.Atanh.
func Atanh(x float64) float64 {
	return hyper.Atanh(x)
}

// Acoth returns the inverse hyperbolic cotangent of x, computed as Atanh(1/x).
// It returns NaN for NaN input and for |x| <= 1.
func Acoth(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x >= -1 && x <= 1 {
		return nan()
	}
	return Atanh(1 / x)
}

// Asech returns the inverse hyperbolic secant of x, computed as Acosh(1/x).
// It returns NaN for NaN input and for x <= 0 or x > 1; Asech(1) is 0.
func Asech(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x <= 0 || x > 1 {
		return nan()
	}
	if x == 1 {
		return 0
	}
	return Acosh(1 / x)
}

// Acsch returns the inverse hyperbolic cosecant of x, computed as Asinh(1/x).
// It returns NaN for NaN input, returns signed infinity at x == 0,
// and returns 0 for infinite input.
func Acsch(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x == 0 {
		// acsch(x) → sign(x) * Inf as x → 0
		return math.Copysign(inf(1), x)
	}
	if isInf(x, 0) {
		return 0
	}
	return Asinh(1 / x)
}

// DegToRad converts an angle in degrees to radians.
func DegToRad(deg float64) float64 {
	return deg * pi / 180
}

// RadToDeg converts an angle in radians to degrees.
func RadToDeg(rad float64) float64 {
	return rad * 180 / pi
}

// SinCos returns sin(x) and cos(x) from a single math.Sincos call.
func SinCos(x float64) (sin, cos float64) {
	return math.Sincos(x)
}

// SinhCosh returns sinh(x) and cosh(x) computed with math.Sinh and math.Cosh,
// bypassing the overflow-protected hyper.Sinh and hyper.Cosh.
func SinhCosh(x float64) (sinh, cosh float64) {
	return math.Sinh(x), math.Cosh(x)
}

// SinBatch returns a new slice containing the sine of each element of x,
// computed through SIMD dispatch (eml.SinSIMD).
func SinBatch(x []float64) []float64 {
	return eml.SinSIMD(x)
}

// CosBatch returns a new slice containing the cosine of each element of x,
// computed through SIMD dispatch (eml.CosSIMD).
func CosBatch(x []float64) []float64 {
	return eml.CosSIMD(x)
}

// SinCosBatch returns two new slices containing the sine and the cosine
// of each element of x, computed through SIMD dispatch (eml.SinCosSIMD).
func SinCosBatch(x []float64) (sin, cos []float64) {
	return eml.SinCosSIMD(x)
}

// TanBatch returns a new slice containing the tangent of each element of x,
// computed through SIMD dispatch (eml.TanSIMD).
func TanBatch(x []float64) []float64 {
	return eml.TanSIMD(x)
}

// SinFast returns the sine of x computed directly with math.Sin.
// It returns NaN for NaN or infinite input and preserves the sign of zero.
func SinFast(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	if x == 0 {
		return x // preserve signed zero
	}
	return math.Sin(x)
}

// CosFast returns the cosine of x computed directly with math.Cos.
// It returns NaN for NaN or infinite input.
func CosFast(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return math.Cos(x)
}

// TanFast returns tan(x) as sin(x)/cos(x) from a single math.Sincos call,
// with no NaN/Inf pre-checks.
func TanFast(x float64) float64 {
	s, c := math.Sincos(x)
	return s / c
}
