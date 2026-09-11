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

func (c *Compiler) Compile(expr string) (Func, error) {
	ast, err := Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	code, err := compileToCode(ast)
	if err != nil {
		return nil, fmt.Errorf("codegen error: %v", err)
	}

	ptr, err := AllocateExecutableMemory(code)
	if err != nil {
		return nil, fmt.Errorf("memory allocation: %v", err)
	}

	return MakeFunc(ptr), nil
}

func MakeFunc(ptr unsafe.Pointer) Func {
	fv := &struct{ fn uintptr }{fn: uintptr(ptr)}
	return *(*Func)(unsafe.Pointer(&fv)) // #nosec G103
}
