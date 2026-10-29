//go:build arm64
// +build arm64

package jit

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

type jitFunc = mathFunc
type jitFunc2 func(float64, float64) float64

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
}

func (e *encoder) emit32(v uint32) {
	e.code = append(e.code, byte(v), byte(v>>8), byte(v>>16), byte(v>>24)) // #nosec G115
}

func (e *encoder) emit64(v uint64) {
	e.emit32(uint32(v))       // #nosec G115
	e.emit32(uint32(v >> 32)) // #nosec G115
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

// loadConstant emits LDR Dt, label (PC-relative 64-bit FP literal load).
func (e *encoder) loadConstant(dst byte, idx int) {
	off := len(e.code)
	e.fixups = append(e.fixups, poolFixup{poolIdx: idx, codeOff: off})
	e.emit32(0x5C000000 | uint32(dst&31))
}

func (e *encoder) fixupConstants() {
	if len(e.code)%8 != 0 {
		e.emit32(0xD503201F) // NOP for 8-byte alignment of literal pool
	}
	poolBase := len(e.code)
	ptrPoolBase := poolBase + len(e.pool)*8

	for _, f := range e.pool {
		e.emit64(math.Float64bits(f))
	}
	for _, p := range e.ptrPool {
		e.emit64(p)
	}

	for _, fx := range e.fixups {
		relWords := ((poolBase + fx.poolIdx*8) - fx.codeOff) >> 2
		insn := binary.LittleEndian.Uint32(e.code[fx.codeOff : fx.codeOff+4])
		insn |= (uint32(relWords) & 0x7FFFF) << 5 // #nosec G115
		binary.LittleEndian.PutUint32(e.code[fx.codeOff:fx.codeOff+4], insn)
	}

	for _, fx := range e.ptrFixups {
		relWords := ((ptrPoolBase + fx.poolIdx*8) - fx.codeOff) >> 2
		insn := binary.LittleEndian.Uint32(e.code[fx.codeOff : fx.codeOff+4])
		insn |= (uint32(relWords) & 0x7FFFF) << 5 // #nosec G115
		binary.LittleEndian.PutUint32(e.code[fx.codeOff:fx.codeOff+4], insn)
	}
}

// AArch64 scalar float64 instructions over D0..D31:
func (e *encoder) fmov(dst, src byte) {
	e.emit32(0x1E604000 | (uint32(src&31) << 5) | uint32(dst&31))
}

func (e *encoder) fadd(dst, lhs, rhs byte) {
	e.emit32(0x1E602800 | (uint32(rhs&31) << 16) | (uint32(lhs&31) << 5) | uint32(dst&31))
}

func (e *encoder) fsub(dst, lhs, rhs byte) {
	e.emit32(0x1E603800 | (uint32(rhs&31) << 16) | (uint32(lhs&31) << 5) | uint32(dst&31))
}

func (e *encoder) fmul(dst, lhs, rhs byte) {
	e.emit32(0x1E600800 | (uint32(rhs&31) << 16) | (uint32(lhs&31) << 5) | uint32(dst&31))
}

func (e *encoder) fdiv(dst, lhs, rhs byte) {
	e.emit32(0x1E601800 | (uint32(rhs&31) << 16) | (uint32(lhs&31) << 5) | uint32(dst&31))
}

func (e *encoder) fmadd(dst, mul1, mul2, addend byte) {
	e.emit32(0x1F400000 | (uint32(mul2&31) << 16) | (uint32(addend&31) << 10) | (uint32(mul1&31) << 5) | uint32(dst&31))
}

func (e *encoder) fsqrt(dst, src byte) {
	e.emit32(0x1E61C000 | (uint32(src&31) << 5) | uint32(dst&31))
}

func (e *encoder) fneg(dst, src byte) {
	e.emit32(0x1E614000 | (uint32(src&31) << 5) | uint32(dst&31))
}

func (e *encoder) fabs(dst, src byte) {
	e.emit32(0x1E60C000 | (uint32(src&31) << 5) | uint32(dst&31))
}

func (e *encoder) ret() {
	e.emit32(0xD65F03C0)
}

func (e *encoder) subSP(imm12 uint32) {
	e.emit32(0xD10003FF | ((imm12 & 0xFFF) << 10))
}

func (e *encoder) addSP(imm12 uint32) {
	e.emit32(0x910003FF | ((imm12 & 0xFFF) << 10))
}

func (e *encoder) saveD(reg byte, byteOffset uint32) {
	e.emit32(0xFD0003E0 | (((byteOffset / 8) & 0xFFF) << 10) | uint32(reg&31))
}

func (e *encoder) restoreD(reg byte, byteOffset uint32) {
	e.emit32(0xFD4003E0 | (((byteOffset / 8) & 0xFFF) << 10) | uint32(reg&31))
}

func (e *encoder) saveX(reg byte, byteOffset uint32) {
	e.emit32(0xF90003E0 | (((byteOffset / 8) & 0xFFF) << 10) | uint32(reg&31))
}

func (e *encoder) restoreX(reg byte, byteOffset uint32) {
	e.emit32(0xF94003E0 | (((byteOffset / 8) & 0xFFF) << 10) | uint32(reg&31))
}

func (e *encoder) push16(reg byte) {
	e.subSP(16)
	e.saveD(reg, 0)
}

func (e *encoder) pop16(reg byte) {
	e.restoreD(reg, 0)
	e.addSP(16)
}

const (
	callFrameSize  = 288
	callSaveOffset = 128
	fpSaveOffset   = 256
	lrSaveOffset   = 264
)

func (e *encoder) callFunc(fn jitFunc, argReg byte) {
	funcAddr := reflect.ValueOf(fn).Pointer()
	idx := e.addUint64Pool(uint64(funcAddr)) // #nosec G115

	e.subSP(callFrameSize)
	e.saveX(29, fpSaveOffset)
	e.saveX(30, lrSaveOffset)
	for i := byte(0); i < 16; i++ {
		e.saveD(i, callSaveOffset+uint32(i)*8)
	}
	if argReg != 0 {
		e.restoreD(0, callSaveOffset+uint32(argReg)*8)
	}

	off := len(e.code)
	e.ptrFixups = append(e.ptrFixups, poolFixup{poolIdx: idx, codeOff: off})
	e.emit32(0x58000008) // LDR X8, label
	e.emit32(0xD63F0100) // BLR X8

	for i := byte(1); i < 16; i++ {
		e.restoreD(i, callSaveOffset+uint32(i)*8)
	}
	e.restoreX(29, fpSaveOffset)
	e.restoreX(30, lrSaveOffset)
	e.addSP(callFrameSize)
}

func (e *encoder) callFunc2(fn jitFunc2, arg0, arg1 byte) {
	funcAddr := reflect.ValueOf(fn).Pointer()
	idx := e.addUint64Pool(uint64(funcAddr)) // #nosec G115

	e.subSP(callFrameSize)
	e.saveX(29, fpSaveOffset)
	e.saveX(30, lrSaveOffset)
	for i := byte(0); i < 16; i++ {
		e.saveD(i, callSaveOffset+uint32(i)*8)
	}
	e.restoreD(0, callSaveOffset+uint32(arg0)*8)
	e.restoreD(1, callSaveOffset+uint32(arg1)*8)

	off := len(e.code)
	e.ptrFixups = append(e.ptrFixups, poolFixup{poolIdx: idx, codeOff: off})
	e.emit32(0x58000008) // LDR X8, label
	e.emit32(0xD63F0100) // BLR X8

	for i := byte(1); i < 16; i++ {
		e.restoreD(i, callSaveOffset+uint32(i)*8)
	}
	e.restoreX(29, fpSaveOffset)
	e.restoreX(30, lrSaveOffset)
	e.addSP(callFrameSize)
}

const xReg byte = 14
const maxSpillDepth = 5

type generator struct {
	enc        encoder
	used       [16]bool
	spillDepth int
}

func (g *generator) alloc() (byte, error) {
	for i := byte(0); i < 14; i++ {
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
		g.enc.fmov(dst, xReg)
	case UnaryOp:
		if err := g.gen(v.Operand, dst); err != nil {
			return err
		}
		g.enc.fneg(dst, dst)
	case BinaryOp:
		if v.Op == '^' {
			return g.genPow(v.Left, v.Right, dst)
		}
		if err := g.gen(v.Left, dst); err != nil {
			return err
		}
		leftSave, err := g.alloc()
		if err == nil {
			g.enc.fmov(leftSave, dst)
			tmp, err := g.alloc()
			if err == nil {
				if err := g.gen(v.Right, tmp); err != nil {
					g.free(tmp)
					g.free(leftSave)
					return err
				}
				g.enc.fmov(dst, leftSave)
				g.free(leftSave)
				switch v.Op {
				case '+':
					g.enc.fadd(dst, dst, tmp)
				case '-':
					g.enc.fsub(dst, dst, tmp)
				case '*':
					g.enc.fmul(dst, dst, tmp)
				case '/':
					g.enc.fdiv(dst, dst, tmp)
				default:
					g.free(tmp)
					return fmt.Errorf("unsupported operator: %c", v.Op)
				}
				g.free(tmp)
				return nil
			}
			g.free(leftSave)
		}

		if g.spillDepth >= maxSpillDepth {
			return fmt.Errorf("out of register resources")
		}
		g.spillDepth++
		defer func() { g.spillDepth-- }()

		g.enc.push16(dst)
		if err := g.gen(v.Right, dst); err != nil {
			return err
		}
		g.enc.fmov(15, dst)
		g.enc.pop16(dst)
		switch v.Op {
		case '+':
			g.enc.fadd(dst, dst, 15)
		case '-':
			g.enc.fsub(dst, dst, 15)
		case '*':
			g.enc.fmul(dst, dst, 15)
		case '/':
			g.enc.fdiv(dst, dst, 15)
		default:
			return fmt.Errorf("unsupported operator: %c", v.Op)
		}
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
	if err := g.gen(arg, dst); err != nil {
		return err
	}
	switch name {
	case "sqrt":
		g.enc.fsqrt(dst, dst)
		return nil
	case "abs":
		g.enc.fabs(dst, dst)
		return nil
	}
	g.enc.callFunc(fn, dst)
	if dst != 0 {
		g.enc.fmov(dst, 0)
	}
	return nil
}

func (g *generator) genPow(base, exp Node, dst byte) error {
	if u, ok := exp.(UnaryOp); ok && u.Op == '-' {
		if n, ok := u.Operand.(Number); ok {
			return g.genPowWithSign(base, n.Value, true, dst)
		}
	}
	num, ok := exp.(Number)
	if !ok {
		return g.genPowCall(base, exp, dst)
	}
	return g.genPowWithSign(base, num.Value, false, dst)
}

func (g *generator) genPowWithSign(base Node, value float64, negate bool, dst byte) error {
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
	if n < -16 || n > 16 {
		return g.genPowCall(base, Number{Value: float64(n)}, dst)
	}
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
		g.enc.loadConstant(dst, idx)
		g.enc.fdiv(dst, dst, tmp)
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
		g.enc.fmov(dst, 0)
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
	if (n & (n - 1)) == 0 {
		for n > 1 {
			g.enc.fmul(dst, dst, dst)
			n >>= 1
		}
		return nil
	}

	tempReg, err := g.alloc()
	if err != nil {
		return err
	}
	defer g.free(tempReg)

	accumInitialized := false
	for n > 0 {
		if n&1 != 0 {
			if !accumInitialized {
				g.enc.fmov(tempReg, dst)
				accumInitialized = true
			} else {
				g.enc.fmul(tempReg, tempReg, dst)
			}
		}
		n >>= 1
		if n > 0 {
			g.enc.fmul(dst, dst, dst)
		}
	}
	if accumInitialized {
		g.enc.fmov(dst, tempReg)
	}
	return nil
}

var maxJITCodeSize = MaxJITCodeSize

func compileToCode(n Node) ([]byte, error) {
	if err := validateASTLimits(n); err != nil {
		return nil, err
	}
	var g generator
	g.used[0] = true
	g.used[xReg] = true
	g.enc.fmov(xReg, 0)
	if err := g.gen(n, 0); err != nil {
		return nil, err
	}
	g.enc.ret()
	g.enc.fixupConstants()
	if len(g.enc.code) > maxJITCodeSize {
		return nil, fmt.Errorf("generated code exceeds maximum size (%d bytes)", maxJITCodeSize)
	}
	return g.enc.code, nil
}

