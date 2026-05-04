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

func HasSSE4() bool { return hasSSE4 }
func HasAVX2() bool { return hasAVX2 }
func HasAVX512() bool { return hasAVX512 }
func HasNeon() bool { return hasNeon }
func HasNeonDot() bool { return hasNeonDot }
func HasSVE() bool { return hasSVE }
func HasFMA() bool { return hasFMA }
func HasAVXVNNI() bool { return hasAVXVNNI }

//go:inline
func FmaScalar(a, b, c float64) float64 { return fmaScalar(a, b, c) }
//go:inline
func SqrtScalar(x float64) float64 { return sqrtScalar(x) }
//go:inline
func AbsScalar(x float64) float64 { return absScalar(x) }
//go:inline
func NegScalar(x float64) float64 { return negScalar(x) }

func EmlSIMD(x, y, result []float64) {
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}
	if len(x) == 0 { return }
	emlSIMD(x, y, result)
}

func ExpSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	ExpSIMDTo(x, result)
	return result
}

func ExpSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	expSIMDTo(x, result)
}

func expSIMDTo(x, result []float64) {
	dispatchExpSIMDTo(x, result)
}

func LogSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	LogSIMDTo(x, result)
	return result
}

func LogSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	logSIMDTo(x, result)
}

func logSIMDTo(x, result []float64) {
	dispatchLogSIMDTo(x, result)
}

func SqrtSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	SqrtSIMDTo(x, result)
	return result
}

func SqrtSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	sqrtSIMDTo(x, result)
}

func sqrtSIMDTo(x, result []float64) {
	dispatchSqrtSIMDTo(x, result)
}

func SinSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	SinSIMDTo(x, result)
	return result
}

func SinSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	sinSIMDTo(x, result)
}

func sinSIMDTo(x, result []float64) {
	dispatchSinSIMDTo(x, result)
}

func CosSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	CosSIMDTo(x, result)
	return result
}

func CosSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	cosSIMDTo(x, result)
}

func cosSIMDTo(x, result []float64) {
	dispatchCosSIMDTo(x, result)
}

func TanSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	TanSIMDTo(x, result)
	return result
}

func TanSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	tanSIMDTo(x, result)
}

func tanSIMDTo(x, result []float64) {
	dispatchTanSIMDTo(x, result)
}

func SinCosSIMD(x []float64) (sin, cos []float64) {
	sin = make([]float64, len(x))
	cos = make([]float64, len(x))
	SinCosSIMDTo(x, sin, cos)
	return
}

func SinCosSIMDTo(x, sin, cos []float64) {
	if len(x) != len(sin) || len(x) != len(cos) { panic("slice length mismatch") }
	sincosSIMDTo(x, sin, cos)
}

func sincosSIMDTo(x, sin, cos []float64) {
	dispatchSinCosSIMDTo(x, sin, cos)
}

func AddSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchAddSIMD(a, b, result)
	return result
}

func SubSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchSubSIMD(a, b, result)
	return result
}

func MulSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchMulSIMD(a, b, result)
	return result
}

func DivSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	dispatchDivSIMD(a, b, result)
	return result
}

func AbsSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	AbsSIMDTo(x, result)
	return result
}

func AbsSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	absSIMD(x, result)
}

func absSIMD(x, result []float64) {
	dispatchAbsSIMD(x, result)
}

func NegSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	NegSIMDTo(x, result)
	return result
}

func NegSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	negSIMD(x, result)
}

func negSIMD(x, result []float64) {
	dispatchNegSIMD(x, result)
}

func InvSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	InvSIMDTo(x, result)
	return result
}

func InvSIMDTo(x, result []float64) {
	if len(x) != len(result) { panic("slice length mismatch") }
	invSIMD(x, result)
}

func invSIMD(x, result []float64) {
	dispatchInvSIMD(x, result)
}

func SinhBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Sinh)
	return res
}

func CoshBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Cosh)
	return res
}

func TanhBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Tanh)
	return res
}

func AsinhBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Asinh)
	return res
}

func AcoshBatch(x []float64) []float64 {
	res := make([]float64, len(x))
	parallelizeGeneric(x, res, Acosh)
	return res
}

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
	SmallCutoff = 256
	LargeCutoff = 4096
)

var ErrLengthMismatch = EMLError("slice length mismatch")

type EMLError string
func (e EMLError) Error() string { return string(e) }

func EmlBatch(x, y []float64, fn VectorFunc) error {
	if len(x) != len(y) { return ErrLengthMismatch }
	result := make([]float64, len(x))
	return fn(x, y, result)
}

type VectorFunc func(x, y, result []float64) error

func AddScalarSIMD(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	AddScalarSIMDTo(a, b, result)
	return result
}

func AddScalarSIMDTo(a []float64, b float64, result []float64) {
	if len(a) != len(result) { panic("slice length mismatch") }
	dispatchAddScalarSIMD(a, b, result)
}

func MulScalarSIMD(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	MulScalarSIMDTo(a, b, result)
	return result
}

func MulScalarSIMDTo(a []float64, b float64, result []float64) {
	if len(a) != len(result) { panic("slice length mismatch") }
	dispatchMulScalarSIMD(a, b, result)
}

const L1TileSize = 32768
