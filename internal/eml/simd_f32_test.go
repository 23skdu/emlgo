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

func TestSubSIMDF32(t *testing.T) {
	for _, n := range []int{0, 1, 7, 16, 31, 32, 35, 64, 73} {
		a := make([]float32, n)
		b := make([]float32, n)
		dst := make([]float32, n)
		for i := 0; i < n; i++ {
			a[i] = float32(i * 3)
			b[i] = float32(i + 1)
		}
		SubSIMDF32To(a, b, dst)
		for i := 0; i < n; i++ {
			want := a[i] - b[i]
			if dst[i] != want {
				t.Fatalf("SubSIMDF32To n=%d i=%d got %f want %f", n, i, dst[i], want)
			}
		}
		got := SubSIMDF32(a, b)
		for i := 0; i < n; i++ {
			if got[i] != dst[i] {
				t.Fatalf("SubSIMDF32 n=%d i=%d got %f want %f", n, i, got[i], dst[i])
			}
		}
	}
}

func TestDivSIMDF32(t *testing.T) {
	for _, n := range []int{0, 1, 7, 16, 31, 32, 35, 64, 73} {
		a := make([]float32, n)
		b := make([]float32, n)
		dst := make([]float32, n)
		for i := 0; i < n; i++ {
			a[i] = float32((i + 1) * 6)
			b[i] = float32(i + 2)
		}
		DivSIMDF32To(a, b, dst)
		for i := 0; i < n; i++ {
			want := a[i] / b[i]
			if math.Abs(float64(dst[i]-want)) > 1e-6 {
				t.Fatalf("DivSIMDF32To n=%d i=%d got %f want %f", n, i, dst[i], want)
			}
		}
		got := DivSIMDF32(a, b)
		for i := 0; i < n; i++ {
			if got[i] != dst[i] {
				t.Fatalf("DivSIMDF32 n=%d i=%d got %f want %f", n, i, got[i], dst[i])
			}
		}
	}
}

func TestAddMulSqrtSIMDF32To(t *testing.T) {
	for _, n := range []int{0, 1, 7, 16, 31, 32, 35, 64, 73} {
		a := make([]float32, n)
		b := make([]float32, n)
		dstAdd := make([]float32, n)
		dstMul := make([]float32, n)
		dstSqrt := make([]float32, n)
		for i := 0; i < n; i++ {
			a[i] = float32(i + 1)
			b[i] = float32(i + 2)
		}
		AddSIMDF32To(a, b, dstAdd)
		MulSIMDF32To(a, b, dstMul)
		SqrtSIMDF32To(a, dstSqrt)
		for i := 0; i < n; i++ {
			if dstAdd[i] != a[i]+b[i] {
				t.Fatalf("AddSIMDF32To n=%d got %f want %f", n, dstAdd[i], a[i]+b[i])
			}
			if dstMul[i] != a[i]*b[i] {
				t.Fatalf("MulSIMDF32To n=%d got %f want %f", n, dstMul[i], a[i]*b[i])
			}
			wantSqrt := float32(math.Sqrt(float64(a[i])))
			if math.Abs(float64(dstSqrt[i]-wantSqrt)) > 1e-6 {
				t.Fatalf("SqrtSIMDF32To n=%d got %f want %f", n, dstSqrt[i], wantSqrt)
			}
		}
	}
}

func TestDotInt8SIMD(t *testing.T) {
	for _, n := range []int{0, 1, 7, 63, 64, 65, 127, 128, 137} {
		a := make([]uint8, n)
		b := make([]int8, n)
		var expected int32
		for i := 0; i < n; i++ {
			a[i] = uint8((i * 3) % 256)
			b[i] = int8((i*5)%256 - 128)
			expected += int32(a[i]) * int32(b[i])
		}
		got := DotInt8SIMD(a, b)
		if got != expected {
			t.Fatalf("DotInt8SIMD n=%d got %d want %d", n, got, expected)
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
