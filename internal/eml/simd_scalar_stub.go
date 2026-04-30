//go:build !amd64 && !arm64
// +build !amd64,!arm64

package eml

func sqrtScalar(x float64) float64     { return nativeSqrt(x) }
func fmaScalar(a, b, c float64) float64 { return a*b + c }
