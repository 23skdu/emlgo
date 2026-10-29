package test

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCheckBenchRegressionScript tests scripts/check_bench_regression.py CLI behavior.
func TestCheckBenchRegressionScript(t *testing.T) {
	root := getRepoRoot(t)
	scriptPath := filepath.Join(root, "scripts", "check_bench_regression.py")

	pythonPath, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not found; skipping python script tests")
	}

	t.Run("PythonSyntaxCheck", func(t *testing.T) {
		cmd := exec.Command(pythonPath, "-m", "py_compile", scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Python syntax check failed: %v\nOutput: %s", err, string(out))
		}
	})

	t.Run("HelpFlag", func(t *testing.T) {
		cmd := exec.Command(pythonPath, scriptPath, "-h")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Help flag failed: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "Fail only when repeated") {
			t.Errorf("Unexpected help text: %s", string(out))
		}
	})

	t.Run("InvalidMaxRegressionFlag", func(t *testing.T) {
		dummy1 := filepath.Join(t.TempDir(), "base.txt")
		dummy2 := filepath.Join(t.TempDir(), "curr.txt")
		_ = os.WriteFile(dummy1, []byte(""), 0644)
		_ = os.WriteFile(dummy2, []byte(""), 0644)

		cmd := exec.Command(pythonPath, scriptPath, "--max-regression", "1.5", dummy1, dummy2)
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected error for --max-regression 1.5, got success:\n%s", string(out))
		}
	})

	t.Run("InsufficientSamplesExit2", func(t *testing.T) {
		tmpDir := t.TempDir()
		baseFile := filepath.Join(tmpDir, "base.txt")
		currFile := filepath.Join(tmpDir, "curr.txt")

		// Only 3 samples (< 5 required)
		benchData := `
BenchmarkExp-8   1000000   100.0 ns/op   0 B/op   0 allocs/op
BenchmarkExp-8   1000000   101.0 ns/op   0 B/op   0 allocs/op
BenchmarkExp-8   1000000   102.0 ns/op   0 B/op   0 allocs/op
`
		if err := os.WriteFile(baseFile, []byte(benchData), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(currFile, []byte(benchData), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := exec.Command(pythonPath, scriptPath, baseFile, currFile)
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected non-zero exit for insufficient samples, got success:\n%s", string(out))
		}
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.ExitCode() != 2 {
			t.Fatalf("Expected exit code 2, got: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "SKIP BenchmarkExp") {
			t.Errorf("Expected output to report SKIP BenchmarkExp, got:\n%s", string(out))
		}
	})

	t.Run("PassingRegressionExit0", func(t *testing.T) {
		tmpDir := t.TempDir()
		baseFile := filepath.Join(tmpDir, "base.txt")
		currFile := filepath.Join(tmpDir, "curr.txt")

		var baseSb, currSb strings.Builder
		for i := 0; i < 6; i++ {
			baseSb.WriteString(fmt.Sprintf("BenchmarkExp-8   1000000   %d.0 ns/op\n", 100+i))
			currSb.WriteString(fmt.Sprintf("BenchmarkExp-8   1000000   %d.0 ns/op\n", 101+i))
		}

		if err := os.WriteFile(baseFile, []byte(baseSb.String()), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(currFile, []byte(currSb.String()), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := exec.Command(pythonPath, scriptPath, baseFile, currFile)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Expected exit code 0 for passing regression, got: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "PASS BenchmarkExp") {
			t.Errorf("Expected PASS BenchmarkExp in output, got:\n%s", string(out))
		}
	})

	t.Run("FailingRegressionExit1", func(t *testing.T) {
		tmpDir := t.TempDir()
		baseFile := filepath.Join(tmpDir, "base.txt")
		currFile := filepath.Join(tmpDir, "curr.txt")

		var baseSb, currSb strings.Builder
		for i := 0; i < 6; i++ {
			baseSb.WriteString(fmt.Sprintf("BenchmarkExp-8   1000000   %d.0 ns/op\n", 100+i))
			currSb.WriteString(fmt.Sprintf("BenchmarkExp-8   1000000   %d.0 ns/op\n", 150+i)) // 50% slowdown
		}

		if err := os.WriteFile(baseFile, []byte(baseSb.String()), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(currFile, []byte(currSb.String()), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := exec.Command(pythonPath, scriptPath, "--max-regression", "0.10", baseFile, currFile)
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected exit code 1 for regression slowdown, got 0:\n%s", string(out))
		}
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.ExitCode() != 1 {
			t.Fatalf("Expected exit code 1, got: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "FAIL BenchmarkExp") {
			t.Errorf("Expected FAIL BenchmarkExp in output, got:\n%s", string(out))
		}
	})
}

// TestRunBenchmarkScript tests scripts/run_benchmark.sh.
func TestRunBenchmarkScript(t *testing.T) {
	root := getRepoRoot(t)
	scriptPath := filepath.Join(root, "scripts", "run_benchmark.sh")

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found; skipping bash script tests")
	}

	t.Run("BashSyntaxCheck", func(t *testing.T) {
		cmd := exec.Command(bashPath, "-n", scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Syntax check for run_benchmark.sh failed: %v\nOutput: %s", err, string(out))
		}
	})

	t.Run("HelpOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "-h")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("run_benchmark.sh -h failed: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "Usage:") || !strings.Contains(string(out), "Run feature parity comparison") {
			t.Errorf("Expected usage instructions in output, got:\n%s", string(out))
		}
	})

	t.Run("InvalidOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "-z")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected non-zero exit for invalid option -z, got success:\n%s", string(out))
		}
		if !strings.Contains(string(out), "illegal option") && !strings.Contains(string(out), "Invalid option") {
			t.Errorf("Expected invalid option error in output, got:\n%s", string(out))
		}
	})
}

// TestRunValidationScript tests scripts/run_validation.sh.
func TestRunValidationScript(t *testing.T) {
	root := getRepoRoot(t)
	scriptPath := filepath.Join(root, "scripts", "run_validation.sh")

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found; skipping bash script tests")
	}

	t.Run("BashSyntaxCheck", func(t *testing.T) {
		cmd := exec.Command(bashPath, "-n", scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Syntax check for run_validation.sh failed: %v\nOutput: %s", err, string(out))
		}
	})

	t.Run("HelpOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "-h")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("run_validation.sh -h failed: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "Usage:") || !strings.Contains(string(out), "Filter by type") {
			t.Errorf("Expected usage instructions in output, got:\n%s", string(out))
		}
	})

	t.Run("InvalidOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "-z")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected non-zero exit for invalid option -z, got success:\n%s", string(out))
		}
		if !strings.Contains(string(out), "illegal option") && !strings.Contains(string(out), "Invalid option") {
			t.Errorf("Expected invalid option error in output, got:\n%s", string(out))
		}
	})
}

// TestBenchCompareScript tests scripts/bench-compare.sh.
func TestBenchCompareScript(t *testing.T) {
	root := getRepoRoot(t)
	scriptPath := filepath.Join(root, "scripts", "bench-compare.sh")

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found; skipping bash script tests")
	}

	t.Run("BashSyntaxCheck", func(t *testing.T) {
		cmd := exec.Command(bashPath, "-n", scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Syntax check for bench-compare.sh failed: %v\nOutput: %s", err, string(out))
		}
	})

	t.Run("ScriptDirectivesAndStructure", func(t *testing.T) {
		content, err := os.ReadFile(scriptPath)
		if err != nil {
			t.Fatalf("Failed to read bench-compare.sh: %v", err)
		}
		s := string(content)
		expectedStrings := []string{
			"BASELINE=${1:-main}",
			"BENCH_PATTERN=${2:-.}",
			"git stash",
			"git checkout \"$BASELINE\"",
			"benchstat",
		}
		for _, expected := range expectedStrings {
			if !strings.Contains(s, expected) {
				t.Errorf("bench-compare.sh missing expected construct: %q", expected)
			}
		}
	})
}

// TestWasmTestScript tests scripts/wasm_test.sh.
func TestWasmTestScript(t *testing.T) {
	root := getRepoRoot(t)
	scriptPath := filepath.Join(root, "scripts", "wasm_test.sh")

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found; skipping bash script tests")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found; skipping wasm_test.sh tests")
	}

	t.Run("BashSyntaxCheck", func(t *testing.T) {
		cmd := exec.Command(bashPath, "-n", scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Syntax check for wasm_test.sh failed: %v\nOutput: %s", err, string(out))
		}
	})

	t.Run("InvalidMode", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "--invalid-mode")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected non-zero exit for invalid mode, got 0:\n%s", string(out))
		}
		if !strings.Contains(string(out), "Usage: ") {
			t.Errorf("Expected usage output for invalid mode, got:\n%s", string(out))
		}
	})

	t.Run("CleanMode", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "clean")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("wasm_test.sh clean failed: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "Cleaning WASM build artifacts") {
			t.Errorf("Expected clean confirmation, got:\n%s", string(out))
		}
	})
}

// TestGenerateMPFRVectorsSource validates scripts/generate_mpfr_vectors.c.
func TestGenerateMPFRVectorsSource(t *testing.T) {
	root := getRepoRoot(t)
	srcPath := filepath.Join(root, "scripts", "generate_mpfr_vectors.c")

	content, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Failed to read generate_mpfr_vectors.c: %v", err)
	}
	s := string(content)

	t.Run("RequiredStructuresAndPrototypes", func(t *testing.T) {
		required := []string{
			"#include <mpfr.h>",
			"typedef struct",
			"vector64",
			"vector32",
			"exp64_inputs",
			"log64_inputs",
			"exp32_inputs",
			"log32_inputs",
			"mpfr_get_version()",
			"4.2.2",
			"emit64",
			"emit32",
		}
		for _, req := range required {
			if !strings.Contains(s, req) {
				t.Errorf("generate_mpfr_vectors.c missing required element: %q", req)
			}
		}
	})

	t.Run("IEEE754BitConstantsValidation", func(t *testing.T) {
		// Verify standard IEEE 754 bit values used in the C file
		tests64 := []struct {
			name     string
			val      float64
			expected uint64
		}{
			{"0.0", 0.0, 0x0000000000000000},
			{"1.0", 1.0, 0x3ff0000000000000},
			{"2.0", 2.0, 0x4000000000000000},
			{"0.5", 0.5, 0x3fe0000000000000},
			{"-1.0", -1.0, 0xbff0000000000000},
			{"-10.0", -10.0, 0xc024000000000000},
			{"10.0", 10.0, 0x4024000000000000},
		}

		for _, tc := range tests64 {
			bits := math.Float64bits(tc.val)
			if bits != tc.expected {
				t.Errorf("Float64 constant %s bits = 0x%016x, expected 0x%016x", tc.name, bits, tc.expected)
			}
			hexPattern := fmt.Sprintf("%016x", tc.expected)
			if !strings.Contains(s, hexPattern) {
				t.Errorf("generate_mpfr_vectors.c missing expected 64-bit constant %s (0x%s)", tc.name, hexPattern)
			}
		}

		tests32 := []struct {
			name     string
			val      float32
			expected uint32
		}{
			{"0.0f", 0.0, 0x00000000},
			{"1.0f", 1.0, 0x3f800000},
			{"2.0f", 2.0, 0x40000000},
			{"0.5f", 0.5, 0x3f000000},
			{"-1.0f", -1.0, 0xbf800000},
		}

		for _, tc := range tests32 {
			bits := math.Float32bits(tc.val)
			if bits != tc.expected {
				t.Errorf("Float32 constant %s bits = 0x%08x, expected 0x%08x", tc.name, bits, tc.expected)
			}
			hexPattern := fmt.Sprintf("%08x", tc.expected)
			if !strings.Contains(s, hexPattern) {
				t.Errorf("generate_mpfr_vectors.c missing expected 32-bit constant %s (0x%s)", tc.name, hexPattern)
			}
		}
	})
}

// TestOracleAccuracyGateScript tests scripts/oracle_accuracy_gate.sh.
func TestOracleAccuracyGateScript(t *testing.T) {
	root := getRepoRoot(t)
	scriptPath := filepath.Join(root, "scripts", "oracle_accuracy_gate.sh")

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found; skipping bash script tests")
	}

	t.Run("BashSyntaxCheck", func(t *testing.T) {
		cmd := exec.Command(bashPath, "-n", scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Syntax check for oracle_accuracy_gate.sh failed: %v\nOutput: %s", err, string(out))
		}
	})

	t.Run("HelpOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "-h")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("oracle_accuracy_gate.sh -h failed: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "Oracle Accuracy Gate") || !strings.Contains(string(out), "--quick") {
			t.Errorf("Expected usage instructions in output, got:\n%s", string(out))
		}
	})

	t.Run("InvalidOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "--invalid-flag")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected non-zero exit for invalid option --invalid-flag, got success:\n%s", string(out))
		}
		if !strings.Contains(string(out), "Invalid option") {
			t.Errorf("Expected Invalid option error in output, got:\n%s", string(out))
		}
	})

	t.Run("QuickModeExecution", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping full quick-mode oracle gate execution under -short")
		}
		cmd := exec.Command(bashPath, scriptPath, "--quick")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("oracle_accuracy_gate.sh --quick failed: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "ORACLE ACCURACY GATE PASSED") {
			t.Errorf("Expected pass banner in output, got:\n%s", string(out))
		}
	})
}

// TestPreCommitHookScript tests scripts/pre-commit.sh and .git/hooks/pre-commit synchronization.
func TestPreCommitHookScript(t *testing.T) {
	root := getRepoRoot(t)
	scriptPath := filepath.Join(root, "scripts", "pre-commit.sh")
	hookPath := filepath.Join(root, ".git", "hooks", "pre-commit")

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found; skipping bash script tests")
	}

	t.Run("BashSyntaxCheck", func(t *testing.T) {
		cmd := exec.Command(bashPath, "-n", scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Syntax check for pre-commit.sh failed: %v\nOutput: %s", err, string(out))
		}
	})

	t.Run("HelpOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "-h")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("pre-commit.sh -h failed: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "--min-coverage") || !strings.Contains(string(out), "linting, and code coverage") {
			t.Errorf("Expected usage instructions in output, got:\n%s", string(out))
		}
	})

	t.Run("InvalidOption", func(t *testing.T) {
		cmd := exec.Command(bashPath, scriptPath, "--unknown-flag")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected non-zero exit for invalid option, got success:\n%s", string(out))
		}
		if !strings.Contains(string(out), "Invalid option") {
			t.Errorf("Expected Invalid option error in output, got:\n%s", string(out))
		}
	})

	t.Run("GitHookInstalledAndSynced", func(t *testing.T) {
		if _, err := os.Stat(filepath.Join(root, ".git")); os.IsNotExist(err) {
			t.Skip(".git directory not present")
		}
		scriptBytes, err := os.ReadFile(scriptPath)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", scriptPath, err)
		}
		hookBytes, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", hookPath, err)
		}
		if string(scriptBytes) != string(hookBytes) {
			t.Errorf(".git/hooks/pre-commit differs from scripts/pre-commit.sh; run 'make install-hooks'")
		}
	})
}
