package jit

import (
	"math"
	"testing"
)

// -- Plan 1: Diff / DiffEval tests --

// finiteDiff computes a numerical derivative via central differences.
func finiteDiff(n *EMLNode, x, h float64) float64 {
	return (EMLEval(n, x+h) - EMLEval(n, x-h)) / (2 * h)
}

func assertDiffApprox(t *testing.T, name string, n *EMLNode, x, tol float64) {
	t.Helper()
	want := finiteDiff(n, x, 1e-7)
	got := DiffEval(n, x)
	if math.Abs(got-want) > tol {
		t.Errorf("%s: DiffEval at x=%v = %v, finite diff = %v (diff=%v)", name, x, got, want, math.Abs(got-want))
	}
}

func TestDiffConst(t *testing.T) {
	n := constNode(5)
	d := Diff(n)
	if d.Kind != EMLConst || d.Value != 0 {
		t.Errorf("Diff(const) should be 0, got %v", d)
	}
	if DiffEval(n, 3) != 0 {
		t.Error("DiffEval(const, 3) should be 0")
	}
}

func TestDiffVar(t *testing.T) {
	n := varNode()
	d := Diff(n)
	if d.Kind != EMLConst || d.Value != 1 {
		t.Errorf("Diff(var) should be 1, got %v", d)
	}
}

func TestDiffExp(t *testing.T) {
	// d/dx exp(x) = exp(x)
	n := CanonicalExp(varNode())
	for _, x := range []float64{0, 1, -1, 2.5} {
		want := math.Exp(x)
		got := DiffEval(n, x)
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("DiffEval(exp(x), %v) = %v, want %v", x, got, want)
		}
	}
}

func TestDiffLog(t *testing.T) {
	// d/dx log(x) = 1/x
	n := CanonicalLog(varNode())
	for _, x := range []float64{0.5, 1, 2, math.E} {
		want := 1 / x
		got := DiffEval(n, x)
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("DiffEval(log(x), %v) = %v, want %v", x, got, want)
		}
	}
}

func TestDiffSin(t *testing.T) {
	// d/dx sin(x) = cos(x)
	n := emlFuncUnary("sin", varNode())
	for _, x := range []float64{0, 1, math.Pi / 4, -0.5} {
		want := math.Cos(x)
		got := DiffEval(n, x)
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("DiffEval(sin(x), %v) = %v, want %v", x, got, want)
		}
	}
}

func TestDiffCos(t *testing.T) {
	// d/dx cos(x) = -sin(x)
	n := emlFuncUnary("cos", varNode())
	for _, x := range []float64{0, 1, math.Pi / 4} {
		want := -math.Sin(x)
		got := DiffEval(n, x)
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("DiffEval(cos(x), %v) = %v, want %v", x, got, want)
		}
	}
}

func TestDiffSqrt(t *testing.T) {
	// d/dx sqrt(x) = 1/(2*sqrt(x))
	n := emlFuncUnary("sqrt", varNode())
	for _, x := range []float64{1, 4, 9, 0.25} {
		want := 1 / (2 * math.Sqrt(x))
		got := DiffEval(n, x)
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("DiffEval(sqrt(x), %v) = %v, want %v", x, got, want)
		}
	}
}

func TestDiffChainRule(t *testing.T) {
	// d/dx sin(exp(x)) = cos(exp(x)) * exp(x)
	n := emlFuncUnary("sin", emlFuncUnary("exp", varNode()))
	for _, x := range []float64{0, 0.5, 1} {
		assertDiffApprox(t, "sin(exp(x))", n, x, 1e-5)
	}
}

func TestDiffProductRule(t *testing.T) {
	// d/dx (x * sin(x)) = sin(x) + x*cos(x)
	n := emlFuncBinary("mul", varNode(), emlFuncUnary("sin", varNode()))
	for _, x := range []float64{0.5, 1, 2} {
		assertDiffApprox(t, "x*sin(x)", n, x, 1e-5)
	}
}

func TestDiffPowerConstant(t *testing.T) {
	// d/dx x^3 = 3*x^2
	n := emlFuncBinary("pow", varNode(), constNode(3))
	for _, x := range []float64{1, 2, -1, 0.5} {
		want := 3 * x * x
		got := DiffEval(n, x)
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("DiffEval(x^3, %v) = %v, want %v", x, got, want)
		}
	}
}

// -- Plan 2: Simplify / Depth tests --

func TestDepthLeaf(t *testing.T) {
	if Depth(constNode(1)) != 1 {
		t.Error("Depth(leaf) should be 1")
	}
	if Depth(varNode()) != 1 {
		t.Error("Depth(var) should be 1")
	}
	if Depth(nil) != 0 {
		t.Error("Depth(nil) should be 0")
	}
}

func TestDepthNested(t *testing.T) {
	n := CanonicalLog(varNode()) // eml(1, eml(eml(1,x),1)) — depth 3
	d := Depth(n)
	if d < 2 {
		t.Errorf("Depth of CanonicalLog should be >=2, got %d", d)
	}
}

func TestSimplifyConstFoldEML(t *testing.T) {
	// eml(0, 1) = exp(0) - log(1) = 1 - 0 = 1
	n := emlNode(constNode(0), constNode(1))
	s := Simplify(n)
	if s.Kind != EMLConst {
		t.Errorf("Simplify(eml(0,1)) should fold to const, got kind %d", s.Kind)
	}
	if math.Abs(s.Value-1) > 1e-12 {
		t.Errorf("Simplify(eml(0,1)) = %v, want 1", s.Value)
	}
}

func TestSimplifyNegConst(t *testing.T) {
	n := emlFuncUnary("neg", constNode(3))
	s := Simplify(n)
	if s.Kind != EMLConst || s.Value != -3 {
		t.Errorf("Simplify(neg(3)) = %v, want -3", s)
	}
}

func TestSimplifyNegNeg(t *testing.T) {
	// neg(neg(x)) = x
	n := emlFuncUnary("neg", emlFuncUnary("neg", varNode()))
	s := Simplify(n)
	if s.Kind != EMLVar {
		t.Errorf("Simplify(neg(neg(x))) should be var, got kind %d", s.Kind)
	}
}

func TestSimplifyMulConst(t *testing.T) {
	n := emlFuncBinary("mul", constNode(2), constNode(3))
	s := Simplify(n)
	if s.Kind != EMLConst || s.Value != 6 {
		t.Errorf("Simplify(mul(2,3)) = %v, want 6", s)
	}
}

func TestSimplifyMulOne(t *testing.T) {
	n := emlFuncBinary("mul", constNode(1), varNode())
	s := Simplify(n)
	if s.Kind != EMLVar {
		t.Errorf("Simplify(mul(1,x)) should be var, got kind %d", s.Kind)
	}
}

func TestSimplifyMulZero(t *testing.T) {
	n := emlFuncBinary("mul", constNode(0), varNode())
	s := Simplify(n)
	if s.Kind != EMLConst || s.Value != 0 {
		t.Errorf("Simplify(mul(0,x)) should be 0, got %v", s)
	}
}

func TestSimplifyAddZero(t *testing.T) {
	n := emlFuncBinary("add", constNode(0), varNode())
	s := Simplify(n)
	if s.Kind != EMLVar {
		t.Errorf("Simplify(add(0,x)) should be var, got kind %d", s.Kind)
	}
}

func TestSimplifyPowZero(t *testing.T) {
	n := emlFuncBinary("pow", varNode(), constNode(0))
	s := Simplify(n)
	if s.Kind != EMLConst || s.Value != 1 {
		t.Errorf("Simplify(pow(x,0)) should be 1, got %v", s)
	}
}

func TestSimplifyPowOne(t *testing.T) {
	n := emlFuncBinary("pow", varNode(), constNode(1))
	s := Simplify(n)
	if s.Kind != EMLVar {
		t.Errorf("Simplify(pow(x,1)) should be var, got kind %d", s.Kind)
	}
}

func TestSimplifyPreservesValue(t *testing.T) {
	// Property: Simplify(n) evaluates to same result as n for all x.
	testExprs := []*EMLNode{
		emlFuncBinary("mul", constNode(2), varNode()),
		emlFuncBinary("add", constNode(0), emlFuncUnary("sin", varNode())),
		emlFuncBinary("pow", varNode(), constNode(2)),
		emlFuncUnary("neg", emlFuncUnary("neg", varNode())),
		emlNode(constNode(0), constNode(1)),
	}
	xs := []float64{0.1, 0.5, 1.0, 1.5, 2.0, math.E, math.Pi / 4}
	for _, n := range testExprs {
		s := Simplify(n)
		for _, x := range xs {
			orig := EMLEval(n, x)
			simp := EMLEval(s, x)
			if math.IsNaN(orig) && math.IsNaN(simp) {
				continue
			}
			if math.Abs(orig-simp) > 1e-10 {
				t.Errorf("Simplify changed value: orig=%v simp=%v at x=%v for %v", orig, simp, x, n)
			}
		}
	}
}

func TestSimplifyDepthNonIncreasing(t *testing.T) {
	// Simplify should not increase depth.
	ns := []*EMLNode{
		emlFuncBinary("mul", constNode(1), emlFuncUnary("sin", varNode())),
		emlFuncBinary("add", constNode(0), varNode()),
		emlFuncBinary("pow", varNode(), constNode(1)),
	}
	for _, n := range ns {
		if Depth(Simplify(n)) > Depth(n) {
			t.Errorf("Simplify increased depth for %v", n)
		}
	}
}
