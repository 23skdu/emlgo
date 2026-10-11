// Package fastmath provides high-performance math approximations using polynomial minimax methods.
//
// Approximation error and supported input ranges vary by function; see the
// function documentation and docs/accuracy.md. These functions are intended
// for workloads that explicitly accept approximation error, not as replacements
// for standard-accuracy math operations.
//
// Sqrt and FMA delegate to hand-tuned assembly (SQRTSD/VFMADD on amd64) for maximum throughput.
package fastmath
