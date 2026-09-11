package jit

import (
	"math"
	"testing"
)

func approxEq(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

func TestCanonicalExp(t *testing.T) {
	for _, x := range []float64{0, 1, -1, 2, 0.5} {
		node := CanonicalExp(constNode(x))
		got := EMLEval(node, 0)
		want := math.Exp(x)
		if !approxEq(got, want, 1e-12) {
			t.Errorf("CanonicalExp(%v): got %v, want %v", x, got, want)
		}
	}
	if node := CanonicalExp(varNode()); node.Kind != EMLOp {
		t.Fatalf("CanonicalExp root kind = %d, want EMLOp", node.Kind)
	}
}

func TestCanonicalLog(t *testing.T) {
	for _, x := range []float64{1, 2, 5, 10, 0.1} {
		node := CanonicalLog(constNode(x))
		got := EMLEval(node, 0)
		want := math.Log(x)
		if !approxEq(got, want, 1e-6) {
			t.Errorf("CanonicalLog(%v): got %v, want %v", x, got, want)
		}
	}
}

func TestEMLSize(t *testing.T) {
	if EMLSize(nil) != 0 {
		t.Fatal("EMLSize(nil) != 0")
	}
	if got := EMLSize(constNode(1)); got != 1 {
		t.Fatalf("EMLSize(const) = %d, want 1", got)
	}
	x := varNode()
	one := constNode(1)
	tree := emlNode(x, one)
	if got := EMLSize(tree); got != 3 {
		t.Fatalf("EMLSize(eml(x,1)) = %d, want 3", got)
	}
}

func TestCanonicalize(t *testing.T) {
	ast, err := Parse("exp(x)")
	if err != nil {
		t.Fatal(err)
	}
	tree := Canonicalize(ast)
	if tree == nil {
		t.Fatal("Canonicalize(exp(x)) returned nil")
	}
	if tree.Kind != EMLOp {
		t.Fatalf("CanonicalExp root kind = %d, want EMLOp", tree.Kind)
	}
	// eml(x,1): left=x (var), right=1 (const)
	if tree.Left.Kind != EMLVar {
		t.Errorf("left child kind = %d, want EMLVar", tree.Left.Kind)
	}
	if tree.Right.Kind != EMLConst || tree.Right.Value != 1 {
		t.Errorf("right child = %v, want const 1", tree.Right)
	}
	// Numerical check
	got := EMLEval(tree, 3)
	want := math.Exp(3)
	if !approxEq(got, want, 1e-12) {
		t.Errorf("EMLEval(Canonicalize(exp(x)), 3) = %v, want %v", got, want)
	}
}

func TestCanonicalizeBinary(t *testing.T) {
	ast, err := Parse("x + 1")
	if err != nil {
		t.Fatal(err)
	}
	tree := Canonicalize(ast)
	if tree == nil {
		t.Fatal("Canonicalize returned nil")
	}
	if tree.Kind != EMLFunc || tree.Name != "add" {
		t.Fatalf("root = %v, want add(...)", tree)
	}
	got := EMLEval(tree, 2)
	if !approxEq(got, 3, 1e-12) {
		t.Errorf("EMLEval = %v, want 3", got)
	}
}

func TestEquiv(t *testing.T) {
	a := emlNode(constNode(1), varNode())
	b := emlNode(constNode(1), varNode())
	if !Equiv(a, b) {
		t.Error("identical trees not equiv")
	}
	c := emlNode(constNode(2), varNode())
	if Equiv(a, c) {
		t.Error("different const trees equiv")
	}
	d := emlNode(varNode(), constNode(1))
	if Equiv(a, d) {
		t.Error("swapped operands equiv")
	}
	if !Equiv(nil, nil) {
		t.Error("nil not equiv nil")
	}
	if Equiv(nil, a) {
		t.Error("nil equiv non-nil")
	}
}

func TestEMLConstFold(t *testing.T) {
	tree := emlNode(constNode(1), constNode(1))
	got := EMLEval(tree, 0)
	want := math.E
	if !approxEq(got, want, 1e-12) {
		t.Errorf("eml(1,1) = %v, want e = %v", got, want)
	}
}

func TestCanonicalSqrt(t *testing.T) {
	for _, x := range []float64{1, 4, 9, 2} {
		node := CanonicalSqrt(constNode(x))
		got := EMLEval(node, 0)
		want := math.Sqrt(x)
		if !approxEq(got, want, 1e-6) {
			t.Errorf("CanonicalSqrt(%v): got %v, want %v", x, got, want)
		}
	}
}

func TestCanonicalizeUnaryNeg(t *testing.T) {
	ast, err := Parse("-x")
	if err != nil {
		t.Fatal(err)
	}
	tree := Canonicalize(ast)
	got := EMLEval(tree, 5)
	if !approxEq(got, -5, 1e-12) {
		t.Errorf("EMLEval(Canonicalize(-x), 5) = %v, want -5", got)
	}
}

func TestEMLNodeString(t *testing.T) {
	cases := []struct {
		node *EMLNode
		want string
	}{
		{nil, "<nil>"},
		{constNode(3), "3"},
		{varNode(), "x"},
		{emlNode(varNode(), constNode(1)), "eml(x, 1)"},
		{&EMLNode{Kind: EMLFunc, Name: "sin", Left: varNode()}, "sin(x)"},
	}
	for _, tc := range cases {
		got := tc.node.String()
		if got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestSimplifyEMLIdentities(t *testing.T) {
	// eml(1, 1) -> e
	eTree := emlNode(constNode(1), constNode(1))
	s1 := Simplify(eTree)
	if s1.Kind != EMLConst || math.Abs(s1.Value-math.E) > 1e-12 {
		t.Errorf("Simplify(eml(1, 1)): got %v, want e", s1)
	}

	// eml(x, 1) -> exp(x)
	expTree := emlNode(varNode(), constNode(1))
	s2 := Simplify(expTree)
	if s2.Kind != EMLFunc || s2.Name != "exp" {
		t.Errorf("Simplify(eml(x, 1)): got %v, want exp(x)", s2)
	}

	// eml(1, eml(eml(1, x), 1)) -> log(x)
	logTree := CanonicalLog(varNode())
	s3 := Simplify(logTree)
	if s3.Kind != EMLFunc || s3.Name != "log" {
		t.Errorf("Simplify(CanonicalLog): got %v, want log(x)", s3)
	}
}

func TestEMLEvalRegularized(t *testing.T) {
	// Even with non-positive values, EMLEvalRegularized should produce finite numbers (no NaNs)
	logTree := CanonicalLog(varNode())
	for _, x := range []float64{-10.0, -1.0, 0.0, 0.5, 1.0, 5.0} {
		val := EMLEvalRegularized(logTree, x, 1e-6)
		if math.IsNaN(val) || math.IsInf(val, 0) {
			t.Errorf("EMLEvalRegularized(log, %v) produced non-finite: %v", x, val)
		}
	}
}

