package main

import (
	"math"
	"testing"
)

func TestValidate100PercentCoverage(t *testing.T) {
	oldExit := exitFunc
	exitCode := 0
	exitFunc = func(code int) {
		exitCode = code
	}
	origAllResults := allResults
	defer func() {
		exitFunc = oldExit
		allResults = origAllResults
		forceFailTol = false
	}()

	// 1. Test main() execution
	allResults = nil
	verbose = false
	failedOnly = false
	typeFilter = ""
	main()

	// 2. Test each type filter
	for _, tf := range []string{"int", "uint", "float", "complex", "fastmath", "quant", "batch", "other"} {
		allResults = nil
		typeFilter = tf
		validateIntTypes()
		validateUintTypes()
		validateFloatTypes()
		validateComplexTypes()
		validateFastMath()
		validateQuant()
		validateBatch()
	}
	typeFilter = ""

	// 3. Test verbose and failedOnly flag combinations
	allResults = nil
	verbose = true
	failedOnly = false
	validateIntTypes()
	validateFastMath()
	validateQuant()
	validateBatch()

	allResults = nil
	verbose = false
	failedOnly = true
	validateIntTypes()
	validateFastMath()
	validateQuant()
	validateBatch()

	allResults = nil
	verbose = true
	failedOnly = true
	validateIntTypes()
	validateFastMath()
	validateQuant()
	validateBatch()

	// 4. Test with forceFailTol = true to hit all failure branches
	allResults = nil
	forceFailTol = true
	verbose = true
	failedOnly = false
	validateIntTypes()
	validateUintTypes()
	validateFloatTypes()
	validateComplexTypes()
	validateFastMath()
	validateQuant()
	validateBatch()

	allResults = nil
	verbose = false
	failedOnly = false
	validateIntTypes()
	validateUintTypes()
	validateFloatTypes()
	validateComplexTypes()
	validateFastMath()
	validateQuant()
	validateBatch()
	forceFailTol = false
	allResults = nil

	// 5. Test summarize with failed tests
	failResults := []ValidationResult{
		{Type: "float64", Function: "Exp", Passed: true, Message: "ok"},
		{Type: "float64", Function: "Log", Passed: false, Message: "mismatch"},
	}
	if summarize(failResults) {
		t.Errorf("summarize should return false when there are failures")
	}

	// 6. Test printSummary when there are failures
	allResults = failResults
	exitCode = 0
	printSummary()
	if exitCode != 1 {
		t.Errorf("expected exitCode 1, got %d", exitCode)
	}
	allResults = nil

	// 7. Exhaustive tolerance branch coverage
	infPos := math.Inf(1)
	infNeg := math.Inf(-1)
	nan := math.NaN()

	// withinTol
	if !withinTol(nan, nan, 0.1) {
		t.Error("withinTol(NaN, NaN) should be true")
	}
	if !withinTol(infPos, infPos, 0.1) {
		t.Error("withinTol(+Inf, +Inf) should be true")
	}
	if !withinTol(infNeg, infNeg, 0.1) {
		t.Error("withinTol(-Inf, -Inf) should be true")
	}
	if withinTol(infPos, infNeg, 0.1) {
		t.Error("withinTol(+Inf, -Inf) should be false")
	}
	if !withinTol(1.0, 1.0, 0.1) {
		t.Error("withinTol(1, 1) should be true")
	}
	if withinTol(1.0, 10.0, 0.1) {
		t.Error("withinTol(1, 10) should be false")
	}

	// withinTolFloat32
	if !withinTolFloat32(float32(nan), float32(nan)) {
		t.Error("withinTolFloat32(NaN, NaN) should be true")
	}
	if !withinTolFloat32(float32(infPos), float32(infPos)) {
		t.Error("withinTolFloat32(+Inf, +Inf) should be true")
	}
	if !withinTolFloat32(float32(infNeg), float32(infNeg)) {
		t.Error("withinTolFloat32(-Inf, -Inf) should be true")
	}
	if withinTolFloat32(float32(infPos), float32(infNeg)) {
		t.Error("withinTolFloat32(+Inf, -Inf) should be false")
	}
	if !withinTolFloat32(1.0, 1.0) {
		t.Error("withinTolFloat32(1, 1) should be true")
	}

	// withinTolComplex64
	if !withinTolComplex64(complex64(complex(1.0, 2.0)), complex(1.0, 2.0)) {
		t.Error("withinTolComplex64 should match")
	}

	// withinTolComplex128
	// Both NaN
	if !withinTolComplex128(complex(nan, nan), complex(nan, nan)) {
		t.Error("withinTolComplex128 NaN both should be true")
	}
	// Real +Inf
	if !withinTolComplex128(complex(infPos, nan), complex(infPos, nan)) {
		t.Error("withinTolComplex128 +Inf, NaN should be true")
	}
	if !withinTolComplex128(complex(infPos, 1.0), complex(infPos, 1.0)) {
		t.Error("withinTolComplex128 +Inf, 1 should be true")
	}
	// Real -Inf
	if !withinTolComplex128(complex(infNeg, nan), complex(infNeg, nan)) {
		t.Error("withinTolComplex128 -Inf, NaN should be true")
	}
	if !withinTolComplex128(complex(infNeg, 1.0), complex(infNeg, 1.0)) {
		t.Error("withinTolComplex128 -Inf, 1 should be true")
	}
	// Imag NaN
	if !withinTolComplex128(complex(1.0, nan), complex(1.0, nan)) {
		t.Error("withinTolComplex128 Imag NaN should be true")
	}
	// Real NaN
	if !withinTolComplex128(complex(nan, 1.0), complex(nan, 1.0)) {
		t.Error("withinTolComplex128 Real NaN should be true")
	}
	if !withinTolComplex128(complex(1.0, 2.0), complex(1.0, 2.0)) {
		t.Error("withinTolComplex128 finite should be true")
	}

	// 8. Complex helpers edge cases
	_ = trigComplexSin(infPos, infPos)
	_ = trigComplexSin(infPos, 1.0)
	_ = trigComplexSin(1.0, infPos)
	_ = trigComplexSin(1.0, 1.0)

	_ = trigComplexCos(infPos, infPos)
	_ = trigComplexCos(infPos, 1.0)
	_ = trigComplexCos(1.0, infPos)
	_ = trigComplexCos(1.0, 1.0)

	_ = complexExp(infPos, infPos)
	_ = complexExp(infPos, 0)
	_ = complexExp(1000.0, 1.0) // overflow exp
	_ = complexExp(1.0, 1.0)

	_ = complexLog(infPos, infPos)
	_ = complexLog(infPos, 1.0)
	_ = complexLog(1.0, infPos)
	_ = complexLog(0, 0)
	_ = complexLog(2.0, 1.0)
	_ = complexLog(1.0, 2.0)

	_ = complexSqrt(infPos, infPos)
	_ = complexSqrt(infPos, infNeg)
	_ = complexSqrt(infPos, 0)
	_ = complexSqrt(infPos, 1.0)
	_ = complexSqrt(1.0, infPos)
	_ = complexSqrt(0, 0)
	_ = complexSqrt(3.0, 4.0)
	_ = complexSqrt(3.0, -4.0)
}
