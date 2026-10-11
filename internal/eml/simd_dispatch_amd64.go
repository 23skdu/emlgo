//go:build amd64 && !purego
// +build amd64,!purego

package eml

import "math"

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
	for i := 0; i < n; i++ {
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
	for i := 0; i < n; i++ {
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
	for i := 0; i < n; i++ {
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
	for i := 0; i < n; i++ {
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

func dispatchSinCosSIMDTo(x, sin, cos []float64) {
	parallelizeSinCos(x, sin, cos)
}
func dispatchSqrtSIMDTo(x, result []float64) { amd64SqrtSIMD(x, result) }
func dispatchAddSIMD(a, b, result []float64) { amd64AddSIMD(a, b, result) }
func dispatchSubSIMD(a, b, result []float64) { amd64SubSIMD(a, b, result) }
func dispatchMulSIMD(a, b, result []float64) { amd64MulSIMD(a, b, result) }
func dispatchDivSIMD(a, b, result []float64) { amd64DivSIMD(a, b, result) }
func dispatchAbsSIMD(x, result []float64)    { amd64AbsSIMD(x, result) }
func dispatchNegSIMD(x, result []float64)    { amd64NegSIMD(x, result) }
func dispatchInvSIMD(x, result []float64)    { amd64InvSIMD(x, result) }

func dispatchAddScalarSIMD(a []float64, b float64, result []float64) {
	amd64AddScalarSIMD(a, b, result)
}

func dispatchMulScalarSIMD(a []float64, b float64, result []float64) {
	amd64MulScalarSIMD(a, b, result)
}

func dispatchAddSatInt8SIMD(a, b, result []int8) {
	n := len(a)
	if hasAVX512 && n >= 64 {
		simdLen := (n / 64) * 64
		addSatInt8AVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			v := int32(a[i]) + int32(b[i])
			if v > 127 {
				v = 127
			} else if v < -128 {
				v = -128
			}
			result[i] = int8(v)
		}
		return
	}
	if hasAVX2 && n >= 32 {
		simdLen := (n / 32) * 32
		addSatInt8AVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			v := int32(a[i]) + int32(b[i])
			if v > 127 {
				v = 127
			} else if v < -128 {
				v = -128
			}
			result[i] = int8(v)
		}
		return
	}
	for i := 0; i < n; i++ {
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
	if hasAVX512 && n >= 64 {
		simdLen := (n / 64) * 64
		subSatInt8AVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			v := int32(a[i]) - int32(b[i])
			if v > 127 {
				v = 127
			} else if v < -128 {
				v = -128
			}
			result[i] = int8(v)
		}
		return
	}
	if hasAVX2 && n >= 32 {
		simdLen := (n / 32) * 32
		subSatInt8AVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			v := int32(a[i]) - int32(b[i])
			if v > 127 {
				v = 127
			} else if v < -128 {
				v = -128
			}
			result[i] = int8(v)
		}
		return
	}
	for i := 0; i < n; i++ {
		v := int32(a[i]) - int32(b[i])
		if v > 127 {
			v = 127
		} else if v < -128 {
			v = -128
		}
		result[i] = int8(v)
	}
}

func dispatchDotInt8SIMD(a []uint8, b []int8) int32 {
	n := len(a)
	var sum int32
	i := 0
	if hasAVX512 && hasAVX512VNNI && n >= 64 {
		simdLen := (n / 64) * 64
		sum = dotInt8AVX512VNNI(a[:simdLen], b[:simdLen])
		i = simdLen
	}
	for ; i < n; i++ {
		sum += int32(a[i]) * int32(b[i])
	}
	return sum
}

func dispatchAddSIMDF32To(a, b, dst []float32) {
	n := len(a)
	i := 0
	if hasAVX512 && n >= 16 {
		simdLen := (n / 16) * 16
		addAVX512F32(a[:simdLen], b[:simdLen], dst[:simdLen])
		i = simdLen
	}
	for ; i+7 < n; i += 8 {
		dst[i] = a[i] + b[i]
		dst[i+1] = a[i+1] + b[i+1]
		dst[i+2] = a[i+2] + b[i+2]
		dst[i+3] = a[i+3] + b[i+3]
		dst[i+4] = a[i+4] + b[i+4]
		dst[i+5] = a[i+5] + b[i+5]
		dst[i+6] = a[i+6] + b[i+6]
		dst[i+7] = a[i+7] + b[i+7]
	}
	for ; i < n; i++ {
		dst[i] = a[i] + b[i]
	}
}

func dispatchSubSIMDF32To(a, b, dst []float32) {
	n := len(a)
	i := 0
	if hasAVX512 && n >= 16 {
		simdLen := (n / 16) * 16
		subAVX512F32(a[:simdLen], b[:simdLen], dst[:simdLen])
		i = simdLen
	}
	for ; i+7 < n; i += 8 {
		dst[i] = a[i] - b[i]
		dst[i+1] = a[i+1] - b[i+1]
		dst[i+2] = a[i+2] - b[i+2]
		dst[i+3] = a[i+3] - b[i+3]
		dst[i+4] = a[i+4] - b[i+4]
		dst[i+5] = a[i+5] - b[i+5]
		dst[i+6] = a[i+6] - b[i+6]
		dst[i+7] = a[i+7] - b[i+7]
	}
	for ; i < n; i++ {
		dst[i] = a[i] - b[i]
	}
}

func dispatchMulSIMDF32To(a, b, dst []float32) {
	n := len(a)
	i := 0
	if hasAVX512 && n >= 16 {
		simdLen := (n / 16) * 16
		mulAVX512F32(a[:simdLen], b[:simdLen], dst[:simdLen])
		i = simdLen
	}
	for ; i+7 < n; i += 8 {
		dst[i] = a[i] * b[i]
		dst[i+1] = a[i+1] * b[i+1]
		dst[i+2] = a[i+2] * b[i+2]
		dst[i+3] = a[i+3] * b[i+3]
		dst[i+4] = a[i+4] * b[i+4]
		dst[i+5] = a[i+5] * b[i+5]
		dst[i+6] = a[i+6] * b[i+6]
		dst[i+7] = a[i+7] * b[i+7]
	}
	for ; i < n; i++ {
		dst[i] = a[i] * b[i]
	}
}

func dispatchDivSIMDF32To(a, b, dst []float32) {
	n := len(a)
	i := 0
	if hasAVX512 && n >= 16 {
		simdLen := (n / 16) * 16
		divAVX512F32(a[:simdLen], b[:simdLen], dst[:simdLen])
		i = simdLen
	}
	for ; i+7 < n; i += 8 {
		dst[i] = a[i] / b[i]
		dst[i+1] = a[i+1] / b[i+1]
		dst[i+2] = a[i+2] / b[i+2]
		dst[i+3] = a[i+3] / b[i+3]
		dst[i+4] = a[i+4] / b[i+4]
		dst[i+5] = a[i+5] / b[i+5]
		dst[i+6] = a[i+6] / b[i+6]
		dst[i+7] = a[i+7] / b[i+7]
	}
	for ; i < n; i++ {
		dst[i] = a[i] / b[i]
	}
}

func dispatchSqrtSIMDF32To(src, dst []float32) {
	n := len(src)
	i := 0
	if hasAVX512 && n >= 16 {
		simdLen := (n / 16) * 16
		sqrtAVX512F32(src[:simdLen], dst[:simdLen])
		i = simdLen
	}
	for ; i < n; i++ {
		dst[i] = float32(math.Sqrt(float64(src[i])))
	}
}

