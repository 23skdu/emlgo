package eml

import (
	"math"
	"testing"
)

func TestExpSIMDTo(t *testing.T) {
	x := []float64{0, 1, 2, -1, 0.5, 10, -10}
	result := make([]float64, len(x))
	ExpSIMDTo(x, result)

	for i := range x {
		expected := math.Exp(x[i])
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("ExpSIMDTo[%d]: got %v, want %v", i, result[i], expected)
		}
	}
}

func TestLogSIMDTo(t *testing.T) {
	x := []float64{1, math.E, 10, 0.5, 100}
	result := make([]float64, len(x))
	LogSIMDTo(x, result)

	for i := range x {
		expected := math.Log(x[i])
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("LogSIMDTo[%d]: got %v, want %v", i, result[i], expected)
		}
	}
}

func TestSinSIMDTo(t *testing.T) {
	x := []float64{0, math.Pi / 4, math.Pi / 2, math.Pi, 3 * math.Pi / 2}
	result := make([]float64, len(x))
	SinSIMDTo(x, result)

	for i := range x {
		expected := math.Sin(x[i])
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("SinSIMDTo[%d]: got %v, want %v", i, result[i], expected)
		}
	}
}

func TestCosSIMDTo(t *testing.T) {
	x := []float64{0, math.Pi / 4, math.Pi / 2, math.Pi, 3 * math.Pi / 2}
	result := make([]float64, len(x))
	CosSIMDTo(x, result)

	for i := range x {
		expected := math.Cos(x[i])
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("CosSIMDTo[%d]: got %v, want %v", i, result[i], expected)
		}
	}
}

func TestSinCosSIMDTo(t *testing.T) {
	x := []float64{0, math.Pi / 4, math.Pi / 2, math.Pi}
	sin := make([]float64, len(x))
	cos := make([]float64, len(x))
	SinCosSIMDTo(x, sin, cos)

	for i := range x {
		expectedSin := math.Sin(x[i])
		expectedCos := math.Cos(x[i])
		if !approximatelyEqual(sin[i], expectedSin, 1e-10) {
			t.Errorf("SinCosSIMDTo sin[%d]: got %v, want %v", i, sin[i], expectedSin)
		}
		if !approximatelyEqual(cos[i], expectedCos, 1e-10) {
			t.Errorf("SinCosSIMDTo cos[%d]: got %v, want %v", i, cos[i], expectedCos)
		}
	}
}

func TestTanSIMDTo(t *testing.T) {
	x := []float64{0, math.Pi / 4, math.Pi / 3}
	result := make([]float64, len(x))
	TanSIMDTo(x, result)

	for i := range x {
		expected := math.Tan(x[i])
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("TanSIMDTo[%d]: got %v, want %v", i, result[i], expected)
		}
	}
}

func TestSqrtSIMDTo(t *testing.T) {
	x := []float64{0, 1, 4, 9, 16, 0.25}
	result := make([]float64, len(x))
	SqrtSIMDTo(x, result)

	for i := range x {
		expected := math.Sqrt(x[i])
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("SqrtSIMDTo[%d]: got %v, want %v", i, result[i], expected)
		}
	}
}

func TestExpMulBatch(t *testing.T) {
	a := []float64{0, 1, 2}
	b := []float64{1, 2, 3}
	result := ExpMulBatch(a, b)

	expected := []float64{math.Exp(0) * 1, math.Exp(1) * 2, math.Exp(2) * 3}
	for i := range result {
		if !approximatelyEqual(result[i], expected[i], 1e-10) {
			t.Errorf("ExpMulBatch[%d]: got %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestExpAddBatch(t *testing.T) {
	a := []float64{0, 1, 2}
	b := []float64{1, 2, 3}
	result := ExpAddBatch(a, b)

	expected := []float64{math.Exp(0) + 1, math.Exp(1) + 2, math.Exp(2) + 3}
	for i := range result {
		if !approximatelyEqual(result[i], expected[i], 1e-10) {
			t.Errorf("ExpAddBatch[%d]: got %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestLogDivBatch(t *testing.T) {
	a := []float64{1, math.E, 10}
	b := []float64{1, 2, 2}
	result := LogDivBatch(a, b)

	expected := []float64{math.Log(1) / 1, math.Log(math.E) / 2, math.Log(10) / 2}
	for i := range result {
		if !approximatelyEqual(result[i], expected[i], 1e-10) {
			t.Errorf("LogDivBatch[%d]: got %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestLogSubBatch(t *testing.T) {
	a := []float64{1, math.E, 10}
	b := []float64{1, 1, 2}
	result := LogSubBatch(a, b)

	expected := []float64{math.Log(1) - 1, math.Log(math.E) - 1, math.Log(10) - 2}
	for i := range result {
		if !approximatelyEqual(result[i], expected[i], 1e-10) {
			t.Errorf("LogSubBatch[%d]: got %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestBranchless(t *testing.T) {
	tests := []struct {
		fn    func(float64) float64
		name string
		x    float64
		want float64
	}{
		{AbsBranchless, "AbsBranchless", -5, 5},
		{AbsBranchless, "AbsBranchless", 5, 5},
		{AbsBranchless, "AbsBranchless", 0, 0},
		{AbsBranchless, "AbsBranchless", -0.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.x)
			if !approximatelyEqual(got, tt.want, 1e-10) {
				t.Errorf("%s(%v): got %v, want %v", tt.name, tt.x, got, tt.want)
			}
		})
	}
}

func TestMinMaxBranchless(t *testing.T) {
	if got := MinBranchless(3, 5); got != 3 {
		t.Errorf("MinBranchless(3, 5): got %v, want 3", got)
	}
	if got := MinBranchless(5, 3); got != 3 {
		t.Errorf("MinBranchless(5, 3): got %v, want 3", got)
	}
	if got := MaxBranchless(3, 5); got != 5 {
		t.Errorf("MaxBranchless(3, 5): got %v, want 5", got)
	}
	if got := MaxBranchless(5, 3); got != 5 {
		t.Errorf("MaxBranchless(5, 3): got %v, want 5", got)
	}
}

func TestSelectBranchless(t *testing.T) {
	if got := SelectBranchless(true, 1, 0); got != 1 {
		t.Errorf("SelectBranchless(true, 1, 0): got %v, want 1", got)
	}
	if got := SelectBranchless(false, 1, 0); got != 0 {
		t.Errorf("SelectBranchless(false, 1, 0): got %v, want 0", got)
	}
}

func TestSelectNaNBranchless(t *testing.T) {
	if got := SelectNaNBranchless(true, math.NaN(), 1); got != 1 {
		t.Errorf("SelectNaNBranchless(true, NaN, 1): got %v, want 1", got)
	}
	if got := SelectNaNBranchless(false, 1, 0); got != 1 {
		t.Errorf("SelectNaNBranchless(false, 1, 0): got %v, want 1", got)
	}
}

func TestGetParallelChunkSize(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{100, 100},
		{256, 256},
		{1000, 256},
		{10000, 256},
	}

	for _, tt := range tests {
		got := GetParallelChunkSize(tt.n)
		if got > tt.n {
			t.Errorf("GetParallelChunkSize(%d): got %d, want <= %d", tt.n, got, tt.n)
		}
	}
}

func TestParallelChunkSizeWithCache(t *testing.T) {
	n := 10000
	chunk := GetParallelChunkSize(n)
	if chunk > L1TileSize/cpuNum {
		t.Errorf("chunk size %d exceeds L1 tile size %d", chunk, L1TileSize/cpuNum)
	}
}

func TestExpSIMDToLarge(t *testing.T) {
	n := 10000
	x := make([]float64, n)
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i) * 0.1
	}
	ExpSIMDTo(x, result)

	for i := 0; i < n; i++ {
		expected := math.Exp(x[i])
		if math.IsInf(expected, 0) && math.IsInf(result[i], 0) {
			continue
		}
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("ExpSIMDToLarge[%d]: got %v, want %v", i, result[i], expected)
			break
		}
	}
}

func TestLogSIMDToLarge(t *testing.T) {
	n := 10000
	x := make([]float64, n)
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i+1) * 0.1
	}
	LogSIMDTo(x, result)

	for i := 0; i < n; i++ {
		expected := math.Log(x[i])
		if math.IsInf(expected, 0) && math.IsInf(result[i], 0) {
			continue
		}
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("LogSIMDToLarge[%d]: got %v, want %v", i, result[i], expected)
			break
		}
	}
}

func TestSinCosSIMDToLarge(t *testing.T) {
	n := 10000
	x := make([]float64, n)
	sin := make([]float64, n)
	cos := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i) * 0.01
	}
	SinCosSIMDTo(x, sin, cos)

	for i := 0; i < n; i++ {
		expectedSin := math.Sin(x[i])
		expectedCos := math.Cos(x[i])
		if !approximatelyEqual(sin[i], expectedSin, 1e-10) {
			t.Errorf("SinCosSIMDToLarge sin[%d]: got %v, want %v", i, sin[i], expectedSin)
			break
		}
		if !approximatelyEqual(cos[i], expectedCos, 1e-10) {
			t.Errorf("SinCosSIMDToLarge cos[%d]: got %v, want %v", i, cos[i], expectedCos)
			break
		}
	}
}

func TestFusedExpMulBatchLarge(t *testing.T) {
	n := 10000
	a := make([]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = float64(i) * 0.1
		b[i] = 2.0
	}
	result := ExpMulBatch(a, b)

	for i := 0; i < n; i++ {
		expected := math.Exp(a[i]) * b[i]
		if math.IsInf(expected, 0) && math.IsInf(result[i], 0) {
			continue
		}
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("ExpMulBatchLarge[%d]: got %v, want %v", i, result[i], expected)
			break
		}
	}
}

func TestFusedLogDivBatchLarge(t *testing.T) {
	n := 10000
	a := make([]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = float64(i+1) * 0.1
		b[i] = 2.0
	}
	result := LogDivBatch(a, b)

	for i := 0; i < n; i++ {
		expected := math.Log(a[i]) / b[i]
		if !approximatelyEqual(result[i], expected, 1e-10) {
			t.Errorf("LogDivBatchLarge[%d]: got %v, want %v", i, result[i], expected)
			break
		}
	}
}

func approximatelyEqual(a, b, epsilon float64) bool {
	delta := math.Abs(a - b)
	if delta < epsilon {
		return true
	}
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale > 0 {
		return delta/scale < epsilon
	}
	return delta < epsilon
}

func TestExpLogAccuracy(t *testing.T) {
	tests := []struct {
		name  string
		fn   func(float64) float64
		x    float64
		relTol float64
	}{
		{"Exp(0)", math.Exp, 0, 1e-10},
		{"Exp(1)", math.Exp, 1, 1e-10},
		{"Exp(-1)", math.Exp, -1, 1e-10},
		{"Log(1)", math.Log, 1, 1e-10},
		{"Log(2)", math.Log, 2, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.x)
			expected := tt.fn(tt.x)
			if result == 0 {
				return
			}
			relErr := math.Abs(result-expected) / math.Abs(expected)
			if relErr > tt.relTol {
				t.Errorf("relative error %e > tolerance %e", relErr, tt.relTol)
			}
		})
	}
}