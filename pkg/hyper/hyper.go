package hyper

import (
	"math"

	"github.com/23skdu/emlgo/internal/constants"
	"github.com/23skdu/emlgo/internal/eml"
	"github.com/23skdu/emlgo/pkg/arithmetic"
)

var (
	isNaN     = eml.IsNaN
	isInf     = eml.IsInf
	inf       = eml.Inf
	nan       = eml.NaN
	nativeExp = eml.Exp
	nativeLog = eml.Log
	nativeAbs = eml.Abs
)

// Overflow threshold for exp-based hyperbolic functions:
// ln(MaxFloat64) ≈ 709.782712893384 — exp(x) overflows to +Inf beyond this.
const expOverflow = constants.ExpOverflow
const sinhOverflow = 710.4758600739439

// Sinh returns the hyperbolic sine of x with explicit overflow protection.
// NaN, infinite, and zero inputs are returned unchanged, and |x| beyond the
// sinh overflow threshold returns +/-Inf; small |x| avoids cancellation via expm1.
func Sinh(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}
	sign := 1.0
	if x < 0 {
		x = -x
		sign = -1.0
	}
	if x > sinhOverflow {
		return sign * inf(1)
	}
	var res float64
	if x < 0x1p-28 {
		res = x
	} else if x < 0.5 {
		// Avoid catastrophic cancellation: (expm1(x) - expm1(-x)) / 2
		res = (math.Expm1(x) - math.Expm1(-x)) / 2
	} else if x > 22.0 {
		eh := nativeExp(x * 0.5)
		res = (0.5 * eh) * eh
	} else {
		ex := nativeExp(x)
		res = 0.5 * (ex - 1.0/ex)
	}
	return sign * res
}

// Cosh returns the hyperbolic cosine of x with explicit overflow protection.
// NaN input is returned unchanged; infinite input and |x| beyond the sinh
// overflow threshold return +Inf. |x| > 22 is computed as (exp(x/2))^2 / 2.
func Cosh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if isInf(x, 0) {
		return inf(1)
	}
	absX := nativeAbs(x)
	if absX > sinhOverflow {
		return inf(1)
	}
	if absX > 22.0 {
		eh := nativeExp(absX * 0.5)
		return (0.5 * eh) * eh
	}
	ex := nativeExp(absX)
	return 0.5 * (ex + 1.0/ex)
}

// Tanh returns the hyperbolic tangent of x.
// NaN and zero inputs are returned unchanged and |x| > 20 saturates to +/-1;
// |x| < 2^-28 returns x, otherwise expm1(2x) / (expm1(2x) + 2) is evaluated.
func Tanh(x float64) float64 {
	if isNaN(x) || x == 0 {
		return x
	}
	if isInf(x, 1) || x >= expOverflow || x > 20.0 {
		return 1.0
	}
	if isInf(x, -1) || x <= -expOverflow || x < -20.0 {
		return -1.0
	}
	if nativeAbs(x) < 0x1p-28 {
		return x
	}
	e2x := math.Expm1(2 * x)
	return e2x / (e2x + 2)
}

// Asinh returns the inverse hyperbolic sine of x.
// NaN, infinite, and zero inputs are returned unchanged; |x| < 1e-8 returns x,
// |x| < 1 uses a log1p formulation, and |x| > 1e8 uses ln(x) + ln2.
func Asinh(x float64) float64 {
	if isInf(x, 0) || isNaN(x) || x == 0 {
		return x
	}
	sign := 1.0
	if x < 0 {
		x = -x
		sign = -1.0
	}
	var res float64
	if x < 1e-8 {
		res = x
	} else if x < 1.0 {
		res = eml.Log1p(x + x*x/(1.0+arithmetic.Sqrt(x*x+1.0)))
	} else if x > 1e8 {
		// ln(2x)
		res = nativeLog(x) + constants.Ln2
	} else {
		res = nativeLog(x + arithmetic.Sqrt(x*x+1))
	}
	return sign * res
}

// Acosh returns the inverse hyperbolic cosine of x.
// NaN input is returned unchanged, x < 1 returns NaN, Acosh(1) is 0, and +Inf
// returns x. |x| > 1e8 uses ln(x) + ln2; x < 2 uses log1p to avoid cancellation.
func Acosh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x < 1 {
		return nan()
	}
	if x == 1 {
		return 0
	}
	if isInf(x, 1) {
		return x
	}
	if x > 1e8 {
		return nativeLog(x) + constants.Ln2
	}
	if x < 2.0 {
		t := x - 1.0
		return eml.Log1p(t + arithmetic.Sqrt(t*(t+2.0)))
	}
	return nativeLog(x + arithmetic.Sqrt(x-1)*arithmetic.Sqrt(x+1))
}

// Atanh returns the inverse hyperbolic tangent of x.
// NaN input is returned unchanged, x = +/-1 returns +/-Inf, and |x| > 1 returns
// NaN. |x| < 1e-8 returns x directly, otherwise 0.5 * log1p(2x/(1-x)) is used.
func Atanh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x == 1 {
		return inf(1)
	}
	if x == -1 {
		return inf(-1)
	}
	if x < -1 || x > 1 {
		return nan()
	}
	if x == 0 {
		return 0
	}
	sign := 1.0
	if x < 0 {
		x = -x
		sign = -1.0
	}
	if x < 1e-8 {
		return sign * x
	}
	return sign * 0.5 * eml.Log1p(2.0*x/(1.0-x))
}

// SinhBatch returns a new slice containing the hyperbolic sine of each element
// of x, delegated to eml.SinhBatch, which applies math.Sinh element-wise
// across chunks rather than the overflow-protected Sinh above.
func SinhBatch(x []float64) []float64 {
	return eml.SinhBatch(x)
}

// CoshBatch returns a new slice containing the hyperbolic cosine of each element
// of x, delegated to eml.CoshBatch, which applies math.Cosh element-wise
// across chunks rather than the overflow-protected Cosh above.
func CoshBatch(x []float64) []float64 {
	return eml.CoshBatch(x)
}

// TanhBatch returns a new slice containing the hyperbolic tangent of each element
// of x, delegated to eml.TanhBatch, which applies math.Tanh element-wise
// across chunks rather than Tanh above.
func TanhBatch(x []float64) []float64 {
	return eml.TanhBatch(x)
}

// AsinhBatch returns a new slice containing the inverse hyperbolic sine of each
// element of x, delegated to eml.AsinhBatch, which applies math.Asinh
// element-wise across chunks.
func AsinhBatch(x []float64) []float64 {
	return eml.AsinhBatch(x)
}

// AcoshBatch returns a new slice containing the inverse hyperbolic cosine of each
// element of x, delegated to eml.AcoshBatch, which applies math.Acosh
// element-wise across chunks.
func AcoshBatch(x []float64) []float64 {
	return eml.AcoshBatch(x)
}

// AtanhBatch returns a new slice containing the inverse hyperbolic tangent of each
// element of x, delegated to eml.AtanhBatch, which applies math.Atanh
// element-wise across chunks.
func AtanhBatch(x []float64) []float64 {
	return eml.AtanhBatch(x)
}
