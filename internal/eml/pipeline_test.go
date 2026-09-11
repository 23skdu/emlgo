package eml

import (
	"math"
	"testing"
)

func TestPipelineExpSingle(t *testing.T) {
	p := NewPipeline(4)
	p.Exp()
	input := []float64{0, 1, -1, 2}
	got := p.Run(input)
	for i, v := range input {
		want := math.Exp(v)
		if math.Abs(got[i]-want) > 1e-12 {
			t.Errorf("Exp()[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestPipelineChainExpMulLog(t *testing.T) {
	// exp(x) * 2 → log → should give x + log(2)
	p := NewPipeline(4)
	p.Exp().MulScalar(2).Log()
	input := []float64{0, 0.5, 1, 2}
	got := p.Run(input)
	for i, v := range input {
		want := v + math.Log(2)
		if math.Abs(got[i]-want) > 1e-12 {
			t.Errorf("chain[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestPipelineAddScalar(t *testing.T) {
	p := NewPipeline(3)
	p.AddScalar(10)
	got := p.Run([]float64{1, 2, 3})
	for i, want := range []float64{11, 12, 13} {
		if got[i] != want {
			t.Errorf("AddScalar[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestPipelineNeg(t *testing.T) {
	p := NewPipeline(3)
	p.Neg()
	got := p.Run([]float64{1, -2, 3})
	want := []float64{-1, 2, -3}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Neg[%d]: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestPipelineAbs(t *testing.T) {
	p := NewPipeline(3)
	p.Abs()
	got := p.Run([]float64{-1, 2, -3})
	for i, want := range []float64{1, 2, 3} {
		if got[i] != want {
			t.Errorf("Abs[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestPipelineSin(t *testing.T) {
	p := NewPipeline(3)
	p.Sin()
	input := []float64{0, math.Pi / 2, math.Pi}
	got := p.Run(input)
	for i, v := range input {
		want := math.Sin(v)
		if math.Abs(got[i]-want) > 1e-12 {
			t.Errorf("Sin[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func TestPipelineRunTo(t *testing.T) {
	p := NewPipeline(4)
	p.Exp().MulScalar(2)
	input := []float64{0, 1, 2, 3}
	output := make([]float64, 4)
	p.RunTo(input, output)
	for i, v := range input {
		want := math.Exp(v) * 2
		if math.Abs(output[i]-want) > 1e-12 {
			t.Errorf("RunTo[%d]: got %v, want %v", i, output[i], want)
		}
	}
}

func TestPipelineGrowsBuffer(t *testing.T) {
	// Start with n=2, then run with larger input.
	p := NewPipeline(2)
	p.MulScalar(3)
	got := p.Run([]float64{1, 2, 3, 4, 5})
	for i, want := range []float64{3, 6, 9, 12, 15} {
		if got[i] != want {
			t.Errorf("grow[%d]: got %v, want %v", i, got[i], want)
		}
	}
}

func BenchmarkPipelineVsManual(b *testing.B) {
	const n = 1024
	input := make([]float64, n)
	for i := range input {
		input[i] = float64(i+1) * 0.001
	}

	p := NewPipeline(n)
	p.Exp().MulScalar(2).Log()

	b.Run("Pipeline", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = p.Run(input)
		}
	})

	b.Run("Manual", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			tmp1 := make([]float64, n)
			tmp2 := make([]float64, n)
			for j, v := range input {
				tmp1[j] = math.Exp(v)
			}
			for j, v := range tmp1 {
				tmp2[j] = v * 2
			}
			out := make([]float64, n)
			for j, v := range tmp2 {
				out[j] = math.Log(v)
			}
			_ = out
		}
	})
}
