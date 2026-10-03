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
// It is an alias of mathFunc so codegen shares the single dispatch table
// declared in functab.go.
type jitFunc = mathFunc

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
	if prefix != 0 {
		e.emit(prefix)
	}
	if rex != 0x40 {
		e.emit(rex)
	}
	e.emit(0x0F, opcode, e.modrm(3, dst&7, src&7))
}

func (e *encoder) addPool(f float64) int {
	bits := math.Float64bits(f)
	for i, v := range e.pool {
		if math.Float64bits(v) == bits {
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
		insLen = 9
		dispOff = 5 // displacement shifts by 1 with REX prefix
	}
	e.fixups = append(e.fixups, poolFixup{poolIdx: idx, codeOff: off, insLen: insLen, dispOff: dispOff})
	xmmReg := dst & 7
	e.emit(0xF2)
	if hasREX {
		e.emit(0x44)
	}
	e.emit(0x0F, 0x10, e.modrm(0, xmmReg, 5))
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
	e.emit(0xF2)
	if reg >= 8 {
		e.emit(0x44)
	}
	e.emit(0x0F, 0x11, e.modrm(0, xmmReg, rsp), e.sib(0, 4, rsp))
}

func (e *encoder) movsdLoad(reg byte) {
	xmmReg := reg & 7
	e.emit(0xF2)
	if reg >= 8 {
		e.emit(0x44)
	}
	e.emit(0x0F, 0x10, e.modrm(0, xmmReg, rsp), e.sib(0, 4, rsp))
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
func (e *encoder) xorpd(dst, src byte)  { e.sse2(0x66, 0x57, dst, src) }

const (
	callFrameSize  = 264
	callSaveOffset = 128
)

func (e *encoder) saveXMM(reg byte, offset int32) {
	e.emit(0xF2)
	if reg >= 8 {
		e.emit(0x44)
	}
	e.emit(0x0F, 0x11, e.modrm(2, reg&7, rsp), e.sib(0, 4, rsp))
	e.emit32(uint32(offset))
}

func (e *encoder) restoreXMM(reg byte, offset int32) {
	e.emit(0xF2)
	if reg >= 8 {
		e.emit(0x44)
	}
	e.emit(0x0F, 0x10, e.modrm(2, reg&7, rsp), e.sib(0, 4, rsp))
	e.emit32(uint32(offset))
}

func (e *encoder) push16(reg byte) {
	e.emit(0x48, 0x83, 0xEC, 16) // sub rsp, 16
	e.movsdStore(reg)
}

func (e *encoder) pop16(reg byte) {
	e.movsdLoad(reg)
	e.emit(0x48, 0x83, 0xC4, 16) // add rsp, 16
}

type jitFunc2 func(float64, float64) float64

// callFunc2 emits code to call a two-argument Go function through an indirect call.
// Uses ABIInternal calling convention: arg0 in xmm0, arg1 in xmm1, result in xmm0.
// Allocates 128 bytes of callee register spill space above RSP to ensure that
// Go ABIInternal callees do not overwrite saved caller registers.
func (e *encoder) callFunc2(fn jitFunc2, arg0, arg1 byte) {
	funcAddr := reflect.ValueOf(fn).Pointer()
	idx := e.addUint64Pool(uint64(funcAddr)) // #nosec G115

	// On entry to JIT: RSP = 16n - 8
	// sub 264: RSP = 16n - 272 = 16(n-17), aligned to 16 bytes
	e.emit(0x48, 0x81, 0xEC)
	e.emit32(callFrameSize)

	// Save xmm0-xmm15 to [rsp+128]..[rsp+248]
	for i := byte(0); i < 16; i++ {
		e.saveXMM(i, callSaveOffset+int32(i)*8)
	}

	// Move arguments: arg0 -> xmm0, arg1 -> xmm1 from saved stack slots.
	e.restoreXMM(0, callSaveOffset+int32(arg0)*8)
	e.restoreXMM(1, callSaveOffset+int32(arg1)*8)

	// Zero xmm15 (Go ABI register invariant)
	e.xorpd(15, 15)

	// Load function pointer into R8
	off := len(e.code)
	e.emit(0x4C, 0x8B, 0x05)
	e.emit32(0)
	e.ptrFixups = append(e.ptrFixups, poolFixup{poolIdx: idx, codeOff: off, insLen: 7, dispOff: 3})

	// CALL R8 — result in xmm0
	e.emit(0x41, 0xFF, 0xD0)

	// Restore xmm1-xmm15 (skip xmm0 which holds the return value)
	for i := byte(1); i < 16; i++ {
		e.restoreXMM(i, callSaveOffset+int32(i)*8)
	}

	// add rsp, 264
	e.emit(0x48, 0x81, 0xC4)
	e.emit32(callFrameSize)
}

// callFunc emits code to call a Go function through an indirect call.
// Uses ABIInternal calling convention: argument in xmm0, result in xmm0.
// Allocates 128 bytes of callee register spill space above RSP to ensure that
// Go ABIInternal callees do not overwrite saved caller registers.
func (e *encoder) callFunc(fn jitFunc, argReg byte) {
	funcAddr := reflect.ValueOf(fn).Pointer()
	idx := e.addUint64Pool(uint64(funcAddr)) // #nosec G115

	// On entry to JIT: RSP = 16n - 8
	// sub 264: RSP = 16n - 272 = 16(n-17), aligned to 16 bytes
	e.emit(0x48, 0x81, 0xEC)
	e.emit32(callFrameSize)

	// Save xmm0-xmm15 to [rsp+128] through [rsp+248]
	for i := byte(0); i < 16; i++ {
		e.saveXMM(i, callSaveOffset+int32(i)*8)
	}

	// Move argument to xmm0 AFTER saving (from saved stack slot)
	if argReg != 0 {
		e.restoreXMM(0, callSaveOffset+int32(argReg)*8)
	}

	// Zero xmm15 (Go ABI register invariant)
	e.xorpd(15, 15)

	// Load function pointer into R8
	off := len(e.code)
	e.emit(0x4C, 0x8B, 0x05)
	e.emit32(0)
	e.ptrFixups = append(e.ptrFixups, poolFixup{poolIdx: idx, codeOff: off, insLen: 7, dispOff: 3})

	// CALL R8 — result in xmm0
	e.emit(0x41, 0xFF, 0xD0)

	// Restore xmm1-xmm15 only (xmm0 holds the return value, skip it)
	for i := byte(1); i < 16; i++ {
		e.restoreXMM(i, callSaveOffset+int32(i)*8)
	}

	// add rsp, 264
	e.emit(0x48, 0x81, 0xC4)
	e.emit32(callFrameSize)
	// Result is in xmm0.
}

const xReg byte = 14
const maxSpillDepth = 5

type generator struct {
	enc        encoder
	used       [16]bool
	spillDepth int
}

func (g *generator) alloc() (byte, error) {
	for i := byte(0); i < 14; i++ { // Skip xReg (14) and xmm15 (scratch/zero)
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
		if v.Name != "" && v.Name != "x" {
			return fmt.Errorf("JIT codegen only supports variable 'x', got %q", v.Name)
		}
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
		// Bitwise negation: flip sign bit via XORPD with 0x8000000000000000 (-0.0)
		idx := g.enc.addPool(math.Float64frombits(0x8000000000000000))
		g.enc.loadConstant(tmp, idx)
		g.enc.xorpd(dst, tmp)
	case BinaryOp:
		if v.Op == '^' {
			return g.genPow(v.Left, v.Right, dst)
		}
		if err := g.gen(v.Left, dst); err != nil {
			return err
		}
		leftSave, err := g.alloc()
		if err == nil {
			g.enc.movsdXmmXmm(leftSave, dst)
			tmp, err := g.alloc()
			if err == nil {
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
				return nil
			}
			g.free(leftSave)
		}

		// Register spill path: all scratch registers are currently occupied.
		if g.spillDepth >= maxSpillDepth {
			return fmt.Errorf("out of register resources")
		}
		g.spillDepth++
		defer func() { g.spillDepth-- }()

		// Spill left operand to stack (preserves 16-byte alignment).
		g.enc.push16(dst)
		// Evaluate right operand into dst.
		if err := g.gen(v.Right, dst); err != nil {
			return err
		}
		// Move right operand to scratch register 15.
		g.enc.movsdXmmXmm(15, dst)
		// Restore left operand from stack into dst.
		g.enc.pop16(dst)
		switch v.Op {
		case '+':
			g.enc.addsd(dst, 15)
		case '-':
			g.enc.subsd(dst, 15)
		case '*':
			g.enc.mulsd(dst, 15)
		case '/':
			g.enc.divsd(dst, 15)
		default:
			return fmt.Errorf("unsupported operator: %c", v.Op)
		}
		g.enc.xorpd(15, 15)
		return nil
	case FunctionCall:
		return g.genFuncCall(v.Name, v.Arg, dst)
	}
	return nil
}

func (g *generator) genFuncCall(name string, arg Node, dst byte) error {
	fn, ok := mathFuncs[name]
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
	// Handle unary minus on constant integer: x^(-n)
	if u, ok := exp.(UnaryOp); ok && u.Op == '-' {
		if n, ok := u.Operand.(Number); ok {
			return g.genPowWithSign(base, n.Value, true, dst)
		}
	}

	num, ok := exp.(Number)
	if !ok {
		// Variable or complex expression exponent: call math.Pow(base, exp)
		return g.genPowCall(base, exp, dst)
	}
	return g.genPowWithSign(base, num.Value, false, dst)
}

func (g *generator) genPowWithSign(base Node, value float64, negate bool, dst byte) error {
	// Non-integer exponent: call math.Pow directly to correctly handle negative bases and 0^0
	if value != float64(int(value)) {
		expVal := value
		if negate {
			expVal = -expVal
		}
		return g.genPowCall(base, Number{Value: expVal}, dst)
	}
	n := int(value)
	if negate {
		n = -n
	}

	// For large integer powers, math.Pow is faster/more accurate
	if n < -16 || n > 16 {
		return g.genPowCall(base, Number{Value: float64(n)}, dst)
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

func (g *generator) genPowCall(base, exp Node, dst byte) error {
	baseReg, err := g.alloc()
	if err != nil {
		return err
	}
	defer g.free(baseReg)
	if err = g.gen(base, baseReg); err != nil {
		return err
	}

	expReg, allocErr := g.alloc()
	if allocErr != nil {
		return allocErr
	}
	defer g.free(expReg)
	if err = g.gen(exp, expReg); err != nil {
		return err
	}

	g.enc.callFunc2(math.Pow, baseReg, expReg)
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
