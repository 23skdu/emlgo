#ifndef EML_CAPI_H
#define EML_CAPI_H

#ifdef __cplusplus
extern "C" {
#endif

// Pure C API for cgo consumption.
// No CUDA types appear in function signatures.
// All functions return 0 on success, nonzero on error.

// ---------- Lifecycle ----------
int eml_init(void);
int eml_cleanup(void);

// ---------- Device Query ----------
int eml_get_device_count(int* count);
int eml_get_device_name(int device_id, char* name, int max_len);
int eml_get_device_props(
    int device_id,
    int* compute_major,
    int* compute_minor,
    long long* memory_bytes,
    int* max_threads_per_block,
    int* warp_size,
    int* clock_rate_khz
);

// ---------- Memory Management ----------
int eml_allocate(void** ptr, long long size);
int eml_free(void* ptr);
int eml_copy_to_device(void* dst, const void* src, long long size);
int eml_copy_to_host(void* dst, const void* src, long long size);
int eml_sync_device(void);
int eml_memset(void* ptr, int value, long long size);
int eml_memset_stream(void* ptr, int value, long long size, long long stream);

// ---------- Pinned Memory ----------
int eml_allocate_pinned(void** ptr, long long size);
int eml_free_pinned(void* ptr);

// ---------- Async Streams ----------
long long eml_create_stream(void);
int eml_destroy_stream(long long stream);
int eml_sync_stream(long long stream);

// ---------- Kernel Launches: Float64 (Synchronous & Streamed) ----------
int eml_launch_exp(const void* x, void* result, int n, int block_size);
int eml_launch_log(const void* x, void* result, int n, int block_size);
int eml_launch_sin(const void* x, void* result, int n, int block_size);
int eml_launch_cos(const void* x, void* result, int n, int block_size);
int eml_launch_tan(const void* x, void* result, int n, int block_size);
int eml_launch_sinh(const void* x, void* result, int n, int block_size);
int eml_launch_cosh(const void* x, void* result, int n, int block_size);
int eml_launch_tanh(const void* x, void* result, int n, int block_size);
int eml_launch_sqrt(const void* x, void* result, int n, int block_size);
int eml_launch_eml(const void* x, const void* y, void* result, int n, int block_size);
int eml_launch_dot(const void* a, const void* b, void* result, int n, int block_size);

int eml_launch_exp_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_log_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sin_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cos_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tan_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sinh_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cosh_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tanh_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sqrt_stream(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_eml_stream(const void* x, const void* y, void* result, int n, int block_size, long long stream);
int eml_launch_dot_stream(const void* a, const void* b, void* result, int n, int block_size, long long stream);

// ---------- Kernel Launches: Float32 (Synchronous & Streamed) ----------
int eml_launch_exp_f32(const void* x, void* result, int n, int block_size);
int eml_launch_log_f32(const void* x, void* result, int n, int block_size);
int eml_launch_sin_f32(const void* x, void* result, int n, int block_size);
int eml_launch_cos_f32(const void* x, void* result, int n, int block_size);
int eml_launch_tan_f32(const void* x, void* result, int n, int block_size);
int eml_launch_sinh_f32(const void* x, void* result, int n, int block_size);
int eml_launch_cosh_f32(const void* x, void* result, int n, int block_size);
int eml_launch_tanh_f32(const void* x, void* result, int n, int block_size);
int eml_launch_sqrt_f32(const void* x, void* result, int n, int block_size);
int eml_launch_eml_f32(const void* x, const void* y, void* result, int n, int block_size);
int eml_launch_dot_f32(const void* a, const void* b, void* result, int n, int block_size);

int eml_launch_exp_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_log_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sin_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cos_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tan_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sinh_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cosh_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tanh_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sqrt_stream_f32(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_eml_stream_f32(const void* x, const void* y, void* result, int n, int block_size, long long stream);
int eml_launch_dot_stream_f32(const void* a, const void* b, void* result, int n, int block_size, long long stream);

// ---------- Kernel Launches: Complex64 (Synchronous & Streamed) ----------
int eml_launch_exp_c64(const void* x, void* result, int n, int block_size);
int eml_launch_log_c64(const void* x, void* result, int n, int block_size);
int eml_launch_sin_c64(const void* x, void* result, int n, int block_size);
int eml_launch_cos_c64(const void* x, void* result, int n, int block_size);
int eml_launch_tan_c64(const void* x, void* result, int n, int block_size);
int eml_launch_sinh_c64(const void* x, void* result, int n, int block_size);
int eml_launch_cosh_c64(const void* x, void* result, int n, int block_size);
int eml_launch_tanh_c64(const void* x, void* result, int n, int block_size);
int eml_launch_sqrt_c64(const void* x, void* result, int n, int block_size);
int eml_launch_eml_c64(const void* x, const void* y, void* result, int n, int block_size);
int eml_launch_dot_c64(const void* a, const void* b, void* result, int n, int block_size);

int eml_launch_exp_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_log_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sin_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cos_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tan_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sinh_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cosh_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tanh_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sqrt_stream_c64(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_eml_stream_c64(const void* x, const void* y, void* result, int n, int block_size, long long stream);
int eml_launch_dot_stream_c64(const void* a, const void* b, void* result, int n, int block_size, long long stream);

// ---------- Kernel Launches: Complex128 (Synchronous & Streamed) ----------
int eml_launch_exp_c128(const void* x, void* result, int n, int block_size);
int eml_launch_log_c128(const void* x, void* result, int n, int block_size);
int eml_launch_sin_c128(const void* x, void* result, int n, int block_size);
int eml_launch_cos_c128(const void* x, void* result, int n, int block_size);
int eml_launch_tan_c128(const void* x, void* result, int n, int block_size);
int eml_launch_sinh_c128(const void* x, void* result, int n, int block_size);
int eml_launch_cosh_c128(const void* x, void* result, int n, int block_size);
int eml_launch_tanh_c128(const void* x, void* result, int n, int block_size);
int eml_launch_sqrt_c128(const void* x, void* result, int n, int block_size);
int eml_launch_eml_c128(const void* x, const void* y, void* result, int n, int block_size);
int eml_launch_dot_c128(const void* a, const void* b, void* result, int n, int block_size);

int eml_launch_exp_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_log_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sin_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cos_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tan_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sinh_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_cosh_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_tanh_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_sqrt_stream_c128(const void* x, void* result, int n, int block_size, long long stream);
int eml_launch_eml_stream_c128(const void* x, const void* y, void* result, int n, int block_size, long long stream);
int eml_launch_dot_stream_c128(const void* a, const void* b, void* result, int n, int block_size, long long stream);

#ifdef __cplusplus
}
#endif

#endif // EML_CAPI_H
