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

// String returns a human-readable representation of the EMLNode.
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

// EMLSize returns the total number of nodes in the expression tree rooted at n.
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

// CanonicalExp returns an EMLNode representing exp(x) = eml(x, 1).
func CanonicalExp(x *EMLNode) *EMLNode {
	return emlNode(x, constNode(1))
}

// CanonicalLog returns an EMLNode representing log(x) using the eml canonical form.
func CanonicalLog(x *EMLNode) *EMLNode {
	return emlNode(constNode(1), emlNode(emlNode(constNode(1), x), constNode(1)))
}

// CanonicalSin returns an EMLNode representing sin(x).
func CanonicalSin(x *EMLNode) *EMLNode {
	return &EMLNode{Kind: EMLFunc, Name: "sin", Left: x}
}

// CanonicalCos returns an EMLNode representing cos(x).
func CanonicalCos(x *EMLNode) *EMLNode {
	return &EMLNode{Kind: EMLFunc, Name: "cos", Left: x}
}

// CanonicalSqrt returns an EMLNode representing sqrt(x) via exp(0.5 * log(x)).
func CanonicalSqrt(x *EMLNode) *EMLNode {
	return CanonicalExp(&EMLNode{
		Kind:  EMLFunc,
		Name:  "mul",
		Left:  constNode(0.5),
		Right: CanonicalLog(x),
	})
}

// EMLEval evaluates the canonical EML expression tree n with variable x set to the given value.
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
		case "exp":
			return math.Exp(arg)
		case "log":
			return math.Log(arg)
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

// EMLEvalRegularized evaluates the canonical EML expression tree n using smooth regularization:
// ln_eps(y) = 0.5 * ln(y^2 + eps^2) and clamped exp.
// This prevents NaN cascading during symbolic fitting and genetic programming.
func EMLEvalRegularized(n *EMLNode, x float64, eps float64) float64 {
	if n == nil {
		return 0
	}
	if eps <= 0 {
		eps = 1e-12
	}
	switch n.Kind {
	case EMLConst:
		return n.Value
	case EMLVar:
		return x
	case EMLOp:
		left := EMLEvalRegularized(n.Left, x, eps)
		right := EMLEvalRegularized(n.Right, x, eps)
		if left > 700.0 {
			left = 700.0
		} else if left < -700.0 {
			return -0.5 * math.Log(right*right+eps*eps)
		}
		return math.Exp(left) - 0.5*math.Log(right*right+eps*eps)
	case EMLFunc:
		arg := EMLEvalRegularized(n.Left, x, eps)
		switch n.Name {
		case "sin":
			return math.Sin(arg)
		case "cos":
			return math.Cos(arg)
		case "exp":
			if arg > 700.0 {
				arg = 700.0
			}
			return math.Exp(arg)
		case "log":
			return 0.5 * math.Log(arg*arg+eps*eps)
		case "neg":
			return -arg
		case "add":
			return arg + EMLEvalRegularized(n.Right, x, eps)
		case "sub":
			return arg - EMLEvalRegularized(n.Right, x, eps)
		case "mul":
			return arg * EMLEvalRegularized(n.Right, x, eps)
		case "div":
			denom := EMLEvalRegularized(n.Right, x, eps)
			return arg * denom / (denom*denom + eps*eps)
		case "pow":
			return math.Pow(arg, EMLEvalRegularized(n.Right, x, eps))
		case "sqrt":
			if arg < 0 {
				arg = -arg
			}
			return math.Sqrt(arg)
		}
	}
	return 0
}

// Canonicalize converts a Node interface value into the canonical EMLNode representation.
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

// Depth returns the height of the EMLNode tree (leaf = 0).
func Depth(n *EMLNode) int {
	if n == nil {
		return 0
	}
	l, r := Depth(n.Left), Depth(n.Right)
	if l > r {
		return l + 1
	}
	return r + 1
}

// -- helpers used by Simplify / Diff --

func emlFuncBinary(name string, l, r *EMLNode) *EMLNode {
	return &EMLNode{Kind: EMLFunc, Name: name, Left: l, Right: r}
}

func emlFuncUnary(name string, arg *EMLNode) *EMLNode {
	return &EMLNode{Kind: EMLFunc, Name: name, Left: arg}
}

// Simplify performs constant folding and algebraic simplification on an EMLNode tree.
// The returned tree evaluates to the same value as n for all inputs.
func Simplify(n *EMLNode) *EMLNode {
	if n == nil {
		return nil
	}
	// Recursively simplify children first.
	l := Simplify(n.Left)
	r := Simplify(n.Right)
	node := &EMLNode{Kind: n.Kind, Value: n.Value, Name: n.Name, Left: l, Right: r}

	switch n.Kind {
	case EMLConst, EMLVar:
		return node

	case EMLOp:
		// eml(u, v) = exp(u) - ln(v) — fold when both children are constants.
		if l != nil && r != nil && l.Kind == EMLConst && r.Kind == EMLConst {
			if l.Value == 1.0 && r.Value == 1.0 {
				return constNode(math.E)
			}
			if l.Value == 0.0 && r.Value == 1.0 {
				return constNode(1.0)
			}
			return constNode(math.Exp(l.Value) - math.Log(r.Value))
		}
		// eml(x, 1) -> exp(x)
		if r != nil && r.Kind == EMLConst && r.Value == 1.0 {
			return &EMLNode{Kind: EMLFunc, Name: "exp", Left: l}
		}
		// eml(0, y) -> 1 - ln(y)
		if l != nil && l.Kind == EMLConst && l.Value == 0.0 {
			return emlFuncBinary("sub", constNode(1), &EMLNode{Kind: EMLFunc, Name: "log", Left: r})
		}
		// 1. eml(1, eml(eml(1, x), 1)) -> log(x)
		// Or after bottom-up reduction: eml(1, exp(eml(1, x))) -> log(x)
		if l != nil && l.Kind == EMLConst && l.Value == 1.0 && r != nil {
			if r.Kind == EMLFunc && r.Name == "exp" && r.Left != nil && r.Left.Kind == EMLOp {
				u := r.Left
				if u.Left != nil && u.Left.Kind == EMLConst && u.Left.Value == 1.0 && u.Right != nil {
					return &EMLNode{Kind: EMLFunc, Name: "log", Left: u.Right}
				}
			}
			if r.Kind == EMLOp && r.Right != nil && r.Right.Kind == EMLConst && r.Right.Value == 1.0 && r.Left != nil && r.Left.Kind == EMLOp {
				u := r.Left
				if u.Left != nil && u.Left.Kind == EMLConst && u.Left.Value == 1.0 && u.Right != nil {
					return &EMLNode{Kind: EMLFunc, Name: "log", Left: u.Right}
				}
			}
		}
		// 2. eml(eml(1, eml(x, 1)), eml(1, 1)) -> -x (neg(x))
		if l != nil && l.Kind == EMLOp && r != nil && r.Kind == EMLOp {
			if r.Left != nil && r.Left.Kind == EMLConst && r.Left.Value == 1.0 &&
				r.Right != nil && r.Right.Kind == EMLConst && r.Right.Value == 1.0 {
				if l.Left != nil && l.Left.Kind == EMLConst && l.Left.Value == 1.0 &&
					l.Right != nil && l.Right.Kind == EMLOp &&
					l.Right.Right != nil && l.Right.Right.Kind == EMLConst && l.Right.Right.Value == 1.0 {
					return &EMLNode{Kind: EMLFunc, Name: "neg", Left: l.Right.Left}
				}
			}
			// 3. eml(eml(1, x), eml(x, 1)) -> 1/x (div(1, x))
			if l.Left != nil && l.Left.Kind == EMLConst && l.Left.Value == 1.0 &&
				r.Right != nil && r.Right.Kind == EMLConst && r.Right.Value == 1.0 &&
				Equiv(l.Right, r.Left) {
				return emlFuncBinary("div", constNode(1), l.Right)
			}
		}
		return node

	case EMLFunc:
		switch n.Name {
		case "neg":
			if l != nil && l.Kind == EMLConst {
				return constNode(-l.Value)
			}
			// neg(neg(x)) = x
			if l != nil && l.Kind == EMLFunc && l.Name == "neg" {
				return l.Left
			}
		case "add":
			if l != nil && r != nil && l.Kind == EMLConst && r.Kind == EMLConst {
				return constNode(l.Value + r.Value)
			}
			if l != nil && l.Kind == EMLConst && l.Value == 0 {
				return r
			}
			if r != nil && r.Kind == EMLConst && r.Value == 0 {
				return l
			}
		case "sub":
			if l != nil && r != nil && l.Kind == EMLConst && r.Kind == EMLConst {
				return constNode(l.Value - r.Value)
			}
			if r != nil && r.Kind == EMLConst && r.Value == 0 {
				return l
			}
		case "mul":
			if l != nil && r != nil && l.Kind == EMLConst && r.Kind == EMLConst {
				return constNode(l.Value * r.Value)
			}
			if l != nil && l.Kind == EMLConst {
				if l.Value == 0 {
					return constNode(0)
				}
				if l.Value == 1 {
					return r
				}
			}
			if r != nil && r.Kind == EMLConst {
				if r.Value == 0 {
					return constNode(0)
				}
				if r.Value == 1 {
					return l
				}
			}
		case "div":
			if l != nil && r != nil && l.Kind == EMLConst && r.Kind == EMLConst && r.Value != 0 {
				return constNode(l.Value / r.Value)
			}
			if r != nil && r.Kind == EMLConst && r.Value == 1 {
				return l
			}
		case "pow":
			if l != nil && r != nil && l.Kind == EMLConst && r.Kind == EMLConst {
				return constNode(math.Pow(l.Value, r.Value))
			}
			if r != nil && r.Kind == EMLConst {
				if r.Value == 0 {
					return constNode(1)
				}
				if r.Value == 1 {
					return l
				}
			}
		case "sqrt":
			if l != nil && l.Kind == EMLConst {
				return constNode(math.Sqrt(l.Value))
			}
		}
		return node
	}
	return node
}

// Diff symbolically differentiates the EMLNode tree with respect to "x" (the variable).
// The returned tree represents d(n)/dx in the EMLFunc representation.
func Diff(n *EMLNode) *EMLNode {
	if n == nil {
		return constNode(0)
	}
	switch n.Kind {
	case EMLConst:
		return constNode(0)
	case EMLVar:
		return constNode(1)
	case EMLOp:
		// eml(u, v) = exp(u) - ln(v)
		// d/dx = exp(u)·u' - v'/v
		u, v := n.Left, n.Right
		u_ := Diff(u)
		v_ := Diff(v)
		// exp(u) * u'
		term1 := emlFuncBinary("mul", emlFuncUnary("exp", u), u_)
		// v' / v
		term2 := emlFuncBinary("div", v_, v)
		return Simplify(emlFuncBinary("sub", term1, term2))
	case EMLFunc:
		arg := n.Left
		arg_ := Diff(arg)
		switch n.Name {
		case "exp":
			// d/dx exp(u) = exp(u) * u'
			return Simplify(emlFuncBinary("mul", emlFuncUnary("exp", arg), arg_))
		case "log":
			// d/dx log(u) = u' / u
			return Simplify(emlFuncBinary("div", arg_, arg))
		case "sin":
			// d/dx sin(u) = cos(u) * u'
			return Simplify(emlFuncBinary("mul", emlFuncUnary("cos", arg), arg_))
		case "cos":
			// d/dx cos(u) = -sin(u) * u'
			return Simplify(emlFuncBinary("mul", emlFuncUnary("neg", emlFuncUnary("sin", arg)), arg_))
		case "sqrt":
			// d/dx sqrt(u) = u' / (2 * sqrt(u))
			denom := emlFuncBinary("mul", constNode(2), emlFuncUnary("sqrt", arg))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "neg":
			// d/dx -u = -u'
			return Simplify(emlFuncUnary("neg", arg_))
		case "add":
			// d/dx (u + v) = u' + v'
			return Simplify(emlFuncBinary("add", Diff(arg), Diff(n.Right)))
		case "sub":
			// d/dx (u - v) = u' - v'
			return Simplify(emlFuncBinary("sub", Diff(arg), Diff(n.Right)))
		case "mul":
			// product rule: u'v + uv'
			v := n.Right
			v_ := Diff(v)
			return Simplify(emlFuncBinary("add",
				emlFuncBinary("mul", arg_, v),
				emlFuncBinary("mul", arg, v_)))
		case "div":
			// quotient rule: (u'v - uv') / v²
			v := n.Right
			v_ := Diff(v)
			num := emlFuncBinary("sub",
				emlFuncBinary("mul", arg_, v),
				emlFuncBinary("mul", arg, v_))
			denom := emlFuncBinary("pow", v, constNode(2))
			return Simplify(emlFuncBinary("div", num, denom))
		case "pow":
			// power rule for constant exponent: d/dx u^c = c * u^(c-1) * u'
			v := n.Right
			if v != nil && v.Kind == EMLConst {
				c := v.Value
				return Simplify(emlFuncBinary("mul",
					constNode(c),
					emlFuncBinary("mul",
						emlFuncBinary("pow", arg, constNode(c-1)),
						arg_)))
			}
			// general: d/dx u^v = u^v * (v'*ln(u) + v*u'/u)
			v_ := Diff(v)
			term1 := emlFuncBinary("mul", v_, emlFuncUnary("log", arg))
			term2 := emlFuncBinary("div", emlFuncBinary("mul", v, arg_), arg)
			return Simplify(emlFuncBinary("mul", n,
				emlFuncBinary("add", term1, term2)))
		}
	}
	return constNode(0)
}

// DiffEval evaluates the symbolic derivative Diff(n) at x.
func DiffEval(n *EMLNode, x float64) float64 {
	return EMLEval(Diff(n), x)
}

// Equiv reports whether two EMLNode trees are structurally equivalent.
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
