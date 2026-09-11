//go:build js && wasm
// +build js,wasm

package jit

import (
	"fmt"
	"unsafe"
)

// Func represents a JIT-compiled function.
type Func func(float64) float64

// Compiler manages the JIT compilation process.
type Compiler struct{}

// NewCompiler creates a new JIT compiler.
func NewCompiler() *Compiler {
	return &Compiler{}
}

// Compile parses the expression into an AST and returns an evaluation closure on WebAssembly targets.
func (c *Compiler) Compile(expr string) (Func, error) {
	node, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	return func(x float64) float64 {
		return Eval(node, x)
	}, nil
}

// AllocateExecutableMemory is not supported on WebAssembly.
func AllocateExecutableMemory(code []byte) (unsafe.Pointer, error) {
	return nil, fmt.Errorf("executable memory allocation is not supported on WebAssembly")
}

// MakeFunc is not supported on WebAssembly.
func MakeFunc(ptr unsafe.Pointer) Func {
	return nil
}
