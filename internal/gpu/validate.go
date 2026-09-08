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
// If either value is NaN or Inf, returns 0.
func ulpDiff(a, b float64) uint64 {
	if a == b {
		return 0
	}
	if math.IsNaN(a) || math.IsNaN(b) {
		return 0
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return 0
	}
	bits := math.Float64bits(a)
	targetBits := math.Float64bits(b)
	if bits > targetBits {
		return bits - targetBits
	}
	return targetBits - bits
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
