//go:build (!amd64 && !arm64 && !wasm) || purego
// +build !amd64,!arm64,!wasm purego

package eml

func stubAddSIMD(a, b, result []float64) {
	for i := range result {
		result[i] = a[i] + b[i]
	}
}

func stubSubSIMD(a, b, result []float64) {
	for i := range result {
		result[i] = a[i] - b[i]
	}
}

func stubMulSIMD(a, b, result []float64) {
	for i := range result {
		result[i] = a[i] * b[i]
	}
}

func stubDivSIMD(a, b, result []float64) {
	for i := range result {
		result[i] = a[i] / b[i]
	}
}

func stubAbsSIMD(a, result []float64) {
	for i := range a {
		result[i] = nativeAbs(a[i])
	}
}

func stubNegSIMD(a, result []float64) {
	for i := range a {
		result[i] = nativeNeg(a[i])
	}
}

func stubInvSIMD(a, result []float64) {
	for i := range a {
		result[i] = nativeInv(a[i])
	}
}

func fmaSIMD(a, b, c, result []float64) {
	for i := range a {
		result[i] = a[i]*b[i] + c[i]
	}
}

func emlSIMD(x, y, result []float64) {
	scalarEml(x, y, result)
}

func detectPlatformSIMD() {
}

func scalarEml(x, y, result []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		result[i] = nativeExp(x[i]) - nativeLog(y[i])
	}
}

func dispatchExpSIMDTo(x, result []float64)      { parallelizeGeneric(x, result, nativeExp) }
func dispatchLogSIMDTo(x, result []float64)      { parallelizeGeneric(x, result, logScalar) }
func dispatchSinSIMDTo(x, result []float64)      { parallelizeGeneric(x, result, nativeSin) }
func dispatchCosSIMDTo(x, result []float64)      { parallelizeGeneric(x, result, nativeCos) }
func dispatchTanSIMDTo(x, result []float64)      { parallelizeGeneric(x, result, nativeTan) }
func dispatchSinCosSIMDTo(x, sin, cos []float64) { parallelizeSinCos(x, sin, cos) }
func dispatchSqrtSIMDTo(x, result []float64)     { parallelizeGeneric(x, result, nativeSqrt) }
func dispatchAddSIMD(a, b, result []float64)     { stubAddSIMD(a, b, result) }
func dispatchSubSIMD(a, b, result []float64)     { stubSubSIMD(a, b, result) }
func dispatchMulSIMD(a, b, result []float64)     { stubMulSIMD(a, b, result) }
func dispatchDivSIMD(a, b, result []float64)     { stubDivSIMD(a, b, result) }
func dispatchAbsSIMD(x, result []float64)        { stubAbsSIMD(x, result) }
func dispatchNegSIMD(x, result []float64)        { stubNegSIMD(x, result) }
func dispatchInvSIMD(x, result []float64)        { stubInvSIMD(x, result) }

func dispatchAddScalarSIMD(a []float64, b float64, result []float64) {
	for i := range a {
		result[i] = a[i] + b
	}
}

func dispatchMulScalarSIMD(a []float64, b float64, result []float64) {
	for i := range a {
		result[i] = a[i] * b
	}
}

func dispatchAddSatInt8SIMD(a, b, result []int8) {
	for i := range a {
		v := int32(a[i]) + int32(b[i])
		if v > 127 {
			v = 127
		} else if v < -128 {
			v = -128
		}
		result[i] = int8(v)
	}
}

func dispatchSubSatInt8SIMD(a, b, result []int8) {
	for i := range a {
		v := int32(a[i]) - int32(b[i])
		if v > 127 {
			v = 127
		} else if v < -128 {
			v = -128
		}
		result[i] = int8(v)
	}
}
