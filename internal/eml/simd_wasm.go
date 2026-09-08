//go:build wasm

package eml

// WASM SIMD128 optimization strategy:
//
// WebAssembly SIMD128 provides 128-bit vector lanes that operate on
// 2×float64 or 4×float32 simultaneously. The Go WASM JIT compiler
// (wazero / the Go compiler's WASM backend) can auto-vectorize loops
// when it detects a predictable, straight-line access pattern over
// contiguous memory.
//
// The 8-wide unrolling in each function serves two purposes:
//  1. It exposes enough independent operations for the JIT to fill
//     two 128-bit float64 lanes (2 elements per SIMD register × 4
//     registers = 8 elements) in a single pipeline depth.
//  2. It amortizes loop overhead and branch prediction cost, which is
//     critical on WASM runtimes where indirect-call overhead is high.
//
// For transcendental operations (exp, log, sin, cos, …) the unrolled
// calls to native math wrappers still execute as scalar calls, but
// the surrounding loop structure is preserved so that a future WASM
// SIMD math library (e.g. wasm-SIMD-accelerated libm) can be dropped
// in by replacing the inner call without restructuring the batch loop.

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

func expWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeExp(a[i])
		result[i+1] = nativeExp(a[i+1])
		result[i+2] = nativeExp(a[i+2])
		result[i+3] = nativeExp(a[i+3])
		result[i+4] = nativeExp(a[i+4])
		result[i+5] = nativeExp(a[i+5])
		result[i+6] = nativeExp(a[i+6])
		result[i+7] = nativeExp(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeExp(a[i])
	}
}

func logWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeLog(a[i])
		result[i+1] = nativeLog(a[i+1])
		result[i+2] = nativeLog(a[i+2])
		result[i+3] = nativeLog(a[i+3])
		result[i+4] = nativeLog(a[i+4])
		result[i+5] = nativeLog(a[i+5])
		result[i+6] = nativeLog(a[i+6])
		result[i+7] = nativeLog(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeLog(a[i])
	}
}

func sinWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeSin(a[i])
		result[i+1] = nativeSin(a[i+1])
		result[i+2] = nativeSin(a[i+2])
		result[i+3] = nativeSin(a[i+3])
		result[i+4] = nativeSin(a[i+4])
		result[i+5] = nativeSin(a[i+5])
		result[i+6] = nativeSin(a[i+6])
		result[i+7] = nativeSin(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeSin(a[i])
	}
}

func cosWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeCos(a[i])
		result[i+1] = nativeCos(a[i+1])
		result[i+2] = nativeCos(a[i+2])
		result[i+3] = nativeCos(a[i+3])
		result[i+4] = nativeCos(a[i+4])
		result[i+5] = nativeCos(a[i+5])
		result[i+6] = nativeCos(a[i+6])
		result[i+7] = nativeCos(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeCos(a[i])
	}
}

func tanWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = nativeTan(a[i])
		result[i+1] = nativeTan(a[i+1])
		result[i+2] = nativeTan(a[i+2])
		result[i+3] = nativeTan(a[i+3])
		result[i+4] = nativeTan(a[i+4])
		result[i+5] = nativeTan(a[i+5])
		result[i+6] = nativeTan(a[i+6])
		result[i+7] = nativeTan(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = nativeTan(a[i])
	}
}

func sinhWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = Sinh(a[i])
		result[i+1] = Sinh(a[i+1])
		result[i+2] = Sinh(a[i+2])
		result[i+3] = Sinh(a[i+3])
		result[i+4] = Sinh(a[i+4])
		result[i+5] = Sinh(a[i+5])
		result[i+6] = Sinh(a[i+6])
		result[i+7] = Sinh(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = Sinh(a[i])
	}
}

func coshWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = Cosh(a[i])
		result[i+1] = Cosh(a[i+1])
		result[i+2] = Cosh(a[i+2])
		result[i+3] = Cosh(a[i+3])
		result[i+4] = Cosh(a[i+4])
		result[i+5] = Cosh(a[i+5])
		result[i+6] = Cosh(a[i+6])
		result[i+7] = Cosh(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = Cosh(a[i])
	}
}

func tanhWasmSIMD(a, result []float64) {
	n := len(a)
	i := 0
	for i <= n-8 {
		result[i] = Tanh(a[i])
		result[i+1] = Tanh(a[i+1])
		result[i+2] = Tanh(a[i+2])
		result[i+3] = Tanh(a[i+3])
		result[i+4] = Tanh(a[i+4])
		result[i+5] = Tanh(a[i+5])
		result[i+6] = Tanh(a[i+6])
		result[i+7] = Tanh(a[i+7])
		i += 8
	}
	for ; i < n; i++ {
		result[i] = Tanh(a[i])
	}
}

func sincosWasmSIMD(x, sin, cos []float64) {
	n := len(x)
	i := 0
	for i <= n-8 {
		sin[i], cos[i] = nativeSincos(x[i])
		sin[i+1], cos[i+1] = nativeSincos(x[i+1])
		sin[i+2], cos[i+2] = nativeSincos(x[i+2])
		sin[i+3], cos[i+3] = nativeSincos(x[i+3])
		sin[i+4], cos[i+4] = nativeSincos(x[i+4])
		sin[i+5], cos[i+5] = nativeSincos(x[i+5])
		sin[i+6], cos[i+6] = nativeSincos(x[i+6])
		sin[i+7], cos[i+7] = nativeSincos(x[i+7])
		i += 8
	}
	for ; i < n; i++ {
		sin[i], cos[i] = nativeSincos(x[i])
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
