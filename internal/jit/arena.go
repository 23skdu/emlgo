package jit

import "math"

type nodeKind uint8

const (
	kindNumber nodeKind = iota
	kindVariable
	kindUnaryOp
	kindBinaryOp
	kindFuncCall
)

type ArenaNode struct {
	Kind  nodeKind
	Op    rune
	Value float64
	Name  [16]byte
	Left  *ArenaNode
	Right *ArenaNode
	_     [7]byte
}

type Arena struct {
	buf []ArenaNode
	off int
}

// NewArena creates a new Arena with the given initial capacity.
func NewArena(capacity int) *Arena {
	if capacity <= 0 {
		capacity = 64
	}
	return &Arena{
		buf: make([]ArenaNode, capacity),
	}
}

// Alloc returns a pointer to the next available ArenaNode, growing the buffer if needed.
func (a *Arena) Alloc() *ArenaNode {
	if a.off >= len(a.buf) {
		a.grow()
	}
	n := &a.buf[a.off]
	a.off++
	*n = ArenaNode{}
	return n
}

func (a *Arena) grow() {
	newBuf := make([]ArenaNode, len(a.buf)*2)
	copy(newBuf, a.buf)
	a.buf = newBuf
}

// NumNodes returns the number of nodes allocated in the arena.
func (a *Arena) NumNodes() int {
	return a.off
}

// Reset resets the arena so that all previously allocated nodes become available for reuse.
func (a *Arena) Reset() {
	a.off = 0
}

// ToInterface converts an ArenaNode to the corresponding Node interface value.
func (a *Arena) ToInterface(n *ArenaNode) Node {
	return a.toInterfaceDepth(n, 0)
}

func (a *Arena) toInterfaceDepth(n *ArenaNode, depth int) Node {
	if n == nil || depth > MaxASTDepth {
		return nil
	}
	switch n.Kind {
	case kindNumber:
		return Number{Value: n.Value}
	case kindVariable:
		return Variable{}
	case kindUnaryOp:
		return UnaryOp{Op: n.Op, Operand: a.toInterfaceDepth(n.Left, depth+1)}
	case kindBinaryOp:
		return BinaryOp{Left: a.toInterfaceDepth(n.Left, depth+1), Op: n.Op, Right: a.toInterfaceDepth(n.Right, depth+1)}
	case kindFuncCall:
		name := ""
		for i := 0; i < len(n.Name) && n.Name[i] != 0; i++ {
			name += string(n.Name[i])
		}
		return FunctionCall{Name: name, Arg: a.toInterfaceDepth(n.Left, depth+1)}
	}
	return nil
}

// FromInterface converts a Node interface value into a new ArenaNode allocated from the arena.
func (a *Arena) FromInterface(n Node) *ArenaNode {
	return a.fromInterfaceDepth(n, 0)
}

func (a *Arena) fromInterfaceDepth(n Node, depth int) *ArenaNode {
	if n == nil || depth > MaxASTDepth {
		return nil
	}
	switch v := n.(type) {
	case Number:
		node := a.Alloc()
		node.Kind = kindNumber
		node.Value = v.Value
		return node
	case Variable:
		node := a.Alloc()
		node.Kind = kindVariable
		return node
	case UnaryOp:
		node := a.Alloc()
		node.Kind = kindUnaryOp
		node.Op = v.Op
		node.Left = a.fromInterfaceDepth(v.Operand, depth+1)
		return node
	case BinaryOp:
		node := a.Alloc()
		node.Kind = kindBinaryOp
		node.Op = v.Op
		node.Left = a.fromInterfaceDepth(v.Left, depth+1)
		node.Right = a.fromInterfaceDepth(v.Right, depth+1)
		return node
	case FunctionCall:
		node := a.Alloc()
		node.Kind = kindFuncCall
		for i := 0; i < len(v.Name) && i < 15; i++ {
			node.Name[i] = v.Name[i]
		}
		node.Left = a.fromInterfaceDepth(v.Arg, depth+1)
		return node
	}
	return nil
}

// EvalArena evaluates the expression tree rooted at n with variable x set to the given value.
func EvalArena(n *ArenaNode, x float64) float64 {
	return evalArenaDepth(n, x, 0)
}

func evalArenaDepth(n *ArenaNode, x float64, depth int) float64 {
	if n == nil || depth > MaxASTDepth {
		return 0
	}
	switch n.Kind {
	case kindNumber:
		return n.Value
	case kindVariable:
		return x
	case kindUnaryOp:
		return -evalArenaDepth(n.Left, x, depth+1)
	case kindBinaryOp:
		left := evalArenaDepth(n.Left, x, depth+1)
		right := evalArenaDepth(n.Right, x, depth+1)
		switch n.Op {
		case '+':
			return left + right
		case '-':
			return left - right
		case '*':
			return left * right
		case '/':
			return left / right
		case '^':
			return math.Pow(left, right)
		}
	case kindFuncCall:
		arg := evalArenaDepth(n.Left, x, depth+1)
		nameLen := 0
		for nameLen < len(n.Name) && n.Name[nameLen] != 0 {
			nameLen++
		}
		return callMathFunc(n.Name[:nameLen], arg)
	}
	return 0
}

// callMathFunc evaluates name(arg) using the shared math function table.
// It returns 0 for names that are not built-in math functions.
func callMathFunc(name []byte, arg float64) float64 {
	result, _ := callMathFuncByName(string(name), arg)
	return result
}
