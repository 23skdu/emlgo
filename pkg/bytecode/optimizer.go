package bytecode

import (
	"math"

	"github.com/emlgo/eml/pkg/fastmath"
)

type optNode struct {
	op       OpCode
	isConst  bool
	val      float64
	varIdx   uint16
	children []*optNode
}

func newConstNode(v float64) *optNode {
	return &optNode{
		op:      OpConst,
		isConst: true,
		val:     v,
	}
}

func newVarNode(idx uint16) *optNode {
	return &optNode{
		op:     OpVar,
		varIdx: idx,
	}
}

func newUnaryNode(op OpCode, child *optNode) *optNode {
	return &optNode{
		op:       op,
		children: []*optNode{child},
	}
}

func newBinaryNode(op OpCode, left, right *optNode) *optNode {
	return &optNode{
		op:       op,
		children: []*optNode{left, right},
	}
}

func (n *optNode) emit(opt *Program) {
	if n.isConst {
		opt.Ops = append(opt.Ops, OpConst)
		opt.Consts = append(opt.Consts, n.val)
		return
	}
	if n.op == OpVar {
		opt.Ops = append(opt.Ops, OpVar)
		opt.VarIndices = append(opt.VarIndices, n.varIdx)
		return
	}
	for _, child := range n.children {
		child.emit(opt)
	}
	opt.Ops = append(opt.Ops, n.op)
}

// Optimize performs constant folding and identity reductions via an AST-based
// tree transformation. It preserves exact operand ordering and commutativity rules,
// simplifies identities like eml(1, 1) -> e and eml(x, 1) -> exp(x), and collapses
// algebraic neutral elements (x+0, x-0, 0-x, x*1, x*0, x/1, x^0, x^1).
func Optimize(p *Program) *Program {
	if p == nil || len(p.Ops) == 0 {
		return p
	}

	var stack []*optNode
	cIdx := 0
	vIdx := 0

	for _, op := range p.Ops {
		switch op {
		case OpConst:
			if cIdx >= len(p.Consts) {
				return p.Clone()
			}
			val := p.Consts[cIdx]
			cIdx++
			stack = append(stack, newConstNode(val))

		case OpVar:
			if vIdx >= len(p.VarIndices) {
				return p.Clone()
			}
			idx := p.VarIndices[vIdx]
			vIdx++
			stack = append(stack, newVarNode(idx))

		case OpNeg:
			if len(stack) < 1 {
				return p.Clone()
			}
			child := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if child.isConst {
				stack = append(stack, newConstNode(-child.val))
			} else if child.op == OpNeg {
				// -(-x) == x
				stack = append(stack, child.children[0])
			} else {
				stack = append(stack, newUnaryNode(OpNeg, child))
			}

		case OpInv:
			if len(stack) < 1 {
				return p.Clone()
			}
			child := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if child.isConst && child.val != 0 {
				stack = append(stack, newConstNode(1.0/child.val))
			} else if child.op == OpInv {
				// 1 / (1 / x) == x
				stack = append(stack, child.children[0])
			} else {
				stack = append(stack, newUnaryNode(OpInv, child))
			}

		case OpExp:
			if len(stack) < 1 {
				return p.Clone()
			}
			child := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if child.isConst {
				stack = append(stack, newConstNode(fastmath.FastExp(child.val)))
			} else {
				stack = append(stack, newUnaryNode(OpExp, child))
			}

		case OpLog:
			if len(stack) < 1 {
				return p.Clone()
			}
			child := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if child.isConst && child.val > 0 {
				stack = append(stack, newConstNode(fastmath.FastLog(child.val)))
			} else {
				stack = append(stack, newUnaryNode(OpLog, child))
			}

		case OpAdd:
			if len(stack) < 2 {
				return p.Clone()
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if left.isConst && right.isConst {
				stack = append(stack, newConstNode(left.val+right.val))
			} else if right.isConst && right.val == 0.0 {
				stack = append(stack, left)
			} else if left.isConst && left.val == 0.0 {
				stack = append(stack, right)
			} else {
				stack = append(stack, newBinaryNode(OpAdd, left, right))
			}

		case OpSub:
			if len(stack) < 2 {
				return p.Clone()
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if left.isConst && right.isConst {
				stack = append(stack, newConstNode(left.val-right.val))
			} else if right.isConst && right.val == 0.0 {
				stack = append(stack, left)
			} else if left.isConst && left.val == 0.0 {
				stack = append(stack, newUnaryNode(OpNeg, right))
			} else {
				stack = append(stack, newBinaryNode(OpSub, left, right))
			}

		case OpMul:
			if len(stack) < 2 {
				return p.Clone()
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if left.isConst && right.isConst {
				stack = append(stack, newConstNode(left.val*right.val))
			} else if right.isConst && right.val == 1.0 {
				stack = append(stack, left)
			} else if left.isConst && left.val == 1.0 {
				stack = append(stack, right)
			} else if right.isConst && right.val == 0.0 {
				stack = append(stack, newConstNode(0.0))
			} else if left.isConst && left.val == 0.0 {
				stack = append(stack, newConstNode(0.0))
			} else {
				stack = append(stack, newBinaryNode(OpMul, left, right))
			}

		case OpDiv:
			if len(stack) < 2 {
				return p.Clone()
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if left.isConst && right.isConst && right.val != 0 {
				stack = append(stack, newConstNode(left.val/right.val))
			} else if right.isConst && right.val == 1.0 {
				stack = append(stack, left)
			} else if left.isConst && left.val == 0.0 {
				stack = append(stack, newConstNode(0.0))
			} else {
				stack = append(stack, newBinaryNode(OpDiv, left, right))
			}

		case OpPow:
			if len(stack) < 2 {
				return p.Clone()
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if left.isConst && right.isConst {
				if !(left.val < 0 && math.Floor(right.val) != right.val) {
					stack = append(stack, newConstNode(math.Pow(left.val, right.val)))
				} else {
					stack = append(stack, newBinaryNode(OpPow, left, right))
				}
			} else if right.isConst && right.val == 0.0 {
				stack = append(stack, newConstNode(1.0))
			} else if right.isConst && right.val == 1.0 {
				stack = append(stack, left)
			} else {
				stack = append(stack, newBinaryNode(OpPow, left, right))
			}

		case OpEML:
			if len(stack) < 2 {
				return p.Clone()
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			if left.isConst && right.isConst {
				var folded float64
				if left.val == 1.0 && right.val == 1.0 {
					folded = math.E
				} else if left.val == 0.0 && right.val == 1.0 {
					folded = 1.0
				} else if right.val > 0 {
					folded = fastmath.FastEml(left.val, right.val)
				} else {
					stack = append(stack, newBinaryNode(OpEML, left, right))
					break
				}
				stack = append(stack, newConstNode(folded))
			} else if right.isConst && right.val == 1.0 {
				// eml(x, 1) = exp(x) - ln(1) = exp(x)
				stack = append(stack, newUnaryNode(OpExp, left))
			} else {
				stack = append(stack, newBinaryNode(OpEML, left, right))
			}

		default:
			if fn := opUnaryTable[op]; fn != nil {
				if len(stack) < 1 {
					return p.Clone()
				}
				child := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if child.isConst {
					stack = append(stack, newConstNode(fn(child.val)))
				} else {
					stack = append(stack, newUnaryNode(op, child))
				}
			} else {
				// Unknown opcode or unsupported multi-arity opcode
				return p.Clone()
			}
		}
	}

	if cIdx != len(p.Consts) || vIdx != len(p.VarIndices) {
		return p.Clone()
	}

	opt := NewProgram()
	for _, n := range stack {
		n.emit(opt)
	}
	opt.CalculateMaxStackDepth()
	return opt
}
