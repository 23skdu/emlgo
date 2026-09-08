// Package eml implements the EML (Exp-Minus-Log) operator and provides SIMD-accelerated
// batch math operations.
//
// The EML operator is defined as eml(x, y) = exp(x) - log(y). Together with the constant 1,
// it can derive all elementary functions. This package provides the core SIMD dispatch layer
// that routes operations to the fastest available implementation on the current platform
// (AVX-512, AVX2, NEON, SVE, WASM SIMD128, or scalar fallback).
//
// The worker pool in this package handles fused batch operations (ExpMul, ExpAdd, LogDiv, LogSub)
// and parallelized generic operations using pre-allocated goroutines.
package eml
