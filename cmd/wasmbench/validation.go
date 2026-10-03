package main

import (
	"fmt"
	"math"

	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/logexp"
	"github.com/emlgo/eml/pkg/trig"
)

var (
	validateInjectFail = false
	validateSizes      = []int{1024, 1027}
	sliceChecker       = checkSlice
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

func checkSlice(name string, got, expected []float64, tol float64) bool {
	for i := range got {
		if !withinTol(got[i], expected[i], tol) {
			fmt.Printf("  [FAIL] %s: index %d, got %f, expected %f\n", name, i, got[i], expected[i])
			return false
		}
	}
	return true
}

// RunValidation tests parity across sizes for batch math functions.
func RunValidation() bool {
	fmt.Println("=== Starting WASM vs Go Math Parity Validation ===")
	failed := false
	if validateInjectFail {
		failed = true
	}

	for _, n := range validateSizes {
		fmt.Printf("Testing array size: %d...\n", n)

		// 1. ExpBatch
		inputs := make([]float64, n)
		expExp := make([]float64, n)
		for i := range inputs {
			inputs[i] = -5.0 + float64(i)*10.0/float64(n)
			expExp[i] = math.Exp(inputs[i])
		}
		resExp := logexp.ExpBatch(inputs)
		if !sliceChecker("ExpBatch", resExp, expExp, 1e-9) {
			failed = true
		}

		// 2. LogBatch
		expLog := make([]float64, n)
		for i := range inputs {
			inputs[i] = 0.001 + float64(i)*100.0/float64(n)
			expLog[i] = math.Log(inputs[i])
		}
		resLog := logexp.LogBatch(inputs)
		if !sliceChecker("LogBatch", resLog, expLog, 1e-9) {
			failed = true
		}

		// 3. SinBatch / CosBatch / TanBatch
		expSin := make([]float64, n)
		expCos := make([]float64, n)
		for i := range inputs {
			inputs[i] = -2.0*math.Pi + float64(i)*4.0*math.Pi/float64(n)
			expSin[i] = math.Sin(inputs[i])
			expCos[i] = math.Cos(inputs[i])
		}
		resSin := trig.SinBatch(inputs)
		resCos := trig.CosBatch(inputs)
		if !sliceChecker("SinBatch", resSin, expSin, 1e-9) {
			failed = true
		}
		if !sliceChecker("CosBatch", resCos, expCos, 1e-9) {
			failed = true
		}

		// 4. SqrtBatch
		expSqrt := make([]float64, n)
		for i := range inputs {
			inputs[i] = float64(i) * 100.0 / float64(n)
			expSqrt[i] = math.Sqrt(inputs[i])
		}
		resSqrt := arithmetic.SqrtBatch(inputs)
		if !sliceChecker("SqrtBatch", resSqrt, expSqrt, 1e-9) {
			failed = true
		}

		// 5. AddBatch, SubBatch, MulBatch, DivBatch
		a := make([]float64, n)
		b := make([]float64, n)
		expAdd := make([]float64, n)
		expSub := make([]float64, n)
		expMul := make([]float64, n)
		expDiv := make([]float64, n)
		for i := range a {
			a[i] = float64(i)
			b[i] = float64(i)*2.0 + 1.0
			expAdd[i] = a[i] + b[i]
			expSub[i] = a[i] - b[i]
			expMul[i] = a[i] * b[i]
			expDiv[i] = a[i] / b[i]
		}
		if !sliceChecker("AddBatch", arithmetic.AddBatch(a, b), expAdd, 1e-9) {
			failed = true
		}
		if !sliceChecker("SubBatch", arithmetic.SubBatch(a, b), expSub, 1e-9) {
			failed = true
		}
		if !sliceChecker("MulBatch", arithmetic.MulBatch(a, b), expMul, 1e-9) {
			failed = true
		}
		if !sliceChecker("DivBatch", arithmetic.DivBatch(a, b), expDiv, 1e-9) {
			failed = true
		}

		// 6. AbsBatch, NegBatch, InvBatch
		expAbs := make([]float64, n)
		expNeg := make([]float64, n)
		expInv := make([]float64, n)
		for i := range inputs {
			inputs[i] = -50.0 + float64(i)*100.0/float64(n)
			if inputs[i] == 0 {
				inputs[i] = 1.0
			}
			expAbs[i] = math.Abs(inputs[i])
			expNeg[i] = -inputs[i]
			expInv[i] = 1.0 / inputs[i]
		}
		if !sliceChecker("AbsBatch", arithmetic.AbsBatch(inputs), expAbs, 1e-9) {
			failed = true
		}
		if !sliceChecker("NegBatch", arithmetic.NegBatch(inputs), expNeg, 1e-9) {
			failed = true
		}
		if !sliceChecker("InvBatch", arithmetic.InvBatch(inputs), expInv, 1e-9) {
			failed = true
		}

		// 7. FMABatch
		c := make([]float64, n)
		expFMA := make([]float64, n)
		for i := range c {
			c[i] = float64(i) * 0.5
			expFMA[i] = a[i]*b[i] + c[i]
		}
		if !sliceChecker("FMABatch", arithmetic.FmaBatch(a, b, c), expFMA, 1e-9) {
			failed = true
		}

		// 8. AddScalarBatch, MulScalarBatch
		expAddScalar := make([]float64, n)
		expMulScalar := make([]float64, n)
		for i := range a {
			expAddScalar[i] = a[i] + 3.14
			expMulScalar[i] = a[i] * 2.5
		}
		if !sliceChecker("AddScalarBatch", arithmetic.AddScalarBatch(a, 3.14), expAddScalar, 1e-9) {
			failed = true
		}
		if !sliceChecker("MulScalarBatch", arithmetic.MulScalarBatch(a, 2.5), expMulScalar, 1e-9) {
			failed = true
		}
	}

	if failed {
		fmt.Println("=== WASM vs Go Math Parity Validation FAILED! ===")
		return false
	}
	fmt.Println("=== WASM vs Go Math Parity Validation PASSED SUCCESSFULLY! ===")
	return true
}
