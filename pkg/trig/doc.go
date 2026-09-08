// Package trig provides trigonometric, inverse trigonometric, and hyperbolic functions.
//
// All functions handle edge cases (NaN, Inf, signed zero) per IEEE 754 semantics.
// Batch operations (SinBatch, CosBatch, TanBatch, SinCosBatch) leverage SIMD dispatch
// for high-throughput vectorized computation.
//
// Fast variants (SinFast, CosFast, TanFast) skip some edge case handling for
// performance-critical paths where NaN/Inf inputs are not expected.
package trig
