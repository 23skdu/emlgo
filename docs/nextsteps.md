# Next Steps: EML Go Library Implementation Plan

## Project Goal

Implement a high-performance Go library that provides all elementary mathematical functions using only the EML operator `eml(x, y) = exp(x) - ln(y)` and the constant `1`. The library must follow Go best practices, have no external dependencies, and utilize SIMD instructions where possible.

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
- [ ] Initialize git repository and create initial commit
- [ ] Write README.md with usage examples
- [ ] Configure goreleaser for releases
- [ ] Add gofmt/gci linting configuration

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
- [ ] Implement actual AVX2/AVX512 assembly (requires Go assembly)
- [ ] Benchmark SIMD vs scalar implementations
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
- [ ] Document performance characteristics in benchmarks/

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
- [ ] Write performance report comparing to standard library
- [ ] Tag v1.0.0 release
- [ ] Add badge for go.dev reference
- [x] Publish to GitHub and verify go install works

---

## Implementation Priority

1. **Phase 1 (Core):** Steps 1-3 - Project setup, EML operator, constants
2. **Phase 2 (Functions):** Steps 4-7 - All elementary functions
3. **Phase 3 (Optimization):** Step 8 - SIMD implementation
4. **Phase 4 (Quality):** Step 9 - Testing and verification
5. **Phase 5 (Release):** Step 10 - Documentation and release

---

## Technical Considerations

- **Accuracy Target:** Within 1 ULP of math library for primary functions
- **Performance Target:** Within 10x of math library (SIMD should bring closer)
- **Compatibility:** Go 1.21+ (for latest stdlib features)
- **Dependencies:** None (pure Go + assembly for SIMD)
- **Platform Support:** linux/amd64, darwin/amd64, windows/amd64, linux/arm64

---

## Related Resources

- Original Paper: arXiv:2603.21852v2
- Reference Implementation: EML_toolkit (Mathematica/Rust)
- Related: Kolmogorov-Arnold Networks (KAN) for similar tree structures