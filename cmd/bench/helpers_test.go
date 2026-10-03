package main

import (
	"math"
	"math/cmplx"
	"strings"
	"testing"
)

// TestWithinTol exercises the tolerance predicate used by the parity checks,
// including the special-value branches.
func TestWithinTol(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		tol  float64
		want bool
	}{
		{"identical", 1.5, 1.5, 1e-10, true},
		{"nan both", math.NaN(), math.NaN(), 1e-10, true},
		{"nan one", math.NaN(), 1, 1e-10, false},
		{"+inf both", math.Inf(1), math.Inf(1), 1e-10, true},
		{"-inf both", math.Inf(-1), math.Inf(-1), 1e-10, true},
		{"+inf mixed sign", math.Inf(1), math.Inf(-1), 1e-10, false},
		{"inf vs finite", math.Inf(1), 1, 1e-10, false},
		{"within absolute tol", 1.0, 1.0 + 1e-12, 1e-10, true},
		{"outside absolute tol", 1.0, 1.0 + 1e-6, 1e-10, false},
		{"within relative tol", 1e6, 1e6 * (1 + 1e-12), 1e-10, true},
		{"outside relative tol", 1.0, 1.5, 1e-10, false},
		{"zero tolerance equal", 3, 3, 0, false},
	}
	for _, tc := range tests {
		if got := withinTol(tc.a, tc.b, tc.tol); got != tc.want {
			t.Errorf("%s: withinTol(%v, %v, %v) = %v, want %v", tc.name, tc.a, tc.b, tc.tol, got, tc.want)
		}
	}
}

// TestULPDiff covers the ULP distance helper, including the early exits.
func TestULPDiff(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want uint64
	}{
		{"equal", 1.5, 1.5, 0},
		{"zero vs zero", 0, 0, 0},
		{"one ulp", 1.0, math.Nextafter(1.0, 2), 1},
		{"two ulp", 1.0, math.Nextafter(math.Nextafter(1.0, 2), 2), 2},
		{"nan either", math.NaN(), 1, math.MaxUint64},
		{"either nan", 1, math.NaN(), math.MaxUint64},
		{"both nan", math.NaN(), math.NaN(), 0},
		{"inf either", math.Inf(1), 1, math.MaxUint64},
		{"either inf", 1, math.Inf(1), math.MaxUint64},
		{"both inf", math.Inf(1), math.Inf(1), 0},
		{"both neg inf", math.Inf(-1), math.Inf(-1), 0},
		{"mismatched inf", math.Inf(1), math.Inf(-1), math.MaxUint64},
		{"neg inf either", math.Inf(-1), 1, math.MaxUint64},
		{"opposite sign", 1.0, -1.0, (math.Float64bits(1.0) & 0x7fffffffffffffff) + (math.Float64bits(-1.0) & 0x7fffffffffffffff)},
	}
	for _, tc := range tests {
		if got := ulpDiff(tc.a, tc.b); got != tc.want {
			t.Errorf("%s: ulpDiff(%v, %v) = %d, want %d", tc.name, tc.a, tc.b, got, tc.want)
		}
	}

	// Symmetry: the distance is unsigned.
	a, b := 1.0, 2.0
	if ulpDiff(a, b) != ulpDiff(b, a) {
		t.Errorf("ulpDiff is not symmetric: %d != %d", ulpDiff(a, b), ulpDiff(b, a))
	}
}

// TestParityChecksPass runs every scalar parity check; each must agree with
// the math library across the sampled range.
func TestParityChecksPass(t *testing.T) {
	verbose = false
	checks := map[string]func() bool{
		"Exp":   testExpParity,
		"Log":   testLogParity,
		"Sin":   testSinParity,
		"Cos":   testCosParity,
		"Tan":   testTanParity,
		"Sinh":  testSinhParity,
		"Cosh":  testCoshParity,
		"Tanh":  testTanhParity,
		"Asinh": testAsinhParity,
		"Acosh": testAcoshParity,
		"Atanh": testAtanhParity,
		"Sqrt":  testSqrtParity,
		"Pow":   testPowParity,
	}
	for name, fn := range checks {
		if !fn() {
			t.Errorf("%s parity check failed", name)
		}
	}
}

// TestAccuracyChecksPass runs every ULP accuracy check; each must stay within
// the 200 ULP budget used by the tool.
func TestAccuracyChecksPass(t *testing.T) {
	verbose = false
	checks := map[string]func() (uint64, bool){
		"Exp":   testExpAccuracy,
		"Log":   testLogAccuracy,
		"Sin":   testSinAccuracy,
		"Cos":   testCosAccuracy,
		"Tan":   testTanAccuracy,
		"Sinh":  testSinhAccuracy,
		"Cosh":  testCoshAccuracy,
		"Tanh":  testTanhAccuracy,
		"Asinh": testAsinhAccuracy,
		"Acosh": testAcoshAccuracy,
		"Atanh": testAtanhAccuracy,
		"Sqrt":  testSqrtAccuracy,
		"Pow":   testPowAccuracy,
	}
	for name, fn := range checks {
		ulps, ok := fn()
		if !ok {
			t.Errorf("%s accuracy check failed at %d ULP (budget 200)", name, ulps)
		}
	}
}

// TestTestOpParityDetectsMismatch verifies the parity harness actually fails
// when the emlgo and math implementations disagree.
func TestTestOpParityDetectsMismatch(t *testing.T) {
	verbose = false
	// Identical functions must pass.
	if !testOpParity(func(x float64) float64 { return x }, func(x float64) float64 { return x }) {
		t.Error("identical functions reported as non-parity")
	}
	// A deliberate mismatch must be caught.
	if testOpParity(func(x float64) float64 { return x }, func(x float64) float64 { return x + 1 }) {
		t.Error("mismatched functions reported as parity")
	}
}

// TestTestOpAccuracyDetectsMismatch verifies the ULP harness rejects a bad
// implementation and reports a plausible ULP count.
func TestTestOpAccuracyDetectsMismatch(t *testing.T) {
	verbose = false
	if ulps, ok := testOpAccuracy(func(x float64) float64 { return x }, func(x float64) float64 { return x }); !ok || ulps != 0 {
		t.Errorf("identical functions: got (%d, %v), want (0, true)", ulps, ok)
	}
	ulps, ok := testOpAccuracy(func(x float64) float64 { return x * 2 }, func(x float64) float64 { return x })
	if ok {
		t.Error("2x mismatch accepted as accurate")
	}
	if ulps == 0 {
		t.Error("2x mismatch reported 0 ULP difference")
	}
}

// TestBenchmarkRunnersProduceResults drives every benchmark runner with a
// single iteration and checks the results are well formed.
func TestBenchmarkRunnersProduceResults(t *testing.T) {
	prev := iterations
	iterations = 1
	t.Cleanup(func() { iterations = prev })

	runners := map[string]func() []BenchmarkResult{
		"int":        runIntBenchmarks,
		"uint":       runUintBenchmarks,
		"float32":    runFloat32Benchmarks,
		"float64":    runFloat64Benchmarks,
		"complex64":  runComplex64Benchmarks,
		"complex128": runComplex128Benchmarks,
		"fastmath":   runFastMathBenchmarks,
		"batch":      runBatchBenchmarks,
	}
	for name, run := range runners {
		results := run()
		if len(results) == 0 {
			t.Errorf("%s runner returned no results", name)
			continue
		}
		for _, r := range results {
			if r.Type == "" {
				t.Errorf("%s: result with empty type: %+v", name, r)
			}
			if r.Name == "" {
				t.Errorf("%s: result with empty name: %+v", name, r)
			}
			if r.EmlgoTime < 0 || r.MathTime < 0 {
				t.Errorf("%s/%s: negative timing %+v", r.Type, r.Name, r)
			}
		}
	}
}

// TestRunAllBenchmarksTypeFilter checks the type filter selects the right
// subset of benchmarks.
func TestRunAllBenchmarksTypeFilter(t *testing.T) {
	prevIterations, prevFilter := iterations, typeFilter
	iterations = 1
	t.Cleanup(func() { iterations, typeFilter = prevIterations, prevFilter })

	for _, filter := range []string{"all", "int", "uint", "float32", "float64", "complex64", "complex128"} {
		typeFilter = filter
		results := runAllBenchmarks()
		if len(results) == 0 {
			t.Errorf("typeFilter=%q produced no results", filter)
			continue
		}
		for _, r := range results {
			if filter != "all" && r.Type != filter {
				t.Errorf("typeFilter=%q returned type %q", filter, r.Type)
			}
		}
	}

	// An unknown filter selects nothing rather than everything.
	typeFilter = "nonexistent"
	if got := runAllBenchmarks(); len(got) != 0 {
		t.Errorf("typeFilter=nonexistent returned %d results, want 0", len(got))
	}
}

// TestPrintResultsEmpty verifies printResults handles an empty result set.
func TestPrintResultsEmpty(t *testing.T) {
	// Must not panic or divide by zero on an empty slice.
	printResults(nil)
	printResults([]BenchmarkResult{})
}

// TestClassifyRegressions covers the baseline comparison logic, which is
// otherwise only reachable through the -regression command-line flag.
func TestClassifyRegressions(t *testing.T) {
	// "float64/Exp" has a baseline of 1.10.
	tests := []struct {
		name        string
		results     []BenchmarkResult
		wantRegress int
		wantImprove int
	}{
		{
			name:        "empty input",
			results:     nil,
			wantRegress: 0,
			wantImprove: 0,
		},
		{
			name:        "no baseline entry is ignored",
			results:     []BenchmarkResult{{Type: "unknown", Name: "Thing", Ratio: 99}},
			wantRegress: 0,
			wantImprove: 0,
		},
		{
			name:        "within tolerance",
			results:     []BenchmarkResult{{Type: "float64", Name: "Exp", Ratio: 1.15}},
			wantRegress: 0,
			wantImprove: 0,
		},
		{
			name:        "regression above threshold",
			results:     []BenchmarkResult{{Type: "float64", Name: "Exp", Ratio: 1.50}},
			wantRegress: 1,
			wantImprove: 0,
		},
		{
			name:        "improvement below negative threshold",
			results:     []BenchmarkResult{{Type: "float64", Name: "Exp", Ratio: 0.50}},
			wantRegress: 0,
			wantImprove: 1,
		},
		{
			name: "mixed",
			results: []BenchmarkResult{
				{Type: "float64", Name: "Exp", Ratio: 1.50},
				{Type: "int", Name: "Add", Ratio: 2.00},
				{Type: "uint", Name: "Mul", Ratio: 0.10},
				{Type: "float64", Name: "Log", Ratio: 1.00},
			},
			wantRegress: 2,
			wantImprove: 1,
		},
	}
	for _, tc := range tests {
		regressions, improvements := classifyRegressions(tc.results)
		if len(regressions) != tc.wantRegress {
			t.Errorf("%s: got %d regressions, want %d (%v)", tc.name, len(regressions), tc.wantRegress, regressions)
		}
		if len(improvements) != tc.wantImprove {
			t.Errorf("%s: got %d improvements, want %d (%v)", tc.name, len(improvements), tc.wantImprove, improvements)
		}
		for _, msg := range regressions {
			if !strings.Contains(msg, "REGRESSION") {
				t.Errorf("%s: regression message %q missing keyword", tc.name, msg)
			}
		}
		for _, msg := range improvements {
			if !strings.Contains(msg, "IMPROVEMENT") {
				t.Errorf("%s: improvement message %q missing keyword", tc.name, msg)
			}
		}
	}
}

// TestClassifyRegressionsUsesBaseline verifies the classifier is keyed on
// "Type/Name" against the recorded baseline map.
func TestClassifyRegressionsUsesBaseline(t *testing.T) {
	for key := range baseline {
		typ, name, ok := strings.Cut(key, "/")
		if !ok {
			t.Errorf("baseline key %q is not in Type/Name form", key)
			continue
		}
		// A ratio far above the baseline must register as a regression.
		regressions, _ := classifyRegressions([]BenchmarkResult{
			{Type: typ, Name: name, Ratio: baseline[key] + 1},
		})
		if len(regressions) != 1 {
			t.Errorf("baseline key %q: got %d regressions, want 1", key, len(regressions))
		}
	}
}

// TestComplexHelpers checks the complex helper closures used by the complex
// benchmark runners against math/cmplx. For z = r + i*im the helpers must
// reproduce sin(z), cos(z), tan(z) and exp(z).
func TestComplexHelpers(t *testing.T) {
	cases := []struct {
		r, im float64
	}{
		{0.75, -1.25},
		{0, 0},
		{-2, 0.5},
		{1e-3, 1e-3},
		{3, -2},
	}
	for _, c := range cases {
		z := complex(c.r, c.im)

		if got, want := trigComplexSin(c.r, c.im), cmplx.Sin(z); !closeComplex(got, want) {
			t.Errorf("trigComplexSin(%v, %v) = %v, want %v", c.r, c.im, got, want)
		}
		if got, want := trigComplexCos(c.r, c.im), cmplx.Cos(z); !closeComplex(got, want) {
			t.Errorf("trigComplexCos(%v, %v) = %v, want %v", c.r, c.im, got, want)
		}
		if got, want := complexExp(c.r, c.im), cmplx.Exp(z); !closeComplex(got, want) {
			t.Errorf("complexExp(%v, %v) = %v, want %v", c.r, c.im, got, want)
		}
		// tan is defined as sin/cos away from the poles.
		if got, want := trigComplexTan(c.r, c.im), cmplx.Tan(z); !closeComplex(got, want) {
			t.Errorf("trigComplexTan(%v, %v) = %v, want %v", c.r, c.im, got, want)
		}

		// sin(z)^2 + cos(z)^2 = 1 for all complex z.
		sum := trigComplexSin(c.r, c.im)*trigComplexSin(c.r, c.im) +
			trigComplexCos(c.r, c.im)*trigComplexCos(c.r, c.im)
		if !closeComplex(sum, 1) {
			t.Errorf("sin^2+cos^2 at %v = %v, want 1", z, sum)
		}
	}
}

func closeComplex(a, b complex128) bool {
	d := cmplx.Abs(a - b)
	return d < 1e-9*(1+cmplx.Abs(b))
}
