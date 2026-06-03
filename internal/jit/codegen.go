//go:build amd64
// +build amd64

package jit

import (
	"fmt"
	"math"
)

const (
	rsp byte = 4
)

type encoder struct {
	code   []byte
	pool   []float64
	fixups []poolFixup
}

type poolFixup struct {
	poolIdx int
	codeOff int
	insLen  int
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

func (e *encoder) loadConstant(dst byte, idx int) {
	off := len(e.code)
	insLen := 8
	hasREX := dst >= 8
	if hasREX {
		e.emit(0x44)
		insLen = 9
	}
	e.fixups = append(e.fixups, poolFixup{poolIdx: idx, codeOff: off, insLen: insLen})
	xmmReg := dst & 7
	e.emit(0xF2, 0x0F, 0x10, e.modrm(0, xmmReg, 5))
	e.emit32(0)
}

func (e *encoder) fixupConstants() {
	codeLen := len(e.code)
	for _, f := range e.pool {
		e.emit64(math.Float64bits(f))
	}
	for _, fx := range e.fixups {
		insnEnd := fx.codeOff + fx.insLen
		rel := (codeLen + fx.poolIdx*8) - insnEnd
		e.code[fx.codeOff+4] = byte(rel)           // #nosec G115
		e.code[fx.codeOff+5] = byte(rel >> 8)      // #nosec G115
		e.code[fx.codeOff+6] = byte(rel >> 16)     // #nosec G115
		e.code[fx.codeOff+7] = byte(rel >> 24)     // #nosec G115
	}
}


func (e *encoder) movsdXmmXmm(dst, src byte) {
	e.sse2(0xF2, 0x10, dst, src)
}

func (e *encoder) movsdStore(reg byte) {
	xmmReg := reg & 7
	hasREX := reg >= 8
	if hasREX {
		e.emit(0x41)
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
	e.emit(0x48, 0x83, 0xEC, 0x08)
	e.movsdStore(0)
}

func (e *encoder) popTo(reg byte) {
	e.movsdLoad(reg)
	e.emit(0x48, 0x83, 0xC4, 0x08)
}

func (e *encoder) addsd(dst, src byte)  { e.sse2(0xF2, 0x58, dst, src) }
func (e *encoder) subsd(dst, src byte)  { e.sse2(0xF2, 0x5C, dst, src) }
func (e *encoder) mulsd(dst, src byte)  { e.sse2(0xF2, 0x59, dst, src) }
func (e *encoder) divsd(dst, src byte)  { e.sse2(0xF2, 0x5E, dst, src) }
func (e *encoder) sqrtsd(dst, src byte) { e.sse2(0xF2, 0x51, dst, src) }

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
		tmp, err := g.alloc()
		if err != nil {
			return err
		}
		if err := g.gen(v.Right, tmp); err != nil {
			g.free(tmp)
			return err
		}
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
		return fmt.Errorf("function calls not supported in JIT codegen: %s", v.Name)
	}
	return nil
}

func (g *generator) genPow(base, exp Node, dst byte) error {
	num, ok := exp.(Number)
	if !ok {
		return fmt.Errorf("only constant integer exponents are supported in JIT codegen")
	}
	e := num.Value
	if e != float64(int(e)) || e < 0 {
		return fmt.Errorf("only non-negative integer exponents are supported in JIT codegen")
	}
	n := int(e)
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
