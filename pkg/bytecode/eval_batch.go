package bytecode

import (
	"math"

	"github.com/emlgo/eml/pkg/fastmath"
)

const defaultBatchChunkSize = 1024

// BatchScratch holds the reusable working memory for columnar batch evaluation.
//
// EvalBatchColumnar used to build its depth x chunkSize stack, and the
// slice-of-slices indexing into it, on every call: 24,656 B and 2 allocations
// per invocation at the default chunk size. In a genetic-programming inner loop
// that is pure GC pressure.
//
// Hold a Scratch across calls to make a steady-state batch loop allocation
// free. A Scratch is not safe for concurrent use; give each goroutine its own,
// which is what the parallel fitness evaluation in ga.go does.
type BatchScratch struct {
	// stack is the flat depth*chunkSize backing array.
	stack []float64
	// cols are stack slices into it, one per stack level.
	cols [][]float64
	// chunkSize is the chunk length the current slices were built for.
	chunkSize int
	// depth is the stack depth they were built for.
	depth int
}

// cap returns a BatchScratch sized for a given stack depth and chunk length.
func newBatchScratch(depth, chunkSize int) BatchScratch {
	return BatchScratch{
		stack:     make([]float64, depth*chunkSize),
		cols:      make([][]float64, depth),
		chunkSize: chunkSize,
		depth:     depth,
	}
}

// ensure resizes the scratch if the program or chunk size grew, and returns the
// per-level column slices.
func (b *BatchScratch) ensure(depth, chunkSize int) [][]float64 {
	needed := depth * chunkSize
	if cap(b.stack) < needed || cap(b.cols) < depth {
		*b = newBatchScratch(depth, chunkSize)
	} else {
		b.depth = depth
		b.chunkSize = chunkSize
		b.stack = b.stack[:needed]
		b.cols = b.cols[:depth]
	}
	for i := 0; i < depth; i++ {
		b.cols[i] = b.stack[i*chunkSize : (i+1)*chunkSize]
	}
	return b.cols
}

// NewBatchScratch allocates scratch for a program evaluated in chunks of the
// given size. Passing 0 selects defaultBatchChunkSize.
func NewBatchScratch(p *Program, chunkSize int) BatchScratch {
	if chunkSize <= 0 {
		chunkSize = defaultBatchChunkSize
	}
	depth := 1
	if p != nil && p.MaxStackDepth > 0 {
		depth = p.MaxStackDepth
	}
	return newBatchScratch(depth, chunkSize)
}

// EvalBatchColumnar evaluates the bytecode program over an N-element batch of data
// where variables are structured in columnar format (data[varIdx] is []float64 of length N).
// Results are stored in dst.
// Interpretation overhead is amortized across vector chunks of size up to 1024.
//
// This form allocates the columnar stack on every call; use
// EvalBatchColumnarScratch in a hot loop.
func (p *Program) EvalBatchColumnar(data [][]float64, dst []float64) {
	s := NewBatchScratch(p, 0)
	p.EvalBatchColumnarScratch(data, dst, &s)
}

// EvalBatchColumnarScratch is EvalBatchColumnar with caller-supplied scratch.
// After the first call the scratch is reused and the call performs no heap
// allocations.
func (p *Program) EvalBatchColumnarScratch(data [][]float64, dst []float64, scratch *BatchScratch) {
	if p == nil || len(p.Ops) == 0 {
		return
	}
	n := len(dst)
	if n == 0 {
		return
	}

	chunkSize := defaultBatchChunkSize
	if n < chunkSize {
		chunkSize = n
	}

	depth := p.MaxStackDepth
	if depth < 1 {
		depth = p.CalculateMaxStackDepth()
	}
	stackCols := scratch.ensure(depth, chunkSize)

	for offset := 0; offset < n; offset += chunkSize {
		currLen := chunkSize
		if offset+currLen > n {
			currLen = n - offset
		}

		sp := 0
		cIdx := 0
		vIdx := 0

		for _, op := range p.Ops {
			switch op {
			case OpConst:
				cVal := p.Consts[cIdx]
				col := stackCols[sp][:currLen]
				for i := range col {
					col[i] = cVal
				}
				sp++
				cIdx++

			case OpVar:
				varID := p.VarIndices[vIdx]
				col := stackCols[sp][:currLen]
				if int(varID) < len(data) {
					srcCol := data[varID][offset : offset+currLen]
					copy(col, srcCol)
				} else {
					for i := range col {
						col[i] = 0
					}
				}
				sp++
				vIdx++

			case OpEML:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				fastmath.FastEmlBatchTo(xCol, yCol, xCol)

			case OpAdd:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] += yCol[i]
				}

			case OpSub:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] -= yCol[i]
				}

			case OpMul:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] *= yCol[i]
				}

			case OpDiv:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] /= yCol[i]
				}

			case OpPow:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] = math.Pow(xCol[i], yCol[i])
				}

			case OpNeg:
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] = -xCol[i]
				}

			case OpInv:
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] = 1.0 / xCol[i]
				}

			case OpSqrt:
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] = math.Sqrt(xCol[i])
				}

			case OpExp:
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] = fastmath.FastExp(xCol[i])
				}

			case OpLog:
				xCol := stackCols[sp-1][:currLen]
				for i := range xCol {
					xCol[i] = fastmath.FastLog(xCol[i])
				}

			default:
				// Every remaining opcode is a unary function; the inner loop is
				// a simple range so Go's auto-vectoriser can widen the ones
				// that have a vectorisable body.
				if fn := opUnaryTable[op]; fn != nil {
					xCol := stackCols[sp-1][:currLen]
					for i := range xCol {
						xCol[i] = fn(xCol[i])
					}
				}
			}
		}

		// Copy result chunk to dst
		copy(dst[offset:offset+currLen], stackCols[0][:currLen])
	}
}

// EvalBatch evaluates the program across a row-based dataset (data[i] is []float64 variables for sample i).
func (p *Program) EvalBatch(data [][]float64, dst []float64) {
	if p == nil || len(p.Ops) == 0 {
		return
	}
	n := len(data)
	if len(dst) < n {
		n = len(dst)
	}

	var localScratch [32]float64
	var scratch []float64
	if p.MaxStackDepth <= 32 {
		scratch = localScratch[:p.MaxStackDepth]
	} else {
		scratch = make([]float64, p.MaxStackDepth)
	}

	for i := 0; i < n; i++ {
		dst[i] = p.Eval(data[i], scratch)
	}
}
