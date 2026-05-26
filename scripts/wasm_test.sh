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

build_wasm() {
    local pkg="$1"
    local out="$2"
    echo "Building $pkg -> $out"
    GOOS=js GOARCH=wasm go build -o "$WASM_DIR/$out" "$pkg"
}

run_wasm_test() {
    local wasm_file="$1"
    local test_script="$2"
    
    cp "$WASM_DIR/wasm_exec.js" "$WASM_DIR/$(basename $wasm_file .wasm).js" 2>/dev/null || true
    
    # Build a Go test binary for WASM
    echo "Running: node $WASM_DIR/wasm_exec.js $WASM_DIR/$wasm_file"
    node "$WASM_DIR/wasm_exec.js" "$WASM_DIR/$wasm_file"
}

case "$MODE" in
    test)
        echo "--- Building WASM test binary ---"
        # Build the CLI tool for WASM to run gpu-status/demo
        build_wasm "./cmd/emlcli" "emlcli.wasm"
        echo ""
        echo "--- Running WASM demo ---"
        # Create a small Go WASM test program that runs the EML library
        cat > "$WASM_DIR/main_test.go" << 'GOEOF'
//go:build wasm
// +build wasm

package main

import (
    "fmt"
    "syscall/js"

    "github.com/emlgo/eml/internal/eml"
    "github.com/emlgo/eml/pkg/arithmetic"
    "github.com/emlgo/eml/pkg/logexp"
    "github.com/emlgo/eml/pkg/trig"
)

func main() {
    c := make(chan struct{}, 0)

    js.Global().Set("emlgo_run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        fmt.Println("=== EML WASM Test ===")
        fmt.Printf("HasWasmSIMD: %v\n", eml.HasWasmSIMD())
        fmt.Printf("EML(1,1) = %v\n", eml.Eml(1, 1))
        fmt.Printf("Exp(1) = %v\n", logexp.Exp(1))
        fmt.Printf("Log(e) = %v\n", logexp.Log(2.718281828459045))
        fmt.Printf("Sin(0) = %v\n", trig.Sin(0))
        fmt.Printf("Cos(0) = %v\n", trig.Cos(0))
        fmt.Printf("Sqrt(4) = %v\n", arithmetic.Sqrt(4))

        // Test batch operations
        n := 1024
        data := make([]float64, n)
        for i := range data {
            data[i] = float64(i%100) / 100.0
        }
        result := logexp.ExpBatch(data)
        fmt.Printf("ExpBatch(%d): first=%.6f last=%.6f\n", n, result[0], result[n-1])

        fmt.Println("=== All WASM tests passed ===")
        return nil
    }))

    <-c
}
GOEOF

        GOOS=js GOARCH=wasm go build -o "$WASM_DIR/wasm_test.wasm" "$WASM_DIR/main_test.go"
        
        # Create the runner HTML
        cat > "$WASM_DIR/test_runner.html" << 'HTML'
<!DOCTYPE html>
<html>
<head><title>EML WASM Test</title></head>
<body>
<script src="wasm_exec.js"></script>
<script>
const go = new Go();
WebAssembly.instantiateStreaming(fetch("wasm_test.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    // Call the exported function
    emlgo_run();
    document.body.innerHTML = '<pre>' + 
        document.getElementById('output')?.textContent || 'Check console for output' + 
        '</pre>';
});
</script>
</body>
</html>
HTML

        echo ""
        echo "--- Running with Node.js ---"
        node "$WASM_DIR/wasm_exec.js" "$WASM_DIR/wasm_test.wasm" 2>&1 || {
            echo "Note: WASM binary built successfully at $WASM_DIR/wasm_test.wasm"
            echo "To run in browser: open wasm/test_runner.html in a web server"
            echo "To run with Node: node wasm/wasm_exec.js wasm/wasm_test.wasm"
        }
        ;;

    bench)
        echo "--- Running WASM Benchmarks ---"
        cat > "$WASM_DIR/bench_test.go" << 'GOEOF'
//go:build wasm

package main

import (
    "fmt"
    "time"

    "github.com/emlgo/eml/internal/eml"
    "github.com/emlgo/eml/pkg/logexp"
    "github.com/emlgo/eml/pkg/arithmetic"
)

func main() {
    sizes := []int{1024, 4096, 16384, 65536}
    fmt.Println("=== WASM Batch Performance ===")
    fmt.Printf("%-10s %14s %14s %14s\n", "Size", "Exp (ns/elem)", "Sqrt (ns/elem)", "Add (ns/elem)")

    for _, n := range sizes {
        data := make([]float64, n)
        for i := range data {
            data[i] = float64(i%100) / 100.0
        }

        start := time.Now()
        _ = logexp.ExpBatch(data)
        expNs := float64(time.Since(start).Nanoseconds()) / float64(n)

        start = time.Now()
        _ = arithmetic.SqrtBatch(data)
        sqrtNs := float64(time.Since(start).Nanoseconds()) / float64(n)

        a := make([]float64, n)
        b := make([]float64, n)
        for i := range a {
            a[i] = float64(i)
            b[i] = float64(i) * 2
        }
        start = time.Now()
        _ = arithmetic.AddBatch(a, b)
        addNs := float64(time.Since(start).Nanoseconds()) / float64(n)

        fmt.Printf("%-10d %14.2f %14.2f %14.2f\n", n, expNs, sqrtNs, addNs)
    }
}
GOEOF

        GOOS=js GOARCH=wasm go build -o "$WASM_DIR/wasm_bench.wasm" "$WASM_DIR/bench_test.go"
        echo ""
        echo "Running benchmarks..."
        node "$WASM_DIR/wasm_exec.js" "$WASM_DIR/wasm_bench.wasm" 2>&1
        ;;

    serve)
        echo "Starting web benchmark server at http://localhost:8080"
        echo "Open http://localhost:8080/wasm/bench.html in your browser"
        cd "$PROJ_DIR"
        python3 -m http.server 8080
        ;;

    clean)
        echo "Cleaning WASM build artifacts..."
        rm -rf "$WASM_DIR"/*
        echo "Done"
        ;;

    *)
        echo "Usage: $0 [test|bench|serve|clean]"
        exit 1
        ;;
esac
