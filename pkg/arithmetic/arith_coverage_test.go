package arithmetic

import (
	"math"
	"testing"
)

func TestArithCoverageRemaining(t *testing.T) {
	if !math.IsNaN(Pow(math.NaN(), 1.0)) || !math.IsNaN(Pow(1.0, math.NaN())) {
		t.Errorf("Pow with NaN should return NaN")
	}
	if Pow(0, 5) != 0 {
		t.Errorf("Pow(0, 5) should be 0")
	}
	if Pow(0, 0) != 1 {
		t.Errorf("Pow(0, 0) should be 1")
	}
	if Pow(5, 0) != 1 {
		t.Errorf("Pow(5, 0) should be 1")
	}
	if !math.IsNaN(Pow(-2, 0.5)) {
		t.Errorf("Pow(-2, 0.5) should be NaN")
	}
	if !math.IsInf(Pow(0, -2), 1) {
		t.Errorf("Pow(0, -2) should be +Inf")
	}
	// x < 0 and isInteger(y)
	if math.Abs(Pow(-2, 2)-4.0) > 1e-5 {
		t.Errorf("Pow(-2, 2) should be ~4")
	}
	if math.Abs(Pow(-2, 3)-(-8.0)) > 1e-5 {
		t.Errorf("Pow(-2, 3) should be ~-8")
	}
	if math.Abs(Pow(-2, -2)-0.25) > 1e-5 {
		t.Errorf("Pow(-2, -2) should be ~0.25")
	}
	if math.Abs(Pow(-2, -3)-(-0.125)) > 1e-5 {
		t.Errorf("Pow(-2, -3) should be ~-0.125")
	}
	if !math.IsNaN(Pow(-2, float64(math.MinInt))) {
		t.Errorf("Pow(-2, MinInt) should be NaN")
	}

	// PowInt
	if PowInt(2.0, 0) != 1.0 {
		t.Errorf("PowInt(2, 0) should be 1")
	}
	if PowInt(2.0, 3) != 8.0 {
		t.Errorf("PowInt(2, 3) should be 8")
	}
	if PowInt(2.0, -3) != 0.125 {
		t.Errorf("PowInt(2, -3) should be 0.125")
	}
	_ = PowInt(2.0, math.MinInt)

	// GCD and IntAbs MinInt cases
	_ = GCD(math.MinInt64, 10)
	_ = GCD(10, math.MinInt64)
	_ = GCD(-10, -20)
	if IntAbs(math.MinInt) != math.MaxInt {
		t.Errorf("IntAbs(MinInt) should be MaxInt")
	}
	if IntAbs(-5) != 5 {
		t.Errorf("IntAbs(-5) should be 5")
	}
	if IntAbs(5) != 5 {
		t.Errorf("IntAbs(5) should be 5")
	}
}

func TestInt8AndUint8Operations(t *testing.T) {
	aInt8 := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bInt8 := []int8{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	// Dot products
	_ = DotProductInt8(aInt8, bInt8)
	_ = DotProductInt8([]int8{1}, []int8{2, 3}) // mismatched len
	_ = DotProductUint8([]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}, []uint8{2, 3, 4, 5, 6, 7, 8, 9, 10})
	_ = DotProductUint8([]uint8{1}, []uint8{2, 3})

	// L2Squared and Euclidean
	_ = L2SquaredInt8(aInt8, bInt8)
	_ = L2SquaredInt8([]int8{1}, []int8{2, 3})
	_ = EuclideanDistanceInt8(aInt8, bInt8)

	// CosineDistanceInt8
	if CosineDistanceInt8(nil, nil) != 1.0 {
		t.Errorf("CosineDistanceInt8(nil, nil) should be 1.0")
	}
	if CosineDistanceInt8([]int8{0, 0}, []int8{0, 0}) != 1.0 {
		t.Errorf("Zero norm should be 1.0")
	}
	_ = CosineDistanceInt8(aInt8, bInt8)
	_ = CosineDistanceInt8([]int8{1}, []int8{2, 3})
	_ = CosineDistanceInt8([]int8{1, 1}, []int8{-1, -1}) // opposite

	// CosineDistanceUint8
	if CosineDistanceUint8(nil, nil) != 1.0 {
		t.Errorf("CosineDistanceUint8(nil, nil) should be 1.0")
	}
	if CosineDistanceUint8([]uint8{0, 0}, []uint8{0, 0}) != 1.0 {
		t.Errorf("Zero norm should be 1.0")
	}
	aUint8 := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bUint8 := []uint8{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	_ = CosineDistanceUint8(aUint8, bUint8)
	_ = CosineDistanceUint8([]uint8{1}, []uint8{2, 3})

	// AddBatchInt8 saturation and length mismatch
	addRes := AddBatchInt8([]int8{100, -100, 10}, []int8{50, -50, 20, 99})
	if addRes[0] != math.MaxInt8 || addRes[1] != math.MinInt8 || addRes[2] != 30 {
		t.Errorf("AddBatchInt8 saturation failed: %v", addRes)
	}

	// SubBatchInt8 saturation and length mismatch
	subRes := SubBatchInt8([]int8{100, -100, 10}, []int8{-50, 50, 5, 99})
	if subRes[0] != math.MaxInt8 || subRes[1] != math.MinInt8 || subRes[2] != 5 {
		t.Errorf("SubBatchInt8 saturation failed: %v", subRes)
	}

	// MulBatchInt8 saturation and length mismatch
	mulRes := MulBatchInt8([]int8{20, -20, 2}, []int8{10, 10, 3, 99})
	if mulRes[0] != math.MaxInt8 || mulRes[1] != math.MinInt8 || mulRes[2] != 6 {
		t.Errorf("MulBatchInt8 saturation failed: %v", mulRes)
	}
}
