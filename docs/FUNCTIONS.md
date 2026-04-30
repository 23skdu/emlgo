# Function Reference

Complete list of all public functions in the emlgo library.

---

## Package: pkg/logexp

Exponential and logarithmic functions using EML operator.

| Function | Signature | Description |
|----------|-----------|-------------|
| Exp | `func Exp(x float64) float64` | Exponential function e^x |
| Log | `func Log(x float64) float64` | Natural logarithm ln(x) |

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
| Pow | `func Pow(x, y float64) float64` | Power x^y |
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
| Max | `func Max(x, y float64) float64` | Maximum of x and y |
| Min | `func Min(x, y float64) float64` | Minimum of x and y |

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
| EmlOne | `func EmlOne(x float64) float64` | Eml(x, 1) = exp(x) - 0 = exp(x) |
| OneEml | `func OneEml(y float64) float64` | Eml(1, y) = e - ln(y) |
| HasAVX2 | `func HasAVX2() bool` | AVX2 detection |
| HasAVX512 | `func HasAVX512() bool` | AVX-512 detection |
| HasNeon | `func HasNeon() bool` | ARM NEON detection |

---

## Package: internal/constants

Mathematical constants.

| Constant | Value | Description |
|----------|-------|-------------|
| One | 1.0 | Unit constant |
| E | 2.718281828459045... | Euler's number |
| Pi | 3.141592653589793... | π |
| I | 0+1i | Imaginary unit |

---

## Implementation Notes

All functions in this library are implemented using the EML operator `eml(x,y) = exp(x) - ln(y)` as the single primitive. This includes:

- **Exp/Log**: Direct EML implementations
- **Sin/Cos/Tan**: Using complex exponentials with EML-based magnitude
- **Sinh/Cosh/Tanh**: Using EML via logexp.Exp
- **Sqrt**: Using exp(log(x)/2) via EML
- **Pow**: Using exp(y * log(x)) via EML

This approach follows the mathematical framework from the original EML paper.