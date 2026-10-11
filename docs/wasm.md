# WebAssembly (WASM) Targets

This document describes emlgo's WebAssembly targets, loop-unrolled batch kernels, build instructions, and test infrastructure.

## Overview

emlgo compiles to WebAssembly via `GOOS=js GOARCH=wasm` (browser/Node.js) or `GOOS=wasip1 GOARCH=wasm` (WASI). Core math APIs and batch operations are available; batch kernels use loop-unrolled Go code. The Go toolchain and runtime determine whether any operations are lowered to SIMD instructions.

The JIT compiler is **not** available on WebAssembly. `js/wasm` falls back to the AST interpreter in `internal/jit/jit_wasm.go`; on `wasip1` `Compile` returns an error, and `AllocateExecutableMemory` reports that executable memory is unsupported. Both targets are covered by the `cross-build` CI job.

## Building for WASM

Build the full library for WASM:

```bash
GOOS=js GOARCH=wasm go build ./...
```

Build the CLI tool as a WASM binary:

```bash
GOOS=js GOARCH=wasm go build -o wasm/emlcli.wasm ./cmd/emlcli
```

## WASM Batch Kernels

The `js/wasm` dispatch uses 8-wide loop-unrolled Go loops rather than
hand-written `wasm_simd128` instructions. Unrolling amortizes loop overhead, but
it does not guarantee that Go or a browser engine emits SIMD instructions.
Transcendental operations call scalar math functions for each element. Measure
the generated build on the target runtime before making performance assumptions.

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
| Exp | `expWasmSIMD` | Element-wise exponential |
| Log | `logWasmSIMD` | Element-wise natural logarithm |
| Sin | `sinWasmSIMD` | Element-wise sine |
| Cos | `cosWasmSIMD` | Element-wise cosine |
| Tan | `tanWasmSIMD` | Element-wise tangent |
| Sinh | `sinhWasmSIMD` | Element-wise hyperbolic sine |
| Cosh | `coshWasmSIMD` | Element-wise hyperbolic cosine |
| Tanh | `tanhWasmSIMD` | Element-wise hyperbolic tangent |
| SinCos | `sincosWasmSIMD` | Simultaneous sin and cos |

### Feature Detection

```go
import "github.com/emlgo/eml/internal/eml"

fmt.Printf("Has WASM SIMD: %v\n", eml.HasWasmSIMD())
```

On WASM builds, `HasWasmSIMD()` reports that the WASM dispatch path is selected;
it is not evidence that a particular loop was vectorized. Non-WASM CPU feature
flags return `false`.

## Memory Alignment

WASM linear memory access benefits from alignment. The library provides two utilities:

### WasmAlign16

Returns a float64 subslice that is 16-byte aligned when the slice has an
element available to skip. This is a low-level helper and does not guarantee
that the compiler emits aligned SIMD loads or stores:

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

- Go 1.27.2+ (for WASM target support)
- Node.js (to execute the test harness; browser/runtime support depends on the target)

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
- `simd_wasm.go` — implements block-unrolled kernels for all operations including transcendentals
- `wasm_utils.go` — memory alignment utilities
- `simd_dispatch_stub.go` — excluded from WASM builds (`!wasm` constraint)

## Performance Considerations

- **Batch size**: Optimal performance is achieved with batch sizes >1024 elements. Small batches incur overhead from memory allocation and function call dispatch.
- **Goroutines**: On WASM, goroutines are cooperatively scheduled on a single thread. Parallelized operations (ExpBatch, LogBatch, etc.) provide concurrency but not true parallelism.
- **SIMD lowering**: The block-unrolled loops may help the runtime optimize regular memory access, but do not guarantee `wasm_simd128` instructions. Inspect generated code or profile the selected runtime.
- **Memory alignment**: Use `WasmAlign16()` for critical paths to avoid alignment-related slowdowns in JIT-compiled code.
- **Transcendentals**: Exp, Log, Sin, Cos, Tan, Sinh, Cosh, Tanh, and SinCos use unrolled batch loops around scalar math calls, not vector transcendental instructions.
