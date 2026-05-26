# WebAssembly (WASM) Support

This document describes emlgo's WebAssembly support, including SIMD-optimized kernels, build instructions, and testing infrastructure.

## Overview

emlgo compiles to WebAssembly via `GOOS=js GOARCH=wasm` (browser/Node.js) or `GOOS=wasip1 GOARCH=wasm` (WASI). All core mathematical functions are available in WASM builds, with SIMD-accelerated batch operations using loop-unrolled kernels designed for WASM JIT auto-vectorization.

## Building for WASM

Build the full library for WASM:

```bash
GOOS=js GOARCH=wasm go build ./...
```

Build the CLI tool as a WASM binary:

```bash
GOOS=js GOARCH=wasm go build -o wasm/emlcli.wasm ./cmd/emlcli
```

## WASM SIMD Kernels

The library provides optimized batch kernels that take advantage of WASM SIMD (wasm_simd128) through JIT auto-vectorization. Rather than architecture-specific assembly, these kernels use 8-wide loop unrolling, enabling V8 (Chrome/Node.js) and SpiderMonkey (Firefox) to emit SIMD instructions automatically.

### Available Operations

| Operation | Function | Description |
|-----------|----------|-------------|
| Add | `addWasmSIMD` | Element-wise addition |
| Sub | `subWasmSIMD` | Element-wise subtraction |
| Mul | `mulWasmSIMD` | Element-wise multiplication |
| Div | `divWasmSIMD` | Element-wise division |
| Sqrt | `sqrtWasmSIMD` | Element-wise square root |
| Abs | `absWasmSIMD` | Element-wise absolute value |
| Neg | `negWasmSIMD` | Element-wise negation |
| Inv | `invWasmSIMD` | Element-wise inverse (1/x) |
| FMA | `fmaWasmSIMD` | Fused multiply-add (a*b + c) |
| AddScalar | `addScalarWasmSIMD` | Add constant to each element |
| MulScalar | `mulScalarWasmSIMD` | Multiply each element by constant |

Transcendental functions (Exp, Log, Sin, Cos, Tan) use parallelized scalar computation via goroutines, which are cooperatively scheduled on WASM's single thread.

### Feature Detection

```go
import "github.com/emlgo/eml/internal/eml"

fmt.Printf("Has WASM SIMD: %v\n", eml.HasWasmSIMD())
```

On WASM builds, `HasWasmSIMD()` returns `true` and all non-WASM feature flags (AVX2, AVX512, NEON, SVE) return `false`.

## Memory Alignment

WASM linear memory access benefits from alignment. The library provides two utilities:

### WasmAlign16

Aligns a float64 slice to a 16-byte boundary, enabling JIT compilers to emit aligned `wasm_simd128` load/store instructions:

```go
import "github.com/emlgo/eml/internal/eml"

data := make([]float64, 1024)
aligned := eml.WasmAlign16(data) // drops first element if misaligned
```

### WasmPageAlign

Rounds a byte count up to the nearest WASM page (64 KiB):

```go
size := eml.WasmPageAlign(100000) // returns 131072 (2 pages)
```

## Test Harness

The `scripts/wasm_test.sh` script builds WASM binaries and runs them with Node.js.

### Prerequisites

- Go 1.21+ (for WASM target support)
- Node.js 16+ (for WASM SIMD support)

### Usage

```bash
# Build and run tests
./scripts/wasm_test.sh test

# Run performance benchmarks
./scripts/wasm_test.sh bench

# Start local web server for benchmark page
./scripts/wasm_test.sh serve

# Clean build artifacts
./scripts/wasm_test.sh clean
```

The `test` mode builds a WASM test binary that exercises all core operations and prints SIMD detection status. The `bench` mode runs batch performance benchmarks across multiple sizes and reports ns/element metrics.

## Web Benchmark

Open `wasm/bench.html` in a browser (via a local web server) for interactive performance benchmarking. The page:

- Loads the compiled WASM binary
- Runs all batch operations across multiple sizes
- Reports timing (ms) and per-element latency (ns/elem)
- Shows pass/fail status for each test

Start the local server:

```bash
./scripts/wasm_test.sh serve
# Open http://localhost:8080/wasm/bench.html
```

## Build Architecture

The WASM build uses `//go:build wasm` build tags to select the appropriate dispatch layer:

- `simd_dispatch_wasm.go` — routes all SIMD operations to WASM-optimized kernels
- `simd_wasm.go` — implements block-unrolled kernels
- `wasm_utils.go` — memory alignment utilities
- `simd_dispatch_stub.go` — excluded from WASM builds (`!wasm` constraint)

## Performance Considerations

- **Batch size**: Optimal performance is achieved with batch sizes >1024 elements. Small batches incur overhead from memory allocation and function call dispatch.
- **Goroutines**: On WASM, goroutines are cooperatively scheduled on a single thread. Parallelized operations (ExpBatch, LogBatch, etc.) provide concurrency but not true parallelism.
- **SIMD auto-vectorization**: The block-unrolled loops are recognized by V8/TurboFan and SpiderMonkey and compiled to `wasm_simd128` instructions. Check the browser's DevTools Performance panel to verify SIMD usage.
- **Memory alignment**: Use `WasmAlign16()` for critical paths to avoid alignment-related slowdowns in JIT-compiled code.
