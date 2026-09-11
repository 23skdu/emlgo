# Makefile for building CUDA C API shared library
# Requirements: nvcc (NVIDIA CUDA Compiler), CUDA Toolkit

NVCC := nvcc
TARGET := libeml_capi.so
SRCS := eml_cuda.cu eml_capi.cu
OBJS := eml_cuda.o eml_capi.o

# CUDA architecture: sm_70 = Volta, sm_75 = Turing, sm_80 = Ampere, sm_89 = Ada Lovelace, sm_90 = Hopper
# Defaults to sm_89 for RTX 40-series; override with ARCH=sm_XX as needed.
ARCH ?= sm_89

CFLAGS := -O3 -arch=$(ARCH) -Xcompiler -fPIC
LDFLAGS := -shared -lcuda -lcudart

.PHONY: all clean install example

all: $(TARGET)

eml_cuda.o: eml_cuda.cu eml_cuda.h
	$(NVCC) $(CFLAGS) -c eml_cuda.cu -o eml_cuda.o

eml_capi.o: eml_capi.cu eml_capi.h eml_cuda.h
	$(NVCC) $(CFLAGS) -c eml_capi.cu -o eml_capi.o

$(TARGET): $(OBJS)
	$(NVCC) $(OBJS) $(LDFLAGS) -o $(TARGET)

clean:
	rm -f $(OBJS) $(TARGET) eml_cuda_example

install: $(TARGET)
	@echo "To use: LD_LIBRARY_PATH=$(CURDIR) go build -tags cuda ./..."

# Also build the standalone example
example: eml_cuda.cu eml_cuda.h
	$(NVCC) -O3 -arch=$(ARCH) eml_cuda.cu -o eml_cuda_example
