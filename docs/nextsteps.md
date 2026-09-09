# Next Steps: Improvement Plan & Known Issues

This document catalogs planned improvements and the current state of the `emlgo` math library.
Plans 1–2, 4–10 have been fully implemented and verified in the codebase. Plan 3 (ARM64
NEON/SVE2 Assembly) requires ARM64 hardware access to complete.

---

## 10-Part Improvement Plan

> **Research context:** The EML operator (`eml(x,y) = exp(x) − ln(y)`) was introduced by
> Odrzywołek (arXiv:2603.21852v2, 2026). It is functionally complete over elementary functions
> given the constant `1` — analogous to NAND in Boolean logic. The canonical grammar is
> `S → 1 | eml(S,S)`. This library implements the operator, its batched/SIMD variants, JIT
> compilation of expression trees, an arbitrary-precision backend (`bigmath`), and a canonical
> EML tree system. The improvements below are grounded in both EML theory and practical Go math
> library best practices (cf. Gonum, Gorgonia).

### Plan 1 — EML Gradient / Symbolic Differentiation Engine ✅

**Status:** Implemented in `internal/jit/canonical.go:303-390`.

The `Diff(n *EMLNode) *EMLNode` function performs symbolic differentiation using the chain rule:
- `d/dx eml(u, v) = exp(u)·u' − v'/v`
- Supports all standard functions: exp, log, sin, cos, sqrt, neg, add, sub, mul, div, pow
- Includes `Simplify()` integration for optimized derivative expressions

---

### Plan 2 — EML Tree Complexity Reduction & Constant Folding ✅

**Status:** Implemented in `internal/jit/canonical.go:205-301`.

The `Simplify(n *EMLNode) *EMLNode` function performs:
- Constant folding: evaluates subtrees where all leaves are `EMLConst`
- Identity reduction: `eml(x, 1)` → `exp(x)` direct node
- Algebraic simplifications: 0+x=x, 1*x=x, 0*x=0, x^0=1, x^1=x, etc.
- Double negation elimination: `neg(neg(x)) = x`

---

### Plan 3 — ARM64 NEON/SVE2 Assembly Kernels ⏳

**Status:** Pending — requires ARM64 hardware (Apple M-series, Ampere, or AWS Graviton CI runner).

`sVE2/simd_arm64.s` contains only the `textflag.h` header (3 bytes). Pure-Go stubs exist in
`sVE2/simd_arm64.go` and `sVE2/simd_sve.go` but compile to scalar loops, not SIMD instructions.

**Remaining work:**
- Write `simd_arm64.s` with Plan 9 assembly for NEON `float64x2` lanes
- Implement `detectSVE()` with `getauxval(AT_HWCAP)` via `golang.org/x/sys/unix`
- Requires Apple M-series or Ampere/Graviton CI runner

---

### Plan 4 — `StopWorkerPool` Concurrency Safety via `sync.Once` ✅

**Status:** Implemented in `internal/eml/simd.go:408-426`.

```go
var stopOnce sync.Once

func StopWorkerPool() {
    stopOnce.Do(func() { close(jobQueue) })
}
```

Replaces the old `workerPoolStopped bool` guard. Concurrent calls are safe.

---

### Plan 5 — JIT Parser: Multi-Variable Support & Structured Errors ✅

**Status:** Implemented in `internal/jit/parser.go:113-175`.

- `ParseWithVars(input string, vars []string) (Node, error)` — multi-variable support
- `ParseError` struct with `Pos`, `Token`, and `Msg` fields
- `round` added to `funcTable` in `codegen.go:37`
- `FuzzParse` and `FuzzEval` fuzz tests in `fuzz_test.go`

---

### Plan 6 — `bigmath` Complete Transcendental Coverage ✅

**Status:** Implemented in `internal/eml/bigmath/bigmath.go:393-516`.

- `Tan(x) = Sin(x)/Cos(x)` with division-by-zero guard
- `Atan(x)` via half-angle reduction + reciprocal identity
- `Asin(x) = Atan(x / Sqrt(1−x²))`
- `Acos(x) = π/2 − Asin(x)`
- `NewFloatFromInt(n int64)` constructor

---

### Plan 7 — Composable Zero-Allocation `Pipeline` API ✅

**Status:** Implemented in `internal/eml/pipeline.go`.

Double-buffered composable pipeline with 9 operations: Exp, Log, Sqrt, Sin, Cos, Abs, Neg, MulScalar, AddScalar. Uses buffer swapping to avoid allocations per step.

```go
p := eml.NewPipeline(n)
p.Exp().MulScalar(2.0).Log().RunTo(input, output)
```

---

### Plan 8 — `float32` SIMD Batch Support ✅

**Status:** Implemented in `internal/eml/simd_f32.go`.

- `ExpSIMDF32`, `LogSIMDF32`, `SqrtSIMDF32` — transcendentals via float64 upcast
- `AddSIMDF32`, `SubSIMDF32`, `MulSIMDF32`, `DivSIMDF32` — arithmetic
- `AbsSIMDF32`, `NegSIMDF32`, `InvSIMDF32` — unary operations
- `SinSIMDF32`, `CosSIMDF32`, `TanSIMDF32` — trigonometric
- `AddScalarSIMDF32`, `MulScalarSIMDF32` — scalar operations

---

### Plan 9 — JIT Expression Cache & `CompileCached` ✅

**Status:** Implemented in `internal/jit/cache.go`.

- Mutex-protected LRU cache with `container/list` (max 1024 entries)
- `CompileCached(expr string) (Func, error)` — returns cached result
- `ClearJITCache()` for long-running programs

---

### Plan 10 — EML → Infix Decompiler & LaTeX Emitter ✅

**Status:** Implemented in `internal/jit/decompile.go`.

- `Decompile(n *EMLNode) string` — parenthesized infix output
- `DecompileLaTeX(n *EMLNode) string` — LaTeX math mode output
- `DecompileNodeToExpr(n *EMLNode) string` — uses JIT formatter
- `--decompile` flag added to `cmd/emlcli`

---

## Pending ARM64 Platform Work

The following items require ARM64 hardware (Apple M-series, Ampere, or AWS Graviton CI runner):

### ARM64-1. NEON Assembly Kernels
* **Status:** Pending — `simd_arm64.s` is empty; pure-Go stubs exist in `simd_arm64.go`.

### ARM64-2. SVE/SVE2 Assembly Kernels
* **Status:** Pending — `simd_sve.go` has Go-level SVE simulation; real SVE requires ARMv8.2-A+.

### ARM64-3. ARM64 Transcendental Batch Kernels
* **Status:** Pending — No `expNEON`/`logNEON`/`sinNEON` assembly; falls back to `parallelizeGeneric`.
