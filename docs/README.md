# emlgo

A high-performance mathematical library for Go, implementing all elementary functions using the EML (Exp-Minus-Log) operator with SIMD acceleration, JIT compilation, GPU backends, and arbitrary-precision verification.

Based on the research of **Andrzej Odrzywołek**: [All elementary functions from a single operator](https://arxiv.org/abs/2603.21852v2) (2026).

## Features

- **EML Operator**: `eml(x, y) = exp(x) - ln(y)` — single primitive from which all elementary functions derive
- **SIMD Batch Operations**: AVX2, AVX-512 (AMD64), NEON, SVE (ARM64), WASM SIMD128
- **JIT Compiler**: x86-64 machine code generation for math expressions (16 functions, non-integer/variable exponents)
- **GPU Backends**: CUDA (Linux/Windows) and Metal (macOS/ARM64)
- **Complex Numbers**: First-class `complex128` support via `math/cmplx`
- **Arbitrary Precision**: `math/big.Float` backend for symbolic verification
- **Zero-Allocation AST**: Arena allocator for JIT parse/eval paths
- **Canonical EML Trees**: Map any expression to minimal EML form
- **Numerically Stable**: Log1p/Expm1 optimization for compound expressions
- **FastMath**: FMA-optimized polynomial approximations (~10% faster than `math.Sin`)

## Installation

```bash
go get github.com/emlgo/eml
```

Requires **Go 1.23+**.

## Quick Start

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/pkg/trig"
    "github.com/emlgo/eml/pkg/logexp"
    "github.com/emlgo/eml/pkg/arithmetic"
    "github.com/emlgo/eml/pkg/hyper"
)

func main() {
    // Trigonometric
    fmt.Printf("sin(π/4) = %.6f\n", trig.Sin(math.Pi/4))

    // Exponential & Logarithmic
    fmt.Printf("exp(1) = %.6f\n", logexp.Exp(1))

    // Arithmetic with numerical stability
    fmt.Printf("pow(1.000001, 1e6) = %.6f\n", arithmetic.Pow(1.000001, 1e6))

    // Hyperbolic
    fmt.Printf("sinh(1) = %.6f\n", hyper.Sinh(1))

    // Batch operations (SIMD-accelerated)
    x := []float64{0, 0.5, 1.0, 1.5, 2.0}
    sin := trig.SinBatch(x)
    exp := logexp.ExpBatch(x)
    fmt.Printf("SinBatch: %v\n", sin)
    fmt.Printf("ExpBatch: %v\n", exp)
}
```

## Packages

| Package | Description |
| :--- | :--- |
| `pkg/arithmetic` | Add, Sub, Mul, Div, Pow, Sqrt, Cbrt, FMA, GCD, LCM, batch ops |
| `pkg/trig` | Sin, Cos, Tan, Cot, Sec, Csc, Asin, Acos, Atan, Atan2, batch ops |
| `pkg/hyper` | Sinh, Cosh, Tanh, Asinh, Acosh, Atanh, batch ops |
| `pkg/logexp` | Exp, Log, ExpBatch, LogBatch, ExpFast, LogFast |
| `pkg/fastmath` | FMA-optimized polynomial approximations for Exp, Sin, Cos, Log |
| `internal/eml` | Core EML operator, SIMD dispatch, worker pool, complex batch ops |
| `internal/eml/bigmath` | Arbitrary-precision EML via `math/big.Float`, identity verifier |
| `internal/jit` | JIT compiler (x86-64 codegen), arena allocator, canonical EML trees |
| `internal/gpu` | CUDA & Metal GPU backends, ULP-based verification |
| `internal/constants` | Mathematical constants (e, π, ln2, √2, φ, etc.) |

## Architecture

```text
emlgo/
├── cmd/
│   ├── bench/           # Benchmark tool
│   ├── validate/        # Validation tool
│   └── emlcli/          # CLI demo
├── internal/
│   ├── eml/             # Core EML operator + SIMD dispatch
│   │   ├── bigmath/     # Arbitrary-precision backend
│   │   └── simd_*.go    # Platform-specific dispatch (amd64/arm64/wasm)
│   ├── jit/             # JIT compiler + arena allocator + canonical trees
│   ├── gpu/             # CUDA & Metal GPU backends
│   └── constants/       # Mathematical constants
├── pkg/
│   ├── arithmetic/      # Basic arithmetic + batch ops
│   ├── trig/            # Trigonometric + batch ops
│   ├── hyper/           # Hyperbolic functions + batch ops
│   ├── logexp/          # Exponential & logarithmic
│   └── fastmath/        # High-performance scalar ops
├── docs/                # Documentation
└── scripts/             # Benchmark & validation scripts
```

## SIMD Support

| Architecture | Instructions | Width |
| :--- | :--- | :--- |
| AMD64 | AVX-512 | 8-wide float64 |
| AMD64 | AVX2 + FMA | 4-wide float64 |
| ARM64 | NEON | 2-wide float64 |
| ARM64 | SVE/SVE2 | Scalable |
| WASM | SIMD128 | 2-wide float64 (8-wide unrolled) |

Batch operations automatically dispatch to the fastest available SIMD path.

## Building & Testing

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
gosec -exclude-generated ./...
./scripts/bench-compare.sh
```

## Performance

| Operation | Scalar | Batch (SIMD) |
| :--- | :--- | :--- |
| Add/Sub/Mul | ~parity with `math` | **1.2-15x faster** |
| Exp/Log/Sin/Cos | ~parity with `math` | **1.1-5x faster** |
| PowInt | **5-6x faster** than `math.Pow` | parallelized |
| FastMath Sin | **10% faster** than `math.Sin` | N/A |
| Fused (ExpMul) | N/A | **20-30% less memory** |

## License

See LICENSE file.
