//go:build amd64 && !purego && !emlasm

package eml

// Transcendental dispatch on amd64.
//
// The hand-written AVX2 kernels for exp/log/sin/cos/tan are NOT used by
// default. Measured on an i7-12650H (AVX2, 16 cores), the generic path below is
// dramatically faster than the assembly it replaces:
//
//	n        expAVX2      generic (this file)   speedup
//	256          57,740 ns              6,150 ns     9.4x
//	4096         872,098 ns             21,907 ns    39.8x
//	65536     14,322,837 ns            106,226 ns   134.8x
//	1048576 211,667,934 ns          1,110,432 ns   190.6x
//
// The cause is in simd_amd64.s: each of those loops re-materialises its
// polynomial coefficients on every iteration via
// MOVQ $imm64 -> MOVQ AX, Xn -> VBROADCASTSD, which is three uops with a
// multi-cycle latency on a serial dependency chain through AX, per coefficient,
// per four elements.
//
// The assembly is therefore strictly dominated: it is both far slower and less
// accurate than parallelising math.Exp/math.Log. Build with `-tags emlasm` to
// select it while it is being reworked; see docs/nextsteps.md Step 1.
//
// The elementwise kernels (add/sub/mul/div/abs/neg/inv/sqrt/fma) are unaffected
// and remain dispatched from simd_dispatch_amd64.go: Go's compiler already
// auto-vectorises those loops to the same width, so hand assembly buys nothing.

func dispatchExpSIMDTo(x, result []float64) { parallelizeGeneric(x, result, nativeExp) }

func dispatchLogSIMDTo(x, result []float64) { parallelizeGeneric(x, result, logScalar) }

func dispatchSinSIMDTo(x, result []float64) { parallelizeGeneric(x, result, nativeSin) }

func dispatchCosSIMDTo(x, result []float64) { parallelizeGeneric(x, result, nativeCos) }

func dispatchTanSIMDTo(x, result []float64) { parallelizeGeneric(x, result, nativeTan) }
