package eml

import (
	"sync"
)

type fusedOps struct{}

func (f fusedOps) ExpMulBatch(a, b []float64, result []float64) {
	n := len(a)

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
	if n == 0 {
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
	if n == 0 {
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
	if n == 0 {
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

// ExpMulBatch returns a new slice containing Exp(a[i]) * b[i].
func ExpMulBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.ExpMulBatch(a, b, result)
	return result
}


// ExpAddBatch returns a new slice containing Exp(a[i]) + b[i].
func ExpAddBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.ExpAddBatch(a, b, result)
	return result
}


// LogDivBatch returns a new slice containing Log(a[i]) / b[i].
func LogDivBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.LogDivBatch(a, b, result)
	return result
}


// LogSubBatch returns a new slice containing Log(a[i]) - b[i].
func LogSubBatch(a, b []float64) []float64 {
	n := len(a)
	if n == 0 || len(a) != len(b) {
		return a
	}
	result := make([]float64, n)
	fused.LogSubBatch(a, b, result)
	return result
}


// ExpMulTo computes Exp(a[i]) * b[i] and stores the result in the provided slice.
func ExpMulTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.ExpMulBatch(a, b, result)
}


// ExpAddTo computes Exp(a[i]) + b[i] and stores the result in the provided slice.
func ExpAddTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.ExpAddBatch(a, b, result)
}


// LogDivTo computes Log(a[i]) / b[i] and stores the result in the provided slice.
func LogDivTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.LogDivBatch(a, b, result)
}


// LogSubTo computes Log(a[i]) - b[i] and stores the result in the provided slice.
func LogSubTo(a, b, result []float64) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("slice length mismatch")
	}
	fused.LogSubBatch(a, b, result)
}


// AbsBranchless returns the absolute value of x without using explicit branches (where possible).
func AbsBranchless(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}


// MinBranchless returns the minimum of a and b.
func MinBranchless(a, b float64) float64 {
	if a <= b {
		return a
	}
	return b
}


// MaxBranchless returns the maximum of a and b.
func MaxBranchless(a, b float64) float64 {
	if a >= b {
		return a
	}
	return b
}


// SelectBranchless returns ifTrue if cond is true, else ifFalse.
func SelectBranchless(cond bool, ifTrue, ifFalse float64) float64 {
	if cond {
		return ifTrue
	}
	return ifFalse
}


// SelectNaNBranchless returns alt if isNaN is true, else val.
func SelectNaNBranchless(isNaN bool, val, alt float64) float64 {
	if isNaN {
		return alt
	}
	return val
}