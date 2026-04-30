# Architecture of emlgo

`emlgo` is built on the principle that all elementary functions can be derived from a single operator, the **EML Operator**, combined with mathematical constants.

## 1. Core Operator: EML

The core of the library is the `eml(x, y)` function, defined as:

`eml(x, y) = exp(x) - ln(y)`

This operator is implemented in `internal/eml` and is the primary target for hardware acceleration.

## 2. Dispatch Layer (`internal/eml`)

The dispatch layer automatically detects the host architecture and selects the most efficient implementation:

- **AMD64**: Targets AVX2 and AVX512.
- **ARM64**: Targets NEON (ASIMD).
- **Generic**: Clean Go fallback for other architectures.

### Optimization Strategy

- **Batch Processing**: Operations on slices are chunked and processed using SIMD assembly kernels.
- **Scalar Kernels**: Latency-sensitive operations use direct assembly for instructions like `Sqrt` and `FMA`.
- **Parallelization**: For very large slices, the library automatically distributes work across multiple CPU cores.

## 3. Package Structure

- **`pkg/arithmetic`**: Basic operations (Add, Sub, Mul, Div, Sqrt, Pow, FMA).
- **`pkg/logexp`**: Exponential and Logarithmic functions.
- **`pkg/trig`**: Trigonometric, Inverse Trigonometric, and Hyperbolic functions.
- **`pkg/fastmath`**: High-performance scalar alternatives with relaxed IEEE 754 compliance.

## 4. Design Principles

1. **Zero Allocations**: Hot paths avoid heap allocations to ensure predictable performance.
2. **Minimal Dependencies**: The library depends only on the Go standard library.
3. **Architecture-Aware**: High-level APIs automatically benefit from hardware acceleration without user intervention.
4. **Correctness First**: The `internal/eml` layer ensures that edge cases (NaN, Inf) are handled consistently across all architectures.
