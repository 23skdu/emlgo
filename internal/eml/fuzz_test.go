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
		got := AddSIMD(fa, fb)
		for i := range got {
			assertFuzzExact(t, "add", got[i], fa[i]+fb[i])
		}
	})
}

func FuzzSubSIMD(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		got := SubSIMD(fa, fb)
		for i := range got {
			assertFuzzExact(t, "sub", got[i], fa[i]-fb[i])
		}
	})
}

func FuzzMulSIMD(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		got := MulSIMD(fa, fb)
		for i := range got {
			assertFuzzExact(t, "mul", got[i], fa[i]*fb[i])
		}
	})
}

func FuzzDivSIMD(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		got := DivSIMD(fa, fb)
		for i := range got {
			assertFuzzExact(t, "div", got[i], fa[i]/fb[i])
		}
	})
}

func FuzzSqrtSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		got := SqrtSIMD(fa)
		if len(got) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(fa))
		}
		for i := range got {
			assertFuzzExact(t, "sqrt", got[i], math.Sqrt(fa[i]))
		}
	})
}

func FuzzExpSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		result := ExpSIMD(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
		for i := range result {
			assertFuzzClose(t, "exp", result[i], math.Exp(fa[i]))
		}
	})
}

func FuzzLogSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		result := LogSIMD(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
		for i := range result {
			assertFuzzClose(t, "log", result[i], ExactLogOracle(fa[i]))
		}
	})
}

func FuzzSinSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		result := SinSIMD(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
		for i := range result {
			assertFuzzClose(t, "sin", result[i], math.Sin(fa[i]))
		}
	})
}

func FuzzCosSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		result := CosSIMD(fa)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
		for i := range result {
			assertFuzzClose(t, "cos", result[i], math.Cos(fa[i]))
		}
	})
}

func FuzzTanSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		got := TanSIMD(fa)
		if len(got) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(fa))
		}
		for i := range got {
			assertFuzzClose(t, "tan", got[i], math.Tan(fa[i]))
		}
	})
}

func FuzzFmaSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f},
		[]byte{0, 0, 0, 0, 0, 0, 0, 0x40},
		[]byte{0, 0, 0, 0, 0, 0, 0xf0, 0xbf})
	f.Fuzz(func(t *testing.T, a, b, c []byte) {
		fa, fb, fc := bytesToFloat64(a), bytesToFloat64(b), bytesToFloat64(c)
		if len(fa) != len(fb) || len(fa) != len(fc) {
			return
		}
		got := FmaSIMD(fa, fb, fc)
		if len(got) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(fa))
		}
		for i := range got {
			assertFuzzClose(t, "fma", got[i], math.FMA(fa[i], fb[i], fc[i]))
		}
	})
}

func FuzzFloat32SIMD(f *testing.F) {
	f.Add(
		[]byte{0, 0, 0, 0, 0, 0, 128, 63, 0, 0, 0, 128, 0, 0, 128, 127},
		[]byte{0, 0, 0, 64, 0, 0, 128, 191, 0, 0, 0, 0},
	)
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa, fb := bytesToFloat32(a), bytesToFloat32(b)
		unary := []struct {
			name string
			fn   func([]float32) []float32
			ref  func(float32) float32
		}{
			{"exp32", ExpSIMDF32, func(x float32) float32 { return float32(math.Exp(float64(x))) }},
			{"log32", LogSIMDF32, func(x float32) float32 { return float32(math.Log(float64(x))) }},
			{"sqrt32", SqrtSIMDF32, func(x float32) float32 { return float32(math.Sqrt(float64(x))) }},
			{"abs32", AbsSIMDF32, func(x float32) float32 { return float32(math.Abs(float64(x))) }},
			{"neg32", NegSIMDF32, func(x float32) float32 { return -x }},
			{"inv32", InvSIMDF32, func(x float32) float32 { return 1 / x }},
			{"add-scalar32", func(x []float32) []float32 { return AddScalarSIMDF32(x, 1.25) },
				func(x float32) float32 { return x + 1.25 }},
			{"mul-scalar32", func(x []float32) []float32 { return MulScalarSIMDF32(x, -2.5) },
				func(x float32) float32 { return x * -2.5 }},
			{"sin32", SinSIMDF32, func(x float32) float32 { return float32(math.Sin(float64(x))) }},
			{"cos32", CosSIMDF32, func(x float32) float32 { return float32(math.Cos(float64(x))) }},
			{"tan32", TanSIMDF32, func(x float32) float32 { return float32(math.Tan(float64(x))) }},
		}
		for _, op := range unary {
			got := op.fn(fa)
			if len(got) != len(fa) {
				t.Fatalf("%s length mismatch: got %d, want %d", op.name, len(got), len(fa))
			}
			for i := range got {
				assertFuzzExact32(t, op.name, got[i], op.ref(fa[i]))
			}
		}
		if len(fa) != len(fb) {
			return
		}
		binary := []struct {
			name string
			fn   func([]float32, []float32) []float32
			to   func([]float32, []float32, []float32)
			ref  func(float32, float32) float32
		}{
			{"add32", AddSIMDF32, AddSIMDF32To, func(x, y float32) float32 { return x + y }},
			{"sub32", SubSIMDF32, SubSIMDF32To, func(x, y float32) float32 { return x - y }},
			{"mul32", MulSIMDF32, MulSIMDF32To, func(x, y float32) float32 { return x * y }},
			{"div32", DivSIMDF32, DivSIMDF32To, func(x, y float32) float32 { return x / y }},
		}
		for _, op := range binary {
			got := op.fn(fa, fb)
			if len(got) != len(fa) {
				t.Fatalf("%s length mismatch: got %d, want %d", op.name, len(got), len(fa))
			}
			for i := range got {
				assertFuzzExact32(t, op.name, got[i], op.ref(fa[i], fb[i]))
			}
			inPlace := append([]float32(nil), fa...)
			op.to(inPlace, fb, inPlace)
			for i := range inPlace {
				assertFuzzExact32(t, op.name+" alias", inPlace[i], op.ref(fa[i], fb[i]))
			}
		}
	})
}

func FuzzUnaryBatchResults(f *testing.F) {
	f.Add([]byte{
		0, 0, 0, 0, 0, 0, 240, 191,
		0, 0, 0, 0, 0, 0, 240, 63,
		0, 0, 0, 0, 0, 0, 0, 0,
	})
	f.Fuzz(func(t *testing.T, data []byte) {
		x := bytesToFloat64(data)
		ops := []struct {
			name string
			fn   func([]float64) []float64
			ref  func(float64) float64
		}{
			{"abs", AbsSIMD, math.Abs},
			{"neg", NegSIMD, func(x float64) float64 { return -x }},
			{"inv", InvSIMD, func(x float64) float64 { return 1 / x }},
			{"sinh", SinhBatch, math.Sinh},
			{"cosh", CoshBatch, math.Cosh},
			{"tanh", TanhBatch, math.Tanh},
			{"asinh", AsinhBatch, math.Asinh},
			{"acosh", AcoshBatch, math.Acosh},
			{"atanh", AtanhBatch, math.Atanh},
		}
		for _, op := range ops {
			got := op.fn(x)
			if len(got) != len(x) {
				t.Fatalf("%s length mismatch: got %d, want %d", op.name, len(got), len(x))
			}
			for i := range got {
				assertFuzzClose(t, op.name, got[i], op.ref(x[i]))
			}
		}

		sin, cos := SinCosSIMD(x)
		if len(sin) != len(x) || len(cos) != len(x) {
			t.Fatalf("sincos length mismatch: got (%d,%d), want %d", len(sin), len(cos), len(x))
		}
		for i := range x {
			wantSin, wantCos := math.Sincos(x[i])
			assertFuzzClose(t, "sincos-sin", sin[i], wantSin)
			assertFuzzClose(t, "sincos-cos", cos[i], wantCos)
		}

		log2, log10 := make([]float64, len(x)), make([]float64, len(x))
		Log2SIMDTo(x, log2)
		Log10SIMDTo(x, log10)
		for i := range x {
			assertFuzzClose(t, "log2", log2[i], nativeLog2(x[i]))
			assertFuzzClose(t, "log10", log10[i], nativeLog10(x[i]))
		}

		addScalar, mulScalar := make([]float64, len(x)), make([]float64, len(x))
		AddScalarSIMDTo(x, 1.25, addScalar)
		MulScalarSIMDTo(x, -2.5, mulScalar)
		for i := range x {
			assertFuzzExact(t, "add-scalar", addScalar[i], x[i]+1.25)
			assertFuzzExact(t, "mul-scalar", mulScalar[i], x[i]*-2.5)
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
		for i := range result {
			assertFuzzClose(t, "exp-mul", result[i], math.Exp(fa[i])*fb[i])
		}
	})
}

func FuzzLogDivBatch(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f}, []byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0}, []byte{0, 0, 0, 0, 0, 0, 0, 0xc0})
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa := bytesToFloat64(a)
		fb := bytesToFloat64(b)
		if len(fa) != len(fb) || len(fa) == 0 {
			return
		}
		result := LogDivBatch(fa, fb)
		if len(result) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(result), len(fa))
		}
		for i := range result {
			assertFuzzClose(t, "log-div", result[i], ExactLogOracle(fa[i])/fb[i])
		}
	})
}

func assertFuzzExact(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.IsNaN(got) && math.IsNaN(want) {
		return
	}
	if math.IsNaN(got) || math.IsNaN(want) {
		t.Fatalf("%s result = %g, want %g", name, got, want)
	}
	if got != want || (got == 0 && math.Signbit(got) != math.Signbit(want)) {
		t.Fatalf("%s result = %g (%016x), want %g (%016x)",
			name, got, math.Float64bits(got), want, math.Float64bits(want))
	}
}

func assertFuzzClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.IsNaN(got) && math.IsNaN(want) {
		return
	}
	if math.IsNaN(got) || math.IsNaN(want) {
		t.Fatalf("%s result = %g, want %g", name, got, want)
	}
	if math.IsInf(got, 0) || math.IsInf(want, 0) {
		if got != want {
			t.Fatalf("%s result = %g, want %g", name, got, want)
		}
		return
	}
	if got == 0 || want == 0 {
		if got != want || math.Signbit(got) != math.Signbit(want) {
			t.Fatalf("%s result = %g, want %g", name, got, want)
		}
		return
	}
	if math.Abs(want) < math.Float64frombits(0x0010000000000000) {
		if distance := ulpDistance(got, want); distance > 4 {
			t.Fatalf("%s result = %g, want %g (%d ULP)", name, got, want, distance)
		}
		return
	}
	err := math.Abs(got - want)
	budget := 1e-14 + 1e-15*math.Abs(want)
	if err > budget {
		t.Fatalf("%s result = %.17g, want %.17g (error %.3g > budget %.3g)",
			name, got, want, err, budget)
	}
}

func assertFuzzExact32(t *testing.T, name string, got, want float32) {
	t.Helper()
	if math.IsNaN(float64(got)) && math.IsNaN(float64(want)) {
		return
	}
	if got != want || (got == 0 && math.Signbit(float64(got)) != math.Signbit(float64(want))) {
		t.Fatalf("%s result = %g (%08x), want %g (%08x)",
			name, got, math.Float32bits(got), want, math.Float32bits(want))
	}
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

func bytesToFloat32(b []byte) []float32 {
	n := len(b) / 4
	if n == 0 {
		return nil
	}
	res := make([]float32, n)
	for i := 0; i < n; i++ {
		res[i] = math.Float32frombits(uint32(b[i*4]) | uint32(b[i*4+1])<<8 |
			uint32(b[i*4+2])<<16 | uint32(b[i*4+3])<<24)
	}
	return res
}

func FuzzFloorSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		got := FloorSIMD(fa)
		if len(got) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(fa))
		}
		for i := range got {
			assertFuzzExact(t, "floor", got[i], math.Floor(fa[i]))
		}
	})
}

func FuzzCeilSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		got := CeilSIMD(fa)
		if len(got) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(fa))
		}
		for i := range got {
			assertFuzzExact(t, "ceil", got[i], math.Ceil(fa[i]))
		}
	})
}

func FuzzTruncSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		got := TruncSIMD(fa)
		if len(got) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(fa))
		}
		for i := range got {
			assertFuzzExact(t, "trunc", got[i], math.Trunc(fa[i]))
		}
	})
}

func FuzzRoundSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f})
	f.Fuzz(func(t *testing.T, a []byte) {
		fa := bytesToFloat64(a)
		got := RoundSIMD(fa)
		if len(got) != len(fa) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(fa))
		}
		for i := range got {
			assertFuzzExact(t, "round", got[i], math.RoundToEven(fa[i]))
		}
	})
}

func FuzzMaxSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f}, []byte{0, 0, 0, 0, 0, 0, 0, 0x40})
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa, fb := bytesToFloat64(a), bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		got := MaxSIMD(fa, fb)
		for i := range got {
			var want float64
			if math.IsNaN(fa[i]) || math.IsNaN(fb[i]) {
				want = math.NaN()
			} else if fa[i] > fb[i] {
				want = fa[i]
			} else {
				want = fb[i]
			}
			assertFuzzExact(t, "max", got[i], want)
		}
	})
}

func FuzzMinSIMD(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0xf0, 0x3f}, []byte{0, 0, 0, 0, 0, 0, 0, 0x40})
	f.Fuzz(func(t *testing.T, a, b []byte) {
		fa, fb := bytesToFloat64(a), bytesToFloat64(b)
		if len(fa) != len(fb) {
			return
		}
		got := MinSIMD(fa, fb)
		for i := range got {
			var want float64
			if math.IsNaN(fa[i]) || math.IsNaN(fb[i]) {
				want = math.NaN()
			} else if fa[i] < fb[i] {
				want = fa[i]
			} else {
				want = fb[i]
			}
			assertFuzzExact(t, "min", got[i], want)
		}
	})
}

func FuzzFmaSIMDF32(f *testing.F) {
	f.Add(
		[]byte{0, 0, 0x80, 0x3f},
		[]byte{0, 0, 0, 0x40},
		[]byte{0, 0, 0x80, 0xbf},
	)
	f.Fuzz(func(t *testing.T, a, b, c []byte) {
		fa, fb, fc := bytesToFloat32(a), bytesToFloat32(b), bytesToFloat32(c)
		if len(fa) != len(fb) || len(fa) != len(fc) {
			return
		}
		got := FmaSIMDF32(fa, fb, fc)
		for i := range got {
			want := float32(math.FMA(float64(fa[i]), float64(fb[i]), float64(fc[i])))
			if math.IsNaN(float64(got[i])) && math.IsNaN(float64(want)) {
				continue
			}
			if math.IsInf(float64(got[i]), 0) && math.IsInf(float64(want), 0) && math.Signbit(float64(got[i])) == math.Signbit(float64(want)) {
				continue
			}
			if !f32Close(got[i], want, 1e-4) {
				t.Fatalf("fma32 i=%d got %v want %v", i, got[i], want)
			}
		}
	})
}
