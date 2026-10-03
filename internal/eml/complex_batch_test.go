package eml

import (
	"math"
	"math/cmplx"
	"testing"
)

// assertPanics runs fn and fails if it does not panic.
func assertPanics(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: expected panic, got none", name)
		}
	}()
	fn()
}

func complex128Samples() []complex128 {
	return []complex128{
		1 + 2i,
		-0.5 + 0.25i,
		3,
		-2i,
		0,
		1e-3 - 1e-3i,
	}
}

func complex64Samples() []complex64 {
	return []complex64{
		1 + 2i,
		-0.5 + 0.25i,
		3,
		-2i,
		0,
	}
}

// TestComplexBatchToVariants covers the in-place *BatchTo entry points for
// complex128 and checks each against its allocating counterpart.
func TestComplexBatchToVariants(t *testing.T) {
	x := complex128Samples()
	y := complex128Samples()
	dst := make([]complex128, len(x))

	ComplexExpBatchTo(x, dst)
	want := ComplexExpBatch(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexExpBatchTo[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexLogBatchTo(x, dst)
	want = ComplexLogBatch(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexLogBatchTo[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexSinBatchTo(x, dst)
	want = ComplexSinBatch(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexSinBatchTo[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexCosBatchTo(x, dst)
	want = ComplexCosBatch(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexCosBatchTo[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexTanBatchTo(x, dst)
	want = ComplexTanBatch(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexTanBatchTo[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexSqrtBatchTo(x, dst)
	want = ComplexSqrtBatch(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexSqrtBatchTo[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	got := make([]complex128, len(y))
	ComplexBatch(x, y, got)
	for i := range got {
		want := cmplx.Exp(x[i]) - cmplx.Log(y[i])
		if got[i] != want {
			t.Errorf("ComplexBatch[%d] = %v, want %v", i, got[i], want)
		}
	}
}

// TestComplexBatchToC64Variants covers the complex64 in-place variants.
func TestComplexBatchToC64Variants(t *testing.T) {
	x := complex64Samples()
	dst := make([]complex64, len(x))

	ComplexExpBatchToC64(x, dst)
	want := ComplexExpBatchC64(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexExpBatchToC64[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexLogBatchToC64(x, dst)
	want = ComplexLogBatchC64(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexLogBatchToC64[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexSinBatchToC64(x, dst)
	want = ComplexSinBatchC64(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexSinBatchToC64[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexCosBatchToC64(x, dst)
	want = ComplexCosBatchC64(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexCosBatchToC64[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexTanBatchToC64(x, dst)
	want = ComplexTanBatchC64(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexTanBatchToC64[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	ComplexSqrtBatchToC64(x, dst)
	want = ComplexSqrtBatchC64(x)
	for i := range want {
		if dst[i] != want[i] {
			t.Errorf("ComplexSqrtBatchToC64[%d] = %v, want %v", i, dst[i], want[i])
		}
	}

	got := make([]complex64, len(x))
	ComplexBatchC64(x, x, got)
	for i := range got {
		if want := Complex64(x[i], x[i]); got[i] != want {
			t.Errorf("ComplexBatchC64[%d] = %v, want %v", i, got[i], want)
		}
	}
}

// TestComplexBatchLengthMismatch verifies every complex batch entry point
// rejects mismatched slice lengths instead of reading out of bounds.
func TestComplexBatchLengthMismatch(t *testing.T) {
	n := 4
	x128 := make([]complex128, n)
	short128 := make([]complex128, n-1)
	x64 := make([]complex64, n)
	short64 := make([]complex64, n-1)

	cases := []struct {
		name string
		fn   func()
	}{
		{"ComplexBatch", func() { ComplexBatch(x128, short128, make([]complex128, n)) }},
		{"ComplexBatchC64", func() { ComplexBatchC64(x64, short64, make([]complex64, n)) }},
		{"ComplexExpBatchTo", func() { ComplexExpBatchTo(x128, short128) }},
		{"ComplexLogBatchTo", func() { ComplexLogBatchTo(x128, short128) }},
		{"ComplexSinBatchTo", func() { ComplexSinBatchTo(x128, short128) }},
		{"ComplexCosBatchTo", func() { ComplexCosBatchTo(x128, short128) }},
		{"ComplexTanBatchTo", func() { ComplexTanBatchTo(x128, short128) }},
		{"ComplexSqrtBatchTo", func() { ComplexSqrtBatchTo(x128, short128) }},
		{"ComplexExpBatchToC64", func() { ComplexExpBatchToC64(x64, short64) }},
		{"ComplexLogBatchToC64", func() { ComplexLogBatchToC64(x64, short64) }},
		{"ComplexSinBatchToC64", func() { ComplexSinBatchToC64(x64, short64) }},
		{"ComplexCosBatchToC64", func() { ComplexCosBatchToC64(x64, short64) }},
		{"ComplexTanBatchToC64", func() { ComplexTanBatchToC64(x64, short64) }},
		{"ComplexSqrtBatchToC64", func() { ComplexSqrtBatchToC64(x64, short64) }},
		{"ComplexDotProduct", func() { ComplexDotProduct(x128, short128) }},
		{"ComplexDotProductC64", func() { ComplexDotProductC64(x64, short64) }},
	}
	for _, c := range cases {
		assertPanics(t, c.name, c.fn)
	}
}

// TestComplexDotProduct checks the Hermitian inner product sum(a[i]*conj(b[i])).
func TestComplexDotProduct(t *testing.T) {
	a := []complex128{1 + 2i, 3 - 1i}
	b := []complex128{4 + 0.5i, -2 + 3i}

	got := ComplexDotProduct(a, b)
	want := a[0]*cmplx.Conj(b[0]) + a[1]*cmplx.Conj(b[1])
	if got != want {
		t.Errorf("ComplexDotProduct = %v, want %v", got, want)
	}

	// Conjugating the first argument negates the result.
	if flipped := ComplexDotProduct(b, a); flipped != cmplx.Conj(got) {
		t.Errorf("ComplexDotProduct(b,a) = %v, want %v", flipped, cmplx.Conj(got))
	}

	// Empty input yields zero.
	if v := ComplexDotProduct(nil, nil); v != 0 {
		t.Errorf("ComplexDotProduct(nil,nil) = %v, want 0", v)
	}

	a64 := []complex64{complex64(a[0]), complex64(a[1])}
	b64 := []complex64{complex64(b[0]), complex64(b[1])}
	got64 := ComplexDotProductC64(a64, b64)
	if diff := cmplx.Abs(complex128(got64) - want); diff > 1e-6 {
		t.Errorf("ComplexDotProductC64 = %v, want ~%v (diff %g)", got64, want, diff)
	}
	if v := ComplexDotProductC64(nil, nil); v != 0 {
		t.Errorf("ComplexDotProductC64(nil,nil) = %v, want 0", v)
	}
}

// TestComplexBatchEmptyInputs verifies zero-length slices are accepted.
func TestComplexBatchEmptyInputs(t *testing.T) {
	if got := ComplexExpBatch(nil); len(got) != 0 {
		t.Errorf("ComplexExpBatch(nil) len = %d, want 0", len(got))
	}
	if got := ComplexExpBatchC64(nil); len(got) != 0 {
		t.Errorf("ComplexExpBatchC64(nil) len = %d, want 0", len(got))
	}
	if got := ComplexLogBatch(nil); len(got) != 0 {
		t.Errorf("ComplexLogBatch(nil) len = %d, want 0", len(got))
	}
	if got := ComplexSqrtBatchC64(nil); len(got) != 0 {
		t.Errorf("ComplexSqrtBatchC64(nil) len = %d, want 0", len(got))
	}
	dst := make([]complex128, 0)
	ComplexExpBatchTo(nil, dst)
	ComplexBatch(nil, nil, nil)
}

// TestComplexTrigIdentitiesBatch re-checks sin^2(z)+cos^2(z)=1 over a batch.
func TestComplexTrigIdentitiesBatch(t *testing.T) {
	x := complex128Samples()
	sin := ComplexSinBatch(x)
	cos := ComplexCosBatch(x)
	for i := range x {
		got := sin[i]*sin[i] + cos[i]*cos[i]
		if cmplx.Abs(got-1) > 1e-12 {
			t.Errorf("x=%v: sin^2+cos^2 = %v, want 1", x[i], got)
		}
	}
}

// TestComplexSqrtBatchBranchCut checks sqrt(z)^2 == z across both square-root
// branch cuts used by cmplx.Sqrt.
func TestComplexSqrtBatchBranchCut(t *testing.T) {
	x := []complex128{-1, 1, -1e-8 + 1e-8i, -1e-8 - 1e-8i, 1e-12 + 0i}
	got := ComplexSqrtBatch(x)
	for i, z := range x {
		if diff := cmplx.Abs(got[i]*got[i] - z); diff > 1e-14 {
			t.Errorf("sqrt(%v)^2 = %v, want %v (diff %g)", z, got[i]*got[i], z, diff)
		}
	}
}

// TestComplexBatchToOverwrites verifies the in-place variants fully overwrite
// the destination rather than accumulating into it.
func TestComplexBatchToOverwrites(t *testing.T) {
	// NaN is not produced by any of these operations for the sample inputs, so
	// it is a safe "not yet written" sentinel.
	x := complex128Samples()
	dst := make([]complex128, len(x))
	for i := range dst {
		dst[i] = complex(math.NaN(), math.NaN())
	}
	ComplexExpBatchTo(x, dst)
	for i := range dst {
		if math.IsNaN(real(dst[i])) {
			t.Errorf("dst[%d] not overwritten: %v", i, dst[i])
		}
	}

	x64 := complex64Samples()
	dst64 := make([]complex64, len(x64))
	for i := range dst64 {
		dst64[i] = complex64(complex(math.NaN(), math.NaN()))
	}
	ComplexLogBatchToC64(x64, dst64)
	for i := range dst64 {
		if math.IsNaN(float64(real(dst64[i]))) {
			t.Errorf("dst64[%d] not overwritten: %v", i, dst64[i])
		}
	}
}
