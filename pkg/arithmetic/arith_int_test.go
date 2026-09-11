package arithmetic

import (
	"math"
	"math/rand"
	"testing"
)

func TestDotProductInt8(t *testing.T) {
	a := []int8{1, 2, 3, -4, 5, 6, -7, 8, 9, 10, -11, 12, 13, 14, -15, 16, 17, -18}
	b := []int8{2, -1, 4, 3, -2, 1, 0, -3, 2, -1, 4, 3, -2, 1, 0, -3, 2, 1}

	var expected int32
	for i := range a {
		expected += int32(a[i]) * int32(b[i])
	}

	got := DotProductInt8(a, b)
	if got != expected {
		t.Errorf("DotProductInt8() = %v, want %v", got, expected)
	}
}

func TestDotProductUint8(t *testing.T) {
	a := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
	b := []uint8{2, 1, 4, 3, 2, 1, 0, 3, 2, 1, 4, 3, 2, 1, 0, 3, 2, 1}

	var expected uint32
	for i := range a {
		expected += uint32(a[i]) * uint32(b[i])
	}

	got := DotProductUint8(a, b)
	if got != expected {
		t.Errorf("DotProductUint8() = %v, want %v", got, expected)
	}
}

func TestCosineDistanceInt8(t *testing.T) {
	a := []int8{10, 20, 30, 40}
	b := []int8{10, 20, 30, 40}

	dist := CosineDistanceInt8(a, b)
	if math.Abs(float64(dist)) > 1e-5 {
		t.Errorf("CosineDistanceInt8 identical vectors = %v, want 0", dist)
	}

	c := []int8{-10, -20, -30, -40}
	distOpposite := CosineDistanceInt8(a, c)
	if math.Abs(float64(distOpposite-2.0)) > 1e-5 {
		t.Errorf("CosineDistanceInt8 opposite vectors = %v, want 2.0", distOpposite)
	}
}

func TestEuclideanDistanceInt8(t *testing.T) {
	a := []int8{0, 3}
	b := []int8{4, 0}
	dist := EuclideanDistanceInt8(a, b)
	if math.Abs(float64(dist-5.0)) > 1e-5 {
		t.Errorf("EuclideanDistanceInt8() = %v, want 5.0", dist)
	}
}

func TestBatchInt8(t *testing.T) {
	a := []int8{100, -100, 50}
	b := []int8{50, -50, 10}

	add := AddBatchInt8(a, b)
	if add[0] != 127 || add[1] != -128 || add[2] != 60 {
		t.Errorf("AddBatchInt8() = %v, want [127 -128 60]", add)
	}

	sub := SubBatchInt8(a, b)
	if sub[0] != 50 || sub[1] != -50 || sub[2] != 40 {
		t.Errorf("SubBatchInt8() = %v, want [50 -50 40]", sub)
	}
}

func BenchmarkDotProductInt8(b *testing.B) {
	dims := 128
	a := make([]int8, dims)
	bVec := make([]int8, dims)
	rng := rand.New(rand.NewSource(42)) // #nosec G404 - deterministic benchmark test data
	for i := range a {
		a[i] = int8(rng.Intn(256) - 128)    // #nosec G115 - value in [-128, 127] fits in int8
		bVec[i] = int8(rng.Intn(256) - 128) // #nosec G115 - value in [-128, 127] fits in int8
	}

	for b.Loop() {
		_ = DotProductInt8(a, bVec)
	}
}

func BenchmarkCosineDistanceInt8(b *testing.B) {
	dims := 128
	a := make([]int8, dims)
	bVec := make([]int8, dims)
	rng := rand.New(rand.NewSource(42)) // #nosec G404 - deterministic benchmark test data
	for i := range a {
		a[i] = int8(rng.Intn(256) - 128)    // #nosec G115 - value in [-128, 127] fits in int8
		bVec[i] = int8(rng.Intn(256) - 128) // #nosec G115 - value in [-128, 127] fits in int8
	}

	for b.Loop() {
		_ = CosineDistanceInt8(a, bVec)
	}
}
