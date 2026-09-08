package eml

import "math"

// pipeStep is a single pipeline step operating on a slice in-place.
type pipeStep func(src, dst []float64)

// Pipeline is a composable, zero-allocation chain of batch operations.
// Steps are applied sequentially using double-buffered scratch slices.
type Pipeline struct {
	n     int
	buf   [2][]float64
	cur   int
	steps []pipeStep
}

// NewPipeline creates a new Pipeline for batches of up to n elements.
func NewPipeline(n int) *Pipeline {
	return &Pipeline{
		n:   n,
		buf: [2][]float64{make([]float64, n), make([]float64, n)},
	}
}

func (p *Pipeline) addStep(step pipeStep) *Pipeline {
	p.steps = append(p.steps, step)
	return p
}

// Exp appends an element-wise exp() step.
func (p *Pipeline) Exp() *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = math.Exp(v)
		}
	})
}

// Log appends an element-wise log() step.
func (p *Pipeline) Log() *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = math.Log(v)
		}
	})
}

// Sqrt appends an element-wise sqrt() step.
func (p *Pipeline) Sqrt() *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = math.Sqrt(v)
		}
	})
}

// Sin appends an element-wise sin() step.
func (p *Pipeline) Sin() *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = math.Sin(v)
		}
	})
}

// Cos appends an element-wise cos() step.
func (p *Pipeline) Cos() *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = math.Cos(v)
		}
	})
}

// Abs appends an element-wise abs() step.
func (p *Pipeline) Abs() *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = math.Abs(v)
		}
	})
}

// Neg appends an element-wise negation step.
func (p *Pipeline) Neg() *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = -v
		}
	})
}

// MulScalar appends a step that multiplies each element by scalar c.
func (p *Pipeline) MulScalar(c float64) *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = v * c
		}
	})
}

// AddScalar appends a step that adds scalar c to each element.
func (p *Pipeline) AddScalar(c float64) *Pipeline {
	return p.addStep(func(src, dst []float64) {
		for i, v := range src {
			dst[i] = v + c
		}
	})
}

// Run executes the pipeline on input and returns the result slice.
// The returned slice is owned by the pipeline buffer; copy it if you need it
// to outlive the next Run call.
func (p *Pipeline) Run(input []float64) []float64 {
	n := len(input)
	// Grow buffers if needed.
	if n > p.n {
		p.n = n
		p.buf[0] = make([]float64, n)
		p.buf[1] = make([]float64, n)
	}
	// Seed first buffer.
	copy(p.buf[0], input)
	src, dst := 0, 1
	for _, step := range p.steps {
		step(p.buf[src][:n], p.buf[dst][:n])
		src, dst = dst, src
	}
	return p.buf[src][:n]
}

// RunTo executes the pipeline on input and writes the result into output.
// output must have length >= len(input). This is the zero-allocation path.
func (p *Pipeline) RunTo(input, output []float64) {
	result := p.Run(input)
	copy(output, result)
}
