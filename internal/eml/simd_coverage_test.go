package eml

import (
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
	}
	
	// Test large enough size to trigger concurrency if implemented
	n := 10000
	a := make([]float64, n)
	b := make([]float64, n)
	_ = AddSIMD(a, b)
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
