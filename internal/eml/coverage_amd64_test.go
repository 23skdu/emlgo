//go:build amd64 && !purego
// +build amd64,!purego

package eml

import (
	"testing"
)

func TestAMD64DispatchBranches(t *testing.T) {
	old512 := hasAVX512
	old2 := hasAVX2
	oldFMA := hasFMA
	defer func() {
		hasAVX512 = old512
		hasAVX2 = old2
		hasFMA = oldFMA
	}()

	const n = 7
	a := []float64{1, 2, 3, 4, 5, 6, 7}
	b := []float64{2, 3, 4, 5, 6, 7, 8}
	c := []float64{1, 1, 1, 1, 1, 1, 1}
	res := make([]float64, n)

	// Mode 1: AVX-512 enabled (simdLen=0 for n=7, tests AVX-512 branch and tail scalar loop)
	hasAVX512 = true
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
	fmaSIMD(a, b, c, res)

	// Mode 2: AVX-512 disabled, AVX2 enabled (simdLen=4 for n=7, tests AVX2 branch and tail scalar loop)
	hasAVX512 = false
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
	fmaSIMD(a, b, c, res)

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
	fmaSIMD(a, b, c, res)
}

func TestInt8SatSIMD(t *testing.T) {
	old2 := hasAVX2
	defer func() { hasAVX2 = old2 }()

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

	// With AVX2
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
}
