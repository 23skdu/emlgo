package bigmath

import (
	"math"
	"testing"
)

func TestStep35_NewFloatFromStringChecked(t *testing.T) {
	// Valid inputs
	validCases := []struct {
		input string
		want  float64
	}{
		{"1.5", 1.5},
		{"-3.1415", -3.1415},
		{"0", 0.0},
		{"1e-5", 1e-5},
		{"2.5e+3", 2500.0},
		{"0x10", 16.0},
		{"0b101", 5.0},
		{"0o77", 63.0},
	}

	for _, tc := range validCases {
		f, err := NewFloatFromStringChecked(tc.input)
		if err != nil {
			t.Errorf("NewFloatFromStringChecked(%q) unexpected error: %v", tc.input, err)
			continue
		}
		got, _ := f.Float64()
		if math.Abs(got-tc.want) > 1e-12 {
			t.Errorf("NewFloatFromStringChecked(%q) = %v, want %v", tc.input, got, tc.want)
		}

		// Also check deprecated NewFloatFromString
		fOld := NewFloatFromString(tc.input)
		if fOld == nil {
			t.Errorf("NewFloatFromString(%q) returned nil, want valid *big.Float", tc.input)
		} else {
			gotOld, _ := fOld.Float64()
			if math.Abs(gotOld-tc.want) > 1e-12 {
				t.Errorf("NewFloatFromString(%q) = %v, want %v", tc.input, gotOld, tc.want)
			}
		}
	}

	// Invalid inputs: must return error and NewFloatFromString must return nil without panicking
	invalidInputs := []string{
		"",
		"abc",
		"1.5x",
		"0xZZ",
		"++1",
		"--1",
		".",
		"e10",
	}

	for _, input := range invalidInputs {
		f, err := NewFloatFromStringChecked(input)
		if err == nil {
			t.Errorf("NewFloatFromStringChecked(%q) expected error, got %v", input, f)
		}

		// NewFloatFromString must not panic and must return nil
		fOld := NewFloatFromString(input)
		if fOld != nil {
			t.Errorf("NewFloatFromString(%q) = %v, want nil", input, fOld)
		}
	}
}

func TestStep35_LogDomainHandling(t *testing.T) {
	// Log(0) -> -Inf
	zero := NewFloat(0.0)
	resZero := Log(zero)
	if !resZero.IsInf() || resZero.Sign() >= 0 {
		t.Errorf("Log(0) = %v, want -Inf", resZero)
	}

	// Log(-1) -> 0 (sentinel)
	neg := NewFloat(-1.0)
	resNeg := Log(neg)
	if resNeg.Sign() != 0 || resNeg.IsInf() {
		t.Errorf("Log(-1) = %v, want 0", resNeg)
	}

	// Log(-100) -> 0 (sentinel)
	neg100 := NewFloat(-100.0)
	resNeg100 := Log(neg100)
	if resNeg100.Sign() != 0 {
		t.Errorf("Log(-100) = %v, want 0", resNeg100)
	}

	// Log(1) -> 0
	oneF := NewFloat(1.0)
	resOne := Log(oneF)
	if resOne.Sign() != 0 {
		t.Errorf("Log(1) = %v, want 0", resOne)
	}

	// Log(e) -> 1
	eVal := Exp(NewFloat(1.0))
	resE := Log(eVal)
	fVal, _ := resE.Float64()
	if math.Abs(fVal-1.0) > 1e-15 {
		t.Errorf("Log(e) = %v, want 1.0", fVal)
	}
}

func TestStep35_AllocationProfile(t *testing.T) {
	// Warm up caches
	_ = piConst(Prec + 128)
	_ = ln2Const(Prec + 64)

	x := NewFloat(0.5)

	// Measure Log allocs
	logAllocs := testing.AllocsPerRun(10, func() {
		_ = Log(x)
	})
	if logAllocs > 25 {
		t.Errorf("Log allocs/op = %v, want <= 25 (previously 322)", logAllocs)
	}

	// Measure Exp allocs (previously 303, now ~187)
	expAllocs := testing.AllocsPerRun(10, func() {
		_ = Exp(x)
	})
	if expAllocs > 200 {
		t.Errorf("Exp allocs/op = %v, want <= 200 (previously 303)", expAllocs)
	}

	// Measure Sin allocs (previously 244, now ~138)
	sinAllocs := testing.AllocsPerRun(10, func() {
		_ = Sin(x)
	})
	if sinAllocs > 150 {
		t.Errorf("Sin allocs/op = %v, want <= 150 (previously 244)", sinAllocs)
	}

	// Measure Cos allocs
	cosAllocs := testing.AllocsPerRun(10, func() {
		_ = Cos(x)
	})
	if cosAllocs > 150 {
		t.Errorf("Cos allocs/op = %v, want <= 150", cosAllocs)
	}

	// Measure Atan allocs
	atanAllocs := testing.AllocsPerRun(10, func() {
		_ = Atan(x)
	})
	if atanAllocs > 400 {
		t.Errorf("Atan allocs/op = %v, want <= 400", atanAllocs)
	}

	// Measure Acos allocs (previously 744, now ~450)
	acosAllocs := testing.AllocsPerRun(10, func() {
		_ = Acos(x)
	})
	if acosAllocs > 480 {
		t.Errorf("Acos allocs/op = %v, want <= 480 (previously 744)", acosAllocs)
	}
}

func TestStep35_PiAndLn2Caching(t *testing.T) {
	p1 := piConst(Prec)
	p2 := piConst(Prec)
	if p1.Cmp(p2) != 0 {
		t.Errorf("piConst mismatch: %v vs %v", p1, p2)
	}

	l1 := ln2Const(Prec)
	l2 := ln2Const(Prec)
	if l1.Cmp(l2) != 0 {
		t.Errorf("ln2Const mismatch: %v vs %v", l1, l2)
	}
}

func BenchmarkStep35_Log(b *testing.B) {
	x := NewFloat(2.5)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = Log(x)
	}
}

func BenchmarkStep35_Exp(b *testing.B) {
	x := NewFloat(1.5)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = Exp(x)
	}
}

func BenchmarkStep35_Sin(b *testing.B) {
	x := NewFloat(1.5)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = Sin(x)
	}
}

func BenchmarkStep35_Acos(b *testing.B) {
	x := NewFloat(0.5)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = Acos(x)
	}
}

func BenchmarkStep35_Tan(b *testing.B) {
	x := NewFloat(0.5)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = Tan(x)
	}
}
