package bigmath

import (
	"math"
	"math/big"
	"testing"
)

func TestVerifyIdentityAtReportsReproducibleMismatch(t *testing.T) {
	verifier := DefaultVerifier()
	points := []*big.Float{
		new(big.Float).SetPrec(256).SetInt64(-1),
		new(big.Float).SetPrec(256).SetInt64(0),
		new(big.Float).SetPrec(256).SetFloat64(1e-300),
		new(big.Float).SetPrec(256).SetFloat64(1e300),
	}
	result, err := verifier.VerifyIdentityAt(
		func(x *big.Float) *big.Float { return new(big.Float).Copy(x) },
		func(x *big.Float) *big.Float {
			if x.Sign() == 0 {
				return new(big.Float).SetPrec(x.Prec()).SetInt64(1)
			}
			return new(big.Float).Copy(x)
		},
		points,
	)
	if err != nil {
		t.Fatalf("VerifyIdentityAt() error = %v", err)
	}
	if result.Passed || result.SamplesChecked != 2 || result.Failure == nil {
		t.Fatalf("result = %+v, want failure at sample 1", result)
	}
	if result.Failure.Sample != 1 || result.Failure.Input != "0" ||
		result.Failure.Left != "0" || result.Failure.Right != "1" ||
		result.Failure.RelativeError == "" || result.Failure.Threshold == "" {
		t.Errorf("failure report = %+v, want reproducible values at zero", result.Failure)
	}
}

func TestVerifyIdentityAtHandlesDomainAndInvalidResults(t *testing.T) {
	verifier := DefaultVerifier()
	points := []*big.Float{
		new(big.Float).SetPrec(256).SetInt64(-2),
		new(big.Float).SetPrec(256).SetInt64(0),
		new(big.Float).SetPrec(256).SetFloat64(1e-300),
		new(big.Float).SetPrec(256).SetFloat64(1e300),
	}
	identity := func(x *big.Float) *big.Float { return new(big.Float).Copy(x) }
	result, err := verifier.VerifyIdentityAt(identity, identity, points)
	if err != nil || !result.Passed || result.SamplesChecked != len(points) {
		t.Fatalf("VerifyIdentityAt() = (%+v, %v), want all points verified", result, err)
	}

	_, err = verifier.VerifyIdentityAt(
		func(*big.Float) *big.Float { return nil },
		identity,
		points[:1],
	)
	if err == nil {
		t.Error("VerifyIdentityAt() accepted a nil expression result")
	}
	_, err = verifier.VerifyIdentityAt(identity, identity, []*big.Float{nil})
	if err == nil {
		t.Error("VerifyIdentityAt() accepted a nil input point")
	}
	_, err = verifier.VerifyIdentityAt(
		func(*big.Float) *big.Float { return new(big.Float).SetInf(false) },
		identity,
		points[:1],
	)
	if err == nil {
		t.Error("VerifyIdentityAt() accepted an infinite expression result")
	}
	verifier.Tolerance = math.NaN()
	if _, err := verifier.VerifyIdentityAt(identity, identity, points[:1]); err == nil {
		t.Error("VerifyIdentityAt() accepted a NaN tolerance")
	}
}

func TestVerifyIdentityDoesNotMutateConfiguration(t *testing.T) {
	verifier := &IdentityVerifier{}
	if !verifier.VerifyIdentity(
		func(x *big.Float) *big.Float { return new(big.Float).Copy(x) },
		func(x *big.Float) *big.Float { return new(big.Float).Copy(x) },
	) {
		t.Fatal("VerifyIdentity() rejected equal expressions")
	}
	if verifier.Precision != 0 || verifier.Samples != 0 || verifier.Tolerance != 0 {
		t.Errorf("VerifyIdentity() mutated configuration: %+v", verifier)
	}
}

func TestVerifyIdentity_SinCosSquared(t *testing.T) {
	v := DefaultVerifier()

	// sin²(x) + cos²(x) = 1
	expr1 := func(x *big.Float) *big.Float {
		s := Sin(x)
		c := Cos(x)
		s2 := new(big.Float).Mul(s, s)
		c2 := new(big.Float).Mul(c, c)
		return new(big.Float).Add(s2, c2)
	}
	expr2 := func(x *big.Float) *big.Float {
		return new(big.Float).SetPrec(v.Precision).SetInt64(1)
	}

	if !v.VerifyIdentity(expr1, expr2) {
		t.Error("Failed to verify sin²(x) + cos²(x) = 1")
	}
}

func TestVerifyIdentity_ExpLog(t *testing.T) {
	v := DefaultVerifier()

	// exp(log(x)) = x
	expr1 := func(x *big.Float) *big.Float {
		return Exp(Log(x))
	}
	expr2 := func(x *big.Float) *big.Float {
		return new(big.Float).Copy(x)
	}

	if !v.VerifyIdentity(expr1, expr2) {
		t.Error("Failed to verify exp(log(x)) = x")
	}
}

func TestVerifyIdentity_LogExp(t *testing.T) {
	v := DefaultVerifier()

	// log(exp(x)) = x
	expr1 := func(x *big.Float) *big.Float {
		return Log(Exp(x))
	}
	expr2 := func(x *big.Float) *big.Float {
		return new(big.Float).Copy(x)
	}

	if !v.VerifyIdentity(expr1, expr2) {
		t.Error("Failed to verify log(exp(x)) = x")
	}
}

func TestVerifyIdentity_NonIdentity(t *testing.T) {
	v := DefaultVerifier()

	// sin(x) ≠ cos(x) — should return false
	expr1 := func(x *big.Float) *big.Float {
		return Sin(x)
	}
	expr2 := func(x *big.Float) *big.Float {
		return Cos(x)
	}

	if v.VerifyIdentity(expr1, expr2) {
		t.Error("Verifier incorrectly reported sin(x) = cos(x)")
	}
}

func TestDefaultVerifier(t *testing.T) {
	v := DefaultVerifier()
	if v.Precision != 256 {
		t.Errorf("Default precision = %d, want 256", v.Precision)
	}
	if v.Samples != 50 {
		t.Errorf("Default samples = %d, want 50", v.Samples)
	}
	if v.Tolerance != 4.0 {
		t.Errorf("Default tolerance = %f, want 4.0", v.Tolerance)
	}
}

func TestEvalAt(t *testing.T) {
	result := EvalAt(func(x *big.Float) *big.Float {
		return Exp(x)
	}, 0)
	if result != 1.0 {
		t.Errorf("EvalAt(exp, 0) = %v, want 1", result)
	}
}
