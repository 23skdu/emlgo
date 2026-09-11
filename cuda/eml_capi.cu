// Pure C API implementation bridging to CUDA kernels.
// Compiled with nvcc to a shared library: libeml_capi.so

#include "eml_capi.h"
#include "eml_cuda.h"
#include <cuda_runtime.h>
#include <cuComplex.h>
#include <string.h>
#include <stdlib.h>

// ---------- Lifecycle ----------

int eml_init(void) {
    int deviceCount;
    cudaError_t err = cudaGetDeviceCount(&deviceCount);
    if (err != cudaSuccess) return (int)err;
    if (deviceCount == 0) return 1;
    err = cudaSetDevice(0);
    return (int)err;
}

int eml_cleanup(void) {
    return (int)cudaDeviceReset();
}

// ---------- Device Query ----------

int eml_get_device_count(int* count) {
    cudaError_t err = cudaGetDeviceCount(count);
    return (int)err;
}

int eml_get_device_name(int device_id, char* name, int max_len) {
    if (name == NULL || max_len <= 0) return -1;
    cudaDeviceProp props;
    cudaError_t err = cudaGetDeviceProperties(&props, device_id);
    if (err != cudaSuccess) return (int)err;
    strncpy(name, props.name, max_len - 1);
    name[max_len - 1] = '\0';
    return 0;
}

int eml_get_device_props(
    int device_id,
    int* compute_major,
    int* compute_minor,
    long long* memory_bytes,
    int* max_threads_per_block,
    int* warp_size,
    int* clock_rate_khz
) {
    cudaDeviceProp props;
    cudaError_t err = cudaGetDeviceProperties(&props, device_id);
    if (err != cudaSuccess) return (int)err;

    if (compute_major) *compute_major = props.major;
    if (compute_minor) *compute_minor = props.minor;
    if (memory_bytes) *memory_bytes = (long long)props.totalGlobalMem;
    if (max_threads_per_block) *max_threads_per_block = props.maxThreadsPerBlock;
    if (warp_size) *warp_size = props.warpSize;
    if (clock_rate_khz) *clock_rate_khz = props.clockRate;

    return 0;
}

// ---------- Memory Management ----------

int eml_allocate(void** ptr, long long size) {
    return (int)cudaMalloc(ptr, (size_t)size);
}

int eml_free(void* ptr) {
    return (int)cudaFree(ptr);
}

int eml_copy_to_device(void* dst, const void* src, long long size) {
    return (int)cudaMemcpy(dst, src, (size_t)size, cudaMemcpyHostToDevice);
}

int eml_copy_to_host(void* dst, const void* src, long long size) {
    return (int)cudaMemcpy(dst, src, (size_t)size, cudaMemcpyDeviceToHost);
}

int eml_sync_device(void) {
    return (int)cudaDeviceSynchronize();
}

int eml_memset(void* ptr, int value, long long size) {
    return (int)cudaMemset(ptr, value, (size_t)size);
}

int eml_memset_stream(void* ptr, int value, long long size, long long stream) {
    return (int)cudaMemsetAsync(ptr, value, (size_t)size, (cudaStream_t)stream);
}

// ---------- Pinned Memory ----------

int eml_allocate_pinned(void** ptr, long long size) {
    return (int)cudaHostAlloc(ptr, (size_t)size, cudaHostAllocDefault);
}

int eml_free_pinned(void* ptr) {
    return (int)cudaFreeHost(ptr);
}

// ---------- Async Streams ----------

long long eml_create_stream(void) {
    cudaStream_t stream;
    cudaError_t err = cudaStreamCreate(&stream);
    if (err != cudaSuccess) return 0;
    return (long long)stream;
}

int eml_destroy_stream(long long stream) {
    if (stream == 0) return 0;
    return (int)cudaStreamDestroy((cudaStream_t)stream);
}

int eml_sync_stream(long long stream) {
    if (stream == 0) return 0;
    return (int)cudaStreamSynchronize((cudaStream_t)stream);
}

// ---------- Helper Template Launchers ----------

template<typename T, typename KernelFunc>
static inline int launch_unary_stream(const void* x, void* result, int n, int block_size, cudaStream_t stream, KernelFunc kernel) {
    if (n <= 0) return 0;
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    kernel<<<grid_size, block_size, 0, stream>>>((const T*)x, (T*)result, n);
    return (int)cudaGetLastError();
}

template<typename T, typename KernelFunc>
static inline int launch_binary_stream(const void* x, const void* y, void* result, int n, int block_size, cudaStream_t stream, KernelFunc kernel) {
    if (n <= 0) return 0;
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    kernel<<<grid_size, block_size, 0, stream>>>((const T*)x, (const T*)y, (T*)result, n);
    return (int)cudaGetLastError();
}

template<typename T, typename KernelFunc>
static inline int launch_dot_stream(const void* a, const void* b, void* result, int n, int block_size, cudaStream_t stream, KernelFunc kernel) {
    if (n <= 0) return 0;
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    if (grid_size > 256) grid_size = 256;
    if (grid_size == 0) grid_size = 1;

    cudaError_t err = cudaMemsetAsync(result, 0, sizeof(T), stream);
    if (err != cudaSuccess) return (int)err;

    kernel<<<grid_size, block_size, 0, stream>>>((const T*)a, (const T*)b, (T*)result, n);
    return (int)cudaGetLastError();
}

// ============================================================================
// Float64 Kernel Launchers
// ============================================================================

int eml_launch_exp(const void* x, void* result, int n, int block_size) {
    return eml_launch_exp_stream(x, result, n, block_size, 0);
}
int eml_launch_log(const void* x, void* result, int n, int block_size) {
    return eml_launch_log_stream(x, result, n, block_size, 0);
}
int eml_launch_sin(const void* x, void* result, int n, int block_size) {
    return eml_launch_sin_stream(x, result, n, block_size, 0);
}
int eml_launch_cos(const void* x, void* result, int n, int block_size) {
    return eml_launch_cos_stream(x, result, n, block_size, 0);
}
int eml_launch_tan(const void* x, void* result, int n, int block_size) {
    return eml_launch_tan_stream(x, result, n, block_size, 0);
}
int eml_launch_sinh(const void* x, void* result, int n, int block_size) {
    return eml_launch_sinh_stream(x, result, n, block_size, 0);
}
int eml_launch_cosh(const void* x, void* result, int n, int block_size) {
    return eml_launch_cosh_stream(x, result, n, block_size, 0);
}
int eml_launch_tanh(const void* x, void* result, int n, int block_size) {
    return eml_launch_tanh_stream(x, result, n, block_size, 0);
}
int eml_launch_sqrt(const void* x, void* result, int n, int block_size) {
    return eml_launch_sqrt_stream(x, result, n, block_size, 0);
}
int eml_launch_eml(const void* x, const void* y, void* result, int n, int block_size) {
    return eml_launch_eml_stream(x, y, result, n, block_size, 0);
}
int eml_launch_dot(const void* a, const void* b, void* result, int n, int block_size) {
    return eml_launch_dot_stream(a, b, result, n, block_size, 0);
}

int eml_launch_exp_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, exp_kernel);
}
int eml_launch_log_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, log_kernel);
}
int eml_launch_sin_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, sin_kernel);
}
int eml_launch_cos_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, cos_kernel);
}
int eml_launch_tan_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, tan_kernel);
}
int eml_launch_sinh_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, sinh_kernel);
}
int eml_launch_cosh_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, cosh_kernel);
}
int eml_launch_tanh_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, tanh_kernel);
}
int eml_launch_sqrt_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<double>(x, result, n, block_size, (cudaStream_t)stream, sqrt_kernel);
}
int eml_launch_eml_stream(const void* x, const void* y, void* result, int n, int block_size, long long stream) {
    return launch_binary_stream<double>(x, y, result, n, block_size, (cudaStream_t)stream, eml_kernel);
}
int eml_launch_dot_stream(const void* a, const void* b, void* result, int n, int block_size, long long stream) {
    return launch_dot_stream<double>(a, b, result, n, block_size, (cudaStream_t)stream, dot_kernel_f64);
}

// ============================================================================
// Float32 Kernel Launchers
// ============================================================================

int eml_launch_exp_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_exp_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_log_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_log_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_sin_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_sin_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_cos_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_cos_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_tan_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_tan_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_sinh_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_sinh_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_cosh_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_cosh_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_tanh_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_tanh_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_sqrt_f32(const void* x, void* result, int n, int block_size) {
    return eml_launch_sqrt_stream_f32(x, result, n, block_size, 0);
}
int eml_launch_eml_f32(const void* x, const void* y, void* result, int n, int block_size) {
    return eml_launch_eml_stream_f32(x, y, result, n, block_size, 0);
}
int eml_launch_dot_f32(const void* a, const void* b, void* result, int n, int block_size) {
    return eml_launch_dot_stream_f32(a, b, result, n, block_size, 0);
}

int eml_launch_exp_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, exp_kernel_f32);
}
int eml_launch_log_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, log_kernel_f32);
}
int eml_launch_sin_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, sin_kernel_f32);
}
int eml_launch_cos_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, cos_kernel_f32);
}
int eml_launch_tan_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, tan_kernel_f32);
}
int eml_launch_sinh_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, sinh_kernel_f32);
}
int eml_launch_cosh_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, cosh_kernel_f32);
}
int eml_launch_tanh_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, tanh_kernel_f32);
}
int eml_launch_sqrt_stream_f32(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<float>(x, result, n, block_size, (cudaStream_t)stream, sqrt_kernel_f32);
}
int eml_launch_eml_stream_f32(const void* x, const void* y, void* result, int n, int block_size, long long stream) {
    return launch_binary_stream<float>(x, y, result, n, block_size, (cudaStream_t)stream, eml_kernel_f32);
}
int eml_launch_dot_stream_f32(const void* a, const void* b, void* result, int n, int block_size, long long stream) {
    return launch_dot_stream<float>(a, b, result, n, block_size, (cudaStream_t)stream, dot_kernel_f32);
}

// ============================================================================
// Complex64 Kernel Launchers
// ============================================================================

int eml_launch_exp_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_exp_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_log_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_log_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_sin_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_sin_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_cos_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_cos_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_tan_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_tan_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_sinh_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_sinh_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_cosh_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_cosh_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_tanh_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_tanh_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_sqrt_c64(const void* x, void* result, int n, int block_size) {
    return eml_launch_sqrt_stream_c64(x, result, n, block_size, 0);
}
int eml_launch_eml_c64(const void* x, const void* y, void* result, int n, int block_size) {
    return eml_launch_eml_stream_c64(x, y, result, n, block_size, 0);
}
int eml_launch_dot_c64(const void* a, const void* b, void* result, int n, int block_size) {
    return eml_launch_dot_stream_c64(a, b, result, n, block_size, 0);
}

int eml_launch_exp_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, exp_kernel_c64);
}
int eml_launch_log_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, log_kernel_c64);
}
int eml_launch_sin_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, sin_kernel_c64);
}
int eml_launch_cos_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, cos_kernel_c64);
}
int eml_launch_tan_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, tan_kernel_c64);
}
int eml_launch_sinh_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, sinh_kernel_c64);
}
int eml_launch_cosh_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, cosh_kernel_c64);
}
int eml_launch_tanh_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, tanh_kernel_c64);
}
int eml_launch_sqrt_stream_c64(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuFloatComplex>(x, result, n, block_size, (cudaStream_t)stream, sqrt_kernel_c64);
}
int eml_launch_eml_stream_c64(const void* x, const void* y, void* result, int n, int block_size, long long stream) {
    return launch_binary_stream<cuFloatComplex>(x, y, result, n, block_size, (cudaStream_t)stream, eml_kernel_c64);
}
int eml_launch_dot_stream_c64(const void* a, const void* b, void* result, int n, int block_size, long long stream) {
    return launch_dot_stream<cuFloatComplex>(a, b, result, n, block_size, (cudaStream_t)stream, dot_kernel_c64);
}

// ============================================================================
// Complex128 Kernel Launchers
// ============================================================================

int eml_launch_exp_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_exp_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_log_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_log_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_sin_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_sin_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_cos_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_cos_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_tan_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_tan_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_sinh_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_sinh_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_cosh_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_cosh_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_tanh_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_tanh_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_sqrt_c128(const void* x, void* result, int n, int block_size) {
    return eml_launch_sqrt_stream_c128(x, result, n, block_size, 0);
}
int eml_launch_eml_c128(const void* x, const void* y, void* result, int n, int block_size) {
    return eml_launch_eml_stream_c128(x, y, result, n, block_size, 0);
}
int eml_launch_dot_c128(const void* a, const void* b, void* result, int n, int block_size) {
    return eml_launch_dot_stream_c128(a, b, result, n, block_size, 0);
}

int eml_launch_exp_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, exp_kernel_c128);
}
int eml_launch_log_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, log_kernel_c128);
}
int eml_launch_sin_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, sin_kernel_c128);
}
int eml_launch_cos_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, cos_kernel_c128);
}
int eml_launch_tan_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, tan_kernel_c128);
}
int eml_launch_sinh_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, sinh_kernel_c128);
}
int eml_launch_cosh_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, cosh_kernel_c128);
}
int eml_launch_tanh_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, tanh_kernel_c128);
}
int eml_launch_sqrt_stream_c128(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_unary_stream<cuDoubleComplex>(x, result, n, block_size, (cudaStream_t)stream, sqrt_kernel_c128);
}
int eml_launch_eml_stream_c128(const void* x, const void* y, void* result, int n, int block_size, long long stream) {
    return launch_binary_stream<cuDoubleComplex>(x, y, result, n, block_size, (cudaStream_t)stream, eml_kernel_c128);
}
int eml_launch_dot_stream_c128(const void* a, const void* b, void* result, int n, int block_size, long long stream) {
    return launch_dot_stream<cuDoubleComplex>(a, b, result, n, block_size, (cudaStream_t)stream, dot_kernel_c128);
}
