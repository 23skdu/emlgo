# Makefile for building CUDA C API shared library
# Requirements: nvcc (NVIDIA CUDA Compiler), CUDA Toolkit

NVCC := nvcc
TARGET := libeml_capi.so
SRCS := eml_capi.cu
OBJS := eml_capi.o

# CUDA architecture: sm_70 = Volta, sm_75 = Turing, sm_80 = Ampere, sm_90 = Hopper
# Adjust ARCH based on your GPU
ARCH ?= sm_70

CFLAGS := -O3 -arch=$(ARCH) -Xcompiler -fPIC
LDFLAGS := -shared -lcuda -lcudart

.PHONY: all clean

all: $(TARGET)

$(OBJS): $(SRCS) eml_cuda.h eml_capi.h
	$(NVCC) $(CFLAGS) -c $(SRCS) -o $(OBJS)

$(TARGET): $(OBJS)
	$(NVCC) $(OBJS) $(LDFLAGS) -o $(TARGET)

clean:
	rm -f $(OBJS) $(TARGET)

install: $(TARGET)
	@echo "To use: LD_LIBRARY_PATH=$(CURDIR) go build -tags cuda ./..."

# Also build the standalone example
example: eml_cuda.cu eml_cuda.h
	$(NVCC) -O3 -arch=$(ARCH) eml_cuda.cu -o eml_cuda_example
