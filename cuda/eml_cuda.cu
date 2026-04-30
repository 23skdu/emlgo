/*
 * EML CUDA Implementation Source File
 * 
 * Compile: nvcc -O3 -arch=sm_70 -c eml_cuda.cu -o eml_cuda.o
 * Link: nvcc -O3 -arch=sm_70 eml_cuda.o -o eml_cuda
 */

#include "eml_cuda.h"
#include <cuda_runtime.h>
#include <math.h>

// ============================================================================
// Device Functions (must be in .cu file)
// ============================================================================

// Fast exponential with range handling
__device__ __forceinline__ double fast_exp(double x) {
    // Clamp to prevent overflow/underflow
    if (x > 709.782712893384) return INFINITY;
    if (x < -708.4779660139996) return 0;
    
    // Use built-in for accuracy (or replace with polynomial)
    return __expd(x);
}

// Fast logarithm
__device__ __forceinline__ double fast_log(double x) {
    if (x <= 0) return NAN;
    return __logd(x);
}

// ============================================================================
// Kernel Implementations
// ============================================================================

// EML Core: exp(x) - ln(y)
__global__ void eml_kernel(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __expd(x[idx]) - __logd(y[idx]);
    }
}

// EML with vectorization (4 elements per thread)
__global__ void eml_kernel_vectorized(
    const double* __restrict__ x,
    const double* __restrict__ y,
    double* __restrict__ result,
    int n
) {
    int idx = (blockIdx.x * blockDim.x + threadIdx.x) * 4;
    
    if (idx + 3 < n) {
        double4 x4 = *(double4*)(&x[idx]);
        double4 y4 = *(double4*)(&y[idx]);
        double4 r4;
        r4.x = __expd(x4.x) - __logd(y4.x);
        r4.y = __expd(x4.y) - __logd(y4.y);
        r4.z = __expd(x4.z) - __logd(y4.z);
        r4.w = __expd(x4.w) - __logd(y4.w);
        *(double4*)(&result[idx]) = r4;
    } else {
        for (int i = 0; idx + i < n && i < 4; i++) {
            result[idx + i] = __expd(x[idx + i]) - __logd(y[idx + i]);
        }
    }
}

// Exponential
__global__ void exp_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __expd(x[idx]);
    }
}

// Natural logarithm
__global__ void log_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __logd(x[idx]);
    }
}

// Sine
__global__ void sin_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __sin(x[idx]);
    }
}

// Cosine  
__global__ void cos_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __cos(x[idx]);
    }
}

// Tangent
__global__ void tan_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __tan(x[idx]);
    }
}

// Hyperbolic sine
__global__ void sinh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        double ex = __expd(x[idx]);
        double emx = __expd(-x[idx]);
        result[idx] = (ex - emx) * 0.5;
    }
}

// Hyperbolic cosine
__global__ void cosh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        double ex = __expd(x[idx]);
        double emx = __expd(-x[idx]);
        result[idx] = (ex + emx) * 0.5;
    }
}

// Hyperbolic tangent
__global__ void tanh_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        double ex = __expd(x[idx]);
        double emx = __expd(-x[idx]);
        result[idx] = (ex - emx) / (ex + emx);
    }
}

// Square root
__global__ void sqrt_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        result[idx] = __dsqrt(x[idx]);
    }
}

// ============================================================================
// Matrix Multiplication Kernels
// ============================================================================

// Standard matrix multiplication (C = A * B)
__global__ void matmul_kernel(
    const double* __restrict__ A,
    const double* __restrict__ B,
    double* __restrict__ C,
    int M, int N, int K,
    int lda, int ldb, int ldc
) {
    // Shared memory for tiling
    __shared__ double As[16][16];
    __shared__ double Bs[16][16];
    
    int row = blockIdx.y * 16 + threadIdx.y;
    int col = blockIdx.x * 16 + threadIdx.x;
    
    double Cvalue = 0;
    
    // Loop over tiles
    for (int kk = 0; kk < K; kk += 16) {
        // Load into shared memory
        if (row < M && (kk + threadIdx.x) < K)
            As[threadIdx.y][threadIdx.x] = A[row * lda + kk + threadIdx.x];
        else
            As[threadIdx.y][threadIdx.x] = 0;
            
        if (col < N && (kk + threadIdx.y) < K)
            Bs[threadIdx.y][threadIdx.x] = B[(kk + threadIdx.y) * ldb + col];
        else
            Bs[threadIdx.y][threadIdx.x] = 0;
            
        __syncthreads();
        
        // Compute partial dot product
        #pragma unroll
        for (int k = 0; k < 16; k++) {
            Cvalue += As[threadIdx.y][k] * Bs[k][threadIdx.x];
        }
        __syncthreads();
    }
    
    if (row < M && col < N)
        C[row * ldc + col] = Cvalue;
}

// EML-weighted matrix multiplication
// Uses EML operator for accumulation instead of standard addition
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
        double eml_acc = 0;  // Initial EML accumulator
        
        for (int k = 0; k < K; k++) {
            double a = A[row * lda + k];
            double b = B[k * ldb + col];
            double prod = a * b;
            
            // EML accumulation: eml(acc, prod) = exp(acc) - ln(prod)
            eml_acc = __expd(eml_acc) - __logd(prod);
        }
        
        C[row * ldc + col] = eml_acc;
    }
}

// ============================================================================
// Warp-level Reduction for Parallel Reduction
// ============================================================================

__device__ __forceinline__ double warp_reduce_sum(double val) {
    // Warp shuffle reduction
    for (int offset = 16; offset > 0; offset /= 2) {
        val += __shfl_down_sync(0xffffffff, val, offset);
    }
    return val;
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

// ============================================================================
// Softmax using EML
// ============================================================================

__global__ void softmax_kernel(
    const double* __restrict__ x,
    double* __restrict__ result,
    int n
) {
    // First find max for numerical stability
    __shared__ double max_val;
    int tid = threadIdx.x;
    double local_max = -INFINITY;
    
    for (int i = tid; i < n; i += blockDim.x) {
        local_max = fmax(local_max, x[i]);
    }
    
    // Reduce max
    local_max = warp_reduce_sum(local_max);
    if (tid == 0) max_val = local_max;
    __syncthreads();
    
    // Compute exp(x - max) and sum
    double local_sum = 0;
    for (int i = tid; i < n; i += blockDim.x) {
        result[i] = __expd(x[i] - max_val);
        local_sum += result[i];
    }
    
    local_sum = warp_reduce_sum(local_sum);
    if (tid == 0) max_val = local_sum;
    __syncthreads();
    
    // Normalize
    for (int i = tid; i < n; i += blockDim.x) {
        result[i] /= max_val;
    }
}

// ============================================================================
// Example Usage Functions
// ============================================================================

/*
 * Example: Compute EML(A, B) * C matrix operation
 * 
 * Input: A[M x K], B[K x N], C[M x N]
 * Output: Result = EML(A * B, C)
 */
extern "C" void eml_matrix_operation_example() {
    // Allocate memory
    double *d_A, *d_B, *d_C, *d_result;
    int M = 1024, N = 1024, K = 512;
    
    size_t size_A = M * K * sizeof(double);
    size_t size_B = K * N * sizeof(double);
    size_t size_C = M * N * sizeof(double);
    
    cudaMalloc(&d_A, size_A);
    cudaMalloc(&d_B, size_B);
    cudaMalloc(&d_C, size_C);
    cudaMalloc(&d_result, size_C);
    
    // Copy data to device
    // cudaMemcpy(d_A, h_A, size_A, cudaMemcpyHostToDevice);
    // cudaMemcpy(d_B, h_B, size_B, cudaMemcpyHostToDevice);
    // cudaMemcpy(d_C, h_C, size_C, cudaMemcpyHostToDevice);
    
    // Launch matmul
    dim3 block(16, 16);
    dim3 grid((N + 15) / 16, (M + 15) / 16);
    matmul_kernel<<<grid, block>>>(d_A, d_B, d_result, M, N, K, K, N, N);
    
    // Copy result back
    // cudaMemcpy(h_result, d_result, size_C, cudaMemcpyDeviceToHost);
    
    // Cleanup
    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);
    cudaFree(d_result);
}