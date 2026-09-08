// Package logexp provides exponential and logarithmic functions with SIMD batch support.
//
// Exp and Log provide standard IEEE 754 compliant implementations. ExpFast and LogFast
// provide overflow/underflow pre-checks for performance-critical paths. Batch operations
// (ExpBatch, LogBatch) use platform-specific SIMD acceleration.
package logexp
