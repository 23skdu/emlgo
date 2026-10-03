//go:build !windows && amd64 && (!js || !wasm)

package jit

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"unsafe"
)

type testCustomNode struct{}

func (testCustomNode) String() string { return "custom" }
func (testCustomNode) nodeSigil()     {}

func TestArenaEdgeCases(t *testing.T) {
	a := NewArena(0)
	if a == nil || len(a.buf) != 64 {
		t.Fatalf("NewArena(0) failed")
	}

	if a.ToInterface(nil) != nil {
		t.Fatal("ToInterface(nil) should be nil")
	}
	if a.ToInterface(&ArenaNode{Kind: 99}) != nil {
		t.Fatal("ToInterface(unknown) should be nil")
	}

	if a.FromInterface(nil) != nil {
		t.Fatal("FromInterface(nil) should be nil")
	}
	if a.FromInterface(testCustomNode{}) != nil {
		t.Fatal("FromInterface(custom) should be nil")
	}

	if EvalArena(nil, 0) != 0 {
		t.Fatal("EvalArena(nil) should be 0")
	}
	if EvalArena(&ArenaNode{Kind: 99}, 0) != 0 {
		t.Fatal("EvalArena(unknown) should be 0")
	}
}

func TestCacheEdgeCases(t *testing.T) {
	c := newJITCache(10)
	shard := c.getShard("bad")
	shard.mu.Lock()
	shard.items["bad"] = nil
	shard.mu.Unlock()

	if fn, ok := c.get("bad"); ok || fn != nil {
		t.Fatal("expected get to return nil, false for corrupted entry")
	}
}

func TestCanonicalEdgeCases(t *testing.T) {
	if got := (&EMLNode{Kind: EMLVar, Name: "y"}).String(); got != "y" {
		t.Fatalf("String(var y) = %q, want y", got)
	}
	if got := (&EMLNode{Kind: 99}).String(); got != "?" {
		t.Fatalf("String(unknown) = %q, want ?", got)
	}

	if EMLEval(nil, 0) != 0 {
		t.Fatal("EMLEval(nil) should be 0")
	}
	if EMLEval(&EMLNode{Kind: 99}, 0) != 0 {
		t.Fatal("EMLEval(unknown) should be 0")
	}

	if EMLEvalRegularized(nil, 0, 1e-6) != 0 {
		t.Fatal("EMLEvalRegularized(nil) should be 0")
	}
	if got := EMLEvalRegularized(&EMLNode{Kind: EMLFunc, Name: "exp", Left: &EMLNode{Kind: EMLConst, Value: 750}}, 0, 1e-6); got != math.Exp(700) {
		t.Fatalf("EMLEvalRegularized exp(750) = %f, want exp(700)", got)
	}
	// non-NaN math func
	if got := EMLEvalRegularized(&EMLNode{Kind: EMLFunc, Name: "sinh", Left: &EMLNode{Kind: EMLConst, Value: 1.0}}, 0, 1e-6); got != math.Sinh(1.0) {
		t.Fatalf("EMLEvalRegularized sinh(1.0) = %f, want %f", got, math.Sinh(1.0))
	}
	// acosh(0.5) is NaN, but arg is not NaN -> returns arg
	if got := EMLEvalRegularized(&EMLNode{Kind: EMLFunc, Name: "acosh", Left: &EMLNode{Kind: EMLConst, Value: 0.5}}, 0, 1e-6); got != 0.5 {
		t.Fatalf("EMLEvalRegularized acosh(0.5) = %f, want 0.5", got)
	}
	// acosh(NaN) -> returns 0
	if got := EMLEvalRegularized(&EMLNode{Kind: EMLFunc, Name: "acosh", Left: &EMLNode{Kind: EMLConst, Value: math.NaN()}}, 0, 1e-6); got != 0 {
		t.Fatalf("EMLEvalRegularized acosh(NaN) = %f, want 0", got)
	}

	if Canonicalize(nil) != nil {
		t.Fatal("Canonicalize(nil) should be nil")
	}
	if got := Canonicalize(Variable{Name: ""}); got == nil || got.Name != "x" {
		t.Fatalf("Canonicalize empty var name failed: %v", got)
	}
	if got := Canonicalize(UnaryOp{Op: '+', Operand: Number{Value: 3}}); got == nil || got.Value != 3 {
		t.Fatalf("Canonicalize unary + failed: %v", got)
	}
	if Canonicalize(testCustomNode{}) != nil {
		t.Fatal("Canonicalize(custom) should be nil")
	}

	if Simplify(&EMLNode{Kind: 99}) == nil {
		t.Fatal("Simplify unknown failed")
	}

	diffCeil := Diff(&EMLNode{Kind: EMLFunc, Name: "ceil", Left: &EMLNode{Kind: EMLVar}})
	if diffCeil == nil || diffCeil.Kind != EMLConst || diffCeil.Value != 0 {
		t.Fatalf("Diff ceil failed: %v", diffCeil)
	}

	sinVar := &EMLNode{Kind: EMLFunc, Name: "sin", Left: &EMLNode{Kind: EMLVar}}
	if !Equiv(sinVar, sinVar) {
		t.Fatal("Equiv(sin, sin) should be true")
	}
	if Equiv(&EMLNode{Kind: 99}, &EMLNode{Kind: 99}) {
		t.Fatal("Equiv(unknown, unknown) should be false")
	}
}

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

func TestDecompileEdgeCases(t *testing.T) {
	if emlPrec(nil) != 5 {
		t.Fatalf("emlPrec(nil) = %d, want 5", emlPrec(nil))
	}
	if emlPrec(&EMLNode{Kind: EMLFunc, Name: "pow"}) != 3 {
		t.Fatalf("emlPrec(pow) = %d, want 3", emlPrec(&EMLNode{Kind: EMLFunc, Name: "pow"}))
	}
	if wrapDecomp(nil, "+", true) != "" {
		t.Fatal("wrapDecomp(nil) should be empty")
	}
	powNode := &EMLNode{Kind: EMLFunc, Name: "pow", Left: &EMLNode{Kind: EMLVar}, Right: &EMLNode{Kind: EMLConst, Value: 2}}
	if got := wrapDecomp(powNode, "pow", true); got != "(x^2.0)" {
		t.Fatalf("wrapDecomp pow in ^ left = %q, want (x^2.0)", got)
	}

	if Decompile(&EMLNode{Kind: 99}) != "" {
		t.Fatal("Decompile(unknown) should be empty")
	}
	if DecompileLaTeX(&EMLNode{Kind: 99}) != "" {
		t.Fatal("DecompileLaTeX(unknown) should be empty")
	}
	if DecompileNodeToExpr(&EMLNode{Kind: 99}) != "" {
		t.Fatal("DecompileNodeToExpr(unknown) should be empty")
	}
	if emlNodeToJITNode(nil) != nil {
		t.Fatal("emlNodeToJITNode(nil) should be nil")
	}
	if emlNodeToJITNode(&EMLNode{Kind: 99}) != nil {
		t.Fatal("emlNodeToJITNode(unknown) should be nil")
	}
}

func TestJITPosixEdgeCases(t *testing.T) {
	if _, err := AllocateExecutableMemory(nil); err == nil {
		t.Fatal("expected error for nil code")
	}
	if _, err := AllocateExecutableMemory([]byte{}); err == nil {
		t.Fatal("expected error for empty code")
	}

	oldMprotect := mprotectFn
	mprotectFn = func([]byte, int) error {
		return fmt.Errorf("injected mprotect error")
	}
	_, err := AllocateExecutableMemory([]byte{0xC3}) // RET
	mprotectFn = oldMprotect
	if err == nil {
		t.Fatal("expected error for failed mprotect")
	}
}

func TestParserEdgeCases(t *testing.T) {
	if err := validateVars(Variable{Name: ""}, map[string]bool{"x": true}); err != nil {
		t.Fatalf("validateVars empty name failed: %v", err)
	}
	if err := validateVars(BinaryOp{Left: Variable{Name: "bad"}}, map[string]bool{"x": true}); err == nil {
		t.Fatal("expected error for bad var in Left")
	}
	if err := validateVars(UnaryOp{Operand: Variable{Name: "bad"}}, map[string]bool{"x": true}); err == nil {
		t.Fatal("expected error for bad var in Operand")
	}

	if _, err := ParseWithVars("x +", []string{"x"}); err == nil {
		t.Fatal("expected ParseWithVars syntax error")
	}

	if got := (Variable{Name: ""}).String(); got != "x" {
		t.Fatalf("Variable{}.String() = %q, want x", got)
	}

	negZero := math.Copysign(0.0, -1.0)
	if got := FormatExpr(Number{Value: negZero}); got != "0" {
		t.Fatalf("FormatExpr(-0.0) = %q, want 0", got)
	}
}

func TestJITRegisterABIStep11(t *testing.T) {
	expr := "1+(2+(3+(4+(sin(x)))))"
	c := NewCompiler()
	fn, err := c.Compile(expr)
	if err != nil {
		t.Fatal(err)
	}
	ast, err := Parse(expr)
	if err != nil {
		t.Fatal(err)
	}
	x := 0.5
	got := fn(x)
	want := Eval(ast, x)
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("MISMATCH: got=%v, want=%v", got, want)
	}

	// Cover +, -, *, /, and invalid operators under spill
	for _, op := range []rune{'+', '-', '*', '/', '%'} {
		gSpill := &generator{}
		for i := 0; i < 14; i++ {
			_, _ = gSpill.alloc()
		}
		_ = gSpill.gen(BinaryOp{Left: Number{Value: 10}, Op: op, Right: Number{Value: 2}}, 0)
	}
}

func TestStep18CacheAndMemoryReclamation(t *testing.T) {
	// 1. FreeExecutableMemory nil / non-positive size
	if err := FreeExecutableMemory(nil, 100); err != nil {
		t.Errorf("FreeExecutableMemory(nil) = %v, want nil", err)
	}
	dummy := byte(0)
	if err := FreeExecutableMemory(unsafe.Pointer(&dummy), 0); err != nil {
		t.Errorf("FreeExecutableMemory(size=0) = %v, want nil", err)
	}

	// 2. CodeHandle.Free edge cases
	var nilHandle *CodeHandle
	if err := nilHandle.Free(); err != nil {
		t.Errorf("nilHandle.Free() = %v, want nil", err)
	}
	emptyHandle := &CodeHandle{}
	if err := emptyHandle.Free(); err != nil {
		t.Errorf("emptyHandle.Free() = %v, want nil", err)
	}

	// 3. CompileWithHandle and double Free
	c := NewCompiler()
	fn, handle, err := c.CompileWithHandle("3 * x + 7")
	if err != nil {
		t.Fatalf("CompileWithHandle: %v", err)
	}
	if fn(2) != 13 {
		t.Fatalf("fn(2) = %v, want 13", fn(2))
	}
	if handle.size <= 0 {
		t.Fatalf("handle size = %d, want > 0", handle.size)
	}
	if err := handle.Free(); err != nil {
		t.Fatalf("handle.Free() = %v", err)
	}
	if err := handle.Free(); err != nil {
		t.Fatalf("second handle.Free() = %v, want nil", err)
	}

	// 4. newCacheShard(0) and newJITCache(0)
	shard := newCacheShard(0)
	if shard.maxEntries != 1 {
		t.Errorf("shard.maxEntries = %d, want 1", shard.maxEntries)
	}
	cache0 := newJITCache(0)
	if cache0.maxSize != maxCacheEntries {
		t.Errorf("cache0.maxSize = %d, want %d", cache0.maxSize, maxCacheEntries)
	}

	// 5. Cache eviction frees handle
	smallCache := newJITCache(2)
	_, h1, _ := c.CompileWithHandle("x + 1")
	_, h2, _ := c.CompileWithHandle("x + 2")
	_, h3, _ := c.CompileWithHandle("x + 3")
	smallCache.put("k1", func(x float64) float64 { return x + 1 }, h1)
	smallCache.put("k2", func(x float64) float64 { return x + 2 }, h2)
	smallCache.put("k1", func(x float64) float64 { return x + 1 }) // recency update
	smallCache.put("k3", func(x float64) float64 { return x + 3 }, h3) // evicts k2
	if !h2.free {
		t.Errorf("expected h2 to be freed on eviction")
	}
	smallCache.clear()
	if !h1.free || !h3.free {
		t.Errorf("expected h1 and h3 to be freed on clear")
	}

	// 6. Canonicalization, singleflight concurrency, and error handling in CompileCached
	ClearJITCache()
	fn1, err := CompileCached("2 * x + 1")
	if err != nil {
		t.Fatalf("CompileCached: %v", err)
	}
	fn2, err := CompileCached(" ( 2 * x + 1 ) ")
	if err != nil {
		t.Fatalf("CompileCached with parens: %v", err)
	}
	if fn1(10) != fn2(10) || fn1(10) != 21 {
		t.Errorf("canonicalization mismatch: fn1=%v fn2=%v", fn1(10), fn2(10))
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f, err := CompileCached("x^2 + 5*x + 6")
			if err != nil {
				t.Errorf("concurrent CompileCached failed: %v", err)
			}
			if f(3) != 30 {
				t.Errorf("concurrent f(3) = %v, want 30", f(3))
			}
		}()
	}
	wg.Wait()

	// Parse error in CompileCached
	if _, err := CompileCached("2 * +"); err == nil {
		t.Errorf("CompileCached(\"2 * +\") should fail")
	}

	// Singleflight error propagation
	var sg singleflightGroup
	_, errSF := sg.do("errKey", func() (Func, error) {
		return nil, fmt.Errorf("injected failure")
	})
	if errSF == nil || errSF.Error() != "injected failure" {
		t.Errorf("singleflight do did not propagate error: %v", errSF)
	}

	// Singleflight wait branch coverage
	var sf singleflightGroup
	blockCh := make(chan struct{})
	enteredCh := make(chan struct{})
	go func() {
		_, _ = sf.do("sameKey", func() (Func, error) {
			close(enteredCh)
			<-blockCh
			return func(x float64) float64 { return x }, nil
		})
	}()
	<-enteredCh
	go func() {
		close(blockCh)
	}()
	fnWait, errWait := sf.do("sameKey", func() (Func, error) {
		return nil, nil
	})
	if errWait != nil || fnWait == nil {
		t.Errorf("expected singleflight wait branch to succeed")
	}

	// Cache recheck inside CompileCached sf.do
	ClearJITCache()
	defaultCache.put("cachedCanon", func(x float64) float64 { return x })
	fnRecheck, err := defaultCache.sf.do("cachedCanon", func() (Func, error) {
		if fn, ok := defaultCache.get("cachedCanon"); ok {
			return fn, nil
		}
		return nil, nil
	})
	if err != nil || fnRecheck == nil {
		t.Errorf("expected recheck in sf.do to hit")
	}

	// Update existing cache entry with a handle
	testShardCache := newJITCache(5)
	testShardCache.put("kUpdate", func(x float64) float64 { return x })
	_, updateH, _ := c.CompileWithHandle("x + 99")
	testShardCache.put("kUpdate", func(x float64) float64 { return x + 99 }, updateH)

	// Test cache recheck inside CompileCached sf.do via compilePreflightHook
	ClearJITCache()
	compilePreflightHook = func() {
		defaultCache.put("x + 5", func(x float64) float64 { return x * 3 })
	}
	fnHook, errHook := CompileCached("x + 5")
	compilePreflightHook = nil
	if errHook != nil || fnHook(4) != 12 {
		t.Errorf("hook_canon failed: %v", errHook)
	}

	// CompileAST invalid node error
	if _, _, err := c.CompileAST(Variable{Name: "y"}); err == nil {
		t.Errorf("CompileAST(Variable{Name: \"y\"}) should fail")
	}

	// CompileCached memory allocation error via mprotect injection
	origMprotect := mprotectFn
	mprotectFn = func(b []byte, prot int) error { return fmt.Errorf("injected mprotect error") }
	_, errMprotect := CompileCached("x + 9999")
	mprotectFn = origMprotect
	if errMprotect == nil {
		t.Errorf("CompileCached with failing mprotect should error")
	}
}
