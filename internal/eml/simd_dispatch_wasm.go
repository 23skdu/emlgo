//go:build wasm

package eml

func emlSIMD(x, y, result []float64) {
	scalarEml(x, y, result)
}

func detectPlatformSIMD() {
	detectWasmSIMD()
}

func scalarEml(x, y, result []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		result[i] = nativeExp(x[i]) - nativeLog(y[i])
	}
}

func fmaSIMD(a, b, c, result []float64) {
	fmaWasmSIMD(a, b, c, result)
}

func dispatchExpSIMDTo(x, result []float64)      { expWasmSIMD(x, result) }
func dispatchLogSIMDTo(x, result []float64)      { logWasmSIMD(x, result) }
func dispatchSinSIMDTo(x, result []float64)      { sinWasmSIMD(x, result) }
func dispatchCosSIMDTo(x, result []float64)      { cosWasmSIMD(x, result) }
func dispatchTanSIMDTo(x, result []float64)      { tanWasmSIMD(x, result) }
func dispatchSinCosSIMDTo(x, sin, cos []float64) { sincosWasmSIMD(x, sin, cos) }
func dispatchSqrtSIMDTo(x, result []float64)     { sqrtWasmSIMD(x, result) }
func dispatchAddSIMD(a, b, result []float64)     { addWasmSIMD(a, b, result) }
func dispatchSubSIMD(a, b, result []float64)     { subWasmSIMD(a, b, result) }
func dispatchMulSIMD(a, b, result []float64)     { mulWasmSIMD(a, b, result) }
func dispatchDivSIMD(a, b, result []float64)     { divWasmSIMD(a, b, result) }
func dispatchAbsSIMD(x, result []float64)        { absWasmSIMD(x, result) }
func dispatchNegSIMD(x, result []float64)        { negWasmSIMD(x, result) }
func dispatchInvSIMD(x, result []float64)        { invWasmSIMD(x, result) }

func dispatchAddScalarSIMD(a []float64, b float64, result []float64) {
	addScalarWasmSIMD(a, b, result)
}

func dispatchMulScalarSIMD(a []float64, b float64, result []float64) {
	mulScalarWasmSIMD(a, b, result)
}

var _ = parallelizeSinCos
