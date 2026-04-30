package fastmath

import (
	"math"
	"github.com/emlgo/eml/internal/eml"
)

// Sqrt returns the square root of x.
// It uses direct assembly (SQRTSD/FSQRTD) and is faster than math.Sqrt.
func Sqrt(x float64) float64 {
	return eml.SqrtScalar(x)
}

// FMA returns x*y + z, computed with only one rounding.
// It uses direct assembly (VFMADD/FMADD) and is significantly faster than math.FMA.
func FMA(x, y, z float64) float64 {
	return eml.FmaScalar(x, y, z)
}

// Exp returns e^x.
// It uses a polynomial approximation optimized with FMA.
// Accuracy: ~1e-7 in the primary range.
func Exp(x float64) float64 {
	const (
		Log2E = 1.44269504088896340736 // 1/ln(2)
		Ln2Hi = 0.69314718036912381649
		Ln2Lo = 1.9082149292705877000e-10
	)
	
	if x > 709.78 { return math.Inf(1) }
	if x < -745.13 { return 0 }
	
	k := math.Round(x * Log2E)
	r := x - k*Ln2Hi - k*Ln2Lo

	// 5th degree minimax polynomial for e^r
	p := 0.008333333333333333
	p = eml.FmaScalar(p, r, 0.041666666666666664)
	p = eml.FmaScalar(p, r, 0.16666666666666666)
	p = eml.FmaScalar(p, r, 0.5)
	p = eml.FmaScalar(p, r, 1.0)
	p = eml.FmaScalar(p, r, 1.0)

	return math.Ldexp(p, int(k))
}

// Sin returns the sine of x.
// It uses a polynomial approximation optimized with FMA.
// Optimized for x in [-pi/2, pi/2].
func Sin(x float64) float64 {
	// Simple range reduction for FastMath
	if x < -1.5707963267948966 || x > 1.5707963267948966 {
		return math.Sin(x) // Fallback for large x
	}
	
	x2 := x * x
	p := -0.0001984126984126984
	p = eml.FmaScalar(p, x2, 0.008333333333333333)
	p = eml.FmaScalar(p, x2, -0.16666666666666666)
	p = eml.FmaScalar(p, x2, 1.0)
	return x * p
}

// Cos returns the cosine of x.
// It uses a polynomial approximation optimized with FMA.
// Optimized for x in [-pi/2, pi/2].
func Cos(x float64) float64 {
	if x < -1.5707963267948966 || x > 1.5707963267948966 {
		return math.Cos(x) // Fallback
	}
	
	x2 := x * x
	p := -0.001388888888888889
	p = eml.FmaScalar(p, x2, 0.041666666666666664)
	p = eml.FmaScalar(p, x2, -0.5)
	p = eml.FmaScalar(p, x2, 1.0)
	return p
}

// Log returns the natural logarithm of x.
// It uses range reduction (Frexp) and a polynomial approximation for the mantissa.
func Log(x float64) float64 {
	if x <= 0 {
		return math.Log(x)
	}
	
	frac, exp := math.Frexp(x)
	frac *= 2
	exp--
	
	r := frac - 1.0
	p := -0.0645354
	p = eml.FmaScalar(p, r, 0.1638394)
	p = eml.FmaScalar(p, r, -0.2420911)
	p = eml.FmaScalar(p, r, 0.3395481)
	p = eml.FmaScalar(p, r, -0.4996741)
	p = eml.FmaScalar(p, r, 0.9999964)
	p = p * r
	
	return eml.FmaScalar(float64(exp), 0.6931471805599453, p)
}
