package eml

import (
	"math"
	"testing"
)

func TestNativeMath(t *testing.T) {
	// Test nan and inf
	n := nan()
	if !math.IsNaN(n) {
		t.Error("nan() did not return NaN")
	}

	i1 := inf(1)
	if !math.IsInf(i1, 1) {
		t.Error("inf(1) did not return +Inf")
	}
	i2 := inf(-1)
	if !math.IsInf(i2, -1) {
		t.Error("inf(-1) did not return -Inf")
	}

	// Test checks
	if !isNaN(n) {
		t.Error("isNaN failed for NaN")
	}
	if isNaN(0.0) {
		t.Error("isNaN failed for 0.0")
	}
	if !isInf(i1, 1) {
		t.Error("isInf failed for +Inf")
	}
	if !isInf(i2, -1) {
		t.Error("isInf failed for -Inf")
	}

	// Test bits
	if f64bits(1.0) != math.Float64bits(1.0) {
		t.Error("f64bits failed")
	}
	if f64frombits(math.Float64bits(1.0)) != 1.0 {
		t.Error("f64frombits failed")
	}

	// Test arithmetic
	if nativeNeg(1.0) != -1.0 {
		t.Error("nativeNeg failed")
	}
	if nativeAbs(-1.0) != 1.0 {
		t.Error("nativeAbs failed")
	}
	if nativeInv(2.0) != 0.5 {
		t.Error("nativeInv failed")
	}
	if nativeMax(1.0, 2.0) != 2.0 {
		t.Error("nativeMax failed")
	}
	if nativeMin(1.0, 2.0) != 1.0 {
		t.Error("nativeMin failed")
	}

	// Test rounding
	if floor(1.9) != 1.0 {
		t.Error("floor failed")
	}
	if ceil(1.1) != 2.0 {
		t.Error("ceil failed")
	}
	if trunc(1.9) != 1.0 {
		t.Error("trunc failed")
	}
	if round(1.5) != 2.0 {
		t.Error("round failed")
	}

	// Test exponential/log
	if math.Abs(nativeExp(1.0)-math.Exp(1.0)) > 1e-15 {
		t.Error("nativeExp failed")
	}
	if math.Abs(nativeLog(math.E)-1.0) > 1e-15 {
		t.Error("nativeLog failed")
	}
	if math.Abs(nativeLog10(100.0)-2.0) > 1e-15 {
		t.Error("nativeLog10 failed")
	}
	if math.Abs(nativeSqrt(4.0)-2.0) > 1e-15 {
		t.Error("nativeSqrt failed")
	}

	// Test trig
	if math.Abs(nativeSin(math.Pi/2)-1.0) > 1e-15 {
		t.Error("nativeSin failed")
	}
	if math.Abs(nativeCos(math.Pi)-(-1.0)) > 1e-15 {
		t.Error("nativeCos failed")
	}
	s, c := nativeSincos(0)
	if s != 0 || c != 1 {
		t.Error("nativeSincos failed")
	}

	// Test hyper
	if math.Abs(nativeAsinh(0.0)-0.0) > 1e-15 {
		t.Error("nativeAsinh failed")
	}
	if math.Abs(nativeAcosh(1.0)-0.0) > 1e-15 {
		t.Error("nativeAcosh failed")
	}
	if math.Abs(nativeAtanh(0.0)-0.0) > 1e-15 {
		t.Error("nativeAtanh failed")
	}

	// Test Misc
	if nativeCbrt(8.0) != 2.0 {
		t.Error("nativeCbrt failed")
	}
	if nativeHypot(3.0, 4.0) != 5.0 {
		t.Error("nativeHypot failed")
	}
	if math.Abs(nativePow(2.0, 3.0)-8.0) > 1e-13 {
		t.Errorf("nativePow(2, 3) = %v, want 8.0", nativePow(2.0, 3.0))
	}
}

func TestNativeEdge(t *testing.T) {
	// nativeSqrt negative
	if !math.IsNaN(nativeSqrt(-1.0)) {
		t.Error("nativeSqrt(-1) should be NaN")
	}
	// nativeSqrt zero
	if nativeSqrt(0) != 0 {
		t.Error("nativeSqrt(0) should be 0")
	}

	// nativeExp large
	if !math.IsInf(nativeExp(1000), 1) {
		t.Error("nativeExp(1000) should be Inf")
	}

	// nativeMod
	if nativeMod(10, 3) != 1 {
		t.Errorf("nativeMod(10, 3) = %v, want 1", nativeMod(10, 3))
	}

	// nativeRemainder
	if nativeRemainder(10, 3) != 1 {
		t.Errorf("nativeRemainder(10, 3) = %v, want 1", nativeRemainder(10, 3))
	}

	// nativePow branches
	if nativePow(1.0, 5.0) != 1.0 { t.Error("nativePow(1, 5) failed") }
	if nativePow(2.0, 0.0) != 1.0 { t.Error("nativePow(2, 0) failed") }
	if nativePow(0.0, 2.0) != 0.0 { t.Error("nativePow(0, 2) failed") }
	if !math.IsInf(nativePow(0.0, -1.0), 1) { t.Error("nativePow(0, -1) failed") }
	if !math.IsNaN(nativePow(-2.0, 0.5)) { t.Error("nativePow(-2, 0.5) should be NaN") }
	if nativePow(-2.0, 2.0) != 4.0 { t.Errorf("nativePow(-2, 2) = %v", nativePow(-2, 2)) }
	if math.Abs(nativePow(-2.0, 3.0)-(-8.0)) > 1e-13 { t.Errorf("nativePow(-2, 3) = %v", nativePow(-2, 3)) }

	// nativeLog1p and nativeExpm1
	if nativeLog1p(0) != 0 { t.Error("nativeLog1p(0) failed") }
	if nativeExpm1(0) != 0 { t.Error("nativeExpm1(0) failed") }

	// nativeMax/Min with NaNs (our implementation returns the non-NaN value)
	if nativeMax(math.NaN(), 1.0) != 1.0 { t.Error("nativeMax with NaN failed") }
	if nativeMin(math.NaN(), 1.0) != 1.0 { t.Error("nativeMin with NaN failed") }
	
	// nativeMax/Min with Infs
	if nativeMax(math.Inf(1), 1.0) != math.Inf(1) { t.Error("nativeMax with Inf failed") }
	if nativeMin(math.Inf(-1), 1.0) != math.Inf(-1) { t.Error("nativeMin with -Inf failed") }

	// nativeAbs/Neg edge cases
	if nativeAbs(math.Inf(-1)) != math.Inf(1) { t.Error("nativeAbs(-Inf) failed") }
	if nativeNeg(math.Inf(1)) != math.Inf(-1) { t.Error("nativeNeg(Inf) failed") }

	// nativeInv edge cases
	if !math.IsInf(nativeInv(0), 1) { t.Error("nativeInv(0) failed") }

	// nativeMod/Remainder edge cases
	if !math.IsNaN(nativeMod(10, 0)) { t.Error("nativeMod(10, 0) should be NaN") }
	if !math.IsNaN(nativeRemainder(10, 0)) { t.Error("nativeRemainder(10, 0) should be NaN") }

	// nativeHypot edge cases
	if math.IsInf(nativeHypot(math.Inf(1), 0), 1) == false { t.Error("nativeHypot(Inf, 0) failed") }

	// nativeCbrt edge cases
	if nativeCbrt(0) != 0 { t.Error("nativeCbrt(0) failed") }
	if nativeCbrt(-8) != -2 { t.Error("nativeCbrt(-8) failed") }
}
