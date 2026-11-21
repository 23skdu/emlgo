package gpu

import (
	"fmt"
	"math"
)

// BatchVerifier holds configuration for GPU vs CPU result validation.
type BatchVerifier struct {
	// MaxULP is the maximum allowed ULP difference for each element.
	MaxULP uint64
}

// DefaultVerifier returns a BatchVerifier with a 1-ULP tolerance.
func DefaultVerifier() *BatchVerifier {
	return &BatchVerifier{MaxULP: 1}
}

// VerifyOp compares GPU batch results against a reference CPU function
// and returns per-element ULP differences plus overall status.
func (v *BatchVerifier) VerifyOp(name string, input, gpuResult []float64, cpuRef func(float64) float64) (maxULP uint64, failed int, err error) {
	if len(input) != len(gpuResult) {
		return 0, 0, fmt.Errorf("input length %d != result length %d", len(input), len(gpuResult))
	}

	var max uint64
	var failCount int

	for i := range input {
		cpuVal := cpuRef(input[i])
		gpuVal := gpuResult[i]
		ulp := ulpDiff(gpuVal, cpuVal)
		if ulp > max {
			max = ulp
		}
		if ulp > v.MaxULP {
			failCount++
		}
	}

	return max, failCount, nil
}

// ulpDiff returns the ULP distance between two float64 values.
// NaN matches only NaN (returning 0); NaN vs non-NaN returns math.MaxUint64.
// Matching infinities (+Inf/+Inf or -Inf/-Inf) return 0; mismatched infinities or Inf vs finite return math.MaxUint64.
func ulpDiff(a, b float64) uint64 {
	if math.IsNaN(a) || math.IsNaN(b) {
		if math.IsNaN(a) && math.IsNaN(b) {
			return 0
		}
		return math.MaxUint64
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		if math.IsInf(a, 1) && math.IsInf(b, 1) {
			return 0
		}
		if math.IsInf(a, -1) && math.IsInf(b, -1) {
			return 0
		}
		return math.MaxUint64
	}
	if a == b {
		return 0
	}
	bitsA := math.Float64bits(a)
	bitsB := math.Float64bits(b)
	magA := bitsA & 0x7fffffffffffffff
	magB := bitsB & 0x7fffffffffffffff
	signA := bitsA >> 63
	signB := bitsB >> 63
	if signA == signB {
		if magA > magB {
			return magA - magB
		}
		return magB - magA
	}
	return magA + magB
}

// CPU reference functions for all GPU batch operations.
var cpuRefs = map[string]func(float64) float64{
	"Exp":  func(x float64) float64 { return math.Exp(x) },
	"Log":  func(x float64) float64 { return math.Log(x) },
	"Sin":  func(x float64) float64 { return math.Sin(x) },
	"Cos":  func(x float64) float64 { return math.Cos(x) },
	"Tan":  func(x float64) float64 { return math.Tan(x) },
	"Sinh": func(x float64) float64 { return math.Sinh(x) },
	"Cosh": func(x float64) float64 { return math.Cosh(x) },
	"Tanh": func(x float64) float64 { return math.Tanh(x) },
	"Sqrt": func(x float64) float64 { return math.Sqrt(x) },
	"Abs":  func(x float64) float64 { return math.Abs(x) },
	"Neg":  func(x float64) float64 { return -x },
	"Inv":  func(x float64) float64 { return 1 / x },
	"Fma":  func(x float64) float64 { return x }, // placeholder, actual FMA needs 3 inputs
}

// BinaryOpRefs holds CPU reference functions for binary (two-input) GPU operations.
var BinaryOpRefs = map[string]func(float64, float64) float64{
	"Eml": func(x, y float64) float64 { return math.Exp(x) - math.Log(y) },
}

// VerifyBinaryOp compares GPU batch results against a reference CPU function
// for binary operations and returns per-element ULP differences plus overall status.
func (v *BatchVerifier) VerifyBinaryOp(name string, inputA, inputB, gpuResult []float64, cpuRef func(float64, float64) float64) (maxULP uint64, failed int, err error) {
	if len(inputA) != len(gpuResult) || len(inputB) != len(gpuResult) {
		return 0, 0, fmt.Errorf("input length mismatch")
	}
	var max uint64
	var failCount int
	for i := range inputA {
		cpuVal := cpuRef(inputA[i], inputB[i])
		gpuVal := gpuResult[i]
		ulp := ulpDiff(gpuVal, cpuVal)
		if ulp > max {
			max = ulp
		}
		if ulp > v.MaxULP {
			failCount++
		}
	}
	return max, failCount, nil
}

// VerifyOpF32 compares float32 GPU batch results against a reference CPU function
// and returns the maximum ULP distance and count of elements exceeding MaxULP.
func (v *BatchVerifier) VerifyOpF32(name string, input, gpuResult []float32, cpuRef func(float64) float64) (maxULP uint64, failed int, err error) {
	if len(input) != len(gpuResult) {
		return 0, 0, fmt.Errorf("input length %d != result length %d", len(input), len(gpuResult))
	}
	var max uint64
	var failCount int
	for i := range input {
		cpuVal := float32(cpuRef(float64(input[i])))
		ulp := ulpDiffF32(gpuResult[i], cpuVal)
		if ulp > max {
			max = ulp
		}
		if ulp > v.MaxULP {
			failCount++
		}
	}
	return max, failCount, nil
}

// VerifyBinaryOpF32 compares float32 binary GPU batch results against a reference CPU function.
func (v *BatchVerifier) VerifyBinaryOpF32(name string, inputA, inputB, gpuResult []float32, cpuRef func(float64, float64) float64) (maxULP uint64, failed int, err error) {
	if len(inputA) != len(gpuResult) || len(inputB) != len(gpuResult) {
		return 0, 0, fmt.Errorf("input length mismatch")
	}
	var max uint64
	var failCount int
	for i := range inputA {
		cpuVal := float32(cpuRef(float64(inputA[i]), float64(inputB[i])))
		ulp := ulpDiffF32(gpuResult[i], cpuVal)
		if ulp > max {
			max = ulp
		}
		if ulp > v.MaxULP {
			failCount++
		}
	}
	return max, failCount, nil
}

func ulpDiffF32(a, b float32) uint64 {
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		if math.IsNaN(float64(a)) && math.IsNaN(float64(b)) {
			return 0
		}
		return math.MaxUint64
	}
	if math.IsInf(float64(a), 0) || math.IsInf(float64(b), 0) {
		if math.IsInf(float64(a), 1) && math.IsInf(float64(b), 1) {
			return 0
		}
		if math.IsInf(float64(a), -1) && math.IsInf(float64(b), -1) {
			return 0
		}
		return math.MaxUint64
	}
	if a == b {
		return 0
	}
	bitsA := uint64(math.Float32bits(a))
	bitsB := uint64(math.Float32bits(b))
	magA := bitsA & 0x7fffffff
	magB := bitsB & 0x7fffffff
	signA := bitsA >> 31
	signB := bitsB >> 31
	if signA == signB {
		if magA > magB {
			return magA - magB
		}
		return magB - magA
	}
	return magA + magB
}
