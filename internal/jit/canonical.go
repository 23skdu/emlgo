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
		if n.Name != "" {
			return n.Name
		}
		return "x"
	case EMLOp:
		return fmt.Sprintf("eml(%s, %s)", n.Left, n.Right)
	case EMLFunc:
		return fmt.Sprintf("%s(%s)", n.Name, n.Left)
	}
	return "?"
}

// MaxASTDepth defines the maximum allowable recursion depth for AST traversals
// to prevent goroutine stack exhaustion on deeply nested expressions.
const MaxASTDepth = 1000

// EMLSize returns the total number of nodes in the expression tree rooted at n.
func EMLSize(n *EMLNode) int {
	return emlSizeDepth(n, 0)
}

func emlSizeDepth(n *EMLNode, depth int) int {
	if n == nil || depth > MaxASTDepth {
		return 0
	}
	return 1 + emlSizeDepth(n.Left, depth+1) + emlSizeDepth(n.Right, depth+1)
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
	return emlEvalDepth(n, x, 0)
}

func emlEvalDepth(n *EMLNode, x float64, depth int) float64 {
	if n == nil || depth > MaxASTDepth {
		return 0
	}
	switch n.Kind {
	case EMLConst:
		return n.Value
	case EMLVar:
		return x
	case EMLOp:
		return math.Exp(emlEvalDepth(n.Left, x, depth+1)) - math.Log(emlEvalDepth(n.Right, x, depth+1))
	case EMLFunc:
		arg := emlEvalDepth(n.Left, x, depth+1)
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
			return arg + emlEvalDepth(n.Right, x, depth+1)
		case "sub":
			return arg - emlEvalDepth(n.Right, x, depth+1)
		case "mul":
			return arg * emlEvalDepth(n.Right, x, depth+1)
		case "div":
			return arg / emlEvalDepth(n.Right, x, depth+1)
		case "pow":
			return math.Pow(arg, emlEvalDepth(n.Right, x, depth+1))
		case "sqrt":
			return math.Sqrt(arg)
		default:
			if fn, ok := callMathFuncByName(n.Name, arg); ok {
				return fn
			}
		}
	}
	return 0
}

// EMLEvalRegularized evaluates the canonical EML expression tree n using smooth regularization:
// ln_eps(y) = 0.5 * ln(y^2 + eps^2) and clamped exp.
// This prevents NaN cascading during symbolic fitting and genetic programming.
func EMLEvalRegularized(n *EMLNode, x float64, eps float64) float64 {
	return emlEvalRegularizedDepth(n, x, eps, 0)
}

func emlEvalRegularizedDepth(n *EMLNode, x float64, eps float64, depth int) float64 {
	if n == nil || depth > MaxASTDepth {
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
		left := emlEvalRegularizedDepth(n.Left, x, eps, depth+1)
		right := emlEvalRegularizedDepth(n.Right, x, eps, depth+1)
		if left > 700.0 {
			left = 700.0
		} else if left < -700.0 {
			return -0.5 * math.Log(right*right+eps*eps)
		}
		return math.Exp(left) - 0.5*math.Log(right*right+eps*eps)
	case EMLFunc:
		arg := emlEvalRegularizedDepth(n.Left, x, eps, depth+1)
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
			return arg + emlEvalRegularizedDepth(n.Right, x, eps, depth+1)
		case "sub":
			return arg - emlEvalRegularizedDepth(n.Right, x, eps, depth+1)
		case "mul":
			return arg * emlEvalRegularizedDepth(n.Right, x, eps, depth+1)
		case "div":
			denom := emlEvalRegularizedDepth(n.Right, x, eps, depth+1)
			return arg * denom / (denom*denom + eps*eps)
		case "pow":
			return math.Pow(arg, emlEvalRegularizedDepth(n.Right, x, eps, depth+1))
		case "sqrt":
			if arg < 0 {
				arg = -arg
			}
			return math.Sqrt(arg)
		default:
			if res, ok := callMathFuncByName(n.Name, arg); ok {
				if !math.IsNaN(res) {
					return res
				}
				if !math.IsNaN(arg) {
					return arg
				}
				return 0
			}
		}
	}
	return 0
}

// Canonicalize converts a Node interface value into the canonical EMLNode representation.
func Canonicalize(n Node) *EMLNode {
	return canonicalizeDepth(n, 0)
}

func canonicalizeDepth(n Node, depth int) *EMLNode {
	if n == nil || depth > MaxASTDepth {
		return nil
	}
	switch v := n.(type) {
	case Number:
		return constNode(v.Value)
	case Variable:
		name := v.Name
		if name == "" {
			name = "x"
		}
		return &EMLNode{Kind: EMLVar, Name: name}
	case UnaryOp:
		if v.Op == '-' {
			child := canonicalizeDepth(v.Operand, depth+1)
			return &EMLNode{Kind: EMLFunc, Name: "neg", Left: child}
		}
		return canonicalizeDepth(v.Operand, depth+1)
	case BinaryOp:
		left := canonicalizeDepth(v.Left, depth+1)
		right := canonicalizeDepth(v.Right, depth+1)
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
		child := canonicalizeDepth(v.Arg, depth+1)
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
	return depthWithLimit(n, 0)
}

func depthWithLimit(n *EMLNode, d int) int {
	if n == nil || d > MaxASTDepth {
		return 0
	}
	l, r := depthWithLimit(n.Left, d+1), depthWithLimit(n.Right, d+1)
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
	return simplifyDepth(n, 0)
}

func simplifyDepth(n *EMLNode, depth int) *EMLNode {
	if n == nil || depth > MaxASTDepth {
		return n
	}
	// Recursively simplify children first.
	l := simplifyDepth(n.Left, depth+1)
	r := simplifyDepth(n.Right, depth+1)
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
		}
		// Note: no negation or reciprocal reductions are performed here. The
		// candidate reductions (e.g. eml(eml(1, eml(x,1)), eml(1,1)) -> -x)
		// are not identities for eml(u,v) = exp(u) - ln(v), and they are in any
		// case unreachable here because the bottom-up recursion above has already
		// folded eml(1,1) to the constant e and eml(x,1) to exp(x).
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
			if l != nil && l.Kind == EMLConst && l.Value == 0 {
				return emlFuncUnary("neg", r)
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

// diffIsZeroDerivative reports whether name is a piecewise-constant function
// whose derivative is zero almost everywhere.
func diffIsZeroDerivative(name string) bool {
	switch name {
	case "ceil", "floor", "trunc", "round":
		return true
	}
	return false
}

// Diff symbolically differentiates the EMLNode tree with respect to "x" (the variable).
// The returned tree represents d(n)/dx in the EMLFunc representation.
func Diff(n *EMLNode) *EMLNode {
	return diffDepth(n, 0)
}

func diffDepth(n *EMLNode, depth int) *EMLNode {
	if n == nil || depth > MaxASTDepth {
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
		u_ := diffDepth(u, depth+1)
		v_ := diffDepth(v, depth+1)
		// exp(u) * u'
		term1 := emlFuncBinary("mul", emlFuncUnary("exp", u), u_)
		// v' / v
		term2 := emlFuncBinary("div", v_, v)
		return Simplify(emlFuncBinary("sub", term1, term2))
	case EMLFunc:
		arg := n.Left
		arg_ := diffDepth(arg, depth+1)
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
			return Simplify(emlFuncBinary("add", diffDepth(arg, depth+1), diffDepth(n.Right, depth+1)))
		case "sub":
			// d/dx (u - v) = u' - v'
			return Simplify(emlFuncBinary("sub", diffDepth(arg, depth+1), diffDepth(n.Right, depth+1)))
		case "mul":
			// product rule: u'v + uv'
			v := n.Right
			v_ := diffDepth(v, depth+1)
			return Simplify(emlFuncBinary("add",
				emlFuncBinary("mul", arg_, v),
				emlFuncBinary("mul", arg, v_)))
		case "div":
			// quotient rule: (u'v - uv') / v²
			v := n.Right
			v_ := diffDepth(v, depth+1)
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
			v_ := diffDepth(v, depth+1)
			term1 := emlFuncBinary("mul", v_, emlFuncUnary("log", arg))
			term2 := emlFuncBinary("div", emlFuncBinary("mul", v, arg_), arg)
			return Simplify(emlFuncBinary("mul", n, emlFuncBinary("add", term1, term2)))
		case "tan":
			// d/dx tan(u) = (1 + tan(u)^2) * u'
			tanU := emlFuncUnary("tan", arg)
			sec2 := emlFuncBinary("add", constNode(1), emlFuncBinary("pow", tanU, constNode(2)))
			return Simplify(emlFuncBinary("mul", sec2, arg_))
		case "sinh":
			// d/dx sinh(u) = cosh(u) * u'
			return Simplify(emlFuncBinary("mul", emlFuncUnary("cosh", arg), arg_))
		case "cosh":
			// d/dx cosh(u) = sinh(u) * u'
			return Simplify(emlFuncBinary("mul", emlFuncUnary("sinh", arg), arg_))
		case "tanh":
			// d/dx tanh(u) = (1 - tanh(u)^2) * u'
			tanhU := emlFuncUnary("tanh", arg)
			sech2 := emlFuncBinary("sub", constNode(1), emlFuncBinary("pow", tanhU, constNode(2)))
			return Simplify(emlFuncBinary("mul", sech2, arg_))
		case "asin":
			// d/dx asin(u) = u' / sqrt(1 - u^2)
			denom := emlFuncUnary("sqrt", emlFuncBinary("sub", constNode(1), emlFuncBinary("pow", arg, constNode(2))))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "acos":
			// d/dx acos(u) = -u' / sqrt(1 - u^2)
			denom := emlFuncUnary("sqrt", emlFuncBinary("sub", constNode(1), emlFuncBinary("pow", arg, constNode(2))))
			return Simplify(emlFuncUnary("neg", emlFuncBinary("div", arg_, denom)))
		case "atan":
			// d/dx atan(u) = u' / (1 + u^2)
			denom := emlFuncBinary("add", constNode(1), emlFuncBinary("pow", arg, constNode(2)))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "asinh":
			// d/dx asinh(u) = u' / sqrt(u^2 + 1)
			denom := emlFuncUnary("sqrt", emlFuncBinary("add", emlFuncBinary("pow", arg, constNode(2)), constNode(1)))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "acosh":
			// d/dx acosh(u) = u' / sqrt(u^2 - 1)
			denom := emlFuncUnary("sqrt", emlFuncBinary("sub", emlFuncBinary("pow", arg, constNode(2)), constNode(1)))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "atanh":
			// d/dx atanh(u) = u' / (1 - u^2)
			denom := emlFuncBinary("sub", constNode(1), emlFuncBinary("pow", arg, constNode(2)))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "abs":
			// d/dx |u| = (u / |u|) * u'
			sgn := emlFuncBinary("div", arg, emlFuncUnary("abs", arg))
			return Simplify(emlFuncBinary("mul", sgn, arg_))
		case "cbrt":
			// d/dx cbrt(u) = u' / (3 * cbrt(u)^2)
			denom := emlFuncBinary("mul", constNode(3), emlFuncBinary("pow", emlFuncUnary("cbrt", arg), constNode(2)))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "log2":
			// d/dx log2(u) = u' / (u * ln 2)
			denom := emlFuncBinary("mul", arg, constNode(math.Ln2))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "log10":
			// d/dx log10(u) = u' / (u * ln 10)
			denom := emlFuncBinary("mul", arg, constNode(math.Ln10))
			return Simplify(emlFuncBinary("div", arg_, denom))
		case "erf":
			// d/dx erf(u) = (2 / sqrt(pi)) * exp(-u^2) * u'
			const twoOverSqrtPi = 1.1283791670955125738961589
			u2 := emlFuncBinary("pow", arg, constNode(2))
			negU2 := emlFuncUnary("neg", u2)
			expNegU2 := emlFuncUnary("exp", negU2)
			return Simplify(emlFuncBinary("mul", constNode(twoOverSqrtPi), emlFuncBinary("mul", expNegU2, arg_)))
		case "gamma":
			// d/dx \Gamma(x) = \Gamma(x) * \psi(x) (where \psi is digamma).
			// Digamma is not implemented, so return NaN instead of an incorrect 0.
			return constNode(math.NaN())
		case "ceil", "floor", "trunc", "round":
			if diffIsZeroDerivative(n.Name) {
				return constNode(0)
			}
		}
	}
	return constNode(0)
}

// DiffEval evaluates the first derivative of expression n at x: EMLEval(Diff(n), x).
func DiffEval(n *EMLNode, x float64) float64 {
	return EMLEval(Diff(n), x)
}

// SecondDerivativeEval evaluates the second derivative of expression n at x: EMLEval(Diff(Diff(n)), x).
func SecondDerivativeEval(n *EMLNode, x float64) float64 {
	return EMLEval(Diff(Diff(n)), x)
}

// Equiv reports whether two EMLNode trees are structurally equivalent.
func Equiv(a, b *EMLNode) bool {
	return equivDepth(a, b, 0)
}

func equivDepth(a, b *EMLNode, depth int) bool {
	if depth > MaxASTDepth {
		return false
	}
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
		return equivDepth(a.Left, b.Left, depth+1) && equivDepth(a.Right, b.Right, depth+1)
	case EMLFunc:
		if a.Name != b.Name {
			return false
		}
		return equivDepth(a.Left, b.Left, depth+1) && equivDepth(a.Right, b.Right, depth+1)
	}
	return false
}
