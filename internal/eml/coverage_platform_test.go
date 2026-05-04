package eml

import (
	"testing"
)

func TestExpMulBatchAll(t *testing.T) {
	for _, n := range []int{1, 7, 8, 15, 16, 31, 32} {
		a := make([]float64, n)
		b := make([]float64, n)
		for i := range a {
			a[i] = float64(i) * 0.1
			b[i] = float64(i) * 0.2
		}
		_ = ExpMulBatch(a, b)
		_ = ExpAddBatch(a, b)
		_ = LogDivBatch(a, b)
		_ = LogSubBatch(a, b)
	}
}

func TestGetParallelChunk(t *testing.T) {
	for _, n := range []int{0, 1, 2, 3, 4, 7, 8, 15, 16, 31, 32, 100, 1000} {
		_ = GetParallelChunkSize(n)
	}
}

func TestDetectCheck2(t *testing.T) {
	detectSIMD()
}
func TestFusedLarge(t *testing.T) {
	n := 1000
	a := make([]float64, n)
	b := make([]float64, n)
	ExpMulBatch(a, b)
	ExpAddBatch(a, b)
	LogDivBatch(a, b)
	LogSubBatch(a, b)
}
