#!/bin/bash
# Benchmark comparison script
# Usage: ./scripts/bench-compare.sh [baseline-ref]

set -e

BASELINE=${1:-main}
BENCH_PATTERN=${2:-.}

echo "=== Benchmarking current code ==="
go test -bench="$BENCH_PATTERN" -benchmem -count=5 ./... > /tmp/bench-current.txt

echo "=== Checking out baseline ($BASELINE) ==="
git stash
git checkout "$BASELINE"

echo "=== Benchmarking baseline ==="
go test -bench="$BENCH_PATTERN" -benchmem -count=5 ./... > /tmp/bench-baseline.txt

echo "=== Restoring current code ==="
git checkout -
git stash pop

echo "=== Comparing results ==="
if command -v benchstat &> /dev/null; then
    benchstat /tmp/bench-baseline.txt /tmp/bench-current.txt
else
    echo "Install benchstat: go install golang.org/x/perf/cmd/benchstat@latest"
    echo "Baseline results:  /tmp/bench-baseline.txt"
    echo "Current results:   /tmp/bench-current.txt"
fi
