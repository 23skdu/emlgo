//go:build !js || !wasm
// +build !js !wasm

package jit

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
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
		_ = unix.Munmap(data)
		return nil, fmt.Errorf("mprotect failed: %v", err)
	}

	return unsafe.Pointer(&data[0]), nil // #nosec G103
}

func MakeFunc(ptr unsafe.Pointer) Func {
	fv := &struct{ fn uintptr }{fn: uintptr(ptr)}
	return *(*Func)(unsafe.Pointer(&fv)) // #nosec G103
}
