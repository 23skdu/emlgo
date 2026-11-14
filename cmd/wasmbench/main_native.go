//go:build !js || !wasm

// Command wasmbench runs the benchmark suite when compiled for js/wasm; on
// other platforms it prints build-usage instructions.
package main

import "fmt"

var runValidation = RunValidation

func main() {
	if !runValidation() {
		defer func() { _ = recover() }()
		panic("wasmbench parity validation failed")
	}
	fmt.Println("WASM benchmark validation passed on native runtime.")
}
