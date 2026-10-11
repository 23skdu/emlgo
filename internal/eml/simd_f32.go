package eml

import (
	"fmt"
	"math"
)

// parallelizeGenericF32 runs fn(src[lo:hi], dst[lo:hi]) in parallel workers
// for float32 slices using the shared worker pool. For short slices it runs inline.
func parallelizeGenericF32(src, dst []float32, fn func(src, dst []float32)) {
	n := len(src)
	ForEachChunk(n, func(start, end int) {
		fn(src[start:end], dst[start:end])
	})
}

// ExpSIMDF32 computes element-wise exp(x) for a float32 slice.
func ExpSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	parallelizeGenericF32(x, dst, func(src, d []float32) {
		for i, v := range src {
			d[i] = float32(math.Exp(float64(v)))
		}
	})
	return dst
}

// LogSIMDF32 computes element-wise log(x) for a float32 slice.
func LogSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	parallelizeGenericF32(x, dst, func(src, d []float32) {
		for i, v := range src {
			d[i] = float32(math.Log(float64(v)))
		}
	})
	return dst
}

// SqrtSIMDF32To computes element-wise sqrt(x) storing into dst.
func SqrtSIMDF32To(src, dst []float32) {
	if len(src) != len(dst) {
		panic("slice length mismatch")
	}
	dispatchSqrtSIMDF32To(src, dst)
}

// SqrtSIMDF32 computes element-wise sqrt(x) for a float32 slice.
func SqrtSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	SqrtSIMDF32To(x, dst)
	return dst
}

// AddSIMDF32To computes element-wise a[i]+b[i] storing into dst using SIMD when available.
func AddSIMDF32To(a, b, dst []float32) {
	if len(a) != len(b) || len(a) != len(dst) {
		panic("slice length mismatch")
	}
	dispatchAddSIMDF32To(a, b, dst)
}

// AddSIMDF32 computes element-wise a[i]+b[i] for float32 slices.
// Panics if len(a) != len(b).
func AddSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("AddSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	AddSIMDF32To(a, b, dst)
	return dst
}

// MulSIMDF32To computes element-wise a[i]*b[i] storing into dst using SIMD when available.
func MulSIMDF32To(a, b, dst []float32) {
	if len(a) != len(b) || len(a) != len(dst) {
		panic("slice length mismatch")
	}
	dispatchMulSIMDF32To(a, b, dst)
}

// MulSIMDF32 computes element-wise a[i]*b[i] for float32 slices.
// Panics if len(a) != len(b).
func MulSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("MulSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	MulSIMDF32To(a, b, dst)
	return dst
}

// AddScalarSIMDF32 adds scalar b to each element of a.
func AddScalarSIMDF32(a []float32, b float32) []float32 {
	dst := make([]float32, len(a))
	for i, v := range a {
		dst[i] = v + b
	}
	return dst
}

// MulScalarSIMDF32 multiplies each element of a by scalar b.
func MulScalarSIMDF32(a []float32, b float32) []float32 {
	dst := make([]float32, len(a))
	for i, v := range a {
		dst[i] = v * b
	}
	return dst
}

// SubSIMDF32To computes element-wise a[i]-b[i] storing into dst using SIMD when available.
func SubSIMDF32To(a, b, dst []float32) {
	if len(a) != len(b) || len(a) != len(dst) {
		panic("slice length mismatch")
	}
	dispatchSubSIMDF32To(a, b, dst)
}

// SubSIMDF32 computes element-wise a[i]-b[i] for float32 slices.
// Panics if len(a) != len(b).
func SubSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("SubSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	SubSIMDF32To(a, b, dst)
	return dst
}

// DivSIMDF32To computes element-wise a[i]/b[i] storing into dst using SIMD when available.
func DivSIMDF32To(a, b, dst []float32) {
	if len(a) != len(b) || len(a) != len(dst) {
		panic("slice length mismatch")
	}
	dispatchDivSIMDF32To(a, b, dst)
}

// DivSIMDF32 computes element-wise a[i]/b[i] for float32 slices.
// Panics if len(a) != len(b).
func DivSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("DivSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	DivSIMDF32To(a, b, dst)
	return dst
}

// AbsSIMDF32 computes element-wise abs(x) for float32 slices.
func AbsSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	for i, v := range x {
		dst[i] = math.Float32frombits(math.Float32bits(v) & 0x7fffffff)
	}
	return dst
}

// NegSIMDF32 computes element-wise -x for float32 slices.
func NegSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	for i, v := range x {
		dst[i] = -v
	}
	return dst
}

// InvSIMDF32 computes element-wise 1/x for float32 slices.
func InvSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	for i, v := range x {
		dst[i] = 1.0 / v
	}
	return dst
}

// SinSIMDF32 computes element-wise sin(x) for float32 slices.
func SinSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	parallelizeGenericF32(x, dst, func(src, d []float32) {
		for i, v := range src {
			d[i] = float32(math.Sin(float64(v)))
		}
	})
	return dst
}

// CosSIMDF32 computes element-wise cos(x) for float32 slices.
func CosSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	parallelizeGenericF32(x, dst, func(src, d []float32) {
		for i, v := range src {
			d[i] = float32(math.Cos(float64(v)))
		}
	})
	return dst
}

// TanSIMDF32 computes element-wise tan(x) for float32 slices.
func TanSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	parallelizeGenericF32(x, dst, func(src, d []float32) {
		for i, v := range src {
			d[i] = float32(math.Tan(float64(v)))
		}
	})
	return dst
}
