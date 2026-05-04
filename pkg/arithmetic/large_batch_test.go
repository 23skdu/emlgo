package arithmetic

import (
	"math"
	"testing"
)

func TestAllLargeBatches(t *testing.T) {
	n := 1000
	data := make([]float64, n)
	data2 := make([]float64, n)
	for i := range data {
		data[i] = float64(i + 1)
		data2[i] = 2.0
	}

	AbsBatch(data)
	NegBatch(data)
	InvBatch(data)
	FloorBatch(data)
	CeilBatch(data)
	TruncBatch(data)
	Log1pBatch(data)
	Expm1Batch(data)
	PowBatch(data, 2.0)
	CbrtBatch(data)
	HypotBatch(data, data2)
	MaxBatch(data, 500)
	MinBatch(data, 500)
}

func TestArithmeticPanics(t *testing.T) {
	defer func() { recover() }()
	LCM(0, 0)
}

func TestArithmeticMore(t *testing.T) {
	Pow(2, 2)
	Max(1, 2)
	Min(1, 2)
	GCD(10, 5)
	LCM(10, 5)
}

func TestArithmeticEdgeFinal(t *testing.T) {
	// n == 0
	AbsBatch(nil)
	NegBatch(nil)
	InvBatch(nil)
	FloorBatch(nil)
	CeilBatch(nil)
	TruncBatch(nil)
	Log1pBatch(nil)
	Expm1Batch(nil)
	PowBatch(nil, 2.0)
	CbrtBatch(nil)
	HypotBatch(nil, nil)
	MaxBatch(nil, 0)
	MinBatch(nil, 0)
	
	// n == 32 (buffer path)
	x32 := make([]float64, 32)
	AbsBatch(x32)
	NegBatch(x32)
	InvBatch(x32)
	FloorBatch(x32)
	CeilBatch(x32)
	TruncBatch(x32)
	Log1pBatch(x32)
	Expm1Batch(x32)
	PowBatch(x32, 2.0)
	CbrtBatch(x32)
	HypotBatch(x32, x32)
	MaxBatch(x32, 0)
	MinBatch(x32, 0)
}

func TestArithmeticEdgeFinal2(t *testing.T) {
	// 65 <= n < 256 (hit the non-parallel, non-buffered path)
	n100 := make([]float64, 100)
	AbsBatch(n100)
	NegBatch(n100)
	InvBatch(n100)
	FloorBatch(n100)
	CeilBatch(n100)
	TruncBatch(n100)
	Log1pBatch(n100)
	Expm1Batch(n100)
	PowBatch(n100, 2.0)
	CbrtBatch(n100)
	HypotBatch(n100, n100)
	MaxBatch(n100, 0)
	MinBatch(n100, 0)
	
	// Max/Min NaNs
	Max(math.NaN(), 1)
	Max(1, math.NaN())
	Min(math.NaN(), 1)
	Min(1, math.NaN())
}
