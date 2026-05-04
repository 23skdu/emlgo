//go:build !js || !wasm
// +build !js !wasm

package jit

import (
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Func represents a JIT-compiled function.
type Func func(float64) float64

// Compiler manages the JIT compilation process.
type Compiler struct {
	// Add fields for parser, code generator, etc.
}

// NewCompiler creates a new JIT compiler.
func NewCompiler() *Compiler {
	return &Compiler{}
}

// Compile compiles a polynomial expression into a JIT function.
// This is a skeleton implementation.
func (c *Compiler) Compile(expr string) (Func, error) {
	// 1. Parse the expression (to be implemented)
	// 2. Generate machine code (to be implemented)
	// 3. Allocate executable memory
	// 4. Copy code and return function

	return nil, fmt.Errorf("jit compilation not yet implemented for: %s", expr)
}

// AllocateExecutableMemory allocates a page of executable memory.
func AllocateExecutableMemory(code []byte) (unsafe.Pointer, error) {
	pageSize := unix.Getpagesize()
	size := ((len(code) + pageSize - 1) / pageSize) * pageSize

	data, err := unix.Mmap(-1, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE)
	if err != nil {
		return nil, fmt.Errorf("mmap failed: %v", err)
	}

	copy(data, code)

	err = unix.Mprotect(data, unix.PROT_READ|unix.PROT_EXEC)
	if err != nil {
		_ = unix.Munmap(data) // #nosec G104
		return nil, fmt.Errorf("mprotect failed: %v", err)
	}

	return unsafe.Pointer(&data[0]), nil // #nosec G103
}

// MakeFunc converts a pointer to executable memory into a JIT Func.
func MakeFunc(ptr unsafe.Pointer) Func {
	var f Func
	// #nosec G103
	structPtr := (*reflect.SliceHeader)(unsafe.Pointer(&f))
	structPtr.Data = uintptr(ptr)
	return f
}
