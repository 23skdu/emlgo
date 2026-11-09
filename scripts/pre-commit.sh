#!/usr/bin/env bash
# Git pre-commit & verification hook (scripts/verify.sh):
#   1. Prevent committing node dependencies and binary/profiling artifacts
#   2. Enforce static analysis, formatting & linting (go vet & golangci-lint)
#   3. Enforce Go security checks (gosec static security audit)
#   4. Enforce numerical oracle & MPFR 4.2.2 accuracy quality gate
#   5. Enforce unit test execution and minimum statement code coverage threshold
#   6. Optional CPU & memory profiling verification (--with-profile)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
if [[ "$SCRIPT_DIR" == *"/.git/hooks" ]]; then
    PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
fi

MIN_COVERAGE="${PRECOMMIT_MIN_COVERAGE:-100.0}"
WITH_PROFILE=false
SKIP_ORACLE=false

usage() {
    cat <<EOF
Usage: $0 [options]

Git pre-commit hook enforcing artifact hygiene, security, oracle accuracy, linting, and code coverage.

Options:
  --min-coverage <pct>  Minimum total statement coverage percentage (default: ${MIN_COVERAGE})
  --with-profile        Run CPU and memory allocation profiling checks
  --skip-oracle         Skip the fast oracle accuracy gate check
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
        --with-profile)
            WITH_PROFILE=true
            shift
            ;;
        --skip-oracle)
            SKIP_ORACLE=true
            shift
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

cd "$PROJECT_DIR"

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

echo "[verify 1/5] Running go vet..."
go vet ./...

echo "[verify 2/5] Running golangci-lint (formatting, staticcheck, govet, errcheck, misspell)..."
if ! command -v golangci-lint >/dev/null 2>&1; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi
golangci-lint run ./...

echo "[verify 3/5] Running Go security check (gosec)..."
golangci-lint run --enable-only=gosec ./...

if [ "$SKIP_ORACLE" = false ]; then
    echo "[verify 4/5] Running Oracle & MPFR 4.2.2 accuracy quality gate (--quick)..."
    ./scripts/oracle_accuracy_gate.sh --quick >/dev/null
else
    echo "[verify 4/5] Skipped Oracle accuracy gate (--skip-oracle)"
fi

echo "[verify 5/5] Running unit tests and verifying statement coverage >= ${MIN_COVERAGE}%..."
TMP_COV="$(mktemp -t emlgo-coverage.XXXXXX.out)"
trap 'rm -f "$TMP_COV"' EXIT

go test -short -count=1 -coverprofile="$TMP_COV" -covermode=atomic ./... >/dev/null
TOTAL_LINE="$(go tool cover -func="$TMP_COV" | tail -n 1)"
TOTAL_PCT="$(echo "$TOTAL_LINE" | awk '{gsub(/%/, "", $3); print $3}')"

if ! awk -v actual="$TOTAL_PCT" -v min="$MIN_COVERAGE" 'BEGIN { exit !(actual + 0.0 >= min + 0.0) }'; then
    echo "ERROR: Total statement coverage ${TOTAL_PCT}% is below minimum required ${MIN_COVERAGE}%" >&2
    exit 1
fi

if [ "$WITH_PROFILE" = true ]; then
    echo "[verify +profile] Running CPU and memory profiling check..."
    TMP_CPU="$(mktemp -t emlgo-cpu.XXXXXX.prof)"
    TMP_MEM="$(mktemp -t emlgo-mem.XXXXXX.prof)"
    trap 'rm -f "$TMP_COV" "$TMP_CPU" "$TMP_MEM"' EXIT
    go test -run='^$' -bench='^Benchmark(ExpSIMD|LogSIMD)$' -benchtime=50ms -benchmem \
        -cpuprofile="$TMP_CPU" -memprofile="$TMP_MEM" ./internal/eml >/dev/null
    go tool pprof -top -nodecount=5 "$TMP_CPU" >/dev/null
    go tool pprof -top -alloc_space -nodecount=5 "$TMP_MEM" >/dev/null
fi

echo "PASS: Verification checks succeeded (security, oracle gate, linting clean, statement coverage ${TOTAL_PCT}% >= ${MIN_COVERAGE}%)."
exit 0
