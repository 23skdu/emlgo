#!/bin/bash

# Benchmark script for emlgo vs Go math library
# Usage: ./scripts/run_benchmark.sh [options]
#
# Options:
#   -n N      Number of iterations (default: 1000000)
#   -c        Run feature parity comparison
#   -a        Run accuracy (ULP) tests
#   -v        Verbose output
#   -h        Show this help

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

BENCHMARK="./bench"
DEFAULT_N=1000000

show_help() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -n N      Number of iterations (default: $DEFAULT_N)"
    echo "  -c        Run feature parity comparison"
    echo "  -a        Run accuracy (ULP) tests"
    echo "  -v        Verbose output"
    echo "  -h        Show this help"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run speed benchmark"
    echo "  $0 -n 500000          # Run with custom iterations"
    echo "  $0 -c                 # Run feature parity test"
    echo "  $0 -a                 # Run accuracy test"
    echo "  $0 -c -a              # Run all tests"
}

# Parse arguments
N=$DEFAULT_N
PARITY=""
ACCURACY=""
VERBOSE=""

while getopts "n:cavh" opt; do
    case $opt in
        n)
            N=$OPTARG
            ;;
        c)
            PARITY="-compare"
            ;;
        a)
            ACCURACY="-accuracy"
            ;;
        v)
            VERBOSE="-v"
            ;;
        h)
            show_help
            exit 0
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            show_help
            exit 1
            ;;
    esac
done

# Build benchmark tool
echo "Building benchmark tool..."
go build -o "$BENCHMARK" ./cmd/bench/ || {
    echo "Error: Failed to build benchmark tool"
    exit 1
}

# Determine what to run
if [ -z "$PARITY" ] && [ -z "$ACCURACY" ]; then
    # Default: run speed benchmark
    echo ""
    echo "========================================="
    echo "Running Speed Benchmark (n=$N)"
    echo "========================================="
    "$BENCHMARK" -n "$N"
elif [ -n "$PARITY" ] || [ -n "$ACCURACY" ]; then
    # Run parity and/or accuracy tests
    if [ -n "$PARITY" ]; then
        echo ""
        echo "========================================="
        echo "Running Feature Parity Test"
        echo "========================================="
        "$BENCHMARK" $PARITY $VERBOSE
    fi
    
    if [ -n "$ACCURACY" ]; then
        echo ""
        echo "========================================="
        echo "Running Accuracy Test (ULP)"
        echo "========================================="
        "$BENCHMARK" $ACCURACY
    fi
fi

echo ""
echo "Benchmark complete!"