package eml

import (
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	cpuNum      = runtime.NumCPU()
	hasSSE4     bool
	hasAVX2     bool
	hasAVX512   bool
	hasNeon     bool
	hasNeonDot  bool
	hasSVE      bool
	hasFMA      bool
	hasAVXVNNI  bool
	hasWasmSIMD bool
)

func init() {
	detectSIMD()
	initWorkerPool()
}

func detectSIMD() {
	detectPlatformSIMD()
}

// GetParallelChunkSize returns the ideal chunk size for parallel processing of n elements.
func GetParallelChunkSize(n int) int {
	if n < SmallCutoff() {
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

// HasWasmSIMD reports whether WASM SIMD is available.
func HasWasmSIMD() bool { return hasWasmSIMD }

// FmaScalar returns a * b + c.
func FmaScalar(a, b, c float64) float64 {
	if hasFMA {
		return fmaScalar(a, b, c)
	}
	return a*b + c
}

// SqrtScalar returns the square root of x.
func SqrtScalar(x float64) float64 { return sqrtScalar(x) }

// AbsScalar returns the absolute value of x.
func AbsScalar(x float64) float64 { return absScalar(x) }

// NegScalar returns the negation of x.
func NegScalar(x float64) float64 { return negScalar(x) }

// SIMD computes Exp(x[i]) - Log(y[i]) for each element and stores it in result.
func SIMD(x, y, result []float64) {
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}
	n := len(x)
	if n == 0 {
		return
	}
	if n < SmallCutoff() || poolClosed.Load() {
		emlSIMD(x, y, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		emlSIMD(x[start:end], y[start:end], result[start:end])
	})
}



// ExpSIMD returns a new slice containing the exponential of each element in x.
func ExpSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	ExpSIMDTo(x, result)
	return result
}

// ExpSIMDTo computes the exponential of each element in x and stores it in result.
func ExpSIMDTo(x, result []float64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
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
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
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
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	sqrtSIMDTo(x, result)
}

func sqrtSIMDTo(x, result []float64) {
	n := len(x)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchSqrtSIMDTo(x, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchSqrtSIMDTo(x[start:end], result[start:end])
	})
}

// SinSIMD returns a new slice containing the sine of each element in x.
func SinSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	SinSIMDTo(x, result)
	return result
}

// SinSIMDTo computes the sine of each element in x and stores it in result.
func SinSIMDTo(x, result []float64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
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
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
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
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
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
	if len(x) != len(sin) || len(x) != len(cos) {
		panic("slice length mismatch")
	}
	sincosSIMDTo(x, sin, cos)
}

func sincosSIMDTo(x, sin, cos []float64) {
	dispatchSinCosSIMDTo(x, sin, cos)
}

// AddSIMD returns a new slice containing the sum of elements in a and b.
func AddSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	AddSIMDTo(a, b, result)
	return result
}

// AddSIMDTo computes the sum of elements in a and b and stores it in result.
func AddSIMDTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchAddSIMD(a, b, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchAddSIMD(a[start:end], b[start:end], result[start:end])
	})
}

// SubSIMD returns a new slice containing the difference of elements in a and b.
func SubSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	SubSIMDTo(a, b, result)
	return result
}

// SubSIMDTo computes the difference of elements in a and b and stores it in result.
func SubSIMDTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchSubSIMD(a, b, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchSubSIMD(a[start:end], b[start:end], result[start:end])
	})
}

// MulSIMD returns a new slice containing the product of elements in a and b.
func MulSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	MulSIMDTo(a, b, result)
	return result
}

// MulSIMDTo computes the product of elements in a and b and stores it in result.
func MulSIMDTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchMulSIMD(a, b, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchMulSIMD(a[start:end], b[start:end], result[start:end])
	})
}

// DivSIMD returns a new slice containing the quotient of elements in a and b.
func DivSIMD(a, b []float64) []float64 {
	result := make([]float64, len(a))
	DivSIMDTo(a, b, result)
	return result
}

// DivSIMDTo computes the quotient of elements in a and b and stores it in result.
func DivSIMDTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchDivSIMD(a, b, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchDivSIMD(a[start:end], b[start:end], result[start:end])
	})
}

// AbsSIMD returns a new slice containing the absolute value of each element in x.
func AbsSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	AbsSIMDTo(x, result)
	return result
}

// AbsSIMDTo computes the absolute value of each element in x and stores it in result.
func AbsSIMDTo(x, result []float64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	absSIMD(x, result)
}

func absSIMD(x, result []float64) {
	n := len(x)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchAbsSIMD(x, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchAbsSIMD(x[start:end], result[start:end])
	})
}

// NegSIMD returns a new slice containing the negation of each element in x.
func NegSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	NegSIMDTo(x, result)
	return result
}

// NegSIMDTo computes the negation of each element in x and stores it in result.
func NegSIMDTo(x, result []float64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	negSIMD(x, result)
}

func negSIMD(x, result []float64) {
	n := len(x)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchNegSIMD(x, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchNegSIMD(x[start:end], result[start:end])
	})
}

// InvSIMD returns a new slice containing the inverse of each element in x.
func InvSIMD(x []float64) []float64 {
	result := make([]float64, len(x))
	InvSIMDTo(x, result)
	return result
}

// InvSIMDTo computes the inverse of each element in x and stores it in result.
func InvSIMDTo(x, result []float64) {
	if len(x) != len(result) {
		panic("slice length mismatch")
	}
	invSIMD(x, result)
}

func invSIMD(x, result []float64) {
	n := len(x)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchInvSIMD(x, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchInvSIMD(x[start:end], result[start:end])
	})
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

type parallelJob struct {
	x        []float64
	a        []float64
	b        []float64
	result   []float64
	sinOut   []float64
	cosOut   []float64
	start    int
	end      int
	fn       func(float64) float64
	chunkFn  func(start, end int)
	fusedOp  int
	isSinCos bool
	wg       *sync.WaitGroup
}

const (
	fusedNone   = 0
	fusedExpMul = 1
	fusedExpAdd = 2
	fusedLogDiv = 3
	fusedLogSub = 4
)

var jobQueue chan parallelJob

var (
	poolMu     sync.Mutex
	stopOnce   sync.Once
	poolClosed atomic.Bool
	workersWg  sync.WaitGroup
)

func initWorkerPool() {
	numWorkers := max(1, cpuNum)
	q := make(chan parallelJob, numWorkers*8)
	jobQueue = q
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go workerPoolWorker(q)
	}
}

// StopWorkerPool gracefully shuts down the worker pool goroutines.
// After calling Stop, no further parallel operations should be submitted.
// Safe to call multiple times or concurrently — protected by sync.Once.
func StopWorkerPool() {
	stopOnce.Do(func() {
		poolClosed.Store(true)
		close(jobQueue)
		workersWg.Wait()
	})
}

func restartWorkerPool() {
	workersWg.Wait()
	stopOnce = sync.Once{}
	poolClosed.Store(false)
	initWorkerPool()
}

func workerPoolWorker(q <-chan parallelJob) {
	defer workersWg.Done()
	for job := range q {
		switch {
		case job.chunkFn != nil:
			job.chunkFn(job.start, job.end)
		case job.isSinCos:
			for j := job.start; j < job.end; j++ {
				job.sinOut[j], job.cosOut[j] = Sincos(job.x[j])
			}
		case job.fusedOp != fusedNone:
			applyFusedOp(job)
		default:
			for j := job.start; j < job.end; j++ {
				job.result[j] = job.fn(job.x[j])
			}
		}
		job.wg.Done()
	}
}

func applyFusedOp(job parallelJob) {
	a, b, result := job.a, job.b, job.result
	switch job.fusedOp {
	case fusedExpMul:
		for j := job.start; j < job.end; j++ {
			result[j] = nativeExp(a[j]) * b[j]
		}
	case fusedExpAdd:
		for j := job.start; j < job.end; j++ {
			result[j] = nativeExp(a[j]) + b[j]
		}
	case fusedLogDiv:
		for j := job.start; j < job.end; j++ {
			if a[j] > 0 && b[j] > 0 {
				result[j] = nativeLog(a[j]) / b[j]
			} else {
				result[j] = nan()
			}
		}
	case fusedLogSub:
		for j := job.start; j < job.end; j++ {
			if a[j] > 0 {
				result[j] = nativeLog(a[j]) - b[j]
			} else {
				result[j] = nan()
			}
		}
	}
}

func parallelizeGeneric(x, result []float64, fn func(float64) float64) {
	n := len(x)
	if n == 0 {
		return
	}

	if n < SmallCutoff() || poolClosed.Load() {
		for j := 0; j < n; j++ {
			result[j] = fn(x[j])
		}
		return
	}

	chunkSize := GetParallelChunkSize(n)
	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		jobQueue <- parallelJob{
			x:      x,
			result: result,
			start:  i,
			end:    end,
			fn:     fn,
			wg:     &wg,
		}
	}
	wg.Wait()
}

func parallelizeSinCos(x, sin, cos []float64) {
	n := len(x)
	if n == 0 {
		return
	}

	if n < SmallCutoff() || poolClosed.Load() {
		for j := 0; j < n; j++ {
			sin[j], cos[j] = Sincos(x[j])
		}
		return
	}

	chunkSize := GetParallelChunkSize(n)
	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		jobQueue <- parallelJob{
			x:        x,
			sinOut:   sin,
			cosOut:   cos,
			start:    i,
			end:      end,
			isSinCos: true,
			wg:       &wg,
		}
	}
	wg.Wait()
}

func parallelizeFused(a, b, result []float64, fusedOp int) {
	n := len(a)
	if n == 0 {
		return
	}

	if n < SmallCutoff() || poolClosed.Load() {
		applyFusedOp(parallelJob{a: a, b: b, result: result, fusedOp: fusedOp, start: 0, end: n})
		return
	}

	chunkSize := GetParallelChunkSize(n)
	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		jobQueue <- parallelJob{
			a:       a,
			b:       b,
			result:  result,
			fusedOp: fusedOp,
			start:   i,
			end:     end,
			wg:      &wg,
		}
	}
	wg.Wait()
}

// ForEachChunk distributes work over [0, n) across the shared worker pool in chunks.
// For small n (n < SmallCutoff()) or when the worker pool is stopped, fn is executed
// synchronously on the caller goroutine.
func ForEachChunk(n int, fn func(start, end int)) {
	if n <= 0 {
		return
	}

	if n < SmallCutoff() || poolClosed.Load() {
		fn(0, n)
		return
	}

	chunkSize := GetParallelChunkSize(n)
	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		jobQueue <- parallelJob{
			start:   i,
			end:     end,
			chunkFn: fn,
			wg:      &wg,
		}
	}
	wg.Wait()
}

// SmallWorkloadFactor is the minimum amount of work, in float64 elements, that
// each worker should receive before fanning out to the pool is worthwhile.
//
// Measured crossover on a 16-core AVX2 host (exp, float64): the pool loses to a
// plain serial loop by 4.4x at n=256, 2.8x at n=512 and 1.9x at n=1024, breaks
// even around n=4096, and only wins beyond that (0.30x at n=65536). The cost is
// goroutine hand-off, so the threshold has to be expressed per worker rather
// than as one fixed slice length.
const SmallWorkloadFactor = 512

// SmallCutoff is the threshold below which operations are performed
// sequentially rather than fanned out to the worker pool.
//
// This used to be a flat 256, which was roughly 16x too low: it made the pool
// spawn one goroutine per core to process as few as 256 elements, which is
// almost entirely scheduling overhead. It is now derived from the core count so
// that a many-core machine needs proportionally more elements before fanning
// out pays, and a single-core machine never fans out at all.
var smallCutoffVal atomic.Int64

func init() {
	smallCutoffVal.Store(int64(SmallWorkloadFactor * runtime.NumCPU()))
}

// SmallCutoff returns the threshold below which operations are performed
// sequentially rather than fanned out to the worker pool.
func SmallCutoff() int {
	return int(smallCutoffVal.Load())
}

// SetSmallCutoff reconfigures the parallel cutoff threshold in elements.
// If n <= 0, it resets to the default derived from runtime.NumCPU().
func SetSmallCutoff(n int) {
	if n <= 0 {
		n = SmallWorkloadFactor * runtime.NumCPU()
	}
	smallCutoffVal.Store(int64(n))
}

// Parallelism returns the number of worker goroutines currently configured.
func Parallelism() int {
	return cpuNum
}

// SetParallelism reconfigures the worker pool size. If workers <= 0, resets to runtime.NumCPU().
func SetParallelism(workers int) {
	poolMu.Lock()
	defer poolMu.Unlock()
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	StopWorkerPool()
	cpuNum = workers
	restartWorkerPool()
}

// LargeCutoff is the maximum chunk size for parallel operations.
const LargeCutoff = 4096

// ErrLengthMismatch is returned when slice lengths do not match.
var ErrLengthMismatch = Error("slice length mismatch")

// Error represents an EML package error.
type Error string

// Error returns the string representation of the error.
func (e Error) Error() string { return string(e) }

// Batch applies a VectorFunc to x and y.
func Batch(x, y []float64, fn VectorFunc) error {
	if len(x) != len(y) {
		return ErrLengthMismatch
	}
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
	if len(a) != len(result) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchAddScalarSIMD(a, b, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchAddScalarSIMD(a[start:end], b, result[start:end])
	})
}

// MulScalarSIMD returns a new slice containing each element in a multiplied by b.
func MulScalarSIMD(a []float64, b float64) []float64 {
	result := make([]float64, len(a))
	MulScalarSIMDTo(a, b, result)
	return result
}

// MulScalarSIMDTo multiplies each element in a by b and stores it in result.
func MulScalarSIMDTo(a []float64, b float64, result []float64) {
	if len(a) != len(result) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n < SmallCutoff() || poolClosed.Load() {
		dispatchMulScalarSIMD(a, b, result)
		return
	}
	ForEachChunk(n, func(start, end int) {
		dispatchMulScalarSIMD(a[start:end], b, result[start:end])
	})
}


// AddSatInt8SIMDTo computes saturating addition of int8 slices element-wise into result.
func AddSatInt8SIMDTo(a, b, result []int8) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	dispatchAddSatInt8SIMD(a, b, result)
}

// SubSatInt8SIMDTo computes saturating subtraction of int8 slices element-wise into result.
func SubSatInt8SIMDTo(a, b, result []int8) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	dispatchSubSatInt8SIMD(a, b, result)
}
