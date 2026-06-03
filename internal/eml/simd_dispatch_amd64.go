//go:build amd64
// +build amd64

package eml

func amd64AddSIMD(a, b, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		addAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b[i]
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		addAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b[i]
		}
		return
	}
	for i := range result {
		result[i] = a[i] + b[i]
	}
}

func amd64SubSIMD(a, b, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		subAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] - b[i]
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		subAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] - b[i]
		}
		return
	}
	for i := range result {
		result[i] = a[i] - b[i]
	}
}

func amd64MulSIMD(a, b, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		mulAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b[i]
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		mulAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b[i]
		}
		return
	}
	for i := range result {
		result[i] = a[i] * b[i]
	}
}

func amd64DivSIMD(a, b, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		divAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] / b[i]
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		divAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] / b[i]
		}
		return
	}
	for i := range result {
		result[i] = a[i] / b[i]
	}
}

func amd64AddScalarSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		addScalarAVX512(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		addScalarAVX2(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b
		}
		return
	}
	for i := range a {
		result[i] = a[i] + b
	}
}

func amd64MulScalarSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		mulScalarAVX512(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		mulScalarAVX2(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b
		}
		return
	}
	for i := range a {
		result[i] = a[i] * b
	}
}

func amd64SqrtSIMD(a, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		sqrtAVX512(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeSqrt(a[i])
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		sqrtAVX2(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeSqrt(a[i])
		}
		return
	}
	for i := range a {
		result[i] = nativeSqrt(a[i])
	}
}

func amd64AbsSIMD(a, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		absAVX512(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeAbs(a[i])
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		absAVX2(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeAbs(a[i])
		}
		return
	}
	for i := range a {
		result[i] = nativeAbs(a[i])
	}
}

func amd64NegSIMD(a, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		negAVX512(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeNeg(a[i])
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		negAVX2(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeNeg(a[i])
		}
		return
	}
	for i := range a {
		result[i] = nativeNeg(a[i])
	}
}

func amd64InvSIMD(a, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		invAVX512(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeInv(a[i])
		}
		return
	}
	if hasAVX2 {
		simdLen := (n / 4) * 4
		invAVX2(a[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = nativeInv(a[i])
		}
		return
	}
	for i := range a {
		result[i] = nativeInv(a[i])
	}
}

func emlSIMD(x, y, result []float64) {
	scalarEml(x, y, result)
}

func scalarEml(x, y, result []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		result[i] = nativeExp(x[i]) - nativeLog(y[i])
	}
}

func detectPlatformSIMD() {
	detectAMD64SIMD()
}

func fmaSIMD(a, b, c, result []float64) {
	n := len(a)
	if hasAVX512 {
		simdLen := (n / 8) * 8
		fmaAVX512(a[:simdLen], b[:simdLen], c[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i]*b[i] + c[i]
		}
		return
	}
	if hasAVX2 && hasFMA {
		simdLen := (n / 4) * 4
		fmaAVX2(a[:simdLen], b[:simdLen], c[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i]*b[i] + c[i]
		}
		return
	}
	for i := 0; i < n; i++ {
		result[i] = a[i]*b[i] + c[i]
	}
}

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
	parallelizeGeneric(x, result, nativeLog)
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

func dispatchSinCosSIMDTo(x, sin, cos []float64) {
	n := len(x)
	if hasAVX2 {
		simdLen := (n / 4) * 4
		sinAVX2(x[:simdLen], sin[:simdLen])
		cosAVX2(x[:simdLen], cos[:simdLen])
		for i := simdLen; i < n; i++ {
			sin[i], cos[i] = Sincos(x[i])
		}
		return
	}
	parallelizeSinCos(x, sin, cos)
}
func dispatchSqrtSIMDTo(x, result []float64) { amd64SqrtSIMD(x, result) }
func dispatchAddSIMD(a, b, result []float64) { amd64AddSIMD(a, b, result) }
func dispatchSubSIMD(a, b, result []float64) { amd64SubSIMD(a, b, result) }
func dispatchMulSIMD(a, b, result []float64) { amd64MulSIMD(a, b, result) }
func dispatchDivSIMD(a, b, result []float64) { amd64DivSIMD(a, b, result) }
func dispatchAbsSIMD(x, result []float64) { amd64AbsSIMD(x, result) }
func dispatchNegSIMD(x, result []float64) { amd64NegSIMD(x, result) }
func dispatchInvSIMD(x, result []float64) { amd64InvSIMD(x, result) }

func dispatchAddScalarSIMD(a []float64, b float64, result []float64) {
	amd64AddScalarSIMD(a, b, result)
}

func dispatchMulScalarSIMD(a []float64, b float64, result []float64) {
	amd64MulScalarSIMD(a, b, result)
}
