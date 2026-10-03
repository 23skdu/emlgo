//go:build amd64 && !purego
// +build amd64,!purego

package eml

import (
	"math"
	"testing"
)

// TestTranscendentalAsmKernelsAreWiredUp checks that the experimental
// transcendental AVX2 kernels in simd_amd64.s are still declared and callable.
//
// They are not on the default dispatch path -- see simd_trans_amd64.go for the
// measurements that moved them behind `-tags emlasm` -- but go vet's asmdecl
// check requires the Go declarations to exist, and referencing them here keeps
// the linker honest about the symbols still existing.
func TestTranscendentalAsmKernelsAreWiredUp(t *testing.T) {
	// Empty slices: every kernel reduces its length by the vector width and
	// returns immediately, so this exercises the entry/exit paths only.
	expAVX2(nil, nil)
	logAVX2(nil, nil)
	sinAVX2(nil, nil)
	cosAVX2(nil, nil)
	tanAVX2(nil, nil)

	// These kernels process exactly one vector per iteration and have no
	// scalar tail: the caller is responsible for the remainder. That is why
	// every dispatcher computes simdLen := (n/4)*4 and finishes in Go.
	in := make([]float64, 4)
	for i := range in {
		in[i] = 0.5 + float64(i)*0.1
	}
	out := make([]float64, 4)

	for name, fn := range map[string]func(a, result []float64){
		"exp": expAVX2,
		"log": logAVX2,
		"sin": sinAVX2,
		"cos": cosAVX2,
		"tan": tanAVX2,
	} {
		for i := range out {
			out[i] = math.NaN()
		}
		fn(in, out)
		for i := range out {
			if math.IsNaN(out[i]) {
				t.Errorf("%sAVX2 left out[%d] unwritten for a full vector", name, i)
			}
		}
	}
}
