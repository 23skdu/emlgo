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

