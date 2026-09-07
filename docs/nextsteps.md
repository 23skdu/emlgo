# Next Steps: 10-Part Improvement Plan

This document outlines the remaining architectural and quality improvements for the `emlgo` math library based on a deep-dive analysis of the CPU, GPU, and JIT compilation subsystems.

---

## 10-Part Improvement Plan

```mermaid
graph TD
    A[Deep Code Analysis] --> B[Platform Parity]
    A --> C[Code Quality]
    A --> D[JIT Enhancement]
    A --> E[Build & CI]
    
    B --> B1[1. ARM64 NEON Assembly Kernels]
    B --> B2[2. ARM64 SVE/SVE2 Assembly Kernels]
    B --> B3[3. ARM64 Transcendental Kernels]
    
    C --> C1[4. Worker Pool for Fused Operations]
    C --> C2[5. Eliminate Branchless Misnomers]
    C --> C3[6. Refactor Parallelization Boilerplate]
    
    D --> D1[7. JIT Function Call Codegen]
    D --> D2[8. JIT Arbitrary Exponent Support]
    
    E --> E1[9. Fix CI/CD & Linter Config]
    E --> E2[10. Go Documentation & Error Consistency]
```

### 1. ARM64 NEON Assembly Kernels for Arithmetic & Unary Ops
* **Current State:** `internal/eml/simd_arm64.s` is empty. All NEON-labeled functions (`addNEON`, `subNEON`, `mulNEON`, `divNEON`, `sqrtNEON`, etc.) in `simd_arm64.go` are plain Go loops, not actual NEON assembly.
* **Proposed Plan:**
  - Write hand-tuned ARM64 NEON assembly in `simd_arm64.s` utilizing 128-bit `V0.2D`-`V7.2D` vector registers for 2-wide `float64` operations (`FADD`, `FMUL`, `FDIV`, `FSQRT`).
  - Implement bitwise NEON kernels for `Abs` (using `BIC` with sign-bit mask) and `Neg` (using `EOR` with sign-bit mask).
  - Add scalar broadcast (`VMOV` + `FADD`/`FMUL`) for `AddScalarNEON` and `MulScalarNEON`.
  - Wire assembly functions via `//go:noescape` declarations in `simd_arm64.go` and update the ARM64 dispatch layer.

### 2. ARM64 SVE/SVE2 Vector-Length Agnostic Assembly Kernels
* **Current State:** `simd_arm64.go` implements SVE via Go-level VL-aware loops (e.g., `addSVE`, `mulSVE`). No actual SVE assembly instructions are used.
* **Proposed Plan:**
  - Write VLA (Vector-Length Agnostic) SVE assembly kernels using predicate registers (`P0-P7`) and scalable vector registers (`Z0-Z31`) for 2-wide double operations.
  - Use `WHILELT` to generate progressive predicates and `FADD_ZZZ`/`FMUL_ZZZ`/`FDIV_ZZZ`/`FSQRT_ZZZ` for arithmetic.
  - Support both SVE2 (Graviton 3/4, Apple M4) and baseline SVE via runtime detection from `/proc/self/auxv`.
  - Add transcendental SVE kernels (Exp, Log, Sin, Cos, Tan) using SVE2 polynomial approximation instructions where available.

### 3. ARM64 Transcendental Batch Kernels (NEON)
* **Current State:** ARM64 transcendental batch operations (`ExpSIMD`, `LogSIMD`, `SinSIMD`, `CosSIMD`, `TanSIMD`) fall through to `parallelizeGeneric` which calls scalar `math.*` in goroutines. The AVX2 hand-coded transcendental kernels have no ARM64 equivalent.
* **Proposed Plan:**
  - Implement NEON-vectorized Cody-Waite range reduction for Sin/Cos/Tan using `FMLS` (fused multiply-subtract) with precomputed constant tables.
  - Port the minimax polynomial evaluation for Exp/Log/Sin/Cos/Tan to NEON using `FMLA` (fused multiply-add) with 2-wide `float64` vectors.
  - Target Apple Silicon (M1/M2/M3/M4) and AWS Graviton3/4 as primary validation platforms.
  - Achieve ≥2x speedup over scalar `math.*` on ARM64 for batch sizes ≥32.

### 4. Route Fused Operations Through Pre-Allocated Worker Pool
* **Current State:** `ExpMulBatch`, `ExpAddBatch`, `LogDivBatch`, `LogSubBatch` in `fused.go` spawn new goroutines via `go func(...)` per chunk instead of sending jobs to the pre-allocated `jobQueue` worker pool initialized in `simd.go`. This defeats the purpose of the pool and creates unnecessary goroutine overhead.
* **Proposed Plan:**
  - Refactor `fused.go` to send chunked work to `jobQueue` using the `parallelJob` type, consistent with `parallelizeGeneric` and `parallelizeSinCos`.
  - Add a `isFused` flag to `parallelJob` (or create a `fusedJob` variant) to carry both input slices and the fused operation selector.
  - Benchmark before/after to verify reduced goroutine allocation overhead for batch sizes in the 256–8192 range.

### 5. Eliminate Branchless Misnomers and Implement True Branchless Variants
* **Current State:** `AbsBranchless`, `MinBranchless`, `MaxBranchless`, `SelectBranchless`, `SelectNaNBranchless` in `fused.go` all use explicit `if` branches despite their names.
* **Proposed Plan:**
  - Rename existing functions to remove the "Branchless" suffix (e.g., `AbsBranchless` → `Abs`) to eliminate misleading API surface.
  - Implement true branchless versions using bitwise tricks: `AbsBranchless` → `(x ^ (x >> 63)) | (x >> 63)` pattern via `math.Float64frombits`/`math.Float64bits`.
  - For `Min`/`Max`, use the branchless `a - ((a - b) & ((a - b) >> 63))` pattern or rely on the compiler's conditional-move optimization.
  - Expose both variants in the public API with clear documentation of tradeoffs.

### 6. Refactor Parallelization Boilerplate into Shared Helper
* **Current State:** The parallelization pattern (check `SmallCutoff`, compute `chunkSize`, iterate chunks, spawn goroutines, `wg.Wait()`) is copy-pasted into ~15 functions in `pkg/arithmetic/arith.go` and `fused.go`.
* **Proposed Plan:**
  - Create a generic `parallelForEach(n int, fn func(start, end int))` helper in `internal/eml/simd.go` that encapsulates the chunking and worker pool dispatch.
  - Refactor all ~15 parallelized functions in `arithmetic.go` and `fused.go` to use this helper.
  - This reduces code duplication from ~300 lines of boilerplate to ~30 lines, making future changes to the parallelization strategy (e.g., work-stealing, adaptive chunking) a single-point edit.

### 7. JIT Function Call Codegen (sin, cos, exp, log, sqrt)
* **Current State:** The JIT codegen (`codegen.go:217`) returns `fmt.Errorf("function calls not supported in JIT codegen: %s", v.Name)` for all `FunctionCall` nodes. Only arithmetic operators (`+`, `-`, `*`, `/`, `^`) and variables are compiled to native code.
* **Proposed Plan:**
  - Implement `CALL` codegen for supported math functions by emitting `MOV` of the function address into a temp register and using `CALL` instruction with proper stack alignment (16-byte RSP before call).
  - Support: `sin`, `cos`, `exp`, `log`, `sqrt` as initial set (all available as Go-internal symbols via `runtime` or via PIC call to shared library).
  - Handle function arguments: unary functions take the argument from `dst` register and return result in `xmm0`.
  - Add comprehensive tests in `codegen_test.go` verifying function call codegen produces correct results.

### 8. JIT Arbitrary (Non-Integer) Exponent Support
* **Current State:** The JIT `genPow` function (`codegen.go:222-278`) only supports constant non-negative integer exponents. Attempting `x^2.5` or `x^y` (where y is a variable) returns an error.
* **Proposed Plan:**
  - For variable exponents: emit codegen for `exp(y * log(x))` as a fallback path (limited to amd64 SSE2 math calls).
  - For fractional constant exponents (e.g., `x^0.5`): decompose into `sqrt(x)` call for `.5`, or general `exp(c * log(x))` for arbitrary fractions.
  - For negative integer exponents: emit `1.0 / genPow(x, |n)` to extend the existing binary exponentiation.
  - Add JIT codegen tests covering `x^0.5`, `x^-1`, `x^2.5`, and `x^y` (variable exponent).

### 9. Fix CI/CD Pipeline and Consolidate Linter Configuration
* **Current State:** 
  - CI workflow (`ci.yml`) tests Go 1.21–1.23, but `go.mod` declares `go 1.26.1` (a future version). CI will fail on `go mod download`.
  - Duplicate golangci-lint configs exist: `.golangci.yml` (with deprecated linters `structcheck`, `varcheck`) and `.golangci.yaml` (different linter set). Only one is used.
* **Proposed Plan:**
  - Update `go.mod` to a valid Go version (e.g., `go 1.23` matching the latest stable).
  - Update CI matrix to test Go 1.22 and 1.23 (drop 1.21).
  - Delete `.golangci.yaml` and keep `.golangci.yml` as the single source of truth.
  - Update `.golangci.yml`: remove deprecated `structcheck`/`varcheck`, add `errcheck` with `check-type-assertions: true`.
  - Add a CI step for `go vet` and `staticcheck` as separate lint stages.

### 10. Go Documentation Coverage and Error Consistency
* **Current State:**
  - Most exported functions in `internal/eml/simd.go`, `fused.go` lack godoc comments.
  - Inconsistent error handling: `Batch()` returns `error`, but `*SIMD()` functions panic on length mismatch. `ExpMulTo` panics but `ExpMulBatch` returns a sentinel.
  - No CHANGELOG or release notes exist.
* **Proposed Plan:**
  - Add godoc comments to all exported functions in `internal/eml/` (SIMD functions, fused operations, worker pool utilities).
  - Standardize error handling: introduce `ValidateSlices(args ...[]float64) error` helper and use it consistently. For API consistency, decide on panic (performance-critical path) vs error (user-facing API) and document the rationale.
  - Create `CHANGELOG.md` with semantic versioning entries starting from current state.
  - Ensure `go doc` renders useful descriptions for all public symbols.
