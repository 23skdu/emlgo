package eml

import (
	"math"
	"math/cmplx"
	"testing"
)

const complexTolerance = 1e-10

func complexEq(a, b complex128) bool {
	return cmplx.Abs(a-b) < complexTolerance
}

func TestComplexBatchBasic(t *testing.T) {
	x := []complex128{1, 2, 3}
	y := []complex128{1, 1, 1}
	result := make([]complex128, 3)
	ComplexBatch(x, y, result)

	for i, v := range x {
		expected := cmplx.Exp(v) - cmplx.Log(y[i])
		if !complexEq(result[i], expected) {
			t.Errorf("ComplexBatch[%d]: got %v, want %v", i, result[i], expected)
		}
	}
}

func TestComplexExpBatch(t *testing.T) {
	input := []complex128{0, 1, complex(0, math.Pi)}
	got := ComplexExpBatch(input)

	for i, v := range input {
		want := cmplx.Exp(v)
		if !complexEq(got[i], want) {
			t.Errorf("ComplexExpBatch[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestComplexLogBatch(t *testing.T) {
	input := []complex128{1, complex(math.E, 0), complex(0, 1)}
	got := ComplexLogBatch(input)

	for i, v := range input {
		want := cmplx.Log(v)
		if !complexEq(got[i], want) {
			t.Errorf("ComplexLogBatch[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestComplexSinBatch(t *testing.T) {
	input := []complex128{0, 1, complex(0, 1)}
	got := ComplexSinBatch(input)

	for i, v := range input {
		want := cmplx.Sin(v)
		if !complexEq(got[i], want) {
			t.Errorf("ComplexSinBatch[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestComplexCosBatch(t *testing.T) {
	input := []complex128{0, 1, complex(0, 1)}
	got := ComplexCosBatch(input)

	for i, v := range input {
		want := cmplx.Cos(v)
		if !complexEq(got[i], want) {
			t.Errorf("ComplexCosBatch[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestComplexTanBatch(t *testing.T) {
	input := []complex128{0, 1, complex(0, 1)}
	got := ComplexTanBatch(input)

	for i, v := range input {
		want := cmplx.Tan(v)
		if !complexEq(got[i], want) {
			t.Errorf("ComplexTanBatch[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestComplexTrigIdentity(t *testing.T) {
	// sin²(z) + cos²(z) = 1
	z := []complex128{
		1,
		complex(2, 3),
		complex(-1, 4),
		complex(0.5, -0.5),
	}
	sin := ComplexSinBatch(z)
	cos := ComplexCosBatch(z)

	for i := range z {
		identity := sin[i]*sin[i] + cos[i]*cos[i]
		if !complexEq(identity, 1) {
			t.Errorf("sin²+cos² identity failed for z=%v: got %v, want 1", z[i], identity)
		}
	}
}

func TestComplexEulerIdentity(t *testing.T) {
	// e^(iπ) + 1 = 0
	piI := complex(0, math.Pi)
	got := cmplx.Exp(piI) + 1
	if cmplx.Abs(got) > complexTolerance {
		t.Errorf("Euler's identity: e^(iπ)+1 = %v, want 0", got)
	}
}

func TestComplexSqrtSquare(t *testing.T) {
	// sqrt(z²) = z or sqrt(z²) = -z (principal branch)
	z := []complex128{
		1,
		complex(2, 3),
		complex(0.5, -0.5),
	}

	for _, v := range z {
		got := cmplx.Sqrt(v * v)
		if !complexEq(got, v) && !complexEq(got, -v) {
			t.Errorf("sqrt(z²) failed for z=%v: got %v", v, got)
		}
	}
}

func TestComplexSqrtSquareBranchCut(t *testing.T) {
	// For Re(z) < 0, sqrt(z²) = -z (principal branch)
	z := []complex128{
		complex(-1, 4),
		complex(-2, -3),
	}

	for _, v := range z {
		got := cmplx.Sqrt(v * v)
		if !complexEq(got, -v) {
			t.Errorf("sqrt(z²) branch cut failed for z=%v: got %v, want %v", v, got, -v)
		}
	}
}

func TestComplexBatchPanicLengthMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for length mismatch, did not get one")
		}
	}()
	x := []complex128{1, 2}
	y := []complex128{1}
	result := make([]complex128, 2)
	ComplexBatch(x, y, result)
}

func TestComplexBatchEmpty(t *testing.T) {
	x := []complex128{}
	y := []complex128{}
	result := []complex128{}
	ComplexBatch(x, y, result)
	// no panic = pass
}

func TestComplexBatchSingleElement(t *testing.T) {
	x := []complex128{complex(1, 2)}
	y := []complex128{complex(3, 4)}
	result := make([]complex128, 1)
	ComplexBatch(x, y, result)

	expected := cmplx.Exp(x[0]) - cmplx.Log(y[0])
	if !complexEq(result[0], expected) {
		t.Errorf("got %v, want %v", result[0], expected)
	}
}
