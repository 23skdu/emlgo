package bigmath

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

// IdentityVerificationResult describes a successful verification or its first
// reproducible mismatch.
type IdentityVerificationResult struct {
	Passed         bool
	SamplesChecked int
	Failure        *IdentityVerificationFailure
}

// IdentityVerificationFailure records the sample and values that exceeded the
// configured relative-error threshold.
type IdentityVerificationFailure struct {
	Sample        int
	Input         string
	Left          string
	Right         string
	RelativeError string
	Threshold     string
}

// IdentityVerifier holds configuration for symbolic identity verification.
type IdentityVerifier struct {
	Precision uint    // big.Float precision in bits
	Samples   int     // number of random test points
	Tolerance float64 // multiplier for the 2^-Precision relative-error threshold
}

// DefaultVerifier returns an IdentityVerifier with default settings.
func DefaultVerifier() *IdentityVerifier {
	return &IdentityVerifier{
		Precision: 256,
		Samples:   50,
		Tolerance: 4.0,
	}
}

// VerifyIdentity evaluates two expressions at reproducible sample points
// converted to the configured precision and checks their relative error.
// expr1 and expr2 are functions that take *big.Float and return *big.Float.
func (v *IdentityVerifier) VerifyIdentity(expr1, expr2 func(*big.Float) *big.Float) bool {
	precision, samples, _ := v.settings()
	if samples < 0 || precision > uint(math.MaxInt) {
		return false
	}
	rng := rand.New(rand.NewSource(42)) //#nosec G404 -- deterministic seed intentional: mathematical verifier, not a security-sensitive RNG
	for i := 0; i < samples; i++ {
		x := 0.1 + rng.Float64()*9.9
		point := new(big.Float).SetPrec(precision).SetFloat64(x)
		result, err := v.VerifyIdentityAt(expr1, expr2, []*big.Float{point})
		if err != nil || !result.Passed {
			return false
		}
	}
	return true
}

// VerifyIdentityAt checks two expressions at the supplied deterministic
// points. Callers should choose points inside both expressions' domains.
// A mismatch is returned in the result; invalid inputs or non-finite results
// are returned as errors.
func (v *IdentityVerifier) VerifyIdentityAt(
	expr1, expr2 func(*big.Float) *big.Float,
	points []*big.Float,
) (IdentityVerificationResult, error) {
	precision, _, tolerance := v.settings()
	if expr1 == nil || expr2 == nil {
		return IdentityVerificationResult{}, fmt.Errorf("identity expressions must not be nil")
	}
	if len(points) == 0 {
		return IdentityVerificationResult{}, fmt.Errorf("identity verification requires at least one point")
	}
	if precision > uint(math.MaxInt) {
		return IdentityVerificationResult{}, fmt.Errorf("precision %d exceeds supported exponent range", precision)
	}
	if math.IsNaN(tolerance) || math.IsInf(tolerance, 0) || tolerance < 0 {
		return IdentityVerificationResult{}, fmt.Errorf("tolerance must be finite and non-negative")
	}
	tolFloat := new(big.Float).SetPrec(precision).SetFloat64(tolerance)
	threshold := new(big.Float).SetPrec(precision).SetMantExp(
		new(big.Float).SetPrec(precision).SetInt64(1), -int(precision))
	threshold.Mul(threshold, tolFloat)

	result := IdentityVerificationResult{Passed: true}
	for i, point := range points {
		if point == nil || point.IsInf() {
			return IdentityVerificationResult{}, fmt.Errorf("verification point %d must be finite and non-nil", i)
		}
		input := new(big.Float).SetPrec(precision).Copy(point)
		left := expr1(new(big.Float).SetPrec(precision).Copy(input))
		right := expr2(new(big.Float).SetPrec(precision).Copy(input))
		if left == nil || right == nil || left.IsInf() || right.IsInf() {
			return IdentityVerificationResult{}, fmt.Errorf("expression returned nil or non-finite result at sample %d", i)
		}
		left = new(big.Float).SetPrec(precision).Copy(left)
		right = new(big.Float).SetPrec(precision).Copy(right)
		result.SamplesChecked++

		diff := new(big.Float).SetPrec(precision).Sub(left, right)
		diff.Abs(diff)
		denominator := new(big.Float).SetPrec(precision).Abs(left)
		if denominator.Sign() == 0 {
			denominator.Abs(right)
		}
		if denominator.Sign() == 0 {
			continue
		}
		relativeError := new(big.Float).SetPrec(precision).Quo(diff, denominator)
		if relativeError.Cmp(threshold) > 0 {
			result.Passed = false
			result.Failure = &IdentityVerificationFailure{
				Sample:        i,
				Input:         input.Text('g', -1),
				Left:          left.Text('g', -1),
				Right:         right.Text('g', -1),
				RelativeError: relativeError.Text('g', -1),
				Threshold:     threshold.Text('g', -1),
			}
			return result, nil
		}
	}
	return result, nil
}

func (v *IdentityVerifier) settings() (uint, int, float64) {
	precision, samples, tolerance := v.Precision, v.Samples, v.Tolerance
	if precision == 0 {
		precision = Prec
	}
	if samples == 0 {
		samples = 50
	}
	if tolerance == 0 {
		tolerance = 2.0
	}
	return precision, samples, tolerance
}

// EvalAt evaluates a function at a specific float64 value using high precision.
func EvalAt(f func(*big.Float) *big.Float, x float64) float64 {
	bx := new(big.Float).SetPrec(Prec).SetFloat64(x)
	return Float64(f(bx))
}
