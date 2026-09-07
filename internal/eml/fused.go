package eml

import "math"

type fusedOps struct{}

func (f fusedOps) ExpMulBatch(a, b []float64, result []float64) {
	parallelizeFused(a, b, result, fusedExpMul)
}

func (f fusedOps) ExpAddBatch(a, b []float64, result []float64) {
	parallelizeFused(a, b, result, fusedExpAdd)
}

func (f fusedOps) LogDivBatch(a, b []float64, result []float64) {
	parallelizeFused(a, b, result, fusedLogDiv)
}

func (f fusedOps) LogSubBatch(a, b []float64, result []float64) {
	parallelizeFused(a, b, result, fusedLogSub)
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

// AbsBranchless returns the absolute value of x using bitwise operations.
// This is a true branchless implementation — the compiler emits no conditional jumps.
func AbsBranchless(x float64) float64 {
	bits := math.Float64bits(x)
	sign := bits >> 63
	return math.Float64frombits(bits ^ (sign<<63 | sign))
}

// MinBranchless returns the minimum of a and b using bitwise selection.
func MinBranchless(a, b float64) float64 {
	diff := a - b
	mask := -int64(math.Float64bits(diff) >> 63)
	return math.Float64frombits(
		(math.Float64bits(a) & uint64(mask)) | (math.Float64bits(b) & ^uint64(mask)),
	)
}

// MaxBranchless returns the maximum of a and b using bitwise selection.
func MaxBranchless(a, b float64) float64 {
	diff := a - b
	mask := -int64(math.Float64bits(diff) >> 63)
	return math.Float64frombits(
		(math.Float64bits(a) & ^uint64(mask)) | (math.Float64bits(b) & uint64(mask)),
	)
}

// SelectBranchless returns ifTrue if cond is true, else ifFalse, without branching.
func SelectBranchless(cond bool, ifTrue, ifFalse float64) float64 {
	mask := -int64(boolToInt(cond))
	return math.Float64frombits(
		(math.Float64bits(ifTrue) & uint64(mask)) | (math.Float64bits(ifFalse) & ^uint64(mask)),
	)
}

// SelectNaNBranchless returns alt if isNaN is true, else val, without branching.
func SelectNaNBranchless(isNaN bool, val, alt float64) float64 {
	mask := -int64(boolToInt(isNaN))
	return math.Float64frombits(
		(math.Float64bits(val) & ^uint64(mask)) | (math.Float64bits(alt) & uint64(mask)),
	)
}

func boolToInt(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}
