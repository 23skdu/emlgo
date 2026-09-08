# Changelog

All notable changes to the emlgo library will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- JIT compilation support for non-integer exponents (x^0.5, x^1.5, etc.)
- JIT compilation support for variable exponents (x^x, x^(x-1), etc.)
- GPU stub methods: AbsBatch, NegBatch, PowBatch, InvBatch, FmaBatch
- GPU BatchVerifier.VerifyBinaryOp for two-input GPU operations
- WASM SIMD batch functions for exp, log, sin, cos, tan, sinh, cosh, tanh
- WASM sincosWasmSIMD for combined sine/cosine computation
- This CHANGELOG

### Changed
- WASM dispatch now uses direct SIMD batch functions instead of parallelizeGeneric
- Improved JIT eval to support cbrt, log2, log10, ceil, floor, trunc

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
