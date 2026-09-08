package jit

import (
	"math"
	"testing"
	"unsafe"
)

func TestArenaAlloc(t *testing.T) {
	a := NewArena(8)
	if a.NumNodes() != 0 {
		t.Fatalf("expected 0 nodes, got %d", a.NumNodes())
	}
	for i := 0; i < 10; i++ {
		a.Alloc()
	}
	if got := a.NumNodes(); got != 10 {
		t.Fatalf("expected 10 nodes, got %d", got)
	}
}

func TestArenaGrow(t *testing.T) {
	a := NewArena(2)
	for i := 0; i < 100; i++ {
		a.Alloc()
	}
	if got := a.NumNodes(); got != 100 {
		t.Fatalf("expected 100 nodes, got %d", got)
	}
}

func TestArenaReset(t *testing.T) {
	a := NewArena(8)
	for i := 0; i < 20; i++ {
		a.Alloc()
	}
	a.Reset()
	if got := a.NumNodes(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	n := a.Alloc()
	n.Kind = kindNumber
	n.Value = 42
	if got := a.NumNodes(); got != 1 {
		t.Fatalf("expected 1 after alloc, got %d", got)
	}
}

func TestArenaFromInterfaceNumber(t *testing.T) {
	a := NewArena(8)
	node := a.FromInterface(Number{Value: 3.14})
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Kind != kindNumber || node.Value != 3.14 {
		t.Fatalf("expected Number(3.14), got kind=%d value=%v", node.Kind, node.Value)
	}
}

func TestArenaFromInterfaceVariable(t *testing.T) {
	a := NewArena(8)
	node := a.FromInterface(Variable{})
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Kind != kindVariable {
		t.Fatalf("expected kindVariable, got %d", node.Kind)
	}
}

func TestArenaFromInterfaceUnaryOp(t *testing.T) {
	a := NewArena(8)
	node := a.FromInterface(UnaryOp{Op: '-', Operand: Number{Value: 5}})
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Kind != kindUnaryOp || node.Op != '-' {
		t.Fatalf("expected UnaryOp(-), got kind=%d op=%c", node.Kind, node.Op)
	}
	if node.Left == nil || node.Left.Kind != kindNumber || node.Left.Value != 5 {
		t.Fatal("expected operand Number(5)")
	}
}

func TestArenaFromInterfaceBinaryOp(t *testing.T) {
	a := NewArena(8)
	expr := BinaryOp{
		Left:  Variable{},
		Op:    '+',
		Right: Number{Value: 1},
	}
	node := a.FromInterface(expr)
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Kind != kindBinaryOp || node.Op != '+' {
		t.Fatalf("expected BinaryOp(+), got kind=%d op=%c", node.Kind, node.Op)
	}
	if node.Left == nil || node.Left.Kind != kindVariable {
		t.Fatal("expected Left=Variable")
	}
	if node.Right == nil || node.Right.Kind != kindNumber || node.Right.Value != 1 {
		t.Fatal("expected Right=Number(1)")
	}
}

func TestArenaFromInterfaceFuncCall(t *testing.T) {
	a := NewArena(8)
	expr := FunctionCall{Name: "sin", Arg: Variable{}}
	node := a.FromInterface(expr)
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.Kind != kindFuncCall {
		t.Fatalf("expected kindFuncCall, got %d", node.Kind)
	}
	name := ""
	for i := 0; i < len(node.Name) && node.Name[i] != 0; i++ {
		name += string(node.Name[i])
	}
	if name != "sin" {
		t.Fatalf("expected name 'sin', got '%s'", name)
	}
}

func TestArenaRoundTrip(t *testing.T) {
	exprs := []string{
		"42",
		"x",
		"-x",
		"x + 1",
		"2 * x - 3",
		"x^2 + 2*x + 1",
		"sin(x)",
		"sqrt(x^2 + 1)",
		"-sin(x) * cos(x)",
		"(x + 1) / (x - 1)",
		"exp(x) + log(x)",
		"x^3 + x^2 + x + 1",
	}
	for _, expr := range exprs {
		original, err := Parse(expr)
		if err != nil {
			t.Fatalf("parse %q: %v", expr, err)
		}
		a := NewArena(8)
		arenaNode := a.FromInterface(original)
		roundTripped := a.ToInterface(arenaNode)
		if roundTripped == nil {
			t.Fatalf("round-trip of %q produced nil", expr)
		}
		for _, x := range []float64{0, 1, 2, 3, -1, 0.5, 10} {
			got := Eval(roundTripped, x)
			want := Eval(original, x)
			if math.Abs(got-want) > 1e-14 {
				t.Errorf("%s at x=%v: round-trip Eval = %v, want %v", expr, x, got, want)
			}
		}
	}
}

func TestEvalArenaMatchesEval(t *testing.T) {
	tests := []struct {
		expr string
		xs   []float64
	}{
		{"42", []float64{0}},
		{"x", []float64{0, 1, -5, 100}},
		{"-x", []float64{0, 3, -7}},
		{"x + 1", []float64{0, 5, -1}},
		{"x - 1", []float64{0, 5, -1}},
		{"2 * x", []float64{0, 3, -4}},
		{"x / 2", []float64{0, 10, -4}},
		{"x^2", []float64{0, 3, -2}},
		{"x^3", []float64{0, 2, -3}},
		{"x^2 + 2*x + 1", []float64{0, 1, 2, 3, -1}},
		{"sin(x)", []float64{0, 1, math.Pi / 2}},
		{"cos(x)", []float64{0, 1, math.Pi}},
		{"exp(x)", []float64{0, 1, 2}},
		{"log(x)", []float64{1, 2, 10}},
		{"sqrt(x)", []float64{1, 4, 9, 16}},
		{"tan(x)", []float64{0, 1}},
		{"abs(-x)", []float64{1, 5, -3}},
		{"x^2 + sin(x) * cos(x)", []float64{0, 1, 2}},
		{"(x + 1) * (x - 1)", []float64{0, 2, 5}},
	}
	for _, tc := range tests {
		original, err := Parse(tc.expr)
		if err != nil {
			t.Fatalf("parse %q: %v", tc.expr, err)
		}
		a := NewArena(8)
		arenaNode := a.FromInterface(original)
		for _, x := range tc.xs {
			got := EvalArena(arenaNode, x)
			want := Eval(original, x)
			if math.Abs(got-want) > 1e-12 {
				t.Errorf("%s at x=%v: EvalArena = %v, Eval = %v", tc.expr, x, got, want)
			}
		}
	}
}

func TestEvalArenaZeroAllocations(t *testing.T) {
	a := NewArena(32)
	expr, err := Parse("x^2 + 2*x + 1")
	if err != nil {
		t.Fatal(err)
	}
	arenaNode := a.FromInterface(expr)
	x := 3.0
	allocs := testing.AllocsPerRun(100, func() {
		EvalArena(arenaNode, x)
	})
	if allocs != 0 {
		t.Errorf("EvalArena had %f allocations per run, want 0", allocs)
	}
}

func TestEvalArenaZeroAllocationsFuncCall(t *testing.T) {
	a := NewArena(32)
	expr, err := Parse("sin(x) + cos(x)")
	if err != nil {
		t.Fatal(err)
	}
	arenaNode := a.FromInterface(expr)
	x := 1.0
	allocs := testing.AllocsPerRun(100, func() {
		EvalArena(arenaNode, x)
	})
	if allocs != 0 {
		t.Errorf("EvalArena had %f allocations per run, want 0", allocs)
	}
}

func TestArenaNodeFixedSize(t *testing.T) {
	sz1 := unsafe.Sizeof(ArenaNode{})
	n1 := &ArenaNode{Kind: kindNumber, Value: 42}
	n2 := &ArenaNode{Kind: kindFuncCall}
	n2.Name[0] = 's'
	n2.Name[1] = 'i'
	n2.Name[2] = 'n'
	sz2 := unsafe.Sizeof(*n2)
	if sz1 != sz2 {
		t.Fatalf("ArenaNode sizes differ: %d vs %d", sz1, sz2)
	}
	_ = n1
}
