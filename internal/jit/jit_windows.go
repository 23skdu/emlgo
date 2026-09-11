//go:build windows

package jit

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func AllocateExecutableMemory(code []byte) (unsafe.Pointer, error) {
	if len(code) == 0 {
		return nil, fmt.Errorf("cannot allocate 0 bytes")
	}
	size := uintptr(len(code))

	// Allocate read-write memory first
	addr, err := windows.VirtualAlloc(0, size, windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if err != nil {
		return nil, fmt.Errorf("VirtualAlloc failed: %v", err)
	}

	// Copy code into allocated memory
	dst := unsafe.Slice((*byte)(unsafe.Pointer(addr)), len(code)) // #nosec G103
	copy(dst, code)

	// Change protection to executable
	var oldProtect uint32
	err = windows.VirtualProtect(addr, size, windows.PAGE_EXECUTE_READ, &oldProtect)
	if err != nil {
		_ = windows.VirtualFree(addr, 0, windows.MEM_RELEASE)
		return nil, fmt.Errorf("VirtualProtect failed: %v", err)
	}

	return unsafe.Pointer(addr), nil // #nosec G103
}
