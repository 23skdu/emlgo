//go:build cuda

package gpu

import (
	"math"
	"testing"
)

// FuzzLaunchConfig tests that various combinations of n and blockSize
// produce valid launch configurations.
func FuzzLaunchConfig(f *testing.F) {
	seeds := []struct {
		n, blockSize int
	}{
		{0, 256},
		{1, 1},
		{1024, 256},
		{1000000, 512},
		{1, 1024},
		{0, 0},
		{-1, 256},
		{1024, -1},
	}

	for _, s := range seeds {
		f.Add(s.n, s.blockSize)
	}

	f.Fuzz(func(t *testing.T, n, blockSize int) {
		lc := DefaultLaunchConfig(n, blockSize)

		// Validate block size bounds
		if lc.BlockDimX <= 0 || lc.BlockDimX > 1024 {
			t.Errorf("invalid BlockDimX %d for n=%d blockSize=%d", lc.BlockDimX, n, blockSize)
		}

		// Validate grid dimensions for valid n
		if n > 0 && lc.GridDimX <= 0 {
			t.Errorf("GridDimX should be >0 for n=%d", n)
		}

		// Validate total threads >= n for normal ranges
		totalThreads := lc.GridDimX * lc.GridDimY * lc.GridDimZ *
			lc.BlockDimX * lc.BlockDimY * lc.BlockDimZ
		if n > 0 && n <= 1<<20 && totalThreads < n {
			t.Errorf("total threads %d < n %d", totalThreads, n)
		}
	})
}

// FuzzMatMulLaunchConfig tests matrix multiply configs with various dimensions.
func FuzzMatMulLaunchConfig(f *testing.F) {
	seeds := []struct {
		m, n, tile int
	}{
		{1, 1, 16},
		{32, 32, 16},
		{1024, 1024, 16},
		{4096, 4096, 32},
		{0, 0, 16},
	}

	for _, s := range seeds {
		f.Add(s.m, s.n, s.tile)
	}

	f.Fuzz(func(t *testing.T, m, n, tile int) {
		lc := MatMulLaunchConfig(m, n, tile)

		if lc.BlockDimX <= 0 || lc.BlockDimX > 1024 {
			t.Errorf("invalid BlockDimX %d for m=%d n=%d tile=%d", lc.BlockDimX, m, n, tile)
		}
		if lc.BlockDimY <= 0 || lc.BlockDimY > 1024 {
			t.Errorf("invalid BlockDimY %d", lc.BlockDimY)
		}
	})
}

// FuzzULPDiff tests that the ULP computation works for extreme values.
func FuzzULPDiff(f *testing.F) {
	seeds := []struct {
		a, b float64
	}{
		{0.0, 0.0},
		{1.0, 1.0},
		{math.NaN(), 1.0},
		{math.Inf(1), math.Inf(1)},
		{math.Inf(-1), math.Inf(1)},
	}

	for _, s := range seeds {
		f.Add(s.a, s.b)
	}

	f.Fuzz(func(t *testing.T, a, b float64) {
		_ = ulpDiff(a, b)
	})
}

// ulpDiff computes the ULP distance between two float64 values.
func ulpDiff(a, b float64) uint64 {
	if a == b {
		return 0
	}
	if math.IsNaN(a) || math.IsNaN(b) {
		return 0
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return 0
	}
	bits, targetBits := math.Float64bits(a), math.Float64bits(b)
	if bits > targetBits {
		return bits - targetBits
	}
	return targetBits - bits
}
