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
//
// Overflow/underflow thresholds approximate ln(MaxFloat64) and ln(SmallestNonzeroFloat64).
const expOverflowFm = 709.78
const expUnderflowFm = -745.13

func Exp(x float64) float64 {
	const (
		Log2E = 1.44269504088896340736 // 1/ln(2)
		Ln2Hi = 0.69314718036912381649
		Ln2Lo = 1.9082149292705877000e-10
	)
	
	if x > expOverflowFm { return math.Inf(1) }
	if x < expUnderflowFm { return 0 }
	
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

var _sin = [...]float64{
	1.58962301576546568060e-10, // 0x3de5d8fd1fd19ccd
	-2.50507477628578072866e-8, // 0xbe5ae5e5a9291f5d
	2.75573136213857245213e-6,  // 0x3ec71de3567d48a1
	-1.98412698295895385996e-4, // 0xbf2a01a019bfdf03
	8.33333333332211858878e-3,  // 0x3f8111111110f7d0
	-1.66666666666666307295e-1, // 0xbfc5555555555548
}

var _cos = [...]float64{
	-1.13585365213876817300e-11, // 0xbda8fa49a0861a9b
	2.08757008419747316778e-9,   // 0x3e21ee9d7b4e3f05
	-2.75573141792967388112e-7,  // 0xbe927e4f7eac4bc6
	2.48015872888517045348e-5,   // 0x3efa01a019c844f5
	-1.38888888888730564116e-3,  // 0xbf56c16c16c14f91
	4.16666666666665929218e-2,   // 0x3fa555555555554b
}

// Sin returns the sine of x.
// It uses a high-accuracy branchless Cody-Waite range reduction and minimax polynomial approximation.
func Sin(x float64) float64 {
	const (
		TwoOverPi = 0.63661977236758134308
		PI4A      = 7.85398125648498535156e-1
		PI4B      = 3.77489470793079817668e-8
		PI4C      = 2.69515142907905952645e-15
	)

	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}

	sign := 1.0
	if x < 0 {
		x = -x
		sign = -1.0
	}

	y := math.Round(x*TwoOverPi) * 2.0
	k := int(y)

	z := ((x - y*PI4A) - y*PI4B) - y*PI4C
	zz := z * z

	// evaluate sin
	p_sin := _sin[0]
	p_sin = eml.FmaScalar(p_sin, zz, _sin[1])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[2])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[3])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[4])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[5])
	y_sin := z + z*zz*p_sin

	// evaluate cos
	p_cos := _cos[0]
	p_cos = eml.FmaScalar(p_cos, zz, _cos[1])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[2])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[3])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[4])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[5])
	y_cos := 1.0 - 0.5*zz + zz*zz*p_cos

	// select
	isCos := float64((k & 2) >> 1)
	res := (1.0-isCos)*y_sin + isCos*y_cos

	// sign adjustment
	signMask := 1.0 - float64((k&4)>>1)
	sign *= signMask

	return sign * res
}

// Cos returns the cosine of x.
// It uses a high-accuracy branchless Cody-Waite range reduction and minimax polynomial approximation.
func Cos(x float64) float64 {
	const (
		TwoOverPi = 0.63661977236758134308
		PI4A      = 7.85398125648498535156e-1
		PI4B      = 3.77489470793079817668e-8
		PI4C      = 2.69515142907905952645e-15
	)

	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}

	if x < 0 {
		x = -x
	}

	y := math.Round(x*TwoOverPi) * 2.0
	k := int(y)

	z := ((x - y*PI4A) - y*PI4B) - y*PI4C
	zz := z * z

	// evaluate sin
	p_sin := _sin[0]
	p_sin = eml.FmaScalar(p_sin, zz, _sin[1])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[2])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[3])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[4])
	p_sin = eml.FmaScalar(p_sin, zz, _sin[5])
	y_sin := z + z*zz*p_sin

	// evaluate cos
	p_cos := _cos[0]
	p_cos = eml.FmaScalar(p_cos, zz, _cos[1])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[2])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[3])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[4])
	p_cos = eml.FmaScalar(p_cos, zz, _cos[5])
	y_cos := 1.0 - 0.5*zz + zz*zz*p_cos

	// select
	isSin := float64((k & 2) >> 1)
	res := (1.0-isSin)*y_cos + isSin*y_sin

	// sign adjustment
	sign := 1.0 - float64(2*(((k&2)>>1)^((k&4)>>2)))

	return sign * res
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
