# Next Steps: Performance Optimization Roadmap

This document outlines a 10-part plan to enhance performance across scalar, batch, and system-level operations.

## Completed Items (from prior roadmap)

1. `[x]` **Scalar Sqrt:** Map directly to `SQRTSD` / `FSQRT`.
2. `[x]` **Scalar FMA:** Implement `FMA(a, b, c)` as a single instruction.
3. `[x]` **Scalar Exp/Log:** Optimized polynomial implementation using FMA.
4. `[x]` **FastScalar Package:** Initial release with relaxed error handling.

---

## 10-Part Performance Improvement Plan

### 1. Expanded Batch SIMD Coverage

**Target:** Exp, Log, Sin, Cos, Tan, Pow for batch sizes > 64 elements  
**Implementation:** Add AVX2/AVX512/NEON vectorized kernels in assembly  
**Benchmark:** Achieve 2x-4x speedup vs parallel Go implementation  
**Files:** `simd_amd64.s`, `simd_arm64.s`, `simd.go`

### 2. Cache-Oblivious Algorithms

**Target:** Implement cache-tiling for batch operations > 1MB  
**Implementation:** Process data in L1/L2/L3 cache-friendly chunks (32KB/256KB/1MB)  
**Benefits:** Reduces cache misses by 40-60% for large batch operations  
**Files:** `simd.go`, `batch.go` (new file)

### 3. Adaptive Parallelization Strategy

**Target:** Use adaptive chunking based on CPU cache topology  
**Implementation:** Query `runtime.NumCPU()`, L1d cache size, and adjust workers  
**Features:** Gang scheduling, work-stealing queue for load balancing  
**Files:** `parallel.go` (new or enhanced in `internal/eml/`)

### 4. Reduced Memory Allocations

**Target:** Eliminate allocations in functions called < 1000x/second  
**Implementation:** Pre-allocate buffers, use stack-allocated temporaries  
**Benefits:** GC pressure reduction, 10-20% performance gain  
**Files:** All `pkg/*` files - audit and optimize allocation hot paths

### 5. Branchless Implementation

**Target:** Convert branch-heavy code to branchless equivalents  
**Implementation:** Use bitwise operations for conditional logic (AND/OR instead of if/else)  
**Examples:** Sign handling in `Sin`, `Cos`, domain checks in `Log`, `Pow`  
**Files:** `pkg/trig/trig.go`, `pkg/logexp/exp.go`, `pkg/arithmetic/arith.go`

### 6. Hardware-Accelerated Transcendentals

**Target:** Use `VGETEXP`, `VGETMANT`, `VSCALEF` where accuracy permits  
**Implementation:** Combine with polynomial correction for full accuracy  
**Files:** `simd_amd64.s`

### 7. Benchmark Infrastructure Enhancement

**Target:** Add latency profiling and cache simulation tools  
**Implementation:** Extend `cmd/bench` with cycle-accurate measurements  
**Features:** ULP tracking, cache miss profiling, branch mispredict tracking  
**Files:** `cmd/bench/main.go`

### 8. Improved Polynomial Evaluation

**Target:** Optimize minimax polynomial coefficients  
**Implementation:** Use Remez algorithm for tighter approximations  
**Benefits:** 5-10% reduction in polynomial degree needed  
**Files:** `internal/eml/math_helpers.go`, `pkg/fastmath/`

### 9. Batch Operation Fusion

**Target:** Fuse multiple operations to reduce memory traffic  
**Implementation:** Combine Exp+Mul, Log+Div into single pass  
**Benefits:** 20-30% memory bandwidth reduction  
**Files:** `simd.go`, batch operation functions

### 10. Microkernel Optimization

**Target:** Optimize inner loop kernels for modern CPU pipelines  
**Implementation:** Software prefetching, loop unrolling, instruction scheduling  
**Benefits:** 10-15% improvement on memory-bound operations  
**Files:** `simd_amd64.s`, `simd_arm64.s`

---

## Implementation Checklist

- [ ] **Item 1:** Expanded Batch SIMD - Add vectorized Exp, Log, Sin, Cos, Tan, Pow
- [ ] **Item 2:** Cache-Oblivious Algorithms - Implement cache tiling
- [ ] **Item 3:** Adaptive Parallelization - Dynamic chunking based on CPU topology
- [ ] **Item 4:** Reduced Memory Allocations - Eliminate allocations in hot paths
- [ ] **Item 5:** Branchless Implementation - Convert conditional logic to bitwise
- [ ] **Item 6:** Hardware-Accelerated Transcendentals - Use CPU intrinsics
- [ ] **Item 7:** Benchmark Infrastructure - Add latency/cachesim profiling
- [ ] **Item 8:** Polynomial Evaluation - Optimize minimax coefficients
- [ ] **Item 9:** Batch Operation Fusion - Fuse multiple ops in single pass
- [ ] **Item 10:** Microkernel Optimization - Software prefetch, loop unroll

---

## Future Considerations (Beyond 10-Part Plan)

- GPU/CUDA Production Readiness
- ARM SVE/SVE2 Support  
- WebAssembly SIMD Intrinsics
- JIT Polynomial Compilation