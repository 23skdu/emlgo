//go:build amd64
// +build amd64

package eml

import "testing"

func TestAvx2Check2(t *testing.T) {
	for _, n := range []int{1, 8, 16} {
		x := make([]float64, n)
		y := make([]float64, n)
		avx2Eml(x, y, make([]float64, n))
	}
}

func TestAvx512Check2(t *testing.T) {
	for _, n := range []int{1, 8, 16} {
		x := make([]float64, n)
		y := make([]float64, n)
		avx512Eml(x, y, make([]float64, n))
	}
}
