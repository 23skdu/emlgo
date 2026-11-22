# EML Go Engineering Roadmap: Remaining Work & Recommendations

An evidence-based engineering roadmap for `emlgo` tracking the remaining open items.
Completed improvements (IMP-02 through IMP-19, and the `amd64` portions of IMP-01 and
IMP-20) have been verified, merged, and archived in Git history and release notes.
IMP-01 (ARM64 JIT atomic restore) and IMP-20 (ARM64 NEON `float64` rounding/extrema
kernels) have **remaining work** detailed below.

---

## 1. Completed Milestones (Removed)

All 18 fully-completed items from the v0.8 performance & accuracy audit (IMP-02 through
IMP-19) — including `//go:noescape` zero-alloc SIMD kernels, AVX-512 `NaN` propagation,
ARM64 `VFMLA` single-rounding FMA, zero-alloc complex batch kernels, `eml.SIMD`
vectorization, IEEE 754 `Pow`/`trig`/`hyper`/`bigmath` fixes, `sincosAVX2`, bytecode VM
& optimizer fixes, `TurboQuant4` pooling, AVX2 `dotInt8AVX2`, pipeline affine
vectorization, fastmath `NaN` guardrails, float32 `*To` APIs, and the build-tagged
Apple Metal (`-tags metal`) GPU backend — have been verified 100% complete and removed
from this backlog.

> **Note:** IMP-01 (`amd64` JIT stackguard0 hardening) and IMP-20 (float32 `*To`
> variants) are **partially complete** — their `amd64`/cross-platform portions are done,
> but ARM64-specific work remains open (see FUT-01, FUT-06 below).

---

## 2. Open Future Work

| ID | Priority | Category | Component / Files | Scope & Acceptance Criteria |
| :--- | :--- | :--- | :--- | :--- |
| **FUT-01** | **P0** | ARM64 JIT Atomic Restore | `internal/jit/codegen_arm64.go` | ✅ **Implemented** — Replaced non-atomic `LDR`/`STR` stackguard0 restore with `CASAL` (ARMv8.1+) in `callFunc` and `callFunc2`, matching `amd64`'s `LOCK CMPXCHGQ`. Encoding verified via clang (`casal x9, x10, [x11]` = `0xC8E9FD6A`). **Verification pending** — requires `linux/arm64` hardware (Graviton/QEMU) under concurrent GC pressure. |
| **FUT-02** | **P2** | ARM64 NEON Rounding & Extrema Assembly | `internal/eml/simd_arm64.s`, `internal/eml/simd_dispatch_arm64.go` | Implement native 2-wide `float64` NEON assembly kernels (`FRINTM.2D`, `FRINTP.2D`, `FRINTZ.2D`, `FRINTN.2D`, `FMAX.2D`, `FMIN.2D` with IEEE 754 `NaN`/signed-zero semantics) for `floorNEON`, `ceilNEON`, `truncNEON`, `roundNEON`, `maxNEON`, and `minNEON`, replacing the current scalar Go loops in `simd_dispatch_arm64.go`. **ARM64 Darwin (Apple Silicon) NEON kernels are unable to be verified at this time** — no ARM64 Darwin CI runner with confirmed NEON assembly execution is available. Marked as **incomplete/unverified** on `darwin/arm64`. |
| **FUT-03** | **P2** | ARM64 NEON Int8 Dot Product Assembly | `internal/eml/simd_arm64.s`, `internal/eml/simd_dispatch_arm64.go` | Implement `dotInt8NEON(a []uint8, b []int8) int32` using `UXTL`/`SXTL` + `SMLAL`/`SMLAL2` + `ADDV` (and optional `SDOT` when available) to accelerate `dispatchDotInt8SIMD` on `arm64`. |
| **FUT-04** | **P2** | ARM64 SVE / SVE2 Assembly Kernels | `internal/eml/simd_sve.go`, `internal/eml/simd_sve_arm64.s` | Add predicate-based SVE/SVE2 assembly kernels on `arm64 && linux` once SVE hardware/CI runners are available for end-to-end verification. |
| **FUT-05** | **P2** | JIT Inline Hardware Math & FMA Fusion | `internal/jit/codegen.go`, `internal/jit/codegen_arm64.go` | Inline `abs` (`ANDPD` on `amd64`), `floor`/`ceil`/`trunc` (`ROUNDSD` on `amd64`, `FRINTM`/`FRINTP`/`FRINTZ` on `arm64`), and fuse `a*b + c` into hardware FMA in the JIT compilers. |
| **FUT-06** | **P2** | Precompiled Metal `.metallib` Fast-Path Loading | `internal/gpu/metal.go`, `internal/gpu/metal_bridge_darwin_arm64.m`, `metal/Makefile` | Support optional ahead-of-time compiled `eml_metal_shaders.metallib` loading with transparent fallback to runtime MSL compilation (`MTLMathModeSafe`). |

---

## 3. New Recommendations

Based on deep analysis of the codebase, CI pipeline, and documentation, the following new recommendations are added:

### 3.1 CI/CD Improvements

- **Address Node.js 20 deprecation in GitHub Actions.** ✅ Done — All actions updated to Node 24-compatible versions: `actions/checkout@v5.0.0`, `actions/setup-go@v6.5.0`, `docker/login-action@v4.0.0`, `softprops/action-gh-release@v3.0.0`, `actions/attest-build-provenance@v4.0.0`, `codecov/codecov-action@v5.0.0`.

- **Upgrade `macos-14` → `macos-15` runner.** ✅ Done — Applied in `ci.yml` for `darwin/arm64` platform test.

- **Add an ARM64 Darwin verification job.** ✅ Done — Added `arm64-darwin-verify` job running on `macos-15` with NEON assembly test commands.

- **Pin `golangci-lint` to v2 in CI.** ✅ Done — `verify` job installs `golangci-lint/v2@v2.15.0`.

### 3.2 Code Quality & Documentation

- **Fix roadmap documentation inconsistency.** The Executive Summary table marks IMP-01 and IMP-20 as "Resolved & Verified" while the detailed sections and actual code show arm64 parity is incomplete. Update the Executive Summary to accurately reflect partial resolution (done in this revision).

- **Add `go:generate` or CI check for godoc coverage.** ✅ Done — Added `scripts/check_godoc.sh` (AST-based scanner) and `make check-godoc` target. CI `verify` job runs `make check-godoc` as a gate.

- **Standardize `TAG` default in Makefile.** ✅ Done — `TAG` default set to `v0.9.0`.

### 3.3 Performance & Accuracy

- **Verify ARM64 NEON float64 kernels on native hardware.** Once an ARM64 Linux (Graviton) or ARM64 Darwin (Apple Silicon) runner is available, implement and benchmark `floorNEON`, `ceilNEON`, `truncNEON`, `roundNEON`, `maxNEON`, `minNEON` and confirm ≥ 2x speedup over scalar loops.

- **Complete IMP-01 arm64 atomic restore.** ✅ Implemented — `CASAL`-based atomic restore added to `codegen_arm64.go` (`callFunc` and `callFunc2`). Verification under concurrent GC pressure on `linux/arm64` (QEMU or Graviton) still pending.

- **Consider SVE/SVE2 support.** `HasSVE()` currently returns `false` and no SVE kernels exist. For ARM64 servers with SVE (AWS Graviton 3+, Fujitsu A64FX), 4-wide or 8-wide SVE float64 kernels would provide additional throughput. This is a future enhancement, not a blocker.

---

## 4. Unit & Fuzz Testing Scope (100.0% Statement Coverage)

Statement coverage across all 17 non-test packages is enforced at **100.0%** (`PRECOMMIT_MIN_COVERAGE=100.0`). The fuzzing matrix includes **64 fuzz targets** across all 11 core and internal packages.

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
