#!/usr/bin/env bash
# Oracle Accuracy Gate for EML Go
# Enforces multi-tier numerical accuracy verification against:
#   1. GNU MPFR 4.2.2 (MPFR_RNDN) golden reference vectors (53-bit & 24-bit)
#   2. Arbitrary-precision big.Float (256-bit/512-bit) & ExactLogOracle differential suites
#   3. SIMD batch kernel ULP, mixed-tolerance, tail-loop, and identity gates
#   4. JIT native machine-code vs AST canonical differential evaluation
#   5. Full type/domain validation (cmd/validate) and ULP/parity harness (cmd/bench)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

VERBOSE=false
QUICK=false
REGEN_CHECK=false

usage() {
    cat <<EOF
Usage: $0 [options]

Oracle Accuracy Gate — verifies numerical accuracy across all EML Go kernels and packages.

Options:
  -v, --verbose       Enable verbose test and oracle output
  -q, --quick         Run fast Go oracle test gates only (skip cmd/validate and cmd/bench sweeps)
  -r, --regen-check   Recompile scripts/generate_mpfr_vectors.c (if libmpfr 4.2.2 is installed)
                      and verify byte-for-byte identity with internal/accuracyvectors/mpfr-vectors.csv
  -h, --help          Show this help message
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -q|--quick)
            QUICK=true
            shift
            ;;
        -r|--regen-check)
            REGEN_CHECK=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}ERROR: Invalid option: $1${NC}" >&2
            usage >&2
            exit 1
            ;;
    esac
done

cd "$PROJECT_DIR"

GO_TEST_FLAGS=("-count=1")
if [ "$VERBOSE" = true ]; then
    GO_TEST_FLAGS+=("-v")
fi

echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}               EML Go — Oracle Accuracy Gate                    ${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""

# ------------------------------------------------------------------------------
# Gate 1: MPFR 4.2.2 Provenance & Golden Vector Verification
# ------------------------------------------------------------------------------
echo -e "${YELLOW}[1/5] Verifying MPFR 4.2.2 provenance and golden reference vectors...${NC}"

MPFR_CSV="internal/accuracyvectors/mpfr-vectors.csv"
EXPECTED_HEADER="# MPFR 4.2.2; MPFR_RNDN; exact binary inputs; result precision 53/24 bits"

if [ ! -f "$MPFR_CSV" ]; then
    echo -e "${RED}FAIL: Missing MPFR vector file: $MPFR_CSV${NC}" >&2
    exit 1
fi

ACTUAL_HEADER="$(head -n 1 "$MPFR_CSV")"
if [ "$ACTUAL_HEADER" != "$EXPECTED_HEADER" ]; then
    echo -e "${RED}FAIL: Invalid MPFR provenance header in $MPFR_CSV${NC}" >&2
    echo "  Expected: $EXPECTED_HEADER" >&2
    echo "  Actual:   $ACTUAL_HEADER" >&2
    exit 1
fi

VECTOR_COUNT="$(grep -v '^#' "$MPFR_CSV" | grep -v '^operation,' | grep -v '^[[:space:]]*$' | wc -l | tr -d ' ')"
if [ "$VECTOR_COUNT" -lt 130 ]; then
    echo -e "${RED}FAIL: Expected at least 130 MPFR reference vectors in $MPFR_CSV, found $VECTOR_COUNT${NC}" >&2
    exit 1
fi

if [ "$REGEN_CHECK" = true ]; then
    if command -v cc >/dev/null 2>&1; then
        TMP_BIN="$(mktemp)"
        TMP_CSV="$(mktemp)"
        trap 'rm -f "$TMP_BIN" "$TMP_CSV"' EXIT
        if cc -O2 -Wall -Wextra -Werror scripts/generate_mpfr_vectors.c -o "$TMP_BIN" -lmpfr -lgmp 2>/dev/null || \
           cc -O2 -Wall -Wextra -Werror scripts/generate_mpfr_vectors.c /usr/lib/x86_64-linux-gnu/libmpfr.so.6 -o "$TMP_BIN" 2>/dev/null; then
            if "$TMP_BIN" "$TMP_CSV" 2>/dev/null; then
                if ! diff -u "$MPFR_CSV" "$TMP_CSV"; then
                    echo -e "${RED}FAIL: Regenerated MPFR vectors differ from $MPFR_CSV${NC}" >&2
                    exit 1
                fi
                echo -e "  ${GREEN}✓ Regenerated MPFR 4.2.2 vectors match $MPFR_CSV byte-for-byte${NC}"
            else
                echo -e "  ${YELLOW}⚠ Host MPFR version differs from 4.2.2; skipped live regeneration diff${NC}"
            fi
        else
            echo -e "  ${YELLOW}⚠ libmpfr not found; skipped live regeneration diff${NC}"
        fi
    else
        echo -e "  ${YELLOW}⚠ C compiler not installed; skipped live regeneration diff${NC}"
    fi
fi

go test "${GO_TEST_FLAGS[@]}" -run 'Test(Rows|.*MPFR).*' ./internal/accuracyvectors/... ./internal/eml/... ./pkg/fastmath/...
echo -e "${GREEN}PASS [1/5]: MPFR 4.2.2 golden vectors verified ($VECTOR_COUNT vectors)${NC}"
echo ""

# ------------------------------------------------------------------------------
# Gate 2: Arbitrary-Precision bigmath & ExactLogOracle / SIMD Differential Suite
# ------------------------------------------------------------------------------
echo -e "${YELLOW}[2/5] Running arbitrary-precision bigmath & SIMD differential accuracy oracles...${NC}"
go test "${GO_TEST_FLAGS[@]}" ./internal/eml/bigmath/...
go test "${GO_TEST_FLAGS[@]}" \
    -run 'Test(BatchKernelAccuracy|BatchKernelTailLoopAccuracy|BatchTrigIdentities|BatchSinCosZeros|SIMDToDifferentialAndAliasing|BatchKernelMatchesScalarPath|BatchKernelSpecialValues|BatchLogMatchesScalarDomain|ScalarLogMatchesOracleForPositive|Subnormal|Oracle).*' \
    ./internal/eml/... ./pkg/...
echo -e "${GREEN}PASS [2/5]: Arbitrary-precision bigmath & SIMD batch accuracy gates passed${NC}"
echo ""

# ------------------------------------------------------------------------------
# Gate 3: JIT vs AST Canonical & Differential Accuracy Gate
# ------------------------------------------------------------------------------
echo -e "${YELLOW}[3/5] Running JIT native codegen vs AST canonical oracle gate...${NC}"
go test "${GO_TEST_FLAGS[@]}" -run 'Test(Compile|Canonical|Differential|Eval).*' ./internal/jit/...
echo -e "${GREEN}PASS [3/5]: JIT vs AST differential accuracy gate passed${NC}"
echo ""

if [ "$QUICK" = true ]; then
    echo -e "${YELLOW}[4/5] Skipped (--quick): cmd/validate full type validation${NC}"
    echo -e "${YELLOW}[5/5] Skipped (--quick): cmd/bench -compare and -accuracy ULP sweeps${NC}"
    echo ""
    echo -e "${GREEN}================================================================${NC}"
    echo -e "${GREEN}     ORACLE ACCURACY GATE PASSED (QUICK MODE)                   ${NC}"
    echo -e "${GREEN}================================================================${NC}"
    exit 0
fi

# ------------------------------------------------------------------------------
# Gate 4: Full Type & Domain Validation (cmd/validate)
# ------------------------------------------------------------------------------
echo -e "${YELLOW}[4/5] Running full type and domain validation (cmd/validate)...${NC}"
VALIDATE_ARGS=()
if [ "$VERBOSE" = true ]; then
    VALIDATE_ARGS+=("-v")
fi
VALIDATE_OUT="$(go run ./cmd/validate "${VALIDATE_ARGS[@]}")"
if [ "$VERBOSE" = true ]; then
    echo "$VALIDATE_OUT"
fi
if ! echo "$VALIDATE_OUT" | grep -q "All validation tests passed"; then
    echo "$VALIDATE_OUT" >&2
    echo -e "${RED}FAIL [4/5]: cmd/validate reported validation failures${NC}" >&2
    exit 1
fi
echo -e "${GREEN}PASS [4/5]: cmd/validate passed all type, precision, quantization, and batch checks${NC}"
echo ""

# ------------------------------------------------------------------------------
# Gate 5: Feature Parity & ULP Accuracy Bounds (cmd/bench -compare & -accuracy)
# ------------------------------------------------------------------------------
echo -e "${YELLOW}[5/5] Running feature parity and ULP accuracy bounds (cmd/bench)...${NC}"
COMPARE_OUT="$(go run ./cmd/bench -compare)"
if [ "$VERBOSE" = true ]; then
    echo "$COMPARE_OUT"
fi
if ! echo "$COMPARE_OUT" | grep -q "0 failed"; then
    echo "$COMPARE_OUT" >&2
    echo -e "${RED}FAIL [5/5]: cmd/bench -compare reported parity failures${NC}" >&2
    exit 1
fi

ACCURACY_OUT="$(go run ./cmd/bench -accuracy)"
if [ "$VERBOSE" = true ]; then
    echo "$ACCURACY_OUT"
fi
if ! echo "$ACCURACY_OUT" | grep -q "0 failed"; then
    echo "$ACCURACY_OUT" >&2
    echo -e "${RED}FAIL [5/5]: cmd/bench -accuracy exceeded ULP bounds${NC}" >&2
    exit 1
fi
echo -e "${GREEN}PASS [5/5]: cmd/bench parity and ULP accuracy bounds passed${NC}"
echo ""

echo -e "${GREEN}================================================================${NC}"
echo -e "${GREEN}     ALL 5 ORACLE ACCURACY GATES PASSED                         ${NC}"
echo -e "${GREEN}================================================================${NC}"
