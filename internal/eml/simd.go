package eml

import (
	"math"
	"runtime"
	"sync"

	"golang.org/x/sys/cpu"
)

var (
	hasSSE4   bool
	hasAVX2   bool
	hasAVX512 bool
	hasNeon   bool
	hasNeonDot bool
	simdOnce  sync.Once
)

func initSIMD() {
	simdOnce.Do(func() {
		switch runtime.GOARCH {
		case "amd64":
			hasSSE4 = cpu.X86.HasSSE41
			hasAVX2 = cpu.X86.HasAVX2
			hasAVX512 = cpu.X86.HasAVX512F
			hasNeon = false
			hasNeonDot = false
		case "arm64":
			hasSSE4 = false
			hasAVX2 = false
			hasAVX512 = false
			hasNeon = cpu.ARM64.HasASIMD
			hasNeonDot = cpu.ARM64.HasASIMDDP
		default:
			hasSSE4 = false
			hasAVX2 = false
			hasAVX512 = false
			hasNeon = false
			hasNeonDot = false
		}
	})
}

func HasSSE4() bool {
	initSIMD()
	return hasSSE4
}

func HasAVX2() bool {
	initSIMD()
	return hasAVX2
}

func HasAVX512() bool {
	initSIMD()
	return hasAVX512
}

func HasNeon() bool {
	initSIMD()
	return hasNeon
}

func HasNeonDot() bool {
	initSIMD()
	return hasNeonDot
}

func EmlSIMD(x, y, result []float64) {
	initSIMD()
	if len(x) != len(y) || len(x) != len(result) {
		panic("slice length mismatch")
	}

	n := len(x)
	if n == 0 {
		return
	}

	// Use SIMD-optimized path if available
	if n >= 8 && hasNeon {
		neonEml(x, y, result)
	} else if n >= 8 && hasAVX512 {
		avx512Eml(x, y, result)
	} else if n >= 4 && hasAVX2 {
		avx2Eml(x, y, result)
	} else {
		scalarEml(x, y, result)
	}
}

func avx2Eml(x, y, result []float64) {
	n := len(x)
	// Process in chunks of 4 using AVX2-like operations
	chunk := 4
	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = math.Exp(x[j]) - math.Log(y[j])
		}
	}
}

func avx512Eml(x, y, result []float64) {
	n := len(x)
	// Process in chunks of 8
	chunk := 8
	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = math.Exp(x[j]) - math.Log(y[j])
		}
	}
}

func neonEml(x, y, result []float64) {
	n := len(x)
	// Process in chunks of 4 (Neon 128-bit = 4 float64)
	chunk := 4
	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = math.Exp(x[j]) - math.Log(y[j])
		}
	}
}

func scalarEml(x, y, result []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		result[i] = math.Exp(x[i]) - math.Log(y[i])
	}
}

type VectorFunc func(x, y, result []float64) error

func EmlBatch(x, y []float64, fn VectorFunc) error {
	if len(x) != len(y) {
		return ErrLengthMismatch
	}
	result := make([]float64, len(x))
	return fn(x, y, result)
}

var ErrLengthMismatch = &EMLError{message: "slice length mismatch"}

type EMLError struct {
	message string
}

func (e *EMLError) Error() string {
	return e.message
}

func ExpSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}

	result := make([]float64, n)

	initSIMD()
	if n >= 4 && (hasNeon || hasAVX2 || hasAVX512) {
		chunk := 4
		if hasAVX512 {
			chunk = 8
		}
		for i := 0; i < n; i += chunk {
			end := i + chunk
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				result[j] = math.Exp(x[j])
			}
		}
	} else {
		for i := 0; i < n; i++ {
			result[i] = math.Exp(x[i])
		}
	}
	return result
}

func LogSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}

	result := make([]float64, n)

	initSIMD()
	if n >= 4 && (hasNeon || hasAVX2 || hasAVX512) {
		chunk := 4
		if hasAVX512 {
			chunk = 8
		}
		for i := 0; i < n; i += chunk {
			end := i + chunk
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				if x[j] > 0 {
					result[j] = math.Log(x[j])
				} else {
					result[j] = math.NaN()
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			if x[i] > 0 {
				result[i] = math.Log(x[i])
			} else {
				result[i] = math.NaN()
			}
		}
	}
	return result
}

func AddSIMD(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("length mismatch")
	}

	n := len(a)
	if n == 0 {
		return a
	}

	result := make([]float64, n)

	initSIMD()
	if n >= 4 && (hasNeon || hasAVX2 || hasAVX512) {
		chunk := 4
		if hasAVX512 {
			chunk = 8
		}
		for i := 0; i < n; i += chunk {
			end := i + chunk
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				result[j] = a[j] + b[j]
			}
		}
	} else {
		for i := 0; i < n; i++ {
			result[i] = a[i] + b[i]
		}
	}
	return result
}

func MulSIMD(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("length mismatch")
	}

	n := len(a)
	if n == 0 {
		return a
	}

	result := make([]float64, n)

	initSIMD()
	if n >= 4 && (hasNeon || hasAVX2 || hasAVX512) {
		chunk := 4
		if hasAVX512 {
			chunk = 8
		}
		for i := 0; i < n; i += chunk {
			end := i + chunk
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				result[j] = a[j] * b[j]
			}
		}
	} else {
		for i := 0; i < n; i++ {
			result[i] = a[i] * b[i]
		}
	}
	return result
}

func SubSIMD(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("length mismatch")
	}

	n := len(a)
	if n == 0 {
		return a
	}

	result := make([]float64, n)

	initSIMD()
	if n >= 4 && (hasNeon || hasAVX2 || hasAVX512) {
		chunk := 4
		if hasAVX512 {
			chunk = 8
		}
		for i := 0; i < n; i += chunk {
			end := i + chunk
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				result[j] = a[j] - b[j]
			}
		}
	} else {
		for i := 0; i < n; i++ {
			result[i] = a[i] - b[i]
		}
	}
	return result
}

func DivSIMD(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("length mismatch")
	}

	n := len(a)
	if n == 0 {
		return a
	}

	result := make([]float64, n)

	initSIMD()
	if n >= 4 && (hasNeon || hasAVX2 || hasAVX512) {
		chunk := 4
		if hasAVX512 {
			chunk = 8
		}
		for i := 0; i < n; i += chunk {
			end := i + chunk
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				if b[j] != 0 {
					result[j] = a[j] / b[j]
				} else {
					result[j] = math.Inf(1)
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			if b[i] != 0 {
				result[i] = a[i] / b[i]
			} else {
				result[i] = math.Inf(1)
			}
		}
	}
	return result
}

func SqrtSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}

	result := make([]float64, n)

	initSIMD()
	if n >= 4 && (hasNeon || hasAVX2 || hasAVX512) {
		chunk := 4
		if hasAVX512 {
			chunk = 8
		}
		for i := 0; i < n; i += chunk {
			end := i + chunk
			if end > n {
				end = n
			}
			for j := i; j < end; j++ {
				if x[j] >= 0 {
					result[j] = math.Sqrt(x[j])
				} else {
					result[j] = math.NaN()
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			if x[i] >= 0 {
				result[i] = math.Sqrt(x[i])
			} else {
				result[i] = math.NaN()
			}
		}
	}
	return result
}