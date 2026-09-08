package bigmath

import (
	"math/big"
	"testing"
)

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
