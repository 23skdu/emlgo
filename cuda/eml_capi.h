#ifndef EML_CAPI_H
#define EML_CAPI_H

// Pure C API for cgo consumption.
// No CUDA types appear in function signatures.
// All functions return 0 on success, nonzero on error.

// Initialize CUDA subsystem. Must be called before any other function.
int eml_init(void);

// Shut down CUDA subsystem.
int eml_cleanup(void);

// Get number of available CUDA devices.
int eml_get_device_count(int* count);

// Get device name (max 256 chars including null terminator).
int eml_get_device_name(int device_id, char* name, int max_len);

// Get device properties.
int eml_get_device_props(
    int device_id,
    int* compute_major,
    int* compute_minor,
    long long* memory_bytes,
    int* max_threads_per_block,
    int* warp_size,
    int* clock_rate_khz
);

// Allocate device memory.
int eml_allocate(void** ptr, long long size);

// Free device memory.
int eml_free(void* ptr);

// Copy from host to device.
int eml_copy_to_device(void* dst, const void* src, long long size);

// Copy from device to host.
int eml_copy_to_host(void* dst, const void* src, long long size);

// Synchronize device.
int eml_sync_device(void);

// ---------- Kernel Launches ----------

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

// ---------- Pinned Memory ----------

int eml_allocate_pinned(void** ptr, long long size);
int eml_free_pinned(void* ptr);

// ---------- Async Streams ----------

// Create a stream, returns handle via stream_out (0 indicates default/null stream).
long long eml_create_stream(void);
int eml_destroy_stream(long long stream);
int eml_sync_stream(long long stream);

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

#endif // EML_CAPI_H
