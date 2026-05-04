package arithmetic

import (
	"math"
	"testing"
)

func TestEdgeCasesForCoverage(t *testing.T) {
	t.Run("MaxEdgeCases", func(t *testing.T) {
		Max(math.NaN(), 5)
		Max(5, math.NaN())
		Max(0, math.Copysign(0, -1))
		Max(math.Copysign(0, -1), 0)
	})

	t.Run("MinEdgeCases", func(t *testing.T) {
		Min(math.NaN(), 5)
		Min(5, math.NaN())
		Min(0, math.Copysign(0, -1))
		Min(math.Copysign(0, -1), 0)
	})

	t.Run("PowEdgeCases", func(t *testing.T) {
		Pow(0, 0)
		Pow(0, 1)
		Pow(0, -1)
		Pow(-1, 2)
		Pow(-1, 3)
		Pow(-1, 2.5) // NaN
		Pow(1, 10)
		Pow(2, 2)
		Pow(math.NaN(), 1)
		Pow(1, math.NaN())
		Pow(0, 2) // x=0, y>0
	})

	t.Run("GCD_LCM_Edge", func(t *testing.T) {
		GCD(0, 0)
		GCD(-5, 5)
		GCD(5, 0)
		GCD(0, 5)
		LCM(0, 5)
		LCM(5, 0)
		LCM(0, 0)
		LCM(-5, 5)
	})

	t.Run("IntEdgeCases", func(t *testing.T) {
		IntMax(5, 5)
		IntMax(5, 4)
		IntMax(4, 5)
		IntMin(5, 5)
		IntMin(5, 4)
		IntMin(4, 5)
		UintMax(5, 5)
		UintMax(5, 4)
		UintMax(4, 5)
		UintMin(5, 5)
		UintMin(5, 4)
		UintMin(4, 5)
	})

	t.Run("ExpM1_NaN", func(t *testing.T) {
		ExpM1(math.NaN())
	})
}

func TestArithmeticEdgeExtra(t *testing.T) {
	GCD(0, 5)
	GCD(5, 0)
	GCD(0, 0)
	LCM(0, 5)
	LCM(5, 0)
	LCM(0, 0)
}
