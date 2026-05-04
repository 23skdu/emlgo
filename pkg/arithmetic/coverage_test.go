package arithmetic

import (
	"testing"
)

func TestArithmeticCoverage(t *testing.T) {
	sizes := []int{0, 1, 63, 64, 65, 255, 256, 512}
	for _, n := range sizes {
		a := make([]float64, n)
		b := make([]float64, n)
		for i := range a {
			a[i] = float64(i)
			b[i] = float64(i + 1)
		}

		AbsBatch(a)
		NegBatch(a)
		InvBatch(a)
		FloorBatch(a)
		CeilBatch(a)
		TruncBatch(a)
		SqrtBatch(a)
		AddBatch(a, b)
		SubBatch(a, b)
		MulBatch(a, b)
		DivBatch(a, b)
		AddScalarBatch(a, 10.0)
		MulScalarBatch(a, 10.0)
		
		// Int and Uint ops
		IntAdd(1, 2)
		IntSub(1, 2)
		IntMul(1, 2)
		IntDiv(1, 2)
		IntMod(1, 2)
		IntAbs(-1)
		IntMax(1, 2)
		IntMin(1, 2)
		
		UintAdd(1, 2)
		UintSub(1, 2)
		UintMul(1, 2)
		UintDiv(1, 2)
		UintMod(1, 2)
		UintMax(1, 2)
		UintMin(1, 2)
	}
}

func TestArithEdgeCases(t *testing.T) {
	// Div by zero
	IntDiv(1, 0)
	IntMod(1, 0)
	UintDiv(1, 0)
	UintMod(1, 0)
	
	// Max/Min/Mod etc
	Max(1.0, 2.0)
	Min(1.0, 2.0)
	Mod(10.0, 3.0)
	Remainder(10.0, 3.0)
	LogBase2(8.0)
	LogBase10(100.0)
	PowInt(2.0, 3)
	Round(1.5)
}
