# Performance

This document provides benchmark results comparing emlgo's EML-based implementations against Go's standard `math` library.

## Benchmark Methodology

- **Platform**: Linux, AMD64
- **Iterations**: 100,000 per operation
- **Comparison**: emlgo vs Go's `math` library
- **Ratio**: `emlgo_time / math_time` (values >1 mean emlgo is slower)

## Summary Results

| Category | Average Ratio | Notes |
|----------|---------------|-------|
| int operations | ~1.5x | Float conversion overhead |
| uint operations | ~1.0x | Similar performance |
| float32 | ~4.5x | Additional type conversion |
| float64 | ~2.8x | Core EML operations |
| complex64 | ~5.0x | Multiple conversions |
| complex128 | ~5.2x | Complex arithmetic |

## Detailed float64 Results

### Exponential & Logarithmic
| Function | Ratio | Notes |
|----------|-------|-------|
| Exp | 2.00x | Uses EML exp via logexp.Exp |
| Log | 2.73x | Uses logexp.Log |
| Log2 | 3.59x | LogBase2 with logexp |
| Log10 | 4.67x | LogBase10 with logexp |

### Trigonometric
| Function | Ratio | Notes |
|----------|-------|-------|
| Sin | 1.53x | EML via complex128 |
| Cos | 1.78x | EML via complex128 |
| Tan | 2.94x | Sin/Cos ratio |
| Cot | 2.95x | 1/tan |
| Sec | 1.77x | 1/cos |
| Csc | 1.61x | 1/sin |
| Asin | 2.41x | Inverse trig |
| Acos | 2.50x | Inverse trig |
| Atan | 1.19x | Best performing trig |
| Atan2 | 1.04x | Two-arg atan |

### Hyperbolic
| Function | Ratio | Notes |
|----------|-------|-------|
| Sinh | 1.87x | (e^x - e^-x)/2 |
| Cosh | 3.26x | (e^x + e^-x)/2 |
| Tanh | 2.06x | sinh/cosh |
| Asinh | 6.08x | ln(x + sqrt(x^2+1)) |
| Acosh | 8.02x | ln(x + sqrt(x^2-1)) |
| Atanh | 1.33x | 0.5*ln((1+x)/(1-x)) |

### Arithmetic
| Function | Ratio | Notes |
|----------|-------|-------|
| Sqrt | 17.14x | EML-based vs hardware |
| Cbrt | 7.04x | Cube root |
| Pow | 1.88x | x^2.5 |
| PowInt | 0.16x | **Faster than math!** |
| Floor | 1.05x | Direct implementation |
| Ceil | 1.05x | Direct implementation |
| Round | 1.12x | Round to nearest |
| Trunc | 1.04x | Truncate fractional |
| Abs | 1.00x | Direct, no conversion |
| Neg | 1.01x | Direct negation |
| Inv | 1.00x | 1/x |
| Square | 1.00x | x*x |
| Cube | 1.02x | x*x*x |
| Max | 0.81x | **Faster than math!** |
| Min | 0.89x | **Faster than math!** |
| Hypot | 1.40x | sqrt(a^2 + b^2) |
| FMA | 0.96x | a*b + c |
| Add | 1.22x | a + b |
| Sub | 0.98x | a - b |
| Mul | 0.92x | a * b |
| Div | 0.95x | a / b |

## Key Observations

1. **Fastest operations**: Max, Min, Mul, Sub, Div, FMA - all within 10% of math library
2. **Slowest operations**: Sqrt (17x), Acosh (8x), PowInt actually faster (0.16x!)
3. **Trig functions**: 1.5-3x slower due to complex128 EML implementation
4. **Hyperbolic**: Higher variance due to multiple EML exp/log calls
5. **Type conversion**: Extra overhead for float32 and complex types

## Accuracy (ULP)

Maximum ULP (Unit in Last Place) difference from math library:

| Function | Max ULP | Notes |
|----------|---------|-------|
| Exp | 0 | Exact match |
| Log | 124 | Different algorithm |
| Sin | 0 | Exact match |
| Cos | 0 | Exact match |
| Tan | 2 | Very close |
| Sinh | 18 | Close |
| Cosh | 1 | Very close |
| Tanh | 17 | Close |
| Asinh | 130 | Different algorithm |
| Acosh | 7 | Very close |
| Atanh | 193 | Different algorithm |
| Sqrt | 2 | Very close |
| Pow | 11 | Close |

Most functions are within 20 ULP of the standard library, which is acceptable for scientific computing.

## Performance Optimization Notes

1. **SIMD**: AVX2/AVX512 batch operations provide ~4-8x speedup for vectorized work
2. **GPU**: CUDA kernels available for massive parallel workloads
3. **Integer ops**: Cast to float64 introduces overhead; consider direct implementations
4. **Complex ops**: Multiple conversions between float64 and complex128 add overhead

## Running Benchmarks

```bash
# Full benchmark
go run cmd/bench/main.go -n 100000

# Feature parity test
go run cmd/bench/main.go -compare

# Accuracy test
go run cmd/bench/main.go -accuracy

# Filter by type
go run cmd/bench/main.go -type float64 -n 100000
```

## Conclusion

emlgo provides mathematically correct implementations of all elementary functions using the EML operator. Performance is within 3x of the standard library for most operations, with some operations (Max, Min, PowInt) actually faster. The trade-off is acceptable for applications requiring EML-based computations.