package trig

import (
	"math"
)

func Sin(x float64) float64 {
	return math.Sin(x)
}

func Cos(x float64) float64 {
	return math.Cos(x)
}

func Tan(x float64) float64 {
	return math.Tan(x)
}

func Asin(x float64) float64 {
	if x < -1 || x > 1 {
		return math.NaN()
	}
	return math.Asin(x)
}

func Acos(x float64) float64 {
	if x < -1 || x > 1 {
		return math.NaN()
	}
	return math.Acos(x)
}

func Atan(x float64) float64 {
	return math.Atan(x)
}