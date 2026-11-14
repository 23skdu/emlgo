package jit

// Node is a node in the parsed expression AST.
type Node interface {
	nodeSigil()
	String() string
}

type (
	// Number is a numeric literal node.
	Number struct{ Value float64 }
	// Variable is a reference to the input variable; Name="" defaults to "x".
	Variable struct{ Name string } // Name="" defaults to "x"
	// UnaryOp applies the unary operator Op to Operand (negation in practice).
	UnaryOp struct {
		Op      rune
		Operand Node
	}
	// BinaryOp is the binary operation Left Op Right, where Op is one of +, -, *, / or ^.
	BinaryOp struct {
		Left  Node
		Op    rune
		Right Node
	}
	// FunctionCall is a call to the named function Name with a single argument Arg.
	FunctionCall struct {
		Name string
		Arg  Node
	}
)

func (Number) nodeSigil()       { _ = 0 }
func (Variable) nodeSigil()     { _ = 0 }
func (UnaryOp) nodeSigil()      { _ = 0 }
func (BinaryOp) nodeSigil()     { _ = 0 }
func (FunctionCall) nodeSigil() { _ = 0 }
