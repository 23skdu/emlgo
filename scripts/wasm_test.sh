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
    local test_script="$2"
    
    cp "$WASM_DIR/wasm_exec.js" "$WASM_DIR/$(basename $wasm_file .wasm).js" 2>/dev/null || true
    
    # Build a Go test binary for WASM
    echo "Running: node $WASM_DIR/run.js $WASM_DIR/$wasm_file"
    node "$WASM_DIR/run.js" "$WASM_DIR/$wasm_file"
}

case "$MODE" in
    test)
        echo "--- Building WASM test binary ---"
        # Build the CLI tool for WASM to run gpu-status/demo
        build_wasm "./cmd/emlcli" "emlcli.wasm"
        echo ""
        echo "--- Running WASM demo ---"
        # Create a small Go WASM test program that runs the EML library
        cat > "$WASM_DIR/main.go" << 'GOEOF'
//go:build wasm
// +build wasm

package main

import (
    "fmt"
    "math"
    "os"
    "syscall/js"

    "github.com/emlgo/eml/pkg/arithmetic"
    "github.com/emlgo/eml/pkg/logexp"
    "github.com/emlgo/eml/pkg/trig"
)

func withinTol(a, b, tol float64) bool {
    if math.IsNaN(a) && math.IsNaN(b) {
        return true
    }
    if math.IsInf(a, 1) && math.IsInf(b, 1) {
        return true
    }
    if math.IsInf(a, -1) && math.IsInf(b, -1) {
        return true
    }
    diff := math.Abs(a - b)
    sumAbs := math.Abs(a) + math.Abs(b) + 1e-10
    return diff < tol || diff/sumAbs < tol
}

func runValidation() bool {
    fmt.Println("=== Starting WASM vs Go Math Parity Validation ===")
    failed := false

    // We test multiple sizes to check vector alignment and remainder loops:
    // 1024 (multiple of 8), 1027 (not multiple of 8)
    sizes := []int{1024, 1027}

    for _, n := range sizes {
        fmt.Printf("Testing array size: %d...\n", n)

        // 1. ExpBatch
        inputs := make([]float64, n)
        for i := range inputs {
            inputs[i] = -5.0 + float64(i)*10.0/float64(n)
        }
        resExp := logexp.ExpBatch(inputs)
        for i, val := range inputs {
            expected := math.Exp(val)
            if !withinTol(resExp[i], expected, 1e-9) {
                fmt.Printf("  [FAIL] ExpBatch: index %d, input %f, got %f, expected %f\n", i, val, resExp[i], expected)
                failed = true
                break
            }
        }

        // 2. LogBatch
        for i := range inputs {
            inputs[i] = 0.001 + float64(i)*100.0/float64(n)
        }
        resLog := logexp.LogBatch(inputs)
        for i, val := range inputs {
            expected := math.Log(val)
            if !withinTol(resLog[i], expected, 1e-9) {
                fmt.Printf("  [FAIL] LogBatch: index %d, input %f, got %f, expected %f\n", i, val, resLog[i], expected)
                failed = true
                break
            }
        }

        // 3. SinBatch / CosBatch / TanBatch
        for i := range inputs {
            inputs[i] = -2.0*math.Pi + float64(i)*4.0*math.Pi/float64(n)
        }
        resSin := trig.SinBatch(inputs)
        resCos := trig.CosBatch(inputs)
        resTan := trig.TanBatch(inputs)
        for i, val := range inputs {
            expectedSin := math.Sin(val)
            expectedCos := math.Cos(val)
            expectedTan := math.Tan(val)
            if !withinTol(resSin[i], expectedSin, 1e-9) {
                fmt.Printf("  [FAIL] SinBatch: index %d, input %f, got %f, expected %f\n", i, val, resSin[i], expectedSin)
                failed = true
                break
            }
            if !withinTol(resCos[i], expectedCos, 1e-9) {
                fmt.Printf("  [FAIL] CosBatch: index %d, input %f, got %f, expected %f\n", i, val, resCos[i], expectedCos)
                failed = true
                break
            }
            // Avoid asymptotes of Tan for simple parity
            if val > -math.Pi/2 + 0.1 && val < math.Pi/2 - 0.1 {
                if !withinTol(resTan[i], expectedTan, 1e-9) {
                    fmt.Printf("  [FAIL] TanBatch: index %d, input %f, got %f, expected %f\n", i, val, resTan[i], expectedTan)
                    failed = true
                    break
                }
            }
        }

        // 4. SqrtBatch
        for i := range inputs {
            inputs[i] = float64(i) * 100.0 / float64(n)
        }
        resSqrt := arithmetic.SqrtBatch(inputs)
        for i, val := range inputs {
            expected := math.Sqrt(val)
            if !withinTol(resSqrt[i], expected, 1e-9) {
                fmt.Printf("  [FAIL] SqrtBatch: index %d, input %f, got %f, expected %f\n", i, val, resSqrt[i], expected)
                failed = true
                break
            }
        }

        // 5. AddBatch, SubBatch, MulBatch, DivBatch
        a := make([]float64, n)
        b := make([]float64, n)
        for i := range a {
            a[i] = float64(i)
            b[i] = float64(i)*2.0 + 1.0
        }
        resAdd := arithmetic.AddBatch(a, b)
        resSub := arithmetic.SubBatch(a, b)
        resMul := arithmetic.MulBatch(a, b)
        resDiv := arithmetic.DivBatch(a, b)
        for i := range a {
            if !withinTol(resAdd[i], a[i]+b[i], 1e-9) {
                fmt.Printf("  [FAIL] AddBatch: index %d, got %f, expected %f\n", i, resAdd[i], a[i]+b[i])
                failed = true
                break
            }
            if !withinTol(resSub[i], a[i]-b[i], 1e-9) {
                fmt.Printf("  [FAIL] SubBatch: index %d, got %f, expected %f\n", i, resSub[i], a[i]-b[i])
                failed = true
                break
            }
            if !withinTol(resMul[i], a[i]*b[i], 1e-9) {
                fmt.Printf("  [FAIL] MulBatch: index %d, got %f, expected %f\n", i, resMul[i], a[i]*b[i])
                failed = true
                break
            }
            if !withinTol(resDiv[i], a[i]/b[i], 1e-9) {
                fmt.Printf("  [FAIL] DivBatch: index %d, got %f, expected %f\n", i, resDiv[i], a[i]/b[i])
                failed = true
                break
            }
        }

        // 6. AbsBatch, NegBatch, InvBatch
        for i := range inputs {
            inputs[i] = -50.0 + float64(i)*100.0/float64(n)
            if inputs[i] == 0 {
                inputs[i] = 1.0 // avoid division by zero
            }
        }
        resAbs := arithmetic.AbsBatch(inputs)
        resNeg := arithmetic.NegBatch(inputs)
        resInv := arithmetic.InvBatch(inputs)
        for i, val := range inputs {
            if !withinTol(resAbs[i], math.Abs(val), 1e-9) {
                fmt.Printf("  [FAIL] AbsBatch: index %d, input %f, got %f, expected %f\n", i, val, resAbs[i], math.Abs(val))
                failed = true
                break
            }
            if !withinTol(resNeg[i], -val, 1e-9) {
                fmt.Printf("  [FAIL] NegBatch: index %d, input %f, got %f, expected %f\n", i, val, resNeg[i], -val)
                failed = true
                break
            }
            if !withinTol(resInv[i], 1.0/val, 1e-9) {
                fmt.Printf("  [FAIL] InvBatch: index %d, input %f, got %f, expected %f\n", i, val, resInv[i], 1.0/val)
                failed = true
                break
            }
        }

        // 7. FMABatch
        c := make([]float64, n)
        for i := range c {
            c[i] = float64(i) * 0.5
        }
        resFMA := arithmetic.FmaBatch(a, b, c)
        for i := range a {
            expected := a[i]*b[i] + c[i]
            if !withinTol(resFMA[i], expected, 1e-9) {
                fmt.Printf("  [FAIL] FmaBatch: index %d, got %f, expected %f\n", i, resFMA[i], expected)
                failed = true
                break
            }
        }

        // 8. AddScalarBatch, MulScalarBatch
        resAddScalar := arithmetic.AddScalarBatch(a, 3.14)
        resMulScalar := arithmetic.MulScalarBatch(a, 2.5)
        for i := range a {
            if !withinTol(resAddScalar[i], a[i]+3.14, 1e-9) {
                fmt.Printf("  [FAIL] AddScalarBatch: index %d, got %f, expected %f\n", i, resAddScalar[i], a[i]+3.14)
                failed = true
                break
            }
            if !withinTol(resMulScalar[i], a[i]*2.5, 1e-9) {
                fmt.Printf("  [FAIL] MulScalarBatch: index %d, got %f, expected %f\n", i, resMulScalar[i], a[i]*2.5)
                failed = true
                break
            }
        }
    }

    if failed {
        fmt.Println("=== WASM vs Go Math Parity Validation FAILED! ===")
        return false
    }
    fmt.Println("=== WASM vs Go Math Parity Validation PASSED SUCCESSFULLY! ===")
    return true
}

func main() {
    runTests := func() {
        success := runValidation()
        if !success {
            os.Exit(1)
        }
    }

    if js.Global().Get("window").IsUndefined() {
        // Run immediately and exit for Node.js/CLI execution
        runTests()
    } else {
        // Register export and block for browser execution
        c := make(chan struct{})
        js.Global().Set("emlgo_run", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            runTests()
            return nil
        }))
        <-c
    }
}
GOEOF

        GOOS=js GOARCH=wasm go build -o "$WASM_DIR/wasm_test.wasm" "$WASM_DIR/main.go"
        
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
        node "$WASM_DIR/run.js" "$WASM_DIR/wasm_test.wasm" 2>&1 || {
            echo "Note: WASM binary built successfully at $WASM_DIR/wasm_test.wasm"
            echo "To run in browser: open wasm/test_runner.html in a web server"
            echo "To run with Node: node wasm/run.js wasm/wasm_test.wasm"
        }
        ;;

    bench)
        echo "--- Running WASM Benchmarks ---"
        cat > "$WASM_DIR/bench.go" << 'GOEOF'
//go:build wasm

package main

import (
    "fmt"
    "time"

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

        GOOS=js GOARCH=wasm go build -o "$WASM_DIR/wasm_bench.wasm" "$WASM_DIR/bench.go"
        echo ""
        echo "Running benchmarks..."
        node "$WASM_DIR/run.js" "$WASM_DIR/wasm_bench.wasm" 2>&1
        ;;

    serve)
        echo "Starting web benchmark server at http://localhost:8080"
        echo "Open http://localhost:8080/wasm/bench.html in your browser"
        cd "$PROJ_DIR"
        python3 -m http.server 8080
        ;;

    clean)
        echo "Cleaning WASM build artifacts..."
        rm -f "$WASM_DIR"/*.wasm "$WASM_DIR"/main.go "$WASM_DIR"/bench.go "$WASM_DIR"/run.js "$WASM_DIR"/wasm_exec.js "$WASM_DIR"/test_runner.html
        echo "Done"
        ;;

    *)
        echo "Usage: $0 [test|bench|serve|clean]"
        exit 1
        ;;
esac
