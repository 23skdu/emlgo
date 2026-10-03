//go:build (!amd64 && !wasm) || purego

package eml

import "testing"

// TestAmd64DispatchStubs exercises the non-amd64 dispatch stubs so that the
// fallback path stays covered. It is excluded on js/wasm, where the amd64
// dispatch symbols are not declared.
func TestAmd64DispatchStubs(t *testing.T) {
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
	addSatInt8AVX2(nil, nil, nil)
	subSatInt8AVX2(nil, nil, nil)
	absAVX2(nil, nil)
	negAVX2(nil, nil)
	invAVX2(nil, nil)
	absAVX512(nil, nil)
	negAVX512(nil, nil)
	invAVX512(nil, nil)
	expAVX2(nil, nil)
	logAVX2(nil, nil)
	sinAVX2(nil, nil)
	cosAVX2(nil, nil)
	tanAVX2(nil, nil)
	_, _, _, _ = cpuid(0, 0)
	detectAMD64SIMD()
}
