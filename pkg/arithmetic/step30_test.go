package arithmetic

import (
	"math"
	"testing"
)

// TestStep30_AddBatchInt8To_ZeroAlloc asserts that AddBatchInt8To and SubBatchInt8To
// write directly into destination without allocating any heap memory.
func TestStep30_AddBatchInt8To_ZeroAlloc(t *testing.T) {
	const n = 1024
	a := make([]int8, n)
	b := make([]int8, n)
	dst := make([]int8, n)

	for i := range a {
		a[i] = int8(i % 100)
		b[i] = int8((i * 3) % 100)
	}

	allocs := testing.AllocsPerRun(100, func() {
		AddBatchInt8To(a, b, dst)
		SubBatchInt8To(a, b, dst)
		MulBatchInt8To(a, b, dst)
	})

	if allocs != 0 {
		t.Fatalf("expected 0 heap allocs, got %v", allocs)
	}
}

// TestStep30_SaturatingArithmeticValues tests boundary saturation conditions
// across varying slice lengths (1 to 257) ensuring AVX2 32-element blocks
// and scalar remainder elements are bit-exact.
func TestStep30_SaturatingArithmeticValues(t *testing.T) {
	lengths := []int{1, 7, 16, 31, 32, 33, 64, 127, 128, 257, 1024}

	for _, n := range lengths {
		a := make([]int8, n)
		b := make([]int8, n)
		wantAdd := make([]int8, n)
		wantSub := make([]int8, n)
		wantMul := make([]int8, n)

		for i := 0; i < n; i++ {
			va := int8((i*17)%200 - 100) // #nosec G115 -- test value within int8 range
			vb := int8((i*31)%200 - 100) // #nosec G115 -- test value within int8 range
			a[i] = va
			b[i] = vb

			// Expected saturating add
			vAdd := int32(va) + int32(vb)
			if vAdd > math.MaxInt8 {
				vAdd = math.MaxInt8
			} else if vAdd < math.MinInt8 {
				vAdd = math.MinInt8
			}
			wantAdd[i] = int8(vAdd) // #nosec G115 -- clamped to [math.MinInt8, math.MaxInt8]

			// Expected saturating sub
			vSub := int32(va) - int32(vb)
			if vSub > math.MaxInt8 {
				vSub = math.MaxInt8
			} else if vSub < math.MinInt8 {
				vSub = math.MinInt8
			}
			wantSub[i] = int8(vSub) // #nosec G115 -- clamped to [math.MinInt8, math.MaxInt8]

			// Expected saturating mul
			vMul := int32(va) * int32(vb)
			if vMul > math.MaxInt8 {
				vMul = math.MaxInt8
			} else if vMul < math.MinInt8 {
				vMul = math.MinInt8
			}
			wantMul[i] = int8(vMul) // #nosec G115 -- clamped to [math.MinInt8, math.MaxInt8]
		}

		gotAdd := AddBatchInt8(a, b)
		gotSub := SubBatchInt8(a, b)
		gotMul := MulBatchInt8(a, b)

		for i := 0; i < n; i++ {
			if gotAdd[i] != wantAdd[i] {
				t.Fatalf("n=%d Add[%d]: got %d, want %d (va=%d, vb=%d)", n, i, gotAdd[i], wantAdd[i], a[i], b[i])
			}
			if gotSub[i] != wantSub[i] {
				t.Fatalf("n=%d Sub[%d]: got %d, want %d (va=%d, vb=%d)", n, i, gotSub[i], wantSub[i], a[i], b[i])
			}
			if gotMul[i] != wantMul[i] {
				t.Fatalf("n=%d Mul[%d]: got %d, want %d (va=%d, vb=%d)", n, i, gotMul[i], wantMul[i], a[i], b[i])
			}
		}
	}
}

// TestStep30_LengthMismatchPanics asserts that passing mismatched slices
// panics across all integer arithmetic functions instead of silent truncation.
func TestStep30_LengthMismatchPanics(t *testing.T) {
	short := []int8{1, 2}
	long := []int8{1, 2, 3}
	dst := make([]int8, 2)

	assertPanic := func(name string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s: expected panic on length mismatch", name)
			}
		}()
		fn()
	}

	assertPanic("AddBatchInt8", func() { _ = AddBatchInt8(short, long) })
	assertPanic("AddBatchInt8To", func() { AddBatchInt8To(short, long, dst) })
	assertPanic("SubBatchInt8", func() { _ = SubBatchInt8(short, long) })
	assertPanic("SubBatchInt8To", func() { SubBatchInt8To(short, long, dst) })
	assertPanic("MulBatchInt8", func() { _ = MulBatchInt8(short, long) })
	assertPanic("MulBatchInt8To", func() { MulBatchInt8To(short, long, dst) })
	assertPanic("DotProductInt8", func() { _ = DotProductInt8(short, long) })
	assertPanic("CosineDistanceInt8", func() { _ = CosineDistanceInt8(short, long) })
	assertPanic("L2SquaredInt8", func() { _ = L2SquaredInt8(short, long) })
}

// BenchmarkStep30_AddBatchInt8_Large benchmarks AddBatchInt8 and AddBatchInt8To at n = 65,536.
func BenchmarkStep30_AddBatchInt8_Large(b *testing.B) {
	const n = 65536
	a := make([]int8, n)
	bSlice := make([]int8, n)
	dst := make([]int8, n)

	for i := range a {
		a[i] = int8(i % 120)
		bSlice[i] = int8((i * 2) % 120)
	}

	b.Run("Allocating", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = AddBatchInt8(a, bSlice)
		}
	})

	b.Run("To_InPlace", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			AddBatchInt8To(a, bSlice, dst)
		}
	})
}

// FuzzAddSubSatInt8 fuzzes random byte pairs through AddBatchInt8 and SubBatchInt8
// verifying exact saturation behavior against reference arithmetic.
func FuzzAddSubSatInt8(f *testing.F) {
	f.Add(int8(100), int8(50))
	f.Add(int8(-56), int8(-106))
	f.Add(int8(127), int8(1))
	f.Add(int8(-128), int8(-1))
	f.Add(int8(0), int8(0))

	f.Fuzz(func(t *testing.T, a, b int8) {
		va := []int8{a}
		vb := []int8{b}

		gotAdd := AddBatchInt8(va, vb)[0]
		gotSub := SubBatchInt8(va, vb)[0]

		// Reference add
		rawAdd := int32(a) + int32(b)
		if rawAdd > math.MaxInt8 {
			rawAdd = math.MaxInt8
		} else if rawAdd < math.MinInt8 {
			rawAdd = math.MinInt8
		}
		if gotAdd != int8(rawAdd) {
			t.Errorf("AddSat(%d, %d) = %d, want %d", a, b, gotAdd, rawAdd)
		}

		// Reference sub
		rawSub := int32(a) - int32(b)
		if rawSub > math.MaxInt8 {
			rawSub = math.MaxInt8
		} else if rawSub < math.MinInt8 {
			rawSub = math.MinInt8
		}
		if gotSub != int8(rawSub) {
			t.Errorf("SubSat(%d, %d) = %d, want %d", a, b, gotSub, rawSub)
		}
	})
}
