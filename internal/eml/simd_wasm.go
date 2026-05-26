//go:build wasm

package eml

// addWasmSIMD implements batch addition with 8-wide block unrolling
// to enable WASM JIT auto-vectorization to wasm_simd128.
func addWasmSIMD(a, b, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = a[i] + b[i]
		result[i+1] = a[i+1] + b[i+1]
		result[i+2] = a[i+2] + b[i+2]
		result[i+3] = a[i+3] + b[i+3]
		result[i+4] = a[i+4] + b[i+4]
		result[i+5] = a[i+5] + b[i+5]
		result[i+6] = a[i+6] + b[i+6]
		result[i+7] = a[i+7] + b[i+7]
		i += 8
	}
	for ; i < n; i++ {
		result[i] = a[i] + b[i]
	}
}

func subWasmSIMD(a, b, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = a[i] - b[i]
		result[i+1] = a[i+1] - b[i+1]
		result[i+2] = a[i+2] - b[i+2]
		result[i+3] = a[i+3] - b[i+3]
		result[i+4] = a[i+4] - b[i+4]
		result[i+5] = a[i+5] - b[i+5]
		result[i+6] = a[i+6] - b[i+6]
		result[i+7] = a[i+7] - b[i+7]
		i += 8
	}
	for ; i < n; i++ {
		result[i] = a[i] - b[i]
	}
}

func mulWasmSIMD(a, b, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = a[i] * b[i]
		result[i+1] = a[i+1] * b[i+1]
		result[i+2] = a[i+2] * b[i+2]
		result[i+3] = a[i+3] * b[i+3]
		result[i+4] = a[i+4] * b[i+4]
		result[i+5] = a[i+5] * b[i+5]
		result[i+6] = a[i+6] * b[i+6]
		result[i+7] = a[i+7] * b[i+7]
		i += 8
	}
	for ; i < n; i++ {
		result[i] = a[i] * b[i]
	}
}

func divWasmSIMD(a, b, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = a[i] / b[i]
		result[i+1] = a[i+1] / b[i+1]
		result[i+2] = a[i+2] / b[i+2]
		result[i+3] = a[i+3] / b[i+3]
		result[i+4] = a[i+4] / b[i+4]
		result[i+5] = a[i+5] / b[i+5]
		result[i+6] = a[i+6] / b[i+6]
		result[i+7] = a[i+7] / b[i+7]
		i += 8
	}
	for ; i < n; i++ {
		result[i] = a[i] / b[i]
	}
}

func sqrtWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeSqrt(a[i])
		result[i+1] = nativeSqrt(a[i+1])
		result[i+2] = nativeSqrt(a[i+2])
		result[i+3] = nativeSqrt(a[i+3])
		result[i+4] = nativeSqrt(a[i+4])
		result[i+5] = nativeSqrt(a[i+5])
		result[i+6] = nativeSqrt(a[i+6])
		result[i+7] = nativeSqrt(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeSqrt(a[i])
	}
}

func absWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeAbs(a[i])
		result[i+1] = nativeAbs(a[i+1])
		result[i+2] = nativeAbs(a[i+2])
		result[i+3] = nativeAbs(a[i+3])
		result[i+4] = nativeAbs(a[i+4])
		result[i+5] = nativeAbs(a[i+5])
		result[i+6] = nativeAbs(a[i+6])
		result[i+7] = nativeAbs(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeAbs(a[i])
	}
}

func negWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeNeg(a[i])
		result[i+1] = nativeNeg(a[i+1])
		result[i+2] = nativeNeg(a[i+2])
		result[i+3] = nativeNeg(a[i+3])
		result[i+4] = nativeNeg(a[i+4])
		result[i+5] = nativeNeg(a[i+5])
		result[i+6] = nativeNeg(a[i+6])
		result[i+7] = nativeNeg(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeNeg(a[i])
	}
}

func invWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeInv(a[i])
		result[i+1] = nativeInv(a[i+1])
		result[i+2] = nativeInv(a[i+2])
		result[i+3] = nativeInv(a[i+3])
		result[i+4] = nativeInv(a[i+4])
		result[i+5] = nativeInv(a[i+5])
		result[i+6] = nativeInv(a[i+6])
		result[i+7] = nativeInv(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeInv(a[i])
	}
}

func addScalarWasmSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = a[i] + b
		result[i+1] = a[i+1] + b
		result[i+2] = a[i+2] + b
		result[i+3] = a[i+3] + b
		result[i+4] = a[i+4] + b
		result[i+5] = a[i+5] + b
		result[i+6] = a[i+6] + b
		result[i+7] = a[i+7] + b
		i += 8
	}
	for ; i < n; i++ {
		result[i] = a[i] + b
	}
}

func mulScalarWasmSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = a[i] * b
		result[i+1] = a[i+1] * b
		result[i+2] = a[i+2] * b
		result[i+3] = a[i+3] * b
		result[i+4] = a[i+4] * b
		result[i+5] = a[i+5] * b
		result[i+6] = a[i+6] * b
		result[i+7] = a[i+7] * b
		i += 8
	}
	for ; i < n; i++ {
		result[i] = a[i] * b
	}
}

func fmaWasmSIMD(a, b, c, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = a[i]*b[i] + c[i]
		result[i+1] = a[i+1]*b[i+1] + c[i+1]
		result[i+2] = a[i+2]*b[i+2] + c[i+2]
		result[i+3] = a[i+3]*b[i+3] + c[i+3]
		result[i+4] = a[i+4]*b[i+4] + c[i+4]
		result[i+5] = a[i+5]*b[i+5] + c[i+5]
		result[i+6] = a[i+6]*b[i+6] + c[i+6]
		result[i+7] = a[i+7]*b[i+7] + c[i+7]
		i += 8
	}
	for ; i < n; i++ {
		result[i] = a[i]*b[i] + c[i]
	}
}

func detectWasmSIMD() {
	hasWasmSIMD = true
	hasSSE4 = false
	hasAVX2 = false
	hasAVX512 = false
	hasNeon = false
	hasNeonDot = false
	hasSVE = false
}
