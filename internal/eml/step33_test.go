package eml

import (
	"math"
	"testing"
)



func TestStep33_SubnormalLogOracle(t *testing.T) {
	// 1. Verify powers of 2 in [2^-1074, 2^-1022] (all subnormal powers of 2)
	// plus 2^-1022 (smallest normal).
	var worstMathErr float64
	var worstNativeErr float64

	for k := -1074; k <= -1022; k++ {
		x := math.Ldexp(1.0, k)
		wantExact := float64(k) * math.Ln2
		oracleVal := ExactLogOracle(x)
		mathVal := math.Log(x)
		nativeVal := nativeLog(x)
		emlLogVal := Log(x)

		// Oracle matches mathematical truth k*ln2 within floating point precision
		if math.Abs(oracleVal-wantExact) > 1e-12 {
			t.Errorf("oracle(2^%d) = %v, want %v", k, oracleVal, wantExact)
		}

		// nativeLog must match exact oracle
		nativeErr := math.Abs(nativeVal - wantExact)
		if nativeErr > worstNativeErr {
			worstNativeErr = nativeErr
		}
		if nativeErr > 1e-12 {
			t.Errorf("nativeLog(2^%d) = %v, want %v (err %.3g)", k, nativeVal, wantExact, nativeErr)
		}

		// eml.Log must also match
		if math.Abs(emlLogVal-wantExact) > 1e-12 {
			t.Errorf("eml.Log(2^%d) = %v, want %v", k, emlLogVal, wantExact)
		}

		// nativeLog2(2^k) must match k
		wantLog2 := float64(k)
		if math.Abs(nativeLog2(x)-wantLog2) > 1e-12 {
			t.Errorf("nativeLog2(2^%d) = %v, want %v", k, nativeLog2(x), wantLog2)
		}

		// nativeLog10(2^k) must match k * log10(2)
		wantLog10 := float64(k) * math.Log10(2)
		if math.Abs(nativeLog10(x)-wantLog10) > 1e-12 {
			t.Errorf("nativeLog10(2^%d) = %v, want %v", k, nativeLog10(x), wantLog10)
		}

		mathErr := math.Abs(mathVal - wantExact)
		if mathErr > worstMathErr {
			worstMathErr = mathErr
		}
	}

	t.Logf("Subnormal sweep [-1074, -1022]: worst math.Log err = %.2f nats, worst nativeLog err = %.3g",
		worstMathErr, worstNativeErr)

	// Confirm that math.Log has the documented ~35.35 nat error at 2^-1074
	if worstMathErr < 30.0 {
		t.Logf("Notice: Go toolchain math.Log error is %.2f (if fixed upstream, update expectation)", worstMathErr)
	}

	// 2. Test smallest normal 2^-1022: math.Log and exact oracle agree
	smallestNormal := math.Ldexp(1.0, -1022)
	if math.Abs(math.Log(smallestNormal)-ExactLogOracle(smallestNormal)) > 1e-14 {
		t.Errorf("math.Log diverges on smallest normal: %v vs %v", math.Log(smallestNormal), ExactLogOracle(smallestNormal))
	}

	// 3. Dense non-power-of-2 subnormal verification
	for i := 1; i <= 20; i++ {
		mant := 1.0 + float64(i)*0.04
		subnorm := mant * math.Ldexp(1.0, -1050)
		want := ExactLogOracle(subnorm)
		if got := nativeLog(subnorm); math.Abs(got-want) > 1e-12 {
			t.Errorf("nativeLog(%g) = %v, want %v", subnorm, got, want)
		}
		if got := logScalar(subnorm); math.Abs(got-want) > 1e-12 {
			t.Errorf("logScalar(%g) = %v, want %v", subnorm, got, want)
		}
	}

	// 4. Batch LogSIMDTo verification on subnormals
	subSlice := make([]float64, 32)
	resSlice := make([]float64, 32)
	for i := range subSlice {
		subSlice[i] = math.Ldexp(1.0, -1074+i)
	}
	LogSIMDTo(subSlice, resSlice)
	for i, x := range subSlice {
		want := ExactLogOracle(x)
		if math.Abs(resSlice[i]-want) > 1e-12 {
			t.Errorf("LogSIMDTo[%d] = %v, want %v", i, resSlice[i], want)
		}
	}

	// 5. Domain boundary coverage
	for _, tc := range []struct {
		in   float64
		want float64
	}{
		{0, math.Inf(-1)},
		{math.Copysign(0, -1), math.Inf(-1)},
		{-1.0, math.NaN()},
		{math.Inf(1), math.Inf(1)},
		{math.NaN(), math.NaN()},
	} {
		gotOracle := ExactLogOracle(tc.in)
		gotScalar := logScalar(tc.in)
		if math.IsNaN(tc.want) {
			if !math.IsNaN(gotOracle) {
				t.Errorf("ExactLogOracle(%v) = %v, want NaN", tc.in, gotOracle)
			}
			if !math.IsNaN(gotScalar) {
				t.Errorf("logScalar(%v) = %v, want NaN", tc.in, gotScalar)
			}
		} else if math.IsInf(tc.want, 0) {
			if gotOracle != tc.want {
				t.Errorf("ExactLogOracle(%v) = %v, want %v", tc.in, gotOracle, tc.want)
			}
			if gotScalar != tc.want {
				t.Errorf("logScalar(%v) = %v, want %v", tc.in, gotScalar, tc.want)
			}
		}
	}
}
