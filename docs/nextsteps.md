# Next Steps: Improvement Plan & Known Issues

This document catalogs planned improvements and the current state of the `emlgo` math library.
All P0 blockers, P1 bugs, P2 issues, and original Improvement Plan items 4–10 have been
fully implemented and verified in the codebase. The pending ARM64 platform items require
hardware access to complete.

---

## 10-Part Improvement Plan

> **Research context:** The EML operator (`eml(x,y) = exp(x) − ln(y)`) was introduced by
> Odrzywołek (arXiv:2603.21852v2, 2026). It is functionally complete over elementary functions
> given the constant `1` — analogous to NAND in Boolean logic. The canonical grammar is
> `S → 1 | eml(S,S)`. This library implements the operator, its batched/SIMD variants, JIT
> compilation of expression trees, an arbitrary-precision backend (`bigmath`), and a canonical
> EML tree system. The improvements below are grounded in both EML theory and practical Go math
> library best practices (cf. Gonum, Gorgonia).

### Plan 1 — EML Gradient / Symbolic Differentiation Engine

**Background:** Odrzywołek's paper explicitly motivates EML trees for gradient-based symbolic
regression (Adam optimizer on tree weights). Currently `emlgo` has no differentiation capability
despite having the canonical tree representation in `internal/jit/canonical.go`.

**Proposed work:**
- Implement `Diff(n *EMLNode, varName string) *EMLNode` using the chain rule:
  - `d/dx eml(u, v) = exp(u)·u' − v'/v`
  - `d/dx eml(x, 1) = exp(x)` (the Exp canonical form)
- Add `AutoDiff` forward-mode accumulation for batch `[]float64` differentiation.
- Test: `Diff(CanonicalExp(varNode()), "x")` equals `CanonicalExp(varNode())`.

**Priority:** High — unlocks the library's stated symbolic regression goal.

---

### Plan 2 — EML Tree Complexity Reduction & Constant Folding

**Background:** Deep EML trees for simple operations are expensive: `ln(x)` requires 3 nested
`eml` calls. `Canonicalize` and `EMLEval` perform no simplification. Constant subtrees are
re-evaluated on every call.

**Proposed work:**
- Implement `Simplify(n *EMLNode) *EMLNode` in `internal/jit/canonical.go`:
  - Constant folding: evaluate subtrees where all leaves are `EMLConst`.
  - Identity reduction: `eml(x, 1)` → `exp(x)` direct node when subtree already simplified.
  - CSE via structural hashing with the existing `Equiv()` function.
- Add `Depth(n *EMLNode) int` metric.
- Property test: `EMLEval(Simplify(n), x) == EMLEval(n, x)` for all x.

**Priority:** High — directly improves JIT and interpreter throughput.

---

### Plan 3 — ARM64 NEON/SVE2 Assembly Kernels

**Background:** `simd_arm64.go` contains pure-Go stubs named `addNEON`, `addSVE`, etc., but
`simd_arm64.s` is essentially empty (24 bytes). AMD64 has full AVX2/AVX512 assembly. ARM64
should achieve parity for the Apple M-series / Graviton cloud market.

**Proposed work:**
- Write `simd_arm64.s` with Plan 9 assembly for NEON `float64x2` lanes:
  - `VFADD`, `VFSUB`, `VFMUL`, `VFDIV`, `FSQRT`, `FABS`, `FNEG` intrinsics.
- Implement `detectSVE()` with `getauxval(AT_HWCAP)` via `golang.org/x/sys/unix`.
- Requires Apple M-series or Ampere/Graviton CI runner (see pending section).

**Priority:** Medium — hardware-gated.

---

### Plan 4 — `StopWorkerPool` Concurrency Safety via `sync.Once`

**Background:** `StopWorkerPool()` in `internal/eml/simd.go` uses a plain `bool`
(`workerPoolStopped`) as a guard against double-close. Two simultaneous callers could both
pass the guard and trigger `close` on an already-closed channel, causing a panic. The race
detector passes today only because no tests call `StopWorkerPool()` concurrently.

**Proposed work:**
- Replace `workerPoolStopped bool` with a `sync.Once`:
  ```go
  var stopOnce sync.Once
  func StopWorkerPool() {
      stopOnce.Do(func() { close(jobQueue) })
  }
  ```
- Remove the `workerPoolStopped` package-level variable.
- Add a concurrent stress test: 10 goroutines call `StopWorkerPool()` simultaneously.

**Priority:** High — eliminates a latent double-close panic.

---

### Plan 5 — JIT Parser: Multi-Variable Support & Structured Errors

**Background:** `internal/jit/parser.go` parses single-variable expressions (variable `x`
only). The `funcTable` is complete (16 functions), but `round` is absent despite
`math.Round` being available. Parser errors are unstructured strings.

**Proposed work:**
- Extend `Parse()` to accept `vars []string`, enabling `"x*y + sin(z)"`.
- Define `ParseError{Line, Column int; Token string; Msg string}` implementing `error`.
- Add `round` to `funcTable` in `codegen.go`.
- Add fuzz corpus: `go test -fuzz=FuzzParse` targeting panics on malformed input.

**Priority:** Medium — broadens practical JIT usability.

---

### Plan 6 — `bigmath` Complete Transcendental Coverage

**Background:** `internal/eml/bigmath/bigmath.go` provides `Exp`, `Log`, `Sin`, `Cos`,
`Sqrt` at arbitrary precision. `Tan`, `Atan`, `Atan2`, `Asin`, and `Acos` are missing.
The current `reduceTrig` for `Sin` has accuracy loss for large arguments (reduces modulo 2π
rather than the more stable Crandall–Payne method).

**Proposed work:**
- Add `Tan(x) = Sin(x)/Cos(x)` with division-by-zero guard.
- Add `Atan(x)` via Machin-like series (same infrastructure as `machinPi`).
- Add `Asin(x) = Atan(x / Sqrt(1−x²))` and `Acos(x) = π/2 − Asin(x)`.
- Improve `reduceTrig` argument reduction to use extended-precision π multiples.
- Add `NewFloatFromInt(n int64)` constructor.

**Priority:** Medium — completes bigmath transcendental coverage.

---

### Plan 7 — Composable Zero-Allocation `Pipeline` API

**Background:** Batch operations exist as independent functions (`ExpSIMD`, `LogSIMD`, etc.)
but chaining them allocates O(n) intermediate slices per operation. There is no composable
pipeline for operator fusion without allocation.

**Proposed work:**
- Define a `Pipeline` type with double-buffered scratch slices:
  ```go
  p := eml.NewPipeline(n)
  p.ExpSIMD(input).MulScalar(2.0).LogSIMD().WriteTo(output)
  ```
- Implement as a slice of `op func(src, dst []float64)` composed at build time.
- Add benchmark: `Pipeline` vs chained individual calls (expected >30% reduction in allocs).

**Priority:** Medium — reduces GC pressure in ML inference workloads.

---

### Plan 8 — `float32` SIMD Batch Support

**Background:** All SIMD and batch operations are `float64`-only. Modern ML and graphics
workloads use `float32` heavily. The existing AVX2/AVX512 assembly already supports 8 or 16
`float32` lanes per instruction (vs 4/8 for `float64`), so throughput can nearly double.

**Proposed work:**
- Add `float32` variants: `ExpSIMDF32`, `LogSIMDF32`, `AddSIMDF32`, `MulSIMDF32`.
- Extend `simd_dispatch_amd64.go` with an `f32` dispatch layer using `vaddps` / `vmulps`.
- Add `Float32Batch` matching the existing `Batch` API signature.
- Benchmark: `float32` vs `float64` throughput on 1M-element arrays.

**Priority:** Medium — significant practical demand in ML/graphics.

---

### Plan 9 — JIT Expression Cache & `CompileCached`

**Background:** `internal/jit/jit.go` compiles expressions on every call. In expression-heavy
workloads (e.g., evaluating the same formula over a large dataset), the JIT compilation cost
(assembly encoding, mmap, fixup) is incurred repeatedly for identical expressions.

**Proposed work:**
- Add a `sync.Map`-backed cache keyed by the canonical expression string in `jit.go`.
- Add `CompileCached(expr string) (func(float64) float64, error)` — returns cached result.
- Add `ClearJITCache()` for long-running programs.
- Cap at `maxCacheEntries = 1024` with LRU eviction to prevent unbounded growth.
- Benchmark: cache hit vs cold compile on a repeated 1000-call workload.

**Priority:** Medium — eliminates repeated JIT overhead.

---

### Plan 10 — EML → Infix Decompiler & LaTeX Emitter

**Background:** `canonical.go` builds EML trees from expressions and `EMLEval` evaluates
them, but there is no inverse: given an `EMLNode`, produce a human-readable expression.
This completes the EML toolchain as a bidirectional compiler, enabling use in symbolic
regression output and documentation generation.

**Proposed work:**
- Implement `Decompile(n *EMLNode) string` — parenthesized infix (e.g. `"exp(x) - log(y)"`).
- Implement `DecompileLaTeX(n *EMLNode) string` — LaTeX math mode output.
- Add `--decompile` flag to `cmd/emlcli` that reads an EML tree JSON and prints the formula.
- Roundtrip property test: `Parse(Decompile(Canonicalize(Parse(expr)))) ≡ Parse(expr)`.

**Priority:** Low — completes the EML toolchain; primarily a developer/research tool.

---

## Pending ARM64 Platform Work

The following items require ARM64 hardware (Apple M-series, Ampere, or AWS Graviton CI runner):

### ARM64-1. NEON Assembly Kernels
* **Status:** Pending — `simd_arm64.s` is empty; pure-Go stubs exist in `simd_arm64.go`.

### ARM64-2. SVE/SVE2 Assembly Kernels
* **Status:** Pending — `simd_sve.go` has Go-level SVE simulation; real SVE requires ARMv8.2-A+.

### ARM64-3. ARM64 Transcendental Batch Kernels
* **Status:** Pending — No `expNEON`/`logNEON`/`sinNEON` assembly; falls back to `parallelizeGeneric`.

---
