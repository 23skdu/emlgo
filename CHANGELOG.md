# Changelog

All notable changes to the emlgo library will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Symbolic differentiation engine (`Diff`, `DiffEval`) with chain rule support
- Expression simplification (`Simplify`) with constant folding, identity reduction, algebraic simplifications
- EML decompiler (`Decompile`, `DecompileLaTeX`, `DecompileNodeToExpr`) for tree → infix/LaTeX conversion
- Composable zero-allocation `Pipeline` API with buffer swapping (Exp, Log, Sqrt, Sin, Cos, Abs, Neg, MulScalar, AddScalar)
- `float32` SIMD batch operations (Exp, Log, Sqrt, Add, Sub, Mul, Div, Abs, Neg, Inv, Sin, Cos, Tan, scalar ops)
- JIT expression LRU cache (`CompileCached`, `ClearJITCache`, max 1024 entries)
- `round` function added to JIT compiler (17 functions total)
- `ParseWithVars` for multi-variable expression parsing
- `ParseError` struct with structured error information (Pos, Token, Msg)
- `NewFloatFromInt(n int64)` constructor for bigmath
- `bigmath.Tan`, `bigmath.Atan`, `bigmath.Asin`, `bigmath.Acos` for complete transcendental coverage
- `StopWorkerPool` concurrency safety via `sync.Once`
- CLI `--decompile` flag for command-line expression decompilation

### Fixed
- Data race in `jitCache.get()`: `MoveToFront` was mutating list under `RLock` (now uses `Lock`)
- CLI `runDecompile()` now actually calls `Decompile`/`DecompileLaTeX` instead of discarding the parsed node

### Changed
- Documentation completely overhauled: README, architecture, functions, JIT, usage guides updated

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
