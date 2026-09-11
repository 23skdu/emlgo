package eml

import "math/cmplx"

// ComplexBatch applies Complex(x[i], y[i]) to each element pair and stores in result.
func ComplexBatch(x, y, result []complex128) {
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i := range x {
		result[i] = cmplx.Exp(x[i]) - cmplx.Log(y[i])
	}
}

// ComplexBatchC64 applies Complex64(x[i], y[i]) to each element pair and stores in result.
func ComplexBatchC64(x, y, result []complex64) {
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i := range x {
		result[i] = Complex64(x[i], y[i])
	}
}

// ComplexExpBatch returns a new slice containing the complex exponential of each element.
func ComplexExpBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	ComplexExpBatchTo(x, result)
	return result
}

func ComplexExpBatchTo(x, result []complex128) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = cmplx.Exp(v)
	}
}

func ComplexExpBatchC64(x []complex64) []complex64 {
	result := make([]complex64, len(x))
	ComplexExpBatchToC64(x, result)
	return result
}

func ComplexExpBatchToC64(x, result []complex64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = complex64(cmplx.Exp(complex128(v)))
	}
}

// ComplexLogBatch returns a new slice containing the complex logarithm of each element.
func ComplexLogBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	ComplexLogBatchTo(x, result)
	return result
}

func ComplexLogBatchTo(x, result []complex128) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = cmplx.Log(v)
	}
}

func ComplexLogBatchC64(x []complex64) []complex64 {
	result := make([]complex64, len(x))
	ComplexLogBatchToC64(x, result)
	return result
}

func ComplexLogBatchToC64(x, result []complex64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = complex64(cmplx.Log(complex128(v)))
	}
}

// ComplexSinBatch returns a new slice containing the complex sine of each element.
func ComplexSinBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	ComplexSinBatchTo(x, result)
	return result
}

func ComplexSinBatchTo(x, result []complex128) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = cmplx.Sin(v)
	}
}

func ComplexSinBatchC64(x []complex64) []complex64 {
	result := make([]complex64, len(x))
	ComplexSinBatchToC64(x, result)
	return result
}

func ComplexSinBatchToC64(x, result []complex64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = complex64(cmplx.Sin(complex128(v)))
	}
}

// ComplexCosBatch returns a new slice containing the complex cosine of each element.
func ComplexCosBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	ComplexCosBatchTo(x, result)
	return result
}

func ComplexCosBatchTo(x, result []complex128) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = cmplx.Cos(v)
	}
}

func ComplexCosBatchC64(x []complex64) []complex64 {
	result := make([]complex64, len(x))
	ComplexCosBatchToC64(x, result)
	return result
}

func ComplexCosBatchToC64(x, result []complex64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = complex64(cmplx.Cos(complex128(v)))
	}
}

// ComplexTanBatch returns a new slice containing the complex tangent of each element.
func ComplexTanBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	ComplexTanBatchTo(x, result)
	return result
}

func ComplexTanBatchTo(x, result []complex128) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = cmplx.Tan(v)
	}
}

func ComplexTanBatchC64(x []complex64) []complex64 {
	result := make([]complex64, len(x))
	ComplexTanBatchToC64(x, result)
	return result
}

func ComplexTanBatchToC64(x, result []complex64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = complex64(cmplx.Tan(complex128(v)))
	}
}

// ComplexSqrtBatch returns a new slice containing the complex square root of each element.
func ComplexSqrtBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	ComplexSqrtBatchTo(x, result)
	return result
}

func ComplexSqrtBatchTo(x, result []complex128) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = cmplx.Sqrt(v)
	}
}

func ComplexSqrtBatchC64(x []complex64) []complex64 {
	result := make([]complex64, len(x))
	ComplexSqrtBatchToC64(x, result)
	return result
}

func ComplexSqrtBatchToC64(x, result []complex64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		result[i] = complex64(cmplx.Sqrt(complex128(v)))
	}
}

// ComplexDotProduct returns the Hermitian inner product of two complex128 slices: sum(a[i] * conj(b[i])).
func ComplexDotProduct(a, b []complex128) complex128 {
	if len(a) != len(b) {
		panic("slice length mismatch")
	}
	var sum complex128
	for i := range a {
		sum += a[i] * cmplx.Conj(b[i])
	}
	return sum
}

// ComplexDotProductC64 returns the Hermitian inner product of two complex64 slices.
func ComplexDotProductC64(a, b []complex64) complex64 {
	if len(a) != len(b) {
		panic("slice length mismatch")
	}
	var sum complex128
	for i := range a {
		sum += complex128(a[i]) * cmplx.Conj(complex128(b[i]))
	}
	return complex64(sum)
}
