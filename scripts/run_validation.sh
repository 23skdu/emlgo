#!/bin/bash

# Validation script for emlgo library
# Tests all Go math data types (int, uint, float, complex)
#
# Usage: ./scripts/run_validation.sh [options]
#
# Options:
#   -v        Verbose output
#   -f        Show only failed tests
#   -t TYPE   Filter by type (int, uint, float, complex, all)
#   -h        Show this help

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

VALIDATE="./validate"

show_help() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -v        Verbose output"
    echo "  -f        Show only failed tests"
    echo "  -t TYPE   Filter by type (int, uint, float, complex, all)"
    echo "  -h        Show this help"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run all validations"
    echo "  $0 -v                 # Run with verbose output"
    echo "  $0 -t float           # Test only float types"
    echo "  $0 -t complex        # Test only complex types"
    echo "  $0 -f                # Show only failed tests"
}

# Parse arguments
VERBOSE=""
FAILED_ONLY=""
TYPE_FILTER=""

while getopts "vfht:" opt; do
    case $opt in
        v)
            VERBOSE="-v"
            ;;
        f)
            FAILED_ONLY="-f"
            ;;
        t)
            TYPE_FILTER="-type $OPTARG"
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

# Build validation tool
echo "Building validation tool..."
go build -o "$VALIDATE" ./cmd/validate/ || {
    echo "Error: Failed to build validation tool"
    exit 1
}

echo ""
echo "========================================="
echo "Running Validation Tests"
echo "========================================="

# Run validation
"$VALIDATE" $VERBOSE $FAILED_ONLY $TYPE_FILTER

echo ""
echo "Validation complete!"