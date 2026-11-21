#include <metal_stdlib>
using namespace metal;

// ============================================================================
// Helper functions for numerically stable complex64 (float2) arithmetic
// ============================================================================

inline float hypot_safe(float a, float b) {
    float aa = abs(a);
    float bb = abs(b);
    float m = max(aa, bb);
    if (m == 0.0f) {
        return 0.0f;
    }
    if (isinf(m)) {
        return INFINITY;
    }
    float ra = aa / m;
    float rb = bb / m;
    return m * sqrt(ra * ra + rb * rb);
}

inline float2 c_add(float2 a, float2 b) {
    return float2(a.x + b.x, a.y + b.y);
}

inline float2 c_sub(float2 a, float2 b) {
    return float2(a.x - b.x, a.y - b.y);
}

inline float2 c_div(float2 a, float2 b) {
    float ar = a.x, ai = a.y;
    float br = b.x, bi = b.y;
    if (abs(br) >= abs(bi)) {
        float r = bi / br;
        float denom = br + bi * r;
        return float2((ar + ai * r) / denom, (ai - ar * r) / denom);
    } else {
        float r = br / bi;
        float denom = bi + br * r;
        return float2((ar * r + ai) / denom, (ai * r - ar) / denom);
    }
}

inline float2 c_exp(float2 z) {
    float ea = exp(z.x);
    float s = sin(z.y);
    float c = cos(z.y);
    return float2(ea * c, ea * s);
}

inline float2 c_log(float2 z) {
    return float2(log(hypot_safe(z.x, z.y)), atan2(z.y, z.x));
}

inline float sinh_accurate(float x) {
    float ax = abs(x);
    if (ax < 0.25f) {
        float x2 = x * x;
        // x * (1 + x^2*(1/6 + x^2*(1/120 + x^2/5040)))
        return fma(x * x2, fma(x2, fma(x2, 1.0f / 5040.0f, 1.0f / 120.0f), 1.0f / 6.0f), x);
    }
    return sinh(x);
}

inline float tanh_accurate(float x) {
    float ax = abs(x);
    if (ax < 0.25f) {
        return sinh_accurate(x) / cosh(x);
    }
    return tanh(x);
}

inline float2 c_sin(float2 z) {
    return float2(sin(z.x) * cosh(z.y), cos(z.x) * sinh_accurate(z.y));
}

inline float2 c_cos(float2 z) {
    return float2(cos(z.x) * cosh(z.y), -sin(z.x) * sinh_accurate(z.y));
}

inline float2 c_tan(float2 z) {
    if (abs(z.y) > 10.0f) {
        return float2(0.0f, copysign(1.0f, z.y));
    }
    return c_div(c_sin(z), c_cos(z));
}

inline float2 c_sinh(float2 z) {
    return float2(sinh_accurate(z.x) * cos(z.y), cosh(z.x) * sin(z.y));
}

inline float2 c_cosh(float2 z) {
    return float2(cosh(z.x) * cos(z.y), sinh_accurate(z.x) * sin(z.y));
}

inline float2 c_tanh(float2 z) {
    if (abs(z.x) > 10.0f) {
        return float2(copysign(1.0f, z.x), 0.0f);
    }
    return c_div(c_sinh(z), c_cosh(z));
}

inline float2 c_sqrt(float2 z) {
    float r = hypot_safe(z.x, z.y);
    if (r == 0.0f) {
        return float2(0.0f, 0.0f);
    }
    if (z.x >= 0.0f) {
        float re = sqrt(0.5f * (r + z.x));
        float im = (re == 0.0f) ? 0.0f : (z.y / (2.0f * re));
        return float2(re, im);
    } else {
        float im = copysign(sqrt(0.5f * (r - z.x)), z.y);
        float re = (im == 0.0f) ? 0.0f : (z.y / (2.0f * im));
        return float2(re, im);
    }
}

// ============================================================================
// Float32 Unary Kernels
// ============================================================================

kernel void kernel_exp_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = exp(x[id]);
}

kernel void kernel_log_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = log(x[id]);
}

kernel void kernel_sin_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = sin(x[id]);
}

kernel void kernel_cos_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = cos(x[id]);
}

kernel void kernel_tan_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = tan(x[id]);
}

kernel void kernel_sinh_f32(device const float *x [[buffer(0)]],
                            device float *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = sinh_accurate(x[id]);
}

kernel void kernel_cosh_f32(device const float *x [[buffer(0)]],
                            device float *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = cosh(x[id]);
}

kernel void kernel_tanh_f32(device const float *x [[buffer(0)]],
                            device float *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = tanh_accurate(x[id]);
}

kernel void kernel_sqrt_f32(device const float *x [[buffer(0)]],
                            device float *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = sqrt(x[id]);
}

kernel void kernel_abs_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = abs(x[id]);
}

kernel void kernel_neg_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = -x[id];
}

kernel void kernel_inv_f32(device const float *x [[buffer(0)]],
                           device float *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = 1.0f / x[id];
}

// ============================================================================
// Float32 Binary, Ternary & Reduction Kernels
// ============================================================================

kernel void kernel_eml_f32(device const float *x [[buffer(0)]],
                           device const float *y [[buffer(1)]],
                           device float *result [[buffer(2)]],
                           constant uint &n [[buffer(3)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = exp(x[id]) - log(y[id]);
}

kernel void kernel_fma_f32(device const float *a [[buffer(0)]],
                           device const float *b [[buffer(1)]],
                           device const float *c [[buffer(2)]],
                           device float *result [[buffer(3)]],
                           constant uint &n [[buffer(4)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = fma(a[id], b[id], c[id]);
}

kernel void kernel_dot_f32(device const float *a [[buffer(0)]],
                           device const float *b [[buffer(1)]],
                           device float *partials [[buffer(2)]],
                           constant uint &n [[buffer(3)]],
                           uint gid [[thread_position_in_grid]],
                           uint tid [[thread_position_in_threadgroup]],
                           uint tgIdx [[threadgroup_position_in_grid]],
                           uint tgSize [[threads_per_threadgroup]]) {
    threadgroup float sdata[256];
    sdata[tid] = (gid < n) ? (a[gid] * b[gid]) : 0.0f;
    threadgroup_barrier(mem_flags::mem_threadgroup);

    for (uint s = tgSize >> 1; s > 0; s >>= 1) {
        if (tid < s) {
            sdata[tid] += sdata[tid + s];
        }
        threadgroup_barrier(mem_flags::mem_threadgroup);
    }

    if (tid == 0) {
        partials[tgIdx] = sdata[0];
    }
}

// ============================================================================
// Complex64 (float2) Unary, Binary & Reduction Kernels
// ============================================================================

kernel void kernel_exp_c64(device const float2 *x [[buffer(0)]],
                           device float2 *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_exp(x[id]);
}

kernel void kernel_log_c64(device const float2 *x [[buffer(0)]],
                           device float2 *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_log(x[id]);
}

kernel void kernel_sin_c64(device const float2 *x [[buffer(0)]],
                           device float2 *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_sin(x[id]);
}

kernel void kernel_cos_c64(device const float2 *x [[buffer(0)]],
                           device float2 *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_cos(x[id]);
}

kernel void kernel_tan_c64(device const float2 *x [[buffer(0)]],
                           device float2 *result [[buffer(1)]],
                           constant uint &n [[buffer(2)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_tan(x[id]);
}

kernel void kernel_sinh_c64(device const float2 *x [[buffer(0)]],
                            device float2 *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_sinh(x[id]);
}

kernel void kernel_cosh_c64(device const float2 *x [[buffer(0)]],
                            device float2 *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_cosh(x[id]);
}

kernel void kernel_tanh_c64(device const float2 *x [[buffer(0)]],
                            device float2 *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_tanh(x[id]);
}

kernel void kernel_sqrt_c64(device const float2 *x [[buffer(0)]],
                            device float2 *result [[buffer(1)]],
                            constant uint &n [[buffer(2)]],
                            uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_sqrt(x[id]);
}

kernel void kernel_eml_c64(device const float2 *x [[buffer(0)]],
                           device const float2 *y [[buffer(1)]],
                           device float2 *result [[buffer(2)]],
                           constant uint &n [[buffer(3)]],
                           uint id [[thread_position_in_grid]]) {
    if (id >= n) return;
    result[id] = c_sub(c_exp(x[id]), c_log(y[id]));
}

// Hermitian inner product: sum(a[i] * conj(b[i]))
kernel void kernel_dot_c64(device const float2 *a [[buffer(0)]],
                           device const float2 *b [[buffer(1)]],
                           device float2 *partials [[buffer(2)]],
                           constant uint &n [[buffer(3)]],
                           uint gid [[thread_position_in_grid]],
                           uint tid [[thread_position_in_threadgroup]],
                           uint tgIdx [[threadgroup_position_in_grid]],
                           uint tgSize [[threads_per_threadgroup]]) {
    threadgroup float2 sdata[256];
    if (gid < n) {
        float2 av = a[gid];
        float2 bv = b[gid];
        // (av.x + i*av.y) * (bv.x - i*bv.y)
        sdata[tid] = float2(av.x * bv.x + av.y * bv.y, av.y * bv.x - av.x * bv.y);
    } else {
        sdata[tid] = float2(0.0f, 0.0f);
    }
    threadgroup_barrier(mem_flags::mem_threadgroup);

    for (uint s = tgSize >> 1; s > 0; s >>= 1) {
        if (tid < s) {
            sdata[tid] = c_add(sdata[tid], sdata[tid + s]);
        }
        threadgroup_barrier(mem_flags::mem_threadgroup);
    }

    if (tid == 0) {
        partials[tgIdx] = sdata[0];
    }
}
