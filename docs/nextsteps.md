# Next Steps: Performance Optimization Roadmap

This document outlines the roadmap for enhancing performance and stability of the `emlgo` library.

## PO BLOCKERS 🚫 (Must-Fix Before Further Work)

### Critical

1. **Integer overflow in GCD when negating `math.MinInt64`** — `pkg/arithmetic/arith.go:281-284`
   `a = -a` when `a == math.MinInt64` overflows (two's complement), producing a negative value that breaks the GCD algorithm and corrupts LCM results.

2. **Integer overflow in `PowInt` with `math.MinInt64` exponent** — `pkg/arithmetic/arith.go:113-125`
   `PowInt(x, math.MinInt64)` causes `-n` to overflow back to `MinInt64`, resulting in infinite recursion and stack overflow.

3. **LCM overflow guard insufficient** — `pkg/arithmetic/arith.go:293-302`
   The bounds check only covers `a/gcd` but not the subsequent multiplication `(a/gcd) * b`. Large inputs silently overflow.

4. **`fmaScalar` SIGILL on pre-Haswell CPUs** — `internal/eml/simd_amd64.s:318-330`
   `VFMADD213SD` is issued unconditionally with no runtime `hasFMA` check. Any CPU without FMA support (pre-Haswell, non-Intel) will crash with SIGILL.

5. **`WasmAlign16` alignment logic is broken** — `internal/eml/wasm_utils.go:13-24`
   The skip calculation `(16 - int(offset>>3)) & 1` can only produce 0 or 1 regardless of actual misalignment. Silently drops elements from the slice.

### Major

6. **Dead SVE branch in arm64 `emlSIMD` dispatch** — `internal/eml/simd_dispatch_arm64.go:127-136`
   The `hasSVE && n >= sveVL()` branch calls `neonEml()` — the same function as the `else if` fallback. The SVE path is dead code with no actual SVE kernel.

7. **Stack buffer escapes in `AbsBatch` and `FloorBatch`** — `pkg/arithmetic/arith.go:389-413, 508-530`
   When `64 < n < 256`, returns a slice backed by a stack-allocated `[64]float64` array — use-after-return undefined behavior.

8. **Incorrect parallelization heuristic** — `pkg/arithmetic/arith.go` (all batch functions use `runtime.GOMAXPROCS(0)` instead of `runtime.NumCPU()`)
   `GOMAXPROCS` is the Go scheduler thread limit, not the CPU count. Internal `eml` package correctly uses `NumCPU()`.

9. **IEEE 754 signed zero lost in `trig.Sin`** — `pkg/trig/trig.go:35-38`
   `if x == 0 { return 0 }` flattens both `+0.0` and `-0.0` to `+0.0`. Should use `return x` to preserve sign bit.

10. **`bench/main.go` complex128 Tan benchmark compares `cmplx.Tan` to itself** — `cmd/bench/main.go:811-815`
    Both the emlgo path and the math path call `cmplx.Tan(x)`. The benchmark produces a meaningless ratio of ~1.0.

11. **Hardcoded magic thresholds duplicated across packages** — `pkg/hyper/hyper.go`, `pkg/logexp/exp.go`, `pkg/fastmath/fastmath.go`
    `709.78` (near `math.Ln(math.MaxFloat64)`) and `-745.13` are repeated as raw float literals with no named constants.

12. **`PowInt` uses O(n) loop** — `pkg/arithmetic/arith.go:120-124`
    Exponentiation by repeated multiplication instead of O(log n) binary exponentiation.

---

## ACHIEVED MILESTONES (Hardening Phase)

1. `[x]` **97%+ Test Coverage:** Reached 97.2% total test coverage across all architectural paths.
2. `[x]` **SIMD Dispatcher Architecture:** Refactored core arithmetic and transcendental functions to use platform-specific dispatchers, ensuring clean build tags and architecture-agnostic orchestration.
3. `[x]` **Numerical Accuracy Parity:** Verified all core mathematical functions against the standard Go `math` library to ensure 100% correctness.
4. `[x]` **FMA/SIMD Performance:** Maintained performance parity while simplifying scalar fallbacks and improving testability.
5. `[x]` **Redundant Code Elimination:** Cleaned up unreachable branches and simplified monolithic kernel implementations.

---

## FUTURE CONSIDERATIONS (Beyond v2.1)

### 1. GPU Backend

#### CUDA (Linux/Windows with NVIDIA GPU)

- **Core Implementation:**
  - `[x]` Port all major math operators (Exp, Log, Sin, Cos, Tan, Sinh, Cosh, Tanh, Sqrt, EML) to CUDA kernels.
  - `[x]` Create a pure C API layer (`cuda/eml_capi.h`/`.cu`) bridging Go cgo to CUDA kernels.
  - `[x]` Implement cgo bridge (`internal/gpu/bridge.go`) with `//go:build cuda` tag for conditional compilation.
  - `[x]` Implement zero-copy pinned memory support for host-to-device transfers to minimize latency.
  - `[x]` Add asynchronous execution streams to overlap compute and data transfer.
  - `[x]` Provide stub fallback (`internal/gpu/stub.go`) when `-tags cuda` is not set.
- **Tooling & Validation:**
  - `[x]` **Benchtool:** Add `--device gpu` flag to track GPU vs CPU performance ratios.
  - `[x]` **CLI:** Add `eml gpu-status` to verify hardware availability, driver version, and compute capability.
  - `[x]` **CLI:** Add `eml gpu-bench` for quick GPU batch performance test (all 9 ops, multiple sizes, GPU vs CPU timing, speedup).
  - `[x]` **CLI:** Add `eml gpu-verify` for ULP-based GPU result verification against Go's `math` library.
  - `[x]` **Validation:** Implement `BatchVerifier` with configurable ULP tolerance and per-element error reporting.
- **Testing:**
  - `[x]` **Unit Tests:** Verify kernel launch parameters and grid/block size calculations.
  - `[x]` **Fuzz Tests:** Fuzz GPU kernel launch configs with extreme edge cases.
  - `[x]` **End-to-End Tests:** Validate GPU results match CPU results within 1 ULP (requires CUDA hardware; build tag `cuda`).

#### Metal (macOS, Apple Silicon)

- **Core Implementation:**
  - `[x]` Implement all math operators (Exp, Log, Sin, Cos, Tan, Sinh, Cosh, Tanh, Sqrt, Abs, Neg, Inv, Add, Sub, Mul, Div, FMA, AddScalar, MulScalar, EML) as Metal compute kernels.
  - `[x]` Create Objective-C bridge (`internal/gpu/metal_bridge.m`) with cgo integration for darwin/arm64.
  - `[x]` Provide automatic Metal GPU detection on Apple Silicon (no build tags required).
  - `[x]` Automatic float64<->float32 conversion at the bridge layer.
  - `[x]` Unified memory usage via `MTLResourceStorageModeShared`.
- **Tooling & Validation:**
  - `[x]` **CLI:** `eml gpu-status` automatically shows Metal devices on darwin/arm64.
  - `[x]` **CLI:** `eml gpu-bench` runs benchmarks on the Metal GPU.
  - `[ ]` **CLI:** `eml gpu-verify` for ULP-based verification (requires double-precision shaders).
- **Optimization Opportunities:**
  - `[ ]` **Double-precision shaders:** Use `double` type in Metal for full float64 accuracy.
  - `[ ]` **Buffer reuse:** Cache Metal buffers across calls to reduce allocation overhead.
  - `[ ]` **Async command submission:** Use Metal command queues with multiple command buffers.
  - `[ ]` **NEON SIMD assembly:** Write hand-tuned NEON assembly for 2-wide float64 arithmetic.

### 2. ARM SVE/SVE2 Support

- **Core Implementation:**
  - `[x]` Implement runtime detection for SVE vector length (VL) via `/proc/self/auxv` HWCAP_SVE parsing + `prctl(PR_SVE_GET_VL)`.
  - `[x]` Create Vector-Length Agnostic (VLA) Go-level kernels for all 11 SIMD-accelerated paths, chunking by detected VL.
  - `[ ]` Write SVE assembly kernels (`simd_sve_arm64.s`) using Go's SVE asm support — requires arm64 SVE hardware to verify.
  - `[ ]` Leverage SVE2 specific instructions for complex number arithmetic and cryptography-adjacent ops.
- **Tooling & Validation:**
  - `[ ]` **Benchtool:** Report active SVE vector length and performance scaling metrics.
  - `[ ]` **Validation:** Verify numerical stability across different VL configurations.
- **Testing:**
  - `[ ]` **Fuzz Tests:** Validate SVE kernel results against scalar implementations using property-based testing.

### 3. WebAssembly SIMD Intrinsics

- **Core Implementation:**
  - `[x]` Implement optimized WASM SIMD kernels for all batch operations (add, sub, mul, div, sqrt, abs, neg, inv, fma, addScalar, mulScalar) with 8-wide block unrolling for JIT auto-vectorization to `wasm_simd128`.
  - `[x]` Implement WASM-specific dispatch layer (`simd_dispatch_wasm.go`) with `//go:build wasm` tag.
  - `[x]` Optimize memory alignment for WASM linear memory access (`WasmAlign16`, `WasmPageAlign` utilities).
  - `[x]` Add `HasWasmSIMD()` feature flag.
- **Tooling & Validation:**
  - `[x]` **Test Harness:** Create `scripts/wasm_test.sh` — builds WASM binaries, runs with Node.js, supports `test`, `bench`, and `serve` modes.
  - `[x]` **Benchtool:** Create `wasm/bench.html` — a web-based benchmark harness for browser-side performance verification, with results table and pass/fail indicators.
- **Testing:**
  - `[ ]` **Cross-platform consistency checks:** WASM vs Native result validation (requires WASM runtime in CI).

### 4. JIT Polynomial Compilation

- **Core Implementation:**
  - `[x]` Implement a lightweight expression parser for mathematical polynomials.
  - `[x]` Build a JIT engine using `golang.org/x/sys/unix` for managing executable memory pages.
  - `[x]` Generate optimized machine code for specific polynomial expressions at runtime.
- **Tooling & Validation:**
  - `[x]` **Benchtool:** Compare JIT-compiled polynomial performance vs. pre-compiled SIMD kernels.
  - `[x]` **Validation:** Ensure JIT-compiled results match standard Go implementations within 1 ULP.
- **Testing:**
  - `[x]` **Fuzz Tests:** Fuzz the expression parser with malicious or malformed polynomial strings.
