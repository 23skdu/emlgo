package eml

import (
	"math"
	"testing"
)

// TestPipelineResetStopsStepGrowth is the regression test for the unbounded
// growth footgun: the builder methods append a step, so a Pipeline rebuilt
// inside a loop used to accumulate 3 steps per iteration forever.
func TestPipelineResetStopsStepGrowth(t *testing.T) {
	const n = 64
	in := make([]float64, n)
	for i := range in {
		in[i] = 0.5 + float64(i)*0.01
	}
	out := make([]float64, n)

	p := NewPipeline(n)
	for i := 0; i < 100; i++ {
		p.Reset().Exp().MulScalar(2.0).Log().RunTo(in, out)
	}
	if got := p.Len(); got != 3 {
		t.Errorf("after 100 rebuilds with Reset, Len() = %d, want 3", got)
	}

	// Without Reset the count grows linearly; assert it does, to document why
	// Reset exists. Use a fresh pipeline so we do not leave it dirty.
	q := NewPipeline(n)
	for i := 0; i < 10; i++ {
		q.Exp().MulScalar(2.0).Log().RunTo(in, out)
	}
	if got := q.Len(); got != 30 {
		t.Errorf("without Reset, Len() = %d, want 30", got)
	}
}

// TestPipelineRunToMatchesReference checks RunTo against the equivalent
// hand-written expression for every supported step, including the odd/even step
// counts that exercise the output/scratch ping-pong parity.
func TestPipelineRunToMatchesReference(t *testing.T) {
	const n = 257 // deliberately not a round number
	in := make([]float64, n)
	for i := range in {
		in[i] = 0.25 + float64(i)*0.02
	}
	out := make([]float64, n)

	tests := []struct {
		name  string
		build func(p *Pipeline) *Pipeline
		want  func(v float64) float64
	}{
		{
			name:  "no steps",
			build: func(p *Pipeline) *Pipeline { return p },
			want:  func(v float64) float64 { return v },
		},
		{
			name:  "one step",
			build: func(p *Pipeline) *Pipeline { return p.Exp() },
			want:  func(v float64) float64 { return math.Exp(v) },
		},
		{
			name:  "two steps",
			build: func(p *Pipeline) *Pipeline { return p.Exp().Log() },
			want:  func(v float64) float64 { return math.Log(math.Exp(v)) },
		},
		{
			name:  "three steps",
			build: func(p *Pipeline) *Pipeline { return p.Exp().MulScalar(2).Log() },
			want:  func(v float64) float64 { return math.Log(math.Exp(v) * 2) },
		},
		{
			name:  "four steps",
			build: func(p *Pipeline) *Pipeline { return p.Exp().MulScalar(2).Log().Abs() },
			want:  func(v float64) float64 { return math.Abs(math.Log(math.Exp(v) * 2)) },
		},
		{
			name:  "five steps",
			build: func(p *Pipeline) *Pipeline { return p.Exp().MulScalar(2).Log().Abs().Neg() },
			want:  func(v float64) float64 { return -math.Abs(math.Log(math.Exp(v) * 2)) },
		},
		{
			name: "add scalar",
			build: func(p *Pipeline) *Pipeline {
				return p.AddScalar(1.5).MulScalar(2).SubScalar(3.5).DivScalar(2).Neg()
			},
			want: func(v float64) float64 { return -(((v+1.5)*2 - 3.5) / 2) },
		},
		{
			name:  "sqrt and trig",
			build: func(p *Pipeline) *Pipeline { return p.Sqrt().Sin().Cos() },
			want:  func(v float64) float64 { return math.Cos(math.Sin(math.Sqrt(v))) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPipeline(n)
			tc.build(p)
			p.RunTo(in, out)
			for i := range in {
				if got, want := out[i], tc.want(in[i]); math.Abs(got-want) > 1e-12 {
					t.Fatalf("out[%d] = %v, want %v", i, got, want)
				}
			}
		})
	}
}

// TestPipelineRunToDoesNotMutateInput verifies RunTo leaves the caller's input
// untouched, which is what the new write-into-output path must guarantee.
func TestPipelineRunToDoesNotMutateInput(t *testing.T) {
	const n = 64
	in := make([]float64, n)
	orig := make([]float64, n)
	for i := range in {
		in[i] = 1 + float64(i)*0.1
		orig[i] = in[i]
	}
	out := make([]float64, n)

	p := NewPipeline(n)
	p.Neg().AddScalar(10)
	p.RunTo(in, out)

	for i := range in {
		if in[i] != orig[i] {
			t.Fatalf("input mutated at %d: %v != %v", i, in[i], orig[i])
		}
		if got, want := out[i], 10-in[i]; got != want {
			t.Fatalf("out[%d] = %v, want %v", i, got, want)
		}
	}
}

// TestPipelineRunToShortOutputPanics documents the contract on output length.
func TestPipelineRunToShortOutputPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic when output is shorter than input")
		}
	}()
	p := NewPipeline(8)
	p.Exp().RunTo(make([]float64, 8), make([]float64, 4))
}

// TestPipelineGrowsBuffers checks a pipeline sized smaller than its input still
// produces correct results.
func TestPipelineGrowsBuffers(t *testing.T) {
	const n = 300
	p := NewPipeline(4)
	p.Exp().MulScalar(2)

	in := make([]float64, n)
	out := make([]float64, n)
	for i := range in {
		in[i] = 0.1 + float64(i)*0.005
	}
	p.RunTo(in, out)
	for i := range in {
		if got, want := out[i], math.Exp(in[i])*2; got != want {
			t.Fatalf("out[%d] = %v, want %v", i, got, want)
		}
	}
}

// TestPipelineRunMatchesRunTo checks the two entry points agree.
func TestPipelineRunMatchesRunTo(t *testing.T) {
	const n = 129
	in := make([]float64, n)
	for i := range in {
		in[i] = 0.3 + float64(i)*0.01
	}

	p := NewPipeline(n)
	p.Exp().AddScalar(1).Sqrt()

	got := make([]float64, n)
	p.RunTo(in, got)
	want := p.Run(in)

	for i := range in {
		if got[i] != want[i] {
			t.Fatalf("index %d: RunTo = %v, Run = %v", i, got[i], want[i])
		}
	}
}
