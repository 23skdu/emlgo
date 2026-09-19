.PHONY: build test test-race test-cover bench lint vet fuzz clean

build:
	go build ./...

test:
	go test -count=1 ./...

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

fuzz:
	go test -run='^$$' -fuzz=FuzzParse -fuzztime=30s ./internal/jit/...
	go test -run='^$$' -fuzz=FuzzEml -fuzztime=30s ./internal/eml/...
	go test -run='^$$' -fuzz=FuzzArithmetic -fuzztime=30s ./pkg/arithmetic/...

gosec:
	gosec -exclude-generated -tags purego ./...

clean:
	rm -f /bench /validate /emlcli coverage.out
	rm -rf dist/
