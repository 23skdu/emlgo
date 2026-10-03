package jit

import "math"

// Eval evaluates a Node tree with variable x.
func Eval(n Node, x float64) float64 {
	return EvalVars(n, map[string]float64{"x": x, "": x})
}

// EvalVars evaluates a Node tree with named variable bindings.
// Variable nodes look up their Name in vars; unknown names evaluate to 0.
func EvalVars(n Node, vars map[string]float64) float64 {
	return evalVarsDepth(n, vars, 0)
}

func evalVarsDepth(n Node, vars map[string]float64, depth int) float64 {
	if n == nil || depth > MaxASTDepth {
		return 0
	}
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
		return -evalVarsDepth(v.Operand, vars, depth+1)
	case BinaryOp:
		left := evalVarsDepth(v.Left, vars, depth+1)
		right := evalVarsDepth(v.Right, vars, depth+1)
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
		arg := evalVarsDepth(v.Arg, vars, depth+1)
		if result, ok := callMathFuncByName(v.Name, arg); ok {
			return result
		}
	}
	return 0
}
