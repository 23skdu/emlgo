package bytecode

import (
	"math"
	"testing"
)

// TestEvalBatchColumnarScratchMatchesAllocating checks the scratch form is
// numerically identical to the allocating one, including across chunk
// boundaries and for a scratch that has to grow.
func TestEvalBatchColumnarScratchMatchesAllocating(t *testing.T) {
	programs := []string{
		"x + 1",
		"2*x*x + 3*x + 1",
		"sqrt(x) + exp(x) - log(x)",
		"sin(x)*cos(x) + abs(x) - round(x)",
		"x / (x + 2) - cbrt(x)",
		"((((x + 1) * x + 2) * x + 3) * x + 4)",
	}

	// Sizes that straddle the default chunk size of 1024.
	sizes := []int{1, 7, 1023, 1024, 1025, 4096}

	for _, expr := range programs {
		p, err := CompileExpr(expr)
		if err != nil {
			t.Fatalf("CompileExpr(%q): %v", expr, err)
		}
		for _, n := range sizes {
			data := make([][]float64, 1)
			data[0] = make([]float64, n)
			for i := range data[0] {
				data[0][i] = 0.5 + float64(i%53)*0.05
			}

			want := make([]float64, n)
			got := make([]float64, n)

			p.EvalBatchColumnar(data, want)
			s := NewBatchScratch(p, 0)
			p.EvalBatchColumnarScratch(data, got, &s)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("%q n=%d: [%d] scratch %v != allocating %v",
						expr, n, i, got[i], want[i])
				}
			}
		}
	}
}

// TestEvalBatchColumnarScratchIsReusedAcrossPrograms checks the scratch
// resizes correctly when a deeper program or a larger chunk needs more room.
func TestEvalBatchColumnarScratchIsReusedAcrossPrograms(t *testing.T) {
	shallow, err := CompileExpr("x + 1")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}
	deep, err := CompileExpr("(((((((x + 1) * x) + 1) * x) + 1) * x) + 1) * x")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}

	const n = 2048
	data := make([][]float64, 1)
	data[0] = make([]float64, n)
	for i := range data[0] {
		data[0][i] = 1 + float64(i%17)*0.1
	}

	// One scratch, alternated between a shallow and a deep program.
	s := NewBatchScratch(shallow, 0)

	for round := 0; round < 3; round++ {
		for _, tc := range []struct {
			name string
			p    *Program
		}{{"shallow", shallow}, {"deep", deep}, {"shallow again", shallow}, {"deep again", deep}} {
			got := make([]float64, n)
			want := make([]float64, n)
			tc.p.EvalBatchColumnar(data, want)
			tc.p.EvalBatchColumnarScratch(data, got, &s)
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("round %d %s: [%d] %v != %v", round, tc.name, i, got[i], want[i])
				}
			}
		}
	}
}

// TestEvalBatchColumnarScratchChunkSizes checks a caller-supplied chunk size is
// honoured and produces the same results.
func TestEvalBatchColumnarScratchChunkSizes(t *testing.T) {
	p, err := CompileExpr("x*x + 1")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}

	const n = 500
	data := make([][]float64, 1)
	data[0] = make([]float64, n)
	for i := range data[0] {
		data[0][i] = 0.25 + float64(i)*0.01
	}

	want := make([]float64, n)
	p.EvalBatchColumnar(data, want)

	for _, chunk := range []int{0, 1, 3, 64, 1024} {
		got := make([]float64, n)
		s := NewBatchScratch(p, chunk)
		p.EvalBatchColumnarScratch(data, got, &s)
		for i := range want {
			if math.Abs(got[i]-want[i]) > 1e-15*math.Max(1, math.Abs(want[i])) {
				t.Fatalf("chunk=%d: [%d] %v != %v", chunk, i, got[i], want[i])
			}
		}
	}
}

// TestEvalBatchColumnarScratchHandlesEmptyAndNil checks the degenerate inputs.
func TestEvalBatchColumnarScratchHandlesEmptyAndNil(t *testing.T) {
	p, err := CompileExpr("x + 1")
	if err != nil {
		t.Fatalf("CompileExpr: %v", err)
	}
	s := NewBatchScratch(p, 0)

	p.EvalBatchColumnarScratch(nil, nil, &s)
	p.EvalBatchColumnarScratch([][]float64{make([]float64, 4)}, nil, &s)

	// An empty Program must be a no-op rather than a panic.
	var empty *Program
	s2 := NewBatchScratch(empty, 0)
	empty.EvalBatchColumnarScratch([][]float64{make([]float64, 4)}, make([]float64, 4), &s2)
	NewProgram().EvalBatchColumnarScratch([][]float64{make([]float64, 4)}, make([]float64, 4), &s2)
}
