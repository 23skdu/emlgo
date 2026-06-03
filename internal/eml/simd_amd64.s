#include "textflag.h"

// func cpuid(op, op2 uint32) (eax, ebx, ecx, edx uint32)
TEXT ·cpuid(SB), NOSPLIT, $0-24
	MOVL op+0(FP), AX
	MOVL op2+4(FP), CX
	CPUID
	MOVL AX, eax+8(FP)
	MOVL BX, ebx+12(FP)
	MOVL CX, ecx+16(FP)
	MOVL DX, edx+20(FP)
	RET

// func addAVX2(a, b, result []float64)
TEXT ·addAVX2(SB), NOSPLIT, $0-72
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX // n / 4
	JZ done

loop_add:
	VMOVUPD (SI), Y0
	VMOVUPD (DI), Y1
	VADDPD Y1, Y0, Y2
	VMOVUPD Y2, (DX)

	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_add

done:
	VZEROUPPER
	RET

// func subAVX2(a, b, result []float64)
TEXT ·subAVX2(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_sub
loop_sub:
    VMOVUPD (SI), Y0
    VMOVUPD (DI), Y1
    VSUBPD Y1, Y0, Y2
    VMOVUPD Y2, (DX)
    ADDQ $32, SI
    ADDQ $32, DI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_sub
done_sub:
    VZEROUPPER
    RET

// func mulAVX2(a, b, result []float64)
TEXT ·mulAVX2(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_mul
loop_mul:
    VMOVUPD (SI), Y0
    VMOVUPD (DI), Y1
    VMULPD Y1, Y0, Y2
    VMOVUPD Y2, (DX)
    ADDQ $32, SI
    ADDQ $32, DI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_mul
done_mul:
    VZEROUPPER
    RET

// func divAVX2(a, b, result []float64)
TEXT ·divAVX2(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_div
loop_div:
    VMOVUPD (SI), Y0
    VMOVUPD (DI), Y1
    VDIVPD Y1, Y0, Y2
    VMOVUPD Y2, (DX)
    ADDQ $32, SI
    ADDQ $32, DI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_div
done_div:
    VZEROUPPER
    RET

// func addScalarAVX2(a []float64, b float64, result []float64)
TEXT ·addScalarAVX2(SB), NOSPLIT, $0-56
    MOVQ a_base+0(FP), SI
    MOVQ a_len+8(FP), CX
    MOVSD b+24(FP), X0
    MOVQ result_base+32(FP), DX
    SHRQ $2, CX
    JZ done_add_scalar
    VINSERTF128 $1, X0, Y0, Y0
    VPERMPD $0, Y0, Y0
loop_add_scalar:
    VMOVUPD (SI), Y1
    VADDPD Y0, Y1, Y2
    VMOVUPD Y2, (DX)
    ADDQ $32, SI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_add_scalar
done_add_scalar:
    VZEROUPPER
    RET

// func mulScalarAVX2(a []float64, b float64, result []float64)
TEXT ·mulScalarAVX2(SB), NOSPLIT, $0-56
    MOVQ a_base+0(FP), SI
    MOVQ a_len+8(FP), CX
    MOVSD b+24(FP), X0
    MOVQ result_base+32(FP), DX
    SHRQ $2, CX
    JZ done_mul_scalar
    VINSERTF128 $1, X0, Y0, Y0
    VPERMPD $0, Y0, Y0
loop_mul_scalar:
    VMOVUPD (SI), Y1
    VMULPD Y0, Y1, Y2
    VMOVUPD Y2, (DX)
    ADDQ $32, SI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_mul_scalar
done_mul_scalar:
    VZEROUPPER
    RET

// AVX512 Kernels

// func addAVX512(a, b, result []float64)
TEXT ·addAVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX // n / 8
    JZ done_add512
loop_add512:
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VADDPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_add512
done_add512:
    RET

// func subAVX512(a, b, result []float64)
TEXT ·subAVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_sub512
loop_sub512:
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VSUBPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_sub512
done_sub512:
    RET

// func mulAVX512(a, b, result []float64)
TEXT ·mulAVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_mul512
loop_mul512:
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VMULPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_mul512
done_mul512:
    RET

// func divAVX512(a, b, result []float64)
TEXT ·divAVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_div512
loop_div512:
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VDIVPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_div512
done_div512:
    RET

// func addScalarAVX512(a []float64, b float64, result []float64)
TEXT ·addScalarAVX512(SB), NOSPLIT, $0-56
    MOVQ a_base+0(FP), SI
    MOVQ a_len+8(FP), CX
    MOVSD b+24(FP), X0
    MOVQ result_base+32(FP), DX
    SHRQ $3, CX
    JZ done_add_scalar512
    VBROADCASTSD X0, Z0
loop_add_scalar512:
    VMOVUPD (SI), Z1
    VADDPD Z0, Z1, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_add_scalar512
done_add_scalar512:
    RET

// func mulScalarAVX512(a []float64, b float64, result []float64)
TEXT ·mulScalarAVX512(SB), NOSPLIT, $0-56
    MOVQ a_base+0(FP), SI
    MOVQ a_len+8(FP), CX
    MOVSD b+24(FP), X0
    MOVQ result_base+32(FP), DX
    SHRQ $3, CX
    JZ done_mul_scalar512
    VBROADCASTSD X0, Z0
loop_mul_scalar512:
    VMOVUPD (SI), Z1
    VMULPD Z0, Z1, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_mul_scalar512
done_mul_scalar512:
    RET

// func sqrtAVX2(a, result []float64)
TEXT ·sqrtAVX2(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_sqrt2
loop_sqrt2:
    VSQRTPD (SI), Y0
    VMOVUPD Y0, (DX)
    ADDQ $32, SI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_sqrt2
done_sqrt2:
    VZEROUPPER
    RET

// func sqrtAVX512(a, result []float64)
TEXT ·sqrtAVX512(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_sqrt512
loop_sqrt512:
    VSQRTPD (SI), Z0
    VMOVUPD Z0, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_sqrt512
done_sqrt512:
    RET

// func sqrtScalar(x float64) float64
TEXT ·sqrtScalar(SB), NOSPLIT, $0-16
	MOVSD x+0(FP), X0
	SQRTSD X0, X0
	MOVSD X0, ret+8(FP)
	RET

// func fmaScalar(a, b, c float64) float64
// Computes: a*b + c
TEXT ·fmaScalar(SB), NOSPLIT, $0
	MOVSD a+0(FP), X0
	MOVSD b+8(FP), X1
	MOVSD c+16(FP), X2
	// If FMA is supported, use VFMADD213SD
	// Otherwise fallback to MULSD + ADDSD (handled by detection in Go if needed, 
	// but here we just use the instruction if we assume HasFMA is checked)
	// Actually, we can use VFMADD213SD directly if we know it's supported.
	VFMADD213SD X2, X1, X0
	MOVSD X0, ret+24(FP)
	RET

// func absScalar(x float64) float64
TEXT ·absScalar(SB), NOSPLIT, $0
	MOVQ x+0(FP), AX
	BTRQ $63, AX // Clear bit 63 (sign bit)
	MOVQ AX, ret+8(FP)
	RET

// func negScalar(x float64) float64
TEXT ·negScalar(SB), NOSPLIT, $0
	MOVQ x+0(FP), AX
	BTCQ $63, AX // Complement bit 63
	MOVQ AX, ret+8(FP)
	RET

// func fmaAVX2(a, b, c, result []float64)
TEXT ·fmaAVX2(SB), NOSPLIT, $0-96
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ c_base+48(FP), BX
	MOVQ result_base+72(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX
	JZ done_fma2

loop_fma2:
	VMOVUPD (SI), Y0
	VMOVUPD (DI), Y1
	VMOVUPD (BX), Y2
	VFMADD213PD Y2, Y1, Y0
	VMOVUPD Y0, (DX)

	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, BX
	ADDQ $32, DX
	DECQ CX
	JNZ loop_fma2

done_fma2:
	VZEROUPPER
	RET

// func fmaAVX512(a, b, c, result []float64)
TEXT ·fmaAVX512(SB), NOSPLIT, $0-96
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ c_base+48(FP), BX
	MOVQ result_base+72(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $3, CX
	JZ done_fma512

loop_fma512:
	VMOVUPD (SI), Z0
	VMOVUPD (DI), Z1
	VMOVUPD (BX), Z2
	VFMADD213PD Z2, Z1, Z0
	VMOVUPD Z0, (DX)

	ADDQ $64, SI
	ADDQ $64, DI
	ADDQ $64, BX
	ADDQ $64, DX
	DECQ CX
	JNZ loop_fma512

done_fma512:
	RET

// func absAVX2(a, result []float64)
TEXT ·absAVX2(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_abs2
    MOVQ $0x7fffffffffffffff, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Y15
loop_abs2:
    VMOVUPD (SI), Y0
    VANDPD Y15, Y0, Y0
    VMOVUPD Y0, (DX)
    ADDQ $32, SI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_abs2
done_abs2:
    VZEROUPPER
    RET

// func negAVX2(a, result []float64)
TEXT ·negAVX2(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_neg2
    MOVQ $0x8000000000000000, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Y15
loop_neg2:
    VMOVUPD (SI), Y0
    VXORPD Y15, Y0, Y0
    VMOVUPD Y0, (DX)
    ADDQ $32, SI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_neg2
done_neg2:
    VZEROUPPER
    RET

// func invAVX2(a, result []float64)
TEXT ·invAVX2(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_inv2
    MOVQ $0x3ff0000000000000, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Y15
loop_inv2:
    VMOVUPD (SI), Y0
    VDIVPD Y0, Y15, Y1
    VMOVUPD Y1, (DX)
    ADDQ $32, SI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_inv2
done_inv2:
    VZEROUPPER
    RET

// func absAVX512(a, result []float64)
TEXT ·absAVX512(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_abs512
    MOVQ $0x7fffffffffffffff, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Z15
loop_abs512:
    VMOVUPD (SI), Z0
    VPANDQ Z15, Z0, Z0
    VMOVUPD Z0, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_abs512
done_abs512:
    RET

// func negAVX512(a, result []float64)
TEXT ·negAVX512(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_neg512
    MOVQ $0x8000000000000000, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Z15
loop_neg512:
    VMOVUPD (SI), Z0
    VPXORQ Z15, Z0, Z0
    VMOVUPD Z0, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_neg512
done_neg512:
    RET

// func invAVX512(a, result []float64)
TEXT ·invAVX512(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_inv512
    MOVQ $0x3ff0000000000000, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Z15
loop_inv512:
    VMOVUPD (SI), Z0
    VDIVPD Z0, Z15, Z1
    VMOVUPD Z1, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_inv512
done_inv512:
    RET

DATA add_1023_val<>+0(SB)/4, $1023
DATA add_1023_val<>+4(SB)/4, $1023
DATA add_1023_val<>+8(SB)/4, $1023
DATA add_1023_val<>+12(SB)/4, $1023
GLOBL add_1023_val<>(SB), RODATA, $16

// func expAVX2(a, result []float64)
TEXT ·expAVX2(SB), NOSPLIT, $0-48
	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX // n / 4
	JZ done_exp_avx2

	// Load constant 1/ln2 into Y7
	MOVQ $0x3ff71547652b82fe, AX
	MOVQ AX, X7
	VBROADCASTSD X7, Y7

	// Load constant C1 (ln2 part 1) into Y5
	MOVQ $0x3fe62e42fefa39ef, AX
	MOVQ AX, X5
	VBROADCASTSD X5, Y5

	// Load constant C2 (ln2 part 2) into Y6
	MOVQ $0x3c7abc9e3b39803f, AX
	MOVQ AX, X6
	VBROADCASTSD X6, Y6

	// Load 1.0 into Y14
	MOVQ $0x3ff0000000000000, AX
	MOVQ AX, X14
	VBROADCASTSD X14, Y14

	// Load 0.5 into Y13
	MOVQ $0x3fe0000000000000, AX
	MOVQ AX, X13
	VBROADCASTSD X13, Y13

	// Load overflow threshold (709.782712893384) into Y8
	MOVQ $0x40862e42fefa39ef, AX
	MOVQ AX, X8
	VBROADCASTSD X8, Y8

	// Load underflow threshold (-745.1332191019412) into Y9
	MOVQ $0xc0874910d52d3052, AX
	MOVQ AX, X9
	VBROADCASTSD X9, Y9

	// Load +Inf into Y10
	MOVQ $0x7ff0000000000000, AX
	MOVQ AX, X10
	VBROADCASTSD X10, Y10

	// Load 0.0 into Y11
	VXORPD Y11, Y11, Y11

loop_exp_avx2:
	VMOVUPD (SI), Y0 // Y0 = x

	// Y1 = round(x * 1/ln2)
	VMULPD Y7, Y0, Y1
	VROUNDPD $0, Y1, Y1 // Y1 = n

	// Y4 = r = x - n*C1 - n*C2
	VMULPD Y5, Y1, Y2
	VSUBPD Y2, Y0, Y4
	VMULPD Y6, Y1, Y2
	VSUBPD Y2, Y4, Y4 // Y4 = r

	// Polynomial evaluation (Horner's method)
	// Y3 = c11
	MOVQ $0x3e5ae64567f544e4, AX
	MOVQ AX, X3
	VBROADCASTSD X3, Y3

	// Y3 = Y3 * r + c10
	VMULPD Y4, Y3, Y3
	MOVQ $0x3e927e4fb7789f5c, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + c9
	VMULPD Y4, Y3, Y3
	MOVQ $0x3ec71de3a556c733, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + c8
	VMULPD Y4, Y3, Y3
	MOVQ $0x3efa01a01a01a01a, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + c7
	VMULPD Y4, Y3, Y3
	MOVQ $0x3f2a01a01a01a01a, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + c6
	VMULPD Y4, Y3, Y3
	MOVQ $0x3f56c16c16c16c17, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + c5
	VMULPD Y4, Y3, Y3
	MOVQ $0x3f81111111111111, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + c4
	VMULPD Y4, Y3, Y3
	MOVQ $0x3fa5555555555555, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + c3
	VMULPD Y4, Y3, Y3
	MOVQ $0x3fc5555555555555, AX
	MOVQ AX, X15
	VBROADCASTSD X15, Y2
	VADDPD Y2, Y3, Y3

	// Y3 = Y3 * r + 0.5
	VMULPD Y4, Y3, Y3
	VADDPD Y13, Y3, Y3

	// Y3 = Y3 * r + 1.0
	VMULPD Y4, Y3, Y3
	VADDPD Y14, Y3, Y3

	// Y3 = Y3 * r + 1.0
	VMULPD Y4, Y3, Y3
	VADDPD Y14, Y3, Y3 // Y3 = P(r)

	// Now compute 2^n.
	// Convert Y1 (n) from float64 to int32 (X15)
	VCVTPD2DQY Y1, X15

	// Add 1023 to X15
	MOVOU add_1023_val<>(SB), X2
	VPADDD X2, X15, X15
	VPMOVZXDQ X15, Y15
	VPSLLQ $52, Y15, Y15 // Y15 now has 2^n in YMM!
	VMULPD Y15, Y3, Y3 // Multiply P(r) by 2^n!

	// Compare x >= 709.782712893384
	VCMPPD $14, Y8, Y0, Y1
	// Compare x < -745.1332191019412
	VCMPPD $1, Y9, Y0, Y2

	// Blend with 0.0 if underflow
	VBLENDVPD Y2, Y11, Y3, Y3

	// Blend with +Inf if overflow
	VBLENDVPD Y1, Y10, Y3, Y3

	// Store result
	VMOVUPD Y3, (DX)

	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_exp_avx2

done_exp_avx2:
	VZEROUPPER
	RET

// func logAVX2(a, result []float64)
TEXT ·logAVX2(SB), NOSPLIT, $0-48
	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX // n / 4
	JZ done_log_avx2

	// Load 2.0 float64 into Y11
	MOVQ $0x4000000000000000, AX
	MOVQ AX, X11
	VBROADCASTSD X11, Y11

	// Load 1.0 float64 into Y12
	MOVQ $0x3ff0000000000000, AX
	MOVQ AX, X12
	VBROADCASTSD X12, Y12

loop_log_avx2:
	VMOVUPD (SI), Y0 // Y0 = x

	// Extract e_bits
	VPSRLQ $52, Y0, Y1
	MOVQ $0x7ff, AX
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VPAND Y8, Y1, Y1

	// Convert e_bits to float64 (e)
	MOVQ $0x4330000000000000, AX
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VPOR Y8, Y1, Y1
	VSUBPD Y8, Y1, Y1
	MOVQ $0x408ff00000000000, AX // 1022.0
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VSUBPD Y8, Y1, Y1 // Y1 = e

	// Extract m
	MOVQ $0x000fffffffffffff, AX // mantissa mask
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VPAND Y8, Y0, Y2
	MOVQ $0x3fe0000000000000, AX // 1022 exponent
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VPOR Y8, Y2, Y2 // Y2 = m

	// Compare m < sqrt(0.5)
	MOVQ $0x3fe6a09e667f3bcd, AX // sqrt(0.5)
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VCMPPD $1, Y8, Y2, Y3 // Y3 = mask

	// Adjust m and e
	// Y8 = mask ? 2.0 : 1.0
	VBLENDVPD Y3, Y11, Y12, Y8
	VMULPD Y8, Y2, Y2 // Y2 = m adjusted
	// Y8 = mask ? 1.0 : 0.0
	VXORPD Y9, Y9, Y9 // Y9 = 0.0
	VBLENDVPD Y3, Y12, Y9, Y8
	VSUBPD Y8, Y1, Y1 // Y1 = e adjusted

	// f = m - 1.0
	VSUBPD Y12, Y2, Y2 // Y2 = f

	// s = f / (2.0 + f)
	VADDPD Y11, Y2, Y3 // Y3 = 2.0 + f
	VDIVPD Y3, Y2, Y3 // Y3 = s

	// s2 = s * s
	VMULPD Y3, Y3, Y4 // Y4 = s2
	// s4 = s2 * s2
	VMULPD Y4, Y4, Y5 // Y5 = s4

	// t1 = s2 * (L1 + s4 * (L3 + s4 * (L5 + s4 * L7)))
	// Compute L5 + s4 * L7 first:
	MOVQ $0x3fc2f112df3e5244, AX // L7
	MOVQ AX, X8
	VBROADCASTSD X8, Y6 // Y6 = L7
	VMULPD Y5, Y6, Y6
	MOVQ $0x3fc7466496cb03de, AX // L5
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VADDPD Y8, Y6, Y6 // Y6 = L5 + s4 * L7

	// Multiply by s4 and add L3:
	VMULPD Y5, Y6, Y6
	MOVQ $0x3fd2492494229359, AX // L3
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VADDPD Y8, Y6, Y6 // Y6 = L3 + s4 * (L5 + s4 * L7)

	// Multiply by s4 and add L1:
	VMULPD Y5, Y6, Y6
	MOVQ $0x3fe5555555555593, AX // L1
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VADDPD Y8, Y6, Y6 // Y6 = L1 + s4 * (L3 + s4 * (L5 + s4 * L7))

	// Multiply by s2:
	VMULPD Y4, Y6, Y6 // Y6 = t1

	// t2 = s4 * (L2 + s4 * (L4 + s4 * L6))
	// Compute L4 + s4 * L6 first:
	MOVQ $0x3fc39a09d078c69f, AX // L6
	MOVQ AX, X8
	VBROADCASTSD X8, Y7 // Y7 = L6
	VMULPD Y5, Y7, Y7
	MOVQ $0x3fcc71c51d8e78af, AX // L4
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VADDPD Y8, Y7, Y7 // Y7 = L4 + s4 * L6

	// Multiply by s4 and add L2:
	VMULPD Y5, Y7, Y7
	MOVQ $0x3fd999999997fa04, AX // L2
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VADDPD Y8, Y7, Y7 // Y7 = L2 + s4 * (L4 + s4 * L6)

	// Multiply by s4:
	VMULPD Y5, Y7, Y7 // Y7 = t2

	// R = t1 + t2
	VADDPD Y7, Y6, Y6 // Y6 = R

	// hfsq = 0.5 * f * f
	MOVQ $0x3fe0000000000000, AX // 0.5
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VMULPD Y8, Y2, Y7 // Y7 = 0.5 * f
	VMULPD Y2, Y7, Y7 // Y7 = hfsq

	// s * (hfsq + R)
	VADDPD Y6, Y7, Y6 // Y6 = hfsq + R
	VMULPD Y3, Y6, Y6 // Y6 = s * (hfsq + R)

	// k * Ln2Lo
	MOVQ $0x3dea39ef35793c76, AX // Ln2Lo
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VMULPD Y1, Y8, Y8 // Y8 = k * Ln2Lo

	// s * (hfsq + R) + k * Ln2Lo
	VADDPD Y8, Y6, Y6

	// hfsq - (s * (hfsq + R) + k * Ln2Lo)
	VSUBPD Y6, Y7, Y7

	// (hfsq - (s * (hfsq + R) + k * Ln2Lo)) - f
	VSUBPD Y2, Y7, Y7

	// k * Ln2Hi
	MOVQ $0x3fe62e42fee00000, AX // Ln2Hi
	MOVQ AX, X8
	VBROADCASTSD X8, Y8
	VMULPD Y1, Y8, Y8 // Y8 = k * Ln2Hi

	// result = k * Ln2Hi - ((hfsq - (s * (hfsq + R) + k * Ln2Lo)) - f)
	VSUBPD Y7, Y8, Y0 // Y0 = result

	// Store result
	VMOVUPD Y0, (DX)

	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_log_avx2

done_log_avx2:
	VZEROUPPER
	RET

// func sinAVX2(a, result []float64)
TEXT ·sinAVX2(SB), NOSPLIT, $0-48
	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX // n / 4
	JZ done_sin_avx2

	// Load 1/pi (0.3183098861837907) into Y5
	MOVQ $0x3fd45f306dc9c883, AX
	MOVQ AX, X5
	VBROADCASTSD X5, Y5

	// Load P1 (3.141592653589793) into Y6
	MOVQ $0x400921fb54442d18, AX
	MOVQ AX, X6
	VBROADCASTSD X6, Y6

	// Load P2 (1.2246467991473532e-16) into Y7
	MOVQ $0x3ca1a62633145c07, AX
	MOVQ AX, X7
	VBROADCASTSD X7, Y7

	// Load 1.0 into Y8
	MOVQ $0x3ff0000000000000, AX
	MOVQ AX, X8
	VBROADCASTSD X8, Y8

loop_sin_avx2:
	VMOVUPD (SI), Y0 // Y0 = x

	// Y1 = round(x / pi)
	VMULPD Y5, Y0, Y1
	VROUNDPD $0, Y1, Y1 // Y1 = n

	// Y2 = r = x - n*P1 - n*P2
	VMULPD Y6, Y1, Y2
	VSUBPD Y2, Y0, Y2
	VMULPD Y7, Y1, Y3
	VSUBPD Y3, Y2, Y2 // Y2 = r

	// r2 = r * r
	VMULPD Y2, Y2, Y3 // Y3 = r2

	// Polynomial for sin(r)/r
	// Y4 = c15
	MOVQ $0xbd6ae7f3e733b81f, AX // c15 = -7.647163731819816e-13
	MOVQ AX, X4
	VBROADCASTSD X4, Y4

	// Y4 = Y4 * r2 + c13
	VMULPD Y3, Y4, Y4
	MOVQ $0x3de6124613a86d0a, AX // c13 = 1.6059043836821616e-10
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c11
	VMULPD Y3, Y4, Y4
	MOVQ $0xbe5ae64567f544e4, AX // c11 = -2.505210838544172e-8
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c9
	VMULPD Y3, Y4, Y4
	MOVQ $0x3ec71de3a556c733, AX // c9 = 2.755731922398589e-6
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c7
	VMULPD Y3, Y4, Y4
	MOVQ $0xbf2a01a01a01a01a, AX // c7 = -0.0001984126984126984
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c5
	VMULPD Y3, Y4, Y4
	MOVQ $0x3f81111111111111, AX // c5 = 0.008333333333333333
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c3
	VMULPD Y3, Y4, Y4
	MOVQ $0xbfc5555555555555, AX // c3 = -0.16666666666666666
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + 1.0
	VMULPD Y3, Y4, Y4
	VADDPD Y8, Y4, Y4

	// Y4 = r * Y4
	VMULPD Y2, Y4, Y4

	// Negate Y4 if n is odd.
	VCVTPD2DQY Y1, X15
	MOVQ $1, AX
	MOVQ AX, X0
	VPBROADCASTD X0, X0
	VPAND X0, X15, X15
	VPMOVZXDQ X15, Y15
	VPSLLQ $63, Y15, Y15
	VXORPD Y15, Y4, Y4

	// Store result
	VMOVUPD Y4, (DX)

	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_sin_avx2

done_sin_avx2:
	VZEROUPPER
	RET

// func cosAVX2(a, result []float64)
TEXT ·cosAVX2(SB), NOSPLIT, $0-48
	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX
	JZ done_cos_avx2

	// Load 1/pi (0.3183098861837907) into Y5
	MOVQ $0x3fd45f306dc9c883, AX
	MOVQ AX, X5
	VBROADCASTSD X5, Y5

	// Load P1 (3.141592653589793) into Y6
	MOVQ $0x400921fb54442d18, AX
	MOVQ AX, X6
	VBROADCASTSD X6, Y6

	// Load P2 (1.2246467991473532e-16) into Y7
	MOVQ $0x3ca1a62633145c07, AX
	MOVQ AX, X7
	VBROADCASTSD X7, Y7

	// Load 1.0 into Y8
	MOVQ $0x3ff0000000000000, AX
	MOVQ AX, X8
	VBROADCASTSD X8, Y8

loop_cos_avx2:
	VMOVUPD (SI), Y0 // Y0 = x

	// Y1 = round(x / pi)
	VMULPD Y5, Y0, Y1
	VROUNDPD $0, Y1, Y1 // Y1 = n

	// Y2 = r = x - n*P1 - n*P2
	VMULPD Y6, Y1, Y2
	VSUBPD Y2, Y0, Y2
	VMULPD Y7, Y1, Y3
	VSUBPD Y3, Y2, Y2 // Y2 = r

	// r2 = r * r
	VMULPD Y2, Y2, Y3 // Y3 = r2

	// Polynomial for cos(r)
	// Y4 = c14
	MOVQ $0xbda93974a8c07c9d, AX // c14 = -1.1470745597729725e-11
	MOVQ AX, X4
	VBROADCASTSD X4, Y4

	// Y4 = Y4 * r2 + c12
	VMULPD Y3, Y4, Y4
	MOVQ $0x3e21eed8eff8d898, AX // c12 = 2.08767569878681e-9
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c10
	VMULPD Y3, Y4, Y4
	MOVQ $0xbe927e4fb7789f5c, AX // c10 = -2.755731922398589e-7
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c8
	VMULPD Y3, Y4, Y4
	MOVQ $0x3efa01a01a01a01a, AX // c8 = 2.48015873015873e-5
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c6
	VMULPD Y3, Y4, Y4
	MOVQ $0xbf56c16c16c16c17, AX // c6 = -0.001388888888888889
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c4
	VMULPD Y3, Y4, Y4
	MOVQ $0x3fa5555555555555, AX // c4 = 0.041666666666666664
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c2
	VMULPD Y3, Y4, Y4
	MOVQ $0xbfe0000000000000, AX // c2 = -0.5
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + 1.0
	VMULPD Y3, Y4, Y4
	VADDPD Y8, Y4, Y4 // Y4 = cos(r)

	// Negate Y4 if n is odd.
	VCVTPD2DQY Y1, X15
	MOVQ $1, AX
	MOVQ AX, X0
	VPBROADCASTD X0, X0
	VPAND X0, X15, X15
	VPMOVZXDQ X15, Y15
	VPSLLQ $63, Y15, Y15
	VXORPD Y15, Y4, Y4

	// Store result
	VMOVUPD Y4, (DX)

	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_cos_avx2

done_cos_avx2:
	VZEROUPPER
	RET

// func tanAVX2(a, result []float64)
TEXT ·tanAVX2(SB), NOSPLIT, $0-48
	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX
	JZ done_tan_avx2

	// Load 1/pi (0.3183098861837907) into Y5
	MOVQ $0x3fd45f306dc9c883, AX
	MOVQ AX, X5
	VBROADCASTSD X5, Y5

	// Load P1 (3.141592653589793) into Y6
	MOVQ $0x400921fb54442d18, AX
	MOVQ AX, X6
	VBROADCASTSD X6, Y6

	// Load P2 (1.2246467991473532e-16) into Y7
	MOVQ $0x3ca1a62633145c07, AX
	MOVQ AX, X7
	VBROADCASTSD X7, Y7

	// Load 1.0 into Y8
	MOVQ $0x3ff0000000000000, AX
	MOVQ AX, X8
	VBROADCASTSD X8, Y8

loop_tan_avx2:
	VMOVUPD (SI), Y0 // Y0 = x

	// Y1 = round(x / pi)
	VMULPD Y5, Y0, Y1
	VROUNDPD $0, Y1, Y1 // Y1 = n

	// Y2 = r = x - n*P1 - n*P2
	VMULPD Y6, Y1, Y2
	VSUBPD Y2, Y0, Y2
	VMULPD Y7, Y1, Y3
	VSUBPD Y3, Y2, Y2 // Y2 = r

	// r2 = r * r
	VMULPD Y2, Y2, Y3 // Y3 = r2

	// We'll compute sin(r)/r in Y4, and cos(r) in Y11.
	
	// --- Compute sin(r)/r in Y4 ---
	MOVQ $0xbd6ae7f3e733b81f, AX // c15 = -7.647163731819816e-13
	MOVQ AX, X4
	VBROADCASTSD X4, Y4

	// Y4 = Y4 * r2 + c13
	VMULPD Y3, Y4, Y4
	MOVQ $0x3de6124613a86d0a, AX // c13 = 1.6059043836821616e-10
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c11
	VMULPD Y3, Y4, Y4
	MOVQ $0xbe5ae64567f544e4, AX // c11 = -2.505210838544172e-8
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c9
	VMULPD Y3, Y4, Y4
	MOVQ $0x3ec71de3a556c733, AX // c9 = 2.755731922398589e-6
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c7
	VMULPD Y3, Y4, Y4
	MOVQ $0xbf2a01a01a01a01a, AX // c7 = -0.0001984126984126984
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c5
	VMULPD Y3, Y4, Y4
	MOVQ $0x3f81111111111111, AX // c5 = 0.008333333333333333
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + c3
	VMULPD Y3, Y4, Y4
	MOVQ $0xbfc5555555555555, AX // c3 = -0.16666666666666666
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y4, Y4

	// Y4 = Y4 * r2 + 1.0
	VMULPD Y3, Y4, Y4
	VADDPD Y8, Y4, Y4

	// Y4 = r * Y4 (this is sin(r))
	VMULPD Y2, Y4, Y4

	// --- Compute cos(r) in Y11 ---
	MOVQ $0xbda93974a8c07c9d, AX // c14 = -1.1470745597729725e-11
	MOVQ AX, X11
	VBROADCASTSD X11, Y11

	// Y11 = Y11 * r2 + c12
	VMULPD Y3, Y11, Y11
	MOVQ $0x3e21eed8eff8d898, AX // c12 = 2.08767569878681e-9
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y11, Y11

	// Y11 = Y11 * r2 + c10
	VMULPD Y3, Y11, Y11
	MOVQ $0xbe927e4fb7789f5c, AX // c10 = -2.755731922398589e-7
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y11, Y11

	// Y11 = Y11 * r2 + c8
	VMULPD Y3, Y11, Y11
	MOVQ $0x3efa01a01a01a01a, AX // c8 = 2.48015873015873e-5
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y11, Y11

	// Y11 = Y11 * r2 + c6
	VMULPD Y3, Y11, Y11
	MOVQ $0xbf56c16c16c16c17, AX // c6 = -0.001388888888888889
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y11, Y11

	// Y11 = Y11 * r2 + c4
	VMULPD Y3, Y11, Y11
	MOVQ $0x3fa5555555555555, AX // c4 = 0.041666666666666664
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y11, Y11

	// Y11 = Y11 * r2 + c2
	VMULPD Y3, Y11, Y11
	MOVQ $0xbfe0000000000000, AX // c2 = -0.5
	MOVQ AX, X15
	VBROADCASTSD X15, Y0
	VADDPD Y0, Y11, Y11

	// Y11 = Y11 * r2 + 1.0
	VMULPD Y3, Y11, Y11
	VADDPD Y8, Y11, Y11 // Y11 = cos(r)

	// Y4 = sin(r) / cos(r)
	VDIVPD Y11, Y4, Y4

	// Store result
	VMOVUPD Y4, (DX)

	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_tan_avx2

done_tan_avx2:
	VZEROUPPER
	RET

