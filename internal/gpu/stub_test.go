//go:build !cuda && !(darwin && arm64)

package gpu

import (
	"strings"
	"testing"
)

func TestStubGetDevices(t *testing.T) {
	devices, err := GetDevices()
	if err != nil {
		t.Fatalf("GetDevices() should not error on linux/windows: %v", err)
	}
	if len(devices) > 0 {
		t.Errorf("expected no devices in stub mode, got %d", len(devices))
	}
}

func TestStubGetDevicesUnsupportedOS(t *testing.T) {
	saved := cudaSupportedOS
	cudaSupportedOS = func() bool { return false }
	defer func() { cudaSupportedOS = saved }()

	_, err := GetDevices()
	if err == nil {
		t.Error("GetDevices() should error on unsupported OS")
	}
	if !strings.Contains(err.Error(), "Linux and Windows") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestStubStatus(t *testing.T) {
	s := Status()
	if s == "" {
		t.Fatal("Status() returned empty string")
	}
	if !strings.Contains(s, "No GPU devices found") {
		t.Errorf("unexpected status: %s", s)
	}
}

func TestStubStatusError(t *testing.T) {
	saved := cudaSupportedOS
	cudaSupportedOS = func() bool { return false }
	defer func() { cudaSupportedOS = saved }()

	s := Status()
	if !strings.Contains(s, "GPU Error") {
		t.Errorf("expected GPU Error in status, got: %s", s)
	}
}

func TestStubStatusWithDevices(t *testing.T) {
	saved := getDevicesFn
	getDevicesFn = func() ([]Device, error) {
		return []Device{{ID: 0, Name: "Fake GPU", MemoryBytes: 1073741824}}, nil
	}
	defer func() { getDevicesFn = saved }()

	s := Status()
	if !strings.Contains(s, "Found 1 GPU device") {
		t.Errorf("expected 'Found 1 GPU device', got: %s", s)
	}
}

func TestStubInitShutdown(t *testing.T) {
	if err := Init(); err != nil {
		t.Errorf("Init() should not error in stub mode: %v", err)
	}
	Shutdown()
}

func TestStubBatchOps(t *testing.T) {
	d := &Device{ID: 0, Name: "stub"}

	ops := []struct {
		name string
		fn   func() error
	}{
		{"ExpBatch", func() error { _, err := d.ExpBatch(nil); return err }},
		{"LogBatch", func() error { _, err := d.LogBatch(nil); return err }},
		{"SinBatch", func() error { _, err := d.SinBatch(nil); return err }},
		{"CosBatch", func() error { _, err := d.CosBatch(nil); return err }},
		{"TanBatch", func() error { _, err := d.TanBatch(nil); return err }},
		{"SinhBatch", func() error { _, err := d.SinhBatch(nil); return err }},
		{"CoshBatch", func() error { _, err := d.CoshBatch(nil); return err }},
		{"TanhBatch", func() error { _, err := d.TanhBatch(nil); return err }},
		{"SqrtBatch", func() error { _, err := d.SqrtBatch(nil); return err }},
	}

	for _, op := range ops {
		err := op.fn()
		if err == nil {
			t.Errorf("%s: expected error in stub mode, got nil", op.name)
		}
	}
}

func TestStubEmlBatch(t *testing.T) {
	d := &Device{ID: 0}
	_, err := d.EmlBatch(nil, nil)
	if err == nil {
		t.Error("EmlBatch: expected error in stub mode")
	}
}

func TestStubStream(t *testing.T) {
	s, err := NewStream()
	if err == nil {
		t.Error("NewStream: expected error in stub mode")
	}
	if s != nil {
		t.Error("NewStream: expected nil stream in stub mode")
	}
}

func TestStubPinnedMemory(t *testing.T) {
	p, err := AllocatePinned(100)
	if err != nil {
		t.Errorf("AllocatePinned: unexpected error: %v", err)
	}
	if len(p) != 100 {
		t.Errorf("AllocatePinned: expected len 100, got %d", len(p))
	}
	FreePinned(p)
}

func TestStubAllocatePinnedZero(t *testing.T) {
	p, err := AllocatePinned(0)
	if err != nil {
		t.Errorf("AllocatePinned(0): unexpected error: %v", err)
	}
	if len(p) != 0 {
		t.Errorf("AllocatePinned(0): expected len 0, got %d", len(p))
	}
	FreePinned(p)
}

func TestStubFreePinnedEmpty(t *testing.T) {
	FreePinned(nil)
	FreePinned([]float64{})
}

func TestStubBatchOpsWithData(t *testing.T) {
	d := &Device{ID: 0}
	data := []float64{1.0, 2.0, 3.0}

	ops := []struct {
		name string
		fn   func() error
	}{
		{"ExpBatch", func() error { _, err := d.ExpBatch(data); return err }},
		{"LogBatch", func() error { _, err := d.LogBatch(data); return err }},
		{"SinBatch", func() error { _, err := d.SinBatch(data); return err }},
		{"CosBatch", func() error { _, err := d.CosBatch(data); return err }},
		{"TanBatch", func() error { _, err := d.TanBatch(data); return err }},
		{"SinhBatch", func() error { _, err := d.SinhBatch(data); return err }},
		{"CoshBatch", func() error { _, err := d.CoshBatch(data); return err }},
		{"TanhBatch", func() error { _, err := d.TanhBatch(data); return err }},
		{"SqrtBatch", func() error { _, err := d.SqrtBatch(data); return err }},
	}

	for _, op := range ops {
		err := op.fn()
		if err == nil {
			t.Errorf("%s: expected error in stub mode", op.name)
		}
	}
}

func TestStubEmlBatchWithData(t *testing.T) {
	d := &Device{ID: 0}
	x := []float64{1.0, 2.0}
	y := []float64{3.0, 4.0}
	_, err := d.EmlBatch(x, y)
	if err == nil {
		t.Error("EmlBatch: expected error in stub mode")
	}
}
