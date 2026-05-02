package eml

import (
	"math"
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

func FuzzExpSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		result := ExpSIMDTo_(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
		for i := range result {
			if math.IsNaN(result[i]) && math.IsNaN(math.Exp(fa[i])) {
				continue
			}
			if math.IsInf(result[i], 0) && math.IsInf(math.Exp(fa[i]), 0) {
				continue
			}
		}
	})
}

func FuzzLogSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		for i := range fa {
			if fa[i] <= 0 {
				fa[i] = 1.0
			}
		}
		result := LogSIMDTo_(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
		for i := range result {
			if math.IsNaN(fa[i]) {
				continue
			}
		}
	})
}

func FuzzSinSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		result := SinSIMDTo_(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
	})
}

func FuzzCosSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		result := CosSIMDTo_(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
	})
}

func FuzzExpMulBatch(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) || len(fa) == 0 {
			return
		}
		result := ExpMulBatch(fa, fb)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
	})
}

func FuzzLogDivBatch(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f}, []byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) || len(fa) == 0 {
			return
		}
		for i := range fb {
			if fb[i] <= 0 {
				fb[i] = 1.0
			}
		}
		result := LogDivBatch(fa, fb)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
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

func ExpSIMDTo_(x []float64) []float64 {
	if len(x) == 0 {
		return x
	}
	result := make([]float64, len(x))
	ExpSIMDTo(x, result)
	return result
}

func LogSIMDTo_(x []float64) []float64 {
	if len(x) == 0 {
		return x
	}
	result := make([]float64, len(x))
	LogSIMDTo(x, result)
	return result
}

func SinSIMDTo_(x []float64) []float64 {
	if len(x) == 0 {
		return x
	}
	result := make([]float64, len(x))
	SinSIMDTo(x, result)
	return result
}

func CosSIMDTo_(x []float64) []float64 {
	if len(x) == 0 {
		return x
	}
	result := make([]float64, len(x))
	CosSIMDTo(x, result)
	return result
}
