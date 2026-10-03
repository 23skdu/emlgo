package jit

import (
	"fmt"
	"math"
	"strconv"
)

func emlOpRune(name string) rune {
	switch name {
	case "add":
		return '+'
	case "sub":
		return '-'
	case "mul":
		return '*'
	case "div":
		return '/'
	case "pow":
		return '^'
	}
	return 0
}

func emlPrec(n *EMLNode) int {
	if n == nil {
		return 5
	}
	if n.Kind == EMLFunc {
		switch n.Name {
		case "add", "sub":
			return 1
		case "mul", "div":
			return 2
		case "pow":
			return 3
		case "neg":
			return 4
		}
	}
	return 5
}

func wrapDecomp(n *EMLNode, parentOp string, left bool) string {
	s := Decompile(n)
	if n == nil {
		return s
	}
	parentRune := emlOpRune(parentOp)
	childPrec := emlPrec(n)
	parentPrec := 0
	switch parentRune {
	case '+', '-':
		parentPrec = 1
	case '*', '/':
		parentPrec = 2
	case '^':
		parentPrec = 3
	}

	if childPrec < parentPrec {
		return "(" + s + ")"
	}
	if childPrec == parentPrec && !left && (parentRune == '-' || parentRune == '/') {
		return "(" + s + ")"
	}
	if parentRune == '^' {
		if left && childPrec <= parentPrec {
			return "(" + s + ")"
		}
		if !left && (childPrec < parentPrec || (n.Kind == EMLFunc && n.Name == "neg")) {
			return "(" + s + ")"
		}
	}
	return s
}

// Decompile converts an EMLNode tree to a parenthesized infix string with correct operator precedence.
func Decompile(n *EMLNode) string {
	if n == nil {
		return ""
	}
	switch n.Kind {
	case EMLConst:
		return formatFloat(n.Value)
	case EMLVar:
		if n.Name != "" {
			return n.Name
		}
		return "x"
	case EMLOp:
		left := Decompile(n.Left)
		right := Decompile(n.Right)
		return fmt.Sprintf("eml(%s, %s)", left, right)
	case EMLFunc:
		switch n.Name {
		case "exp":
			return fmt.Sprintf("exp(%s)", Decompile(n.Left))
		case "log":
			return fmt.Sprintf("log(%s)", Decompile(n.Left))
		case "sin":
			return fmt.Sprintf("sin(%s)", Decompile(n.Left))
		case "cos":
			return fmt.Sprintf("cos(%s)", Decompile(n.Left))
		case "sqrt":
			return fmt.Sprintf("sqrt(%s)", Decompile(n.Left))
		case "neg":
			if emlPrec(n.Left) <= 2 {
				return fmt.Sprintf("-(%s)", Decompile(n.Left))
			}
			return fmt.Sprintf("-%s", wrapDecomp(n.Left, "neg", true))
		case "add":
			return fmt.Sprintf("%s + %s", wrapDecomp(n.Left, "add", true), wrapDecomp(n.Right, "add", false))
		case "sub":
			return fmt.Sprintf("%s - %s", wrapDecomp(n.Left, "sub", true), wrapDecomp(n.Right, "sub", false))
		case "mul":
			return fmt.Sprintf("%s * %s", wrapDecomp(n.Left, "mul", true), wrapDecomp(n.Right, "mul", false))
		case "div":
			return fmt.Sprintf("%s / %s", wrapDecomp(n.Left, "div", true), wrapDecomp(n.Right, "div", false))
		case "pow":
			return fmt.Sprintf("%s^%s", wrapDecomp(n.Left, "pow", true), wrapDecomp(n.Right, "pow", false))
		default:
			return fmt.Sprintf("%s(%s)", n.Name, Decompile(n.Left))
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
		if n.Name != "" {
			return n.Name
		}
		return "x"
	case EMLOp:
		left := DecompileLaTeX(n.Left)
		right := DecompileLaTeX(n.Right)
		return fmt.Sprintf("\\operatorname{eml}(%s, %s)", left, right)
	case EMLFunc:
		arg := DecompileLaTeX(n.Left)
		switch n.Name {
		case "exp":
			return fmt.Sprintf("e^{%s}", wrapLatexNode(n.Left, arg))
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
			return fmt.Sprintf("-%s", wrapLatexNode(n.Left, arg))
		case "add":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("%s + %s", arg, rightArg)
		case "sub":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("%s - %s", arg, rightArg)
		case "mul":
			rightArg := DecompileLaTeX(n.Right)
			return fmt.Sprintf("%s \\cdot %s", wrapLatexNode(n.Left, arg), wrapLatexNode(n.Right, rightArg))
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

func wrapLatexNode(n *EMLNode, s string) string {
	if n != nil && emlPrec(n) <= 2 {
		return "(" + s + ")"
	}
	return s
}

func wrapLatexPow(n *EMLNode) string {
	s := DecompileLaTeX(n)
	if n != nil && (n.Kind == EMLOp || (n.Kind == EMLFunc && emlPrec(n) < 5)) {
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
