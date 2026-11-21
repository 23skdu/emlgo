# EML Go Engineering Roadmap

An evidence-based engineering roadmap for `emlgo`. Completed items are removed once 100% verified (unit tests, fuzz tests, numerical accuracy oracle gates, `100.0%` statement coverage, and zero dead code).

> **Completed Milestones (Removed):**
> All 20 items from the v0.8 performance & accuracy audit (`IMP-01` through `IMP-20`) — including JIT `g.stackguard0` hardening on `amd64`/`arm64`, `//go:noescape` zero-alloc SIMD kernels, AVX-512 `NaN` propagation, ARM64 `VFMLA` single-rounding FMA, zero-alloc complex batch kernels, `eml.SIMD` vectorization, IEEE 754 `Pow`/`trig`/`hyper`/`bigmath` fixes, `sincosAVX2`, bytecode VM & optimizer fixes, `TurboQuant4` pooling, AVX2 `dotInt8AVX2`, pipeline affine vectorization, fastmath `NaN` guardrails, float32 `*To` APIs, and the build-tagged Apple Metal (`-tags metal`) GPU backend — have been verified 100% complete and removed from this backlog.

---

## 1. Open Future Work

| ID | Priority | Category | Component / Files | Scope & Acceptance Criteria |
| :--- | :--- | :--- | :--- | :--- |
| **FUT-01** | **P2** | ARM64 NEON Rounding & Extrema Assembly | `internal/eml/simd_arm64.s`, `internal/eml/simd_dispatch_arm64.go` | Implement native 2-wide `float64` NEON assembly kernels (`FRINTM.2D`, `FRINTP.2D`, `FRINTZ.2D`, `FRINTN.2D`, `FMAX.2D`, `FMIN.2D` with IEEE 754 `NaN`/signed-zero semantics) for `floorNEON`, `ceilNEON`, `truncNEON`, `roundNEON`, `maxNEON`, and `minNEON`, replacing the current scalar Go loops in `simd_dispatch_arm64.go`. |
| **FUT-02** | **P2** | ARM64 NEON Int8 Dot Product Assembly | `internal/eml/simd_arm64.s`, `internal/eml/simd_dispatch_arm64.go` | Implement `dotInt8NEON(a []uint8, b []int8) int32` using `UXTL`/`SXTL` + `SMLAL`/`SMLAL2` + `ADDV` (and optional `SDOT` when available) to accelerate `dispatchDotInt8SIMD` on `arm64`. |
| **FUT-03** | **P2** | ARM64 SVE / SVE2 Assembly Kernels | `internal/eml/simd_sve.go`, `internal/eml/simd_sve_arm64.s` | Add predicate-based SVE/SVE2 assembly kernels on `arm64 && linux` once SVE hardware/CI runners are available for end-to-end verification. |
| **FUT-04** | **P2** | JIT Inline Hardware Math & FMA Fusion | `internal/jit/codegen.go`, `internal/jit/codegen_arm64.go` | Inline `abs` (`ANDPD` on `amd64`), `floor`/`ceil`/`trunc` (`ROUNDSD` on `amd64`, `FRINTM`/`FRINTP`/`FRINTZ` on `arm64`), and fuse `a*b + c` into hardware FMA in the JIT compilers. |
| **FUT-05** | **P2** | Precompiled Metal `.metallib` Fast-Path Loading | `internal/gpu/metal.go`, `internal/gpu/metal_bridge_darwin_arm64.m`, `metal/Makefile` | Support optional ahead-of-time compiled `eml_metal_shaders.metallib` loading with transparent fallback to runtime MSL compilation (`MTLMathModeSafe`). |

---

## 2. Unit & Fuzz Testing Matrix (`100.0%` Statement Coverage)

Statement coverage across all 17 non-test packages in `emlgo` is enforced at **100.0%** (`PRECOMMIT_MIN_COVERAGE=100.0`) on both `amd64` and `darwin/arm64`. The fuzzing matrix spans **67 fuzz targets** across all 11 core and internal packages:

| Package | Statement Coverage | Fuzz Targets (`go test -fuzz=...`) | Oracle & Differential Invariants Tested |
| :--- | :---: | :--- | :--- |
| `internal/eml` | **100.0%** | `FuzzAddSIMD`, `FuzzSubSIMD`, `FuzzMulSIMD`, `FuzzDivSIMD`, `FuzzSqrtSIMD`, `FuzzExpSIMD`, `FuzzLogSIMD`, `FuzzSinSIMD`, `FuzzCosSIMD`, `FuzzTanSIMD`, `FuzzFmaSIMD`, `FuzzFloat32SIMD`, `FuzzUnaryBatchResults`, `FuzzExpMulBatch`, `FuzzLogDivBatch`, `FuzzFloorSIMD`, `FuzzCeilSIMD`, `FuzzTruncSIMD`, `FuzzRoundSIMD`, `FuzzMaxSIMD`, `FuzzMinSIMD`, `FuzzFmaSIMDF32`, `FuzzPipelineDifferential`, `FuzzComplexDotProduct`, `FuzzEmlSIMD`, `FuzzSinCosSIMD`, `FuzzComplexBatchDifferential`, `FuzzDotInt8AndSatInt8`, `FuzzARM64NEONKernels` (`arm64`) | SIMD vs scalar differential parity, odd tail lengths, in-place aliasing, subnormals, complex128/complex64 batch parity, `int8` saturating & VNNI dot products, ARM64 NEON assembly parity |
| `internal/eml/bigmath` | **100.0%** | `FuzzExpLnRoundtrip`, `FuzzBigMathPrecision`, `FuzzBigMathTrigAndInverseIdentities` | 256-bit $\exp(\ln(x)) = x$, $\sin^2(x)+\cos^2(x)=1$, $\arctan(\tan(x))=x$, `IdentityVerifier` |
| `internal/jit` | **100.0%** | `FuzzParse`, `FuzzParseNoPanic`, `FuzzEval`, `FuzzCompileMatchesEval`, `FuzzJITDifferential`, `FuzzScientificFloatRoundTrip`, `FuzzCanonicalSimplifyAndDiff` | Native machine-code (`amd64`/`arm64`) vs AST interpreter differential parity, parser robustness, symbolic `Simplify` & `DiffEval` |
| `internal/gpu` | **100.0%** (stub) / **98.7%** (`-tags metal`) | `FuzzLaunchConfig`, `FuzzMatMulLaunchConfig`, `FuzzULPDiff`, `FuzzBatchVerifierVerifyOp`, `FuzzMetalFloat32Kernels` (`-tags metal`), `FuzzMetalComplex64Kernels` (`-tags metal`) | 1D/2D grid/block bounds, ULP distance symmetry & NaN/Inf handling, `BatchVerifier.VerifyOp`/`VerifyOpF32`, Apple Silicon Metal float32/complex64 GPU kernel accuracy |
| `pkg/arithmetic` | **100.0%** | `FuzzAbs`, `FuzzAdd`, `FuzzMod`, `FuzzRemainder`, `FuzzPow`, `FuzzLogBase`, `FuzzAddSubSatInt8`, `FuzzGCDAndLCM`, `FuzzInt8DistanceAndDot` | IEEE 754 `Mod`/`Remainder`/`Pow`, Euclidean `GCD`/`LCM` divisibility on `int64`, `int8`/`uint8` dot & L2/cosine distances |
| `pkg/bytecode` | **100.0%** | `FuzzOptimizerDifferential`, `FuzzCompileEMLRoundTrip`, `FuzzEvalBatchColumnarDifferential` | RPN bytecode `Optimize(p)` vs `p` equivalence, columnar batch VM `EvalBatchColumnarScratch` vs row `Eval` |
| `pkg/fastmath` | **100.0%** | `FuzzFastMathDifferential`, `FuzzGuardrailsAndRegularizedEml` | Minimax/Chebyshev `FastExp`/`FastLog` error bounds, `FastEmlClipped`/`FastEmlRegularized` finiteness across $\mathbb{R}^2$ |
| `pkg/hyper` | **100.0%** | `FuzzHyperbolicAccuracy`, `FuzzHyperbolicInverseRoundtrip` | $\cosh^2(x)-\sinh^2(x)=1$, $\operatorname{asinh}(\sinh(x))=x$, cancellation-free small-$x$ & asymptotic large-$x$ |
| `pkg/logexp` | **100.0%** | `FuzzLogExpDifferential`, `FuzzLogExpBatchAndSubnormalOracle` | `Exp`/`Log`/`ExpBatch`/`LogBatch` vs `ExactLogOracle` across subnormals $2^{-1074}..2^{-1022}$ |
| `pkg/quant` | **100.0%** | `FuzzPackUnpack4Bit`, `FuzzTurboQuant4EncodeDistance` | 4-bit nibble pack/unpack bijectivity, `EncodeTurboQuant4` & `TurboQuant4Distance` non-negativity and error bounds |
| `pkg/trig` | **100.0%** | `FuzzTrigPythagoreanIdentity`, `FuzzTrigReciprocalAndInverse` | $\sin^2(x)+\cos^2(x)=1$, $\arcsin(x)+\arccos(x)=\pi/2$, degree/radian round-trip, batch trig parity |

---

## 3. Accuracy & Oracle Code Quality Gates (`scripts/oracle_accuracy_gate.sh`)

The repository enforces a 5-stage automated numerical oracle gate via `./scripts/oracle_accuracy_gate.sh` (`make oracle-gate` / `make accuracy-gate`):

1. **Gate 1 — GNU MPFR 4.2.2 (`MPFR_RNDN`) Golden Reference Vectors (`internal/accuracyvectors/mpfr-vectors.csv`)**:
   - Verifies provenance header and **134 exact binary64 (53-bit) & binary32 (24-bit) golden vectors** across 17 operations (`exp`, `log`, `expm1`, `log1p`, `sin`, `cos`, `tan`, `sinh`, `cosh`, `tanh`, `asinh`, `acosh`, `atanh`, `sqrt`, `cbrt`, `eml`, `pow`).
   - Supports live C recompilation and byte-for-byte diffing via `./scripts/oracle_accuracy_gate.sh --regen-check` (including Homebrew `mpfr`/`gmp` on macOS).
2. **Gate 2 — Arbitrary-Precision `bigmath` (256/512-bit) & Package-Wide `TestOracle*` Differential Suites**:
   - Executes `internal/eml/bigmath`, `internal/eml` SIMD differential accuracy gates, and dedicated `TestOracle*` suites across every `pkg/*` package.
3. **Gate 3 — JIT Native Machine-Code vs AST Canonical & Symbolic Derivative Oracle**:
   - Verifies 3,000 random expression trees (`TestJITDifferentialStep11`), all 25 built-in math functions, and symbolic differentiation (`DiffEval` / `SecondDerivativeEval`).
4. **Gate 4 — Full Type, Precision, Quantization & Batch Domain Validation (`cmd/validate`)**:
   - Runs all 435 validation checks across `int*`, `uint*`, `float32`, `float64`, `complex64`, `complex128`, `fastmath`, `quant`, and `batch`.
5. **Gate 5 — Feature Parity & Tightened Per-Operation ULP Sweeps (`cmd/bench -compare` & `cmd/bench -accuracy`)**:
   - Enforces `0 failed` across domain-valid parity and ULP sweeps using canonical `orderedFloat64Bits` ULP distance.

---

## 4. Makefile Targets & `scripts/verify.sh` Pre-Commit Hook Infrastructure

### 4.1 Makefile Target Matrix

| Category | Makefile Target | Command / Action Performed |
| :--- | :--- | :--- |
| **Target Discovery** | `make help` | Prints every Makefile target with a one-line description and overridable variables (`FUZZTIME`, `IMAGE_REPO`, `TAG`, `CUDA_ARCH`). |
| **Git Hooks** | `make makehooks` (alias: `make install-hooks`) | Marks `scripts/verify.sh` executable and installs it into `.git/hooks/pre-commit` (and syncs `scripts/pre-commit.sh`). |
| **Verification Gate** | `make verify` | Runs `./scripts/verify.sh` (staged artifact hygiene, `go vet`, `golangci-lint`, `gosec`, fast oracle accuracy gate, and `100.0%` statement coverage gate). |
| **Full CI Quality Gate** | `make check` | Runs `vet`, `fmt-check`, `lint`, `security` (`gosec` + `govulncheck`), `oracle-gate`, and `test-cover`. |
| **Apple Metal & ARM64** | `make test-metal` / `make test-arm64-simd` / `make fuzz-metal` | Runs build-tagged Metal GPU unit/accuracy/race/fuzz tests (`-tags metal`) and ARM64 NEON + `-tags purego` tests on `darwin/arm64`. |
| **Dead Code Detection** | `make deadcode` | Runs `deadcode -test ./...` (and `-tags metal` on `darwin/arm64`) to enforce zero unreachable functions. |
| **Go Security Checks** | `make security` (`gosec` + `govulncheck`) | Runs static security analysis (`golangci-lint --enable-only=gosec`) and Go vulnerability DB scanning (`govulncheck ./...`). |
| **Linting & Formatting** | `make lint` / `make lint-tools` | Runs `fmt-check`, `vet`, `golangci-lint`, plus external Homebrew linters (`lint-shell` via `shellcheck`/`shfmt`, `lint-docker` via `hadolint`, `lint-workflows` via `actionlint`, `lint-yaml`, `lint-docs`). |
| **Homebrew Tooling** | `make brew-tools` | Installs all macOS Homebrew validation tools (`golangci-lint`, `govulncheck`, `gosec`, `shellcheck`, `shfmt`, `hadolint`, `actionlint`, `yamllint`, `markdownlint-cli`, `mpfr`). |
| **CPU & Memory Profiling** | `make profile` (`profile-cpu` + `profile-mem`) | Generates CPU and memory allocation profiles on core SIMD, `Pipeline`, and `Complex` batch kernels. |
| **Allocation Benchmark** | `make bench-mem` | Runs `go test -bench=. -benchmem` across `./internal/eml/...` and `./pkg/...`. |
| **Oracle Accuracy Gate** | `make oracle-gate` (alias: `make accuracy-gate`) | Executes the 5-stage numerical accuracy gate `./scripts/oracle_accuracy_gate.sh`. |
| **Fuzz Testing** | `make fuzz` | Runs all anchored fuzz targets across `internal/jit`, `internal/eml`, `internal/eml/bigmath`, `internal/gpu`, and `pkg/*` (`FUZZTIME ?= 30s`). |
| **Artifact Cleanup** | `make clean` / `clean-docker` / `clean-cache` / `clean-all` | Removes build, test, coverage, profiling, CUDA/Metal, WASM, Node, Python, Docker, and Go test-cache artifacts. |

### 4.2 Installing and Running `scripts/verify.sh`

```bash
# Install scripts/verify.sh as .git/hooks/pre-commit:
make makehooks

# Run the pre-commit verification suite manually:
./scripts/verify.sh

# Run deadcode, external brew linters, and Metal GPU tests:
make deadcode lint-tools test-metal
```
