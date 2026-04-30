package logexp

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
)

func Exp(x float64) float64 {
	return eml.EmlOne(x)
}

func Log(x float64) float64 {
	if x <= 0 {
		return math.NaN()
	}
	return eml.Eml(1, eml.Eml(eml.Eml(1, x), 1))
}