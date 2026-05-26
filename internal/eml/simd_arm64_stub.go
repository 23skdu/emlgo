//go:build !arm64
// +build !arm64

package eml

func addNEON(_, _, _ []float64)       {}
func subNEON(_, _, _ []float64)       {}
func mulNEON(_, _, _ []float64)       {}
func divNEON(_, _, _ []float64)       {}
func addScalarNEON(_ []float64, _ float64, _ []float64) {}
func mulScalarNEON(_ []float64, _ float64, _ []float64) {}
func sqrtNEON(_, _ []float64)         {}
func detectARM64SIMD()                {}
func addSVE(_, _, _ []float64)        {}
func subSVE(_, _, _ []float64)        {}
func mulSVE(_, _, _ []float64)        {}
func divSVE(_, _, _ []float64)        {}
func addScalarSVE(_ []float64, _ float64, _ []float64) {}
func mulScalarSVE(_ []float64, _ float64, _ []float64) {}
func sqrtSVE(_, _ []float64)          {}
func absSVE(_, _ []float64)           {}
func negSVE(_, _ []float64)           {}
func invSVE(_, _ []float64)           {}
func fmaSVE(_, _, _, _ []float64)     {}
func sveVL() int                      { return 16 }
