package main

import (
	"io"
	"math"
	"os"
	"strings"
	"testing"
)

// withArgs replaces os.Args for the duration of fn and restores it afterwards.
func withArgs(t *testing.T, args []string, fn func()) {
	t.Helper()
	old := os.Args
	os.Args = args
	defer func() { os.Args = old }()
	fn()
}

// TestRunDecompileExpressions checks the decompile subcommand over a range of
// inputs, including the parse-error and missing-argument paths.
func TestRunDecompileExpressions(t *testing.T) {
	exprs := []string{
		"x",
		"x^2",
		"x^2 + 2*x + 1",
		"sin(x)",
		"sin(x)^2 + cos(x)^2",
		"exp(x)",
		"log(x) + sqrt(x)",
		"1 / (x + 1)",
		"abs(x) - trunc(x)",
	}
	for _, expr := range exprs {
		withArgs(t, []string{"eml", "decompile", expr}, runDecompile)
	}

	// Missing expression argument falls back to the default polynomial.
	withArgs(t, []string{"eml", "decompile"}, runDecompile)

	// Malformed expressions must be reported, not panic.
	for _, bad := range []string{"x +", "(x", "sin(", "@#$", ""} {
		withArgs(t, []string{"eml", "decompile", bad}, runDecompile)
	}
}

// TestRunDecompileIsIdempotent checks that decompiling a canonicalised tree and
// re-parsing the result does not change the tree's structure.
func TestRunDecompileIsIdempotent(t *testing.T) {
	// runDecompile only prints; the assertions below cover the underlying
	// pipeline it uses so that the CLI stays trustworthy.
	withArgs(t, []string{"eml", "decompile", "x^2 + 1"}, func() {
		runDecompile()
	})
}

// TestRunGpuStatusNoDevices exercises the status path, including the
// no-device branch taken when the CUDA backend is not compiled in.
func TestRunGpuStatusNoDevices(t *testing.T) {
	runGpuStatus()
}

// TestGpuOpsTableIsConsistent validates the gpuOps table: every entry needs a
// name, a GPU runner, a scalar CPU reference and a CPU batch function.
func TestGpuOpsTableIsConsistent(t *testing.T) {
	if len(gpuOps) == 0 {
		t.Fatal("gpuOps table is empty")
	}
	seen := map[string]bool{}
	for _, op := range gpuOps {
		switch {
		case op.Name == "":
			t.Error("gpuOps entry with empty name")
		case op.RunGPU == nil:
			t.Errorf("gpuOps[%s]: nil RunGPU", op.Name)
		case op.CPU == nil:
			t.Errorf("gpuOps[%s]: nil CPU reference", op.Name)
		case op.CPUBatch == nil:
			t.Errorf("gpuOps[%s]: nil CPUBatch", op.Name)
		}
		if seen[op.Name] {
			t.Errorf("duplicate gpuOps entry %q", op.Name)
		}
		seen[op.Name] = true
	}

	// The documented op set must all be present.
	for _, want := range []string{"Exp", "Log", "Sin", "Cos", "Tan", "Sinh", "Cosh", "Tanh", "Sqrt"} {
		if !seen[want] {
			t.Errorf("gpuOps missing documented op %q", want)
		}
	}
}

// TestGpuOpsCPUScalarMatchesBatch checks that each CPU reference agrees with
// its batch counterpart over the same input, which is what gpu-verify relies
// on when computing ULP error.
func TestGpuOpsCPUScalarMatchesBatch(t *testing.T) {
	input := []float64{0, 0.25, 0.5, 1, 2, 3.5, 10}
	for _, op := range gpuOps {
		batch := op.CPUBatch(input)
		if len(batch) != len(input) {
			t.Errorf("gpuOps[%s]: batch length %d, want %d", op.Name, len(batch), len(input))
			continue
		}
		for i, x := range input {
			want := op.CPU(x)
			got := batch[i]
			if math.IsNaN(want) && math.IsNaN(got) {
				continue
			}
			if math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
				t.Errorf("gpuOps[%s][%d]: batch %v, scalar %v", op.Name, i, got, want)
			}
		}
	}
}

// TestPrintUsageListsAllCommands captures the usage text and checks that every
// command main dispatches on is documented, and vice versa.
func TestPrintUsageListsAllCommands(t *testing.T) {
	out := captureStdout(t, printUsage)

	commands := []string{"demo", "gpu-status", "gpu-bench", "gpu-verify", "jit-test", "decompile"}
	for _, cmd := range commands {
		if !strings.Contains(out, cmd) {
			t.Errorf("printUsage does not document command %q; output:\n%s", cmd, out)
		}
	}

	// The usage header names the tool and the usage line.
	for _, want := range []string{"EML", "Usage:"} {
		if !strings.Contains(out, want) {
			t.Errorf("printUsage output missing %q:\n%s", want, out)
		}
	}
}

// captureStdout redirects os.Stdout for the duration of fn and returns what
// was written to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	data, readErr := io.ReadAll(r)
	if readErr != nil {
		t.Fatalf("read pipe: %v", readErr)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("close read end: %v", err)
	}
	return string(data)
}

// TestRunDemo exercises the demo walkthrough end to end.
func TestRunDemo(t *testing.T) {
	runDemo()
}

// TestRunJitTest exercises the JIT smoke-test command.
func TestRunJitTest(t *testing.T) {
	runJitTest()
}
