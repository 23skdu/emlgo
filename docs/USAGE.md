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
)

func main() {
    fmt.Printf("sinh(1) = %.6f (expected %.6f)\n", hyper.Sinh(1), math.Sinh(1))
    fmt.Printf("cosh(1) = %.6f (expected %.6f)\n", hyper.Cosh(1), math.Cosh(1))
    fmt.Printf("tanh(1) = %.6f (expected %.6f)\n", hyper.Tanh(1), math.Tanh(1))

    // Inverse hyperbolic
    fmt.Printf("asinh(1) = %.6f (expected %.6f)\n", hyper.Asinh(1), math.Asinh(1))
    fmt.Printf("acosh(2) = %.6f (expected %.6f)\n", hyper.Acosh(2), math.Acosh(2))
    fmt.Printf("atanh(0.5) = %.6f (expected %.6f)\n", hyper.Atanh(0.5), math.Atanh(0.5))
}
```

## Complex Number Operations

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    // Single complex operations
    z := complex(1.0, 2.0)
    result := eml.Complex(z, 1) // cmplx.Exp(z) - cmplx.Log(1) = cmplx.Exp(z)
    fmt.Printf("eml(1+2i, 1) = %v\n", result)

    // Batch complex operations
    inputs := []complex128{
        complex(1, 0),
        complex(0, 1),   // i
        complex(-1, 0),  // -1
        complex(0, -1),  // -i
    }

    expResults := eml.ComplexExpBatch(inputs)
    fmt.Println("Exp batch results:")
    for i, z := range inputs {
        fmt.Printf("  exp(%v) = %v\n", z, expResults[i])
    }

    // Trigonometric identity: sin²(z) + cos²(z) = 1
    sinResults := eml.ComplexSinBatch(inputs)
    cosResults := eml.ComplexCosBatch(inputs)
    fmt.Println("\nTrig identity sin²(z) + cos²(z) = 1:")
    for i, z := range inputs {
        s := sinResults[i]
        c := cosResults[i]
        identity := s*s + c*c
        fmt.Printf("  sin²(%v) + cos²(%v) = %v\n", z, z, identity)
    }

    // Euler's identity: e^(iπ) + 1 = 0
    euler := eml.ComplexExpBatch([]complex128{complex(0, math.Pi)})
    fmt.Printf("\ne^(iπ) = %v (should be ≈ -1)\n", euler[0])
    fmt.Printf("e^(iπ) + 1 = %v (should be ≈ 0)\n", euler[0]+1)
}
```

## JIT Compilation with Non-Integer and Variable Exponents

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    c := jit.NewCompiler()

    // Integer exponent (binary exponentiation)
    f, _ := c.Compile("x^3")
    fmt.Printf("x^3 at x=2: %.6f\n", f(2)) // 8.0

    // Negative integer exponent
    g, _ := c.Compile("x^-2")
    fmt.Printf("x^-2 at x=4: %.6f\n", g(4)) // 0.0625

    // Non-integer constant exponent (exp(y * log(x)))
    h, _ := c.Compile("x^0.5")
    fmt.Printf("x^0.5 at x=4: %.6f\n", h(4)) // 2.0

    // Non-integer constant exponent
    k, _ := c.Compile("x^1.5")
    fmt.Printf("x^1.5 at x=4: %.6f\n", k(4)) // 8.0

    // Variable exponent (exp(y * log(x)))
    l, _ := c.Compile("x^x")
    fmt.Printf("x^x at x=2: %.6f\n", l(2)) // 4.0
    fmt.Printf("x^x at x=3: %.6f\n", l(3)) // 27.0

    // Variable exponent expression
    m, _ := c.Compile("x^(x-1)")
    fmt.Printf("x^(x-1) at x=3: %.6f\n", m(3)) // 9.0

    // Complex expressions with 17 supported functions
    n, _ := c.Compile("sin(x)^2 + cos(x)^2")
    fmt.Printf("sin²(x) + cos²(x) = %.6f\n", n(1.23)) // 1.0

    // All 17 functions work
    o, _ := c.Compile("sqrt(abs(x)) + cbrt(x)")
    fmt.Printf("sqrt(|-8|) + cbrt(-8) = %.6f\n", o(-8)) // 2.0 + (-2.0) = 0.0

    // round function
    r, _ := c.Compile("round(x)")
    fmt.Printf("round(3.7) = %.6f\n", r(3.7)) // 4.0
}
```

## JIT Expression Cache

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    // First call compiles and caches
    fn, err := jit.CompileCached("sin(x)^2 + cos(x)^2")
    if err != nil {
        panic(err)
    }
    fmt.Printf("f(1.23) = %.6f\n", fn(1.23)) // 1.0

    // Subsequent calls return cached result (no recompilation)
    fn2, _ := jit.CompileCached("sin(x)^2 + cos(x)^2")
    fmt.Printf("Same function: %v\n", fn == fn2) // true

    // Clear cache for long-running programs
    jit.ClearJITCache()
}
```

## Symbolic Differentiation

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    // Parse and canonicalize
    node := jit.Canonicalize(jit.Parse("x^2"))

    // Differentiate: d/dx x² = 2x
    dx := jit.Diff(node)
    fmt.Printf("d/dx x² = %s\n", jit.Decompile(dx)) // mul(2.0, x)

    // Evaluate derivative at a point
    result := jit.DiffEval(node, 3.0) // 2*3 = 6.0
    fmt.Printf("d/dx x² at x=3: %.6f\n", result)

    // More complex: d/dx sin(x) = cos(x)
    sinNode := jit.Canonicalize(jit.Parse("sin(x")))
    dsin := jit.Diff(sinNode)
    fmt.Printf("d/dx sin(x) = %s\n", jit.Decompile(dsin)) // cos(x)
}
```

## Expression Simplification

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    // Constant folding
    node := jit.Simplify(jit.Canonicalize(jit.Parse("2 + 3")))
    fmt.Printf("2 + 3 = %s\n", jit.Decompile(node)) // 5.0

    // Identity reduction
    node2 := jit.Simplify(jit.Canonicalize(jit.Parse("x + 0")))
    fmt.Printf("x + 0 = %s\n", jit.Decompile(node2)) // x

    // Double negation
    node3 := jit.Simplify(jit.Canonicalize(jit.Parse("-(-x)")))
    fmt.Printf("-(-x) = %s\n", jit.Decompile(node3)) // x

    // Multiplication identities
    node4 := jit.Simplify(jit.Canonicalize(jit.Parse("x * 1")))
    fmt.Printf("x * 1 = %s\n", jit.Decompile(node4)) // x

    node5 := jit.Simplify(jit.Canonicalize(jit.Parse("x * 0")))
    fmt.Printf("x * 0 = %s\n", jit.Decompile(node5)) // 0.0
}
```

## EML Decompiler & LaTeX

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    // Parse and canonicalize
    emlNode := jit.Canonicalize(jit.Parse("sin(x)^2 + cos(x)^2"))

    // Infix notation
    fmt.Println("Infix:", jit.Decompile(emlNode))

    // LaTeX math mode
    fmt.Println("LaTeX:", jit.DecompileLaTeX(emlNode))

    // CLI usage:
    // emlcli decompile "sin(x)^2 + cos(x)^2"
}
```

## Composable Pipeline

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    input := []float64{0.5, 1.0, 1.5, 2.0}
    output := make([]float64, len(input))

    // Zero-allocation composable pipeline with buffer swapping
    p := eml.NewPipeline(len(input))
    p.Exp().MulScalar(2.0).Log().RunTo(input, output)

    fmt.Printf("input:  %v\n", input)
    fmt.Printf("output: %v\n", output)
    // output[i] = log(2 * exp(input[i]))
}
```

## Arena Zero-Allocation Parsing

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    // Create arena with initial capacity
    arena := jit.NewArena(64)

    // Parse an expression and convert to arena nodes
    c := jit.NewCompiler()
    ast, _ := c.Parse("sin(x)^2 + cos(x)^2")

    // Convert AST to arena-allocated nodes (zero heap allocation)
    arenaNode := arena.FromInterface(ast)
    fmt.Printf("Arena nodes allocated: %d\n", arena.NumNodes())

    // Evaluate without touching the heap
    result := jit.EvalArena(arenaNode, 3.14)
    fmt.Printf("Result: %.6f\n", result) // 1.0

    // Reset and reuse
    arena.Reset()
    fmt.Printf("After reset: %d nodes\n", arena.NumNodes()) // 0
}
```

## Canonical EML Tree Construction

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    c := jit.NewCompiler()

    // Parse and canonicalize
    ast, _ := c.Parse("exp(x)")
    canonical := jit.Canonicalize(ast)
    fmt.Printf("exp(x) canonical: %v\n", canonical)
    // eml(x, 1)

    // All exp/log/sqrt map to canonical EML form
    ast2, _ := c.Parse("log(x)")
    canonical2 := jit.Canonicalize(ast2)
    fmt.Printf("log(x) canonical: %v\n", canonical2)
    // eml(1, eml(eml(1, x), 1))

    // Structural equivalence
    ast3, _ := c.Parse("exp(x)")
    canonical3 := jit.Canonicalize(ast3)
    fmt.Printf("exp(x) ≡ exp(x): %v\n", jit.Equiv(canonical, canonical3)) // true

    // Tree size
    fmt.Printf("Canonical exp(x) size: %d nodes\n", jit.EMLSize(canonical)) // 2
    fmt.Printf("Canonical log(x) size: %d nodes\n", jit.EMLSize(canonical2)) // 5

    // Canonical sqrt(x) = exp(0.5 * log(x))
    ast4, _ := c.Parse("sqrt(x)")
    canonical4 := jit.Canonicalize(ast4)
    fmt.Printf("sqrt(x) canonical: %v\n", canonical4)
    fmt.Printf("sqrt(x) size: %d nodes\n", jit.EMLSize(canonical4))
}
```

## Arbitrary Precision Verification

```go
package main

import (
    "fmt"
    "math/big"
    "github.com/emlgo/eml/internal/eml/bigmath"
)

func main() {
    // Basic arbitrary-precision EML
    x := bigmath.NewFloat(1.0)
    y := bigmath.NewFloat(1.0)
    result := bigmath.Eml(x, y) // exp(1) - log(1) = e
    fmt.Printf("eml(1, 1) = %.6f\n", bigmath.Float64(result)) // 2.718282

    // Verify sin²(x) + cos²(x) = 1
    v := bigmath.DefaultVerifier()
    sin2pluscos2 := func(x *big.Float) *big.Float {
        s := bigmath.Sin(x)
        c := bigmath.Cos(x)
        ss := new(big.Float).Mul(s, s)
        cc := new(big.Float).Mul(c, c)
        return new(big.Float).Add(ss, cc)
    }
    one := func(x *big.Float) *big.Float { return bigmath.NewFloat(1.0) }

    identityHolds := v.VerifyIdentity(sin2pluscos2, one)
    fmt.Printf("sin²(x) + cos²(x) = 1 verified: %v\n", identityHolds)

    // Evaluate at specific point
    f := func(x *big.Float) *big.Float { return bigmath.Exp(x) }
    val := bigmath.EvalAt(f, 1.0)
    fmt.Printf("exp(1) at high precision: %.10f\n", val) // 2.7182818285
}
```

## Numerically Stable Pow Near x=1

```go
package main

import (
    "fmt"
    "math"
    "github.com/emlgo/eml/pkg/arithmetic"
)

func main() {
    // Near x=1, standard Log suffers from catastrophic cancellation
    x := 1.0 + 1e-15
    y := 1e6

    // emlgo uses Log1p for stability
    emlResult := arithmetic.Pow(x, y)
    mathResult := math.Pow(x, y)

    fmt.Printf("emlgo Pow(1+1e-15, 1e6): %.12f\n", emlResult)
    fmt.Printf("math   Pow(1+1e-15, 1e6): %.12f\n", mathResult)
    fmt.Printf("Difference: %.2e\n", math.Abs(emlResult-mathResult))

    // Exp near x=0 uses Expm1
    smallX := 1e-15
    fmt.Printf("\nemlgo Exp(1e-15): %.25f\n", arithmetic.Exp(smallX))
    fmt.Printf("math   Exp(1e-15): %.25f\n", math.Exp(smallX))
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

    // Batch operations (SIMD-accelerated)
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

## float32 SIMD Batch Operations

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    // float32 batch operations for memory-constrained workloads
    x := []float32{0.5, 1.0, 1.5, 2.0}
    sin := eml.SinSIMDF32(x)
    exp := eml.ExpSIMDF32(x)
    fmt.Printf("SinSIMDF32: %v\n", sin)
    fmt.Printf("ExpSIMDF32: %v\n", exp)
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

    // Infinity handling
    fmt.Printf("sin(+Inf) = %v\n", trig.Sin(math.Inf(1)))
    fmt.Printf("cos(-Inf) = %v\n", trig.Cos(math.Inf(-1)))

    // Domain errors
    fmt.Printf("asin(2) = %v (outside domain)\n", trig.Asin(2))
    fmt.Printf("acosh(0.5) = %v (outside domain)\n", trig.Acosh(0.5))

    // Division by zero
    fmt.Printf("Div(1, 0) = %v\n", arithmetic.Div(1, 0))
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
    fmt.Printf("Has FMA: %v\n", eml.HasFMA())
    fmt.Printf("Has WASM SIMD: %v\n", eml.HasWasmSIMD())
}
```

## Worker Pool Shutdown

```go
package main

import (
    "github.com/emlgo/eml/internal/eml"
)

func main() {
    // Use batch operations (workers are auto-started)...

    // When done, gracefully shut down the worker pool
    eml.StopWorkerPool()
}
```

## Performance Considerations

- Single function calls have similar overhead to math library
- Batch operations (SinBatch, CosBatch, etc.) are optimized for large slices
- SIMD chunk size: 4 for AVX2/NEON, 8 for AVX-512
- For small slices (<8 elements), scalar implementation may be faster
- Arena allocator eliminates GC pressure in JIT hot paths
- Canonical EML trees enable expression-level optimization and comparison
- Pipeline API provides zero-allocation chained operations
- JIT expression cache avoids recompilation for repeated expressions
