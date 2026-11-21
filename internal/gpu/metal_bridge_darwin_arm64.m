//go:build darwin && arm64 && metal && cgo && !cuda

#import <Foundation/Foundation.h>
#import <Metal/Metal.h>
#include <pthread.h>
#include <stdint.h>
#include <string.h>

typedef struct {
    int id;
    char name[256];
    int64_t memoryBytes;
    int maxThreadsPerThreadgroup;
    int warpSize;
    int clockRateMHz;
} MetalDeviceProps;

static pthread_mutex_t g_metal_lock = PTHREAD_MUTEX_INITIALIZER;
static id<MTLDevice> g_device = nil;
static id<MTLCommandQueue> g_queue = nil;
static id<MTLLibrary> g_library = nil;
static NSMutableDictionary<NSString *, id<MTLComputePipelineState>> *g_pipelines = nil;

static id<MTLComputePipelineState> get_pipeline_locked(const char *kernel_name) {
    if (!g_device || !g_library || !g_pipelines || !kernel_name) {
        return nil;
    }
    NSString *key = [NSString stringWithUTF8String:kernel_name];
    id<MTLComputePipelineState> pso = g_pipelines[key];
    if (pso) {
        return pso;
    }
    id<MTLFunction> fn = [g_library newFunctionWithName:key];
    if (!fn) {
        return nil;
    }
    NSError *err = nil;
    pso = [g_device newComputePipelineStateWithFunction:fn error:&err];
    if (!pso || err) {
        return nil;
    }
    g_pipelines[key] = pso;
    return pso;
}

int eml_metal_init(const char *msl_source) {
    if (!msl_source) {
        return 6; // ErrInvalidConfig
    }
    pthread_mutex_lock(&g_metal_lock);
    @autoreleasepool {
        if (g_device && g_queue && g_library) {
            pthread_mutex_unlock(&g_metal_lock);
            return 0;
        }
        id<MTLDevice> dev = MTLCreateSystemDefaultDevice();
        if (!dev) {
            pthread_mutex_unlock(&g_metal_lock);
            return 1; // ErrNoDevice
        }
        id<MTLCommandQueue> queue = [dev newCommandQueue];
        if (!queue) {
            pthread_mutex_unlock(&g_metal_lock);
            return 2; // ErrNotInit
        }
        MTLCompileOptions *opts = [[MTLCompileOptions alloc] init];
#if defined(__MAC_15_0) && __MAC_OS_X_VERSION_MAX_ALLOWED >= __MAC_15_0
        if (@available(macOS 15.0, *)) {
            opts.mathMode = MTLMathModeSafe;
        } else {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
            opts.fastMathEnabled = NO;
#pragma clang diagnostic pop
        }
#else
        opts.fastMathEnabled = NO;
#endif
        NSString *src = [NSString stringWithUTF8String:msl_source];
        NSError *err = nil;
        id<MTLLibrary> lib = [dev newLibraryWithSource:src options:opts error:&err];
        if (!lib) {
            pthread_mutex_unlock(&g_metal_lock);
            return 2; // ErrNotInit
        }
        g_device = dev;
        g_queue = queue;
        g_library = lib;
        g_pipelines = [[NSMutableDictionary alloc] init];
    }
    pthread_mutex_unlock(&g_metal_lock);
    return 0;
}

void eml_metal_cleanup(void) {
    pthread_mutex_lock(&g_metal_lock);
    @autoreleasepool {
        g_pipelines = nil;
        g_library = nil;
        g_queue = nil;
        g_device = nil;
    }
    pthread_mutex_unlock(&g_metal_lock);
}

int eml_metal_get_device_count(void) {
    @autoreleasepool {
        NSArray<id<MTLDevice>> *devices = MTLCopyAllDevices();
        if (devices && devices.count > 0) {
            return (int)devices.count;
        }
        id<MTLDevice> dev = MTLCreateSystemDefaultDevice();
        return dev ? 1 : 0;
    }
}

int eml_metal_get_device_props(int device_id, MetalDeviceProps *props) {
    if (!props || device_id < 0) {
        return 6;
    }
    @autoreleasepool {
        NSArray<id<MTLDevice>> *devices = MTLCopyAllDevices();
        id<MTLDevice> dev = nil;
        if (devices && (NSUInteger)device_id < devices.count) {
            dev = devices[(NSUInteger)device_id];
        } else if (device_id == 0) {
            dev = MTLCreateSystemDefaultDevice();
        }
        if (!dev) {
            return 1;
        }
        memset(props, 0, sizeof(MetalDeviceProps));
        props->id = device_id;
        const char *utf8 = [dev.name UTF8String];
        if (utf8) {
            strncpy(props->name, utf8, sizeof(props->name) - 1);
        }
        props->memoryBytes = (int64_t)dev.recommendedMaxWorkingSetSize;
        props->maxThreadsPerThreadgroup = (int)dev.maxThreadsPerThreadgroup.width;
        props->warpSize = 32; // Apple Silicon SIMD execution width
        props->clockRateMHz = 0;
        return 0;
    }
}

int eml_metal_launch_unary(const char *kernel_name, const void *x, void *result, int n, int elem_size) {
    if (n <= 0 || elem_size <= 0 || !x || !result || !kernel_name) {
        return 6;
    }
    pthread_mutex_lock(&g_metal_lock);
    int rc = 0;
    @autoreleasepool {
        id<MTLComputePipelineState> pso = get_pipeline_locked(kernel_name);
        if (!pso || !g_device || !g_queue) {
            rc = 2;
        } else {
            NSUInteger bytes = (NSUInteger)n * (NSUInteger)elem_size;
            id<MTLBuffer> xBuf = [g_device newBufferWithBytes:x length:bytes options:MTLResourceStorageModeShared];
            id<MTLBuffer> resBuf = [g_device newBufferWithLength:bytes options:MTLResourceStorageModeShared];
            if (!xBuf || !resBuf) {
                rc = 3;
            } else {
                id<MTLCommandBuffer> cmdBuf = [g_queue commandBuffer];
                id<MTLComputeCommandEncoder> enc = [cmdBuf computeCommandEncoder];
                uint32_t count = (uint32_t)n;
                [enc setComputePipelineState:pso];
                [enc setBuffer:xBuf offset:0 atIndex:0];
                [enc setBuffer:resBuf offset:0 atIndex:1];
                [enc setBytes:&count length:sizeof(uint32_t) atIndex:2];

                NSUInteger tgWidth = pso.maxTotalThreadsPerThreadgroup;
                if (tgWidth > 256) tgWidth = 256;
                if (tgWidth > (NSUInteger)n) tgWidth = (NSUInteger)n;
                if (tgWidth == 0) tgWidth = 1;

                MTLSize gridSize = MTLSizeMake((NSUInteger)n, 1, 1);
                MTLSize tgSize = MTLSizeMake(tgWidth, 1, 1);
                [enc dispatchThreads:gridSize threadsPerThreadgroup:tgSize];
                [enc endEncoding];

                [cmdBuf commit];
                [cmdBuf waitUntilCompleted];
                if (cmdBuf.status != MTLCommandBufferStatusCompleted) {
                    rc = 5;
                } else {
                    memcpy(result, resBuf.contents, bytes);
                }
            }
        }
    }
    pthread_mutex_unlock(&g_metal_lock);
    return rc;
}

int eml_metal_launch_binary(const char *kernel_name, const void *x, const void *y, void *result, int n, int elem_size) {
    if (n <= 0 || elem_size <= 0 || !x || !y || !result || !kernel_name) {
        return 6;
    }
    pthread_mutex_lock(&g_metal_lock);
    int rc = 0;
    @autoreleasepool {
        id<MTLComputePipelineState> pso = get_pipeline_locked(kernel_name);
        if (!pso || !g_device || !g_queue) {
            rc = 2;
        } else {
            NSUInteger bytes = (NSUInteger)n * (NSUInteger)elem_size;
            id<MTLBuffer> xBuf = [g_device newBufferWithBytes:x length:bytes options:MTLResourceStorageModeShared];
            id<MTLBuffer> yBuf = [g_device newBufferWithBytes:y length:bytes options:MTLResourceStorageModeShared];
            id<MTLBuffer> resBuf = [g_device newBufferWithLength:bytes options:MTLResourceStorageModeShared];
            if (!xBuf || !yBuf || !resBuf) {
                rc = 3;
            } else {
                id<MTLCommandBuffer> cmdBuf = [g_queue commandBuffer];
                id<MTLComputeCommandEncoder> enc = [cmdBuf computeCommandEncoder];
                uint32_t count = (uint32_t)n;
                [enc setComputePipelineState:pso];
                [enc setBuffer:xBuf offset:0 atIndex:0];
                [enc setBuffer:yBuf offset:0 atIndex:1];
                [enc setBuffer:resBuf offset:0 atIndex:2];
                [enc setBytes:&count length:sizeof(uint32_t) atIndex:3];

                NSUInteger tgWidth = pso.maxTotalThreadsPerThreadgroup;
                if (tgWidth > 256) tgWidth = 256;
                if (tgWidth > (NSUInteger)n) tgWidth = (NSUInteger)n;
                if (tgWidth == 0) tgWidth = 1;

                MTLSize gridSize = MTLSizeMake((NSUInteger)n, 1, 1);
                MTLSize tgSize = MTLSizeMake(tgWidth, 1, 1);
                [enc dispatchThreads:gridSize threadsPerThreadgroup:tgSize];
                [enc endEncoding];

                [cmdBuf commit];
                [cmdBuf waitUntilCompleted];
                if (cmdBuf.status != MTLCommandBufferStatusCompleted) {
                    rc = 5;
                } else {
                    memcpy(result, resBuf.contents, bytes);
                }
            }
        }
    }
    pthread_mutex_unlock(&g_metal_lock);
    return rc;
}

int eml_metal_launch_dot_f32(const float *a, const float *b, float *out, int n) {
    if (n <= 0 || !a || !b || !out) {
        return 6;
    }
    pthread_mutex_lock(&g_metal_lock);
    int rc = 0;
    @autoreleasepool {
        id<MTLComputePipelineState> pso = get_pipeline_locked("kernel_dot_f32");
        if (!pso || !g_device || !g_queue) {
            rc = 2;
        } else {
            const NSUInteger tgWidth = 256;
            NSUInteger numGroups = ((NSUInteger)n + tgWidth - 1) / tgWidth;
            NSUInteger bytes = (NSUInteger)n * sizeof(float);
            NSUInteger partialBytes = numGroups * sizeof(float);

            id<MTLBuffer> aBuf = [g_device newBufferWithBytes:a length:bytes options:MTLResourceStorageModeShared];
            id<MTLBuffer> bBuf = [g_device newBufferWithBytes:b length:bytes options:MTLResourceStorageModeShared];
            id<MTLBuffer> pBuf = [g_device newBufferWithLength:partialBytes options:MTLResourceStorageModeShared];
            if (!aBuf || !bBuf || !pBuf) {
                rc = 3;
            } else {
                id<MTLCommandBuffer> cmdBuf = [g_queue commandBuffer];
                id<MTLComputeCommandEncoder> enc = [cmdBuf computeCommandEncoder];
                uint32_t count = (uint32_t)n;
                [enc setComputePipelineState:pso];
                [enc setBuffer:aBuf offset:0 atIndex:0];
                [enc setBuffer:bBuf offset:0 atIndex:1];
                [enc setBuffer:pBuf offset:0 atIndex:2];
                [enc setBytes:&count length:sizeof(uint32_t) atIndex:3];

                [enc dispatchThreadgroups:MTLSizeMake(numGroups, 1, 1)
                    threadsPerThreadgroup:MTLSizeMake(tgWidth, 1, 1)];
                [enc endEncoding];

                [cmdBuf commit];
                [cmdBuf waitUntilCompleted];
                if (cmdBuf.status != MTLCommandBufferStatusCompleted) {
                    rc = 5;
                } else {
                    const float *partials = (const float *)pBuf.contents;
                    double sum = 0.0;
                    for (NSUInteger i = 0; i < numGroups; i++) {
                        sum += (double)partials[i];
                    }
                    *out = (float)sum;
                }
            }
        }
    }
    pthread_mutex_unlock(&g_metal_lock);
    return rc;
}

int eml_metal_launch_dot_c64(const float *a, const float *b, float *out_re, float *out_im, int n) {
    if (n <= 0 || !a || !b || !out_re || !out_im) {
        return 6;
    }
    pthread_mutex_lock(&g_metal_lock);
    int rc = 0;
    @autoreleasepool {
        id<MTLComputePipelineState> pso = get_pipeline_locked("kernel_dot_c64");
        if (!pso || !g_device || !g_queue) {
            rc = 2;
        } else {
            const NSUInteger tgWidth = 256;
            NSUInteger numGroups = ((NSUInteger)n + tgWidth - 1) / tgWidth;
            NSUInteger bytes = (NSUInteger)n * 2 * sizeof(float);
            NSUInteger partialBytes = numGroups * 2 * sizeof(float);

            id<MTLBuffer> aBuf = [g_device newBufferWithBytes:a length:bytes options:MTLResourceStorageModeShared];
            id<MTLBuffer> bBuf = [g_device newBufferWithBytes:b length:bytes options:MTLResourceStorageModeShared];
            id<MTLBuffer> pBuf = [g_device newBufferWithLength:partialBytes options:MTLResourceStorageModeShared];
            if (!aBuf || !bBuf || !pBuf) {
                rc = 3;
            } else {
                id<MTLCommandBuffer> cmdBuf = [g_queue commandBuffer];
                id<MTLComputeCommandEncoder> enc = [cmdBuf computeCommandEncoder];
                uint32_t count = (uint32_t)n;
                [enc setComputePipelineState:pso];
                [enc setBuffer:aBuf offset:0 atIndex:0];
                [enc setBuffer:bBuf offset:0 atIndex:1];
                [enc setBuffer:pBuf offset:0 atIndex:2];
                [enc setBytes:&count length:sizeof(uint32_t) atIndex:3];

                [enc dispatchThreadgroups:MTLSizeMake(numGroups, 1, 1)
                    threadsPerThreadgroup:MTLSizeMake(tgWidth, 1, 1)];
                [enc endEncoding];

                [cmdBuf commit];
                [cmdBuf waitUntilCompleted];
                if (cmdBuf.status != MTLCommandBufferStatusCompleted) {
                    rc = 5;
                } else {
                    const float *partials = (const float *)pBuf.contents;
                    double sumRe = 0.0;
                    double sumIm = 0.0;
                    for (NSUInteger i = 0; i < numGroups; i++) {
                        sumRe += (double)partials[2 * i];
                        sumIm += (double)partials[2 * i + 1];
                    }
                    *out_re = (float)sumRe;
                    *out_im = (float)sumIm;
                }
            }
        }
    }
    pthread_mutex_unlock(&g_metal_lock);
    return rc;
}
