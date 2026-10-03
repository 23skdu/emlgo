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

// CodeHandle manages executable code memory on WebAssembly (no-op).
type CodeHandle struct{}

// Free is a no-op on WebAssembly.
func (h *CodeHandle) Free() error {
	return nil
}

// CompileAST compiles an AST into an evaluation closure on WebAssembly.
func (c *Compiler) CompileAST(node Node) (Func, *CodeHandle, error) {
	return func(x float64) float64 {
		return Eval(node, x)
	}, &CodeHandle{}, nil
}

// CompileWithHandle parses and compiles an expression on WebAssembly.
func (c *Compiler) CompileWithHandle(expr string) (Func, *CodeHandle, error) {
	node, err := Parse(expr)
	if err != nil {
		return nil, nil, err
	}
	return c.CompileAST(node)
}

// Compile parses the expression into an AST and returns an evaluation closure on WebAssembly targets.
func (c *Compiler) Compile(expr string) (Func, error) {
	fn, _, err := c.CompileWithHandle(expr)
	return fn, err
}

// AllocateExecutableMemory is not supported on WebAssembly.
func AllocateExecutableMemory(code []byte) (unsafe.Pointer, error) {
	return nil, fmt.Errorf("executable memory allocation is not supported on WebAssembly")
}

// MakeFunc is not supported on WebAssembly.
func MakeFunc(ptr unsafe.Pointer) Func {
	return nil
}

// FreeExecutableMemory is a no-op on WebAssembly.
func FreeExecutableMemory(ptr unsafe.Pointer, size int) error {
	return nil
}
