//go:build amd64
// +build amd64

package eml

func cpuid(op, op2 uint32) (eax, ebx, ecx, edx uint32)

func addAVX2(a, b, result []float64)
func subAVX2(a, b, result []float64)
func mulAVX2(a, b, result []float64)
func divAVX2(a, b, result []float64)

func addScalarAVX2(a []float64, b float64, result []float64)
func mulScalarAVX2(a []float64, b float64, result []float64)

func addAVX512(a, b, result []float64)
func subAVX512(a, b, result []float64)
func mulAVX512(a, b, result []float64)
func divAVX512(a, b, result []float64)

func addScalarAVX512(a []float64, b float64, result []float64)
func mulScalarAVX512(a []float64, b float64, result []float64)

func sqrtAVX2(a, result []float64)
func sqrtAVX512(a, result []float64)

func sqrtScalar(x float64) float64
func fmaScalar(a, b, c float64) float64

func negScalar(x float64) float64 {
	if x == 0 {
		return copysign(0, -1)
	}
	return -x
}

func absScalar(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func detectAMD64SIMD() {
	_, _, ecx, _ := cpuid(1, 0)
	hasSSE4 = (ecx & (1 << 19)) != 0

	_, ebx, _, _ := cpuid(7, 0)
	hasAVX2 = (ebx & (1 << 5)) != 0
	hasAVX512 = (ebx & (1 << 16)) != 0
}
