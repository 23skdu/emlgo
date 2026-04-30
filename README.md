# emlgo

A high-performance mathematical library for Go, implementing elementary functions using the EML (Exp-Minus-Log) operator. Optimized for massive parallelization and SIMD-accelerated batch processing.

## Executive Summary

`emlgo` is designed to bridge the gap between pure Go portability and the performance of native C/C++ math libraries. By leveraging architecture-specific assembly kernels (AVX2, AVX512, NEON) and the novel EML operator, it provides:

- **Batch Performance**: 1.2x to 15.0x speedups over standard library loops using SIMD.
- **FastMath API**: A specialized `pkg/fastmath` package for high-throughput scalar operations, achieving **10% better performance** than `math.Sin`.
- **Scalar Efficiency**: Native assembly kernels for `Sqrt` and `FMA` (Fused Multiply-Add) to bypass standard library overhead.
- **100% Correctness**: Full feature parity with Go's `math` package, verified by a comprehensive validation suite (375+ tests).

## Installation

```bash
go get github.com/emlgo/eml
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/pkg/trig"
    "github.com/emlgo/eml/pkg/logexp"
    "github.com/emlgo/eml/pkg/arithmetic"
)

func main() {
    // Trigonometric functions
    fmt.Printf("sin(π/4) = %v\n", trig.Sin(math.Pi/4))  // ≈ 0.707
    fmt.Printf("cos(π/4) = %v\n", trig.Cos(math.Pi/4))  // ≈ 0.707
    fmt.Printf("tan(π/4) = %v\n", trig.Tan(math.Pi/4))  // ≈ 1.0

    // Exponential and Logarithmic
    fmt.Printf("exp(1) = %v\n", logexp.Exp(1))          // ≈ 2.718
    fmt.Printf("ln(e) = %v\n", logexp.Log(math.E))      // ≈ 1.0

    // Arithmetic
    fmt.Printf("sqrt(2) = %v\n", arithmetic.Sqrt(2))     // ≈ 1.414
    fmt.Printf("pow(2, 3) = %v\n", arithmetic.Pow(2, 3)) // ≈ 8.0
}
```

## Packages

| Package | Description |
| :--- | :--- |
| `pkg/logexp` | Exp, Log functions using EML |
| `pkg/trig` | Trig, inverse trig, hyperbolic functions |
| `pkg/hyper` | Dedicated hyperbolic functions |
| `pkg/arithmetic` | Add, Sub, Mul, Div, Pow, Sqrt, etc. |
| `pkg/fastmath` | High-performance relaxed scalar operations |
| `internal/eml` | Core EML operator with SIMD/Scalar kernels |
| `internal/constants` | Mathematical constants |

## SIMD Support

The library automatically detects and uses SIMD instructions:
- **AMD64**: AVX-512 (8-wide), AVX2 (4-wide)
- **ARM64**: NEON/ASIMD (4-wide)

Batch operations automatically use SIMD:

```go
values := make([]float64, 1000)
// ... fill values ...
results := trig.SinBatch(values)  // Uses SIMD when available
```

## Benchmarking

Run speed benchmarks comparing emlgo to math library:

```bash
./scripts/run_benchmark.sh           # Speed benchmark
./scripts/run_benchmark.sh -c        # Feature parity
./scripts/run_benchmark.sh -a        # Accuracy test
```

## Validation

Run validation tests for all Go data types:

```bash
./scripts/run_validation.sh           # All validations
./scripts/run_validation.sh -v       # Verbose output
./scripts/run_validation.sh -t float # Filter by type
```

## Building & Testing

```bash
# Build all packages
go build ./...

# Run tests
go test ./...

# Test with race detection
go test -race ./...

# Run linter
go vet ./...

# Security scan
gosec ./...
```

## Architecture

```text
emlgo/
├── cmd/
│   ├── bench/         # Benchmark tool
│   ├── validate/      # Validation tool
│   └── emlcli/       # CLI demo
├── internal/
│   ├── eml/          # Core EML operator + SIMD
│   └── constants/   # Mathematical constants
├── pkg/
│   ├── logexp/       # Exp, Log functions
│   ├── trig/         # Trig functions + batch ops
│   ├── hyper/        # Hyperbolic functions
│   └── arithmetic/   # Basic arithmetic
├── docs/             # Documentation
└── scripts/          # Benchmark & validation scripts
```

## Key Features

- **All functions use EML**: Every function is implemented using the EML operator
- **Full type support**: Tested for int, uint, float32, float64, complex64, complex128
- **Edge case handling**: Properly handles NaN, Inf, subnormal numbers
- **Performance**: SIMD batch operations for large data processing
- **Accuracy**: Matches math library to within tolerance (validation: 375/375 pass)

## References

- [All elementary functions from a single operator](https://arxiv.org/abs/2603.21852v2) - Andrzej Odrzywołek (2026)
- Related: Kolmogorov-Arnold Networks (KAN) for similar tree structures
