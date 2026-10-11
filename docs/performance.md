# Performance Measurement

## Status of published results

This repository does not currently publish portable performance results or
cross-platform speedup claims. Earlier tables lacked retained raw output or
consistent environment metadata and conflict with one another; they have been
withdrawn. Treat the results from your own workload and target machine as
authoritative.

## Reproduce batch benchmarks

Run the same commit on an otherwise idle machine and retain the complete output:

```sh
go version
go env GOOS GOARCH
go test -run='^$' -bench='^(BenchmarkSIMD|BenchmarkExpSIMD|BenchmarkLogSIMD|BenchmarkAddSIMD)$' \
  -benchmem -benchtime=250ms -count=10 ./internal/eml
```

Record the commit, CPU model, OS/kernel, Go version, `GOMAXPROCS`, build tags,
benchmark names, input sizes, and raw output. Compare only equivalent runs on
the same runner; do not combine ratios across architectures or infer speedups
from a single sample. For package-specific questions, select the relevant
benchmarks and use the same settings for both revisions.

## Pull-request regression gate

Pull requests compare the base and proposed revisions on the same
`ubuntu-24.04` runner using ten 250 ms samples for the selected batch
benchmarks. [`scripts/check_bench_regression.py`](../scripts/check_bench_regression.py)
uses a deterministic bootstrap 95% lower confidence bound and fails only when
that bound indicates a slowdown greater than 10%. This gate covers those
benchmarks on that runner; it is not a general performance certification.

## Useful local commands

```sh
# All package benchmarks, three samples each
make bench

# Repeated benchmark set with statistical summary (installs benchstat if absent)
make bench-stat

# Compare against a baseline ref with the repository helper
./scripts/bench-compare.sh <baseline-ref>
```

For numerical error budgets and accuracy evidence, see
[`accuracy.md`](accuracy.md). A fast benchmark result does not establish
numerical correctness; report performance and accuracy evidence separately.
