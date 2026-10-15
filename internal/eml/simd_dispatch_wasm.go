//go:build wasm && !purego

package eml

import (
	"math"
)

func emlSIMD(x, y, result []float64) {
	scalarEml(x, y, result)
}

func detectPlatformSIMD() {
	detectWasmSIMD()
}

func scalarEml(x, y, result []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		result[i] = nativeExp(x[i]) - nativeLog(y[i])
	}
}

func fmaSIMD(a, b, c, result []float64) {
	fmaWasmSIMD(a, b, c, result)
}

func dispatchExpSIMDTo(x, result []float64)      { expWasmSIMD(x, result) }
func dispatchLogSIMDTo(x, result []float64)      { logWasmSIMD(x, result) }
func dispatchSinSIMDTo(x, result []float64)      { sinWasmSIMD(x, result) }
func dispatchCosSIMDTo(x, result []float64)      { cosWasmSIMD(x, result) }
func dispatchTanSIMDTo(x, result []float64)      { tanWasmSIMD(x, result) }
func dispatchSinCosSIMDTo(x, sin, cos []float64) { sincosWasmSIMD(x, sin, cos) }
func dispatchSqrtSIMDTo(x, result []float64)     { sqrtWasmSIMD(x, result) }
func dispatchAddSIMD(a, b, result []float64)     { addWasmSIMD(a, b, result) }
func dispatchSubSIMD(a, b, result []float64)     { subWasmSIMD(a, b, result) }
func dispatchMulSIMD(a, b, result []float64)     { mulWasmSIMD(a, b, result) }
func dispatchDivSIMD(a, b, result []float64)     { divWasmSIMD(a, b, result) }
func dispatchAbsSIMD(x, result []float64)        { absWasmSIMD(x, result) }
func dispatchNegSIMD(x, result []float64)        { negWasmSIMD(x, result) }
func dispatchInvSIMD(x, result []float64)        { invWasmSIMD(x, result) }

func dispatchAddScalarSIMD(a []float64, b float64, result []float64) {
	addScalarWasmSIMD(a, b, result)
}

func dispatchMulScalarSIMD(a []float64, b float64, result []float64) {
	mulScalarWasmSIMD(a, b, result)
}

func dispatchDotInt8SIMD(a []uint8, b []int8) int32 {
	var sum int32
	for i := range a {
		sum += int32(a[i]) * int32(b[i])
	}
	return sum
}

func dispatchAddSIMDF32To(a, b, dst []float32) {
	for i := range a {
		dst[i] = a[i] + b[i]
	}
}

func dispatchSubSIMDF32To(a, b, dst []float32) {
	for i := range a {
		dst[i] = a[i] - b[i]
	}
}

func dispatchMulSIMDF32To(a, b, dst []float32) {
	for i := range a {
		dst[i] = a[i] * b[i]
	}
}

func dispatchDivSIMDF32To(a, b, dst []float32) {
	for i := range a {
		dst[i] = a[i] / b[i]
	}
}

func dispatchSqrtSIMDF32To(src, dst []float32) {
	for i := range src {
		dst[i] = float32(nativeSqrt(float64(src[i])))
	}
}

func dispatchAddScalarSIMDF32To(a []float32, b float32, dst []float32) {
	for i := range a {
		dst[i] = a[i] + b
	}
}

func dispatchMulScalarSIMDF32To(a []float32, b float32, dst []float32) {
	for i := range a {
		dst[i] = a[i] * b
	}
}

func dispatchAbsSIMDF32To(src, dst []float32) {
	for i := range src {
		dst[i] = math.Float32frombits(math.Float32bits(src[i]) & 0x7fffffff)
	}
}

func dispatchNegSIMDF32To(src, dst []float32) {
	for i := range src {
		dst[i] = -src[i]
	}
}

func dispatchInvSIMDF32To(src, dst []float32) {
	for i := range src {
		dst[i] = 1.0 / src[i]
	}
}

func dispatchFmaSIMDF32To(a, b, c, dst []float32) {
	for i := range a {
		dst[i] = float32(math.FMA(float64(a[i]), float64(b[i]), float64(c[i])))
	}
}

func dispatchFloorSIMDTo(src, dst []float64) {
	for i := range src {
		dst[i] = math.Floor(src[i])
	}
}

func dispatchCeilSIMDTo(src, dst []float64) {
	for i := range src {
		dst[i] = math.Ceil(src[i])
	}
}

func dispatchTruncSIMDTo(src, dst []float64) {
	for i := range src {
		dst[i] = math.Trunc(src[i])
	}
}

func dispatchRoundSIMDTo(src, dst []float64) {
	for i := range src {
		dst[i] = math.RoundToEven(src[i])
	}
}

func dispatchMaxSIMDTo(a, b, dst []float64) {
	for i := range a {
		if math.IsNaN(a[i]) || math.IsNaN(b[i]) {
			dst[i] = math.NaN()
		} else if a[i] > b[i] {
			dst[i] = a[i]
		} else {
			dst[i] = b[i]
		}
	}
}

func dispatchMinSIMDTo(a, b, dst []float64) {
	for i := range a {
		if math.IsNaN(a[i]) || math.IsNaN(b[i]) {
			dst[i] = math.NaN()
		} else if a[i] < b[i] {
			dst[i] = a[i]
		} else {
			dst[i] = b[i]
		}
	}
}

var _ = parallelizeSinCos

