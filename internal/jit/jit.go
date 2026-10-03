//go:build !js || !wasm

package jit

import (
	"fmt"
	"unsafe"
)

type Func func(float64) float64

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

// CodeHandle manages the lifecycle of native executable code memory allocated for JIT functions.
type CodeHandle struct {
	ptr  unsafe.Pointer
	size int
	free bool
}

// Free releases the executable memory held by this handle.
func (h *CodeHandle) Free() error {
	if h == nil || h.free || h.ptr == nil {
		return nil
	}
	h.free = true
	return FreeExecutableMemory(h.ptr, h.size)
}

// CompileAST compiles an expression AST directly into executable machine code, returning a CodeHandle.
func (c *Compiler) CompileAST(ast Node) (Func, *CodeHandle, error) {
	code, err := compileToCode(ast)
	if err != nil {
		return nil, nil, fmt.Errorf("codegen error: %v", err)
	}

	ptr, err := AllocateExecutableMemory(code)
	if err != nil {
		return nil, nil, fmt.Errorf("memory allocation: %v", err)
	}

	handle := &CodeHandle{
		ptr:  ptr,
		size: len(code),
	}
	return MakeFunc(ptr), handle, nil
}

// CompileWithHandle parses and compiles an expression, returning the executable function and its memory handle.
func (c *Compiler) CompileWithHandle(expr string) (Func, *CodeHandle, error) {
	ast, err := Parse(expr)
	if err != nil {
		return nil, nil, fmt.Errorf("parse error: %v", err)
	}
	return c.CompileAST(ast)
}

func (c *Compiler) Compile(expr string) (Func, error) {
	fn, _, err := c.CompileWithHandle(expr)
	return fn, err
}

func MakeFunc(ptr unsafe.Pointer) Func {
	fv := &struct{ fn uintptr }{fn: uintptr(ptr)}
	return *(*Func)(unsafe.Pointer(&fv)) // #nosec G103
}
