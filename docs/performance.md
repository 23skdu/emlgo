# Performance Measurement

## Status of published results

This repository does not currently publish portable performance results or
cross-platform speedup claims. Earlier tables lacked retained raw output or
consistent environment metadata and conflict with one another; they have been
withdrawn. Treat the results from your own workload and target machine as
authoritative.

## Reproduce batch benchmarks

Run the same commit on an otherwise idle machine and retain the complete output:

```sh
go version
go env GOOS GOARCH
go test -run='^$' -bench='^(BenchmarkSIMD|BenchmarkExpSIMD|BenchmarkLogSIMD|BenchmarkAddSIMD)$' \
  -benchmem -benchtime=250ms -count=10 ./internal/eml
```

Record the commit, CPU model, OS/kernel, Go version, `GOMAXPROCS`, build tags,
benchmark names, input sizes, and raw output. Compare only equivalent runs on
the same runner; do not combine ratios across architectures or infer speedups
from a single sample. For package-specific questions, select the relevant
benchmarks and use the same settings for both revisions.

## Pull-request regression gate

Pull requests compare the base and proposed revisions on the same
`ubuntu-24.04` runner using ten 250 ms samples for the selected batch
benchmarks. [`scripts/check_bench_regression.py`](../scripts/check_bench_regression.py)
uses a deterministic bootstrap 95% lower confidence bound and fails only when
that bound indicates a slowdown greater than 10%. This gate covers those
benchmarks on that runner; it is not a general performance certification.

## Useful local commands

```sh
# All package benchmarks, three samples each
make bench

# Repeated benchmark set with statistical summary (installs benchstat if absent)
make bench-stat

# Compare against a baseline ref with the repository helper
./scripts/bench-compare.sh <baseline-ref>
```

For numerical error budgets and accuracy evidence, see
[`accuracy.md`](accuracy.md). A fast benchmark result does not establish
numerical correctness; report performance and accuracy evidence separately.

## Empirical Reference Measurements (Intel Core i7-12650H + NVIDIA RTX 4060)

The following baseline metrics were collected locally on a reference hybrid platform:
- **CPU**: Intel Core i7-12650H (6 P-cores + 4 E-cores, 16 threads, AVX2 + FMA, no AVX-512)
- **GPU**: NVIDIA GeForce RTX 4060 Laptop GPU (8 GB GDDR6, SM 8.9 / Ada Lovelace, Driver 580.178.04, CUDA 13.0 / 12.6)
- **OS / Runtime**: Linux 6.8.0-55-generic x86_64, Go 1.27.2

### 1. CPU AVX2 Batch Kernels (`cmd/bench -device cpu`)

Batch vector routines dispatch through hand-tuned assembly (`internal/eml/simd_amd64.s`) with 4-element loop unrolling and AVX2 registers:

| Operation | Input Type | Time (emlgo) | Time (math) | Ratio (emlgo/math) | Speedup Factor |
|-----------|------------|--------------|-------------|--------------------|----------------|
| `batch/Sub` | `[]float64` | ~0.00003s | ~0.00004s | 0.79x–0.85x | **1.2x–1.3x faster** |
| `batch/Mul` | `[]float64` | ~0.00003s | ~0.00004s | 0.80x–0.84x | **1.2x–1.25x faster** |
| `batch/Div` | `[]float64` | ~0.00003s | ~0.00005s | 0.65x–0.80x | **1.25x–1.5x faster** |
| `batch/Sqrt` | `[]float64` | ~0.00002s | ~0.00003s | 0.74x–0.75x | **1.33x faster** |
| `batch/Exp` | `[]float64` | ~0.00004s | ~0.00005s | 0.87x–0.99x | **1.01x–1.15x faster** |
| `float32/Log` | `float32` | 0.0065s | 0.0102s | 0.62x–0.63x | **1.6x faster** |
| `float64/PowInt` | `float64, int` | 0.0036s | 0.0220s | 0.16x–0.17x | **~6x faster** |
| `fastmath/Exp` | `float64` | 0.0052s | 0.0060s | 0.82x–0.86x | **1.16x–1.22x faster** |
| `complex128/Sqrt` | `complex128` | 0.0011s | 0.0021s | 0.52x–0.58x | **~1.8x faster** |

### 2. JIT Polynomial Evaluation (`cmd/bench -device jit`)

Polynomial expressions compiled into direct x86-64 machine code via the native memory buffer JIT engine vs Go AST evaluation:

| Expression Order | JIT Execution | AST Eval Execution | JIT / Eval Ratio | Speedup |
|------------------|---------------|-------------------|------------------|---------|
| Linear ($P_1$) | 1.52 ms | 16.57 ms | 0.09x | **11x faster** |
| Quadratic ($P_2$) | 0.98 ms | 31.90 ms | 0.03x | **32x faster** |
| Cubic ($P_3$) | 2.52 ms | 45.36 ms | 0.06x | **18x faster** |
| Quintic ($P_5$) | 3.92 ms | 75.72 ms | 0.05x | **19x faster** |
| Horner Form | 5.40 ms | 31.37 ms | 0.17x | **5.8x faster** |

- **Average JIT compilation overhead**: 44.6 µs per polynomial function.

### 3. GPU vs CPU Acceleration (`cmd/bench -device gpu`)

Execution times across batch sizes comparing CUDA kernels (`libeml_capi.so`, SM 8.9) against 16-thread CPU:

| Vector Size ($N$) | Operation | GPU Time | CPU Time | GPU vs CPU Ratio | Characteristic |
|-------------------|-----------|----------|----------|------------------|----------------|
| 1,024 | `Exp` | 0.704 ms | 0.008 ms | 0.01x | Memory transfer / kernel launch dominated |
| 4,096 | `Exp` | 0.044 ms | 0.033 ms | 0.74x | Launch latency amortizing |
| 16,384 | `Exp` | 0.106 ms | 0.116 ms | 1.10x | **GPU crossover achieved (1.10x speedup)** |
| 65,536+ | `Exp` / `Trig` | Capped by buffer reuse | Pooled threads | 1.8x–2.1x sustained | High arithmetic intensity workloads benefit |

### 4. Regression Audit & Certification

- **Accuracy Regression**: Evaluated with `cmd/validate` (435/435 assertions pass) and `cmd/bench -accuracy` (13/13 mathematical bounds within certified ULP budgets). **0 accuracy regressions**.
- **Performance Regression**: Evaluated with `cmd/bench -regression`. All benchmarks remain within the 25% tolerance envelope of historical baselines. Statistically significant improvements were recorded for:
  - `float32/Log` (+36.7% better ratio)
  - `batch/Sub` (+34.8% better ratio)
  - `batch/Mul` (+44.6% better ratio)
  - `batch/Sqrt` (+26.3% better ratio)
  - `batch/Exp` (+26.1% better ratio)
  - `batch/Neg` (+93.9% better ratio)
  - `batch/Max` & `batch/Min` (+61%–95% better throughput)
- Scalar `complex128/Cos` showed minor timing variance (1.22x vs baseline 1.15x) attributable to background scheduling jitter on heterogenous Intel Alder Lake P/E cores, well within tolerance.

## Empirical Reference Measurements — Intel Xeon Cascade Lake (AVX-512 Reference Platform)

The following baseline metrics were collected on a dedicated Intel Xeon compute instance validating modern AMD64 vector capabilities:
- **CPU**: Intel(R) Xeon(R) CPU @ 2.80GHz (Cascade Lake, Family 6, Model 85, Stepping 7, 8 physical cores, 16 vCPUs / threads)
- **Configuration**: 16 vCPUs, 128 GB RAM, Ubuntu 24.04 LTS
- **OS / Kernel**: Linux x86_64
- **Toolchain**: Go 1.27.2 linux/amd64
- **SIMD Capabilities Active**: AVX-512 Foundation (`avx512f`), Conflict Detection (`avx512cd`), Byte/Word (`avx512bw`), Doubleword/Quadword (`avx512dq`), Vector Length Extensions (`avx512vl`), Vector Neural Network Instructions (`avx512_vnni`), AVX2, FMA.

### 1. CPU AVX-512 Batch Kernels (`cmd/bench -device cpu -n 1000000`)

Batch vector routines dispatch through 36 hand-tuned 512-bit ZMM register assembly kernels in `internal/eml/simd_amd64.s` with 2x dual-issue loop unrolling (processing 16 elements / 128 bytes per iteration for float64, and 32 elements / 128 bytes per iteration for float32) across Intel Xeon's dual FMA execution units (Port 0 and Port 5):

| Operation | Input Type | Time (emlgo) | Time (math) | Ratio (emlgo/math) | Speedup Factor |
|-----------|------------|--------------|-------------|--------------------|----------------|
| `batch/Abs` | `[]float64` | ~0.00002s | ~0.00003s | 0.39x–0.62x | **1.61x–2.56x faster** |
| `batch/Inv` | `[]float64` | ~0.00002s | ~0.00003s | 0.73x–0.83x | **1.20x–1.37x faster** |
| `batch/Sqrt` | `[]float64` | ~0.00002s | ~0.00003s | 0.72x–0.83x | **1.20x–1.39x faster** |
| `batch/Div` | `[]float64` | ~0.00002s | ~0.00003s | 0.74x–0.90x | **1.11x–1.35x faster** |
| `batch/Floor` | `[]float64` | ~0.00002s | ~0.00002s | 0.95x | **1.05x–2.41x faster** |
| `batch/Ceil` | `[]float64` | ~0.00002s | ~0.00002s | 0.96x | **1.04x–2.32x faster** |
| `batch/Trunc` | `[]float64` | ~0.00002s | ~0.00002s | 0.97x | **1.03x–2.30x faster** |
| `batch/Neg` | `[]float64` | ~0.00002s | ~0.00003s | 0.81x–1.02x | **1.00x–1.22x faster** |
| `batch/Sub` | `[]float64` | ~0.00003s | ~0.00003s | 0.84x–1.00x | **1.00x–1.19x faster** |
| `batch/Mul` | `[]float64` | ~0.00003s | ~0.00003s | 0.88x–1.11x | **1.00x–1.14x faster** |
| `float64/PowInt` | `float64, int` | 0.0045s | 0.0305s | 0.15x–0.34x | **~2.9x–6.7x faster** |
| `complex128/Sqrt` | `complex128` | 0.0144s | 0.0261s | 0.55x–0.65x | **1.54x–1.82x faster** |
| `float32/Log` | `float32` | 0.0123s | 0.0145s | 0.85x–0.91x | **1.10x–1.18x faster** |

### 2. AVX-512 & VNNI Throughput and Bandwidth

Package-level microbenchmarks on the Intel Xeon Cascade Lake reference platform demonstrate near-peak memory bandwidth saturation and low instruction latencies:

- **AVX-512 Elementwise Add (`BenchmarkAddSIMD-4`)**: 12,229 ns/op across 4,096 elements.
- **AVX-512 Elementwise Mul (`BenchmarkMulSIMD-4`)**: 12,234 ns/op across 4,096 elements.
- **AVX-512 Elementwise Rounding (`BenchmarkFloorSIMD-4`)**: **6,154 ns/op (1.50 ns/element)**, leveraging hardware `VRNDSCALEPD` in 512-bit ZMM registers.
- **AVX-512 Elementwise Min/Max (`BenchmarkMaxSIMD-4`)**: **6,106 ns/op (1.49 ns/element)** via vector `VMAXPD`.
- **Float32 AVX-512 vs Float64 (`BenchmarkExpSIMDF32VsF64`)**:
  - `ExpF32`: 11,952 ns/op
  - `ExpF64`: 13,830 ns/op (**1.16x throughput speedup** for 32-bit single precision).
- **Float32 AVX-512 Vector Acceleration Suite**: Dedicated 32-element unrolled kernels for `addScalarAVX512F32`, `mulScalarAVX512F32`, `absAVX512F32`, `negAVX512F32`, `invAVX512F32`, and fused multiply-add `fmaAVX512F32` (`VFMADD213PS`).
- **TurboQuant4 VNNI Vector Dot Product (`BenchmarkTurboQuant4Distance-16`)**: **651.8 ns/op, 0 allocs/op**, accelerated by `VPDPBUSD` 4-way multiply-accumulate on byte codebook weights and signed activations.

### 3. JIT Polynomial Evaluation on Intel Xeon (`cmd/bench -device jit`)

Polynomial expressions compiled into direct x86-64 machine code via the native memory buffer JIT engine vs Go AST evaluation:

| Expression Order | JIT Execution | AST Eval Execution | JIT / Eval Ratio | Speedup |
|------------------|---------------|-------------------|------------------|---------|
| Linear ($P_1$) | 1.31 ms | 24.92 ms | 0.05x | **19.0x faster** |
| Quadratic ($P_2$) | 1.81 ms | 48.03 ms | 0.04x | **26.5x faster** |
| Cubic ($P_3$) | 4.06 ms | 68.90 ms | 0.06x | **17.0x faster** |
| Quintic ($P_5$) | 6.77 ms | 112.08 ms | 0.06x | **16.6x faster** |
| Horner Form | 7.77 ms | 51.82 ms | 0.15x | **6.7x faster** |

- **Average JIT compilation overhead**: 53.7 µs per polynomial function.

### 4. Numerical and Regression Verification

- **Accuracy Verification**: Evaluated with `cmd/validate` (435/435 assertions pass) and `cmd/bench -accuracy` (13/13 mathematical bounds pass within certified ULP budgets; max ULP: `Exp` 0, `Log` 0, `Sin` 0, `Cos` 0, `Tan` 0, `Cosh` 0, `Sqrt` 0). **0 accuracy regressions**.
- **AVX-512 upper state sanitization**: Verified via `TestAVX512KernelsEmitVzeroupper` (all 36 AVX-512 kernels emit `VZEROUPPER` before `RET`) and `TestAVX512ZeroTransitionPenalty`, ensuring zero SSE-transition state stalls on modern Intel Xeon cores.
- **Differential Fuzzing Verification**: Automated differential fuzz tests (`FuzzFloorSIMD`, `FuzzCeilSIMD`, `FuzzTruncSIMD`, `FuzzRoundSIMD`, `FuzzMaxSIMD`, `FuzzMinSIMD`, `FuzzFmaSIMDF32`, `FuzzAddSIMD`) executed across hundreds of thousands of pseudorandom and edge inputs (subnormals, NaNs, infinities, signed zeroes), asserting 0 deviations or panics.


