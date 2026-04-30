package eml

import (
	"math"
	"testing"
)

func TestEmlSIMD(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4}
	y := []float64{1, 1, 1, 1, 1}
	result := make([]float64, len(x))

	EmlSIMD(x, y, result)

	for i, v := range result {
		want := Eml(x[i], y[i])
		if math.Abs(v-want) > 1e-10 {
			t.Errorf("EmlSIMD[%d] = %v, want %v", i, v, want)
		}
	}
}

func TestEmlSIMDLenMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should have panicked")
		}
	}()

	x := []float64{1, 2}
	y := []float64{1}
	result := make([]float64, 2)
	EmlSIMD(x, y, result)
}

func TestHasAVX2(t *testing.T) {
	_ = HasAVX2()
}

func TestHasAVX512(t *testing.T) {
	_ = HasAVX512()
}

func BenchmarkEmlSIMD(b *testing.B) {
	size := 1000
	x := make([]float64, size)
	y := make([]float64, size)
	result := make([]float64, size)
	for i := 0; i < size; i++ {
		x[i] = float64(i)
		y[i] = float64(i) + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EmlSIMD(x, y, result)
	}
}

func BenchmarkEmlScalar(b *testing.B) {
	size := 1000
	x := make([]float64, size)
	y := make([]float64, size)
	result := make([]float64, size)
	for i := 0; i < size; i++ {
		x[i] = float64(i)
		y[i] = float64(i) + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < size; j++ {
			result[j] = Eml(x[j], y[j])
		}
	}
}