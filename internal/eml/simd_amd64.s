//go:build amd64 && !purego

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

// func xgetbv(cx uint32) (eax, edx uint32)
TEXT ·xgetbv(SB), NOSPLIT, $0-16
	MOVL cx+0(FP), CX
	BYTE $0x0F; BYTE $0x01; BYTE $0xD0 // XGETBV
	MOVL AX, eax+8(FP)
	MOVL DX, edx+12(FP)
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

    MOVQ CX, AX
    SHRQ $4, CX // n / 16
    JZ tail_add512

loop_add512_16:
    VMOVUPD (SI), Z0
    VMOVUPD 64(SI), Z1
    VMOVUPD (DI), Z2
    VMOVUPD 64(DI), Z3
    VADDPD Z2, Z0, Z0
    VADDPD Z3, Z1, Z1
    VMOVUPD Z0, (DX)
    VMOVUPD Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_add512_16

tail_add512:
    TESTQ $8, AX
    JZ done_add512
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VADDPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)

done_add512:
    VZEROUPPER
    RET

// func subAVX512(a, b, result []float64)
TEXT ·subAVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $4, CX
    JZ tail_sub512

loop_sub512_16:
    VMOVUPD (SI), Z0
    VMOVUPD 64(SI), Z1
    VMOVUPD (DI), Z2
    VMOVUPD 64(DI), Z3
    VSUBPD Z2, Z0, Z0
    VSUBPD Z3, Z1, Z1
    VMOVUPD Z0, (DX)
    VMOVUPD Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_sub512_16

tail_sub512:
    TESTQ $8, AX
    JZ done_sub512
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VSUBPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)

done_sub512:
    VZEROUPPER
    RET

// func mulAVX512(a, b, result []float64)
TEXT ·mulAVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $4, CX
    JZ tail_mul512

loop_mul512_16:
    VMOVUPD (SI), Z0
    VMOVUPD 64(SI), Z1
    VMOVUPD (DI), Z2
    VMOVUPD 64(DI), Z3
    VMULPD Z2, Z0, Z0
    VMULPD Z3, Z1, Z1
    VMOVUPD Z0, (DX)
    VMOVUPD Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_mul512_16

tail_mul512:
    TESTQ $8, AX
    JZ done_mul512
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VMULPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)

done_mul512:
    VZEROUPPER
    RET

// func divAVX512(a, b, result []float64)
TEXT ·divAVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $4, CX
    JZ tail_div512

loop_div512_16:
    VMOVUPD (SI), Z0
    VMOVUPD 64(SI), Z1
    VMOVUPD (DI), Z2
    VMOVUPD 64(DI), Z3
    VDIVPD Z2, Z0, Z0
    VDIVPD Z3, Z1, Z1
    VMOVUPD Z0, (DX)
    VMOVUPD Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_div512_16

tail_div512:
    TESTQ $8, AX
    JZ done_div512
    VMOVUPD (SI), Z0
    VMOVUPD (DI), Z1
    VDIVPD Z1, Z0, Z2
    VMOVUPD Z2, (DX)

done_div512:
    VZEROUPPER
    RET

// func addScalarAVX512(a []float64, b float64, result []float64)
TEXT ·addScalarAVX512(SB), NOSPLIT, $0-56
    MOVQ a_len+8(FP), CX
    MOVQ CX, AX
    TESTQ CX, CX
    JZ done_add_scalar512

    MOVQ a_base+0(FP), SI
    MOVSD b+24(FP), X0
    MOVQ result_base+32(FP), DX
    VBROADCASTSD X0, Z0

    SHRQ $4, CX
    JZ tail_add_scalar512

loop_add_scalar512_16:
    VMOVUPD (SI), Z1
    VMOVUPD 64(SI), Z2
    VADDPD Z0, Z1, Z3
    VADDPD Z0, Z2, Z4
    VMOVUPD Z3, (DX)
    VMOVUPD Z4, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_add_scalar512_16

tail_add_scalar512:
    TESTQ $8, AX
    JZ done_add_scalar512
    VMOVUPD (SI), Z1
    VADDPD Z0, Z1, Z2
    VMOVUPD Z2, (DX)

done_add_scalar512:
    VZEROUPPER
    RET

// func mulScalarAVX512(a []float64, b float64, result []float64)
TEXT ·mulScalarAVX512(SB), NOSPLIT, $0-56
    MOVQ a_len+8(FP), CX
    MOVQ CX, AX
    TESTQ CX, CX
    JZ done_mul_scalar512

    MOVQ a_base+0(FP), SI
    MOVSD b+24(FP), X0
    MOVQ result_base+32(FP), DX
    VBROADCASTSD X0, Z0

    SHRQ $4, CX
    JZ tail_mul_scalar512

loop_mul_scalar512_16:
    VMOVUPD (SI), Z1
    VMOVUPD 64(SI), Z2
    VMULPD Z0, Z1, Z3
    VMULPD Z0, Z2, Z4
    VMOVUPD Z3, (DX)
    VMOVUPD Z4, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_mul_scalar512_16

tail_mul_scalar512:
    TESTQ $8, AX
    JZ done_mul_scalar512
    VMOVUPD (SI), Z1
    VMULPD Z0, Z1, Z2
    VMOVUPD Z2, (DX)

done_mul_scalar512:
    VZEROUPPER
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

    MOVQ CX, AX
    SHRQ $4, CX
    JZ tail_sqrt512

loop_sqrt512_16:
    VSQRTPD (SI), Z0
    VSQRTPD 64(SI), Z1
    VMOVUPD Z0, (DX)
    VMOVUPD Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_sqrt512_16

tail_sqrt512:
    TESTQ $8, AX
    JZ done_sqrt512
    VSQRTPD (SI), Z0
    VMOVUPD Z0, (DX)

done_sqrt512:
    VZEROUPPER
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

	MOVQ CX, AX
	SHRQ $4, CX // n / 16
	JZ tail_fma512

loop_fma512_16:
	VMOVUPD (SI), Z0
	VMOVUPD 64(SI), Z1
	VMOVUPD (DI), Z2
	VMOVUPD 64(DI), Z3
	VMOVUPD (BX), Z4
	VMOVUPD 64(BX), Z5
	VFMADD213PD Z4, Z2, Z0
	VFMADD213PD Z5, Z3, Z1
	VMOVUPD Z0, (DX)
	VMOVUPD Z1, 64(DX)

	ADDQ $128, SI
	ADDQ $128, DI
	ADDQ $128, BX
	ADDQ $128, DX
	DECQ CX
	JNZ loop_fma512_16

tail_fma512:
	TESTQ $8, AX
	JZ done_fma512
	VMOVUPD (SI), Z0
	VMOVUPD (DI), Z1
	VMOVUPD (BX), Z2
	VFMADD213PD Z2, Z1, Z0
	VMOVUPD Z0, (DX)

done_fma512:
	VZEROUPPER
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
    MOVQ a_len+8(FP), CX
    TESTQ CX, CX
    JZ done_abs512

    MOVQ $0x7fffffffffffffff, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Z15

    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX

    MOVQ CX, AX
    SHRQ $4, CX
    JZ tail_abs512

loop_abs512_16:
    VMOVUPD (SI), Z0
    VMOVUPD 64(SI), Z1
    VPANDQ Z15, Z0, Z0
    VPANDQ Z15, Z1, Z1
    VMOVUPD Z0, (DX)
    VMOVUPD Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_abs512_16

tail_abs512:
    TESTQ $8, AX
    JZ done_abs512
    VMOVUPD (SI), Z0
    VPANDQ Z15, Z0, Z0
    VMOVUPD Z0, (DX)

done_abs512:
    VZEROUPPER
    RET

// func negAVX512(a, result []float64)
TEXT ·negAVX512(SB), NOSPLIT, $0-48
    MOVQ a_len+8(FP), CX
    TESTQ CX, CX
    JZ done_neg512

    MOVQ $0x8000000000000000, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Z15

    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX

    MOVQ CX, AX
    SHRQ $4, CX
    JZ tail_neg512

loop_neg512_16:
    VMOVUPD (SI), Z0
    VMOVUPD 64(SI), Z1
    VPXORQ Z15, Z0, Z0
    VPXORQ Z15, Z1, Z1
    VMOVUPD Z0, (DX)
    VMOVUPD Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_neg512_16

tail_neg512:
    TESTQ $8, AX
    JZ done_neg512
    VMOVUPD (SI), Z0
    VPXORQ Z15, Z0, Z0
    VMOVUPD Z0, (DX)

done_neg512:
    VZEROUPPER
    RET

// func invAVX512(a, result []float64)
TEXT ·invAVX512(SB), NOSPLIT, $0-48
    MOVQ a_len+8(FP), CX
    TESTQ CX, CX
    JZ done_inv512

    MOVQ $0x3ff0000000000000, AX
    MOVQ AX, X15
    VBROADCASTSD X15, Z15

    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX

    MOVQ CX, AX
    SHRQ $4, CX
    JZ tail_inv512


loop_inv512_16:
    VMOVUPD (SI), Z0
    VMOVUPD 64(SI), Z1
    VDIVPD Z0, Z15, Z2
    VDIVPD Z1, Z15, Z3
    VMOVUPD Z2, (DX)
    VMOVUPD Z3, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_inv512_16

tail_inv512:
    TESTQ $8, AX
    JZ done_inv512
    VMOVUPD (SI), Z0
    VDIVPD Z0, Z15, Z1
    VMOVUPD Z1, (DX)

done_inv512:
    VZEROUPPER
    RET


DATA add_1023_val<>+0(SB)/4, $1023
DATA add_1023_val<>+4(SB)/4, $1023
DATA add_1023_val<>+8(SB)/4, $1023
DATA add_1023_val<>+12(SB)/4, $1023
GLOBL add_1023_val<>(SB), RODATA|NOPTR, $16

DATA add_1022_val<>+0(SB)/4, $1022
DATA add_1022_val<>+4(SB)/4, $1022
DATA add_1022_val<>+8(SB)/4, $1022
DATA add_1022_val<>+12(SB)/4, $1022
GLOBL add_1022_val<>(SB), RODATA|NOPTR, $16

DATA exp_consts<>+0(SB)/8, $0x3ff71547652b82fe  // 0: LOG2E
DATA exp_consts<>+8(SB)/8, $0x3fe62e42fefa3800  // 8: LN2U
DATA exp_consts<>+16(SB)/8, $0x3d2ef35793c76730 // 16: LN2L
DATA exp_consts<>+24(SB)/8, $0x3fb0000000000000 // 24: 0.0625
DATA exp_consts<>+32(SB)/8, $0x3efa01a01a01a01a // 32: c8 = 1/8!
DATA exp_consts<>+40(SB)/8, $0x3f2a01a01a01a01a // 40: c7 = 1/7!
DATA exp_consts<>+48(SB)/8, $0x3f56c16c16c16c17 // 48: c6 = 1/6!
DATA exp_consts<>+56(SB)/8, $0x3f81111111111111 // 56: c5 = 1/5!
DATA exp_consts<>+64(SB)/8, $0x3fa5555555555555 // 64: c4 = 1/4!
DATA exp_consts<>+72(SB)/8, $0x3fc5555555555555 // 72: c3 = 1/3!
DATA exp_consts<>+80(SB)/8, $0x3fe0000000000000 // 80: 0.5
DATA exp_consts<>+88(SB)/8, $0x3ff0000000000000 // 88: 1.0
DATA exp_consts<>+96(SB)/8, $0x4000000000000000 // 96: 2.0
DATA exp_consts<>+104(SB)/8, $0x40862e42fefa39ef // 104: overflow (709.782712893384)
DATA exp_consts<>+112(SB)/8, $0xc0874910d52d3052 // 112: underflow (-745.1332191019412)
DATA exp_consts<>+120(SB)/8, $0x7ff0000000000000 // 120: +Inf
DATA exp_consts<>+128(SB)/8, $0xc090cc0000000000 // 128: k_min (-1075.0)
DATA exp_consts<>+136(SB)/8, $0x4090000000000000 // 136: k_max (1024.0)
GLOBL exp_consts<>(SB), RODATA|NOPTR, $144

DATA log_consts<>+0(SB)/8, $0x4000000000000000  // 0: 2.0
DATA log_consts<>+8(SB)/8, $0x3ff0000000000000  // 8: 1.0
DATA log_consts<>+16(SB)/8, $0x3fe0000000000000 // 16: 0.5
DATA log_consts<>+24(SB)/8, $0x3fe6a09e667f3bcd // 24: sqrt(0.5)
DATA log_consts<>+32(SB)/8, $0x3fe62e42fee00000 // 32: Ln2Hi
DATA log_consts<>+40(SB)/8, $0x3dea39ef35793c76 // 40: Ln2Lo
DATA log_consts<>+48(SB)/8, $0x3fe5555555555593 // 48: L1
DATA log_consts<>+56(SB)/8, $0x3fd999999997fa04 // 56: L2
DATA log_consts<>+64(SB)/8, $0x3fd2492494229359 // 64: L3
DATA log_consts<>+72(SB)/8, $0x3fcc71c51d8e78af // 72: L4
DATA log_consts<>+80(SB)/8, $0x3fc7466496cb03de // 80: L5
DATA log_consts<>+88(SB)/8, $0x3fc39a09d078c69f // 88: L6
DATA log_consts<>+96(SB)/8, $0x3fc2f112df3e5244 // 96: L7
DATA log_consts<>+104(SB)/8, $0x000fffffffffffff // 104: mantissa mask
DATA log_consts<>+112(SB)/8, $0x408ff00000000000 // 112: 1022.0
DATA log_consts<>+120(SB)/8, $0x00000000000007ff // 120: 0x7ff exponent mask
DATA log_consts<>+128(SB)/8, $0x4330000000000000 // 128: uint52-to-float64 magic
DATA log_consts<>+136(SB)/8, $0x0010000000000000 // 136: min normal 2^-1022
DATA log_consts<>+144(SB)/8, $0x4350000000000000 // 144: 2^54 subnormal scale
DATA log_consts<>+152(SB)/8, $0x4042b708872320e2 // 152: 54 * Ln2
DATA log_consts<>+160(SB)/8, $0xfff0000000000000 // 160: -Inf
DATA log_consts<>+168(SB)/8, $0x7ff8000000000001 // 168: NaN
DATA log_consts<>+176(SB)/8, $0x7ff0000000000000 // 176: +Inf
GLOBL log_consts<>(SB), RODATA|NOPTR, $184

DATA trig_consts<>+0(SB)/8, $0x3fe45f306dc9c883  // 0: 2/pi
DATA trig_consts<>+8(SB)/8, $0x3ff921fb40000000  // 8: PI2A
DATA trig_consts<>+16(SB)/8, $0x3e74442d00000000 // 16: PI2B
DATA trig_consts<>+24(SB)/8, $0x3cf8469898cc5170 // 24: PI2C
DATA trig_consts<>+32(SB)/8, $0x3de5d8fd1fd19ccd // 32: s0
DATA trig_consts<>+40(SB)/8, $0xbe5ae5e5a9291f5d // 40: s1
DATA trig_consts<>+48(SB)/8, $0x3ec71de3567d48a1 // 48: s2
DATA trig_consts<>+56(SB)/8, $0xbf2a01a019bfdf03 // 56: s3
DATA trig_consts<>+64(SB)/8, $0x3f8111111110f7d0 // 64: s4
DATA trig_consts<>+72(SB)/8, $0xbfc5555555555548 // 72: s5
DATA trig_consts<>+80(SB)/8, $0xbda8fa49a0861a9b // 80: c0
DATA trig_consts<>+88(SB)/8, $0x3e21ee9d7b4e3f05 // 88: c1
DATA trig_consts<>+96(SB)/8, $0xbe927e4f7eac4bc6 // 96: c2
DATA trig_consts<>+104(SB)/8, $0x3efa01a019c844f5 // 104: c3
DATA trig_consts<>+112(SB)/8, $0xbf56c16c16c14f91 // 112: c4
DATA trig_consts<>+120(SB)/8, $0x3fa555555555554b // 120: c5
DATA trig_consts<>+128(SB)/8, $0x41d0000000000000 // 128: 2^30 reduce threshold
DATA trig_consts<>+136(SB)/8, $0x7fffffffffffffff // 136: abs mask
DATA trig_consts<>+144(SB)/8, $0x8000000000000000 // 144: sign mask
GLOBL trig_consts<>(SB), RODATA|NOPTR, $152

// func expAVX2(a, result []float64)
TEXT ·expAVX2(SB), NOSPLIT, $0-48
	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $2, CX // n / 4
	JZ done_exp_avx2

	VBROADCASTSD exp_consts<>+0(SB), Y5  // LOG2E
	VBROADCASTSD exp_consts<>+8(SB), Y6  // LN2U
	VBROADCASTSD exp_consts<>+16(SB), Y7 // LN2L
	VBROADCASTSD exp_consts<>+24(SB), Y8 // 0.0625
	VBROADCASTSD exp_consts<>+32(SB), Y9 // c8
	VBROADCASTSD exp_consts<>+40(SB), Y10 // c7
	VBROADCASTSD exp_consts<>+48(SB), Y11 // c6
	VBROADCASTSD exp_consts<>+56(SB), Y12 // c5
	VBROADCASTSD exp_consts<>+64(SB), Y13 // c4
	VBROADCASTSD exp_consts<>+72(SB), Y14 // c3
	VBROADCASTSD exp_consts<>+80(SB), Y15 // 0.5

loop_exp_avx2:
	VMOVUPD (SI), Y0 // x
	VMOVUPD Y0, Y4   // save original x

	VMULPD Y5, Y0, Y1
	VROUNDPD $0, Y1, Y1 // k = roundToEven(x * LOG2E)

	VFNMADD231PD Y6, Y1, Y0 // x - k * LN2U
	VFNMADD231PD Y7, Y1, Y0 // - k * LN2L
	VMULPD Y8, Y0, Y0       // r = (x - k*LN2) * 0.0625

	// Taylor polynomial on r = (x - k*LN2)/16
	VMOVUPD Y9, Y2
	VFMADD213PD Y10, Y0, Y2
	VFMADD213PD Y11, Y0, Y2
	VFMADD213PD Y12, Y0, Y2
	VFMADD213PD Y13, Y0, Y2
	VFMADD213PD Y14, Y0, Y2
	VFMADD213PD Y15, Y0, Y2
	VBROADCASTSD exp_consts<>+88(SB), Y3 // 1.0
	VFMADD213PD Y3, Y0, Y2
	VMULPD Y2, Y0, Y0

	// 4 squaring steps: exp(16*r)
	VBROADCASTSD exp_consts<>+96(SB), Y2 // 2.0
	VADDPD Y2, Y0, Y3
	VMULPD Y3, Y0, Y0
	VADDPD Y2, Y0, Y3
	VMULPD Y3, Y0, Y0
	VADDPD Y2, Y0, Y3
	VMULPD Y3, Y0, Y0
	VADDPD Y2, Y0, Y3
	VBROADCASTSD exp_consts<>+88(SB), Y2 // 1.0
	VFMADD213PD Y2, Y3, Y0

	// Clamp k to [-1075, 1024] before int32 conversion
	VBROADCASTSD exp_consts<>+128(SB), Y2
	VMAXPD Y2, Y1, Y1
	VBROADCASTSD exp_consts<>+136(SB), Y2
	VMINPD Y2, Y1, Y1

	VCVTPD2DQY Y1, X2
	MOVOU add_1023_val<>(SB), X3
	VPADDD X3, X2, X2 // BX = k + 1023

	// Handle denormal exponent (BX <= 0) identically to math/exp_amd64.s
	VPXOR X1, X1, X1
	VPCMPGTD X2, X1, X1 // X1 = (0 > BX) ? -1 : 0
	MOVOU add_1022_val<>(SB), X3
	VPAND X3, X1, X1    // X1 = (BX <= 0) ? 1022 : 0
	VPADDD X1, X2, X2   // BX' = (BX <= 0) ? BX + 1022 : BX
	VPMOVZXDQ X2, Y2
	VPSLLQ $52, Y2, Y2
	VMULPD Y2, Y0, Y0

	MOVOU add_1023_val<>(SB), X3
	VPSUBD X1, X3, X1   // (BX <= 0) ? 1 : 1023
	VPMOVZXDQ X1, Y1
	VPSLLQ $52, Y1, Y1
	VMULPD Y1, Y0, Y0

	// Underflow: x < -745.1332191019412 -> 0.0
	VBROADCASTSD exp_consts<>+112(SB), Y2
	VCMPPD $1, Y2, Y4, Y1
	VXORPD Y3, Y3, Y3
	VBLENDVPD Y1, Y3, Y0, Y0

	// Overflow: x > 709.782712893384 -> +Inf
	VBROADCASTSD exp_consts<>+104(SB), Y2
	VCMPPD $14, Y2, Y4, Y1
	VBROADCASTSD exp_consts<>+120(SB), Y3
	VBLENDVPD Y1, Y3, Y0, Y0

	// NaN: x != x -> NaN
	VCMPPD $3, Y4, Y4, Y1
	VBLENDVPD Y1, Y4, Y0, Y0

	VMOVUPD Y0, (DX)

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

	VBROADCASTSD log_consts<>+0(SB), Y11  // 2.0
	VBROADCASTSD log_consts<>+8(SB), Y12  // 1.0
	VBROADCASTSD log_consts<>+24(SB), Y14 // sqrt(0.5)
	VBROADCASTSD log_consts<>+104(SB), Y15 // mantissa mask

loop_log_avx2:
	VMOVUPD (SI), Y0 // x
	VMOVUPD Y0, Y10  // save original x

	// Subnormal pre-scaling: if 0 < x < 2^-1022, scale x by 2^54
	VXORPD Y9, Y9, Y9
	VCMPPD $14, Y9, Y0, Y1                  // x > 0
	VBROADCASTSD log_consts<>+136(SB), Y2   // 2^-1022
	VCMPPD $1, Y2, Y0, Y3                   // x < 2^-1022
	VANDPD Y3, Y1, Y13                      // Y13 = subnormal mask
	VBROADCASTSD log_consts<>+144(SB), Y2   // 2^54
	VBLENDVPD Y13, Y2, Y12, Y2
	VMULPD Y2, Y0, Y0

	// Extract e_bits
	VPSRLQ $52, Y0, Y1
	VBROADCASTSD log_consts<>+120(SB), Y8   // 0x7ff
	VPAND Y8, Y1, Y1

	// Convert e_bits to float64 (e)
	VBROADCASTSD log_consts<>+128(SB), Y8   // 2^52 magic
	VPOR Y8, Y1, Y1
	VSUBPD Y8, Y1, Y1
	VBROADCASTSD log_consts<>+112(SB), Y8   // 1022.0
	VSUBPD Y8, Y1, Y1                       // Y1 = e

	// Extract m
	VPAND Y15, Y0, Y2
	VBROADCASTSD log_consts<>+16(SB), Y8    // 0x3fe0000000000000 (0.5 / 1022 bias)
	VPOR Y8, Y2, Y2                         // Y2 = m

	// Compare m < sqrt(0.5)
	VCMPPD $1, Y14, Y2, Y3                  // Y3 = mask
	VBLENDVPD Y3, Y11, Y12, Y8              // mask ? 2.0 : 1.0
	VMULPD Y8, Y2, Y2                       // m adjusted
	VBLENDVPD Y3, Y12, Y9, Y8               // mask ? 1.0 : 0.0
	VSUBPD Y8, Y1, Y1                       // e adjusted

	// f = m - 1.0
	VSUBPD Y12, Y2, Y2                      // Y2 = f

	// s = f / (2.0 + f)
	VADDPD Y11, Y2, Y3
	VDIVPD Y3, Y2, Y3                       // Y3 = s

	// s2 = s * s, s4 = s2 * s2
	VMULPD Y3, Y3, Y4                       // Y4 = s2
	VMULPD Y4, Y4, Y5                       // Y5 = s4

	// t1 = s2 * (L1 + s4 * (L3 + s4 * (L5 + s4 * L7)))
	VBROADCASTSD log_consts<>+96(SB), Y6    // L7
	VMULPD Y5, Y6, Y6
	VBROADCASTSD log_consts<>+80(SB), Y8    // L5
	VADDPD Y8, Y6, Y6
	VMULPD Y5, Y6, Y6
	VBROADCASTSD log_consts<>+64(SB), Y8    // L3
	VADDPD Y8, Y6, Y6
	VMULPD Y5, Y6, Y6
	VBROADCASTSD log_consts<>+48(SB), Y8    // L1
	VADDPD Y8, Y6, Y6
	VMULPD Y4, Y6, Y6                       // Y6 = t1

	// t2 = s4 * (L2 + s4 * (L4 + s4 * L6))
	VBROADCASTSD log_consts<>+88(SB), Y7    // L6
	VMULPD Y5, Y7, Y7
	VBROADCASTSD log_consts<>+72(SB), Y8    // L4
	VADDPD Y8, Y7, Y7
	VMULPD Y5, Y7, Y7
	VBROADCASTSD log_consts<>+56(SB), Y8    // L2
	VADDPD Y8, Y7, Y7
	VMULPD Y5, Y7, Y7                       // Y7 = t2

	// R = t1 + t2
	VADDPD Y7, Y6, Y6                       // Y6 = R

	// hfsq = 0.5 * f * f
	VBROADCASTSD log_consts<>+16(SB), Y8    // 0.5
	VMULPD Y8, Y2, Y7
	VMULPD Y2, Y7, Y7                       // Y7 = hfsq

	// s * (hfsq + R)
	VADDPD Y6, Y7, Y6
	VMULPD Y3, Y6, Y6

	// k * Ln2Lo
	VBROADCASTSD log_consts<>+40(SB), Y8    // Ln2Lo
	VMULPD Y1, Y8, Y8
	VADDPD Y8, Y6, Y6
	VSUBPD Y6, Y7, Y7
	VSUBPD Y2, Y7, Y7

	// k * Ln2Hi
	VBROADCASTSD log_consts<>+32(SB), Y8    // Ln2Hi
	VMULPD Y1, Y8, Y8
	VSUBPD Y7, Y8, Y0                       // Y0 = result

	// Subtract 54*Ln2 for subnormal inputs
	VBROADCASTSD log_consts<>+152(SB), Y2   // 54 * Ln2
	VSUBPD Y2, Y0, Y2
	VBLENDVPD Y13, Y2, Y0, Y0

	// Special cases on original x (Y10)
	// x == 1.0 -> +0.0
	VCMPPD $0, Y12, Y10, Y1
	VBLENDVPD Y1, Y9, Y0, Y0
	// x == ±0.0 -> -Inf
	VCMPPD $0, Y9, Y10, Y1
	VBROADCASTSD log_consts<>+160(SB), Y2
	VBLENDVPD Y1, Y2, Y0, Y0
	// x < 0.0 -> NaN
	VCMPPD $1, Y9, Y10, Y1
	VBROADCASTSD log_consts<>+168(SB), Y2
	VBLENDVPD Y1, Y2, Y0, Y0
	// x == +Inf -> +Inf
	VBROADCASTSD log_consts<>+176(SB), Y2
	VCMPPD $0, Y2, Y10, Y1
	VBLENDVPD Y1, Y2, Y0, Y0
	// x is NaN -> NaN
	VCMPPD $3, Y10, Y10, Y1
	VBLENDVPD Y1, Y10, Y0, Y0

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

	VBROADCASTSD trig_consts<>+0(SB), Y8   // 2/pi
	VBROADCASTSD trig_consts<>+8(SB), Y9   // PI2A
	VBROADCASTSD trig_consts<>+16(SB), Y10 // PI2B
	VBROADCASTSD trig_consts<>+24(SB), Y11 // PI2C
	VBROADCASTSD trig_consts<>+136(SB), Y12 // abs mask

loop_sin_avx2:
	VMOVUPD (SI), Y0
	VBROADCASTSD trig_consts<>+144(SB), Y14
	VANDPD Y0, Y14, Y14                     // Y14 = sign bit of x
	VANDPD Y12, Y0, Y0                      // Y0 = |x|

	// Check |x| < 2^30
	VBROADCASTSD trig_consts<>+128(SB), Y13
	VCMPPD $1, Y13, Y0, Y13                 // Y13 = (|x| < 2^30)

	// k = roundToEven(|x| * (2/pi))
	VMULPD Y8, Y0, Y1
	VROUNDPD $0, Y1, Y1

	// Extended precision modular reduction: z = ((|x| - k*PI2A) - k*PI2B) - k*PI2C
	VMULPD Y9, Y1, Y2
	VSUBPD Y2, Y0, Y2
	VMULPD Y10, Y1, Y3
	VSUBPD Y3, Y2, Y2
	VMULPD Y11, Y1, Y3
	VSUBPD Y3, Y2, Y2                       // Y2 = z

	VMULPD Y2, Y2, Y3                       // Y3 = zz = z*z

	// S(z) = z + z * zz * P_sin(zz) in Y4
	VBROADCASTSD trig_consts<>+32(SB), Y4   // s0
	VBROADCASTSD trig_consts<>+40(SB), Y6   // s1
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+48(SB), Y6   // s2
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+56(SB), Y6   // s3
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+64(SB), Y6   // s4
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+72(SB), Y6   // s5
	VFMADD213PD Y6, Y3, Y4
	VMULPD Y3, Y4, Y4
	VFMADD213PD Y2, Y2, Y4                  // Y4 = S(z)

	// C(z) = (1.0 - 0.5*zz) + zz*zz*P_cos(zz) in Y5
	VBROADCASTSD trig_consts<>+80(SB), Y5   // c0
	VBROADCASTSD trig_consts<>+88(SB), Y6   // c1
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+96(SB), Y6   // c2
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+104(SB), Y6  // c3
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+112(SB), Y6  // c4
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+120(SB), Y6  // c5
	VFMADD213PD Y6, Y3, Y5
	VMULPD Y3, Y5, Y5
	VMULPD Y3, Y5, Y5
	VBROADCASTSD exp_consts<>+80(SB), Y6    // 0.5
	VMULPD Y6, Y3, Y6
	VBROADCASTSD exp_consts<>+88(SB), Y7    // 1.0
	VSUBPD Y6, Y7, Y7
	VADDPD Y7, Y5, Y5                       // Y5 = C(z)

	// Quadrant selection from k
	VCVTPD2DQY Y1, X15
	VPMOVZXDQ X15, Y15
	VPSLLQ $63, Y15, Y6                     // bit 63 = k & 1
	VBLENDVPD Y6, Y5, Y4, Y4                // (k & 1) ? C(z) : S(z)

	VPSLLQ $62, Y15, Y6                     // bit 63 = (k >> 1) & 1
	VBROADCASTSD trig_consts<>+144(SB), Y7
	VANDPD Y7, Y6, Y6
	VXORPD Y14, Y6, Y6                      // combine with sign of x
	VXORPD Y6, Y4, Y4

	// Mark lanes with !( |x| < 2^30 ) as NaN for scalar Payne-Hanek fallback
	VBROADCASTSD log_consts<>+168(SB), Y6
	VBLENDVPD Y13, Y4, Y6, Y4

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

	VBROADCASTSD trig_consts<>+0(SB), Y8   // 2/pi
	VBROADCASTSD trig_consts<>+8(SB), Y9   // PI2A
	VBROADCASTSD trig_consts<>+16(SB), Y10 // PI2B
	VBROADCASTSD trig_consts<>+24(SB), Y11 // PI2C
	VBROADCASTSD trig_consts<>+136(SB), Y12 // abs mask

loop_cos_avx2:
	VMOVUPD (SI), Y0
	VANDPD Y12, Y0, Y0                      // Y0 = |x|

	// Check |x| < 2^30
	VBROADCASTSD trig_consts<>+128(SB), Y13
	VCMPPD $1, Y13, Y0, Y13

	// k = roundToEven(|x| * (2/pi))
	VMULPD Y8, Y0, Y1
	VROUNDPD $0, Y1, Y1

	// Extended precision modular reduction: z = ((|x| - k*PI2A) - k*PI2B) - k*PI2C
	VMULPD Y9, Y1, Y2
	VSUBPD Y2, Y0, Y2
	VMULPD Y10, Y1, Y3
	VSUBPD Y3, Y2, Y2
	VMULPD Y11, Y1, Y3
	VSUBPD Y3, Y2, Y2                       // Y2 = z

	VMULPD Y2, Y2, Y3                       // Y3 = zz = z*z

	// S(z) in Y4
	VBROADCASTSD trig_consts<>+32(SB), Y4   // s0
	VBROADCASTSD trig_consts<>+40(SB), Y6   // s1
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+48(SB), Y6   // s2
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+56(SB), Y6   // s3
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+64(SB), Y6   // s4
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+72(SB), Y6   // s5
	VFMADD213PD Y6, Y3, Y4
	VMULPD Y3, Y4, Y4
	VFMADD213PD Y2, Y2, Y4                  // Y4 = S(z)

	// C(z) in Y5
	VBROADCASTSD trig_consts<>+80(SB), Y5   // c0
	VBROADCASTSD trig_consts<>+88(SB), Y6   // c1
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+96(SB), Y6   // c2
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+104(SB), Y6  // c3
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+112(SB), Y6  // c4
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+120(SB), Y6  // c5
	VFMADD213PD Y6, Y3, Y5
	VMULPD Y3, Y5, Y5
	VMULPD Y3, Y5, Y5
	VBROADCASTSD exp_consts<>+80(SB), Y6    // 0.5
	VMULPD Y6, Y3, Y6
	VBROADCASTSD exp_consts<>+88(SB), Y7    // 1.0
	VSUBPD Y6, Y7, Y7
	VADDPD Y7, Y5, Y5                       // Y5 = C(z)

	// Quadrant selection from k
	VCVTPD2DQY Y1, X15
	VPMOVZXDQ X15, Y15
	VPSLLQ $63, Y15, Y6                     // bit 63 = k & 1
	VBLENDVPD Y6, Y4, Y5, Y4                // (k & 1) ? S(z) : C(z)

	VPSRLQ $1, Y15, Y6
	VPXOR Y15, Y6, Y6
	VPSLLQ $63, Y6, Y6                      // bit 63 = (k ^ (k >> 1)) & 1
	VXORPD Y6, Y4, Y4

	// Mark lanes with !( |x| < 2^30 ) as NaN for scalar Payne-Hanek fallback
	VBROADCASTSD log_consts<>+168(SB), Y6
	VBLENDVPD Y13, Y4, Y6, Y4

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

	VBROADCASTSD trig_consts<>+0(SB), Y8   // 2/pi
	VBROADCASTSD trig_consts<>+8(SB), Y9   // PI2A
	VBROADCASTSD trig_consts<>+16(SB), Y10 // PI2B
	VBROADCASTSD trig_consts<>+24(SB), Y11 // PI2C
	VBROADCASTSD trig_consts<>+136(SB), Y12 // abs mask

loop_tan_avx2:
	VMOVUPD (SI), Y0
	VBROADCASTSD trig_consts<>+144(SB), Y14
	VANDPD Y0, Y14, Y14                     // Y14 = sign bit of x
	VANDPD Y12, Y0, Y0                      // Y0 = |x|

	// Check |x| < 2^30
	VBROADCASTSD trig_consts<>+128(SB), Y13
	VCMPPD $1, Y13, Y0, Y13

	// k = roundToEven(|x| * (2/pi))
	VMULPD Y8, Y0, Y1
	VROUNDPD $0, Y1, Y1

	// Extended precision modular reduction: z = ((|x| - k*PI2A) - k*PI2B) - k*PI2C
	VMULPD Y9, Y1, Y2
	VSUBPD Y2, Y0, Y2
	VMULPD Y10, Y1, Y3
	VSUBPD Y3, Y2, Y2
	VMULPD Y11, Y1, Y3
	VSUBPD Y3, Y2, Y2                       // Y2 = z

	VMULPD Y2, Y2, Y3                       // Y3 = zz = z*z

	// S(z) in Y4
	VBROADCASTSD trig_consts<>+32(SB), Y4   // s0
	VBROADCASTSD trig_consts<>+40(SB), Y6   // s1
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+48(SB), Y6   // s2
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+56(SB), Y6   // s3
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+64(SB), Y6   // s4
	VFMADD213PD Y6, Y3, Y4
	VBROADCASTSD trig_consts<>+72(SB), Y6   // s5
	VFMADD213PD Y6, Y3, Y4
	VMULPD Y3, Y4, Y4
	VFMADD213PD Y2, Y2, Y4                  // Y4 = S(z)

	// C(z) in Y5
	VBROADCASTSD trig_consts<>+80(SB), Y5   // c0
	VBROADCASTSD trig_consts<>+88(SB), Y6   // c1
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+96(SB), Y6   // c2
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+104(SB), Y6  // c3
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+112(SB), Y6  // c4
	VFMADD213PD Y6, Y3, Y5
	VBROADCASTSD trig_consts<>+120(SB), Y6  // c5
	VFMADD213PD Y6, Y3, Y5
	VMULPD Y3, Y5, Y5
	VMULPD Y3, Y5, Y5
	VBROADCASTSD exp_consts<>+80(SB), Y6    // 0.5
	VMULPD Y6, Y3, Y6
	VBROADCASTSD exp_consts<>+88(SB), Y7    // 1.0
	VSUBPD Y6, Y7, Y7
	VADDPD Y7, Y5, Y5                       // Y5 = C(z)

	// Quadrant selection: k odd -> -C(z)/S(z), k even -> S(z)/C(z)
	VCVTPD2DQY Y1, X15
	VPMOVZXDQ X15, Y15
	VPSLLQ $63, Y15, Y6                     // bit 63 = k & 1
	VBROADCASTSD trig_consts<>+144(SB), Y7
	VXORPD Y7, Y5, Y7                       // Y7 = -C(z)
	VBLENDVPD Y6, Y7, Y4, Y7                // num = (k & 1) ? -C(z) : S(z)
	VBLENDVPD Y6, Y4, Y5, Y5                // den = (k & 1) ? S(z) : C(z)
	VDIVPD Y5, Y7, Y4
	VXORPD Y14, Y4, Y4                      // restore sign of x

	// Mark lanes with !( |x| < 2^30 ) as NaN for scalar Payne-Hanek fallback
	VBROADCASTSD log_consts<>+168(SB), Y6
	VBLENDVPD Y13, Y4, Y6, Y4

	VMOVUPD Y4, (DX)

	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_tan_avx2

done_tan_avx2:
	VZEROUPPER
	RET

// func addSatInt8AVX2(a, b, result []int8)
TEXT ·addSatInt8AVX2(SB), NOSPLIT, $0-72
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $5, CX // n / 32
	JZ done_add_sat

loop_add_sat:
	VMOVDQU (SI), Y0
	VMOVDQU (DI), Y1
	VPADDSB Y1, Y0, Y2
	VMOVDQU Y2, (DX)

	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_add_sat

done_add_sat:
	VZEROUPPER
	RET

// func subSatInt8AVX2(a, b, result []int8)
TEXT ·subSatInt8AVX2(SB), NOSPLIT, $0-72
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX
	MOVQ a_len+8(FP), CX

	SHRQ $5, CX // n / 32
	JZ done_sub_sat

loop_sub_sat:
	VMOVDQU (SI), Y0
	VMOVDQU (DI), Y1
	VPSUBSB Y1, Y0, Y2
	VMOVDQU Y2, (DX)

	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_sub_sat

done_sub_sat:
	VZEROUPPER
	RET

// func expAVX512(a, result []float64)
TEXT ·expAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_exp512
	JMP ·expAVX2(SB)
done_exp512:
	VZEROUPPER
	RET

// func logAVX512(a, result []float64)
TEXT ·logAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_log512
	JMP ·logAVX2(SB)
done_log512:
	VZEROUPPER
	RET

// func sinAVX512(a, result []float64)
TEXT ·sinAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_sin512
	JMP ·sinAVX2(SB)
done_sin512:
	VZEROUPPER
	RET

// func cosAVX512(a, result []float64)
TEXT ·cosAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_cos512
	JMP ·cosAVX2(SB)
done_cos512:
	VZEROUPPER
	RET

// func tanAVX512(a, result []float64)
TEXT ·tanAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_tan512
	JMP ·tanAVX2(SB)
done_tan512:
	VZEROUPPER
	RET

// func addAVX512F32(a, b, result []float32)
TEXT ·addAVX512F32(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $5, CX // n / 32
    JZ tail_add512_f32

loop_add512_f32:
    VMOVUPS (SI), Z0
    VMOVUPS 64(SI), Z1
    VMOVUPS (DI), Z2
    VMOVUPS 64(DI), Z3
    VADDPS Z2, Z0, Z0
    VADDPS Z3, Z1, Z1
    VMOVUPS Z0, (DX)
    VMOVUPS Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_add512_f32

tail_add512_f32:
    TESTQ $16, AX
    JZ done_add512_f32
    VMOVUPS (SI), Z0
    VMOVUPS (DI), Z1
    VADDPS Z1, Z0, Z2
    VMOVUPS Z2, (DX)

done_add512_f32:
    VZEROUPPER
    RET

// func subAVX512F32(a, b, result []float32)
TEXT ·subAVX512F32(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $5, CX
    JZ tail_sub512_f32

loop_sub512_f32:
    VMOVUPS (SI), Z0
    VMOVUPS 64(SI), Z1
    VMOVUPS (DI), Z2
    VMOVUPS 64(DI), Z3
    VSUBPS Z2, Z0, Z0
    VSUBPS Z3, Z1, Z1
    VMOVUPS Z0, (DX)
    VMOVUPS Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_sub512_f32

tail_sub512_f32:
    TESTQ $16, AX
    JZ done_sub512_f32
    VMOVUPS (SI), Z0
    VMOVUPS (DI), Z1
    VSUBPS Z1, Z0, Z2
    VMOVUPS Z2, (DX)

done_sub512_f32:
    VZEROUPPER
    RET

// func mulAVX512F32(a, b, result []float32)
TEXT ·mulAVX512F32(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $5, CX
    JZ tail_mul512_f32

loop_mul512_f32:
    VMOVUPS (SI), Z0
    VMOVUPS 64(SI), Z1
    VMOVUPS (DI), Z2
    VMOVUPS 64(DI), Z3
    VMULPS Z2, Z0, Z0
    VMULPS Z3, Z1, Z1
    VMOVUPS Z0, (DX)
    VMOVUPS Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_mul512_f32

tail_mul512_f32:
    TESTQ $16, AX
    JZ done_mul512_f32
    VMOVUPS (SI), Z0
    VMOVUPS (DI), Z1
    VMULPS Z1, Z0, Z2
    VMOVUPS Z2, (DX)

done_mul512_f32:
    VZEROUPPER
    RET

// func divAVX512F32(a, b, result []float32)
TEXT ·divAVX512F32(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $5, CX
    JZ tail_div512_f32

loop_div512_f32:
    VMOVUPS (SI), Z0
    VMOVUPS 64(SI), Z1
    VMOVUPS (DI), Z2
    VMOVUPS 64(DI), Z3
    VDIVPS Z2, Z0, Z0
    VDIVPS Z3, Z1, Z1
    VMOVUPS Z0, (DX)
    VMOVUPS Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_div512_f32

tail_div512_f32:
    TESTQ $16, AX
    JZ done_div512_f32
    VMOVUPS (SI), Z0
    VMOVUPS (DI), Z1
    VDIVPS Z1, Z0, Z2
    VMOVUPS Z2, (DX)

done_div512_f32:
    VZEROUPPER
    RET

// func sqrtAVX512F32(a, result []float32)
TEXT ·sqrtAVX512F32(SB), NOSPLIT, $0-48
    MOVQ a_base+0(FP), SI
    MOVQ result_base+24(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $5, CX
    JZ tail_sqrt512_f32

loop_sqrt512_f32:
    VSQRTPS (SI), Z0
    VSQRTPS 64(SI), Z1
    VMOVUPS Z0, (DX)
    VMOVUPS Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_sqrt512_f32

tail_sqrt512_f32:
    TESTQ $16, AX
    JZ done_sqrt512_f32
    VSQRTPS (SI), Z0
    VMOVUPS Z0, (DX)

done_sqrt512_f32:
    VZEROUPPER
    RET

// func addSatInt8AVX512(a, b, result []int8)
TEXT ·addSatInt8AVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $7, CX // n / 128
    JZ tail_add_sat512

loop_add_sat512:
    VMOVDQU64 (SI), Z0
    VMOVDQU64 64(SI), Z1
    VMOVDQU64 (DI), Z2
    VMOVDQU64 64(DI), Z3
    VPADDSB Z2, Z0, Z0
    VPADDSB Z3, Z1, Z1
    VMOVDQU64 Z0, (DX)
    VMOVDQU64 Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_add_sat512

tail_add_sat512:
    TESTQ $64, AX
    JZ done_add_sat512
    VMOVDQU64 (SI), Z0
    VMOVDQU64 (DI), Z1
    VPADDSB Z1, Z0, Z2
    VMOVDQU64 Z2, (DX)

done_add_sat512:
    VZEROUPPER
    RET

// func subSatInt8AVX512(a, b, result []int8)
TEXT ·subSatInt8AVX512(SB), NOSPLIT, $0-72
    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ result_base+48(FP), DX
    MOVQ a_len+8(FP), CX

    MOVQ CX, AX
    SHRQ $7, CX // n / 128
    JZ tail_sub_sat512

loop_sub_sat512:
    VMOVDQU64 (SI), Z0
    VMOVDQU64 64(SI), Z1
    VMOVDQU64 (DI), Z2
    VMOVDQU64 64(DI), Z3
    VPSUBSB Z2, Z0, Z0
    VPSUBSB Z3, Z1, Z1
    VMOVDQU64 Z0, (DX)
    VMOVDQU64 Z1, 64(DX)
    ADDQ $128, SI
    ADDQ $128, DI
    ADDQ $128, DX
    DECQ CX
    JNZ loop_sub_sat512

tail_sub_sat512:
    TESTQ $64, AX
    JZ done_sub_sat512
    VMOVDQU64 (SI), Z0
    VMOVDQU64 (DI), Z1
    VPSUBSB Z1, Z0, Z2
    VMOVDQU64 Z2, (DX)

done_sub_sat512:
    VZEROUPPER
    RET

// func dotInt8AVX512VNNI(a []uint8, b []int8) int32
TEXT ·dotInt8AVX512VNNI(SB), NOSPLIT, $0-52
	MOVQ a_len+8(FP), CX
	SHRQ $6, CX // n / 64
	JZ zero_vnni

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI

	VPXORD Z0, Z0, Z0 // accumulator

loop_vnni:
	VMOVDQU64 (SI), Z1
	VMOVDQU64 (DI), Z2
	VPDPBUSD Z2, Z1, Z0

	ADDQ $64, SI
	ADDQ $64, DI
	DECQ CX
	JNZ loop_vnni

hsum_vnni:
	VEXTRACTI64X4 $1, Z0, Y1
	VPADDD Y1, Y0, Y0
	VEXTRACTI128 $1, Y0, X1
	VPADDD X1, X0, X0
	VPSHUFD $0x4E, X0, X1
	VPADDD X1, X0, X0
	VPSHUFD $0xB1, X0, X1
	VPADDD X1, X0, X0
	MOVL X0, ret+48(FP)

	VZEROUPPER
	RET

zero_vnni:
	MOVL $0, ret+48(FP)
	VZEROUPPER
	RET

// func addScalarAVX512F32(a []float32, b float32, result []float32)
TEXT ·addScalarAVX512F32(SB), NOSPLIT, $0-56
	MOVQ a_len+8(FP), CX
	MOVQ CX, AX
	TESTQ CX, CX
	JZ done_add_scalar512_f32

	MOVQ a_base+0(FP), SI
	MOVSS b+24(FP), X0
	MOVQ result_base+32(FP), DX
	VBROADCASTSS X0, Z0

	SHRQ $5, CX // n / 32
	JZ tail_add_scalar512_f32

loop_add_scalar512_f32_32:
	VMOVUPS (SI), Z1
	VMOVUPS 64(SI), Z2
	VADDPS Z0, Z1, Z3
	VADDPS Z0, Z2, Z4
	VMOVUPS Z3, (DX)
	VMOVUPS Z4, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_add_scalar512_f32_32

tail_add_scalar512_f32:
	TESTQ $16, AX
	JZ done_add_scalar512_f32
	VMOVUPS (SI), Z1
	VADDPS Z0, Z1, Z2
	VMOVUPS Z2, (DX)

done_add_scalar512_f32:
	VZEROUPPER
	RET

// func mulScalarAVX512F32(a []float32, b float32, result []float32)
TEXT ·mulScalarAVX512F32(SB), NOSPLIT, $0-56
	MOVQ a_len+8(FP), CX
	MOVQ CX, AX
	TESTQ CX, CX
	JZ done_mul_scalar512_f32

	MOVQ a_base+0(FP), SI
	MOVSS b+24(FP), X0
	MOVQ result_base+32(FP), DX
	VBROADCASTSS X0, Z0

	SHRQ $5, CX // n / 32
	JZ tail_mul_scalar512_f32

loop_mul_scalar512_f32_32:
	VMOVUPS (SI), Z1
	VMOVUPS 64(SI), Z2
	VMULPS Z0, Z1, Z3
	VMULPS Z0, Z2, Z4
	VMOVUPS Z3, (DX)
	VMOVUPS Z4, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_mul_scalar512_f32_32

tail_mul_scalar512_f32:
	TESTQ $16, AX
	JZ done_mul_scalar512_f32
	VMOVUPS (SI), Z1
	VMULPS Z0, Z1, Z2
	VMOVUPS Z2, (DX)

done_mul_scalar512_f32:
	VZEROUPPER
	RET

// func absAVX512F32(a, result []float32)
TEXT ·absAVX512F32(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	TESTQ CX, CX
	JZ done_abs512_f32

	MOVL $0x7fffffff, AX
	MOVL AX, X15
	VBROADCASTSS X15, Z15

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

	MOVQ CX, AX
	SHRQ $5, CX // n / 32
	JZ tail_abs512_f32

loop_abs512_f32_32:
	VMOVUPS (SI), Z0
	VMOVUPS 64(SI), Z1
	VPANDD Z15, Z0, Z0
	VPANDD Z15, Z1, Z1
	VMOVUPS Z0, (DX)
	VMOVUPS Z1, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_abs512_f32_32

tail_abs512_f32:
	TESTQ $16, AX
	JZ done_abs512_f32
	VMOVUPS (SI), Z0
	VPANDD Z15, Z0, Z0
	VMOVUPS Z0, (DX)

done_abs512_f32:
	VZEROUPPER
	RET

// func negAVX512F32(a, result []float32)
TEXT ·negAVX512F32(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	TESTQ CX, CX
	JZ done_neg512_f32

	MOVL $0x80000000, AX
	MOVL AX, X15
	VBROADCASTSS X15, Z15

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

	MOVQ CX, AX
	SHRQ $5, CX // n / 32
	JZ tail_neg512_f32

loop_neg512_f32_32:
	VMOVUPS (SI), Z0
	VMOVUPS 64(SI), Z1
	VPXORD Z15, Z0, Z0
	VPXORD Z15, Z1, Z1
	VMOVUPS Z0, (DX)
	VMOVUPS Z1, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_neg512_f32_32

tail_neg512_f32:
	TESTQ $16, AX
	JZ done_neg512_f32
	VMOVUPS (SI), Z0
	VPXORD Z15, Z0, Z0
	VMOVUPS Z0, (DX)

done_neg512_f32:
	VZEROUPPER
	RET

// func invAVX512F32(a, result []float32)
TEXT ·invAVX512F32(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	TESTQ CX, CX
	JZ done_inv512_f32

	MOVL $0x3f800000, AX // 1.0f
	MOVL AX, X15
	VBROADCASTSS X15, Z15

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

	MOVQ CX, AX
	SHRQ $5, CX // n / 32
	JZ tail_inv512_f32

loop_inv512_f32_32:
	VMOVUPS (SI), Z0
	VMOVUPS 64(SI), Z1
	VDIVPS Z0, Z15, Z2
	VDIVPS Z1, Z15, Z3
	VMOVUPS Z2, (DX)
	VMOVUPS Z3, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_inv512_f32_32

tail_inv512_f32:
	TESTQ $16, AX
	JZ done_inv512_f32
	VMOVUPS (SI), Z0
	VDIVPS Z0, Z15, Z1
	VMOVUPS Z1, (DX)

done_inv512_f32:
	VZEROUPPER
	RET

// func fmaAVX512F32(a, b, c, result []float32)
TEXT ·fmaAVX512F32(SB), NOSPLIT, $0-96
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ c_base+48(FP), BX
	MOVQ result_base+72(FP), DX
	MOVQ a_len+8(FP), CX

	MOVQ CX, AX
	SHRQ $5, CX // n / 32
	JZ tail_fma512_f32

loop_fma512_f32_32:
	VMOVUPS (SI), Z0
	VMOVUPS 64(SI), Z1
	VMOVUPS (DI), Z2
	VMOVUPS 64(DI), Z3
	VMOVUPS (BX), Z4
	VMOVUPS 64(BX), Z5
	VFMADD213PS Z4, Z2, Z0
	VFMADD213PS Z5, Z3, Z1
	VMOVUPS Z0, (DX)
	VMOVUPS Z1, 64(DX)

	ADDQ $128, SI
	ADDQ $128, DI
	ADDQ $128, BX
	ADDQ $128, DX
	DECQ CX
	JNZ loop_fma512_f32_32

tail_fma512_f32:
	TESTQ $16, AX
	JZ done_fma512_f32
	VMOVUPS (SI), Z0
	VMOVUPS (DI), Z1
	VMOVUPS (BX), Z2
	VFMADD213PS Z2, Z1, Z0
	VMOVUPS Z0, (DX)

done_fma512_f32:
	VZEROUPPER
	RET

// func floorAVX512(a, result []float64)
TEXT ·floorAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	TESTQ CX, CX
	JZ done_floor512

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

	MOVQ CX, AX
	SHRQ $4, CX // n / 16
	JZ tail_floor512

loop_floor512_16:
	VMOVUPD (SI), Z0
	VMOVUPD 64(SI), Z1
	VRNDSCALEPD $1, Z0, Z2
	VRNDSCALEPD $1, Z1, Z3
	VMOVUPD Z2, (DX)
	VMOVUPD Z3, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_floor512_16

tail_floor512:
	TESTQ $8, AX
	JZ done_floor512
	VMOVUPD (SI), Z0
	VRNDSCALEPD $1, Z0, Z1
	VMOVUPD Z1, (DX)

done_floor512:
	VZEROUPPER
	RET

// func ceilAVX512(a, result []float64)
TEXT ·ceilAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	TESTQ CX, CX
	JZ done_ceil512

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

	MOVQ CX, AX
	SHRQ $4, CX // n / 16
	JZ tail_ceil512

loop_ceil512_16:
	VMOVUPD (SI), Z0
	VMOVUPD 64(SI), Z1
	VRNDSCALEPD $2, Z0, Z2
	VRNDSCALEPD $2, Z1, Z3
	VMOVUPD Z2, (DX)
	VMOVUPD Z3, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_ceil512_16

tail_ceil512:
	TESTQ $8, AX
	JZ done_ceil512
	VMOVUPD (SI), Z0
	VRNDSCALEPD $2, Z0, Z1
	VMOVUPD Z1, (DX)

done_ceil512:
	VZEROUPPER
	RET

// func truncAVX512(a, result []float64)
TEXT ·truncAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	TESTQ CX, CX
	JZ done_trunc512

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

	MOVQ CX, AX
	SHRQ $4, CX // n / 16
	JZ tail_trunc512

loop_trunc512_16:
	VMOVUPD (SI), Z0
	VMOVUPD 64(SI), Z1
	VRNDSCALEPD $3, Z0, Z2
	VRNDSCALEPD $3, Z1, Z3
	VMOVUPD Z2, (DX)
	VMOVUPD Z3, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_trunc512_16

tail_trunc512:
	TESTQ $8, AX
	JZ done_trunc512
	VMOVUPD (SI), Z0
	VRNDSCALEPD $3, Z0, Z1
	VMOVUPD Z1, (DX)

done_trunc512:
	VZEROUPPER
	RET

// func roundAVX512(a, result []float64)
TEXT ·roundAVX512(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	TESTQ CX, CX
	JZ done_round512

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

	MOVQ CX, AX
	SHRQ $4, CX // n / 16
	JZ tail_round512

loop_round512_16:
	VMOVUPD (SI), Z0
	VMOVUPD 64(SI), Z1
	VRNDSCALEPD $0, Z0, Z2
	VRNDSCALEPD $0, Z1, Z3
	VMOVUPD Z2, (DX)
	VMOVUPD Z3, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_round512_16

tail_round512:
	TESTQ $8, AX
	JZ done_round512
	VMOVUPD (SI), Z0
	VRNDSCALEPD $0, Z0, Z1
	VMOVUPD Z1, (DX)

done_round512:
	VZEROUPPER
	RET

// func maxAVX512(a, b, result []float64)
TEXT ·maxAVX512(SB), NOSPLIT, $0-72
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX
	MOVQ a_len+8(FP), CX

	MOVQ CX, AX
	SHRQ $4, CX // n / 16
	JZ tail_max512

loop_max512_16:
	VMOVUPD (SI), Z0
	VMOVUPD 64(SI), Z1
	VMOVUPD (DI), Z2
	VMOVUPD 64(DI), Z3
	VMAXPD Z2, Z0, Z0
	VMAXPD Z3, Z1, Z1
	VMOVUPD Z0, (DX)
	VMOVUPD Z1, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_max512_16

tail_max512:
	TESTQ $8, AX
	JZ done_max512
	VMOVUPD (SI), Z0
	VMOVUPD (DI), Z1
	VMAXPD Z1, Z0, Z2
	VMOVUPD Z2, (DX)

done_max512:
	VZEROUPPER
	RET

// func minAVX512(a, b, result []float64)
TEXT ·minAVX512(SB), NOSPLIT, $0-72
	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX
	MOVQ a_len+8(FP), CX

	MOVQ CX, AX
	SHRQ $4, CX // n / 16
	JZ tail_min512

loop_min512_16:
	VMOVUPD (SI), Z0
	VMOVUPD 64(SI), Z1
	VMOVUPD (DI), Z2
	VMOVUPD 64(DI), Z3
	VMINPD Z2, Z0, Z0
	VMINPD Z3, Z1, Z1
	VMOVUPD Z0, (DX)
	VMOVUPD Z1, 64(DX)
	ADDQ $128, SI
	ADDQ $128, DI
	ADDQ $128, DX
	DECQ CX
	JNZ loop_min512_16

tail_min512:
	TESTQ $8, AX
	JZ done_min512
	VMOVUPD (SI), Z0
	VMOVUPD (DI), Z1
	VMINPD Z1, Z0, Z2
	VMOVUPD Z2, (DX)

done_min512:
	VZEROUPPER
	RET

// func floorAVX2(a, result []float64)
TEXT ·floorAVX2(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $2, CX // n / 4
	JZ done_floor_avx2

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_floor_avx2:
	VROUNDPD $1, (SI), Y0
	VMOVUPD Y0, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_floor_avx2

done_floor_avx2:
	VZEROUPPER
	RET

// func ceilAVX2(a, result []float64)
TEXT ·ceilAVX2(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $2, CX // n / 4
	JZ done_ceil_avx2

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_ceil_avx2:
	VROUNDPD $2, (SI), Y0
	VMOVUPD Y0, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_ceil_avx2

done_ceil_avx2:
	VZEROUPPER
	RET

// func truncAVX2(a, result []float64)
TEXT ·truncAVX2(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $2, CX // n / 4
	JZ done_trunc_avx2

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_trunc_avx2:
	VROUNDPD $3, (SI), Y0
	VMOVUPD Y0, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_trunc_avx2

done_trunc_avx2:
	VZEROUPPER
	RET

// func roundAVX2(a, result []float64)
TEXT ·roundAVX2(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $2, CX // n / 4
	JZ done_round_avx2

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_round_avx2:
	VROUNDPD $0, (SI), Y0
	VMOVUPD Y0, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_round_avx2

done_round_avx2:
	VZEROUPPER
	RET

// func maxAVX2(a, b, result []float64)
TEXT ·maxAVX2(SB), NOSPLIT, $0-72
	MOVQ a_len+8(FP), CX
	SHRQ $2, CX // n / 4
	JZ done_max_avx2

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX

loop_max_avx2:
	VMOVUPD (SI), Y0
	VMOVUPD (DI), Y1
	VCMPPD $3, Y1, Y0, Y3
	VMAXPD Y1, Y0, Y2
	VPOR Y3, Y2, Y2
	VMOVUPD Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_max_avx2

done_max_avx2:
	VZEROUPPER
	RET

// func minAVX2(a, b, result []float64)
TEXT ·minAVX2(SB), NOSPLIT, $0-72
	MOVQ a_len+8(FP), CX
	SHRQ $2, CX // n / 4
	JZ done_min_avx2

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX

loop_min_avx2:
	VMOVUPD (SI), Y0
	VMOVUPD (DI), Y1
	VCMPPD $3, Y1, Y0, Y3
	VMINPD Y1, Y0, Y2
	VPOR Y3, Y2, Y2
	VMOVUPD Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_min_avx2

done_min_avx2:
	VZEROUPPER
	RET

// func addAVX2F32(a, b, result []float32)
TEXT ·addAVX2F32(SB), NOSPLIT, $0-72
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_add_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX

loop_add_avx2_f32:
	VMOVUPS (SI), Y0
	VMOVUPS (DI), Y1
	VADDPS Y1, Y0, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_add_avx2_f32

done_add_avx2_f32:
	VZEROUPPER
	RET

// func subAVX2F32(a, b, result []float32)
TEXT ·subAVX2F32(SB), NOSPLIT, $0-72
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_sub_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX

loop_sub_avx2_f32:
	VMOVUPS (SI), Y0
	VMOVUPS (DI), Y1
	VSUBPS Y1, Y0, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_sub_avx2_f32

done_sub_avx2_f32:
	VZEROUPPER
	RET

// func mulAVX2F32(a, b, result []float32)
TEXT ·mulAVX2F32(SB), NOSPLIT, $0-72
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_mul_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX

loop_mul_avx2_f32:
	VMOVUPS (SI), Y0
	VMOVUPS (DI), Y1
	VMULPS Y1, Y0, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_mul_avx2_f32

done_mul_avx2_f32:
	VZEROUPPER
	RET

// func divAVX2F32(a, b, result []float32)
TEXT ·divAVX2F32(SB), NOSPLIT, $0-72
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_div_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ result_base+48(FP), DX

loop_div_avx2_f32:
	VMOVUPS (SI), Y0
	VMOVUPS (DI), Y1
	VDIVPS Y1, Y0, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_div_avx2_f32

done_div_avx2_f32:
	VZEROUPPER
	RET

// func sqrtAVX2F32(a, result []float32)
TEXT ·sqrtAVX2F32(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_sqrt_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_sqrt_avx2_f32:
	VSQRTPS (SI), Y0
	VMOVUPS Y0, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_sqrt_avx2_f32

done_sqrt_avx2_f32:
	VZEROUPPER
	RET

// func addScalarAVX2F32(a []float32, b float32, result []float32)
TEXT ·addScalarAVX2F32(SB), NOSPLIT, $0-56
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_add_scalar_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVSS b+24(FP), X0
	MOVQ result_base+32(FP), DX
	VBROADCASTSS X0, Y0

loop_add_scalar_avx2_f32:
	VMOVUPS (SI), Y1
	VADDPS Y0, Y1, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_add_scalar_avx2_f32

done_add_scalar_avx2_f32:
	VZEROUPPER
	RET

// func mulScalarAVX2F32(a []float32, b float32, result []float32)
TEXT ·mulScalarAVX2F32(SB), NOSPLIT, $0-56
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_mul_scalar_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVSS b+24(FP), X0
	MOVQ result_base+32(FP), DX
	VBROADCASTSS X0, Y0

loop_mul_scalar_avx2_f32:
	VMOVUPS (SI), Y1
	VMULPS Y0, Y1, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_mul_scalar_avx2_f32

done_mul_scalar_avx2_f32:
	VZEROUPPER
	RET

// func absAVX2F32(a, result []float32)
TEXT ·absAVX2F32(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_abs_avx2_f32

	MOVL $0x7fffffff, AX
	MOVL AX, X1
	VBROADCASTSS X1, Y1

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_abs_avx2_f32:
	VMOVUPS (SI), Y0
	VANDPS Y1, Y0, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_abs_avx2_f32

done_abs_avx2_f32:
	VZEROUPPER
	RET

// func negAVX2F32(a, result []float32)
TEXT ·negAVX2F32(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_neg_avx2_f32

	MOVL $0x80000000, AX
	MOVL AX, X1
	VBROADCASTSS X1, Y1

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_neg_avx2_f32:
	VMOVUPS (SI), Y0
	VXORPS Y1, Y0, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_neg_avx2_f32

done_neg_avx2_f32:
	VZEROUPPER
	RET

// func invAVX2F32(a, result []float32)
TEXT ·invAVX2F32(SB), NOSPLIT, $0-48
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_inv_avx2_f32

	MOVL $0x3f800000, AX // 1.0f
	MOVL AX, X1
	VBROADCASTSS X1, Y1

	MOVQ a_base+0(FP), SI
	MOVQ result_base+24(FP), DX

loop_inv_avx2_f32:
	VMOVUPS (SI), Y0
	VDIVPS Y0, Y1, Y2
	VMOVUPS Y2, (DX)
	ADDQ $32, SI
	ADDQ $32, DX
	DECQ CX
	JNZ loop_inv_avx2_f32

done_inv_avx2_f32:
	VZEROUPPER
	RET

// func fmaAVX2F32(a, b, c, result []float32)
TEXT ·fmaAVX2F32(SB), NOSPLIT, $0-96
	MOVQ a_len+8(FP), CX
	SHRQ $3, CX // n / 8
	JZ done_fma_avx2_f32

	MOVQ a_base+0(FP), SI
	MOVQ b_base+24(FP), DI
	MOVQ c_base+48(FP), BX
	MOVQ result_base+72(FP), DX

loop_fma_avx2_f32:
	VMOVUPS (SI), Y0
	VMOVUPS (DI), Y1
	VMOVUPS (BX), Y2
	VFMADD213PS Y2, Y1, Y0
	VMOVUPS Y0, (DX)
	ADDQ $32, SI
	ADDQ $32, DI
	ADDQ $32, BX
	ADDQ $32, DX
	DECQ CX
	JNZ loop_fma_avx2_f32

done_fma_avx2_f32:
	VZEROUPPER
	RET



