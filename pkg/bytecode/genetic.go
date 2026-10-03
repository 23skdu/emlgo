package bytecode

import (
	"fmt"
	"math"
	"math/rand"
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

// mutationPools groups opcodes by arity so that mutation always swaps an opcode
// for one of the same arity, which is what keeps a candidate's stack sequence
// valid.
//
// The arity-1 pool covers every unary elementary function. Before the VM grew
// its function set it only offered exp/log/sqrt/neg/inv, so a search could never
// mutate into a program containing sin or gamma at all.
//
// Arity 0 is deliberately absent. OpConst and OpVar both have arity 0 but draw
// from *different* positional operand streams, so swapping one for the other
// without also editing Consts or VarIndices produces a program whose opcode
// counts no longer match its operands -- which panics during evaluation. Leaf
// opcodes are handled explicitly in Mutate instead, which does edit the stream.
var mutationPools = buildMutationPools()

func buildMutationPools() map[int][]OpCode {
	unary := []OpCode{OpNeg, OpInv}
	// Iterate the opcode space rather than the map so the pool is in a stable
	// order, which keeps a seeded search reproducible.
	for op := funcOpFirst; op <= funcOpLast; op++ {
		if _, ok := funcName(op); ok {
			unary = append(unary, op)
		}
	}
	return map[int][]OpCode{
		1: unary,
		2: {OpEML, OpAdd, OpSub, OpMul, OpDiv, OpPow},
	}
}

// poolFor returns a replacement opcode drawn from the mutation pool for op's
// arity, picking the entry at index i (wrapping). It returns op unchanged when
// no pool exists for that arity, which is the case for the arity-0 leaves.
func poolFor(op OpCode, i int) OpCode {
	pool, ok := mutationPools[op.Arity()]
	if !ok || len(pool) == 0 {
		return op
	}
	return pool[i%len(pool)]
}

// varCount reports how many distinct variable slots a program references.
func varCount(p *Program) int {
	max := -1
	for _, idx := range p.VarIndices {
		if int(idx) > max {
			max = int(idx)
		}
	}
	return max + 1
}

// insertIndex inserts v into s at position i, growing s as needed.
func insertIndex(s []uint16, i int, v uint16) []uint16 {
	if i >= len(s) {
		return append(s, v)
	}
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}

// insertConst inserts v into s at position i, growing s as needed.
func insertConst(s []float64, i int, v float64) []float64 {
	if i >= len(s) {
		return append(s, v)
	}
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}

// dropAt removes the element at position i, returning the shortened slice.
func dropAt[T any](s []T, i int) []T {
	if i >= len(s) {
		return s
	}
	copy(s[i:], s[i+1:])
	return s[:len(s)-1]
}

// splitPoint marks an opcode index at which the evaluation stack is back to
// height 1 -- that is, the prefix ending there is a complete expression -- along
// with how many constants and variable indices that prefix has consumed.
type splitPoint struct {
	ops     int
	consts  int
	indices int
}

// splitPoints returns every position in p where the stack height is 1, i.e.
// every complete-subexpression boundary.
//
// The earlier version of this function looked for height 0. A well-formed RPN
// expression only reaches height 0 at index 0 and ends at height 1, so that
// yielded exactly one split point and Crossover degenerated into swapping the
// two parents whole. Height 1 is the right boundary: splicing a prefix that
// ends at height 1 onto a suffix that starts at height 1 still ends at 1.
func splitPoints(p *Program) []splitPoint {
	if p == nil {
		return nil
	}
	out := make([]splitPoint, 0, len(p.Ops))
	height := 0
	cIdx, vIdx := 0, 0
	for i, op := range p.Ops {
		if height == 1 && i > 0 {
			out = append(out, splitPoint{ops: i, consts: cIdx, indices: vIdx})
		}
		switch op {
		case OpConst:
			cIdx++
		case OpVar:
			vIdx++
		}
		height = height - op.Arity() + 1
	}
	return out
}

// splice concatenates the prefix of a up to point pa with the suffix of b from
// point pb. Because both cut points are stack-balanced, the result is always a
// well-formed RPN program, and the positional operand streams line up because
// they are spliced at the matching positions.
func splice(a, b *Program, pa, pb splitPoint) *Program {
	out := NewProgram()
	out.Ops = make([]OpCode, 0, pa.ops+len(b.Ops)-pb.ops)
	out.Ops = append(out.Ops, a.Ops[:pa.ops]...)
	out.Ops = append(out.Ops, b.Ops[pb.ops:]...)

	out.Consts = make([]float64, 0, pa.consts+len(b.Consts)-pb.consts)
	out.Consts = append(out.Consts, a.Consts[:pa.consts]...)
	out.Consts = append(out.Consts, b.Consts[pb.consts:]...)

	out.VarIndices = make([]uint16, 0, pa.indices+len(b.VarIndices)-pb.indices)
	out.VarIndices = append(out.VarIndices, a.VarIndices[:pa.indices]...)
	out.VarIndices = append(out.VarIndices, b.VarIndices[pb.indices:]...)

	out.CalculateMaxStackDepth()
	return out
}

// Crossover performs single-point linear crossover between two programs,
// returning two children. The crossover points are chosen at stack-balanced
// positions, so both children are always evaluable and their constant and
// variable-index streams stay consistent with their opcode streams.
//
// rng must be non-nil. Use a seeded *rand.Rand for a reproducible search.
func Crossover(p1, p2 *Program, rng *rand.Rand) (*Program, *Program, error) {
	if p1 == nil || p2 == nil || len(p1.Ops) == 0 || len(p2.Ops) == 0 {
		return nil, nil, fmt.Errorf("invalid inputs to crossover")
	}
	if rng == nil {
		return nil, nil, fmt.Errorf("crossover requires a random source")
	}

	points1 := splitPoints(p1)
	points2 := splitPoints(p2)
	if len(points1) == 0 || len(points2) == 0 {
		return p1.Clone(), p2.Clone(), nil
	}

	pa := points1[rng.Intn(len(points1))]
	pb := points2[rng.Intn(len(points2))]

	c1 := splice(p1, p2, pa, pb)
	c2 := splice(p2, p1, pb, pa)

	if !validateStackFn(c1.Ops) || !validateStackFn(c2.Ops) {
		// Should be unreachable given stack-balanced cut points, but never
		// hand back a program that would panic during evaluation.
		return p1.Clone(), p2.Clone(), nil
	}
	return c1, c2, nil
}

// CrossoverOne produces a single offspring by splicing a subtree from p2 into p1.
// Unlike Crossover, it avoids allocating and constructing the unused second child.
func CrossoverOne(p1, p2 *Program, rng *rand.Rand) (*Program, error) {
	if p1 == nil || p2 == nil || len(p1.Ops) == 0 || len(p2.Ops) == 0 {
		return nil, fmt.Errorf("invalid inputs to crossover")
	}
	if rng == nil {
		return nil, fmt.Errorf("crossover requires a random source")
	}

	points1 := splitPoints(p1)
	points2 := splitPoints(p2)
	if len(points1) == 0 || len(points2) == 0 {
		return p1.Clone(), nil
	}

	pa := points1[rng.Intn(len(points1))]
	pb := points2[rng.Intn(len(points2))]

	c1 := splice(p1, p2, pa, pb)
	if !validateStackFn(c1.Ops) {
		return p1.Clone(), nil
	}
	return c1, nil
}

// Mutate performs random point mutation on a program with probability
// mutationRate per op. Opcodes are replaced by ones of the same arity, and a
// literal may be turned into a variable reference or vice versa, so a candidate
// keeps both a valid stack sequence and consistent operand streams. Constant
// operands are jittered by the same probability, otherwise a search can never
// move away from the constants a candidate happened to start with.
//
// rng must be non-nil.
func Mutate(p *Program, mutationRate float64, rng *rand.Rand) *Program {
	if p == nil || len(p.Ops) == 0 {
		return p
	}
	if rng == nil {
		return p.Clone()
	}
	mut := p.Clone()
	nVars := varCount(mut)
	if nVars < 1 {
		nVars = 1
	}

	// Walk the opcode stream tracking the positional operand cursors, so a leaf
	// can be converted without desynchronising Consts and VarIndices.
	// The two cursors track how many operands the opcodes seen so far consume.
	// Converting a leaf changes which stream it draws from, so each cursor must
	// be advanced or held to match:
	//
	//   CONST -> VAR: the constant is removed and the opcode no longer
	//     consumes one, so the element that shifted into this slot is the next
	//     one a CONST should read. cIdx stays put.
	//   VAR -> CONST: a constant is inserted at cIdx and this opcode reads it,
	//     so the next CONST reads cIdx+1. cIdx advances; vIdx stays put.
	cIdx, vIdx := 0, 0
	for i, op := range mut.Ops {
		switch op {
		case OpConst:
			if rng.Float64() < mutationRate {
				mut.Consts = dropAt(mut.Consts, cIdx)
				mut.VarIndices = insertIndex(mut.VarIndices, vIdx, uint16(rng.Intn(nVars))) // #nosec G115 -- < nVars
				mut.Ops[i] = OpVar
				vIdx++
			} else {
				cIdx++
			}
		case OpVar:
			if rng.Float64() < mutationRate {
				mut.VarIndices = dropAt(mut.VarIndices, vIdx)
				mut.Consts = insertConst(mut.Consts, cIdx, 2*rng.Float64()-1)
				mut.Ops[i] = OpConst
				cIdx++
			} else {
				vIdx++
			}
		default:
			if rng.Float64() < mutationRate {
				if pool, ok := mutationPools[op.Arity()]; ok && len(pool) > 0 {
					mut.Ops[i] = pool[rng.Intn(len(pool))]
				}
			}
		}
	}

	constJitter := 0.5
	for i := range mut.Consts {
		if rng.Float64() < mutationRate {
			mut.Consts[i] += constJitter * (2*rng.Float64() - 1)
		}
	}

	mut.CalculateMaxStackDepth()
	return mut
}

// exprNode is a node of the random expression tree that RandomProgram builds
// before flattening it to RPN.
//
// The tree matters: generating RPN directly (push a leaf, apply some unaries,
// maybe add a binary op) yields a *flat* chain whose stack height is 1 after the
// very first opcode. Such a program has exactly one stack-balanced split point,
// at index 0, so Crossover can only swap whole programs and never recombines
// subexpressions. Building a tree first gives every complete subtree a balanced
// RPN prefix, which is where the useful crossover points come from.
type exprNode struct {
	op    OpCode
	left  *exprNode
	right *exprNode
	value float64 // OpConst
	varID int     // OpVar
}

// RandomProgram builds a random, guaranteed-valid RPN program over nVars
// variables with at least minOps and at most maxOps opcodes.
//
// The program is generated as a random expression tree and then flattened, so
// ValidateStack holds by construction and no rejection sampling is needed.
func RandomProgram(nVars int, minOps, maxOps int, rng *rand.Rand) *Program {
	// #nosec G404 -- genetic operators need a reproducible source, not a
	// cryptographic one; callers pass their own seeded generator.
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	if nVars < 1 {
		nVars = 1
	}
	if maxOps < 1 {
		maxOps = 1
	}
	if minOps < 1 {
		minOps = 1
	}
	if minOps > maxOps {
		minOps = maxOps
	}
	target := minOps
	if maxOps > minOps {
		target = minOps + rng.Intn(maxOps-minOps+1)
	}

	unary := mutationPools[1]
	binary := mutationPools[2]

	// Each tree node costs one opcode, so the budget bounds the total.
	// Always build a composite root when the budget allows it, otherwise a
	// target of several opcodes could still yield a single-leaf program.
	budget := target - 1
	root := randomExpr(rng, nVars, &budget, unary, binary, true)

	p := NewProgram()
	emitRPN(root, p)
	p.CalculateMaxStackDepth()
	return p
}

// randomExpr builds one tree node, spending from the opcode budget. Leaves are
// returned for free once the budget is exhausted, which terminates the recursion
// without a depth limit.
func randomExpr(rng *rand.Rand, nVars int, budget *int, unary, binary []OpCode, root bool) *exprNode {
	if *budget <= 0 {
		return randomLeaf(rng, nVars)
	}
	*budget--

	r := rng.Float64()
	// The root must be composite when there is budget for it.
	if !root && r < 0.25 {
		return randomLeaf(rng, nVars)
	}
	if r < 0.55 {
		return &exprNode{
			op:   unary[rng.Intn(len(unary))],
			left: randomExpr(rng, nVars, budget, unary, binary, false),
		}
	}
	return &exprNode{
		op:    binary[rng.Intn(len(binary))],
		left:  randomExpr(rng, nVars, budget, unary, binary, false),
		right: randomExpr(rng, nVars, budget, unary, binary, false),
	}
}

func randomLeaf(rng *rand.Rand, nVars int) *exprNode {
	if rng.Intn(2) == 0 {
		return &exprNode{op: OpConst, value: 2*rng.Float64() - 1}
	}
	return &exprNode{op: OpVar, varID: rng.Intn(nVars)} // #nosec G115 -- rng.Intn(nVars) < nVars
}

// emitRPN appends the postfix encoding of the subtree to p, keeping the
// positional operand streams consistent as it goes.
func emitRPN(n *exprNode, p *Program) {
	switch n.op {
	case OpConst:
		p.Ops = append(p.Ops, OpConst)
		p.Consts = append(p.Consts, n.value)
		return
	case OpVar:
		p.Ops = append(p.Ops, OpVar)
		p.VarIndices = append(p.VarIndices, uint16(n.varID)) // #nosec G115 -- bounded by nVars
		return
	}

	if n.right == nil {
		emitRPN(n.left, p)
		p.Ops = append(p.Ops, n.op)
		return
	}
	emitRPN(n.left, p)
	emitRPN(n.right, p)
	p.Ops = append(p.Ops, n.op)
}

// ProgramDepth returns the maximum evaluation stack depth of p.
func ProgramDepth(p *Program) int {
	if p == nil {
		return 0
	}
	return p.CalculateMaxStackDepth()
}

// ProgramCost returns a rough size measure used to penalise bloated candidates.
func ProgramCost(p *Program) float64 {
	if p == nil {
		return math.Inf(1)
	}
	return float64(len(p.Ops)) + 0.25*float64(len(p.Consts))
}
