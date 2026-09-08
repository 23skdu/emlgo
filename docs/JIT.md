# JIT Compiler

emlgo includes a JIT (Just-In-Time) compiler that generates x86-64 machine code from mathematical expression strings.

## Overview

The JIT compiler (`internal/jit`) parses a mathematical expression string, builds an AST, and emits native x86-64 SSE2 machine code that can be called as a regular Go function. This provides near-native performance for compiled expressions without interpretation overhead.

**Platform**: AMD64 only (non-amd64 builds use the interpreter fallback in `arena.go`).

## Supported Expressions

### Operators

| Operator | Example | Description |
|----------|---------|-------------|
| `+` | `x + y` | Addition |
| `-` | `x - y` | Subtraction |
| `*` | `x * y` | Multiplication |
| `/` | `x / y` | Division |
| `^` | `x ^ n` | Power (see Exponent Support below) |
| unary `-` | `-x` | Negation |

### Functions (16 total)

| Function | Example | Description |
|----------|---------|-------------|
| `sin` | `sin(x)` | Sine |
| `cos` | `cos(x)` | Cosine |
| `tan` | `tan(x)` | Tangent |
| `exp` | `exp(x)` | Exponential e^x |
| `log` | `log(x)` | Natural logarithm ln(x) |
| `sqrt` | `sqrt(x)` | Square root |
| `asin` | `asin(x)` | Arcsine |
| `acos` | `acos(x)` | Arccosine |
| `atan` | `atan(x)` | Arctangent |
| `abs` | `abs(x)` | Absolute value |
| `cbrt` | `cbrt(x)` | Cube root |
| `log2` | `log2(x)` | Log base 2 |
| `log10` | `log10(x)` | Log base 10 |
| `ceil` | `ceil(x)` | Ceiling |
| `floor` | `floor(x)` | Floor |
| `trunc` | `trunc(x)` | Truncate |

## Exponent Support

The JIT compiler supports three exponent strategies:

### Integer Exponents (Binary Exponentiation)

For integer constant exponents, the compiler uses binary exponentiation (O(log n) multiplications):

```go
c := jit.NewCompiler()
f, _ := c.Compile("x^3")
// Emits: base * base * base using squaring
```

Negative integer exponents compute `1.0 / x^|n|`:

```go
g, _ := c.Compile("x^-2")
// Emits: 1.0 / (x * x)
```

### Non-Integer Constant Exponents

Non-integer constant exponents (e.g., `x^0.5`, `x^1.5`) are compiled as `exp(y * log(x))`:

```go
h, _ := c.Compile("x^0.5")   // sqrt(x)
k, _ := c.Compile("x^1.5")   // x * sqrt(x)
l, _ := c.Compile("x^0.333") // approximate cbrt(x)
```

### Variable Exponents

Variable exponents (e.g., `x^x`, `x^(x-1)`) are compiled as `exp(y * log(x))` where `y` is evaluated at runtime:

```go
m, _ := c.Compile("x^x")     // exp(x * log(x))
n, _ := c.Compile("x^(x-1)") // exp((x-1) * log(x))
```

## Usage Examples

### Basic Compilation

```go
package main

import (
    "fmt"
    "github.com/emlgo/eml/internal/jit"
)

func main() {
    c := jit.NewCompiler()

    // Compile and evaluate
    f, err := c.Compile("sin(x)^2 + cos(x)^2")
    if err != nil {
        panic(err)
    }

    result := f(1.23) // Always 1.0
    fmt.Printf("sin²(1.23) + cos²(1.23) = %.6f\n", result)
}
```

### Non-Integer Exponents

```go
c := jit.NewCompiler()

sqrt, _ := c.Compile("x^0.5")
fmt.Printf("sqrt(4) = %.6f\n", sqrt(4))   // 2.0
fmt.Printf("sqrt(9) = %.6f\n", sqrt(9))   // 3.0

pow15, _ := c.Compile("x^1.5")
fmt.Printf("4^1.5 = %.6f\n", pow15(4))    // 8.0
```

### Variable Exponents

```go
c := jit.NewCompiler()

xx, _ := c.Compile("x^x")
fmt.Printf("2^2 = %.6f\n", xx(2))     // 4.0
fmt.Printf("3^3 = %.6f\n", xx(3))     // 27.0

xpow, _ := c.Compile("x^(x-1)")
fmt.Printf("3^2 = %.6f\n", xpow(3))   // 9.0
```

### Complex Expressions

```go
c := jit.NewCompiler()

// Polynomial
poly, _ := c.Compile("x^3 - 2*x^2 + x - 1")

// Nested functions
nested, _ := c.Compile("sqrt(sin(x)^2 + cos(x)^2)")

// Multiple operations
expr, _ := c.Compile("(x + 1) / (x - 1)")

// All 16 functions work
allFuncs, _ := c.Compile("log2(abs(x)) + log10(abs(x))")
```

## Arena Allocator

The arena allocator provides zero-allocation paths for expression parsing and evaluation:

### Creating an Arena

```go
arena := jit.NewArena(64) // Initial capacity of 64 nodes
```

### Allocating Nodes

```go
node := arena.Alloc() // Bump-pointer allocation, grows if needed
```

### Converting Between Node Types

```go
// Convert Node interface to arena node
arenaNode := arena.FromInterface(ast)

// Convert arena node back to Node interface
nodeInterface := arena.ToInterface(arenaNode)
```

### Evaluating Arena Trees

```go
result := jit.EvalArena(arenaNode, 3.14)
```

### Resetting

```go
arena.Reset() // All previously allocated nodes become available for reuse
```

### Node Size

Each `ArenaNode` is exactly 86 bytes:
- `Kind` (1 byte)
- `Op` (8 bytes)
- `Value` (8 bytes)
- `Name` (16 bytes)
- `Left`, `Right` (16 bytes each)
- Padding (7 bytes)

## Canonical EML Trees

The canonical EML tree representation normalizes any expression to its minimal EML form:

### Canonical Constructors

```go
// exp(x) = eml(x, 1)
expNode := jit.CanonicalExp(x)

// log(x) = eml(1, eml(eml(1, x), 1))
logNode := jit.CanonicalLog(x)

// sin(x) (function node)
sinNode := jit.CanonicalSin(x)

// cos(x) (function node)
cosNode := jit.CanonicalCos(x)

// sqrt(x) = exp(0.5 * log(x))
sqrtNode := jit.CanonicalSqrt(x)
```

### Canonicalizing Any Expression

```go
c := jit.NewCompiler()
ast, _ := c.Parse("exp(x)")
canonical := jit.Canonicalize(ast)
// canonical is eml(x, 1)
```

### Structural Equivalence

```go
ast1, _ := c.Parse("exp(x)")
ast2, _ := c.Parse("exp(x)")
can1 := jit.Canonicalize(ast1)
can2 := jit.Canonicalize(ast2)
equiv := jit.Equiv(can1, can2) // true
```

### Tree Size

```go
size := jit.EMLSize(canonical) // Number of nodes in the tree
```

## Register Allocation

The JIT compiler uses 15 XMM registers:
- **xmm0**: Return register (also holds initial input x)
- **xmm1-xmm14**: General purpose temporaries
- **xmm15**: Reserved for the input variable x

Expressions that require more than 14 live temporaries will return an error.

## Function Call Mechanism

JIT-compiled functions are called through indirect calls using ABIInternal calling convention:
1. xmm0-xmm7 are saved to stack
2. Argument is moved to xmm0
3. Function pointer is loaded and called via `CALL R8`
4. xmm1-xmm7 are restored (xmm0 holds return value)
5. Return value is moved to the destination register

## Performance

- **Integer powers**: Binary exponentiation, O(log n) multiplications
- **Non-integer powers**: `exp(y * log(x))` — single `log` + `exp` call
- **Function calls**: Indirect call overhead (~2-3ns per call)
- **Overall**: Typically 1.5-3x faster than interpreting the AST

## Limitations

- Only available on AMD64 (uses SSE2 instruction set)
- Maximum 15 live temporaries (expressions with more will fail to compile)
- Function names limited to 15 characters
- No support for user-defined functions
