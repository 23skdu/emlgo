//go:build arm64 && !purego
// +build arm64,!purego

package eml

// Go implementations for scalar ops that the compiler inlines.
func absScalar(x float64) float64  { return nativeAbs(x) }
func negScalar(x float64) float64  { return -x }
func sqrtScalar(x float64) float64 { return nativeSqrt(x) }
func fmaScalar(a, b, c float64) float64

// Real NEON vector elementwise kernels implemented in simd_arm64.s
//
//go:noescape
func addNEON(a, b, result []float64)

//go:noescape
func subNEON(a, b, result []float64)

//go:noescape
func mulNEON(a, b, result []float64)

//go:noescape
func divNEON(a, b, result []float64)

//go:noescape
func addScalarNEON(a []float64, b float64, result []float64)

//go:noescape
func mulScalarNEON(a []float64, b float64, result []float64)

//go:noescape
func sqrtNEON(a, result []float64)

//go:noescape
func absNEON(a, result []float64)

//go:noescape
func negNEON(a, result []float64)

//go:noescape
func invNEON(a, result []float64)

//go:noescape
func fmaNEON(a, b, c, result []float64)

// Float32 NEON vector kernels
//
//go:noescape
func addNEONF32(a, b, result []float32)

//go:noescape
func subNEONF32(a, b, result []float32)

//go:noescape
func mulNEONF32(a, b, result []float32)

//go:noescape
func divNEONF32(a, b, result []float32)

//go:noescape
func sqrtNEONF32(a, result []float32)

//go:noescape
func addScalarNEONF32(a []float32, b float32, result []float32)

//go:noescape
func mulScalarNEONF32(a []float32, b float32, result []float32)

//go:noescape
func absNEONF32(a, result []float32)

//go:noescape
func negNEONF32(a, result []float32)

//go:noescape
func invNEONF32(a, result []float32)

//go:noescape
func fmaNEONF32(a, b, c, result []float32)

// Signed 8-bit saturating vector kernels
//
//go:noescape
func addSatInt8NEON(a, b, result []int8)

//go:noescape
func subSatInt8NEON(a, b, result []int8)

// addSVE is retained as a stub for backwards-compatibility with test suites
func addSVE(a, b, result []float64) {
	_, _, _ = a, b, result
}

func detectARM64SIMD() {
	hasSSE4 = false
	hasAVX2 = false
	hasAVX512 = false
	hasNeon = true
	hasFMA = true
	hasSVE = false
	hasNeonDot = false

	detectSVE()
}
