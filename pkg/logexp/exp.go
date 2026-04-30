package logexp

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
)

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
	return eml.EmlOne(x)
}

func Log(x float64) float64 {
	if x <= 0 {
		return math.NaN()
	}
	if x < math.SmallestNonzeroFloat64 {
		// Handle subnormal numbers
		// For very small x, use approximation: log(x) ≈ log(2^(-1074)) + (x - 2^(-1074)) / 2^(-1074)
		// But actually for subnormals, we can use: log(x) = log(2^-149) + log(x * 2^149)
		// Simplest: use math.Log for subnormals since they're rare
		return math.Log(x)
	}
	if x < 1e-100 {
		// Very small but normal - use math.Log
		return math.Log(x)
	}
	return eml.Eml(1, eml.Eml(eml.Eml(1, x), 1))
}