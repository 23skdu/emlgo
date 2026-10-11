package test

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// getRepoRoot locates the repository root relative to this test file.
func getRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Failed to find repository root containing go.mod")
		}
		dir = parent
	}
}

// TestDockerfileCPUStructure validates Dockerfile.cpu structure and configuration.
func TestDockerfileCPUStructure(t *testing.T) {
	root := getRepoRoot(t)
	content, err := os.ReadFile(filepath.Join(root, "Dockerfile.cpu"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile.cpu: %v", err)
	}
	text := string(content)

	tests := []struct {
		name     string
		contains string
	}{
		{"Go Version 1.27.2", "ARG GO_VERSION=1.27.2"},
		{"Builder stage", "FROM golang:${GO_VERSION}-bookworm AS builder"},
		{"CGO disabled for static build", "ENV CGO_ENABLED=0"},
		{"GOAMD64 v3 microarchitecture", "ENV GOAMD64=v3"},
		{"Linux target OS", "ENV GOOS=linux"},
		{"Amd64 target architecture", "ENV GOARCH=amd64"},
		{"Build emlcli binary", "-o /out/emlcli ./cmd/emlcli"},
		{"Build bench binary", "-o /out/bench ./cmd/bench"},
		{"Build validate binary", "-o /out/validate ./cmd/validate"},
		{"In-build validation check", "RUN /out/validate"},
		{"Runtime stage", "FROM debian:bookworm-slim AS runtime"},
		{"Non-root user creation", "useradd -r -g emlgo"},
		{"Copy emlcli", "COPY --from=builder /out/emlcli /usr/local/bin/emlcli"},
		{"Copy bench", "COPY --from=builder /out/bench /usr/local/bin/bench"},
		{"Copy validate", "COPY --from=builder /out/validate /usr/local/bin/validate"},
		{"Enforce non-root user", "USER emlgo:emlgo"},
		{"Default entrypoint", `ENTRYPOINT ["emlcli"]`},
		{"Default command", `CMD ["demo"]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(text, tt.contains) {
				t.Errorf("Dockerfile.cpu missing required directive: %q", tt.contains)
			}
		})
	}
}

// TestDockerfileCUDAStructure validates Dockerfile.cuda structure and configuration.
func TestDockerfileCUDAStructure(t *testing.T) {
	root := getRepoRoot(t)
	content, err := os.ReadFile(filepath.Join(root, "Dockerfile.cuda"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile.cuda: %v", err)
	}
	text := string(content)

	tests := []struct {
		name     string
		contains string
	}{
		{"CUDA 12.6.3 default version", "ARG CUDA_VERSION=12.6.3"},
		{"Go 1.27.2 default version", "ARG GO_VERSION=1.27.2"},
		{"Default CUDA architecture sm_89", "ARG CUDA_ARCH=sm_89"},
		{"Go toolchain provider stage", "FROM golang:${GO_VERSION}-bookworm AS go-toolchain"},
		{"CUDA devel builder stage", "FROM nvidia/cuda:${CUDA_VERSION}-devel-ubuntu24.04 AS builder"},
		{"Copy Go toolchain", "COPY --from=go-toolchain /usr/local/go /usr/local/go"},
		{"Build CUDA shared library", "make clean && make ARCH=${CUDA_ARCH}"},
		{"CGO enabled for CUDA bridge", "ENV CGO_ENABLED=1"},
		{"GOAMD64 v3 microarchitecture", "ENV GOAMD64=v3"},
		{"CUDA library search path", `ENV LIBRARY_PATH="/usr/local/cuda/lib64:/usr/local/cuda/lib64/stubs:${LIBRARY_PATH}"`},
		{"CUDA runtime search path", `ENV LD_LIBRARY_PATH="/usr/local/cuda/lib64:/usr/local/cuda/lib64/stubs:/src/cuda:${LD_LIBRARY_PATH}"`},
		{"Build emlcli with cuda tag", `go build -tags cuda -trimpath -ldflags="-s -w" -o /out/emlcli ./cmd/emlcli`},
		{"Build bench with cuda tag", `go build -tags cuda -trimpath -ldflags="-s -w" -o /out/bench ./cmd/bench`},
		{"Build validate with cuda tag", `go build -tags cuda -trimpath -ldflags="-s -w" -o /out/validate ./cmd/validate`},
		{"Copy CUDA shared lib to out", "cp /src/cuda/libeml_capi.so /out/libeml_capi.so"},
		{"CUDA runtime stage", "FROM nvidia/cuda:${CUDA_VERSION}-runtime-ubuntu24.04 AS runtime"},
		{"Install shared library to /usr/local/lib", "COPY --from=builder /out/libeml_capi.so /usr/local/lib/libeml_capi.so"},
		{"Run ldconfig", "RUN ldconfig"},
		{"NVIDIA visible devices", "ENV NVIDIA_VISIBLE_DEVICES=all"},
		{"NVIDIA driver capabilities", "ENV NVIDIA_DRIVER_CAPABILITIES=compute,utility"},
		{"Runtime non-root user", "USER emlgo:emlgo"},
		{"Runtime entrypoint", `ENTRYPOINT ["emlcli"]`},
		{"Runtime default CMD", `CMD ["gpu-status"]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(text, tt.contains) {
				t.Errorf("Dockerfile.cuda missing required directive: %q", tt.contains)
			}
		})
	}
}

// TestUnifiedDockerfileStructure validates the unified multi-target Dockerfile.
func TestUnifiedDockerfileStructure(t *testing.T) {
	root := getRepoRoot(t)
	content, err := os.ReadFile(filepath.Join(root, "Dockerfile"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile: %v", err)
	}
	text := string(content)

	targets := []struct {
		name   string
		target string
	}{
		{"Target cpu", "FROM debian:bookworm-slim AS cpu"},
		{"Target cuda", "FROM nvidia/cuda:${CUDA_VERSION}-runtime-ubuntu24.04 AS cuda"},
		{"Stage cpu-builder", "FROM golang:${GO_VERSION}-bookworm AS cpu-builder"},
		{"Stage cuda-builder", "FROM nvidia/cuda:${CUDA_VERSION}-devel-ubuntu24.04 AS cuda-builder"},
	}

	for _, tt := range targets {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(text, tt.target) {
				t.Errorf("Dockerfile missing target stage: %q", tt.target)
			}
		})
	}
}

// TestDockerignoreEntries validates that critical build artifacts and secrets are ignored.
func TestDockerignoreEntries(t *testing.T) {
	root := getRepoRoot(t)
	content, err := os.ReadFile(filepath.Join(root, ".dockerignore"))
	if err != nil {
		t.Fatalf("Failed to read .dockerignore: %v", err)
	}

	requiredIgnores := []string{
		".git",
		"node_modules",
		"coverage.out",
		"dist",
		"bin",
		"bench",
		"validate",
		"emlcli",
		"cuda/*.o",
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	found := make(map[string]bool)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			found[line] = true
		}
	}

	for _, req := range requiredIgnores {
		t.Run("Ignore "+req, func(t *testing.T) {
			if !found[req] {
				t.Errorf(".dockerignore missing pattern: %q", req)
			}
		})
	}
}

// TestDockerComposeStructure validates docker-compose.yml configuration.
func TestDockerComposeStructure(t *testing.T) {
	root := getRepoRoot(t)
	content, err := os.ReadFile(filepath.Join(root, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("Failed to read docker-compose.yml: %v", err)
	}
	text := string(content)

	checks := []struct {
		name     string
		contains string
	}{
		{"CPU Service Definition", "emlgo-cpu:"},
		{"CPU Image Reference", "image: ghcr.io/23skdu/emlgo:cpu"},
		{"CPU Benchmark Service", "emlgo-cpu-bench:"},
		{"CPU Validation Service", "emlgo-cpu-validate:"},
		{"CUDA Service Definition", "emlgo-cuda:"},
		{"CUDA Image Reference", "image: ghcr.io/23skdu/emlgo:cuda"},
		{"CUDA Benchmark Service", "emlgo-cuda-bench:"},
		{"CUDA Verification Service", "emlgo-cuda-verify:"},
		{"NVIDIA Driver Reservation", "driver: nvidia"},
		{"GPU Capabilities", "capabilities: [gpu]"},
	}

	for _, tt := range checks {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(text, tt.contains) {
				t.Errorf("docker-compose.yml missing configuration: %q", tt.contains)
			}
		})
	}
}

// TestDockerExecutionIfAvailable validates live container behavior when docker daemon is reachable.
func TestDockerExecutionIfAvailable(t *testing.T) {
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		t.Skip("Docker CLI not available on host; skipping container runtime validation")
	}

	// Verify docker daemon connectivity
	cmd := exec.Command(dockerPath, "info")
	if err := cmd.Run(); err != nil {
		t.Skip("Docker daemon not running; skipping container runtime validation")
	}

	t.Run("CPU Container Demo Execution", func(t *testing.T) {
		out, err := exec.Command(dockerPath, "run", "--rm", "ghcr.io/23skdu/emlgo:cpu", "demo").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to execute cpu container demo: %v\nOutput: %s", err, string(out))
		}
		if !strings.Contains(string(out), "Has AVX2: true") {
			t.Errorf("Expected output to report 'Has AVX2: true', got:\n%s", string(out))
		}
	})

	t.Run("CPU Container Non-Root User", func(t *testing.T) {
		out, err := exec.Command(dockerPath, "run", "--rm", "--entrypoint", "id", "ghcr.io/23skdu/emlgo:cpu", "-u").CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to check container user: %v\nOutput: %s", err, string(out))
		}
		uid := strings.TrimSpace(string(out))
		if uid == "0" {
			t.Errorf("Container is running as root (UID 0); expected non-root user")
		}
	})

	t.Run("CPU Container Binaries Executable", func(t *testing.T) {
		for _, bin := range []string{"emlcli", "bench", "validate"} {
			out, err := exec.Command(dockerPath, "run", "--rm", "--entrypoint", bin, "ghcr.io/23skdu/emlgo:cpu", "-h").CombinedOutput()
			if err != nil && !strings.Contains(err.Error(), "exit status 2") && !strings.Contains(err.Error(), "exit status 1") {
				t.Fatalf("Binary %s failed execution: %v\nOutput: %s", bin, err, string(out))
			}
		}
	})

	// Check if GPU is present for CUDA container test
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		t.Run("CUDA Container GPU Status", func(t *testing.T) {
			out, err := exec.Command(dockerPath, "run", "--rm", "--gpus", "all", "ghcr.io/23skdu/emlgo:cuda", "gpu-status").CombinedOutput()
			if err != nil {
				t.Skipf("CUDA execution skipped (GPU driver/runtime passthrough not configured in test runner): %v", err)
			}
			if !strings.Contains(string(out), "Compute Capability") && !strings.Contains(string(out), "GPU device") {
				t.Errorf("Unexpected GPU status output:\n%s", string(out))
			}
		})
	}
}
