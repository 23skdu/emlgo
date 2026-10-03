package eml

import (
	"math"
	"testing"
)

func TestStep32_AllocsPerRun(t *testing.T) {
	n := 1024
	a := make([]float64, n)
	b := make([]float64, n)
	c := make([]float64, n)
	res := make([]float64, n)
	for i := range a {
		a[i] = float64(i + 1)
		b[i] = float64(i + 2)
		c[i] = 1.0
	}

	allocsAdd := testing.AllocsPerRun(10, func() {
		AddSIMDTo(a, b, res)
	})
	if allocsAdd > 0 {
		t.Errorf("AddSIMDTo allocated %v allocs/op, want 0", allocsAdd)
	}

	allocsSub := testing.AllocsPerRun(10, func() {
		SubSIMDTo(a, b, res)
	})
	if allocsSub > 0 {
		t.Errorf("SubSIMDTo allocated %v allocs/op, want 0", allocsSub)
	}

	allocsMul := testing.AllocsPerRun(10, func() {
		MulSIMDTo(a, b, res)
	})
	if allocsMul > 0 {
		t.Errorf("MulSIMDTo allocated %v allocs/op, want 0", allocsMul)
	}

	allocsDiv := testing.AllocsPerRun(10, func() {
		DivSIMDTo(a, b, res)
	})
	if allocsDiv > 0 {
		t.Errorf("DivSIMDTo allocated %v allocs/op, want 0", allocsDiv)
	}

	allocsFma := testing.AllocsPerRun(10, func() {
		FmaSIMDTo(a, b, c, res)
	})
	if allocsFma > 0 {
		t.Errorf("FmaSIMDTo allocated %v allocs/op, want 0", allocsFma)
	}

	allocsSqrt := testing.AllocsPerRun(10, func() {
		SqrtSIMDTo(a, res)
	})
	if allocsSqrt > 0 {
		t.Errorf("SqrtSIMDTo allocated %v allocs/op, want 0", allocsSqrt)
	}

	allocsAbs := testing.AllocsPerRun(10, func() {
		AbsSIMDTo(a, res)
	})
	if allocsAbs > 0 {
		t.Errorf("AbsSIMDTo allocated %v allocs/op, want 0", allocsAbs)
	}

	allocsNeg := testing.AllocsPerRun(10, func() {
		NegSIMDTo(a, res)
	})
	if allocsNeg > 0 {
		t.Errorf("NegSIMDTo allocated %v allocs/op, want 0", allocsNeg)
	}

	allocsInv := testing.AllocsPerRun(10, func() {
		InvSIMDTo(a, res)
	})
	if allocsInv > 0 {
		t.Errorf("InvSIMDTo allocated %v allocs/op, want 0", allocsInv)
	}

	allocsSIMD := testing.AllocsPerRun(10, func() {
		SIMD(a, b, res)
	})
	if allocsSIMD > 0 {
		t.Errorf("SIMD allocated %v allocs/op, want 0", allocsSIMD)
	}
}

func TestStep32_LengthMismatches(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{1, 2}
	res := []float64{1, 2, 3}

	assertPanics := func(name string, fn func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s expected panic on length mismatch, got none", name)
			}
		}()
		fn()
	}

	assertPanics("AddSIMDTo a!=b", func() { AddSIMDTo(a, b, res) })
	assertPanics("AddSIMDTo a!=res", func() { AddSIMDTo(a, a, b) })
	assertPanics("SubSIMDTo a!=b", func() { SubSIMDTo(a, b, res) })
	assertPanics("SubSIMDTo a!=res", func() { SubSIMDTo(a, a, b) })
	assertPanics("MulSIMDTo a!=b", func() { MulSIMDTo(a, b, res) })
	assertPanics("MulSIMDTo a!=res", func() { MulSIMDTo(a, a, b) })
	assertPanics("DivSIMDTo a!=b", func() { DivSIMDTo(a, b, res) })
	assertPanics("DivSIMDTo a!=res", func() { DivSIMDTo(a, a, b) })
}

func TestStep32_CrossoverSweep(t *testing.T) {
	sizes := []int{256, 512, 1024, 4096, 16384, 65536}
	for _, n := range sizes {
		x := make([]float64, n)
		y := make([]float64, n)
		resSIMD := make([]float64, n)
		resAdd := make([]float64, n)
		resSqrt := make([]float64, n)

		for i := range x {
			x[i] = 1.0 + float64(i%100)*0.01
			y[i] = 2.0 + float64(i%100)*0.01
		}

		SIMD(x, y, resSIMD)
		AddSIMDTo(x, y, resAdd)
		SqrtSIMDTo(x, resSqrt)

		// Spot check correctness
		for i := 0; i < n; i += n / 10 {
			wantSIMD := nativeExp(x[i]) - nativeLog(y[i])
			if math.Abs(resSIMD[i]-wantSIMD) > 1e-12 {
				t.Fatalf("n=%d SIMD[%d] = %v, want %v", n, i, resSIMD[i], wantSIMD)
			}
			wantAdd := x[i] + y[i]
			if resAdd[i] != wantAdd {
				t.Fatalf("n=%d Add[%d] = %v, want %v", n, i, resAdd[i], wantAdd)
			}
			wantSqrt := math.Sqrt(x[i])
			if math.Abs(resSqrt[i]-wantSqrt) > 1e-12 {
				t.Fatalf("n=%d Sqrt[%d] = %v, want %v", n, i, resSqrt[i], wantSqrt)
			}
		}
	}
}

func TestStep32_LargeSliceParallelOps(t *testing.T) {
	n := SmallCutoff() * 2
	x := make([]float64, n)
	y := make([]float64, n)
	z := make([]float64, n)
	res := make([]float64, n)
	for i := range x {
		x[i] = float64(i + 1)
		y[i] = float64(i + 2)
		z[i] = 1.0
	}

	AbsSIMDTo(x, res)
	NegSIMDTo(x, res)
	InvSIMDTo(x, res)
	SubSIMDTo(x, y, res)
	MulSIMDTo(x, y, res)
	DivSIMDTo(x, y, res)
	FmaSIMDTo(x, y, z, res)
	AddScalarSIMDTo(x, 2.0, res)
	MulScalarSIMDTo(x, 2.0, res)
}

