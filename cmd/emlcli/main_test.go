package main

import (
	"os"
	"testing"
)

func TestEmlCliCommands(t *testing.T) {
	// Test usage
	printUsage()

	// Test demo
	runDemo()

	// Test gpu status
	runGpuStatus()

	// Test jit-test
	runJitTest()

	// Test decompile
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"eml", "decompile", "x + 1"}
	runDecompile()

	// gpu verify and gpu bench
	runGpuVerify()
	runGpuBench()
}
