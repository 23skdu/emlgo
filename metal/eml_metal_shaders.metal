#include <metal_stdlib>
using namespace metal;

kernel void kernel_exp(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = exp(x[id]);
}

kernel void kernel_log(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = log(x[id]);
}

kernel void kernel_sin(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = sin(x[id]);
}

kernel void kernel_cos(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = cos(x[id]);
}

kernel void kernel_tan(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = tan(x[id]);
}

kernel void kernel_sinh(device const float *x [[buffer(0)]],
                        device float *result [[buffer(1)]],
                        uint id [[thread_position_in_grid]]) {
    result[id] = sinh(x[id]);
}

kernel void kernel_cosh(device const float *x [[buffer(0)]],
                        device float *result [[buffer(1)]],
                        uint id [[thread_position_in_grid]]) {
    result[id] = cosh(x[id]);
}

kernel void kernel_tanh(device const float *x [[buffer(0)]],
                        device float *result [[buffer(1)]],
                        uint id [[thread_position_in_grid]]) {
    result[id] = tanh(x[id]);
}

kernel void kernel_sqrt(device const float *x [[buffer(0)]],
                        device float *result [[buffer(1)]],
                        uint id [[thread_position_in_grid]]) {
    result[id] = sqrt(x[id]);
}

kernel void kernel_abs(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = abs(x[id]);
}

kernel void kernel_neg(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = -x[id];
}

kernel void kernel_inv(device const float *x [[buffer(0)]],
                       device float *result [[buffer(1)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = 1.0 / x[id];
}

kernel void kernel_add(device const float *a [[buffer(0)]],
                       device const float *b [[buffer(1)]],
                       device float *result [[buffer(2)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = a[id] + b[id];
}

kernel void kernel_sub(device const float *a [[buffer(0)]],
                       device const float *b [[buffer(1)]],
                       device float *result [[buffer(2)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = a[id] - b[id];
}

kernel void kernel_mul(device const float *a [[buffer(0)]],
                       device const float *b [[buffer(1)]],
                       device float *result [[buffer(2)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = a[id] * b[id];
}

kernel void kernel_div(device const float *a [[buffer(0)]],
                       device const float *b [[buffer(1)]],
                       device float *result [[buffer(2)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = a[id] / b[id];
}

kernel void kernel_fma(device const float *a [[buffer(0)]],
                       device const float *b [[buffer(1)]],
                       device const float *c [[buffer(2)]],
                       device float *result [[buffer(3)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = fma(a[id], b[id], c[id]);
}

kernel void kernel_addScalar(device const float *a [[buffer(0)]],
                             constant float &b [[buffer(1)]],
                             device float *result [[buffer(2)]],
                             uint id [[thread_position_in_grid]]) {
    result[id] = a[id] + b;
}

kernel void kernel_mulScalar(device const float *a [[buffer(0)]],
                             constant float &b [[buffer(1)]],
                             device float *result [[buffer(2)]],
                             uint id [[thread_position_in_grid]]) {
    result[id] = a[id] * b;
}

kernel void kernel_eml(device const float *x [[buffer(0)]],
                       device const float *y [[buffer(1)]],
                       device float *result [[buffer(2)]],
                       uint id [[thread_position_in_grid]]) {
    result[id] = exp(x[id]) - log(y[id]);
}
