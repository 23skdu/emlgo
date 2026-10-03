# Changelog

All notable changes to the emlgo library will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.0] - 2026-10-02

### Added
- **Vectorized Pipeline Engine & Complex SIMD (Step 28)**: `Pipeline` now dispatches steps directly to SIMD kernels (`ExpSIMDTo`, `LogSIMDTo`, `SinSIMDTo`, `CosSIMDTo`, `SqrtSIMDTo`), achieving a **28.2x speedup** at $n=65,536$. Consecutive affine steps are fused into single-pass FMA transforms ($a_2(a_1 x + b_1) + b_2$). Buffer allocation in `NewPipeline`/`RunTo` reduced by 50%. `ComplexDotProduct` and `ComplexDotProductC64` hooked to the shared worker pool via `ForEachChunk` with 4-way unrolled Neumaier compensated summation (**5.1x speedup** at $n=1\text{M}$).
- **Vectorized Integer Saturating SIMD (Step 30)**: AVX2 assembly kernels `VPADDSB` and `VPSUBSB` in `internal/eml/simd_amd64.s` with pure Go fallback, exposed via `AddSatInt8SIMDTo` and `SubSatInt8SIMDTo` (**17.0x speedup** at $n=65,536$). Added destination-passing `AddBatchInt8To`, `SubBatchInt8To`, `MulBatchInt8To` in `pkg/arithmetic`.
- **Flat VM Jump Table & Scratch Reuse (Step 29)**: Replaced bytecode opcode hash map lookups with a static direct-indexed array `opUnaryTable [256]mathFunc` in `pkg/bytecode/functab.go` and `eval.go`. Updated `BatchScratch.ensure` to check capacity rather than exact length, guaranteeing 0 allocs across varying batch sizes.
- **Complete `CompileEML` Function Set (Step 25)**: All 27 supported math functions wired through `unaryOpcode()` in `pkg/bytecode/compiler.go`. Canonical evaluation (`EMLEval`, `EMLEvalRegularized`) and symbolic differentiation (`Diff`) in `internal/jit/canonical.go` wired to all 27 operations.
- Native Go fuzz tests (`testing.F`) added across the repository: `FuzzOptimizerDifferential`, `FuzzScientificFloatRoundTrip`, `FuzzHyperbolicAccuracy`, `FuzzBigMathPrecision`, `FuzzCompileEMLRoundTrip`, `FuzzPipelineDifferential`, `FuzzComplexDotProduct`, and `FuzzAddSubSatInt8` (over 10 million executions with 0 failures).

### Fixed
- **Bytecode Optimizer RPN Corruption (Step 21)**: Rewrote `pkg/bytecode/optimizer.go` with AST-based constant folding and algebraic identity reductions, fixing RPN operand inversion ($x - 2 \to x - 2$, $x + 0 \to x$, $x \cdot 1 \to x$, $0 - x \to -x$).
- **Parser Scientific Notation & Whitespace (Step 22)**: `internal/jit/parser.go` now correctly lexes scientific floats (`1e-5`, `2.5E+3`) and tab/newline whitespace; added variable scoping in `ParseWithVars`.
- **Decompiler Precedence & Associativity (Step 26)**: Fixed parenthesization in `internal/jit/decompile.go` to respect operator precedence and left-associativity, ensuring mathematical round-trip invariance for `(x+1)*2`, `x/(y*z)`, `(x+1)^2`, and `x-(y-z)`.
- **JIT Power Codegen & Unary Negation (Step 27)**: Emitted ABIInternal 2-argument JIT call to `math.Pow(x, y)` for non-constant power expressions, correctly handling negative bases (e.g. $(-2)^3 = -8$, $(-2)^4 = 16$) and $0^0 = 1$ without NaN corruption. Implemented 1-cycle bitwise `xorpd` unary minus with a 16-byte sign mask (`0x8000000000000000`).
- **Hyperbolic Cancellation & Log Underflow (Step 23)**: `pkg/hyper/hyper.go` updated to use odd symmetry and logarithmic asymptotic reduction in `Asinh`, eliminating the $-10^8$ `-Inf` cliff; `Sinh` and `Tanh` use `math.Expm1` near zero. Fixed `pkg/logexp/exp.go` to return `-Inf` for `Log(0)` per IEEE 754.
- **Arbitrary-Precision Limits in `bigmath` (Step 24)**: Removed $-750$ underflow clamp in `internal/eml/bigmath/bigmath.go`; range reduction power $k$ computed via `MantExp` directly; memoized $\pi$ and $\ln 2$ via `sync.Map`; `Asin` computed purely in `big.Float` without precision loss.
- **Normalized Length Mismatch Policy (Step 30)**: Replaced silent slice truncation in `pkg/arithmetic/arith_int.go` with standardized panics on slice length mismatch.
- `internal/eml/simd_accuracy_test.go` gates every CPU batch kernel against the
  standard library, using a mixed absolute+relative tolerance
  (`1e-14 + 1e-15*|want|`) because a purely absolute budget is unsatisfiable for
  large results and a purely relative one is meaningless where the true value
  crosses zero. Coverage includes dense grids, the scalar tail for lengths 1-63,
  `sin^2+cos^2`, the +-pi/2 range-reduction boundaries, agreement with the scalar
  path, and domain edges. The gates reported 6170 real failures before the
  dispatch fix and pass on every path now.
- `pkg/bytecode`: the VM now supports all 27 function names (`eml`, `sin`, `cos`,
  `tan`, `exp`, `log`, `ln`, `sqrt`, `asin`, `acos`, `atan`, `abs`, `cbrt`,
  `log2`, `log10`, `ceil`, `floor`, `trunc`, `round`, `sinh`, `cosh`, `tanh`,
  `asinh`, `acosh`, `atanh`, `erf`, `gamma`) instead of 5. Both dispatch sites
  now read from a single table (`pkg/bytecode/functab.go`); `emitAST` carried a
  third, which accepted `sin`/`cos` and then failed at emit time.
  `Optimize` folds every function's constant subtree.
- `pkg/bytecode/ga.go`: a genetic search -- `Dataset`, `Fitness`, `Config`,
  `Result`, `Search` and a ready-made `LeastSquares` objective -- completing the
  symbolic-regression path the operators alone left unfinished. Scoring is
  parallel across candidates. `genetic.go` gained `RandomProgram` plus
  seedable `Crossover`/`Mutate`, so searches are reproducible via `Config.Seed`.
  Fitting y = 3x^2+2x+1 from 24 samples reaches RMS error 4.5e-13.
- `Config.LocalSearchSteps` adds coordinate descent over a candidate's constants
  (default 30). Random jitter is a random walk, so pinning a coefficient to 3.0
  takes many generations; this is the difference between RMS 0.21 and 4e-13 on
  the same seed.
- `internal/eml`: `Pipeline.Reset()`, `Pipeline.Len()`, `Pipeline.SubScalar` and
  `Pipeline.DivScalar`.
- `SmallWorkloadFactor` documents the per-worker minimum that sets `SmallCutoff`.

- Symbolic differentiation engine (`Diff`, `DiffEval`) with chain rule support
- Expression simplification (`Simplify`) with constant folding, identity reduction, algebraic simplifications
- EML decompiler (`Decompile`, `DecompileLaTeX`, `DecompileNodeToExpr`) for tree → infix/LaTeX conversion
- Composable zero-allocation `Pipeline` API with buffer swapping (Exp, Log, Sqrt, Sin, Cos, Abs, Neg, MulScalar, AddScalar)
- `float32` SIMD batch operations (Exp, Log, Sqrt, Add, Sub, Mul, Div, Abs, Neg, Inv, Sin, Cos, Tan, scalar ops)
- JIT expression LRU cache (`CompileCached`, `ClearJITCache`, max 1024 entries)
- `round` function added to JIT compiler (25 functions total)
- `ParseWithVars`; unknown identifiers are parsed as named variables (the `vars` list is currently ignored)
- `ParseError` struct with structured error information (Pos, Token, Msg)
- `NewFloatFromInt(n int64)` constructor for bigmath
- `bigmath.Tan`, `bigmath.Atan`, `bigmath.Asin`, `bigmath.Acos` for complete transcendental coverage
- `StopWorkerPool` concurrency safety via `sync.Once`
- `emlcli decompile <expr>` subcommand for command-line expression decompilation

### Changed
- `go.mod`: `go 1.25.0` → `go 1.26.0`; `golang.org/x/sys` `v0.43.0` → `v0.48.0`.
  Minimum supported Go is now 1.26.
- CI: Go matrix `1.24/1.25` → `1.26/1.27`; `actions/checkout@v4` → `v7`,
  `actions/setup-go@v5` → `v7`, `codecov/codecov-action@v4` → `v7`,
  `golangci/golangci-lint-action@v6` → `v9` (pinned to golangci-lint `v2.14.0`),
  `softprops/action-gh-release@v2` → `v3`; `gosec` pinned to `v2.29.0`.
- CI: new `cross-build` job covering linux/amd64, linux/arm64, darwin/amd64,
  darwin/arm64, windows/amd64, js/wasm and wasip1/wasm.
- Removed the no-op `detectCacheTopology` stub in `internal/eml`, which existed only
  to support a "cache-aware parallelization" claim that the code never implemented.
- `cmd/validate.printSummary` now delegates to a pure `summarize` helper so the
  reporting logic is testable; `cmd/bench.checkRegression` delegates to
  `classifyRegressions` for the same reason. Both still exit non-zero on failure.

### Fixed
- `internal/jit`: the arena interpreter silently evaluated 9 of the 25 built-in
  functions to zero. `callMathFunc` handled only 16 names (`round`, `sinh`, `cosh`,
  `tanh`, `asinh`, `acosh`, `atanh`, `erf`, `gamma` returned 0). The three dispatch
  tables (codegen, `EvalVars`, `EvalArena`) are now consolidated into a single
  platform-independent table in `internal/jit/functab.go`, so they cannot drift apart
  again.
- `internal/jit`: removed two `Simplify` rules that were unreachable *and* not
  identities for `eml(u,v) = exp(u) - ln(v)`. They claimed
  `eml(eml(1,eml(x,1)),eml(1,1)) = -x` (actually `e^(e-x) - 1`) and
  `eml(eml(1,x),eml(x,1)) = 1/x` (actually `e^e/x - x`). Bottom-up recursion had
  already folded `eml(1,1)` to `e` and `eml(x,1)` to `exp(x)`, so neither rule could
  ever fire.
- `pkg/bytecode`: `CompileExpr` accepted malformed expressions such as `"1 +"`,
  `"*"`, `"exp"` and `"1 * * 2"`, producing programs that **panicked** with
  `index out of range` during `Eval`. All three compile entry points
  (`CompileExpr`, `CompileAST`, `CompileEML`) now run the existing `ValidateStack`
  check and report an error instead.
- `cmd/validate`: `complexLog(0, 0)` returned `NaN` instead of `-Inf`, and the
  large-magnitude path dropped the `(i/r)^2` term (7e-4 relative error at
  `1e200 + 1e200i`). Replaced with an overflow-free `logMagnitude` helper that
  factors out the larger component.
- `cmd/validate`: `complexSqrt` used the approximation `sqrt(x+iy) ≈ sqrt(x) +
  iy/(2 sqrt(x))`, which is only valid for `|y| << |x|`; it was ~11% wrong at
  `1e200 + 1e200i`. Replaced with an exact formulation that scales by
  `max(|r|,|i|)` and therefore cannot overflow. Infinite-imaginary inputs no longer
  produce `NaN`.
- `cmd/bench`: `benchmarkBatch` hardcoded 100000 iterations, ignoring the documented
  `-n` flag.
- `cmd/bench`: `runComplex64Benchmarks` labelled results `Complex64` while every other
  runner used lowercase, so `-type complex64` results did not match the `Type`/`Name`
  keys used elsewhere.
- WASI: `GOOS=wasip1 GOARCH=wasm go build ./...` failed with 10 errors in
  `internal/jit/jit_posix.go` (`unix.Mmap` and friends are unavailable there),
  despite WASI support being documented. `jit_posix.go` now excludes `wasip1` and a
  new `internal/jit/jit_wasi.go` reports unsupported executable memory.
- `internal/eml`: `GOOS=js GOARCH=wasm go vet ./...` failed because
  `coverage_extra_test.go` referenced amd64-only dispatch symbols without a build
  tag. Moved to `coverage_amd64stub_test.go` with the correct constraint.
- `checkRegression` reported ">15% slower" while the threshold was 10%; the message
  and threshold now agree.
- Repository-wide `gofmt -s` violations (6 files, including non-test
  `pkg/fastmath/fast_exp.go`).

- `pkg/bytecode`: `CompileExpr` accepted malformed input such as `"1 +"`, `"*"`
  and `"exp"`, producing programs that panicked with an index-out-of-range during
  `Eval`. All three compile entry points now validate the opcode stream.
- `internal/eml`: the scalar logarithm returns NaN for non-positive arguments
  while the batch path returned `-Inf` at zero, so `Log(0)` and `LogBatch([0])`
  disagreed. Added `logScalar` and routed every platform's `dispatchLogSIMDTo`
  through it.
- `pkg/bytecode`: `RandomProgram` generated unbalanced stacks (final heights of 4,
  7 and 11) because its balancing binary operator was appended only half the
  time.
- `pkg/bytecode`: `splitPoints` looked for stack height 0, which a well-formed RPN
  expression only reaches at index 0, so `Crossover` could only swap whole
  programs and never recombined subexpressions. Split points are now the
  complete-subexpression boundaries (height 1); programs are also built as
  expression trees and flattened, giving 3.69 usable crossover points per program
  instead of 1.
- `pkg/bytecode`: `Mutate` swapped `OpConst` for `OpVar` (both arity 0) without
  editing the positional operand streams, desynchronising them and panicking
  during evaluation. Leaves are now converted explicitly, with the operand
  cursors tracked.
- `pkg/bytecode`: crossover spliced opcode ranges while leaving the positional
  constant and variable-index streams untouched.
- `internal/eml`: the AVX2 transcendental kernels have no scalar tail and require
  a multiple of the vector width; documented and covered by a test.

### Performance
- The AVX2 transcendental kernels (`expAVX2`, `logAVX2`, `sinAVX2`, `cosAVX2`,
  `tanAVX2`) re-materialised every polynomial coefficient *inside* the loop via
  `MOVQ $imm64 -> MOVQ AX, Xn -> VBROADCASTSD`, three uops on a serial dependency
  chain per coefficient per four elements. `ExpSIMDTo` was therefore **190x slower
  than the generic path already in the same package** and up to 41x slower than a
  plain `math.Exp` loop. The five transcendental dispatchers now use
  `parallelizeGeneric`, matching what arm64, wasm and the generic stub already
  did. Measured `ExpSIMDTo`: n=4096 872,098 ns -> 20,961 ns; n=65,536
  14,322,837 ns -> 97,848 ns; n=1,048,576 211,667,934 ns -> 1,051,251 ns.
  The assembly is retained behind a new `-tags emlasm` build tag
  (`internal/eml/simd_trans_asm_amd64.go`) rather than deleted, so it stays
  auditable while being reworked; `TestTranscendentalAsmKernelsAreWiredUp` keeps
  the symbols referenced.
- `SmallCutoff` was a flat 256, about 16x too low: the worker pool spawned one
  goroutine per core for as few as 256 elements. Measured the pool lost to a
  plain serial loop by 4.38x at n=256, 2.76x at n=512 and 1.87x at n=1024, and
  only won above ~4096. It is now `SmallWorkloadFactor * runtime.NumCPU()`
  (512 per core) and every inline-vs-pool guard routes through it. Worst ratio is
  now 1.03 at n=1024.
- `Pipeline.RunTo` performed two full-length copies via `Run`; it now writes the
  first step's output straight into `output` and ping-pongs between `output` and a
  single scratch buffer. 0 allocs at n=65,536.
- `EvalBatchColumnar` allocated 24,657 B and 2 objects per call, rebuilding the
  columnar stack and its slice-of-slices index every time. Added `BatchScratch`
  and `EvalBatchColumnarScratch`; measured 0 allocs and 14,543 ns -> 11,034 ns at
  n=1024. The 10 inner element loops were also converted to `range` form so the
  compiler can widen them.

### Known limitations
- `bytecode.Search` recovers shallow closed forms built from the available
  operators essentially exactly, but approximates transcendental targets rather
  than recovering them: fitting y = sin(x) over [0, 2pi) plateaus at RMS ~0.17
  (versus ~0.71 for the best constant fit) and does not improve with budget
  (800x600 gives 0.169, 1500x1500 gives 0.177). Approximating sin requires
  discovering a composition of elementary functions, which the source paper does
  with gradient-based training of a parameterised tree rather than a genetic
  algorithm. See docs/nextsteps.md.
- `EvalBatchColumnarScratch` removes the per-call allocation but its opcode bodies
  still run scalar loops over columns: 10.7 ns/sample against 33.6 ns for scalar
  `Eval`. Real widening needs generated kernels, since Go cannot auto-vectorise a
  call to `math.Sin`.
- The `-tags emlasm` path is expected to fail the accuracy gates. That is the
  contract a reworked kernel must meet before it can be re-enabled by default.


## [0.3.0] - 2026-09-08

### Added
- Core EML operator: eml(x, y) = exp(x) - log(y)
- Complex number support via math/cmplx
- SIMD dispatch for AMD64 (AVX2, AVX-512), ARM64 (NEON, SVE), WASM
- JIT compiler for amd64 with x86-64 machine code generation
- GPU backends: CUDA (Linux/Windows) and Metal (macOS/ARM64)
- Arithmetic package: Add, Sub, Mul, Div, Mod, Pow, Sqrt, Cbrt, Hypot, etc.
- Trig package: Sin, Cos, Tan, Cot, Sec, Csc, Asin, Acos, Atan, Atan2, etc.
- Hyperbolic package: Sinh, Cosh, Tanh, Asinh, Acosh, Atanh
- LogExp package: Exp, Log, ExpBatch, LogBatch, ExpFast, LogFast
- FastMath package: polynomial approximations for Exp, Sin, Cos, Log
- Worker pool for fused batch operations (ExpMul, ExpAdd, LogDiv, LogSub)
- Branchless Abs, Min, Max, Select operations
- Mathematical constants package (e, pi, ln2, sqrt2, phi, etc.)
- CLI tools: emlcli, bench, validate
- Comprehensive test suite with fuzz testing
- GPU result verification with ULP-based comparison
