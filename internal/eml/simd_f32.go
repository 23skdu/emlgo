package eml

import (
	"fmt"
	"math"
	"runtime"
)

// parallelizeGenericF32 runs fn(src[lo:hi], dst[lo:hi]) in parallel workers
// for float32 slices. For short slices it runs inline.
func parallelizeGenericF32(src, dst []float32, fn func(src, dst []float32)) {
	n := len(src)
	if n < SmallCutoff {
		fn(src, dst)
		return
	}
	numWorkers := runtime.NumCPU()
	if numWorkers > n {
		numWorkers = n
	}
	chunkSize := (n + numWorkers - 1) / numWorkers
	done := make(chan struct{}, numWorkers)
	for w := 0; w < numWorkers; w++ {
		lo := w * chunkSize
		if lo >= n {
			done <- struct{}{}
			continue
		}
		hi := lo + chunkSize
		if hi > n {
			hi = n
		}
		go func(lo, hi int) {
			fn(src[lo:hi], dst[lo:hi])
			done <- struct{}{}
		}(lo, hi)
	}
	for w := 0; w < numWorkers; w++ {
		<-done
	}
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

// SqrtSIMDF32 computes element-wise sqrt(x) for a float32 slice.
func SqrtSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	parallelizeGenericF32(x, dst, func(src, d []float32) {
		for i, v := range src {
			d[i] = float32(math.Sqrt(float64(v)))
		}
	})
	return dst
}

// AddSIMDF32 computes element-wise a[i]+b[i] for float32 slices.
// Panics if len(a) != len(b).
func AddSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("AddSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	for i := range a {
		dst[i] = a[i] + b[i]
	}
	return dst
}

// MulSIMDF32 computes element-wise a[i]*b[i] for float32 slices.
// Panics if len(a) != len(b).
func MulSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("MulSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	for i := range a {
		dst[i] = a[i] * b[i]
	}
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

// SubSIMDF32 computes element-wise a[i]-b[i] for float32 slices.
// Panics if len(a) != len(b).
func SubSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("SubSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	for i := range a {
		dst[i] = a[i] - b[i]
	}
	return dst
}

// DivSIMDF32 computes element-wise a[i]/b[i] for float32 slices.
// Panics if len(a) != len(b).
func DivSIMDF32(a, b []float32) []float32 {
	if len(a) != len(b) {
		panic(fmt.Sprintf("DivSIMDF32: length mismatch %d vs %d", len(a), len(b)))
	}
	dst := make([]float32, len(a))
	for i := range a {
		dst[i] = a[i] / b[i]
	}
	return dst
}

// AbsSIMDF32 computes element-wise abs(x) for float32 slices.
func AbsSIMDF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	for i, v := range x {
		if v < 0 {
			dst[i] = -v
		} else {
			dst[i] = v
		}
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
