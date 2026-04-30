package logexp

import (
	"github.com/emlgo/eml/internal/eml"
)

var (
	nativeExp = eml.Exp
	nativeLog = eml.Log
)

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