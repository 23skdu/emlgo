# emlgo

A high-performance mathematical library for Go built around the EML (Exp-Minus-Log) operator, with SIMD acceleration, JIT compilation, GPU backends, and arbitrary-precision verification.

The scalar elementary functions delegate to the Go standard library (`math`, `math/cmplx`) so they stay bit-exact with `math`. EML provides the theoretical foundation and the machinery built on it: canonical EML expression trees, identity reductions, symbolic differentiation, and a JIT compiler that emits x86-64 machine code directly.

Based on the research of **Andrzej Odrzywołek**: [All elementary functions from a single operator](https://arxiv.org/abs/2603.21852v2) (2026).

## Features

- **EML Operator**: `eml(x, y) = exp(x) - ln(y)` — single primitive from which all elementary functions derive
- **SIMD Batch Operations**: AVX2, AVX-512 (AMD64), NEON, SVE (ARM64), WASM SIMD128
- **float32 SIMD**: Dedicated `float32` batch operations for memory-constrained workloads
- **JIT Compiler**: x86-64 machine code generation for math expressions (25 functions, non-integer/variable exponents)
- **JIT Expression Cache**: LRU cache with `CompileCached` for repeated compilations
- **GPU Backends**: CUDA (Linux/Windows) and Metal (macOS/ARM64)
- **Complex Numbers**: First-class `complex128` support via `math/cmplx`
- **Arbitrary Precision**: `math/big.Float` backend for symbolic verification
- **Symbolic Regression**: `pkg/bytecode` GA recovers closed forms from data
- **Zero-Allocation AST**: Arena allocator for JIT parse/eval paths
- **Canonical EML Trees**: Map any expression to minimal EML form
- **Symbolic Differentiation**: `Diff()` with chain rule and simplification
- **Expression Simplification**: Constant folding, identity reduction, algebraic simplifications
- **EML Decompiler**: `Decompile()` and `DecompileLaTeX()` for EML tree → infix/LaTeX conversion
- **Composable Pipeline**: Zero-allocation `Pipeline` API with buffer swapping
- **Numerically Stable**: Log1p/Expm1 optimization for compound expressions
- **FastMath**: Minimax polynomial approximations for `exp` and `log` (`Sin`/`Cos`/`Sqrt` delegate to `math`), plus domain guardrails

## Installation

```bash
go get github.com/emlgo/eml
```

Requires **Go 1.26+**.

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
| `pkg/fastmath` | Minimax polynomial approximations for Exp and Log, domain guardrails |
| `pkg/bytecode` | Bytecode VM (all 27 functions), compiler, optimizer, genetic search for symbolic regression |
| `pkg/quant` | TurboQuant4: 4-bit polar-quantized vector search and distance |
| `internal/eml` | Core EML operator, SIMD dispatch, worker pool, complex batch ops, Pipeline |
| `internal/eml/bigmath` | Arbitrary-precision EML via `math/big.Float`, identity verifier |
| `internal/jit` | JIT compiler, arena allocator, canonical trees, Diff, Simplify, Decompile, Cache |
| `internal/gpu` | CUDA & Metal GPU backends, ULP-based verification |
| `internal/constants` | Mathematical constants (e, π, ln2, √2, φ, etc.) |

## Architecture

```text
emlgo/
├── cmd/
│   ├── bench/           # Benchmark tool
│   ├── validate/        # Validation tool
│   └── emlcli/          # CLI demo (has a `decompile` subcommand)
├── internal/
│   ├── eml/             # Core EML operator + SIMD dispatch
│   │   ├── bigmath/     # Arbitrary-precision backend
│   │   └── simd_*.go    # Platform-specific dispatch (amd64/arm64/wasm)
│   ├── jit/             # JIT compiler + arena + canonical trees + Diff/Simplify/Decompile
│   ├── gpu/             # CUDA & Metal GPU backends
│   └── constants/       # Mathematical constants
├── pkg/
│   ├── arithmetic/      # Basic arithmetic + batch ops
│   ├── bytecode/        # Bytecode VM, compiler, optimizer, genetic programming
│   ├── fastmath/        # Fast polynomial exp/log + guardrails
│   ├── hyper/           # Hyperbolic functions + batch ops
│   ├── logexp/          # Exponential & logarithmic
│   ├── quant/           # TurboQuant4 vector quantization
│   └── trig/            # Trigonometric + batch ops
├── cuda/                # CUDA sources and C API shim
├── metal/               # Metal shaders
├── wasm/                # WebAssembly benchmark page
├── docs/                # Documentation
├── scripts/             # Benchmark & validation scripts
└── .github/workflows/   # CI and release automation
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

## Symbolic Differentiation

```go
import "github.com/emlgo/eml/internal/jit"

// Parse returns (Node, error); Canonicalize turns it into an *EMLNode.
n, _ := jit.Parse("x")
v := jit.Canonicalize(n)

// Build the canonical EML tree for exp(x) = eml(x, 1)
x := jit.CanonicalExp(v)

// Differentiate symbolically
dx := jit.Diff(x)
fmt.Println(jit.Decompile(dx)) // exp(x)

// Evaluate the derivative
result := jit.DiffEval(x, 2.0) // ≈ exp(2)
```

`jit.Parse` returns `(Node, error)`, and `Canonicalize` maps it to the `*EMLNode` form that `Diff`, `Simplify` and `EMLSize` operate on.

## Expression Simplification

```go
import "github.com/emlgo/eml/internal/jit"

// Constant folding: Simplify works on *EMLNode trees, so parse first and
// canonicalise the result.
n, _ := jit.Parse("2 + 3")
node := jit.Simplify(jit.Canonicalize(n))
// Result: constant 5

// Identity reduction
n2, _ := jit.Parse("x + 0")
node2 := jit.Simplify(jit.Canonicalize(n2))
// Result: x
```

## EML Decompiler & LaTeX

```go
import "github.com/emlgo/eml/internal/jit"

n, _ := jit.Parse("sin(x)^2 + cos(x)^2")
emlNode := jit.Canonicalize(n)
fmt.Println(jit.Decompile(emlNode))       // sin(x)^2.0 + cos(x)^2.0
fmt.Println(jit.DecompileLaTeX(emlNode)) // (\sin(x))^{2.0} + (\cos(x))^{2.0}
```

`Decompile` emits infix notation with infix operators (`a + b`, `a * b`, `a^b`); `DecompileLaTeX` emits LaTeX math-mode markup.

## Symbolic Regression

```go
import "github.com/emlgo/eml/pkg/bytecode"

// Fit a closed form to data. Deterministic for a given Seed.
d := bytecode.Dataset{X: [][]float64{x}, Y: y}
res, err := bytecode.Search(bytecode.Config{
    PopulationSize:   400,
    Generations:      400,
    LocalSearchSteps: 60,
    Seed:             1,
}, d, bytecode.LeastSquares(1e-6))

fmt.Printf("RMS error %.3g: %s\n", -res.BestFitness, res.Best)
```

Fitting `3x² + 2x + 1` reaches an RMS error of 4e-13. Transcendental targets are
approximated rather than recovered — fitting `sin(x)` plateaus around RMS 0.17
regardless of budget. See [docs/nextsteps.md](docs/nextsteps.md).

## Composable Pipeline

```go
import "github.com/emlgo/eml/internal/eml"

// Zero-allocation composable pipeline with buffer swapping.
// Each builder call appends a step, so either build the chain once and reuse it
// (as below) or call Reset() before rebuilding.
p := eml.NewPipeline(len(input))
p.Exp().MulScalar(2.0).Log()
p.RunTo(input, output)
```

## JIT Compiler

```go
import "github.com/emlgo/eml/internal/jit"

c := jit.NewCompiler()

// Integer and non-integer exponents
f, _ := c.Compile("x^0.5")    // sqrt(x)
g, _ := c.Compile("x^(2*x)")  // variable exponent

// 25 math functions: sin cos exp log sqrt tan asin acos atan abs cbrt
// log2 log10 ceil floor trunc round sinh cosh tanh asinh acosh atanh erf gamma
h, _ := c.Compile("sin(x)^2 + cos(x)^2") // = 1.0

// Cached compilation
fn, _ := jit.CompileCached("x^2 + 1") // cached after first call
```

## Complex Numbers

```go
import "github.com/emlgo/eml/internal/eml"

z := []complex128{1 + 2i, 3 + 4i, 5 + 6i}
exp := eml.ComplexExpBatch(z)
sin := eml.ComplexSinBatch(z)

// Trigonometric identities hold
// sin²(z) + cos²(z) = 1 verified at complex128 precision
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
import (
    "math/big"

    "github.com/emlgo/eml/internal/eml/bigmath"
)

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

## float32 SIMD

```go
import "github.com/emlgo/eml/internal/eml"

// Dedicated float32 batch operations
x32 := []float32{1.0, 2.0, 3.0, 4.0}
sin32 := eml.SinSIMDF32(x32)
exp32 := eml.ExpSIMDF32(x32)
```

## Building & Testing

```bash
make build                                  # Build all packages
make test                                   # Run tests
make test-race                              # Race detection
make test-cover                             # Coverage report
make bench                                  # Run benchmarks
make lint                                   # golangci-lint
make fuzz                                   # Fuzz testing (30s per target)
make gosec                                  # Security scan
./scripts/bench-compare.sh                  # Benchmark regression
```

## Performance

Ratio is `emlgo time / reference time`, so **lower is faster**.

| Operation | Scalar | Batch (SIMD) |
| :--- | :--- | :--- |
| Add/Sub/Mul | Bit-exact with `math` | 0.75x - 0.84x |
| Exp/Log/Sin/Cos | Bit-exact with `math` | 0.31x - 0.98x |
| PowInt | **0.15x - 0.19x** (5-6x faster than `math.Pow`) | parallelized |
| FastMath Exp / Log | Minimax polynomials, ~1.5e-6 / ~1e-7 relative error | N/A |
| Fused (ExpMul) | N/A | One pass instead of two, **20-30% less memory traffic** |

Measured on the two hosts recorded in [docs/performance.md](docs/performance.md)
(`linux/arm64` NEON and `linux/amd64` AVX2, n=1e6). In those runs the SIMD batch
paths were **slower** than the naive reference loops — the AVX2/AVX-512 kernels are
selected by runtime CPU detection and the generic scalar fallback otherwise
dominates. Integer exponentiation is the clear win. Re-run
`./scripts/bench-compare.sh` on your own hardware before relying on the batch paths.

## License

See LICENSE file.
