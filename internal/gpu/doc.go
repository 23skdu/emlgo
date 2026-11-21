// Package gpu provides GPU batch computation backends for CUDA and Metal.
//
// The package uses build tags to conditionally compile GPU support:
//   - "cuda" tag (with cgo): NVIDIA CUDA backend for Linux/Windows
//   - "metal" tag (with cgo on darwin/arm64): Apple Metal backend for macOS Apple Silicon
//
// Without GPU support, all batch operations return errors and Init/Shutdown are no-ops.
// The package includes a BatchVerifier for comparing GPU results against CPU reference
// implementations using ULP (Units in the Last Place) distance metrics.
package gpu
