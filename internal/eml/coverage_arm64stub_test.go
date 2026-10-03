//go:build (!arm64 && !wasm) || purego

package eml

import "testing"

// TestArm64DispatchStubs exercises the non-arm64 dispatch stubs so that the
// fallback path stays covered on amd64 / non-arm64 hosts.
func TestArm64DispatchStubs(t *testing.T) {
	addNEON(nil, nil, nil)
	subNEON(nil, nil, nil)
	mulNEON(nil, nil, nil)
	divNEON(nil, nil, nil)
	addScalarNEON(nil, 0, nil)
	mulScalarNEON(nil, 0, nil)
	sqrtNEON(nil, nil)
	absNEON(nil, nil)
	negNEON(nil, nil)
	invNEON(nil, nil)
	fmaNEON(nil, nil, nil, nil)
	addSatInt8NEON(nil, nil, nil)
	subSatInt8NEON(nil, nil, nil)
	addSVE(nil, nil, nil)
	detectARM64SIMD()
}
