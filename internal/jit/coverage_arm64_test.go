//go:build arm64

package jit

import (
	"fmt"
	"math"
	"testing"
)

func TestNodeSigilARM64(t *testing.T) {
	(Number{}).nodeSigil()
	(Variable{}).nodeSigil()
	(UnaryOp{}).nodeSigil()
	(BinaryOp{}).nodeSigil()
	(FunctionCall{}).nodeSigil()
}

func TestCodegenARM64EdgeCases100(t *testing.T) {
	var enc encoder
	enc.push16(0)
	enc.pop16(0)
	if len(enc.code) == 0 {
		t.Fatal("expected emitted arm64 instructions")
	}

	g := &generator{}
	if err := g.gen(Variable{Name: "y"}, 0); err == nil {
		t.Fatal("expected error for variable 'y'")
	}
	if err := g.gen(UnaryOp{Operand: Variable{Name: "y"}}, 0); err == nil {
		t.Fatal("expected error for UnaryOp with variable 'y'")
	}
	if err := g.gen(BinaryOp{Left: Variable{Name: "y"}, Op: '+', Right: Number{Value: 1}}, 0); err == nil {
		t.Fatal("expected error for BinaryOp left error")
	}
	if err := g.gen(BinaryOp{Left: Number{Value: 1}, Op: '+', Right: Variable{Name: "y"}}, 0); err == nil {
		t.Fatal("expected error for BinaryOp right error")
	}
	if err := g.gen(BinaryOp{Left: Number{Value: 1}, Op: '%', Right: Number{Value: 2}}, 0); err == nil {
		t.Fatal("expected error for unsupported operator")
	}
	if err := g.gen(FunctionCall{Name: "unknown_func", Arg: Number{Value: 1}}, 0); err == nil {
		t.Fatal("expected error for unknown function")
	}
	if err := g.gen(FunctionCall{Name: "sin", Arg: Variable{Name: "y"}}, 0); err == nil {
		t.Fatal("expected error for sin(y)")
	}
	if err := g.genFuncCall("sin", Number{Value: 1}, 2); err != nil {
		t.Fatalf("genFuncCall dst!=0 failed: %v", err)
	}
	if err := g.gen(testCustomNode{}, 0); err != nil {
		t.Fatalf("gen unknown node returned unexpected error: %v", err)
	}

	// Spill paths in BinaryOp (all operators + error paths + maxSpillDepth)
	for _, op := range []rune{'+', '-', '*', '/'} {
		gs := &generator{}
		for range 14 {
			_, _ = gs.alloc()
		}
		if err := gs.gen(BinaryOp{Left: Number{Value: 1}, Op: op, Right: Number{Value: 2}}, 0); err != nil {
			t.Fatalf("spill path for %c failed: %v", op, err)
		}
	}
	gsBadOp := &generator{}
	for range 14 {
		_, _ = gsBadOp.alloc()
	}
	if err := gsBadOp.gen(BinaryOp{Left: Number{Value: 1}, Op: '%', Right: Number{Value: 2}}, 0); err == nil {
		t.Fatal("expected unsupported operator error on spill path")
	}
	gsBadRight := &generator{}
	for range 14 {
		_, _ = gsBadRight.alloc()
	}
	if err := gsBadRight.gen(BinaryOp{Left: Number{Value: 1}, Op: '+', Right: Variable{Name: "y"}}, 0); err == nil {
		t.Fatal("expected right operand error on spill path")
	}

	// Spill exhaustion + pow edge cases
	gEx := &generator{spillDepth: maxSpillDepth}
	for range 14 {
		_, _ = gEx.alloc()
	}
	if err := gEx.gen(BinaryOp{Left: Number{Value: 1}, Op: '+', Right: Number{Value: 2}}, 0); err == nil {
		t.Fatal("expected spill exhaustion error")
	}
	if err := gEx.genPowWithSign(Number{Value: 1}, 2, true, 0); err == nil {
		t.Fatal("expected alloc error in genPowWithSign")
	}
	if err := gEx.genPowCall(Number{Value: 1}, Number{Value: 2}, 0); err == nil {
		t.Fatal("expected alloc error in genPowCall")
	}
	if err := gEx.genPowUint(Number{Value: 1}, 3, 0); err == nil {
		t.Fatal("expected alloc error in genPowUint")
	}

	// Second alloc failure in BinaryOp and genPowCall
	gOne := &generator{spillDepth: maxSpillDepth}
	for range 13 {
		_, _ = gOne.alloc()
	}
	if err := gOne.gen(BinaryOp{Left: Number{Value: 1}, Op: '+', Right: Number{Value: 2}}, 0); err == nil {
		t.Fatal("expected second alloc exhaustion in BinaryOp")
	}
	gPowOne := &generator{spillDepth: maxSpillDepth}
	for range 13 {
		_, _ = gPowOne.alloc()
	}
	if err := gPowOne.genPowCall(Number{Value: 1}, Number{Value: 2}, 0); err == nil {
		t.Fatal("expected second alloc exhaustion in genPowCall")
	}

	// Pow error & special exponent cases
	gPow := &generator{}
	if err := gPow.genPowWithSign(Variable{Name: "y"}, 2, true, 0); err == nil {
		t.Fatal("expected error for negative pow with variable y")
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

	if _, err := CompileCached("x ^ (-2.5)"); err != nil {
		t.Fatalf("CompileCached x ^ (-2.5) failed: %v", err)
	}
	if _, err := CompileCached("x ^ 17"); err != nil {
		t.Fatalf("CompileCached x ^ 17 failed: %v", err)
	}
	if _, err := CompileCached("x ^ (-17)"); err != nil {
		t.Fatalf("CompileCached x ^ (-17) failed: %v", err)
	}

	// compileToCode maxJITCodeSize exceeded
	oldMax := maxJITCodeSize
	maxJITCodeSize = 1
	_, err := compileToCode(Number{Value: 1})
	maxJITCodeSize = oldMax
	if err == nil {
		t.Fatal("expected error when generated code exceeds maxJITCodeSize")
	}

	// mmapFn error path in AllocateExecutableMemory
	oldMmap := mmapFn
	mmapFn = func(int, int64, int, int, int) ([]byte, error) {
		return nil, fmt.Errorf("injected mmap error")
	}
	_, err = AllocateExecutableMemory([]byte{0x01})
	mmapFn = oldMmap
	if err == nil {
		t.Fatal("expected mmap error")
	}

	// Parser edge cases on arm64
	if _, err := Parse("@"); err == nil {
		t.Fatal("expected error for unknown char")
	}
	if tokTypeFromRune('@') != tokEOF {
		t.Fatal("expected tokEOF for '@'")
	}
	if _, err := Parse("- *"); err == nil {
		t.Fatal("expected unary error")
	}
	if _, err := Parse("12.34.56"); err == nil {
		t.Fatal("expected invalid number error")
	}
	if _, err := Parse("( + 1)"); err == nil {
		t.Fatal("expected paren error")
	}
	if _, err := Parse("sin( + 1)"); err == nil {
		t.Fatal("expected func arg error")
	}
	if _, err := Parse("sin(1"); err == nil {
		t.Fatal("expected missing ')' error")
	}
	for _, expr := range []string{
		"x + (x - 1)", "x - (x + 1)", "(x + 1) * (x - 1)", "x / (x * 2)",
		"(x + 1)^2", "x^(x + 1)", "(-x)^2", "-(x + 1)", "x^2^3", "(x^2)^3",
	} {
		n, err := Parse(expr)
		if err != nil {
			t.Fatalf("Parse(%q) failed: %v", expr, err)
		}
		if FormatExpr(n) == "" {
			t.Fatalf("FormatExpr(%q) returned empty", expr)
		}
	}
	if FormatExpr(UnaryOp{Op: '-', Operand: BinaryOp{Left: Variable{Name: "x"}, Op: '+', Right: Number{Value: 1}}}) == "" {
		t.Fatal("FormatExpr unary binop failed")
	}
	if got := FormatExpr(Number{Value: math.Copysign(0, -1)}); got != "0" {
		t.Fatalf("FormatExpr(-0.0) = %q, want 0", got)
	}
	if got := FormatExpr(Number{Value: 1.5}); got != "1.5" {
		t.Fatalf("FormatExpr(1.5) = %q, want 1.5", got)
	}
	if got := FormatExpr(testCustomNode{}); got != "custom" {
		t.Fatalf("FormatExpr(custom) = %q, want custom", got)
	}
	if got := wrapBinOp(BinaryOp{Left: Variable{Name: "x"}, Op: '%', Right: Number{Value: 2}}, '+', true); got == "" {
		t.Fatal("wrapBinOp unknown op failed")
	}
}
