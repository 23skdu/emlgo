package bytecode

import (
	"math"

	"github.com/emlgo/eml/pkg/fastmath"
)

// Optimize performs constant folding and identity reductions directly on the linear bytecode.
// It simplifies identities like eml(1, 1) -> e and collapses canonical EML sequences.
func Optimize(p *Program) *Program {
	if p == nil || len(p.Ops) == 0 {
		return p
	}

	opt := NewProgram()
	var stackVal []float64
	var stackIsConst []bool

	cIdx := 0
	vIdx := 0

	for _, op := range p.Ops {
		switch op {
		case OpConst:
			val := p.Consts[cIdx]
			cIdx++
			stackVal = append(stackVal, val)
			stackIsConst = append(stackIsConst, true)

		case OpVar:
			idx := p.VarIndices[vIdx]
			vIdx++
			// Flush any pending constants
			stackVal = append(stackVal, 0)
			stackIsConst = append(stackIsConst, false)
			opt.Ops = append(opt.Ops, OpVar)
			opt.VarIndices = append(opt.VarIndices, idx)

		case OpEML:
			n := len(stackIsConst)
			if n >= 2 && stackIsConst[n-1] && stackIsConst[n-2] {
				// Fold constant eml(c1, c2)
				y := stackVal[n-1]
				x := stackVal[n-2]
				stackVal = stackVal[:n-2]
				stackIsConst = stackIsConst[:n-2]

				var folded float64
				if x == 1.0 && y == 1.0 {
					folded = math.E
				} else if x == 0.0 && y == 1.0 {
					folded = 1.0
				} else {
					folded = fastmath.FastEml(x, y)
				}
				stackVal = append(stackVal, folded)
				stackIsConst = append(stackIsConst, true)
			} else if n >= 1 && stackIsConst[n-1] && stackVal[n-1] == 1.0 {
				// x, CONST(1), EML -> x, EXP
				stackVal = stackVal[:n-1]
				stackIsConst = stackIsConst[:n-1]
				opt.Ops = append(opt.Ops, OpExp)
			} else {
				opt.Ops = append(opt.Ops, OpEML)
			}

		case OpAdd:
			n := len(stackIsConst)
			if n >= 2 && stackIsConst[n-1] && stackIsConst[n-2] {
				y := stackVal[n-1]
				x := stackVal[n-2]
				stackVal = stackVal[:n-2]
				stackIsConst = stackIsConst[:n-2]
				stackVal = append(stackVal, x+y)
				stackIsConst = append(stackIsConst, true)
			} else {
				opt.Ops = append(opt.Ops, OpAdd)
			}

		case OpSub:
			n := len(stackIsConst)
			if n >= 2 && stackIsConst[n-1] && stackIsConst[n-2] {
				y := stackVal[n-1]
				x := stackVal[n-2]
				stackVal = stackVal[:n-2]
				stackIsConst = stackIsConst[:n-2]
				stackVal = append(stackVal, x-y)
				stackIsConst = append(stackIsConst, true)
			} else {
				opt.Ops = append(opt.Ops, OpSub)
			}

		case OpMul:
			n := len(stackIsConst)
			if n >= 2 && stackIsConst[n-1] && stackIsConst[n-2] {
				y := stackVal[n-1]
				x := stackVal[n-2]
				stackVal = stackVal[:n-2]
				stackIsConst = stackIsConst[:n-2]
				stackVal = append(stackVal, x*y)
				stackIsConst = append(stackIsConst, true)
			} else {
				opt.Ops = append(opt.Ops, OpMul)
			}

		case OpDiv:
			n := len(stackIsConst)
			if n >= 2 && stackIsConst[n-1] && stackIsConst[n-2] && stackVal[n-1] != 0 {
				y := stackVal[n-1]
				x := stackVal[n-2]
				stackVal = stackVal[:n-2]
				stackIsConst = stackIsConst[:n-2]
				stackVal = append(stackVal, x/y)
				stackIsConst = append(stackIsConst, true)
			} else {
				opt.Ops = append(opt.Ops, OpDiv)
			}

		case OpNeg:
			n := len(stackIsConst)
			if n >= 1 && stackIsConst[n-1] {
				stackVal[n-1] = -stackVal[n-1]
			} else {
				opt.Ops = append(opt.Ops, OpNeg)
			}

		case OpExp:
			n := len(stackIsConst)
			if n >= 1 && stackIsConst[n-1] {
				stackVal[n-1] = fastmath.FastExp(stackVal[n-1])
			} else {
				opt.Ops = append(opt.Ops, OpExp)
			}

		case OpLog:
			n := len(stackIsConst)
			if n >= 1 && stackIsConst[n-1] {
				stackVal[n-1] = fastmath.FastLog(stackVal[n-1])
			} else {
				opt.Ops = append(opt.Ops, OpLog)
			}

		default:
			opt.Ops = append(opt.Ops, op)
		}
	}

	// Any remaining folded constants
	for i, isC := range stackIsConst {
		if isC {
			opt.Ops = append([]OpCode{OpConst}, opt.Ops...)
			opt.Consts = append([]float64{stackVal[i]}, opt.Consts...)
		}
	}

	opt.CalculateMaxStackDepth()
	return opt
}
