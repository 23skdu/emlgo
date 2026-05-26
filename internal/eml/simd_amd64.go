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

func fmaAVX2(a, b, c, result []float64)
func fmaAVX512(a, b, c, result []float64)

func sqrtScalar(x float64) float64
func fmaScalar(a, b, c float64) float64

func negScalar(x float64) float64
func absScalar(x float64) float64

func detectAMD64SIMD() {
	_, _, ecx, _ := cpuid(1, 0)
	hasSSE4 = (ecx & (1 << 19)) != 0
	hasFMA = (ecx & (1 << 12)) != 0

	_, ebx, _, _ := cpuid(7, 0)
	hasAVX2 = (ebx & (1 << 5)) != 0
	hasAVX512 = (ebx & (1 << 16)) != 0

	// CPUID.7.1:EAX
	eax71, _, _, _ := cpuid(7, 1)
	hasAVXVNNI = (eax71 & (1 << 4)) != 0
}
