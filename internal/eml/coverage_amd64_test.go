//go:build amd64 && !purego
// +build amd64,!purego

package eml

import (
	"math"
	"testing"
)

func TestAMD64DispatchBranches(t *testing.T) {
	old512 := hasAVX512
	old512VNNI := hasAVX512VNNI
	old2 := hasAVX2
	oldFMA := hasFMA
	defer func() {
		hasAVX512 = old512
		hasAVX512VNNI = old512VNNI
		hasAVX2 = old2
		hasFMA = oldFMA
	}()

	_ = HasAVX512VNNI()

	const n = 7
	a := []float64{1, math.NaN(), 3, 4, 5, 6, math.NaN()}
	b := []float64{2, 1, 4, 3, 6, 5, 8}
	c := []float64{1, 1, 1, 1, 1, 1, 1}
	res := make([]float64, n)

	const nf32 = 9
	af32 := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	bf32 := []float32{2, 3, 4, 5, 6, 7, 8, 9, 10}
	cf32 := []float32{1, 1, 1, 1, 1, 1, 1, 1, 1}
	resf32 := make([]float32, nf32)

	u8 := []uint8{1, 2, 3, 4, 5, 6, 7}
	i8 := []int8{1, -2, 3, -4, 5, -6, 7}

	// Mode 1: AVX-512 enabled (simdLen=0 for n=7 / nf32=9, tests AVX-512 branch and tail scalar loop)
	hasAVX512 = true
	hasAVX512VNNI = true
	hasAVX2 = true
	hasFMA = true

	amd64AddSIMD(a, b, res)
	amd64SubSIMD(a, b, res)
	amd64MulSIMD(a, b, res)
	amd64DivSIMD(a, b, res)
	amd64AddScalarSIMD(a, 2.0, res)
	amd64MulScalarSIMD(a, 2.0, res)
	amd64SqrtSIMD(a, res)
	amd64AbsSIMD(a, res)
	amd64NegSIMD(a, res)
	amd64InvSIMD(a, res)
	amd64ExpSIMD(a, res)
	amd64LogSIMD(a, res)
	amd64SinSIMD(a, res)
	amd64CosSIMD(a, res)
	amd64TanSIMD(a, res)
	fmaSIMD(a, b, c, res)
	dispatchFloorSIMDTo(a, res)
	dispatchCeilSIMDTo(a, res)
	dispatchTruncSIMDTo(a, res)
	dispatchRoundSIMDTo(a, res)
	dispatchMaxSIMDTo(a, b, res)
	dispatchMinSIMDTo(a, b, res)

	dispatchAddSIMDF32To(af32, bf32, resf32)
	dispatchSubSIMDF32To(af32, bf32, resf32)
	dispatchMulSIMDF32To(af32, bf32, resf32)
	dispatchDivSIMDF32To(af32, bf32, resf32)
	dispatchSqrtSIMDF32To(af32, resf32)
	dispatchAddScalarSIMDF32To(af32, 2.0, resf32)
	dispatchMulScalarSIMDF32To(af32, 2.0, resf32)
	dispatchAbsSIMDF32To(af32, resf32)
	dispatchNegSIMDF32To(af32, resf32)
	dispatchInvSIMDF32To(af32, resf32)
	dispatchFmaSIMDF32To(af32, bf32, cf32, resf32)
	_ = dispatchDotInt8SIMD(u8, i8)

	// Mode 2: AVX-512 disabled, AVX2 enabled (simdLen=4 for n=7, simdLen=8 for nf32=9)
	hasAVX512 = false
	hasAVX512VNNI = false
	hasAVX2 = true
	hasFMA = true

	amd64AddSIMD(a, b, res)
	amd64SubSIMD(a, b, res)
	amd64MulSIMD(a, b, res)
	amd64DivSIMD(a, b, res)
	amd64AddScalarSIMD(a, 2.0, res)
	amd64MulScalarSIMD(a, 2.0, res)
	amd64SqrtSIMD(a, res)
	amd64AbsSIMD(a, res)
	amd64NegSIMD(a, res)
	amd64InvSIMD(a, res)
	amd64ExpSIMD(a, res)
	amd64LogSIMD(a, res)
	amd64SinSIMD(a, res)
	amd64CosSIMD(a, res)
	amd64TanSIMD(a, res)
	fmaSIMD(a, b, c, res)
	dispatchFloorSIMDTo(a, res)
	dispatchCeilSIMDTo(a, res)
	dispatchTruncSIMDTo(a, res)
	dispatchRoundSIMDTo(a, res)
	dispatchMaxSIMDTo(a, b, res)
	dispatchMinSIMDTo(a, b, res)

	dispatchAddSIMDF32To(af32, bf32, resf32)
	dispatchSubSIMDF32To(af32, bf32, resf32)
	dispatchMulSIMDF32To(af32, bf32, resf32)
	dispatchDivSIMDF32To(af32, bf32, resf32)
	dispatchSqrtSIMDF32To(af32, resf32)
	dispatchAddScalarSIMDF32To(af32, 2.0, resf32)
	dispatchMulScalarSIMDF32To(af32, 2.0, resf32)
	dispatchAbsSIMDF32To(af32, resf32)
	dispatchNegSIMDF32To(af32, resf32)
	dispatchInvSIMDF32To(af32, resf32)
	dispatchFmaSIMDF32To(af32, bf32, cf32, resf32)

	// Mode 3: AVX-512 disabled, AVX2 disabled (tests fallback scalar loop)
	hasAVX512 = false
	hasAVX2 = false
	hasFMA = false

	amd64AddSIMD(a, b, res)
	amd64SubSIMD(a, b, res)
	amd64MulSIMD(a, b, res)
	amd64DivSIMD(a, b, res)
	amd64AddScalarSIMD(a, 2.0, res)
	amd64MulScalarSIMD(a, 2.0, res)
	amd64SqrtSIMD(a, res)
	amd64AbsSIMD(a, res)
	amd64NegSIMD(a, res)
	amd64InvSIMD(a, res)
	amd64ExpSIMD(a, res)
	amd64LogSIMD(a, res)
	amd64SinSIMD(a, res)
	amd64CosSIMD(a, res)
	amd64TanSIMD(a, res)
	fmaSIMD(a, b, c, res)
	dispatchFloorSIMDTo(a, res)
	dispatchCeilSIMDTo(a, res)
	dispatchTruncSIMDTo(a, res)
	dispatchRoundSIMDTo(a, res)
	dispatchMaxSIMDTo(a, b, res)
	dispatchMinSIMDTo(a, b, res)

	dispatchAddSIMDF32To(af32, bf32, resf32)
	dispatchSubSIMDF32To(af32, bf32, resf32)
	dispatchMulSIMDF32To(af32, bf32, resf32)
	dispatchDivSIMDF32To(af32, bf32, resf32)
	dispatchSqrtSIMDF32To(af32, resf32)
	dispatchAddScalarSIMDF32To(af32, 2.0, resf32)
	dispatchMulScalarSIMDF32To(af32, 2.0, resf32)
	dispatchAbsSIMDF32To(af32, resf32)
	dispatchNegSIMDF32To(af32, resf32)
	dispatchInvSIMDF32To(af32, resf32)
	dispatchFmaSIMDF32To(af32, bf32, cf32, resf32)
}

func TestInt8SatSIMD(t *testing.T) {
	old512 := hasAVX512
	old2 := hasAVX2
	defer func() {
		hasAVX512 = old512
		hasAVX2 = old2
	}()

	const n = 35
	a := make([]int8, n)
	b := make([]int8, n)
	res := make([]int8, n)

	// Tail clamping cases for add: 100+50>127, -100-50<-128
	a[32] = 100
	b[32] = 50
	a[33] = -100
	b[33] = -50
	a[34] = 10
	b[34] = 20

	// Tail clamping cases for sub: 100-(-50)>127, -100-50<-128
	subA := make([]int8, n)
	subB := make([]int8, n)
	subA[32] = 100
	subB[32] = -50
	subA[33] = -100
	subB[33] = 50
	subA[34] = 20
	subB[34] = 10

	// With AVX512 (simdLen=0 for n=35, tests AVX512 tail clamping)
	hasAVX512 = true
	hasAVX2 = true
	dispatchAddSatInt8SIMD(a, b, res)
	dispatchSubSatInt8SIMD(subA, subB, res)

	// With AVX2
	hasAVX512 = false
	hasAVX2 = true
	dispatchAddSatInt8SIMD(a, b, res)
	dispatchSubSatInt8SIMD(subA, subB, res)
	AddSatInt8SIMDTo(a, b, res)
	SubSatInt8SIMDTo(subA, subB, res)

	// Without AVX2 (fallback scalar loop)
	hasAVX2 = false
	a5 := a[30:]
	b5 := b[30:]
	res5 := res[30:]
	dispatchAddSatInt8SIMD(a5, b5, res5)
	subA5 := subA[30:]
	subB5 := subB[30:]
	dispatchSubSatInt8SIMD(subA5, subB5, res5)

	// Panics on length mismatch
	func() {
		defer func() { _ = recover() }()
		AddSatInt8SIMDTo(a, b[:10], res)
	}()
	func() {
		defer func() { _ = recover() }()
		SubSatInt8SIMDTo(a, b[:10], res)
	}()
	func() {
		defer func() { _ = recover() }()
		_ = DotInt8SIMD(make([]uint8, 2), make([]int8, 3))
	}()
	func() {
		defer func() { _ = recover() }()
		SqrtSIMDF32To(make([]float32, 2), make([]float32, 3))
	}()
	func() {
		defer func() { _ = recover() }()
		FloorSIMDTo(make([]float64, 2), make([]float64, 3))
	}()
	func() {
		defer func() { _ = recover() }()
		CeilSIMDTo(make([]float64, 2), make([]float64, 3))
	}()
	func() {
		defer func() { _ = recover() }()
		TruncSIMDTo(make([]float64, 2), make([]float64, 3))
	}()
	func() {
		defer func() { _ = recover() }()
		RoundSIMDTo(make([]float64, 2), make([]float64, 3))
	}()
	func() {
		defer func() { _ = recover() }()
		MaxSIMDTo(make([]float64, 2), make([]float64, 3), make([]float64, 2))
	}()
	func() {
		defer func() { _ = recover() }()
		MinSIMDTo(make([]float64, 2), make([]float64, 3), make([]float64, 2))
	}()

	// Multi-chunk ForEachChunk branch for Floor/Ceil/Trunc/Round/Max/Min (n >= SmallCutoff)
	largeN := SmallCutoff() * 2
	la := make([]float64, largeN)
	lb := make([]float64, largeN)
	lres := make([]float64, largeN)
	FloorSIMDTo(la, lres)
	CeilSIMDTo(la, lres)
	TruncSIMDTo(la, lres)
	RoundSIMDTo(la, lres)
	MaxSIMDTo(la, lb, lres)
	MinSIMDTo(la, lb, lres)
}
