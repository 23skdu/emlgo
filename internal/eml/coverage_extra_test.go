package eml

import (
	"math"
	"testing"
)

func TestExpandedSIMD(t *testing.T) {
	n := 1000
	x := make([]float64, n)
	result := make([]float64, n)
	for i := range x {
		x[i] = float64(i + 1)
	}

	t.Run("Log2SIMDTo", func(t *testing.T) {
		Log2SIMDTo(x, result)
		// Basic check
		if result[0] != 0 {
			t.Errorf("Log2SIMDTo(1) = %v, want 0", result[0])
		}
	})

	t.Run("Log10SIMDTo", func(t *testing.T) {
		Log10SIMDTo(x, result)
		// Basic check
		if result[9] != 1 {
			t.Errorf("Log10SIMDTo(10) = %v, want 1", result[9])
		}
	})
}

func TestFmaSIMD(t *testing.T) {
	n := 1000
	a := make([]float64, n)
	b := make([]float64, n)
	c := make([]float64, n)
	result := make([]float64, n)
	for i := range a {
		a[i] = 2
		b[i] = 3
		c[i] = 1
	}

	t.Run("FmaSIMD", func(t *testing.T) {
		got := FmaSIMD(a, b, c)
		if len(got) != n || got[0] != 7 {
			t.Errorf("FmaSIMD failed: %v", got[0])
		}
	})

	t.Run("FmaSIMDTo", func(t *testing.T) {
		FmaSIMDTo(a, b, c, result)
		if result[0] != 7 {
			t.Errorf("FmaSIMDTo failed: %v", result[0])
		}
	})
}

func TestSIMDHelpers(t *testing.T) {
	// Triggering small paths
	x := []float64{1, 2, 3}
	res := make([]float64, 3)
	
	ExpSIMDTo(x, res)
	LogSIMDTo(x, res)
	SinSIMDTo(x, res)
	CosSIMDTo(x, res)
	TanSIMDTo(x, res)
	SqrtSIMDTo(x, res)
	
	s := make([]float64, 3)
	c := make([]float64, 3)
	SinCosSIMDTo(x, s, c)

	AbsSIMDTo(x, res)
	NegSIMDTo(x, res)
	InvSIMDTo(x, res)
}

func TestNativeMathExtra(t *testing.T) {
	t.Run("nativeMax", func(t *testing.T) {
		nativeMax(math.NaN(), 5)
		nativeMax(5, math.NaN())
		nativeMax(0, math.Copysign(0, -1))
		nativeMax(math.Copysign(0, -1), 0)
	})
	t.Run("nativeMin", func(t *testing.T) {
		nativeMin(math.NaN(), 5)
		nativeMin(5, math.NaN())
		nativeMin(0, math.Copysign(0, -1))
		nativeMin(math.Copysign(0, -1), 0)
	})
	t.Run("nativeInv", func(t *testing.T) {
		nativeInv(0)
		nativeInv(math.NaN())
	})
	t.Run("nativeNeg", func(t *testing.T) {
		nativeNeg(math.NaN())
	})
	t.Run("nativeAbs", func(t *testing.T) {
		nativeAbs(math.NaN())
	})
}

func TestGetters(t *testing.T) {
	HasSSE4()
	HasAVX2()
	HasAVX512()
	HasNeon()
	HasNeonDot()
	HasSVE()
	HasFMA()
	HasAVXVNNI()
	
	FmaScalar(1, 2, 3)
	SqrtScalar(4)
	AbsScalar(-5)
	NegScalar(5)
}

func TestEmlBatchError(t *testing.T) {
	err := EmlBatch(make([]float64, 5), make([]float64, 4), nil)
	if err == nil {
		t.Errorf("Expected error for length mismatch in EmlBatch")
	}
	// Also test Error() method of EMLError
	_ = err.Error()
}

func TestFusedEdgeCases(t *testing.T) {
	x := []float64{-1, 0, 1}
	res := make([]float64, 3)
	LogDivTo(x, x, res)
	LogSubTo(x, x, res)
	ExpMulTo(x, x, res)
	ExpAddTo(x, x, res)
	
	AbsBranchless(-5)
	MinBranchless(1, 2)
	MinBranchless(2, 1)
	MaxBranchless(1, 2)
	MaxBranchless(2, 1)
	SelectBranchless(true, 1, 2)
	SelectBranchless(false, 1, 2)
	SelectNaNBranchless(true, 1, 2)
	SelectNaNBranchless(false, 1, 2)
}

func TestSIMDMixed(t *testing.T) {
	n := 1000
	x := make([]float64, n)
	res := make([]float64, n)
	AbsSIMDTo(x, res)
	NegSIMDTo(x, res)
	InvSIMDTo(x, res)
	
	_ = AbsSIMD(x)
	_ = NegSIMD(x)
	_ = InvSIMD(x)
}

func TestExpandedLarge(t *testing.T) {
	n := 1000
	x := make([]float64, n)
	res := make([]float64, n)
	Log2SIMDTo(x, res)
	Log10SIMDTo(x, res)
}

func TestNativeMissing(t *testing.T) {
	nativeAtan(0)
	nativeAtan2(0, 0)
	nativeAsin(0)
	nativeAcos(0)
}

func TestFusedLargeEdge(t *testing.T) {
	n := 1000
	x := make([]float64, n)
	y := make([]float64, n)
	res := make([]float64, n)
	for i := range x {
		x[i] = -1.0
		y[i] = -1.0
	}
	LogDivTo(x, y, res)
	LogSubTo(x, y, res)
}

func TestChunkLarge(t *testing.T) {
	GetParallelChunkSize(100000)
}

func TestFusedSmallEdge(t *testing.T) {
	x := []float64{-1, 0, 1}
	res := make([]float64, 3)
	LogDivTo(x, x, res)
	LogSubTo(x, x, res)
}
