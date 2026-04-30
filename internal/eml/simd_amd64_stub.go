//go:build !amd64
// +build !amd64

package eml

func cpuid(_, _ uint32) (eax, ebx, ecx, edx uint32) { return 0, 0, 0, 0 }

func addAVX2(a, b, result []float64)           {}
func subAVX2(a, b, result []float64)           {}
func mulAVX2(a, b, result []float64)           {}
func divAVX2(a, b, result []float64)           {}
func addScalarAVX2(a []float64, b float64, result []float64) {}
func mulScalarAVX2(a []float64, b float64, result []float64) {}

func addAVX512(a, b, result []float64)           {}
func subAVX512(a, b, result []float64)           {}
func mulAVX512(a, b, result []float64)           {}
func divAVX512(a, b, result []float64)           {}
func addScalarAVX512(a []float64, b float64, result []float64) {}
func mulScalarAVX512(a []float64, b float64, result []float64) {}

func detectAMD64SIMD() {}
