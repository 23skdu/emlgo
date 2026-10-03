package eml

type stepKind int

const (
	stepExp stepKind = iota
	stepLog
	stepSqrt
	stepSin
	stepCos
	stepAbs
	stepAffine // dst[i] = src[i] * scale + bias
	stepCustom
)

type pipeStep struct {
	kind  stepKind
	scale float64
	bias  float64
	fn    func(src, dst []float64)
}

func (s pipeStep) run(src, dst []float64) {
	switch s.kind {
	case stepExp:
		ExpSIMDTo(src, dst)
	case stepLog:
		LogSIMDTo(src, dst)
	case stepSqrt:
		SqrtSIMDTo(src, dst)
	case stepSin:
		SinSIMDTo(src, dst)
	case stepCos:
		CosSIMDTo(src, dst)
	case stepAbs:
		AbsSIMDTo(src, dst)
	case stepAffine:
		executeAffine(src, s.scale, s.bias, dst)
	case stepCustom:
		s.fn(src, dst)
	}
}

func executeAffine(src []float64, scale, bias float64, dst []float64) {
	if scale == 1 && bias == 0 {
		copy(dst, src)
		return
	}
	if scale == 1 {
		AddScalarSIMDTo(src, bias, dst)
		return
	}
	if bias == 0 {
		MulScalarSIMDTo(src, scale, dst)
		return
	}
	// Fused affine pass: dst[i] = src[i] * scale + bias
	n := len(src)
	i := 0
	for ; i+3 < n; i += 4 {
		dst[i] = src[i]*scale + bias
		dst[i+1] = src[i+1]*scale + bias
		dst[i+2] = src[i+2]*scale + bias
		dst[i+3] = src[i+3]*scale + bias
	}
	for ; i < n; i++ {
		dst[i] = src[i]*scale + bias
	}
}

// Pipeline is a composable, zero-allocation chain of batch operations.
// Steps are applied sequentially using double-buffered scratch slices.
type Pipeline struct {
	n     int
	buf   [2][]float64
	steps []pipeStep
}

// NewPipeline creates a new Pipeline for batches of up to n elements.
// Allocates only one scratch buffer initially; the second buffer is allocated
// lazily only if Run() is called instead of RunTo(), reducing memory usage by 50%.
func NewPipeline(n int) *Pipeline {
	return &Pipeline{
		n:   n,
		buf: [2][]float64{make([]float64, n), nil},
	}
}

func (p *Pipeline) addAffine(scale, bias float64) *Pipeline {
	if len(p.steps) > 0 {
		last := &p.steps[len(p.steps)-1]
		if last.kind == stepAffine {
			// Step fusion: combine consecutive affine transforms
			// (last.scale * x + last.bias) * scale + bias
			// = (last.scale * scale) * x + (last.bias * scale + bias)
			last.scale *= scale
			last.bias = last.bias*scale + bias
			return p
		}
	}
	p.steps = append(p.steps, pipeStep{
		kind:  stepAffine,
		scale: scale,
		bias:  bias,
	})
	return p
}

// Custom appends a custom transformation step to the pipeline.
func (p *Pipeline) Custom(step func(src, dst []float64)) *Pipeline {
	p.steps = append(p.steps, pipeStep{
		kind: stepCustom,
		fn:   step,
	})
	return p
}

// Reset removes every step from the pipeline, keeping its scratch buffers.
//
// The builder methods (Exp, MulScalar, ...) each append a step, so a Pipeline
// that is rebuilt inside a loop accumulates steps without bound. Call Reset
// before rebuilding, or build the chain once and reuse it -- both patterns are
// shown on NewPipeline.
//
// Reset makes the first pattern safe:
//
//	p := eml.NewPipeline(n)
//	for _, batch := range batches {
//	    p.Reset().Exp().MulScalar(2.0).Log().RunTo(batch, out)
//	}
func (p *Pipeline) Reset() *Pipeline {
	clear(p.steps)
	p.steps = p.steps[:0]
	return p
}

// Len reports how many steps the pipeline currently holds. It is intended for
// tests and diagnostics.
func (p *Pipeline) Len() int { return len(p.steps) }

// Exp appends an element-wise exp() step dispatched via SIMD.
func (p *Pipeline) Exp() *Pipeline {
	p.steps = append(p.steps, pipeStep{kind: stepExp})
	return p
}

// Log appends an element-wise log() step dispatched via SIMD.
func (p *Pipeline) Log() *Pipeline {
	p.steps = append(p.steps, pipeStep{kind: stepLog})
	return p
}

// Sqrt appends an element-wise sqrt() step dispatched via SIMD.
func (p *Pipeline) Sqrt() *Pipeline {
	p.steps = append(p.steps, pipeStep{kind: stepSqrt})
	return p
}

// Sin appends an element-wise sin() step dispatched via SIMD.
func (p *Pipeline) Sin() *Pipeline {
	p.steps = append(p.steps, pipeStep{kind: stepSin})
	return p
}

// Cos appends an element-wise cos() step dispatched via SIMD.
func (p *Pipeline) Cos() *Pipeline {
	p.steps = append(p.steps, pipeStep{kind: stepCos})
	return p
}

// Abs appends an element-wise abs() step dispatched via SIMD.
func (p *Pipeline) Abs() *Pipeline {
	p.steps = append(p.steps, pipeStep{kind: stepAbs})
	return p
}

// Neg appends an element-wise negation step.
func (p *Pipeline) Neg() *Pipeline {
	return p.addAffine(-1, 0)
}

// MulScalar appends a step that multiplies each element by scalar c.
func (p *Pipeline) MulScalar(c float64) *Pipeline {
	return p.addAffine(c, 0)
}

// AddScalar appends a step that adds scalar c to each element.
func (p *Pipeline) AddScalar(c float64) *Pipeline {
	return p.addAffine(1, c)
}

// SubScalar appends a step that subtracts scalar c from each element.
func (p *Pipeline) SubScalar(c float64) *Pipeline {
	return p.addAffine(1, -c)
}

// DivScalar appends a step that divides each element by scalar c.
func (p *Pipeline) DivScalar(c float64) *Pipeline {
	return p.addAffine(1.0/c, 0)
}

// Run executes the pipeline on input and returns the result slice.
// The returned slice is owned by the pipeline buffer; copy it if you need it
// to outlive the next Run call.
//
// Running a pipeline with no steps returns a copy of the input.
func (p *Pipeline) Run(input []float64) []float64 {
	n := len(input)
	p.grow(n)
	if len(p.buf[1]) < n {
		p.buf[1] = make([]float64, n)
	}
	if len(p.steps) == 0 {
		res := make([]float64, n)
		copy(res, input)
		return res
	}
	// Seed first buffer.
	copy(p.buf[0], input)
	src, dst := 0, 1
	for _, step := range p.steps {
		step.run(p.buf[src][:n], p.buf[dst][:n])
		src, dst = dst, src
	}
	return p.buf[src][:n]
}

// RunTo executes the pipeline on input and writes the result into output.
// output must have length >= len(input).
//
// This is the zero-allocation path and, unlike Run, does not copy the result out
// of the pipeline: the first step writes straight into output and the remaining
// steps ping-pong between output and the scratch buffer.
func (p *Pipeline) RunTo(input, output []float64) {
	n := len(input)
	if n > len(output) {
		panic("pipeline: output shorter than input")
	}
	p.grow(n)

	if len(p.steps) == 0 {
		copy(output, input)
		return
	}

	// The first step reads the input and writes straight into output, so no
	// copy of the input is needed. Each later step alternates between output
	// and the single scratch buffer.
	scratch := p.buf[0]
	p.steps[0].run(input, output[:n])

	inOutput := true
	for _, step := range p.steps[1:] {
		if inOutput {
			step.run(output[:n], scratch[:n])
		} else {
			step.run(scratch[:n], output[:n])
		}
		inOutput = !inOutput
	}
	if !inOutput {
		// An even number of remaining steps left the result in scratch.
		copy(output[:n], scratch[:n])
	}
}

// grow enlarges the scratch buffers to hold at least n elements.
func (p *Pipeline) grow(n int) {
	if n <= p.n {
		return
	}
	p.n = n
	p.buf[0] = make([]float64, n)
	if p.buf[1] != nil {
		p.buf[1] = make([]float64, n)
	}
}
