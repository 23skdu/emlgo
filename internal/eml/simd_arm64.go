//go:build arm64
// +build arm64

package eml

func addNEON(a, b, result []float64) {
	for i := range a {
		result[i] = a[i] + b[i]
	}
}

func subNEON(a, b, result []float64) {
	for i := range a {
		result[i] = a[i] - b[i]
	}
}

func mulNEON(a, b, result []float64) {
	for i := range a {
		result[i] = a[i] * b[i]
	}
}

func divNEON(a, b, result []float64) {
	for i := range a {
		result[i] = a[i] / b[i]
	}
}

func addScalarNEON(a []float64, b float64, result []float64) {
	for i := range a {
		result[i] = a[i] + b
	}
}

func mulScalarNEON(a []float64, b float64, result []float64) {
	for i := range a {
		result[i] = a[i] * b
	}
}

func sqrtNEON(a, result []float64) {
	for i := range a {
		result[i] = sqrtScalar(a[i])
	}
}

func sqrtScalar(x float64) float64
func fmaScalar(a, b, c float64) float64

func detectARM64SIMD() {
	// NEON is always available on our target ARM64 platforms, 
	// but we use scalar fallback for now due to assembler limitations.
	hasSSE4 = false
	hasAVX2 = false
	hasAVX512 = false
	hasNeon = true
	hasNeonDot = true
}
