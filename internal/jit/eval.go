package jit

import "math"

// Eval evaluates a Node tree with variable x.
func Eval(n Node, x float64) float64 {
	return EvalVars(n, map[string]float64{"x": x, "": x})
}

// EvalVars evaluates a Node tree with named variable bindings.
// Variable nodes look up their Name in vars; unknown names evaluate to 0.
func EvalVars(n Node, vars map[string]float64) float64 {
	switch v := n.(type) {
	case Number:
		return v.Value
	case Variable:
		name := v.Name
		if name == "" {
			name = "x"
		}
		return vars[name]
	case UnaryOp:
		return -EvalVars(v.Operand, vars)
	case BinaryOp:
		left := EvalVars(v.Left, vars)
		right := EvalVars(v.Right, vars)
		switch v.Op {
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
	case FunctionCall:
		arg := EvalVars(v.Arg, vars)
		if result, ok := callMathFuncByName(v.Name, arg); ok {
			return result
		}
	}
	return 0
}
