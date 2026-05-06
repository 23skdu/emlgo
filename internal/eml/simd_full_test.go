package eml

import (
	"math"
	"runtime"
	"testing"
)

func TestSIMDDetection(t *testing.T) {
	switch runtime.GOARCH {
	case "amd64":
		t.Logf("AMD64: AVX2=%v, AVX512=%v", hasAVX2, hasAVX512)
	case "arm64":
		t.Logf("ARM64: Neon=%v", hasNeon)
	default:
		t.Logf("Architecture: %s (no SIMD)", runtime.GOARCH)
	}
}

func TestSIMD(t *testing.T) {
	sizes := []int{1, 7, 8, 15, 16, 31, 32, 64, 127, 128, 255, 256}

	for _, n := range sizes {
		t.Run("size_"+string(rune('0'+n%10)), func(t *testing.T) {
			x := make([]float64, n)
			y := make([]float64, n)
			result := make([]float64, n)

			for i := 0; i < n; i++ {
				x[i] = float64(i) * 0.1
				y[i] = float64(i+1) * 0.1
			}

			SIMD(x, y, result)

			for i := 0; i < n; i++ {
				want := math.Exp(x[i]) - math.Log(y[i])
				if math.Abs(result[i]-want) > 1e-10 {
					t.Errorf("SIMD[%d] = %v, want %v", i, result[i], want)
				}
			}
		})
	}
}


func TestExpSIMD(t *testing.T) {
	x := []float64{0, 0.5, 1, 1.5, 2, math.E, 10}
	result := ExpSIMD(x)

	for i, v := range result {
		want := math.Exp(x[i])
		if math.Abs(v-want) > 1e-10 {
			t.Errorf("ExpSIMD[%d] = %v, want %v", i, v, want)
		}
	}
}

func TestLogSIMD(t *testing.T) {
	x := []float64{0.1, 0.5, 1, 2, math.E, 10}
	result := LogSIMD(x)

	for i, v := range result {
		want := math.Log(x[i])
		if math.Abs(v-want) > 1e-10 {
			t.Errorf("LogSIMD[%d] = %v, want %v", i, v, want)
		}
	}
}

func TestAddSIMD(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{10, 20, 30, 40, 50}
	result := AddSIMD(a, b)

	expected := []float64{11, 22, 33, 44, 55}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("AddSIMD[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestMulSIMD(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{2, 3, 4, 5, 6}
	result := MulSIMD(a, b)

	expected := []float64{2, 6, 12, 20, 30}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("MulSIMD[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSubSIMD(t *testing.T) {
	a := []float64{10, 20, 30, 40, 50}
	b := []float64{1, 2, 3, 4, 5}
	result := SubSIMD(a, b)

	expected := []float64{9, 18, 27, 36, 45}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("SubSIMD[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestDivSIMD(t *testing.T) {
	a := []float64{10, 20, 30, 40, 50}
	b := []float64{2, 4, 5, 8, 10}
	result := DivSIMD(a, b)

	expected := []float64{5, 5, 6, 5, 5}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("DivSIMD[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSIMDEmptySlice(t *testing.T) {
	x := make([]float64, 0)
	y := make([]float64, 0)
	result := make([]float64, 0)

	SIMD(x, y, result)

	if len(result) != 0 {
		t.Error("Expected empty result")
	}
}


func TestSIMDLargeSlice(t *testing.T) {
	n := 10000
	x := make([]float64, n)
	y := make([]float64, n)
	result := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = float64(i%100) * 0.1
		y[i] = float64(i%100+1) * 0.1
	}

	SIMD(x, y, result)

	for i := 0; i < n; i++ {
		want := math.Exp(x[i]) - math.Log(y[i])
		if math.Abs(result[i]-want) > 1e-10 {
			t.Errorf("SIMD[%d] = %v, want %v", i, result[i], want)
			break
		}
	}
}


func TestSIMDNaNHandling(t *testing.T) {
	x := []float64{math.NaN(), 1, 2}
	y := []float64{1, math.NaN(), 2}
	result := make([]float64, 3)

	SIMD(x, y, result)

	if !math.IsNaN(result[0]) {
		t.Error("Expected NaN for NaN input")
	}
	if !math.IsNaN(result[1]) {
		t.Error("Expected NaN for NaN y")
	}
}


func TestSIMDInfHandling(t *testing.T) {
	x := []float64{math.Inf(1), math.Inf(-1), 0}
	y := []float64{1, 1, 0}
	result := make([]float64, 3)

	SIMD(x, y, result)

	if !math.IsInf(result[0], 1) {
		t.Error("Expected +Inf for +Inf input")
	}
	if result[1] != 0 {
		t.Error("Expected 0 for -Inf input")
	}
	if !math.IsInf(result[2], 1) {
		t.Error("Expected +Inf for y=0 (log(0)=-Inf, exp(0)-(-Inf)=+Inf)")
	}
}


func TestBatch(t *testing.T) {
	x := []float64{1, 2, 3}
	y := []float64{2, 3, 4}

	err := Batch(x, y, func(x, y, result []float64) error {
		SIMD(x, y, result)
		return nil
	})

	if err != nil {
		t.Errorf("Batch error: %v", err)
	}
}


func TestBatchLengthMismatch(t *testing.T) {
	x := []float64{1, 2, 3}
	y := []float64{2, 3}

	err := Batch(x, y, func(x, y, result []float64) error {
		return nil
	})

	if err != ErrLengthMismatch {
		t.Error("Expected ErrLengthMismatch")
	}
}


func BenchmarkSIMD(b *testing.B) {
	sizes := []int{64, 256, 1024, 4096}

	for _, n := range sizes {
		b.Run("size_"+string(rune('0'+n%10)), func(b *testing.B) {
			x := make([]float64, n)
			y := make([]float64, n)
			result := make([]float64, n)

			for i := 0; i < n; i++ {
				x[i] = float64(i) * 0.1
				y[i] = float64(i+1) * 0.1
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SIMD(x, y, result)
			}
		})
	}
}


func BenchmarkExpSIMD(b *testing.B) {
	n := 4096
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i) * 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExpSIMD(x)
	}
}

func BenchmarkLogSIMD(b *testing.B) {
	n := 4096
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i+1) * 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LogSIMD(x)
	}
}

func BenchmarkAddSIMD(b *testing.B) {
	n := 4096
	a := make([]float64, n)
	c := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = float64(i)
		c[i] = float64(n - i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddSIMD(a, c)
	}
}

func BenchmarkMulSIMD(b *testing.B) {
	n := 4096
	a := make([]float64, n)
	c := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = float64(i + 1)
		c[i] = float64(i + 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MulSIMD(a, c)
	}
}