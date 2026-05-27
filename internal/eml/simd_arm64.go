//go:build arm64
// +build arm64

package eml

// Go implementations for scalar ops that the compiler inlines.
// On arm64, hasFMA is false so FmaScalar already uses a*b+c; and
// abs/neg/sqrt are single-instruction intrinsics from the runtime.
func absScalar(x float64) float64  { return nativeAbs(x) }
func negScalar(x float64) float64  { return -x }
func sqrtScalar(x float64) float64 { return nativeSqrt(x) }
func fmaScalar(a, b, c float64) float64 { return a*b + c }

var sveVectorLength int

func sveVL() int {
	vl := sveVectorLength / 8
	if vl == 0 {
		return 16
	}
	return vl
}

func addSVE(a, b, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j] + b[i+j]
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = a[i] + b[i]
	}
}

func subSVE(a, b, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j] - b[i+j]
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = a[i] - b[i]
	}
}

func mulSVE(a, b, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j] * b[i+j]
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = a[i] * b[i]
	}
}

func divSVE(a, b, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j] / b[i+j]
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = a[i] / b[i]
	}
}

func addScalarSVE(a []float64, b float64, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j] + b
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = a[i] + b
	}
}

func mulScalarSVE(a []float64, b float64, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j] * b
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = a[i] * b
	}
}

func sqrtSVE(a, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = nativeSqrt(a[i+j])
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = nativeSqrt(a[i])
	}
}

func absSVE(a, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = nativeAbs(a[i+j])
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = nativeAbs(a[i])
	}
}

func negSVE(a, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = nativeNeg(a[i+j])
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = nativeNeg(a[i])
	}
}

func invSVE(a, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = nativeInv(a[i+j])
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = nativeInv(a[i])
	}
}

func fmaSVE(a, b, c, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j]*b[i+j] + c[i+j]
		}
	}
	for i := (n / vl) * vl; i < n; i++ {
		result[i] = a[i]*b[i] + c[i]
	}
}

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
		result[i] = nativeSqrt(a[i])
	}
}

func detectARM64SIMD() {
	hasSSE4 = false
	hasAVX2 = false
	hasAVX512 = false
	hasNeon = true
	hasNeonDot = true
	detectSVE()
}


