# Complex Number Support

emlgo provides first-class `complex128` support, enabling all elementary operations on complex arguments.

## Overview

Complex number operations are implemented in `internal/eml/complex_batch.go` using the `math/cmplx` package. The complex EML operator is defined as:

```math
Complex(x, y) = cmplx.Exp(x) - cmplx.Log(y)
```

When `y = 1`, this simplifies to `cmplx.Exp(x)`.

## Architecture

The library uses a dual-path design:

- **float64 fast path**: All scalar `float64` operations use native EML + SIMD paths (AVX2/AVX512/NEON/WASM).
- **complex128 path**: Uses `math/cmplx` for `Exp`, `Log`, `Sin`, `Cos`, `Tan` on complex arguments.

Batch complex operations (`ComplexBatch`, `ComplexExpBatch`, etc.) process slices of `complex128` in parallel using the worker pool.

## Batch Operations

### ComplexBatch

Applies the complex EML operator to each element pair:

```go
func ComplexBatch(x, y, result []complex128)
```

Computes `result[i] = cmplx.Exp(x[i]) - cmplx.Log(y[i])` for each `i`.

### ComplexExpBatch

Returns the complex exponential of each element:

```go
func ComplexExpBatch(x []complex128) []complex128
```

### ComplexLogBatch

Returns the complex logarithm of each element:

```go
func ComplexLogBatch(x []complex128) []complex128
```

### ComplexSinBatch

Returns the complex sine of each element:

```go
func ComplexSinBatch(x []complex128) []complex128
```

### ComplexCosBatch

Returns the complex cosine of each element:

```go
func ComplexCosBatch(x []complex128) []complex128
```

### ComplexTanBatch

Returns the complex tangent of each element:

```go
func ComplexTanBatch(x []complex128) []complex128
```

## Trigonometric Identity Verification

The fundamental trigonometric identity holds for complex arguments:

```go
// sin²(z) + cos²(z) = 1 for all complex z
z := []complex128{1 + 2i, 3 + 4i, 5 + 6i}
sin := eml.ComplexSinBatch(z)
cos := eml.ComplexCosBatch(z)
for i := range z {
    identity := sin[i]*sin[i] + cos[i]*cos[i]
    // identity ≈ 1+0i for each element
}
```

## Euler's Identity Verification

Euler's identity `e^(iπ) + 1 = 0` is verified at complex128 precision:

```go
z := []complex128{complex(0, math.Pi)}
result := eml.ComplexExpBatch(z)
// result[0] ≈ -1+0i
// result[0] + 1 ≈ 0+0i
```

## Usage Examples

### Basic Complex Operations

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    // Single complex EML
    z := eml.Complex(complex(1, 2), 1) // exp(1+2i)
    fmt.Printf("exp(1+2i) = %v\n", z)

    // Batch exponential
    inputs := []complex128{1 + 0i, 0 + 1i, -1 + 0i, 0 - 1i}
    exps := eml.ComplexExpBatch(inputs)
    for i, v := range inputs {
        fmt.Printf("exp(%v) = %v\n", v, exps[i])
    }
}
```

### Trigonometric Identities

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    z := []complex128{complex(0.5, 1.0), complex(2.0, -3.0)}
    sin := eml.ComplexSinBatch(z)
    cos := eml.ComplexCosBatch(z)

    for i := range z {
        s2 := sin[i] * sin[i]
        c2 := cos[i] * cos[i]
        fmt.Printf("sin²(%v) + cos²(%v) = %v\n", z[i], z[i], s2+c2)
    }

    // Euler's identity
    euler := eml.ComplexExpBatch([]complex128{complex(0, math.Pi)})
    fmt.Printf("e^(iπ) + 1 = %v\n", euler[0]+1)
}
```

### Complex Trigonometric Functions

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    z := []complex128{1 + 1i, 2 - 3i}

    sin := eml.ComplexSinBatch(z)
    cos := eml.ComplexCosBatch(z)
    tan := eml.ComplexTanBatch(z)

    for i, v := range z {
        fmt.Printf("z = %v\n", v)
        fmt.Printf("  sin(z) = %v\n", sin[i])
        fmt.Printf("  cos(z) = %v\n", cos[i])
        fmt.Printf("  tan(z) = %v\n", tan[i])
        fmt.Printf("  sin/cos = %v\n", sin[i]/cos[i])
    }
}
```

## Internal Constants for Complex Numbers

The `internal/constants` package provides complex constants:

| Constant | Value | Description |
|----------|-------|-------------|
| I | 0+1i | Imaginary unit (var) |
| ComplexOne | 1+0i | Complex unit |
| ComplexI | 0+1i | Complex imaginary unit |
| ComplexNegI | 0-1i | Complex negative imaginary unit |

## Notes

- All complex batch operations panic on slice length mismatch.
- Complex operations use `math/cmplx` which handles branch cuts per IEEE 754.
- The worker pool parallelizes complex batch operations for large slices.
