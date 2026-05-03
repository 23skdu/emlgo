# Next Steps: Performance Optimization Roadmap

This document outlines a 10-part plan to enhance performance across scalar, batch, and system-level operations.

## Completed Items (from prior roadmap)

1. `[x]` **Scalar Sqrt:** Map directly to `SQRTSD` / `FSQRT`.
2. `[x]` **Scalar FMA:** Implement `FMA(a, b, c)` as a single instruction.
3. `[x]` **Scalar Exp/Log:** Optimized polynomial implementation using FMA.
4. `[x]` **FastScalar Package:** Initial release with relaxed error handling.
5. `[x]` **Batch SIMD:** Add ExpSIMD, LogSIMD, SinSIMD, CosSIMD with zero-allocation To variants.
6. `[x]` **Fused Batch Operations:** ExpMulBatch, ExpAddBatch, LogDivBatch, LogSubBatch.
7. `[x]` **Adaptive Parallelization:** GetParallelChunkSize with L1/L2 cache awareness.
8. `[x]` **Branchless Utilities:** AbsBranchless, MinBranchless, MaxBranchless.

---

## Scalar Performance Issues (v2.1 Priority)

Based on benchmark results, the following scalar operations have significant overhead compared to Go's stdlib:

| Function | Current Ratio | Issue |
|----------|--------------|-------|
| Abs | 7.32x | Function call overhead |
| Neg | 2.08x | Function call overhead |
| Inv | 2.02x | Function call overhead |
| Sin/Cos/Tan | ~1.2x | Compiler intrinsics faster |

---

## 10-Part Scalar Optimization Plan

### 1. Inline Critical Scalar Functions

**Target:** Eliminating function call overhead for hot-path scalar operations  
**Implementation:** Use `//go:inline` pragma or mark functions as `static` for inlining  
**Functions:** `Neg`, `Inv`, `Abs`, `Square`, `Cube`, `Floor`, `Ceil`, `Round`, `Trunc`  
**Benchmark:** Target < 1.5x ratio (abs), < 1.2x ratio (others)

### 2. Direct Intrinsic Mapping

**Target:** Map scalar operations directly to single CPU instructions  
**Implementation:** Use platform-specific assembly or `runtime_arch` conditional compilation  
**Mappings:** `Abs` → `PNABS`/`FNABS`, `Neg` → `NEG`, `Inv` → `1/x` approximation + Newton-Raphson  
**Files:** `pkg/arithmetic/arith.go`, `internal/eml/native_math.go`

### 3. Compiler Intrinsic Integration

**Target:** Leverage Go compiler's intrinsic knowledge for trig functions  
**Implementation:** Ensure trig functions use `math.Sin`, `math.Cos`, `math.Tan` internally  
**Fallback:** Implement separate `SinFast`, `CosFast` with relaxed accuracy  
**Files:** `pkg/trig/trig.go`

### 4. Reduced Function Call Chains

**Target:** Eliminate intermediate call layers  
**Implementation:** Flatten wrapper functions into direct implementations  
**Current:** `Abs` → `eml.Abs` → `native.Abs` → platform implementation  
**Target:** `Abs` → platform implementation (single call)

### 5. Batch Scalar Operations

**Target:** Process multiple scalars in tight loops without function call overhead  
**Implementation:** Provide batch versions: `AbsBatch`, `NegBatch`, `InvBatch`  
**Benefits:** Amortize function call overhead across N elements  
**Benchmark:** Target < 0.5x ratio for batch size ≥ 64

### 6. Stack-Allocated Temporary Buffers

**Target:** Eliminate heap allocations in scalar wrappers  
**Implementation:** Use `[64]float64` stack array instead of `make([]float64, N)`  
**Functions:** All batch operations with internal allocation  
**Files:** `internal/eml/simd.go`, `pkg/*`

### 7. Fast Path / Slow Path Separation

**Target:** Optimize common case without accuracy checks  
**Implementation:** Add `Fast` variants: `SinFast`, `CosFast`, `ExpFast`, `LogFast`  
**Trade-off:** Relax error from 0 ULP to 1-2 ULP for 2-3x speedup  
**Files:** `pkg/trig/trig.go`, `pkg/logexp/exp.go`

### 8. Architecture-Specific Optimizations

**Target:** Use CPU-specific instructions on AMD64/ARM64  
**Implementation:** Conditional compile with `GOARCH=amd64` / `GOARCH=arm64`  
**Options:** `VABS`, `VFNEG` (AVX), `FABSG`, `FNEG` (NEON)  
**Files:** `simd_amd64.s`, `simd_arm64.s`

### 9. Profiling-Guided Optimization

**Target:** Identify actual hot paths in production workloads  
**Implementation:** Add profile-guided annotation with `pprof` integration  
**Tools:** Add CPU/memory profile endpoints for batch processing jobs  
**Files:** `cmd/bench/main.go`

### 10. Benchmark-Driven Iteration

**Target:** Establish baselines and track improvements  
**Implementation:** Add microbenchmarks for each scalar function  
**Automation:** CI regression detection when ratio changes > 10%  
**Files:** Microbenchmarks in `internal/eml/*_test.go`

---

## Implementation Checklist

- [x] **Item 1:** Inline Critical Functions - Add inline pragma to hot functions
- [x] **Item 2:** Direct Intrinsic Mapping - Platform-specific assembly for Abs, Neg
- [x] **Item 3:** Compiler Intrinsic Integration - Use stdlib trig internally
- [x] **Item 4:** Reduced Function Call Chains - Flatten wrapper layers
- [x] **Item 5:** Batch Scalar Operations - Add batch Abs, Neg, Inv
- [x] **Item 6:** Stack-Allocated Buffers - Eliminate heap allocations
- [x] **Item 7:** Fast/Slow Path Separation - Add Fast variants
- [x] **Item 8:** Architecture-Specific Optimizations - Use CPU intrinsics
- [x] **Item 9:** Profiling-Guided Optimization - Add pprof integration
- [x] **Item 10:** Benchmark-Driven Iteration - Track regressions

---

## Target Ratios (v2.1)

| Function | Current | Target |
|----------|---------|--------|
| Abs | 7.32x | < 1.5x |
| Neg | 2.08x | < 1.2x |
| Inv | 2.02x | < 1.2x |
| Sin | 1.05x | < 1.0x |
| Cos | 1.20x | < 1.0x |
| Tan | 1.25x | < 1.0x |

---

## Future Considerations (Beyond v2.1)

- GPU/CUDA Production Readiness
- ARM SVE/SVE2 Support
- WebAssembly SIMD Intrinsics
- JIT Polynomial Compilation