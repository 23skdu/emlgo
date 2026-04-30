package eml

import (
	"testing"
)

func FuzzAddSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		_ = AddSIMD(fa, fb)
	})
}

func FuzzSubSIMD(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		_ = SubSIMD(fa, fb)
	})
}

func FuzzMulSIMD(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		_ = MulSIMD(fa, fb)
	})
}

func FuzzDivSIMD(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		_ = DivSIMD(fa, fb)
	})
}

func bytesToFloat64(b []byte) []float64 {
	n := len(b) / 8
	if n == 0 {
		return nil
	}
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = f64frombits(uint64(b[i*8]) | uint64(b[i*8+1])<<8 | uint64(b[i*8+2])<<16 | uint64(b[i*8+3])<<24 |
			uint64(b[i*8+4])<<32 | uint64(b[i*8+5])<<40 | uint64(b[i*8+6])<<48 | uint64(b[i*8+7])<<56)
	}
	return res
}
