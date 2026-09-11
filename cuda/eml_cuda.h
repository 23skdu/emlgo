/*
 * EML CUDA Kernels for NVIDIA GPUs
 * 
 * This header declares CUDA kernels for EML operations across
 * Float32, Float64, Complex64, and Complex128 data types.
 */

#ifndef EML_CUDA_H
#define EML_CUDA_H

#include <cuda_runtime.h>
#include <cuComplex.h>
#include <stdio.h>
#include <stdlib.h>

// ============================================================================
// Double Precision (FP64) Kernels
// ============================================================================

__global__ void eml_kernel(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
);

__global__ void eml_kernel_vectorized(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
);

__global__ void exp_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void log_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void sin_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void cos_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void tan_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void sinh_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void cosh_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void tanh_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void sqrt_kernel(const double* __restrict__ x, double* __restrict__ result, int n);
__global__ void dot_kernel_f64(const double* __restrict__ a, const double* __restrict__ b, double* __restrict__ result, int n);
__global__ void matmul_kernel(
    const double* __restrict__ A,
    const double* __restrict__ B,
    double* __restrict__ C,
    int M, int N, int K,
    int lda, int ldb, int ldc
);
__global__ void eml_matmul_kernel(
    const double* __restrict__ A,
    const double* __restrict__ B,
    double* __restrict__ C,
    int M, int N, int K,
    int lda, int ldb, int ldc
);

// ============================================================================
// Single Precision (FP32) Kernels
// ============================================================================

__global__ void eml_kernel_f32(
    const float* __restrict__ x,
    const float* __restrict__ y,
    float* __restrict__ result,
    int n
);

__global__ void eml_kernel_vectorized_f32(
    const float* __restrict__ x,
    const float* __restrict__ y,
    float* __restrict__ result,
    int n
);

__global__ void exp_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void log_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void sin_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void cos_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void tan_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void sinh_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void cosh_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void tanh_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void sqrt_kernel_f32(const float* __restrict__ x, float* __restrict__ result, int n);
__global__ void dot_kernel_f32(const float* __restrict__ a, const float* __restrict__ b, float* __restrict__ result, int n);

// ============================================================================
// Complex64 (cuFloatComplex) Kernels
// ============================================================================

__global__ void eml_kernel_c64(const cuFloatComplex* __restrict__ x, const cuFloatComplex* __restrict__ y, cuFloatComplex* __restrict__ result, int n);
__global__ void exp_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void log_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void sin_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void cos_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void tan_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void sinh_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void cosh_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void tanh_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void sqrt_kernel_c64(const cuFloatComplex* __restrict__ x, cuFloatComplex* __restrict__ result, int n);
__global__ void dot_kernel_c64(const cuFloatComplex* __restrict__ a, const cuFloatComplex* __restrict__ b, cuFloatComplex* __restrict__ result, int n);

// ============================================================================
// Complex128 (cuDoubleComplex) Kernels
// ============================================================================

__global__ void eml_kernel_c128(const cuDoubleComplex* __restrict__ x, const cuDoubleComplex* __restrict__ y, cuDoubleComplex* __restrict__ result, int n);
__global__ void exp_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void log_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void sin_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void cos_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void tan_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void sinh_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void cosh_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void tanh_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void sqrt_kernel_c128(const cuDoubleComplex* __restrict__ x, cuDoubleComplex* __restrict__ result, int n);
__global__ void dot_kernel_c128(const cuDoubleComplex* __restrict__ a, const cuDoubleComplex* __restrict__ b, cuDoubleComplex* __restrict__ result, int n);

// ============================================================================
// Utility Functions (C++ / CUDA)
// ============================================================================

inline cudaError_t eml_cuda_init() {
    int deviceCount;
    cudaError_t err = cudaGetDeviceCount(&deviceCount);
    if (err != cudaSuccess) return err;
    if (deviceCount == 0) return cudaErrorNoDevice;
    return cudaSetDevice(0);
}

inline cudaError_t eml_cuda_cleanup() {
    return cudaDeviceReset();
}

inline cudaError_t eml_cuda_allocate(void** ptr, size_t size) {
    return cudaMalloc(ptr, size);
}

inline cudaError_t eml_cuda_free(void* ptr) {
    return cudaFree(ptr);
}

inline cudaError_t eml_cuda_copy_to_device(void* dst, const void* src, size_t size) {
    return cudaMemcpy(dst, src, size, cudaMemcpyHostToDevice);
}

inline cudaError_t eml_cuda_copy_to_host(void* dst, const void* src, size_t size) {
    return cudaMemcpy(dst, src, size, cudaMemcpyDeviceToHost);
}

inline cudaError_t eml_cuda_launch_exp(const double* x, double* result, int n, cudaStream_t stream = 0) {
    int blockSize = 256;
    int gridSize = (n + blockSize - 1) / blockSize;
    exp_kernel<<<gridSize, blockSize, 0, stream>>>(x, result, n);
    return cudaGetLastError();
}

inline cudaError_t eml_cuda_launch_log(const double* x, double* result, int n, cudaStream_t stream = 0) {
    int blockSize = 256;
    int gridSize = (n + blockSize - 1) / blockSize;
    log_kernel<<<gridSize, blockSize, 0, stream>>>(x, result, n);
    return cudaGetLastError();
}

inline cudaError_t eml_cuda_launch_eml(const double* x, const double* y, double* result, int n, cudaStream_t stream = 0) {
    int blockSize = 256;
    int gridSize = (n + blockSize - 1) / blockSize;
    eml_kernel<<<gridSize, blockSize, 0, stream>>>(x, y, result, n);
    return cudaGetLastError();
}

inline cudaError_t eml_cuda_launch_matmul(const double* A, const double* B, double* C, int M, int N, int K, cudaStream_t stream = 0) {
    dim3 blockSize(16, 16);
    dim3 gridSize((N + 15) / 16, (M + 15) / 16);
    matmul_kernel<<<gridSize, blockSize, 0, stream>>>(A, B, C, M, N, K, K, N, N);
    return cudaGetLastError();
}

#endif // EML_CUDA_H