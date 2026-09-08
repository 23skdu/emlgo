// Package fastmath provides high-performance math approximations using polynomial minimax methods.
//
// Functions in this package use Cody-Waite range reduction and minimax polynomial evaluation
// with FMA (fused multiply-add) instructions for reduced rounding error. Accuracy is typically
// within 1e-7 for the primary input range, making these suitable for graphics, gaming, and
// other applications where performance is prioritized over full IEEE 754 precision.
//
// Sqrt and FMA delegate to hand-tuned assembly (SQRTSD/VFMADD on amd64) for maximum throughput.
package fastmath
