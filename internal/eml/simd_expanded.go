package eml

func Log2SIMDTo(x, result []float64) {
	if len(x) != len(result) {
		panic("length mismatch")
	}
	if len(x) < SmallCutoff {
		for i := range x {
			result[i] = nativeLog2(x[i])
		}
		return
	}
	parallelizeGeneric(x, result, nativeLog2)
}

func Log10SIMDTo(x, result []float64) {
	if len(x) != len(result) {
		panic("length mismatch")
	}
	if len(x) < SmallCutoff {
		for i := range x {
			result[i] = nativeLog10(x[i])
		}
		return
	}
	parallelizeGeneric(x, result, nativeLog10)
}
