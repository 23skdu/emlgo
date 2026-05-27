#import <Foundation/Foundation.h>
#import <Metal/Metal.h>

typedef struct {
    int id;
    char name[256];
    int64_t memoryBytes;
    int maxThreadsPerThreadgroup;
    int clockRateMHz;
} MetalDeviceProps;

int eml_metal_init(void);
void eml_metal_cleanup(void);
int eml_metal_get_device_count(void);
int eml_metal_get_device_props(int id, MetalDeviceProps *props);
int eml_metal_launch_unary(const char *name, const double *x,
                           double *result, int n);
int eml_metal_launch_binary(const char *name, const double *a,
                            const double *b, double *result, int n);
int eml_metal_launch_ternary(const char *name, const double *a,
                             const double *b, const double *c,
                             double *result, int n);
int eml_metal_launch_scalar(const char *name, const double *a,
                            double scalar, double *result, int n);
int eml_metal_launch_eml(const double *x, const double *y,
                         double *result, int n);
void eml_metal_sync(void);

static id<MTLDevice> _device = nil;
static id<MTLCommandQueue> _queue = nil;
static id<MTLLibrary> _library = nil;
static bool _initialized = false;
static int _deviceCount = 0;
static MetalDeviceProps _deviceProps[16];

static int initMetal(void) {
    if (_initialized) return 0;

    @autoreleasepool {
        NSArray<id<MTLDevice>> *devices = MTLCopyAllDevices();
        _deviceCount = (int)[devices count];
        if (_deviceCount == 0) { [devices release]; return 1; }

        _device = [devices firstObject];
        [_device retain];
        [devices release];

        NSString *libPaths[] = {
            [[NSBundle mainBundle] pathForResource:@"eml_metal_shaders"
                                            ofType:@"metallib"],
            [[[NSBundle mainBundle] resourcePath]
                stringByAppendingPathComponent:@"eml_metal_shaders.metallib"],
            @"./eml_metal_shaders.metallib",
            @"./metal/eml_metal_shaders.metallib",
            [[[NSProcessInfo processInfo] environment][@"PWD"]
                stringByAppendingPathComponent:@"metal/eml_metal_shaders.metallib"],
        };

        NSError *err = nil;
        _library = nil;
        for (int i = 0; i < 5 && _library == nil; i++) {
            if (libPaths[i]) {
                NSURL *url = [NSURL fileURLWithPath:libPaths[i]];
                _library = [_device newLibraryWithURL:url error:&err];
            }
        }
        if (!_library) {
            fprintf(stderr, "Failed to load Metal shader library\n");
            return 2;
        }
        [_library retain];

        _queue = [_device newCommandQueue];
        [_queue retain];

        for (int i = 0; i < _deviceCount && i < 16; i++) {
            id<MTLDevice> dev = (i == 0) ? _device : MTLCopyAllDevices()[i];
            _deviceProps[i].id = i;
            strncpy(_deviceProps[i].name, [[dev name] UTF8String], 255);
            _deviceProps[i].name[255] = '\0';
            _deviceProps[i].memoryBytes = (int64_t)[dev recommendedMaxWorkingSetSize];
            _deviceProps[i].maxThreadsPerThreadgroup =
                (int)[dev maxThreadsPerThreadgroup].width;
            _deviceProps[i].clockRateMHz = 0;
            if (i != 0) [dev release];
        }
        _initialized = true;
    }
    return 0;
}

int eml_metal_init(void) { return initMetal(); }

void eml_metal_cleanup(void) {
    if (_queue) { [_queue release]; _queue = nil; }
    if (_library) { [_library release]; _library = nil; }
    if (_device) { [_device release]; _device = nil; }
    _initialized = false; _deviceCount = 0;
}

int eml_metal_get_device_count(void) {
    if (!_initialized) initMetal();
    return _deviceCount;
}

int eml_metal_get_device_props(int id, MetalDeviceProps *props) {
    if (!_initialized) initMetal();
    if (id < 0 || id >= _deviceCount) return 1;
    *props = _deviceProps[id];
    return 0;
}

static int launchKernel(const char *kernelName,
                        void (^setup)(id<MTLComputeCommandEncoder>,
                                      id<MTLComputePipelineState>,
                                      id<MTLCommandBuffer>)) {
    if (!_initialized && initMetal() != 0) return 1;

    @autoreleasepool {
        NSString *name = [NSString stringWithUTF8String:kernelName];
        id<MTLFunction> func = [_library newFunctionWithName:name];
        if (!func) return 2;

        NSError *err = nil;
        id<MTLComputePipelineState> pipeline =
            [_device newComputePipelineStateWithFunction:func error:&err];
        [func release];
        if (!pipeline) return 3;

        id<MTLCommandBuffer> cmdBuf = [_queue commandBuffer];
        id<MTLComputeCommandEncoder> enc = [cmdBuf computeCommandEncoder];
        [enc setComputePipelineState:pipeline];
        setup(enc, pipeline, cmdBuf);
        [enc endEncoding];
        [cmdBuf commit];
        [cmdBuf waitUntilCompleted];

        int status = (cmdBuf.status == MTLCommandBufferStatusCompleted) ? 0 : 4;
        [pipeline release];
        return status;
    }
}

int eml_metal_launch_unary(const char *name,
                            const double *x, double *result, int n) {
    if (n <= 0) return 0;
    float *xf = (float *)malloc(n * sizeof(float));
    float *rf = (float *)malloc(n * sizeof(float));
    if (!xf || !rf) { free(xf); free(rf); return 5; }
    for (int i = 0; i < n; i++) xf[i] = (float)x[i];

    int err = launchKernel(name, ^(id<MTLComputeCommandEncoder> enc,
                                    id<MTLComputePipelineState> pipeline,
                                    id<MTLCommandBuffer> cmdBuf) {
        id<MTLBuffer> bufX = [_device newBufferWithBytes:xf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        id<MTLBuffer> bufR = [_device newBufferWithBytes:rf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        [enc setBuffer:bufX offset:0 atIndex:0];
        [enc setBuffer:bufR offset:0 atIndex:1];
        [enc dispatchThreads:MTLSizeMake(n, 1, 1)
            threadsPerThreadgroup:MTLSizeMake(256, 1, 1)];
        memcpy(rf, [bufR contents], n * sizeof(float));
        [bufX release]; [bufR release];
    });

    for (int i = 0; i < n; i++) result[i] = (double)rf[i];
    free(xf); free(rf);
    return err;
}

int eml_metal_launch_binary(const char *name,
                             const double *a, const double *b,
                             double *result, int n) {
    if (n <= 0) return 0;
    float *af = (float *)malloc(n * sizeof(float));
    float *bf = (float *)malloc(n * sizeof(float));
    float *rf = (float *)malloc(n * sizeof(float));
    if (!af || !bf || !rf) { free(af); free(bf); free(rf); return 5; }
    for (int i = 0; i < n; i++) { af[i] = (float)a[i]; bf[i] = (float)b[i]; }

    int err = launchKernel(name, ^(id<MTLComputeCommandEncoder> enc,
                                    id<MTLComputePipelineState> pipeline,
                                    id<MTLCommandBuffer> cmdBuf) {
        id<MTLBuffer> bufA = [_device newBufferWithBytes:af
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        id<MTLBuffer> bufB = [_device newBufferWithBytes:bf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        id<MTLBuffer> bufR = [_device newBufferWithBytes:rf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        [enc setBuffer:bufA offset:0 atIndex:0];
        [enc setBuffer:bufB offset:0 atIndex:1];
        [enc setBuffer:bufR offset:0 atIndex:2];
        [enc dispatchThreads:MTLSizeMake(n, 1, 1)
            threadsPerThreadgroup:MTLSizeMake(256, 1, 1)];
        memcpy(rf, [bufR contents], n * sizeof(float));
        [bufA release]; [bufB release]; [bufR release];
    });

    for (int i = 0; i < n; i++) result[i] = (double)rf[i];
    free(af); free(bf); free(rf);
    return err;
}

int eml_metal_launch_ternary(const char *name,
                              const double *a, const double *b, const double *c,
                              double *result, int n) {
    if (n <= 0) return 0;
    float *af = (float *)malloc(n * sizeof(float));
    float *bf = (float *)malloc(n * sizeof(float));
    float *cf = (float *)malloc(n * sizeof(float));
    float *rf = (float *)malloc(n * sizeof(float));
    if (!af || !bf || !cf || !rf) {
        free(af); free(bf); free(cf); free(rf); return 5;
    }
    for (int i = 0; i < n; i++) {
        af[i] = (float)a[i]; bf[i] = (float)b[i]; cf[i] = (float)c[i];
    }

    int err = launchKernel(name, ^(id<MTLComputeCommandEncoder> enc,
                                    id<MTLComputePipelineState> pipeline,
                                    id<MTLCommandBuffer> cmdBuf) {
        id<MTLBuffer> bufA = [_device newBufferWithBytes:af
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        id<MTLBuffer> bufB = [_device newBufferWithBytes:bf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        id<MTLBuffer> bufC = [_device newBufferWithBytes:cf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        id<MTLBuffer> bufR = [_device newBufferWithBytes:rf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        [enc setBuffer:bufA offset:0 atIndex:0];
        [enc setBuffer:bufB offset:0 atIndex:1];
        [enc setBuffer:bufC offset:0 atIndex:2];
        [enc setBuffer:bufR offset:0 atIndex:3];
        [enc dispatchThreads:MTLSizeMake(n, 1, 1)
            threadsPerThreadgroup:MTLSizeMake(256, 1, 1)];
        memcpy(rf, [bufR contents], n * sizeof(float));
        [bufA release]; [bufB release]; [bufC release]; [bufR release];
    });

    for (int i = 0; i < n; i++) result[i] = (double)rf[i];
    free(af); free(bf); free(cf); free(rf);
    return err;
}

int eml_metal_launch_scalar(const char *name,
                             const double *a, double scalar,
                             double *result, int n) {
    if (n <= 0) return 0;
    float *af = (float *)malloc(n * sizeof(float));
    float *rf = (float *)malloc(n * sizeof(float));
    float sf = (float)scalar;
    if (!af || !rf) { free(af); free(rf); return 5; }
    for (int i = 0; i < n; i++) af[i] = (float)a[i];

    int err = launchKernel(name, ^(id<MTLComputeCommandEncoder> enc,
                                    id<MTLComputePipelineState> pipeline,
                                    id<MTLCommandBuffer> cmdBuf) {
        id<MTLBuffer> bufA = [_device newBufferWithBytes:af
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        id<MTLBuffer> bufR = [_device newBufferWithBytes:rf
            length:n * sizeof(float) options:MTLResourceStorageModeShared];
        [enc setBuffer:bufA offset:0 atIndex:0];
        [enc setBytes:&sf length:sizeof(float) atIndex:1];
        [enc setBuffer:bufR offset:0 atIndex:2];
        [enc dispatchThreads:MTLSizeMake(n, 1, 1)
            threadsPerThreadgroup:MTLSizeMake(256, 1, 1)];
        memcpy(rf, [bufR contents], n * sizeof(float));
        [bufA release]; [bufR release];
    });

    for (int i = 0; i < n; i++) result[i] = (double)rf[i];
    free(af); free(rf);
    return err;
}

int eml_metal_launch_eml(const double *x, const double *y,
                          double *result, int n) {
    return eml_metal_launch_binary("kernel_eml", x, y, result, n);
}

void eml_metal_sync(void) {}
