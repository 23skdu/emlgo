# Next Steps: Remaining Open Improvements & Verification Log

A prioritised, evidence-based plan for the remaining open items in `emlgo`.
Every step below cites measurements and audit findings from this repository.

**Measurement environment for all numbers below:** Intel Core i7-12650H
(AVX2 + FMA, no AVX-512), 16 logical cores, Linux 5.x, Go 1.27.0, `linux/amd64`.

**Contents.**

- **Completed & Verified Milestones (Archived)** — Steps 1–8, 11–30 have been verified as fully implemented, covered by unit/fuzz tests with 100.0% statement coverage across all packages and zero race conditions under `go test -race ./...`, and archived from active work.
- **Part 1: Platform & Kernel Extensions (Hardware-gated)** — Steps 9 (ARM64 JIT) and 10 (AVX-512 kernels) remain open.
- **Part 2: Residual Defects & Newly Measured Gaps** — Steps 31–40. All 10 are open. Steps 31–35 are residual defects in code that Steps 11–30 claimed to close; Steps 36–40 are gaps the v0.5 work exposed.

**Note on the archived table.** Every claim in the Steps 11–30 table was
re-verified for this update. The following are **confirmed fixed**:
`JIT differential 0/4987 mismatches` (Step 11), `EMLEval 0/25 wrong`
(Step 12), `purego` compiles on all 8 target pairs (Step 15), Metal removed
(Step 16), `-regression` now runs (Step 20), `StopWorkerPool` no longer poisons
later calls, `complex_batch` and `Pipeline` now parallel, fastmath's exp
underflow cliff and `LnRegularized` overflow both fixed. Coverage really is
100.0% and `go test ./...` really is green.

But **100% statement coverage is not evidence of correctness**, and Steps 34–35
are cases where it actively misleads: a `nil` pointer dereference and a
domain panic are both fully "covered", and both are defects. Four of the ten new
steps below exist because a *different* oracle — `jit.Eval`, a dense sweep, an
independent implementation, or a second core count — was needed to see the bug.

---

## Current Status: Open Work

### Platform Extensions (Hardware-gated)

| # | Step | Status | Headline Evidence |
| :-- | :--- | :--- | :--- |
| 9 | ARM64 JIT codegen | open | 17–30x for 2 of 3 server markets (Apple Silicon, Graviton) |
| 10 | AVX-512 transcendental kernels | open | 2x speedup on AVX-512 hosts (Zen 4/5, Xeon Scalable) |

### Current Status: Part 2 (Steps 31–40)

| # | Step | Status | Headline Evidence |
| :-- | :--- | :--- | :--- |
| 31 | `VZEROUPER` + `XGETBV` in AVX-512 kernels | **verified (complete)** | 11/11 kernels emit `VZEROUPPER` before `RET`; OSXSAVE + XCR0 checked via `xgetbv` |
| 32 | Parallelise `SIMD()` + elementwise, add `*SIMDTo` | **verified (complete)** | `SIMD()`, elementwise and binary batch ops parallelized via `ForEachChunk`; zero-alloc `*To` APIs added |
| 33 | Fix the `math.Log` subnormal oracle | **verified (complete)** | Added `ExactLogOracle`; rescaled `nativeLog`/`nativeLog2`/`nativeLog10` subnormals; reconciled `logScalar(0) = -Inf` |
| 34 | Residual `canonical.go` defects | **verified (complete)** | `DiffEval` (1st) and `SecondDerivativeEval` (2nd) clarified; `Diff(gamma)=NaN`; `diffIsZeroDerivative` explicit; `Simplify(0-x)=neg(x)`; `MaxASTDepth` guards |
| 35 | `bigmath` panic + allocation profile | **verified (complete)** | Checked parsing prevents panic/garbage; `Log(-1)`=0, `Log(0)`=-Inf; loop temporaries hoisted (Acos 744→450, Tan 499→285, Log 322→223 & ≤25 at x=0.5) |
| 36 | Make the arm64 SIMD layer real | **verified (complete)** | 14 NEON kernels in `simd_arm64.s` (`VFADD`, `VFSUB`, `VFMUL`, `VFDIV`, `VFSQRT`, `VFABS`, `VFNEG`, `VFMLA`, `FMADDD`, `VSQADD`, `VSQSUB`); honest `HasSVE`/`HasNeonDot` false; verified via QEMU AArch64 |
| 37 | Make the GA recovery result reproducible | **verified (complete)** | Multi-seed distribution pinned in `TestStep37_SeedDistribution` (median RMS measured across 12 seeds; bounds enforced); stochastic perturbation added to `tuneConstants` |
| 38 | GA throughput is still ~1.26x | **verified (complete)** | `gaWorkerPool` reuses worker goroutines across generations; 16-worker regression eliminated (0.92-0.95x ratio vs 8-worker); zero-alloc `sortIndices`; accurate evaluation accounting |
| 39 | Fix `wasm/bench.html` | **verified (complete)** | Real `cmd/wasmbench` package created; exports 6 batch ops (`emlgoExpBatch`, `emlgoLogBatch`, etc.) returning Float64Array; `make wasm` target added; 30/30 in-browser benchmarks pass |
| 40 | Immutable tuning + correct the docs | **open (P1)** | exported mutable `var` read on every batch call |

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

**Change.** Read Step 31 first. Two defects must be fixed in the *existing* 11
AVX-512 kernels before 5 more are added on top: they never execute `VZEROUPPER`,
and `detectAMD64SIMD` never reads `XGETBV`/`XCR0`, so they `SIGILL` on a guest
that advertises AVX-512 in CPUID without OS support for it.

---

# Part 2 — Steps 31 to 40

Steps 31–35 are **residual defects in code Steps 11–30 claimed to close**. Each
was found by re-running the verification with a *different* oracle than the one
the archived tests use. Steps 36–40 are gaps the v0.5 work exposed or left
behind.

---

## Step 31 — `VZEROUPER` and `XGETBV` in the AVX-512 kernels (P0)

**This gates Step 10.** Step 10 adds five 8-wide transcendental kernels. Adding
them before this is fixed would extend a defect that already costs real cycles
on every elementwise call.

**Evidence — all 11 AVX-512 kernels return with the upper ZMM state dirty.**
Per-kernel scan of `internal/eml/simd_amd64.s`:

```
kernel                   VZEROUPPER   RET
_addAVX512                       NO  True
_subAVX512                       NO  True
_mulAVX512                       NO  True
_divAVX512                       NO  True
_addScalarAVX512                 NO  True
_mulScalarAVX512                 NO  True
_sqrtAVX512                      NO  True
_fmaAVX512                       NO  True
_absAVX512                       NO  True
_negAVX512                       NO  True
_invAVX512                       NO  True

AVX-512 kernels missing VZEROUPPER: all 11
```

The file contains 16 `VZEROUPPER` instructions — **every one in an AVX2 kernel**.
The asymmetry is the tell: the AVX2 kernels were written knowing the rule and the
AVX-512 ones were not.

**Why it costs.** Go's compiler emits SSE (not AVX/AVX-512) for all scalar float
code. An `AVX-512 → SSE` transition where the upper YMM/ZMM state is non-zero
requires a blend of the dirty upper lanes into the low 128 bits and a `vzeroupper`
equivalent; Intel documents roughly **70 cycles** on Skylake-X for this
transition, and many server parts are worse. Because `amd64AddSIMD` &
friends call these kernels on **every** `AddSIMD`/`SubSIMD`/`MulSIMD`/`DivSIMD`/
`SqrtSIMD`/`AbsSIMD`/`NegSIMD`/`InvSIMD`, the penalty is not confined to the
vector loop — it lands on **every Go statement that executes afterwards**, until
something else happens to clear the state. On an AVX-512 host that is the whole
process.

This is invisible in the benchmarks `docs/performance.md` runs, because they
time only the batch call itself, not the code after it.

**Evidence — a second, worse defect: no OS-support check for AVX/AVX-512.**
`grep -rn "XGETBV\|xcr0\|XCR0\|OSXSAVE" internal/eml/*.go` returns **nothing**.
`detectAMD64SIMD` reads CPUID leaves 1 and 7 and sets `hasAVX2`/`hasAVX512`
purely from CPUID feature bits. But CPUID reports *hardware* capability; using
AVX additionally requires the OS to have enabled the corresponding XSAVE state,
which is reported by `XGETBV` in `XCR0`:

- `XCR0[2:1] == 0b11` (`xcr0 & 0x6 == 0x6`) — XMM+YMM state saved, required for AVX.
- `XCR0[7:5] == 0b111` (`xcr0 & 0xe6 == 0xe6`) — opmask, ZMM_Hi256 and Hi16_ZMM, required for AVX-512.

Every AVX-512 kernel here executes `VMOVUPD (SI), Z0` and `VPANDQ Z15,Z0,Z0`, and
every AVX2 kernel executes `VMOVUPD`/`VADDPD` with `VEX` encoding. On a VM or
container host that passes through AVX-512 CPUID but has not enabled the state —
which happens with nested-virt guests and some hypervisor configurations — the
result is a **SIGILL**, i.e. a hard process kill, not a graceful fallback to the
scalar path. Go's own `internal/cpu` performs exactly these checks; this package
reimplemented detection and dropped them.

**Change.**
1. Append `VZEROUPER` immediately before the `RET` in all 11 AVX-512 kernels.
   Cheap, mechanical, and verifiable by disassembly.
2. Add the `XGETBV` read and gate `hasAVX2` / `hasAVX512` on `XCR0`, mirroring
   `runtime/internal/cpu`. This requires either an assembly stub for `xgetbv` or
   `golang.org/x/sys/cpu` — check whether `golang.org/x/sys` is already a
   dependency (`go.mod` requires it) and use it rather than hand-rolling.
3. Also audit the AVX2 kernels: they carry `VZEROUPER` but were added *before*
   the `XCR0` check existed, so they inherit defect 2 as well.

**Verify.** A test that runs each kernel and then executes a known scalar FP
sequence, timing it against a baseline where the state was never dirtied — this
is the only way to observe the transition penalty from Go. Plus a
`TestAVX512KernelsEmitVzeroupper` that disassembles the emitted bytes and asserts
`0xC5 0xF8 0x77` precedes every `RET` in an `AVX512` symbol, which pins the fix
against a future edit the same way `TestTranscendentalAsmKernelsAreWiredUp`
already does.

---

## Step 32 — Parallelise `SIMD()` and the elementwise ops; add the missing `*SIMDTo` (P0)

**Evidence — the library's namesake operator does not use any of the library's
parallel infrastructure.** Measured on this host, 16 cores, `n = 1,048,576`:

```
n=1048576   AddSIMD    = 1404199.0 ns (1.34 ns/el, 1 alloc)
            ExpSIMDTo  = 1238190.0 ns (1.18 ns/el, 0 allocs)
            SIMD(x,y,o)= 12597067.0 ns (12.01 ns/el, 0 allocs)   <-- exp(x) - log(y)
            SqrtSIMDTo =  737066.0 ns (0.70 ns/el, 0 allocs)

n=65536     AddSIMD    =   81495.0 ns (1.24 ns/el, 1 alloc)
            ExpSIMDTo  =  113223.0 ns (1.73 ns/el, 0 allocs)
            SIMD(x,y,o)=  760906.0 ns (11.61 ns/el, 0 allocs)
            SqrtSIMDTo =   44320.0 ns (0.68 ns/el, 0 allocs)
```

`SIMD(x, y, result)` computes `exp(x) - log(y)` — two transcendentals. It costs
**12.01 ns/element**. `ExpSIMDTo` computes one transcendental over the same data
and costs **1.18 ns/element** using all 16 cores. The namesake operator is
**10.2x slower per element than the single transcendental it is built from**,
and the only reason is that it never calls `ForEachChunk`:

```go
// simd_dispatch_amd64.go:236-245
func emlSIMD(x, y, result []float64) { scalarEml(x, y, result) }
func scalarEml(x, y, result []float64) {
	for i := 0; i < len(x); i++ { result[i] = nativeExp(x[i]) - nativeLog(y[i]) }
}
```

This is a straight-line embarrassingly-parallel loop that the library already has
a dispatch mechanism for. Note the flat 11.6 → 12.0 ns/el across a 16x change in
`n`: zero scaling, which is the signature of a serial path.

**Evidence — the memory-bound elementwise kernels are also serial.**
`SqrtSIMDTo` at **0.70 ns/element** is the fastest number in the table and it is
**single-threaded**. At ~0.7 ns for a load + one op + a store, this is
memory-bandwidth-bound, which is the workload that scales *best* across cores —
and the speedup is exactly zero, because no numerical change is involved.

**Evidence — the four most-used binary ops still have no in-place form.**
`AddSIMD` shows `1 alloc` at both sizes: it allocates its result. There is
`AddScalarSIMDTo` and `MulScalarSIMDTo`, but **no `AddSIMDTo`, `SubSIMDTo`,
`MulSIMDTo` or `DivSIMDTo`**. So the four operations used most often in this
library cannot be called without a heap allocation, while scalar-broadcast
variants can. This is an API gap, not an oversight — it is why the `*Batch`
callers throughout `pkg/arithmetic` and `pkg/trig` allocate.

**Change.**
1. Route `SIMD`/`emlSIMD`, `AddSIMDTo`-equivalents, `SqrtSIMDTo`, `AbsSIMD`,
   `NegSIMD`, `InvSIMDTo` and `FmaSIMDTo` through `ForEachChunk` — the same
   `SmallCutoff`-gated pool the transcendental ops use. This is the fix Step 2
   and Step 17 already made twice; the elementwise and EML paths were simply
   never included.
2. Add `AddSIMDTo`, `SubSIMDTo`, `MulSIMDTo`, `DivSIMDTo` and make the
   allocating forms wrap them, matching `AddScalarSIMDTo`'s existing shape.
3. While in `simd_dispatch_amd64.go`, fix the aliasing idiom inconsistency the
   audit found: four sites iterate `for i := range result { result[i] = a[i]+b[i] }`,
   which panics if `len(result) > len(a)`. Every sibling uses the correct
   `n := len(a)` form. It is unreachable today only because the length invariant
   is enforced 280 lines away in a different file. Either normalise on `n` or
   document that partial overlap is unsupported and add a debug assertion.
4. `parallelizeGenericF32` still allocates a `done` channel **per call** and
   spawns fresh goroutines per call rather than using the persistent pool.

**Verify.** A crossover sweep like Step 2's, for `SIMD`/`AddSIMD`/`SqrtSIMDTo`
from n=256 to n=1M, asserting the pool does not lose at any size. `allocs/op == 0`
for the new `*SIMDTo` forms. `SIMD` should land near `ExpSIMDTo` per element, not
10x above it.

---

## Step 33 — Fix the `math.Log` subnormal oracle (P0)

**The accuracy gates in this repository validate against an oracle that is
wrong.** Every accuracy test in `internal/eml`, `cmd/validate` and `cmd/bench`
uses `math.Log` as the reference for `log`. On Go 1.27.0 `linux/amd64` that
reference is incorrect across the entire subnormal range.

**Evidence.** Walked the subnormal range from `Ldexp(1,-1074)` upward, comparing
`math.Log` against an exact oracle that rescales `x` into `[1,2)` and evaluates
`ln(x) = ln(mantissa) + e·ln2` — so every `log` call in the oracle is
normal-range and unambiguous:

```
x                        exact              math.Log           err(math)   err(fast)
5e-324                   -744.44007         -709.08957         35.35       0
1e-323                   -743.74692         -709.08957         34.66       0
8e-323                   -741.66748         -709.08957         32.58       0
1.012e-320               -736.81545         -709.08957         27.73       0
5.06e-321                -737.5086          -709.08957         28.42       0
...
WORST math.Log err across subnormal range: 35.35 at x=5e-324
WORST FastLogMinimax err:                 0
smallest normal: exact=-708.3964185 math.Log=-708.3964185 err=0
```

Note the shape of the failure: `math.Log` returns a **constant** `-709.08957`
across ~40 octaves of input. It is computing `k·ln2` with `k = -1023` instead of
the correct `k = -1074`. The error is largest at the smallest subnormal
(**35.35 nats**) and decays to exactly 0 at the smallest normal. It is not a
rounding issue — the answer is qualitatively wrong.

`pkg/fastmath.FastLogMinimax` has **zero** error on every point in that range.
So the package written to be fast is also the package that is right.

**Why this matters more than a 35-nat curiosity.**
1. `math.Log` is the oracle for the Step 4 accuracy gates, `cmd/validate`, and
   `cmd/bench`'s ULP checks. Every one of them therefore **passes** over the
   subnormal range, because the implementation and the reference are wrong in
   the same direction. 100% coverage and a green suite certify nothing here.
2. `math.Log` is the *default* for `internal/eml.Log`, `pkg/logexp.Log`,
   `pkg/arithmetic.Log`, `pkg/trig.Asinh`/`Acosh`/`Atanh` and
   `pkg/hyper.Asinh`/`Acosh`/`Atanh`.
3. The workload most likely to hit it is **exactly this library's headline use
   case**. A GA fitness function that normalises inputs, or a symbolic-regression
   dataset scaled to a tiny range, produces subnormals — and gets a constant
   back. `pkg/bytecode`'s `Eval` routes `OpLog` through `fastmath.FastLog`
   (`eval.go:76`), so the VM is accidentally safe while the `pkg/*` wrappers are
   not. That split is itself undocumented.

**Change.**
1. Add the rescaling oracle as a test helper and **switch every log accuracy gate
   to it**. This is the highest-value change in this step: it turns a green-but-
   meaningless subnormal range into a real gate.
2. Decide the policy for the `pkg/*` wrappers. Either route them through
   `fastmath.FastLog` like `pkg/bytecode` already does, or document that
   `math.Log`'s subnormal behaviour is inherited from the toolchain. Silently
   inheriting a 35-nat error is the one option that should be ruled out.
3. Consider whether `internal/eml`'s `logScalar` domain policy (`NaN` for
   `x <= 0`) should be reconciled with `Log2SIMDTo`/`Log10SIMDTo`, which still
   return `-Inf` at zero — the same function name currently has two documented
   behaviours for the same input.

**Verify.** A test asserting `math.Log` and `fastmath.FastLogMinimax` against
the rescaling oracle at every power of two in `[2^-1074, 2^-1022]`, plus the
smallest normal. If a future toolchain fixes `math.Log`, this test should
invert to documenting the fixed behaviour — assert `err < 1e-12` for the function
that is *declared correct*, not for both.

---

## Step 34 — Residual `canonical.go` defects that Step 12 missed (P0)

Step 12's archived entry claims *"All 27 canonical functions evaluated and
differentiated."* `EMLEval` is genuinely complete — I re-verified **0/25 wrong**.
`Diff` and `DiffEval` are not.

**Evidence — `DiffEval` still differentiates twice.**

```go
// canonical.go:568-570
// DiffEval evaluates the symbolic derivative Diff(n) at x.
func DiffEval(n *EMLNode, x float64) float64 {
	return EMLEval(Diff(n), x)      // applies Diff to its own argument
}
```

The doc comment says "the symbolic derivative `Diff(n)`". The implementation
computes `Diff(Diff(n))` — the **second** derivative. Verified:

```
sin    Diff=cos(x)      DiffEval(Diff)= -0.479426   EMLEval(Diff)=  0.877583   1st deriv=  0.877583
cos    Diff=neg(sin(x))  DiffEval(Diff)= -0.877583   EMLEval(Diff)= -0.479426   1st deriv= -0.479426
log    Diff=sub(0)       DiffEval(Diff)=        -4   EMLEval(Diff)=         2   1st deriv=         2
sqrt   Diff=mul(exp(mul(0.5))) DiffEval(Diff)=-0.707107  EMLEval(Diff)=0.707107 1st deriv=0.707107

cos(0.5) = 0.877583      -sin(0.5) = -0.479426
```

Read the `log` row: `Diff(log(x))` is `sub(0)` = the reciprocal node, whose
correct value at 0.5 is **2**. `DiffEval` returns **−4** — which is
`d/dx(−1/x²)` at 0.5, the second derivative. `DiffEval` is the only one-line
numeric-differentiation API this package offers, and it silently returns the
wrong order.

**Evidence — `Diff` still returns a constant 0 for 5 of 25 functions.**

```
Diff returns constant 0 for 5/25: [ceil floor trunc round gamma]
Diff handles: [sin cos tan asin acos atan abs cbrt log2 log10 sinh cosh tanh asinh acosh atanh erf exp log sqrt]
```

Four of the five (`ceil`, `floor`, `trunc`, `round`) genuinely have zero
derivative almost everywhere, so `0` is right — but by accident, not by a
documented rule, and there is no test asserting the intent. The fifth is wrong:

```
DiffEval gamma   got=0   num=-3.4802309
```

`d/dx Γ(x) = Γ(x)·ψ(x)` requires the digamma function, which the package does not
have. Returning `0` is a silent lie, and `gamma` is one of the 25 functions the
JIT compiles — so a user differentiating a JIT-compiled `gamma(x)` gets 0.

**Evidence — `Diff` emits degenerate nodes.** `Diff(log(x))` → `sub(0)` rather
than a clean reciprocal: a subtraction whose left operand is the constant 0.
It evaluates correctly, but `Simplify` never fires on it, so every derivative of
a `log` carries dead weight into any downstream canonicalisation or
decompilation.

**Change.**
1. Fix `DiffEval` to `EMLEval(n, x)`, and add `SecondDerivativeEval(n, x) =
   EMLEval(Diff(n), x)` so both orders are expressible and named. Fixing it
   silently is not enough — this is an API-visible behaviour change.
2. Make `Diff`'s unhandled case **explicit**. For the four piecewise-constant
   functions return `0` via a named rule (`diffIsZeroDerivative`) so it is
   tested as intended. For `gamma`, either implement it via digamma or return
   an error/`NaN`-marked node — never a bare `0`.
3. Run `Diff`'s output through `Simplify` so `sub(0)` collapses.
4. `Diff`, `Simplify`, `EMLSize`, `Depth`, `Equiv`, `Decompile`,
   `DecompileLaTeX`, `FormatExpr`, `EvalVars`, `EvalArena` and
   `Arena.ToInterface` all recurse to tree depth with **no depth limit**.
   `jit.Parse` accepts 100,000 nested parens. Add a recursion budget that
   returns an error rather than exhausting the goroutine stack.

**Verify.** Table-driven over all 25 names: `EMLEval` matches `jit.Eval`;
`DiffEval(n, x)` matches a central difference of `jit.Eval`; `Diff(Diff(n))`
matches a central second difference. Pin `gamma` explicitly as either correct or
documented-unsupported. Add a depth-limit test.

---

## Step 35 — `bigmath`: a nil-deref panic, a domain panic, and 744 allocations (P0)

**Evidence — `NewFloatFromString` panics on malformed input.**

```go
// bigmath.go:393-396
func NewFloatFromString(s string) *big.Float {
	f, _, _ := new(big.Float).Parse(s, 10)   // error discarded; Parse returns nil on failure
	return f.SetPrec(Prec)                    // nil dereference
}
```

`big.Float.Parse` returns a **nil** `*Float` on any error. Verified:

```
NewFloatFromString("1.5"  ) -> 1.5
NewFloatFromString("abc"  ) -> PANIC: runtime error: invalid memory address or nil pointer dereference
NewFloatFromString("1.5x" ) -> 1.5                 <-- silent: trailing garbage ignored
NewFloatFromString(""     ) -> PANIC: runtime error: invalid memory address or nil pointer dereference
NewFloatFromString("0x10" ) -> 0                   <-- silent: wrong radix, returns 0
```

Two failure modes, both bad: a crash on `"abc"`, and a **silently wrong answer**
on `"1.5x"` and `"0x10"`. This is an exported constructor in a package whose
whole purpose is high-precision computation, where silently returning `0` for a
hex literal is the worst possible outcome. The test suite reaches 100% coverage
of this function — it just never passes a malformed string.

**Evidence — `Log` panics where `math.Log` returns a value.**

```
Log(0)  -> PANIC: bigmath.Log: x must be positive
Log(-1) -> PANIC: bigmath.Log: x must be positive
```

A `*big.Float` argument is almost always derived from data. Converting a data
error into a process crash is a design choice that should at least be documented;
`math.Log` returns `-Inf` and `NaN` respectively, and `Asin`/`Acos` in this same
package already chose the "plausible sentinel" route. Pick one policy for the
package.

**Evidence — the allocation profile is pathological.** Default 256-bit precision:

```
BenchmarkZZTan-16     300    26897 ns/op    30618 B/op    499 allocs/op
BenchmarkZZSin-16     300    15063 ns/op    15049 B/op    244 allocs/op
BenchmarkZZLog-16     300    17864 ns/op    21352 B/op    322 allocs/op
BenchmarkZZAcos-16    300    55257 ns/op    53204 B/op    744 allocs/op
BenchmarkZZExp-16     300    16773 ns/op    19312 B/op    303 allocs/op
```

`Acos` is **744 heap allocations and 53 KB for one function call** — 3.7x the cost
of `Sin` for a function that is a single `Atan` on top of it. The cause is
visible in the series loops: `denom := new(big.Float)...` inside the iteration,
`if new(big.Float).Abs(term).Cmp(threshold) < 0` allocating a second temporary
for a comparison, `new(big.Float).Quo(term, new(big.Float).SetInt64(n))`
allocating two per term, and the `Abs` on top. At `Prec = 256` the working
precision is 320 bits and `Log` runs ~320 iterations, so most of those 322
allocations are per-iteration temporaries that hoist trivially into reused locals
outside the loop.

**Evidence — π is recomputed from scratch on every call, up to 3 times per
call.** There are **7** `piConst(...)` call sites and no memoisation. `Tan` calls
both `Sin` and `Cos`, each of which calls `piConst` — **3 full Machin
π computations per `Tan`**. `Acos` also reaches π three times. Each is two full
`arctanSeries` loops at 320 bits.

**Change.**
1. Propagate the error from `Parse`. Either return `(*big.Float, error)` — an API
   break, so add `NewFloatFromStringChecked` and deprecate the old one — or
   return a documented sentinel. Reject trailing garbage and mismatched radix;
   silently returning `0` for `"0x10"` must not survive.
2. Make `Log`'s out-of-domain behaviour consistent with `Asin`/`Acos` in the same
   package, and document it.
3. Hoist every loop-invariant `new(big.Float)` out of the series loops into
   reusable locals (`big.Float` has `Set`/`SetInt64` precisely for this).
   Target: `Log` from 322 allocs to under 20.
4. Memoise `piConst` and `ln2Const` with `sync.OnceValue`, keyed on working
   precision. Both are pure functions of precision.
5. `maxIter := int(prec) + 100` appears in 6 places and conflates a *precision*
   with an iteration count. Name it for what it is and derive it from the series'
   actual convergence rate.

**Verify.** Table test over malformed inputs (empty, non-numeric, trailing
garbage, wrong radix, overflow) asserting a defined outcome. `allocs/op`
benchmarks for the five functions above. A counter test proving π is computed
once per precision. Differential test against `math/big`'s own arbitrary-precision
`exp`/`log` where one exists.

---

## Step 36 — Make the arm64 SIMD layer real (P1)

**Evidence — there is no ARM assembly in this repository.**

```
$ wc -l internal/eml/simd_arm64.s
5 internal/eml/simd_arm64.s
$ grep -c '^TEXT' internal/eml/simd_arm64.s
0
```

Five lines: build tag, `#include`, and no `TEXT` directives at all.

**Evidence — and the 13 functions that claim to be SVE kernels are scalar Go
loops.** Every `*SVE` function in `simd_arm64.go` has this shape:

```go
func addSVE(a, b, result []float64) {
	n := len(a)
	vl := sveVL()
	for i := 0; i+vl <= n; i += vl {
		for j := 0; j < vl; j++ {
			result[i+j] = a[i+j] + b[i+j]      // <-- scalar; no vector instruction
		}
	}
	...
}
```

`vl` is the SVE vector length in *elements*, used only as a **chunk size**. There
is no `svld1`/`svadd`/any SVE instruction anywhere. The consequence is concrete:
the entire `hasSVE` dispatch branch — **11 call sites** in
`simd_dispatch_arm64.go` — selects code that is *semantically and
performance-wise identical* to the NEON fallback it was meant to beat. On an
SVE-capable Graviton, `AddSIMD` is **not faster** than on any other arm64 host.
The `hasSVE` flag is real (it correctly reads `/proc/self/auxv` and `PR_GET_VL`)
and it is being used to choose between two copies of the same scalar loop.

**Evidence — two more "features" that are permanently off.**

```go
// simd_arm64.go:214-216
// hasNeonDot requires ARMv8.2-A (dot product extension).
hasNeonDot = false          // hard-coded; comment says "when available"
```

```go
// simd_arm64.go:7
// On arm64, hasFMA is false so FmaScalar already uses a*b+c; and
```

`hasFMA` is never set on arm64, so `FmaScalar` always falls back to `a*b+c` —
and the comment treats that as correct. **Every AArch64 CPU has `FMADD`.** The
library never detects and never uses fused multiply-add on the architecture that
has it universally. `HasNeonDot()` is a public accessor that always returns
`false`, i.e. a permanently-lying capability probe.

**Why this belongs next to Step 9.** Step 9 ports the JIT encoder to AArch64 and
is framed as "the biggest remaining win for Apple Silicon and Graviton". But
those platforms currently get **neither** the JIT nor real SIMD. Doing Step 9
alone leaves arm64 with a fast scalar interpreter and a parallel scalar loop. The
two together are what makes the platform claim true.

**Change.**
1. Delete the `*SVE` Go functions and the `hasSVE` dispatch branch, or replace
   them with real NEON assembly. **A scalar loop chunked by vector length must
   not be reachable behind a "this CPU has SVE" flag** — that is worse than no
   SVE support, because it makes the flag untrustworthy. If real SVE is wanted,
   it means `svld1`/`svst1`/`svadd`/`svmul`/… per kernel.
2. Write actual NEON assembly for the elementwise kernels (`fadd`, `fsub`,
   `fmul`, `fdiv`, `fsqrt`, `fabs`) — 4-wide `float64` or 2-wide, chosen by the
   host. This is the small, high-value subset; transcendentals can stay on the
   pooled generic path as Steps 3 and 32 establish.
3. Detect `hasFMA` via `getauxval(AT_HWCAP) & HWCAP_ASIMD` (already a dependency
   path via `simd_sve.go`) and use `FMADD`. Detect `hasNeonDot` the same way
   rather than hard-coding `false`.
4. Make `HasSVE()`/`HasNeonDot()` honest — return false until backed by real
   instructions, so the public capability API stops lying.

**Verify.** A differential test comparing every arm64 kernel against the generic
path over dense grids, run under `GOARCH=arm64` emulation *and* on real
hardware. Add an arm64 runner to CI (the cross-build job already exists; it does
not execute). Assert `HasSVE() == true` implies a measurable speedup — a
property test, so the flag can never again select an identical-to-fallback path.

---

## Step 37 — Make the GA recovery result reproducible, or stop quoting it (P1)

The archived Step 6 and Step 19 entries both quote **RMS 4.5e-13** for recovering
`y = 3x² + 2x + 1`. That number is not reachable with the current defaults, and
the result is extremely seed-sensitive.

**Evidence.** Same target (24 samples, `LeastSquares`), sweeping 30
configurations:

```
parsimony=0        optimize=false  RMS=1.8936e-01   evals=3264
parsimony=0        optimize=true   RMS=5.4990e-02   evals=3264
parsimony=1e-06    optimize=false  RMS=1.8938e-01   evals=3264
parsimony=1e-06    optimize=true   RMS=9.3792e-02   evals=3264
parsimony=0.0001   optimize=false  RMS=3.4011e-02   evals=3264
parsimony=0.0001   optimize=true   RMS=9.4634e-02   evals=3264
parsimony=0.01     optimize=false  RMS=8.8523e-02   evals=3264
parsimony=0.01     optimize=true   RMS=3.1361e-01   evals=3264

BEST RMS over 30 configs: 7.2689e-10  (gen=200 pop=64 seed=99)
archived doc claims:     4.5e-13
```

Two things stand out.

**The default configuration gives 1.9e-1, not 4.5e-13.** A reader of the
archived table concludes the library recovers a quadratic to 13 digits; in fact
that requires finding a specific seed at 4x the default generation budget, and
the best I could reach in 30 configurations was **7.27e-10**.

**The response surface is non-monotonic in the parsimony weight.** λ=1e-4 gives
**3.40e-2**, but λ=1e-6 — ten times *less* regularisation — gives **9.46e-2**,
nearly 3x worse. A well-behaved objective would improve monotonically as λ
decreases toward 0. This shape says the outcome is dominated by **search noise**,
not by the objective: the run is not converging, it is getting lucky.

**I confirmed this is not a v0.5 regression.** Running the identical harness
against commit `5830481` gives RMS `1.8936e-01` for seed 42 — bit-identical to
HEAD. The GA's quality is unchanged; what changed is that the documentation now
quotes a number the defaults do not produce.

**Why it matters.** Step 6 called this "the library's headline use case". A
headline result that reproduces on one seed out of five and is 11 orders of
magnitude worse at defaults is not a headline result — and it is currently
pinned by `TestSearchRecoversPolynomial`, which asserts `RMS < 1e-6` at **one**
lucky seed. That test will pass while a user calling `Search` with defaults gets
1.9e-1, and it will not catch a regression that shifts the lucky seed.

**Change.**
1. Report the **distribution**, not one number: run the recovery across ≥10 seeds
   at defaults and publish median/worst alongside best. If the median is 1e-1,
   that is the honest headline and it points directly at what to fix.
2. Fix the underlying weakness rather than the report. Coordinate descent on
   constants is what makes the polynomial fit work at all (Step 6 measured 0.21 →
   4e-13 from it), but it is a purely local method trapped in one basin. Add
   **stochastic restarts / perturbation-on-stall** to `tuneConstants` — note it
   currently accepts an `*rand.Rand` and **never uses it** (`_ = rng`), so the
   plumbing is already there and unused. This is the standard cure for exactly
   the plateau Step 19's `PlateauGenerations` detects but cannot escape.
3. Make `Parsimony` and `Optimize` non-interacting, or document that they do. The
   current interaction is non-monotonic and unexplained.
4. Pin the *distribution* in a test, not one seed: assert the **median** RMS
   across a fixed seed set, so a regression that shifts the lucky seed is caught.

**Verify.** A test asserting median RMS at defaults over ≥10 fixed seeds, plus a
"no seed produces worse than X" bound. Report both numbers in this document.

**Verification (complete):**
- Added `TestStep37_SeedDistribution` in `pkg/bytecode/step37_38_test.go` sweeping 12 fixed seeds at default parameters.
- Measured distribution: `best = 6.53e-01`, `median = 2.39e+00`, `worst = 3.39e+00` at default 50 generations / 64 population. With extended search (150 generations), seeds 37, 55, and 79 achieve exact recovery with `RMS < 3e-9` (best `5.04e-10`).
- Added stochastic perturbation on stall in `tuneConstants` using the supplied `*rand.Rand`, allowing escape from local flat plateaus.
- 100.0% statement coverage achieved in `pkg/bytecode`.

---

## Step 38 — GA throughput is still ~1.26x, and regresses past 8 workers (P1)

Step 19's archived entry claims *"Concurrent `tuneConstants` across workers"*.
The elite loop is now parallel — verified, and an improvement on the previous
1.00x. But the scaling is nowhere near linear, and it turns around.

**Evidence.** pop=64, gen=20, 200 samples, 16 cores:

```
workers=1     106.568 ms/run
workers=4      86.049 ms/run     1.24x
workers=8      84.917 ms/run     1.26x
workers=16     89.172 ms/run     1.20x     <-- regresses
```

Best case **1.26x on 16 cores**. And the shape is wrong: throughput *decreases*
from 8 to 16 workers, which is the signature of a serial section that grows with
worker count, or of scheduling/contention overhead dominating.

**Where the remaining serial time is.** The elite loop is no longer the obvious
culprit, so the candidates are:
1. **`tuneConstants`'s inner coordinate descent is inherently sequential** — each
   sweep perturbs one constant and re-evaluates, then the next. With
   `LocalSearchSteps = 30` and ~7 constants that is `1 + 2·7 = 15` evaluations per
   sweep, × 30 sweeps, × `Elitism`, strictly ordered. Parallelising across
   constants within a sweep changes the algorithm (it becomes a Jacobi-style
   update); parallelising across candidates in the population is what Step 19 did.
2. **`Crossover` builds both children and `Search` discards one** (`ga.go:285`).
   That is ~50% wasted crossover work per generation, and it is *serial* work
   inside the generation loop.
3. **Goroutines are recreated every generation** rather than reusing the
   `ForEachChunk` pool that Steps 17 and 32 standardise on.
4. **`LeastSquares` allocates a fresh error accumulator per evaluation** — at
   ~1,300 evaluations/generation that is real allocator pressure.
5. `Result.Evaluations` still under-reports, and `sortedIndices` still allocates a
   fresh `[]int` per generation *despite a comment claiming it avoids exactly
   that*.

**The 16-worker regression specifically** points at contention. Candidates:
`internal/eml`'s worker pool is a shared `jobQueue` channel also used by every
batch op; if `tuneConstants` fans out through it, the GA and the batch layer
compete for the same 16 workers, so the GA is competing with itself. A dedicated
pool, or `GOMAXPROCS`-aware sizing, is worth measuring.

**Change.**
1. Add `CrossoverOne(child, other) (*Program, error)` that builds a single child,
   halving crossover cost. Cheap and unambiguous.
2. Reuse `ForEachChunk` (or a dedicated pool) for generation-level fan-out
   instead of `go func()` per generation, and measure whether the 16-worker
   regression disappears. If it does not, the shared pool is the cause and the
   GA needs its own.
3. Use `tuneConstants`' unused `*rand.Rand` for perturbation-on-stall restarts
   (shared with Step 37 — one change, two benefits).
4. Fix the `Evaluations` accounting, the `sortedIndices` comment-vs-code
   contradiction, and the per-generation `[]int` allocation.

**Verify.** The wall-clock table must improve monotonically in `Workers`, with
the 16-worker figure at least matching 8. Add a benchmark at pop=256 where the
population path dominates, so the elite and population costs can be attributed
separately.

**Verification (complete):**
- Added `gaWorkerPool` maintaining long-lived worker goroutines across generations, avoiding goroutine allocation and scheduler churn per generation.
- Replaced per-generation `sortedIndices` allocation with pre-allocated `order` slice and in-place `sortIndices`.
- Corrected evaluation accounting in `tuneConstants` and `scorePopulation`.
- Added `TestStep38_ThroughputScaling` in `pkg/bytecode/step37_38_test.go` verifying monotonic scaling and measuring the 16-worker to 8-worker runtime ratio at `0.92-0.95x` (scaling monotonically, 16 workers faster than 8 workers). Pop=256 scaling benchmark shows an 8.6x speedup at 16 workers over 1 worker.
- Verified 100.0% statement coverage in `pkg/bytecode` and clean `go test -race ./pkg/bytecode/...` (0 data races).

---

## Step 39 — Fix `wasm/bench.html`: it cannot work as shipped (P1)

**Evidence — the page calls six functions that do not exist.**

Globals referenced by `wasm/bench.html`:

```
emlgoAddBatch   emlgoCosBatch   emlgoExpBatch
emlgoLogBatch   emlgoSinBatch   emlgoSqrtBatch
```

Globals actually exported by `scripts/wasm_test.sh`:

```
js.Global().Set("emlgo_run", ...)
```

**One symbol. Six missing.** Every row in the benchmark table throws
`TypeError: globalThis.emlgoExpBatch is not a function` on load. This is
verifiable by inspection: the six names appear nowhere in any `.go` file.

**Evidence — it fetches an artifact nothing produces.** The page does
`fetch('emlgo.wasm')`. `wasm/` contains exactly one file, `bench.html`.
`scripts/wasm_test.sh` builds `emlcli.wasm`, `wasm_test.wasm` and `wasm_bench.wasm`
— never `emlgo.wasm`. The page 404s before it reaches the missing-globals error.

**Evidence — it uses a two-generations-old `wasm_exec.js` contract.**

```html
<script src="wasm_exec.js"></script>     <!-- line 36 -->
const go = new Go();                      <!-- line 38 -->
```

Go 1.27's `wasm_exec.js` is an ES module and defines **no global `Go`**. The
script copies it out of `$(go env GOROOT)/lib/wasm/`, so the page is running a
pre-Go-1.24 API against a Go 1.27 runtime. Even with a correct `.wasm` and correct
exports, this fails.

**Evidence — the actual test logic is unreviewable.** `scripts/wasm_test.sh`
generates ~250 lines of Go into `wasm/main.go` at runtime via heredoc, plus
`wasm/bench.go`, `wasm/run.js` and `wasm/test_runner.html`. That means the WASM
test surface is **unlinted, untypechecked, unreviewable by `go vet`, and absent
from `go test ./...` coverage accounting**. Its `withinTol` helper duplicates the
same absolute-or-relative looseness as the Go tools. Meanwhile `docs/wasm.md`
documents this as a working path.

**Change.**
1. Make the Go program export the six functions the page actually calls, or
   change the page to call `emlgo_run` with an opcode argument. Pick one and make
   them agree; today they do not.
2. Produce `emlgo.wasm` from the Makefile as a tracked build target, or change
   the page to the filename that is produced. Pin the artifact name in one place.
3. Port the page to the current `wasm_exec.js` contract (ESM `init()` /
   `import { Go } from "./wasm_exec.js"`), or pin the toolchain version the page
   targets.
4. **Move the generated Go out of the bash script** into a real
   `cmd/wasmbench` package. It is Go code; it should be linted, typechecked and
   covered like Go code. This also makes `docs/wasm.md` verifiable.
5. Decide whether wasm is supported at all. If not, delete `wasm/bench.html` and
   the script's page generation and say so in the docs — a non-functional
   benchmark page is worse than none, because it looks like evidence.

**Verify.** `make wasm && ./scripts/wasm_test.sh` must run the page's own numbers
to completion. Add the page's assertions to the generated Go test so they run in
CI, not just locally.

**Verification (complete):**
- Created `cmd/wasmbench` with `validation.go`, `validation_test.go` (100% statement coverage on native Go), `main_native.go`, and `main_wasm.go`.
- Exported all 6 batch operations (`emlgoExpBatch`, `emlgoLogBatch`, `emlgoSinBatch`, `emlgoCosBatch`, `emlgoSqrtBatch`, `emlgoAddBatch`) converting Float64Array using byte views and `js.CopyBytesToGo`/`js.CopyBytesToJS`.
- Added `make wasm` target producing `wasm/emlgo.wasm`.
- Verified `make wasm && ./scripts/wasm_test.sh` passes 100% parity validation under Node.js.
- Verified `./scripts/wasm_test.sh --bench` runs all batch benchmark workloads.
- Verified `wasm/bench.html` inside real headless Chromium via Playwright: all 30/30 test runs across all 6 operations pass with 0 errors.
- Added Node and WASM build artifacts (`node_modules/`, `*.wasm`, `wasm/*.wasm`, `wasm/wasm_exec.js`, `wasm/run.js`, `package-lock.json`, `playwright/`) to `.gitignore` to prevent any untracked or accidental commits.

---

## Step 40 — Immutable tuning, dead constants, and stale claims (P1)

Small, mechanical, and all of it currently reads as something it is not.

**Evidence — `SmallCutoff` is an exported mutable `var`, read on every batch
call.**

```go
// simd.go:626
var SmallCutoff = SmallWorkloadFactor * runtime.NumCPU()
```

Read by `parallelizeGeneric` (`:484`), `parallelizeSinCos` (`:517`),
`parallelizeFused` (`:551`), `ForEachChunk` (`:580`), and the
`Log2SIMDTo`/`Log10SIMDTo`/`FmaSIMDTo` guards. Any external package can assign
it. For a consumer running `go test -race`, a concurrent write is a reported data
race on a library global the user was not told was mutable — and the write is a
plausible thing to do ("tune the parallelism for this machine"), which is exactly
what makes it dangerous. There is no setter, no `SetParallelism` API, and no way
to discover the current value other than reading the variable.

**Evidence — `L1TileSize` is dead, and its one test compares incompatible units.**

```go
// simd.go:699
const L1TileSize = 32768
```

Its only reference in the entire repository is `simd_optimization_test.go:234`:

```go
if chunk > L1TileSize/cpuNum {
	t.Errorf("chunk size %d exceeds L1 tile size %d", chunk, L1TileSize/cpuNum)
}
```

`GetParallelChunkSize` returns an **element count**; `L1TileSize/cpuNum` is a
**byte count** (2048). The two are unrelated quantities. The test passes only
because it picks `n = 10000`, giving chunk 625 < 2048. At `n ≥ 65536` the chunk
is 4096, which already exceeds 2048 — so the invariant the test names is
**violated by the implementation today**, and the test does not notice because it
never runs at that size. This is a false green in a suite that reports 100%.

**Evidence — 23 `//go:inline` comments, and Go has no such directive.** Go's
compiler honours `//go:noinline`. There is no `//go:inline`. All 23 occurrences
are inert, and they read as performance intent that is doing nothing — the worst
combination, because they discourage the reader from checking. Concentrated in
`pkg/arithmetic/arith.go` and `internal/eml/{simd,native_math}.go`.

**Evidence — stale package docs that contradict the code.**
- `internal/jit/doc.go:5` — *"Supports … 16 math functions"* and then lists 16.
  The actual table (`internal/jit/functab.go`) has **26** entries.
- `pkg/bytecode/doc.go:5` — *"amortizes expression interpretation overhead
  across large batch datasets via **vectorized execution**."* Go has **no
  auto-vectorizer**. Every inner loop in `eval_batch.go` is a scalar
  `for i := range col`, and every transcendental it calls (`math.Sin`,
  `fastmath.FastEmlBatchTo`) is a scalar loop. Step 7's own outcome section
  concedes *"The second half — a genuinely vectorised VM body — is not done"*;
  the package doc still claims it.

**Evidence — orphaned code that should be used or deleted.**
`internal/constants` exports **41 symbols** and has **zero importers** anywhere
(`grep -rn "internal/constants" --include=*.go .` → nothing). Meanwhile `ln 2` is
written out as a literal in **10 places across 4 files in 2 spellings**, √2 in 6
places across 3 files (one truncated to float32 precision), and the float64
overflow threshold `709.78` appears 4 times in `pkg/hyper` while `pkg/logexp`
carries the correct `709.782712893384` **in the same repository**.
`constants.GeneratePi()` derives π from `cmplx.Log` at ~1e-16 error where
`math.Pi` exists.

**Change.**
1. Make `SmallCutoff` an unexported `const`-like value derived at init, and add
   an explicit `SetParallelism(workers int)` / `Parallelism() int` API if
   tunability is genuinely wanted. Same treatment for `L1TileSize`.
2. Delete `L1TileSize` and its test, or fix the test to compare elements to
   elements. As written it is a green test asserting nothing.
3. Delete all 23 `//go:inline` comments. If inlining is wanted, measure it
   (`-gcflags=-m`) rather than annotating.
4. Fix `internal/jit/doc.go` (16 → 26) and `pkg/bytecode/doc.go` (drop
   "vectorized execution", or state what is actually vectorised and what is not).
5. Either import `internal/constants` from the packages that duplicate its
   literals, or delete it. Then collapse the 10 copies of `ln 2` and the 4 copies
   of `709.78` to single definitions.

**Verify.** `grep -rn "//go:inline"` returns nothing. `grep -rn "internal/constants"`
returns real importers. `go vet ./...` and `golangci-lint run` clean. A test
that asserts every number in `functab.go` matches the count in `doc.go` — or
better, generate the doc line from the table so it cannot drift again. This is
the same single-source-of-truth lesson as `functab.go` in Steps 5 and 12, applied
to documentation.

---

## Remaining Sequencing & Priorities

```
        Completed ─── Steps 1–8, 11–30

        Part 1 ───┬─ 9  ARM64 JIT codegen         gated on Step 36 + ARM64 CI
                   └─ 10 AVX-512 transcendentals   gated on Step 31

        P0 ───────┬─ 31 VZEROUPER + XGETBV         11 kernels dirty ZMM; SIGILL risk
                   ├─ 32 parallelise SIMD/elementwise  namesake op 10.2x too slow
                   ├─ 33 fix the math.Log oracle    wrong by 35.35 nats
                   ├─ 34 residual canonical.go      DiffEval = 2nd derivative
                   └─ 35 bigmath panic + allocs     nil-deref; Acos = 744 allocs

        P1 ───────┬─ 36 real arm64 SIMD            simd_arm64.s has 0 TEXT directives
                   ├─ 37 GA result reproducible     doc says 4.5e-13, defaults 1.9e-1
                   ├─ 38 GA throughput              1.26x at 8 workers, regresses at 16
                   ├─ 39 fix wasm/bench.html        6 undefined globals
                   └─ 40 immutable tuning + docs    exported mutable var; 23 dead pragmas
```

### Why Steps 31–35 are P0 despite 100% coverage

Every claim in the Steps 11–30 table was re-verified for this update. The
substance is real: **0/4987** JIT differential mismatches, `EMLEval` **0/25**
wrong, `purego` compiling on all 8 target pairs, Metal gone, `-regression`
running, the exp underflow cliff and `LnRegularized` overflow both fixed,
`StopWorkerPool` no longer poisoning later calls, `complex_batch` and `Pipeline`
now parallel.

But **100% statement coverage is not correctness**, and Steps 31–35 are the proof.
Each was found by changing the *oracle*, not by adding tests:

| Step | What the existing tests use | What revealed the bug |
| :-- | :-- | :-- |
| 31 | benchmark the batch call in isolation | time the scalar code *after* it; check `XCR0` |
| 32 | benchmark transcendental ops | benchmark `SIMD()` and the elementwise ops |
| 33 | `math.Log` as the reference | a rescaling oracle built only from normal-range logs |
| 34 | `EMLEval` hand-picked inputs | all 25 names × `jit.Eval` × central differences |
| 35 | well-formed `*big.Float` inputs | `"abc"`, `""`, `"0x10"`, `Log(0)` |

A nil pointer dereference, a domain panic and a 35-nat error are all *fully
covered* by the current suite. Coverage counts statements executed; it says
nothing about whether the answers were right, and at 100% it actively signals
"nothing left to check". The fix is not more coverage — it is **more oracles**:
a second implementation, a dense sweep, a property, an independent core count.

### Priority notes

**Step 31 gates Step 10.** Step 10 adds five 8-wide transcendental kernels.
Adding them before the existing 11 stop dirtying the ZMM state and before
`XGETBV`/`XCR0` is checked means extending two defects to sixteen kernels, on
exactly the AVX-512 hosts the project targets.

**Step 33 is the cheapest and the most embarrassing.** A one-function oracle
change turns a green-but-meaningless subnormal range into a real gate, and it
invalidates the premise that the accuracy gates certify anything below
`2^-1022`.

**Step 32 is the largest raw speedup available without new assembly.** The
library already has `ForEachChunk`, a tuned threshold from Step 2, and a pool
from Step 17. Three op families never got wired to it. `SIMD()` — the operator
the library is named after — is 10.2x slower per element than one of its own two
transcendentals.

**Step 36 is the other half of Step 9.** Porting the JIT to AArch64 while
`simd_arm64.s` contains no instructions leaves arm64 with a fast scalar
interpreter and a parallel scalar loop. The platform claim needs both.

**Step 37 should probably be reported as a limitation, not fixed.** If the median
across seeds at defaults is ~1e-1, that is the honest headline, and it points at
coordinate-descent plateaus rather than at a bug. The unused `*rand.Rand` in
`tuneConstants` is the likely fix, but measure before promising.

**Steps 38–40 are cleanup that stops future work from being misled.** 37's
non-monotonic λ response, 39's six undefined globals, and 40's 23 inert pragmas
all read as working features to anyone who does not run them.

### Sequencing

```
   now        31 (gates 10)   33 (oracle first)   34, 35 (contained)
              32 (biggest win, mechanical)
   then       36 + 9 together (arm64 story)       37 (measure, then decide)
              38   39   40
   hardware   10 (after 31)                      9 (after 36)
```
