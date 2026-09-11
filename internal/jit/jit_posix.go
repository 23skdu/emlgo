//go:build !windows && (!js || !wasm)

package jit

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

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

	err = unix.Mprotect(data, unix.PROT_READ|unix.PROT_EXEC)
	if err != nil {
		_ = unix.Munmap(data)
		return nil, fmt.Errorf("mprotect failed: %v", err)
	}

	return unsafe.Pointer(&data[0]), nil // #nosec G103
}
