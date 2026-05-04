package hyper

import (
	"testing"
)

func TestLargeHyperBatchForCoverage(t *testing.T) {
	n := 1000
	data := make([]float64, n)
	for i := range data {
		data[i] = float64(i + 1) / 1000.0
	}

	t.Run("SinhBatch_Large", func(t *testing.T) {
		SinhBatch(data)
	})
	t.Run("CoshBatch_Large", func(t *testing.T) {
		CoshBatch(data)
	})
	t.Run("TanhBatch_Large", func(t *testing.T) {
		TanhBatch(data)
	})
	t.Run("AsinhBatch_Large", func(t *testing.T) {
		AsinhBatch(data)
	})
	t.Run("AcoshBatch_Large", func(t *testing.T) {
		AcoshBatch(data)
	})
	t.Run("AtanhBatch_Large", func(t *testing.T) {
		AtanhBatch(data)
	})
}