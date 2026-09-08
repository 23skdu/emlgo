# Changelog

All notable changes to the emlgo library will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complex number batch operations (ComplexBatch, ComplexExpBatch, ComplexLogBatch, ComplexSinBatch, ComplexCosBatch, ComplexTanBatch)
- Arbitrary precision backend (internal/eml/bigmath) with Taylor series Exp/Log/Sin/Cos
- IdentityVerifier for symbolic identity verification at arbitrary precision
- Zero-allocation arena allocator for JIT AST nodes (ArenaNode, Arena, EvalArena)
- Canonical EML tree constructors (CanonicalExp, CanonicalLog, CanonicalSin, CanonicalCos, CanonicalSqrt)
- Canonicalize function to convert any AST to minimal EML form
- Equiv function for structural equivalence checking
- EMLSize function for tree complexity measurement
- StopWorkerPool for graceful worker pool shutdown
- CI/CD pipeline (.github/workflows/ci.yml) with Go 1.22/1.23 matrix
- Benchmark comparison script (scripts/bench-compare.sh)
- Numerically stable Pow using Log1p for x near 1
- Numerically stable Log using Log1p for x near 1
- Numerically stable Exp using Expm1 for x near 0
- JIT support for non-integer exponents (x^0.5, x^1.5, etc.)
- JIT support for variable exponents (x^x, x^(x-1), etc.)
- WASM SIMD batch functions for exp, log, sin, cos, tan, sinh, cosh, tanh
- WASM sincosWasmSIMD for combined sine/cosine computation
- GPU stub methods: AbsBatch, NegBatch, PowBatch, InvBatch, FmaBatch
- GPU BatchVerifier.VerifyBinaryOp for two-input GPU operations
- Documentation: COMPLEX.md, JIT.md, updated all existing docs

### Fixed
- AbsBranchless LSB corruption (was flipping sign bit + LSB, now sign bit only)
- ComplexCos wrong formula (was (Exp(z)+Exp(iz))/2, now (Exp(iz)+Exp(-iz))/2)
- Missing length validation on AddSIMD/SubSIMD/MulSIMD/DivSIMD
- movsdStore wrong REX prefix (0x41 → 0x44 for high XMM registers)
- MinBranchless/MaxBranchless now NaN-aware (returns non-NaN operand)
- IntAbs(MinInt) overflow guard (returns MaxInt)
- hasNeonDot defaults to false (requires ARMv8.2-A runtime detection)
- Pow integer conversion overflow guard
- Sec/Csc NaN/Inf guards added
- Acsch Inf guard added
- dispatchSinCosSIMDTo uses parallelizeSinCos (eliminates double range-reduction)
- go.mod updated from go 1.26.1 to go 1.23

### Changed
- WASM dispatch uses direct SIMD batch functions instead of parallelizeGeneric
- Improved JIT eval to support cbrt, log2, log10, ceil, floor, trunc
- GCD(MinInt64) documented as known limitation

## [0.3.0] - 2026-09-08

### Added
- Core EML operator: eml(x, y) = exp(x) - log(y)
- Complex number support via math/cmplx
- SIMD dispatch for AMD64 (AVX2, AVX-512), ARM64 (NEON, SVE), WASM
- JIT compiler for amd64 with x86-64 machine code generation
- GPU backends: CUDA (Linux/Windows) and Metal (macOS/ARM64)
- Arithmetic package: Add, Sub, Mul, Div, Mod, Pow, Sqrt, Cbrt, Hypot, etc.
- Trig package: Sin, Cos, Tan, Cot, Sec, Csc, Asin, Acos, Atan, Atan2, etc.
- Hyperbolic package: Sinh, Cosh, Tanh, Asinh, Acosh, Atanh
- LogExp package: Exp, Log, ExpBatch, LogBatch, ExpFast, LogFast
- FastMath package: polynomial approximations for Exp, Sin, Cos, Log
- Worker pool for fused batch operations (ExpMul, ExpAdd, LogDiv, LogSub)
- Branchless Abs, Min, Max, Select operations
- Mathematical constants package (e, pi, ln2, sqrt2, phi, etc.)
- CLI tools: emlcli, bench, validate
- Comprehensive test suite with fuzz testing
- GPU result verification with ULP-based comparison
