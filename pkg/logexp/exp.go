package logexp

import (
	"math"

	"github.com/23skdu/emlgo/internal/constants"
	"github.com/23skdu/emlgo/internal/eml"
)

var (
	nativeExp = eml.Exp
	nativeLog = eml.Log
)

// Overflow/underflow thresholds for Exp:
// expOverflow  ≈ ln(MaxFloat64)           — exp(x) → +Inf beyond this
// expUnderflow ≈ ln(SmallestNonzeroFloat64) — exp(x) → 0 below this
const expOverflow = constants.ExpOverflow
const expUnderflow = -745.133224101734

// Exp returns e**x via the internal eml.Exp implementation, which follows the
// standard math behaviour: NaN for NaN input and +Inf beyond ~709.78.
func Exp(x float64) float64 {
	return nativeExp(x)
}

// Log returns the natural logarithm of x.
// It returns NaN for x < 0 and -Inf for x == 0 before delegating to eml.Log.
func Log(x float64) float64 {
	if x < 0 {
		return eml.NaN()
	}
	if x == 0 {
		return eml.Inf(-1)
	}
	return nativeLog(x)
}

// ExpBatch returns a new slice containing e**x[i] for each element of x,
// computed through SIMD dispatch (eml.ExpSIMD).
func ExpBatch(x []float64) []float64 {
	return eml.ExpSIMD(x)
}

// LogBatch returns a new slice containing ln(x[i]) for each element of x,
// computed through SIMD dispatch (eml.LogSIMD).
func LogBatch(x []float64) []float64 {
	return eml.LogSIMD(x)
}

// ExpFast returns e**x after explicit range checks: +Inf above expOverflow,
// 0 below expUnderflow, and math.Exp(x) otherwise.
func ExpFast(x float64) float64 {
	if x > expOverflow {
		return math.Inf(1)
	}
	if x < expUnderflow {
		return 0
	}
	return math.Exp(x)
}

// LogFast returns the natural logarithm of x.
// It returns NaN for x < 0 and -Inf for x == 0, otherwise eml.Log(x).
func LogFast(x float64) float64 {
	if x < 0 {
		return math.NaN()
	}
	if x == 0 {
		return math.Inf(-1)
	}
	return nativeLog(x)
}
