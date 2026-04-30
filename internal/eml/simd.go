package eml

import (
	"runtime"
	"sync"

	"golang.org/x/sys/cpu"
)

var (
	hasAVX2   bool
	hasAVX512 bool
	simdOnce  sync.Once
)

func initSIMD() {
	simdOnce.Do(func() {
		if runtime.GOARCH == "amd64" {
			hasAVX2 = cpu.X86.HasAVX2
			hasAVX512 = cpu.X86.HasAVX512F
		} else if runtime.GOARCH == "arm64" {
			hasAVX2 = true
			hasAVX512 = false
		}
	})
}

func HasAVX2() bool {
	initSIMD()
	return hasAVX2
}

func HasAVX512() bool {
	initSIMD()
	return hasAVX512
}

func EmlSIMD(x, y []float64, result []float64) {
	initSIMD()
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}

	if hasAVX512 || hasAVX2 {
		emlSIMDImpl(x, y, result)
	} else {
		for i := range x {
			result[i] = Eml(x[i], y[i])
		}
	}
}

func emlSIMDImpl(x, y, result []float64) {
	for i := range x {
		result[i] = Eml(x[i], y[i])
	}
}