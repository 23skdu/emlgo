package eml

import (
	"math"
	"sync/atomic"
	"testing"
)

func TestFusedEdgeCases(t *testing.T) {
	nan := math.NaN()
	if !math.IsNaN(MinBranchless(nan, nan)) {
		t.Fatal("MinBranchless(NaN, NaN) should be NaN")
	}
	if !math.IsNaN(MaxBranchless(nan, nan)) {
		t.Fatal("MaxBranchless(NaN, NaN) should be NaN")
	}

	if MinBranchless(1.0, nan) != 1.0 {
		t.Fatal("MinBranchless(1, NaN) should be 1")
	}
	if MinBranchless(nan, 1.0) != 1.0 {
		t.Fatal("MinBranchless(NaN, 1) should be 1")
	}

	if MaxBranchless(1.0, nan) != 1.0 {
		t.Fatal("MaxBranchless(1, NaN) should be 1")
	}
	if MaxBranchless(nan, 1.0) != 1.0 {
		t.Fatal("MaxBranchless(NaN, 1) should be 1")
	}
}

func TestPipelineEdgeCases(t *testing.T) {
	// Empty pipeline Run
	pEmpty := NewPipeline(10)
	got := pEmpty.Run([]float64{1, 2, 3})
	if len(got) != 3 || got[0] != 1 {
		t.Fatalf("empty pipeline Run failed: %v", got)
	}

	// Empty pipeline RunTo
	dst := make([]float64, 3)
	pEmpty.RunTo([]float64{1, 2, 3}, dst)
	if dst[0] != 1 {
		t.Fatalf("empty pipeline RunTo failed: %v", dst)
	}

	// Custom step
	pCustom := NewPipeline(10)
	pCustom.Custom(func(src, dst []float64) {
		for i := range src {
			dst[i] = src[i] * 5
		}
	})
	resCustom := pCustom.Run([]float64{2})
	if resCustom[0] != 10 {
		t.Fatalf("Custom step failed: %v", resCustom)
	}

	// Affine scale=1, bias=0
	pAff1 := NewPipeline(10)
	pAff1.addAffine(1, 0)
	resAff1 := pAff1.Run([]float64{3})
	if resAff1[0] != 3 {
		t.Fatalf("addAffine(1, 0) failed: %v", resAff1)
	}

	// Affine scale=1, bias=2
	pAff2 := NewPipeline(10)
	pAff2.addAffine(1, 2)
	resAff2 := pAff2.Run([]float64{3})
	if resAff2[0] != 5 {
		t.Fatalf("addAffine(1, 2) failed: %v", resAff2)
	}

	// Affine scale=2, bias=0
	pAff3 := NewPipeline(10)
	pAff3.addAffine(2, 0)
	resAff3 := pAff3.Run([]float64{3})
	if resAff3[0] != 6 {
		t.Fatalf("addAffine(2, 0) failed: %v", resAff3)
	}

	// grow with p.buf[1] != nil
	pGrow := NewPipeline(2)
	pGrow.Exp()
	pGrow.Run([]float64{1, 2})
	pGrow.Run([]float64{1, 2, 3, 4})

	// RunTo with multiple steps
	pMulti := NewPipeline(10)
	pMulti.Exp().Log()
	outMulti := make([]float64, 2)
	pMulti.RunTo([]float64{1, 2}, outMulti)
}

func TestSIMDEdgeCases(t *testing.T) {
	addSVE(nil, nil, nil)

	oldFMA := hasFMA
	hasFMA = false
	fmaRes := FmaScalar(2, 3, 4)
	hasFMA = oldFMA
	if fmaRes != 10 {
		t.Fatalf("FmaScalar without FMA = %f, want 10", fmaRes)
	}

	parallelizeGeneric(nil, nil, nil)
	parallelizeSinCos(nil, nil, nil)

	ForEachChunk(0, nil)
	var ran atomic.Bool
	ForEachChunk(100000, func(start, end int) {
		ran.Store(true)
	})
	if !ran.Load() {
		t.Fatal("ForEachChunk did not run")
	}
}

func TestParallelChunkingLarge(t *testing.T) {
	const n = 100000
	f32a := make([]float32, n)
	f32b := make([]float32, n)
	f32res := make([]float32, n)
	f32a[0] = -5.0 // Test v < 0 in AbsSIMDF32
	for i := 1; i < n; i++ {
		f32a[i] = float32(i + 1)
		f32b[i] = 2.0
	}

	AddSIMDF32To(f32a, f32b, f32res)
	MulSIMDF32To(f32a, f32b, f32res)
	SubSIMDF32To(f32a, f32b, f32res)
	_ = SubSIMDF32(f32a, f32b)
	DivSIMDF32To(f32a, f32b, f32res)
	_ = DivSIMDF32(f32a, f32b)
	_ = AbsSIMDF32(f32a)

	parallelizeGenericF32(f32a, f32res, func(src, dst []float32) {
		copy(dst, src)
	})

	// Test parallelizeGenericF32 with SmallCutoff=1 to hit numWorkers > n, hi > n, and lo >= n
	oldCutoff := SmallCutoff
	SmallCutoff = 1
	short1 := make([]float32, 1)
	parallelizeGenericF32(short1, short1, func(s, d []float32) {})
	short17 := make([]float32, 17)
	parallelizeGenericF32(short17, short17, func(s, d []float32) {})
	SmallCutoff = oldCutoff

	// Slice mismatch panics in simd_f32.go
	func() { defer func() { _ = recover() }(); AddSIMDF32To(f32a, f32b[:10], f32res) }()
	func() { defer func() { _ = recover() }(); MulSIMDF32To(f32a, f32b[:10], f32res) }()
	func() { defer func() { _ = recover() }(); SubSIMDF32To(f32a, f32b[:10], f32res) }()
	func() { defer func() { _ = recover() }(); SubSIMDF32(f32a, f32b[:10]) }()
	func() { defer func() { _ = recover() }(); DivSIMDF32To(f32a, f32b[:10], f32res) }()
	func() { defer func() { _ = recover() }(); DivSIMDF32(f32a, f32b[:10]) }()

	f64a := make([]float64, n)
	f64b := make([]float64, n)
	f64c := make([]float64, n)
	f64res := make([]float64, n)
	for i := 0; i < n; i++ {
		f64a[i] = float64(i + 1)
		f64b[i] = 2.0
		f64c[i] = 1.0
	}

	FmaSIMDTo(f64a, f64b, f64c, f64res)
	Log2SIMDTo(f64a, f64res)
	Log10SIMDTo(f64a, f64res)
	sinOut := make([]float64, n)
	cosOut := make([]float64, n)
	SinCosSIMDTo(f64a, sinOut, cosOut)
	ExpSIMDTo(f64a, f64res)

	c128a := make([]complex128, n)
	c128b := make([]complex128, n)
	for i := 0; i < n; i++ {
		c128a[i] = complex(float64(i), 1)
		c128b[i] = complex(1, float64(i))
	}
	dot128 := ComplexDotProduct(c128a, c128b)
	if dot128 == 0 {
		t.Fatal("ComplexDotProduct returned 0")
	}
	if ComplexDotProduct(nil, nil) != 0 {
		t.Fatal("ComplexDotProduct(nil, nil) should be 0")
	}

	func() {
		defer func() { _ = recover() }()
		ComplexDotProduct(c128a, c128b[:10])
	}()

	c64a := make([]complex64, n)
	c64b := make([]complex64, n)
	for i := 0; i < n; i++ {
		c64a[i] = complex(float32(i), 1)
		c64b[i] = complex(1, float32(i))
	}
	dot64 := ComplexDotProductC64(c64a, c64b)
	if dot64 == 0 {
		t.Fatal("ComplexDotProductC64 returned 0")
	}
	if ComplexDotProductC64(nil, nil) != 0 {
		t.Fatal("ComplexDotProductC64(nil, nil) should be 0")
	}

	func() {
		defer func() { _ = recover() }()
		ComplexDotProductC64(c64a, c64b[:10])
	}()
}
