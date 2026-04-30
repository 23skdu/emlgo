# Next Steps: Surpassing Go Standard Library Scalar Performance

While `emlgo` currently provides significant speedups for **batch** operations via SIMD, standard scalar operations often trail behind Go's `math` library due to compiler intrinsics. This document outlines the roadmap to achieve scalar parity and eventual superiority.

## 1. Direct Assembly Scalar Kernels

The primary reason for `math` library superiority is the use of compiler intrinsics that emit direct assembly instructions (e.g., `FSQRT`).

- **Action:** Implement scalar versions of core arithmetic and transcendental functions in `simd_amd64.s` and `simd_arm64.s`.
- **Target:** Bypass Go function call overhead and leverage the same underlying hardware instructions used by the compiler.

## 2. Algorithmic Optimizations

Go's standard library prioritize portability and strict IEEE 754 compliance. We can optimize for speed on modern architectures:

- **FMA (Fused Multiply-Add):** Use `VFMADD` instructions for polynomial evaluations (Sin, Cos, Exp). This reduces latency and improves precision.
- **Polynomial Refinement:** Use lower-degree minimax polynomials for specialized ranges, reducing the number of multiplications needed.
- **Fast Paths:** Implement branchless checks for common values (0, 1, integers) to avoid complex math logic for simple cases.

## 3. Lighter Robustness (FastMath Mode)

Standard `math` performs extensive checks for signed zeros, subnormal numbers, and specific NaN patterns.

- **Action:** Introduce a `fastscalar` sub-package that relaxes non-critical IEEE 754 edge-case handling in exchange for a 2x-3x speedup.
- **Benefit:** High-performance workloads often don't require bit-perfect handling of signed zero or subnormal behavior.

## 4. Inlining and Link-Time Optimization

- **Action:** Ensure all scalar wrappers in `pkg/arithmetic` are simple enough for the Go compiler to inline automatically.
- **Action:** Use `//go:noescape` and `//go:nosplit` annotations on assembly stubs to minimize runtime overhead.

## 5. Benchmarking and Iteration

- **Matrix:** Expand `cmd/bench` to track scalar latency (cycles/op) in addition to throughput.
- **Profiling:** Use `pprof` to identify remaining overhead in the dispatch logic.

---

### Implementation Status:

1. `[x]` **Scalar Sqrt:** Map directly to `SQRTSD` / `FSQRT`.
2. `[x]` **Scalar FMA:** Implement `FMA(a, b, c)` as a single instruction.
3. `[x]` **Scalar Exp/Log:** Optimized polynomial implementation using FMA.
4. `[x]` **FastScalar Package:** Initial release with relaxed error handling for maximum throughput.