package main

import (
	"testing"
)

func TestBenchParityAndAccuracy(t *testing.T) {
	iterations = 1
	verbose = false

	testAllParity()
	testAllAccuracy()
	runJitBenchmarks()
}
