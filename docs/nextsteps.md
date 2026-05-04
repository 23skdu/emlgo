# Next Steps: Performance Optimization Roadmap

This document outlines the roadmap for enhancing performance and stability of the `emlgo` library.

## ACHIEVED MILESTONES (Hardening Phase)

1. `[x]` **97%+ Test Coverage:** Reached 97.2% total test coverage across all architectural paths.
2. `[x]` **SIMD Dispatcher Architecture:** Refactored core arithmetic and transcendental functions to use platform-specific dispatchers, ensuring clean build tags and architecture-agnostic orchestration.
3. `[x]` **Numerical Accuracy Parity:** Verified all core mathematical functions against the standard Go `math` library to ensure 100% correctness.
4. `[x]` **FMA/SIMD Performance:** Maintained performance parity while simplifying scalar fallbacks and improving testability.
5. `[x]` **Redundant Code Elimination:** Cleaned up unreachable branches and simplified monolithic kernel implementations.

---

## FUTURE CONSIDERATIONS (Beyond v2.1)

### 1. GPU/CUDA Production Readiness

- **Core Implementation:**
  - [ ] Port all major math operators (Exp, Log, Sin, Cos, Add, Mul) to CUDA kernels.
  - [ ] Implement zero-copy pinned memory support for host-to-device transfers to minimize latency.
  - [ ] Add asynchronous execution streams to overlap compute and data transfer.
- **Tooling & Validation:**
  - [ ] **Benchtool:** Add `--device gpu` flag to track GPU vs CPU performance ratios.
  - [ ] **CLI:** Add `eml gpu-status` to verify hardware availability, driver version, and compute capability.
  - [ ] **Validation:** Implement ULP-based verification specifically for GPU results against Go's `math` library.
- **Testing:**
  - [ ] **Unit Tests:** Verify kernel launch parameters and grid/block size calculations.
  - [ ] **Fuzz Tests:** Fuzz GPU kernels with extreme edge cases (NaN, Inf, subnormal numbers).

### 2. ARM SVE/SVE2 Support

- **Core Implementation:**
  - [ ] Implement runtime detection for SVE vector length (VL).
  - [ ] Create Vector-Length Agnostic (VLA) assembly kernels for all SIMD-accelerated paths.
  - [ ] Leverage SVE2 specific instructions for complex number arithmetic and cryptography-adjacent ops.
- **Tooling & Validation:**
  - [ ] **Benchtool:** Report active SVE vector length and performance scaling metrics.
  - [ ] **Validation:** Verify numerical stability across different VL configurations.
- **Testing:**
  - [ ] **Fuzz Tests:** Validate SVE kernel results against scalar implementations using property-based testing.

### 3. WebAssembly SIMD Intrinsics

- **Core Implementation:**
  - [ ] Implement `wasm_simd128` kernels for all batch operations.
  - [ ] Optimize memory alignment for WASM linear memory access.
- **Tooling & Validation:**
  - [ ] **Test Harness:** Set up a `node` or `d8` based environment with SIMD enabled for CI.
  - [ ] **Benchtool:** Create a web-based benchmark harness for browser-side performance verification.
- **Testing:**
  - [ ] **Unit Tests:** Cross-platform consistency checks between WASM and Native results.

### 4. JIT Polynomial Compilation

- **Core Implementation:**
  - [ ] Implement a lightweight expression parser for mathematical polynomials.
  - [ ] Build a JIT engine using `golang.org/x/sys/unix` for managing executable memory pages.
  - [ ] Generate optimized machine code for specific polynomial expressions at runtime.
- **Tooling & Validation:**
  - [ ] **Benchtool:** Compare JIT-compiled polynomial performance vs. pre-compiled SIMD kernels.
  - [ ] **Validation:** Ensure JIT-compiled results match standard Go implementations within 1 ULP.
- **Testing:**
  - [ ] **Fuzz Tests:** Fuzz the expression parser with malicious or malformed polynomial strings.
