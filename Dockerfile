# ==============================================================================
# Unified Multi-Target Dockerfile for EMLGo
#
# Targets:
#   - cpu  : CPU AVX2-only build (default)
#   - cuda : NVIDIA CUDA GPU build
#
# Examples:
#   docker build --target cpu -t emlgo:cpu .
#   docker build --target cuda -t emlgo:cuda .
# ==============================================================================

ARG GO_VERSION=1.27.2
ARG CUDA_VERSION=12.6.3
ARG CUDA_ARCH=sm_89

# ------------------------------------------------------------------------------
# Stage: Go Toolchain
# ------------------------------------------------------------------------------
FROM golang:${GO_VERSION}-bookworm AS go-toolchain

# ------------------------------------------------------------------------------
# Stage: CPU AVX2 Builder
# ------------------------------------------------------------------------------
FROM golang:${GO_VERSION}-bookworm AS cpu-builder

ENV CGO_ENABLED=0
ENV GOAMD64=v3
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN mkdir -p /out && \
    go build -trimpath -ldflags="-s -w" -o /out/emlcli ./cmd/emlcli && \
    go build -trimpath -ldflags="-s -w" -o /out/bench ./cmd/bench && \
    go build -trimpath -ldflags="-s -w" -o /out/validate ./cmd/validate

RUN /out/validate

# ------------------------------------------------------------------------------
# Target: cpu (Default)
# ------------------------------------------------------------------------------
FROM debian:bookworm-slim AS cpu

RUN groupadd -r emlgo && useradd -r -g emlgo -m -d /home/emlgo emlgo

WORKDIR /home/emlgo

COPY --from=cpu-builder /out/emlcli /usr/local/bin/emlcli
COPY --from=cpu-builder /out/bench /usr/local/bin/bench
COPY --from=cpu-builder /out/validate /usr/local/bin/validate

USER emlgo:emlgo

ENTRYPOINT ["emlcli"]
CMD ["demo"]

# ------------------------------------------------------------------------------
# Stage: CUDA & CGO Builder
# ------------------------------------------------------------------------------
FROM nvidia/cuda:${CUDA_VERSION}-devel-ubuntu24.04 AS cuda-builder

ARG CUDA_ARCH
ENV CUDA_ARCH=${CUDA_ARCH}
ENV DEBIAN_FRONTEND=noninteractive

COPY --from=go-toolchain /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /src
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

WORKDIR /src/cuda
RUN make clean && make ARCH=${CUDA_ARCH}

WORKDIR /src
ENV CGO_ENABLED=1
ENV GOAMD64=v3
ENV GOOS=linux
ENV GOARCH=amd64
ENV LIBRARY_PATH="/usr/local/cuda/lib64:/usr/local/cuda/lib64/stubs:${LIBRARY_PATH}"
ENV LD_LIBRARY_PATH="/usr/local/cuda/lib64:/usr/local/cuda/lib64/stubs:/src/cuda:${LD_LIBRARY_PATH}"
ENV CGO_LDFLAGS="-L/usr/local/cuda/lib64 -L/usr/local/cuda/lib64/stubs"

RUN mkdir -p /out && \
    go build -tags cuda -trimpath -ldflags="-s -w" -o /out/emlcli ./cmd/emlcli && \
    go build -tags cuda -trimpath -ldflags="-s -w" -o /out/bench ./cmd/bench && \
    go build -tags cuda -trimpath -ldflags="-s -w" -o /out/validate ./cmd/validate && \
    cp /src/cuda/libeml_capi.so /out/libeml_capi.so

# ------------------------------------------------------------------------------
# Target: cuda
# ------------------------------------------------------------------------------
FROM nvidia/cuda:${CUDA_VERSION}-runtime-ubuntu24.04 AS cuda

ENV DEBIAN_FRONTEND=noninteractive

RUN groupadd -r emlgo && useradd -r -g emlgo -m -d /home/emlgo emlgo

COPY --from=cuda-builder /out/libeml_capi.so /usr/local/lib/libeml_capi.so
RUN ldconfig

COPY --from=cuda-builder /out/emlcli /usr/local/bin/emlcli
COPY --from=cuda-builder /out/bench /usr/local/bin/bench
COPY --from=cuda-builder /out/validate /usr/local/bin/validate

ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES=compute,utility
ENV LD_LIBRARY_PATH="/usr/local/lib:${LD_LIBRARY_PATH}"

WORKDIR /home/emlgo
USER emlgo:emlgo

ENTRYPOINT ["emlcli"]
CMD ["gpu-status"]
