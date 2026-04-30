package hyper

import (
	"math"

	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/logexp"
)

func Sinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex - emx) / 2
}

func Cosh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return math.Abs(x)
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex + emx) / 2
}

func Tanh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 1) {
		return 1
	}
	if math.IsInf(x, -1) {
		return -1
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex - emx) / (ex + emx)
}

func Asinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	if x == 0 {
		return 0
	}
	return logexp.Log(x + arithmetic.Sqrt(x*x+1))
}

func Acosh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x < 1 {
		return math.NaN()
	}
	if x == 1 {
		return 0
	}
	if math.IsInf(x, 0) {
		return x
	}
	return logexp.Log(x + arithmetic.Sqrt(x-1)*arithmetic.Sqrt(x+1))
}

func Atanh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x <= -1 || x >= 1 {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	return logexp.Log((1+x)/(1-x)) / 2
}