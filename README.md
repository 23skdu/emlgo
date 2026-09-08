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
| `internal/constants` | Mathematical constants (e, π, ln2, √2, φ) |

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

## Complex Numbers

```go
import "github.com/emlgo/eml/internal/eml"

// Batch complex operations
z := []complex128{1 + 2i, 3 + 4i, 5 + 6i}
exp := eml.ComplexExpBatch(z)
sin := eml.ComplexSinBatch(z)

// Trigonometric identities hold
// sin²(z) + cos²(z) = 1 verified at complex128 precision
```

## JIT Compiler

```go
import "github.com/emlgo/eml/internal/jit"

c := jit.NewCompiler()

// Integer and non-integer exponents
f, _ := c.Compile("x^0.5")    // sqrt(x)
g, _ := c.Compile("x^(2*x)")  // variable exponent

// 16 math functions
h, _ := c.Compile("sin(x)^2 + cos(x)^2") // = 1.0

// Zero-allocation arena
arena := jit.NewArena(64)
node := arena.FromInterface(ast)
result := jit.EvalArena(node, 3.14)
```

## Numerical Stability

```go
// Pow near x=1 uses Log1p for stability
arithmetic.Pow(1.0+1e-15, 1e6) // accurate, not catastrophic cancellation

// Exp near x=0 uses Expm1
logexp.Exp(1e-15) // accurate to ~1e-25

// Log near x=1 uses Log1p
arithmetic.Log(1.0+1e-15) // accurate to ~1e-25
```

## Arbitrary Precision

```go
import "github.com/emlgo/eml/internal/eml/bigmath"

x := bigmath.NewFloat(1.0)
y := bigmath.NewFloat(1.0)
result := bigmath.Eml(x, y) // exp(1) - log(1) = e

// Verify symbolic identities
v := bigmath.DefaultVerifier()
sin2pluscos2 := func(x *big.Float) *big.Float {
    s := bigmath.Sin(x)
    c := bigmath.Cos(x)
    return new(big.Float).Add(new(big.Float).Mul(s, s), new(big.Float).Mul(c, c))
}
one := func(x *big.Float) *big.Float { return bigmath.NewFloat(1.0) }
v.VerifyIdentity(sin2pluscos2, one) // true
```

## Building & Testing

```bash
go build ./...                              # Build
go test ./...                               # Test
go test -race ./...                         # Race detection
go vet ./...                                # Lint
gosec -exclude-generated ./...              # Security scan
./scripts/bench-compare.sh                  # Benchmark regression
```

## CI

GitHub Actions workflow (`.github/workflows/ci.yml`):
- Go 1.22 and 1.23 matrix
- `go vet`, `go test -race`, `gosec`, `golangci-lint`

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
