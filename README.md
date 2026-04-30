# emlgo

A high-performance Go library implementing all elementary mathematical functions using the EML (Exp-Minus-Log) operator discovered in arXiv:2603.21852v2.

## Overview

The EML operator `eml(x, y) = exp(x) - ln(y)` can reconstruct all elementary functions when combined with the constant `1`. This library implements:

- Core EML operator with SIMD support
- Exponential and logarithmic functions
- Trigonometric functions (sin, cos, tan, asin, acos, atan)
- Hyperbolic functions (sinh, cosh, tanh, asinh, acosh, atanh)
- Arithmetic operations (add, mul, div, pow, sqrt)

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
)

func main() {
    fmt.Printf("sin(0) = %v\n", trig.Sin(0))
    fmt.Printf("cos(0) = %v\n", trig.Cos(0))
}
```

## Architecture

```
internal/
  eml/          # Core EML operator + SIMD
  constants/    # Mathematical constants
pkg/
  logexp/       # Exp, Log functions
  trig/         # Trig functions
  hyper/        # Hyperbolic functions
  arithmetic/   # Basic arithmetic
cmd/
  emlcli/       # CLI demo
```

## Building

```bash
go build ./...
go test ./...
go test -race ./...
```

## SIMD Support

The library automatically detects and uses AVX2/AVX512 instructions on supported platforms.

## References

- [All elementary functions from a single operator](https://arxiv.org/abs/2603.21852) - Andrzej Odrzywołek (2026)