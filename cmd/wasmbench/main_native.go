//go:build !js || !wasm

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
