//go:build wasm

package eml

import "unsafe"

// WasmAlign16 returns a slice whose underlying array is aligned to 16 bytes.
// This enables WASM JIT to emit aligned wasm_simd128 load/store instructions.
// Go guarantees float64 slices are 8-byte aligned; this skips at most one
// element to reach 16-byte alignment.
func WasmAlign16(x []float64) []float64 {
	if len(x) == 0 {
		return x
	}
	if uintptr(unsafe.Pointer(&x[0]))&15 == 0 {
		return x
	}
	if len(x) > 1 {
		return x[1:]
	}
	return x
}

// WasmPageSize is the WASM page size (64 KiB).
const WasmPageSize = 65536

// WasmPageAlign rounds n up to the nearest WASM page boundary.
func WasmPageAlign(n int) int {
	if n <= 0 {
		return WasmPageSize
	}
	return (n + WasmPageSize - 1) / WasmPageSize * WasmPageSize
}
