//go:build !amd64 && (!js || !wasm)

package jit

import "fmt"

func compileToCode(_ Node) ([]byte, error) {
	return nil, fmt.Errorf("JIT codegen requires amd64")
}
