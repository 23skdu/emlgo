package eml

import (
	"math"
	"testing"
)

func f32Close(a, b float32, tol float64) bool {
	return math.Abs(float64(a)-float64(b)) < tol
}

func TestExpSIMDF32(t *testing.T) {
	input := []float32{0, 1, -1, 2, 0.5}
	got := ExpSIMDF32(input)
	for i, v := range input {
		want := float32(math.Exp(float64(v)))
		if !f32Close(got[i], want, 1e-5) {
			t.Errorf("ExpSIMDF32[%d](%v) = %v, want %v", i, v, got[i], want)
		}
	}
}

func TestLogSIMDF32(t *testing.T) {
	input := []float32{1, math.E, 10, 0.5}
	got := LogSIMDF32(input)
	for i, v := range input {
		want := float32(math.Log(float64(v)))
		if !f32Close(got[i], want, 1e-5) {
			t.Errorf("LogSIMDF32[%d](%v) = %v, want %v", i, v, got[i], want)
		}
	}
}

func TestSqrtSIMDF32(t *testing.T) {
	input := []float32{0, 1, 4, 9, 2}
	got := SqrtSIMDF32(input)
	for i, v := range input {
		want := float32(math.Sqrt(float64(v)))
		if !f32Close(got[i], want, 1e-6) {
			t.Errorf("SqrtSIMDF32[%d](%v) = %v, want %v", i, v, got[i], want)
		}
	}
}

func TestAddSIMDF32(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{4, 5, 6}
	got := AddSIMDF32(a, b)
	want := []float32{5, 7, 9}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("AddSIMDF32[%d]: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestAddSIMDF32LengthMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on length mismatch")
		}
	}()
	AddSIMDF32([]float32{1, 2}, []float32{1})
}

func TestMulSIMDF32(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{4, 5, 6}
	got := MulSIMDF32(a, b)
	want := []float32{4, 10, 18}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("MulSIMDF32[%d]: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestMulSIMDF32LengthMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on length mismatch")
		}
	}()
	MulSIMDF32([]float32{1, 2}, []float32{1})
}

func TestAddScalarSIMDF32(t *testing.T) {
	got := AddScalarSIMDF32([]float32{1, 2, 3}, 10)
	for i, want := range []float32{11, 12, 13} {
		if got[i] != want {
			t.Errorf("AddScalarSIMDF32[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestMulScalarSIMDF32(t *testing.T) {
	got := MulScalarSIMDF32([]float32{1, 2, 3}, 3)
	for i, want := range []float32{3, 6, 9} {
		if got[i] != want {
			t.Errorf("MulScalarSIMDF32[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func BenchmarkExpSIMDF32VsF64(b *testing.B) {
	const n = 1024
	f32in := make([]float32, n)
	f64in := make([]float64, n)
	for i := range f32in {
		f32in[i] = float32(i) * 0.001
		f64in[i] = float64(i) * 0.001
	}

	b.Run("ExpF32", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = ExpSIMDF32(f32in)
		}
	})

	b.Run("ExpF64", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = ExpSIMD(f64in)
		}
	})
}
