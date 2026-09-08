package jit

import (
	"fmt"
	"math"
)

// EMLNode is a node in the EML expression tree.
type EMLNode struct {
	Kind  EMLNodeKind
	Value float64
	Name  string
	Left  *EMLNode
	Right *EMLNode
}

type EMLNodeKind uint8

const (
	EMLConst EMLNodeKind = iota
	EMLVar
	EMLOp
	EMLFunc
)

func (n *EMLNode) String() string {
	if n == nil {
		return "<nil>"
	}
	switch n.Kind {
	case EMLConst:
		return fmt.Sprintf("%v", n.Value)
	case EMLVar:
		return "x"
	case EMLOp:
		return fmt.Sprintf("eml(%s, %s)", n.Left, n.Right)
	case EMLFunc:
		return fmt.Sprintf("%s(%s)", n.Name, n.Left)
	}
	return "?"
}

func EMLSize(n *EMLNode) int {
	if n == nil {
		return 0
	}
	return 1 + EMLSize(n.Left) + EMLSize(n.Right)
}

func emlNode(x, y *EMLNode) *EMLNode {
	return &EMLNode{Kind: EMLOp, Left: x, Right: y}
}

func constNode(v float64) *EMLNode {
	return &EMLNode{Kind: EMLConst, Value: v}
}

func varNode() *EMLNode {
	return &EMLNode{Kind: EMLVar}
}

func CanonicalExp(x *EMLNode) *EMLNode {
	return emlNode(x, constNode(1))
}

func CanonicalLog(x *EMLNode) *EMLNode {
	return emlNode(constNode(1), emlNode(emlNode(constNode(1), x), constNode(1)))
}

func CanonicalSin(x *EMLNode) *EMLNode {
	return &EMLNode{Kind: EMLFunc, Name: "sin", Left: x}
}

func CanonicalCos(x *EMLNode) *EMLNode {
	return &EMLNode{Kind: EMLFunc, Name: "cos", Left: x}
}

func CanonicalSqrt(x *EMLNode) *EMLNode {
	return CanonicalExp(&EMLNode{
		Kind:  EMLFunc,
		Name:  "mul",
		Left:  constNode(0.5),
		Right: CanonicalLog(x),
	})
}

func EMLEval(n *EMLNode, x float64) float64 {
	if n == nil {
		return 0
	}
	switch n.Kind {
	case EMLConst:
		return n.Value
	case EMLVar:
		return x
	case EMLOp:
		return math.Exp(EMLEval(n.Left, x)) - math.Log(EMLEval(n.Right, x))
	case EMLFunc:
		arg := EMLEval(n.Left, x)
		switch n.Name {
		case "sin":
			return math.Sin(arg)
		case "cos":
			return math.Cos(arg)
		case "neg":
			return -arg
		case "add":
			return arg + EMLEval(n.Right, x)
		case "sub":
			return arg - EMLEval(n.Right, x)
		case "mul":
			return arg * EMLEval(n.Right, x)
		case "div":
			return arg / EMLEval(n.Right, x)
		case "pow":
			return math.Pow(arg, EMLEval(n.Right, x))
		case "sqrt":
			return math.Sqrt(arg)
		}
	}
	return 0
}

func Canonicalize(n Node) *EMLNode {
	if n == nil {
		return nil
	}
	switch v := n.(type) {
	case Number:
		return constNode(v.Value)
	case Variable:
		return varNode()
	case UnaryOp:
		if v.Op == '-' {
			child := Canonicalize(v.Operand)
			return &EMLNode{Kind: EMLFunc, Name: "neg", Left: child}
		}
		return Canonicalize(v.Operand)
	case BinaryOp:
		left := Canonicalize(v.Left)
		right := Canonicalize(v.Right)
		switch v.Op {
		case '+':
			return &EMLNode{Kind: EMLFunc, Name: "add", Left: left, Right: right}
		case '-':
			return &EMLNode{Kind: EMLFunc, Name: "sub", Left: left, Right: right}
		case '*':
			return &EMLNode{Kind: EMLFunc, Name: "mul", Left: left, Right: right}
		case '/':
			return &EMLNode{Kind: EMLFunc, Name: "div", Left: left, Right: right}
		case '^':
			return &EMLNode{Kind: EMLFunc, Name: "pow", Left: left, Right: right}
		}
	case FunctionCall:
		child := Canonicalize(v.Arg)
		switch v.Name {
		case "exp":
			return CanonicalExp(child)
		case "log":
			return CanonicalLog(child)
		case "sqrt":
			return CanonicalSqrt(child)
		default:
			return &EMLNode{Kind: EMLFunc, Name: v.Name, Left: child}
		}
	}
	return nil
}

func Equiv(a, b *EMLNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Kind != b.Kind {
		return false
	}
	switch a.Kind {
	case EMLConst:
		return a.Value == b.Value
	case EMLVar:
		return true
	case EMLOp:
		return Equiv(a.Left, b.Left) && Equiv(a.Right, b.Right)
	case EMLFunc:
		if a.Name != b.Name {
			return false
		}
		return Equiv(a.Left, b.Left) && Equiv(a.Right, b.Right)
	}
	return false
}
