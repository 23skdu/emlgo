//go:build cuda

package gpu

import (
	"math"
	"math/cmplx"
	"testing"
)

// e2eTestData returns representative test values for GPU e2e validation.
func e2eTestData() []float64 {
	return []float64{
		0.0, 0.5, 1.0, 2.0, 10.0, 100.0,
		-0.5, -1.0, -2.0, -10.0,
		0.1, 0.01, 1e-10, 1e10,
		math.Pi, math.E,
	}
}

func requireGPU(t *testing.T) *Device {
	t.Helper()
	devices, err := GetDevices()
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}
	if len(devices) == 0 {
		t.Fatal("no GPU devices available")
	}
	return &devices[0]
}

func TestE2EGPUExpWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	data := e2eTestData()

	result, err := d.ExpBatch(data)
	if err != nil {
		t.Fatalf("ExpBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Exp", data, result, math.Exp)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("ExpBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("ExpBatch: max ULP = %d, all within tolerance", maxULP)
}

func TestE2EGPULogWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	// Only positive values for Log
	data := []float64{0.5, 1.0, 2.0, 10.0, 100.0, 0.1, 0.01, 1e-10, 1e10, math.E, math.Pi}

	result, err := d.LogBatch(data)
	if err != nil {
		t.Fatalf("LogBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Log", data, result, math.Log)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("LogBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("LogBatch: max ULP = %d", maxULP)
}

func TestE2EGPUSinWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	data := e2eTestData()

	result, err := d.SinBatch(data)
	if err != nil {
		t.Fatalf("SinBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Sin", data, result, math.Sin)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("SinBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("SinBatch: max ULP = %d", maxULP)
}

func TestE2EGPUCosWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	data := e2eTestData()

	result, err := d.CosBatch(data)
	if err != nil {
		t.Fatalf("CosBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Cos", data, result, math.Cos)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("CosBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("CosBatch: max ULP = %d", maxULP)
}

func TestE2EGPUTanWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	// Avoid poles for tan
	data := []float64{0.0, 0.5, 1.0, -0.5, -1.0, 0.1, -0.1}

	result, err := d.TanBatch(data)
	if err != nil {
		t.Fatalf("TanBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Tan", data, result, math.Tan)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("TanBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("TanBatch: max ULP = %d", maxULP)
}

func TestE2EGPUSinhWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	data := e2eTestData()

	result, err := d.SinhBatch(data)
	if err != nil {
		t.Fatalf("SinhBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Sinh", data, result, math.Sinh)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("SinhBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("SinhBatch: max ULP = %d", maxULP)
}

func TestE2EGPUCoshWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	data := e2eTestData()

	result, err := d.CoshBatch(data)
	if err != nil {
		t.Fatalf("CoshBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Cosh", data, result, math.Cosh)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("CoshBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("CoshBatch: max ULP = %d", maxULP)
}

func TestE2EGPUTanhWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	data := e2eTestData()

	result, err := d.TanhBatch(data)
	if err != nil {
		t.Fatalf("TanhBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Tanh", data, result, math.Tanh)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("TanhBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("TanhBatch: max ULP = %d", maxULP)
}

func TestE2EGPUSqrtWithin1ULP(t *testing.T) {
	d := requireGPU(t)
	verifier := DefaultVerifier()
	// Only non-negative for sqrt
	data := []float64{0.0, 0.5, 1.0, 2.0, 10.0, 100.0, 0.1, 0.01, 1e-10, 1e10, math.Pi, math.E}

	result, err := d.SqrtBatch(data)
	if err != nil {
		t.Fatalf("SqrtBatch failed: %v", err)
	}

	maxULP, failed, err := verifier.VerifyOp("Sqrt", data, result, math.Sqrt)
	if err != nil {
		t.Fatalf("VerifyOp failed: %v", err)
	}
	if failed > 0 {
		t.Errorf("SqrtBatch: %d/%d values exceeded 1 ULP (max ULP: %d)",
			failed, len(data), maxULP)
	}
	t.Logf("SqrtBatch: max ULP = %d", maxULP)
}

func TestE2EGPUBatchLarge(t *testing.T) {
	d := requireGPU(t)

	sizes := []int{1024, 16384, 262144}
	for _, n := range sizes {
		data := make([]float64, n)
		for i := range data {
			data[i] = float64(i%100+1) / 100.0
		}

		result, err := d.ExpBatch(data)
		if err != nil {
			t.Errorf("ExpBatch(n=%d) failed: %v", n, err)
			continue
		}
		if len(result) != n {
			t.Errorf("ExpBatch(n=%d) returned %d results, want %d", n, len(result), n)
			continue
		}
		// Verify results are finite and positive
		for _, v := range result {
			if math.IsNaN(v) || v <= 0 {
				t.Errorf("ExpBatch(n=%d): unexpected value %v", n, v)
				break
			}
		}
	}
}

func TestE2EGPUStream(t *testing.T) {
	d := requireGPU(t)
	data := []float64{0.0, 1.0, 2.0, 3.0}

	stream, err := NewStream()
	if err != nil {
		t.Fatalf("NewStream failed: %v", err)
	}
	defer stream.Destroy()

	result, err := d.ExpBatch(data)
	if err != nil {
		t.Fatalf("ExpBatch failed: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 results, got %d", len(result))
	}
	for i, v := range result {
		expected := math.Exp(data[i])
		ulp := ulpDiff(v, expected)
		if ulp > 10 {
			t.Errorf("result[%d] = %v, expected %v (ULP: %d)", i, v, expected, ulp)
		}
	}
}

func TestE2EGPUPinnedMemory(t *testing.T) {
	d := requireGPU(t)

	pinned, err := AllocatePinned(1024)
	if err != nil {
		t.Fatalf("AllocatePinned failed: %v", err)
	}
	defer FreePinned(pinned)

	for i := range pinned {
		pinned[i] = float64(i%100) / 100.0
	}

	result, err := d.ExpBatch(pinned)
	if err != nil {
		t.Fatalf("ExpBatch with pinned memory failed: %v", err)
	}
	if len(result) != 1024 {
		t.Fatalf("expected 1024 results, got %d", len(result))
	}
}

func TestE2EGPUMultipleDevices(t *testing.T) {
	devices, err := GetDevices()
	if err != nil {
		t.Fatalf("GetDevices failed: %v", err)
	}
	if len(devices) == 0 {
		t.Fatal("no GPU devices")
	}

	data := []float64{1.0, 2.0, 3.0}
	for i := range devices {
		result, err := devices[i].ExpBatch(data)
		if err != nil {
			t.Errorf("device %d ExpBatch failed: %v", i, err)
			continue
		}
		if len(result) != 3 {
			t.Errorf("device %d: expected 3 results, got %d", i, len(result))
		}
	}
}

func TestE2EGPUFloat32(t *testing.T) {
	d := requireGPU(t)

	f32Data := []float32{0.0, 0.5, 1.0, 2.0, 3.0, -1.0}
	res, err := d.ExpBatchF32(f32Data)
	if err != nil {
		t.Fatalf("ExpBatchF32 failed: %v", err)
	}
	for i, v := range res {
		expected := float32(math.Exp(float64(f32Data[i])))
		if math.Abs(float64(v-expected)) > 1e-5 {
			t.Errorf("ExpBatchF32[%d] = %v, expected %v", i, v, expected)
		}
	}

	dst := make([]float32, len(f32Data))
	if err := d.ExpBatchToF32(dst, f32Data); err != nil {
		t.Fatalf("ExpBatchToF32 failed: %v", err)
	}
	for i, v := range dst {
		if v != res[i] {
			t.Errorf("ExpBatchToF32[%d] = %v, want %v", i, v, res[i])
		}
	}

	// Test DotBatchF32
	a := []float32{1.0, 2.0, 3.0, 4.0}
	b := []float32{2.0, 0.5, 1.0, 2.0}
	dot, err := d.DotBatchF32(a, b)
	if err != nil {
		t.Fatalf("DotBatchF32 failed: %v", err)
	}
	var expectedDot float32 = 1.0*2.0 + 2.0*0.5 + 3.0*1.0 + 4.0*2.0 // 2 + 1 + 3 + 8 = 14
	if math.Abs(float64(dot-expectedDot)) > 1e-5 {
		t.Errorf("DotBatchF32 = %v, want %v", dot, expectedDot)
	}

	// Test Pinned F32
	pinned, err := AllocatePinnedF32(64)
	if err != nil {
		t.Fatalf("AllocatePinnedF32 failed: %v", err)
	}
	defer FreePinnedF32(pinned)
}

func TestE2EGPUComplex64(t *testing.T) {
	d := requireGPU(t)

	c64Data := []complex64{
		complex(1.0, 0.0),
		complex(0.0, 1.0),
		complex(1.0, 1.0),
		complex(-0.5, 0.5),
	}

	res, err := d.ExpBatchC64(c64Data)
	if err != nil {
		t.Fatalf("ExpBatchC64 failed: %v", err)
	}
	for i, v := range res {
		expected := complex64(cmplx.Exp(complex128(c64Data[i])))
		diff := cmplx.Abs(complex128(v - expected))
		if diff > 1e-5 {
			t.Errorf("ExpBatchC64[%d] = %v, expected %v (diff: %v)", i, v, expected, diff)
		}
	}

	// Test DotBatchC64 (Hermitian inner product: sum(a_i * conj(b_i)))
	a := []complex64{complex(1, 2), complex(3, 4)}
	b := []complex64{complex(2, -1), complex(1, 1)}
	// a[0] * conj(b[0]) = (1+2i)*(2+1i) = 2 + i + 4i - 2 = 5i
	// a[1] * conj(b[1]) = (3+4i)*(1-1i) = 3 - 3i + 4i + 4 = 7 + i
	// sum = 7 + 6i
	dot, err := d.DotBatchC64(a, b)
	if err != nil {
		t.Fatalf("DotBatchC64 failed: %v", err)
	}
	expectedDot := complex64(complex(7, 6))
	if cmplx.Abs(complex128(dot-expectedDot)) > 1e-5 {
		t.Errorf("DotBatchC64 = %v, want %v", dot, expectedDot)
	}

	pinned, err := AllocatePinnedC64(64)
	if err != nil {
		t.Fatalf("AllocatePinnedC64 failed: %v", err)
	}
	defer FreePinnedC64(pinned)
}

func TestE2EGPUComplex128(t *testing.T) {
	d := requireGPU(t)

	c128Data := []complex128{
		complex(1.0, 0.0),
		complex(0.0, 1.0),
		complex(1.0, 1.0),
		complex(-0.5, 0.5),
	}

	res, err := d.ExpBatchC128(c128Data)
	if err != nil {
		t.Fatalf("ExpBatchC128 failed: %v", err)
	}
	for i, v := range res {
		expected := cmplx.Exp(c128Data[i])
		diff := cmplx.Abs(v - expected)
		if diff > 1e-12 {
			t.Errorf("ExpBatchC128[%d] = %v, expected %v (diff: %v)", i, v, expected, diff)
		}
	}

	// Test DotBatchC128
	a := []complex128{complex(1, 2), complex(3, 4)}
	b := []complex128{complex(2, -1), complex(1, 1)}
	dot, err := d.DotBatchC128(a, b)
	if err != nil {
		t.Fatalf("DotBatchC128 failed: %v", err)
	}
	expectedDot := complex(7, 6)
	if cmplx.Abs(dot-expectedDot) > 1e-12 {
		t.Errorf("DotBatchC128 = %v, want %v", dot, expectedDot)
	}

	pinned, err := AllocatePinnedC128(64)
	if err != nil {
		t.Fatalf("AllocatePinnedC128 failed: %v", err)
	}
	defer FreePinnedC128(pinned)
}

func TestE2EGPUDotBatchF64(t *testing.T) {
	d := requireGPU(t)

	a := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	b := []float64{2.0, 3.0, 4.0, 5.0, 6.0}
	// 2 + 6 + 12 + 20 + 30 = 70
	dot, err := d.DotBatch(a, b)
	if err != nil {
		t.Fatalf("DotBatch failed: %v", err)
	}
	if math.Abs(dot-70.0) > 1e-10 {
		t.Errorf("DotBatch = %v, want 70.0", dot)
	}
}
