//go:build !windows && !wasip1 && (!js || !wasm)

package jit

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

var mprotectFn = unix.Mprotect

func AllocateExecutableMemory(code []byte) (unsafe.Pointer, error) {
	if len(code) == 0 {
		return nil, fmt.Errorf("cannot allocate 0 bytes")
	}
	pageSize := unix.Getpagesize()
	size := ((len(code) + pageSize - 1) / pageSize) * pageSize

	data, err := unix.Mmap(-1, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE)
	if err != nil {
		return nil, fmt.Errorf("mmap failed: %v", err)
	}

	copy(data, code)

	err = mprotectFn(data, unix.PROT_READ|unix.PROT_EXEC)
	if err != nil {
		_ = unix.Munmap(data)
		return nil, fmt.Errorf("mprotect failed: %v", err)
	}

	return unsafe.Pointer(&data[0]), nil // #nosec G103
}

// FreeExecutableMemory unmaps executable memory previously allocated by AllocateExecutableMemory.
func FreeExecutableMemory(ptr unsafe.Pointer, size int) error {
	if ptr == nil || size <= 0 {
		return nil
	}
	pageSize := unix.Getpagesize()
	allocSize := ((size + pageSize - 1) / pageSize) * pageSize
	data := unsafe.Slice((*byte)(ptr), allocSize)
	return unix.Munmap(data)
}
