// Pure C API implementation bridging to CUDA kernels.
// Compiled with nvcc to a shared library: libeml_capi.so

#include "eml_capi.h"
#include "eml_cuda.h"
#include <cuda_runtime.h>
#include <string.h>
#include <stdlib.h>

// Forward declarations of kernels defined in eml_cuda.cu
// (these are global functions in the .cu file, not in the header)
extern __global__ void exp_kernel(const double*, double*, int);
extern __global__ void log_kernel(const double*, double*, int);
extern __global__ void sin_kernel(const double*, double*, int);
extern __global__ void cos_kernel(const double*, double*, int);
extern __global__ void tan_kernel(const double*, double*, int);
extern __global__ void sinh_kernel(const double*, double*, int);
extern __global__ void cosh_kernel(const double*, double*, int);
extern __global__ void tanh_kernel(const double*, double*, int);
extern __global__ void sqrt_kernel(const double*, double*, int);
extern __global__ void eml_kernel(const double*, const double*, double*, int);

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

// ---------- Kernel Launchers (synchronous, stream=0) ----------

int eml_launch_exp(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    exp_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_log(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    log_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_sin(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    sin_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_cos(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    cos_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_tan(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    tan_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_sinh(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    sinh_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_cosh(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    cosh_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_tanh(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    tanh_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_sqrt(const void* x, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    sqrt_kernel<<<grid_size, block_size>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_eml(const void* x, const void* y, void* result, int n, int block_size) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    eml_kernel<<<grid_size, block_size>>>((const double*)x, (const double*)y, (double*)result, n);
    return (int)cudaGetLastError();
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

// Streamed kernel launchers
static int launch_exp_streamed(const void* x, void* result, int n, int block_size, cudaStream_t stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    exp_kernel<<<grid_size, block_size, 0, stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_exp_stream(const void* x, void* result, int n, int block_size, long long stream) {
    return launch_exp_streamed(x, result, n, block_size, (cudaStream_t)stream);
}

int eml_launch_log_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    log_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_sin_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    sin_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_cos_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    cos_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_tan_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    tan_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_sinh_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    sinh_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_cosh_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    cosh_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_tanh_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    tanh_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_sqrt_stream(const void* x, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    sqrt_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (double*)result, n);
    return (int)cudaGetLastError();
}

int eml_launch_eml_stream(const void* x, const void* y, void* result, int n, int block_size, long long stream) {
    if (block_size <= 0 || block_size > 1024) block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    eml_kernel<<<grid_size, block_size, 0, (cudaStream_t)stream>>>((const double*)x, (const double*)y, (double*)result, n);
    return (int)cudaGetLastError();
}
