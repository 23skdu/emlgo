.PHONY: build test test-race test-cover bench lint vet fuzz clean wasm

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

clean:
	rm -f bench validate emlcli wasmbench coverage.out *.out *.test *.prof *.cov
	rm -rf dist/ wasm/emlgo.wasm wasm/wasm_exec.js wasm/run.js node_modules/ .playwright/ playwright-report/ test-results/
