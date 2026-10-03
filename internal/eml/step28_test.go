package eml

import (
	"math"
	"math/cmplx"
	"testing"
)

// TestStep28_PipelineStepFusion verifies that consecutive affine steps
// (MulScalar, AddScalar, SubScalar, DivScalar, Neg) are fused into a single
// execution step, reducing p.Len() and accelerating throughput.
func TestStep28_PipelineStepFusion(t *testing.T) {
	p := NewPipeline(16)

	// Mul(2) then Add(3) -> 1 fused step (y = 2x + 3)
	p.MulScalar(2).AddScalar(3)
	if p.Len() != 1 {
		t.Fatalf("expected 1 fused step, got %d", p.Len())
	}

	input := []float64{0, 1, 2, 5, 10}
	output := make([]float64, len(input))
	p.RunTo(input, output)

	for i, x := range input {
		want := 2*x + 3
		if math.Abs(output[i]-want) > 1e-12 {
			t.Errorf("[%d]: got %v, want %v", i, output[i], want)
		}
	}

	// Now add Neg and SubScalar(1) -> still 1 fused step (y = -(2x+3) - 1 = -2x - 4)
	p.Neg().SubScalar(1)
	if p.Len() != 1 {
		t.Fatalf("expected 1 fused step after chaining, got %d", p.Len())
	}
	p.RunTo(input, output)
	for i, x := range input {
		want := -2*x - 4
		if math.Abs(output[i]-want) > 1e-12 {
			t.Errorf("chained [%d]: got %v, want %v", i, output[i], want)
		}
	}
}

// TestStep28_PipelineRunToLargeSIMD verifies Pipeline.Exp().MulScalar(2).Log()
// on n = 65,536 elements matches x + ln(2).
func TestStep28_PipelineRunToLargeSIMD(t *testing.T) {
	const n = 65536
	input := make([]float64, n)
	for i := range input {
		input[i] = float64(i%100) * 0.05
	}
	output := make([]float64, n)

	p := NewPipeline(n)
	p.Exp().MulScalar(2).Log()
	p.RunTo(input, output)

	ln2 := math.Log(2)
	for i, x := range input {
		want := x + ln2
		if math.Abs(output[i]-want) > 1e-10 {
			t.Fatalf("at index %d: got %v, want %v", i, output[i], want)
		}
	}
}

// TestStep28_ComplexDotProductCompensated verifies Neumaier compensated summation
// preserves precision in the presence of large cancelling magnitudes.
func TestStep28_ComplexDotProductCompensated(t *testing.T) {
	// Construct an array where standard floating-point addition loses precision:
	// a[0]*conj(b[0]) = 1e16 + 0i
	// a[1..10000]*conj(b[1..10000]) = 1.0 + 1.0i
	// a[10001]*conj(b[10001]) = -1e16 + 0i
	// Exact mathematical sum = 10000 + 10000i.
	// Standard sum: (1e16 + 1.0) rounds to 1e16, so the 10000 small terms are obliterated!
	const count = 10000
	a := make([]complex128, count+2)
	b := make([]complex128, count+2)

	a[0] = complex(1e15, 0)
	b[0] = complex(1, 0)

	for i := 1; i <= count; i++ {
		a[i] = complex(1, 1)
		b[i] = complex(1, 0) // a[i]*conj(b[i]) = 1 + 1i
	}

	a[count+1] = complex(-1e15, 0)
	b[count+1] = complex(1, 0)

	res := ComplexDotProduct(a, b)

	if math.Abs(real(res)-float64(count)) > 1e-5 {
		t.Errorf("real part compensated error: got %v, want %v", real(res), count)
	}
	if math.Abs(imag(res)-float64(count)) > 1e-5 {
		t.Errorf("imag part compensated error: got %v, want %v", imag(res), count)
	}
}

// BenchmarkStep28_PipelineExpMulLarge benchmarks Pipeline.Exp().MulScalar(2).RunTo at n = 65,536.
func BenchmarkStep28_PipelineExpMulLarge(b *testing.B) {
	const n = 65536
	input := make([]float64, n)
	for i := range input {
		input[i] = float64(i%50) * 0.02
	}
	output := make([]float64, n)

	p := NewPipeline(n)
	p.Exp().MulScalar(2)

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		p.RunTo(input, output)
	}
}

// BenchmarkStep28_ComplexDotProductLarge benchmarks ComplexDotProduct at n = 1,000,000.
func BenchmarkStep28_ComplexDotProductLarge(b *testing.B) {
	const n = 1000000
	a := make([]complex128, n)
	bSlice := make([]complex128, n)
	for i := range a {
		a[i] = complex(float64(i)*0.001, float64(i)*0.002)
		bSlice[i] = complex(1.0, 0.5)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = ComplexDotProduct(a, bSlice)
	}
}

// FuzzPipelineDifferential fuzzes random values through a multi-step Pipeline
// comparing against scalar math evaluations.
func FuzzPipelineDifferential(f *testing.F) {
	f.Add(0.5, 2.0, 1.5)
	f.Add(1.0, -1.0, 0.0)
	f.Add(10.0, 0.5, -2.5)

	f.Fuzz(func(t *testing.T, x0, scale, bias float64) {
		if math.IsNaN(x0) || math.IsInf(x0, 0) || x0 < -50 || x0 > 50 {
			return
		}
		if math.IsNaN(scale) || math.IsInf(scale, 0) || math.Abs(scale) > 100 {
			return
		}
		if math.IsNaN(bias) || math.IsInf(bias, 0) || math.Abs(bias) > 100 {
			return
		}

		in := []float64{x0}
		out := make([]float64, 1)

		p := NewPipeline(1)
		p.Exp().MulScalar(scale).AddScalar(bias)
		p.RunTo(in, out)

		want := math.Exp(x0)*scale + bias
		if math.IsNaN(want) || math.IsInf(want, 0) {
			return
		}

		if math.Abs(out[0]-want) > 1e-9*(math.Abs(want)+1.0) {
			t.Errorf("pipeline differential mismatch: got %v, want %v", out[0], want)
		}
	})
}

// FuzzComplexDotProduct verifies Hermitian dot product consistency.
func FuzzComplexDotProduct(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0)
	f.Add(-1.5, 0.5, 2.0, -3.0)
	f.Add(0.0, 0.0, 1.0, 1.0)

	f.Fuzz(func(t *testing.T, ar, ai, br, bi float64) {
		if math.IsNaN(ar) || math.IsNaN(ai) || math.IsNaN(br) || math.IsNaN(bi) {
			return
		}
		if math.IsInf(ar, 0) || math.IsInf(ai, 0) || math.IsInf(br, 0) || math.IsInf(bi, 0) {
			return
		}
		if math.Abs(ar) > 1e6 || math.Abs(ai) > 1e6 || math.Abs(br) > 1e6 || math.Abs(bi) > 1e6 {
			return
		}

		a := []complex128{complex(ar, ai)}
		b := []complex128{complex(br, bi)}

		got := ComplexDotProduct(a, b)
		want := a[0] * cmplx.Conj(b[0])

		if math.Abs(real(got)-real(want)) > 1e-10*(math.Abs(real(want))+1.0) {
			t.Errorf("real mismatch: got %v, want %v", real(got), real(want))
		}
		if math.Abs(imag(got)-imag(want)) > 1e-10*(math.Abs(imag(want))+1.0) {
			t.Errorf("imag mismatch: got %v, want %v", imag(got), imag(want))
		}
	})
}
