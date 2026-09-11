package bytecode

import (
	"fmt"
	"strings"
)

// OpCode represents a virtual machine operation in 1 byte.
type OpCode uint8

const (
	// OpConst pushes Consts[constIdx++] onto the evaluation stack.
	OpConst OpCode = iota

	// OpVar pushes vars[VarIndices[varIdx++]] onto the evaluation stack.
	OpVar

	// OpEML pops y and x, computing eml(x, y) = exp(x) - ln(y).
	OpEML

	// Native elementary operations (used for collapsed EML identities)
	OpAdd  // Pop y, Pop x, Push x + y
	OpSub  // Pop y, Pop x, Push x - y
	OpMul  // Pop y, Pop x, Push x * y
	OpDiv  // Pop y, Pop x, Push x / y
	OpPow  // Pop y, Pop x, Push x ^ y
	OpNeg  // Pop x, Push -x
	OpInv  // Pop x, Push 1.0 / x
	OpSqrt // Pop x, Push sqrt(x)
	OpExp  // Pop x, Push exp(x)
	OpLog  // Pop x, Push log(x)
)

// String returns a readable representation of the OpCode.
func (op OpCode) String() string {
	switch op {
	case OpConst:
		return "CONST"
	case OpVar:
		return "VAR"
	case OpEML:
		return "EML"
	case OpAdd:
		return "ADD"
	case OpSub:
		return "SUB"
	case OpMul:
		return "MUL"
	case OpDiv:
		return "DIV"
	case OpPow:
		return "POW"
	case OpNeg:
		return "NEG"
	case OpInv:
		return "INV"
	case OpSqrt:
		return "SQRT"
	case OpExp:
		return "EXP"
	case OpLog:
		return "LOG"
	default:
		return fmt.Sprintf("OP(%d)", op)
	}
}

// Arity returns the number of stack operands popped by this opcode.
func (op OpCode) Arity() int {
	switch op {
	case OpConst, OpVar:
		return 0
	case OpNeg, OpInv, OpSqrt, OpExp, OpLog:
		return 1
	case OpEML, OpAdd, OpSub, OpMul, OpDiv, OpPow:
		return 2
	default:
		return 0
	}
}

// Program represents a linear bytecode buffer in Structure of Arrays (SoA) layout.
// Memory is contiguous, prefetcher-friendly, and completely free of heap pointer chasing.
type Program struct {
	Ops           []OpCode  // Contiguous opcode stream
	Consts        []float64 // Numerical constants referenced by OpConst
	VarIndices    []uint16  // Variable slot indices referenced by OpVar
	MaxStackDepth int       // Precomputed maximum evaluation stack depth
}

// NewProgram creates an empty Program.
func NewProgram() *Program {
	return &Program{
		Ops:        make([]OpCode, 0, 16),
		Consts:     make([]float64, 0, 4),
		VarIndices: make([]uint16, 0, 4),
	}
}

// Clone creates an exact deep copy of the program.
func (p *Program) Clone() *Program {
	if p == nil {
		return nil
	}
	ops := make([]OpCode, len(p.Ops))
	copy(ops, p.Ops)

	consts := make([]float64, len(p.Consts))
	copy(consts, p.Consts)

	varIndices := make([]uint16, len(p.VarIndices))
	copy(varIndices, p.VarIndices)

	return &Program{
		Ops:           ops,
		Consts:        consts,
		VarIndices:    varIndices,
		MaxStackDepth: p.MaxStackDepth,
	}
}

// CalculateMaxStackDepth computes and updates the peak stack depth required to execute the program.
func (p *Program) CalculateMaxStackDepth() int {
	depth := 0
	maxDepth := 0
	for _, op := range p.Ops {
		switch op.Arity() {
		case 0:
			depth++
		case 1:
			// net change is 0
		case 2:
			depth--
		}
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	if maxDepth < 1 {
		maxDepth = 1
	}
	p.MaxStackDepth = maxDepth
	return maxDepth
}

// String returns a disassembly view of the bytecode program.
func (p *Program) String() string {
	if p == nil {
		return "<nil>"
	}
	var b strings.Builder
	cIdx := 0
	vIdx := 0
	for i, op := range p.Ops {
		if i > 0 {
			b.WriteString(" ")
		}
		switch op {
		case OpConst:
			if cIdx < len(p.Consts) {
				fmt.Fprintf(&b, "CONST(%v)", p.Consts[cIdx])
				cIdx++
			} else {
				b.WriteString("CONST(?)")
			}
		case OpVar:
			if vIdx < len(p.VarIndices) {
				fmt.Fprintf(&b, "VAR(%d)", p.VarIndices[vIdx])
				vIdx++
			} else {
				b.WriteString("VAR(?)")
			}
		default:
			b.WriteString(op.String())
		}
	}
	return b.String()
}
