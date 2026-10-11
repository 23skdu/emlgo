# Architecture of emlgo

`emlgo` is built on the principle that all elementary functions can be derived from a single operator, the **EML Operator**, combined with mathematical constants.

## 1. Core Operator: EML

The core of the library is the `eml(x, y)` function, defined as:

`eml(x, y) = exp(x) - ln(y)`

This operator is implemented in `internal/eml` and is the primary target for hardware acceleration.

## 2. Dispatch Layer (`internal/eml`)

The dispatch layer selects an available implementation for the target architecture; selection does not guarantee the fastest path for every operation or input size:

- **AMD64**: AVX2 and AVX-512 kernels are used for selected arithmetic operations; other paths use Go implementations.
- **ARM64**: NEON kernels are used for selected operations. SVE/SVE2 kernels are not implemented, and `HasSVE()` remains false.
- **WASM**: Batch functions use loop-unrolled Go kernels. The toolchain/runtime may lower suitable code to SIMD instructions, but the implementation does not guarantee this.
- **Generic**: Clean Go fallback for other architectures.

### Optimization Strategy

- **Batch Processing**: Selected operations on slices are chunked across cores. Real-valued work uses a worker pool above a configurable cutoff; see `SmallCutoff()` and the benchmark evidence in [nextsteps.md](nextsteps.md).
- **Scalar Kernels**: Latency-sensitive operations use direct assembly for instructions like `Sqrt` and `FMA`.
- **Transcendentals**: The default batch paths use Go math implementations. Experimental AVX2 transcendental kernels are available behind `-tags emlasm`; historical tests reported substantial performance and accuracy regressions, so measure and validate that opt-in path before use (see [nextsteps.md](nextsteps.md), Step 3).
- **float32 batch operations**: Dedicated `float32` APIs are provided; architecture acceleration and accuracy vary by operation and target.

## 3. Composable Pipeline (`internal/eml/pipeline.go`)

Double-buffered composable pipeline with zero-allocation output:

- **Buffer swapping**: Uses two pre-allocated buffers, swapping on each step to avoid per-step allocations
- **9 operations**: Exp, Log, Sqrt, Sin, Cos, Abs, Neg, MulScalar, AddScalar
- **Chainable API**: `p.Exp().MulScalar(2.0).Log().RunTo(input, output)`

## 4. Complex Number Dual-Path

The library provides first-class `complex128` support via `internal/eml/complex_batch.go`:

- **Real path (float64)**: Scalar operations use the package's Go math wrappers and selected low-level kernels; batch acceleration varies by operation and target.
- **Complex path (complex128)**: Uses `math/cmplx` for `Exp`, `Log`, `Sin`, `Cos`, `Tan` on complex arguments. The complex EML operator is defined as `Complex(x, y) = cmplx.Exp(x) - cmplx.Log(y)`.

Batch operations (`ComplexBatch`, `ComplexExpBatch`, etc.) apply `math/cmplx` element-wise in a single serial loop. The worker pool is used by the real-valued fused and parallel batch paths, not by the complex ones.

## 5. Arbitrary Precision Backend (`internal/eml/bigmath`)

The `bigmath` package provides a `math/big.Float`-based backend for symbolic verification:

- **Taylor series** implementations of Exp, Log, Sin, Cos, Tan, Atan, Asin, Acos with range reduction
- π is generated as `4·atan(1)` via the complex EML operator in `internal/constants.GeneratePi`
- **IdentityVerifier**: Tests symbolic identities at deterministic points. `VerifyIdentity` uses a reproducible positive sample set; `VerifyIdentityAt` lets callers supply domain-appropriate points and reports the first mismatch, including its input, values, relative error, and threshold.

Default precision: 256 bits. All functions accept and return `*big.Float`.

## 6. Arena Allocator (`internal/jit/arena.go`)

The arena allocator provides zero-allocation paths for JIT expression parsing and evaluation:

- **ArenaNode**: 56-byte fixed-size nodes (Kind, Op, Value, Name, Left, Right pointers)
- **Bump allocator**: `Arena.Alloc()` returns the next available node, growing the buffer exponentially when needed
- **Conversion**: `FromInterface`/`ToInterface` convert between the `Node` interface and arena nodes
- **Evaluation**: `EvalArena` evaluates arena-allocated trees without touching the heap

This eliminates garbage collector pressure in hot paths where expressions are parsed and evaluated repeatedly.
The parser limits source expressions to 1 MiB and bounds AST depth and node
count; native JIT compilation additionally limits emitted code size.
Executable code is allocated writable and then changed to read/execute before
being exposed, and `CodeHandle.Free` releases that allocation after callers
have stopped invoking the compiled function.

## 7. Canonical EML Tree Representation (`internal/jit/canonical.go`)

Canonical EML trees map any mathematical expression to a normalized form based on the EML operator:

- **EMLNode tree**: Nodes tagged as `EMLConst`, `EMLVar`, `EMLOp` (eml operator), or `EMLFunc`
- **Canonical constructors**: `CanonicalExp`, `CanonicalLog`, `CanonicalSin`, `CanonicalCos`, `CanonicalSqrt`
- **`Canonicalize`**: Converts any `Node` AST into canonical form, rewriting `exp(x)` as `eml(x, 1)`, `log(x)` as `eml(1, eml(eml(1, x), 1))`, etc.
- **`Equiv`**: Structural equivalence check for two EMLNode trees
- **`EMLSize`**: Count nodes in a canonical tree

## 8. Symbolic Differentiation (`internal/jit/canonical.go`, `Diff`)

The `Diff(n *EMLNode) *EMLNode` function performs symbolic differentiation using the chain rule:

- `d/dx eml(u, v) = exp(u)·u' − v'/v`
- Supports all standard functions: exp, log, sin, cos, sqrt, neg, add, sub, mul, div, pow
- Includes `Simplify()` integration for optimized derivative expressions
- `DiffEval(n, x)` evaluates the derivative at a point

## 9. Expression Simplification (`internal/jit/canonical.go`, `Simplify`)

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
- Work is fanned out only above `SmallCutoff` (= `SmallWorkloadFactor` × `runtime.NumCPU()`); below it the operation runs inline. Measured, the pool loses to a plain serial loop by up to 4.4x at small sizes, which is why the threshold is derived from the core count rather than fixed.

The pool is also used to score a generation of genetic-programming candidates in
parallel (`pkg/bytecode.Search`), since fitness evaluation is embarrassingly
parallel and was previously serial.

## 14. JIT Compiler (`internal/jit`)

The JIT compiler generates native AMD64 machine code from string expressions; JS/WASM uses an AST interpreter, while other targets do not provide native code generation. JIT performance depends on the expression and platform; no portable speedup is claimed.

- **25 supported functions**: sin, cos, exp, log, sqrt, tan, asin, acos, atan, abs, cbrt, log2, log10, ceil, floor, trunc, round, sinh, cosh, tanh, asinh, acosh, atanh, erf, gamma
- All three dispatch paths (codegen, tree evaluator, arena interpreter) share the single table in `internal/jit/functab.go`
- **Non-integer exponents**: `x^0.5` compiled as `exp(0.5 * log(x))`
- **Variable exponents**: `x^x` compiled as `exp(x * log(x))`
- **Binary exponentiation**: Integer powers use squaring for O(log n) multiplications
- **Register allocation**: 15 XMM registers (xmm15 reserved for input variable x)

## 15. Symbolic Regression (`pkg/bytecode`)

A linear bytecode VM with a genetic search on top:

- **Program** is a flat `Structure of Arrays` buffer (`Ops`, `Consts`,
  `VarIndices`) plus a precomputed `MaxStackDepth`.
- **All 27 function names** the JIT supports are available, resolved through one
  shared table (`functab.go`) so the tokenizer, code generator, optimizer and
  printer cannot drift apart.
- **Operators**: `ValidateStack`, `RandomProgram`, `Crossover`, `Mutate`, all
  taking an explicit `*rand.Rand` so searches are reproducible.
- **Search**: `Search(Config, Dataset, Fitness)` with tournament selection,
  elitism, crossover/mutation and coordinate descent over constant operands
  (`LocalSearchSteps`).
- **Batch evaluation** is allocation-free at steady state via `BatchScratch`.

Fitting `3x² + 2x + 1` from 24 samples reaches an RMS error of 4e-13.
Transcendental targets are approximated rather than recovered; see
docs/nextsteps.md for the measured limit.

## 16. Package Structure

- **`pkg/arithmetic`**: Basic operations (Add, Sub, Mul, Div, Sqrt, Pow, FMA).
- **`pkg/logexp`**: Exponential and Logarithmic functions.
- **`pkg/trig`**: Trigonometric, Inverse Trigonometric, and Hyperbolic functions.
- **`pkg/hyper`**: Hyperbolic functions (dedicated package).
- **`pkg/fastmath`**: High-performance scalar alternatives with relaxed IEEE 754 compliance.
- **`internal/eml/bigmath`**: Arbitrary-precision backend.
- **`internal/jit`**: JIT compiler, arena allocator, canonical trees, Diff, Simplify, Decompile, Cache.
- **`internal/gpu`**: CUDA and Metal GPU backends.
- **`internal/constants`**: Mathematical constants (e, π, ln2, √2, φ, etc.).

## 17. Design Principles

1. **Zero Allocations**: Hot paths avoid heap allocations to ensure predictable performance.
2. **Minimal Dependencies**: The library depends only on the Go standard library and `golang.org/x/sys`.
3. **Architecture-Aware**: High-level APIs automatically benefit from hardware acceleration without user intervention.
4. **Correctness First**: Tests cover documented edge-case semantics and selected backend comparisons; see [accuracy.md](accuracy.md) for the limits of that evidence.

Packages under `internal/` are restricted by Go to this module; they are not
importable by downstream modules or part of the supported external API.
