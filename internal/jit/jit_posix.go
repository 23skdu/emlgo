//go:build !windows && !wasip1 && (!js || !wasm)

package jit

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	mmapFn     = unix.Mmap
	mprotectFn = unix.Mprotect
)

// AllocateExecutableMemory mmaps a page-aligned anonymous region large enough
// for code, copies code into it and mprotects it read/execute, returning the
// region's base address.
func AllocateExecutableMemory(code []byte) (unsafe.Pointer, error) {
	if len(code) == 0 {
		return nil, fmt.Errorf("cannot allocate 0 bytes")
	}
	pageSize := unix.Getpagesize()
	size := ((len(code) + pageSize - 1) / pageSize) * pageSize

	data, err := mmapFn(-1, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE)
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
	// munmap requires a page-aligned base address; reject unaligned
	// pointers without constructing an unsafe slice, which would trip
	// the runtime's checkptr validation under -race/-d=checkptr.
	if uintptr(ptr)%uintptr(pageSize) != 0 {
		return fmt.Errorf("munmap: invalid argument")
	}
	allocSize := ((size + pageSize - 1) / pageSize) * pageSize
	data := unsafe.Slice((*byte)(ptr), allocSize) // #nosec G103 -- audited low-level mmap/munmap interface
	return unix.Munmap(data)
}
