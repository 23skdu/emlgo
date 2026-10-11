# emlgo Documentation

Welcome to the emlgo documentation. This directory contains detailed documentation for the library.

Examples and references for packages under `internal/` describe implementation
components only; Go prevents downstream modules from importing those packages.

## Table of Contents

| Document | Description |
| :--- | :--- |
| [../README.md](../README.md) | Project overview, features, and quick start |
| [architecture.md](architecture.md) | System architecture and design principles |
| [EML.md](EML.md) | EML operator theory (Odrzywołek, arXiv:2603.21852) |
| [FUNCTIONS.md](FUNCTIONS.md) | Function reference for public packages and selected internal APIs |
| [JIT.md](JIT.md) | JIT compiler, arena allocator, canonical trees, Diff/Simplify/Decompile |
| [COMPLEX.md](COMPLEX.md) | Complex128 number support |
| [USAGE.md](USAGE.md) | Practical usage examples |
| [performance.md](performance.md) | Reproducible benchmark methodology and regression gate |
| [comparison.md](comparison.md) | Status of historical emlgo vs math performance claims |
| [wasm.md](wasm.md) | WebAssembly targets, batch kernels, and test harness |
| [accuracy.md](accuracy.md) | Numerical accuracy contracts and validation limits |
| [nextsteps.md](nextsteps.md) | Roadmap and verification status |

## Quick Links

- **Getting Started**: See [../README.md](../README.md#quick-start)
- **API Reference**: See [FUNCTIONS.md](FUNCTIONS.md)
- **Architecture**: See [architecture.md](architecture.md)
- **Benchmarks**: See [performance.md](performance.md)
- **Roadmap / known bottlenecks**: See [nextsteps.md](nextsteps.md)
