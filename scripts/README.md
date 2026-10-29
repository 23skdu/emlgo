# EML Go Scripts Reference

This directory contains validation, oracle accuracy gate, benchmarking, WebAssembly testing, and Git hook scripts for `emlgo`. All scripts are regression-tested by [`test/scripts_test.go`](../test/scripts_test.go) via `make test-scripts`.

---

## Script Inventory

| Script | Makefile Target | Purpose |
| :--- | :--- | :--- |
| [`oracle_accuracy_gate.sh`](oracle_accuracy_gate.sh) | `make oracle-gate` / `make accuracy-gate` | 5-stage numerical accuracy gate (MPFR 4.2.2 golden vectors, `bigmath` & `ExactLogOracle` differential suites, JIT vs AST canonical parity, `cmd/validate`, and `cmd/bench -compare` / `-accuracy`). |
| [`pre-commit.sh`](pre-commit.sh) | `make install-hooks` | Git pre-commit hook enforcing binary/artifact hygiene, `go vet`, `golangci-lint`, and **100.0% statement coverage** (`PRECOMMIT_MIN_COVERAGE=100.0`). |
| [`generate_mpfr_vectors.c`](generate_mpfr_vectors.c) | — | C generator using GNU MPFR 4.2.2 (`MPFR_RNDN`) to emit exact `binary64` (53-bit) and `binary32` (24-bit) reference vectors to `internal/accuracyvectors/mpfr-vectors.csv`. |
| [`run_validation.sh`](run_validation.sh) | — | Wrapper around `cmd/validate` supporting type filtering (`-t`), verbose output (`-v`), and post-run binary cleanup (`-c`). |
| [`run_benchmark.sh`](run_benchmark.sh) | — | Wrapper around `cmd/bench` supporting feature parity (`-c`), ULP accuracy (`-a`), regression checks (`-r`), and CSV/JSON export (`-o`, `-f`). |
| [`bench-compare.sh`](bench-compare.sh) | — | Compares current working tree benchmarks against a baseline Git ref (`main` by default) using `benchstat`. |
| [`check_bench_regression.py`](check_bench_regression.py) | — | Statistical Mann-Whitney U + median slowdown gate used in CI to detect benchmark regressions across repeated samples. |
| [`wasm_test.sh`](wasm_test.sh) | `make wasm` | Builds `wasm/emlgo.wasm` and runs headless Node.js / Playwright browser validation and benchmarks. |

---

## Usage Examples

### 1. Oracle Accuracy Gate (`oracle_accuracy_gate.sh`)

Run the full 5-stage oracle accuracy gate:
```bash
./scripts/oracle_accuracy_gate.sh
# or via Makefile:
make oracle-gate
```

Run fast Go-only oracle gates (skipping `cmd/validate` and `cmd/bench` CLI sweeps):
```bash
./scripts/oracle_accuracy_gate.sh --quick
```

Verify MPFR 4.2.2 vector regeneration byte-for-byte (when `libmpfr` 4.2.2 headers are installed):
```bash
./scripts/oracle_accuracy_gate.sh --regen-check --verbose
```

### 2. Git Pre-Commit Hook (`pre-commit.sh`)

Install or sync the tracked pre-commit hook into `.git/hooks/pre-commit`:
```bash
make install-hooks
```

Run the pre-commit gate manually:
```bash
./scripts/pre-commit.sh
```

Override the minimum statement coverage threshold if needed during local debugging:
```bash
./scripts/pre-commit.sh --min-coverage 99.0
```

