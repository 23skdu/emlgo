# emlgo - Elementary Functions from the EML Operator

A pure Go implementation of all elementary mathematical functions using only the EML (ExpMinusLog) operator, based on the research presented in [arXiv:2603.21852v2](https://arxiv.org/html/2603.21852v2).

## Overview

The emlgo library provides a complete set of elementary mathematical functions (trigonometric, hyperbolic, exponential, logarithmic, and arithmetic operations) implemented using a single primitive operator:

```math
eml(x, y) = exp(x) - ln(y)
```

This approach enables:

- **No external dependencies** - Pure Go with optional SIMD support
- **Unified implementation** - All functions derived from a single primitive
- **SIMD optimizations** - Batch processing with AVX2/AVX512/NEON support

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
| `pkg/logexp` | Exponential and logarithmic functions |
| `pkg/trig` | Trigonometric and inverse trigonometric functions |
| `pkg/hyper` | Hyperbolic and inverse hyperbolic functions |
| `pkg/arithmetic` | Basic arithmetic operations, roots, powers |
| `pkg/fastmath` | High-performance scalar alternatives |
| `internal/eml` | Core EML operator + SIMD/Scalar kernels |
| `internal/constants` | Mathematical constants |

## Architecture

For details on the internal design and SIMD dispatch logic, see [Architecture](architecture.md).

## Features

- **Full float64 domain support** - Handles NaN, Inf, edge cases correctly
- **Comprehensive test coverage** - All functions verified against math library
- **SIMD batch operations** - Efficient processing of slice inputs
- **Race-condition safe** - Tested with `-race` flag
- **Zero security issues** - Verified with gosec

## Performance

The library provides comparable accuracy to the standard math library while being implemented entirely through the EML operator. For batch operations, SIMD optimizations provide significant speedups on supported architectures.

## Requirements

- Go 1.21 or later
- Optional: `golang.org/x/sys/cpu` for SIMD detection

## License

MIT License - See LICENSE file for details

## References

- [EML Operator Paper (arXiv:2603.21852v2)](https://arxiv.org/html/2603.21852v2)
- Related: Kolmogorov-Arnold Networks (KAN) for similar tree structures
