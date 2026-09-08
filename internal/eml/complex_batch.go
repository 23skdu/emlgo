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

// ComplexExpBatch returns a new slice containing the complex exponential of each element.
func ComplexExpBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	for i, v := range x {
		result[i] = cmplx.Exp(v)
	}
	return result
}

// ComplexLogBatch returns a new slice containing the complex logarithm of each element.
func ComplexLogBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	for i, v := range x {
		result[i] = cmplx.Log(v)
	}
	return result
}

// ComplexSinBatch returns a new slice containing the complex sine of each element.
func ComplexSinBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	for i, v := range x {
		result[i] = cmplx.Sin(v)
	}
	return result
}

// ComplexCosBatch returns a new slice containing the complex cosine of each element.
func ComplexCosBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	for i, v := range x {
		result[i] = cmplx.Cos(v)
	}
	return result
}

// ComplexTanBatch returns a new slice containing the complex tangent of each element.
func ComplexTanBatch(x []complex128) []complex128 {
	result := make([]complex128, len(x))
	for i, v := range x {
		result[i] = cmplx.Tan(v)
	}
	return result
}
