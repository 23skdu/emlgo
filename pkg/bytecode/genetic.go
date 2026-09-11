package bytecode

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// ValidateStack checks if a sequence of OpCodes represents a syntactically valid RPN expression.
// A valid program must:
// 1. Maintain stack height >= 1 after every prefix.
// 2. Terminate with a final stack height of exactly 1.
func ValidateStack(ops []OpCode) bool {
	if len(ops) == 0 {
		return false
	}
	height := 0
	for _, op := range ops {
		arity := op.Arity()
		if height < arity {
			return false
		}
		height = height - arity + 1
	}
	return height == 1
}

// Crossover performs single-point linear crossover between two programs, returning two children.
// The crossover point is chosen such that the resulting programs maintain valid stack semantics.
func Crossover(p1, p2 *Program) (*Program, *Program, error) {
	if p1 == nil || p2 == nil || len(p1.Ops) == 0 || len(p2.Ops) == 0 {
		return nil, nil, fmt.Errorf("invalid inputs to crossover")
	}

	c1 := p1.Clone()
	c2 := p2.Clone()

	// Find valid sub-slice cut points in p1 and p2 that have equal net stack delta
	cut1 := findValidSubtree(p1.Ops)
	cut2 := findValidSubtree(p2.Ops)

	if cut1.end > cut1.start && cut2.end > cut2.start {
		// Slice out subtrees and swap
		sub1 := make([]OpCode, cut1.end-cut1.start)
		copy(sub1, p1.Ops[cut1.start:cut1.end])

		sub2 := make([]OpCode, cut2.end-cut2.start)
		copy(sub2, p2.Ops[cut2.start:cut2.end])

		// Child 1: p1 prefix + sub2 + p1 suffix
		newOps1 := append([]OpCode{}, p1.Ops[:cut1.start]...)
		newOps1 = append(newOps1, sub2...)
		newOps1 = append(newOps1, p1.Ops[cut1.end:]...)

		// Child 2: p2 prefix + sub1 + p2 suffix
		newOps2 := append([]OpCode{}, p2.Ops[:cut2.start]...)
		newOps2 = append(newOps2, sub1...)
		newOps2 = append(newOps2, p2.Ops[cut2.end:]...)

		if ValidateStack(newOps1) && ValidateStack(newOps2) {
			c1.Ops = newOps1
			c2.Ops = newOps2
			c1.CalculateMaxStackDepth()
			c2.CalculateMaxStackDepth()
			return c1, c2, nil
		}
	}

	return c1, c2, nil
}

type subtreeRange struct {
	start int
	end   int
}

// findValidSubtree finds a sub-slice in ops that evaluates to a single net result (net delta +1).
func findValidSubtree(ops []OpCode) subtreeRange {
	n := len(ops)
	if n == 0 {
		return subtreeRange{0, 0}
	}

	// Pick a random end position
	rIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
	end := int(rIdx.Int64()) + 1

	// Walk backwards to find where net stack delta equals +1
	net := 0
	start := end - 1
	for start >= 0 {
		op := ops[start]
		// Scanning backwards: an operator with arity A requires A inputs and produces 1 output.
		// Moving backwards: operator increases deficit by (A - 1).
		net += (op.Arity() - 1)
		if net == -1 {
			return subtreeRange{start: start, end: end}
		}
		start--
	}

	return subtreeRange{0, n}
}

// Mutate performs random point mutation on a program with probability mutationRate per op.
func Mutate(p *Program, mutationRate float64) *Program {
	if p == nil || len(p.Ops) == 0 {
		return p
	}
	mut := p.Clone()

	for i, op := range mut.Ops {
		r, _ := rand.Int(rand.Reader, big.NewInt(1000))
		if float64(r.Int64())/1000.0 < mutationRate {
			// Mutate to an operator with matching arity to preserve stack validity
			mut.Ops[i] = mutateOp(op)
		}
	}

	mut.CalculateMaxStackDepth()
	return mut
}

func mutateOp(op OpCode) OpCode {
	switch op.Arity() {
	case 0:
		if op == OpConst {
			return OpVar
		}
		return OpConst
	case 1:
		ops := []OpCode{OpNeg, OpInv, OpSqrt, OpExp, OpLog}
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ops))))
		return ops[r.Int64()]
	case 2:
		ops := []OpCode{OpEML, OpAdd, OpSub, OpMul, OpDiv, OpPow}
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ops))))
		return ops[r.Int64()]
	default:
		return op
	}
}
