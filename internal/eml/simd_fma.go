package eml

// FmaSIMD computes a*b + c for each element and returns the result.
func FmaSIMD(a, b, c []float64) []float64 {
	if len(a) != len(b) || len(a) != len(c) {
		panic("length mismatch")
	}
	n := len(a)
	if n == 0 {
		return a
	}
	result := make([]float64, n)
	FmaSIMDTo(a, b, c, result)
	return result
}

// FmaSIMDTo computes a*b + c for each element and stores it in result.
func FmaSIMDTo(a, b, c, result []float64) {
	n := len(a)
	if n != len(b) || n != len(c) || n != len(result) {
		panic("length mismatch")
	}
	if n == 0 {
		return
	}

	if n < SmallCutoff {
		for i := 0; i < n; i++ {
			result[i] = a[i]*b[i] + c[i]
		}
		return
	}

	fmaSIMD(a, b, c, result)
}
