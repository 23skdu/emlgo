package eml

import (
	"math"
	"testing"
)

func TestComplexFullCoverage(t *testing.T) {
	c128 := Complex(1.0, 2.0)
	if math.Abs(real(c128)-2.0251346) > 1e-5 {
		t.Errorf("Complex failed: %v", c128)
	}
	c64 := Complex64(1.0, 2.0)
	if math.Abs(float64(real(c64))-2.0251346) > 1e-4 {
		t.Errorf("Complex64 failed: %v", c64)
	}
	_ = ComplexOne(1.0 + 2i)
	_ = OneComplex(1.0 + 2i)

	c64Slice := []complex64{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i, 11 + 12i, 13 + 14i, 15 + 16i, 17 + 18i}
	dstC64 := make([]complex64, len(c64Slice))

	ComplexBatchC64(c64Slice, c64Slice, dstC64)
	_ = ComplexExpBatchC64(c64Slice)
	ComplexExpBatchToC64(c64Slice, dstC64)
	_ = ComplexLogBatchC64(c64Slice)
	ComplexLogBatchToC64(c64Slice, dstC64)
	_ = ComplexSinBatchC64(c64Slice)
	ComplexSinBatchToC64(c64Slice, dstC64)
	_ = ComplexCosBatchC64(c64Slice)
	ComplexCosBatchToC64(c64Slice, dstC64)
	_ = ComplexTanBatchC64(c64Slice)
	ComplexTanBatchToC64(c64Slice, dstC64)

	c128Slice := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i}
	dstC128 := make([]complex128, len(c128Slice))
	_ = ComplexSqrtBatch(c128Slice)
	ComplexSqrtBatchTo(c128Slice, dstC128)
	_ = ComplexSqrtBatchC64(c64Slice)
	ComplexSqrtBatchToC64(c64Slice, dstC64)

	_ = ComplexDotProduct(c128Slice, c128Slice)
	_ = ComplexDotProductC64(c64Slice, c64Slice)
}

func TestSimdF32Coverage(t *testing.T) {
	n := 20
	a := make([]float32, n)
	b := make([]float32, n)
	dst := make([]float32, n)
	for i := 0; i < n; i++ {
		a[i] = float32(i + 1)
		b[i] = float32(i + 2)
	}

	AddSIMDF32To(a, b, dst)
	AddSIMDF32To(a[:3], b[:3], dst[:3]) // small slice < 8

	_ = SubSIMDF32(a, b)
	SubSIMDF32To(a, b, dst)
	SubSIMDF32To(a[:3], b[:3], dst[:3])

	MulSIMDF32To(a, b, dst)
	MulSIMDF32To(a[:3], b[:3], dst[:3])

	_ = DivSIMDF32(a, b)
	DivSIMDF32To(a, b, dst)
	DivSIMDF32To(a[:3], b[:3], dst[:3])

	_ = AbsSIMDF32(a)
	_ = NegSIMDF32(a)
	_ = InvSIMDF32(a)
	_ = SinSIMDF32(a)
	_ = CosSIMDF32(a)
	_ = TanSIMDF32(a)

	// Test parallelizeGenericF32 with large slice >= SmallCutoff
	large := make([]float32, SmallCutoff+100)
	for i := range large {
		large[i] = float32(i%10) + 1.0
	}
	_ = AbsSIMDF32(large)
}

func TestPipelineCoverage(t *testing.T) {
	data := []float64{1.0, 4.0, 9.0}
	p := NewPipeline(len(data))
	p.Sqrt()
	p.Cos()
	res := p.Run(data)
	if len(res) != 3 {
		t.Errorf("Pipeline failed")
	}
}

func TestSimdExtraCoverage(t *testing.T) {
	_ = HasWasmSIMD()
	_ = FmaScalar(1.0, 2.0, 3.0)

	// Min/Max branchless NaN branches
	_ = MinBranchless(math.NaN(), 1.0)
	_ = MinBranchless(1.0, math.NaN())
	_ = MaxBranchless(math.NaN(), 1.0)
	_ = MaxBranchless(1.0, math.NaN())

	// Panics on length mismatch in simd
	func() {
		defer func() { _ = recover() }()
		_ = AddSIMD([]float64{1}, []float64{1, 2})
	}()
	func() {
		defer func() { _ = recover() }()
		_ = SubSIMD([]float64{1}, []float64{1, 2})
	}()
	func() {
		defer func() { _ = recover() }()
		_ = MulSIMD([]float64{1}, []float64{1, 2})
	}()
	func() {
		defer func() { _ = recover() }()
		_ = DivSIMD([]float64{1}, []float64{1, 2})
	}()
}
