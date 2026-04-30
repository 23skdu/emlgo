package hyper

import (
	"math"

	"github.com/emlgo/eml/pkg/logexp"
)

func Sinh(x float64) float64 {
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex - emx) / 2
}

func Cosh(x float64) float64 {
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex + emx) / 2
}

func Tanh(x float64) float64 {
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex - emx) / (ex + emx)
}

func Asinh(x float64) float64 {
	return math.Asinh(x)
}

func Acosh(x float64) float64 {
	if x < 1 {
		return math.NaN()
	}
	return math.Acosh(x)
}

func Atanh(x float64) float64 {
	if x <= -1 || x >= 1 {
		return math.NaN()
	}
	return math.Atanh(x)
}