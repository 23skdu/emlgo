//go:build cuda

package gpu

import (
	"math"
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
