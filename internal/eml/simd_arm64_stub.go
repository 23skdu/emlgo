//go:build (!arm64 && !wasm) || purego
// +build !arm64,!wasm purego

package eml

func addNEON(a, b, result []float64)                         { _, _, _ = a, b, result }
func subNEON(a, b, result []float64)                         { _, _, _ = a, b, result }
func mulNEON(a, b, result []float64)                         { _, _, _ = a, b, result }
func divNEON(a, b, result []float64)                         { _, _, _ = a, b, result }
func addScalarNEON(a []float64, b float64, result []float64) { _, _, _ = a, b, result }
func mulScalarNEON(a []float64, b float64, result []float64) { _, _, _ = a, b, result }
func sqrtNEON(a, result []float64)                           { _, _ = a, result }
func absNEON(a, result []float64)                            { _, _ = a, result }
func negNEON(a, result []float64)                            { _, _ = a, result }
func invNEON(a, result []float64)                            { _, _ = a, result }
func fmaNEON(a, b, c, result []float64)                      { _, _, _, _ = a, b, c, result }
func addSatInt8NEON(a, b, result []int8)                     { _, _, _ = a, b, result }
func subSatInt8NEON(a, b, result []int8)                     { _, _, _ = a, b, result }

func addSVE(a, b, result []float64) {
	_, _, _ = a, b, result
}

func detectARM64SIMD() {
	_ = 0
}

