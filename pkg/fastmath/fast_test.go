package fastmath

import (
	"math"
	"math/rand"
	"testing"
)

func TestFastExpMinimaxAccuracy(t *testing.T) {
	testPoints := []float64{
		-100.0, -10.0, -5.0, -1.0, -0.5, -0.1, -0.01, 0.0,
		0.01, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 50.0, 100.0, 500.0,
	}
	for _, x := range testPoints {
		got := FastExpMinimax(x)
		want := math.Exp(x)
		relErr := math.Abs(got-want) / want
		if relErr > 1.5e-6 {
			t.Errorf("FastExpMinimax(%v): got %v, want %v, relErr %v > 1.5e-6", x, got, want, relErr)
		}
	}
}

func TestFastLogMinimaxAccuracy(t *testing.T) {
	testPoints := []float64{
		1e-10, 1e-5, 0.01, 0.1, 0.5, 0.707, 1.0, 1.414, 2.0, math.E,
		5.0, 10.0, 100.0, 1000.0, 1e6, 1e15,
	}
	for _, y := range testPoints {
		got := FastLogMinimax(y)
		want := math.Log(y)
		relErr := math.Abs(got-want) / math.Abs(want)
		if math.Abs(want) < 1e-6 {
			relErr = math.Abs(got - want)
		}
		if relErr > 5e-5 {
			t.Errorf("FastLogMinimax(%v): got %v, want %v, relErr %v > 5e-5", y, got, want, relErr)
		}
	}
}

func TestFastExpBitCast(t *testing.T) {
	for _, x := range []float64{-2.0, -1.0, 0.0, 0.5, 1.0, 2.0} {
		got := FastExpBitCast(x)
		want := math.Exp(x)
		relErr := math.Abs(got-want) / want
		if relErr > 0.05 { // Schraudolph relative error is ~1-3%
			t.Errorf("FastExpBitCast(%v): got %v, want %v, relErr %v > 5%%", x, got, want, relErr)
		}
	}
}

func TestApproximationAccuracyAcrossDomains(t *testing.T) {
	const (
		exp64RelBudget = 1.5e-6
		log64AbsBudget = 1e-7
		exp32RelBudget = 2e-4
		log32AbsBudget = 1e-5
	)

	type accuracySample struct {
		input, error, budget, ratio float64
	}
	worst := make(map[string]accuracySample)
	check64 := func(name string, got, want, x, absBudget, relBudget float64) {
		t.Helper()
		if math.IsNaN(want) {
			if !math.IsNaN(got) {
				t.Errorf("%s(%g) = %g, want NaN", name, x, got)
			}
			return
		}
		if math.IsInf(want, 0) {
			if got != want {
				t.Errorf("%s(%g) = %g, want %g", name, x, got, want)
			}
			return
		}
		err := math.Abs(got - want)
		budget := absBudget + relBudget*math.Abs(want)
		ratio := err / budget
		if ratio > worst[name].ratio {
			worst[name] = accuracySample{input: x, error: err, budget: budget, ratio: ratio}
		}
		if err > budget {
			t.Errorf("%s(%g) = %.17g, want %.17g (error %.3g > budget %.3g)",
				name, x, got, want, err, budget)
		}
	}

	check32 := func(name string, got, want, x float32, absBudget, relBudget float64) {
		t.Helper()
		check64(name, float64(got), float64(want), float64(x), absBudget, relBudget)
	}

	expInputs := make([]float64, 0, 25000)
	for x := -700.0; x <= 700.0; x += 0.125 {
		expInputs = append(expInputs, x)
	}
	for k := -1022; k <= 1023; k++ {
		x := float64(k) * math.Ln2
		expInputs = append(expInputs, math.Nextafter(x, math.Inf(-1)), x, math.Nextafter(x, math.Inf(1)))
	}
	expInputs = append(expInputs, -745.133224101734, math.Nextafter(-745.133224101734, math.Inf(1)),
		709.782712893384, math.Nextafter(709.782712893384, math.Inf(-1)))
	for _, x := range expInputs {
		got, want := FastExpMinimax(x), math.Exp(x)
		// Relative error is not a useful bound after exp enters the subnormal range.
		if want > 0 && want < math.Float64frombits(0x0010000000000000) {
			if got != want && ulpDistance64(got, want) > 2 {
				t.Errorf("FastExpMinimax(%g) = %g, want %g (%d ULP)", x, got, want, ulpDistance64(got, want))
			}
			continue
		}
		check64("FastExpMinimax", got, want, x, 0, exp64RelBudget)
		if x >= -700 && x <= 700 {
			check64("FastExpBitCast", FastExpBitCast(x), want, x, 0, 0.05)
		}
	}

	rng := rand.New(rand.NewSource(43))
	for i := 0; i < 32768; i++ {
		// Cover every normal exponent and a stratified selection of mantissas.
		exp := rng.Intn(2046) + 1
		mantissa := rng.Uint64() & ((uint64(1) << 52) - 1)
		y := math.Float64frombits(uint64(exp)<<52 | mantissa)
		if i%2 == 0 {
			y = math.Float64frombits(uint64(exp)<<52 | mantissa | uint64(1)<<63)
		}
		// Log is defined for positive inputs; the sign bit is intentionally cleared.
		y = math.Float64frombits(math.Float64bits(y) &^ (uint64(1) << 63))
		want := math.Log(y)
		check64("FastLogMinimax", FastLogMinimax(y), want, y, log64AbsBudget, 0)
		check64("FastLogChebyshev", FastLogChebyshev(y), want, y, 5e-4, 0)
	}

	for k := -1074; k <= 1023; k++ {
		y := math.Ldexp(1, k)
		if y == 0 || math.IsInf(y, 0) {
			continue
		}
		want := float64(k) * math.Ln2
		check64("FastLogMinimax power of two", FastLogMinimax(y), want, y, log64AbsBudget, 0)
	}
	for i := 0; i < 8192; i++ {
		mantissa := rng.Uint64() & ((uint64(1) << 52) - 1)
		if mantissa == 0 {
			mantissa = 1
		}
		y := math.Float64frombits(mantissa)
		want := math.Log(y*0x1p52) - 52*math.Ln2
		check64("FastLogMinimax subnormal", FastLogMinimax(y), want, y, log64AbsBudget, 0)
	}

	for i := 0; i < 32768; i++ {
		x := float32(rng.Float64()*192 - 103)
		want := float32(math.Exp(float64(x)))
		got := FastExpF32(x)
		if float64(want) >= float64(math.Float32frombits(0x00800000)) {
			check32("FastExpF32", got, want, x, 0, exp32RelBudget)
		} else if want == 0 {
			if got != 0 {
				t.Errorf("FastExpF32(%g) = %g, want 0", x, got)
			}
		} else if got != want && ulpDistance32(got, want) > 20 {
			t.Errorf("FastExpF32(%g) = %g, want %g (%d ULP)", x, got, want, ulpDistance32(got, want))
		}

		y := math.Float32frombits(rng.Uint32() & 0x7fffffff)
		if y == 0 || math.IsInf(float64(y), 0) || math.IsNaN(float64(y)) {
			continue
		}
		wantLog := float32(math.Log(float64(y)))
		check32("FastLogF32", FastLogF32(y), wantLog, y, log32AbsBudget, 0)
	}
	for k := -149; k <= 128; k++ {
		x := float32(float64(k) * math.Ln2)
		for _, boundary := range []float32{
			math.Nextafter32(x, float32(math.Inf(-1))),
			x,
			math.Nextafter32(x, float32(math.Inf(1))),
		} {
			want := float32(math.Exp(float64(boundary)))
			got := FastExpF32(boundary)
			if want >= math.Float32frombits(0x00800000) {
				check32("FastExpF32 range boundary", got, want, boundary, 0, exp32RelBudget)
			} else if want == 0 {
				if got != 0 {
					t.Errorf("FastExpF32(%g) = %g, want 0", boundary, got)
				}
			} else if got != want && ulpDistance32(got, want) > 20 {
				t.Errorf("FastExpF32(%g) = %g, want %g (%d ULP)", boundary, got, want, ulpDistance32(got, want))
			}
		}
		for name, sample := range worst {
			t.Logf("%s worst sampled error %.3g at %.17g (budget %.3g, %.3g of budget)",
				name, sample.error, sample.input, sample.budget, sample.ratio)
		}
	}
}

func TestGuardrails(t *testing.T) {
	// LnRegularized at zero should be finite and equal to ln(eps)
	eps := 1e-6
	valAtZero := LnRegularized(0, eps)
	wantAtZero := math.Log(eps)
	if math.Abs(valAtZero-wantAtZero) > 1e-4 {
		t.Errorf("LnRegularized(0) = %v, want %v", valAtZero, wantAtZero)
	}

	// Negative values should be smooth and positive
	negVal := LnRegularized(-5.0, eps)
	posVal := LnRegularized(5.0, eps)
	if math.Abs(negVal-posVal) > 1e-6 {
		t.Errorf("LnRegularized should be symmetric: f(-5)=%v, f(5)=%v", negVal, posVal)
	}

	// Clamped Exp should not overflow to +Inf
	hugeExp := ExpClamped(1000.0)
	if math.IsInf(hugeExp, 1) || math.IsNaN(hugeExp) {
		t.Errorf("ExpClamped(1000) returned Inf or NaN: %v", hugeExp)
	}

	// FastEmlRegularized should never produce NaN across negative domain
	for _, y := range []float64{-100.0, -1.0, -0.001, 0.0, 0.001, 1.0, 100.0} {
		res := FastEmlRegularized(1.0, y, DefaultEpsilon)
		if math.IsNaN(res) || math.IsInf(res, 0) {
			t.Errorf("FastEmlRegularized(1.0, %v) produced non-finite: %v", y, res)
		}
	}
}

func BenchmarkFastExp(b *testing.B) {
	x := 1.234
	for b.Loop() {
		_ = FastExp(x)
	}
}

func BenchmarkStdExp(b *testing.B) {
	x := 1.234
	for b.Loop() {
		_ = math.Exp(x)
	}
}

func BenchmarkFastLog(b *testing.B) {
	y := 2.345
	for b.Loop() {
		_ = FastLog(y)
	}
}

func BenchmarkStdLog(b *testing.B) {
	y := 2.345
	for b.Loop() {
		_ = math.Log(y)
	}
}
