//go:build arm64 && !purego

#include "textflag.h"

// func addNEONF32(a, b, result []float32)
TEXT ·addNEONF32(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VLD1.P 16(R1), [V2.S4]
	VFADD V2.S4, V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func subNEONF32(a, b, result []float32)
TEXT ·subNEONF32(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VLD1.P 16(R1), [V2.S4]
	VFSUB V2.S4, V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func mulNEONF32(a, b, result []float32)
TEXT ·mulNEONF32(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VLD1.P 16(R1), [V2.S4]
	VFMUL V2.S4, V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func divNEONF32(a, b, result []float32)
TEXT ·divNEONF32(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VLD1.P 16(R1), [V2.S4]
	VFDIV V2.S4, V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func sqrtNEONF32(a, result []float32)
TEXT ·sqrtNEONF32(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VFSQRT V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func addScalarNEONF32(a []float32, b float32, result []float32)
TEXT ·addScalarNEONF32(SB), NOSPLIT, $0-56
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	FMOVS b+24(FP), F0
	MOVD result_base+32(FP), R2

	VDUP V0.S[0], V6.S4

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VFADD V6.S4, V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func mulScalarNEONF32(a []float32, b float32, result []float32)
TEXT ·mulScalarNEONF32(SB), NOSPLIT, $0-56
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	FMOVS b+24(FP), F0
	MOVD result_base+32(FP), R2

	VDUP V0.S[0], V6.S4

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VFMUL V6.S4, V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func absNEONF32(a, result []float32)
TEXT ·absNEONF32(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VFABS V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func negNEONF32(a, result []float32)
TEXT ·negNEONF32(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VFNEG V0.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func invNEONF32(a, result []float32)
TEXT ·invNEONF32(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R2

	FMOVS $1.0, F0
	VDUP V0.S[0], V6.S4

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VFDIV V0.S4, V6.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET

// func fmaNEONF32(a, b, c, result []float32)
TEXT ·fmaNEONF32(SB), NOSPLIT, $0-96
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD c_base+48(FP), R5
	MOVD result_base+72(FP), R2

	LSR $2, R3, R3
	CBZ R3, done

loop:
	VLD1.P 16(R0), [V0.S4]
	VLD1.P 16(R1), [V1.S4]
	VLD1.P 16(R5), [V2.S4]
	VFMUL V1.S4, V0.S4, V3.S4
	VFADD V2.S4, V3.S4, V4.S4
	VST1.P [V4.S4], 16(R2)
	SUB $1, R3, R3
	CBNZ R3, loop

done:
	RET
