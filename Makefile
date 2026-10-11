.PHONY: build test test-race test-cover bench lint vet fuzz clean wasm gosec govulncheck docker docker-build docker-build-cpu docker-build-cuda

build:
	go build ./...

wasm:
	@mkdir -p wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" wasm/wasm_exec.js 2>/dev/null || cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o wasm/emlgo.wasm ./cmd/wasmbench

test:
	go test -count=1 ./...

# Fast subset: skips the long genetic-search convergence runs.
test-short:
	go test -count=1 -short ./...

test-race:
	go test -race -count=1 ./...

test-cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

bench:
	go test -bench=. -benchmem -count=3 ./...

bench-stat:
	@command -v benchstat >/dev/null 2>&1 || go install golang.org/x/perf/cmd/benchstat@latest
	go test -bench=. -benchmem -count=6 ./... | benchstat

lint:
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

vet:
	go vet ./...

FUZZTIME ?= 30s

# Targets are anchored so they cannot match more than one fuzz test.
fuzz:
	go test -run='^$$' -fuzz='^FuzzParseNoPanic$$' -fuzztime=$(FUZZTIME) ./internal/jit/...
	go test -run='^$$' -fuzz='^FuzzEval$$' -fuzztime=$(FUZZTIME) ./internal/jit/...
	go test -run='^$$' -fuzz='^FuzzExpSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzLogSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzSinSIMD$$' -fuzztime=$(FUZZTIME) ./internal/eml
	go test -run='^$$' -fuzz='^FuzzPow$$' -fuzztime=$(FUZZTIME) ./pkg/arithmetic
	go test -run='^$$' -fuzz='^FuzzAdd$$' -fuzztime=$(FUZZTIME) ./pkg/arithmetic

gosec:
	gosec -exclude-generated -tags purego ./...

govulncheck:
	govulncheck ./... || go run golang.org/x/vuln/cmd/govulncheck@latest ./...

IMAGE_REPO ?= ghcr.io/23skdu/emlgo
TAG ?= 0.6
CUDA_ARCH ?= sm_89

docker-build-cpu:
	docker build -f Dockerfile.cpu \
		-t $(IMAGE_REPO)-cpu:$(TAG) \
		-t $(IMAGE_REPO)-cpu:latest \
		-t $(IMAGE_REPO):$(TAG)-cpu \
		-t $(IMAGE_REPO):latest-cpu \
		-t $(IMAGE_REPO):$(TAG) \
		-t $(IMAGE_REPO):latest .

docker-build-cuda:
	docker build -f Dockerfile.cuda --build-arg CUDA_ARCH=$(CUDA_ARCH) \
		-t $(IMAGE_REPO)-cuda:$(TAG) \
		-t $(IMAGE_REPO)-cuda:latest \
		-t $(IMAGE_REPO):$(TAG)-cuda \
		-t $(IMAGE_REPO):latest-cuda .

docker-build: docker-build-cpu docker-build-cuda

docker: docker-build

clean:
	rm -f bench validate emlcli wasmbench coverage.out *.out *.test *.prof *.cov
	rm -rf dist/ wasm/emlgo.wasm wasm/wasm_exec.js wasm/run.js node_modules/ .playwright/ playwright-report/ test-results/
