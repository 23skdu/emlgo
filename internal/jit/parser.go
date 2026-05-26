package jit

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type lexer struct {
	input string
	pos   int
}

type token struct {
	typ   tokenType
	value string
}

type tokenType int

const (
	tokNumber tokenType = iota
	tokX
	tokPlus
	tokMinus
	tokStar
	tokSlash
	tokCaret
	tokLParen
	tokRParen
	tokFunction
	tokEOF
)

func isFuncName(s string) bool {
	switch s {
	case "sin", "cos", "exp", "log", "sqrt", "tan", "asin", "acos", "atan", "abs":
		return true
	}
	return false
}

func (l *lexer) next() token {
	for l.pos < len(l.input) && l.input[l.pos] == ' ' {
		l.pos++
	}
	if l.pos >= len(l.input) {
		return token{typ: tokEOF}
	}
	c := l.input[l.pos]
	if c == 'x' {
		l.pos++
		return token{typ: tokX, value: "x"}
	}
	if c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '(' || c == ')' {
		l.pos++
		return token{typ: tokTypeFromRune(rune(c)), value: string(c)}
	}
	if c >= '0' && c <= '9' || c == '.' {
		start := l.pos
		for l.pos < len(l.input) && (l.input[l.pos] >= '0' && l.input[l.pos] <= '9' || l.input[l.pos] == '.') {
			l.pos++
		}
		return token{typ: tokNumber, value: l.input[start:l.pos]}
	}
	if unicode.IsLetter(rune(c)) {
		start := l.pos
		for l.pos < len(l.input) && unicode.IsLetter(rune(l.input[l.pos])) {
			l.pos++
		}
		name := l.input[start:l.pos]
		if isFuncName(name) {
			return token{typ: tokFunction, value: name}
		}
		return token{typ: tokX, value: name}
	}
	return token{typ: tokEOF}
}

func tokTypeFromRune(c rune) tokenType {
	switch c {
	case '+':
		return tokPlus
	case '-':
		return tokMinus
	case '*':
		return tokStar
	case '/':
		return tokSlash
	case '^':
		return tokCaret
	case '(':
		return tokLParen
	case ')':
		return tokRParen
	}
	return tokEOF
}

type parser struct {
	l     *lexer
	peek  token
	ready bool
}

func newParser(input string) *parser {
	return &parser{l: &lexer{input: input}}
}

func (p *parser) next() token {
	if p.ready {
		p.ready = false
		return p.peek
	}
	return p.l.next()
}

func (p *parser) peekToken() token {
	if !p.ready {
		p.peek = p.l.next()
		p.ready = true
	}
	return p.peek
}

func (p *parser) consume(typ tokenType) (token, error) {
	t := p.next()
	if t.typ != typ {
		return t, fmt.Errorf("expected %d but got %s", typ, t.value)
	}
	return t, nil
}

func Parse(input string) (Node, error) {
	p := newParser(input)
	node, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peekToken().typ != tokEOF {
		return nil, fmt.Errorf("unexpected token: %s", p.peekToken().value)
	}
	return node, nil
}

func (p *parser) parseExpr() (Node, error) {
	return p.parseAddSub()
}

func (p *parser) parseAddSub() (Node, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peekToken()
		if t.typ != tokPlus && t.typ != tokMinus {
			break
		}
		p.next()
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		left = BinaryOp{Left: left, Op: rune(t.value[0]), Right: right}
	}
	return left, nil
}

func (p *parser) parseMulDiv() (Node, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peekToken()
		if t.typ != tokStar && t.typ != tokSlash {
			break
		}
		p.next()
		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		left = BinaryOp{Left: left, Op: rune(t.value[0]), Right: right}
	}
	return left, nil
}

func (p *parser) parsePower() (Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.peekToken().typ == tokCaret {
		p.next()
		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		left = BinaryOp{Left: left, Op: '^', Right: right}
	}
	return left, nil
}

func (p *parser) parseUnary() (Node, error) {
	t := p.peekToken()
	if t.typ == tokMinus {
		p.next()
		operand, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		return UnaryOp{Op: '-', Operand: operand}, nil
	}
	return p.parseBase()
}

func (p *parser) parseBase() (Node, error) {
	t := p.next()
	switch t.typ {
	case tokNumber:
		val, err := strconv.ParseFloat(t.value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", t.value)
		}
		return Number{Value: val}, nil
	case tokX:
		return Variable{}, nil
	case tokLParen:
		node, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokRParen)
		if err != nil {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		return node, nil
	case tokFunction:
		_, err := p.consume(tokLParen)
		if err != nil {
			return nil, fmt.Errorf("expected ( after function %s", t.value)
		}
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokRParen)
		if err != nil {
			return nil, fmt.Errorf("missing ) after function argument")
		}
		return FunctionCall{Name: t.value, Arg: arg}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", t.value)
	}
}

func (n Number) String() string       { return strconv.FormatFloat(n.Value, 'g', -1, 64) }
func (Variable) String() string       { return "x" }
func (n UnaryOp) String() string      { return fmt.Sprintf("(-%s)", n.Operand) }
func (n BinaryOp) String() string     { return fmt.Sprintf("(%s %c %s)", n.Left, n.Op, n.Right) }
func (n FunctionCall) String() string { return fmt.Sprintf("%s(%s)", n.Name, n.Arg) }
func FormatAST(n Node) string {
	return n.String()
}

func FormatExpr(n Node) string {
	switch v := n.(type) {
	case Number:
		s := strconv.FormatFloat(v.Value, 'f', -1, 64)
		if strings.ContainsRune(s, '.') {
			s = strings.TrimRight(s, "0")
			s = strings.TrimRight(s, ".")
		}
		if s == "" || s == "-" {
			s = "0"
		}
		return s
	case Variable:
		return "x"
	case UnaryOp:
		return fmt.Sprintf("-%s", wrapParen(v.Operand))
	case BinaryOp:
		l, r := wrapBinOp(v.Left, v.Op, true), wrapBinOp(v.Right, v.Op, false)
		if v.Op == '^' {
			return fmt.Sprintf("%s^%s", wrapPower(v.Left), wrapPower(v.Right))
		}
		return fmt.Sprintf("%s %c %s", l, v.Op, r)
	case FunctionCall:
		return fmt.Sprintf("%s(%s)", v.Name, FormatExpr(v.Arg))
	}
	return n.String()
}

func wrapParen(n Node) string {
	_, isBin := n.(BinaryOp)
	_, isUnary := n.(UnaryOp)
	if isBin || isUnary {
		return "(" + FormatExpr(n) + ")"
	}
	return FormatExpr(n)
}

func wrapBinOp(n Node, parentOp rune, left bool) string {
	bin, ok := n.(BinaryOp)
	if !ok {
		return wrapParen(n)
	}
	prec := func(op rune) int {
		switch op {
		case '+', '-':
			return 1
		case '*', '/':
			return 2
		case '^':
			return 3
		}
		return 0
	}
	if prec(bin.Op) < prec(parentOp) {
		return "(" + FormatExpr(n) + ")"
	}
	if prec(bin.Op) == prec(parentOp) && !left && (parentOp == '-' || parentOp == '/') {
		return "(" + FormatExpr(n) + ")"
	}
	if prec(bin.Op) == prec(parentOp) && parentOp == '^' {
		return "(" + FormatExpr(n) + ")"
	}
	return FormatExpr(n)
}

func wrapPower(n Node) string {
	_, isBin := n.(BinaryOp)
	_, isUnary := n.(UnaryOp)
	if isBin || isUnary {
		return "(" + FormatExpr(n) + ")"
	}
	return FormatExpr(n)
}
