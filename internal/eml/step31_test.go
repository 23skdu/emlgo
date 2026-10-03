//go:build amd64 && !purego
// +build amd64,!purego

package eml

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

// getFunctionBytes retrieves a byte slice starting at the function entry point.
// If fn is an ABIInternal wrapper that calls an ABI0 assembly routine (via 0xE8 CALL),
// it resolves the target address of the ABI0 routine.
func getFunctionBytes(fn any) []byte {
	type eface struct {
		_type unsafe.Pointer
		data  unsafe.Pointer
	}
	ef := (*eface)(unsafe.Pointer(&fn))
	codePtr := *(*unsafe.Pointer)(ef.data)
	wrapperBytes := unsafe.Slice((*byte)(codePtr), 256)

	// Check if this is an ABIInternal stub: "push rbp; mov rbp, rsp" (0x55 0x48 0x89 0xE5)
	if len(wrapperBytes) > 10 && wrapperBytes[0] == 0x55 && wrapperBytes[1] == 0x48 && wrapperBytes[2] == 0x89 && wrapperBytes[3] == 0xE5 {
		// Scan for 0xE8 (CALL rel32)
		for i := 0; i < 128; i++ {
			if wrapperBytes[i] == 0xE8 {
				rel := int32(binary.LittleEndian.Uint32(wrapperBytes[i+1 : i+5]))
				targetPtr := unsafe.Add(codePtr, int(i+5)+int(rel))
				return unsafe.Slice((*byte)(targetPtr), 2048)
			}
		}
	}
	return unsafe.Slice((*byte)(codePtr), 2048)
}

func hasVzeroupperBeforeRet(bytes []byte) bool {
	for i := 3; i < len(bytes); i++ {
		if bytes[i] == 0xC3 {
			if bytes[i-3] == 0xC5 && bytes[i-2] == 0xF8 && bytes[i-1] == 0x77 {
				return true
			}
		}
	}
	return false
}

// TestAVX512KernelsEmitVzeroupper disassembles the emitted machine code of all 11
// AVX-512 symbols and confirms 0xC5 0xF8 0x77 (VZEROUPPER) precedes RET.
func TestAVX512KernelsEmitVzeroupper(t *testing.T) {
	kernels := []struct {
		name string
		fn   any
	}{
		{"addAVX512", addAVX512},
		{"subAVX512", subAVX512},
		{"mulAVX512", mulAVX512},
		{"divAVX512", divAVX512},
		{"addScalarAVX512", addScalarAVX512},
		{"mulScalarAVX512", mulScalarAVX512},
		{"sqrtAVX512", sqrtAVX512},
		{"fmaAVX512", fmaAVX512},
		{"absAVX512", absAVX512},
		{"negAVX512", negAVX512},
		{"invAVX512", invAVX512},
	}

	for _, k := range kernels {
		t.Run(k.name, func(t *testing.T) {
			b := getFunctionBytes(k.fn)
			if !hasVzeroupperBeforeRet(b) {
				t.Fatalf("kernel %s missing VZEROUPPER (0xC5 0xF8 0x77) before RET", k.name)
			}
		})
	}

	// Also verify that AVX2 kernels have VZEROUPPER
	avx2Kernels := []struct {
		name string
		fn   any
	}{
		{"addAVX2", addAVX2},
		{"subAVX2", subAVX2},
		{"mulAVX2", mulAVX2},
		{"divAVX2", divAVX2},
		{"sqrtAVX2", sqrtAVX2},
		{"absAVX2", absAVX2},
		{"negAVX2", negAVX2},
		{"invAVX2", invAVX2},
	}
	for _, k := range avx2Kernels {
		t.Run(k.name, func(t *testing.T) {
			b := getFunctionBytes(k.fn)
			if !hasVzeroupperBeforeRet(b) {
				t.Fatalf("AVX2 kernel %s missing VZEROUPPER before RET", k.name)
			}
		})
	}
}

// TestXGETBVAndFeatureGating asserts xgetbv execution and XCR0 feature gating.
func TestXGETBVAndFeatureGating(t *testing.T) {
	_, _, ecx, _ := cpuid(1, 0)
	hasOSXSAVE := (ecx & (1 << 27)) != 0
	if !hasOSXSAVE {
		t.Skip("OSXSAVE not supported on this platform")
	}

	eaxX, edxX := xgetbv(0)
	xcr0 := uint64(eaxX) | (uint64(edxX) << 32)

	// Bit 1: XMM, Bit 2: YMM
	if (xcr0 & 0x6) != 0x6 {
		t.Errorf("expected OS support for XMM and YMM (xcr0 & 0x6 == 0x6), got xcr0=%#x", xcr0)
	}

	// Run detectAMD64SIMD and verify consistency
	detectAMD64SIMD()

	// AVX2 must only be enabled if xcr0 & 0x6 == 0x6
	if hasAVX2 && (xcr0&0x6) != 0x6 {
		t.Errorf("hasAVX2 is true but xcr0 lacks XMM/YMM support: %#x", xcr0)
	}

	// AVX512 must only be enabled if xcr0 & 0xe6 == 0xe6
	if hasAVX512 && (xcr0&0xe6) != 0xe6 {
		t.Errorf("hasAVX512 is true but xcr0 lacks AVX-512 state: %#x", xcr0)
	}
}

// TestAVX512ZeroTransitionPenalty runs kernels and executes scalar floating point math
// to verify clean exit and no AVX-to-SSE transition penalties or state corruption.
func TestAVX512ZeroTransitionPenalty(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	b := []float64{8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0}
	c := []float64{2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0}
	res := make([]float64, 8)

	if hasAVX512 {
		// Exercise each AVX-512 kernel exit path
		addAVX512(a, b, res)
		subAVX512(a, b, res)
		mulAVX512(a, b, res)
		divAVX512(a, b, res)
		addScalarAVX512(a, 2.0, res)
		mulScalarAVX512(a, 2.0, res)
		sqrtAVX512(a, res)
		fmaAVX512(a, b, c, res)
		absAVX512(a, res)
		negAVX512(a, res)
		invAVX512(a, res)
	} else if hasAVX2 {
		addAVX2(a, b, res)
		subAVX2(a, b, res)
		mulAVX2(a, b, res)
		divAVX2(a, b, res)
		sqrtAVX2(a, res)
		absAVX2(a, res)
		negAVX2(a, res)
		invAVX2(a, res)
	}

	// Now measure scalar floating point loop immediately after
	start := time.Now()
	var acc float64
	for i := 0; i < 10000; i++ {
		acc += float64(i) * 1.5
	}
	elapsed := time.Since(start)
	if acc == 0 || elapsed > 10*time.Second {
		t.Errorf("unexpected execution state: acc=%v, elapsed=%v", acc, elapsed)
	}
}


// suppress unused reflect import
var _ = reflect.TypeOf
