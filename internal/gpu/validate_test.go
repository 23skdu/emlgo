package gpu

import (
	"math"
	"strings"
	"testing"
)

func TestDefaultVerifier(t *testing.T) {
	v := DefaultVerifier()
	if v.MaxULP != 1 {
		t.Errorf("MaxULP = %d, want 1", v.MaxULP)
	}
}

func TestVerifyOpExactMatch(t *testing.T) {
	v := DefaultVerifier()
	input := []float64{1.0, 2.0, 3.0}
	result := []float64{1.0, 2.0, 3.0}
	ref := func(x float64) float64 { return x }

	maxULP, failed, err := v.VerifyOp("test", input, result, ref)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if maxULP != 0 {
		t.Errorf("maxULP = %d, want 0", maxULP)
	}
	if failed != 0 {
		t.Errorf("failed = %d, want 0", failed)
	}
}

func TestVerifyOpLengthMismatch(t *testing.T) {
	v := DefaultVerifier()
	_, _, err := v.VerifyOp("test", []float64{1.0}, []float64{1.0, 2.0}, math.Exp)
	if err == nil {
		t.Fatal("expected error for length mismatch")
	}
	if !strings.Contains(err.Error(), "length") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifyOpSmallULP(t *testing.T) {
	v := DefaultVerifier()
	input := []float64{1.0}
	// Slightly offset to produce a measurable ULP difference
	result := []float64{1.0 + 4.44e-16} // ~2 ULP for values near 1.0
	ref := func(x float64) float64 { return x }

	maxULP, failed, err := v.VerifyOp("test", input, result, ref)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if maxULP == 0 {
		t.Error("expected non-zero ULP")
	}
	if failed == 0 {
		t.Log("ULP difference within tolerance (MaxULP=1)")
	}
}

func TestVerifyOpLargeULP(t *testing.T) {
	v := &BatchVerifier{MaxULP: 100}
	input := []float64{1.0}
	result := []float64{2.0} // Deliberately wrong
	ref := func(x float64) float64 { return 1.0 }

	maxULP, failed, err := v.VerifyOp("test", input, result, ref)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if maxULP == 0 {
		t.Error("expected large ULP for deliberately wrong result")
	}
	if failed == 0 {
		t.Error("expected failures for deliberately wrong result")
	}
	t.Logf("ULP diff for 1.0 vs 2.0: %d", maxULP)
}

func TestULPDiffIdentical(t *testing.T) {
	if d := ulpDiff(1.0, 1.0); d != 0 {
		t.Errorf("ulpDiff(1.0, 1.0) = %d, want 0", d)
	}
}

func TestULPDiffZero(t *testing.T) {
	if d := ulpDiff(0.0, 0.0); d != 0 {
		t.Errorf("ulpDiff(0, 0) = %d, want 0", d)
	}
	if d := ulpDiff(0.0, -0.0); d != 0 {
		t.Errorf("ulpDiff(0, -0) = %d, want 0", d)
	}
}

func TestULPDiffNaN(t *testing.T) {
	if d := ulpDiff(math.NaN(), math.NaN()); d != 0 {
		t.Errorf("ulpDiff(NaN, NaN) = %d, want 0", d)
	}
	if d := ulpDiff(math.NaN(), 1.0); d != 0 {
		t.Errorf("ulpDiff(NaN, 1.0) = %d, want 0", d)
	}
}

func TestULPDiffInf(t *testing.T) {
	if d := ulpDiff(math.Inf(1), math.Inf(1)); d != 0 {
		t.Errorf("ulpDiff(+Inf, +Inf) = %d, want 0", d)
	}
	if d := ulpDiff(math.Inf(-1), math.Inf(-1)); d != 0 {
		t.Errorf("ulpDiff(-Inf, -Inf) = %d, want 0", d)
	}
	if d := ulpDiff(math.Inf(1), math.Inf(-1)); d != 0 {
		t.Errorf("ulpDiff(+Inf, -Inf) = %d, want 0", d)
	}
	if d := ulpDiff(math.Inf(1), 1.0); d != 0 {
		t.Errorf("ulpDiff(+Inf, 1.0) = %d, want 0", d)
	}
}

func TestULPDiffSmall(t *testing.T) {
	a := 1.0
	b := math.Nextafter(1.0, 2.0) // 1 ULP larger
	if d := ulpDiff(a, b); d != 1 {
		t.Errorf("ulpDiff(%v, %v) = %d, want 1", a, b, d)
	}
}

func TestULPDiffReverse(t *testing.T) {
	a := 1.0
	b := math.Nextafter(1.0, 2.0)
	if d := ulpDiff(b, a); d != 1 {
		t.Errorf("ulpDiff(%v, %v) = %d, want 1 (symmetric)", b, a, d)
	}
}

func TestULPDiffLarge(t *testing.T) {
	a := 0.0
	b := math.MaxFloat64
	d := ulpDiff(a, b)
	if d == 0 {
		t.Error("expected non-zero ULP between 0 and MaxFloat64")
	}
	t.Logf("ulpDiff(0, MaxFloat64) = %d", d)
}

func TestULPDiffSubnormal(t *testing.T) {
	a := math.SmallestNonzeroFloat64
	b := math.Nextafter(math.SmallestNonzeroFloat64, 0)
	d := ulpDiff(a, b)
	if d != 1 {
		t.Errorf("ulpDiff(SmallestNonzero, nextafter) = %d, want 1", d)
	}
}

func TestCPURefs(t *testing.T) {
	names := []string{"Exp", "Log", "Sin", "Cos", "Tan", "Sinh", "Cosh", "Tanh", "Sqrt"}
	for _, name := range names {
		ref, ok := cpuRefs[name]
		if !ok {
			t.Errorf("cpuRefs missing entry for %s", name)
			continue
		}
		if ref == nil {
			t.Errorf("cpuRefs[%s] is nil", name)
		}
	}
}

func TestCPURefsExp(t *testing.T) {
	ref := cpuRefs["Exp"]
	expected := math.Exp(1.0)
	if ref(1.0) != expected {
		t.Errorf("Exp(1) = %v, want %v", ref(1.0), expected)
	}
}

func TestCPURefsLog(t *testing.T) {
	ref := cpuRefs["Log"]
	expected := math.Log(2.0)
	if ref(2.0) != expected {
		t.Errorf("Log(2) = %v, want %v", ref(2.0), expected)
	}
}

func TestCPURefsSin(t *testing.T) {
	ref := cpuRefs["Sin"]
	expected := math.Sin(0.5)
	if ref(0.5) != expected {
		t.Errorf("Sin(0.5) = %v, want %v", ref(0.5), expected)
	}
}

func TestCPURefsCos(t *testing.T) {
	ref := cpuRefs["Cos"]
	expected := math.Cos(0.5)
	if ref(0.5) != expected {
		t.Errorf("Cos(0.5) = %v, want %v", ref(0.5), expected)
	}
}

func TestCPURefsTan(t *testing.T) {
	ref := cpuRefs["Tan"]
	expected := math.Tan(0.5)
	if ref(0.5) != expected {
		t.Errorf("Tan(0.5) = %v, want %v", ref(0.5), expected)
	}
}

func TestCPURefsSinh(t *testing.T) {
	ref := cpuRefs["Sinh"]
	expected := math.Sinh(0.5)
	if ref(0.5) != expected {
		t.Errorf("Sinh(0.5) = %v, want %v", ref(0.5), expected)
	}
}

func TestCPURefsCosh(t *testing.T) {
	ref := cpuRefs["Cosh"]
	expected := math.Cosh(0.5)
	if ref(0.5) != expected {
		t.Errorf("Cosh(0.5) = %v, want %v", ref(0.5), expected)
	}
}

func TestCPURefsTanh(t *testing.T) {
	ref := cpuRefs["Tanh"]
	expected := math.Tanh(0.5)
	if ref(0.5) != expected {
		t.Errorf("Tanh(0.5) = %v, want %v", ref(0.5), expected)
	}
}

func TestCPURefsSqrt(t *testing.T) {
	ref := cpuRefs["Sqrt"]
	expected := math.Sqrt(4.0)
	if ref(4.0) != expected {
		t.Errorf("Sqrt(4) = %v, want %v", ref(4.0), expected)
	}
}

func TestCPURefsAll(t *testing.T) {
	x := 1.5
	expected := map[string]float64{
		"Exp":  math.Exp(x),
		"Log":  math.Log(x),
		"Sin":  math.Sin(x),
		"Cos":  math.Cos(x),
		"Tan":  math.Tan(x),
		"Sinh": math.Sinh(x),
		"Cosh": math.Cosh(x),
		"Tanh": math.Tanh(x),
		"Sqrt": math.Sqrt(x),
	}
	for name, ref := range cpuRefs {
		got := ref(x)
		want := expected[name]
		if got != want {
			t.Errorf("cpuRefs[%s](%v) = %v, want %v", name, x, got, want)
		}
	}
}

func BenchmarkULPDiff(b *testing.B) {
	a, c := 1.0, math.Nextafter(1.0, 2.0)
	for i := 0; i < b.N; i++ {
		ulpDiff(a, c)
	}
}

func BenchmarkVerifyOp(b *testing.B) {
	v := DefaultVerifier()
	input := make([]float64, 1000)
	result := make([]float64, 1000)
	for i := range input {
		input[i] = float64(i)
		result[i] = float64(i)
	}
	ref := func(x float64) float64 { return x }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = v.VerifyOp("bench", input, result, ref)
	}
}
