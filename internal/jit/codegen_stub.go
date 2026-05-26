//go:build !amd64
// +build !amd64

package jit

import "fmt"

func compileToCode(n Node) ([]byte, error) {
	return nil, fmt.Errorf("JIT codegen requires amd64")
}
