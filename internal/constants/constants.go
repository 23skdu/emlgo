package constants

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
)

const (
	One     = 1.0
	E       = math.E
	Pi      = math.Pi
	I       = complex(0, 1)
	NegOne  = -1.0
	Two     = 2.0
	Half    = 0.5
	Sqrt2   = 1.4142135623730951
)

func GenerateE() float64 {
	return eml.Eml(eml.Eml(1.0, 1.0), 1.0)
}

func GeneratePi() float64 {
	return math.Pi
}

func GenerateI() complex128 {
	return I
}