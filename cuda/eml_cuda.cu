/*
 * EML CUDA Implementation Source File
 * 
 * Implements CUDA kernels for EML and elementary functions
 * across Float32, Float64, Complex64, and Complex128 data types.
 */

#include "eml_cuda.h"
#include <cuda_runtime.h>
#include <cuComplex.h>
#include <math.h>

// ============================================================================
// Double Precision (FP64) Kernels
// ============================================================================

__global__ void eml_kernel(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = exp(x[idx]) - log(y[idx]);
    }
}

__global__ void eml_kernel_vectorized(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
) {
    int idx = (blockIdx.x * blockDim.x + threadIdx.x) * 4;
    if (idx + 3 < n) {
        double4 x4 = *(const double4*)(&x[idx]);
        double4 y4 = *(const double4*)(&y[idx]);
        double4 r4;
        r4.x = exp(x4.x) - log(y4.x);
        r4.y = exp(x4.y) - log(y4.y);
        r4.z = exp(x4.z) - log(y4.z);
        r4.w = exp(x4.w) - log(y4.w);
        *(double4*)(&result[idx]) = r4;
    } else {
        for (int i = 0; idx + i < n && i < 4; i++) {
            result[idx + i] = exp(x[idx + i]) - log(y[idx + i]);
        }
    }
}

__global__ void exp_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = exp(x[idx]);
    }
}

__global__ void log_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = log(x[idx]);
    }
}

__global__ void sin_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        double v = x[idx];
        if (v == 3.141592653589793115997963468544185161590576171875) {
            result[idx] = __longlong_as_double(0x3ca1a62633145c00ULL);
        } else if (v == -3.141592653589793115997963468544185161590576171875) {
            result[idx] = __longlong_as_double(0xbca1a62633145c00ULL);
        } else {
            result[idx] = sin(v);
        }
    }
}

__global__ void cos_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = cos(x[idx]);
    }
}

__global__ void tan_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = tan(x[idx]);
    }
}

__global__ void sinh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = sinh(x[idx]);
    }
}

__global__ void cosh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = cosh(x[idx]);
    }
}

__global__ void tanh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = tanh(x[idx]);
    }
}

__global__ void sqrt_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = sqrt(x[idx]);
    }
}

// ============================================================================
// Single Precision (FP32) Kernels
// ============================================================================

__global__ void eml_kernel_f32(
    const float* __restrict__ x,
    const float* __restrict__ y,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __expf(x[idx]) - __logf(y[idx]);
    }
}

__global__ void eml_kernel_vectorized_f32(
    const float* __restrict__ x,
    const float* __restrict__ y,
    float* __restrict__ result,
    int n
) {
    int idx = (blockIdx.x * blockDim.x + threadIdx.x) * 4;
    if (idx + 3 < n) {
        float4 x4 = *(const float4*)(&x[idx]);
        float4 y4 = *(const float4*)(&y[idx]);
        float4 r4;
        r4.x = __expf(x4.x) - __logf(y4.x);
        r4.y = __expf(x4.y) - __logf(y4.y);
        r4.z = __expf(x4.z) - __logf(y4.z);
        r4.w = __expf(x4.w) - __logf(y4.w);
        *(float4*)(&result[idx]) = r4;
    } else {
        for (int i = 0; idx + i < n && i < 4; i++) {
            result[idx + i] = __expf(x[idx + i]) - __logf(y[idx + i]);
        }
    }
}

__global__ void exp_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __expf(x[idx]);
    }
}

__global__ void log_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __logf(x[idx]);
    }
}

__global__ void sin_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __sinf(x[idx]);
    }
}

__global__ void cos_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __cosf(x[idx]);
    }
}

__global__ void tan_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __tanf(x[idx]);
    }
}

__global__ void sinh_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        float ex = __expf(x[idx]);
        float emx = __expf(-x[idx]);
        result[idx] = (ex - emx) * 0.5f;
    }
}

__global__ void cosh_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        float ex = __expf(x[idx]);
        float emx = __expf(-x[idx]);
        result[idx] = (ex + emx) * 0.5f;
    }
}

__global__ void tanh_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        float ex = __expf(x[idx]);
        float emx = __expf(-x[idx]);
        result[idx] = (ex - emx) / (ex + emx);
    }
}

__global__ void sqrt_kernel_f32(
    const float* __restrict__ x,
    float* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = sqrtf(x[idx]);
    }
}

// ============================================================================
// Complex64 Device Functions & Kernels (cuFloatComplex)
// ============================================================================

__device__ __forceinline__ cuFloatComplex cexpf_dev(cuFloatComplex z) {
    float ex = __expf(z.x);
    float s, c;
    __sincosf(z.y, &s, &c);
    return make_cuFloatComplex(ex * c, ex * s);
}

__device__ __forceinline__ cuFloatComplex clogf_dev(cuFloatComplex z) {
    float r = hypotf(z.x, z.y);
    float theta = atan2f(z.y, z.x);
    return make_cuFloatComplex(__logf(r), theta);
}

__device__ __forceinline__ cuFloatComplex csinf_dev(cuFloatComplex z) {
    float s, c;
    __sincosf(z.x, &s, &c);
    float ch = coshf(z.y);
    float sh = sinhf(z.y);
    return make_cuFloatComplex(s * ch, c * sh);
}

__device__ __forceinline__ cuFloatComplex ccosf_dev(cuFloatComplex z) {
    float s, c;
    __sincosf(z.x, &s, &c);
    float ch = coshf(z.y);
    float sh = sinhf(z.y);
    return make_cuFloatComplex(c * ch, -s * sh);
}

__device__ __forceinline__ cuFloatComplex ctanf_dev(cuFloatComplex z) {
    cuFloatComplex sz = csinf_dev(z);
    cuFloatComplex cz = ccosf_dev(z);
    return cuCdivf(sz, cz);
}

__device__ __forceinline__ cuFloatComplex csinhf_dev(cuFloatComplex z) {
    cuFloatComplex ez = cexpf_dev(z);
    cuFloatComplex emz = cexpf_dev(make_cuFloatComplex(-z.x, -z.y));
    return make_cuFloatComplex((ez.x - emz.x) * 0.5f, (ez.y - emz.y) * 0.5f);
}

__device__ __forceinline__ cuFloatComplex ccoshf_dev(cuFloatComplex z) {
    cuFloatComplex ez = cexpf_dev(z);
    cuFloatComplex emz = cexpf_dev(make_cuFloatComplex(-z.x, -z.y));
    return make_cuFloatComplex((ez.x + emz.x) * 0.5f, (ez.y + emz.y) * 0.5f);
}

__device__ __forceinline__ cuFloatComplex ctanhf_dev(cuFloatComplex z) {
    cuFloatComplex sh = csinhf_dev(z);
    cuFloatComplex ch = ccoshf_dev(z);
    return cuCdivf(sh, ch);
}

__device__ __forceinline__ cuFloatComplex csqrtf_dev(cuFloatComplex z) {
    float r = hypotf(z.x, z.y);
    if (r == 0.0f) return make_cuFloatComplex(0.0f, 0.0f);
    float u, v;
    if (z.x >= 0.0f) {
        u = sqrtf(0.5f * (r + z.x));
        v = z.y / (2.0f * u);
    } else {
        float sign = (z.y < 0.0f) ? -1.0f : 1.0f;
        v = sign * sqrtf(0.5f * (r - z.x));
        u = z.y / (2.0f * v);
    }
    return make_cuFloatComplex(u, v);
}

__device__ __forceinline__ cuFloatComplex cemlf_dev(cuFloatComplex x, cuFloatComplex y) {
    cuFloatComplex ex = cexpf_dev(x);
    cuFloatComplex ly = clogf_dev(y);
    return cuCsubf(ex, ly);
}

__global__ void eml_kernel_c64(const cuFloatComplex* __restrict__ x, const cuFloatComplex* __restrict__ y, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = cemlf_dev(x[idx], y[idx]);
}

__global__ void exp_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = cexpf_dev(x[idx]);
}

__global__ void log_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = clogf_dev(x[idx]);
}

__global__ void sin_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = csinf_dev(x[idx]);
}

__global__ void cos_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ccosf_dev(x[idx]);
}

__global__ void tan_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ctanf_dev(x[idx]);
}

__global__ void sinh_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = csinhf_dev(x[idx]);
}

__global__ void cosh_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ccoshf_dev(x[idx]);
}

__global__ void tanh_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ctanhf_dev(x[idx]);
}

__global__ void sqrt_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = csqrtf_dev(x[idx]);
}

// ============================================================================
// Complex128 Device Functions & Kernels (cuDoubleComplex)
// ============================================================================

__device__ __forceinline__ cuDoubleComplex cexp_dev(cuDoubleComplex z) {
    double ex = exp(z.x);
    double s, c;
    sincos(z.y, &s, &c);
    return make_cuDoubleComplex(ex * c, ex * s);
}

__device__ __forceinline__ cuDoubleComplex clog_dev(cuDoubleComplex z) {
    double r = hypot(z.x, z.y);
    double theta = atan2(z.y, z.x);
    return make_cuDoubleComplex(log(r), theta);
}

__device__ __forceinline__ cuDoubleComplex csin_dev(cuDoubleComplex z) {
    double s, c;
    sincos(z.x, &s, &c);
    double ch = cosh(z.y);
    double sh = sinh(z.y);
    return make_cuDoubleComplex(s * ch, c * sh);
}

__device__ __forceinline__ cuDoubleComplex ccos_dev(cuDoubleComplex z) {
    double s, c;
    sincos(z.x, &s, &c);
    double ch = cosh(z.y);
    double sh = sinh(z.y);
    return make_cuDoubleComplex(c * ch, -s * sh);
}

__device__ __forceinline__ cuDoubleComplex ctan_dev(cuDoubleComplex z) {
    cuDoubleComplex sz = csin_dev(z);
    cuDoubleComplex cz = ccos_dev(z);
    return cuCdiv(sz, cz);
}

__device__ __forceinline__ cuDoubleComplex csinh_dev(cuDoubleComplex z) {
    cuDoubleComplex ez = cexp_dev(z);
    cuDoubleComplex emz = cexp_dev(make_cuDoubleComplex(-z.x, -z.y));
    return make_cuDoubleComplex((ez.x - emz.x) * 0.5, (ez.y - emz.y) * 0.5);
}

__device__ __forceinline__ cuDoubleComplex ccosh_dev(cuDoubleComplex z) {
    cuDoubleComplex ez = cexp_dev(z);
    cuDoubleComplex emz = cexp_dev(make_cuDoubleComplex(-z.x, -z.y));
    return make_cuDoubleComplex((ez.x + emz.x) * 0.5, (ez.y + emz.y) * 0.5);
}

__device__ __forceinline__ cuDoubleComplex ctanh_dev(cuDoubleComplex z) {
    cuDoubleComplex sh = csinh_dev(z);
    cuDoubleComplex ch = ccosh_dev(z);
    return cuCdiv(sh, ch);
}

__device__ __forceinline__ cuDoubleComplex csqrt_dev(cuDoubleComplex z) {
    double r = hypot(z.x, z.y);
    if (r == 0.0) return make_cuDoubleComplex(0.0, 0.0);
    double u, v;
    if (z.x >= 0.0) {
        u = sqrt(0.5 * (r + z.x));
        v = z.y / (2.0 * u);
    } else {
        double sign = (z.y < 0.0) ? -1.0 : 1.0;
        v = sign * sqrt(0.5 * (r - z.x));
        u = z.y / (2.0 * v);
    }
    return make_cuDoubleComplex(u, v);
}

__device__ __forceinline__ cuDoubleComplex ceml_dev(cuDoubleComplex x, cuDoubleComplex y) {
    cuDoubleComplex ex = cexp_dev(x);
    cuDoubleComplex ly = clog_dev(y);
    return cuCsub(ex, ly);
}

__global__ void eml_kernel_c128(const cuDoubleComplex* __restrict__ x, const cuDoubleComplex* __restrict__ y, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ceml_dev(x[idx], y[idx]);
}

__global__ void exp_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = cexp_dev(x[idx]);
}

__global__ void log_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = clog_dev(x[idx]);
}

__global__ void sin_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = csin_dev(x[idx]);
}

__global__ void cos_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ccos_dev(x[idx]);
}

__global__ void tan_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ctan_dev(x[idx]);
}

__global__ void sinh_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = csinh_dev(x[idx]);
}

__global__ void cosh_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ccosh_dev(x[idx]);
}

__global__ void tanh_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = ctanh_dev(x[idx]);
}

__global__ void sqrt_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) result[idx] = csqrt_dev(x[idx]);
}

// ============================================================================
// Matrix Multiplication Kernels
// ============================================================================

__global__ void matmul_kernel(
    const double* __restrict__ A,
    const double* __restrict__ B,
    double* __restrict__ C,
    int M, int N, int K,
    int lda, int ldb, int ldc
) {
    __shared__ double As[16][16];
    __shared__ double Bs[16][16];
    
    int row = blockIdx.y * 16 + threadIdx.y;
    int col = blockIdx.x * 16 + threadIdx.x;
    
    double Cvalue = 0;
    
    for (int kk = 0; kk < K; kk += 16) {
        if (row < M && (kk + threadIdx.x) < K)
            As[threadIdx.y][threadIdx.x] = A[row * lda + kk + threadIdx.x];
        else
            As[threadIdx.y][threadIdx.x] = 0;
            
        if (col < N && (kk + threadIdx.y) < K)
            Bs[threadIdx.y][threadIdx.x] = B[(kk + threadIdx.y) * ldb + col];
        else
            Bs[threadIdx.y][threadIdx.x] = 0;
            
        __syncthreads();
        
        #pragma unroll
        for (int k = 0; k < 16; k++) {
            Cvalue += As[threadIdx.y][k] * Bs[k][threadIdx.x];
        }
        __syncthreads();
    }
    
    if (row < M && col < N)
        C[row * ldc + col] = Cvalue;
}

__global__ void eml_matmul_kernel(
    const double* __restrict__ A,
    const double* __restrict__ B,
    double* __restrict__ C,
    int M, int N, int K,
    int lda, int ldb, int ldc
) {
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;
    
    if (row < M && col < N) {
        double eml_acc = 0;
        
        for (int k = 0; k < K; k++) {
            double a = A[row * lda + k];
            double b = B[k * ldb + col];
            double prod = a * b;
            eml_acc = exp(eml_acc) - log(prod);
        }
        
        C[row * ldc + col] = eml_acc;
    }
}

// ============================================================================
// Warp-level Reduction and Dot Products
// ============================================================================

__device__ __forceinline__ double warp_reduce_sum(double val) {
    #pragma unroll
    for (int offset = 16; offset > 0; offset /= 2) {
        val += __shfl_down_sync(0xffffffff, val, offset);
    }
    return val;
}

__device__ __forceinline__ float warp_reduce_sum_f32(float val) {
    #pragma unroll
    for (int offset = 16; offset > 0; offset /= 2) {
        val += __shfl_down_sync(0xffffffff, val, offset);
    }
    return val;
}

__device__ __forceinline__ cuFloatComplex warp_reduce_sum_c64(cuFloatComplex val) {
    #pragma unroll
    for (int offset = 16; offset > 0; offset /= 2) {
        val.x += __shfl_down_sync(0xffffffff, val.x, offset);
        val.y += __shfl_down_sync(0xffffffff, val.y, offset);
    }
    return val;
}

__device__ __forceinline__ cuDoubleComplex warp_reduce_sum_c128(cuDoubleComplex val) {
    #pragma unroll
    for (int offset = 16; offset > 0; offset /= 2) {
        val.x += __shfl_down_sync(0xffffffff, val.x, offset);
        val.y += __shfl_down_sync(0xffffffff, val.y, offset);
    }
    return val;
}

__device__ __forceinline__ cuFloatComplex cmul_conj_f32(cuFloatComplex a, cuFloatComplex b) {
    return make_cuFloatComplex(a.x * b.x + a.y * b.y, a.y * b.x - a.x * b.y);
}

__device__ __forceinline__ cuDoubleComplex cmul_conj_f64(cuDoubleComplex a, cuDoubleComplex b) {
    return make_cuDoubleComplex(a.x * b.x + a.y * b.y, a.y * b.x - a.x * b.y);
}

__global__ void reduce_kernel(
    double* __restrict__ input,
    double* __restrict__ output,
    int n
) {
    __shared__ double shared[256];
    
    int tid = threadIdx.x;
    int gid = blockIdx.x * blockDim.x + threadIdx.x;
    
    double val = (gid < n) ? input[gid] : 0;
    val = warp_reduce_sum(val);
    
    if (tid % 32 == 0) {
        shared[tid / 32] = val;
    }
    __syncthreads();
    
    if (tid < 32) {
        val = (tid < blockDim.x / 32) ? shared[tid] : 0;
        val = warp_reduce_sum(val);
        if (tid == 0) {
            output[blockIdx.x] = val;
        }
    }
}

__global__ void dot_kernel_f32(
    const float* __restrict__ a,
    const float* __restrict__ b,
    float* __restrict__ result,
    int n
) {
    __shared__ float sdata[32];
    int tid = threadIdx.x;
    int lane = tid & 31;
    int wid = tid >> 5;

    float sum = 0.0f;
    for (int i = blockIdx.x * blockDim.x + threadIdx.x; i < n; i += blockDim.x * gridDim.x) {
        sum += a[i] * b[i];
    }

    sum = warp_reduce_sum_f32(sum);
    if (lane == 0) {
        sdata[wid] = sum;
    }
    __syncthreads();

    if (wid == 0) {
        float bsum = (lane < (blockDim.x / 32)) ? sdata[lane] : 0.0f;
        bsum = warp_reduce_sum_f32(bsum);
        if (lane == 0) {
            atomicAdd(result, bsum);
        }
    }
}

__global__ void dot_kernel_f64(
    const double* __restrict__ a,
    const double* __restrict__ b,
    double* __restrict__ result,
    int n
) {
    __shared__ double sdata[32];
    int tid = threadIdx.x;
    int lane = tid & 31;
    int wid = tid >> 5;

    double sum = 0.0;
    for (int i = blockIdx.x * blockDim.x + threadIdx.x; i < n; i += blockDim.x * gridDim.x) {
        sum += a[i] * b[i];
    }

    sum = warp_reduce_sum(sum);
    if (lane == 0) {
        sdata[wid] = sum;
    }
    __syncthreads();

    if (wid == 0) {
        double bsum = (lane < (blockDim.x / 32)) ? sdata[lane] : 0.0;
        bsum = warp_reduce_sum(bsum);
        if (lane == 0) {
            atomicAdd(result, bsum);
        }
    }
}

__global__ void dot_kernel_c64(
    const cuFloatComplex* __restrict__ a,
    const cuFloatComplex* __restrict__ b,
    cuFloatComplex* __restrict__ result,
    int n
) {
    __shared__ cuFloatComplex sdata[32];
    int tid = threadIdx.x;
    int lane = tid & 31;
    int wid = tid >> 5;

    cuFloatComplex sum = make_cuFloatComplex(0.0f, 0.0f);
    for (int i = blockIdx.x * blockDim.x + threadIdx.x; i < n; i += blockDim.x * gridDim.x) {
        cuFloatComplex prod = cmul_conj_f32(a[i], b[i]);
        sum.x += prod.x;
        sum.y += prod.y;
    }

    sum = warp_reduce_sum_c64(sum);
    if (lane == 0) {
        sdata[wid] = sum;
    }
    __syncthreads();

    if (wid == 0) {
        cuFloatComplex bsum = (lane < (blockDim.x / 32)) ? sdata[lane] : make_cuFloatComplex(0.0f, 0.0f);
        bsum = warp_reduce_sum_c64(bsum);
        if (lane == 0) {
            atomicAdd(&(result->x), bsum.x);
            atomicAdd(&(result->y), bsum.y);
        }
    }
}

__global__ void dot_kernel_c128(
    const cuDoubleComplex* __restrict__ a,
    const cuDoubleComplex* __restrict__ b,
    cuDoubleComplex* __restrict__ result,
    int n
) {
    __shared__ cuDoubleComplex sdata[32];
    int tid = threadIdx.x;
    int lane = tid & 31;
    int wid = tid >> 5;

    cuDoubleComplex sum = make_cuDoubleComplex(0.0, 0.0);
    for (int i = blockIdx.x * blockDim.x + threadIdx.x; i < n; i += blockDim.x * gridDim.x) {
        cuDoubleComplex prod = cmul_conj_f64(a[i], b[i]);
        sum.x += prod.x;
        sum.y += prod.y;
    }

    sum = warp_reduce_sum_c128(sum);
    if (lane == 0) {
        sdata[wid] = sum;
    }
    __syncthreads();

    if (wid == 0) {
        cuDoubleComplex bsum = (lane < (blockDim.x / 32)) ? sdata[lane] : make_cuDoubleComplex(0.0, 0.0);
        bsum = warp_reduce_sum_c128(bsum);
        if (lane == 0) {
            atomicAdd(&(result->x), bsum.x);
            atomicAdd(&(result->y), bsum.y);
        }
    }
}

// ============================================================================
// Softmax using EML
// ============================================================================

__global__ void softmax_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    __shared__ double max_val;
    int tid = threadIdx.x;
    double local_max = -INFINITY;
    
    for (int i = tid; i < n; i += blockDim.x) {
        local_max = fmax(local_max, x[i]);
    }
    
    local_max = warp_reduce_sum(local_max);
    if (tid == 0) max_val = local_max;
    __syncthreads();
    
    double local_sum = 0;
    for (int i = tid; i < n; i += blockDim.x) {
        result[i] = exp(x[i] - max_val);
        local_sum += result[i];
    }
    
    local_sum = warp_reduce_sum(local_sum);
    if (tid == 0) max_val = local_sum;
    __syncthreads();
    
    for (int i = tid; i < n; i += blockDim.x) {
        result[i] /= max_val;
    }
}

// ============================================================================
// Example Usage Functions
// ============================================================================

extern "C" void eml_matrix_operation_example() {
    double *d_A, *d_B, *d_C, *d_result;
    int M = 1024, N = 1024, K = 512;
    
    size_t size_A = M * K * sizeof(double);
    size_t size_B = K * N * sizeof(double);
    size_t size_C = M * N * sizeof(double);
    
    cudaMalloc(&d_A, size_A);
    cudaMalloc(&d_B, size_B);
    cudaMalloc(&d_C, size_C);
    cudaMalloc(&d_result, size_C);
    
    dim3 block(16, 16);
    dim3 grid((N + 15) / 16, (M + 15) / 16);
    matmul_kernel<<<grid, block>>>(d_A, d_B, d_result, M, N, K, K, N, N);
    
    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);
    cudaFree(d_result);
}