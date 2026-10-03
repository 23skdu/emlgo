package hyper

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/arithmetic"
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
// ln(MaxFloat64) ≈ 709.78 — exp(x) overflows to +Inf beyond this.
const expOverflow = 709.78

func Sinh(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}
	sign := 1.0
	if x < 0 {
		x = -x
		sign = -1.0
	}
	if x > expOverflow {
		return sign * inf(1)
	}
	var res float64
	if x < 1e-8 {
		res = x
	} else if x < 0.5 {
		// Avoid catastrophic cancellation: (expm1(x) - expm1(-x)) / 2
		res = (math.Expm1(x) - math.Expm1(-x)) / 2
	} else if x > 22.0 {
		res = nativeExp(x) / 2
	} else {
		res = (nativeExp(x) - nativeExp(-x)) / 2
	}
	return sign * res
}

func Cosh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if isInf(x, 0) {
		return inf(1)
	}
	absX := nativeAbs(x)
	if absX > expOverflow {
		return inf(1)
	}
	if absX > 22.0 {
		return nativeExp(absX) / 2
	}
	ex := nativeExp(absX)
	return 0.5 * (ex + 1.0/ex)
}

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
	if nativeAbs(x) < 1e-8 {
		return x
	}
	// For small to medium x, use expm1(2x)/(expm1(2x)+2) to preserve accuracy
	if nativeAbs(x) < 0.5 {
		e2x := math.Expm1(2 * x)
		return e2x / (e2x + 2)
	}
	ex := nativeExp(x)
	emx := nativeExp(-x)
	return (ex - emx) / (ex + emx)
}

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
	} else if x > 1e8 {
		// ln(2x)
		res = nativeLog(x) + 0.693147180559945309417232121458
	} else {
		res = nativeLog(x + arithmetic.Sqrt(x*x+1))
	}
	return sign * res
}

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
		return nativeLog(x) + 0.693147180559945309417232121458
	}
	return nativeLog(x + arithmetic.Sqrt(x-1)*arithmetic.Sqrt(x+1))
}

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
	return sign * nativeLog((1+x)/(1-x)) / 2
}

func SinhBatch(x []float64) []float64 {
	return eml.SinhBatch(x)
}

func CoshBatch(x []float64) []float64 {
	return eml.CoshBatch(x)
}

func TanhBatch(x []float64) []float64 {
	return eml.TanhBatch(x)
}

func AsinhBatch(x []float64) []float64 {
	return eml.AsinhBatch(x)
}

func AcoshBatch(x []float64) []float64 {
	return eml.AcoshBatch(x)
}

func AtanhBatch(x []float64) []float64 {
	return eml.AtanhBatch(x)
}
