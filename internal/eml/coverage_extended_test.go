package eml

import (
	"math"
	"runtime"
	"testing"
)

func TestSIMDFunctionsExpanded(t *testing.T) {
	smallData := []float64{1, 2, 3}
	emptyData := []float64{}
	result := make([]float64, len(smallData))

	t.Run("ExpSIMD_empty", func(t *testing.T) {
		r := ExpSIMD(emptyData)
		_ = r
	})
	t.Run("ExpSIMDTo_small", func(t *testing.T) {
		ExpSIMDTo(smallData, result)
	})
	t.Run("LogSIMD_empty", func(t *testing.T) {
		r := LogSIMD(emptyData)
		_ = r
	})
	t.Run("LogSIMDTo_small", func(t *testing.T) {
		LogSIMDTo(smallData, result)
	})
	t.Run("SinSIMD_empty", func(t *testing.T) {
		r := SinSIMD(emptyData)
		if r != nil && len(r) != 0 {
			t.Errorf("SinSIMD empty = %v", r)
		}
	})
	t.Run("SinSIMDTo_small", func(t *testing.T) {
		SinSIMDTo(smallData, result)
	})
	t.Run("CosSIMD_empty", func(t *testing.T) {
		r := CosSIMD(emptyData)
		if r != nil && len(r) != 0 {
			t.Errorf("CosSIMD empty = %v", r)
		}
	})
	t.Run("CosSIMDTo_small", func(t *testing.T) {
		CosSIMDTo(smallData, result)
	})
	t.Run("SinCosSIMD_empty", func(t *testing.T) {
		sin, cos := SinCosSIMD(emptyData)
		if (sin != nil && len(sin) != 0) || (cos != nil && len(cos) != 0) {
			t.Errorf("SinCosSIMD empty = %v, %v", sin, cos)
		}
	})
	t.Run("SinCosSIMDTo_small", func(t *testing.T) {
		sin := make([]float64, len(smallData))
		cos := make([]float64, len(smallData))
		SinCosSIMDTo(smallData, sin, cos)
	})
	t.Run("TanSIMD_empty", func(t *testing.T) {
		r := TanSIMD(emptyData)
		if r != nil && len(r) != 0 {
			t.Errorf("TanSIMD empty = %v", r)
		}
	})
	t.Run("TanSIMDTo_small", func(t *testing.T) {
		TanSIMDTo(smallData, result)
	})
	t.Run("SqrtSIMD_empty", func(t *testing.T) {
		r := SqrtSIMD(emptyData)
		if r != nil && len(r) != 0 {
			t.Errorf("SqrtSIMD empty = %v", r)
		}
	})
	t.Run("SqrtSIMDTo_small", func(t *testing.T) {
		SqrtSIMDTo(smallData, result)
	})
}

func TestFusedBatchOperationsExpanded(t *testing.T) {
	dataA := []float64{1, 2, 3, 4}
	dataC := []float64{1, 1, 1, 1}
	empty := []float64{}

	t.Run("ExpMulBatch_empty", func(t *testing.T) {
		r := ExpMulBatch(empty, dataC)
		if len(r) != 0 {
			t.Errorf("ExpMulBatch empty len = %v", len(r))
		}
	})
	t.Run("ExpMulBatch_small", func(t *testing.T) {
		r := ExpMulBatch(dataA[:2], dataC[:2])
		if len(r) != 2 {
			t.Errorf("ExpMulBatch small len = %v", len(r))
		}
	})
	t.Run("ExpAddBatch_empty", func(t *testing.T) {
		r := ExpAddBatch(empty, empty)
		if len(r) != 0 {
			t.Errorf("ExpAddBatch empty len = %v", len(r))
		}
	})
	t.Run("ExpAddBatch_small", func(t *testing.T) {
		r := ExpAddBatch(dataA[:2], dataC[:2])
		if len(r) != 2 {
			t.Errorf("ExpAddBatch small len = %v", len(r))
		}
	})
	t.Run("LogDivBatch_empty", func(t *testing.T) {
		r := LogDivBatch(empty, empty)
		if len(r) != 0 {
			t.Errorf("LogDivBatch empty len = %v", len(r))
		}
	})
	t.Run("LogDivBatch_small", func(t *testing.T) {
		r := LogDivBatch(dataA[:2], dataC[:2])
		if len(r) != 2 {
			t.Errorf("LogDivBatch small len = %v", len(r))
		}
	})
	t.Run("LogSubBatch_empty", func(t *testing.T) {
		r := LogSubBatch(empty, empty)
		if len(r) != 0 {
			t.Errorf("LogSubBatch empty len = %v", len(r))
		}
	})
	t.Run("LogSubBatch_small", func(t *testing.T) {
		r := LogSubBatch(dataA[:2], dataC[:2])
		if len(r) != 2 {
			t.Errorf("LogSubBatch small len = %v", len(r))
		}
	})
}

func TestParallelization(t *testing.T) {
	smallData := []float64{1, 2}
	_ = smallData

	t.Run("GetParallelChunkSize_small", func(t *testing.T) {
		c := GetParallelChunkSize(10)
		if c <= 0 {
			t.Errorf("GetParallelChunkSize(10) = %v, want > 0", c)
		}
	})
	t.Run("GetParallelChunkSize_medium", func(t *testing.T) {
		c := GetParallelChunkSize(1000)
		if c <= 0 {
			t.Errorf("GetParallelChunkSize(1000) = %v, want > 0", c)
		}
	})
}

func TestSIMDFeatureDetection(t *testing.T) {
	t.Run("HasSSE4", func(t *testing.T) {
		_ = HasSSE4()
	})
	t.Run("HasAVX2", func(t *testing.T) {
		_ = HasAVX2()
	})
	t.Run("HasAVX512", func(t *testing.T) {
		_ = HasAVX512()
	})
	t.Run("HasNeon", func(t *testing.T) {
		_ = HasNeon()
	})
	t.Run("HasNeonDot", func(t *testing.T) {
		_ = HasNeonDot()
	})
}

func TestScalarFunctionsExpanded(t *testing.T) {
	t.Run("FmaScalar", func(t *testing.T) {
		if got := FmaScalar(2, 3, 1); got != 7 {
			t.Errorf("FmaScalar(2,3,1) = %v, want 7", got)
		}
	})
	t.Run("SqrtScalar_zero", func(t *testing.T) {
		if got := SqrtScalar(0); got != 0 {
			t.Errorf("SqrtScalar(0) = %v, want 0", got)
		}
	})
	t.Run("SqrtScalar_pos", func(t *testing.T) {
		if got := SqrtScalar(4); got != 2 {
			t.Errorf("SqrtScalar(4) = %v, want 2", got)
		}
	})
	t.Run("SqrtScalar_neg", func(t *testing.T) {
		if !math.IsNaN(SqrtScalar(-1)) {
			t.Error("SqrtScalar(-1) should be NaN")
		}
	})
}

func TestComplexMathFunctions(t *testing.T) {
	t.Run("Copysign_pos", func(t *testing.T) {
		if got := Copysign(5, 1); got != 5 {
			t.Errorf("Copysign(5, 1) = %v, want 5", got)
		}
	})
	t.Run("Copysign_neg", func(t *testing.T) {
		if got := Copysign(5, -1); got != -5 {
			t.Errorf("Copysign(5, -1) = %v, want -5", got)
		}
	})
}

func TestNegSIMDFunctions(t *testing.T) {
	smallData := []float64{1, -2, 3}
	emptyData := []float64{}

	t.Run("NegSIMD_empty", func(t *testing.T) {
		r := NegSIMD(emptyData)
		if r != nil && len(r) != 0 {
			t.Errorf("NegSIMD empty = %v", r)
		}
	})
	t.Run("NegSIMD_small", func(t *testing.T) {
		r := NegSIMD(smallData)
		if len(r) != len(smallData) {
			t.Errorf("NegSIMD len = %v, want %v", len(r), len(smallData))
		}
	})
}

func TestInvSIMDFunctions(t *testing.T) {
	smallData := []float64{1, 2, 0.5}
	emptyData := []float64{}

	t.Run("InvSIMD_empty", func(t *testing.T) {
		r := InvSIMD(emptyData)
		if r != nil && len(r) != 0 {
			t.Errorf("InvSIMD empty = %v", r)
		}
	})
	t.Run("InvSIMD_small", func(t *testing.T) {
		r := InvSIMD(smallData)
		if len(r) != len(smallData) {
			t.Errorf("InvSIMD len = %v, want %v", len(r), len(smallData))
		}
	})
}

func TestBatchParallelism(t *testing.T) {
	smallData := []float64{1, 2, 3}
	largeData := make([]float64, 10000)
	for i := range largeData {
		largeData[i] = float64(i) * 0.1
	}

	t.Run("ExpSIMD_large", func(t *testing.T) {
		r := ExpSIMD(largeData)
		if len(r) != len(largeData) {
			t.Errorf("ExpSIMD large len = %v", len(r))
		}
	})
	t.Run("LogSIMD_large", func(t *testing.T) {
		r := LogSIMD(largeData)
		if len(r) != len(largeData) {
			t.Errorf("LogSIMD large len = %v", len(r))
		}
	})
	t.Run("AddSIMD", func(t *testing.T) {
		r := AddSIMD(smallData, smallData)
		if len(r) != len(smallData) {
			t.Errorf("AddSIMD len = %v", len(r))
		}
	})
	t.Run("SubSIMD", func(t *testing.T) {
		r := SubSIMD(smallData, smallData)
		if len(r) != len(smallData) {
			t.Errorf("SubSIMD len = %v", len(r))
		}
	})
	t.Run("MulSIMD", func(t *testing.T) {
		r := MulSIMD(smallData, smallData)
		if len(r) != len(smallData) {
			t.Errorf("MulSIMD len = %v", len(r))
		}
	})
	t.Run("DivSIMD", func(t *testing.T) {
		r := DivSIMD(smallData, smallData)
		if len(r) != len(smallData) {
			t.Errorf("DivSIMD len = %v", len(r))
		}
	})
}

func TestAbsSIMDFunctions(t *testing.T) {
	smallData := []float64{1, -2, 3}
	emptyData := []float64{}

	t.Run("AbsSIMD_empty", func(t *testing.T) {
		r := AbsSIMD(emptyData)
		if r != nil && len(r) != 0 {
			t.Errorf("AbsSIMD empty = %v", r)
		}
	})
	t.Run("AbsSIMD_small", func(t *testing.T) {
		r := AbsSIMD(smallData)
		if len(r) != len(smallData) {
			t.Errorf("AbsSIMD len = %v, want %v", len(r), len(smallData))
		}
	})
}

func TestModf(t *testing.T) {
	intPart, fracPart := Modf(3.5)
	if intPart != 3 {
		t.Errorf("Modf(3.5) int = %v, want 3", intPart)
	}
	if fracPart != 0.5 {
		t.Errorf("Modf(3.5) frac = %v, want 0.5", fracPart)
	}

	intPart, fracPart = Modf(-3.5)
	_ = intPart
	_ = fracPart
}

func TestArchitectureInfo(t *testing.T) {
	arch := runtime.GOARCH
	if arch == "" {
		t.Error("GOARCH should not be empty")
	}
}