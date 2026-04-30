//go:build !arm64
// +build !arm64

package eml

func addNEON(a, b, result []float64)           {}
func subNEON(a, b, result []float64)           {}
func mulNEON(a, b, result []float64)           {}
func divNEON(a, b, result []float64)           {}
func addScalarNEON(a []float64, b float64, result []float64) {}
func mulScalarNEON(a []float64, b float64, result []float64) {}
func sqrtNEON(a, result []float64)                          {}

func detectARM64SIMD() {}
