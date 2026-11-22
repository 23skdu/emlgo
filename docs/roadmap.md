# EML Go Engineering Roadmap: Remaining Work & Recommendations

An evidence-based engineering roadmap for `emlgo` tracking the remaining open items
after the IMP-01–IMP-20 improvement cycle. All completed improvements (IMP-02 through
IMP-19) have been verified, merged, and archived in Git history and release notes.

---

## 1. Executive Summary: Remaining Improvements

| ID | Priority | Category | Component / Files | Headline Defect, Bottleneck or Gap & Empirical Evidence |
| :--- | :--- | :--- | :--- | :--- |
| **IMP-01** | **P0** | JIT Stability & Runtime Safety | `internal/jit/codegen.go`, `internal/jit/codegen_arm64.go`, `internal/jit/codegen_test.go` | **Partially resolved:** Hardened indirect `CALL R8` transitions on `amd64` via atomic `g.stackguard0` (`16(R14)`) save/restore (`LOCK CMPXCHGQ`). **`arm64` parity is still pending** — `codegen_arm64.go` uses non-atomic `LDR`/`STR` for the restore, which can clobber concurrent GC `stackPreempt` writes. |
| **IMP-20** | **P2** | Cross-Platform SIMD Completeness (`arm64` & `float32`) | `internal/eml/simd_arm64.s`, `internal/eml/simd_dispatch_arm64.go`, `internal/eml/simd_f32.go` | **Partially resolved:** Zero-allocation `float32` `*To` variants (`ExpSIMDF32To`, `LogSIMDF32To`, `SinSIMDF32To`, `CosSIMDF32To`, `TanSIMDF32To`) are implemented. **ARM64 NEON `float64` rounding/extrema kernels (`floorNEON`, `ceilNEON`, `truncNEON`, `roundNEON`, `maxNEON`, `minNEON`) are not implemented** — dispatch falls back to scalar Go loops. ARM64 Darwin (Apple Silicon) NEON kernels **cannot be verified at this time** (no ARM64 Darwin CI runner with NEON assembly execution available). |

---

## 2. Detailed Technical Specifications

### IMP-01: JIT Go-Call Stack-Unwind Hardening (`fatal error: unknown caller pc`)

**Priority:** P0 | **Category:** JIT Stability & Runtime Safety | **Status:** Resolved on `amd64`; `arm64` atomic restore (CASAL) pending

**Affected Files:**
- `internal/jit/codegen.go:228-328` (amd64 — completed)
- `internal/jit/codegen_arm64.go:177-246` (arm64 — pending atomic restore)

**Evidence & Root Cause:**
When JIT-compiled machine code residing in anonymous `mmap` memory executes an indirect call (`CALL R8` on `amd64`, `BLR X8` on `arm64`) to a Go `math.*` function, the CPU pushes the `mmap` return PC onto the goroutine stack. If concurrent GC mark termination sets `g.stackguard0 = stackPreempt` to shrink or preempt the goroutine stack, the Go `math.*` prologue jumps into `runtime.morestack()` → `runtime.newstack()` → `runtime.shrinkstack()` → `runtime.copystack()`. When `runtime.(*unwinder).next()` encounters the `mmap` return PC (absent from Go's `pclntab`), the runtime aborts with `fatal error: unknown caller pc`.

On `amd64`, this is fixed by saving `g.stackguard0` to the stack, setting it to `g.stack.lo` (NOSPLIT semantics) before `CALL R8`, then atomically restoring via `LOCK CMPXCHGQ` so concurrent `stackPreempt` writes are preserved.

On `arm64`, the save/restore exists but uses non-atomic `LDR`/`STR`, which unconditionally overwrites `g.stackguard0` and can clobber a concurrent `stackPreempt` signal.

**Required Changes:**
1. In `internal/jit/codegen_arm64.go` (`callFunc` and `callFunc2`), replace the non-atomic restore:
   ```assembly
   LDR X9, [sp+272]
   STR X9, [X28, #16]
   ```
   with an atomic compare-and-swap using `CASAL` (ARMv8.1+) or a `LDXR`/`STXR` retry loop:
   ```assembly
   // Atomic restore: only write saved value if g.stackguard0 still equals g.stack.lo
   LDR  X9, [X28, #0]        // X9 = g.stack.lo (expected current value)
   LDR  X10, [sp+272]        // X10 = saved original stackguard0
   CASAL X9, X10, [X28, #16] // atomically swap if [X28,#16] == X9
   ```
2. Ensure all high-iteration JIT tests (`TestJITDifferentialStep11`, `FuzzJITDifferential`) release executable pages via `handle.Free()`.

**Verification:**
- `go test -short -count=1 -covermode=atomic ./internal/jit/...` passes 100/100 runs under heavy concurrent GC pressure on `linux/arm64` without stack-unwinder crashes.

---

### IMP-20: ARM64 NEON `float64` Rounding/Extrema Kernels & Vectorized `float32` Transcendental `*To` Kernels

**Priority:** P2 | **Category:** Cross-Platform SIMD Completeness (`arm64` & `float32`) | **Status:** Partially resolved — float32 `*To` variants done; ARM64 NEON float64 rounding/extrema pending; ARM64 Darwin NEON kernels unable to be verified

**Affected Files:**
- `internal/eml/simd_arm64.s` (missing: `floorNEON`, `ceilNEON`, `truncNEON`, `roundNEON`, `maxNEON`, `minNEON`)
- `internal/eml/simd_dispatch_arm64.go:398-445` (scalar Go loop fallbacks)
- `internal/eml/simd_f32.go` (float32 `*To` variants — completed)

**Evidence & Root Cause:**
1. In `internal/eml/simd_dispatch_arm64.go:398-445`, `dispatchFloorSIMDTo`, `dispatchCeilSIMDTo`, `dispatchTruncSIMDTo`, `dispatchRoundSIMDTo`, `dispatchMaxSIMDTo`, and `dispatchMinSIMDTo` execute scalar Go loops on `arm64`, even though ARMv8-A NEON provides single-instruction 2-wide `float64` vector instructions `FRINTM.2D`, `FRINTP.2D`, `FRINTZ.2D`, `FRINTN.2D`, `FMAX.2D`, and `FMIN.2D`.
2. The `float32` `*To` variants (`ExpSIMDF32To`, `LogSIMDF32To`, `SinSIMDF32To`, `CosSIMDF32To`, `TanSIMDF32To`) have been implemented in `internal/eml/simd_f32.go` and are verified on `amd64`. The ARM64 NEON float32 assembly kernels (`addNEONF32`, `subNEONF32`, `mulNEONF32`, `divNEONF32`, `sqrtNEONF32`, `fmaNEONF32`) are implemented in `internal/eml/simd_arm64_f32.s`.
3. **ARM64 Darwin (Apple Silicon) NEON kernels cannot be verified at this time** — no ARM64 Darwin CI runner with NEON assembly execution is available in the current GitHub Actions configuration (`macos-15` runners provide ARM64 but NEON assembly verification requires native execution which has not been confirmed).

**Required Changes:**
1. Implement `floorNEON`, `ceilNEON`, `truncNEON`, `roundNEON`, `maxNEON`, and `minNEON` in `internal/eml/simd_arm64.s` using `FRINTM.2D`, `FRINTP.2D`, `FRINTZ.2D`, `FRINTN.2D`, `FMAX.2D`, `FMIN.2D`.
2. Wire them into `simd_dispatch_arm64.go` following the SIMD-lanes-then-scalar-tail pattern (process 2-wide NEON lanes, then scalar fallback for the tail).
3. Verify on ARM64 Linux (Graviton/QEMU) and, when available, on ARM64 Darwin (Apple Silicon).

**Verification:**
- `go test -bench=. ./internal/eml/...` on ARM64 Linux verifies `0 allocs/op` on the new NEON dispatch paths and ≤ 1 ULP accuracy against `math.Floor`/`math.Ceil`/`math.Trunc`/`math.RoundToEven`/`math.Max`/`math.Min`.
- ARM64 Darwin NEON verification is **blocked** — no verified ARM64 Darwin test environment available.

---

## 3. New Recommendations

Based on deep analysis of the codebase, CI pipeline, and documentation, the following new recommendations are added:

### 3.1 CI/CD Improvements

1. **Address Node.js 20 deprecation in GitHub Actions.** All pinned actions (`actions/checkout@v4.2.2`, `actions/setup-go@v5.6.0`, `docker/login-action@v3.7.0`, etc.) target Node.js 20, which GitHub is deprecating. Update to Node.js 24-compatible action versions (e.g., `actions/checkout@v5`, `actions/setup-go@v6`) to eliminate deprecation warnings and future-proof the pipeline.

2. **Upgrade `macos-14` → `macos-15` runner.** Already applied in the current `ci.yml` for the `darwin/arm64` platform test. Ensure all macOS-dependent jobs use `macos-15` or `macos-latest` going forward.

3. **Add an ARM64 Darwin verification job.** The `platform-test` matrix includes `darwin/arm64` on `macos-15`, but NEON assembly execution has not been confirmed. Add a dedicated job that runs `go test -v ./internal/eml/...` on `macos-15` with `-tags` for NEON assembly and reports pass/fail explicitly.

4. **Pin `golangci-lint` to v2 in CI.** The `verify` job already installs `github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0`. Ensure the `golangci/golangci-lint-action` (if re-added) uses `version: v2.15.0` instead of `latest` to avoid v1/v2 incompatibility with `.golangci.yml` schema v2.

### 3.2 Code Quality & Documentation

5. **Fix roadmap documentation inconsistency.** The Executive Summary table marks IMP-01 and IMP-20 as "Resolved & Verified" while the detailed sections and actual code show arm64 parity is incomplete. Update the Executive Summary to accurately reflect partial resolution (done in this revision).

6. **Add `go:generate` or CI check for godoc coverage.** The project enforces 100% statement coverage via `scripts/verify.sh` but has no automated check for exported-symbol documentation coverage. Consider adding a `golangci-lint` `revive` rule (`exported`) or a custom CI step that runs the godoc scanner (`/tmp/opencode/godocscan.go` or equivalent) and fails if any exported symbol lacks a doc comment.

7. **Standardize `TAG` default in Makefile.** The `TAG` variable default was updated from `0.6` to `v0.9.0` to match the current release. Ensure `docker-compose.yml` and `Dockerfile` references stay in sync with the latest release tag.

### 3.3 Performance & Accuracy

8. **Verify ARM64 NEON float64 kernels on native hardware.** Once an ARM64 Linux (Graviton) or ARM64 Darwin (Apple Silicon) runner is available, implement and benchmark `floorNEON`, `ceilNEON`, `truncNEON`, `roundNEON`, `maxNEON`, `minNEON` and confirm ≥ 2x speedup over scalar loops.

9. **Complete IMP-01 arm64 atomic restore.** Implement `CASAL`-based atomic restore in `codegen_arm64.go` and verify under concurrent GC pressure on `linux/arm64` (QEMU or Graviton). This is the last remaining JIT stability gap.

10. **Consider SVE/SVE2 support.** `HasSVE()` currently returns `false` and no SVE kernels exist. For ARM64 servers with SVE (AWS Graviton 3+, Fujitsu A64FX), 4-wide or 8-wide SVE float64 kernels would provide additional throughput. This is a future enhancement, not a blocker.

---

## 4. Unit & Fuzz Testing Scope (100.0% Statement Coverage)

Statement coverage across all 17 non-test packages is enforced at **100.0%** (`PRECOMMIT_MIN_COVERAGE=100.0`). The fuzzing matrix includes **64 fuzz targets** across all 11 core and internal packages. See `docs/roadmap.md` (Section 3 of the previous version) for the full fuzz target inventory.

---

## 5. Accuracy & Oracle Code Quality Gates (`scripts/oracle_accuracy_gate.sh`)

The repository enforces a 5-stage automated numerical oracle gate via `./scripts/oracle_accuracy_gate.sh` (`make oracle-gate` / `make accuracy-gate`):

1. **Gate 1 — GNU MPFR 4.2.2 (`MPFR_RNDN`) Golden Reference Vectors** (134 vectors, 17 operations)
2. **Gate 2 — Arbitrary-Precision `bigmath` (256/512-bit) & Package-Wide `TestOracle*` Differential Suites**
3. **Gate 3 — JIT Native Machine-Code vs AST Canonical & Symbolic Derivative Oracle** (3,000 random expression trees)
4. **Gate 4 — Full Type, Precision, Quantization & Batch Domain Validation (`cmd/validate`)** (435 checks)
5. **Gate 5 — Feature Parity & Tightened Per-Operation ULP Sweeps (`cmd/bench -compare` & `cmd/bench -accuracy`)**

---

## 6. Makefile Targets & `scripts/verify.sh` Pre-Commit Hook Infrastructure

### 6.1 Makefile Target Matrix

| Category | Makefile Target | Command / Action Performed |
| :--- | :--- | :--- |
| **Target Discovery** | `make help` | Prints every Makefile target with a one-line description, followed by the overridable variables (`FUZZTIME`, `IMAGE_REPO`, `TAG`, `CUDA_ARCH`) and their defaults. |
| **Git Hooks** | `make makehooks` (alias: `make install-hooks`) | Marks `scripts/verify.sh` executable and installs it into `.git/hooks/pre-commit` (and syncs `scripts/pre-commit.sh`). |
| **Verification Gate** | `make verify` | Runs `./scripts/verify.sh` (staged artifact hygiene, `go vet`, `golangci-lint`, `gosec`, fast oracle accuracy gate, and 100.0% statement coverage gate). |
| **Full CI Quality Gate** | `make check` | Runs `vet`, `fmt-check`, `lint`, `security` (`gosec` + `govulncheck`), `oracle-gate`, and `test-cover`. |
| **Go Security Checks** | `make security` | Runs both `make gosec` and `make govulncheck`. |
| **Go Security (Static)** | `make gosec` | Runs `golangci-lint run --enable-only=gosec ./...` across all packages. |
| **Go Security (Vuln DB)** | `make govulncheck` | Scans module and stdlib call stacks against the Go vulnerability database via `govulncheck ./...`. |
| **Linting & Formatting** | `make lint` | Runs `make fmt-check` (`gofmt -s -l .`), `make vet` (`go vet ./...`), and `golangci-lint run ./...` (`staticcheck`, `govet` shadow, `errcheck`, `ineffassign`, `unused`, `misspell`, `gosec`). |
| **CPU & Memory Profiling** | `make profile` | Runs both `make profile-cpu` and `make profile-mem`. |
| **CPU Profiling** | `make profile-cpu` | Generates `cpu.prof` on core SIMD kernels (`ExpSIMD`, `LogSIMD`, `SinSIMD`, `CosSIMD`, `AddSIMD`, `FmaSIMD`) and prints top 20 hot symbols via `go tool pprof -top`. |
| **Memory Profiling** | `make profile-mem` | Generates `mem.prof` on SIMD, `Pipeline`, and `Complex` batch kernels and prints top 20 allocation sites via `go tool pprof -top -alloc_space`. |
| **Allocation Benchmark** | `make bench-mem` | Runs `go test -bench=. -benchmem` across `./internal/eml/...` and `./pkg/...`. |
| **Oracle Accuracy Gate** | `make oracle-gate` (alias: `make accuracy-gate`) | Executes the 5-stage numerical accuracy gate `./scripts/oracle_accuracy_gate.sh`. |
| **Fuzz Testing** | `make fuzz` | Runs all 32 anchored fuzz targets across `internal/jit`, `internal/eml`, `internal/eml/bigmath`, `internal/gpu`, and `pkg/*` (`FUZZTIME ?= 30s`). |
| **Artifact Cleanup** | `make clean` | Removes build, test, coverage, profiling, WASM, Node, and Python venv/`__pycache__` artifacts. |
| **Docker Cleanup** | `make clean-docker` | Removes tagged `$(IMAGE_REPO)` Docker images and prunes the Docker build cache. |
| **Test Cache Cleanup** | `make clean-cache` | Clears the Go test cache (`go clean -testcache`). |
| **Full Cleanup** | `make clean-all` | Runs `clean`, `clean-docker`, and `clean-cache`. |

### 6.2 Installing and Running `scripts/verify.sh`

```bash
# Install scripts/verify.sh as .git/hooks/pre-commit:
make makehooks

# Run the pre-commit verification suite manually:
./scripts/verify.sh

# Run pre-commit verification including CPU & memory pprof profiling checks:
./scripts/verify.sh --with-profile
```