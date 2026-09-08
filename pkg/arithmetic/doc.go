// Package arithmetic provides basic and advanced arithmetic operations with SIMD acceleration.
//
// The package supports scalar operations (Add, Sub, Mul, Div, etc.), integer and unsigned
// integer arithmetic, and batch operations that leverage platform-specific SIMD instructions
// (AVX2, AVX-512, NEON, SVE, WASM SIMD128) for high-throughput vectorized computation.
//
// Batch operations automatically dispatch to the fastest available SIMD path on the
// current platform, falling back to parallelized scalar code when SIMD is unavailable.
package arithmetic
