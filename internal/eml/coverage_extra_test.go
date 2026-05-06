package eml

import (
	"math"
	"testing"
)

func TestExtraEmlCoverage(t *testing.T) {
	empty := []float64{}
	small := []float64{1, 2, 3}
	large := make([]float64, 1000)
	huge := make([]float64, 1000000)
	resSmall := make([]float64, 3)
	resLarge := make([]float64, 1000)
	resHuge := make([]float64, 1000000)
	_ = resHuge
	
	// Fused
	f := fusedOps{}
	f.ExpAddBatch(empty, empty, empty)
	f.ExpAddBatch(small, small, resSmall)
	f.ExpAddBatch(large, large, resLarge)
	
	f.LogDivBatch(empty, empty, empty)
	f.LogDivBatch(small, small, resSmall)
	f.LogDivBatch(large, large, resLarge)
	
	f.LogSubBatch(empty, empty, empty)
	f.LogSubBatch(small, small, resSmall)
	f.LogSubBatch(large, large, resLarge)
	
	f.ExpMulBatch(empty, empty, empty)
	f.ExpMulBatch(small, small, resSmall)
	f.ExpMulBatch(large, large, resLarge)
	
	// FMA
	FmaSIMD(empty, empty, empty)
	FmaSIMD(small, small, small)
	FmaSIMD(large, large, large)
	FmaSIMDTo(empty, empty, empty, empty)
	
	// Expanded
	Log2SIMDTo(empty, empty)
	Log2SIMDTo(small, resSmall)
	Log2SIMDTo(large, resLarge)
	
	Log10SIMDTo(empty, empty)
	Log10SIMDTo(small, resSmall)
	Log10SIMDTo(large, resLarge)

	
	// Large chunks
	GetParallelChunkSize(1000000)
	ExpAddBatch(huge, huge)
}

func TestEmlPanics(t *testing.T) {
	a5 := make([]float64, 5)
	a4 := make([]float64, 4)
	
	tests := []func(){
		func() { ExpSIMDTo(a5, a4) },
		func() { LogSIMDTo(a5, a4) },
		func() { SqrtSIMDTo(a5, a4) },
		func() { SinSIMDTo(a5, a4) },
		func() { CosSIMDTo(a5, a4) },
		func() { TanSIMDTo(a5, a4) },
		func() { AbsSIMDTo(a5, a4) },
		func() { NegSIMDTo(a5, a4) },
		func() { InvSIMDTo(a5, a4) },
		func() { SinCosSIMDTo(a5, a4, a5) },
		func() { SinCosSIMDTo(a5, a5, a4) },
		func() { AddScalarSIMDTo(a5, 1.0, a4) },
		func() { MulScalarSIMDTo(a5, 1.0, a4) },
		func() { SIMD(a5, a4, a5) },
		func() { SIMD(a5, a5, a4) },
		func() { ExpMulTo(a5, a5, a4) },
		func() { ExpAddTo(a5, a5, a4) },
		func() { LogDivTo(a5, a5, a4) },
		func() { LogSubTo(a5, a5, a4) },
		func() { Log2SIMDTo(a5, a4) },
		func() { Log10SIMDTo(a5, a4) },
		func() { FmaSIMD(a5, a4, a5) },
		func() { FmaSIMD(a5, a5, a4) },
		func() { FmaSIMDTo(a5, a4, a5, a5) },
		func() { FmaSIMDTo(a5, a5, a4, a5) },
		func() { FmaSIMDTo(a5, a5, a5, a4) },
	}
	
	for _, fn := range tests {
		assertPanic(t, fn)
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestNativeWrappersAndStubs(t *testing.T) {
	nativeAtan(0)
	nativeAtan2(0, 0)
	nativeAsin(0)
	nativeAcos(0)
	_ = ErrLengthMismatch.Error()
	nativeMax(5, math.NaN())
	nativeMax(5, 10)
	nativeMin(5, math.NaN())
	nativeMin(10, 5)
	
	// AMD64 stubs
	addAVX2(nil, nil, nil)
	subAVX2(nil, nil, nil)
	mulAVX2(nil, nil, nil)
	divAVX2(nil, nil, nil)
	addScalarAVX2(nil, 0, nil)
	mulScalarAVX2(nil, 0, nil)
	addAVX512(nil, nil, nil)
	subAVX512(nil, nil, nil)
	mulAVX512(nil, nil, nil)
	divAVX512(nil, nil, nil)
	addScalarAVX512(nil, 0, nil)
	mulScalarAVX512(nil, 0, nil)
	sqrtAVX2(nil, nil)
	sqrtAVX512(nil, nil)
	fmaAVX2(nil, nil, nil, nil)
	fmaAVX512(nil, nil, nil, nil)
	detectAMD64SIMD()
	
	// SVE stubs
	addSVE(nil, nil, nil)
	detectSVE()
	
	// SIMD wrappers
	small := []float64{1, 2, 3}
	TanhBatch(small)
	AsinhBatch(small)
	AcoshBatch(small)
	AtanhBatch(small)
	AbsScalar(1.0)
	NegScalar(1.0)
	SqrtScalar(1.0)
	FmaScalar(1, 2, 3)
	HasSVE()
	HasFMA()
	HasAVXVNNI()
}
