//go:build js && wasm

package main

import (
	"os"
	"syscall/js"
	"unsafe"

	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

func jsToFloat64Slice(v js.Value) []float64 {
	n := v.Get("length").Int()
	if n <= 0 {
		return nil
	}
	res := make([]float64, n)
	u8 := js.Global().Get("Uint8Array").New(v.Get("buffer"), v.Get("byteOffset"), v.Get("byteLength"))
	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(&res[0])), n*8)
	js.CopyBytesToGo(byteSlice, u8)
	return res
}

func float64SliceToJS(slice []float64) js.Value {
	n := len(slice)
	if n == 0 {
		return js.Global().Get("Float64Array").New(0)
	}
	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), n*8)
	u8 := js.Global().Get("Uint8Array").New(n * 8)
	js.CopyBytesToJS(u8, byteSlice)
	return js.Global().Get("Float64Array").New(u8.Get("buffer"), u8.Get("byteOffset"), n)
}

func registerExports() {
	js.Global().Set("emlgoExpBatch", js.FuncOf(func(this js.Value, args []js.Value) any {
		in := jsToFloat64Slice(args[0])
		out := logexp.ExpBatch(in)
		return float64SliceToJS(out)
	}))

	js.Global().Set("emlgoLogBatch", js.FuncOf(func(this js.Value, args []js.Value) any {
		in := jsToFloat64Slice(args[0])
		out := logexp.LogBatch(in)
		return float64SliceToJS(out)
	}))

	js.Global().Set("emlgoSinBatch", js.FuncOf(func(this js.Value, args []js.Value) any {
		in := jsToFloat64Slice(args[0])
		out := trig.SinBatch(in)
		return float64SliceToJS(out)
	}))

	js.Global().Set("emlgoCosBatch", js.FuncOf(func(this js.Value, args []js.Value) any {
		in := jsToFloat64Slice(args[0])
		out := trig.CosBatch(in)
		return float64SliceToJS(out)
	}))

	js.Global().Set("emlgoSqrtBatch", js.FuncOf(func(this js.Value, args []js.Value) any {
		in := jsToFloat64Slice(args[0])
		out := arithmetic.SqrtBatch(in)
		return float64SliceToJS(out)
	}))

	js.Global().Set("emlgoAddBatch", js.FuncOf(func(this js.Value, args []js.Value) any {
		a := jsToFloat64Slice(args[0])
		b := jsToFloat64Slice(args[1])
		out := arithmetic.AddBatch(a, b)
		return float64SliceToJS(out)
	}))

	js.Global().Set("emlgo_run", js.FuncOf(func(this js.Value, args []js.Value) any {
		ok := RunValidation()
		return ok
	}))
}

func main() {
	registerExports()

	if js.Global().Get("window").IsUndefined() {
		// Running in Node.js / CLI environment: run validation and exit
		if !RunValidation() {
			os.Exit(1)
		}
	} else {
		// Browser environment: keep event loop alive for interactive benchmarks
		c := make(chan struct{})
		<-c
	}
}
