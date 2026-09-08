package jit

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Decompile converts an EMLNode tree to a parenthesized infix string.
func Decompile(n *EMLNode) string {
	if n == nil {
		return ""
	}
	switch n.Kind {
	case EMLConst:
		return formatFloat(n.Value)
	case EMLVar:
		return "x"
	case EMLOp:
		left := Decompile(n.Left)
		right := Decompile(n.Right)
		return fmt.Sprintf("eml(%s, %s)", left, right)
	case EMLFunc:
		arg := Decompile(n.Left)
		switch n.Name {
		case "exp":
			return fmt.Sprintf("exp(%s)", arg)
		case "log":
			return fmt.Sprintf("log(%s)", arg)
		case "sin":
			return fmt.Sprintf("sin(%s)", arg)
		case "cos":
			return fmt.Sprintf("cos(%s)", arg)
		case "sqrt":
			return fmt.Sprintf("sqrt(%s)", arg)
		case "neg":
			return fmt.Sprintf("-%s", wrapParenDecomp(n.Left))
		case "add":
			rightArg := Decompile(n.Right)
			return fmt.Sprintf("%s + %s", arg, rightArg)
		case "sub":
			rightArg := Decompile(n.Right)
			return fmt.Sprintf("%s - %s", arg, rightArg)
		case "mul":
			rightArg := Decompile(n.Right)
			return fmt.Sprintf("%s * %s", arg, rightArg)
		case "div":
			rightArg := Decompile(n.Right)
			return fmt.Sprintf("%s / %s", arg, rightArg)
		case "pow":
			rightArg := Decompile(n.Right)
			return fmt.Sprintf("%s^%s", arg, rightArg)
		default:
			return fmt.Sprintf("%s(%s)", n.Name, arg)
		}
	}
	return ""
}

// DecompileLaTeX converts an EMLNode tree to LaTeX math mode output.
func DecompileLaTeX(n *EMLNode) string {
	if n == nil {
		return ""
	}
	switch n.Kind {
	case EMLConst:
		return formatFloatLatex(n.Value)
	case EMLVar:
		return "x"
	case EMLOp:
		left := DecompileLaTeX(n.Left)
		right := DecompileLaTeX(n.Right)
		return fmt.Sprintf("\\operatorname{eml}(%s, %s)", left, right)
	case EMLFunc:
		arg := DecompileLaTeX(n.Left)
		switch n.Name {
		case "exp":
			return fmt.Sprintf("e^{%s}", wrapLatex(arg))
		case "log":
			return fmt.Sprintf("\\ln(%s)", arg)
		case "sin":
			return fmt.Sprintf("\\sin(%s)", arg)
		case "cos":
			return fmt.Sprintf("\\cos(%s)", arg)
		case "tan":
			return fmt.Sprintf("\\tan(%s)", arg)
		case "sqrt":
			return fmt.Sprintf("\\sqrt{%s}", arg)
		case "neg":
			return fmt.Sprintf("-%s", wrapLatex(arg))
		case "add":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("%s + %s", arg, rightArg)
		case "sub":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("%s - %s", arg, rightArg)
		case "mul":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("%s \\cdot %s", wrapLatex(arg), wrapLatex(rightArg))
		case "div":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("\\frac{%s}{%s}", arg, rightArg)
		case "pow":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("%s^{%s}", wrapLatexPow(n.Left), rightArg)
		default:
			return fmt.Sprintf("\\operatorname{%s}(%s)", n.Name, arg)
		}
	}
	return ""
}

func formatFloat(v float64) string {
	if v == math.Trunc(v) && !math.IsInf(v, 0) && !math.IsNaN(v) {
		return strconv.FormatFloat(v, 'f', 1, 64)
	}
	s := strconv.FormatFloat(v, 'g', -1, 64)
	return s
}

func formatFloatLatex(v float64) string {
	if v == math.Trunc(v) && !math.IsInf(v, 0) && !math.IsNaN(v) {
		return strconv.FormatFloat(v, 'f', 1, 64)
	}
	s := strconv.FormatFloat(v, 'g', -1, 64)
	return s
}

func wrapParenDecomp(n *EMLNode) string {
	if n.Kind == EMLFunc && (n.Name == "add" || n.Name == "sub") {
		return "(" + Decompile(n) + ")"
	}
	return Decompile(n)
}

func wrapSubDecomp(n *EMLNode) string {
	if n.Kind == EMLFunc && (n.Name == "add" || n.Name == "sub") {
		return "(" + Decompile(n) + ")"
	}
	return Decompile(n)
}

func wrapSubDecompR(n *EMLNode) string {
	if n.Kind == EMLFunc && (n.Name == "add" || n.Name == "sub") {
		return "(" + Decompile(n) + ")"
	}
	return Decompile(n)
}

func wrapMulDecomp(n *EMLNode) string {
	if n.Kind == EMLFunc && (n.Name == "add" || n.Name == "sub") {
		return "(" + Decompile(n) + ")"
	}
	return Decompile(n)
}

func wrapDivDecomp(n *EMLNode) string {
	if n.Kind == EMLFunc && (n.Name == "add" || n.Name == "sub" || n.Name == "mul" || n.Name == "div") {
		return "(" + Decompile(n) + ")"
	}
	return Decompile(n)
}

func wrapPowerDecomp(n *EMLNode) string {
	if n.Kind == EMLOp || (n.Kind == EMLFunc && n.Name != "pow" && n.Name != "neg") {
		return "(" + Decompile(n) + ")"
	}
	return Decompile(n)
}

func wrapLatex(s string) string {
	if strings.Contains(s, "+") || strings.Contains(s, "-") {
		return "(" + s + ")"
	}
	return s
}

func wrapLatexPow(n *EMLNode) string {
	s := DecompileLaTeX(n)
	if n.Kind == EMLOp || (n.Kind == EMLFunc && n.Name != "pow") {
		return "(" + s + ")"
	}
	return s
}

// DecompileNodeToExpr converts an EMLNode to an infix expression using the JIT formatter.
func DecompileNodeToExpr(n *EMLNode) string {
	if n == nil {
		return ""
	}
	conv := emlNodeToJITNode(n)
	if conv == nil {
		return ""
	}
	return FormatExpr(conv)
}

func emlNodeToJITNode(n *EMLNode) Node {
	if n == nil {
		return nil
	}
	switch n.Kind {
	case EMLConst:
		return Number{Value: n.Value}
	case EMLVar:
		return Variable{Name: "x"}
	case EMLOp:
		left := emlNodeToJITNode(n.Left)
		right := emlNodeToJITNode(n.Right)
		return BinaryOp{
			Left:  FunctionCall{Name: "exp", Arg: left},
			Op:    '-',
			Right: FunctionCall{Name: "log", Arg: right},
		}
	case EMLFunc:
		arg := emlNodeToJITNode(n.Left)
		switch n.Name {
		case "add":
			return BinaryOp{Left: arg, Op: '+', Right: emlNodeToJITNode(n.Right)}
		case "sub":
			return BinaryOp{Left: arg, Op: '-', Right: emlNodeToJITNode(n.Right)}
		case "mul":
			return BinaryOp{Left: arg, Op: '*', Right: emlNodeToJITNode(n.Right)}
		case "div":
			return BinaryOp{Left: arg, Op: '/', Right: emlNodeToJITNode(n.Right)}
		case "pow":
			return BinaryOp{Left: arg, Op: '^', Right: emlNodeToJITNode(n.Right)}
		case "neg":
			return UnaryOp{Op: '-', Operand: arg}
		default:
			return FunctionCall{Name: n.Name, Arg: arg}
		}
	}
	return nil
}
