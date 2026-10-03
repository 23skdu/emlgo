//go:build amd64 && !purego && emlasm

package eml

// Experimental: dispatch the transcendental batch kernels to the hand-written
// AVX2 assembly instead of the generic parallel path.
//
// Enabled with `-tags emlasm`. This path is currently both slower and less
// accurate than the default; see simd_trans_amd64.go for the measurements and
// docs/nextsteps.md Step 1 for the root cause and the rework plan.
//
// The accuracy gates in simd_accuracy_test.go are expected to FAIL under this
// tag. That is intentional: they are what a reworked kernel must pass before
// this path can be re-enabled by default.

func dispatchExpSIMDTo(x, result []float64) {
	n := len(x)
	if hasAVX2 {
		simdLen := (n / 4) * 4
		expAVX2(x[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeExp(x[i])
		}
		return
	}
	parallelizeGeneric(x, result, nativeExp)
}

func dispatchLogSIMDTo(x, result []float64) {
	n := len(x)
	if hasAVX2 {
		simdLen := (n / 4) * 4
		logAVX2(x[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeLog(x[i])
		}
		return
	}
	parallelizeGeneric(x, result, logScalar)
}

func dispatchSinSIMDTo(x, result []float64) {
	n := len(x)
	if hasAVX2 {
		simdLen := (n / 4) * 4
		sinAVX2(x[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeSin(x[i])
		}
		return
	}
	parallelizeGeneric(x, result, nativeSin)
}

func dispatchCosSIMDTo(x, result []float64) {
	n := len(x)
	if hasAVX2 {
		simdLen := (n / 4) * 4
		cosAVX2(x[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeCos(x[i])
		}
		return
	}
	parallelizeGeneric(x, result, nativeCos)
}

func dispatchTanSIMDTo(x, result []float64) {
	n := len(x)
	if hasAVX2 {
		simdLen := (n / 4) * 4
		tanAVX2(x[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeTan(x[i])
		}
		return
	}
	parallelizeGeneric(x, result, nativeTan)
}
