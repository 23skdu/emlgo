package eml

import (
	"runtime"
	"sync"
)

var (
	cpuNum = runtime.NumCPU()
	hasSSE4   bool
	hasAVX2   bool
	hasAVX512 bool
	hasNeon   bool
	hasNeonDot bool
	hasSVE    bool
	hasFMA    bool
	hasAVXVNNI bool
)

func init() {
	detectSIMD()
	detectCacheTopology()
}

func detectSIMD() {
	detectPlatformSIMD()
}

func detectCacheTopology() {
	_ = cpuNum
}

// GetParallelChunkSize returns the ideal chunk size for parallel processing of n elements.
func GetParallelChunkSize(n int) int {
	if n < SmallCutoff {
		return n
	}
	chunkSize := (n + cpuNum - 1) / cpuNum
	if chunkSize > LargeCutoff {
		return LargeCutoff
	}
	return chunkSize
}


// HasSSE4 reports whether the CPU supports SSE4 instructions.
func HasSSE4() bool { return hasSSE4 }
// HasAVX2 reports whether the CPU supports AVX2 instructions.
func HasAVX2() bool { return hasAVX2 }
// HasAVX512 reports whether the CPU supports AVX-512 instructions.
func HasAVX512() bool { return hasAVX512 }
// HasNeon reports whether the CPU supports ARM Neon instructions.
func HasNeon() bool { return hasNeon }
// HasNeonDot reports whether the CPU supports ARM Neon dot product instructions.
func HasNeonDot() bool { return hasNeonDot }
// HasSVE reports whether the CPU supports ARM SVE instructions.
func HasSVE() bool { return hasSVE }
// HasFMA reports whether the CPU supports FMA instructions.
func HasFMA() bool { return hasFMA }
// HasAVXVNNI reports whether the CPU supports AVX-VNNI instructions.
func HasAVXVNNI() bool { return hasAVXVNNI }


// FmaScalar returns a * b + c.
//go:inline
func FmaScalar(a, b, c float64) float64 { return fmaScalar(a, b, c) }

// SqrtScalar returns the square root of x.
//go:inline
func SqrtScalar(x float64) float64 { return sqrtScalar(x) }

// AbsScalar returns the absolute value of x.
//go:inline
func AbsScalar(x float64) float64 { return absScalar(x) }

// NegScalar returns the negation of x.
//go:inline
func NegScalar(x float64) float64 { return negScalar(x) }


// SIMD computes Exp(x[i]) - Log(y[i]) for each element and stores it in result.
func SIMD(x, y, result []float64) {
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}
	if len(x) == 0 { return }
	emlSIMD(x, y, result)
}


// ExpSIMD returns a new slice containing the exponential of each element in x.
func ExpSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	ExpSIMDTo(x, result)
	return result
}


// ExpSIMDTo computes the exponential of each element in x and stores it in result.
func ExpSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	expSIMDTo(x, result)
}


func expSIMDTo(x, result []float64) {
	dispatchExpSIMDTo(x, result)
}

// LogSIMD returns a new slice containing the natural logarithm of each element in x.
func LogSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	LogSIMDTo(x, result)
	return result
}


// LogSIMDTo computes the natural logarithm of each element in x and stores it in result.
func LogSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	logSIMDTo(x, result)
}


func logSIMDTo(x, result []float64) {
	dispatchLogSIMDTo(x, result)
}

// SqrtSIMD returns a new slice containing the square root of each element in x.
func SqrtSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	SqrtSIMDTo(x, result)
	return result
}


// SqrtSIMDTo computes the square root of each element in x and stores it in result.
func SqrtSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	sqrtSIMDTo(x, result)
}


func sqrtSIMDTo(x, result []float64) {
	dispatchSqrtSIMDTo(x, result)
}

// SinSIMD returns a new slice containing the sine of each element in x.
func SinSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	SinSIMDTo(x, result)
	return result
}


// SinSIMDTo computes the sine of each element in x and stores it in result.
func SinSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	sinSIMDTo(x, result)
}


func sinSIMDTo(x, result []float64) {
	dispatchSinSIMDTo(x, result)
}

// CosSIMD returns a new slice containing the cosine of each element in x.
func CosSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	CosSIMDTo(x, result)
	return result
}


// CosSIMDTo computes the cosine of each element in x and stores it in result.
func CosSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	cosSIMDTo(x, result)
}


func cosSIMDTo(x, result []float64) {
	dispatchCosSIMDTo(x, result)
}

// TanSIMD returns a new slice containing the tangent of each element in x.
func TanSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	TanSIMDTo(x, result)
	return result
}


// TanSIMDTo computes the tangent of each element in x and stores it in result.
func TanSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	tanSIMDTo(x, result)
}


func tanSIMDTo(x, result []float64) {
	dispatchTanSIMDTo(x, result)
}

// SinCosSIMD returns two new slices containing the sine and cosine of each element in x.
func SinCosSIMD(x []float64) (sin, cos []float64) {
	sin = make([]float64, len(x))
	cos = make([]float64, len(x))
	SinCosSIMDTo(x, sin, cos)
	return
}


// SinCosSIMDTo computes the sine and cosine of each element in x and stores them in sin and cos respectively.
func SinCosSIMDTo(x, sin, cos []float64) {
	if len(x) != len(sin) || len(x) != len(cos) { panic("slice length mismatch") }
	sincosSIMDTo(x, sin, cos)
}


func sincosSIMDTo(x, sin, cos []float64) {
	dispatchSinCosSIMDTo(x, sin, cos)
}

// AddSIMD returns a new slice containing the sum of elements in a and b.
func AddSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchAddSIMD(a, b, result)
	return result
}


// SubSIMD returns a new slice containing the difference of elements in a and b.
func SubSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchSubSIMD(a, b, result)
	return result
}


// MulSIMD returns a new slice containing the product of elements in a and b.
func MulSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchMulSIMD(a, b, result)
	return result
}


// DivSIMD returns a new slice containing the quotient of elements in a and b.
func DivSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchDivSIMD(a, b, result)
	return result
}


// AbsSIMD returns a new slice containing the absolute value of each element in x.
func AbsSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	AbsSIMDTo(x, result)
	return result
}


// AbsSIMDTo computes the absolute value of each element in x and stores it in result.
func AbsSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	absSIMD(x, result)
}


func absSIMD(x, result []float64) {
	dispatchAbsSIMD(x, result)
}

// NegSIMD returns a new slice containing the negation of each element in x.
func NegSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	NegSIMDTo(x, result)
	return result
}


// NegSIMDTo computes the negation of each element in x and stores it in result.
func NegSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	negSIMD(x, result)
}


func negSIMD(x, result []float64) {
	dispatchNegSIMD(x, result)
}

// InvSIMD returns a new slice containing the inverse of each element in x.
func InvSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	InvSIMDTo(x, result)
	return result
}


// InvSIMDTo computes the inverse of each element in x and stores it in result.
func InvSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	invSIMD(x, result)
}


func invSIMD(x, result []float64) {
	dispatchInvSIMD(x, result)
}

// SinhBatch returns a new slice containing the hyperbolic sine of each element in x.
func SinhBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Sinh)
	return res
}


// CoshBatch returns a new slice containing the hyperbolic cosine of each element in x.
func CoshBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Cosh)
	return res
}


// TanhBatch returns a new slice containing the hyperbolic tangent of each element in x.
func TanhBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Tanh)
	return res
}


// AsinhBatch returns a new slice containing the inverse hyperbolic sine of each element in x.
func AsinhBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Asinh)
	return res
}


// AcoshBatch returns a new slice containing the inverse hyperbolic cosine of each element in x.
func AcoshBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Acosh)
	return res
}


// AtanhBatch returns a new slice containing the inverse hyperbolic tangent of each element in x.
func AtanhBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Atanh)
	return res
}


func parallelizeGeneric(x, result []float64, fn func(float64) float64) {
	n := len(x)
	if n == 0 { return }
	chunkSize := GetParallelChunkSize(n)
	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n { end = n }
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = fn(x[j])
			}
		}(i, end)
	}
	wg.Wait()
}

func parallelizeSinCos(x, sin, cos []float64) {
	n := len(x)
	if n == 0 { return }
	chunkSize := GetParallelChunkSize(n)
	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n { end = n }
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				sin[j], cos[j] = Sincos(x[j])
			}
		}(i, end)
	}
	wg.Wait()
}

const (
	// SmallCutoff is the threshold below which operations are performed sequentially.
	SmallCutoff = 256
	// LargeCutoff is the maximum chunk size for parallel operations.
	LargeCutoff = 4096
)


// ErrLengthMismatch is returned when slice lengths do not match.
var ErrLengthMismatch = Error("slice length mismatch")

// Error represents an EML package error.
type Error string

func (e Error) Error() string { return string(e) }


// Batch applies a VectorFunc to x and y.
func Batch(x, y []float64, fn VectorFunc) error {
	if len(x) != len(y) { return ErrLengthMismatch }
	result := make([]float64, len(x))
	return fn(x, y, result)
}


// VectorFunc is a function that operates on two input slices and one result slice.
type VectorFunc func(x, y, result []float64) error


// AddScalarSIMD returns a new slice containing each element in a plus b.
func AddScalarSIMD(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	AddScalarSIMDTo(a, b, result)
	return result
}


// AddScalarSIMDTo adds b to each element in a and stores it in result.
func AddScalarSIMDTo(a []float64, b float64, result []float64) {
	if len(a) != len(result) { panic("slice length mismatch") }
	dispatchAddScalarSIMD(a, b, result)
}


// MulScalarSIMD returns a new slice containing each element in a multiplied by b.
func MulScalarSIMD(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	MulScalarSIMDTo(a, b, result)
	return result
}


// MulScalarSIMDTo multiplies each element in a by b and stores it in result.
func MulScalarSIMDTo(a []float64, b float64, result []float64) {
	if len(a) != len(result) { panic("slice length mismatch") }
	dispatchMulScalarSIMD(a, b, result)
}


// L1TileSize is the suggested tile size for L1 cache optimizations.
const L1TileSize = 32768

