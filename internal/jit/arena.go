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

func NewArena(capacity int) *Arena {
	if capacity <= 0 {
		capacity = 64
	}
	return &Arena{
		buf: make([]ArenaNode, capacity),
	}
}

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

func (a *Arena) NumNodes() int {
	return a.off
}

func (a *Arena) Reset() {
	a.off = 0
}

func (a *Arena) ToInterface(n *ArenaNode) Node {
	if n == nil {
		return nil
	}
	switch n.Kind {
	case kindNumber:
		return Number{Value: n.Value}
	case kindVariable:
		return Variable{}
	case kindUnaryOp:
		return UnaryOp{Op: n.Op, Operand: a.ToInterface(n.Left)}
	case kindBinaryOp:
		return BinaryOp{Left: a.ToInterface(n.Left), Op: n.Op, Right: a.ToInterface(n.Right)}
	case kindFuncCall:
		name := ""
		for i := 0; i < len(n.Name) && n.Name[i] != 0; i++ {
			name += string(n.Name[i])
		}
		return FunctionCall{Name: name, Arg: a.ToInterface(n.Left)}
	}
	return nil
}

func (a *Arena) FromInterface(n Node) *ArenaNode {
	if n == nil {
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
		node.Left = a.FromInterface(v.Operand)
		return node
	case BinaryOp:
		node := a.Alloc()
		node.Kind = kindBinaryOp
		node.Op = v.Op
		node.Left = a.FromInterface(v.Left)
		node.Right = a.FromInterface(v.Right)
		return node
	case FunctionCall:
		node := a.Alloc()
		node.Kind = kindFuncCall
		for i := 0; i < len(v.Name) && i < 15; i++ {
			node.Name[i] = v.Name[i]
		}
		node.Left = a.FromInterface(v.Arg)
		return node
	}
	return nil
}

func EvalArena(n *ArenaNode, x float64) float64 {
	if n == nil {
		return 0
	}
	switch n.Kind {
	case kindNumber:
		return n.Value
	case kindVariable:
		return x
	case kindUnaryOp:
		return -EvalArena(n.Left, x)
	case kindBinaryOp:
		left := EvalArena(n.Left, x)
		right := EvalArena(n.Right, x)
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
		arg := EvalArena(n.Left, x)
		nameLen := 0
		for nameLen < len(n.Name) && n.Name[nameLen] != 0 {
			nameLen++
		}
		return callMathFunc(n.Name[:nameLen], arg)
	}
	return 0
}

func callMathFunc(name []byte, arg float64) float64 {
	switch {
	case equalStr(name, "sin"):
		return math.Sin(arg)
	case equalStr(name, "cos"):
		return math.Cos(arg)
	case equalStr(name, "exp"):
		return math.Exp(arg)
	case equalStr(name, "log"):
		return math.Log(arg)
	case equalStr(name, "sqrt"):
		return math.Sqrt(arg)
	case equalStr(name, "tan"):
		return math.Tan(arg)
	case equalStr(name, "asin"):
		return math.Asin(arg)
	case equalStr(name, "acos"):
		return math.Acos(arg)
	case equalStr(name, "atan"):
		return math.Atan(arg)
	case equalStr(name, "abs"):
		return math.Abs(arg)
	case equalStr(name, "cbrt"):
		return math.Cbrt(arg)
	case equalStr(name, "log2"):
		return math.Log2(arg)
	case equalStr(name, "log10"):
		return math.Log10(arg)
	case equalStr(name, "ceil"):
		return math.Ceil(arg)
	case equalStr(name, "floor"):
		return math.Floor(arg)
	case equalStr(name, "trunc"):
		return math.Trunc(arg)
	}
	return 0
}

func equalStr(b []byte, s string) bool {
	if len(b) != len(s) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if b[i] != s[i] {
			return false
		}
	}
	return true
}
