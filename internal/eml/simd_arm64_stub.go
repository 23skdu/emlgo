//go:build !arm64
// +build !arm64

package eml

func addNEON(_, _, _ []float64)       {}
func subNEON(_, _, _ []float64)       {}
func mulNEON(_, _, _ []float64)       {}
func divNEON(_, _, _ []float64)       {}
func addScalarNEON(_, _ float64, _ []float64) {}
func mulScalarNEON(_, _ float64, _ []float64) {}
func sqrtNEON(_, _ []float64)         {}
func detectARM64SIMD()                {}
