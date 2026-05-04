package hyper

import (
	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/arithmetic"
)

var (
	isNaN   = eml.IsNaN
	isInf   = eml.IsInf
	inf     = eml.Inf
	nan     = eml.NaN
	nativeExp = eml.Exp
	nativeLog = eml.Log
	nativeAbs  = eml.Abs
)

func Sinh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if isInf(x, 0) {
		return x
	}
	if x > 709.78 || x < -709.78 {
		if x > 0 {
			return inf(1)
		}
		return inf(-1)
	}
	ex := nativeExp(x)
	emx := nativeExp(-x)
	return (ex - emx) / 2
}

func Cosh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if isInf(x, 0) {
		return nativeAbs(x)
	}
	if x > 709.78 || x < -709.78 {
		return inf(1)
	}
	ex := nativeExp(x)
	emx := nativeExp(-x)
	return (ex + emx) / 2
}

func Tanh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x > 709.78 || isInf(x, 1) {
		return 1
	}
	if x < -709.78 || isInf(x, -1) {
		return -1
	}
	ex := nativeExp(x)
	emx := nativeExp(-x)
	sum := ex + emx
	if isInf(sum, 1) {
		if x > 0 {
			return 1
		}
		return -1
	}
	return (ex - emx) / sum
}

func Asinh(x float64) float64 {
	if isInf(x, 0) || isNaN(x) {
		return x
	}
	if x == 0 {
		return 0
	}
	absX := nativeAbs(x)
	if absX > 1e150 || absX > eml.MaxFloat64/2 {
		logx := nativeLog(absX)
		approx := logx + 0.693147180559945309417232121458
		if x < 0 {
			return -approx
		}
		return approx
	}
	term := arithmetic.Sqrt(x*x + 1)
	return nativeLog(x + term)
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
	if isInf(x, 0) {
		return x
	}
	if x > eml.MaxFloat64/2 {
		return nativeLog(2*x) - 0.693147180559945309417232121458
	}
	return nativeLog(x + arithmetic.Sqrt(x-1)*arithmetic.Sqrt(x+1))
}

func Atanh(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x <= -1 || x >= 1 {
		return nan()
	}
	if x == 0 {
		return 0
	}
	if nativeAbs(x) > 1-1e-16 {
		if x > 0 {
			return inf(1)
		}
		return inf(-1)
	}
	return nativeLog((1+x)/(1-x)) / 2
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
