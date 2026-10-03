package eml

import (
	"math"
	"math/cmplx"
	"sync"
)

// ComplexBatch applies Complex(x[i], y[i]) to each element pair and stores in result.
func ComplexBatch(x, y, result []complex128) {
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = cmplx.Exp(x[i]) - cmplx.Log(y[i])
		}
	})
}

// ComplexBatchC64 applies Complex64(x[i], y[i]) to each element pair and stores in result.
func ComplexBatchC64(x, y, result []complex64) {
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = Complex64(x[i], y[i])
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = cmplx.Exp(x[i])
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = complex64(cmplx.Exp(complex128(x[i])))
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = cmplx.Log(x[i])
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = complex64(cmplx.Log(complex128(x[i])))
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = cmplx.Sin(x[i])
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = complex64(cmplx.Sin(complex128(x[i])))
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = cmplx.Cos(x[i])
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = complex64(cmplx.Cos(complex128(x[i])))
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = cmplx.Tan(x[i])
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = complex64(cmplx.Tan(complex128(x[i])))
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = cmplx.Sqrt(x[i])
		}
	})
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
	ForEachChunk(len(x), func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = complex64(cmplx.Sqrt(complex128(x[i])))
		}
	})
}

type neumaierAccumulator struct {
	sum float64
	c   float64
}

func (acc *neumaierAccumulator) add(val float64) {
	t := acc.sum + val
	if math.Abs(acc.sum) >= math.Abs(val) {
		acc.c += (acc.sum - t) + val
	} else {
		acc.c += (val - t) + acc.sum
	}
	acc.sum = t
}

func (acc *neumaierAccumulator) total() float64 {
	return acc.sum + acc.c
}

func complexDotChunk(a, b []complex128, start, end int) (float64, float64) {
	var s0Re, s1Re, s2Re, s3Re neumaierAccumulator
	var s0Im, s1Im, s2Im, s3Im neumaierAccumulator

	i := start
	for ; i+3 < end; i += 4 {
		a0, b0 := a[i], b[i]
		u0, v0 := real(a0), imag(a0)
		x0, y0 := real(b0), imag(b0)
		s0Re.add(u0*x0 + v0*y0)
		s0Im.add(v0*x0 - u0*y0)

		a1, b1 := a[i+1], b[i+1]
		u1, v1 := real(a1), imag(a1)
		x1, y1 := real(b1), imag(b1)
		s1Re.add(u1*x1 + v1*y1)
		s1Im.add(v1*x1 - u1*y1)

		a2, b2 := a[i+2], b[i+2]
		u2, v2 := real(a2), imag(a2)
		x2, y2 := real(b2), imag(b2)
		s2Re.add(u2*x2 + v2*y2)
		s2Im.add(v2*x2 - u2*y2)

		a3, b3 := a[i+3], b[i+3]
		u3, v3 := real(a3), imag(a3)
		x3, y3 := real(b3), imag(b3)
		s3Re.add(u3*x3 + v3*y3)
		s3Im.add(v3*x3 - u3*y3)
	}

	var totRe, totIm neumaierAccumulator
	totRe.add(s0Re.total())
	totRe.add(s1Re.total())
	totRe.add(s2Re.total())
	totRe.add(s3Re.total())

	totIm.add(s0Im.total())
	totIm.add(s1Im.total())
	totIm.add(s2Im.total())
	totIm.add(s3Im.total())

	for ; i < end; i++ {
		ai, bi := a[i], b[i]
		ui, vi := real(ai), imag(ai)
		xi, yi := real(bi), imag(bi)
		totRe.add(ui*xi + vi*yi)
		totIm.add(vi*xi - ui*yi)
	}

	return totRe.total(), totIm.total()
}

// ComplexDotProduct returns the Hermitian inner product of two complex128 slices: sum(a[i] * conj(b[i]))
// using 4-way unrolled Neumaier compensated summation fanned out across the shared worker pool.
func ComplexDotProduct(a, b []complex128) complex128 {
	if len(a) != len(b) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n == 0 {
		return 0
	}
	if n < SmallCutoff || poolClosed.Load() {
		re, im := complexDotChunk(a, b, 0, n)
		return complex(re, im)
	}

	chunkSize := GetParallelChunkSize(n)
	numChunks := (n + chunkSize - 1) / chunkSize
	partialRe := make([]float64, numChunks)
	partialIm := make([]float64, numChunks)

	var wg sync.WaitGroup
	chunkIdx := 0
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		idx := chunkIdx
		chunkIdx++
		wg.Add(1)
		jobQueue <- parallelJob{
			start: i,
			end:   end,
			chunkFn: func(s, e int) {
				partialRe[idx], partialIm[idx] = complexDotChunk(a, b, s, e)
			},
			wg: &wg,
		}
	}
	wg.Wait()

	var totRe, totIm neumaierAccumulator
	for i := 0; i < numChunks; i++ {
		totRe.add(partialRe[i])
		totIm.add(partialIm[i])
	}
	return complex(totRe.total(), totIm.total())
}

func complexDotChunkC64(a, b []complex64, start, end int) (float64, float64) {
	var s0Re, s1Re, s2Re, s3Re neumaierAccumulator
	var s0Im, s1Im, s2Im, s3Im neumaierAccumulator

	i := start
	for ; i+3 < end; i += 4 {
		a0, b0 := a[i], b[i]
		u0, v0 := float64(real(a0)), float64(imag(a0))
		x0, y0 := float64(real(b0)), float64(imag(b0))
		s0Re.add(u0*x0 + v0*y0)
		s0Im.add(v0*x0 - u0*y0)

		a1, b1 := a[i+1], b[i+1]
		u1, v1 := float64(real(a1)), float64(imag(a1))
		x1, y1 := float64(real(b1)), float64(imag(b1))
		s1Re.add(u1*x1 + v1*y1)
		s1Im.add(v1*x1 - u1*y1)

		a2, b2 := a[i+2], b[i+2]
		u2, v2 := float64(real(a2)), float64(imag(a2))
		x2, y2 := float64(real(b2)), float64(imag(b2))
		s2Re.add(u2*x2 + v2*y2)
		s2Im.add(v2*x2 - u2*y2)

		a3, b3 := a[i+3], b[i+3]
		u3, v3 := float64(real(a3)), float64(imag(a3))
		x3, y3 := float64(real(b3)), float64(imag(b3))
		s3Re.add(u3*x3 + v3*y3)
		s3Im.add(v3*x3 - u3*y3)
	}

	var totRe, totIm neumaierAccumulator
	totRe.add(s0Re.total())
	totRe.add(s1Re.total())
	totRe.add(s2Re.total())
	totRe.add(s3Re.total())

	totIm.add(s0Im.total())
	totIm.add(s1Im.total())
	totIm.add(s2Im.total())
	totIm.add(s3Im.total())

	for ; i < end; i++ {
		ai, bi := a[i], b[i]
		ui, vi := float64(real(ai)), float64(imag(ai))
		xi, yi := float64(real(bi)), float64(imag(bi))
		totRe.add(ui*xi + vi*yi)
		totIm.add(vi*xi - ui*yi)
	}

	return totRe.total(), totIm.total()
}

// ComplexDotProductC64 returns the Hermitian inner product of two complex64 slices
// using 4-way unrolled Neumaier compensated summation in float64 precision.
func ComplexDotProductC64(a, b []complex64) complex64 {
	if len(a) != len(b) {
		panic("slice length mismatch")
	}
	n := len(a)
	if n == 0 {
		return 0
	}
	if n < SmallCutoff || poolClosed.Load() {
		re, im := complexDotChunkC64(a, b, 0, n)
		return complex64(complex(re, im))
	}

	chunkSize := GetParallelChunkSize(n)
	numChunks := (n + chunkSize - 1) / chunkSize
	partialRe := make([]float64, numChunks)
	partialIm := make([]float64, numChunks)

	var wg sync.WaitGroup
	chunkIdx := 0
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		idx := chunkIdx
		chunkIdx++
		wg.Add(1)
		jobQueue <- parallelJob{
			start: i,
			end:   end,
			chunkFn: func(s, e int) {
				partialRe[idx], partialIm[idx] = complexDotChunkC64(a, b, s, e)
			},
			wg: &wg,
		}
	}
	wg.Wait()

	var totRe, totIm neumaierAccumulator
	for i := 0; i < numChunks; i++ {
		totRe.add(partialRe[i])
		totIm.add(partialIm[i])
	}
	return complex64(complex(totRe.total(), totIm.total()))
}
