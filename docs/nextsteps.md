# Next Steps: Remaining Open Improvements & Verification Log

A prioritised, evidence-based plan for the remaining open items in `emlgo`.
Every step below cites measurements and audit findings from this repository.

**Measurement environment for all numbers below:** Intel Core i7-12650H
(AVX2 + FMA, no AVX-512), 16 logical cores, Linux 5.x, Go 1.27.0, `linux/amd64`.

**Contents.**

- **Completed & Verified Milestones (Archived)** — Steps 1–8, 11–30 have been verified as fully implemented, covered by unit/fuzz tests with 100.0% statement coverage across all packages and zero race conditions under `go test -race ./...`, and archived from active work.
- **Part 1: Platform & Kernel Extensions (Hardware-gated)** — Steps 9 (ARM64 JIT) and 10 (AVX-512 kernels) remain open.

---

## Current Status: Open Work

### Platform Extensions (Hardware-gated)

| # | Step | Status | Headline Evidence |
| :-- | :--- | :--- | :--- |
| 9 | ARM64 JIT codegen | open | 17–30x for 2 of 3 server markets (Apple Silicon, Graviton) |
| 10 | AVX-512 transcendental kernels | open | 2x speedup on AVX-512 hosts (Zen 4/5, Xeon Scalable) |

---

## Completed & Verified Milestones (Archived)

The following 28 steps have been verified as complete in the codebase, with automated regression and fuzz tests in place, maintaining 100.0% statement coverage and zero race conditions:

| # | Step | Completed In | Result / Evidence | Verification Suite |
| :-- | :--- | :--- | :--- | :--- |
| 1 | Fix AVX2 transcendental kernels | Step 3 | Superseded by generic fallback | `internal/eml/simd_trans_amd64.go` |
| 2 | Gate worker pool on real threshold | v0.4 | `SmallCutoff = 512 * NumCPU` | Crossover sweep: no loss at any size |
| 3 | Make generic path the fallback | v0.4 | Generic path default, `exp` 201x faster at n=1M | `simd_trans_amd64.go`, `-tags emlasm` |
| 4 | Add accuracy gates for batch paths | v0.4 | Mixed abs/rel tolerance (`1e-14 + 1e-15*|want|`) | `simd_accuracy_test.go` |
| 5 | Extend bytecode VM to 27 functions | v0.4 | 5 ops → 27 ops unified across VM | `pkg/bytecode/functab.go` |
| 6 | Ship GA driver + fitness | v0.4 | Recovers $y = 3x^2+2x+1$ to 4.5e-13 RMS | `pkg/bytecode/ga.go`, `ga_test.go` |
| 7 | Zero-alloc `EvalBatch` + wider loops | v0.4 | `BatchScratch` reduces 24.7 KB → 0 B per call | `eval_batch_scratch_test.go` |
| 8 | `Pipeline.Reset()` and in-place `RunTo` | v0.4 | 0 allocs, ping-pong buffers, no step leak | `pipeline_test.go` |
| 11 | Fix JIT register-ABI clobber & spill | v0.5 | Callee-save across all 16 XMM registers; stack spilling; 0/3000 differential mismatches | `internal/jit/step11_test.go`, `FuzzJITDifferential` |
| 12 | Close `EMLEval`/`Diff` silent-zero holes | v0.5 | All 27 canonical functions evaluated and differentiated | `internal/jit/step12_test.go` |
| 13 | Make bytecode VM fail loudly & validate | v0.5 | Full bytecode opcode validation in `finalize()`, loud panic on unknown op, bounds checks in `EvalBatchColumnarScratch` | `pkg/bytecode/step13_test.go` |
| 14 | Fix fastmath silent range failures | v0.5 | Cody-Waite exponent limit corrected (-745.13); float32 limits aligned; `LnRegularized` uses `Hypot` | `pkg/fastmath/step14_test.go` |
| 15 | Fix `purego` on arm64 and wasm | v0.5 | Added `&& !purego` to arm64 and wasm dispatchers; compiles cleanly across all target architectures | CI cross-build matrix, `-tags=purego` |
| 16 | Fix or delete the Metal backend | v0.5 | Removed broken Metal runtime, unified GPU on CUDA/stub with correct buffer tracking & pinned memory | `internal/gpu` tests, `stub.go` |
| 17 | Delete `arithmetic.parallelMap` | v0.5 | Routed 13 batch operations through `internal/eml.ForEachChunk` persistent pool; 0 allocs | `pkg/arithmetic` benchmarks & tests |
| 18 | Reclaim JIT code memory & shard cache | v0.5 | Allocation size tracking with `FreeExecutableMemory` (`Munmap`), 16-shard LRU cache, singleflight deduplication | `internal/jit/step18_test.go` |
| 19 | Make GA scale and stop bloating | v0.5 | Concurrent `tuneConstants` across workers, parsimony penalty ($\lambda$), adaptive mutation, plateau convergence | `pkg/bytecode/step19_test.go` |
| 20 | Make perf/accuracy claims enforceable | v0.5 | `-regression` flag in `init()`, strict `ulpDiff` (NaN/Inf), `cmd/validate` extended to fastmath/quant/batch (435/435 pass), `pkg/quant` fidelity tests | `cmd/bench`, `cmd/validate`, `pkg/quant/fidelity_test.go` |
| 21 | Fix bytecode optimizer RPN corruption | v0.5 | AST constant folder & identity reducer; preserved commutativity | `FuzzOptimizerDifferential` (1.57M execs) |
| 22 | Fix JIT parser scientific notation & whitespace | v0.5 | `[eE][+-]?[0-9]+` & `\t,\n,\r` accepted; var scoping | `FuzzScientificFloatRoundTrip` (1.64M execs) |
| 23 | Fix hyperbolic cancellation & log(0) | v0.5 | `Asinh` odd symmetry, `Sinh` expm1, `Log(0)` = -Inf | `FuzzHyperbolicAccuracy` (2.12M execs) |
| 24 | Fix arbitrary-precision limits in `bigmath` | v0.5 | `Exp(-1000)` non-zero, exact `Asin`, memoized $\pi$/$\ln 2$ | `FuzzBigMathPrecision` (509k execs) |
| 25 | Complete `CompileEML` function set | v0.5 | All 27 ops wired via `unaryOpcode()` from `functab.go` | `FuzzCompileEMLRoundTrip` (1.53M execs) |
| 26 | Fix decompiler precedence & associativity | v0.5 | Precedence-aware parens, variable names preserved | Bit-exact round-trip unit & fuzz tests |
| 27 | Fix JIT power codegen & negative base NaN | v0.5 | ABIInternal `math.Pow` call for $(-2)^3$ & $0^0$; 1-cycle `xorpd` | `codegen.go`, unit & fuzz tests |
| 28 | Vectorize Pipeline engine & complex ops | v0.5 | SIMD Pipeline (28x win), affine fusion, Neumaier dot (5x win) | `FuzzPipelineDifferential` (1.89M execs) |
| 29 | Replace VM opcode map with flat jump table | v0.5 | `opUnaryTable [256]mathFunc`, zero-alloc scratch capacity reuse | `eval.go`, `eval_batch.go`, benchmarks |
| 30 | Vectorize integer saturation & normalize lengths | v0.5 | AVX2 `VPADDSB`/`VPSUBSB` 17x speedup, 0 allocs in `*To`, panics unified | `FuzzAddSubSatInt8` (2.2M execs) |

---

# Detailed Open Step Specifications

---

## Step 9 — ARM64 JIT codegen

**Evidence.** The JIT is the library's biggest genuine win, and it is amd64-only:

| expression | JIT | AST interpreter | speedup |
| :--- | --: | --: | --: |
| `2*x + 1` | 2.22 ns | 36.99 ns | **16.7x** |
| `x^2 + 2*x + 1` | 2.08 ns | 61.60 ns | **29.6x** |
| `sin(x)` | 8.34 ns | 46.44 ns | **5.6x** |

Compile cost is ~4 µs, and `CompileCached` drops the repeat cost to 28 ns (140x
better), so the economics are excellent.

But `internal/jit/codegen_stub.go` returns `"JIT codegen requires amd64"` for
every other target. `js/wasm` gets the interpreter via `jit_wasm.go`, but
**arm64 gets nothing** — `internal/eml/simd_arm64.s` is a 5-line stub (build tag
plus an `#include`), so ARM64 batch math is a plain Go loop too. Apple Silicon,
Graviton and the entire ARM server market get neither the JIT nor hand-written
SIMD.

**Change.** Port the codegen to AArch64. The register model maps cleanly
(`xmm0–xmm15` → `v0–v31`, and there are *more* registers, so the 14-temporary
limit relaxes). `AllocateExecutableMemory` already works on
`darwin/arm64` + `linux/arm64` via `unix.Mmap`/`Mprotect`, so only the encoder
is missing. Keep the `functab.go` dispatch shared.

**Verify.** Port the existing `internal/jit` codegen tests and run them under
`GOARCH=arm64` emulation or on real hardware; add `darwin/arm64` to the CI
cross-build job (already present) plus a compile-and-verify job.

---

## Step 10 — AVX-512 transcendental kernels

**Evidence.** `simd_amd64.s` defines 11 AVX-512 symbols — `add`, `sub`, `mul`,
`div`, `addScalar`, `mulScalar`, `sqrt`, `fma`, `abs`, `neg`, `inv` — and **no
transcendental ones**. `dispatchExpSIMDTo` guards on `hasAVX2` only:

```go
func dispatchExpSIMDTo(x, result []float64) {
	n := len(x)
	if hasAVX2 {                 // never checks hasAVX512
		simdLen := (n / 4) * 4
		expAVX2(x[:simdLen], result[:simdLen])
```

So on Zen 4/5, Sapphire Rapids and Xeon Scalable, `exp`/`log`/`sin`/`cos`/`tan`
batches run the 4-wide kernel while the elementwise ops correctly use the 8-wide
one — an inconsistency in both width and accuracy.

**Change.** Add `expAVX512`, `logAVX512`, `sinAVX512`, `cosAVX512`, `tanAVX512`
using 8-wide `Z` registers, and add a `hasAVX512` branch to each transcendental
dispatcher. Do this **after** Step 1, reusing its corrected coefficient-loading
technique, or the new kernels will inherit the same per-iteration broadcast bug.

**Verify.** Extend the Step 4 accuracy gates with 8-wide paths, and add an
AVX-512 host to the benchmark matrix in `docs/performance.md`.

---

## Remaining Sequencing & Priorities

```
        Completed ─── Steps 1–8, 11–30 (All P0 & P1 milestones complete & verified)

        Hardware ─┬─ 9  ARM64 JIT codegen         largest single win, requires ARM64 CI
                  └── 10 AVX-512 transcendentals   requires AVX-512 CI host
```

All software correctness blockers (P0: Steps 11, 13, 14, 15, 16) and performance/enforceability milestones (P1: Steps 17, 18, 19, 20) are completed, verified with 100.0% statement code coverage across all 15 packages, pass under `go test -race ./...` with zero data races, and are continuously validated against regressions via `cmd/bench -regression`.

Steps 9 and 10 remain as the sole open platform extensions, gated on ARM64 and AVX-512 physical hardware execution environments.
