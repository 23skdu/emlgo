//go:build amd64
// +build amd64

package jit

import (
	"fmt"
	"math"
	"reflect"
)

const (
	rsp byte = 4
)

// jitFunc is a Go function that takes a float64 and returns a float64.
type jitFunc func(float64) float64

// funcTable maps function names to their Go implementations for JIT codegen.
var funcTable = map[string]jitFunc{
	"sin":   math.Sin,
	"cos":   math.Cos,
	"exp":   math.Exp,
	"log":   math.Log,
	"sqrt":  math.Sqrt,
	"tan":   math.Tan,
	"asin":  math.Asin,
	"acos":  math.Acos,
	"atan":  math.Atan,
	"abs":   math.Abs,
	"cbrt":  math.Cbrt,
	"log2":  math.Log2,
	"log10": math.Log10,
	"ceil":  math.Ceil,
	"floor": math.Floor,
	"trunc": math.Trunc,
	"round": math.Round,
}

type encoder struct {
	code      []byte
	pool      []float64
	fixups    []poolFixup
	ptrPool   []uint64
	ptrFixups []poolFixup
}

type poolFixup struct {
	poolIdx int
	codeOff int
	insLen  int
	dispOff int // byte offset within instruction where displacement starts
}

func (e *encoder) emit(b ...byte) {
	e.code = append(e.code, b...)
}

func (e *encoder) emit32(v uint32) {
	e.emit(byte(v), byte(v>>8), byte(v>>16), byte(v>>24)) // #nosec G115
}

func (e *encoder) emit64(v uint64) {
	e.emit32(uint32(v))       // #nosec G115
	e.emit32(uint32(v >> 32)) // #nosec G115
}

func (e *encoder) rex(w, r, x, b byte) byte {
	return 0x40 | (w << 3) | (r << 2) | (x << 1) | b
}

func (e *encoder) modrm(mod, reg, rm byte) byte {
	return (mod << 6) | ((reg & 7) << 3) | (rm & 7)
}

func (e *encoder) sib(scale, index, base byte) byte {
	return (scale << 6) | ((index & 7) << 3) | (base & 7)
}

func (e *encoder) sse2(prefix, opcode byte, dst, src byte) {
	rx := byte(0)
	if dst >= 8 {
		rx |= 1 << 2
	}
	if src >= 8 {
		rx |= 1
	}
	rex := e.rex(0, rx>>2, 0, rx&1)
	if rex != 0x40 {
		e.emit(rex)
	}
	e.emit(prefix, 0x0F, opcode, e.modrm(3, dst&7, src&7))
}

func (e *encoder) addPool(f float64) int {
	for i, v := range e.pool {
		if v == f {
			return i
		}
	}
	e.pool = append(e.pool, f)
	return len(e.pool) - 1
}

func (e *encoder) addUint64Pool(v uint64) int {
	for i, p := range e.ptrPool {
		if p == v {
			return i
		}
	}
	e.ptrPool = append(e.ptrPool, v)
	return len(e.ptrPool) - 1
}

func (e *encoder) loadConstant(dst byte, idx int) {
	off := len(e.code)
	insLen := 8
	dispOff := 4 // displacement at byte 4 for MOVSD without REX
	hasREX := dst >= 8
	if hasREX {
		e.emit(0x44)
		insLen = 9
		dispOff = 5 // displacement shifts by 1 with REX prefix
	}
	e.fixups = append(e.fixups, poolFixup{poolIdx: idx, codeOff: off, insLen: insLen, dispOff: dispOff})
	xmmReg := dst & 7
	e.emit(0xF2, 0x0F, 0x10, e.modrm(0, xmmReg, 5))
	e.emit32(0)
}

func (e *encoder) fixupConstants() {
	codeLen := len(e.code)
	float64PoolEnd := codeLen + len(e.pool)*8
	// Emit float64 pool
	for _, f := range e.pool {
		e.emit64(math.Float64bits(f))
	}
	// Emit pointer pool (function pointers)
	for _, p := range e.ptrPool {
		e.emit64(p)
	}
	// Fixup float64 pool references
	for _, fx := range e.fixups {
		insnEnd := fx.codeOff + fx.insLen
		rel := (codeLen + fx.poolIdx*8) - insnEnd
		e.code[fx.codeOff+fx.dispOff] = byte(rel)         // #nosec G115
		e.code[fx.codeOff+fx.dispOff+1] = byte(rel >> 8)  // #nosec G115
		e.code[fx.codeOff+fx.dispOff+2] = byte(rel >> 16) // #nosec G115
		e.code[fx.codeOff+fx.dispOff+3] = byte(rel >> 24) // #nosec G115
	}
	// Fixup pointer pool references
	for _, fx := range e.ptrFixups {
		insnEnd := fx.codeOff + fx.insLen
		rel := (float64PoolEnd + fx.poolIdx*8) - insnEnd
		e.code[fx.codeOff+fx.dispOff] = byte(rel)         // #nosec G115
		e.code[fx.codeOff+fx.dispOff+1] = byte(rel >> 8)  // #nosec G115
		e.code[fx.codeOff+fx.dispOff+2] = byte(rel >> 16) // #nosec G115
		e.code[fx.codeOff+fx.dispOff+3] = byte(rel >> 24) // #nosec G115
	}
}

func (e *encoder) movsdXmmXmm(dst, src byte) {
	e.sse2(0xF2, 0x10, dst, src)
}

func (e *encoder) movsdStore(reg byte) {
	xmmReg := reg & 7
	hasREX := reg >= 8
	if hasREX {
		e.emit(0x44)
	}
	e.emit(0xF2, 0x0F, 0x11, e.modrm(0, xmmReg, rsp), e.sib(0, 4, rsp))
}

func (e *encoder) movsdLoad(reg byte) {
	xmmReg := reg & 7
	hasREX := reg >= 8
	if hasREX {
		e.emit(0x44)
	}
	e.emit(0xF2, 0x0F, 0x10, e.modrm(0, xmmReg, rsp), e.sib(0, 4, rsp))
}

func (e *encoder) push() {
	e.emit(0x48, 0x83, 0xEC, 0x08) // sub rsp, 8
	e.movsdStore(0)
}

func (e *encoder) popTo(reg byte) {
	e.movsdLoad(reg)
	e.emit(0x48, 0x83, 0xC4, 0x08) // add rsp, 8
}

func (e *encoder) addsd(dst, src byte)  { e.sse2(0xF2, 0x58, dst, src) }
func (e *encoder) subsd(dst, src byte)  { e.sse2(0xF2, 0x5C, dst, src) }
func (e *encoder) mulsd(dst, src byte)  { e.sse2(0xF2, 0x59, dst, src) }
func (e *encoder) divsd(dst, src byte)  { e.sse2(0xF2, 0x5E, dst, src) }
func (e *encoder) sqrtsd(dst, src byte) { e.sse2(0xF2, 0x51, dst, src) }

// callFunc emits code to call a Go function through an indirect call.
// Uses ABIInternal calling convention: argument in xmm0, result in xmm0.
// Saves xmm0-xmm7 before the call and restores xmm1-xmm7 after (not xmm0,
// since it holds the return value). If the caller needs xmm0's pre-call value,
// it must save/restore it around the callFunc call.
func (e *encoder) callFunc(fn jitFunc, argReg byte) {
	funcAddr := reflect.ValueOf(fn).Pointer()
	idx := e.addUint64Pool(uint64(funcAddr)) // #nosec G115

	// On entry to JIT: RSP = 16n - 8 (after Go's CALL pushed ret addr)
	// sub 72: RSP = 16n - 80 = 16(n-5), aligned ✓
	e.emit(0x48, 0x83, 0xEC, 72) // sub rsp, 72

	// Save xmm0-xmm7 to [rsp+0] through [rsp+56]
	for i := byte(0); i < 8; i++ {
		xmmReg := i & 7
		e.emit(0xF2, 0x0F, 0x11, e.modrm(1, xmmReg, rsp), e.sib(0, 4, rsp), i*8)
	}

	// Move argument to xmm0 AFTER saving (so original xmm0 is preserved on stack)
	if argReg != 0 {
		e.movsdXmmXmm(0, argReg)
	}

	// Load function pointer into R8
	off := len(e.code)
	e.emit(0x4C, 0x8B, 0x05)
	e.emit32(0)
	e.ptrFixups = append(e.ptrFixups, poolFixup{poolIdx: idx, codeOff: off, insLen: 7, dispOff: 3})

	// CALL R8 — result in xmm0
	e.emit(0x41, 0xFF, 0xD0)

	// Restore xmm1-xmm7 only (xmm0 holds the return value, skip it)
	for i := byte(1); i < 8; i++ {
		xmmReg := i & 7
		e.emit(0xF2, 0x0F, 0x10, e.modrm(1, xmmReg, rsp), e.sib(0, 4, rsp), i*8)
	}

	// add rsp, 72
	e.emit(0x48, 0x83, 0xC4, 72)
	// Result is in xmm0.
}

const xReg byte = 15

type generator struct {
	enc  encoder
	used [16]bool
}

func (g *generator) alloc() (byte, error) {
	for i := byte(0); i < 15; i++ { // Skip xReg (15)
		if !g.used[i] {
			g.used[i] = true
			return i, nil
		}
	}
	return 0, fmt.Errorf("out of register resources")
}

func (g *generator) free(reg byte) {
	if reg < 16 {
		g.used[reg] = false
	}
}

func (g *generator) gen(n Node, dst byte) error {
	switch v := n.(type) {
	case Number:
		idx := g.enc.addPool(v.Value)
		g.enc.loadConstant(dst, idx)
	case Variable:
		g.enc.movsdXmmXmm(dst, xReg)
	case UnaryOp:
		if err := g.gen(v.Operand, dst); err != nil {
			return err
		}
		tmp, err := g.alloc()
		if err != nil {
			return err
		}
		defer g.free(tmp)
		idx := g.enc.addPool(-1)
		g.enc.loadConstant(tmp, idx)
		g.enc.mulsd(dst, tmp)
	case BinaryOp:
		if v.Op == '^' {
			return g.genPow(v.Left, v.Right, dst)
		}
		if err := g.gen(v.Left, dst); err != nil {
			return err
		}
		leftSave, err := g.alloc()
		if err != nil {
			return err
		}
		g.enc.movsdXmmXmm(leftSave, dst)
		tmp, err := g.alloc()
		if err != nil {
			g.free(leftSave)
			return err
		}
		if err := g.gen(v.Right, tmp); err != nil {
			g.free(tmp)
			g.free(leftSave)
			return err
		}
		g.enc.movsdXmmXmm(dst, leftSave)
		g.free(leftSave)
		switch v.Op {
		case '+':
			g.enc.addsd(dst, tmp)
		case '-':
			g.enc.subsd(dst, tmp)
		case '*':
			g.enc.mulsd(dst, tmp)
		case '/':
			g.enc.divsd(dst, tmp)
		default:
			g.free(tmp)
			return fmt.Errorf("unsupported operator: %c", v.Op)
		}
		g.free(tmp)
	case FunctionCall:
		return g.genFuncCall(v.Name, v.Arg, dst)
	}
	return nil
}

func (g *generator) genFuncCall(name string, arg Node, dst byte) error {
	fn, ok := funcTable[name]
	if !ok {
		return fmt.Errorf("unsupported function in JIT codegen: %s", name)
	}
	// Evaluate the argument into dst.
	if err := g.gen(arg, dst); err != nil {
		return err
	}
	// Call the function. dst holds the argument; result goes to xmm0.
	// callFunc saves/restores xmm1-xmm7 but xmm0 gets the return value.
	g.enc.callFunc(fn, dst)
	// Move the result from xmm0 to dst.
	if dst != 0 {
		g.enc.movsdXmmXmm(dst, 0)
	}
	return nil
}

func (g *generator) genPow(base, exp Node, dst byte) error {
	// Handle unary minus: x^(-n) where n is a constant
	if u, ok := exp.(UnaryOp); ok && u.Op == '-' {
		if n, ok := u.Operand.(Number); ok {
			return g.genPowWithSign(base, n.Value, true, dst)
		}
	}

	num, ok := exp.(Number)
	if !ok {
		// Variable or complex exponent: x^y = exp(y * log(x))
		return g.genPowExpLog(base, exp, 0, false, dst)
	}
	return g.genPowWithSign(base, num.Value, false, dst)
}

func (g *generator) genPowWithSign(base Node, value float64, negate bool, dst byte) error {
	// Non-integer exponent: x^y = exp(y * log(x))
	if value != float64(int(value)) {
		expVal := value
		if negate {
			expVal = -expVal
		}
		return g.genPowExpLog(base, nil, expVal, false, dst)
	}
	n := int(value)
	if negate {
		n = -n
	}

	// Handle negative exponents: x^n = 1.0 / x^|n|
	if n < 0 {
		tmp, err := g.alloc()
		if err != nil {
			return err
		}
		defer g.free(tmp)
		if err := g.genPowUint(base, -n, tmp); err != nil {
			return err
		}
		idx := g.enc.addPool(1)
		g.enc.loadConstant(dst, idx) // dst = 1.0
		g.enc.divsd(dst, tmp)        // dst = 1.0 / x^|n|
		return nil
	}
	return g.genPowUint(base, n, dst)
}

// genPowExpLog implements x^y via exp(y * log(x)).
// If expNode is non-nil, it is a variable/complex expression exponent evaluated into a register.
// If expNode is nil, constVal is used as a compile-time constant exponent.
func (g *generator) genPowExpLog(base Node, expNode Node, constVal float64, negate bool, dst byte) error {
	// 1. Generate base into dst
	if err := g.gen(base, dst); err != nil {
		return err
	}

	if expNode != nil {
		// Variable/complex exponent path: x^y = exp(y * log(x))
		tmp, err := g.alloc()
		if err != nil {
			return err
		}
		defer g.free(tmp)
		if err := g.gen(expNode, tmp); err != nil {
			return err
		}
		// dst = log(base) — callFunc saves/restores xmm1-xmm7 so tmp is safe
		g.enc.callFunc(math.Log, dst)
		if dst != 0 {
			g.enc.movsdXmmXmm(dst, 0)
		}
		// dst = y * log(base)
		g.enc.mulsd(dst, tmp)
	} else {
		// Constant exponent path: x^c = exp(c * log(x))
		if negate {
			constVal = -constVal
		}
		// dst = log(base)
		g.enc.callFunc(math.Log, dst)
		if dst != 0 {
			g.enc.movsdXmmXmm(dst, 0)
		}
		// Load constant into a temp and multiply
		tmp, err := g.alloc()
		if err != nil {
			return err
		}
		defer g.free(tmp)
		idx := g.enc.addPool(constVal)
		g.enc.loadConstant(tmp, idx)
		g.enc.mulsd(dst, tmp)
	}

	// dst = exp(y * log(x))
	g.enc.callFunc(math.Exp, dst)
	if dst != 0 {
		g.enc.movsdXmmXmm(dst, 0)
	}
	return nil
}

func (g *generator) genPowUint(base Node, n int, dst byte) error {
	if n == 0 {
		idx := g.enc.addPool(1)
		g.enc.loadConstant(dst, idx)
		return nil
	}
	if err := g.gen(base, dst); err != nil {
		return err
	}

	// If n is a power of two, we can avoid allocating any accumulator register.
	if (n & (n - 1)) == 0 {
		for n > 1 {
			g.enc.mulsd(dst, dst)
			n >>= 1
		}
		return nil
	}

	// Otherwise, allocate a temporary accumulator register.
	tempReg, err := g.alloc()
	if err != nil {
		return err
	}
	defer g.free(tempReg)

	accumInitialized := false

	for n > 0 {
		if n&1 != 0 {
			if !accumInitialized {
				g.enc.movsdXmmXmm(tempReg, dst)
				accumInitialized = true
			} else {
				g.enc.mulsd(tempReg, dst)
			}
		}
		n >>= 1
		if n > 0 {
			g.enc.mulsd(dst, dst)
		}
	}

	if accumInitialized {
		g.enc.movsdXmmXmm(dst, tempReg)
	}
	return nil
}

func compileToCode(n Node) ([]byte, error) {
	var g generator
	g.used[0] = true    // xmm0 is the return register and holds initial x
	g.used[xReg] = true // xmm15 (xReg) holds input variable x
	g.enc.movsdXmmXmm(xReg, 0)
	if err := g.gen(n, 0); err != nil {
		return nil, err
	}
	g.enc.emit(0xC3)
	g.enc.fixupConstants()
	return g.enc.code, nil
}
