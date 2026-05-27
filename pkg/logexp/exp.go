package logexp

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
)

var (
	nativeExp = eml.Exp
	nativeLog = eml.Log
)

// Overflow/underflow thresholds for Exp:
// expOverflow  ≈ ln(MaxFloat64)           — exp(x) → +Inf beyond this
// expUnderflow ≈ ln(SmallestNonzeroFloat64) — exp(x) → 0 below this
const expOverflow = 709.782712893384
const expUnderflow = -745.133224101734

func Exp(x float64) float64 {
	return nativeExp(x)
}

func Log(x float64) float64 {
	if x <= 0 {
		return eml.NaN()
	}
	return nativeLog(x)
}

func ExpBatch(x []float64) []float64 {
	return eml.ExpSIMD(x)
}

func LogBatch(x []float64) []float64 {
	return eml.LogSIMD(x)
}

func ExpFast(x float64) float64 {
	if x > expOverflow {
		return math.Inf(1)
	}
	if x < expUnderflow {
		return 0
	}
	return math.Exp(x)
}

func LogFast(x float64) float64 {
	if x <= 0 {
		return math.NaN()
	}
	return math.Log(x)
}