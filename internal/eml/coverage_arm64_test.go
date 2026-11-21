//go:build arm64 && !purego

package eml

import (
	"math"
	"testing"
)

func TestARM64NEONFloat64KernelsExact(t *testing.T) {
	sizes := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 15, 16, 17, 33}
	for _, n := range sizes {
		a := make([]float64, n)
		b := make([]float64, n)
		c := make([]float64, n)
		dst := make([]float64, n)
		for i := range n {
			a[i] = float64(i)*1.25 - 4.5
			b[i] = float64(i+1) * 0.75
			c[i] = float64(i)*0.5 - 2.0
		}
		if n >= 4 {
			a[0] = math.NaN()
			a[1] = math.Inf(1)
			a[2] = math.Inf(-1)
			a[3] = math.Copysign(0, -1)
		}

		check := func(name string, got, want []float64) {
			t.Helper()
			for i := range got {
				if math.IsNaN(want[i]) {
					if !math.IsNaN(got[i]) {
						t.Errorf("%s(n=%d)[%d] = %v, want NaN", name, n, i, got[i])
					}
					continue
				}
				if math.Float64bits(got[i]) != math.Float64bits(want[i]) && got[i] != want[i] {
					t.Errorf("%s(n=%d)[%d] = %v (%x), want %v (%x)", name, n, i, got[i], math.Float64bits(got[i]), want[i], math.Float64bits(want[i]))
				}
			}
		}

		want := make([]float64, n)

		// Add
		arm64AddSIMD(a, b, dst)
		for i := range n {
			want[i] = a[i] + b[i]
		}
		check("arm64AddSIMD", dst, want)

		// Sub
		arm64SubSIMD(a, b, dst)
		for i := range n {
			want[i] = a[i] - b[i]
		}
		check("arm64SubSIMD", dst, want)

		// Mul
		arm64MulSIMD(a, b, dst)
		for i := range n {
			want[i] = a[i] * b[i]
		}
		check("arm64MulSIMD", dst, want)

		// Div
		arm64DivSIMD(a, b, dst)
		for i := range n {
			want[i] = a[i] / b[i]
		}
		check("arm64DivSIMD", dst, want)

		// AddScalar
		arm64AddScalarSIMD(a, 2.5, dst)
		for i := range n {
			want[i] = a[i] + 2.5
		}
		check("arm64AddScalarSIMD", dst, want)

		// MulScalar
		arm64MulScalarSIMD(a, -3.25, dst)
		for i := range n {
			want[i] = a[i] * -3.25
		}
		check("arm64MulScalarSIMD", dst, want)

		// Sqrt (non-negative + special values)
		arm64SqrtSIMD(b, dst)
		for i := range n {
			want[i] = math.Sqrt(b[i])
		}
		check("arm64SqrtSIMD", dst, want)

		// Abs
		arm64AbsSIMD(a, dst)
		for i := range n {
			want[i] = math.Abs(a[i])
		}
		check("arm64AbsSIMD", dst, want)

		// Neg
		arm64NegSIMD(a, dst)
		for i := range n {
			want[i] = -a[i]
		}
		check("arm64NegSIMD", dst, want)

		// Inv
		arm64InvSIMD(b, dst)
		for i := range n {
			want[i] = 1.0 / b[i]
		}
		check("arm64InvSIMD", dst, want)

		// FMA
		fmaSIMD(a, b, c, dst)
		for i := range n {
			want[i] = math.FMA(a[i], b[i], c[i])
		}
		check("fmaSIMD", dst, want)
	}

	// Direct zero-length assembly entry/exit paths
	addNEON(nil, nil, nil)
	subNEON(nil, nil, nil)
	mulNEON(nil, nil, nil)
	divNEON(nil, nil, nil)
	addScalarNEON(nil, 1.0, nil)
	mulScalarNEON(nil, 1.0, nil)
	sqrtNEON(nil, nil)
	absNEON(nil, nil)
	negNEON(nil, nil)
	invNEON(nil, nil)
	fmaNEON(nil, nil, nil, nil)
	addSVE(nil, nil, nil)

	// fmaScalar
	if got, want := fmaScalar(1.5, 2.5, 3.25), math.FMA(1.5, 2.5, 3.25); got != want {
		t.Fatalf("fmaScalar = %v, want %v", got, want)
	}
}

func TestARM64NEONFloat32AndInt8KernelsExact(t *testing.T) {
	sizes := []int{0, 1, 2, 3, 4, 5, 7, 8, 15, 16, 17, 31, 32, 33, 65}
	for _, n := range sizes {
		a := make([]float32, n)
		b := make([]float32, n)
		c := make([]float32, n)
		dst := make([]float32, n)
		for i := range n {
			a[i] = float32(i)*1.25 - 4.5
			b[i] = float32(i+1) * 0.75
			c[i] = float32(i)*0.5 - 2.0
		}

		checkF32 := func(name string, got, want []float32) {
			t.Helper()
			for i := range got {
				if math.IsNaN(float64(want[i])) {
					if !math.IsNaN(float64(got[i])) {
						t.Errorf("%s(n=%d)[%d] = %v, want NaN", name, n, i, got[i])
					}
					continue
				}
				if got[i] != want[i] {
					t.Errorf("%s(n=%d)[%d] = %v, want %v", name, n, i, got[i], want[i])
				}
			}
		}

		want := make([]float32, n)
		dispatchAddSIMDF32To(a, b, dst)
		for i := range n {
			want[i] = a[i] + b[i]
		}
		checkF32("addNEONF32", dst, want)

		dispatchSubSIMDF32To(a, b, dst)
		for i := range n {
			want[i] = a[i] - b[i]
		}
		checkF32("subNEONF32", dst, want)

		dispatchMulSIMDF32To(a, b, dst)
		for i := range n {
			want[i] = a[i] * b[i]
		}
		checkF32("mulNEONF32", dst, want)

		dispatchDivSIMDF32To(a, b, dst)
		for i := range n {
			want[i] = a[i] / b[i]
		}
		checkF32("divNEONF32", dst, want)

		dispatchSqrtSIMDF32To(b, dst)
		for i := range n {
			want[i] = float32(math.Sqrt(float64(b[i])))
		}
		checkF32("sqrtNEONF32", dst, want)

		dispatchAddScalarSIMDF32To(a, 1.5, dst)
		for i := range n {
			want[i] = a[i] + 1.5
		}
		checkF32("addScalarNEONF32", dst, want)

		dispatchMulScalarSIMDF32To(a, -2.5, dst)
		for i := range n {
			want[i] = a[i] * -2.5
		}
		checkF32("mulScalarNEONF32", dst, want)

		dispatchAbsSIMDF32To(a, dst)
		for i := range n {
			want[i] = float32(math.Abs(float64(a[i])))
		}
		checkF32("absNEONF32", dst, want)

		dispatchNegSIMDF32To(a, dst)
		for i := range n {
			want[i] = -a[i]
		}
		checkF32("negNEONF32", dst, want)

		dispatchInvSIMDF32To(b, dst)
		for i := range n {
			want[i] = 1.0 / b[i]
		}
		checkF32("invNEONF32", dst, want)

		dispatchFmaSIMDF32To(a, b, c, dst)
		for i := range n {
			want[i] = float32(math.FMA(float64(a[i]), float64(b[i]), float64(c[i])))
		}
		checkF32("fmaNEONF32", dst, want)

		// Saturating Int8 NEON
		ia := make([]int8, n)
		ib := make([]int8, n)
		idst := make([]int8, n)
		for i := range n {
			ia[i] = int8((i*37)%256 - 128)
			ib[i] = int8((i*53)%256 - 128)
		}
		dispatchAddSatInt8SIMD(ia, ib, idst)
		for i := range n {
			v := max(-128, min(127, int32(ia[i])+int32(ib[i])))
			if idst[i] != int8(v) {
				t.Fatalf("addSatInt8NEON(n=%d)[%d] = %d, want %d", n, i, idst[i], v)
			}
		}
		dispatchSubSatInt8SIMD(ia, ib, idst)
		for i := range n {
			v := max(-128, min(127, int32(ia[i])-int32(ib[i])))
			if idst[i] != int8(v) {
				t.Fatalf("subSatInt8NEON(n=%d)[%d] = %d, want %d", n, i, idst[i], v)
			}
		}
	}

	// Direct zero-length F32 & Int8 assembly entry/exit paths
	addNEONF32(nil, nil, nil)
	subNEONF32(nil, nil, nil)
	mulNEONF32(nil, nil, nil)
	divNEONF32(nil, nil, nil)
	sqrtNEONF32(nil, nil)
	addScalarNEONF32(nil, 1, nil)
	mulScalarNEONF32(nil, 1, nil)
	absNEONF32(nil, nil)
	negNEONF32(nil, nil)
	invNEONF32(nil, nil)
	fmaNEONF32(nil, nil, nil, nil)
	addSatInt8NEON(nil, nil, nil)
	subSatInt8NEON(nil, nil, nil)
}

func TestARM64DispatchAndSIMDCoverage100(t *testing.T) {
	if HasAVX512VNNI() {
		t.Error("HasAVX512VNNI should be false on arm64")
	}
	if !HasNeon() || !HasFMA() {
		t.Error("expected HasNeon and HasFMA on arm64")
	}

	// Cover large (> SmallCutoff) parallel paths for Floor/Ceil/Trunc/Round/Max/Min
	n := SmallCutoff() + 16
	a := make([]float64, n)
	b := make([]float64, n)
	dst := make([]float64, n)
	for i := range n {
		a[i] = float64(i)*0.3 - 10.0
		b[i] = float64(n-i)*0.2 - 5.0
	}
	a[0] = math.NaN()
	b[1] = math.NaN()

	FloorSIMDTo(a, dst)
	CeilSIMDTo(a, dst)
	TruncSIMDTo(a, dst)
	RoundSIMDTo(a, dst)
	MaxSIMDTo(a, b, dst)
	if !math.IsNaN(dst[0]) || !math.IsNaN(dst[1]) {
		t.Fatal("MaxSIMDTo failed to propagate NaN")
	}
	MinSIMDTo(a, b, dst)
	if !math.IsNaN(dst[0]) || !math.IsNaN(dst[1]) {
		t.Fatal("MinSIMDTo failed to propagate NaN")
	}

	// Cover panic paths for mismatched slice lengths
	expectPanic := func(name string, fn func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic on length mismatch", name)
			}
		}()
		fn()
	}

	short := make([]float64, 1)
	two := make([]float64, 2)
	expectPanic("FloorSIMDTo", func() { FloorSIMDTo(short, two) })
	expectPanic("CeilSIMDTo", func() { CeilSIMDTo(short, two) })
	expectPanic("TruncSIMDTo", func() { TruncSIMDTo(short, two) })
	expectPanic("RoundSIMDTo", func() { RoundSIMDTo(short, two) })
	expectPanic("MaxSIMDTo", func() { MaxSIMDTo(short, two, two) })
	expectPanic("MinSIMDTo", func() { MinSIMDTo(short, two, two) })
	expectPanic("DotInt8SIMD", func() { _ = DotInt8SIMD([]uint8{1}, []int8{1, 2}) })
	expectPanic("SqrtSIMDF32To", func() { SqrtSIMDF32To([]float32{1}, []float32{1, 2}) })
}

func FuzzARM64NEONKernels(f *testing.F) {
	f.Add(1.5, 2.5, -3.0, int8(100), int8(50))
	f.Add(-1e-10, 1e10, 0.0, int8(-100), int8(-50))
	f.Add(0.0, 1.0, -1.0, int8(127), int8(127))

	f.Fuzz(func(t *testing.T, x, y, z float64, i8a, i8b int8) {
		if math.IsNaN(x) || math.IsNaN(y) || math.IsNaN(z) {
			return
		}
		a := []float64{x, y, z, x + 1, y - 1}
		b := []float64{y, z, x, 1.5, 2.5}
		c := []float64{z, x, y, -0.5, 0.25}
		dst := make([]float64, len(a))

		arm64AddSIMD(a, b, dst)
		for i := range a {
			want := a[i] + b[i]
			if !math.IsNaN(want) && dst[i] != want {
				t.Fatalf("arm64AddSIMD[%d] = %v, want %v", i, dst[i], want)
			}
		}

		fmaSIMD(a, b, c, dst)
		for i := range a {
			want := math.FMA(a[i], b[i], c[i])
			if !math.IsNaN(want) && dst[i] != want {
				t.Fatalf("fmaSIMD[%d] = %v, want %v", i, dst[i], want)
			}
		}

		// Saturating Int8 across 17 elements (16 NEON + 1 scalar tail)
		ia := make([]int8, 17)
		ib := make([]int8, 17)
		idst := make([]int8, 17)
		for i := range ia {
			ia[i] = i8a
			ib[i] = i8b
		}
		dispatchAddSatInt8SIMD(ia, ib, idst)
		wantSat := int8(max(-128, min(127, int32(i8a)+int32(i8b))))
		for i, got := range idst {
			if got != wantSat {
				t.Fatalf("addSatInt8[%d] = %d, want %d", i, got, wantSat)
			}
		}
	})
}
