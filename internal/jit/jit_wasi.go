//go:build wasip1 && wasm

package jit

import (
	"fmt"
	"unsafe"
)

// AllocateExecutableMemory is not supported on WASI, which has no mmap-based
// executable memory API. Code generation is unavailable there as well, so
// Compile reports an error before reaching this point.
func AllocateExecutableMemory(code []byte) (unsafe.Pointer, error) {
	return nil, fmt.Errorf("executable memory allocation is not supported on WASI")
}

// FreeExecutableMemory is a no-op on WASI.
func FreeExecutableMemory(ptr unsafe.Pointer, size int) error {
	return nil
}
