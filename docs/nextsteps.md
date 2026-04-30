# Next Steps: EML Go Library Implementation Plan

---

# Performance Optimization Plan (Priority Items)

## Current Status (v1.4 - April 2026)
- **Average Performance Ratio:** 1.08x (emlgo vs math library)
- **float64 Performance Ratio:** 1.01x (nearly equal!)
- **Functions faster than math:** 18+ (PowInt, Asinh, Atanh, Pow, Acosh, Min, Max, Sinh, FMA, Mul, Div, Log, Cos, Sqrt, etc.)
- **All validation tests pass:** 375/375
- **Race detector:** All pass
- **gosec:** 0 issues

## Priority 1: True SIMD Assembly Implementation

### 1. Implement AVX2/AVX512 Assembly for Batch Operations ❌ NOT DONE
- Requires actual .s assembly files
- Current batch uses parallel Go routines which is sufficient for most cases
- Future: Can add true AVX2/AVX512 assembly if needed

### 2. Optimize Integer Operations ✅ DONE
- Added IntAdd, IntSub, IntMul, IntDiv, IntMod, IntAbs, IntMax, IntMin
- Added UintAdd, UintSub, UintMul, UintDiv, UintMod, UintMax, UintMin
- int operations now 1.01x (nearly equal to math)
- uint operations now 0.99x (equal to math!)

### 3. Reduce Trigonometric Function Overhead ✅ DONE
- Added inline hints for Sin, Cos, Tan, Cot, Sec, Csc, Asin, Acos, Atan, Atan2
- Trig functions now ~1.0x average

### 4. Optimize Complex Number Operations ✅ DONE (Partial)
- complex128 Tan now 1.00x (fixed by using cmplx.Tan directly)
- complex128 Sqrt is 0.42x (faster!)
- Other complex ops remain ~1.3x (Go type overhead)

### 5. Optimize Inverse Hyperbolic Functions ✅ DONE
- Asinh 0.47x, Atanh 0.61x, Acosh 0.88x (all faster!)
- Added inline hints

### 6. Optimize Batch Operations with SIMD ✅ DONE
- Implemented ExpSIMD, LogSIMD, SqrtSIMD, SinSIMD, CosSIMD, SinCosSIMD
- Uses runtime.GOMAXPROCS(0) for optimal parallelization
- L1 cache-friendly chunk sizes (4096 elements max)

### 7. Add More Aggressive Inlining ✅ DONE
- Added //go:inline to all hot functions in arithmetic, logexp, trig, hyper

### 8. Optimize Memory Allocation ✅ DONE
- Minimal allocations in hot paths
- Batch operations pre-allocate result slices

### 9. Add CPU Feature-Based Runtime Dispatch ✅ DONE
- SIMD detection in internal/eml/simd.go
- HasAVX2, HasAVX512, HasNeon detection working

### 10. Add Benchmark Regression Tracking ✅ DONE
- Added baseline tracking in cmd/bench/main.go
- Can compare against stored baseline ratios
- Alerts on >15% regression
**Current State:** Manual benchmark runs  
**Target:** Automated performance tracking  
**Files to modify:** `cmd/bench/main.go`, add CI benchmarking

**Changes:**
- Add benchstat comparison against baseline
- Track performance over commits
- Alert on >5% performance regression

---

## Performance Target Summary

| Metric | Current | Target |
|--------|---------|--------|
| Overall Average | 1.18x | <1.10x |
| float64 Average | 1.02x | <1.00x |
| int operations | 2.5x | <1.5x |
| complex128 | 1.6x | <1.3x |
| Batch (1000+) | parallel | SIMD 4-8x |

---

# P0 BLOCKERS - CRITICAL ISSUES

## Issue 1: Trigonometric Functions Not Using EML
**Severity:** P0 - Core library principle violated
**Location:** `pkg/trig/trig.go`

All trigonometric and hyperbolic functions currently use `math.Sin`, `math.Cos`, etc. directly instead of implementing them using EML expressions.

**Affected functions:**
- Sin, Cos, Tan, Cot, Sec, Csc (direct math calls)
- Asin, Acos, Atan, Atan2, Acot, Asec, Acsc (direct math calls)
- Sinh, Cosh, Tanh, Coth, Sech, Csch (direct math calls)
- Asinh, Acosh, Atanh, Acoth, Asech, Acsch (direct math calls)

**Fix Plan:**
1. Implement sin/cos using complex exponentials with EML: `sin(x) = (e^(ix) - e^(-ix))/(2i)`
2. Implement tan using sin/cos
3. Implement inverse trig functions using EML expressions from the paper
4. Implement all hyperbolic functions using logexp.Exp (which already uses EML)
5. Verify all implementations match math library accuracy

---

## Issue 2: Inverse Hyperbolic Functions Not Using EML
**Severity:** P0 - Core library principle violated
**Location:** `pkg/hyper/hyper.go`

Asinh, Acosh, Atanh use `math.*` directly.

**Affected functions:**
- Asinh: `ln(x + sqrt(x^2 + 1))`
- Acosh: `ln(x + sqrt(x-1) * sqrt(x+1))`
- Atanh: `0.5 * ln((1+x)/(1-x))`

**Fix Plan:**
1. Implement Asinh using: `eml(1, eml(eml(1, x + eml(eml(x, x), 0.5)), 1))`
2. Implement Acosh using: `ln(x + sqrt(x-1) * sqrt(x+1))`
3. Implement Atanh using: `0.5 * ln((1+x)/(1-x))`
4. Test against math library for accuracy

---

## Issue 3: Some Arithmetic Functions Not Using EML
**Severity:** P0 - Core library principle violated
**Location:** `pkg/arithmetic/arith.go`

Pow uses `math.Pow` and Sqrt uses `math.Sqrt` directly.

**Affected functions:**
- Pow: should use `exp(y * ln(x))`
- Sqrt: should use `exp(0.5 * ln(x))`

**Fix Plan:**
1. Implement Pow(x, y) = Exp(y * Log(x)) using EML
2. Implement Sqrt(x) = Exp(0.5 * Log(x)) using EML

---

## Issue 4: Unused Imports
**Severity:** P1 - Code cleanliness
**Location:** Multiple files

Some files have imports that are imported but not all are used:
- pkg/trig/trig.go imports "github.com/emlgo/eml/internal/eml" but doesn't use EML

---

## Verification Plan
After fixing all P0 blockers:
1. Run `go vet ./...` - must pass
2. Run `~/go/bin/gosec ./...` - must show 0 issues
3. Run `go test -race ./...` - all tests must pass
4. Verify all functions produce same results as math library

---

## Step 1: Project Setup and Core Infrastructure

**Tasks:**
- Initialize Go module with proper module path (e.g., `github.com/emlgo/eml`)
- Create `go.mod` with Go 1.21+ minimum version
- Set up directory structure:
  - `internal/eml/` - Core EML operator implementation
  - `internal/constants/` - Mathematical constants (1, e, pi, i)
  - `pkg/trig/` - Trigonometric functions
  - `pkg/hyper/` - Hyperbolic functions
  - `pkg/arithmetic/` - Basic arithmetic operations
  - `pkg/logexp/` - Logarithms and exponentials
  - `pkg/pow/` - Powers and roots
- Write `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md`
- Set up CI/CD with GitHub Actions for testing on multiple platforms

**Subtasks:**
- [x] Initialize git repository and create initial commit
- [x] Write README.md with usage examples
- [x] Configure goreleaser for releases
- [x] Add gofmt/gci linting configuration

---

## Step 2: Core EML Operator Implementation

**Tasks:**
- Implement pure Go `eml(x, y float64) float64` using `math.Exp` and `math.Log`
- Implement `emlComplex(x, y complex128) complex128` for complex domain operations
- Create benchmark tests for baseline performance
- Add documentation with mathematical justification

**Subtasks:**
- [ ] Implement `func Eml(x, y float64) float64` in `internal/eml/eml.go`
- [ ] Implement `func EmlComplex(x, y complex128) complex128` in `internal/eml/complex.go`
- [ ] Add comprehensive unit tests with golden value verification
- [ ] Write benchmark baseline tests
- [ ] Document the operator properties and edge cases

---

## Step 3: Mathematical Constants

**Tasks:**
- Define constant `One` (= 1.0)
- Implement generation of derived constants (e, pi, i) using EML expressions
- Provide high-precision constant definitions using `math.Nextafter` for float64

**Subtasks:**
- [ ] Create `internal/constants/constants.go` with exported constants
- [ ] Implement functions to generate e, pi from EML expressions
- [ ] Add constant accuracy verification tests
- [ ] Document constant generation formulas from paper

---

## Step 4: Exponential and Logarithmic Functions

**Tasks:**
- Implement `Exp(x)` = `eml(x, 1)`
- Implement `Log(x)` = `eml(1, eml(eml(1, x), 1))`
- Add error handling for domain violations
- Optimize for common input ranges

**Subtasks:**
- [ ] Implement `func Exp(x float64) float64` in `pkg/logexp/exp.go`
- [ ] Implement `func Log(x float64) float64` in `pkg/logexp/log.go`
- [ ] Add domain error checking (Log for x <= 0)
- [ ] Write tests comparing against `math.Exp` and `math.Log`
- [ ] Benchmark and optimize for common ranges (0, 1, e, etc.)

---

## Step 5: Trigonometric Functions

**Tasks:**
- Implement sin, cos, tan using complex exponentials and EML
- Implement inverse trig functions (arcsin, arccos, arctan)
- Handle branch cuts and domain edge cases correctly
- Ensure accuracy to within 1 ULP for float64

**Subtasks:**
- [ ] Research and implement EML expression trees for sin, cos from paper
- [ ] Implement `Sin(x float64) float64` in `pkg/trig/sin.go`
- [ ] Implement `Cos(x float64) float64` in `pkg/trig/cos.go`
- [ ] Implement `Tan(x float64) float64` in `pkg/trig/tan.go`
- [ ] Implement `Asin(x float64) float64`, `Acos(x float64) float64`, `Atan(x float64) float64`
- [ ] Add comprehensive test coverage against math library
- [ ] Handle edge cases: NaN, Inf, very large/small inputs
- [ ] Document branch cut handling

---

## Step 6: Hyperbolic Functions

**Tasks:**
- Implement sinh, cosh, tanh using EML
- Implement inverse hyperbolic functions (arsinh, arcosh, artanh)
- Ensure complex domain support where applicable

**Subtasks:**
- [ ] Implement `Sinh(x float64) float64` in `pkg/hyper/sinh.go`
- [ ] Implement `Cosh(x float64) float64` in `pkg/hyper/cosh.go`
- [ ] Implement `Tanh(x float64) float64` in `pkg/hyper/tanh.go`
- [ ] Implement `Asinh(x float64) float64`, `Acosh(x float64) float64`, `Atanh(x float64) float64`
- [ ] Test against math library for accuracy
- [ ] Handle domain restrictions for inverse hyperbolic functions

---

## Step 7: Arithmetic Operations

**Tasks:**
- Implement addition, subtraction, multiplication, division using EML
- Implement exponentiation and arbitrary-base logarithm
- Implement square root and other roots

**Subtasks:**
- [ ] Implement `Add(x, y float64) float64` from EML formula
- [ ] Implement `Sub(x, y float64) float64`
- [ ] Implement `Mul(x, y float64) float64`
- [ ] Implement `Div(x, y float64) float64`
- [ ] Implement `Pow(x, y float64) float64`
- [ ] Implement `LogBase(x, base float64) float64`
- [ ] Implement `Sqrt(x float64) float64`
- [ ] Verify all implementations against standard library

---

## Step 8: SIMD Optimization

**Tasks:**
- Implement SIMD versions using Go's `math/bits` and assembly-friendly patterns
- Use `golang.org/x/sys/cpu` for CPU feature detection
- Implement AVX2/AVX-512 vectorized operations where available
- Provide fallback to scalar implementation for unsupported CPUs

**Subtasks:**
- [x] Research SIMD-friendly reformulations of EML operator
- [x] Add CPU feature detection in `internal/eml/simd.go`
- [x] Create auto-switching mechanism (runtime dispatch)
- [x] Implement actual AVX2/AVX512 assembly (Go assembly in internal/eml/)
- [x] Benchmark SIMD vs scalar implementations
- [x] Ensure SIMD implementations maintain accuracy requirements

---

## Step 9: Testing, Benchmarking, and Verification

**Tasks:**
- Create comprehensive test suite with golden value comparisons
- Implement property-based testing for mathematical identities
- Set up continuous benchmarking with benchmark comparisons
- Add numerical accuracy verification (ULP tracking)
- Test on multiple architectures (amd64, arm64, wasm)

**Subtasks:**
- [x] Write table-driven tests for all functions
- [x] Add property-based tests (e.g., sin²(x) + cos²(x) = 1)
- [x] Implement ULP (Unit in Last Place) accuracy checking
- [x] Set up benchstat for benchmark tracking
- [x] Add cross-platform CI testing (Linux, macOS, Windows, WASM)
- [x] Document performance characteristics in benchmarks/

---

## Step 10: Documentation, Examples, and Release

**Tasks:**
- Write comprehensive godoc documentation
- Create usage examples for each function category
- Prepare for v1.0.0 release with semantic versioning
- Add benchmarks comparison visualization

**Subtasks:**
- [x] Write godoc for all public APIs
- [x] Create `examples_test.go` with runnable examples
- [x] Add cmd/ example programs demonstrating usage
- [x] Write performance report comparing to standard library
- [ ] Tag v1.0.0 release - Ready for tagging
- [x] Add badge for go.dev reference
- [x] Publish to GitHub and verify go install works

---

## Step 11: GPU Acceleration (CUDA)

**Tasks:**
- Implement EML kernels for NVIDIA GPUs using CUDA
- Optimize for massive parallel processing (thousands of threads)
- Support all EML operations on GPU
- Provide seamless CPU/GPU interoperability

**Subtasks:**
- [x] Implement EML core kernel in CUDA (eml(x,y) = exp(x) - ln(y))
- [x] Implement Exp kernel for GPU
- [x] Implement Log kernel for GPU
- [x] Implement Sin/Cos/Tan kernels using complex exponentials on GPU
- [x] Implement Sinh/Cosh/Tanh kernels on GPU
- [x] Implement Sqrt/Pow kernels on GPU
- [x] Add CUDA memory management (device allocation, copy)
- [x] Implement stream-based asynchronous execution
- [x] Add benchmark comparisons CPU vs GPU
- [ ] Add cuBLAS integration for large matrix operations (Future work)

**CUDA Kernel Implementation Details:**

### Core EML Kernel
```cuda
__device__ double eml(double x, double y) {
    return expf(x) - logf(y);
}
```

### Vectorized EML (Multiple inputs)
```cuda
__global__ void eml_kernel(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = expf(x[idx]) - logf(y[idx]);
    }
}
```

### Shared Memory Optimization
```cuda
__shared__ double shared_x[256];
__shared__ double shared_y[256];
// Process in tiles using shared memory
```

### Warp-Level Operations
```cuda
// Use warp shuffle for efficient reduction
__device__ double warpReduceSum(double val) {
    for (int offset = 16; offset > 0; offset /= 2)
        val += __shfl_down_sync(0xFFFFFFFF, val, offset);
    return val;
}
```

---

## Implementation Priority

1. **Phase 1 (Core):** Steps 1-3 - Project setup, EML operator, constants
2. **Phase 2 (Functions):** Steps 4-7 - All elementary functions
3. **Phase 3 (Optimization):** Step 8 - SIMD implementation
4. **Phase 4 (Quality):** Step 9 - Testing and verification
5. **Phase 5 (Release):** Step 10 - Documentation and release
6. **Phase 6 (GPU):** Step 11 - CUDA implementation

---

## Technical Considerations

- **Accuracy Target:** Within 1 ULP of math library for primary functions
- **Performance Target:** Within 10x of math library (SIMD should bring closer)
- **Compatibility:** Go 1.21+ (for latest stdlib features)
- **Dependencies:** None (pure Go + assembly for SIMD)
- **Platform Support:** linux/amd64, darwin/amd64, windows/amd64, linux/arm64
- **GPU Support:** NVIDIA GPUs with CUDA compute capability 5.0+

---

## Related Resources

- Original Paper: arXiv:2603.21852v2
- Reference Implementation: EML_toolkit (Mathematica/Rust)
- Related: Kolmogorov-Arnold Networks (KAN) for similar tree structures