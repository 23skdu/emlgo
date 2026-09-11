package bytecode

import (
	"math"

	"github.com/emlgo/eml/pkg/fastmath"
)

const defaultBatchChunkSize = 1024

// EvalBatchColumnar evaluates the bytecode program over an N-element batch of data
// where variables are structured in columnar format (data[varIdx] is []float64 of length N).
// Results are stored in dst.
// Interpretation overhead is amortized across vector chunks of size up to 1024.
func (p *Program) EvalBatchColumnar(data [][]float64, dst []float64) {
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

	// Pre-allocate columnar stack buffer: depth x chunkSize
	depth := p.MaxStackDepth
	if depth < 1 {
		depth = p.CalculateMaxStackDepth()
	}
	stackBuf := make([]float64, depth*chunkSize)
	stackCols := make([][]float64, depth)
	for i := 0; i < depth; i++ {
		stackCols[i] = stackBuf[i*chunkSize : (i+1)*chunkSize]
	}

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
				for i := 0; i < currLen; i++ {
					xCol[i] += yCol[i]
				}

			case OpSub:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] -= yCol[i]
				}

			case OpMul:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] *= yCol[i]
				}

			case OpDiv:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] /= yCol[i]
				}

			case OpPow:
				sp--
				yCol := stackCols[sp][:currLen]
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] = math.Pow(xCol[i], yCol[i])
				}

			case OpNeg:
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] = -xCol[i]
				}

			case OpInv:
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] = 1.0 / xCol[i]
				}

			case OpSqrt:
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] = math.Sqrt(xCol[i])
				}

			case OpExp:
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] = fastmath.FastExp(xCol[i])
				}

			case OpLog:
				xCol := stackCols[sp-1][:currLen]
				for i := 0; i < currLen; i++ {
					xCol[i] = fastmath.FastLog(xCol[i])
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
