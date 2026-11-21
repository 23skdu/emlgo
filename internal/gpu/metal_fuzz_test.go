//go:build darwin && arm64 && metal && cgo && !cuda

package gpu

import (
	"math"
	"math/cmplx"
	"testing"
)

// FuzzMetalFloat32Kernels fuzzes Metal GPU float32 compute kernels against
// IEEE-754 float64->float32 CPU reference implementations on Apple Silicon.
func FuzzMetalFloat32Kernels(f *testing.F) {
	f.Add(float32(0.0), float32(1.0))
	f.Add(float32(0.5), float32(2.0))
	f.Add(float32(-1.5), float32(0.5))
	f.Add(float32(3.1415927), float32(2.7182817))

	d := &Device{ID: 0}
	if err := Init(); err != nil {
		f.Fatalf("Metal Init failed: %v", err)
	}

	f.Fuzz(func(t *testing.T, x, y float32) {
		if math.IsNaN(float64(x)) || math.IsInf(float64(x), 0) ||
			math.IsNaN(float64(y)) || math.IsInf(float64(y), 0) {
			return
		}
		// Keep inputs in well-conditioned non-overflow float32 ranges
		if x < -10 || x > 10 {
			return
		}
		in := []float32{x}

		expRes, err := d.ExpBatchF32(in)
		if err != nil {
			t.Fatalf("ExpBatchF32 failed: %v", err)
		}
		wantExp := float32(math.Exp(float64(x)))
		if ulp := ulpDiffF32(expRes[0], wantExp); ulp > 8 {
			t.Errorf("ExpBatchF32(%v) = %v, want %v (ulp=%d)", x, expRes[0], wantExp, ulp)
		}

		sinRes, err := d.SinBatchF32(in)
		if err != nil {
			t.Fatalf("SinBatchF32 failed: %v", err)
		}
		wantSin := float32(math.Sin(float64(x)))
		if ulp := ulpDiffF32(sinRes[0], wantSin); ulp > 8 {
			t.Errorf("SinBatchF32(%v) = %v, want %v (ulp=%d)", x, sinRes[0], wantSin, ulp)
		}

		cosRes, err := d.CosBatchF32(in)
		if err != nil {
			t.Fatalf("CosBatchF32 failed: %v", err)
		}
		wantCos := float32(math.Cos(float64(x)))
		if ulp := ulpDiffF32(cosRes[0], wantCos); ulp > 8 {
			t.Errorf("CosBatchF32(%v) = %v, want %v (ulp=%d)", x, cosRes[0], wantCos, ulp)
		}

		if y > 1e-4 && y < 1e4 {
			logRes, err := d.LogBatchF32([]float32{y})
			if err != nil {
				t.Fatalf("LogBatchF32 failed: %v", err)
			}
			wantLog := float32(math.Log(float64(y)))
			if ulp := ulpDiffF32(logRes[0], wantLog); ulp > 8 {
				t.Errorf("LogBatchF32(%v) = %v, want %v (ulp=%d)", y, logRes[0], wantLog, ulp)
			}

			sqrtRes, err := d.SqrtBatchF32([]float32{y})
			if err != nil {
				t.Fatalf("SqrtBatchF32 failed: %v", err)
			}
			wantSqrt := float32(math.Sqrt(float64(y)))
			if ulp := ulpDiffF32(sqrtRes[0], wantSqrt); ulp > 2 {
				t.Errorf("SqrtBatchF32(%v) = %v, want %v (ulp=%d)", y, sqrtRes[0], wantSqrt, ulp)
			}

			emlRes, err := d.EmlBatchF32([]float32{x}, []float32{y})
			if err != nil {
				t.Fatalf("EmlBatchF32 failed: %v", err)
			}
			wantEml := float32(math.Exp(float64(x)) - math.Log(float64(y)))
			if !math.IsInf(float64(wantEml), 0) && math.Abs(float64(emlRes[0]-wantEml)) > 1e-3*math.Max(1, math.Abs(float64(wantEml))) {
				t.Errorf("EmlBatchF32(%v, %v) = %v, want %v", x, y, emlRes[0], wantEml)
			}
		}
	})
}

// FuzzMetalComplex64Kernels fuzzes Metal GPU complex64 kernels against math/cmplx.
func FuzzMetalComplex64Kernels(f *testing.F) {
	f.Add(float32(0.5), float32(-0.5))
	f.Add(float32(1.0), float32(1.0))
	f.Add(float32(-2.0), float32(0.25))

	d := &Device{ID: 0}
	if err := Init(); err != nil {
		f.Fatalf("Metal Init failed: %v", err)
	}

	f.Fuzz(func(t *testing.T, re, im float32) {
		if math.IsNaN(float64(re)) || math.IsInf(float64(re), 0) ||
			math.IsNaN(float64(im)) || math.IsInf(float64(im), 0) {
			return
		}
		if re < -5 || re > 5 || im < -5 || im > 5 {
			return
		}
		z := complex(re, im)
		res, err := d.ExpBatchC64([]complex64{z})
		if err != nil {
			t.Fatalf("ExpBatchC64 failed: %v", err)
		}
		want := complex64(cmplx.Exp(complex128(z)))
		tol := 1e-4 * math.Max(1.0, cmplx.Abs(complex128(want)))
		if diff := cmplx.Abs(complex128(res[0] - want)); diff > tol {
			t.Errorf("ExpBatchC64(%v) = %v, want %v (diff=%v)", z, res[0], want, diff)
		}
	})
}
