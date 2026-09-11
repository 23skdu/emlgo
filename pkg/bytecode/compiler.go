package bytecode

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/emlgo/eml/internal/jit"
)

// CompileAST compiles a jit.Node AST into a linear bytecode Program in SoA layout.
func CompileAST(n jit.Node) (*Program, error) {
	if n == nil {
		return nil, fmt.Errorf("cannot compile nil AST node")
	}
	p := NewProgram()
	varMap := make(map[string]uint16)

	if err := emitAST(n, p, varMap); err != nil {
		return nil, err
	}
	p.CalculateMaxStackDepth()
	return p, nil
}

func getOrAssignVar(name string, varMap map[string]uint16) (uint16, error) {
	if idx, exists := varMap[name]; exists {
		return idx, nil
	}
	if len(varMap) >= math.MaxUint16 {
		return 0, fmt.Errorf("too many distinct variables (exceeds %d)", math.MaxUint16)
	}
	// #nosec G115 - bounded by MaxUint16 check above
	idx := uint16(len(varMap))
	varMap[name] = idx
	return idx, nil
}

func emitAST(n jit.Node, p *Program, varMap map[string]uint16) error {
	switch v := n.(type) {
	case jit.Number:
		p.Ops = append(p.Ops, OpConst)
		p.Consts = append(p.Consts, v.Value)
		return nil

	case jit.Variable:
		name := v.Name
		if name == "" {
			name = "x"
		}
		idx, err := getOrAssignVar(name, varMap)
		if err != nil {
			return err
		}
		p.Ops = append(p.Ops, OpVar)
		p.VarIndices = append(p.VarIndices, idx)
		return nil

	case jit.UnaryOp:
		if err := emitAST(v.Operand, p, varMap); err != nil {
			return err
		}
		if v.Op == '-' {
			p.Ops = append(p.Ops, OpNeg)
			return nil
		}
		return fmt.Errorf("unsupported unary op: %c", v.Op)

	case jit.BinaryOp:
		if err := emitAST(v.Left, p, varMap); err != nil {
			return err
		}
		if err := emitAST(v.Right, p, varMap); err != nil {
			return err
		}
		switch v.Op {
		case '+':
			p.Ops = append(p.Ops, OpAdd)
		case '-':
			p.Ops = append(p.Ops, OpSub)
		case '*':
			p.Ops = append(p.Ops, OpMul)
		case '/':
			p.Ops = append(p.Ops, OpDiv)
		case '^':
			p.Ops = append(p.Ops, OpPow)
		default:
			return fmt.Errorf("unsupported binary op: %c", v.Op)
		}
		return nil

	case jit.FunctionCall:
		if err := emitAST(v.Arg, p, varMap); err != nil {
			return err
		}
		switch v.Name {
		case "exp":
			p.Ops = append(p.Ops, OpExp)
		case "log":
			p.Ops = append(p.Ops, OpLog)
		case "sqrt":
			p.Ops = append(p.Ops, OpSqrt)
		default:
			return fmt.Errorf("unsupported function in bytecode compiler: %s", v.Name)
		}
		return nil

	default:
		return fmt.Errorf("unknown AST node type: %T", n)
	}
}

// CompileEML compiles a canonical jit.EMLNode tree into a linear bytecode Program.
func CompileEML(n *jit.EMLNode) (*Program, error) {
	if n == nil {
		return nil, fmt.Errorf("cannot compile nil EMLNode")
	}
	p := NewProgram()
	varMap := make(map[string]uint16)

	if err := emitEML(n, p, varMap); err != nil {
		return nil, err
	}
	p.CalculateMaxStackDepth()
	return p, nil
}

func emitEML(n *jit.EMLNode, p *Program, varMap map[string]uint16) error {
	if n == nil {
		return nil
	}
	switch n.Kind {
	case jit.EMLConst:
		p.Ops = append(p.Ops, OpConst)
		p.Consts = append(p.Consts, n.Value)
		return nil

	case jit.EMLVar:
		name := n.Name
		if name == "" {
			name = "x"
		}
		idx, err := getOrAssignVar(name, varMap)
		if err != nil {
			return err
		}
		p.Ops = append(p.Ops, OpVar)
		p.VarIndices = append(p.VarIndices, idx)
		return nil

	case jit.EMLOp:
		// eml(x, y): evaluate left, evaluate right, then OpEML
		if err := emitEML(n.Left, p, varMap); err != nil {
			return err
		}
		if err := emitEML(n.Right, p, varMap); err != nil {
			return err
		}
		p.Ops = append(p.Ops, OpEML)
		return nil

	case jit.EMLFunc:
		if err := emitEML(n.Left, p, varMap); err != nil {
			return err
		}
		if n.Right != nil {
			if err := emitEML(n.Right, p, varMap); err != nil {
				return err
			}
		}
		switch n.Name {
		case "add":
			p.Ops = append(p.Ops, OpAdd)
		case "sub":
			p.Ops = append(p.Ops, OpSub)
		case "mul":
			p.Ops = append(p.Ops, OpMul)
		case "div":
			p.Ops = append(p.Ops, OpDiv)
		case "pow":
			p.Ops = append(p.Ops, OpPow)
		case "neg":
			p.Ops = append(p.Ops, OpNeg)
		case "exp":
			p.Ops = append(p.Ops, OpExp)
		case "log":
			p.Ops = append(p.Ops, OpLog)
		case "sqrt":
			p.Ops = append(p.Ops, OpSqrt)
		default:
			return fmt.Errorf("unsupported EML function: %s", n.Name)
		}
		return nil

	default:
		return fmt.Errorf("unknown EMLNode kind: %d", n.Kind)
	}
}

// CompileExpr parses a mathematical expression string directly into a bytecode Program
// using Dijkstra's Shunting-Yard algorithm.
// Supports +, -, *, /, ^, parentheses, functions (exp, log, sqrt, eml), and named variables.
func CompileExpr(expr string) (*Program, error) {
	node, err := jit.Parse(expr)
	if err == nil {
		return CompileAST(node)
	}
	// Fallback to direct Shunting-Yard parser supporting multi-arg functions like eml(x, y)
	return parseShuntingYard(expr)
}

// parseShuntingYard parses expressions directly into bytecode, supporting comma-separated multi-arg functions.
func parseShuntingYard(expr string) (*Program, error) {
	p := NewProgram()
	varMap := make(map[string]uint16)

	type opToken struct {
		isFunc   bool
		op       rune
		name     string
		preced   int
		rAssoc   bool
		argCount int
	}

	var opStack []opToken
	tokens := tokenize(expr)

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		switch {
		case tok.isNum:
			val, err := strconv.ParseFloat(tok.str, 64)
			if err != nil {
				return nil, err
			}
			p.Ops = append(p.Ops, OpConst)
			p.Consts = append(p.Consts, val)

		case tok.isIdent:
			lower := strings.ToLower(tok.str)
			if isFunc(lower) {
				opStack = append(opStack, opToken{isFunc: true, name: lower, preced: 10, argCount: 1})
			} else {
				idx, err := getOrAssignVar(tok.str, varMap)
				if err != nil {
					return nil, err
				}
				p.Ops = append(p.Ops, OpVar)
				p.VarIndices = append(p.VarIndices, idx)
			}

		case tok.str == ",":
			// Pop operators until '('
			for len(opStack) > 0 && opStack[len(opStack)-1].op != '(' {
				top := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if err := emitOp(p, top); err != nil {
					return nil, err
				}
			}
			if len(opStack) > 1 && opStack[len(opStack)-2].isFunc {
				opStack[len(opStack)-2].argCount++
			}

		case tok.str == "(":
			opStack = append(opStack, opToken{op: '('})

		case tok.str == ")":
			for len(opStack) > 0 && opStack[len(opStack)-1].op != '(' {
				top := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if err := emitOp(p, top); err != nil {
					return nil, err
				}
			}
			if len(opStack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			opStack = opStack[:len(opStack)-1] // Pop '('
			if len(opStack) > 0 && opStack[len(opStack)-1].isFunc {
				fn := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if err := emitFunc(p, fn); err != nil {
					return nil, err
				}
			}

		default:
			// Operator
			r := rune(tok.str[0])
			preced, rAssoc := getPrecedence(r)
			for len(opStack) > 0 && opStack[len(opStack)-1].op != '(' {
				top := opStack[len(opStack)-1]
				if (rAssoc && top.preced > preced) || (!rAssoc && top.preced >= preced) {
					opStack = opStack[:len(opStack)-1]
					if err := emitOp(p, top); err != nil {
						return nil, err
					}
				} else {
					break
				}
			}
			opStack = append(opStack, opToken{op: r, preced: preced, rAssoc: rAssoc})
		}
	}

	for len(opStack) > 0 {
		top := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]
		if top.op == '(' || top.op == ')' {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		if top.isFunc {
			if err := emitFunc(p, top); err != nil {
				return nil, err
			}
		} else {
			if err := emitOp(p, top); err != nil {
				return nil, err
			}
		}
	}

	p.CalculateMaxStackDepth()
	return p, nil
}

func emitOp(p *Program, tok struct {
	isFunc   bool
	op       rune
	name     string
	preced   int
	rAssoc   bool
	argCount int
}) error {
	switch tok.op {
	case '+':
		p.Ops = append(p.Ops, OpAdd)
	case '-':
		p.Ops = append(p.Ops, OpSub)
	case '*':
		p.Ops = append(p.Ops, OpMul)
	case '/':
		p.Ops = append(p.Ops, OpDiv)
	case '^':
		p.Ops = append(p.Ops, OpPow)
	default:
		return fmt.Errorf("unsupported operator: %c", tok.op)
	}
	return nil
}

func emitFunc(p *Program, tok struct {
	isFunc   bool
	op       rune
	name     string
	preced   int
	rAssoc   bool
	argCount int
}) error {
	switch tok.name {
	case "eml":
		p.Ops = append(p.Ops, OpEML)
	case "exp":
		p.Ops = append(p.Ops, OpExp)
	case "log", "ln":
		p.Ops = append(p.Ops, OpLog)
	case "sqrt":
		p.Ops = append(p.Ops, OpSqrt)
	default:
		return fmt.Errorf("unsupported function: %s", tok.name)
	}
	return nil
}

func isFunc(s string) bool {
	switch s {
	case "eml", "exp", "log", "ln", "sqrt", "sin", "cos":
		return true
	}
	return false
}

func getPrecedence(op rune) (int, bool) {
	switch op {
	case '+', '-':
		return 1, false
	case '*', '/':
		return 2, false
	case '^':
		return 3, true
	}
	return 0, false
}

type tokenInfo struct {
	str     string
	isNum   bool
	isIdent bool
}

func tokenize(expr string) []tokenInfo {
	var tokens []tokenInfo
	i := 0
	n := len(expr)

	for i < n {
		c := expr[i]
		if unicode.IsSpace(rune(c)) {
			i++
			continue
		}

		if (c >= '0' && c <= '9') || (c == '.' && i+1 < n && expr[i+1] >= '0' && expr[i+1] <= '9') {
			start := i
			for i < n && ((expr[i] >= '0' && expr[i] <= '9') || expr[i] == '.') {
				i++
			}
			tokens = append(tokens, tokenInfo{str: expr[start:i], isNum: true})
			continue
		}

		if unicode.IsLetter(rune(c)) || c == '_' {
			start := i
			for i < n && (unicode.IsLetter(rune(expr[i])) || (expr[i] >= '0' && expr[i] <= '9') || expr[i] == '_') {
				i++
			}
			tokens = append(tokens, tokenInfo{str: expr[start:i], isIdent: true})
			continue
		}

		tokens = append(tokens, tokenInfo{str: string(c)})
		i++
	}
	return tokens
}
