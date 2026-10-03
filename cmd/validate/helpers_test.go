package main

import (
	"math"
	"math/cmplx"
	"testing"
)

// TestValidateWithinTol covers the float64 tolerance predicate including its
// special-value branches.
func TestValidateWithinTol(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		tol  float64
		want bool
	}{
		{"identical", 2.5, 2.5, 1e-10, true},
		{"nan both", math.NaN(), math.NaN(), 1e-10, true},
		{"nan one", math.NaN(), 1, 1e-10, false},
		{"+inf both", math.Inf(1), math.Inf(1), 1e-10, true},
		{"-inf both", math.Inf(-1), math.Inf(-1), 1e-10, true},
		{"inf opposite signs", math.Inf(1), math.Inf(-1), 1e-10, false},
		{"inf vs finite", math.Inf(1), 1, 1e-10, false},
		{"within absolute", 1, 1 + 1e-12, 1e-10, true},
		{"outside tolerance", 1, 1.5, 1e-10, false},
		{"within relative", 1e8, 1e8 * (1 + 1e-12), 1e-10, true},
	}
	for _, tc := range tests {
		if got := withinTol(tc.a, tc.b, tc.tol); got != tc.want {
			t.Errorf("%s: withinTol(%v, %v, %v) = %v, want %v", tc.name, tc.a, tc.b, tc.tol, got, tc.want)
		}
	}
}

// TestValidateWithinTolFloat32 covers the single-precision tolerance used by
// the float32 and complex64 validation paths.
func TestValidateWithinTolFloat32(t *testing.T) {
	tests := []struct {
		name string
		a, b float32
		want bool
	}{
		{"identical", 1.5, 1.5, true},
		{"nan both", float32(math.NaN()), float32(math.NaN()), true},
		{"nan one", float32(math.NaN()), 1, false},
		{"+inf both", float32(math.Inf(1)), float32(math.Inf(1)), true},
		{"-inf both", float32(math.Inf(-1)), float32(math.Inf(-1)), true},
		{"inf vs finite", float32(math.Inf(1)), 1, false},
		{"within tolerance", 1, 1 + 1e-8, true},
		{"outside tolerance", 1, 1.5, false},
	}
	for _, tc := range tests {
		if got := withinTolFloat32(tc.a, tc.b); got != tc.want {
			t.Errorf("%s: withinTolFloat32(%v, %v) = %v, want %v", tc.name, tc.a, tc.b, got, tc.want)
		}
	}
}

// TestValidateWithinTolComplex128 covers the complex128 tolerance predicate,
// which has separate handling for NaN and infinity in each component.
func TestValidateWithinTolComplex128(t *testing.T) {
	nan := math.NaN()
	inf := math.Inf(1)
	ninf := math.Inf(-1)

	tests := []struct {
		name string
		a, b complex128
		want bool
	}{
		{"identical", complex(1, 2), complex(1, 2), true},
		{"nan both parts", complex(nan, nan), complex(nan, nan), true},
		{"nan real only", complex(nan, 1), complex(nan, 1), true},
		{"imag nan only", complex(1, nan), complex(1, nan), true},
		{"real nan but imag differs", complex(nan, 1), complex(nan, 2), false},
		{"imag nan but real differs", complex(1, nan), complex(2, nan), false},
		{"+inf real both", complex(inf, 0), complex(inf, 0), true},
		{"+inf real both nan imag", complex(inf, nan), complex(inf, nan), true},
		{"+inf real both, imag differs", complex(inf, 1), complex(inf, 2), false},
		{"-inf real both", complex(ninf, 0), complex(ninf, 0), true},
		{"-inf real both nan imag", complex(ninf, nan), complex(ninf, nan), true},
		{"plain close values", complex(1, 1), complex(1, 1+1e-15), true},
		{"plain distant values", complex(1, 1), complex(5, 5), false},
	}
	for _, tc := range tests {
		if got := withinTolComplex128(tc.a, tc.b); got != tc.want {
			t.Errorf("%s: withinTolComplex128(%v, %v) = %v, want %v", tc.name, tc.a, tc.b, got, tc.want)
		}
	}
}

// TestValidateWithinTolComplex64 checks the mixed complex64/complex128
// comparison used for the complex64 validation path.
func TestValidateWithinTolComplex64(t *testing.T) {
	tests := []struct {
		name string
		a    complex64
		b    complex128
		want bool
	}{
		{"equivalent", complex64(1 + 2i), complex128(1 + 2i), true},
		{"rounding within tolerance", complex64(1 + 2i), complex128(1 + 2.0000001i), true},
		{"distant real", complex64(1), complex128(5), false},
		{"distant imag", complex64(complex(0, 1)), complex128(complex(0, 5)), false},
		{"nan both", complex64(complex(float32(math.NaN()), 0)), complex128(complex(math.NaN(), 0)), true},
	}
	for _, tc := range tests {
		if got := withinTolComplex64(tc.a, tc.b); got != tc.want {
			t.Errorf("%s: withinTolComplex64(%v, %v) = %v, want %v", tc.name, tc.a, tc.b, got, tc.want)
		}
	}
}

// TestFilterPassed covers the pass/fail filter used by every validator.
func TestFilterPassed(t *testing.T) {
	verbose, failedOnly = false, false

	// Empty and all-passing slices report success.
	if !filterPassed(nil) {
		t.Error("filterPassed(nil) = false, want true")
	}
	if !filterPassed([]ValidationResult{{Passed: true}, {Passed: true}}) {
		t.Error("filterPassed on all-passing = false, want true")
	}

	// A single failure fails the whole set.
	if filterPassed([]ValidationResult{{Passed: true}, {Passed: false, Type: "int", Function: "Add"}}) {
		t.Error("filterPassed on a failing entry = true, want false")
	}
}

// TestFilterPassedVerboseOutput exercises both formatting branches of
// filterPassed's failure reporting.
func TestFilterPassedVerboseOutput(t *testing.T) {
	t.Cleanup(func() { verbose, failedOnly = false, false })

	failing := []ValidationResult{{Passed: false, Type: "int", Function: "Add", Message: "mismatch"}}

	// Default output omits the message.
	verbose, failedOnly = false, false
	if filterPassed(failing) {
		t.Error("expected failure")
	}

	// verbose prints the message.
	verbose, failedOnly = true, false
	if filterPassed(failing) {
		t.Error("expected failure")
	}

	// failedOnly also prints the message.
	verbose, failedOnly = false, true
	if filterPassed(failing) {
		t.Error("expected failure")
	}
}

// TestSummarize covers the summary reporting for passing, failing and empty
// result sets. summarize reports rather than exiting, so it is testable; the
// printSummary wrapper terminates the process.
func TestSummarize(t *testing.T) {
	tests := []struct {
		name    string
		results []ValidationResult
		want    bool
	}{
		{"empty", nil, true},
		{"all passing", []ValidationResult{{Passed: true}, {Passed: true}}, true},
		{
			"one failing",
			[]ValidationResult{
				{Passed: true},
				{Passed: false, Type: "float64", Function: "Exp", Message: "bad"},
			},
			false,
		},
		{
			"all failing",
			[]ValidationResult{
				{Passed: false, Type: "int", Function: "Add", Message: "x"},
				{Passed: false, Type: "int", Function: "Sub", Message: "y"},
			},
			false,
		},
	}
	for _, tc := range tests {
		if got := summarize(tc.results); got != tc.want {
			t.Errorf("%s: summarize = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestPrintSummaryAllPassing verifies printSummary runs to completion (and
// does not exit) when every result passed.
func TestPrintSummaryAllPassing(t *testing.T) {
	prev := allResults
	t.Cleanup(func() { allResults = prev })

	verbose, failedOnly = false, false
	allResults = []ValidationResult{{Passed: true}, {Passed: true}}
	printSummary()
}

// TestValidateComplexHelpers checks the hand-rolled complex helpers used by
// the complex64/complex128 validators against math/cmplx. Each helper takes
// (real, imag) and must reproduce the corresponding cmplx function.
func TestValidateComplexHelpers(t *testing.T) {
	cases := []struct {
		r, i float64
	}{
		{0.75, -1.25},
		{0, 0},
		{-2, 0.5},
		{1e-3, 1e-3},
		{3, -2},
		{1e200, 1e200},
	}
	for _, c := range cases {
		z := complex(c.r, c.i)
		if got, want := trigComplexSin(c.r, c.i), cmplx.Sin(z); !closeEnough(got, want) {
			t.Errorf("trigComplexSin(%v, %v) = %v, want %v", c.r, c.i, got, want)
		}
		if got, want := trigComplexCos(c.r, c.i), cmplx.Cos(z); !closeEnough(got, want) {
			t.Errorf("trigComplexCos(%v, %v) = %v, want %v", c.r, c.i, got, want)
		}
		if got, want := complexExp(c.r, c.i), cmplx.Exp(z); !closeEnough(got, want) {
			t.Errorf("complexExp(%v, %v) = %v, want %v", c.r, c.i, got, want)
		}
		if got, want := complexLog(c.r, c.i), cmplx.Log(z); !closeEnough(got, want) {
			t.Errorf("complexLog(%v, %v) = %v, want %v", c.r, c.i, got, want)
		}
		if got, want := complexSqrt(c.r, c.i), cmplx.Sqrt(z); !closeEnough(got, want) {
			t.Errorf("complexSqrt(%v, %v) = %v, want %v", c.r, c.i, got, want)
		}
	}
}

// TestValidateComplexHelpersSpecialValues checks the infinity and overflow
// handling in the complex helpers, which the finite-value test cannot reach.
func TestValidateComplexHelpersSpecialValues(t *testing.T) {
	inf := math.Inf(1)

	// Infinity handling in complexLog.
	for _, c := range []struct{ r, i float64 }{
		{inf, inf},
		{inf, 0},
		{inf, 1},
		{0, inf},
		{1, inf},
	} {
		if got := complexLog(c.r, c.i); math.IsNaN(real(got)) {
			t.Errorf("complexLog(%v, %v) = %v has NaN real part", c.r, c.i, got)
		}
	}

	// Infinity handling in complexSqrt.
	for _, c := range []struct{ r, i float64 }{
		{inf, inf},
		{inf, -inf},
		{inf, 0},
		{-inf, 0},
		{0, inf},
	} {
		if got := complexSqrt(c.r, c.i); math.IsNaN(real(got)) {
			t.Errorf("complexSqrt(%v, %v) = %v has NaN real part", c.r, c.i, got)
		}
	}

	// Very large magnitudes must not overflow to NaN.
	for _, c := range []struct{ r, i float64 }{
		{1e200, 0},
		{1e200, 1e-150},
		{-1e200, 1e-150},
	} {
		if got := complexLog(c.r, c.i); math.IsNaN(real(got)) || math.IsNaN(imag(got)) {
			t.Errorf("complexLog(%v, %v) = %v contains NaN", c.r, c.i, got)
		}
	}
}

// TestValidateSqrtIdentity checks sqrt(z)^2 == z across both branch cuts for
// the hand-rolled complex square root.
func TestValidateSqrtIdentity(t *testing.T) {
	for _, c := range []struct{ r, i float64 }{
		{4, 0},
		{-1, 0},
		{-1e-8, 1e-8},
		{-1e-8, -1e-8},
		{3, 4},
		{-3, 4},
		{0, 1},
	} {
		s := complexSqrt(c.r, c.i)
		if diff := cmplx.Abs(s*s - complex(c.r, c.i)); diff > 1e-9*(1+cmplx.Abs(complex(c.r, c.i))) {
			t.Errorf("complexSqrt(%v, %v)^2 = %v, want %v", c.r, c.i, s*s, complex(c.r, c.i))
		}
	}
}

// TestValidatePerTypeValidators runs each per-type validator directly and
// checks that every produced result carries the type name it was asked for.
func TestValidatePerTypeValidators(t *testing.T) {
	verbose, failedOnly = false, false

	cases := []struct {
		typeName string
		fn       func(string) []ValidationResult
	}{
		{"int", validateIntType},
		{"uint", validateUintType},
		{"float32", validateFloatType},
		{"complex128", validateComplexType},
	}
	for _, c := range cases {
		results := c.fn(c.typeName)
		if len(results) == 0 {
			t.Errorf("%s: validator produced no results", c.typeName)
			continue
		}
		for _, r := range results {
			if r.Type != c.typeName {
				t.Errorf("%s: result has type %q", c.typeName, r.Type)
			}
			if r.Function == "" {
				t.Errorf("%s: result with empty function name: %+v", c.typeName, r)
			}
		}
	}
}

// TestValidateIndividualSuites checks the standalone type suites report
// results and that none of them regress.
func TestValidateIndividualSuites(t *testing.T) {
	verbose, failedOnly = false, false

	intSuites := map[string]func(string) []ValidationResult{
		"int":   testInt[int],
		"int8":  testInt[int8],
		"int16": testInt[int16],
		"int32": testInt[int32],
		"int64": testInt[int64],
	}
	uintSuites := map[string]func(string) []ValidationResult{
		"uint":    testUint[uint],
		"uint8":   testUint[uint8],
		"uint16":  testUint[uint16],
		"uint32":  testUint[uint32],
		"uint64":  testUint[uint64],
		"uintptr": testUint[uintptr],
	}
	floatSuites := map[string]func() []ValidationResult{
		"float32": testFloat32,
		"float64": testFloat64,
	}
	complexSuites := map[string]func() []ValidationResult{
		"complex64":  testComplex64,
		"complex128": testComplex128,
	}

	for name, suite := range intSuites {
		checkSuiteResults(t, name, suite(name))
	}
	for name, suite := range uintSuites {
		checkSuiteResults(t, name, suite(name))
	}
	for name, suite := range floatSuites {
		checkSuiteResults(t, name, suite())
	}
	for name, suite := range complexSuites {
		checkSuiteResults(t, name, suite())
	}
}

// checkSuiteResults asserts a validation suite produced results and that none
// of them regressed.
func checkSuiteResults(t *testing.T, name string, results []ValidationResult) {
	t.Helper()
	if len(results) == 0 {
		t.Errorf("%s: suite produced no results", name)
		return
	}
	for _, r := range results {
		if !r.Passed {
			t.Errorf("%s/%s failed: %s", name, r.Function, r.Message)
		}
	}
}

// closeEnough compares complex values with a relative tolerance. Matching
// infinities and matching NaNs are treated as equal component-wise.
func closeEnough(a, b complex128) bool {
	return closeComponent(real(a), real(b)) && closeComponent(imag(a), imag(b))
}

func closeComponent(a, b float64) bool {
	switch {
	case math.IsNaN(a) && math.IsNaN(b):
		return true
	case math.IsNaN(a) || math.IsNaN(b):
		return false
	case math.IsInf(a, 0) || math.IsInf(b, 0):
		return a == b
	}
	scale := 1 + math.Abs(b)
	return math.Abs(a-b) <= 1e-9*scale
}
