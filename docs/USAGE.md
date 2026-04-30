# Usage Examples

This document provides practical examples of using the emlgo library.

## Basic Usage

### Trigonometric Functions

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/pkg/trig"
)

func main() {
    // Basic trig functions
    fmt.Printf("sin(π/6) = %.6f (expected %.6f)\n", trig.Sin(math.Pi/6), math.Sin(math.Pi/6))
    fmt.Printf("cos(π/3) = %.6f (expected %.6f)\n", trig.Cos(math.Pi/3), math.Cos(math.Pi/3))
    fmt.Printf("tan(π/4) = %.6f (expected %.6f)\n", trig.Tan(math.Pi/4), math.Tan(math.Pi/4))

    // Inverse trig functions
    fmt.Printf("asin(0.5) = %.6f (expected %.6f)\n", trig.Asin(0.5), math.Asin(0.5))
    fmt.Printf("acos(0.5) = %.6f (expected %.6f)\n", trig.Acos(0.5), math.Acos(0.5))
    fmt.Printf("atan(1.0) = %.6f (expected %.6f)\n", trig.Atan(1.0), math.Atan(1.0))

    // Angle conversion
    fmt.Printf("90° in radians = %.6f\n", trig.DegToRad(90))
    fmt.Printf("π rad in degrees = %.6f\n", trig.RadToDeg(math.Pi))
}
```

### Exponential and Logarithmic

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/pkg/logexp"
)

func main() {
    // Exponential
    fmt.Printf("exp(1) = %.6f (expected %.6f)\n", logexp.Exp(1), math.Exp(1))
    fmt.Printf("exp(0) = %.6f\n", logexp.Exp(0))

    // Logarithm
    fmt.Printf("ln(e) = %.6f (expected %.6f)\n", logexp.Log(math.E), math.Log(math.E))
    fmt.Printf("ln(1) = %.6f\n", logexp.Log(1))
}
```

### Arithmetic Operations

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/pkg/arithmetic"
)

func main() {
    // Powers and roots
    fmt.Printf("sqrt(2) = %.6f (expected %.6f)\n", arithmetic.Sqrt(2), math.Sqrt(2))
    fmt.Printf("cbrt(8) = %.6f (expected %.6f)\n", arithmetic.Cbrt(8), math.Cbrt(8))
    fmt.Printf("pow(2, 10) = %.6f (expected %.6f)\n", arithmetic.Pow(2, 10), math.Pow(2, 10))
    fmt.Printf("pow(-2, 3) = %.6f (expected %.6f)\n", arithmetic.Pow(-2, 3), math.Pow(-2, 3))

    // Logarithms with different bases
    fmt.Printf("log2(8) = %.6f\n", arithmetic.LogBase2(8))
    fmt.Printf("log10(100) = %.6f\n", arithmetic.LogBase10(100))
    fmt.Printf("log3(9) = %.6f\n", arithmetic.LogBase(9, 3))

    // Rounding
    fmt.Printf("floor(3.7) = %.6f\n", arithmetic.Floor(3.7))
    fmt.Printf("ceil(3.2) = %.6f\n", arithmetic.Ceil(3.2))
    fmt.Printf("round(3.5) = %.6f\n", arithmetic.Round(3.5))
    fmt.Printf("trunc(-3.7) = %.6f\n", arithmetic.Trunc(-3.7))
}
```

### Hyperbolic Functions

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/pkg/hyper"
    "github.com/emlgo/eml/pkg/trig"
)

func main() {
    // Using hyper package
    fmt.Printf("sinh(1) = %.6f (expected %.6f)\n", hyper.Sinh(1), math.Sinh(1))
    fmt.Printf("cosh(1) = %.6f (expected %.6f)\n", hyper.Cosh(1), math.Cosh(1))
    fmt.Printf("tanh(1) = %.6f (expected %.6f)\n", hyper.Tanh(1), math.Tanh(1))

    // Inverse hyperbolic
    fmt.Printf("asinh(1) = %.6f (expected %.6f)\n", hyper.Asinh(1), math.Asinh(1))
    fmt.Printf("acosh(2) = %.6f (expected %.6f)\n", hyper.Acosh(2), math.Acosh(2))
    fmt.Printf("atanh(0.5) = %.6f (expected %.6f)\n", hyper.Atanh(0.5), math.Atanh(0.5))

    // Also available in trig package
    fmt.Printf("tanh (via trig) = %.6f\n", trig.Tanh(1))
}
```

## Batch Operations

For processing multiple values efficiently, use batch operations:

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/pkg/trig"
)

func main() {
    // Create input slice
    inputs := []float64{0, 0.5, 1.0, 1.5, 2.0}

    // Batch operations
    sinResults := trig.SinBatch(inputs)
    cosResults := trig.CosBatch(inputs)

    // SinCosBatch returns both sin and cos in one pass
    sin, cos := trig.SinCosBatch(inputs)

    fmt.Println("SinBatch results:", sinResults)
    fmt.Println("CosBatch results:", cosResults)
    fmt.Println("SinCos results:")
    for i, x := range inputs {
        fmt.Printf("  x=%.1f: sin=%.6f, cos=%.6f\n", x, sin[i], cos[i])
    }

    // TanBatch
    tanResults := trig.TanBatch(inputs)
    fmt.Println("TanBatch results:", tanResults)
}
```

## Error Handling

All functions handle edge cases correctly:

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/pkg/trig"
    "github.com/emlgo/eml/pkg/arithmetic"
)

func main() {
    // NaN handling
    fmt.Printf("sin(NaN) = %v\n", trig.Sin(math.NaN()))
    fmt.Printf("logexp.Log(-1) = %v\n", math.NaN()) // Returns NaN for invalid input

    // Infinity handling
    fmt.Printf("sin(+Inf) = %v\n", trig.Sin(math.Inf(1)))
    fmt.Printf("cos(-Inf) = %v\n", trig.Cos(math.Inf(-1)))

    // Domain errors
    fmt.Printf("asin(2) = %v (outside domain)\n", trig.Asin(2))
    fmt.Printf("acosh(0.5) = %v (outside domain)\n", trig.Acosh(0.5))
    fmt.Printf("logexp.Log(-1) = %v (invalid input)\n", math.NaN())

    // Division by zero
    fmt.Printf("Div(1, 0) = %v\n", arithmetic.Div(1, 0))
    fmt.Printf("Sin(π/2) = %v (for sec at boundary)\n", trig.Sec(math.Pi/2))
}
```

## SIMD Detection

The library automatically detects SIMD capabilities:

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    fmt.Printf("Has AVX2: %v\n", eml.HasAVX2())
    fmt.Printf("Has AVX-512: %v\n", eml.HasAVX512())
    fmt.Printf("Has NEON: %v\n", eml.HasNeon())
}
```

## Performance Considerations

- Single function calls have similar overhead to math library
- Batch operations (SinBatch, CosBatch, etc.) are optimized for large slices
- SIMD chunk size: 4 for AVX2/NEON, 8 for AVX-512
- For small slices (<8 elements), scalar implementation may be faster