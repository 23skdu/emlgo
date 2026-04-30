# emlgo

A pure Go library implementing all elementary mathematical functions using the EML (Exp-Minus-Log) operator discovered in [arXiv:2603.21852v2](https://arxiv.org/abs/2603.21852) by Andrzej Odrzywołek.

## Overview

The EML operator `eml(x, y) = exp(x) - ln(y)` can reconstruct all elementary functions when combined with the constant `1`. This library provides:

- **100% feature parity** with Go's `math` package for all data types
- **SIMD support** with automatic detection for AVX2/AVX512/NEON
- **Comprehensive testing** - 375 validation tests pass for all Go types
- **No external dependencies** - Pure Go implementation

## Installation

```bash
go get github.com/emlgo/eml
```

## Quick Start

```go
package main

import (
    "fmt"
    "math"

    "github.com/emlgo/eml/pkg/arithmetic"
    "github.com/emlgo/eml/pkg/hyper"
    "github.com/emlgo/eml/pkg/logexp"
    "github.com/emlgo/eml/pkg/trig"
)

func main() {
    // Trigonometric
    fmt.Printf("sin(π/4) = %.6f (math: %.6f)\n", trig.Sin(math.Pi/4), math.Sin(math.Pi/4))
    fmt.Printf("cos(π/4) = %.6f (math: %.6f)\n", trig.Cos(math.Pi/4), math.Cos(math.Pi/4))
    fmt.Printf("tan(π/4) = %.6f (math: %.6f)\n", trig.Tan(math.Pi/4), math.Tan(math.Pi/4))

    // Exponential & Logarithmic
    fmt.Printf("exp(1) = %.6f (math: %.6f)\n", logexp.Exp(1), math.Exp(1))
    fmt.Printf("ln(e) = %.6f (math: %.6f)\n", logexp.Log(math.E), math.Log(math.E))

    // Hyperbolic
    fmt.Printf("sinh(1) = %.6f (math: %.6f)\n", hyper.Sinh(1), math.Sinh(1))
    fmt.Printf("cosh(1) = %.6f (math: %.6f)\n", hyper.Cosh(1), math.Cosh(1))
    fmt.Printf("tanh(1) = %.6f (math: %.6f)\n", hyper.Tanh(1), math.Tanh(1))

    // Arithmetic
    fmt.Printf("sqrt(2) = %.6f (math: %.6f)\n", arithmetic.Sqrt(2), math.Sqrt(2))
    fmt.Printf("pow(2, 10) = %.0f (math: %.0f)\n", arithmetic.Pow(2, 10), math.Pow(2, 10))
}
```

## Packages

| Package | Description |
|---------|-------------|
| `pkg/logexp` | Exp, Log functions using EML |
| `pkg/trig` | Trig, inverse trig, hyperbolic functions |
| `pkg/hyper` | Dedicated hyperbolic functions |
| `pkg/arithmetic` | Add, Sub, Mul, Div, Pow, Sqrt, etc. |
| `internal/eml` | Core EML operator with SIMD |
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

```
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
