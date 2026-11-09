package arithmetic

import (
	"math"
	"testing"
)

func FuzzAbs(f *testing.F) {
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(0.0)
	f.Fuzz(func(t *testing.T, a float64) {
		got := Abs(a)
		want := math.Abs(a)
		if math.IsNaN(a) {
			if !math.IsNaN(got) {
				t.Fatalf("Abs(NaN) = %v, want NaN", got)
			}
			return
		}
		if got != want {
			t.Fatalf("Abs(%g) = %g, want %g", a, got, want)
		}
	})
}

func FuzzAdd(f *testing.F) {
	f.Add(1.0, 2.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		got := Add(a, b)
		want := a + b
		if math.IsNaN(want) && math.IsNaN(got) {
			return
		}
		if got != want {
			t.Fatalf("Add(%g, %g) = %g, want %g", a, b, got, want)
		}
	})
}

func FuzzMod(f *testing.F) {
	f.Add(10.0, 3.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		got := Mod(a, b)
		want := math.Mod(a, b)
		if math.IsNaN(want) && math.IsNaN(got) {
			return
		}
		if got != want {
			t.Fatalf("Mod(%g, %g) = %g, want %g", a, b, got, want)
		}
	})
}

func FuzzRemainder(f *testing.F) {
	f.Add(10.0, 3.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		got := Remainder(a, b)
		want := math.Remainder(a, b)
		if math.IsNaN(want) && math.IsNaN(got) {
			return
		}
		if got != want {
			t.Fatalf("Remainder(%g, %g) = %g, want %g", a, b, got, want)
		}
	})
}

func FuzzPow(f *testing.F) {
	f.Add(2.0, 3.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		got := Pow(a, b)
		if a > 0.1 && a < 10.0 && b >= -5.0 && b <= 5.0 {
			want := math.Pow(a, b)
			if math.Abs(got-want) > 1e-11*math.Max(1.0, math.Abs(want)) {
				t.Fatalf("Pow(%g, %g) = %.17g, want %.17g", a, b, got, want)
			}
		}
	})
}

func FuzzLogBase(f *testing.F) {
	f.Add(8.0)
	f.Fuzz(func(t *testing.T, a float64) {
		l2 := LogBase2(a)
		l10 := LogBase10(a)
		if a <= 0 || math.IsNaN(a) {
			if !math.IsNaN(l2) || !math.IsNaN(l10) {
				t.Fatalf("LogBase2/10(%g) = (%v, %v), want NaN", a, l2, l10)
			}
		}
	})
}

func FuzzGCDAndLCM(f *testing.F) {
	f.Add(int64(12), int64(18))
	f.Add(int64(-24), int64(36))
	f.Add(int64(math.MinInt64), int64(2))
	f.Fuzz(func(t *testing.T, a, b int64) {
		g := GCD(a, b)
		if g < 0 {
			t.Fatalf("GCD(%d, %d) = %d < 0", a, b, g)
		}
		if a != 0 && b != 0 && a != math.MinInt64 && b != math.MinInt64 && g > 0 {
			if a%g != 0 || b%g != 0 {
				t.Fatalf("GCD(%d, %d) = %d does not divide inputs", a, b, g)
			}
		}
		_ = LCM(a, b)
	})
}

func FuzzInt8DistanceAndDot(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, []byte{9, 8, 7, 6, 5, 4, 3, 2, 1})
	f.Fuzz(func(t *testing.T, a, b []byte) {
		if len(a) == 0 || len(a) != len(b) || len(a) > 1024 {
			return
		}
		sa := make([]int8, len(a))
		sb := make([]int8, len(b))
		var wantDot, wantL2 int32
		var wantUDot uint32
		for i := range a {
			sa[i] = int8(a[i])
			sb[i] = int8(b[i])
			wantDot += int32(sa[i]) * int32(sb[i])
			d := int32(sa[i]) - int32(sb[i])
			wantL2 += d * d
			wantUDot += uint32(a[i]) * uint32(b[i])
		}
		if got := DotProductInt8(sa, sb); got != wantDot {
			t.Fatalf("DotProductInt8 = %d, want %d", got, wantDot)
		}
		if got := DotProductUint8(a, b); got != wantUDot {
			t.Fatalf("DotProductUint8 = %d, want %d", got, wantUDot)
		}
		if got := L2SquaredInt8(sa, sb); got != wantL2 {
			t.Fatalf("L2SquaredInt8 = %d, want %d", got, wantL2)
		}
		if d := CosineDistanceInt8(sa, sb); d < -1e-5 || d > 2.0001 {
			t.Fatalf("CosineDistanceInt8 out of range [0, 2]: %v", d)
		}
		if d := CosineDistanceUint8(a, b); d < -1e-5 || d > 1.0001 {
			t.Fatalf("CosineDistanceUint8 out of range [0, 1]: %v", d)
		}
	})
}

func TestOracleArithmeticSuite(t *testing.T) {
	if got := GCD(math.MinInt64, 2); got != 2 {
		t.Fatalf("GCD(MinInt64, 2) = %d, want 2", got)
	}
	if !math.IsInf(Log1p(-1.0), -1) {
		t.Fatalf("Log1p(-1.0) = %v, want -Inf", Log1p(-1.0))
	}
	xs := []float64{0.25, 0.5, 1.0, 2.0, 4.0, 8.0, 16.0}
	gotSqrt := SqrtBatch(xs)
	gotLog2 := Log2Batch(xs)
	gotLog10 := Log10Batch(xs)
	for i, x := range xs {
		if gotSqrt[i] != math.Sqrt(x) {
			t.Fatalf("SqrtBatch[%d] = %v, want %v", i, gotSqrt[i], math.Sqrt(x))
		}
		if math.Abs(gotLog2[i]-math.Log2(x)) > 1e-14 {
			t.Fatalf("Log2Batch[%d] = %v, want %v", i, gotLog2[i], math.Log2(x))
		}
		if math.Abs(gotLog10[i]-math.Log10(x)) > 1e-14 {
			t.Fatalf("Log10Batch[%d] = %v, want %v", i, gotLog10[i], math.Log10(x))
		}
	}
}
