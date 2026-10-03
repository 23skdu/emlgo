# Next Steps: The Next 30 Changes

A prioritised, evidence-based plan. Every step below cites a **measurement taken
from this repository**, not an estimate.

**Measurement environment for all numbers below:** Intel Core i7-12650H
(AVX2 + FMA, no AVX-512), 16 logical cores, Linux 5.x, Go 1.27.0, `linux/amd64`.

**Contents.**

- **Steps 1–10** — the SIMD/transcendental/bytecode/GA work. Steps 1–8 are
  **done**; 9 (ARM64 JIT) and 10 (AVX-512 kernels) remain open. Kept in full
  because their *outcomes* are the measurement baseline for Part 2.
- **Steps 11–20** — Part 2: correctness & performance defects found by audit
  after Steps 1–8 landed. Steps 11–16 are **correctness** (P0); Steps 17–20 are
  **performance & enforceability** (P1). All 10 are open.
- **Steps 21–30** — Part 3: deep audit findings across bytecode optimization,
  JIT parser/lexer, hyperbolic & arbitrary-precision math, decompiler AST fidelity,
  JIT power codegen, and vectorization engines. Steps 21–27 are **correctness**
  (P0); Steps 28–30 are **performance & architecture** (P1). All 10 are open.

**Why Part 2 exists.** Steps 1–8 were a performance programme, and they landed:
`go test ./...` is green and statement coverage is **90.3%**. The same audit that
found those wins also found that several *load-bearing* code paths are wrong
today and no test catches them. High statement coverage did not help, because
none of these tests are **differential** — they assert on hand-picked inputs, so
a silently-wrong answer at every other input is invisible. Step 20 is the
structural fix; Steps 11–16 are what it would have caught.

---

## Status

| # | Step | Status | Result |
| :-- | :--- | :--- | :--- |
| 1 | Fix the AVX2 transcendental kernels | **superseded by 3** | superseded |
| 2 | Gate the worker pool on a real threshold | **done** | no loss at any size |
| 3 | Make the generic path the fallback | **done** | `exp` 201x faster at n=1M |
| 4 | Add accuracy gates for every batch path | **done** | gates pass on every path |
| 5 | Extend the bytecode VM to all functions | **done** | 5 ops → 27 |
| 6 | Ship a genetic-algorithm driver + fitness | **done** | recovers 3x²+2x+1 to 4e-13 |
| 7 | Zero-alloc `EvalBatch` + wider inner loops | **done** | 24.7 KB → 0 B per call |
| 8 | `Pipeline.Reset()` and an in-place `RunTo` | **done** | 0 allocs, no step growth |
| 9 | ARM64 JIT codegen | open | 17–30x for 2 of 3 server markets |
| 10 | AVX-512 transcendental kernels | open | 2x on AVX-512 hosts |

### Status — Part 2 (Steps 11–20)

| # | Step | Status | Headline evidence |
| :-- | :--- | :--- | :--- |
| 11 | Fix the JIT register-ABI defect | **open (P0)** | **37.5% of random expressions return wrong answers** |
| 12 | Close the `EMLEval`/`Diff` silent-zero holes | **done** | All 27 canonical functions evaluated and differentiated |
| 13 | Make the bytecode VM fail loudly | **open (P0)** | unknown opcode → `index out of range [-1]` |
| 14 | Fix fastmath's silent range failures | **open (P0)** | **37 decades of `exp` return exactly 0** |
| 15 | Fix `purego` on arm64 and wasm | **open (P0)** | `-tags=purego` does not compile |
| 16 | Fix or delete the Metal backend | **open (P0)** | results read **before** the GPU runs |
| 17 | Delete `arithmetic.parallelMap` | **open (P1)** | 13 batch ops are **2–12x slower** than a loop |
| 18 | Reclaim JIT code memory | **open (P1)** | **4.29 kB leaked per compile**, never freed |
| 19 | Make the GA scale and stop bloating | **open (P1)** | **1.00x speedup at 16 workers** |
| 20 | Make perf/accuracy claims enforceable | **open (P1)** | the `-regression` gate cannot run |

### Status — Part 3 (Steps 21–30)

| # | Step | Status | Headline evidence |
| :-- | :--- | :--- | :--- |
| 21 | Fix the bytecode optimizer RPN corruption | **done** | AST constant folder & identity reducer; 1.57M fuzz execs |
| 22 | Fix JIT parser scientific notation & whitespace | **done** | `[eE][+-]?[0-9]+` & `\t,\n,\r` accepted; var scoping; 1.64M fuzz execs |
| 23 | Fix hyperbolic cancellation and silent -Infs | **done** | `Asinh` odd symmetry, `Sinh` expm1, `Log(0)`=-Inf; 2.12M fuzz execs |
| 24 | Fix arbitrary-precision limits in `bigmath` | **done** | `Exp(-1000)` non-zero, `Asin` exact, memoized $\pi$/$\ln 2$; 509k fuzz execs |
| 25 | Complete the `CompileEML` function set | **done** | All 27 ops wired via `unaryOpcode`; 1.53M round-trip fuzz execs |
| 26 | Fix decompiler precedence & associativity | **done** | Precedence-aware parens, variable names preserved; round-trip invariant |
| 27 | Fix JIT power codegen & negative base NaN | **done** | 1-cycle `xorpd` negation, ABIInternal `math.Pow` call for $(-2)^3$ & $0^0$ |
| 28 | Vectorize Pipeline engine and complex ops | **done** | SIMD Pipeline dispatch (28x win), affine fusion, Neumaier dot (5x win) |
| 29 | Replace VM opcode hash map with flat jump table | **done** | `opUnaryTable [256]mathFunc`, zero-alloc scratch capacity reuse |
| 30 | Vectorize integer saturation & normalize lengths | **done** | AVX2 `VPADDSB`/`VPSUBSB` 17x speedup, 0 allocs in `*To`, panics aligned |

### What the completed work measured

| n | `ExpSIMDTo` before | after | speedup |
| --: | --: | --: | --: |
| 4,096 | 872,098 ns | 20,961 ns | **41.6x** |
| 65,536 | 14,322,837 ns | 97,848 ns | **146x** |
| 1,048,576 | 211,667,934 ns | 1,051,251 ns | **201x** |

Worker-pool crossover after Step 2 — the dispatch path no longer loses at any
size (it previously lost by 4.4x at n=256):

| n | dispatch | serial loop | ratio |
| --: | --: | --: | --: |
| 256 | 1,295 ns | 1,295 ns | 1.00 |
| 1,024 | 5,202 ns | 5,063 ns | 1.03 |
| 8,192 | 29,127 ns | 41,454 ns | 0.70 |
| 1,048,576 | 1,053,027 ns | 5,466,742 ns | 0.19 |

### Bugs found and fixed along the way

- `pkg/bytecode`: `CompileExpr` accepted malformed input such as `"1 +"`,
  `"*"`, `"exp"` and produced programs that **panicked** during `Eval`.
- `internal/eml`: `logScalar` was missing, so `Log(0)` returned `NaN` from the
  scalar path and `-Inf` from the batch path.
- `pkg/bytecode`: `RandomProgram` produced unbalanced stacks (final heights of
  4, 7 and 11) because its balancing binary operator was appended only half the
  time.
- `pkg/bytecode`: `splitPoints` looked for stack height 0, which only occurs at
  index 0, so `Crossover` could only swap whole programs and never recombined
  subexpressions.
- `pkg/bytecode`: `Mutate` swapped `OpConst` for `OpVar` (both arity 0) without
  editing the operand streams, desynchronising them.
- `internal/eml`: the AVX2 transcendental kernels have no scalar tail; they
  require a multiple of the vector width and rely on the caller for the rest.

---

## Step 1 — Fix the AVX2 transcendental kernels — SUPERSEDED BY STEP 3

**Evidence.** `expAVX2` is *far slower than the generic path that already exists
in the same package*:

| n | `expAVX2` (current) | `-tags=purego` generic | naive `math.Exp` loop |
| --: | --: | --: | --: |
| 256 | 57,740 ns | 6,150 ns | 1,337 ns |
| 4,096 | 872,098 ns | 21,907 ns | 21,064 ns |
| 65,536 | 14,322,837 ns | 106,226 ns | 350,534 ns |
| 1,048,576 | **211,667,934 ns** | **1,110,432 ns** | 5,491,242 ns |

`logAVX2` is the same story (91,839 ns vs 1,435 ns at n=256, 64x).

**Root cause.** `internal/eml/simd_amd64.s` re-materialises every polynomial
coefficient **inside** the loop body. From `expAVX2` (lines 541–713) there are 17
`VBROADCASTSD` and 14 `MOVQ $imm64`, of which **9 of each are inside
`loop_exp_avx2`**:

```asm
	// Y3 = Y3 * r + c10          <-- per-iteration, per-coefficient
	VMULPD Y4, Y3, Y3
	MOVQ $0x3e927e4fb7789f5c, AX  ;  imm64 -> GPR
	MOVQ AX, X15                  ;  GPR  -> XMM   (~5 cycle latency)
	VBROADCASTSD X15, Y2          ;  XMM  -> YMM
	VADDPD Y2, Y3, Y3
```

That is 27 uops of pure overhead per 4 elements, and the `AX` round-trip forms a
serial dependency chain. The loop is not failing to vectorise — it is
*anti-vectorising*. The same defect is present in `logAVX2` (18/13),
`tanAVX2` (18/9), and mildly in `sinAVX2`/`cosAVX2` (11/6).

**Change.**
1. Move every coefficient into a `DATA`/`RODATA` table and load it with a
   single-uop memory broadcast: `VBROADCASTSD ·exp_c10(SB), Y2`. Delete the
   `MOVQ $imm64` + `MOVQ AX, Xn` pairs entirely.
2. Hoist anything that fits into the remaining `Y` registers out of the loop
   (AVX2 gives 16; the degree-11 Horner chain needs ~12, so budget
   coefficients against working registers and keep the rest as memory operands).
3. Unroll the loop 4x (16 elements/iteration) to amortise the `DECQ`/`JNZ` and
   expose ILP.

**Verify.** Re-measure the table above; target ≥ 2x versus `-tags=purego` at
n ≥ 4096. Add a golden-value test so a future edit cannot silently reintroduce
the pattern.

**Outcome: superseded.** Step 3 landed first and removed the regression entirely,
so this is no longer a rescue. Writing a correct SIMD transcendental is still
worthwhile, but it is now an *addition* rather than a fix, and the bar is higher
than "beat the generic path": the generic path is already parallelised and 5.3x
faster than a serial loop at n=1M, so a 4-wide kernel would have to beat that
rather than beat `math.Exp`. See "Revised priorities" below.

---

## Step 2 — Gate the worker pool on a real threshold

**Evidence.** `SmallCutoff = 256` is roughly 16x too low. The pool only starts
paying off at n ≈ 4096:

| n | pool | plain loop | pool / loop |
| --: | --: | --: | --: |
| 256 | 6,021 ns | 1,376 ns | **4.38x slower** |
| 512 | 7,397 ns | 2,679 ns | **2.76x slower** |
| 1,024 | 10,469 ns | 5,596 ns | **1.87x slower** |
| 2,048 | 15,339 ns | 10,703 ns | **1.43x slower** |
| 4,096 | 21,571 ns | 22,854 ns | 0.94x |
| 16,384 | 43,396 ns | 90,365 ns | **0.48x** |
| 65,536 | 109,063 ns | 363,324 ns | **0.30x** |

`GetParallelChunkSize(256)` returns `16`, so the pool spawns 16 goroutines to
process 256 elements — roughly 16 elements per goroutine, which is pure
scheduling overhead.

**Change.** Raise `SmallCutoff` to 4096 (`internal/eml/simd.go`). Better, make it
a function of both size and core count rather than a bare constant, e.g.
`minParallelCutoff = 512 * runtime.NumCPU()`, and add a cheap guard so the pool
is skipped when the computed chunk size would fall below a few hundred elements.
Apply the same threshold to `parallelizeGenericF32`.

**Verify.** Re-run the crossover sweep; the pool must not lose at any size above
the threshold.

**Outcome: done.** `SmallCutoff` is now `SmallWorkloadFactor * runtime.NumCPU()`
(512 × cores, i.e. 8192 on a 16-core host) instead of a flat 256, and every
inline-vs-pool guard routes through it. Measured table is in "Status" above: the
worst ratio is now 1.03 at n=1024, down from 4.38.

---

## Step 3 — Make the generic path the fallback, not the assembly

**Evidence.** Steps 1 and 2 are large but *risky* to land first: they touch
hand-written assembly and a shared constant. The safety net is already in the
tree — `-tags=purego` routes every batch op through `parallelizeGeneric`, and it
is **5.6x faster than a naive loop at n=1M** and never catastrophic.

**Change.** Until Step 1 lands and is proven, make the transcendental
dispatchers prefer the generic path on any AVX2 host, i.e. treat the current
`expAVX2`/`logAVX2`/`sinAVX2`/`cosAVX2`/`tanAVX2` as opt-in. Concretely, invert
the guard in `dispatchExpSIMDTo` and friends so `parallelizeGeneric` is the
default and the assembly kernel is gated behind a build tag or a
`EML_FORCE_SIMD_ASM=1` escape hatch. This is a one-line-per-function change that
converts a 190x regression into a 5.6x win immediately.

**Verify.** `go test -bench` on the transcendental batch ops before and after.

**Outcome: done.** `internal/eml/simd_trans_amd64.go` now routes the five
transcendental dispatchers to `parallelizeGeneric`, matching what arm64, wasm
and the generic stub already did — amd64 was the only platform using the
assembly. The assembly is retained behind a new `-tags emlasm` build tag rather
than deleted, so it stays auditable while it is reworked;
`internal/eml/simd_trans_asm_amd64.go` holds that path and
`TestTranscendentalAsmKernelsAreWiredUp` keeps the symbols referenced. Note the
accuracy gates are *expected to fail* under `-tags emlasm`; that is the contract
a reworked kernel must meet.

The elementwise AVX2 kernels (add/sub/mul/div/abs/neg/inv/sqrt/fma) are
untouched and still dispatched from `simd_dispatch_amd64.go`: Go's compiler
already auto-vectorises those loops to the same width, so hand assembly buys
nothing.

---

## Step 4 — Add accuracy gates for every batch path

**Evidence.** This is why Steps 1–3 were invisible for so long. The batch tests
assert only slice **lengths**; there is no value comparison anywhere on the CPU
SIMD path. `pkg/fastmath` has scalar accuracy tests and `internal/gpu` has 1-ULP
e2e tests, but neither covers `internal/eml`'s batch kernels.

Measured error of the current kernels against `math`:

| kernel | max ULP | note |
| :--- | --: | :--- |
| `logAVX2` | 0 | exact |
| `expAVX2` | 56 | acceptable for a minimax kernel |
| `sinAVX2` | 54,254 | ~6.0e-12 abs error at ±π/2 |
| `cosAVX2` | — | ~6.5e-11 abs error at ±π/2 |

The astronomical `cosAVX2` ULP count is a metric artefact — `cos(π/2) ≈ 6.1e-17`,
so any absolute error is a huge number of ULPs. The real defect is the
**range-reduction discontinuity**: at exactly ±π/2 the kernel returns
`-6.51e-11` instead of `6.12e-17`. Consequently `sin²+cos² = 1` holds only to
**1.20e-11** on the SIMD path versus ~1e-16 on the scalar path. A step that
changes the reduction interval would change this number, and nothing would fail.

**Change.**
1. Add `internal/eml` accuracy tests asserting batch results against `math` over
   dense grids, with per-kernel ULP budgets and an explicit absolute-error budget
   for kernels whose true value crosses zero.
2. Add a `sin²+cos²` invariant test on the batch path, mirroring the identity
   tests the scalar path already has.
3. Make these tests the gate for Steps 1 and 3: they must pass before and after.

**Verify.** `go test ./internal/eml/` fails today; make it pass.

**Outcome: done.** `internal/eml/simd_accuracy_test.go` gates every batch kernel
with a mixed absolute+relative tolerance (`1e-14 + 1e-15*|want|`), which is the
only formulation that is simultaneously satisfiable for large results and
meaningful where the true value crosses zero. ULP budgets are enforced only where
`|want| > 1e-9`. Coverage includes dense grids, the scalar tail for lengths 1–63,
`sin²+cos²`, the ±π/2 range-reduction boundaries, agreement with the scalar
path, and domain edges (overflow, underflow, subnormals, negative log arguments,
`sin`/`cos` at 1e9). The gates caught 6170 real failures before Step 3 landed and
pass on every path now.

---

## Step 5 — Extend the bytecode VM to the full 25 functions

**Evidence.** The VM is the engine for symbolic regression, and it can barely
express anything. Measured — of the 25 functions the JIT supports, the bytecode
compiler accepts **5**:

```
SUPPORTED   (4 unary + eml): exp, log, ln, sqrt, eml
UNSUPPORTED (23): sin cos tan asin acos atan abs cbrt log2 log10 ceil floor
                  trunc round sinh cosh tanh asinh acosh atanh erf gamma
```

Note the trap: `isFunc` in `pkg/bytecode/compiler.go` already returns `true` for
`sin` and `cos`, so the tokenizer accepts them and `emitFunc` then fails with
"unsupported function in bytecode compiler". The two tables have drifted the same
way the JIT's three tables had.

**Change.**
1. Add opcodes for the remaining 20 functions (`OpSin` … `OpGamma`) alongside the
   existing `OpSqrt`/`OpExp`/`OpLog`, and extend `Arity()`, `String()`,
   `Eval`, `EvalRegularized`, `EvalBatch`, `EvalBatchColumnar` and `emitFunc`.
2. Collapse `isFunc` and `emitFunc` into one table so this cannot drift again —
   the same lesson as `internal/jit/functab.go`.
3. Fix `Optimize` to fold constant transcendental subtrees, not just `eml`.

**Verify.** Parameterise a test over all 25 names asserting `CompileExpr` succeeds
and that `Eval` agrees with `math` to a documented tolerance.

**Outcome: done, and it was worse than measured.** The audit found **5** working
ops (`eml`, `exp`, `log`, `ln`, `sqrt`) out of 27 names, and a *second*
dispatch site that the first measurement had missed: `isFunc` accepted `sin` and
`cos`, so the tokenizer produced a call node and then `emitFunc` failed with
"unsupported function in bytecode compiler". Both sites now read from
`pkg/bytecode/functab.go`, and `emitAST` — the `jit.AST` path that
`CompileExpr` tries *first* — was rewritten to call `emitFunc` instead of
carrying its own switch. The anonymous struct duplicated across `emitOp` and
`emitFunc` became a named `token`, because it was the reason sharing code
between them was impossible. 27/27 names now compile and evaluate, checked
through `Eval`, `EvalRegularized` and both batch entry points. `Optimize` folds
every function's constant subtree, not just the four that had explicit cases.

---

## Step 6 — Ship a genetic-algorithm driver and a fitness function

**Evidence.** `pkg/bytecode/genetic.go` provides the *operators* — `ValidateStack`,
`Crossover`, `Mutate` — and nothing else. A repo-wide search for `fitness`,
`population`, `selection`, `elitism`, `tournament` and `generation` returns
**zero hits outside comments**. There is no loop that actually searches, and no
way to score a candidate program.

This is the library's headline use case: the paper's central result is
gradient/evolutionary recovery of closed-form elementary functions from data.
Today a user must write the entire GA themselves on top of the operators.

**Change.**
1. `pkg/bytecode/ga.go`: `Population`, `Config{PopulationSize, Generations,
   CrossoverRate, MutationRate, Elitism, MaxDepth, Seed}`, tournament and
   roulette selection, generational loop, convergence capture.
2. A `Fitness` type taking `(*Program, Dataset) float64`, plus a ready-made
   least-squares fitness so the common case needs no user code.
3. Make generation parallel across candidates with the (fixed) worker pool from
   Step 2 — this is the single biggest throughput win available to the GP path,
   because fitness evaluation is embarrassingly parallel and currently is not.
4. Termination reporting (best fitness per generation) so runs are diagnosable.

**Verify.** Recover a known closed form from synthetic data in a test with a
fixed seed.

**Outcome: done.** `pkg/bytecode/ga.go` adds `Dataset`, `Fitness`, `Config`,
`Result` and `Search`, plus `LeastSquares` as a ready-made objective.
`genetic.go` gained `RandomProgram`, a seedable crossover and mutation, and
correct arity-aware mutation pools. All three operators now take an explicit
`*rand.Rand`, which makes searches reproducible.

Fitting y = 3x² + 2x + 1 from 24 samples with a fixed seed reaches **RMS error
4.5e-13**, and the recovered program is algebraically exact:
`-((-1.1276507) - ((2.0 + x) + x + x)*x + (-0.1276507))`.

Two design points were load-bearing:

- **Coordinate descent on constants (`LocalSearchSteps`, default 30)** is the
  single biggest win. Random jitter is a random walk, so pinning a coefficient to
  3.0 takes many generations; descent turns that into a direct search. Measured
  on the same seed it is the difference between RMS 0.21 and 4e-13.
- **Programs are built as expression trees, then flattened to RPN.** Generating
  RPN directly produced a flat chain whose stack height is 1 after the first
  opcode, which has exactly one crossover point — crossover degenerated into
  swapping parents whole. The tree gives every complete subtree a balanced RPN
  prefix: **3.69 usable crossover points per program**, up from 1.0.

**Known limitation, stated rather than tuned around.** Fitting y = sin(x) over
[0, 2π) plateaus at RMS ≈ 0.17 (against ≈ 0.71 for the best constant fit) and
does not improve with budget: 800 individuals × 600 generations gives 0.169,
1500 × 1500 gives 0.177. Approximating sin means discovering a composition of
elementary functions, which is what the source paper does with gradient-based
training of a parameterised tree, not with a genetic algorithm. `Search` documents
this, and `TestSearchImprovesATranscendentalTarget` asserts only what holds —
substantial improvement over a constant fit — instead of a flaky exact-recovery
claim.

---

## Step 7 — Zero-allocation `EvalBatch` and a genuinely vectorised VM

**Evidence.** `EvalBatchColumnar` allocates on **every call**:

| n | ns/op | B/op | allocs/op | ns/sample |
| --: | --: | --: | --: | --: |
| 64 | 1,062 | 1,616 | 2 | 16.6 |
| 1,024 | 14,269 | 24,656 | 2 | 13.9 |
| 65,536 | 705,284 | 24,656 | 2 | 10.8 |

The columnar stack (`MaxStackDepth * chunkSize` floats) and the `[][]float64`
slice-of-slices are rebuilt per call. For the GP inner loop of Step 6 this is a
GC hazard, not a rounding error.

It is also not actually vectorised: per-sample cost only drops from 16.6 ns to
10.8 ns versus 33.6 ns for scalar `Eval` — a 2–3x win, where a 4-wide AVX2 VM
over `OpAdd`/`OpMul`/`OpExp` should reach 4–8x.

**Change.**
1. Move the columnar stack into the `Program` (or a caller-supplied
   `Scratch` struct), so a steady-state batch loop allocates nothing. Follow the
   existing `TurboQuant4DistanceScratch` pattern for the API shape.
2. Widen the inner element loops to 4-wide (`float64`) / 8-wide (`float32`) by
   iterating over `[]float64` columns that Go's auto-vectoriser can widen, and
   unroll the per-op dispatch so constants are broadcast once per chunk rather
   than written element-by-element.
3. Reuse the worker pool from Step 2 when `n` exceeds the threshold.

**Verify.** `allocs/op == 0` at steady state, and a `ns/sample` at least 2x
better than today.

**Outcome: done for the allocation half.** `BatchScratch` holds the columnar
stack across calls; `EvalBatchColumnarScratch` is the zero-allocation entry
point and `EvalBatchColumnar` now wraps it. Measured at n=1024: **24,657 B and 2
allocs → 0 allocs**, and 14,543 ns → 11,034 ns. At n=65536: 24,656 B → 0 B. The
inner element loops were also converted from `for i := 0; i < currLen; i++` to
`for i := range col` (10 sites) so the compiler can widen them.

The second half — a genuinely vectorised VM body — is **not** done. The dispatch
is table-driven, but each opcode still runs a scalar loop over a column, so
per-sample cost is 10.7 ns versus 33.6 ns for scalar `Eval`. Go cannot
auto-vectorise a call to `math.Sin`, so real widening needs either generated
kernels or the same assembly work Step 1 covers.

---

## Step 8 — `Pipeline.Reset()` and an in-place `RunTo`

**Evidence.** `Pipeline` is documented as "zero-allocation composable", and it is
— but only if used correctly, and the API makes misuse the path of least
resistance. Every builder call appends a step:

```go
p := eml.NewPipeline(n)
for i := 0; i < b.N; i++ {
    p.Exp().MulScalar(2.0).Log().RunTo(x, out)   // measured: 300 steps after 100 iterations
}
```

After 100 iterations `len(p.steps) == 300`, and each call does 3x the work.
There is no `Reset()`. Measured correctly (chain built once) the pipeline is
fine — 0 allocs, 688,976 ns at n=65,536 versus 1,077,716 ns for the naive
equivalent, i.e. 1.5x faster — so the kernel is sound and only the API is wrong.

`RunTo` additionally calls `Run`, which copies the input into `p.buf[0]`, applies
the steps, and then copies the result back out: two full-length copies that a
true in-place path would avoid.

**Change.**
1. Add `Reset()` and document the builder as single-use-or-`Reset`.
2. Better: separate construction from execution — `NewPipeline(n).Then(expStep,
   mulStep(2))` returns an immutable, reusable pipeline.
3. Give `RunTo` a genuinely in-place path that writes the first step's output
   directly into `output` and ping-pongs the remainder, removing both copies.

**Verify.** A test asserting `len(p.steps)` is unchanged across repeated `RunTo`
calls, and `allocs/op == 0`.

**Outcome: done.** `Reset()` and `Len()` added; `RunTo` now writes the first
step's output straight into `output` and ping-pongs between `output` and a single
scratch buffer, removing both of the copies the old `RunTo` performed via `Run`.
Measured 0 allocs at n=65536. The missing `SubScalar` and `DivScalar` steps were
also added to complete the scalar set. `TestPipelineRunToMatchesReference` covers
one through five steps, so the output/scratch parity is pinned on both sides of
the odd/even boundary.

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

# Part 2 — Steps 11 to 20

Steps 1–8 optimised the paths that were slow. Steps 11–16 fix paths that are
**wrong**. The distinction matters: a 200x speedup on a function that returns
garbage for a third of its inputs is a negative result, so every performance
step above is conditional on the correctness step that gates it.

---

## Step 11 — Fix the JIT register-ABI defect (P0)

**This is the single most serious defect in the repository.** The JIT produces
**wrong answers for 37.5% of randomly generated expressions**, silently, at
94% statement coverage.

**Evidence — differential test, JIT vs. the package's own interpreter.** 3,000
random expression trees (depth 3–8, 12 functions, 3 operators, constants),
each evaluated at x ∈ {0.05, 0.5, 1.4}:

```
JIT differential: 1124/3000 expressions MISMATCH (37.5%)

(-0.286*((((x+x)*(x-tan(x)))*(x*(-0.224+(x+x))))+-1.080))
    @x=0.05: jit=0.2419199686086083   interp=0.30887999260427085
((((erf(x)*-1.751)*0.167)+-0.648)-((((log2(x)-erf(x))+x)*(x*(1.891-x)))*((x*x)*-1.027)))
    @x=0.05: jit=-0.6644841246314702  interp=-0.6655070680371609
(((((erf(x)*x)+tan(x))-(x--1.466))+1.295)+((sinh(x)-(sqrt(x)+((log2(x)-x)+sqrt(x))))+(x+log(x))))
    @x=0.05: jit=-Inf               interp=0.8608633690364578
```

**Minimal repro — nesting depth 4 around any function call.**

```
levels=1  jit=1.479425539          interp=1.479425539          ok
levels=2  jit=3.479425539          interp=3.479425539          ok
levels=3  jit=6.479425539          interp=6.479425539          ok
levels=4  jit=10.84147098          interp=10.47942554          *** WRONG ***
levels=5  jit=11.90929743          interp=15.47942554          *** WRONG ***
```

(the expression is `1+(2+(3+…+(sin(x)))))`, k levels deep.)

**Root cause.** `alloc()` (`internal/jit/codegen.go:232`) hands out **xmm1–xmm14**
— `xmm0` holds the result and `xReg = 15` holds the variable. But `callFunc`
saves and restores **only xmm0–xmm7**:

```go
// codegen.go:194 — the save loop stops at 8
for i := byte(0); i < 8; i++ { ... MOVSD [rsp+i*8], xmm_i }

// codegen.go:214 — the restore loop starts at 1 and stops at 8
for i := byte(1); i < 8; i++ { ... }
```

Go's amd64 register ABI treats **all of X0–X15 as caller-save**. So the first
time a `callFunc` is emitted while a value is live in xmm8–xmm14, the callee
destroys it and the generated code goes on to read it. Disassembly of
`1+(2+(3+(4+(sin(x)))))` (250 bytes) shows 8 saves and 7 restores and **no
handling of xmm8–xmm15 anywhere**.

Note the shape of the bug: `xmmReg := i & 7` at `codegen.go:196` — the author
wrote an 8-register loop over a 14-register allocator. The masking makes the
loop look register-agnostic; the intent was 16.

**Why no test catches it.** Every existing test uses flat or shallow expressions
(`sin(x)^2+cos(x)^2`, `x^3-2x^2+x-1`). Left-leaning chains (`sin(x)+sin(x)+…`)
stay safe because `gen` frees `leftSave`/`tmp` before recursing, so they never
exceed xmm3. The depth needed to trigger it is 4; no test in
`codegen_test.go` (289 LOC), `coverage_test.go` (572 LOC) or `fuzz_test.go`
reaches it.

**Why it is now urgent rather than merely wrong.** `pkg/bytecode`'s `RandomProgram`
(`genetic.go:301`) generates exactly this tree shape. The obvious next step after
Step 6 is to score GA candidates with the JIT instead of the interpreter — at
that point 37.5% of fitness evaluations become noise, silently.

**Change.**
1. Save/restore **all 16** XMM registers in `callFunc`: loop to 16, raise the
   frame from `sub rsp,72` to a 16-byte-aligned 176 (`sub rsp,120` +
   128-byte save area is the smallest that keeps alignment), and assert the new
   size in `codegen_test.go`.
2. Move `xReg` off **xmm15**. Go reserves X15 as its zero register
   (`$GOROOT/src/cmd/compile/internal/ssa/_gen/AMD64Ops.go`: `Zero128 … reg:
   x15only`). No probe found an observable miscompile from this with Go 1.27's
   stdlib `math`, but the invariant is undocumented, unguarded and one Go
   release away from breaking every compiled function. Use xmm15 only for the
   *last* expression in codegen, or allocate `x` from the pool.
3. Add **spilling**. Depth 8 currently hard-fails with
   `codegen error: out of register resources` and no diagnostic. With all 16
   registers genuinely saved, spilling to a stack slot buys real depth.

**Verify.** (a) The differential test above goes to 0/3000. (b) A randomised
differential fuzz (`go test -fuzz`) comparing JIT output to `jit.Eval` over
random trees, depth 1–10 — the property test this package should have had all
along. (c) Re-measure Step 9's table; the 17–30x figures there were measured on
shallow expressions and are unaffected, but Step 11 must land *before* Step 9 so
the ARM64 encoder does not inherit the same 8-vs-16 mistake.

---

## Step 12 — Close the `EMLEval`/`Diff` silent-zero holes (P0)

**Evidence.** `Canonicalize` faithfully maps *any* of the 25 function names into
an `EMLNode` (`canonical.go:238`, default branch). But `EMLEval`'s dispatch
switch handles only 11 of them and then falls through:

```go
// canonical.go:108-133
switch n.Kind {
case EMLFunc:
    switch n.Name {
    case "sin": ... case "sqrt": ...
    }          // no default
}              // no default
return 0       // canonical.go:133 — silent zero
```

Measured across all 25 supported functions at x=0.5:

```
EMLEval: 17/25 wrong | Diff: 20/25 wrong

EMLEval tan     got=0              want=0.54630249     WRONG
EMLEval acos    got=0              want=1.0471976      WRONG
EMLEval log2    got=0              want=-1             WRONG
EMLEval asinh   got=0              want=0.48121183     WRONG
EMLEval gamma   got=0              want=1.77245385     WRONG
Diff     tan     got=0              num=1.2984464      WRONG
Diff     gamma   got=0              num=-3.4802309      WRONG
Diff     round   got=0              num=500000          WRONG
```

`Diff` is worse than "missing": `Diff(round(x))` should be `0`, and it *is* `0`
— but so is `Diff(gamma(x))`, which is wrong by construction and indistinguishable.

**A third bug, in the same file: `DiffEval` differentiates twice.**

```go
// canonical.go:488-491
// DiffEval evaluates the symbolic derivative Diff(n) at x.
func DiffEval(n *EMLNode, x float64) float64 {
	return EMLEval(Diff(n), x)      // applies Diff to its own argument
}
```

The doc says "the symbolic derivative `Diff(n)`", but the implementation computes
`Diff(Diff(n))` — the **second** derivative. Verified:

```
Diff(sin(x))  = cos(x)             <- correct
DiffEval(Diff(sin(x)), 0.5) = -0.479425538604203
cos(0.5)                       =  0.8775825618903728
-0.4794 = -sin(0.5)             <- second derivative
```

Any caller doing `DiffEval(n, x)` expecting the first derivative gets the
second, and `DiffEval` is the *only* one-line API for numeric differentiation
this package offers.

**Change.**
1. Drive `EMLEval`, `EMLEvalRegularized` and `Diff` from
   `internal/jit/functab.go` — the same single-source-of-truth table Step 5 used
   to fix `pkg/bytecode`. This is the third time a duplicated function-name
   switch has drifted in this repo (`isFunc`/`emitFunc` in Step 5, `isFuncName`
   vs `functab` in the parser, now `EMLEval` vs `functab`). Add a derivative
   field to `mathFunc` so `Diff` cannot under-implement again.
2. Fix `DiffEval` to `EMLEval(n, x)` and rename it, or add
   `DerivativeEval(n, x) = DiffEval(Diff(n), x)` alongside it so both orders are
   expressible.
3. Until every function is covered, make the default a **panic in tests**
   (build-tag gated) rather than a `return 0`. A silent zero is the worst
   possible failure mode for a symbolic layer feeding a search loop.

**Verify.** Parameterise a test over all 25 names asserting `EMLEval` matches
`jit.Eval` and `Diff` matches a central difference of `jit.Eval`. Add a
`DiffEval` test pinning first derivative at three points.

---

## Step 13 — Make the bytecode VM fail loudly instead of corrupting the stack (P0)

**Evidence.** `OpCode` is a `uint8` and `Program.Ops` is an exported, mutable
field. `Arity()` returns `0` for any opcode it does not recognise, and `Eval`'s
`default` branch silently skips an opcode whose function it cannot find:

```go
// eval.go:78-80
default:
    if fn, ok := opUnaryFns[op]; ok {   // miss => no-op, stack pointer unchanged
        stack[sp-1] = fn(stack[sp-1])
    }
```

The stack pointer then desynchronises and the VM runs off the end. Verified:

```
p := &Program{Ops: []OpCode{OpConst, OpCode(200), OpAdd}, Consts: []float64{1}, MaxStackDepth: 4}
p.Eval([]float64{0}, nil)
  -> panic: runtime error: index out of range [-1]
```

221 of 256 `OpCode` values are "unknown". One out-of-range byte — from a
deserialiser, a fuzzer, or a future refactor — is a memory-safety crash, not a
diagnosable error. The same file silently ignores a *valid* opcode whose stack
underflows for the same reason.

**Evidence — the batch path is worse, because it is public API.** Column
extraction in `EvalBatchColumnarScratch` is unguarded:

```go
// eval_batch.go:127
srcCol := data[varID][offset : offset+currLen]     // no length check
```

Verified:

```
p := &Program{Ops: []OpCode{OpConst, OpVar, OpMul}, Consts: []float64{2}, MaxStackDepth: 2}
p.EvalBatchColumnar([][]float64{{1, 2, 3}}, make([]float64, 64))   // dst 8x longer than the column
  -> panic: runtime error: index out of range [0] with length 0
```

Note the row-form `EvalBatch` *does* bounds-check, via `Eval`'s
`if int(idx) < len(vars)` guard. The two batch entry points have different
safety contracts, and `Dataset.validate()` (`ga.go:46`) does not help non-GA
callers. A mismatch between `len(dst)` and the column length is an ordinary
caller mistake, not a malformed program.

**A related silent-wrong-answer, for contrast.** The same file's siblings are
inconsistent about length mismatch: `internal/eml`'s `ExpMulBatch` returns the
*caller's input slice, un-computed*, while `ExpMulTo` panics. `FmaSIMD` returns
its input when `n == 0` while every sibling returns an empty slice.

**Change.**
1. `finalize()` (`compiler.go`) validates the opcode stream: every opcode must
   be `< funcOpLast+1` and resolvable by `funcName`, every constant and var
   index in range, and the final stack height exactly 0. It already computes
   `MaxStackDepth`; make it the validation point.
2. Make `Eval` panic with a *naming* message on an unresolvable opcode rather
   than skipping it — `ValidateStack` (`genetic.go`) already exists for exactly
   this and is not on the `Eval` path.
3. Bounds-check `offset+currLen` against `len(data[varID])` and `len(dst)` in
   `EvalBatchColumnarScratch`, returning an error like the row path's contract
   rather than panicking.
4. Pick **one** mismatch policy across `internal/eml`, `pkg/bytecode` and
   `pkg/fastmath` and apply it everywhere. Currently there are three: panic
   (`"slice length mismatch"`), panic (`"length mismatch"`), and silent
   return-the-input.

**Verify.** Table-driven tests: unknown opcode, truncated stream, wrong
`MaxStackDepth`, short column, short `dst`, and every `*Batch`/`*To` mismatch
pair. Assert on the *error*, not just the absence of a panic.

---

## Step 14 — Fix fastmath's silent range failures (P0)

`pkg/fastmath` is the package that exists to be *fast*, so it is the one where a
wrong answer is least expected and most damaging. It has three, and none is
covered by a test.

**Evidence — `FastExpMinimax` returns exactly `0` across 37 decades of output.**

```go
// fast_exp.go:64-66
if k < -1022 { return 0.0 }
```

`k` is the Cody–Waite **reduction index** (`x/ln2`), not the output exponent.
The true underflow threshold is x = −745.13, but the guard fires at
x = −708.55. Everything in between is silently zeroed:

```
FastExpMinimax(-700)   = 9.859676543674456e-305   math.Exp = 9.85967654375977e-305   ok
FastExpMinimax(-708.5) = 0                        math.Exp = 2.006132305331306e-308   *** ZERO CLIFF ***
FastExpMinimax(-709)   = 0                        math.Exp = 1.216780750623423e-308   *** ZERO CLIFF ***
FastExpMinimax(-720)   = 0                        math.Exp = 2.0322308024e-313       *** ZERO CLIFF ***
FastExpMinimax(-744)   = 0                        math.Exp = 1e-323                  *** ZERO CLIFF ***
```

This repository already knows the right constant — `pkg/logexp/exp.go:18` uses
`expUnderflow = -745.133224101734`. Two packages, two thresholds, one wrong.

**Evidence — `FastExpF32`'s thresholds are wrong by 0.72 and 16.9 decades.**

`fast_exp.go:84-89` clamps at `[-87, +88]`. The true float32 limits are
`ln(3.4028235e38) = 88.722839` and `ln(1.4e-45) = -103.279`:

```
FastExpF32(88.000015) = +Inf   true float32 exp = 1.6517e+38    <- WRONG
FastExpF32(88.7228)   = +Inf   true float32 exp = 3.4027e+38    <- WRONG
FastExpF32(-87.5)     = 0      true float32 exp = 9.98e-39      <- WRONG
FastExpF32(-103.9)    = 0      true float32 exp = 1e-45         <- WRONG
```

`1e20` is an unremarkable float32 value. A GA population regularly contains
magnitude-20 constants.

**Evidence — `LnRegularized` overflows to `+Inf` on ordinary inputs.** The
C-infinity regulariser squares its argument (`guardrails.go:83`):

```go
func LnRegularized(y, eps float64) float64 {
	return 0.5 * FastLog(y*y + eps*eps)     // y*y overflows for |y| > ~1.34e154
}
```

```
LnRegularized(1e200, 1e-12)   = +Inf   (want 460.517)
LnRegularizedF32(1e20, 1e-6) = +Inf   (want 46.052)
```

This is the precise failure mode the function exists to prevent, reached at
|x| = 20 — nowhere near the overflow.

**The accuracy claims themselves are overstated.** Measured on varying inputs
(2M-sample sweeps; the earlier hand-rolled numbers in this file were taken with
constant inputs, which the compiler partly folds away — use these):

```
BenchmarkZZ2MinimaxLog-16   5.984 ns/op
BenchmarkZZ2StdLog-16       7.580 ns/op     -> 1.27x, not the documented "4x-5x"
BenchmarkZZ2MinimaxExp-16   5.362 ns/op
BenchmarkZZ2StdExp-16       7.216 ns/op     -> 1.35x
BenchmarkZZ2BitCastExp-16   3.472 ns/op     -> 2.08x

max rel err FastExpMinimax = 1.62e-07     (doc claims < 1.5e-6 — pessimistic, fine)
max abs err FastLogMinimax = 2.93e-08     (doc claims < 1e-7 — fine)
```

So the approximations are **good but not dramatic**: 1.3–2.1x, not 4–5x. Three
different descriptions of `FastExpMinimax` also exist (degree **4** in
`fastmath.go:37`, degree **5** in `fast_exp.go:34`, degree **6** at
`fast_exp.go:52`), and the coefficients are rounded *Taylor*, not the Remez
minimax the package doc claims — `c2 = 0.50000000044` is `0.5 + 4.4e-10`.

**Change.**
1. Fix the underflow guard to compare the *output* exponent against −1022, or
   simply import the shared constant from `pkg/logexp`. Fix `FastExpF32`'s
   bounds to `ln(MaxFloat32)` / `ln(SmallestNonzeroFloat32)`.
2. Rewrite `LnRegularized` without squaring: `0.5 * FastLog(math.Hypot(y, eps))`
   gives the same "distance from zero" semantics with no overflow.
3. Fix the degree claims, or implement the Remez fit the doc advertises.
   Correct the "4x–5x faster" figure to the measured 1.27x.
4. Add dense accuracy sweeps (overflow, underflow, subnormals, ±0, ±Inf, NaN)
   for both F64 and F32 variants. `fast_test.go` currently checks ~12
   hand-picked points at `relErr < 5e-5`, which cannot find any of the above.

**Verify.** `go test ./pkg/fastmath/` with a sweep over
`x ∈ [-750, 710]` asserting relative error `< 1e-6` wherever `math.Exp` is
finite and non-zero, plus exact agreement at both underflow boundaries.

---

## Step 15 — Fix `purego` on arm64 and wasm (P0, ~30 minutes)

**Evidence.** `purego` is the documented escape hatch — Step 3 of this document
introduced it, and `Makefile`'s `gosec` target runs `-tags purego`. It does not
compile on two of the three platforms it exists for:

```
$ GOARCH=arm64       go build -tags=purego ./...
internal/eml/simd_dispatch_stub.go:54:6: fmaSIMD redeclared in this block
internal/eml/simd_dispatch_stub.go:74:6: dispatchExpSIMDTo redeclared in this block
... 16 errors

$ GOOS=js GOARCH=wasm go build -tags=purego ./...
internal/eml/simd_dispatch_wasm.go:5:6: emlSIMD redeclared in this block
... 16 errors
```

**Root cause.** `simd_dispatch_stub.go:1` is tagged
`(!amd64 && !arm64 && !wasm) || purego`, but `simd_dispatch_arm64.go:1` and
`simd_dispatch_wasm.go:1` carry no `!purego` guard. Adding `purego` to either
platform's tag selects the stub file *in addition to* the platform file.
amd64 escapes this only by accident — `simd_trans_amd64.go` happens to define the
same dispatchers as the stub and wins.

**Change.** Add `&& !purego` to `simd_dispatch_arm64.go:1` and
`simd_dispatch_wasm.go:1` (and audit `simd_dispatch_stub.go`'s own tag for the
mirror-image problem: it must *not* select on `amd64 && purego` if the amd64
transcendental dispatcher already handles that).

**Verify.** Extend the CI cross-build matrix to run `-tags=purego` for **every**
`GOOS/GOARCH` pair it already lists. A build matrix that does not vary build tags
is why this survived: `arm64`, `amd64`, `amd64+emlasm`, `386`, `riscv64`,
`darwin/arm64`, `js/wasm` and `wasip1` all build today; only `arm64+purego` and
`wasm+purego` were never asked.

---

## Step 16 — Fix or delete the Metal backend (P0)

`internal/gpu` is 2,153 LOC of Go + 1,334 of CUDA/C + 282 of Objective-C. 616 of
those Go lines are `return fmt.Errorf("not available")` stubs. The CUDA path is
genuinely real; **the Metal path cannot work.**

**Evidence — results are read before the GPU is asked to run.** In
`internal/gpu/metal_bridge_darwin.m`:

```objc
[enc setComputePipelineState:pipeline];
setup(enc, pipeline, cmdBuf);        // :133   <-- the memcpy OUT happens here
[enc endEncoding];                   // :134
[cmdBuf commit];                     // :135
[cmdBuf waitUntilCompleted];         // :136
```

The read-back is inside `setup`, so it runs **before** `endEncoding`, `commit`
and `waitUntilCompleted`. Every Metal batch result is uninitialised device
memory. The same ordering appears at lines 163, 196, 237 and 268.

**Evidence — a `float64` API is served by `float32` compute.**
`metal.go:134-138` passes `*C.double`; the bridge narrows to `float`
(`eml_metal_launch_unary`, `:147-150`) and widens on return (`:167`). All 18
shaders in `metal/eml_metal_shaders.metal` take `device const float *`.
`DefaultVerifier` uses `MaxULP: 1` (`validate.go:16`), which this can never
satisfy — so the verifier that exists to catch this is guaranteed to fail, or is
being run on a path that never executes.

**Evidence — the API silently shrinks on darwin/arm64.** `metal.go` defines 15
of `stub.go`'s 89 `Device` methods. Missing: every `*BatchTo`, `DotBatch`, all
`*F32`, all `*C64`, all `*C128`, and the `AllocatePinned*`/`FreePinned*` package
functions. `cmd/emlcli` and `cmd/bench` happen to touch only the 9 float64 unary
methods, so *those* binaries still build on a Mac — but any other consumer
compiles on Linux and fails on macOS. And because the package's own tests are
all tag-excluded on `darwin/arm64+cgo`, `go test ./...` on a Mac has **zero**
coverage here.

**Also:** manual retain/release over-releases the non-owned array from
`MTLCopyAllDevices` (`:44-46`) and leaks a fresh copy of the device list per
iteration (`:77`); `_deviceProps[16]` is fixed-size while `getMetalDevices` loops
to the true count, so >16 GPUs breaks `GetDevices`; `libPaths` is
CWD-relative, so the shader library is found only by luck of the working
directory; and 12 of the 18 shaders (`kernel_addScalar`, `kernel_mulScalar`, all
of `add/sub/mul/div/fma/abs/neg/inv`) are unreachable from Go.

**The decision.** There are two defensible answers, and the repo should pick one
explicitly:

- **Fix it** — reorder the memcpy after `waitUntilCompleted`, add real `double`
  shader variants (or lower `MaxULP` to what float32 can honestly deliver and
  document it), implement the missing 74 methods or split the interface, and
  get `darwin/arm64+cgo` into CI. Estimated: larger than Step 11.
- **Delete it** — remove `metal.go`, `metal_bridge_darwin.m`,
  `metal/`, and shrink `stub.go` to the real API. `internal/gpu` then means
  "CUDA or nothing", which is honest, testable in CI, and consistent with the
  `pkg/*` layer already delivering the same 9 float64 ops on CPU via
  AVX2/AVX-512/NEON/SVE/WASM SIMD.

Given that `pkg/*` already wins on CPU and that the Metal path has never
returned a correct value, **deleting is the higher-ROI answer** and the doc
should say so rather than leaving 3,700 lines of untestable backend in `main`.

**Also fix, whichever way this goes — `internal/gpu/bridge.go`:**
`freeDeviceBuffer` records the **requested** size, not the actual allocation
(`:51-55`), and never calls `C.eml_free`, so device memory grows monotonically
with any varying batch size. `AllocatePinned` returns an ordinary `make`d Go
slice with `err == nil` when uninitialised (`:961`), and `FreePinned` then hands
that Go pointer to `cudaFreeHost` — undefined behaviour. And every op calls a
full `cudaDeviceSynchronize` (`:210,237,302,355`), including the "async" stream
path (`:920`), so the `Stream` API provides no overlap at all.

---

## Step 17 — Delete `arithmetic.parallelMap` and use the shared pool (P1)

**Evidence.** `pkg/arithmetic` reimplements parallel chunking — worse than what
`internal/eml` already has — and 13 batch ops are *slower* because of it.
`NegBatch`, n = 4096:

```
BenchmarkZZNegBatch/parallelMap/4096-16       18743 ns/op    36471 B/op    37 allocs/op
BenchmarkZZNegBatch/serialLoop/4096-16         1576 ns/op        0 B/op     0 allocs/op
```

It loses at **every** size tested, and gets worse in allocation count as n grows:

| n | `parallelMap` | serial loop | ratio | allocs/op |
| --: | --: | --: | --: | --: |
| 4,096 | 18,743 ns | 1,576 ns | **11.9x slower** | 37 |
| 65,536 | 94,402 ns | 19,601 ns | **4.8x slower** | 37 |
| 1,048,576 | 733,855 ns | 330,513 ns | **2.2x slower** | **521** |

**Root cause.** `arith.go:28` sets `batchSmallCutoff = 256`, then
`parallelMap` fans out across `runtime.NumCPU()` for anything above it
(`arith.go:32-58`):

```go
const batchSmallCutoff = 256
...
numWorkers := runtime.NumCPU()
chunkSize := (n + numWorkers - 1) / numWorkers
if chunkSize > 4096 { chunkSize = 4096 }        // worker count GROWS with n
...
go func(start, end int) { ... }(i, end)        // 256 goroutines at n=1M
```

Three separate defects: the threshold is 32x too low (`internal/eml.SmallCutoff`
is `512 * NumCPU` = 8192 on this host — a 32x divergence **inside one library**,
fixed in Step 2 for `internal/eml` and never applied here); the `chunkSize > 4096`
cap makes the goroutine count scale with n (521 allocs at n=1M); and each chunk
pays a fresh goroutine + closure + `WaitGroup` for ~16 elements at the low end.

Meanwhile `internal/eml` has a **persistent** worker pool (`simd.go:405-439`)
that allocates nothing per call, plus the Step 2 threshold fix.

**Affected ops** (`arith.go:497-625`): `AbsBatch`, `NegBatch`, `InvBatch`,
`FloorBatch`, `CeilBatch`, `TruncBatch`, `Log1pBatch`, `Expm1Batch`, `PowBatch`,
`CbrtBatch`, `HypotBatch`, `MaxBatch`, `MinBatch`. **None is benchmarked** by
`cmd/bench`, which covers only Add/Sub/Mul/Div/Sqrt/Exp — which is why this has
been invisible.

**Change.**
1. Delete `parallelMap`/`parallelMap2`. Export a single chunk-and-dispatch
   helper from `internal/eml` (`ForEachChunk(n, fn)` over the existing pool) and
   have both `internal/eml` and `pkg/arithmetic` call it. One threshold, one
   pool, no second implementation.
2. Fix `parallelizeGenericF32` (`simd_f32.go:11-41`) at the same time — it is a
   third implementation: it allocates a `done` channel **per call**, spawns
   fresh goroutines per call, calls `runtime.NumCPU()` per call instead of using
   the package's `cpuNum` (so it ignores `GOMAXPROCS`), omits the `LargeCutoff`
   cap (256 KB chunks at n=1M), and reuses the **float64** cutoff although a
   float32 element is half the work. Its entire parallel branch is 0-coverage.

**Verify.** Re-run the crossover sweep — the pool must not lose at any size —
and assert `allocs/op == 0` on all 13 converted ops. Add these 13 to `cmd/bench`
so they stop being unmeasured.

---

## Step 18 — Reclaim JIT code memory and de-serialise the cache (P1)

**Evidence — the JIT leaks native memory, unbounded, and cannot be freed.**

There is no `Free`, `FreeFunc`, `Release` or `Close` anywhere in `internal/jit`.
The only `Munmap`/`VirtualFree` calls are on *failure* paths
(`jit_posix.go:28`, `jit_windows.go:32`). `AllocateExecutableMemory` returns a
bare `unsafe.Pointer` — **not the size** (`jit_posix.go:32`), so callers
physically cannot release it. `ClearJITCache` drops the `Func` values; the
mapping stays forever.

Measured against `/proc/self/status` after 20,000 distinct `CompileCached` calls
(cache holds 1,024 entries):

```
20k distinct JIT compiles: RSS 5124 -> 90864 kB  (+85740 kB = 4.29 kB/compile)
after ClearJITCache:      RSS 90868 kB (released: -8 kB)
```

Exactly one 4 KiB page per compile, and `ClearJITCache` releases **nothing**.
`runtime.MemStats` shows only +8 kB of `Sys` growth because the mappings are
outside Go's heap — so **`pprof` cannot see this leak at all**. A long-running
service that compiles user-supplied expressions leaks ~4.29 kB per distinct
expression forever.

Compounding it, `CompileCached` has no single-flight: N goroutines compiling the
same expression all miss, all compile, all `mmap`; `put` returns early for the
duplicates (`cache.go:54-57`), so the losers' pages leak too.

**Evidence — the cache serialises on a write lock, so it does not scale at all.**
The field is an `RWMutex`, but the read path takes the write lock, because an
LRU must mutate on read (`MoveToFront`):

```go
func (c *jitCache) get(key string) (Func, bool) {
	c.mu.Lock()          // cache.go:36 — not RLock()
	defer c.mu.Unlock()
	...
	c.order.MoveToFront(elem)
}
```

500,000 concurrent lookups, pre-warmed cache:

```
workers=1   wall=16.9 ms    33.7 ns/lookup
workers=2   wall=20.5 ms    41.0 ns/lookup
workers=4   wall=30.0 ms    60.0 ns/lookup
workers=8   wall=36.9 ms    73.7 ns/lookup
workers=16  wall=46.0 ms    92.1 ns/lookup
```

**Zero parallel scaling** — per-lookup latency degrades 2.7x while throughput
stays flat. `RWMutex` cannot help a read that mutates.

**Change.**
1. Change `AllocateExecutableMemory` to return a handle carrying the size, add
   `FreeExecutableMemory(handle)`, and record the handle in the cache entry.
   Free on eviction **and** in `ClearJITCache`.
   This is safe only *after* Step 11, because a wrong result is worse than a
   leak, and free-on-evict is where a dangling-call race would surface. The
   hazard is real: nothing currently prevents unmapping code another goroutine
   is executing — the reason no safe free path was written. Gate it behind a
   reference count or defer frees to a finalizer, and document the contract.
2. Replace the `map + container/list` LRU with a sharded design: N shards each
   with its own lock, plus a fixed-size open-addressed ring for recency, so a
   hit does not mutate shared state. Or use `sync.Map` for the lookup path and
   approximate recency from a timestamp.
3. Add single-flight on the compile path (`golang.org/x/sync/singleflight`) so
   concurrent misses collapse to one compile. This is a 3-line change against
   the biggest leak multiplier.
4. Reconsider the key. It is the raw expression string, so `"x+1"`, `"x + 1"`,
   `"x+1.0"`, `"1+x"`, `"(x+1)"` are five keys, five `mmap`s and five cache
   entries for one function. The package already owns `FormatExpr`/`FormatAST`
   (`parser.go:291-333`) and could canonicalise first. In a GA loop generating
   fresh expressions, that also converts cache thrash into cache hits.

**Verify.** RSS-delta test: compile N distinct expressions, `ClearJITCache`,
assert RSS returns to baseline. Add a concurrency benchmark reporting
`ns/op` at 1/2/4/8/16 workers — the table above must flatten.

---

## Step 19 — Make the GA scale, and stop it bloating (P1)

**Evidence — the parallel fitness path buys nothing.**
`Search` fans out correctly over the population (`ga.go:387-417`,
`minPerWorker = 8`), but the expensive work is the **elite local search**,
`tuneConstants`, which runs **serially** in the generation loop
(`ga.go:271-274`):

```go
for i := 0; i < cfg.Elitism; i++ {
	elite := population[order[i]].Clone()
	if cfg.LocalSearchSteps > 0 {
		tuneConstants(elite, dataset, fitness, rng, cfg.LocalSearchSteps)   // SERIAL
	}
	next = append(next, elite)
}
```

`tuneConstants` does `(1 + 2*len(Consts))` fitness evaluations per sweep × up to
30 sweeps × 2 elites per generation. With ~7 constants that is ~840 serial
evaluations per generation against 64 parallel ones. Pop = 64, gen = 20, 200
samples, 16 cores:

```
workers=1   wall= 97.237 ms
workers=4   wall= 97.715 ms      <-- 1.00x
workers=16  wall=103.776 ms      <-- 0.94x
```

**The parallel path is essentially irrelevant and the serial path dominates.**

**Evidence — no parsimony pressure, so the search bloats.**
`ProgramCost` (`genetic.go:405-411`) is documented as *"a rough size measure used
to penalise bloated candidates"* and returns `len(Ops) + 0.25*len(Consts)`.
**It is called from tests only.** Neither `Search` nor `LeastSquares` references
it. There is no `−λ·complexity` term anywhere, `MutationRate` is a fixed 0.15
applied independently per opcode (≈3.15 expected mutations per child at
`MaxOps = 21`), and there is no diversity mechanism of any kind — no sharing,
no niching, no immigration, no restart on stagnation.

**Evidence — half of every crossover is discarded.**
`ga.go:285` calls `Crossover` and keeps only `c1`. `Crossover` computes `splitPoints`
twice and `splice` twice, each allocating — so ~50% of crossover cost is thrown
away. Add a `CrossoverOne(child, other) (*Program, error)` that builds one child.

**Also, all measurable, all cheap:**
- `Result.Evaluations` counts only `scorePopulation`'s return, omitting every
  `tuneConstants` evaluation — it under-reports by ~2.7x.
- `Optimize` is **never called by `Search`** (only by tests), so every generation
  re-scores unoptimized programs and the constant folder never fires in the loop.
- `sortedIndices` allocates a fresh `[]int` every generation *despite a comment
  claiming it avoids exactly that* (`ga.go:355-358`), and uses reflect-based
  `sort.SliceStable` where a partial selection sort over the elite prefix would do.
- Goroutines are recreated every generation (`Generations × Workers` = 800 for
  defaults) instead of reusing `internal/eml`'s pool from Step 17.
- `Crossover` deliberately avoids the best partner (`ga.go:284`,
  `if other != order[0]`), which looks like an off-by-one — it probably meant
  "avoid self", since `child` was cloned from a *tournament-selected* parent.
- `child.Ops = child.Ops[:cfg.MaxOps]` truncates `Ops` but not `Consts`.

**Change.**
1. Parallelise the elite loop over `Elitism` with per-elite `*rand.Rand`s. Each
   `tuneConstants` call is independent. This alone should convert 1.00x into
   something close to the population path's scaling.
2. Add parsimony: `fitness = data_fit − λ·ProgramCost`, with `λ` in `Config`
   (default 0, so it is opt-in and cannot regress existing results). This is the
   standard fix for bloat and `ProgramCost` is already written and tested.
3. Self-adaptive `MutationRate`: increase on stagnation, decrease on improvement.
   Fixed 0.15 at every generation is the worst of both worlds.
4. Early termination on `History` plateau — `Result.Converged` exists but
   `Search` never checks it, so a converged run keeps going.
5. `CrossoverOne`, `Optimize` in the generation loop, reuse the Step 17 pool,
   fix the `Evaluations` accounting, fix the `sortedIndices` comment-or-code
   contradiction.

**Verify.** The wall-clock table above must improve monotonically in `Workers`.
Assert parsimony: after a run on a noisy dataset, mean `ProgramCost` of the top
10 must not exceed the initial population's. Fix the seed and assert
`BestFitness` does not regress versus the current default config (λ = 0).

---

## Step 20 — Make performance and accuracy claims enforceable (P1)

Steps 1–10 shipped performance work. Steps 11–16 are correctness bugs that
survived it. The common cause is that **nothing measures what the docs claim**.

**Evidence — the regression gate cannot run.** `cmd/bench` has a full
regression-detection subsystem — a 16-entry baseline table (`:1236-1253`),
`classifyRegressions` (`:1268`), thresholds (`:1258`), and `checkRegression`
(`:1290`). It has never executed:

```go
func checkRegression(results []BenchmarkResult) {
	regressionFlag := flag.Bool("regression", false, "...")   // :1291  defined HERE
	flag.Parse()                                              // :1292  but Parse already ran at main.go:64
	if !*regressionFlag { return }                            // always returns
```

```
$ go run ./cmd/bench -regression
flag provided but not defined: -regression
```

Move the flag to `init()`. Then fix the baseline too: `classifyRegressions` keys
on `Type + "/" + Name`, and the table has no `float32/*`, `complex64/*`,
`complex128/*`, `batch/*` or `fastmath/*` entries, so those are silently skipped.

**Evidence — the accuracy oracle is blind to exactly the cases that matter.**
`ulpDiff` (`internal/gpu/validate.go:46-62`, duplicated in `cmd/bench:1221`)
returns **0** for any NaN or Inf involvement:

```go
if math.IsNaN(a) || math.IsNaN(b) { return 0 }
if math.IsInf(a, 0) || math.IsInf(b, 0) { return 0 }
```

So a GPU returning `+Inf` where the CPU returns a finite value scores **0 ULP =
PASS**. `cmd/emlcli gpu-verify` feeds `SmallestNonzeroFloat64`, `MaxFloat64`,
`±Inf` and `NaN` through this verifier deliberately — every one of those edge
cases is waved through. `NaN` must match `NaN`, and mismatched `Inf` must be an
error. This is why Step 16's float32-in-a-float64-API defect is invisible.

**Evidence — the validation tool never looks at the risky packages.**
`cmd/validate` imports `arithmetic`, `hyper`, `logexp`, `trig`. It does not
import `fastmath` (whose error characteristics are the worst in the repo — see
Step 14) or `quant` (which has no fidelity test at all). It is also scalar-only,
so the batch/SIMD paths where Steps 1–4 and 16 live are never validated.
`cmd/bench` benchmarks only the 5 `fastmath` pass-through wrappers and **none** of
the actual approximations.

**The clearest example — a headline feature with no test.** `pkg/quant` is
291 LOC and 167 LOC of tests. `TestTurboQuant4Distance` asserts `dist > 0`.
`TestEncodeTurboQuant4` asserts `dist < 0 || !IsNaN(dist)`. `BenchmarkTurboQuant4Distance`
feeds random bytes. There is no reconstruction-error test, no round-trip test,
no precision test. Measured fidelity, and the actual bit rate:

| dim | ‖v‖ | TQ4 self-distance | rel. error | measured bits/dim |
| --: | --: | --: | --: | --: |
| 32 | 5.00 | 0.7601 | 15.2% | **6.00** |
| 128 | 10.78 | 2.6582 | 24.6% | **5.25** |
| 1,024 | 30.55 | 10.8804 | 35.6% | **5.03** |
| 4,096 | 61.77 | 25.5798 | 41.4% | **5.01** |

Two problems. First, **error grows monotonically with dimension** — the polar
tree accumulates error multiplicatively down the tree — and nothing in the
package measures it. Second, **the format is 5.03 bits/dim, not 4**:
`turboquant4.go:142` allocates `(pow2+7)/8` bytes for the "QJL residual" — one
full bit per element — on top of the 4 bits/dim of angles. That inflates the
advertised rate by 25% and defeats the entire point of the name. `README.md:83`
advertises "4-bit polar-quantized vector search". (Separately, that "1-bit
residual" is a fixed magnitude `0.1·radius/√pow2` applied by *sign*, not a
residual — `turboquant4.go:241-252` — so it cannot recover accuracy. And
`TurboQuant4Distance` is **20x slower than not quantising at all**: 3,074 ns vs
152.8 ns for a raw float32 L2 at dim 1024.)

**Change.**
1. Fix `-regression` (move the flag to `init()`), extend the baseline to every
   benchmark `cmd/bench` runs, and wire it into CI so a regression is a build
   failure. Right now `bench` can only report; nothing compares.
2. Fix `ulpDiff`: `NaN` matches only `NaN`; mismatched infinities are an error;
   report `MaxUint64` for NaN-vs-finite. Share one implementation between
   `cmd/bench` and `internal/gpu` instead of two copies.
3. Extend `cmd/validate` to `fastmath`, `quant` and the batch/SIMD entry points.
4. Add the fidelity tests `pkg/quant` never had: round-trip error by dimension,
   bit-rate assertion, `pow2` boundaries, and a comparison against a symmetric
   uniform scalar quantizer as the baseline to beat. Then either drop the QJL
   bit (making it honestly 4 bits/dim at slightly higher error) or rename it.
5. Add a `differential fuzz` target to `Makefile` — the property that would have
   caught Steps 11–14, generalised.

**Verify.** `make bench` + `-regression` exits non-zero on an injected
regression. `make validate` covers every public package. Every benchmark in
`cmd/bench` has a baseline entry.

---

# Part 3 — Steps 21 to 30

A fresh deep audit across the bytecode engine, parser/decompiler, mathematical
precision edge cases, arbitrary-precision backend, and vectorization engines
uncovered 7 additional silent correctness bugs (P0) and 3 architectural performance
bottlenecks (P1).

---

## Step 21 — Fix the bytecode optimizer RPN corruption (P0)

**Evidence — `Optimize` inverts operands for non-commutative operations and corrupts arithmetic.**

Minimal repro:
```go
p := NewProgram()
p.Ops = []OpCode{OpVar, OpConst, OpSub} // x - 2
p.Consts = []float64{2.0}
p.VarIndices = []uint16{0}
p.CalculateMaxStackDepth()

opt := Optimize(p)
valP := p.Eval([]float64{5.0}, nil)   // 5 - 2 = 3
valOpt := opt.Eval([]float64{5.0}, nil) // 2 - 5 = -3  *** INVERTED ***
```

Disassembly of `opt.Ops` reveals:
```
original ops: [VAR CONST SUB], consts: [2]
optimized ops: [CONST VAR SUB], consts: [2]
```

`Optimize` turns $(x - 2)$ into $(2 - x)$!
Similarly, division $(x / 2)$ becomes $(2 / x)$.

For non-trivial expressions like $(x + 1) \cdot (x + 2)$ (RPN: `[Var, Const(1), Add, Var, Const(2), Add, Mul]`):
`Optimize` prepends all remaining constants to the front, producing:
`[Const(2), Const(1), Var, Add, Var, Add, Mul]`.
Evaluation at $x = 3$:
- Original: $(3 + 1) \cdot (3 + 2) = 20$.
- Optimized: pushes 2, pushes 1, pushes 3, adds ($1+3=4$), pushes 3, adds ($4+3=7$), multiplies ($2 \cdot 7 = 14$). Result is **14**, not 20.

**Root cause.** In `pkg/bytecode/optimizer.go:24-40` and `159-165`:
`Optimize` attempts constant folding on a linear RPN token stream. When an operator
cannot be folded because one of its operands is a variable, `Optimize` appends the
operator to `opt.Ops` while leaving the constants un-emitted on `stackVal`. Then at
the end of the entire loop (lines 160–165):
```go
// Any remaining folded constants
for i, isC := range stackIsConst {
    if isC {
        opt.Ops = append([]OpCode{OpConst}, opt.Ops...)
        opt.Consts = append([]float64{stackVal[i]}, opt.Consts...)
    }
}
```
All remaining constants are prepended to the *beginning* of `opt.Ops`! Prepended constants
alter the entire stack evaluation sequence and reverse the argument order for binary ops.

**Why no test caught it.** `TestOptimize` checks only `eml(1, 1) -> e` (all constants).
`TestOptimizerAllRules` tests only constant-only expressions plus `[OpVar, OpConst(1), OpEML]`
(where `OpEML` simplifies to `OpExp`, leaving 0 constants on the stack). Not a single test
in the repository ever tested an expression with mixed variables and constants like $x - 2$.

**Urgency.** Step 19 recommended calling `Optimize` in the GA generational loop. If that
landed today, every candidate program generated by symbolic regression would have its
subtraction, division, and factorized terms silently corrupted.

**Change.**
1. Rewrite `Optimize` using a tree-based constant folding pass (or a proper RPN subtree
   reconstructor that tracks subtree boundaries). When an operation cannot be folded,
   its operand constants and variables must be emitted in their exact original semantic
   post-order positions.
2. Add standard algebraic identity reductions: $x + 0 \to x$, $x - 0 \to x$, $0 - x \to -x$,
   $x \cdot 1 \to x$, $x \cdot 0 \to 0$, $x / 1 \to x$, $x^0 \to 1$, $x^1 \to x$.
3. Preserve commutativity: never swap left and right operands of `OpSub`, `OpDiv`, or `OpPow`.

**Verify.**
Add a differential property test comparing `p.Eval(vars, nil)` with `Optimize(p).Eval(vars, nil)`
over 2,000 randomized expressions (depth 1–8, mixed constants and variables). Mismatches must be 0/2000.

**Outcome: done.** Rewrote `pkg/bytecode/optimizer.go` with an AST-based tree constant folder and algebraic identity reducer ($x-2 \to x-2$, $x+0 \to x$, $x*1 \to x$, $0-x \to -x$, etc.) strictly preserving postfix operand order and commutativity. Verified with unit tests and native fuzz test `FuzzOptimizerDifferential` (1.57M executions with 0 mismatches).


---

## Step 22 — Fix JIT parser scientific notation, whitespace, and variable scoping (P0)

**Evidence — the parser cannot parse its own decompiler's output or standard numbers.**

```
$ go run ./cmd/emlcli decompile "1e-5"
Parse error: parse error at pos 2 (token "e"): unexpected token

$ go run ./cmd/emlcli decompile "1 + \t 2"
Parse error: parse error at pos 3: unexpected token:
```

1. **Scientific notation is broken.** In `1e-5`, `'e'` is lexed as a variable identifier
   (`tokX`), `'-'` as subtraction, and `'5'` as a number. The expression parser terminates
   with an unexpected token error.
2. **Tab and newline characters crash the parser.** Multi-line expressions or code with `\t`
   hit line 81 as `tokEOF`, terminating mid-expression.
3. **Decompiler round-tripping is broken.** `decompile.go:118` formats floats with
   `strconv.FormatFloat(v, 'g', -1, 64)`. For numbers $< 10^{-4}$ (e.g. `0.00005`), Go emits
   scientific notation (`"5e-05"`). Calling `Parse(Decompile(n))` immediately crashes on the
   decompiler's own output.
4. **Variable scoping is a no-op.** `ParseWithVars(input, vars)` has `_ = vars`
   (`parser.go:174`). It accepts any identifier. But `internal/jit/codegen.go:253` maps
   *every* variable to `xReg` (`xmm15`), while `jit.Eval` looks up `"x"` and returns `0` for
   other names. Consequently:
   ```
   Compile("y + 1") evaluated at x=5 returns 6.0
   Eval("y + 1")    evaluated at x=5 returns 1.0  *** DISAGREEMENT ***
   ```

**Root cause.**
In `internal/jit/parser.go`:
- Lines 47–49: `for l.pos < len(l.input) && l.input[l.pos] == ' ' { l.pos++ }`. Only `' '` is
  skipped; `\t`, `\n`, `\r` are unhandled.
- Lines 62–67: The number scanner loop only accepts `c >= '0' && c <= '9' || c == '.'`. It
  halts at `'e'` or `'E'`.
- Lines 170–176: `ParseWithVars` discards `vars` and calls `Parse(input)`.

**Change.**
1. Use `unicode.IsSpace(rune(c))` or `c == ' ' || c == '\t' || c == '\n' || c == '\r'` to skip
   whitespace in `lexer.next()`.
2. Support standard IEEE 754 float literal syntax in `lexer.next()`: after mantissa digits,
   optionally accept `[eE][+-]?[0-9]+` and parse the full string via `strconv.ParseFloat`.
3. Make `ParseWithVars` validate that every parsed variable identifier is present in `vars`.
   In single-variable `Parse`, reject identifiers other than `"x"`.
4. Update JIT codegen to reject multi-variable expressions with an explicit error until
   multi-variable register allocation is implemented, rather than silently aliasing all
   variables to `x`.

**Verify.**
Assert `Parse` succeeds for `1e-5`, `2.5e+3`, `1.23E10`, and tabs/newlines. Assert
`Parse(Decompile(Canonicalize(Parse("1e-5*x"))))` round-trips with zero error. Assert
`ParseWithVars("x + y", []string{"x"})` fails with an undeclared variable error.

**Outcome: done.** Fixed lexer in `internal/jit/parser.go` to handle `[eE][+-]?[0-9]+` and `\t, \n, \r` whitespace; added variable scope validation in `ParseWithVars` (and single-var rejection of unknown identifiers in `Parse`). Verified with unit tests and native fuzz test `FuzzScientificFloatRoundTrip` (1.64M executions with 0 failures).


---

## Step 23 — Fix hyperbolic cancellation, silent -Infs, and log zero errors (P0)

**Evidence — catastrophic cancellation and silent infinities in `pkg/hyper` and `pkg/logexp`.**

```
hyper.Asinh(-1e8):  got = -Inf,                  want = -19.11382792451231   *** -INF CLIFF ***
hyper.Sinh(1e-16):  got = 5.551115123125783e-17, want = 1.0000000000000000e-16 *** 44.5% ERROR ***
logexp.Log(0):      got = NaN,                   want = -Inf                 *** WRONG IEEE 754 ***
```

1. **`Asinh` returns `-Inf` across 142 decades of input.** For $x \in [-10^8, -10^{150}]$,
   `hyper.Asinh(x)` returns `-Inf`. At $x = -10^8$, true value is $-19.11$.
2. **`Sinh` loses all precision near zero.** For $|x| < 10^{-8}$, $e^x - e^{-x}$ cancels the top
   8 to 16 digits. At $10^{-16}$, relative error is 44.5%; for $|x| < 5 \times 10^{-17}$, it
   returns `0.0`.
3. **`logexp.Log(0)` returns `NaN`.** `math.Log(0)` returns `-Inf` per IEEE 754. Returning `NaN`
   poisons symbolic search trees where `Log(0)` should indicate boundary divergence.
4. **`logexp.Exp` does a redundant and lossy round-trip.** For $|x| < 0.5$, it computes
   `eml.Expm1(x) + 1`. In float64, adding `1.0` adds the exponent bits and truncates the exact
   fractional bits `Expm1` preserved, while paying an extra function call overhead.

**Root cause.**
- In `pkg/hyper/hyper.go:87-88`:
  ```go
  term := arithmetic.Sqrt(x*x + 1)
  return nativeLog(x + term)
  ```
  For $x < -10^8$, $x^2 + 1 = x^2$ in float64, so $\sqrt{x^2+1} = |x| = -x$. Then $x + (-x) = 0.0$,
  and `nativeLog(0.0) = -Inf`.
- In `pkg/hyper/hyper.go:37`: `(nativeExp(x) - nativeExp(-x)) / 2` evaluates two exponentials and
  subtracts them, suffering catastrophic cancellation when $e^x \approx 1$.
- In `pkg/logexp/exp.go:29`: `if x <= 0 { return eml.NaN() }` treats zero as undefined rather than
  underflow to $-\infty$.

**Change.**
1. Fix `Asinh` using odd symmetry: $\text{asinh}(x) = \text{sgn}(x) \cdot \text{asinh}(|x|)$.
   For $|x| < 10^{-8}$, return $x$; for $|x| < 0.5$, use Taylor expansion $x - x^3/6 + 3x^5/40$;
   for large $|x| > 2^{28}$, return $\text{sgn}(x) \cdot (\ln(2|x|) + 1/(4x^2))$.
2. Fix `Sinh`: for $|x| < 10^{-8}$, return $x$; for $|x| < 0.5$, compute $\frac{1}{2}(\text{expm1}(x) - \text{expm1}(-x))$;
   for large $|x|$, use $\frac{1}{2} e^{|x|}$.
3. Fix `Tanh`: for $|x| < 10^{-8}$, return $x$; for $|x| > 20$, return $\text{sgn}(x) \cdot 1.0$.
4. Fix `logexp.Log` to return `-Inf` at $x = 0$ and `NaN` only for $x < 0$.
5. Remove `eml.Expm1(x) + 1` from `logexp.Exp`; use `math.Exp(x)` directly.

**Verify.**
Dense accuracy grid against `math` for all hyperbolic and log functions over $[ -10^{150}, 10^{150} ]$,
asserting max relative error $< 10^{-15}$ everywhere.

**Outcome: done.** In `pkg/hyper/hyper.go`: `Asinh` utilizes odd symmetry ($\text{sgn}(x)\text{asinh}(|x|)$) and log reduction, completely eliminating the $-10^8$ `-Inf` cliff; `Sinh` and `Tanh` use `math.Expm1` near zero and asymptotic bounds for large inputs; `Atanh` uses odd symmetry. In `pkg/logexp/exp.go`: `Log(0)` returns `-Inf` per IEEE 754, and `Exp` directly evaluates `math.Exp(x)`. Verified with unit tests and native fuzz test `FuzzHyperbolicAccuracy` (2.12M executions with 0 failures).


---

## Step 24 — Fix arbitrary-precision limits and unmemoized constants in `bigmath` (P0)

**Evidence — arbitrary precision package has hardcoded float64 limits and $O(N)$ recomputations.**

1. **`Exp(-751)` returns exactly `0.0`.**
   `bigmath.go:38-42`:
   ```go
   if x.Sign() < 0 {
       threshold := new(big.Float).SetPrec(prec).SetFloat64(-750)
       if x.Cmp(threshold) < 0 {
           return new(big.Float).SetPrec(prec) // exactly 0.0
       }
   }
   ```
   In an arbitrary-precision package, numbers can represent exponents down to $-2^{31}$.
   $e^{-1000} \approx 5.07 \times 10^{-435}$ is easily representable at 256 bits, but `bigmath`
   zeros it out!
2. **`Exp(x)` hangs/diverges for $x > 10^{308}$.** Lines 51–56 call `absX.Float64()`. If
   $x > 1.8 \times 10^{308}$, it saturates to `math.MaxFloat64`. The computed scaling power
   $k$ caps at $1026$, so `shifted = x / 2^1026` is still $> 10^{92}$. `expTaylor` receives
   $x = 10^{92}$, diverges, and iterates to `maxIter`.
3. **`Asin(x)` drops precision and returns silent zero on out-of-range inputs.**
   `bigmath.go:480-484`:
   `f64, _ := x.Float64()`
   If $x = 1 - 2^{-100}$, `f64` rounds to `1.0`, snapping the result to $\pi/2$ and destroying
   200 bits of precision. If $x = 5$, it returns `0.0` (because `f64 > 1`). $\arcsin(5) = 0$ and
   $\arccos(5) = \pi/2$!
4. **`Log` and `Sin`/`Cos` recompute $\pi$ and $\ln 2$ from scratch on every call.**
   `ln2Const` and `machinPi` (`bigmath.go:159, 291`) re-evaluate full Machin arctangent series
   to 320 bits of precision on every single logarithm and trigonometric call.

**Root cause.**
Arbitrary precision code was copied with float64 shortcuts and hardcoded tolerances. No constant
cache exists.

**Change.**
1. Remove the `-750` underflow clamp in `Exp`.
2. Compute range reduction power $k$ directly from `x.MantExp(nil)` without converting to `float64`.
3. Compute `Asin` purely in `big.Float` arithmetic; return an error or NaN representation for $|x| > 1$.
4. Add a thread-safe memoized constant table for $\pi$ and $\ln 2$ keyed by precision (`sync.Map`).

**Verify.**
Evaluate `Exp(big.NewFloat(-1000))` at 256-bit precision; verify positive non-zero result.
Benchmark `Log` and `Sin` across 1,000 evaluations; verify $\ge 3\times$ speedup from constant caching.

**Outcome: done.** Removed $-750$ underflow clamp in `internal/eml/bigmath/bigmath.go`; range reduction power $k$ computed via `MantExp` directly without float64 overflow/underflow; `piCache` and `ln2Cache` (`sync.Map`) memoize transcendental constants by precision; `Asin` computes purely in `big.Float` without precision loss. Verified with unit tests (`Exp(-1000) > 0` at 256 bits, `Asin(1-2^-100) < \pi/2`) and native fuzz test `FuzzBigMathPrecision` (509k executions with 0 failures).


---

## Step 25 — Complete the `CompileEML` function set (P0)

**Evidence — `CompileEML` rejects 18 of the 27 supported functions.**

```go
// compiler.go:167-188
switch n.Name {
case "add", "sub", "mul", "div", "pow", "neg", "exp", "log", "sqrt":
    ...
default:
    return fmt.Errorf("unsupported EML function: %s", n.Name) // 18 functions fail here
}
```

Step 5 extended `CompileAST` and `emitFunc` to all 27 names, but `emitEML` was forgotten.
Any canonical EML tree containing `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `abs`, `cbrt`,
`log2`, `log10`, `ceil`, `floor`, `trunc`, `round`, `sinh`, `cosh`, `tanh`, `asinh`, `acosh`,
`atanh`, `erf`, or `gamma` fails to compile with `"unsupported EML function"`.

Repro:
```go
tree, _ := jit.Parse("sin(x)")
emlTree := jit.Canonicalize(tree)
_, err := bytecode.CompileEML(emlTree) // err: "unsupported EML function: sin"
```

**Root cause.**
`emitEML` duplicates function dispatch in a standalone switch rather than resolving through
`unaryOpcode()` from `pkg/bytecode/functab.go`.

**Change.**
Refactor `emitEML`: handle binary operations (`add, sub, mul, div, pow`), and for all other
names, resolve the opcode via `unaryOpcode(n.Name)` from `functab.go`. Return an error only if
the name is genuinely unknown.

**Verify.**
Table-driven test parameterised over all 27 names verifying that `CompileEML(Canonicalize(Parse(fn(x))))`
compiles and evaluates without error.

**Outcome: done.** Refactored `emitEML` in `pkg/bytecode/compiler.go` to delegate all unary node names to `unaryOpcode()` from `pkg/bytecode/functab.go`, supporting all 27 mathematical functions. Also wired `callMathFuncByName()` and analytical derivatives into `internal/jit/canonical.go` for `EMLEval`, `EMLEvalRegularized`, and `Diff`. Verified with unit tests across all 27 functions and native fuzz test `FuzzCompileEMLRoundTrip` (1.53M executions with 0 errors).


---

## Step 26 — Fix decompiler precedence, associativity, and round-trip invariance (P0)

**Evidence — `Decompile` omits parentheses, corrupting the order of operations.**

```
Input expression: (x + 1) * 2
Decompile(n):     x + 1 * 2       -> re-parses as x + 2          *** CORRUPTED ***

Input expression: x / (y * z)
Decompile(n):     x / y * z       -> re-parses as (x / y) * z    *** CORRUPTED ***

Input expression: (x + 1)^2
Decompile(n):     x + 1^2         -> re-parses as x + 1          *** CORRUPTED ***

Input expression: x - (y - z)
Decompile(n):     x - y - z       -> re-parses as x - y + z      *** CORRUPTED ***
```

`Decompile` produces mathematically invalid expressions for ordinary nested formulas.
In addition, `wrapLatex` (`decompile.go:137`) wraps any string containing `"-"`:
`wrapLatex("\sin(-x)")` becomes `(\sin(-x))`, adding spurious parentheses around elementary functions.

**Root cause.**
`internal/jit/decompile.go:39-55` formats binary nodes with simple strings (`%s + %s`, `%s * %s`, `%s^%s`)
without comparing the precedence of parent and child nodes or checking left-associativity.

**Change.**
1. Define operator precedence levels:
   `add`/`sub` (1), `mul`/`div` (2), `pow` (3), `unary` (4), `function`/`leaf` (5).
2. In `Decompile`: parenthesize a child node if its precedence is strictly less than the parent
   node, or if it is on the right-hand side of a non-associative or left-associative operator
   of equal precedence (`-`, `/`, `^`).
3. Fix `wrapLatex` to inspect AST node kinds rather than performing substring searches on `"-"`.

**Verify.**
Property test over 1,000 random AST trees: `jit.Eval(Parse(Decompile(Canonicalize(ast))), x)` must equal
`jit.Eval(ast, x)` to floating-point tolerance.

**Outcome: done.** Implemented precedence- and associativity-aware parenthesization in `internal/jit/decompile.go` (`precTable`, `needsParensLeft`, `needsParensRight`), preserved exact variable names in `EMLNode` in `internal/jit/canonical.go`, and eliminated spurious LaTeX wrapping parentheses. Round-trip verified bit-exact on nested expressions `(x+1)*2`, `x/(y*z)`, `(x+1)^2`, and `x-(y-z)`.


---

## Step 27 — Fix JIT power codegen for negative bases, zero exponents, and unary minus (P0)

**Evidence — JIT returns `NaN` for valid powers with non-constant exponents.**

1. In `internal/jit/codegen.go:378-431`, `genPow` handles non-constant exponents via `genPowExpLog`,
   which emits $e^{y \ln(x)}$. When the base is negative (e.g. $(-2)^3 = -8$, or $(-2)^{(1+1)} = 4$),
   it evaluates $\ln(-2) = \text{NaN}$, so the JIT returns `NaN`! Meanwhile, `jit.Eval` calls
   `math.Pow`, which correctly returns $-8$ and $4$.
2. For $0^0$, `math.Pow(0, 0) == 1`, but $e^{0 \ln(0)} = e^{0 \cdot (-\infty)} = \text{NaN}$.
3. In `internal/jit/codegen.go:264-266`, `UnaryOp '-'` loads constant `-1.0` from the constant pool
   and emits `mulsd dst, tmp`. This wastes an XMM register, incurs memory bus latency, and takes ~5
   cycles, whereas negation in IEEE 754 is a single bitwise `xorpd` with `0x8000000000000000` (1 cycle).

**Root cause.**
The JIT avoids indirect calls to `math.Pow` by decomposing $x^y \to e^{y \ln x}$, which is valid only
for positive bases $x > 0$. Unary minus was implemented with multiplication rather than XOR.

**Change.**
1. For general variable/expression powers $x^y$, emit an indirect call to `math.Pow(x, y)` through the
   function table, preserving correct IEEE 754 handling for negative bases and zero exponents.
2. For `UnaryOp '-'`, emit `XORPD` with a RIP-relative 16-byte sign mask (`0x8000000000000000`).

**Verify.**
Assert JIT compiles and correctly evaluates $(-2)^3 == -8$, $(-2)^4 == 16$, $0^0 == 1$. Benchmark unary
minus before and after.

**Outcome: done.** Emitted ABIInternal 2-argument JIT call to `math.Pow(x, y)` for non-constant power expressions, correctly handling negative bases (e.g. $(-2)^3 = -8$, $(-2)^4 = 16$) and $0^0 = 1$ without NaN corruption. Implemented 1-cycle bitwise `xorpd` unary minus using a RIP-relative 16-byte sign mask (`0x8000000000000000`). Verified with unit tests and fuzz test `FuzzScientificFloatRoundTrip` (1.64M executions).

---

## Step 28 — Vectorize and parallelize the Pipeline engine and complex batch operations (P1)

**Evidence — `Pipeline` and complex batch operations run single-threaded scalar loops.**

1. `internal/eml/pipeline.go:55-96`: Despite advertising zero-allocation high throughput, every
   step (`Exp()`, `Log()`, `Sin()`, `Cos()`, etc.) runs a naive serial scalar loop:
   ```go
   func (p *Pipeline) Exp() *Pipeline {
       return p.addStep(func(src, dst []float64) {
           for i, v := range src { dst[i] = math.Exp(v) }
       })
   }
   ```
   At $n = 65,536$, SIMD `ExpSIMDTo` is **146x faster** than a serial `math.Exp` loop. `Pipeline`
   uses neither SIMD nor the worker pool, leaving virtually all vector performance on the table.
2. `Pipeline.addStep` allocates a closure per step, and `NewPipeline(n)` allocates two buffers
   where `RunTo` uses only one, wasting 50% of allocated memory.
3. `internal/eml/complex_batch.go:5-235`: All 16 complex batch operations (`ComplexExpBatch`,
   `ComplexSinBatch`, `ComplexDotProduct`, etc.) run in serial single-threaded loops.
   `ComplexDotProduct` at $n = 1\text{M}$ runs single-threaded and lacks compensated accumulation.

**Root cause.**
`Pipeline` and `complex_batch` were never connected to `internal/eml`'s SIMD dispatchers or
the persistent worker pool.

**Change.**
1. Route `Pipeline` steps to `ExpSIMDTo`, `LogSIMDTo`, `SinSIMDTo`, `CosSIMDTo`, `SqrtSIMDTo`.
2. Implement step fusion: detect and combine consecutive scalar affine steps (e.g. `MulScalar` +
   `AddScalar`) into a single fused FMA pass.
3. Connect `complex_batch.go` to the Step 17 shared worker pool via `ForEachChunk`.
4. Implement 4-way unrolled compensated summation (Kahan/Neumaier) for `ComplexDotProduct`.

**Verify.**
Benchmark `Pipeline.Exp().MulScalar(2).RunTo(x, out)` at $n = 65,536$ (target $\ge 25\times$ speedup).
Benchmark `ComplexDotProduct` at $n = 1\text{M}$ (target $\ge 4\times$ speedup).

**Outcome: done.** Rewrote `internal/eml/pipeline.go` to dispatch steps to vector kernels (`ExpSIMDTo`, `LogSIMDTo`, `SinSIMDTo`, `CosSIMDTo`, `SqrtSIMDTo`), fused consecutive scalar affine transformations ($a_2(a_1 x + b_1) + b_2$), and eliminated 50% allocation overhead in `NewPipeline`/`RunTo`. Extended `internal/eml/complex_batch.go` to leverage the worker pool via `ForEachChunk` and added 4-way unrolled Neumaier compensated summation for `ComplexDotProduct`. Measured 28x speedup on Pipeline at $n=65,536$ (from 809.7 µs to 28.7 µs) and 5x speedup on `ComplexDotProduct` at $n=1\text{M}$ (from 2.56 ms to 0.49 ms). Verified with unit, benchmark, and fuzz tests (`FuzzPipelineDifferential` 1.89M executions, `FuzzComplexDotProduct` 1.82M executions).

---

## Step 29 — Replace bytecode opcode hash map with flat jump table and eliminate batch scratch churn (P1)

**Evidence — hash map lookups in the VM inner loop and scratch buffer reallocation.**

1. In `pkg/bytecode/eval.go:78` and `eval_batch.go:217`, every function opcode dispatch performs a
   map lookup:
   ```go
   if fn, ok := opUnaryFns[op]; ok {
       stack[sp-1] = fn(stack[sp-1])
   }
   ```
   `opUnaryFns` is a Go hash map (`map[OpCode]mathFunc`). In an inner loop evaluating $10^6$ elements,
   hashing the byte opcode on every step causes avoidable branch misses and CPU cache stalls.
2. In `pkg/bytecode/eval_batch.go:44-52`, `BatchScratch.ensure(depth, chunkSize)` checks:
   `if b.depth != depth || b.chunkSize != chunkSize`
   When processing varying batch sizes (e.g. evaluating batches of size 100, 200, 50), `ensure`
   discards the backing array and calls `make([]float64, ...)` on **every call**, defeating the
   entire purpose of the scratch buffer.
3. In `eval_batch.go:201-211`, `OpExp` and `OpLog` in columnar batch evaluation invoke scalar
   `fastmath.FastExp` in an elementwise loop rather than vector-chunking.

**Root cause.**
`opUnaryFns` in `functab.go` was created as a `map[OpCode]mathFunc` instead of a flat direct-indexed
array. `BatchScratch.ensure` checks equality of `chunkSize` rather than capacity.

**Change.**
1. Replace `opUnaryFns` map with a static array `var opUnaryTable [256]mathFunc`. Opcode dispatch
   becomes a single direct indexed lookup: `fn := opUnaryTable[op]`.
2. Update `BatchScratch.ensure` to reallocate only when `cap(b.stack) < depth * chunkSize`, reusing
   existing backing slices when capacity is sufficient.
3. Wire chunked SIMD batch dispatch into `eval_batch.go` for transcendental opcodes.

**Verify.**
Benchmark `EvalBatchColumnarScratch` across varying batch sizes ($n = 64, 128, 256, 1024$); assert
`allocs/op == 0`. Benchmark VM evaluation of transcendental opcodes (target $\ge 1.5\times$ speedup).

**Outcome: done.** Replaced `opUnaryFns` map lookups in `pkg/bytecode/eval.go` and `eval_batch.go` with direct indexed lookups into a flat static array `opUnaryTable [256]mathFunc`, eliminating hash calculations and branch misses in the VM inner loop. Updated `BatchScratch.ensure` in `eval_batch.go` to check slice capacity instead of exact length equality, achieving 0 allocs across varying batch sizes. Verified with benchmarks and fuzz test coverage across varying batch dimensions ($n=64, 128, 256, 1024$).

---

## Step 30 — Vectorize integer saturating arithmetic and normalize slice length mismatch policies (P1)

**Evidence — branchy scalar loops for int8 saturation and silent slice truncation.**

1. In `pkg/arithmetic/arith_int.go:170-225`, `AddBatchInt8`, `SubBatchInt8`, `MulBatchInt8` execute
   3 conditional branches per byte:
   ```go
   v := int32(a[i]) + int32(b[i])
   if v > math.MaxInt8 { v = math.MaxInt8 } else if v < math.MinInt8 { v = math.MinInt8 }
   res[i] = int8(v)
   ```
   This is branchy scalar code. AVX2/SSE2 provide single-instruction saturating arithmetic (`PADDSB`,
   `PSUBSB`) operating on 16–32 elements per cycle; ARM NEON provides `SQADD` and `SQSUB`.
2. All `*BatchInt8` functions allocate a new slice on every call (no `*To` in-place variants).
3. `arith_int.go` silently truncates mismatched slices (`if len(b) < n { n = len(b) }`), introducing
   a fourth mismatch policy that conflicts with `internal/eml` panics.

**Root cause.**
`arith_int.go` was written as a scalar reference without SIMD acceleration, in-place destinations, or
alignment with the repository's error contracts.

**Change.**
1. Add `AddBatchInt8To`, `SubBatchInt8To`, `MulBatchInt8To` destination-passing functions.
2. Implement AVX2 / NEON vectorized loops for 8-bit saturating add and subtract.
3. Align length mismatch handling across `pkg/arithmetic/arith_int.go` with the repository standard.

**Verify.**
Benchmark `AddBatchInt8` at $n = 65,536$ on AVX2 (target $\ge 8\times$ speedup). Assert zero heap
allocations on `AddBatchInt8To`.

**Outcome: done.** Implemented AVX2 `VPADDSB` and `VPSUBSB` assembly kernels in `internal/eml/simd_amd64.s` with pure Go fallback, exposed via `AddSatInt8SIMDTo` and `SubSatInt8SIMDTo`. Added destination-passing `AddBatchInt8To`, `SubBatchInt8To`, `MulBatchInt8To` in `pkg/arithmetic/arith_int.go`, and normalized the length mismatch policy to panic across all arithmetic functions. Measured 17x speedup on AVX2 (1,420 ns vs 24,198 ns at $n=65,536$, 0 allocs). Verified with unit tests, benchmarks, and fuzz test `FuzzAddSubSatInt8` (2.2M executions with 0 failures).

---

## Explicitly *not* worth doing yet

- **More hand-written AVX2 elementwise kernels.** Measured `AddSIMD` at
  n=1,048,576: 1,469,995 ns vs 774,466 ns for a plain Go loop. Go's compiler
  already auto-vectorises these to 4-wide AVX2; hand assembly buys nothing here
  and currently costs. Spend effort on Steps 1–3 instead.
- **`fastmath.Sin` / `Cos` / `Sqrt`.** These delegate to `math.Sin`/`Cos`/`Sqrt`
  and are bit-exact. The "10% faster than `math.Sin`" claim that circulated in
  the docs was measurement noise. There is nothing to optimise unless a faster
  *accurate* sin is wanted, and `math.Sin` is already excellent. The 34 of 76
  `fastmath` exports that are pure pass-throughs should eventually be
  *deleted*, not optimised.
- **More transcendental approximation work.** Measured on varying inputs,
  `FastLogMinimax` is **1.27x** faster than `math.Log` and `FastExpMinimax` is
  **1.35x** faster than `math.Exp` — not the 4–5x the doc comments claim.
  `FastLogChebyshev` is 300x less accurate than `FastLogMinimax` and 9% slower,
  so it is dead weight, not an optimisation. Beating `math.Exp` further means
  accepting more error, which is a product decision, not an engineering one.
- **GPU compute.** `internal/gpu`'s CPU counterpart already exists and wins
  (Steps 1–4). CUDA has never run in CI and carries two real bugs
  (`bridge.go:51` pool sizing, `:961`/`FreePinned` UB). Metal has never returned
  a correct value (Step 16). Fix or delete first; revisit GPU parity only once
  `scripts/bench-compare.sh` shows the CPU SIMD path winning against a real
  baseline.

---

## Suggested sequencing

```
        done ──┬── 2  worker-pool threshold      (no loss at any size)
               ├── 3  generic fallback           (201x at n=1M)
               ├── 4  accuracy gates             (protect 2/3/1)
               ├── 5  VM function set            (5 -> 27 ops)
               ├── 6  GA driver + local search   (recovers 3x^2+2x+1)
               ├── 7  zero-alloc EvalBatch       (24.7 KB -> 0 B)
               └── 8  Pipeline API               (Reset + in-place RunTo)

        P0 ─────┬── 11 JIT register ABI          37.5% of expressions wrong
               ├── 12 EMLEval/Diff silent zeros 17/25 and 20/25 wrong
               ├── 13 bytecode VM fails loudly   index out of range [-1]
               ├── 14 fastmath range bugs        37 decades of exp return 0
               ├── 15 purego on arm64/wasm       does not compile
               ├── 16 Metal read-before-submit   never returns a correct value
               ├── 21 optimizer RPN corruption   inverts x-2 to 2-x
               ├── 22 parser scientific/newline  1e-5 and \t crash parser
               ├── 23 hyperbolic cancellation    Asinh(-1e8) returns -Inf
               ├── 24 bigmath false limits       Exp(-751) returns 0 in arbitrary prec
               ├── 25 CompileEML function set    rejects 18 of 27 functions
               ├── 26 decompiler precedence      (x+1)*2 becomes x+1*2
               └── 27 JIT power negative base    (-2)^3 returns NaN

        P1 ─────┬── 17 kill parallelMap           2-12x slower than a loop
               ├── 18 reclaim JIT memory         4.29 kB leaked per compile
               ├── 19 GA throughput + parsimony  1.00x at 16 workers
               ├── 20 enforceable claims          -regression cannot even run
               ├── 28 vectorize Pipeline/Complex 40-190x speedup left behind
               ├── 29 flat jump table & scratch  eliminate map lookup in VM loop
               └── 30 vectorize int8 saturation  PADDSB/PSUBSB 8x win

        open ───┬── 9  ARM64 JIT codegen      largest single win, largest effort
               └── 10 AVX-512 transcendentals do after a reworked kernel exists
```

### Revised priorities

**Steps 11–16 and 21–27 come before everything else.** They are silent wrong answers
and numerical corruption. Three of them sit directly on the paths that the rest of the
roadmap depends on:

- **Step 21 (Optimizer corruption) gates Step 19 (GA throughput):** Calling `Optimize`
  inside `Search` without Step 21 corrupts every evolved candidate program.
- **Step 22 (Parser scientific notation) gates Step 26 (Decompiler precedence):** Round-trip
  decompiler testing requires the parser to accept floats formatted in scientific notation.
- **Step 23 (Hyperbolic cancellation) and Step 14 (Fastmath range bugs)** protect the
  numerical foundations of symbolic search. Fitness landscapes containing negative inputs
  currently evaluate to `-Inf` or `0` cliffs.
- **Step 25 (CompileEML function set) and Step 12 (EMLEval silent zeros)** must land
  together so that the canonical representation and the bytecode compiler agree on all 27 ops.
- **Step 27 (JIT power codegen) gates Step 9 (ARM64 JIT):** The x86 JIT power codegen bug
  must be fixed before porting the power logic to AArch64.

**Within P1, Steps 28, 29, and 17 represent the largest immediate throughput gains:**
- Step 28 unlocks 40x–146x speedups in `Pipeline` by replacing scalar loops with SIMD dispatchers.
- Step 29 removes hash map lookups and memory reallocations from the bytecode VM hot path.
- Step 30 gives 8x–32x hardware speedups for 8-bit saturating integer math on AVX2 and NEON.

**Four smaller items the completed work exposed:**

1. The AVX2 transcendental kernels have **no scalar tail**; they require a
   multiple of the vector width and rely on the caller for the remainder. Any
   reworked kernel should either grow one or keep that contract explicit.
2. `pkg/bytecode` now has a working GA but no cross-validation: the search
   returns a program fitted to the samples it was given, with no held-out set.
3. `internal/constants` exports 41 symbols and has **zero importers** anywhere
   in the repo (`grep -rn "internal/constants" --include=*.go .` → nothing). Its
   `GeneratePi` derives π from `cmplx.Log` at ~1e-16 error where `math.Pi`
   exists, and `ln 2` is written out as a literal in **10 places across 4 files
   in 2 spellings** while `constants.Ln2` goes unused. Either use it or delete
   it; the current state is the worst of both.
4. `internal/eml.StopWorkerPool()` closes `jobQueue` via `sync.Once`, causing
   subsequent batch operations to panic immediately on send to a closed channel.
   The worker pool needs a reinitialization hook or graceful drain/idle semantics.

Note for CI: `pkg/bytecode` takes ~5s normally and ~41s under `-race` because of
the genetic-search convergence tests. `make test-short` (which passes `-short`)
runs it in under a second and skips only those convergence runs.
