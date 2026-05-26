//go:build wasm

package eml

import "unsafe"

// WasmAlign16 returns a slice whose underlying array is aligned to 16 bytes.
// This enables WASM JIT to emit aligned wasm_simd128 load/store instructions.
func WasmAlign16(x []float64) []float64 {
	if len(x) == 0 {
		return x
	}
	ptr := uintptr(unsafe.Pointer(&x[0]))
	offset := ptr & 15
	if offset == 0 {
		return x
	}
	// Drop up to 1 element to realign
	skip := (16 - int(offset>>3)) & 1
	if skip >= len(x) {
		return x
	}
	return x[skip:]
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
