package eml

import "math"

const One = 1.0

func Eml(x, y float64) float64 {
	return math.Exp(x) - math.Log(y)
}

func EmlOne(x float64) float64 {
	return Eml(x, One)
}

func OneEml(x float64) float64 {
	return Eml(One, x)
}