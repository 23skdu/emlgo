.PHONY: build test test-short test-race test-metal metal-shaders test-arm64-simd test-docker test-scripts test-cover bench bench-mem bench-stat profile profile-cpu profile-mem lint deadcode lint-shell lint-docker lint-workflows lint-yaml lint-docs lint-tools brew-tools vet fmt-check fuzz fuzz-metal clean clean-docker clean-cache clean-all wasm security gosec govulncheck oracle-gate accuracy-gate verify check makehooks install-hooks docker docker-build docker-build-cpu docker-build-cuda docker-push docker-push-cpu docker-push-cuda help

build: ## Compile all packages
	go build ./...

wasm: ## Build the WebAssembly binary (wasm/emlgo.wasm)
	@mkdir -p wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/wasm_exec.js 2>/dev/null || cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o wasm/emlgo.wasm ./cmd/wasmbench

test: ## Run the full test suite
	go test -count=1 ./...

# Fast subset: skips the long genetic-search convergence runs.
test-short: ## Run fast tests only (skips long convergence runs)
	go test -count=1 -short ./...

test-race: ## Run tests with the race detector
	go test -race -count=1 ./...

metal-shaders: ## Compile and validate Metal shaders with strict IEEE-754 flags (macOS)
	$(MAKE) -C metal all

test-metal: ## Run Apple Metal GPU unit & numerical accuracy tests (darwin/arm64, -tags metal)
	@if command -v xcrun >/dev/null 2>&1 && xcrun -sdk macosx metal --version >/dev/null 2>&1; then \
		$(MAKE) -C metal all && $(MAKE) -C metal clean; \
	fi
	go test -tags metal -race -count=1 -coverprofile=coverage-metal.out -covermode=atomic ./internal/gpu
	go tool cover -func=coverage-metal.out | tail -n 1
	rm -f coverage-metal.out

test-arm64-simd: ## Run ARM64 NEON SIMD tests and purego build-tag fallback check
	go test -race -count=1 -run 'Test(ARM64|Step36|BatchKernel|SIMD)' ./internal/eml
	go test -tags purego -count=1 ./internal/eml

test-docker: ## Run the Docker integration tests
	go test -v -run 'TestDocker.*' ./test/...

test-scripts: ## Run script/hook/gate tests
	go test -v -run 'Test(CheckBenchRegression|RunBenchmark|RunValidation|BenchCompare|WasmTest|GenerateMPFR|OracleAccuracyGate|PreCommit|VerifyHook|MakefileTargets).*' ./test/...

test-cover: ## Run tests with a coverage report
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

# ==============================================================================
# Performance & Memory Profiling Targets
# ==============================================================================

bench: ## Run all benchmarks (3 runs each)
	go test -bench=. -benchmem -count=3 ./...

bench-mem: ## Run memory benchmarks for ./internal/eml and ./pkg
	go test -run='^$$' -bench=. -benchmem -count=1 ./internal/eml/... ./pkg/...

bench-stat: ## Run benchmarks and compare them with benchstat
	@command -v benchstat >/dev/null 2>&1 || go install golang.org/x/perf/cmd/benchstat@latest
	go test -bench=. -benchmem -count=6 ./... | benchstat

profile-cpu: ## CPU profile of the SIMD benchmarks
	go test -run='^$$' -bench='^Benchmark(ExpSIMD|LogSIMD|SinSIMD|CosSIMD|AddSIMD|FmaSIMD)' -benchtime=200ms -benchmem -cpuprofile=cpu.prof ./internal/eml
	go tool pprof -top -nodecount=20 cpu.prof

profile-mem: ## Memory profile of the SIMD benchmarks
	go test -run='^$$' -bench='^Benchmark(ExpSIMD|LogSIMD|Pipeline|Complex)' -benchtime=200ms -benchmem -memprofile=mem.prof ./internal/eml
	go tool pprof -top -alloc_space -nodecount=20 mem.prof

profile: profile-cpu profile-mem ## Run both the CPU and memory profiles

# ==============================================================================
# Linting & Static Analysis Targets
# ==============================================================================

fmt-check: ## Fail if any Go file is not gofmt'd
	@unformatted="$$(gofmt -s -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "ERROR: Unformatted Go files detected (run 'gofmt -s -w .'):"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet: ## Run go vet
	go vet ./...

lint: fmt-check vet ## fmt-check + vet + golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
	golangci-lint run ./...

deadcode: ## Check for unreachable functions across packages and tests
	@DEADCODE="$$(command -v deadcode 2>/dev/null || echo "$$(go env GOPATH)/bin/deadcode")"; \
	if [ ! -x "$$DEADCODE" ]; then go install golang.org/x/tools/cmd/deadcode@latest; fi; \
	out="$$("$$DEADCODE" -test ./...)"; \
	if [ -n "$$out" ]; then \
		echo "ERROR: Unreachable functions detected:"; \
		echo "$$out"; \
		exit 1; \
	fi; \
	if [ "$$(uname -s)" = "Darwin" ] && [ "$$(uname -m)" = "arm64" ]; then \
		out="$$("$$DEADCODE" -test -tags metal ./internal/gpu)"; \
		if [ -n "$$out" ]; then \
			echo "ERROR: Unreachable Metal functions detected:"; \
			echo "$$out"; \
			exit 1; \
		fi; \
	fi; \
	echo "PASS: 0 unreachable functions detected"

lint-shell: ## Lint shell scripts with shellcheck and shfmt
	@command -v shellcheck >/dev/null 2>&1 || { echo "shellcheck not installed (run 'make brew-tools')"; exit 1; }
	shellcheck -S warning scripts/*.sh
	@command -v shfmt >/dev/null 2>&1 || { echo "shfmt not installed (run 'make brew-tools')"; exit 1; }
	shfmt -i 4 -d scripts/verify.sh scripts/oracle_accuracy_gate.sh >/dev/null 2>&1 || true

lint-docker: ## Lint Dockerfiles with hadolint
	@command -v hadolint >/dev/null 2>&1 || { echo "hadolint not installed (run 'make brew-tools')"; exit 1; }
	hadolint --failure-threshold error Dockerfile Dockerfile.cpu Dockerfile.cuda

lint-workflows: ## Lint GitHub Actions workflows with actionlint
	@command -v actionlint >/dev/null 2>&1 || { echo "actionlint not installed (run 'make brew-tools')"; exit 1; }
	actionlint

lint-yaml: ## Lint YAML configuration files with yamllint
	@command -v yamllint >/dev/null 2>&1 || { echo "yamllint not installed (run 'make brew-tools')"; exit 1; }
	yamllint -d "{extends: relaxed, rules: {line-length: {max: 160}}}" .github/workflows docker-compose.yml .golangci.yml

lint-docs: ## Lint Markdown documentation with markdownlint
	@command -v markdownlint >/dev/null 2>&1 || { echo "markdownlint not installed (run 'make brew-tools')"; exit 1; }
	markdownlint --disable MD013 MD033 MD041 README.md CONTRIBUTING.md CHANGELOG.md docs/*.md

lint-tools: lint-shell lint-docker lint-workflows ## Run external brew-installed codebase linters

brew-tools: ## Install Homebrew validation tools (macOS)
	brew install golangci-lint govulncheck gosec shellcheck shfmt hadolint actionlint yamllint markdownlint-cli mpfr

# ==============================================================================
# Fuzz Testing Targets
# ==============================================================================

FUZZTIME ?= 30s

# Targets are anchored so they cannot match more than one fuzz test.
fuzz: ## Run all fuzz tests (override duration with FUZZTIME=60s)
	go test -run='^$$' -fuzz='^FuzzParseNoPanic$$' -fuzztime=$(FUZZTIME) ./internal/jit/...
	go test -run='^$$' -fuzz='^FuzzEval$$' -fuzztime=$(FUZZTIME) ./internal/jit/...
	go test -run='^$$' -fuzz='^FuzzCompileMatchesEval$$' -fuzztime=$(FUZZTIME) ./internal/jit/...
	go test -run='^$$' -fuzz='^FuzzCanonicalSimplifyAndDiff$$' -fuzztime=$(FUZZTIME) ./internal/jit/...
	go test -run='^$$' -fuzz='^FuzzExpSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzLogSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzSinSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzEmlSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzSinCosSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzComplexBatchDifferential$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzDotInt8AndSatInt8$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzExpLnRoundtrip$$' -fuzztime=$(FUZZTIME) ./internal/eml/bigmath
	go test -run='^$$' -fuzz='^FuzzBigMathTrigAndInverseIdentities$$' -fuzztime=$(FUZZTIME) ./internal/eml/bigmath
	go test -run='^$$' -fuzz='^FuzzLaunchConfig$$' -fuzztime=$(FUZZTIME) ./internal/gpu
	go test -run='^$$' -fuzz='^FuzzBatchVerifierVerifyOp$$' -fuzztime=$(FUZZTIME) ./internal/gpu
	go test -run='^$$' -fuzz='^FuzzPow$$' -fuzztime=$(FUZZTIME) ./pkg/arithmetic
	go test -run='^$$' -fuzz='^FuzzAdd$$' -fuzztime=$(FUZZTIME) ./pkg/arithmetic
	go test -run='^$$' -fuzz='^FuzzGCDAndLCM$$' -fuzztime=$(FUZZTIME) ./pkg/arithmetic
	go test -run='^$$' -fuzz='^FuzzInt8DistanceAndDot$$' -fuzztime=$(FUZZTIME) ./pkg/arithmetic
	go test -run='^$$' -fuzz='^FuzzOptimizerDifferential$$' -fuzztime=$(FUZZTIME) ./pkg/bytecode
	go test -run='^$$' -fuzz='^FuzzEvalBatchColumnarDifferential$$' -fuzztime=$(FUZZTIME) ./pkg/bytecode
	go test -run='^$$' -fuzz='^FuzzFastMathDifferential$$' -fuzztime=$(FUZZTIME) ./pkg/fastmath
	go test -run='^$$' -fuzz='^FuzzGuardrailsAndRegularizedEml$$' -fuzztime=$(FUZZTIME) ./pkg/fastmath
	go test -run='^$$' -fuzz='^FuzzHyperbolicAccuracy$$' -fuzztime=$(FUZZTIME) ./pkg/hyper
	go test -run='^$$' -fuzz='^FuzzHyperbolicInverseRoundtrip$$' -fuzztime=$(FUZZTIME) ./pkg/hyper
	go test -run='^$$' -fuzz='^FuzzLogExpDifferential$$' -fuzztime=$(FUZZTIME) ./pkg/logexp
	go test -run='^$$' -fuzz='^FuzzLogExpBatchAndSubnormalOracle$$' -fuzztime=$(FUZZTIME) ./pkg/logexp
	go test -run='^$$' -fuzz='^FuzzPackUnpack4Bit$$' -fuzztime=$(FUZZTIME) ./pkg/quant
	go test -run='^$$' -fuzz='^FuzzTurboQuant4EncodeDistance$$' -fuzztime=$(FUZZTIME) ./pkg/quant
	go test -run='^$$' -fuzz='^FuzzTrigPythagoreanIdentity$$' -fuzztime=$(FUZZTIME) ./pkg/trig
	go test -run='^$$' -fuzz='^FuzzTrigReciprocalAndInverse$$' -fuzztime=$(FUZZTIME) ./pkg/trig

fuzz-metal: ## Run Apple Metal GPU & ARM64 NEON SIMD fuzz tests (darwin/arm64)
	go test -run='^$$' -fuzz='^FuzzARM64NEONKernels$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -tags metal -run='^$$' -fuzz='^FuzzMetalFloat32Kernels$$' -fuzztime=$(FUZZTIME) ./internal/gpu
	go test -tags metal -run='^$$' -fuzz='^FuzzMetalComplex64Kernels$$' -fuzztime=$(FUZZTIME) ./internal/gpu

# ==============================================================================
# Go Security & Vulnerability Targets
# ==============================================================================

gosec: ## Run gosec static analysis
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
	golangci-lint run --enable-only=gosec ./...

govulncheck: ## Check dependencies for known vulnerabilities
	govulncheck ./... || go run golang.org/x/vuln/cmd/govulncheck@latest ./...

security: gosec govulncheck ## gosec + govulncheck

# ==============================================================================
# Accuracy, Oracle & Pre-Commit Quality Gates
# ==============================================================================

oracle-gate: ## Run the oracle accuracy gate script
	./scripts/oracle_accuracy_gate.sh

accuracy-gate: oracle-gate ## Alias for oracle-gate

verify: ## Run scripts/verify.sh (full verification)
	./scripts/verify.sh

check: vet fmt-check lint security oracle-gate test-cover ## Full quality gate: vet, fmt-check, lint, security, oracle-gate, test-cover

makehooks: ## Install the git pre-commit hook
	chmod +x scripts/verify.sh scripts/pre-commit.sh
	cp scripts/verify.sh scripts/pre-commit.sh
	cp scripts/verify.sh .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

install-hooks: makehooks ## Alias for makehooks

IMAGE_REPO ?= ghcr.io/23skdu/emlgo
TAG ?= v0.8.0
CUDA_ARCH ?= sm_89

docker-build-cpu: ## Build the CPU Docker image
	docker build -f Dockerfile.cpu \
		-t $(IMAGE_REPO):$(TAG) \
		-t $(IMAGE_REPO):latest \
		-t $(IMAGE_REPO):$(TAG)-cpu \
		-t $(IMAGE_REPO):latest-cpu \
		-t $(IMAGE_REPO):cpu .

docker-build-cuda: ## Build the CUDA Docker image (override with CUDA_ARCH=sm_XX)
	docker build -f Dockerfile.cuda --build-arg CUDA_ARCH=$(CUDA_ARCH) \
		-t $(IMAGE_REPO):$(TAG)-cuda \
		-t $(IMAGE_REPO):latest-cuda \
		-t $(IMAGE_REPO):cuda .

docker-build: docker-build-cpu docker-build-cuda ## Build the CPU and CUDA Docker images

docker: docker-build ## Alias for docker-build

docker-push-cpu: ## Push the CPU Docker image tags
	docker push $(IMAGE_REPO):$(TAG)
	docker push $(IMAGE_REPO):latest
	docker push $(IMAGE_REPO):$(TAG)-cpu
	docker push $(IMAGE_REPO):latest-cpu
	docker push $(IMAGE_REPO):cpu

docker-push-cuda: ## Push the CUDA Docker image tags
	docker push $(IMAGE_REPO):$(TAG)-cuda
	docker push $(IMAGE_REPO):latest-cuda
	docker push $(IMAGE_REPO):cuda

docker-push: docker-push-cpu docker-push-cuda ## Push all Docker image tags

# ==============================================================================
# Cleanup Targets
# ==============================================================================
#
# `clean` deliberately performs no Docker or network work: the CUDA Dockerfiles
# run `make clean` inside their build stage, so it must stay a pure
# filesystem cleanup. Docker images/caches live in `clean-docker`.

clean: ## Remove test, build, coverage, profiling and Python venv artifacts
	rm -f bench validate emlcli wasmbench coverage.out *.out *.test *.prof *.cov *.coverprofile
	rm -f *.o *.a *.bin test_*.o
	rm -f cuda/*.o cuda/*.a cuda/libeml_capi.so cuda/eml_cuda_example
	rm -f metal/*.air metal/*.metallib
	rm -rf dist/ bin/ wasm/emlgo.wasm wasm/wasm_exec.js wasm/run.js
	rm -rf node_modules/ .playwright/ playwright/ playwright-report/ test-results/
	rm -rf .venv/ venv/ env/ .pytest_cache/ .mypy_cache/ *.egg-info/
	find . -path ./.git -prune -o -name '__pycache__' -type d -prune -exec rm -rf {} +
	find . -path ./.git -prune -o -name '*.pyc' -type f -exec rm -f {} +

clean-docker: ## Remove this repository's tagged Docker images and the build cache
	@command -v docker >/dev/null 2>&1 || { echo "docker not found; nothing to clean"; exit 0; }
	@docker images --format '{{.Repository}}:{{.Tag}}' 2>/dev/null \
		| grep -F "$(IMAGE_REPO):" \
		| while read -r img; do docker rmi -f "$$img" || true; done
	@docker builder prune -f || true
	@docker image prune -f || true

clean-cache: ## Remove the Go test cache
	go clean -testcache

clean-all: clean clean-docker clean-cache ## Remove artifacts, venvs, Docker images and caches

# ==============================================================================
# Help
# ==============================================================================

help: ## List all available targets and variables
	@echo "Usage: make <target> [VARIABLE=value ...]"
	@echo ""
	@echo "Targets:"
	@grep -hE '^[A-Za-z0-9_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## "} {t = $$1; sub(/:.*/, "", t); printf "  %-18s %s\n", t, $$2}'
	@echo ""
	@echo "Variables (override on the command line):"
	@echo "  FUZZTIME=$(FUZZTIME)          Duration per fuzz test (default 30s)"
	@echo "  IMAGE_REPO=$(IMAGE_REPO)      Docker image repository (default ghcr.io/23skdu/emlgo)"
	@echo "  TAG=$(TAG)                    Docker image tag (default v0.8.0)"
	@echo "  CUDA_ARCH=$(CUDA_ARCH)        CUDA arch for cuda images (default sm_89)"
