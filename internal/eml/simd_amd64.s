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
	MOVQ a_ptr+0(FP), SI
	MOVQ b_ptr+24(FP), DI
	MOVQ res_ptr+48(FP), DX
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
    MOVQ a_ptr+0(FP), SI
    MOVQ b_ptr+24(FP), DI
    MOVQ res_ptr+48(FP), DX
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
    MOVQ a_ptr+0(FP), SI
    MOVQ b_ptr+24(FP), DI
    MOVQ res_ptr+48(FP), DX
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
    MOVQ a_ptr+0(FP), SI
    MOVQ b_ptr+24(FP), DI
    MOVQ res_ptr+48(FP), DX
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
    MOVQ a_ptr+0(FP), SI
    MOVSD b+24(FP), X0
    VBROADCASTSD X0, Y0
    MOVQ res_ptr+32(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_add_scalar
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
    MOVQ a_ptr+0(FP), SI
    MOVSD b+24(FP), X0
    VBROADCASTSD X0, Y0
    MOVQ res_ptr+32(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_mul_scalar
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
	MOVQ a_ptr+0(FP), SI
	MOVQ b_ptr+24(FP), DI
	MOVQ res_ptr+48(FP), DX
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
	VZEROUPPER
	RET

// func subAVX512(a, b, result []float64)
TEXT ·subAVX512(SB), NOSPLIT, $0-72
    MOVQ a_ptr+0(FP), SI
    MOVQ b_ptr+24(FP), DI
    MOVQ res_ptr+48(FP), DX
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
    VZEROUPPER
    RET

// func mulAVX512(a, b, result []float64)
TEXT ·mulAVX512(SB), NOSPLIT, $0-72
    MOVQ a_ptr+0(FP), SI
    MOVQ b_ptr+24(FP), DI
    MOVQ res_ptr+48(FP), DX
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
    VZEROUPPER
    RET

// func divAVX512(a, b, result []float64)
TEXT ·divAVX512(SB), NOSPLIT, $0-72
    MOVQ a_ptr+0(FP), SI
    MOVQ b_ptr+24(FP), DI
    MOVQ res_ptr+48(FP), DX
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
    VZEROUPPER
    RET

// func addScalarAVX512(a []float64, b float64, result []float64)
TEXT ·addScalarAVX512(SB), NOSPLIT, $0-56
    MOVQ a_ptr+0(FP), SI
    MOVSD b+24(FP), X0
    VBROADCASTSD X0, Z0
    MOVQ res_ptr+32(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_add_scalar512
loop_add_scalar512:
    VMOVUPD (SI), Z1
    VADDPD Z0, Z1, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_add_scalar512
done_add_scalar512:
    VZEROUPPER
    RET

// func mulScalarAVX512(a []float64, b float64, result []float64)
TEXT ·mulScalarAVX512(SB), NOSPLIT, $0-56
    MOVQ a_ptr+0(FP), SI
    MOVSD b+24(FP), X0
    VBROADCASTSD X0, Z0
    MOVQ res_ptr+32(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_mul_scalar512
loop_mul_scalar512:
    VMOVUPD (SI), Z1
    VMULPD Z0, Z1, Z2
    VMOVUPD Z2, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_mul_scalar512
done_mul_scalar512:
    VZEROUPPER
    RET

// func sqrtAVX2(a, result []float64)
TEXT ·sqrtAVX2(SB), NOSPLIT, $0-48
    MOVQ a_ptr+0(FP), SI
    MOVQ res_ptr+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $2, CX
    JZ done_sqrt2
loop_sqrt2:
    VMOVUPD (SI), Y0
    VSQRTPD Y0, Y1
    VMOVUPD Y1, (DX)
    ADDQ $32, SI
    ADDQ $32, DX
    DECQ CX
    JNZ loop_sqrt2
done_sqrt2:
    VZEROUPPER
    RET

// func sqrtAVX512(a, result []float64)
TEXT ·sqrtAVX512(SB), NOSPLIT, $0-48
    MOVQ a_ptr+0(FP), SI
    MOVQ res_ptr+24(FP), DX
    MOVQ a_len+8(FP), CX
    SHRQ $3, CX
    JZ done_sqrt512
loop_sqrt512:
    VMOVUPD (SI), Z0
    VSQRTPD Z0, Z1
    VMOVUPD Z1, (DX)
    ADDQ $64, SI
    ADDQ $64, DX
    DECQ CX
    JNZ loop_sqrt512
done_sqrt512:
    VZEROUPPER
    RET
