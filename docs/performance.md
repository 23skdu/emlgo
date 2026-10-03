# Comprehensive Performance Benchmark: emlgo vs Go math Library

## Executive Summary

This document provides a cross-platform performance comparison between the `emlgo` library and the Go standard `math` library. Tests were conducted on both Apple Silicon (`darwin/arm64`) and Intel/AMD (`linux/amd64`) architectures.

### High-Level Comparison

| Metric | emlgo | math | Winner |
| :--- | :--- | :--- | :--- |
| **Scalar Speed** | **0.9x - 2.0x slower** | Baseline | **math** (Intrinsics) |
| **Batch Speed (SIMD)** | Mixed: faster on some ops, slower on others | Baseline | **emlgo** (SIMD) |
| **Memory (Scalar)** | 0 allocations | 0 allocations | **Tie** |
| **Memory (Batch)** | 1 allocation (Result) | 1 allocation (Result) | **Tie** |
| **Feature Parity** | 100% | 100% | **Tie** |

**Key Findings:**

- **Go `math` is generally faster for single scalar operations** due to compiler intrinsics. `emlgo/pkg/fastmath` provides faster *approximate* `exp` (~1.5e-6 relative error) and `log` (~1e-7); its `Sin`/`Cos`/`Sqrt` simply call `math` and are bit-exact, not faster.
- **The CPU SIMD batch kernels are currently a net loss.** The hand-written AVX2 transcendental kernels re-broadcast their polynomial coefficients inside the loop and run up to 190x slower than the generic path already in the package. See [nextsteps.md](nextsteps.md) Steps 1-3, which include the measured tables and the root cause.
- **Memory usage is identical** for performance-critical paths; both libraries avoid heap allocations for scalar math.

---

## New Features

### Zero-Allocation APIs

New in-place batch operations (v0.4) eliminate allocation overhead for hot paths:

```go
// Before: allocates result slice
result := eml.ExpSIMD(x)

// After: reuses pre-allocated buffer
dst := make([]float64, len(x))
eml.ExpSIMDTo(x, dst)  // Zero allocation
```

### Fused Batch Operations

Combined operations reduce memory bandwidth by 20-30%:

| Operation | Traditional | Fused | Speedup |
| :--- | :--- | :--- | :--- |
| Exp + Mul | ExpSIMD + MulSIMD | ExpMulBatch | **1.3x** |
| Exp + Add | ExpSIMD + AddSIMD | ExpAddBatch | **1.2x** |
| Log + Div | LogSIMD + DivSIMD | LogDivBatch | **1.2x** |
| Log + Sub | LogSIMD + SubSIMD | LogSubBatch | **1.2x** |

### Parallel Chunk Sizing

Work is split into `runtime.NumCPU()` chunks, capped at `LargeCutoff`:

```go
const (
    SmallWorkloadFactor = 512  // minimum elements per worker
    LargeCutoff        = 4096 // maximum elements per worker chunk
)

// SmallCutoff is derived from the core count: 512 * NumCPU, i.e. 8192 on a
// 16-core host.
var SmallCutoff = SmallWorkloadFactor * runtime.NumCPU()

func GetParallelChunkSize(n int) int {
    if n < SmallCutoff {
        return n
    }
    chunkSize := (n + cpuNum - 1) / cpuNum
    if chunkSize > LargeCutoff {
        return LargeCutoff
    }
    return chunkSize
}
```

Chunking is based on slice size and core count only — it is **not** tuned to
cache topology.

`SmallCutoff` was a flat 256 before, which was roughly 16x too low: the pool
spawned one goroutine per core for as few as 256 elements and lost to a plain
serial loop by 4.4x. See [nextsteps.md](nextsteps.md) Step 2 for the full
crossover measurement.

---

## 1. Platform-Specific Results

### Host A: Localhost (Apple Silicon M2 - `darwin/arm64`)

#### Test Environment A

Tested with n=1,000,000 iterations

| Type | Function | emlgo (s) | math (s) | Ratio |
| :--- | :--- | :--- | :--- | :--- |
| float64 | Exp | 0.0065 | 0.0058 | 1.13x |
| float64 | PowInt | 0.0031 | 0.0162 | **0.19x** (5x Faster) |
| float64 | fastmath.Sin | 0.0154 | 0.0170 | **0.90x** (measurement noise; `fastmath.Sin` calls `math.Sin`) |
| **Batch** | **ExpBatch** | 0.0004 | 0.0005 | **0.91x** |
| **Batch** | **AddBatch** | 0.0002 | 0.0003 | **0.84x** |
| **Fused** | **ExpMulBatch** | 0.0003 | 0.0005 | **0.65x** |
| **In-Place** | **ExpSIMDTo** | 0.0002 | 0.0005 | **0.45x** |

### Host B: Remote Host (`ancalagon` - `linux/amd64` AVX2)

#### Test Environment B

Tested with n=1,000,000 iterations

| Type | Function | emlgo (s) | math (s) | Ratio |
| :--- | :--- | :--- | :--- | :--- |
| float64 | Exp | 0.0060 | 0.0058 | 1.05x |
| float64 | PowInt | 0.0031 | 0.0202 | **0.15x** (6x Faster) |
| float64 | fastmath.Sin | 0.0189 | 0.0205 | **0.92x** (measurement noise; `fastmath.Sin` calls `math.Sin`) |
| **Batch** | **ExpBatch** | 0.0003 | 0.0004 | **0.98x** |
| **Batch** | **AddBatch** | 0.0001 | 0.0002 | **0.75x** |
| **Fused** | **ExpMulBatch** | 0.0002 | 0.0005 | **0.52x** |
| **In-Place** | **ExpSIMDTo** | 0.0001 | 0.0005 | **0.31x** |

---

## 2. Feature Parity & Correctness

100% feature parity was verified between `emlgo` and `math` for all core mathematical operations.

| Test Category | Local (arm64) | Remote (amd64) |
| :--- | :--- | :--- |
| Basic Arithmetic | ✓ PASSED | ✓ PASSED |
| Trigonometric | ✓ PASSED | ✓ PASSED |
| Hyperbolic | ✓ PASSED | ✓ PASSED |
| Exponential/Log | ✓ PASSED | ✓ PASSED |
| Power/Roots | ✓ PASSED | ✓ PASSED |
| Fused Operations | ✓ PASSED | ✓ PASSED |

### Accuracy Highlights

- **Standard API (`pkg/arithmetic`):**
  - **Exact matches:** `Exp`, `Log`, `Sin`, `Cos`, `Tan`, `Sqrt`.
  - **Near-exact:** `Pow` (within 10 ULP), `Cosh` (1 ULP).
  - **Acceptable:** Inverse hyperbolic functions (within 100 ULP).
- **FastMath API (`pkg/fastmath`):**
  - **Relaxed Accuracy:** Targeted at ~1e-7 absolute error (roughly 24-bit precision).
  - **Optimized for Speed:** Achieves performance gains by using FMA-based polynomial approximations and skipping IEEE 754 edge-case handling (NaN/Inf) in the primary calculation path.

---

## 3. Memory Profile

Both libraries prioritize zero-allocation paths for performance.

- **Scalar operations:** Both use stack-only execution (0 allocs/op).
- **Batch operations:** Both require a single allocation for the result slice (1 alloc/op).
- **In-place operations (v0.4):** For pre-allocated buffers, zero additional allocations.

```go
// Scalar: 0 allocations
x := eml.Exp(1.0)

// Traditional batch: 1 allocation (result)
x := eml.ExpSIMD(input)  // allocates result

// In-place batch: 0 additional allocations
dst := make([]float64, len(input))
eml.ExpSIMDTo(input, dst)  // reuses pre-allocated
```

---

## 4. Comparative Advantage: Scalar vs Batch

The following table highlights where `emlgo` provides the most significant advantages, specifically contrasting scalar latency against batch throughput (SIMD).

| Operation | Scalar Winner | Scalar Note | Batch Winner | Batch Note |
| :--- | :--- | :--- | :--- | :--- |
| **Addition** | **math** | Compiler Intrinsic | **emlgo** | SIMD Vectorized |
| **Subtraction** | **math** | Compiler Intrinsic | **emlgo** | SIMD Vectorized |
| **Multiplication** | **math** | Compiler Intrinsic | **emlgo** | SIMD Vectorized |
| **Division** | **math** | Compiler Intrinsic | **emlgo** | SIMD Vectorized |
| **Sqrt** | **math** | Compiler Intrinsic | **emlgo** | SIMD (AVX512/NEON) |
| **Exp** | **math** | Near-parity | **emlgo** | SIMD (AVX512/NEON) |
| **Sin/Cos** | **emlgo*** | FastMath optimized | **emlgo** | SIMD (AVX512/NEON) |
| **PowInt** | **emlgo** | Optimized loop | **emlgo** | SIMD Parallelized |
| **Fused** | N/A | N/A | **emlgo** | Exp+Mul/Log+Div fused |

*\*Requires `pkg/fastmath` for peak scalar performance.*

---

## 5. Conclusion & Recommendations

### When to use `math`

- Single scalar calculations (e.g., `x = math.Sin(y)`).
- Simple logic where compiler intrinsics provide maximum speed.

### When to use `emlgo`

- **Batch processing** of large arrays or tensors.
- When **SIMD acceleration** (AVX2/AVX512/NEON) is required.
- Mathematical operations involving complex **edge-case handling** (NaN/Inf) that are optimized in `emlgo`.
- Optimized **Integer Powers** (`PowInt`), which is significantly faster than the standard `math.Pow`.
- When **memory efficiency** is critical (use in-place `*SIMDTo` functions).
- When **memory bandwidth** is limited (use fused operations like `ExpMulBatch`).

---

## Benchmark Commands

```bash
# Comprehensive benchmark (1M iterations)
go run cmd/bench/main.go -n 1000000

# Accuracy test
go run cmd/bench/main.go -accuracy

# Feature parity check
go run cmd/bench/main.go -compare

# Batch operations benchmark
go test ./internal/eml/... -bench=Benchmark -benchmem

# Run new fused operation benchmarks
go test ./internal/eml/... -run="Fused" -bench=Benchmark -benchmem
```
