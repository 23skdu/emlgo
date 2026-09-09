# Architecture of emlgo

`emlgo` is built on the principle that all elementary functions can be derived from a single operator, the **EML Operator**, combined with mathematical constants.

## 1. Core Operator: EML

The core of the library is the `eml(x, y)` function, defined as:

`eml(x, y) = exp(x) - ln(y)`

This operator is implemented in `internal/eml` and is the primary target for hardware acceleration.

## 2. Dispatch Layer (`internal/eml`)

The dispatch layer automatically detects the host architecture and selects the most efficient implementation:

- **AMD64**: Targets AVX2 and AVX512.
- **ARM64**: Targets NEON (ASIMD), SVE/SVE2.
- **WASM**: SIMD128 with 8-wide unrolled kernels.
- **Generic**: Clean Go fallback for other architectures.

### Optimization Strategy

- **Batch Processing**: Operations on slices are chunked and processed using SIMD assembly kernels.
- **Scalar Kernels**: Latency-sensitive operations use direct assembly for instructions like `Sqrt` and `FMA`.
- **Parallelization**: For very large slices, the library automatically distributes work across multiple CPU cores via a worker pool.
- **float32 SIMD**: Dedicated `float32` batch operations for memory-constrained workloads.

## 3. Composable Pipeline (`internal/eml/pipeline.go`)

Double-buffered composable pipeline with zero-allocation output:

- **Buffer swapping**: Uses two pre-allocated buffers, swapping on each step to avoid per-step allocations
- **9 operations**: Exp, Log, Sqrt, Sin, Cos, Abs, Neg, MulScalar, AddScalar
- **Chainable API**: `p.Exp().MulScalar(2.0).Log().RunTo(input, output)`

## 4. Complex Number Dual-Path

The library provides first-class `complex128` support via `internal/eml/complex_batch.go`:

- **Fast path (float64)**: All scalar `float64` operations use native EML+SIMD paths.
- **Complex path (complex128)**: Uses `math/cmplx` for `Exp`, `Log`, `Sin`, `Cos`, `Tan` on complex arguments. The complex EML operator is defined as `Complex(x, y) = cmplx.Exp(x) - cmplx.Log(y)`.

Batch operations (`ComplexBatch`, `ComplexExpBatch`, etc.) process slices of `complex128` in parallel using the worker pool.

## 5. Arbitrary Precision Backend (`internal/eml/bigmath`)

The `bigmath` package provides a `math/big.Float`-based backend for symbolic verification:

- **Taylor series** implementations of Exp, Log, Sin, Cos, Tan, Atan, Asin, Acos with range reduction
- **Machin's formula** for computing π at arbitrary precision
- **IdentityVerifier**: Tests symbolic identities (e.g., sin²(x) + cos²(x) = 1) by evaluating at random high-precision points

Default precision: 256 bits. All functions accept and return `*big.Float`.

## 6. Arena Allocator (`internal/jit/arena.go`)

The arena allocator provides zero-allocation paths for JIT expression parsing and evaluation:

- **ArenaNode**: 86-byte fixed-size nodes (Kind, Op, Value, Name, Left, Right pointers)
- **Bump allocator**: `Arena.Alloc()` returns the next available node, growing the buffer exponentially when needed
- **Conversion**: `FromInterface`/`ToInterface` convert between the `Node` interface and arena nodes
- **Evaluation**: `EvalArena` evaluates arena-allocated trees without touching the heap

This eliminates garbage collector pressure in hot paths where expressions are parsed and evaluated repeatedly.

## 7. Canonical EML Tree Representation (`internal/jit/canonical.go`)

Canonical EML trees map any mathematical expression to a normalized form based on the EML operator:

- **EMLNode tree**: Nodes tagged as `EMLConst`, `EMLVar`, `EMLOp` (eml operator), or `EMLFunc`
- **Canonical constructors**: `CanonicalExp`, `CanonicalLog`, `CanonicalSin`, `CanonicalCos`, `CanonicalSqrt`
- **`Canonicalize`**: Converts any `Node` AST into canonical form, rewriting `exp(x)` as `eml(x, 1)`, `log(x)` as `eml(1, eml(eml(1, x), 1))`, etc.
- **`Equiv`**: Structural equivalence check for two EMLNode trees
- **`EMLSize`**: Count nodes in a canonical tree

## 8. Symbolic Differentiation (`internal/jit/canonical.go:303-390`)

The `Diff(n *EMLNode) *EMLNode` function performs symbolic differentiation using the chain rule:

- `d/dx eml(u, v) = exp(u)·u' − v'/v`
- Supports all standard functions: exp, log, sin, cos, sqrt, neg, add, sub, mul, div, pow
- Includes `Simplify()` integration for optimized derivative expressions
- `DiffEval(n, x)` evaluates the derivative at a point

## 9. Expression Simplification (`internal/jit/canonical.go:205-301`)

The `Simplify(n *EMLNode) *EMLNode` function performs:

- Constant folding: evaluates subtrees where all leaves are `EMLConst`
- Identity reduction: `eml(x, 1)` → `exp(x)` direct node
- Algebraic simplifications: 0+x=x, 1*x=x, 0*x=0, x^0=1, x^1=x, etc.
- Double negation elimination: `neg(neg(x)) = x`

## 10. EML Decompiler & LaTeX Emitter (`internal/jit/decompile.go`)

Converts canonical EML trees back to human-readable forms:

- **`Decompile(n *EMLNode) string`**: Parenthesized infix output
- **`DecompileLaTeX(n *EMLNode) string`**: LaTeX math mode output
- **`DecompileNodeToExpr(n *EMLNode) string`**: Uses JIT formatter
- **CLI**: `emlcli decompile <expr>` for command-line decompilation

## 11. JIT Expression Cache (`internal/jit/cache.go`)

LRU cache for compiled JIT expressions:

- Mutex-protected with `container/list` for LRU eviction (max 1024 entries)
- `CompileCached(expr string) (Func, error)` — returns cached result
- `ClearJITCache()` for long-running programs

## 12. Numerically Stable Compound Forms

The library uses `Log1p`/`Expm1` for improved accuracy near critical points:

- **`Pow(x, y)` near x=1**: Uses `Exp(y * Log1p(x-1))` instead of `Exp(y * Log(x))` to avoid catastrophic cancellation in `Log(x)` when x ≈ 1.
- **`Exp(x)` near x=0**: Uses `Expm1(x) + 1` for |x| < 0.5 to preserve precision.
- **`Log(x)` near x=1**: Uses `Log1p(x-1)` directly.

## 13. Worker Pool with Graceful Shutdown

The parallel batch processing system uses a fixed-size worker pool:

- Pre-spawned goroutines consume from a buffered job channel
- `StopWorkerPool()` closes the channel (protected by `sync.Once`)
- After shutdown, no further parallel operations should be submitted

## 14. JIT Compiler (`internal/jit`)

The JIT compiler generates x86-64 machine code from string expressions:

- **17 supported functions**: sin, cos, exp, log, sqrt, tan, asin, acos, atan, abs, cbrt, log2, log10, ceil, floor, trunc, round
- **Non-integer exponents**: `x^0.5` compiled as `exp(0.5 * log(x))`
- **Variable exponents**: `x^x` compiled as `exp(x * log(x))`
- **Binary exponentiation**: Integer powers use squaring for O(log n) multiplications
- **Register allocation**: 15 XMM registers (xmm15 reserved for input variable x)

## 15. Package Structure

- **`pkg/arithmetic`**: Basic operations (Add, Sub, Mul, Div, Sqrt, Pow, FMA).
- **`pkg/logexp`**: Exponential and Logarithmic functions.
- **`pkg/trig`**: Trigonometric, Inverse Trigonometric, and Hyperbolic functions.
- **`pkg/hyper`**: Hyperbolic functions (dedicated package).
- **`pkg/fastmath`**: High-performance scalar alternatives with relaxed IEEE 754 compliance.
- **`internal/eml/bigmath`**: Arbitrary-precision backend.
- **`internal/jit`**: JIT compiler, arena allocator, canonical trees, Diff, Simplify, Decompile, Cache.
- **`internal/gpu`**: CUDA and Metal GPU backends.
- **`internal/constants`**: Mathematical constants (e, π, ln2, √2, φ, etc.).

## 16. Design Principles

1. **Zero Allocations**: Hot paths avoid heap allocations to ensure predictable performance.
2. **Minimal Dependencies**: The library depends only on the Go standard library and `golang.org/x/sys`.
3. **Architecture-Aware**: High-level APIs automatically benefit from hardware acceleration without user intervention.
4. **Correctness First**: The `internal/eml` layer ensures that edge cases (NaN, Inf) are handled consistently across all architectures.
