# Contributing to emlgo

Thank you for your interest in contributing to emlgo!

## Development Setup

1. Clone the repository (the Go module path is `github.com/emlgo/eml`):
   ```bash
   git clone https://github.com/23skdu/emlgo.git
   cd emlgo
   ```

2. Install development dependencies:
   ```bash
   go install github.com/securego/gosec/v2/cmd/gosec@latest
   go install golang.org/x/vuln/cmd/govulncheck@latest
   ```

## Code Style

- Run `go fmt` before committing
- Run `gofmt -s` to check formatting
- Ensure all public functions have godoc comments
- Follow standard Go naming conventions

## Testing

Run all tests including race detection:
```bash
go test -race ./...
```

Run benchmarks:
```bash
go test -bench=. ./...
```

## Security Scanning

Run gosec (use the same invocation as the Makefile and CI):
```bash
gosec -exclude-generated -tags purego ./...
```

Run govulncheck:
```bash
govulncheck ./...
```

## Dependency Management

Dependencies are managed via Go modules. Dependabot is configured to automatically create PRs for:
- **GitHub Actions**: weekly updates to action versions
- **Go modules**: weekly patch-level updates (minor/major bumps require manual review)

When updating `golang.org/x/sys` or other dependencies manually:
1. Run `go test -race ./...` after updating to verify no regressions
2. Run `gosec -tags purego ./...` to check for security issues
3. Check cross-compilation for every target in the CI `cross-build` job:
   ```bash
   for t in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 js/wasm wasip1/wasm; do
     GOOS=${t%/*} GOARCH=${t#*/} go build ./... || echo "FAILED: $t"
   done
   ```
   Note that `GOOS=wasip1` and `GOOS=js` cannot run the JIT: executable memory is
   unavailable there, so `internal/jit` falls back to the AST interpreter on
   `js/wasm` and returns an error from `Compile` on `wasip1`.

## Pull Request Process

1. Ensure all tests pass (`make test-race`)
2. Run `make lint` and fix any issues
3. Run `make gosec` and ensure no security issues
4. Run `make test-cover` and keep package coverage from regressing
5. Update documentation if needed — `README.md` and `docs/*.md` are held to a high
   standard: identifiers, signatures, code samples and measured numbers must match
   the code. Prefer adding a regression test over adding a claim.
6. Submit a pull request with a clear description

## Commit Messages

- Use clear, descriptive commit messages
- Reference issues where applicable
- Keep commits focused and atomic
- Prefix commits: `ci:` for CI changes, `deps:` for dependency updates, `docs:` for documentation

## License

By contributing, you agree that your contributions will be licensed under The Unlicense (public domain).