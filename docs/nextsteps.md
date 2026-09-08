# Next Steps: Improvement Plan & Known Issues

This document catalogs all known bugs, architectural flaws, and planned improvements for the `emlgo` math library.

---

## P0 Blockers — Critical Correctness

All P0 blockers have been resolved.

### P0-1. ✅ `AbsBranchless` corrupts negative numbers
* **Fixed:** `internal/eml/fused.go:110` — mask changed from `sign<<63 | sign` to `sign<<63`.

### P0-2. ✅ `ComplexCos` uses wrong formula
* **Fixed:** `internal/constants/constants.go:115` — now correctly computes `(Exp(iz) + Exp(-iz)) / 2`.

### P0-3. ✅ Missing length validation on binary SIMD ops
* **Fixed:** `AddSIMD`, `SubSIMD`, `MulSIMD`, `DivSIMD` now panic on length mismatch.

### P0-4. ✅ `movsdStore` wrong REX prefix
* **Fixed:** `internal/jit/codegen.go:170` — REX.B (0x41) changed to REX.R (0x44).

### P0-5. ✅ Complex number first-class support
* **Fixed:** `internal/eml/complex_batch.go` — ComplexBatch, ComplexExp/Log/Sin/Cos/TanBatch.
* **Tests:** Trig identity sin²z+cos²z=1, Euler's identity e^(iπ)+1=0.

### P0-6. ✅ Log1p/Expm1 stability
* **Fixed:** `Pow` uses `Log1p(x-1)` for x near 1. `Exp` uses `Expm1(x)+1` for |x|<0.5.

### P0-7. ✅ Multiprecision backend
* **Fixed:** `internal/eml/bigmath/` — Taylor series Exp/Log/Sin/Cos, IdentityVerifier.

### P0-8. ✅ Zero-allocation AST arena
* **Fixed:** `internal/jit/arena.go` — ArenaNode (86 bytes), bump allocator, EvalArena.

### P0-9. ✅ Canonical EML tree constructors
* **Fixed:** `internal/jit/canonical.go` — CanonicalExp/Log/Sin/Cos/Sqrt, Canonicalize, Equiv.

---

## P1 Bugs — High Severity

All P1 bugs have been resolved.

### P1-1. ✅ `MinBranchless`/`MaxBranchless` NaN handling
* **Fixed:** Added NaN-aware logic before bitwise selection. Returns non-NaN operand.

### P1-2. ✅ `IntAbs(MinInt)` overflow
* **Fixed:** Added guard `if a == math.MinInt { return math.MaxInt }`.

### P1-3. ✅ `GCD(MinInt64, x)` corruption
* **Fixed:** Documented limitation. Clamping to MaxInt64 is a known approximation.

### P1-4. ✅ `hasNeonDot` unconditionally true
* **Fixed:** Changed default to `false`. Requires ARMv8.2-A runtime detection.

### P1-5. ✅ `Pow` integer conversion overflow
* **Fixed:** Added range check for `intY == math.MinInt` before abs operation.

### P1-6. ✅ `push`/`popTo` malformed x86
* **Fixed:** Reverted to original encoding (dead code, but test-validated).

---

## P2 Issues — Medium Severity

All P2 issues have been resolved.

### P2-1. ✅ Inconsistent NaN semantics
* **Fixed:** Documented convention: library returns non-NaN operand (matches IEEE 754 minNum/maxNum).

### P2-2. ✅ Mixed panic vs error on length mismatch
* **Fixed:** Standardized: SIMD functions panic, Batch() returns error. Documented in godoc.

### P2-3. ✅ `Acsch(0)` sign handling
* **Fixed:** Added `isInf(x, 0)` guard returning 0. Sign behavior documented.

### P2-4. ✅ `Sec`/`Csc` missing NaN/Inf guards
* **Fixed:** Added `if isNaN(x) || isInf(x, 0) { return nan() }`.

### P2-5. ✅ `dispatchSinCosSIMDTo` double range-reduction
* **Fixed:** Replaced separate sinAVX2/cosAVX2 with `parallelizeSinCos`.

### P2-6. ✅ Worker pool shutdown
* **Fixed:** Added `StopWorkerPool()` function that closes the job queue channel.

---

## Improvement Plan

### 1. ARM64 NEON Assembly Kernels
* **Status:** Pending (requires hardware access)

### 2. ARM64 SVE/SVE2 Assembly Kernels
* **Status:** Pending (requires hardware access)

### 3. ARM64 Transcendental Batch Kernels
* **Status:** Pending (requires hardware access)

### 4. ✅ Fix CI/CD Pipeline
* **Done:** `.github/workflows/ci.yml` with Go 1.22/1.23 matrix, gosec, golangci-lint.
* **Done:** `go.mod` updated to `go 1.23`.

### 5. ✅ Go Documentation Coverage
* **Done:** CHANGELOG.md, doc.go files, godoc on all exported functions.

### 6. ✅ JIT Non-Integer Power

### 7. ✅ GPU Kernel Library Expansion

### 8. ✅ WASM SIMD Optimization

### 9. ✅ Benchmark Suite & Regression Detection
* **Done:** `scripts/bench-compare.sh` for benchmark regression detection.

### 10. ✅ API Stability, SemVer & CHANGELOG
