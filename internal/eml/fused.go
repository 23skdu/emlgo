package eml

import (
	"sync"
)

type fusedOps struct{}

func (f fusedOps) ExpMulBatch(a, b []float64, result []float64) {
	n := len(a)
	if n == 0 || len(a) != len(b) || len(a) != len(result) {
		return
	}

	if n < SmallCutoff {
		for i := 0; i < n; i++ {
			result[i] = nativeExp(a[i]) * b[i]
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
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = nativeExp(a[j]) * b[j]
			}
		}(i, end)
	}
	wg.Wait()
}

func (f fusedOps) ExpAddBatch(a, b []float64, result []float64) {
	n := len(a)
	if n == 0 || len(a) != len(b) || len(a) != len(result) {
		return
	}

	if n < SmallCutoff {
		for i := 0; i < n; i++ {
			result[i] = nativeExp(a[i]) + b[i]
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
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = nativeExp(a[j]) + b[j]
			}
		}(i, end)
	}
	wg.Wait()
}

func (f fusedOps) LogDivBatch(a, b []float64, result []float64) {
	n := len(a)
	if n == 0 || len(a) != len(b) || len(a) != len(result) {
		return
	}

	if n < SmallCutoff {
		for i := 0; i < n; i++ {
			if a[i] > 0 && b[i] > 0 {
				result[i] = nativeLog(a[i]) / b[i]
			} else {
				result[i] = nan()
			}
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
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if a[j] > 0 && b[j] > 0 {
					result[j] = nativeLog(a[j]) / b[j]
				} else {
					result[j] = nan()
				}
			}
		}(i, end)
	}
	wg.Wait()
}

func (f fusedOps) LogSubBatch(a, b []float64, result []float64) {
	n := len(a)
	if n == 0 || len(a) != len(b) || len(a) != len(result) {
		return
	}

	if n < SmallCutoff {
		for i := 0; i < n; i++ {
			if a[i] > 0 {
				result[i] = nativeLog(a[i]) - b[i]
			} else {
				result[i] = nan()
			}
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
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if a[j] > 0 {
					result[j] = nativeLog(a[j]) - b[j]
				} else {
					result[j] = nan()
				}
			}
		}(i, end)
	}
	wg.Wait()
}

var fused fusedOps

func ExpMulBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.ExpMulBatch(a, b, result)
	return result
}

func ExpAddBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.ExpAddBatch(a, b, result)
	return result
}

func LogDivBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.LogDivBatch(a, b, result)
	return result
}

func LogSubBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.LogSubBatch(a, b, result)
	return result
}

func ExpMulTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.ExpMulBatch(a, b, result)
}

func ExpAddTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.ExpAddBatch(a, b, result)
}

func LogDivTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.LogDivBatch(a, b, result)
}

func LogSubTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.LogSubBatch(a, b, result)
}

func AbsBranchless(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func MinBranchless(a, b float64) float64 {
	if a <= b {
		return a
	}
	return b
}

func MaxBranchless(a, b float64) float64 {
	if a >= b {
		return a
	}
	return b
}

func SelectBranchless(cond bool, ifTrue, ifFalse float64) float64 {
	if cond {
		return ifTrue
	}
	return ifFalse
}

func SelectNaNBranchless(isNaN bool, val, alt float64) float64 {
	if isNaN {
		return alt
	}
	return val
}