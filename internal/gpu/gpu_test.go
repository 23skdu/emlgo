package gpu

import (
	"testing"
)

func TestGridSize(t *testing.T) {
	tests := []struct {
		n, blockSize, want int
	}{
		{0, 256, 0},
		{1, 256, 1},
		{256, 256, 1},
		{257, 256, 2},
		{512, 256, 2},
		{1024, 256, 4},
		{100, 0, 1}, // blockSize=0 -> DefaultBlockSize=256
		{300, 32, 10},
		{1, 1024, 1},
		{100, -1, 1},
	}

	for _, tt := range tests {
		got := GridSize(tt.n, tt.blockSize)
		if got != tt.want {
			t.Errorf("GridSize(%d, %d) = %d, want %d", tt.n, tt.blockSize, got, tt.want)
		}
	}
}

func TestDefaultLaunchConfig(t *testing.T) {
	lc := DefaultLaunchConfig(1024, 256)
	if lc.GridDimX != 4 {
		t.Errorf("GridDimX = %d, want 4", lc.GridDimX)
	}
	if lc.BlockDimX != 256 {
		t.Errorf("BlockDimX = %d, want 256", lc.BlockDimX)
	}
	if lc.BlockDimY != 1 || lc.BlockDimZ != 1 {
		t.Errorf("expected 1D config, got %dx%dx%d", lc.BlockDimX, lc.BlockDimY, lc.BlockDimZ)
	}
}

func TestDefaultLaunchConfigZeroBlockSize(t *testing.T) {
	lc := DefaultLaunchConfig(512, 0)
	if lc.BlockDimX != DefaultBlockSize {
		t.Errorf("BlockDimX = %d, want %d", lc.BlockDimX, DefaultBlockSize)
	}
}

func TestDefaultLaunchConfigLargeBlockSize(t *testing.T) {
	lc := DefaultLaunchConfig(512, 2048)
	if lc.BlockDimX != DefaultBlockSize {
		t.Errorf("BlockDimX should be clamped, got %d", lc.BlockDimX)
	}
}

func TestDefaultLaunchConfigNegativeBlockSize(t *testing.T) {
	lc := DefaultLaunchConfig(512, -64)
	if lc.BlockDimX != DefaultBlockSize {
		t.Errorf("BlockDimX should be default, got %d", lc.BlockDimX)
	}
}

func TestDefaultLaunchConfigZeroN(t *testing.T) {
	lc := DefaultLaunchConfig(0, 256)
	if lc.GridDimX != 0 {
		t.Errorf("GridDimX = %d, want 0 for n=0", lc.GridDimX)
	}
}

func TestMatMulLaunchConfig(t *testing.T) {
	lc := MatMulLaunchConfig(128, 256, 16)
	if lc.GridDimX != 16 { // 256/16 = 16
		t.Errorf("GridDimX = %d, want 16", lc.GridDimX)
	}
	if lc.GridDimY != 8 { // 128/16 = 8
		t.Errorf("GridDimY = %d, want 8", lc.GridDimY)
	}
	if lc.BlockDimX != 16 || lc.BlockDimY != 16 {
		t.Errorf("expected 16x16 block, got %dx%d", lc.BlockDimX, lc.BlockDimY)
	}
}

func TestMatMulLaunchConfigDefaultTile(t *testing.T) {
	lc := MatMulLaunchConfig(32, 32, 0)
	if lc.BlockDimX != 16 {
		t.Errorf("expected default 16 tile, got %d", lc.BlockDimX)
	}
}

func TestMatMulLaunchConfigNegativeTile(t *testing.T) {
	lc := MatMulLaunchConfig(32, 32, -8)
	if lc.BlockDimX != 16 {
		t.Errorf("expected default 16 tile, got %d", lc.BlockDimX)
	}
}

func TestMatMulLaunchConfigEdgeCases(t *testing.T) {
	lc := MatMulLaunchConfig(1, 1, 16)
	if lc.GridDimX != 1 || lc.GridDimY != 1 {
		t.Errorf("expected 1x1 grid, got %dx%d", lc.GridDimX, lc.GridDimY)
	}
}

func TestMatMulLaunchConfigNonSquare(t *testing.T) {
	lc := MatMulLaunchConfig(64, 128, 16)
	if lc.GridDimX != 8 { // 128/16 = 8
		t.Errorf("GridDimX = %d, want 8", lc.GridDimX)
	}
	if lc.GridDimY != 4 { // 64/16 = 4
		t.Errorf("GridDimY = %d, want 4", lc.GridDimY)
	}
}

func TestMatMulLaunchConfigLarge(t *testing.T) {
	lc := MatMulLaunchConfig(4096, 4096, 16)
	if lc.GridDimX != 256 {
		t.Errorf("GridDimX = %d, want 256", lc.GridDimX)
	}
	if lc.GridDimY != 256 {
		t.Errorf("GridDimY = %d, want 256", lc.GridDimY)
	}
}

func TestDefaultLaunchConfigTotalThreads(t *testing.T) {
	totalThreads := func(lc LaunchConfig) int {
		return lc.GridDimX * lc.GridDimY * lc.GridDimZ *
			lc.BlockDimX * lc.BlockDimY * lc.BlockDimZ
	}

	tests := []struct {
		n, blockSize int
	}{
		{1, 256},
		{100, 64},
		{1000, 128},
		{100000, 256},
		{1000000, 512},
	}

	for _, tt := range tests {
		lc := DefaultLaunchConfig(tt.n, tt.blockSize)
		threads := totalThreads(lc)
		if threads < tt.n {
			t.Errorf("total threads %d < n %d", threads, tt.n)
		}
	}
}

func TestDeviceStruct(t *testing.T) {
	d := Device{
		ID:                 0,
		Name:               "Test GPU",
		ComputeMajor:       8,
		ComputeMinor:       0,
		MemoryBytes:        8589934592, // 8 GB
		MaxThreadsPerBlock: 1024,
		WarpSize:           32,
		ClockRateKHz:       1500000,
	}

	if d.ID != 0 {
		t.Errorf("ID = %d, want 0", d.ID)
	}
	if d.Name != "Test GPU" {
		t.Errorf("Name = %s, want Test GPU", d.Name)
	}
	if d.MemoryBytes != 8589934592 {
		t.Errorf("MemoryBytes = %d, want 8589934592", d.MemoryBytes)
	}
}

func TestStreamStruct(t *testing.T) {
	s := &Stream{handle: 0}
	if err := s.Sync(); err != nil {
		t.Errorf("Sync on zero stream should not error: %v", err)
	}
	if err := s.Destroy(); err != nil {
		t.Errorf("Destroy on zero stream should not error: %v", err)
	}
}

func TestNilStream(t *testing.T) {
	var s *Stream
	if err := s.Sync(); err != nil {
		t.Errorf("Sync on nil stream should not error: %v", err)
	}
	if err := s.Destroy(); err != nil {
		t.Errorf("Destroy on nil stream should not error: %v", err)
	}
}

func TestDefaultBlockSizeConstant(t *testing.T) {
	if DefaultBlockSize != 256 {
		t.Errorf("DefaultBlockSize = %d, want 256", DefaultBlockSize)
	}
}

func BenchmarkGridSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GridSize(1000000, 256)
	}
}

func BenchmarkDefaultLaunchConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultLaunchConfig(1000000, 256)
	}
}

func BenchmarkMatMulLaunchConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MatMulLaunchConfig(4096, 4096, 16)
	}
}
