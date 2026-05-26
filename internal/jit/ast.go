package jit

type Node interface {
	nodeSigil()
	String() string
}

type (
	Number struct{ Value float64 }
	Variable struct{}
	UnaryOp struct {
		Op      rune
		Operand Node
	}
	BinaryOp struct {
		Left  Node
		Op    rune
		Right Node
	}
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
