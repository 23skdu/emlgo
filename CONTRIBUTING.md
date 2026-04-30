# Contributing to emlgo

Thank you for your interest in contributing to emlgo!

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/emlgo/emlgo.git
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

Run gosec:
```bash
gosec ./...
```

Run govulncheck:
```bash
govulncheck ./...
```

## Pull Request Process

1. Ensure all tests pass
2. Run `go vet ./...` and fix any issues
3. Run gosec and ensure no security issues
4. Update documentation if needed
5. Submit a pull request with a clear description

## Commit Messages

- Use clear, descriptive commit messages
- Reference issues where applicable
- Keep commits focused and atomic

## License

By contributing, you agree that your contributions will be licensed under the MIT License.