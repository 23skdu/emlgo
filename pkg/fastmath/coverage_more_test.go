package fastmath

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestFastmathComprehensive(t *testing.T) {
	// Guardrails float64 & float32
	if ExpClamped(-800.0) != 0.0 {
		t.Errorf("ExpClamped(-800) should be 0")
	}
	_ = ExpClamped(800.0)
	_ = ExpClamped(1.0)

	if ExpClampedF32(-100.0) != 0.0 {
		t.Errorf("ExpClampedF32(-100) should be 0")
	}
	_ = ExpClampedF32(100.0)
	_ = ExpClampedF32(1.0)

	_ = LnClipped(0.0, -1.0) // test eps <= 0
	_ = LnClipped(1.0, 0.5)  // test y >= eps
	_ = LnClipped(0.1, 0.5)  // test y < eps

	_ = LnClippedF32(0.0, -1.0)
	_ = LnClippedF32(1.0, 0.5)
	_ = LnClippedF32(0.1, 0.5)

	_ = LnRegularized(2.0, -1.0) // eps <= 0
	_ = LnRegularizedF32(2.0, -1.0)
	_ = LnRegularizedF32(2.0, 1e-4)

	// fast_eml.go
	_ = FastEml(1.0, 2.0)
	_ = FastEmlF32(1.0, 2.0)
	_ = FastEmlBitCast(1.0, 2.0)
	_ = FastEmlClipped(1.0, 2.0, 1e-6)
	_ = FastEmlClippedF32(1.0, 2.0, 1e-6)
	_ = FastEmlRegularizedF32(1.0, 2.0, 1e-6)

	x64 := []float64{1.0, 2.0}
	y64 := []float64{2.0, 3.0}
	dst64 := make([]float64, 2)
	bRes := FastEmlBatch(x64, y64)
	FastEmlBatchTo(x64, y64, dst64)
	if len(bRes) != 2 {
		t.Errorf("FastEmlBatch length: %d", len(bRes))
	}

	x32 := []float32{1.0, 2.0}
	y32 := []float32{2.0, 3.0}
	dst32 := make([]float32, 2)
	bRes32 := FastEmlBatchF32(x32, y32)
	FastEmlBatchToF32(x32, y32, dst32)
	if len(bRes32) != 2 {
		t.Errorf("FastEmlBatchF32 length: %d", len(bRes32))
	}

	// Panics on length mismatch in fast_eml
	func() {
		defer func() { _ = recover() }()
		_ = FastEmlBatch([]float64{1.0}, []float64{1.0, 2.0})
	}()
	func() {
		defer func() { _ = recover() }()
		FastEmlBatchTo([]float64{1.0}, []float64{1.0, 2.0}, dst64)
	}()
	func() {
		defer func() { _ = recover() }()
		_ = FastEmlBatchF32([]float32{1.0}, []float32{1.0, 2.0})
	}()
	func() {
		defer func() { _ = recover() }()
		FastEmlBatchToF32([]float32{1.0}, []float32{1.0, 2.0}, dst32)
	}()
	func() {
		defer func() { _ = recover() }()
		SqrtBatchToF32([]float32{1.0}, dst32)
	}()

	// fast_exp.go boundaries
	_ = FastExpBitCast(math.NaN())
	_ = FastExpBitCast(-1000.0)
	_ = FastExpBitCast(-650.0) // val < 0 branch
	_ = FastExpBitCast(1000.0)
	_ = FastExpMinimax(math.NaN())
	_ = FastExpMinimax(-800.0)
	_ = FastExpMinimax(800.0)
	_ = FastExpF32(float32(math.NaN()))
	_ = FastExpF32(-100.0)
	_ = FastExpF32(100.0)
	_ = FastExpF32(1.0)

	// fast_log.go boundaries
	_ = FastLogMinimax(math.NaN())
	_ = FastLogMinimax(math.Inf(1))
	if !math.IsNaN(FastLogMinimax(-1.0)) {
		t.Errorf("FastLogMinimax(-1) should be NaN")
	}
	if !math.IsInf(FastLogMinimax(0.0), -1) {
		t.Errorf("FastLogMinimax(0) should be -Inf")
	}
	_ = FastLogMinimax(1e-315) // subnormal path
	_ = FastLogMinimax(1.5)   // m > sqrt2 branch
	_ = FastLogChebyshev(math.NaN())
	_ = FastLogChebyshev(math.Inf(1))
	_ = FastLogChebyshev(1.0)
	_ = FastLogChebyshev(1.5) // m > sqrt2 branch
	_ = FastLogChebyshev(0.0)
	_ = FastLogChebyshev(-1.0)
	_ = FastLogChebyshev(1e-315)

	_ = FastLogF32(float32(math.NaN()))
	_ = FastLogF32(-1.0)
	_ = FastLogF32(0.0)
	_ = FastLogF32(float32(math.Inf(1)))
	_ = FastLogF32(1e-40) // subnormal float32
	_ = FastLogF32(1.5)   // m > sqrt2 branch
	_ = FastLogF32(1.0)

	// fastmath.go functions
	_ = SqrtBatch(x64)
	SqrtBatchTo(x64, dst64)
	_ = SqrtF32(4.0)
	_ = SqrtBatchF32(x32)
	SqrtBatchToF32(x32, dst32)
	_ = ExpF32(1.0)
	_ = SinF32(1.0)
	_ = CosF32(1.0)
	_ = LogF32(1.0)

	// Complex64 batch ops
	c64 := []complex64{1 + 2i, 3 + 4i}
	dstC64 := make([]complex64, 2)
	_ = SqrtC64(c64[0])
	_ = SqrtBatchC64(c64)
	SqrtBatchToC64(c64, dstC64)
	_ = ExpC64(c64[0])
	_ = ExpBatchC64(c64)
	ExpBatchToC64(c64, dstC64)
	_ = LogC64(c64[0])
	_ = LogBatchC64(c64)
	LogBatchToC64(c64, dstC64)
	_ = SinC64(c64[0])
	_ = SinBatchC64(c64)
	SinBatchToC64(c64, dstC64)
	_ = CosC64(c64[0])
	_ = CosBatchC64(c64)
	CosBatchToC64(c64, dstC64)
	_ = TanC64(c64[0])
	_ = TanBatchC64(c64)
	TanBatchToC64(c64, dstC64)

	// Complex128 batch ops
	c128 := []complex128{1 + 2i, 3 + 4i}
	dstC128 := make([]complex128, 2)
	_ = SqrtC128(c128[0])
	_ = SqrtBatchC128(c128)
	SqrtBatchToC128(c128, dstC128)
	_ = ExpC128(c128[0])
	_ = ExpBatchC128(c128)
	ExpBatchToC128(c128, dstC128)
	_ = LogC128(c128[0])
	_ = LogBatchC128(c128)
	LogBatchToC128(c128, dstC128)
	_ = SinC128(c128[0])
	_ = SinBatchC128(c128)
	SinBatchToC128(c128, dstC128)
	_ = CosC128(c128[0])
	_ = CosBatchC128(c128)
	CosBatchToC128(c128, dstC128)
	_ = TanC128(c128[0])
	_ = TanBatchC128(c128)
	TanBatchToC128(c128, dstC128)
	_ = cmplx.Abs(1 + 1i)
}
