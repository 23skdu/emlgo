package bigmath

import (
	"math/big"
	"testing"
)

func TestBigmathCoverageEdgeCases(t *testing.T) {
	// Sqrt with prec == 0
	zeroPrec := new(big.Float).SetFloat64(4.0)
	zeroPrec.SetPrec(0)
	_ = Sqrt(zeroPrec)

	// Tan with prec == 0
	_ = Tan(zeroPrec)

	_ = Tan(NewFloat(1.0))

	// Asin edge cases
	_ = Asin(zeroPrec)
	outOfRange := NewFloat(2.5)
	asinOut := Asin(outOfRange)
	if asinOut.Sign() != 0 {
		t.Errorf("Asin(2.5) should be 0, got %v", asinOut)
	}
	asin1 := Asin(NewFloat(1.0))
	if asin1.Sign() <= 0 {
		t.Errorf("Asin(1.0) should be positive")
	}
	asinNeg1 := Asin(NewFloat(-1.0))
	if asinNeg1.Sign() >= 0 {
		t.Errorf("Asin(-1.0) should be negative")
	}

	// Acos with prec == 0
	_ = Acos(zeroPrec)

	// VerifyIdentity edge cases
	vZero := &IdentityVerifier{Precision: 0, Samples: 0, Tolerance: 0}
	matched := vZero.VerifyIdentity(
		func(x *big.Float) *big.Float { return new(big.Float).SetInt64(0) },
		func(x *big.Float) *big.Float { return new(big.Float).SetInt64(0) },
	)
	if !matched {
		t.Errorf("VerifyIdentity with both zero should match")
	}

	// Failure case
	failed := vZero.VerifyIdentity(
		func(x *big.Float) *big.Float { return new(big.Float).SetInt64(0) },
		func(x *big.Float) *big.Float { return new(big.Float).SetInt64(100) },
	)
	if failed {
		t.Errorf("VerifyIdentity should fail for mismatched expressions")
	}
}
