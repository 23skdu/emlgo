//go:build (!cuda || !cgo) && (!darwin || !arm64 || !cgo || cuda)

package gpu

import "testing"

func TestStubRemainingOpsCoverage(t *testing.T) {
	d := &Device{ID: 0, Name: "test-stub"}
	f64 := []float64{1.0, 2.0}
	f32 := []float32{1.0, 2.0}
	c64 := []complex64{1 + 2i, 3 + 4i}
	c128 := []complex128{1 + 2i, 3 + 4i}

	// Float64 batch ops
	_ = d.ExpBatchTo(f64, f64)
	_ = d.LogBatchTo(f64, f64)
	_ = d.SinBatchTo(f64, f64)
	_ = d.CosBatchTo(f64, f64)
	_ = d.TanBatchTo(f64, f64)
	_ = d.SinhBatchTo(f64, f64)
	_ = d.CoshBatchTo(f64, f64)
	_ = d.TanhBatchTo(f64, f64)
	_ = d.SqrtBatchTo(f64, f64)
	_ = d.EmlBatchTo(f64, f64, f64)
	_, _ = d.DotBatch(f64, f64)

	// Float32 batch ops
	_, _ = d.ExpBatchF32(f32)
	_ = d.ExpBatchToF32(f32, f32)
	_, _ = d.LogBatchF32(f32)
	_ = d.LogBatchToF32(f32, f32)
	_, _ = d.SinBatchF32(f32)
	_ = d.SinBatchToF32(f32, f32)
	_, _ = d.CosBatchF32(f32)
	_ = d.CosBatchToF32(f32, f32)
	_, _ = d.TanBatchF32(f32)
	_ = d.TanBatchToF32(f32, f32)
	_, _ = d.SinhBatchF32(f32)
	_ = d.SinhBatchToF32(f32, f32)
	_, _ = d.CoshBatchF32(f32)
	_ = d.CoshBatchToF32(f32, f32)
	_, _ = d.TanhBatchF32(f32)
	_ = d.TanhBatchToF32(f32, f32)
	_, _ = d.SqrtBatchF32(f32)
	_ = d.SqrtBatchToF32(f32, f32)
	_, _ = d.EmlBatchF32(f32, f32)
	_ = d.EmlBatchToF32(f32, f32, f32)
	_, _ = d.DotBatchF32(f32, f32)

	// Complex64 batch ops
	_, _ = d.ExpBatchC64(c64)
	_ = d.ExpBatchToC64(c64, c64)
	_, _ = d.LogBatchC64(c64)
	_ = d.LogBatchToC64(c64, c64)
	_, _ = d.SinBatchC64(c64)
	_ = d.SinBatchToC64(c64, c64)
	_, _ = d.CosBatchC64(c64)
	_ = d.CosBatchToC64(c64, c64)
	_, _ = d.TanBatchC64(c64)
	_ = d.TanBatchToC64(c64, c64)
	_, _ = d.SinhBatchC64(c64)
	_ = d.SinhBatchToC64(c64, c64)
	_, _ = d.CoshBatchC64(c64)
	_ = d.CoshBatchToC64(c64, c64)
	_, _ = d.TanhBatchC64(c64)
	_ = d.TanhBatchToC64(c64, c64)
	_, _ = d.SqrtBatchC64(c64)
	_ = d.SqrtBatchToC64(c64, c64)
	_, _ = d.EmlBatchC64(c64, c64)
	_ = d.EmlBatchToC64(c64, c64, c64)
	_, _ = d.DotBatchC64(c64, c64)

	// Complex128 batch ops
	_, _ = d.ExpBatchC128(c128)
	_ = d.ExpBatchToC128(c128, c128)
	_, _ = d.LogBatchC128(c128)
	_ = d.LogBatchToC128(c128, c128)
	_, _ = d.SinBatchC128(c128)
	_ = d.SinBatchToC128(c128, c128)
	_, _ = d.CosBatchC128(c128)
	_ = d.CosBatchToC128(c128, c128)
	_, _ = d.TanBatchC128(c128)
	_ = d.TanBatchToC128(c128, c128)
	_, _ = d.SinhBatchC128(c128)
	_ = d.SinhBatchToC128(c128, c128)
	_, _ = d.CoshBatchC128(c128)
	_ = d.CoshBatchToC128(c128, c128)
	_, _ = d.TanhBatchC128(c128)
	_ = d.TanhBatchToC128(c128, c128)
	_, _ = d.SqrtBatchC128(c128)
	_ = d.SqrtBatchToC128(c128, c128)
	_, _ = d.EmlBatchC128(c128, c128)
	_ = d.EmlBatchToC128(c128, c128, c128)
	_, _ = d.DotBatchC128(c128, c128)

	// Pinned memory allocation & free
	p32, err := AllocatePinnedF32(10)
	if err != nil || len(p32) != 10 {
		t.Errorf("AllocatePinnedF32 failed: %v", err)
	}
	FreePinnedF32(p32)
	FreePinnedF32(nil)

	pc64, err := AllocatePinnedC64(10)
	if err != nil || len(pc64) != 10 {
		t.Errorf("AllocatePinnedC64 failed: %v", err)
	}
	FreePinnedC64(pc64)
	FreePinnedC64(nil)

	pc128, err := AllocatePinnedC128(10)
	if err != nil || len(pc128) != 10 {
		t.Errorf("AllocatePinnedC128 failed: %v", err)
	}
	FreePinnedC128(pc128)
	FreePinnedC128(nil)
}
