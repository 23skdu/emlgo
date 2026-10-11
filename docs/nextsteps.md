# Next Steps: Prioritized Engineering Roadmap & Open Work Log

An evidence-based plan for the remaining open engineering items in `emlgo`.
Every step below cites measurements, architectural audits, and hardware evidence from this repository.

---

## Active Roadmap Overview

All historical foundational work (Steps 1–8, 11–30 from v0.4/v0.5), residual defect fixes (Steps 31–50), and **all AMD64 AVX-512 & AVX-512 VNNI vector acceleration work** have been **verified and completed**. Their technical resolutions are archived in Git history, release notes, and regression suites.

The remaining open work is consolidated, prioritized, and tracked below:

| # | Priority | Step | Status | Target Platforms | Headline Objective & Evidence |
| :-- | :--- | :--- | :--- | :--- | :--- |
| **1** | **P0** | Native ARM64 JIT Codegen | **open** | ARM64 / AArch64 | Port polynomial JIT machine code compiler to AArch64 (`v0–v31`); brings 17x–26x speedup to Apple Silicon & Graviton |
| **2** | **P1** | Float32 ARM64 NEON Assembly Acceleration | **open** | ARM64 / AArch64 | 4-wide NEON assembly kernels (`FADD.4S`, `FMUL.4S`, `FSUB.4S`, `FDIV.4S`, `FSQRT.4S`) replacing scalar loops in `simd_f32.go` |
| **3** | **P2** | GPU CUDA Asynchronous Stream Overlap & Pinned Memory | **open** | Linux / CUDA | Pinned host memory (`cudaHostAlloc`) and stream double-buffering to lower GPU crossover to $N \le 4,096$ |

---

## Detailed Step Specifications

---

### Step 1 — Native ARM64 JIT Codegen (P0)

**Status:** Open

**Evidence & Context:**
- The JIT polynomial compiler is the library's largest performance differentiator on x86-64, compiling mathematical AST expressions directly into native machine code:
  - Linear ($P_1$): **19.0x faster** than AST evaluation
  - Quadratic ($P_2$): **26.5x faster** than AST evaluation
  - Cubic ($P_3$): **17.0x faster** than AST evaluation
  - Average JIT compilation overhead: ~45–54 µs.
- `internal/jit/codegen_amd64.go` is currently amd64-only. On all other architectures, `codegen_stub.go` returns `"JIT codegen requires amd64"`, causing ARM64 environments (Apple Silicon M-series, AWS Graviton, Ampere Altra) to fall back to the AST interpreter.
- Executable memory management (`AllocateExecutableMemory`) already supports `linux/arm64` and `darwin/arm64` using `unix.Mmap` and `unix.Mprotect`. Only the instruction encoder is missing.

**Change:**
1. Implement `internal/jit/codegen_arm64.go`:
   - Map expression AST variables and temporaries to AArch64 SIMD/floating-point registers `v0–v31` (scalar float64 lanes `d0–d31`).
   - Implement machine code binary encoder for standard AArch64 arithmetic instructions:
     - `FADD` (`0x1e602800`): $D_d = D_n + D_m$
     - `FSUB` (`0x1e603800`): $D_d = D_n - D_m$
     - `FMUL` (`0x1e600800`): $D_d = D_n \times D_m$
     - `FDIV` (`0x1e601800`): $D_d = D_n / D_m$
     - `FMADD` (`0x1f000000`): $D_d = D_a + (D_n \times D_m)$ (Fused Multiply-Add)
     - `FSQRT` (`0x1e61c000`): $D_d = \sqrt{D_n}$
     - `FNEG` (`0x1e614000`): $D_d = -D_n$
     - `FABS` (`0x1e60c000`): $D_d = |D_n|$
     - `RET` (`0xd65f03c0`): Subroutine return
   - Adhere to Go ABIInternal conventions: pass floating-point arguments in `F0` (`d0`), return in `F0`, and preserve callee-saved registers `d8–d15`.
2. Connect `Compile()` and `CompileCached()` to the AArch64 codegen backend.

**Verify:**
- Run differential fuzzing against AST evaluation (`FuzzJITDifferential`) on native ARM64 runners and QEMU AArch64 emulation (`qemu-aarch64`), asserting 0 mismatches across 100,000 generated expressions.
- Assert memory cleanup via `FreeExecutableMemory` without leaks.
- Benchmark speedup on Apple Silicon / Graviton in `docs/performance.md`.

---

### Step 2 — Float32 ARM64 NEON Assembly Acceleration (P1)

**Status:** Open

**Evidence & Context:**
- Float32 precision is critical for machine learning workloads, neural network layers, and computer vision pipelines.
- AMD64 Float32 AVX-512 hardware acceleration has been implemented and certified (`addAVX512F32`, `mulAVX512F32`, `subAVX512F32`, `divAVX512F32`, `sqrtAVX512F32` processing 16 elements per instruction).
- ARM64 NEON processes **4 float32 values per instruction** in a 128-bit vector register (`FADD.4S`, `FMUL.4S`, `FSUB.4S`, `FDIV.4S`, `FSQRT.4S`).
- `internal/eml/simd_f32.go` currently falls back to unrolled Go scalar loops on ARM64 platforms.

**Change:**
1. Add NEON float32 assembly kernels to `internal/eml/simd_arm64.s` using `FADD`, `FSUB`, `FMUL`, `FDIV`, `FSQRT` over `Vn.4S`.
2. Connect platform dispatch in `internal/eml/simd_dispatch_arm64.go` (`arm64AddSIMDF32`, `arm64MulSIMDF32`, etc.) to invoke the NEON routines.
3. Update `simd_f32.go` to dispatch to ARM64 hardware vector routines when supported.

**Verify:**
- Validate numerical equivalence against Go `math` float32 conversions across subnormals and infinity in `simd_f32_test.go`.
- Assert zero heap allocations in `*To` variants.
- Benchmark throughput scaling via `BenchmarkExpSIMDF32VsF64`.

---

### Step 3 — GPU CUDA Stream Overlap & Pinned Host Memory Optimization (P2)

**Status:** Open

**Evidence & Context:**
- Empirical benchmarks in `docs/performance.md` show GPU speedups of 1.10x–2.1x over 16-thread CPU execution for large vectors ($N \ge 16,384$).
- For smaller vectors ($N \le 4,096$), host-to-device memory copy overhead (`cudaMemcpy`) accounts for over 90% of total elapsed time, making CPU execution faster.
- Pageable host memory requires the CUDA driver to copy data into an internal staging buffer before transferring over PCIe.

**Change:**
1. Add pinned host buffer allocation APIs in `internal/gpu` wrapping `cudaHostAlloc` / `cudaFreeHost`.
2. Implement dual-buffered asynchronous execution using multiple CUDA streams (`cudaStreamCreateWithFlags`):
   - Stream 0: Transfer chunk $k+1$ host $\to$ device (`cudaMemcpyAsync`).
   - Stream 1: Execute kernel computation on chunk $k$.
   - Stream 0: Transfer chunk $k-1$ device $\to$ host.
3. Expose zero-copy unified memory buffers for workloads running iterative pipeline operations.

**Verify:**
- Measure end-to-end latency vs vector size $N$ on CUDA hardware (SM 8.9 / Ada Lovelace).
- Verify the CPU-to-GPU speedup crossover shifts from $N \approx 16,384$ down to $N \approx 4,096$.
- Ensure graceful fallback to CPU when CUDA device memory is exhausted or unavailable.

---

## Historical Archive & Completed Milestones Reference

All prior roadmap steps (1–8, 11–30, 31–50, and AVX-512 optimization items) have been completed, verified, and archived:
- **v0.4**: Bytecode VM expanded to 27 functions, Genetic Algorithm symbolic recovery driver, zero-alloc batch evaluation, and pipeline memory reuse.
- **v0.5**: JIT register clobber fixes, 100% statement coverage across packages, `purego` portability across all architectures, subnormal accuracy oracles, 14 native ARM64 NEON kernels, and WebAssembly browser benchmark harness (`cmd/wasmbench`).
- **v0.6**: Automated GHCR package publishing, Go 1.27.2 toolchain, CPU AVX2 & CUDA multi-stage Docker containers, supply chain security attestations (SLSA provenance, SPDX SBOM, Govulncheck).
- **v0.7 (AVX-512 & VNNI Optimization)**:
  - **24 AVX-512 Assembly Kernels in `simd_amd64.s`**:
    - 11 elementwise float64 kernels (`add`, `sub`, `mul`, `div`, `addScalar`, `mulScalar`, `sqrt`, `fma`, `abs`, `neg`, `inv`) unrolled 2x (16 elements / 128 bytes per iteration) with 8-element odd tail.
    - 5 transcendental kernels (`expAVX512`, `logAVX512`, `sinAVX512`, `cosAVX512`, `tanAVX512`) with loop-invariant minimax polynomial coefficients hoisted to dedicated registers `Z16–Z31`.
    - 5 float32 kernels (`addAVX512F32`, `mulAVX512F32`, `subAVX512F32`, `divAVX512F32`, `sqrtAVX512F32`) processing 16 elements per instruction.
    - 2 saturated int8 kernels (`addSatInt8AVX512`, `subSatInt8AVX512`) processing 64 bytes per iteration.
    - 1 VNNI dot product kernel (`dotInt8AVX512VNNI`) using `VPDPBUSD` 4-way multiply-accumulate on byte codebook weights and signed activations in `pkg/quant` (TurboQuant4).
  - All 24 AVX-512 kernels emit `VZEROUPPER` before `RET` (`TestAVX512KernelsEmitVzeroupper` passes 24/24).
  - Verified on dedicated Intel Xeon Cascade Lake AVX-512 reference platform: `cmd/validate` (435/435 pass), `cmd/bench -accuracy` (13/13 bounds pass), and full project test suites pass with 0 regressions.

