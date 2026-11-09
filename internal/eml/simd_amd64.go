//go:build amd64 && !purego
// +build amd64,!purego

package eml

import _ "unsafe"

func cpuid(op, op2 uint32) (eax, ebx, ecx, edx uint32)
func xgetbv(cx uint32) (eax, edx uint32)

//go:noescape
func addAVX2(a, b, result []float64)

//go:noescape
func subAVX2(a, b, result []float64)

//go:noescape
func mulAVX2(a, b, result []float64)

//go:noescape
func divAVX2(a, b, result []float64)

//go:noescape
func addScalarAVX2(a []float64, b float64, result []float64)

//go:noescape
func mulScalarAVX2(a []float64, b float64, result []float64)

//go:noescape
func addAVX512(a, b, result []float64)

//go:noescape
func subAVX512(a, b, result []float64)

//go:noescape
func mulAVX512(a, b, result []float64)

//go:noescape
func divAVX512(a, b, result []float64)

//go:noescape
func addScalarAVX512(a []float64, b float64, result []float64)

//go:noescape
func mulScalarAVX512(a []float64, b float64, result []float64)

//go:noescape
func sqrtAVX2(a, result []float64)

//go:noescape
func sqrtAVX512(a, result []float64)

//go:noescape
func absAVX2(a, result []float64)

//go:noescape
func negAVX2(a, result []float64)

//go:noescape
func invAVX2(a, result []float64)

//go:noescape
func absAVX512(a, result []float64)

//go:noescape
func negAVX512(a, result []float64)

//go:noescape
func invAVX512(a, result []float64)

//go:noescape
func expAVX2(a, result []float64)

//go:noescape
func logAVX2(a, result []float64)

//go:noescape
func sinAVX2(a, result []float64)

//go:noescape
func cosAVX2(a, result []float64)

//go:noescape
func tanAVX2(a, result []float64)

//go:noescape
func sincosAVX2(a, sinOut, cosOut []float64)

//go:noescape
func expAVX512(a, result []float64)

//go:noescape
func logAVX512(a, result []float64)

//go:noescape
func sinAVX512(a, result []float64)

//go:noescape
func cosAVX512(a, result []float64)

//go:noescape
func tanAVX512(a, result []float64)

//go:noescape
func fmaAVX2(a, b, c, result []float64)

//go:noescape
func fmaAVX512(a, b, c, result []float64)

//go:noescape
func addAVX2F32(a, b, result []float32)

//go:noescape
func subAVX2F32(a, b, result []float32)

//go:noescape
func mulAVX2F32(a, b, result []float32)

//go:noescape
func divAVX2F32(a, b, result []float32)

//go:noescape
func sqrtAVX2F32(a, result []float32)

//go:noescape
func addScalarAVX2F32(a []float32, b float32, result []float32)

//go:noescape
func mulScalarAVX2F32(a []float32, b float32, result []float32)

//go:noescape
func absAVX2F32(a, result []float32)

//go:noescape
func negAVX2F32(a, result []float32)

//go:noescape
func invAVX2F32(a, result []float32)

//go:noescape
func fmaAVX2F32(a, b, c, result []float32)

//go:noescape
func addAVX512F32(a, b, result []float32)

//go:noescape
func subAVX512F32(a, b, result []float32)

//go:noescape
func mulAVX512F32(a, b, result []float32)

//go:noescape
func divAVX512F32(a, b, result []float32)

//go:noescape
func sqrtAVX512F32(a, result []float32)

//go:noescape
func addScalarAVX512F32(a []float32, b float32, result []float32)

//go:noescape
func mulScalarAVX512F32(a []float32, b float32, result []float32)

//go:noescape
func absAVX512F32(a, result []float32)

//go:noescape
func negAVX512F32(a, result []float32)

//go:noescape
func invAVX512F32(a, result []float32)

//go:noescape
func fmaAVX512F32(a, b, c, result []float32)

//go:noescape
func floorAVX2(a, result []float64)

//go:noescape
func ceilAVX2(a, result []float64)

//go:noescape
func truncAVX2(a, result []float64)

//go:noescape
func roundAVX2(a, result []float64)

//go:noescape
func maxAVX2(a, b, result []float64)

//go:noescape
func minAVX2(a, b, result []float64)

//go:noescape
func floorAVX512(a, result []float64)

//go:noescape
func ceilAVX512(a, result []float64)

//go:noescape
func truncAVX512(a, result []float64)

//go:noescape
func roundAVX512(a, result []float64)

//go:noescape
func maxAVX512(a, b, result []float64)

//go:noescape
func minAVX512(a, b, result []float64)

//go:noescape
func addSatInt8AVX2(a, b, result []int8)

//go:noescape
func subSatInt8AVX2(a, b, result []int8)

//go:noescape
func dotInt8AVX2(a []uint8, b []int8) int32

//go:noescape
func addSatInt8AVX512(a, b, result []int8)

//go:noescape
func subSatInt8AVX512(a, b, result []int8)

//go:noescape
func dotInt8AVX512VNNI(a []uint8, b []int8) int32

func sqrtScalar(x float64) float64
func fmaScalar(a, b, c float64) float64

func negScalar(x float64) float64
func absScalar(x float64) float64

func detectAMD64SIMD() {
	_, _, ecx, _ := cpuid(1, 0)
	hasSSE4 = (ecx & (1 << 19)) != 0

	hasOSXSAVE := (ecx & (1 << 27)) != 0
	var osAVX, osAVX512 bool
	if hasOSXSAVE {
		eaxX, edxX := xgetbv(0)
		xcr0 := uint64(eaxX) | (uint64(edxX) << 32)
		osAVX = (xcr0 & 0x6) == 0x6
		osAVX512 = (xcr0 & 0xe6) == 0xe6
	}

	hasFMA = (ecx&(1<<12)) != 0 && osAVX

	_, ebx, ecx7, _ := cpuid(7, 0)
	hasAVX2 = (ebx&(1<<5)) != 0 && osAVX
	hasAVX512 = (ebx&(1<<16)) != 0 && osAVX512
	hasAVX512VNNI = (ecx7&(1<<11)) != 0 && osAVX512

	// CPUID.7.1:EAX
	eax71, _, _, _ := cpuid(7, 1)
	hasAVXVNNI = ((eax71&(1<<4)) != 0 && osAVX) || hasAVX512VNNI
}
