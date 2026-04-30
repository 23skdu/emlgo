#include "textflag.h"

// func sqrtScalar(x float64) float64
TEXT ·sqrtScalar(SB), NOSPLIT, $0-16
	FMOVD x+0(FP), F0
	FSQRTD F0, F0
	FMOVD F0, ret+8(FP)
	RET

// func fmaScalar(a, b, c float64) float64
TEXT ·fmaScalar(SB), NOSPLIT, $0-32
	FMOVD a+0(FP), F0
	FMOVD b+8(FP), F1
	FMOVD c+16(FP), F2
	FMADDD F0, F1, F2, F0 // F0 = F0 * F1 + F2
	FMOVD F0, ret+24(FP)
	RET
