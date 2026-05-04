package arithmetic

import (
	"testing"
)

func TestExtraBatchFunctions(t *testing.T) {
	testData := []float64{1, 2, 3, 4, 5}
	
	t.Run("FmaBatch", func(t *testing.T) {
		a := []float64{1, 2, 3}
		b := []float64{2, 2, 2}
		c := []float64{1, 1, 1}
		result := FmaBatch(a, b, c)
		if len(result) != 3 || result[0] != 3 || result[1] != 5 || result[2] != 7 {
			t.Errorf("FmaBatch failed: %v", result)
		}
	})

	t.Run("Log2Batch", func(t *testing.T) {
		result := Log2Batch(testData)
		if len(result) != len(testData) {
			t.Errorf("Log2Batch length mismatch")
		}
	})

	t.Run("Log10Batch", func(t *testing.T) {
		result := Log10Batch(testData)
		if len(result) != len(testData) {
			t.Errorf("Log10Batch length mismatch")
		}
	})

	t.Run("TanBatch", func(t *testing.T) {
		result := TanBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("TanBatch length mismatch")
		}
	})
}

func TestLargeBatchForCoverage(t *testing.T) {
	// SmallCutoff is usually around 256 or similar.
	// We use 1000 to be safe and trigger parallel paths.
	n := 1000
	data := make([]float64, n)
	data2 := make([]float64, n)
	for i := range data {
		data[i] = float64(i + 1)
		data2[i] = 2.0
	}

	t.Run("AbsBatch_Large", func(t *testing.T) {
		AbsBatch(data)
	})
	t.Run("NegBatch_Large", func(t *testing.T) {
		NegBatch(data)
	})
	t.Run("InvBatch_Large", func(t *testing.T) {
		InvBatch(data)
	})
	t.Run("FloorBatch_Large", func(t *testing.T) {
		FloorBatch(data)
	})
	t.Run("CeilBatch_Large", func(t *testing.T) {
		CeilBatch(data)
	})
	t.Run("TruncBatch_Large", func(t *testing.T) {
		TruncBatch(data)
	})
	t.Run("Log1pBatch_Large", func(t *testing.T) {
		Log1pBatch(data)
	})
	t.Run("Expm1Batch_Large", func(t *testing.T) {
		Expm1Batch(data)
	})
	t.Run("PowBatch_Large", func(t *testing.T) {
		PowBatch(data, 2)
	})
	t.Run("CbrtBatch_Large", func(t *testing.T) {
		CbrtBatch(data)
	})
	t.Run("HypotBatch_Large", func(t *testing.T) {
		HypotBatch(data, data2)
	})
	t.Run("MaxBatch_Large", func(t *testing.T) {
		MaxBatch(data, 500)
	})
	t.Run("MinBatch_Large", func(t *testing.T) {
		MinBatch(data, 500)
	})
}
