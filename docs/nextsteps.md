# Next Steps: Improvement Plan & Known Issues

This document catalogs all known bugs, architectural flaws, and planned improvements for the `emlgo` math library.

---

## P0 Blockers — Critical Correctness

These must be fixed before any release.

### P0-1. `AbsBranchless` corrupts negative numbers
* **File:** `internal/eml/fused.go:110`
* **Bug:** The mask `sign<<63 | sign` flips both the sign bit AND the LSB when `sign=1`. For `x = -5.0` (0xC014000000000000), this produces 0x4014000000000001 ≈ 5.0000000000000009 instead of exactly 5.0.
* **Why:** `sign` is `bits >> 63`, which is 0 or 1. `1<<63 | 1` = `0x8000000000000001`. XOR with that flips sign bit and LSB.
* **Fix:** Change to `bits ^ (sign << 63)` — only flip the sign bit.
* **Test gap:** Tests use `1e-10` tolerance which masks the LSB corruption.

### P0-2. `ComplexCos` uses wrong formula
* **File:** `internal/constants/constants.go:115-119`
* **Bug:** Computes `(Exp(z) + Exp(iz)) / 2`. Correct formula: `cos(z) = (Exp(iz) + Exp(-iz)) / 2`.
* **Why:** Variable `conj` is actually `iz` (not conjugate). Code passes `z` instead of `-iz` to the first `ComplexOne`.
* **Fix:** `iz = complex(-imag(z), real(z))`, `niz = complex(imag(z), -real(z))`, return `(Exp(iz) + Exp(niz)) / 2`.

### P0-3. Missing length validation on binary SIMD ops
* **File:** `internal/eml/simd.go:235-262`
* **Bug:** `AddSIMD`, `SubSIMD`, `MulSIMD`, `DivSIMD` never check `len(a) == len(b)`. If mismatched, SIMD dispatch reads `b[i]` out of bounds → runtime panic.
* **Why:** Every other public SIMD function (`SIMD`, `ExpSIMDTo`, `AbsSIMDTo`, etc.) validates. These four were added later without the check.
* **Fix:** Add `if len(a) != len(b) { panic("slice length mismatch") }` at the top of each.

### P0-4. `movsdStore` wrong REX prefix for high XMM registers
* **File:** `internal/jit/codegen.go:170`
* **Bug:** Uses REX.B (0x41) instead of REX.R (0x44) for high XMM registers. `MOVSD [mem], xmm_reg` puts `xmm_reg` in the ModR/M reg field — needs REX.R to extend the reg field. REX.B extends rm/SIB base (RSP), which doesn't need extension.
* **Why:** Copy-paste error from `movsdLoad` which uses REX.B correctly for the XMM load.
* **Fix:** Change `e.emit(0x41)` to `e.emit(0x44)`. Update `coverage_test.go:57` to expect 0x44.

### P0-5. Complex number first-class support and real-domain fast path
* **Current State:** All math operations treat arguments as `float64`. Complex numbers are only supported in `internal/eml/complex.go` and `internal/constants/constants.go`. The JIT, SIMD batch ops, and public API have no complex path. Generating trigonometric identities, negative roots, and constants like `i` or `π` requires manual complex assembly.
* **Plan:**
  - Design a dual-path API: fast `float64` path for real numbers, `complex128` path via `math/cmplx` for complex expressions.
  - Extend JIT to emit `complex128` operations (or at minimum, support `cmplx.Exp`, `cmplx.Log`, `cmplx.Sin`, `cmplx.Cos` as JIT-callable functions).
  - Add batch complex operations: `AddComplexSIMD`, `MulComplexSIMD`, `ExpComplexSIMD`, etc.
  - Ensure all trigonometric identities (`sin²z + cos²z = 1`, `e^{iπ} = -1`) hold in tests for complex inputs.

### P0-6. Log1p/Expm1 asymmetric sensitivity and compound form stability
* **Current State:** `Log1p` and `Expm1` exist but are not used in compound evaluations. `Exp(y * Log(x))` for `x` near 1 or `y` near 0 suffers catastrophic cancellation. `Log(1 + x)` for small `x` loses all significant digits.
* **Plan:**
  - Audit all places where `Log(1+x)` or `Exp(x)-1` appears and replace with `Log1p(x)` / `Expm1(x)`.
  - Add optimized stable evaluation paths for compound derived forms (e.g., `Pow(x,y)` near `x=1`, `LogRatio(a,b)` for `a≈b`).
  - Use `math.Expm1` and `math.Log1p` internally wherever subtraction/addition with 1.0 occurs.
  - Add numerical stability tests: verify ULP accuracy for inputs near identity points (`x=1` for Log, `x=0` for Exp, etc.).

### P0-7. Optional multiprecision backend for symbolic verification and deep tree reduction
* **Current State:** All computation is `float64`. No way to verify symbolic identities to arbitrary precision. No deep tree reduction / algebraic simplification.
* **Plan:**
  - Add an optional `big` backend using `math/big.Float` for arbitrary-precision evaluation.
  - Evaluate MPFR bindings (`github.com/cznic/mathutil` or `github.com/mattetti/mp`) for high-performance multiprecision.
  - Implement tree reduction: simplify EML expression trees algebraically (constant folding, identity elimination, strength reduction).
  - Provide a `VerifyIdentity(expr1, expr2 string, precision uint) bool` function that evaluates both expressions at random high-precision points and compares ULP distance.

### P0-8. Core AST as explicit interface with zero-allocation design
* **Current State:** JIT uses `Node` interface with 5 concrete types. Each `BinaryOp`, `UnaryOp`, `FunctionCall` heap-allocates. The tree-walking `Eval` function allocates nothing but the JIT path doesn't benefit. No arena/pool for AST nodes.
* **Plan:**
  - Design AST nodes as fixed-size structs with a `nodeKind` discriminator byte instead of interface dispatch.
  - Implement arena allocator: `type Arena struct { buf []byte; off int }` — bump-allocate all nodes in a contiguous buffer.
  - Add `ParseToArena(expr string, arena *Arena) *Node` for zero-GC-parse paths.
  - Benchmark allocation count via `testing.AllocsPerRun` for parse and eval paths.
  - Target: 0 allocations for `Parse` + `Eval` of any expression.

### P0-9. Canonical constructors for minimal EML tree equivalents
* **Current State:** No way to map a standard expression back to its canonical EML tree form. No simplification rules. `x + y` stays as `BinaryOp(+, x, y)` instead of potentially being reduced via EML identities.
* **Plan:**
  - Implement canonical form rules:
    - `exp(x) = eml(x, 1)`
    - `log(x) = eml(1, eml(eml(1, x), 1))`
    - `x + y = eml(log(x), exp(y))` (for the full EML derivation)
    - `sin(x)` as a composition of EML nodes
  - Add a `Canonicalize(n Node) Node` function that rewrites any AST to minimal EML tree form.
  - Add equivalence checker: two expressions are equal if their canonical forms are structurally identical.
  - Provide `EMLSize(n Node) int` to measure the complexity of the EML decomposition.

---

## P1 Bugs — High Severity

### P1-1. `MinBranchless`/`MaxBranchless` broken for NaN in second argument
* **File:** `internal/eml/fused.go:114,123`
* **Bug:** When `b` is NaN, `diff = a - NaN = NaN`, sign bit of NaN is 0, so `mask = 0`, result = `b = NaN`. IEEE 754 `minNum` returns the non-NaN operand.
* **Fix:** Add NaN-aware logic: `if isNaN(a) { return b }; if isNaN(b) { return a }` before the bitwise operation.

### P1-2. `IntAbs(MinInt)` overflows
* **File:** `pkg/arithmetic/arith.go:415`
* **Bug:** `-MinInt64` overflows back to `MinInt64`. The adjacent `GCD` function handles this but `IntAbs` doesn't.
* **Fix:** Add overflow guard: `if a == math.MinInt { return math.MaxInt }` or use unsigned arithmetic.

### P1-3. `GCD(MinInt64, x)` silently corrupts result
* **File:** `pkg/arithmetic/arith.go:356-361`
* **Bug:** Clamping `MinInt64` to `MaxInt64` changes the mathematical result. `GCD(MinInt64, 2) = 2` but returns `GCD(MaxInt64, 2) = 1`.
* **Fix:** Use unsigned arithmetic for the absolute value, or document the limitation clearly.

### P1-4. `hasNeonDot` unconditionally true on all ARM64
* **File:** `internal/eml/simd_arm64.go:214`
* **Bug:** Dot product requires ARMv8.2-A. Cortex-A53/A55 (ARMv8.0-A) will SIGILL when executing dot product instructions.
* **Fix:** Runtime detection via `getauxval(AT_HWCAP)` checking `HWCAP_ATOMICS` or `HWCAP2_I8MM`.

### P1-5. `Pow` integer conversion overflow for large exponents
* **File:** `pkg/arithmetic/arith.go:170`
* **Bug:** `int(y)` silently truncates for `y > 1e19`. The subsequent `intY%2 == 0` check is meaningless.
* **Fix:** Range-check before casting: `if y > float64(math.MaxInt) || y < float64(math.MinInt) { ... }`.

### P1-6. `push`/`popTo` emit malformed x86 instructions
* **File:** `internal/jit/codegen.go:184-192`
* **Bug:** `movsdStore` with mod=00 and SIB base=RSP requires a disp32 per x86-64 ISA. No displacement is emitted, so the CPU reads the next code bytes as the displacement field.
* **Fix:** Use `modrm(1, ...)` with a disp8 of 0, or add the missing disp32. (Currently dead code, but the encoding is invalid.)

---

## P2 Issues — Medium Severity

### P2-1. Inconsistent NaN semantics across packages
* `arithmetic.Min`/`Max` return the non-NaN operand.
* Go 1.21+ `math.Min`/`math.Max` return NaN if either operand is NaN.
* `trig.Cot` returns NaN when sin(x)==0, even at multiples of π (should be ±Inf).
* **Decision needed:** Document which convention the library follows and enforce it consistently.

### P2-2. Mixed panic vs error on length mismatch
* Public SIMD functions (`AddSIMD`, etc.): `panic` on length mismatch.
* `Batch()`: returns `error` on length mismatch.
* `AddBatch()`, `NegBatch()`, etc. in `arithmetic`: no length check at all.
* **Fix:** Standardize: panic for internal/performance-critical paths, error for user-facing API. Document the rationale.

### P2-3. `Acsch(0)` ignores sign
* **File:** `pkg/trig/trig.go:236`
* **Bug:** Returns `+Inf` for both `x→0+` and `x→0-`. Should be `inf(sign(x))`.
* **Fix:** `if x == 0 { return inf(math.Copysign(1, x)) }` or handle via `1/x` approach.

### P2-4. `Sec`/`Csc` missing NaN/Inf guards
* **File:** `pkg/trig/trig.go:64-69`
* **Bug:** Every other trig function checks NaN/Inf, but these two don't. Inconsistent API contract.
* **Fix:** Add `if isNaN(x) || isInf(x, 0) { return nan() }` before the computation.

### P2-5. `dispatchSinCosSIMDTo` doubles range-reduction work
* **File:** `internal/eml/simd_dispatch_amd64.go:339-351`
* **Bug:** Calls `sinAVX2` and `cosAVX2` as independent passes, doubling the Cody-Waite range reduction work.
* **Fix:** Implement a fused `sincosAVX2` kernel that shares the range reduction and evaluates sin/cos polynomials together.

### P2-6. Worker pool goroutines never shut down
* **File:** `internal/eml/simd.go:396-405`
* **Bug:** No `Stop()` or context cancellation. Goroutines leak for the process lifetime.
* **Fix:** Add `Stop()` function that closes the `jobQueue` channel. Workers exit when channel is closed.

---

## Improvement Plan

### 1. ARM64 NEON Assembly Kernels
* **Status:** Pending
* **Current State:** `simd_arm64.s` is empty. NEON-labeled functions are plain Go loops.
* **Plan:** Write hand-tuned ARM64 NEON assembly for 2-wide `float64` ops, bitwise abs/neg, scalar broadcast.

### 2. ARM64 SVE/SVE2 Assembly Kernels
* **Status:** Pending
* **Current State:** SVE implemented via Go-level VL-aware loops. No actual SVE instructions.
* **Plan:** Write VLA SVE assembly with predicate registers and scalable vector registers.

### 3. ARM64 Transcendental Batch Kernels
* **Status:** Pending
* **Current State:** ARM64 transcendental batch ops fall through to `parallelizeGeneric`.
* **Plan:** NEON-vectorized Cody-Waite range reduction and minimax polynomial evaluation.

### 4. Fix CI/CD Pipeline
* **Status:** Pending
* **Current State:** `go.mod` declares `go 1.26.1` (future version). CI tests Go 1.21-1.23. Duplicate linter configs.
* **Plan:** Update `go.mod`, CI matrix, consolidate linter config.

### 5. Go Documentation Coverage
* **Status:** Partially done
* **Done:** CHANGELOG.md created. doc.go files added for all packages.
* **Remaining:** Add godoc to individual exported functions in `internal/eml/`.

### 6. JIT Non-Integer Power
* **Status:** ✅ Done

### 7. GPU Kernel Library Expansion
* **Status:** ✅ Done

### 8. WASM SIMD Optimization
* **Status:** ✅ Done

### 9. Benchmark Suite & Regression Detection
* **Status:** Pending
* **Plan:** `benchstat` integration, ULP tracking, CI artifact collection.

### 10. API Stability, SemVer & CHANGELOG
* **Status:** ✅ Done
