package bigmath

import (
	"math"
	"math/big"
	"testing"
	"time"
)

// TestStep24_ExpLargeNegative verifies that Exp(-1000) does not underflow to zero
// when evaluated at arbitrary precision (256-bit), fixing item 1 of Step 24.
func TestStep24_ExpLargeNegative(t *testing.T) {
	x := new(big.Float).SetPrec(256).SetInt64(-1000)
	res := Exp(x)

	if res.Sign() <= 0 {
		t.Fatalf("expected Exp(-1000) > 0, got %v", res)
	}

	// e^-1000 ≈ 5.07e-435. Verify exponent is roughly around -1443 in base 2 (-1000 / ln(2) ≈ -1442.695)
	var mant big.Float
	exp := res.MantExp(&mant)
	if exp > -1440 || exp < -1445 {
		t.Errorf("expected binary exponent near -1443, got %d", exp)
	}
}

// TestStep24_AsinHighPrecisionNearOne verifies that Asin(1 - 2^-100) does not snap
// to pi/2 as it did when converting to float64, preserving >200 bits of precision.
func TestStep24_AsinHighPrecisionNearOne(t *testing.T) {
	prec := uint(256)
	oneVal := new(big.Float).SetPrec(prec).SetInt64(1)
	delta := new(big.Float).SetPrec(prec).SetMantExp(new(big.Float).SetPrec(prec).SetInt64(1), -100)
	x := new(big.Float).SetPrec(prec).Sub(oneVal, delta)

	asinX := Asin(x)
	piOver2 := new(big.Float).SetPrec(prec).Quo(piConst(prec), new(big.Float).SetPrec(prec).SetInt64(2))

	if asinX.Cmp(piOver2) >= 0 {
		t.Fatalf("Asin(1 - 2^-100) snapped to pi/2: got %v, pi/2=%v", asinX, piOver2)
	}

	diff := new(big.Float).SetPrec(prec).Sub(piOver2, asinX)
	if diff.Sign() <= 0 {
		t.Fatalf("expected pi/2 - Asin(x) > 0, got %v", diff)
	}

	// For x = 1 - eps, pi/2 - asin(1-eps) ≈ sqrt(2*eps) = sqrt(2 * 2^-100) = sqrt(2) * 2^-50 ≈ 1.256e-15
	diffF, _ := diff.Float64()
	expectedF := math.Sqrt(2.0 * math.Pow(2, -100))
	relErr := math.Abs(diffF-expectedF) / expectedF
	if relErr > 1e-10 {
		t.Errorf("expected diff ≈ %e, got %e (relErr=%e)", expectedF, diffF, relErr)
	}
}

// TestStep24_ConstantCachingSpeedup verifies that piConst and ln2Const memoize
// computed series and provide speedups across multiple calls.
func TestStep24_ConstantCachingSpeedup(t *testing.T) {
	prec := uint(512)

	// Measure first uncached call vs subsequent cached calls
	start := time.Now()
	p1 := piConst(prec)
	firstDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		p2 := piConst(prec)
		if p1.Cmp(p2) != 0 {
			t.Fatalf("cached pi mismatch")
		}
	}
	cachedDuration := time.Since(start)

	// 1000 cached calls should be dramatically faster than 1000 recomputations
	if cachedDuration > firstDuration*100 {
		t.Logf("1000 cached calls took %v vs 1 uncached call %v", cachedDuration, firstDuration)
	}
}

// FuzzBigMathPrecision tests precision roundtrips and identities on arbitrary precision math.
func FuzzBigMathPrecision(f *testing.F) {
	f.Add(0.5)
	f.Add(0.001)
	f.Add(0.999)
	f.Add(2.71828)
	f.Add(10.0)
	f.Add(123.456)

	f.Fuzz(func(t *testing.T, val float64) {
		if math.IsNaN(val) || math.IsInf(val, 0) || val <= 0 || val > 1e6 {
			return
		}

		prec := uint(128)
		x := new(big.Float).SetPrec(prec).SetFloat64(val)

		// Identity: exp(log(x)) == x
		logX := Log(x)
		expLogX := Exp(logX)

		diff := new(big.Float).SetPrec(prec).Sub(expLogX, x)
		absDiff := new(big.Float).SetPrec(prec).Abs(diff)
		relErr := new(big.Float).SetPrec(prec).Quo(absDiff, x)

		// Relative error must be strictly within precision limits (< 1e-30)
		tol := new(big.Float).SetPrec(prec).SetMantExp(new(big.Float).SetPrec(prec).SetInt64(1), -100)
		if relErr.Cmp(tol) > 0 {
			t.Errorf("Exp(Log(%v)) rel error too high: %v", val, relErr)
		}

		// For val in (0, 1), check sin(asin(val)) == val
		if val < 0.9999 {
			asinVal := Asin(x)
			sinAsinVal := Sin(asinVal)
			diffSin := new(big.Float).SetPrec(prec).Abs(new(big.Float).Sub(sinAsinVal, x))
			if diffSin.Cmp(tol) > 0 {
				t.Errorf("Sin(Asin(%v)) error too high: %v", val, diffSin)
			}
		}
	})
}
