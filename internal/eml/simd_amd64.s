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
