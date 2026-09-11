package fastmath

// GuardrailMode defines the numeric regularization strategy applied to prevent
// NaN/Inf cascading during symbolic search, numeric fitting, or genetic programming.
type GuardrailMode uint8

const (
	// GuardrailStrict adheres strictly to standard domain rules (non-positive log returns NaN/Inf).
	GuardrailStrict GuardrailMode = iota

	// GuardrailClip clamps values at domain boundaries: ln(max(y, eps)).
	GuardrailClip

	// GuardrailSmooth applies infinitely differentiable (C-infinity) regularization:
	// 0.5 * ln(y^2 + eps^2). Ensures continuous gradients and zero singularities at y=0.
	GuardrailSmooth
)

const (
	// DefaultEpsilon is the default regularization epsilon for float64.
	DefaultEpsilon = 1e-12

	// DefaultEpsilonF32 is the default regularization epsilon for float32.
	DefaultEpsilonF32 = float32(1e-6)

	// MaxExpThreshold is the upper clamping bound for Exp to prevent IEEE +Inf overflow.
	MaxExpThreshold = 700.0

	// MinExpThreshold is the lower clamping bound for Exp to prevent subnormal denormal stalls.
	MinExpThreshold = -700.0
)

// ExpClamped computes e^x with input clamped to [-700.0, 700.0] to prevent +Inf overflow.
func ExpClamped(x float64) float64 {
	if x > MaxExpThreshold {
		x = MaxExpThreshold
	} else if x < MinExpThreshold {
		return 0.0
	}
	return FastExp(x)
}

// ExpClampedF32 computes e^x for float32 clamped to [-87.0, 87.0].
func ExpClampedF32(x float32) float32 {
	if x > 87.0 {
		x = 87.0
	} else if x < -87.0 {
		return 0.0
	}
	return FastExpF32(x)
}

// LnClipped computes ln(max(y, eps)) ensuring non-negative arguments and preventing NaN/Inf.
func LnClipped(y, eps float64) float64 {
	if eps <= 0 {
		eps = DefaultEpsilon
	}
	if y < eps {
		y = eps
	}
	return FastLog(y)
}

// LnClippedF32 computes ln(max(y, eps)) for float32.
func LnClippedF32(y, eps float32) float32 {
	if eps <= 0 {
		eps = DefaultEpsilonF32
	}
	if y < eps {
		y = eps
	}
	return FastLogF32(y)
}

// LnRegularized computes the smooth, infinitely differentiable regularization:
// 0.5 * ln(y^2 + eps^2).
// This function is C-infinity everywhere on R, eliminates singularities at y=0,
// and has continuous derivative: d/dy [0.5 * ln(y^2 + eps^2)] = y / (y^2 + eps^2).
func LnRegularized(y, eps float64) float64 {
	if eps <= 0 {
		eps = DefaultEpsilon
	}
	return 0.5 * FastLog(y*y+eps*eps)
}

// LnRegularizedF32 computes 0.5 * ln(y^2 + eps^2) for float32.
func LnRegularizedF32(y, eps float32) float32 {
	if eps <= 0 {
		eps = DefaultEpsilonF32
	}
	return 0.5 * FastLogF32(y*y+eps*eps)
}
