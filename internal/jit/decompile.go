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
	return wrapDecompDepth(n, parentOp, left, 0)
}

func wrapDecompDepth(n *EMLNode, parentOp string, left bool, depth int) string {
	s := decompileDepth(n, depth)
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
	return decompileDepth(n, 0)
}

func decompileDepth(n *EMLNode, depth int) string {
	if n == nil || depth > MaxASTDepth {
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
		left := decompileDepth(n.Left, depth+1)
		right := decompileDepth(n.Right, depth+1)
		return fmt.Sprintf("eml(%s, %s)", left, right)
	case EMLFunc:
		switch n.Name {
		case "exp":
			return fmt.Sprintf("exp(%s)", decompileDepth(n.Left, depth+1))
		case "log":
			return fmt.Sprintf("log(%s)", decompileDepth(n.Left, depth+1))
		case "sin":
			return fmt.Sprintf("sin(%s)", decompileDepth(n.Left, depth+1))
		case "cos":
			return fmt.Sprintf("cos(%s)", decompileDepth(n.Left, depth+1))
		case "sqrt":
			return fmt.Sprintf("sqrt(%s)", decompileDepth(n.Left, depth+1))
		case "neg":
			if emlPrec(n.Left) <= 2 {
				return fmt.Sprintf("-(%s)", decompileDepth(n.Left, depth+1))
			}
			return fmt.Sprintf("-%s", wrapDecompDepth(n.Left, "neg", true, depth+1))
		case "add":
			return fmt.Sprintf("%s + %s", wrapDecompDepth(n.Left, "add", true, depth+1), wrapDecompDepth(n.Right, "add", false, depth+1))
		case "sub":
			return fmt.Sprintf("%s - %s", wrapDecompDepth(n.Left, "sub", true, depth+1), wrapDecompDepth(n.Right, "sub", false, depth+1))
		case "mul":
			return fmt.Sprintf("%s * %s", wrapDecompDepth(n.Left, "mul", true, depth+1), wrapDecompDepth(n.Right, "mul", false, depth+1))
		case "div":
			return fmt.Sprintf("%s / %s", wrapDecompDepth(n.Left, "div", true, depth+1), wrapDecompDepth(n.Right, "div", false, depth+1))
		case "pow":
			return fmt.Sprintf("%s^%s", wrapDecompDepth(n.Left, "pow", true, depth+1), wrapDecompDepth(n.Right, "pow", false, depth+1))
		default:
			return fmt.Sprintf("%s(%s)", n.Name, decompileDepth(n.Left, depth+1))
		}
	}
	return ""
}

// DecompileLaTeX converts an EMLNode tree to LaTeX math mode output.
func DecompileLaTeX(n *EMLNode) string {
	return decompileLaTeXDepth(n, 0)
}

func decompileLaTeXDepth(n *EMLNode, depth int) string {
	if n == nil || depth > MaxASTDepth {
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
		left := decompileLaTeXDepth(n.Left, depth+1)
		right := decompileLaTeXDepth(n.Right, depth+1)
		return fmt.Sprintf("\\operatorname{eml}(%s, %s)", left, right)
	case EMLFunc:
		arg := decompileLaTeXDepth(n.Left, depth+1)
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
			rightArg := decompileLaTeXDepth(n.Right, depth+1)
			return fmt.Sprintf("%s + %s", arg, rightArg)
		case "sub":
			rightArg := decompileLaTeXDepth(n.Right, depth+1)
			return fmt.Sprintf("%s - %s", arg, rightArg)
		case "mul":
			rightArg := decompileLaTeXDepth(n.Right, depth+1)
			return fmt.Sprintf("%s \\cdot %s", wrapLatexNode(n.Left, arg), wrapLatexNode(n.Right, rightArg))
		case "div":
			rightArg := decompileLaTeXDepth(n.Right, depth+1)
			return fmt.Sprintf("\\frac{%s}{%s}", arg, rightArg)
		case "pow":
			rightArg := decompileLaTeXDepth(n.Right, depth+1)
			return fmt.Sprintf("%s^{%s}", wrapLatexPowDepth(n.Left, depth+1), rightArg)
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
	return wrapLatexPowDepth(n, 0)
}

func wrapLatexPowDepth(n *EMLNode, depth int) string {
	s := decompileLaTeXDepth(n, depth)
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
	return emlNodeToJITNodeDepth(n, 0)
}

func emlNodeToJITNodeDepth(n *EMLNode, depth int) Node {
	if n == nil || depth > MaxASTDepth {
		return nil
	}
	switch n.Kind {
	case EMLConst:
		return Number{Value: n.Value}
	case EMLVar:
		return Variable{Name: "x"}
	case EMLOp:
		left := emlNodeToJITNodeDepth(n.Left, depth+1)
		right := emlNodeToJITNodeDepth(n.Right, depth+1)
		return BinaryOp{
			Left:  FunctionCall{Name: "exp", Arg: left},
			Op:    '-',
			Right: FunctionCall{Name: "log", Arg: right},
		}
	case EMLFunc:
		arg := emlNodeToJITNodeDepth(n.Left, depth+1)
		switch n.Name {
		case "add":
			return BinaryOp{Left: arg, Op: '+', Right: emlNodeToJITNodeDepth(n.Right, depth+1)}
		case "sub":
			return BinaryOp{Left: arg, Op: '-', Right: emlNodeToJITNodeDepth(n.Right, depth+1)}
		case "mul":
			return BinaryOp{Left: arg, Op: '*', Right: emlNodeToJITNodeDepth(n.Right, depth+1)}
		case "div":
			return BinaryOp{Left: arg, Op: '/', Right: emlNodeToJITNodeDepth(n.Right, depth+1)}
		case "pow":
			return BinaryOp{Left: arg, Op: '^', Right: emlNodeToJITNodeDepth(n.Right, depth+1)}
		case "neg":
			return UnaryOp{Op: '-', Operand: arg}
		default:
			return FunctionCall{Name: n.Name, Arg: arg}
		}
	}
	return nil
}
