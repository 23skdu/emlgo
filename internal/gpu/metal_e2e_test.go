//go:build darwin && arm64 && metal && cgo && !cuda

package gpu

import (
	"errors"
	"math"
	"math/cmplx"
	"strings"
	"testing"
)

func requireMetalDevice(t *testing.T) *Device {
	t.Helper()
	if err := Init(); err != nil {
		t.Fatalf("Metal Init() failed: %v", err)
	}
	devices, err := GetDevices()
	if err != nil {
		t.Fatalf("GetDevices() failed: %v", err)
	}
	if len(devices) == 0 {
		t.Fatal("expected at least 1 Apple Metal GPU device on darwin/arm64")
	}
	return &devices[0]
}

func TestMetalDeviceDiscoveryAndStatus(t *testing.T) {
	d := requireMetalDevice(t)
	if d.Name == "" {
		t.Error("expected non-empty Metal device name")
	}
	if d.MaxThreadsPerBlock <= 0 {
		t.Errorf("expected MaxThreadsPerBlock > 0, got %d", d.MaxThreadsPerBlock)
	}
	if d.WarpSize != 32 {
		t.Errorf("expected Apple Silicon SIMD width 32, got %d", d.WarpSize)
	}
	status := Status()
	if !strings.Contains(status, "Found") || !strings.Contains(status, "Metal GPU") {
		t.Errorf("unexpected Metal Status(): %q", status)
	}
	t.Logf("Metal Status:\n%s", status)

	// Test Status() with mocked error and empty devices
	restore := SetGetDevicesForTesting(func() ([]Device, error) {
		return nil, errors.New("mock error")
	})
	if s := Status(); !strings.Contains(s, "GPU Error") {
		t.Errorf("expected GPU Error status, got %q", s)
	}
	restore()

	restore = SetGetDevicesForTesting(func() ([]Device, error) {
		return nil, nil
	})
	if s := Status(); !strings.Contains(s, "No Metal GPU") {
		t.Errorf("expected No Metal GPU status, got %q", s)
	}
	restore()
}

func TestMetalInitShutdownIdempotent(t *testing.T) {
	Shutdown()
	if err := Init(); err != nil {
		t.Fatalf("Init after Shutdown failed: %v", err)
	}
	if err := Init(); err != nil {
		t.Fatalf("second Init failed: %v", err)
	}
	d := requireMetalDevice(t)
	res, err := d.ExpBatchF32([]float32{0.0, 1.0})
	if err != nil {
		t.Fatalf("ExpBatchF32 after re-init failed: %v", err)
	}
	if len(res) != 2 || ulpDiffF32(res[0], 1.0) > 1 {
		t.Fatalf("unexpected ExpBatchF32 result after re-init: %v", res)
	}
}

func TestMetalFloat32UnaryKernelsAccuracy(t *testing.T) {
	d := requireMetalDevice(t)

	general := []float32{0.0, 0.25, 0.5, 1.0, 1.5, 2.0, -0.25, -0.5, -1.0, -2.0, 0.01, -0.01, 3.0, -3.0}
	positive := []float32{0.001, 0.01, 0.1, 0.5, 1.0, 2.0, 10.0, 100.0, float32(math.E), float32(math.Pi)}
	nonNeg := []float32{0.0, 0.01, 0.25, 1.0, 2.0, 4.0, 9.0, 100.0, 1e6}
	tanSafe := []float32{0.0, 0.1, 0.5, 1.0, -0.1, -0.5, -1.0, 1.2, -1.2}

	tests := []struct {
		name   string
		input  []float32
		maxULP uint64
		batch  func([]float32) ([]float32, error)
		to     func([]float32, []float32) error
		cpuRef func(float64) float64
	}{
		{"ExpF32", general, 4, d.ExpBatchF32, d.ExpBatchToF32, math.Exp},
		{"LogF32", positive, 4, d.LogBatchF32, d.LogBatchToF32, math.Log},
		{"SinF32", general, 4, d.SinBatchF32, d.SinBatchToF32, math.Sin},
		{"CosF32", general, 4, d.CosBatchF32, d.CosBatchToF32, math.Cos},
		{"TanF32", tanSafe, 6, d.TanBatchF32, d.TanBatchToF32, math.Tan},
		{"SinhF32", general, 6, d.SinhBatchF32, d.SinhBatchToF32, math.Sinh},
		{"CoshF32", general, 6, d.CoshBatchF32, d.CoshBatchToF32, math.Cosh},
		{"TanhF32", general, 6, d.TanhBatchF32, d.TanhBatchToF32, math.Tanh},
		{"SqrtF32", nonNeg, 1, d.SqrtBatchF32, d.SqrtBatchToF32, math.Sqrt},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.batch(tc.input)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			v := &BatchVerifier{MaxULP: tc.maxULP}
			maxULP, failed, err := v.VerifyOpF32(tc.name, tc.input, res, tc.cpuRef)
			if err != nil {
				t.Fatalf("VerifyOpF32 failed: %v", err)
			}
			if failed > 0 {
				t.Errorf("%s: %d/%d exceeded %d ULP (max ULP=%d)", tc.name, failed, len(tc.input), tc.maxULP, maxULP)
			}
			t.Logf("%s: max ULP = %d (limit %d)", tc.name, maxULP, tc.maxULP)

			// Verify *BatchToF32 matches *BatchF32 and handles in-place aliasing
			dst := make([]float32, len(tc.input))
			if toErr := tc.to(dst, tc.input); toErr != nil {
				t.Fatalf("%s To failed: %v", tc.name, toErr)
			}
			for i := range dst {
				if dst[i] != res[i] {
					t.Fatalf("%s To[%d] = %v, want %v", tc.name, i, dst[i], res[i])
				}
			}

			// In-place aliasing
			alias := append([]float32(nil), tc.input...)
			if aliasErr := tc.to(alias, alias); aliasErr != nil {
				t.Fatalf("%s in-place failed: %v", tc.name, aliasErr)
			}
			for i := range alias {
				if alias[i] != res[i] {
					t.Fatalf("%s in-place[%d] = %v, want %v", tc.name, i, alias[i], res[i])
				}
			}

			// Empty slice
			emptyRes, err := tc.batch(nil)
			if err != nil || emptyRes != nil {
				t.Fatalf("%s(nil) = (%v, %v), want (nil, nil)", tc.name, emptyRes, err)
			}
			if err := tc.to(nil, nil); err != nil {
				t.Fatalf("%s To(nil, nil) failed: %v", tc.name, err)
			}

			// Length mismatch
			if err := tc.to(make([]float32, 1), tc.input); err == nil {
				t.Fatalf("%s To length mismatch should error", tc.name)
			}
		})
	}
}

func TestMetalFloat32EmlAndDotAccuracy(t *testing.T) {
	d := requireMetalDevice(t)

	x := []float32{0.0, 0.5, 1.0, -0.5, 2.0, -1.0, 0.25}
	y := []float32{1.0, 2.0, float32(math.E), 0.5, 4.0, 10.0, 0.25}

	res, err := d.EmlBatchF32(x, y)
	if err != nil {
		t.Fatalf("EmlBatchF32 failed: %v", err)
	}
	v := &BatchVerifier{MaxULP: 8}
	maxULP, failed, err := v.VerifyBinaryOpF32("EmlF32", x, y, res, func(a, b float64) float64 {
		return math.Exp(a) - math.Log(b)
	})
	if err != nil {
		t.Fatalf("VerifyBinaryOpF32 failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("EmlBatchF32: %d/%d exceeded 8 ULP (max ULP=%d)", failed, len(x), maxULP)
	}
	t.Logf("EmlBatchF32: max ULP = %d", maxULP)

	dst := make([]float32, len(x))
	if err := d.EmlBatchToF32(dst, x, y); err != nil {
		t.Fatalf("EmlBatchToF32 failed: %v", err)
	}
	for i := range dst {
		if dst[i] != res[i] {
			t.Fatalf("EmlBatchToF32[%d] = %v, want %v", i, dst[i], res[i])
		}
	}

	// Empty & mismatch checks
	if r, err := d.EmlBatchF32(nil, nil); err != nil || r != nil {
		t.Fatalf("EmlBatchF32(nil, nil) = (%v, %v)", r, err)
	}
	if _, err := d.EmlBatchF32(x, y[:2]); err == nil {
		t.Fatal("EmlBatchF32 mismatch should error")
	}
	if err := d.EmlBatchToF32(nil, nil, nil); err != nil {
		t.Fatalf("EmlBatchToF32(nil, nil, nil) failed: %v", err)
	}
	if err := d.EmlBatchToF32(dst[:2], x, y); err == nil {
		t.Fatal("EmlBatchToF32 mismatch should error")
	}

	// DotBatchF32 across multiple sizes including multi-threadgroup (>256)
	for _, n := range []int{0, 1, 7, 255, 256, 257, 1025, 4096} {
		a := make([]float32, n)
		b := make([]float32, n)
		var expected float64
		for i := range n {
			a[i] = float32((i%13)-6) * 0.25
			b[i] = float32((i%11)-5) * 0.5
			expected += float64(a[i]) * float64(b[i])
		}
		got, err := d.DotBatchF32(a, b)
		if err != nil {
			t.Fatalf("DotBatchF32(n=%d) failed: %v", n, err)
		}
		if math.Abs(float64(got)-expected) > 1e-3 {
			t.Errorf("DotBatchF32(n=%d) = %v, want %v", n, got, expected)
		}
	}
	if _, err := d.DotBatchF32([]float32{1}, []float32{1, 2}); err == nil {
		t.Fatal("DotBatchF32 length mismatch should error")
	}
}

func TestMetalComplex64KernelsAccuracy(t *testing.T) {
	d := requireMetalDevice(t)

	in := []complex64{
		complex(1.0, 0.0),
		complex(0.0, 1.0),
		complex(1.0, 1.0),
		complex(-0.5, 0.5),
		complex(0.25, -0.75),
		complex(-1.2, -0.3),
	}

	tests := []struct {
		name   string
		tol    float64
		batch  func([]complex64) ([]complex64, error)
		to     func([]complex64, []complex64) error
		cpuRef func(complex128) complex128
	}{
		{"ExpC64", 1e-5, d.ExpBatchC64, d.ExpBatchToC64, cmplx.Exp},
		{"LogC64", 1e-5, d.LogBatchC64, d.LogBatchToC64, cmplx.Log},
		{"SinC64", 1e-5, d.SinBatchC64, d.SinBatchToC64, cmplx.Sin},
		{"CosC64", 1e-5, d.CosBatchC64, d.CosBatchToC64, cmplx.Cos},
		{"TanC64", 2e-5, d.TanBatchC64, d.TanBatchToC64, cmplx.Tan},
		{"SinhC64", 1e-5, d.SinhBatchC64, d.SinhBatchToC64, cmplx.Sinh},
		{"CoshC64", 1e-5, d.CoshBatchC64, d.CoshBatchToC64, cmplx.Cosh},
		{"TanhC64", 2e-5, d.TanhBatchC64, d.TanhBatchToC64, cmplx.Tanh},
		{"SqrtC64", 1e-5, d.SqrtBatchC64, d.SqrtBatchToC64, cmplx.Sqrt},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.batch(in)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			for i, got := range res {
				want := complex64(tc.cpuRef(complex128(in[i])))
				diff := cmplx.Abs(complex128(got - want))
				if diff > tc.tol {
					t.Errorf("%s[%d](%v) = %v, want %v (diff=%v > %v)", tc.name, i, in[i], got, want, diff, tc.tol)
				}
			}

			dst := make([]complex64, len(in))
			if err := tc.to(dst, in); err != nil {
				t.Fatalf("%s To failed: %v", tc.name, err)
			}
			for i := range dst {
				if dst[i] != res[i] {
					t.Fatalf("%s To[%d] = %v, want %v", tc.name, i, dst[i], res[i])
				}
			}
		})
	}

	// EmlBatchC64 & EmlBatchToC64
	y := []complex64{
		complex(1.0, 0.5),
		complex(2.0, 0.0),
		complex(0.5, 1.0),
		complex(1.5, -0.5),
		complex(2.0, 2.0),
		complex(1.0, 1.0),
	}
	emlRes, err := d.EmlBatchC64(in, y)
	if err != nil {
		t.Fatalf("EmlBatchC64 failed: %v", err)
	}
	for i, got := range emlRes {
		want := complex64(cmplx.Exp(complex128(in[i])) - cmplx.Log(complex128(y[i])))
		if diff := cmplx.Abs(complex128(got - want)); diff > 1e-5 {
			t.Errorf("EmlBatchC64[%d] = %v, want %v (diff=%v)", i, got, want, diff)
		}
	}
	emlDst := make([]complex64, len(in))
	if err := d.EmlBatchToC64(emlDst, in, y); err != nil {
		t.Fatalf("EmlBatchToC64 failed: %v", err)
	}
	for i := range emlDst {
		if emlDst[i] != emlRes[i] {
			t.Fatalf("EmlBatchToC64[%d] = %v, want %v", i, emlDst[i], emlRes[i])
		}
	}

	// DotBatchC64 (Hermitian inner product: sum(a * conj(b))) across small and multi-threadgroup sizes
	for _, n := range []int{0, 2, 17, 513} {
		a := make([]complex64, n)
		b := make([]complex64, n)
		var want complex128
		for i := range n {
			a[i] = complex(float32((i%5)-2), float32((i%7)-3))
			b[i] = complex(float32((i%3)-1), float32((i%4)-2))
			want += complex128(a[i]) * cmplx.Conj(complex128(b[i]))
		}
		got, err := d.DotBatchC64(a, b)
		if err != nil {
			t.Fatalf("DotBatchC64(n=%d) failed: %v", n, err)
		}
		if diff := cmplx.Abs(complex128(got) - want); diff > 1e-2 {
			t.Errorf("DotBatchC64(n=%d) = %v, want %v (diff=%v)", n, got, want, diff)
		}
	}
	if _, err := d.DotBatchC64(in, in[:2]); err == nil {
		t.Fatal("DotBatchC64 length mismatch should error")
	}
}

func TestMetalLargeBatchesAndNonUniformSizes(t *testing.T) {
	d := requireMetalDevice(t)
	for _, n := range []int{1, 3, 17, 255, 256, 257, 1024, 16384, 65537} {
		x := make([]float32, n)
		for i := range x {
			x[i] = float32(i%100+1) / 100.0
		}
		res, err := d.ExpBatchF32(x)
		if err != nil {
			t.Fatalf("ExpBatchF32(n=%d) failed: %v", n, err)
		}
		if len(res) != n {
			t.Fatalf("ExpBatchF32(n=%d) len=%d", n, len(res))
		}
		// Spot-check start, middle, end
		for _, idx := range []int{0, n / 2, n - 1} {
			want := float32(math.Exp(float64(x[idx])))
			if ulp := ulpDiffF32(res[idx], want); ulp > 4 {
				t.Errorf("ExpBatchF32(n=%d)[%d] = %v, want %v (ulp=%d)", n, idx, res[idx], want, ulp)
			}
		}
	}
}

func TestMetalFloat64AndComplex128ExplicitlyUnsupported(t *testing.T) {
	d := requireMetalDevice(t)
	f64 := []float64{1.0, 2.0}
	dst64 := make([]float64, 2)
	c128 := []complex128{1 + 2i, 3 + 4i}
	dst128 := make([]complex128, 2)

	check := func(name string, err error) {
		t.Helper()
		if !errors.Is(err, ErrMetalFloat64Unsupported) {
			t.Errorf("%s: got err=%v, want ErrMetalFloat64Unsupported", name, err)
		}
	}

	_, err := d.ExpBatch(f64)
	check("ExpBatch", err)
	check("ExpBatchTo", d.ExpBatchTo(dst64, f64))
	_, err = d.LogBatch(f64)
	check("LogBatch", err)
	check("LogBatchTo", d.LogBatchTo(dst64, f64))
	_, err = d.SinBatch(f64)
	check("SinBatch", err)
	check("SinBatchTo", d.SinBatchTo(dst64, f64))
	_, err = d.CosBatch(f64)
	check("CosBatch", err)
	check("CosBatchTo", d.CosBatchTo(dst64, f64))
	_, err = d.TanBatch(f64)
	check("TanBatch", err)
	check("TanBatchTo", d.TanBatchTo(dst64, f64))
	_, err = d.SinhBatch(f64)
	check("SinhBatch", err)
	check("SinhBatchTo", d.SinhBatchTo(dst64, f64))
	_, err = d.CoshBatch(f64)
	check("CoshBatch", err)
	check("CoshBatchTo", d.CoshBatchTo(dst64, f64))
	_, err = d.TanhBatch(f64)
	check("TanhBatch", err)
	check("TanhBatchTo", d.TanhBatchTo(dst64, f64))
	_, err = d.SqrtBatch(f64)
	check("SqrtBatch", err)
	check("SqrtBatchTo", d.SqrtBatchTo(dst64, f64))
	_, err = d.EmlBatch(f64, f64)
	check("EmlBatch", err)
	check("EmlBatchTo", d.EmlBatchTo(dst64, f64, f64))
	_, err = d.DotBatch(f64, f64)
	check("DotBatch", err)
	_, err = d.AbsBatch(f64)
	check("AbsBatch", err)
	_, err = d.NegBatch(f64)
	check("NegBatch", err)
	_, err = d.PowBatch(f64, 2.0)
	check("PowBatch", err)
	_, err = d.InvBatch(f64)
	check("InvBatch", err)
	_, err = d.FmaBatch(f64, f64, f64)
	check("FmaBatch", err)

	_, err = d.ExpBatchC128(c128)
	check("ExpBatchC128", err)
	check("ExpBatchToC128", d.ExpBatchToC128(dst128, c128))
	_, err = d.LogBatchC128(c128)
	check("LogBatchC128", err)
	check("LogBatchToC128", d.LogBatchToC128(dst128, c128))
	_, err = d.SinBatchC128(c128)
	check("SinBatchC128", err)
	check("SinBatchToC128", d.SinBatchToC128(dst128, c128))
	_, err = d.CosBatchC128(c128)
	check("CosBatchC128", err)
	check("CosBatchToC128", d.CosBatchToC128(dst128, c128))
	_, err = d.TanBatchC128(c128)
	check("TanBatchC128", err)
	check("TanBatchToC128", d.TanBatchToC128(dst128, c128))
	_, err = d.SinhBatchC128(c128)
	check("SinhBatchC128", err)
	check("SinhBatchToC128", d.SinhBatchToC128(dst128, c128))
	_, err = d.CoshBatchC128(c128)
	check("CoshBatchC128", err)
	check("CoshBatchToC128", d.CoshBatchToC128(dst128, c128))
	_, err = d.TanhBatchC128(c128)
	check("TanhBatchC128", err)
	check("TanhBatchToC128", d.TanhBatchToC128(dst128, c128))
	_, err = d.SqrtBatchC128(c128)
	check("SqrtBatchC128", err)
	check("SqrtBatchToC128", d.SqrtBatchToC128(dst128, c128))
	_, err = d.EmlBatchC128(c128, c128)
	check("EmlBatchC128", err)
	check("EmlBatchToC128", d.EmlBatchToC128(dst128, c128, c128))
	_, err = d.DotBatchC128(c128, c128)
	check("DotBatchC128", err)
}

func TestMetalStreamsAndPinnedMemory(t *testing.T) {
	s, err := NewStream()
	if err != nil {
		t.Fatalf("NewStream failed: %v", err)
	}
	if syncErr := s.Sync(); syncErr != nil {
		t.Errorf("Stream.Sync failed: %v", syncErr)
	}
	if destroyErr := s.Destroy(); destroyErr != nil {
		t.Errorf("Stream.Destroy failed: %v", destroyErr)
	}
	var nilStream *Stream
	if nilErr := nilStream.Destroy(); nilErr != nil {
		t.Errorf("nil Stream.Destroy failed: %v", nilErr)
	}

	p64, err := AllocatePinned(16)
	if err != nil || len(p64) != 16 {
		t.Fatalf("AllocatePinned failed: len=%d err=%v", len(p64), err)
	}
	FreePinned(p64)

	p32, err := AllocatePinnedF32(16)
	if err != nil || len(p32) != 16 {
		t.Fatalf("AllocatePinnedF32 failed: len=%d err=%v", len(p32), err)
	}
	FreePinnedF32(p32)

	pc64, err := AllocatePinnedC64(16)
	if err != nil || len(pc64) != 16 {
		t.Fatalf("AllocatePinnedC64 failed: len=%d err=%v", len(pc64), err)
	}
	FreePinnedC64(pc64)

	pc128, err := AllocatePinnedC128(16)
	if err != nil || len(pc128) != 16 {
		t.Fatalf("AllocatePinnedC128 failed: len=%d err=%v", len(pc128), err)
	}
	FreePinnedC128(pc128)

	if err := CopyToDeviceAsync(nil, nil, 0, s); err == nil {
		t.Error("expected CopyToDeviceAsync error on Metal unified memory")
	}
	if err := CopyToHostAsync(nil, nil, 0, s); err == nil {
		t.Error("expected CopyToHostAsync error on Metal unified memory")
	}
}

func TestMetalErrorPaths(t *testing.T) {
	// Unknown kernel name triggers C launch failure
	if _, err := metalLaunchUnary("nonexistent_kernel", []float32{1.0}, 4); err == nil {
		t.Error("expected error for unknown unary kernel")
	}
	if _, err := metalLaunchBinary("nonexistent_kernel", []float32{1.0}, []float32{1.0}, 4); err == nil {
		t.Error("expected error for unknown binary kernel")
	}

	// Simulate Metal Init failure with invalid shader source
	Shutdown()
	savedSrc := metalShaderSource
	metalShaderSource = "this is not valid metal shader source !!!"
	defer func() {
		metalShaderSource = savedSrc
		Shutdown()
		_ = Init()
	}()

	if err := Init(); err == nil {
		t.Error("expected Init error for invalid MSL source")
	}
	if _, err := NewStream(); err == nil {
		t.Error("expected NewStream error when Init fails")
	}
	d := &Device{ID: 0}
	if err := d.ExpBatchToF32(make([]float32, 1), []float32{1.0}); err == nil {
		t.Error("expected ExpBatchToF32 error when Init fails")
	}
	if err := d.EmlBatchToF32(make([]float32, 1), []float32{1.0}, []float32{1.0}); err == nil {
		t.Error("expected EmlBatchToF32 error when Init fails")
	}
	if _, err := d.DotBatchF32([]float32{1.0}, []float32{1.0}); err == nil {
		t.Error("expected DotBatchF32 error when Init fails")
	}
	if _, err := d.DotBatchC64([]complex64{1 + 0i}, []complex64{1 + 0i}); err == nil {
		t.Error("expected DotBatchC64 error when Init fails")
	}
}
