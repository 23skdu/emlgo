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
		switch v.Name {
		case "sin":
			return math.Sin(arg)
		case "cos":
			return math.Cos(arg)
		case "exp":
			return math.Exp(arg)
		case "log":
			return math.Log(arg)
		case "sqrt":
			return math.Sqrt(arg)
		case "tan":
			return math.Tan(arg)
		case "asin":
			return math.Asin(arg)
		case "acos":
			return math.Acos(arg)
		case "atan":
			return math.Atan(arg)
		case "abs":
			return math.Abs(arg)
		case "cbrt":
			return math.Cbrt(arg)
		case "log2":
			return math.Log2(arg)
		case "log10":
			return math.Log10(arg)
		case "ceil":
			return math.Ceil(arg)
		case "floor":
			return math.Floor(arg)
		case "trunc":
			return math.Trunc(arg)
		case "round":
			return math.Round(arg)
		}
	}
	return 0
}
