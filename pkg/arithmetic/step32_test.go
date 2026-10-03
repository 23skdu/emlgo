package arithmetic

import (
	"testing"
)

func TestStep32_BatchTo(t *testing.T) {
	n := 16
	x := make([]float64, n)
	y := make([]float64, n)
	res := make([]float64, n)
	for i := range x {
		x[i] = float64(i + 1)
		y[i] = float64(i + 2)
	}

	AddBatchTo(x, y, res)
	for i := range res {
		if res[i] != x[i]+y[i] {
			t.Errorf("AddBatchTo[%d] = %v, want %v", i, res[i], x[i]+y[i])
		}
	}

	SubBatchTo(x, y, res)
	for i := range res {
		if res[i] != x[i]-y[i] {
			t.Errorf("SubBatchTo[%d] = %v, want %v", i, res[i], x[i]-y[i])
		}
	}

	MulBatchTo(x, y, res)
	for i := range res {
		if res[i] != x[i]*y[i] {
			t.Errorf("MulBatchTo[%d] = %v, want %v", i, res[i], x[i]*y[i])
		}
	}

	DivBatchTo(x, y, res)
	for i := range res {
		if res[i] != x[i]/y[i] {
			t.Errorf("DivBatchTo[%d] = %v, want %v", i, res[i], x[i]/y[i])
		}
	}

	allocs := testing.AllocsPerRun(10, func() {
		AddBatchTo(x, y, res)
		SubBatchTo(x, y, res)
		MulBatchTo(x, y, res)
		DivBatchTo(x, y, res)
	})
	if allocs > 0 {
		t.Errorf("BatchTo allocated %v allocs/op, want 0", allocs)
	}
}
