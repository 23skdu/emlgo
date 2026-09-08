package bigmath

import (
	"math/big"
	"math/rand"
)

// IdentityVerifier holds configuration for symbolic identity verification.
type IdentityVerifier struct {
	Precision uint    // big.Float precision in bits
	Samples   int     // number of random test points
	Tolerance float64 // max relative error allowed
}

// DefaultVerifier returns an IdentityVerifier with default settings.
func DefaultVerifier() *IdentityVerifier {
	return &IdentityVerifier{
		Precision: 256,
		Samples:   50,
		Tolerance: 4.0,
	}
}

// VerifyIdentity evaluates two expressions at random high-precision points
// and checks if they agree within the tolerance.
// expr1 and expr2 are functions that take *big.Float and return *big.Float.
func (v *IdentityVerifier) VerifyIdentity(expr1, expr2 func(*big.Float) *big.Float) bool {
	if v.Precision == 0 {
		v.Precision = Prec
	}
	if v.Samples == 0 {
		v.Samples = 50
	}
	if v.Tolerance == 0 {
		v.Tolerance = 2.0
	}

	rng := rand.New(rand.NewSource(42)) //#nosec G404 -- deterministic seed intentional: mathematical verifier, not a security-sensitive RNG
	for i := 0; i < v.Samples; i++ {
		// Generate random test point in (0.1, 10)
		x := 0.1 + rng.Float64()*9.9
		bx := new(big.Float).SetPrec(v.Precision).SetFloat64(x)

		r1 := expr1(bx)
		r2 := expr2(bx)

		diff := new(big.Float).Sub(r1, r2)
		diff.Abs(diff)

		// Check relative error
		absR1 := new(big.Float).Abs(r1)
		// Avoid division by zero: if both are near zero, use absolute tolerance
		if absR1.Sign() == 0 {
			absR2 := new(big.Float).Abs(r2)
			if absR2.Sign() == 0 {
				continue // both zero, matches
			}
			absR1 = absR2
		}
		relErr := new(big.Float).Quo(diff, absR1)

		// threshold = tolerance * 2^(-Precision)
		// This represents a few ULPs at the given precision.
		tolFloat := new(big.Float).SetPrec(v.Precision).SetFloat64(v.Tolerance)
		threshold := new(big.Float).SetPrec(v.Precision).SetMantExp(
			new(big.Float).SetPrec(v.Precision).SetInt64(1), -int(v.Precision))
		threshold.Mul(threshold, tolFloat)
		if relErr.Cmp(threshold) > 0 {
			return false
		}
	}
	return true
}

// EvalAt evaluates a function at a specific float64 value using high precision.
func EvalAt(f func(*big.Float) *big.Float, x float64) float64 {
	bx := new(big.Float).SetPrec(Prec).SetFloat64(x)
	return Float64(f(bx))
}
