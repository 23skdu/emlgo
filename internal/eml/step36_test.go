package eml

import (
	"math"
	"runtime"
	"testing"
)

// TestStep36_HonestCapabilities asserts that HasSVE() and HasNeonDot() do not lie:
// they must return false unless backed by actual assembly instructions.
// HasFMA() and HasNeon() must be true on ARM64 platforms.
func TestStep36_HonestCapabilities(t *testing.T) {
	if runtime.GOARCH == "arm64" {
		if !HasNeon() {
			t.Errorf("ARM64 requires HasNeon() == true")
		}
		if !HasFMA() {
			t.Errorf("ARM64 requires HasFMA() == true (hardware FMADD)")
		}
	}

	// SVE and NEON Dot have no real instruction implementations in this library,
	// so the capability probes must be honest and return false.
	if HasSVE() {
		t.Errorf("HasSVE() must be false until backed by genuine SVE instructions")
	}
	if HasNeonDot() {
		t.Errorf("HasNeonDot() must be false until backed by genuine NEON dot product instructions")
	}
}

// TestStep36_ScalarFMA verifies hardware FMADD behavior on arm64 and fallback elsewhere.
func TestStep36_ScalarFMA(t *testing.T) {
	cases := []struct {
		a, b, c float64
	}{
		{2.0, 3.0, 4.0},
		{0.0, 100.0, 5.0},
		{-2.5, 4.0, 10.0},
		{1e10, 1e-10, 1.0},
		{math.Pi, math.E, 2.0},
		{1e-300, 1e-10, 0.0},
	}

	for _, tc := range cases {
		got := FmaScalar(tc.a, tc.b, tc.c)
		want := math.FMA(tc.a, tc.b, tc.c)
		if math.Abs(got-want) > 1e-15*math.Abs(want)+1e-15 {
			t.Errorf("FmaScalar(%g, %g, %g) = %g, want %g", tc.a, tc.b, tc.c, got, want)
		}
	}
}

// TestStep36_NEONDifferentialSweep performs an exhaustive differential sweep
// comparing all elementwise ops across varied lengths and edge values.
func TestStep36_NEONDifferentialSweep(t *testing.T) {
	lengths := []int{0, 1, 2, 3, 4, 5, 7, 8, 9, 15, 16, 17, 31, 32, 33, 64, 127, 128, 255, 256, 1024}

	for _, n := range lengths {
		a := make([]float64, n)
		b := make([]float64, n)
		c := make([]float64, n)
		for i := 0; i < n; i++ {
			a[i] = float64(i)*0.75 + 1.25
			b[i] = float64(i)*0.25 + 0.5
			c[i] = float64(i)*0.1 + 0.05
		}

		// Add
		resAdd := make([]float64, n)
		dispatchAddSIMD(a, b, resAdd)
		for i := 0; i < n; i++ {
			want := a[i] + b[i]
			if resAdd[i] != want {
				t.Fatalf("AddSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resAdd[i], want)
			}
		}

		// Sub
		resSub := make([]float64, n)
		dispatchSubSIMD(a, b, resSub)
		for i := 0; i < n; i++ {
			want := a[i] - b[i]
			if resSub[i] != want {
				t.Fatalf("SubSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resSub[i], want)
			}
		}

		// Mul
		resMul := make([]float64, n)
		dispatchMulSIMD(a, b, resMul)
		for i := 0; i < n; i++ {
			want := a[i] * b[i]
			if resMul[i] != want {
				t.Fatalf("MulSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resMul[i], want)
			}
		}

		// Div
		resDiv := make([]float64, n)
		dispatchDivSIMD(a, b, resDiv)
		for i := 0; i < n; i++ {
			want := a[i] / b[i]
			if math.Abs(resDiv[i]-want) > 1e-15*math.Abs(want) {
				t.Fatalf("DivSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resDiv[i], want)
			}
		}

		// AddScalar
		resAddS := make([]float64, n)
		dispatchAddScalarSIMD(a, 3.5, resAddS)
		for i := 0; i < n; i++ {
			want := a[i] + 3.5
			if resAddS[i] != want {
				t.Fatalf("AddScalarSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resAddS[i], want)
			}
		}

		// MulScalar
		resMulS := make([]float64, n)
		dispatchMulScalarSIMD(a, 2.25, resMulS)
		for i := 0; i < n; i++ {
			want := a[i] * 2.25
			if resMulS[i] != want {
				t.Fatalf("MulScalarSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resMulS[i], want)
			}
		}

		// Sqrt
		resSqrt := make([]float64, n)
		dispatchSqrtSIMDTo(a, resSqrt)
		for i := 0; i < n; i++ {
			want := math.Sqrt(a[i])
			if resSqrt[i] != want {
				t.Fatalf("SqrtSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resSqrt[i], want)
			}
		}

		// Abs
		negA := make([]float64, n)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				negA[i] = -a[i]
			} else {
				negA[i] = a[i]
			}
		}
		resAbs := make([]float64, n)
		dispatchAbsSIMD(negA, resAbs)
		for i := 0; i < n; i++ {
			want := math.Abs(negA[i])
			if resAbs[i] != want {
				t.Fatalf("AbsSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resAbs[i], want)
			}
		}

		// Neg
		resNeg := make([]float64, n)
		dispatchNegSIMD(a, resNeg)
		for i := 0; i < n; i++ {
			want := -a[i]
			if resNeg[i] != want {
				t.Fatalf("NegSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resNeg[i], want)
			}
		}

		// Inv
		resInv := make([]float64, n)
		dispatchInvSIMD(a, resInv)
		for i := 0; i < n; i++ {
			want := 1.0 / a[i]
			if math.Abs(resInv[i]-want) > 1e-15*math.Abs(want) {
				t.Fatalf("InvSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resInv[i], want)
			}
		}

		// FMA vector: c + a * b
		resFma := make([]float64, n)
		fmaSIMD(a, b, c, resFma)
		for i := 0; i < n; i++ {
			want := a[i]*b[i] + c[i]
			if math.Abs(resFma[i]-want) > 1e-14*math.Abs(want)+1e-15 {
				t.Fatalf("FmaSIMD mismatch at n=%d, i=%d: got %g, want %g", n, i, resFma[i], want)
			}
		}
	}
}

// TestStep36_SaturatingInt8 verifies saturating arithmetic under NEON and fallback.
func TestStep36_SaturatingInt8(t *testing.T) {
	lengths := []int{0, 1, 7, 15, 16, 17, 31, 32, 33, 64, 100}

	for _, n := range lengths {
		a := make([]int8, n)
		b := make([]int8, n)
		for i := 0; i < n; i++ {
			a[i] = int8((i * 17) % 256)
			b[i] = int8((i * 29) % 256)
		}

		resAdd := make([]int8, n)
		dispatchAddSatInt8SIMD(a, b, resAdd)
		resAddTo := make([]int8, n)
		AddSatInt8SIMDTo(a, b, resAddTo)
		for i := 0; i < n; i++ {
			v := int32(a[i]) + int32(b[i])
			if v > 127 {
				v = 127
			} else if v < -128 {
				v = -128
			}
			if resAdd[i] != int8(v) || resAddTo[i] != int8(v) {
				t.Fatalf("AddSatInt8 mismatch at n=%d, i=%d: got %d, want %d", n, i, resAdd[i], int8(v))
			}
		}

		resSub := make([]int8, n)
		dispatchSubSatInt8SIMD(a, b, resSub)
		resSubTo := make([]int8, n)
		SubSatInt8SIMDTo(a, b, resSubTo)
		for i := 0; i < n; i++ {
			v := int32(a[i]) - int32(b[i])
			if v > 127 {
				v = 127
			} else if v < -128 {
				v = -128
			}
			if resSub[i] != int8(v) || resSubTo[i] != int8(v) {
				t.Fatalf("SubSatInt8 mismatch at n=%d, i=%d: got %d, want %d", n, i, resSub[i], int8(v))
			}
		}
	}

	// Length mismatch panics
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on length mismatch")
		}
	}()
	AddSatInt8SIMDTo(make([]int8, 2), make([]int8, 3), make([]int8, 2))
}

func TestStep36_SubSatLengthMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on length mismatch")
		}
	}()
	SubSatInt8SIMDTo(make([]int8, 2), make([]int8, 3), make([]int8, 2))
}

// TestStep36_InPlaceExecution tests that in-place vector operations (dst == src)
// execute correctly without corrupting overlapping elements.
func TestStep36_InPlaceExecution(t *testing.T) {
	n := 128
	a := make([]float64, n)
	b := make([]float64, n)
	for i := range a {
		a[i] = float64(i + 1)
		b[i] = 2.0
	}

	dispatchAddSIMD(a, b, a)
	for i := range a {
		want := float64(i+1) + 2.0
		if a[i] != want {
			t.Fatalf("In-place AddSIMD mismatch at %d: got %g, want %g", i, a[i], want)
		}
	}

	dispatchMulSIMD(a, b, a)
	for i := range a {
		want := (float64(i+1) + 2.0) * 2.0
		if a[i] != want {
			t.Fatalf("In-place MulSIMD mismatch at %d: got %g, want %g", i, a[i], want)
		}
	}

	dispatchSqrtSIMDTo(a, a)
	for i := range a {
		want := math.Sqrt((float64(i+1) + 2.0) * 2.0)
		if math.Abs(a[i]-want) > 1e-15 {
			t.Fatalf("In-place SqrtSIMD mismatch at %d: got %g, want %g", i, a[i], want)
		}
	}
}

// BenchmarkStep36_NEONAdd benchmarks vector addition.
func BenchmarkStep36_NEONAdd(b *testing.B) {
	n := 65536
	x := make([]float64, n)
	y := make([]float64, n)
	out := make([]float64, n)
	for i := range x {
		x[i] = float64(i)
		y[i] = float64(i * 2)
	}

	b.SetBytes(int64(n * 8 * 2))
	b.ResetTimer()
	for b.Loop() {
		dispatchAddSIMD(x, y, out)
	}
}

// BenchmarkStep36_NEONFma benchmarks fused multiply-add.
func BenchmarkStep36_NEONFma(b *testing.B) {
	n := 65536
	x := make([]float64, n)
	y := make([]float64, n)
	z := make([]float64, n)
	out := make([]float64, n)
	for i := range x {
		x[i] = float64(i)
		y[i] = 1.001
		z[i] = float64(i) * 0.5
	}

	b.SetBytes(int64(n * 8 * 3))
	b.ResetTimer()
	for b.Loop() {
		fmaSIMD(x, y, z, out)
	}
}
