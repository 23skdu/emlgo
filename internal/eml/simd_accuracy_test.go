package eml

import (
	"math"
	"testing"
)

// Accuracy budgets for the hand-written SIMD batch kernels.
//
// These are ULP tolerances against the Go standard library, chosen so that a
// regression in a kernel is caught by the test suite rather than discovered by a
// user. They are deliberately tighter than the ~1e-7 accuracy of the
// fastmath approximations, because the batch kernels are expected to be
// full-float64 implementations.
//
// A kernel whose true value passes through zero (cos near π/2) cannot be judged
// in ULPs: cos(π/2) is ~6.1e-17, so any absolute error of 1e-11 is ~1e5 ULPs.
// Those kernels are therefore additionally checked with an absolute tolerance.

const (
	// Mixed tolerance: |got - want| <= absTol + relTol*|want|.
	//
	// A purely absolute budget is unsatisfiable for large results (1e-13
	// absolute on exp(9.4) ~ 1.2e4 is below float64 epsilon), and a purely
	// relative budget is meaningless where the true value crosses zero
	// (cos(pi/2) ~ 6.1e-17). The mixed form handles both ends.
	//
	// absTol also covers the range-reduction discontinuity of the hand-written
	// sin/cos kernels at +-pi/2, which is the dominant error term there.
	budgetAbsTol = 1e-14
	budgetRelTol = 1e-15 // about 5 ULP

	// ULP budgets are enforced only where the true value is large enough for
	// ULP distance to be meaningful.
	budgetExpULP  = 64
	budgetLogULP  = 8
	budgetSinULP  = 1024
	budgetCosULP  = 1024
	budgetTanULP  = 1024
	budgetSqrtULP = 4

	// Below this magnitude a ULP count is dominated by the near-zero crossing
	// and the mixed tolerance is the meaningful gate.
	ulpMeaningfulAbove = 1e-9
)

// withinBudget reports whether err is within the mixed absolute/relative budget
// for a true value of want.
func withinBudget(err, want float64) bool {
	return err <= budgetAbsTol+budgetRelTol*math.Abs(want)
}

// ulpDistance returns the number of representable float64 values between a and
// b. It returns 0 for NaN/Inf pairs and for bit-identical values.
func ulpDistance(a, b float64) uint64 {
	if a == b {
		return 0
	}
	if math.IsNaN(a) || math.IsNaN(b) || math.IsInf(a, 0) || math.IsInf(b, 0) {
		return 0
	}
	x, y := math.Float64bits(a), math.Float64bits(b)
	if x > y {
		return x - y
	}
	return y - x
}

// denseGrid returns n evenly spaced values over [lo, hi].
func denseGrid(lo, hi float64, n int) []float64 {
	x := make([]float64, n)
	if n == 1 {
		x[0] = lo
		return x
	}
	for i := range x {
		x[i] = lo + (hi-lo)*float64(i)/float64(n-1)
	}
	return x
}

// TestBatchKernelAccuracy checks every SIMD transcendental kernel against the
// standard library over a dense grid.
//
// This test is the gate for the kernel rewrites: it must pass on every path
// (AVX2, AVX-512, generic), so a fast-but-wrong kernel cannot be shipped.
func TestBatchKernelAccuracy(t *testing.T) {
	// A grid size that is a multiple of 8 exercises both the vectorised body
	// and the scalar tail loop.
	const n = 4096
	dst := make([]float64, n)

	cases := []struct {
		name      string
		in        []float64
		kernel    func(in, dst []float64)
		reference func(float64) float64
		ulpBudget uint64
		absBudget float64
	}{
		{
			name:      "exp",
			in:        denseGrid(-10, 10, n),
			kernel:    ExpSIMDTo,
			reference: math.Exp,
			ulpBudget: budgetExpULP,
		},
		{
			name:      "log",
			in:        denseGrid(1e-6, 100, n),
			kernel:    LogSIMDTo,
			reference: ExactLogOracle,
			ulpBudget: budgetLogULP,
		},
		{
			name:      "sin",
			in:        denseGrid(-math.Pi, math.Pi, n),
			kernel:    SinSIMDTo,
			reference: math.Sin,
			ulpBudget: budgetSinULP,
		},
		{
			name:      "cos",
			in:        denseGrid(-math.Pi, math.Pi, n),
			kernel:    CosSIMDTo,
			reference: math.Cos,
			ulpBudget: budgetCosULP,
		},
		{
			name:      "tan",
			in:        denseGrid(-1.5, 1.5, n),
			kernel:    TanSIMDTo,
			reference: math.Tan,
			ulpBudget: budgetTanULP,
		},
		{
			name:      "sqrt",
			in:        denseGrid(0, 1e6, n),
			kernel:    SqrtSIMDTo,
			reference: math.Sqrt,
			ulpBudget: budgetSqrtULP,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.kernel(tc.in, dst)

			var (
				worstULP  uint64
				worstULPX float64
				worstRel  float64
				worstRelX float64
				worstAbs  float64
				worstAbsX float64
			)

			for i, x := range tc.in {
				want := tc.reference(x)
				got := dst[i]
				if math.IsNaN(want) || math.IsInf(want, 0) {
					// Outside the kernel's domain; require matching class.
					if math.IsNaN(got) != math.IsNaN(want) ||
						math.IsInf(got, 0) != math.IsInf(want, 0) {
						t.Fatalf("%s(%v) = %v, want %v", tc.name, x, got, want)
					}
					continue
				}

				absErr := math.Abs(got - want)
				if absErr > worstAbs {
					worstAbs, worstAbsX = absErr, x
				}
				if den := math.Abs(want); den > 0 {
					if rel := absErr / den; rel > worstRel {
						worstRel, worstRelX = rel, x
					}
				}

				// Primary gate: mixed absolute/relative tolerance.
				if !withinBudget(absErr, want) {
					t.Errorf("%s(%v) = %v, want %v (abs err %.3g, budget %.3g)",
						tc.name, x, got, want, absErr,
						budgetAbsTol+budgetRelTol*math.Abs(want))
				}

				// Secondary gate: ULP distance, only where it is meaningful.
				if math.Abs(want) > ulpMeaningfulAbove {
					if u := ulpDistance(got, want); u > worstULP {
						worstULP, worstULPX = u, x
					}
				}
			}

			if worstULP > tc.ulpBudget {
				t.Errorf("%s: worst ULP %d at x=%v exceeds budget %d",
					tc.name, worstULP, worstULPX, tc.ulpBudget)
			}
			t.Logf("%s: maxULP=%d (x=%v) maxRel=%.3g (x=%v) maxAbs=%.3g (x=%v)",
				tc.name, worstULP, worstULPX, worstRel, worstRelX, worstAbs, worstAbsX)
		})
	}
}

// TestBatchKernelTailLoopAccuracy covers lengths that are not a multiple of the
// SIMD width, exercising the scalar remainder path for n-1 down to n-8.
func TestBatchKernelTailLoopAccuracy(t *testing.T) {
	for _, n := range []int{1, 2, 3, 5, 7, 8, 9, 15, 16, 17, 31, 33, 63} {
		in := denseGrid(0.1, 4.0, n)
		dst := make([]float64, n)

		ExpSIMDTo(in, dst)
		for i, x := range in {
			if got, want := dst[i], math.Exp(x); !withinBudget(math.Abs(got-want), want) {
				t.Errorf("n=%d: exp(%v) = %v, want %v (abs err %.3g)", n, x, got, want, math.Abs(got-want))
			}
		}

		LogSIMDTo(in, dst)
		for i, x := range in {
			if got, want := dst[i], math.Log(x); !withinBudget(math.Abs(got-want), want) {
				t.Errorf("n=%d: log(%v) = %v, want %v (abs err %.3g)", n, x, got, want, math.Abs(got-want))
			}
		}
	}
}

// TestBatchTrigIdentities checks the defining identities on the batch path.
// The scalar path is bit-exact against math; the batch path must be close.
func TestBatchTrigIdentities(t *testing.T) {
	const n = 4096
	in := denseGrid(-math.Pi, math.Pi, n)
	sin := make([]float64, n)
	cos := make([]float64, n)

	SinSIMDTo(in, sin)
	CosSIMDTo(in, cos)

	for i, x := range in {
		sum := sin[i]*sin[i] + cos[i]*cos[i]
		if d := math.Abs(sum - 1); !withinBudget(d, 1) {
			t.Fatalf("sin^2+cos^2 at x=%v = %v (diff %.3g)", x, sum, d)
		}
	}
}

// TestBatchSinCosZeros pins the argument-reduction behaviour at the reduction
// boundaries, where the hand-written kernels lost the most accuracy: the range
// reduction is discontinuous at ±π/2, and cos is exactly zero there.
func TestBatchSinCosZeros(t *testing.T) {
	const n = 4096
	// Sweep densely across the ±π/2 boundaries.
	in := make([]float64, 0, 4*n)
	for _, centre := range []float64{-math.Pi / 2, math.Pi / 2} {
		in = append(in, denseGrid(centre-1e-3, centre+1e-3, n)...)
	}
	cos := make([]float64, len(in))
	CosSIMDTo(in, cos)

	var worst float64
	var worstX float64
	for i, x := range in {
		if d := math.Abs(cos[i] - math.Cos(x)); d > worst {
			worst, worstX = d, x
		}
	}
	if !withinBudget(worst, 0) {
		t.Errorf("cos near ±π/2: worst absolute error %.3g at x=%.17g, budget %.3g",
			worst, worstX, budgetAbsTol)
	}
	t.Logf("cos near ±π/2: worst absolute error %.3g at x=%.17g", worst, worstX)
}

// TestBatchKernelMatchesScalarPath checks that the vectorised path agrees with
// the scalar reference implementation, which is what most callers actually use
// via pkg/trig and pkg/logexp.
func TestBatchKernelMatchesScalarPath(t *testing.T) {
	const n = 2048
	in := denseGrid(0.05, 5.0, n)
	dst := make([]float64, n)

	pairs := []struct {
		name   string
		kernel func(in, dst []float64)
		scalar func(float64) float64
	}{
		{"exp", ExpSIMDTo, nativeExp},
		{"log", LogSIMDTo, nativeLog},
	}
	for _, p := range pairs {
		p.kernel(in, dst)
		for i, x := range in {
			if got, want := dst[i], p.scalar(x); got != want {
				t.Errorf("%s(%v): batch %v != scalar %v", p.name, x, got, want)
			}
		}
	}
}

// TestBatchKernelSpecialValues checks that the kernels reproduce the standard
// library's behaviour at the edges of the domain, where a naive polynomial
// kernel is most likely to diverge.
func TestBatchKernelSpecialValues(t *testing.T) {
	// exp: underflow, overflow, and the largest finite input.
	for _, x := range []float64{-1000, -800, -745, -708, 0, 709, 710, 1000} {
		in := []float64{x}
		got := make([]float64, 1)
		ExpSIMDTo(in, got)
		want := math.Exp(x)
		if got[0] != want {
			t.Errorf("ExpSIMDTo(%v) = %v, want %v", x, got[0], want)
		}
	}

	// sqrt of negative and zero.
	for _, tc := range []struct{ in, want float64 }{
		{0, 0},
		{4, 2},
		{-1, math.NaN()},
	} {
		in := []float64{tc.in}
		got := make([]float64, 1)
		SqrtSIMDTo(in, got)
		if !sameFloatClass(got[0], tc.want) {
			t.Errorf("SqrtSIMDTo(%v) = %v, want %v", tc.in, got[0], tc.want)
		}
	}

	// log follows the reconciled domain rule (logScalar, matching pkg/logexp.Log):
	// -Inf for zero and NaN for negative arguments.
	for _, in := range []float64{-1, -1e300} {
		got := make([]float64, 1)
		LogSIMDTo([]float64{in}, got)
		if !math.IsNaN(got[0]) {
			t.Errorf("LogSIMDTo(%v) = %v, want NaN", in, got[0])
		}
	}
	for _, in := range []float64{0, math.Copysign(0, -1)} {
		got := make([]float64, 1)
		LogSIMDTo([]float64{in}, got)
		if !math.IsInf(got[0], -1) {
			t.Errorf("LogSIMDTo(%v) = %v, want -Inf", in, got[0])
		}
	}

	// sin/cos at large magnitudes need range reduction, not a raw polynomial.
	for _, in := range []float64{1e6, 1e9, -1e9} {
		gotSin := make([]float64, 1)
		gotCos := make([]float64, 1)
		SinSIMDTo([]float64{in}, gotSin)
		CosSIMDTo([]float64{in}, gotCos)
		wantSin, wantCos := math.Sin(in), math.Cos(in)
		if math.Abs(gotSin[0]-wantSin) > 1e-9 || math.Abs(gotCos[0]-wantCos) > 1e-9 {
			t.Errorf("sin/cos(%v) = %v/%v, want %v/%v", in, gotSin[0], gotCos[0], wantSin, wantCos)
		}
	}
}

// TestBatchLogMatchesScalarDomain checks that the batch logarithm applies the
// same domain rule as the scalar one, so Log and LogBatch can never disagree.
//
// The batch path used to call math.Log directly while pkg/logexp.Log returned
// NaN for x <= 0, which made Log(0) NaN and LogBatch([0]) -Inf.
func TestBatchLogMatchesScalarDomain(t *testing.T) {
	inputs := []float64{
		0, math.Copysign(0, -1), math.SmallestNonzeroFloat64, 1e-300,
		0.5, math.Nextafter(1, 0), 1, math.Nextafter(1, 2),
		2, math.MaxFloat64, -1e-300, -1, -math.MaxFloat64,
		math.NaN(), math.Inf(1), math.Inf(-1),
	}

	batch := LogSIMD(inputs)
	for i, x := range inputs {
		scalar := logScalar(x)
		got := batch[i]
		same := scalar == got || (math.IsNaN(scalar) && math.IsNaN(got))
		if !same {
			t.Errorf("logScalar(%v) = %v but LogSIMD = %v", x, scalar, got)
		}
	}
}

// TestScalarLogMatchesOracleForPositive checks that the guarded scalar logarithm matches
// the exact rescaling oracle across the positive domain, including subnormals.
func TestScalarLogMatchesOracleForPositive(t *testing.T) {
	for _, x := range []float64{5e-324, 1e-300, 1e-10, 0.5, 1, 1.5, 2, 1e300, math.MaxFloat64} {
		want := ExactLogOracle(x)
		if got := logScalar(x); math.Abs(got-want) > 1e-12 {
			t.Errorf("logScalar(%v) = %v, want ExactLogOracle = %v", x, got, want)
		}
	}
}

// sameFloatClass reports whether got and want agree, treating NaN as equal to
// NaN and requiring identical infinities.
func sameFloatClass(got, want float64) bool {
	switch {
	case math.IsNaN(want):
		return math.IsNaN(got)
	case math.IsInf(want, 0):
		return got == want
	default:
		return got == want
	}
}
