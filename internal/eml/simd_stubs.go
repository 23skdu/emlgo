//go:build !amd64 && !arm64
// +build !amd64,!arm64

package eml

import "math"

func absScalar(x float64) float64 {
	return math.Abs(x)
}

func negScalar(x float64) float64 {
	return -x
}

func sqrtScalar(x float64) float64 {
	return math.Sqrt(x)
}

func fmaScalar(a, b, c float64) float64 {
	return a*b + c
}
