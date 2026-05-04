package eml

import (
	"math"
	"testing"
)

func TestBatchCoverage(t *testing.T) {
	data := []float64{0.1, 0.5, 0.9, -0.1, -0.5, -0.9}
	largeData := make([]float64, 512)
	for i := range largeData {
		largeData[i] = float64(i) / 512.0
	}

	t.Run("TanhBatch", func(t *testing.T) {
		TanhBatch(data)
		TanhBatch(largeData)
		TanhBatch([]float64{})
	})

	t.Run("AsinhBatch", func(t *testing.T) {
		AsinhBatch(data)
		AsinhBatch(largeData)
		AsinhBatch([]float64{})
	})

	t.Run("AcoshBatch", func(t *testing.T) {
		AcoshBatch([]float64{1.1, 2.0, 5.0})
		l2 := make([]float64, 512)
		for i := range l2 {
			l2[i] = 1.0 + float64(i)
		}
		AcoshBatch(l2)
		AcoshBatch([]float64{})
	})

	t.Run("AtanhBatch", func(t *testing.T) {
		AtanhBatch([]float64{0.1, 0.5, -0.5})
		AtanhBatch(largeData)
		AtanhBatch([]float64{})
	})

	t.Run("AbsSIMDTo", func(t *testing.T) {
		res := make([]float64, len(data))
		AbsSIMDTo(data, res)
		AbsSIMDTo([]float64{}, []float64{})
		
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("AbsSIMDTo did not panic on length mismatch")
			}
		}()
		AbsSIMDTo(data, make([]float64, 1))
	})

	t.Run("NegSIMDTo", func(t *testing.T) {
		res := make([]float64, len(data))
		NegSIMDTo(data, res)
		NegSIMDTo([]float64{}, []float64{})
		
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NegSIMDTo did not panic on length mismatch")
			}
		}()
		NegSIMDTo(data, make([]float64, 1))
	})

	t.Run("InvSIMDTo", func(t *testing.T) {
		res := make([]float64, len(data))
		InvSIMDTo(data, res)
		InvSIMDTo([]float64{}, []float64{})
		
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("InvSIMDTo did not panic on length mismatch")
			}
		}()
		InvSIMDTo(data, make([]float64, 1))
	})

	t.Run("SinCosSIMDTo", func(t *testing.T) {
		s := make([]float64, len(data))
		c := make([]float64, len(data))
		SinCosSIMDTo(data, s, c)
	})
	
	t.Run("OtherSIMDTo", func(t *testing.T) {
		res := make([]float64, len(data))
		ExpSIMDTo(data, res)
		LogSIMDTo(data, res)
		SinSIMDTo(data, res)
		CosSIMDTo(data, res)
		TanSIMDTo(data, res)
		SqrtSIMDTo(data, res)
	})
	t.Run("EMLError", func(t *testing.T) {
		err := EMLError("test error")
		if err.Error() != "test error" {
			t.Errorf("EMLError.Error() = %v, want 'test error'", err.Error())
		}
	})
}

func TestSIMDEdgeCases(t *testing.T) {
	// SqrtSIMD small n
	SqrtSIMD([]float64{4.0, 9.0})
	
	// SqrtSIMD negative
	res := SqrtSIMD([]float64{-1.0})
	if !math.IsNaN(res[0]) {
		t.Errorf("SqrtSIMD(-1) = %v, want NaN", res[0])
	}
	
	// AddSIMD/MulSIMD etc with small and large n
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	AddSIMD(a, b)
	SubSIMD(a, b)
	MulSIMD(a, b)
	DivSIMD(a, b)
	
	la := make([]float64, 512)
	lb := make([]float64, 512)
	AddSIMD(la, lb)
	SubSIMD(la, lb)
	MulSIMD(la, lb)
	DivSIMD(la, lb)
	
	AddScalarSIMD(a, 10.0)
	MulScalarSIMD(a, 10.0)
	AddScalarSIMD(la, 10.0)
	MulScalarSIMD(la, 10.0)
}
