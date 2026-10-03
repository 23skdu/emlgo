//go:build arm64 && !purego
// +build arm64,!purego

package eml

func arm64AddSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		addNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = a[i] + b[i]
	}
}

func arm64SubSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		subNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = a[i] - b[i]
	}
}

func arm64MulSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		mulNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = a[i] * b[i]
	}
}

func arm64DivSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		divNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = a[i] / b[i]
	}
}

func arm64AddScalarSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		addScalarNEON(a[:simdLen], b, result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = a[i] + b
	}
}

func arm64MulScalarSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		mulScalarNEON(a[:simdLen], b, result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = a[i] * b
	}
}

func arm64SqrtSIMD(a, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		sqrtNEON(a[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeSqrt(a[i])
	}
}

func arm64AbsSIMD(a, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		absNEON(a[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeAbs(a[i])
	}
}

func arm64NegSIMD(a, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		negNEON(a[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeNeg(a[i])
	}
}

func arm64InvSIMD(a, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		invNEON(a[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = nativeInv(a[i])
	}
}

func emlSIMD(x, y, result []float64) {
	n := len(x)
	if n >= 8 {
		neonEml(x, y, result)
	} else {
		scalarEml(x, y, result)
	}
}

func neonEml(x, y, result []float64) {
	n := len(x)
	chunk := 4
	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = nativeExp(x[j]) - nativeLog(y[j])
		}
	}
}

func scalarEml(x, y, result []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		result[i] = nativeExp(x[i]) - nativeLog(y[i])
	}
}

func detectPlatformSIMD() {
	detectARM64SIMD()
}

func fmaSIMD(a, b, c, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	if simdLen > 0 {
		fmaNEON(a[:simdLen], b[:simdLen], c[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		result[i] = a[i]*b[i] + c[i]
	}
}

func dispatchExpSIMDTo(x, result []float64) {
	parallelizeGeneric(x, result, nativeExp)
}

func dispatchLogSIMDTo(x, result []float64) {
	parallelizeGeneric(x, result, logScalar)
}

func dispatchSinSIMDTo(x, result []float64) {
	parallelizeGeneric(x, result, nativeSin)
}

func dispatchCosSIMDTo(x, result []float64) {
	parallelizeGeneric(x, result, nativeCos)
}

func dispatchTanSIMDTo(x, result []float64) {
	parallelizeGeneric(x, result, nativeTan)
}

func dispatchSinCosSIMDTo(x, sin, cos []float64) {
	parallelizeSinCos(x, sin, cos)
}

func dispatchSqrtSIMDTo(x, result []float64) {
	arm64SqrtSIMD(x, result)
}

func dispatchAddSIMD(a, b, result []float64) {
	arm64AddSIMD(a, b, result)
}

func dispatchSubSIMD(a, b, result []float64) {
	arm64SubSIMD(a, b, result)
}

func dispatchMulSIMD(a, b, result []float64) {
	arm64MulSIMD(a, b, result)
}

func dispatchDivSIMD(a, b, result []float64) {
	arm64DivSIMD(a, b, result)
}

func dispatchAbsSIMD(x, result []float64) {
	arm64AbsSIMD(x, result)
}

func dispatchNegSIMD(x, result []float64) {
	arm64NegSIMD(x, result)
}

func dispatchInvSIMD(x, result []float64) {
	arm64InvSIMD(x, result)
}

func dispatchAddScalarSIMD(a []float64, b float64, result []float64) {
	arm64AddScalarSIMD(a, b, result)
}

func dispatchMulScalarSIMD(a []float64, b float64, result []float64) {
	arm64MulScalarSIMD(a, b, result)
}

func dispatchAddSatInt8SIMD(a, b, result []int8) {
	n := len(a)
	simdLen := (n / 16) * 16
	if simdLen > 0 {
		addSatInt8NEON(a[:simdLen], b[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
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
	n := len(a)
	simdLen := (n / 16) * 16
	if simdLen > 0 {
		subSatInt8NEON(a[:simdLen], b[:simdLen], result[:simdLen])
	}
	for i := simdLen; i < n; i++ {
		v := int32(a[i]) - int32(b[i])
		if v > 127 {
			v = 127
		} else if v < -128 {
			v = -128
		}
		result[i] = int8(v)
	}
}
