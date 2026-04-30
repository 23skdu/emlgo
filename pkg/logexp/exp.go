package logexp

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
)

//go:inline
func Exp(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 1) {
		return math.Inf(1)
	}
	if math.IsInf(x, -1) {
		return 0
	}
	if x > 709.782712893384 {
		return math.Inf(1)
	}
	if x < -708.4779660139996 {
		return 0
	}
	return math.Exp(x)
}

//go:inline
func Log(x float64) float64 {
	if x <= 0 {
		return math.NaN()
	}
	return math.Log(x)
}

func ExpBatch(x []float64) []float64 {
	return eml.ExpSIMD(x)
}

func LogBatch(x []float64) []float64 {
	return eml.LogSIMD(x)
}