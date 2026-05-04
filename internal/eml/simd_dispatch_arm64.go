//go:build arm64
// +build arm64

package eml



func arm64AddSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	addNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	for i := simdLen; i < n; i++ {
		result[i] = a[i] + b[i]
	}
}

func arm64SubSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	subNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	for i := simdLen; i < n; i++ {
		result[i] = a[i] - b[i]
	}
}

func arm64MulSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	mulNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	for i := simdLen; i < n; i++ {
		result[i] = a[i] * b[i]
	}
}

func arm64DivSIMD(a, b, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	divNEON(a[:simdLen], b[:simdLen], result[:simdLen])
	for i := simdLen; i < n; i++ {
		result[i] = a[i] / b[i]
	}
}

func arm64AddScalarSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	addScalarNEON(a[:simdLen], b, result[:simdLen])
	for i := simdLen; i < n; i++ {
		result[i] = a[i] + b
	}
}

func arm64MulScalarSIMD(a []float64, b float64, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	mulScalarNEON(a[:simdLen], b, result[:simdLen])
	for i := simdLen; i < n; i++ {
		result[i] = a[i] * b
	}
}

func arm64SqrtSIMD(a, result []float64) {
	n := len(a)
	simdLen := (n / 2) * 2
	sqrtNEON(a[:simdLen], result[:simdLen])
	for i := simdLen; i < n; i++ {
		result[i] = nativeSqrt(a[i])
	}
}

func arm64AbsSIMD(a, result []float64) {
	for i := range a {
		result[i] = nativeAbs(a[i])
	}
}

func arm64NegSIMD(a, result []float64) {
	for i := range a {
		result[i] = nativeNeg(a[i])
	}
}

func arm64InvSIMD(a, result []float64) {
	for i := range a {
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
	for i := 0; i < n; i++ {
		result[i] = a[i]*b[i] + c[i]
	}
}

func dispatchExpSIMDTo(x, result []float64) {
	parallelizeGeneric(x, result, nativeExp)
}

func dispatchLogSIMDTo(x, result []float64) {
	parallelizeGeneric(x, result, nativeLog)
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
