#!/usr/bin/env bash
# Git pre-commit hook:
#   1. Prevent committing node dependencies and binary artifacts
#   2. Enforce static analysis & linting (go vet & golangci-lint)
#   3. Enforce test execution and minimum statement code coverage threshold

set -euo pipefail

MIN_COVERAGE="${PRECOMMIT_MIN_COVERAGE:-100.0}"

usage() {
    cat <<EOF
Usage: $0 [options]

Git pre-commit hook enforcing artifact hygiene, linting, and code coverage.

Options:
  --min-coverage <pct>  Minimum total statement coverage percentage (default: ${MIN_COVERAGE})
  -h, --help            Show this help message
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --min-coverage)
            if [[ $# -lt 2 ]]; then
                echo "ERROR: --min-coverage requires a numeric argument" >&2
                exit 1
            fi
            MIN_COVERAGE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "ERROR: Invalid option: $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

STAGED_FILES=$(git diff --cached --name-only 2>/dev/null || true)

# Patterns for forbidden files
FORBIDDEN_PATTERNS=(
    "node_modules/"
    "package-lock.json"
    "yarn.lock"
    "pnpm-lock.yaml"
    "bun.lockb"
    "\.wasm$"
    "\.exe$"
    "\.dll$"
    "\.so(\..*)?$"
    "\.dylib$"
    "\.a$"
    "\.o$"
    "\.test$"
    "\.prof$"
    "\.bin$"
    "^wasmbench$"
    "^validate$"
    "^bench$"
    "^emlcli$"
)

if [ -n "$STAGED_FILES" ]; then
    for file in $STAGED_FILES; do
        for pattern in "${FORBIDDEN_PATTERNS[@]}"; do
            if echo "$file" | grep -Eq "$pattern"; then
                echo "ERROR: Cannot commit forbidden file matching '$pattern': $file" >&2
                exit 1
            fi
        done

        # Check for binary MIME type on existing staged files
        if [ -f "$file" ]; then
            mime=$(file -b --mime-encoding "$file" 2>/dev/null || true)
            if [ "$mime" = "binary" ]; then
                echo "ERROR: Cannot commit binary file: $file (mime: $mime)" >&2
                exit 1
            fi
        fi
    done
fi

echo "[pre-commit 1/3] Running go vet..."
go vet ./...

echo "[pre-commit 2/3] Running golangci-lint..."
if ! command -v golangci-lint >/dev/null 2>&1; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi
golangci-lint run ./...

echo "[pre-commit 3/3] Running tests and verifying statement coverage >= ${MIN_COVERAGE}%..."
TMP_COV="$(mktemp -t emlgo-coverage.XXXXXX.out)"
trap 'rm -f "$TMP_COV"' EXIT

go test -short -count=1 -coverprofile="$TMP_COV" -covermode=atomic ./... >/dev/null
TOTAL_LINE="$(go tool cover -func="$TMP_COV" | tail -n 1)"
TOTAL_PCT="$(echo "$TOTAL_LINE" | awk '{gsub(/%/, "", $3); print $3}')"

if ! awk -v actual="$TOTAL_PCT" -v min="$MIN_COVERAGE" 'BEGIN { exit !(actual + 0.0 >= min + 0.0) }'; then
    echo "ERROR: Total statement coverage ${TOTAL_PCT}% is below minimum required ${MIN_COVERAGE}%" >&2
    exit 1
fi

echo "PASS: Pre-commit checks succeeded (linting clean, statement coverage ${TOTAL_PCT}% >= ${MIN_COVERAGE}%)."
exit 0

