package eml

import (
	"runtime"
	"testing"
)

func TestSIMDStubs(t *testing.T) {
	// Explicitly call all internal functions to get coverage on stubs and implementations
	a, b, res := make([]float64, 8), make([]float64, 8), make([]float64, 8)
	
	// If we are on ARM64, calling AVX functions will hit the stubs in simd_amd64_stub.go
	if runtime.GOARCH == "arm64" {
		addAVX2(a, b, res)
		subAVX2(a, b, res)
		mulAVX2(a, b, res)
		divAVX2(a, b, res)
		addScalarAVX2(a, 1.0, res)
		mulScalarAVX2(a, 1.0, res)
		sqrtAVX2(a, res)
		
		addAVX512(a, b, res)
		subAVX512(a, b, res)
		mulAVX512(a, b, res)
		divAVX512(a, b, res)
		addScalarAVX512(a, 1.0, res)
		mulScalarAVX512(a, 1.0, res)
		sqrtAVX512(a, res)
		
		detectAMD64SIMD()
		// cpuid is also a stub on arm64
		_, _, _, _ = cpuid(0, 0)
	}
	
	// If we are on AMD64, calling Neon functions will hit the stubs in simd_arm64_stub.go
	if runtime.GOARCH == "amd64" {
		addNEON(a, b, res)
		subNEON(a, b, res)
		mulNEON(a, b, res)
		divNEON(a, b, res)
		addScalarNEON(a, 1.0, res)
		mulScalarNEON(a, 1.0, res)
		sqrtNEON(a, res)
		detectARM64SIMD()
		
		// Also call the real AMD64 functions IF the hardware supports them
		if hasAVX2 {
			addAVX2(a, b, res)
			subAVX2(a, b, res)
			mulAVX2(a, b, res)
			divAVX2(a, b, res)
			addScalarAVX2(a, 1.0, res)
			mulScalarAVX2(a, 1.0, res)
			sqrtAVX2(a, res)
		}
		if hasAVX512 {
			addAVX512(a, b, res)
			subAVX512(a, b, res)
			mulAVX512(a, b, res)
			divAVX512(a, b, res)
			addScalarAVX512(a, 1.0, res)
			mulScalarAVX512(a, 1.0, res)
			sqrtAVX512(a, res)
		}
	}
}

func TestSIMDEdgeCases(t *testing.T) {
	// Test very small sizes to hit the n < 4 branches
	sizes := []int{0, 1, 2, 3}
	for _, n := range sizes {
		a := make([]float64, n)
		b := make([]float64, n)
		res := AddSIMD(a, b)
		if len(res) != n {
			t.Errorf("size %d: got length %d", n, len(res))
		}
		
		_ = SubSIMD(a, b)
		_ = MulSIMD(a, b)
		_ = DivSIMD(a, b)
		_ = AddScalarSIMD(a, 1.0)
		_ = MulScalarSIMD(a, 1.0)
		_ = SqrtSIMD(a)
		_ = ExpSIMD(a)
		_ = LogSIMD(a)
		_ = SinSIMD(a)
		_ = CosSIMD(a)
		_, _ = SinCosSIMD(a)
		_ = TanSIMD(a)
	}
	
	// Test large enough size to trigger concurrency and chunkSize limit
	n := 1000000
	a := make([]float64, n)
	b := make([]float64, n)
	_ = AddSIMD(a, b)
	_ = ExpSIMD(a)
	_ = SqrtSIMD(a)
}

func TestForceScalarFallbacks(t *testing.T) {
	// Temporarily disable SIMD flags to hit scalar fallbacks
	oldAVX2, oldAVX512, oldNeon := hasAVX2, hasAVX512, hasNeon
	hasAVX2, hasAVX512, hasNeon = false, false, false
	defer func() {
		hasAVX2, hasAVX512, hasNeon = oldAVX2, oldAVX512, oldNeon
	}()
	
	n := 10
	a, b, res := make([]float64, n), make([]float64, n), make([]float64, n)
	EmlSIMD(a, b, res) // hits scalarEml
	_ = AddSIMD(a, b)  // hits concurrent fallback
}

func TestRemainingFunctions(t *testing.T) {
	// Call remaining batch functions
	a := make([]float64, 4)
	_ = SinhBatch(a)
	_ = CoshBatch(a)
	_ = TanhBatch(a)
	_ = AsinhBatch(a)
	_ = AcoshBatch(a)
	_ = AtanhBatch(a)
	
	// Has* functions
	_ = HasSSE4()
	_ = HasAVX2()
	_ = HasAVX512()
	_ = HasNeon()
	_ = HasNeonDot()
	
	// Error type
	err := &EMLError{message: "test"}
	_ = err.Error()
	
	// EmlBatch success
	_ = EmlBatch(a, a, func(x, y, r []float64) error { return nil })
	
	// EmlBatch length mismatch
	err2 := EmlBatch(a, make([]float64, 5), func(x, y, r []float64) error { return nil })
	if err2 != ErrLengthMismatch {
		t.Errorf("expected ErrLengthMismatch, got %v", err2)
	}
	
	// EmlBatch callback error
	err3 := EmlBatch(a, a, func(x, y, r []float64) error { return &EMLError{message: "fail"} })
	if err3 == nil || err3.Error() != "fail" {
		t.Errorf("expected test error, got %v", err3)
	}
}

func TestSIMDMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for length mismatch")
		}
	}()
	a := make([]float64, 4)
	b := make([]float64, 5)
	AddSIMD(a, b)
}
