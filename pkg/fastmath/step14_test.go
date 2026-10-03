package fastmath

import (
	"math"
	"testing"
)

func TestStep14FastExpRanges(t *testing.T) {
	// Subnormal and deep underflow checks
	tests := []struct {
		x       float64
		wantPos bool // true if math.Exp(x) is finite and > 0
	}{
		{-700.0, true},
		{-708.5, true},
		{-709.0, true},
		{-720.0, true},
		{-744.0, true},
		{-746.0, false},
	}

	for _, tt := range tests {
		got := FastExpMinimax(tt.x)
		want := math.Exp(tt.x)
		if tt.wantPos {
			if got <= 0 || math.IsInf(got, 0) || math.IsNaN(got) {
				t.Errorf("FastExpMinimax(%v) = %v, want positive finite (math.Exp=%v)", tt.x, got, want)
			}
			// Check relative magnitude (within order of magnitude in subnormal territory)
			ratio := got / want
			if ratio < 0.5 || ratio > 2.0 {
				t.Errorf("FastExpMinimax(%v) = %v, math.Exp = %v (ratio=%v)", tt.x, got, want, ratio)
			}
		} else {
			if got != 0.0 {
				t.Errorf("FastExpMinimax(%v) = %v, want 0.0", tt.x, got)
			}
		}
	}

	// Overflow checks
	if !math.IsInf(FastExpMinimax(711.0), 1) {
		t.Errorf("FastExpMinimax(711.0) = %v, want +Inf", FastExpMinimax(711.0))
	}
	if math.IsInf(FastExpMinimax(709.0), 1) {
		t.Errorf("FastExpMinimax(709.0) should be finite")
	}
	gotNearMax := FastExpMinimax(709.7)
	if math.IsInf(gotNearMax, 1) || gotNearMax <= 0 {
		t.Errorf("FastExpMinimax(709.7) = %v, want finite positive", gotNearMax)
	}
}

func TestStep14FastExpF32Ranges(t *testing.T) {
	// 88.000015 and 88.7 must be finite
	got88 := FastExpF32(88.000015)
	if math.IsInf(float64(got88), 1) || got88 <= 0 {
		t.Errorf("FastExpF32(88.000015) = %v, want finite positive", got88)
	}

	got88_7 := FastExpF32(88.7)
	if math.IsInf(float64(got88_7), 1) || got88_7 <= 0 {
		t.Errorf("FastExpF32(88.7) = %v, want finite positive", got88_7)
	}

	// > 88.722839 must be +Inf
	if !math.IsInf(float64(FastExpF32(89.0)), 1) {
		t.Errorf("FastExpF32(89.0) = %v, want +Inf", FastExpF32(89.0))
	}

	// -87.5 and -103.0 must be finite and > 0
	gotNeg87 := FastExpF32(-87.5)
	if gotNeg87 <= 0 || math.IsNaN(float64(gotNeg87)) {
		t.Errorf("FastExpF32(-87.5) = %v, want positive finite", gotNeg87)
	}

	gotNeg103 := FastExpF32(-103.0)
	if gotNeg103 <= 0 || math.IsNaN(float64(gotNeg103)) {
		t.Errorf("FastExpF32(-103.0) = %v, want positive subnormal", gotNeg103)
	}

	// < -103.27893 must be 0
	if FastExpF32(-104.0) != 0.0 {
		t.Errorf("FastExpF32(-104.0) = %v, want 0.0", FastExpF32(-104.0))
	}

	// NaN and overflow check
	if !math.IsNaN(float64(FastExpF32(float32(math.NaN())))) {
		t.Errorf("FastExpF32(NaN) should be NaN")
	}
}

func TestStep14LnRegularizedNoOverflow(t *testing.T) {
	// Ordinary large inputs must not overflow to +Inf
	got := LnRegularized(1e200, 1e-12)
	want := 460.51701859880914 // ln(1e200)
	if math.IsInf(got, 0) || math.Abs(got-want) > 1e-4 {
		t.Errorf("LnRegularized(1e200, 1e-12) = %v, want ~%v", got, want)
	}

	gotF32 := LnRegularizedF32(1e20, 1e-6)
	wantF32 := float32(46.05170186) // ln(1e20)
	if math.IsInf(float64(gotF32), 0) || math.Abs(float64(gotF32-wantF32)) > 1e-2 {
		t.Errorf("LnRegularizedF32(1e20, 1e-6) = %v, want ~%v", gotF32, wantF32)
	}

	// At 0
	gotZero := LnRegularized(0, 1e-6)
	wantZero := math.Log(1e-6)
	if math.Abs(gotZero-wantZero) > 1e-4 {
		t.Errorf("LnRegularized(0, 1e-6) = %v, want %v", gotZero, wantZero)
	}
}

func TestStep14FastExpDenseSweep(t *testing.T) {
	// Sweep from -700 to +700 in steps of 10.0
	for x := -700.0; x <= 700.0; x += 10.0 {
		got := FastExp(x)
		want := math.Exp(x)
		relErr := math.Abs(got-want) / want
		if relErr > 1e-5 {
			t.Errorf("FastExp(%v) = %v, want %v, relErr = %v", x, got, want, relErr)
		}
	}
}
