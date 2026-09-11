package fastmath

import (
	"math"
	"math/cmplx"

	"github.com/emlgo/eml/internal/eml"
)

// ============================================================================
// Float64 Functions
// ============================================================================

// Sqrt returns the square root of x.
// Delegates to math.Sqrt which the Go compiler handles as an intrinsic on amd64,
// producing a single SQRTSD instruction with zero function-call overhead.
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

// FMA returns x*y + z, computed with only one rounding.
// It uses direct assembly (VFMADD/FMADD) and is significantly faster than math.FMA.
func FMA(x, y, z float64) float64 {
	return eml.FmaScalar(x, y, z)
}

// SqrtBatch computes element-wise square root for float64 slices using SIMD.
func SqrtBatch(x []float64) []float64 {
	return eml.SqrtSIMD(x)
}

// SqrtBatchTo computes element-wise square root for float64 slices storing results in dst.
func SqrtBatchTo(x, dst []float64) {
	eml.SqrtSIMDTo(x, dst)
}

// Exp returns e^x using fast degree-4 Remez minimax polynomial approximation.
func Exp(x float64) float64 {
	return FastExp(x)
}

// Sin returns the sine of x.
func Sin(x float64) float64 {
	return math.Sin(x)
}

// Cos returns the cosine of x.
func Cos(x float64) float64 {
	return math.Cos(x)
}

// Log returns the natural logarithm of x using fast Chebyshev polynomial approximation.
func Log(x float64) float64 {
	return FastLog(x)
}

// ============================================================================
// Float32 Functions
// ============================================================================

// SqrtF32 returns the square root of float32 x using compiler intrinsic.
func SqrtF32(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// SqrtBatchF32 computes element-wise square root for float32 slices.
func SqrtBatchF32(x []float32) []float32 {
	dst := make([]float32, len(x))
	SqrtBatchToF32(x, dst)
	return dst
}

// SqrtBatchToF32 computes element-wise square root for float32 slices storing in dst.
func SqrtBatchToF32(x, dst []float32) {
	if len(x) != len(dst) {
		panic("slice length mismatch")
	}
	for i, v := range x {
		dst[i] = float32(math.Sqrt(float64(v)))
	}
}

// ExpF32 returns e^x for float32 using fast polynomial approximation.
func ExpF32(x float32) float32 {
	return FastExpF32(x)
}

// SinF32 returns sin(x) for float32.
func SinF32(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

// CosF32 returns cos(x) for float32.
func CosF32(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

// LogF32 returns ln(x) for float32 using fast polynomial approximation.
func LogF32(x float32) float32 {
	return FastLogF32(x)
}

// ============================================================================
// Complex64 Functions
// ============================================================================

// SqrtC64 returns the complex square root of complex64 z.
func SqrtC64(z complex64) complex64 {
	return complex64(cmplx.Sqrt(complex128(z)))
}

// SqrtBatchC64 computes element-wise square root for complex64 slices.
func SqrtBatchC64(x []complex64) []complex64 {
	return eml.ComplexSqrtBatchC64(x)
}

// SqrtBatchToC64 computes element-wise square root for complex64 slices storing in dst.
func SqrtBatchToC64(x, dst []complex64) {
	eml.ComplexSqrtBatchToC64(x, dst)
}

// ExpC64 returns e^z for complex64.
func ExpC64(z complex64) complex64 {
	return complex64(cmplx.Exp(complex128(z)))
}

// ExpBatchC64 computes element-wise exp for complex64 slices.
func ExpBatchC64(x []complex64) []complex64 {
	return eml.ComplexExpBatchC64(x)
}

// ExpBatchToC64 computes element-wise exp for complex64 slices storing in dst.
func ExpBatchToC64(x, dst []complex64) {
	eml.ComplexExpBatchToC64(x, dst)
}

// LogC64 returns ln(z) for complex64.
func LogC64(z complex64) complex64 {
	return complex64(cmplx.Log(complex128(z)))
}

// LogBatchC64 computes element-wise log for complex64 slices.
func LogBatchC64(x []complex64) []complex64 {
	return eml.ComplexLogBatchC64(x)
}

// LogBatchToC64 computes element-wise log for complex64 slices storing in dst.
func LogBatchToC64(x, dst []complex64) {
	eml.ComplexLogBatchToC64(x, dst)
}

// SinC64 returns sin(z) for complex64.
func SinC64(z complex64) complex64 {
	return complex64(cmplx.Sin(complex128(z)))
}

// SinBatchC64 computes element-wise sin for complex64 slices.
func SinBatchC64(x []complex64) []complex64 {
	return eml.ComplexSinBatchC64(x)
}

// SinBatchToC64 computes element-wise sin for complex64 slices storing in dst.
func SinBatchToC64(x, dst []complex64) {
	eml.ComplexSinBatchToC64(x, dst)
}

// CosC64 returns cos(z) for complex64.
func CosC64(z complex64) complex64 {
	return complex64(cmplx.Cos(complex128(z)))
}

// CosBatchC64 computes element-wise cos for complex64 slices.
func CosBatchC64(x []complex64) []complex64 {
	return eml.ComplexCosBatchC64(x)
}

// CosBatchToC64 computes element-wise cos for complex64 slices storing in dst.
func CosBatchToC64(x, dst []complex64) {
	eml.ComplexCosBatchToC64(x, dst)
}

// TanC64 returns tan(z) for complex64.
func TanC64(z complex64) complex64 {
	return complex64(cmplx.Tan(complex128(z)))
}

// TanBatchC64 computes element-wise tan for complex64 slices.
func TanBatchC64(x []complex64) []complex64 {
	return eml.ComplexTanBatchC64(x)
}

// TanBatchToC64 computes element-wise tan for complex64 slices storing in dst.
func TanBatchToC64(x, dst []complex64) {
	eml.ComplexTanBatchToC64(x, dst)
}

// ============================================================================
// Complex128 Functions
// ============================================================================

// SqrtC128 returns the complex square root of complex128 z.
func SqrtC128(z complex128) complex128 {
	return cmplx.Sqrt(z)
}

// SqrtBatchC128 computes element-wise square root for complex128 slices.
func SqrtBatchC128(x []complex128) []complex128 {
	return eml.ComplexSqrtBatch(x)
}

// SqrtBatchToC128 computes element-wise square root for complex128 slices storing in dst.
func SqrtBatchToC128(x, dst []complex128) {
	eml.ComplexSqrtBatchTo(x, dst)
}

// ExpC128 returns e^z for complex128.
func ExpC128(z complex128) complex128 {
	return cmplx.Exp(z)
}

// ExpBatchC128 computes element-wise exp for complex128 slices.
func ExpBatchC128(x []complex128) []complex128 {
	return eml.ComplexExpBatch(x)
}

// ExpBatchToC128 computes element-wise exp for complex128 slices storing in dst.
func ExpBatchToC128(x, dst []complex128) {
	eml.ComplexExpBatchTo(x, dst)
}

// LogC128 returns ln(z) for complex128.
func LogC128(z complex128) complex128 {
	return cmplx.Log(z)
}

// LogBatchC128 computes element-wise log for complex128 slices.
func LogBatchC128(x []complex128) []complex128 {
	return eml.ComplexLogBatch(x)
}

// LogBatchToC128 computes element-wise log for complex128 slices storing in dst.
func LogBatchToC128(x, dst []complex128) {
	eml.ComplexLogBatchTo(x, dst)
}

// SinC128 returns sin(z) for complex128.
func SinC128(z complex128) complex128 {
	return cmplx.Sin(z)
}

// SinBatchC128 computes element-wise sin for complex128 slices.
func SinBatchC128(x []complex128) []complex128 {
	return eml.ComplexSinBatch(x)
}

// SinBatchToC128 computes element-wise sin for complex128 slices storing in dst.
func SinBatchToC128(x, dst []complex128) {
	eml.ComplexSinBatchTo(x, dst)
}

// CosC128 returns cos(z) for complex128.
func CosC128(z complex128) complex128 {
	return cmplx.Cos(z)
}

// CosBatchC128 computes element-wise cos for complex128 slices.
func CosBatchC128(x []complex128) []complex128 {
	return eml.ComplexCosBatch(x)
}

// CosBatchToC128 computes element-wise cos for complex128 slices storing in dst.
func CosBatchToC128(x, dst []complex128) {
	eml.ComplexCosBatchTo(x, dst)
}

// TanC128 returns tan(z) for complex128.
func TanC128(z complex128) complex128 {
	return cmplx.Tan(z)
}

// TanBatchC128 computes element-wise tan for complex128 slices.
func TanBatchC128(x []complex128) []complex128 {
	return eml.ComplexTanBatch(x)
}

// TanBatchToC128 computes element-wise tan for complex128 slices storing in dst.
func TanBatchToC128(x, dst []complex128) {
	eml.ComplexTanBatchTo(x, dst)
}
