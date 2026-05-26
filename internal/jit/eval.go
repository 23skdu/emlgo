package jit

import "math"

func Eval(n Node, x float64) float64 {
	switch v := n.(type) {
	case Number:
		return v.Value
	case Variable:
		return x
	case UnaryOp:
		return -Eval(v.Operand, x)
	case BinaryOp:
		left := Eval(v.Left, x)
		right := Eval(v.Right, x)
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
		arg := Eval(v.Arg, x)
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
		}
	}
	return 0
}
