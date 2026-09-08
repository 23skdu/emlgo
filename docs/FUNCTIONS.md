# Function Reference

Complete list of all public functions in the emlgo library.

---

## Package: pkg/logexp

Exponential and logarithmic functions using EML operator.

| Function | Signature | Description |
|----------|-----------|-------------|
| Exp | `func Exp(x float64) float64` | Exponential function e^x |
| Log | `func Log(x float64) float64` | Natural logarithm ln(x) |
| ExpBatch | `func ExpBatch(x []float64) []float64` | Batch exponential (SIMD) |
| LogBatch | `func LogBatch(x []float64) []float64` | Batch logarithm (SIMD) |
| ExpFast | `func ExpFast(x float64) float64` | Fast approximate exp |
| LogFast | `func LogFast(x float64) float64` | Fast approximate log |

---

## Package: pkg/trig

Trigonometric and inverse trigonometric functions.

### Trigonometric Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| Sin | `func Sin(x float64) float64` | Sine of x (radians) |
| Cos | `func Cos(x float64) float64` | Cosine of x (radians) |
| Tan | `func Tan(x float64) float64` | Tangent of x (radians) |
| Cot | `func Cot(x float64) float64` | Cotangent of x (radians) |
| Sec | `func Sec(x float64) float64` | Secant of x (radians) |
| Csc | `func Csc(x float64) float64` | Cosecant of x (radians) |

### Inverse Trigonometric Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| Asin | `func Asin(x float64) float64` | Arcsine, inverse of sin |
| Acos | `func Acos(x float64) float64` | Arccosine, inverse of cos |
| Atan | `func Atan(x float64) float64` | Arctangent, inverse of tan |
| Atan2 | `func Atan2(y, x float64) float64` | Two-argument arctangent |
| Acot | `func Acot(x float64) float64` | Arccotangent |
| Asec | `func Asec(x float64) float64` | Arcsecant |
| Acsc | `func Acsc(x float64) float64` | Arccosecant |

### Hyperbolic Functions (also in trig)

| Function | Signature | Description |
|----------|-----------|-------------|
| Sinh | `func Sinh(x float64) float64` | Hyperbolic sine |
| Cosh | `func Cosh(x float64) float64` | Hyperbolic cosine |
| Tanh | `func Tanh(x float64) float64` | Hyperbolic tangent |
| Coth | `func Coth(x float64) float64` | Hyperbolic cotangent |
| Sech | `func Sech(x float64) float64` | Hyperbolic secant |
| Csch | `func Csch(x float64) float64` | Hyperbolic cosecant |

### Inverse Hyperbolic Functions (also in trig)

| Function | Signature | Description |
|----------|-----------|-------------|
| Asinh | `func Asinh(x float64) float64` | Inverse hyperbolic sine |
| Acosh | `func Acosh(x float64) float64` | Inverse hyperbolic cosine |
| Atanh | `func Atanh(x float64) float64` | Inverse hyperbolic tangent |
| Acoth | `func Acoth(x float64) float64` | Inverse hyperbolic cotangent |
| Asech | `func Asech(x float64) float64` | Inverse hyperbolic secant |
| Acsch | `func Acsch(x float64) float64` | Inverse hyperbolic cosecant |

### Utility Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| SinCos | `func SinCos(x float64) (sin, cos float64)` | Simultaneous sin and cos |
| SinhCosh | `func SinhCosh(x float64) (sinh, cosh float64)` | Simultaneous sinh and cosh |
| DegToRad | `func DegToRad(deg float64) float64` | Degrees to radians |
| RadToDeg | `func RadToDeg(rad float64) float64` | Radians to degrees |

### Batch Operations (SIMD Optimized)

| Function | Signature | Description |
|----------|-----------|-------------|
| SinBatch | `func SinBatch(x []float64) []float64` | Batch sin with SIMD |
| CosBatch | `func CosBatch(x []float64) []float64` | Batch cos with SIMD |
| SinCosBatch | `func SinCosBatch(x []float64) (sin, cos []float64)` | Batch sin/cos with SIMD |
| TanBatch | `func TanBatch(x []float64) []float64` | Batch tan with SIMD |

---

## Package: pkg/hyper

Hyperbolic functions (dedicated package).

| Function | Signature | Description |
|----------|-----------|-------------|
| Sinh | `func Sinh(x float64) float64` | Hyperbolic sine |
| Cosh | `func Cosh(x float64) float64` | Hyperbolic cosine |
| Tanh | `func Tanh(x float64) float64` | Hyperbolic tangent |
| Asinh | `func Asinh(x float64) float64` | Inverse hyperbolic sine |
| Acosh | `func Acosh(x float64) float64` | Inverse hyperbolic cosine |
| Atanh | `func Atanh(x float64) float64` | Inverse hyperbolic tangent |

---

## Package: pkg/arithmetic

Basic arithmetic operations, roots, and powers.

### Basic Arithmetic

| Function | Signature | Description |
|----------|-----------|-------------|
| Add | `func Add(x, y float64) float64` | Addition x + y |
| Sub | `func Sub(x, y float64) float64` | Subtraction x - y |
| Mul | `func Mul(x, y float64) float64` | Multiplication x * y |
| Div | `func Div(x, y float64) float64` | Division x / y |
| Mod | `func Mod(x, y float64) float64` | Modulo x % y |
| Remainder | `func Remainder(x, y float64) float64` | IEEE remainder |

### Powers and Roots

| Function | Signature | Description |
|----------|-----------|-------------|
| Pow | `func Pow(x, y float64) float64` | Power x^y (uses Log1p near x=1) |
| PowInt | `func PowInt(x float64, n int) float64` | Integer power x^n |
| Sqrt | `func Sqrt(x float64) float64` | Square root |
| Cbrt | `func Cbrt(x float64) float64` | Cube root |
| Hypot | `func Hypot(x, y float64) float64` | sqrt(x² + y²) |

### Logarithms

| Function | Signature | Description |
|----------|-----------|-------------|
| Log | `func Log(x float64) float64` | Natural logarithm ln(x) |
| Log1p | `func Log1p(x float64) float64` | ln(1+x), for small x |
| LogBase | `func LogBase(x, base float64) float64` | Log base b: log_b(x) |
| LogBase2 | `func LogBase2(x float64) float64` | Log base 2: log_2(x) |
| LogBase10 | `func LogBase10(x float64) float64` | Log base 10: log_10(x) |

### Exponential

| Function | Signature | Description |
|----------|-----------|-------------|
| Exp | `func Exp(x float64) float64` | Exponential e^x |
| ExpM1 | `func ExpM1(x float64) float64` | e^x - 1, for small x |

### Comparison

| Function | Signature | Description |
|----------|-----------|-------------|
| Max | `func Max(x, y float64) float64` | Maximum of x and y (NaN-aware) |
| Min | `func Min(x, y float64) float64` | Minimum of x and y (NaN-aware) |

### Rounding

| Function | Signature | Description |
|----------|-----------|-------------|
| Floor | `func Floor(x float64) float64` | Floor, greatest integer ≤ x |
| Ceil | `func Ceil(x float64) float64` | Ceiling, smallest integer ≥ x |
| Trunc | `func Trunc(x float64) float64` | Truncate decimal part |
| Round | `func Round(x float64) float64` | Round to nearest integer |

### Unary Operations

| Function | Signature | Description |
|----------|-----------|-------------|
| Abs | `func Abs(x float64) float64` | Absolute value |
| Neg | `func Neg(x float64) float64` | Negation -x |
| Inv | `func Inv(x float64) float64` | Reciprocal 1/x |
| Square | `func Square(x float64) float64` | Square x² |
| Cube | `func Cube(x float64) float64` | Cube x³ |

### Special

| Function | Signature | Description |
|----------|-----------|-------------|
| FMA | `func FMA(x, y, z float64) float64` | Fused multiply-add: x*y+z |
| GCD | `func GCD(a, b int64) int64` | Greatest common divisor |
| LCM | `func LCM(a, b int64) int64` | Least common multiple |

---

## Package: internal/eml

Core EML operator and SIMD utilities (internal).

| Function | Signature | Description |
|----------|-----------|-------------|
| Eml | `func Eml(x, y float64) float64` | Core EML operator: exp(x) - ln(y) |
| EmlOne | `func EmlOne(x float64) float64` | Eml(x, 1) = exp(x) |
| OneEml | `func OneEml(y float64) float64` | Eml(1, y) = e - ln(y) |
| Complex | `func Complex(x, y complex128) complex128` | Complex EML: cmplx.Exp(x) - cmplx.Log(y) |
| ComplexOne | `func ComplexOne(x complex128) complex128` | Complex EML(x, 1) = cmplx.Exp(x) |
| ComplexBatch | `func ComplexBatch(x, y, result []complex128)` | Batch complex EML operation |
| ComplexExpBatch | `func ComplexExpBatch(x []complex128) []complex128` | Batch complex exponential |
| ComplexLogBatch | `func ComplexLogBatch(x []complex128) []complex128` | Batch complex logarithm |
| ComplexSinBatch | `func ComplexSinBatch(x []complex128) []complex128` | Batch complex sine |
| ComplexCosBatch | `func ComplexCosBatch(x []complex128) []complex128` | Batch complex cosine |
| ComplexTanBatch | `func ComplexTanBatch(x []complex128) []complex128` | Batch complex tangent |
| StopWorkerPool | `func StopWorkerPool()` | Gracefully shut down worker pool |
| HasAVX2 | `func HasAVX2() bool` | AVX2 detection |
| HasAVX512 | `func HasAVX512() bool` | AVX-512 detection |
| HasNeon | `func HasNeon() bool` | ARM NEON detection |

---

## Package: internal/eml/bigmath

Arbitrary-precision EML operations using `math/big.Float`.

| Function | Signature | Description |
|----------|-----------|-------------|
| Eml | `func Eml(x, y *big.Float) *big.Float` | exp(x) - log(y) at arbitrary precision |
| Exp | `func Exp(x *big.Float) *big.Float` | Arbitrary-precision exponential (Taylor series) |
| Log | `func Log(x *big.Float) *big.Float` | Arbitrary-precision natural log |
| Sin | `func Sin(x *big.Float) *big.Float` | Arbitrary-precision sine |
| Cos | `func Cos(x *big.Float) *big.Float` | Arbitrary-precision cosine |
| Sqrt | `func Sqrt(x *big.Float) *big.Float` | Arbitrary-precision square root |
| NewFloat | `func NewFloat(x float64) *big.Float` | Create big.Float from float64 (256-bit) |
| Float64 | `func Float64(x *big.Float) float64` | Convert big.Float to float64 |
| DefaultVerifier | `func DefaultVerifier() *IdentityVerifier` | Identity verifier with default settings |
| VerifyIdentity | `func (v *IdentityVerifier) VerifyIdentity(expr1, expr2 func(*big.Float) *big.Float) bool` | Test symbolic identity at high precision |

---

## Package: internal/jit

JIT compiler for math expressions (amd64 only).

### Compiler Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| NewCompiler | `func NewCompiler() *Compiler` | Create a new JIT compiler |
| Compile | `func (c *Compiler) Compile(expr string) (jitFunc, error)` | Compile expression string to native function |

### Supported Operators

| Operator | Example | Description |
|----------|---------|-------------|
| `+` | `x + y` | Addition |
| `-` | `x - y` | Subtraction |
| `*` | `x * y` | Multiplication |
| `/` | `x / y` | Division |
| `^` | `x ^ 3` | Power (integer: binary exponentiation, non-integer/variable: exp(y*log(x))) |
| unary `-` | `-x` | Negation |

### Supported Functions (JIT)

| Function | Example | Description |
|----------|---------|-------------|
| sin | `sin(x)` | Sine |
| cos | `cos(x)` | Cosine |
| tan | `tan(x)` | Tangent |
| exp | `exp(x)` | Exponential e^x |
| log | `log(x)` | Natural logarithm ln(x) |
| sqrt | `sqrt(x)` | Square root |
| asin | `asin(x)` | Arcsine |
| acos | `acos(x)` | Arccosine |
| atan | `atan(x)` | Arctangent |
| abs | `abs(x)` | Absolute value |
| cbrt | `cbrt(x)` | Cube root |
| log2 | `log2(x)` | Log base 2 |
| log10 | `log10(x)` | Log base 10 |
| ceil | `ceil(x)` | Ceiling |
| floor | `floor(x)` | Floor |
| trunc | `trunc(x)` | Truncate |

### Exponent Support

| Expression | Strategy | Description |
|------------|----------|-------------|
| `x^3` | Binary exponentiation | O(log n) integer power |
| `x^-2` | Reciprocal | `1.0 / x^2` |
| `x^0.5` | `exp(0.5 * log(x))` | Non-integer constant exponent |
| `x^1.5` | `exp(1.5 * log(x))` | Non-integer constant exponent |
| `x^x` | `exp(x * log(x))` | Variable exponent |
| `x^(x-1)` | `exp((x-1) * log(x))` | Variable exponent expression |

### Arena Allocator

| Function | Signature | Description |
|----------|-----------|-------------|
| NewArena | `func NewArena(capacity int) *Arena` | Create arena with initial capacity |
| Alloc | `func (a *Arena) Alloc() *ArenaNode` | Allocate a node (bump pointer, grows if needed) |
| Reset | `func (a *Arena) Reset()` | Reset arena for reuse |
| NumNodes | `func (a *Arena) NumNodes() int` | Number of allocated nodes |
| FromInterface | `func (a *Arena) FromInterface(n Node) *ArenaNode` | Convert Node to arena node |
| ToInterface | `func (a *Arena) ToInterface(n *ArenaNode) Node` | Convert arena node to Node |
| EvalArena | `func EvalArena(n *ArenaNode, x float64) float64` | Evaluate arena tree |

### Canonical EML Trees

| Function | Signature | Description |
|----------|-----------|-------------|
| CanonicalExp | `func CanonicalExp(x *EMLNode) *EMLNode` | exp(x) = eml(x, 1) |
| CanonicalLog | `func CanonicalLog(x *EMLNode) *EMLNode` | log(x) = eml(1, eml(eml(1, x), 1)) |
| CanonicalSin | `func CanonicalSin(x *EMLNode) *EMLNode` | sin(x) |
| CanonicalCos | `func CanonicalCos(x *EMLNode) *EMLNode` | cos(x) |
| CanonicalSqrt | `func CanonicalSqrt(x *EMLNode) *EMLNode` | sqrt(x) = exp(0.5 * log(x)) |
| Canonicalize | `func Canonicalize(n Node) *EMLNode` | Convert any Node to canonical EML form |
| Equiv | `func Equiv(a, b *EMLNode) bool` | Structural equivalence check |
| EMLSize | `func EMLSize(n *EMLNode) int` | Count nodes in tree |

---

## Package: internal/constants

Mathematical constants.

| Constant | Value | Description |
|----------|-------|-------------|
| One | 1.0 | Unit constant |
| E | 2.718281828459045... | Euler's number |
| Pi | 3.141592653589793... | π |
| Sqrt2 | 1.4142135623730951... | √2 |
| Sqrt3 | 1.7320508075688772... | √3 |
| Ln2 | 0.6931471805599453... | ln(2) |
| Ln10 | 2.302585092994046... | ln(10) |
| SqrtPi | 1.7724538509055159... | √π |
| Phi | 1.618033988749895... | Golden ratio φ |
| I | 0+1i | Imaginary unit |
| NegOne | -1.0 | Negative unit |
| Two | 2.0 | Two |
| Half | 0.5 | One half |
| ComplexOne | 1+0i | Complex unit |
| ComplexI | 0+1i | Complex imaginary unit |
| ComplexNegI | 0-1i | Complex negative imaginary unit |

---

## Implementation Notes

The core EML operator `eml(x,y) = exp(x) - ln(y)` serves as the theoretical foundation for all elementary functions. However, the actual implementations use platform-optimized paths:

- **Scalar operations**: Direct implementations using `math.*` functions or hand-coded assembly (AVX2/AVX512 on AMD64).
- **Batch operations**: SIMD-vectorized kernels that process 4-8 elements per cycle using architecture-specific assembly.
- **FastMath**: FMA-optimized polynomial approximations with relaxed IEEE 754 compliance.
- **JIT compiler**: x86-64 SSE2 codegen for math expressions parsed from strings, supporting 16 built-in functions, non-integer exponents (via `exp(y*log(x))`), and variable exponents.
- **Complex numbers**: `math/cmplx`-based operations on `complex128` slices, parallelized via the worker pool.
- **Arbitrary precision**: `math/big.Float` Taylor series for symbolic verification at 256-bit precision.
- **Arena allocator**: Zero-allocation JIT parse/eval paths using bump-pointer allocation.
- **Canonical EML trees**: Normalize any expression to minimal EML form for comparison and optimization.
- **GPU backends**: CUDA and Metal kernels for massive parallel workloads.

The EML operator is used in the `Eml()` scalar function and the GPU `EmlBatch` kernel. The mathematical framework from the original EML paper (arXiv:2603.21852v2) demonstrates that all elementary functions can be derived from this single operator, which is the theoretical basis for the library's unified design.
