package main

import (
	"math"
	"testing"
)

func TestValidation(t *testing.T) {
	if !RunValidation() {
		t.Fatal("RunValidation failed")
	}

	validateInjectFail = true
	defer func() { validateInjectFail = false }()
	if RunValidation() {
		t.Fatal("RunValidation should fail when validateInjectFail is true")
	}
}

func TestValidationFailureBranches(t *testing.T) {
	oldChecker := sliceChecker
	defer func() { sliceChecker = oldChecker }()
	sliceChecker = func(name string, got, expected []float64, tol float64) bool { return false }

	oldSizes := validateSizes
	defer func() { validateSizes = oldSizes }()
	validateSizes = []int{2}

	if RunValidation() {
		t.Fatal("expected failure when sliceChecker returns false")
	}
}

func TestWithinTol(t *testing.T) {
	cases := []struct {
		a, b, tol float64
		want      bool
	}{
		{1.0, 1.0, 1e-6, true},
		{0, 0, 1e-6, true},
		{math.NaN(), math.NaN(), 1e-6, true},
		{math.Inf(1), math.Inf(1), 1e-6, true},
		{math.Inf(-1), math.Inf(-1), 1e-6, true},
		{1.0, 2.0, 1e-6, false},
		{math.NaN(), 1.0, 1e-6, false},
		{1.0, math.NaN(), 1e-6, false},
		{math.Inf(1), 1.0, 1e-6, false},
		{math.Inf(-1), 1.0, 1e-6, false},
	}
	for _, c := range cases {
		if got := withinTol(c.a, c.b, c.tol); got != c.want {
			t.Errorf("withinTol(%v, %v, %v) = %v, want %v", c.a, c.b, c.tol, got, c.want)
		}
	}
}

func TestCheckSlice(t *testing.T) {
	got := []float64{1.0, 2.0}
	expGood := []float64{1.0, 2.0}
	expBad := []float64{1.0, 3.0}

	if !checkSlice("good", got, expGood, 1e-6) {
		t.Errorf("checkSlice good failed")
	}
	if checkSlice("bad", got, expBad, 1e-6) {
		t.Errorf("checkSlice bad should have failed")
	}
}

func TestMainNative(t *testing.T) {
	main()

	oldFn := runValidation
	defer func() { runValidation = oldFn }()
	runValidation = func() bool { return false }
	main()
}
