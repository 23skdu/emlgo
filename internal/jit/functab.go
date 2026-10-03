package jit

import "math"

// mathFunc is a unary math function taking float64 and returning float64.
type mathFunc func(float64) float64

// mathFuncs is the single source of truth for the built-in math functions
// recognised by the parser. The JIT codegen table, the tree evaluator and the
// arena interpreter all dispatch through it, so the supported set can never
// drift apart between them.
var mathFuncs = map[string]mathFunc{
	"sin":   math.Sin,
	"cos":   math.Cos,
	"exp":   math.Exp,
	"log":   math.Log,
	"sqrt":  math.Sqrt,
	"tan":   math.Tan,
	"asin":  math.Asin,
	"acos":  math.Acos,
	"atan":  math.Atan,
	"abs":   math.Abs,
	"cbrt":  math.Cbrt,
	"log2":  math.Log2,
	"log10": math.Log10,
	"ceil":  math.Ceil,
	"floor": math.Floor,
	"trunc": math.Trunc,
	"round": math.Round,
	"sinh":  math.Sinh,
	"cosh":  math.Cosh,
	"tanh":  math.Tanh,
	"asinh": math.Asinh,
	"acosh": math.Acosh,
	"atanh": math.Atanh,
	"erf":   math.Erf,
	"gamma": math.Gamma,
}

// callMathFuncByName evaluates name(arg). It returns false if name is not a
// known math function.
func callMathFuncByName(name string, arg float64) (float64, bool) {
	fn, ok := mathFuncs[name]
	if !ok {
		return 0, false
	}
	return fn(arg), true
}
