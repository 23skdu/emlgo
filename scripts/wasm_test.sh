#!/usr/bin/env bash
# Build and run WASM SIMD tests with Node.js.
# Requires: Go 1.21+, Node.js 16+ (for WASM SIMD support)
#
# Usage:
#   ./scripts/wasm_test.sh            # build + test
#   ./scripts/wasm_test.sh --bench    # build + benchmark
#   ./scripts/wasm_test.sh --serve    # start web benchmark server

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJ_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
WASM_DIR="$PROJ_DIR/wasm"
WASM_EXEC_JS="$(go env GOROOT)/lib/wasm/wasm_exec.js"
[ ! -f "$WASM_EXEC_JS" ] && WASM_EXEC_JS="$(go env GOROOT)/misc/wasm/wasm_exec.js"

MODE="${1:-test}"
MODE="${MODE#--}"

echo "=== EML WASM SIMD Test Harness ==="
echo "Go version: $(go version)"
echo "Node version: $(node --version 2>/dev/null || echo 'not found')"
echo ""

# Check prerequisites
if ! command -v node &>/dev/null; then
    echo "ERROR: Node.js is required. Install Node.js 16+ with WASM SIMD support."
    exit 1
fi

mkdir -p "$WASM_DIR"

# Copy wasm_exec.js helper from GOROOT if not present
if [ ! -f "$WASM_DIR/wasm_exec.js" ]; then
    if [ -f "$WASM_EXEC_JS" ]; then
        cp "$WASM_EXEC_JS" "$WASM_DIR/wasm_exec.js"
        echo "Copied wasm_exec.js from Go toolchain"
    else
        echo "ERROR: wasm_exec.js not found at $WASM_EXEC_JS"
        exit 1
    fi
fi

# Create Node.js runner script if not present
cat > "$WASM_DIR/run.js" << 'JSEOF'
const fs = require('fs');
const path = require('path');

// Load Go's wasm_exec polyfills
require(path.join(__dirname, 'wasm_exec.js'));

if (process.argv.length < 3) {
    console.error('Usage: node run.js <wasm_file>');
    process.exit(1);
}

const wasmFile = process.argv[2];
const wasmBuffer = fs.readFileSync(wasmFile);

const go = new Go();
WebAssembly.instantiate(wasmBuffer, go.importObject).then((result) => {
    return go.run(result.instance);
}).catch((err) => {
    console.error(err);
    process.exit(1);
});
JSEOF

build_wasm() {
    local pkg="$1"
    local out="$2"
    echo "Building $pkg -> $out"
    GOOS=js GOARCH=wasm go build -o "$WASM_DIR/$out" "$pkg"
}

run_wasm_test() {
    local wasm_file="$1"
    local test_script="${2:-}"
    
    echo "Running: node $WASM_DIR/run.js $WASM_DIR/$wasm_file"
    node "$WASM_DIR/run.js" "$WASM_DIR/$wasm_file"
}

case "$MODE" in
    test)
        echo "--- Building WASM binary from cmd/wasmbench ---"
        build_wasm "./cmd/wasmbench" "emlgo.wasm"
        echo ""
        echo "--- Running WASM validation suite with Node.js ---"
        run_wasm_test "emlgo.wasm"
        ;;

    bench)
        echo "--- Building WASM binary from cmd/wasmbench ---"
        build_wasm "./cmd/wasmbench" "emlgo.wasm"
        echo ""
        echo "--- Running WASM Benchmarks via Node.js ---"
        node -e '
const fs = require("fs");
const path = require("path");
globalThis.window = globalThis;
require(path.join("'"$WASM_DIR"'", "wasm_exec.js"));
const go = new Go();
const wasmBuffer = fs.readFileSync(path.join("'"$WASM_DIR"'", "emlgo.wasm"));
WebAssembly.instantiate(wasmBuffer, go.importObject).then((result) => {
    go.run(result.instance);
    const sizes = [1024, 4096, 16384, 65536];
    console.log("=== WASM Batch Performance ===");
    console.log("Size       Exp (ns/elem)  Sqrt (ns/elem) Add (ns/elem)");
    for (const n of sizes) {
        const a = new Float64Array(n);
        const b = new Float64Array(n);
        for (let i = 0; i < n; i++) { a[i] = (i % 100) / 100.0; b[i] = a[i] * 2; }
        
        let t0 = performance.now();
        globalThis.emlgoExpBatch(a);
        let expNs = ((performance.now() - t0) * 1e6) / n;

        t0 = performance.now();
        globalThis.emlgoSqrtBatch(a);
        let sqrtNs = ((performance.now() - t0) * 1e6) / n;

        t0 = performance.now();
        globalThis.emlgoAddBatch(a, b);
        let addNs = ((performance.now() - t0) * 1e6) / n;

        console.log(`${n.toString().padEnd(10)} ${expNs.toFixed(2).padStart(14)} ${sqrtNs.toFixed(2).padStart(14)} ${addNs.toFixed(2).padStart(14)}`);
    }
    process.exit(0);
}).catch(err => { console.error(err); process.exit(1); });
'
        ;;

    serve)
        echo "Starting web benchmark server at http://localhost:8080"
        echo "Open http://localhost:8080/wasm/bench.html in your browser"
        cd "$PROJ_DIR"
        python3 -m http.server 8080
        ;;

    clean)
        echo "Cleaning WASM build artifacts..."
        rm -f "$WASM_DIR"/*.wasm "$WASM_DIR"/run.js "$WASM_DIR"/wasm_exec.js
        echo "Done"
        ;;

    *)
        echo "Usage: $0 [test|bench|serve|clean]"
        exit 1
        ;;
esac
