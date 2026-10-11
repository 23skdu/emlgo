#include <mpfr.h>

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
	const char *name;
	uint64_t bits;
} vector64;

typedef struct {
	const char *name;
	uint32_t bits;
} vector32;

static const vector64 exp64_inputs[] = {
	{"exp64", UINT64_C(0xc085e00000000000)}, // -700
	{"exp64", UINT64_C(0xc024000000000000)}, // -10
	{"exp64", UINT64_C(0xbff0000000000000)}, // -1
	{"exp64", UINT64_C(0x0000000000000000)}, // 0
	{"exp64", UINT64_C(0x3fb999999999999a)}, // 0.1
	{"exp64", UINT64_C(0x3ff0000000000000)}, // 1
	{"exp64", UINT64_C(0x4024000000000000)}, // 10
	{"exp64", UINT64_C(0x4085e00000000000)}, // 700
	{"exp64", UINT64_C(0xc087480000000000)}, // -745
};

static const vector64 log64_inputs[] = {
	{"log64", UINT64_C(0x0000000000000001)}, // minimum subnormal
	{"log64", UINT64_C(0x000fffffffffffff)}, // maximum subnormal
	{"log64", UINT64_C(0x0010000000000000)}, // minimum normal
	{"log64", UINT64_C(0x3fe0000000000000)}, // 0.5
	{"log64", UINT64_C(0x3ff0000000000000)}, // 1
	{"log64", UINT64_C(0x3ff0000000000001)}, // next float above 1
	{"log64", UINT64_C(0x4000000000000000)}, // 2
	{"log64", UINT64_C(0x7fefffffffffffff)}, // maximum finite
};

static const vector32 exp32_inputs[] = {
	{"exp32", UINT32_C(0xc2dc0000)}, // -110
	{"exp32", UINT32_C(0xc2ae0000)}, // -87
	{"exp32", UINT32_C(0xbf800000)}, // -1
	{"exp32", UINT32_C(0x00000000)}, // 0
	{"exp32", UINT32_C(0x3dcccccd)}, // 0.1
	{"exp32", UINT32_C(0x3f800000)}, // 1
	{"exp32", UINT32_C(0x42b00000)}, // 88
};

static const vector32 log32_inputs[] = {
	{"log32", UINT32_C(0x00000001)}, // minimum subnormal
	{"log32", UINT32_C(0x007fffff)}, // maximum subnormal
	{"log32", UINT32_C(0x00800000)}, // minimum normal
	{"log32", UINT32_C(0x3f000000)}, // 0.5
	{"log32", UINT32_C(0x3f800000)}, // 1
	{"log32", UINT32_C(0x3f800001)}, // next float above 1
	{"log32", UINT32_C(0x40000000)}, // 2
	{"log32", UINT32_C(0x7f7fffff)}, // maximum finite
};

static void emit64(FILE *out, const vector64 *v, int is_exp) {
	double input;
	uint64_t result_bits;
	mpfr_t x, y;
	memcpy(&input, &v->bits, sizeof(input));
	mpfr_init2(x, 256);
	mpfr_init2(y, 53);
	mpfr_set_d(x, input, MPFR_RNDN);
	if (is_exp) {
		mpfr_exp(y, x, MPFR_RNDN);
	} else {
		mpfr_log(y, x, MPFR_RNDN);
	}
	double result = mpfr_get_d(y, MPFR_RNDN);
	memcpy(&result_bits, &result, sizeof(result_bits));
	fprintf(out, "%s,%016llx,%016llx,53\n", v->name,
		(unsigned long long)v->bits, (unsigned long long)result_bits);
	mpfr_clear(y);
	mpfr_clear(x);
}

static void emit32(FILE *out, const vector32 *v, int is_exp) {
	float input;
	uint32_t result_bits;
	mpfr_t x, y;
	memcpy(&input, &v->bits, sizeof(input));
	mpfr_init2(x, 256);
	mpfr_init2(y, 24);
	mpfr_set_flt(x, input, MPFR_RNDN);
	if (is_exp) {
		mpfr_exp(y, x, MPFR_RNDN);
	} else {
		mpfr_log(y, x, MPFR_RNDN);
	}
	float result = mpfr_get_flt(y, MPFR_RNDN);
	memcpy(&result_bits, &result, sizeof(result_bits));
	fprintf(out, "%s,%08x,%08x,24\n", v->name, v->bits, result_bits);
	mpfr_clear(y);
	mpfr_clear(x);
}

int main(int argc, char **argv) {
	if (argc != 2) {
		fprintf(stderr, "usage: %s OUTPUT.csv\n", argv[0]);
		return 2;
	}
	if (strcmp(mpfr_get_version(), "4.2.2") != 0) {
		fprintf(stderr, "expected MPFR 4.2.2, got %s\n", mpfr_get_version());
		return 1;
	}

	FILE *out = fopen(argv[1], "w");
	if (out == NULL) {
		perror(argv[1]);
		return 1;
	}
	fprintf(out, "# MPFR %s; MPFR_RNDN; exact binary inputs; result precision 53/24 bits\n",
		mpfr_get_version());
	fprintf(out, "operation,input_bits,reference_bits,result_precision\n");
	for (size_t i = 0; i < sizeof(exp64_inputs) / sizeof(exp64_inputs[0]); i++) {
		emit64(out, &exp64_inputs[i], 1);
	}
	for (size_t i = 0; i < sizeof(log64_inputs) / sizeof(log64_inputs[0]); i++) {
		emit64(out, &log64_inputs[i], 0);
	}
	for (size_t i = 0; i < sizeof(exp32_inputs) / sizeof(exp32_inputs[0]); i++) {
		emit32(out, &exp32_inputs[i], 1);
	}
	for (size_t i = 0; i < sizeof(log32_inputs) / sizeof(log32_inputs[0]); i++) {
		emit32(out, &log32_inputs[i], 0);
	}
	if (fclose(out) != 0) {
		perror(argv[1]);
		return 1;
	}
	return 0;
}
