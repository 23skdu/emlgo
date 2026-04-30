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
	nativeSqrt = eml.Sqrt
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
	if x > 709.78 {
		return 1
	}
	if x < -709.78 {
		return -1
	}
	if isInf(x, 1) {
		return 1
	}
	if isInf(x, -1) {
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
	if isNaN(x) {
		return x
	}
	if isInf(x, 0) {
		return x
	}
	if x == 0 {
		return 0
	}
	absX := nativeAbs(x)
	ln10 := 2.3025850929940456840179914546844
	if absX > 1e150 {
		logx := nativeLog(absX)
		approx := 0.693147180559945309417232121458 + ln10*logx/ln10
		if x > 0 {
			return approx
		}
		return -approx
	}
	if absX > eml.MaxFloat64/2 {
		if x > 0 {
			logx := nativeLog(2 * absX)
			return ln10 * logx / ln10
		}
		logx := nativeLog(2 * absX)
		return -ln10 * logx / ln10
	}
	term := arithmetic.Sqrt(x*x + 1)
	if isInf(term, 1) {
		logx := nativeLog(2 * absX)
		if x > 0 {
			return ln10 * logx / ln10
		}
		return -ln10 * logx / ln10
	}
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
	if nativeAbs(x) > 0.9999999999999999 {
		if x > 0 {
			return inf(1)
		}
		return inf(-1)
	}
	return nativeLog((1+x)/(1-x)) / 2
}