package fastmath

// FastEml computes the canonical EML operator: eml(x, y) = exp(x) - ln(y)
// using high-speed minimax polynomial approximations for both exp and log.
func FastEml(x, y float64) float64 {
	return FastExp(x) - FastLog(y)
}

// FastEmlF32 computes eml(x, y) for float32.
func FastEmlF32(x, y float32) float32 {
	return FastExpF32(x) - FastLogF32(y)
}

// FastEmlBitCast computes an ultra-fast approximation of eml(x, y)
// using Schraudolph's bitcast exponential and bit-cast Chebyshev logarithm.
func FastEmlBitCast(x, y float64) float64 {
	return FastExpBitCast(x) - FastLogMinimax(y)
}

// FastEmlClipped computes eml(x, y) with domain clipping: exp_clamped(x) - ln(max(y, eps)).
// Prevents NaN cascading and Inf overflows.
func FastEmlClipped(x, y, eps float64) float64 {
	return ExpClamped(x) - LnClipped(y, eps)
}

// FastEmlClippedF32 computes clipped eml(x, y) for float32.
func FastEmlClippedF32(x, y, eps float32) float32 {
	return ExpClampedF32(x) - LnClippedF32(y, eps)
}

// FastEmlRegularized computes smooth regularized eml(x, y):
// exp_clamped(x) - 0.5 * ln(y^2 + eps^2).
// Infinitely differentiable (C-infinity) across the entire real plane R^2.
func FastEmlRegularized(x, y, eps float64) float64 {
	return ExpClamped(x) - LnRegularized(y, eps)
}

// FastEmlRegularizedF32 computes smooth regularized eml(x, y) for float32.
func FastEmlRegularizedF32(x, y, eps float32) float32 {
	return ExpClampedF32(x) - LnRegularizedF32(y, eps)
}

// FastEmlBatch computes element-wise FastEml for float64 slices.
func FastEmlBatch(x, y []float64) []float64 {
	if len(x) != len(y) {
		panic("slice length mismatch")
	}
	dst := make([]float64, len(x))
	FastEmlBatchTo(x, y, dst)
	return dst
}

// FastEmlBatchTo computes element-wise FastEml storing results in dst.
func FastEmlBatchTo(x, y, dst []float64) {
	if len(x) != len(y) || len(x) != len(dst) {
		panic("slice length mismatch")
	}
	for i := range x {
		dst[i] = FastEml(x[i], y[i])
	}
}

// FastEmlBatchF32 computes element-wise FastEml for float32 slices.
func FastEmlBatchF32(x, y []float32) []float32 {
	if len(x) != len(y) {
		panic("slice length mismatch")
	}
	dst := make([]float32, len(x))
	FastEmlBatchToF32(x, y, dst)
	return dst
}

// FastEmlBatchToF32 computes element-wise FastEml for float32 slices storing in dst.
func FastEmlBatchToF32(x, y, dst []float32) {
	if len(x) != len(y) || len(x) != len(dst) {
		panic("slice length mismatch")
	}
	for i := range x {
		dst[i] = FastEmlF32(x[i], y[i])
	}
}
