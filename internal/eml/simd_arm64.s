//go:build arm64 && !purego

#include "textflag.h"

// func addNEON(a, b, result []float64)
TEXT ·addNEON(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $1, R3, R3       // R3 = pairs = n / 2
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VLD1.P 32(R1), [V2.D2, V3.D2]
	VFADD V2.D2, V0.D2, V4.D2
	VFADD V3.D2, V1.D2, V5.D2
	VST1.P [V4.D2, V5.D2], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VLD1.P 16(R1), [V2.D2]
	VFADD V2.D2, V0.D2, V4.D2
	VST1.P [V4.D2], 16(R2)

done:
	RET

// func subNEON(a, b, result []float64)
TEXT ·subNEON(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VLD1.P 32(R1), [V2.D2, V3.D2]
	VFSUB V2.D2, V0.D2, V4.D2
	VFSUB V3.D2, V1.D2, V5.D2
	VST1.P [V4.D2, V5.D2], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VLD1.P 16(R1), [V2.D2]
	VFSUB V2.D2, V0.D2, V4.D2
	VST1.P [V4.D2], 16(R2)

done:
	RET

// func mulNEON(a, b, result []float64)
TEXT ·mulNEON(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VLD1.P 32(R1), [V2.D2, V3.D2]
	VFMUL V2.D2, V0.D2, V4.D2
	VFMUL V3.D2, V1.D2, V5.D2
	VST1.P [V4.D2, V5.D2], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VLD1.P 16(R1), [V2.D2]
	VFMUL V2.D2, V0.D2, V4.D2
	VST1.P [V4.D2], 16(R2)

done:
	RET

// func divNEON(a, b, result []float64)
TEXT ·divNEON(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VLD1.P 32(R1), [V2.D2, V3.D2]
	VFDIV V2.D2, V0.D2, V4.D2
	VFDIV V3.D2, V1.D2, V5.D2
	VST1.P [V4.D2, V5.D2], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VLD1.P 16(R1), [V2.D2]
	VFDIV V2.D2, V0.D2, V4.D2
	VST1.P [V4.D2], 16(R2)

done:
	RET

// func addScalarNEON(a []float64, b float64, result []float64)
TEXT ·addScalarNEON(SB), NOSPLIT, $0-56
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	FMOVD b+24(FP), F0
	MOVD result_base+32(FP), R2

	VDUP V0.D[0], V6.D2

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VFADD V6.D2, V0.D2, V2.D2
	VFADD V6.D2, V1.D2, V3.D2
	VST1.P [V2.D2, V3.D2], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VFADD V6.D2, V0.D2, V2.D2
	VST1.P [V2.D2], 16(R2)

done:
	RET

// func mulScalarNEON(a []float64, b float64, result []float64)
TEXT ·mulScalarNEON(SB), NOSPLIT, $0-56
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	FMOVD b+24(FP), F0
	MOVD result_base+32(FP), R2

	VDUP V0.D[0], V6.D2

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VFMUL V6.D2, V0.D2, V2.D2
	VFMUL V6.D2, V1.D2, V3.D2
	VST1.P [V2.D2, V3.D2], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VFMUL V6.D2, V0.D2, V2.D2
	VST1.P [V2.D2], 16(R2)

done:
	RET

// func sqrtNEON(a, result []float64)
TEXT ·sqrtNEON(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R1

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VFSQRT V0.D2, V2.D2
	VFSQRT V1.D2, V3.D2
	VST1.P [V2.D2, V3.D2], 32(R1)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VFSQRT V0.D2, V2.D2
	VST1.P [V2.D2], 16(R1)

done:
	RET

// func absNEON(a, result []float64)
TEXT ·absNEON(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R1

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VFABS V0.D2, V2.D2
	VFABS V1.D2, V3.D2
	VST1.P [V2.D2, V3.D2], 32(R1)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VFABS V0.D2, V2.D2
	VST1.P [V2.D2], 16(R1)

done:
	RET

// func negNEON(a, result []float64)
TEXT ·negNEON(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R1

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VFNEG V0.D2, V2.D2
	VFNEG V1.D2, V3.D2
	VST1.P [V2.D2, V3.D2], 32(R1)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VFNEG V0.D2, V2.D2
	VST1.P [V2.D2], 16(R1)

done:
	RET

// func invNEON(a, result []float64)
TEXT ·invNEON(SB), NOSPLIT, $0-48
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD result_base+24(FP), R1

	FMOVD $1.0, F0
	VDUP V0.D[0], V6.D2

	LSR $1, R3, R3
	CBZ R3, done

	CMP $2, R3
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VFDIV V0.D2, V6.D2, V2.D2
	VFDIV V1.D2, V6.D2, V3.D2
	VST1.P [V2.D2, V3.D2], 32(R1)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop4

tail2:
	CBZ R3, done
	VLD1.P 16(R0), [V0.D2]
	VFDIV V0.D2, V6.D2, V2.D2
	VST1.P [V2.D2], 16(R1)

done:
	RET

// func fmaNEON(a, b, c, result []float64)
TEXT ·fmaNEON(SB), NOSPLIT, $0-96
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R4
	MOVD b_base+24(FP), R1
	MOVD c_base+48(FP), R2
	MOVD result_base+72(FP), R3

	LSR $1, R4, R4
	CBZ R4, done

	CMP $2, R4
	BLT tail2

loop4:
	VLD1.P 32(R0), [V0.D2, V1.D2]
	VLD1.P 32(R1), [V2.D2, V3.D2]
	VLD1.P 32(R2), [V4.D2, V5.D2]
	VFMLA V2.D2, V0.D2, V4.D2
	VFMLA V3.D2, V1.D2, V5.D2
	VST1.P [V4.D2, V5.D2], 32(R3)
	SUB $2, R4, R4
	CMP $2, R4
	BGE loop4

tail2:
	CBZ R4, done
	VLD1.P 16(R0), [V0.D2]
	VLD1.P 16(R1), [V2.D2]
	VLD1.P 16(R2), [V4.D2]
	VFMLA V2.D2, V0.D2, V4.D2
	VST1.P [V4.D2], 16(R3)

done:
	RET

// func fmaScalar(a, b, c float64) float64
TEXT ·fmaScalar(SB), NOSPLIT, $0-32
	FMOVD a+0(FP), F0
	FMOVD b+8(FP), F1
	FMOVD c+16(FP), F2
	FMADDD F0, F2, F1, F0
	FMOVD F0, ret+24(FP)
	RET

// func addSatInt8NEON(a, b, result []int8)
TEXT ·addSatInt8NEON(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $4, R3, R3       // R3 = chunks of 16
	CBZ R3, done

	CMP $2, R3
	BLT tail16

loop32:
	VLD1.P 32(R0), [V0.B16, V1.B16]
	VLD1.P 32(R1), [V2.B16, V3.B16]
	VSQADD V2.B16, V0.B16, V4.B16
	VSQADD V3.B16, V1.B16, V5.B16
	VST1.P [V4.B16, V5.B16], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop32

tail16:
	CBZ R3, done
	VLD1.P 16(R0), [V0.B16]
	VLD1.P 16(R1), [V2.B16]
	VSQADD V2.B16, V0.B16, V4.B16
	VST1.P [V4.B16], 16(R2)

done:
	RET

// func subSatInt8NEON(a, b, result []int8)
TEXT ·subSatInt8NEON(SB), NOSPLIT, $0-72
	MOVD a_base+0(FP), R0
	MOVD a_len+8(FP), R3
	MOVD b_base+24(FP), R1
	MOVD result_base+48(FP), R2

	LSR $4, R3, R3       // R3 = chunks of 16
	CBZ R3, done

	CMP $2, R3
	BLT tail16

loop32:
	VLD1.P 32(R0), [V0.B16, V1.B16]
	VLD1.P 32(R1), [V2.B16, V3.B16]
	VSQSUB V2.B16, V0.B16, V4.B16
	VSQSUB V3.B16, V1.B16, V5.B16
	VST1.P [V4.B16, V5.B16], 32(R2)
	SUB $2, R3, R3
	CMP $2, R3
	BGE loop32

tail16:
	CBZ R3, done
	VLD1.P 16(R0), [V0.B16]
	VLD1.P 16(R1), [V2.B16]
	VSQSUB V2.B16, V0.B16, V4.B16
	VST1.P [V4.B16], 16(R2)

done:
	RET
