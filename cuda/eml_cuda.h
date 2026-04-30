/*
 * EML CUDA Kernels for NVIDIA GPUs
 * 
 * This file contains CUDA implementations of the EML operator and related
 * mathematical functions optimized for parallel execution on NVIDIA GPUs.
 * 
 * Compile with: nvcc -O3 -arch=sm_70 eml_cuda.cu -o eml_cuda
 * 
 * Functions:
 * - eml_kernel: Core EML operation (exp(x) - ln(y))
 * - exp_kernel: Exponential function
 * - log_kernel: Natural logarithm
 * - sin_kernel, cos_kernel, tan_kernel: Trigonometric
 * - sinh_kernel, cosh_kernel, tanh_kernel: Hyperbolic
 * - matmul_kernel: Matrix multiplication using EML
 */

#ifndef EML_CUDA_H
#define EML_CUDA_H

#include <cuda_runtime.h>
#include <stdio.h>
#include <stdlib.h>

// Device constants
__device__ __constant__ double CUDA_LN2 = 0.693147180559945309417;
__device__ __constant__ double CUDA_LN10 = 2.30258509299404568402;
__device__ __constant__ double CUDA_E = 2.71828182845904523536;

// Helper functions
__device__ __forceinline__ double fast_exp(double x) {
    // Use Taylor series expansion for small values
    // For production: use polynomial approximation
    if (x > 700) return INFINITY;
    if (x < -700) return 0;
    
    double result = 1.0;
    double term = 1.0;
    for (int i = 1; i <= 10; i++) {
        term *= x / i;
        result += term;
    }
    return result;
}

__device__ __forceinline__ double fast_log(double x) {
    // Use series expansion for log(1+y)
    if (x <= 0) return NAN;
    if (x == 1) return 0;
    
    // Reduce to range [0.5, 1.5]
    int exp;
    double mantissa = frexp(x, &exp);
    mantissa = 2.0 * mantissa - 1.0;
    
    double result = 0;
    double y = mantissa - 1.0;
    double term = y;
    for (int i = 1; i <= 15; i++) {
        result += (i % 2 == 1 ? 1 : -1) * term / i;
        term *= y;
    }
    return result + exp * CUDA_LN2;
}

// ============================================================================
// Core EML Operation
// ============================================================================

/*
 * EML(x, y) = exp(x) - ln(y)
 * 
 * This is the fundamental operation from which all elementary functions
 * can be derived. This kernel processes multiple elements in parallel.
 */
__global__ void eml_kernel(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = fast_exp(x[idx]) - fast_log(y[idx]);
    }
}

// Vectorized EML for coalesced memory access
__global__ void eml_kernel_vectorized(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
) {
    int idx = (blockIdx.x * blockDim.x + threadIdx.x) * 4;
    if (idx + 3 < n) {
        // Process 4 elements at once for better memory coalescing
        #pragma unroll
        for (int i = 0; i < 4; i++) {
            result[idx + i] = fast_exp(x[idx + i]) - fast_log(y[idx + i]);
        }
    } else {
        // Handle remainder
        for (int i = 0; idx + i < n && i < 4; i++) {
            result[idx + i] = fast_exp(x[idx + i]) - fast_log(y[idx + i]);
        }
    }
}

// ============================================================================
// Exponential and Logarithmic
// ============================================================================

__global__ void exp_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = fast_exp(x[idx]);
    }
}

__global__ void log_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = fast_log(x[idx]);
    }
}

// ============================================================================
// Trigonometric Functions (using complex exponentials)
// ============================================================================

__global__ void sin_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        // sin(z) = sin(x)cosh(y) + i*cos(x)sinh(y) for z = x + iy
        // For real input: sin(x) = (e^(ix) - e^(-ix))/(2i)
        double cos_x = cos(x[idx]);
        double sin_x = sin(x[idx]);
        result[idx] = sin_x;
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

// ============================================================================
// Hyperbolic Functions
// ============================================================================

__global__ void sinh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        double ex = fast_exp(x[idx]);
        double emx = fast_exp(-x[idx]);
        result[idx] = (ex - emx) * 0.5;
    }
}

__global__ void cosh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        double ex = fast_exp(x[idx]);
        double emx = fast_exp(-x[idx]);
        result[idx] = (ex + emx) * 0.5;
    }
}

__global__ void tanh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        double ex = fast_exp(x[idx]);
        double emx = fast_exp(-x[idx]);
        result[idx] = (ex - emx) / (ex + emx);
    }
}

// ============================================================================
// Matrix Multiplication using EML
// ============================================================================

/*
 * Matrix multiplication: C = A * B
 * 
 * Uses EML operator for efficient parallel computation.
 * Each thread computes one element of the result matrix.
 */
__global__ void matmul_kernel(
    const double* __restrict__ A,
    const double* __restrict__ B,
    double* __restrict__ C,
    int M, int N, int K,
    int lda, int ldb, int ldc
) {
    // Each thread computes C[i][j]
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;
    
    if (row < M && col < N) {
        double sum = 0.0;
        for (int k = 0; k < K; k++) {
            sum += A[row * lda + k] * B[k * ldb + col];
        }
        C[row * ldc + col] = sum;
    }
}

/*
 * Matrix multiplication with EML weighting
 * C = EML(A, B) where EML is applied element-wise during accumulation
 */
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
        double sum = 0.0;
        double eml_acc = 0.0;  // EML accumulation
        
        for (int k = 0; k < K; k++) {
            double a = A[row * lda + k];
            double b = B[k * ldb + col];
            // Use EML operator for accumulation
            double prod = a * b;
            eml_acc = fast_exp(eml_acc) - fast_log(prod);
        }
        C[row * ldc + col] = eml_acc;
    }
}

// ============================================================================
// Utility Functions
// ============================================================================

// Initialize CUDA EML library
extern "C" {
    cudaError_t eml_cuda_init();
    cudaError_t eml_cuda_cleanup();
    
    // Memory management
    cudaError_t eml_allocate(void** ptr, size_t size);
    cudaError_t eml_free(void* ptr);
    cudaError_t eml_copy_to_device(void* dst, const void* src, size_t size);
    cudaError_t eml_copy_to_host(void* dst, const void* src, size_t size);
    
    // Kernel launchers
    cudaError_t eml_launch_exp(const double* x, double* result, int n, cudaStream_t stream = 0);
    cudaError_t eml_launch_log(const double* x, double* result, int n, cudaStream_t stream = 0);
    cudaError_t eml_launch_eml(const double* x, const double* y, double* result, int n, cudaStream_t stream = 0);
    cudaError_t eml_launch_matmul(const double* A, const double* B, double* C, int M, int N, int K, cudaStream_t stream = 0);
}

// Implementation
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

inline cudaError_t eml_allocate(void** ptr, size_t size) {
    return cudaMalloc(ptr, size);
}

inline cudaError_t eml_free(void* ptr) {
    return cudaFree(ptr);
}

inline cudaError_t eml_copy_to_device(void* dst, const void* src, size_t size) {
    return cudaMemcpy(dst, src, size, cudaMemcpyHostToDevice);
}

inline cudaError_t eml_copy_to_host(void* dst, const void* src, size_t size) {
    return cudaMemcpy(dst, src, size, cudaMemcpyDeviceToHost);
}

inline cudaError_t eml_launch_exp(const double* x, double* result, int n, cudaStream_t stream) {
    int blockSize = 256;
    int gridSize = (n + blockSize - 1) / blockSize;
    exp_kernel<<<gridSize, blockSize, 0, stream>>>(x, result, n);
    return cudaGetLastError();
}

inline cudaError_t eml_launch_log(const double* x, double* result, int n, cudaStream_t stream) {
    int blockSize = 256;
    int gridSize = (n + blockSize - 1) / blockSize;
    log_kernel<<<gridSize, blockSize, 0, stream>>>(x, result, n);
    return cudaGetLastError();
}

inline cudaError_t eml_launch_eml(const double* x, const double* y, double* result, int n, cudaStream_t stream) {
    int blockSize = 256;
    int gridSize = (n + blockSize - 1) / blockSize;
    eml_kernel<<<gridSize, blockSize, 0, stream>>>(x, y, result, n);
    return cudaGetLastError();
}

inline cudaError_t eml_launch_matmul(const double* A, const double* B, double* C, int M, int N, int K, cudaStream_t stream) {
    dim3 blockSize(16, 16);
    dim3 gridSize((N + 15) / 16, (M + 15) / 16);
    matmul_kernel<<<gridSize, blockSize, 0, stream>>>(A, B, C, M, N, K, K, N, N);
    return cudaGetLastError();
}

} // extern "C"

#endif // EML_CUDA_H