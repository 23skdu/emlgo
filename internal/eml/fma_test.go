package eml

import (
	"fmt"
	"testing"
)

func TestFmaDetection(t *testing.T) {
	fmt.Printf("HasFMA: %v\n", HasFMA())
	fmt.Printf("HasAVX2: %v\n", HasAVX2())
	fmt.Printf("HasAVXVNNI: %v\n", HasAVXVNNI())
}

func TestFmaScalar(t *testing.T) {
	if !HasFMA() {
		t.Skip("FMA not supported")
	}
	a, b, c := 2.0, 3.0, 1.0
	got := FmaScalar(a, b, c)
	want := 7.0
	if got != want {
		t.Errorf("FmaScalar(%v, %v, %v) = %v, want %v", a, b, c, got, want)
	}
}

func TestFmaBatch(t *testing.T) {
	if !HasFMA() || !HasAVX2() {
		t.Skip("FMA or AVX2 not supported")
	}
	a := []float64{1, 2, 3, 4}
	b := []float64{2, 2, 2, 2}
	c := []float64{0, 1, 2, 3}
	got := FmaSIMD(a, b, c)
	want := []float64{2, 5, 8, 11}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("FmaBatch[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
