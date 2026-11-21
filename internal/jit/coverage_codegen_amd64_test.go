//go:build amd64

package jit

import (
	"math"
	"testing"
)

func TestCodegenEdgeCases(t *testing.T) {
	enc := &encoder{}
	idx1 := enc.addUint64Pool(12345)
	idx2 := enc.addUint64Pool(12345)
	if idx1 != idx2 {
		t.Fatalf("addUint64Pool duplicate index = %d, want %d", idx2, idx1)
	}

	// callFunc2 with registers >= 8
	enc.callFunc2(math.Pow, 8, 9)

	g := &generator{}
	if err := g.gen(Variable{Name: "y"}, 0); err == nil {
		t.Fatal("expected error for variable 'y'")
	}
	if err := g.gen(UnaryOp{Operand: Variable{Name: "y"}}, 0); err == nil {
		t.Fatal("expected error for UnaryOp with variable 'y'")
	}
	if err := g.gen(FunctionCall{Name: "unknown_func", Arg: Number{Value: 1}}, 0); err == nil {
		t.Fatal("expected error for unknown function")
	}
	if err := g.gen(FunctionCall{Name: "sin", Arg: Variable{Name: "y"}}, 0); err == nil {
		t.Fatal("expected error for sin(y)")
	}

	// Negative non-integer exponent
	if _, err := CompileCached("x ^ (-2.5)"); err != nil {
		t.Fatalf("CompileCached x ^ (-2.5) failed: %v", err)
	}
	// Exponents > 16 and < -16
	if _, err := CompileCached("x ^ 17"); err != nil {
		t.Fatalf("CompileCached x ^ 17 failed: %v", err)
	}
	if _, err := CompileCached("x ^ (-17)"); err != nil {
		t.Fatalf("CompileCached x ^ (-17) failed: %v", err)
	}

	// Generator pow error cases
	gPow := &generator{}
	if err := gPow.genPowWithSign(Variable{Name: "y"}, 2, true, 0); err == nil {
		t.Fatal("expected error for negative pow with base y")
	}
	if err := gPow.genPowCall(Variable{Name: "y"}, Number{Value: 17}, 0); err == nil {
		t.Fatal("expected error for genPowCall base y")
	}
	if err := gPow.genPowCall(Number{Value: 2}, Variable{Name: "y"}, 0); err == nil {
		t.Fatal("expected error for genPowCall exp y")
	}
	if err := gPow.genPowCall(Number{Value: 2}, Number{Value: 3}, 1); err != nil {
		t.Fatalf("genPowCall dst!=0 failed: %v", err)
	}
	if err := gPow.genPowUint(Variable{Name: "y"}, 3, 0); err == nil {
		t.Fatal("expected error for genPowUint base y")
	}

	// Register exhaustion
	gExhausted := &generator{spillDepth: maxSpillDepth}
	for i := 0; i < 14; i++ {
		_, _ = gExhausted.alloc()
	}
	if err := gExhausted.gen(UnaryOp{Operand: Number{Value: 1}}, 0); err == nil {
		t.Fatal("expected register exhaustion error in UnaryOp")
	}
	if err := gExhausted.gen(BinaryOp{Left: Number{Value: 1}, Op: '+', Right: Number{Value: 2}}, 0); err == nil {
		t.Fatal("expected register exhaustion error in BinaryOp")
	}
	if err := gExhausted.genPowCall(Number{Value: 1}, Number{Value: 2}, 0); err == nil {
		t.Fatal("expected register exhaustion error in genPowCall")
	}
	if err := gExhausted.genPowWithSign(Number{Value: 1}, 2, true, 0); err == nil {
		t.Fatal("expected register exhaustion error in genPowWithSign")
	}
	if err := gExhausted.genPowUint(Number{Value: 1}, 3, 0); err == nil {
		t.Fatal("expected register exhaustion error in genPowUint")
	}

	// Exhaust all except 1 register to hit binary op second alloc failure
	gOneLeft := &generator{spillDepth: maxSpillDepth}
	for i := 0; i < 13; i++ {
		_, _ = gOneLeft.alloc()
	}
	if err := gOneLeft.gen(BinaryOp{Left: Number{Value: 1}, Op: '+', Right: Number{Value: 2}}, 0); err == nil {
		t.Fatal("expected register exhaustion on second alloc in BinaryOp")
	}
	// Exhaust all except 1 register to hit genPowCall second alloc failure
	gPowOneLeft := &generator{spillDepth: maxSpillDepth}
	for i := 0; i < 13; i++ {
		_, _ = gPowOneLeft.alloc()
	}
	if err := gPowOneLeft.genPowCall(Number{Value: 1}, Number{Value: 2}, 0); err == nil {
		t.Fatal("expected register exhaustion on second alloc in genPowCall")
	}
}
