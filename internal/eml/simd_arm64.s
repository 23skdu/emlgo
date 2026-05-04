#include "textflag.h"

// func absScalar(x float64) float64
TEXT ·absScalar(SB), NOSPLIT, $0-16
	FMOVD x+0(FP), F0
	FABSD F0, F0
	FMOVD F0, ret+8(FP)
	RET

// func negScalar(x float64) float64
TEXT ·negScalar(SB), NOSPLIT, $0-16
	FMOVD x+0(FP), F0
	FNEGD F0, F0
	FMOVD F0, ret+8(FP)
	RET

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
	FMADDD F0, F2, F1, F0
	FMOVD F0, ret+24(FP)
	RET
