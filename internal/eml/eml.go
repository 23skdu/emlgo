package eml

import "math"

// one is the default value used for single-argument variants.
const one = 1.0

// Eml returns Exp(x) - Log(y).
func Eml(x, y float64) float64 {
	return math.Exp(x) - math.Log(y)
}

// One returns Eml(x, 1).
func One(x float64) float64 {
	return Eml(x, one)
}

// OneEml returns Eml(1, x).
func OneEml(x float64) float64 {
	return Eml(one, x)
}