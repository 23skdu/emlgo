// Package hyper provides hyperbolic and inverse hyperbolic functions with overflow protection.
//
// All functions include explicit overflow thresholds based on ln(MaxFloat64) ≈ 709.78
// to prevent intermediate overflow in exp-based computations. Batch operations delegate
// to the internal SIMD dispatch layer for parallel execution.
package hyper
