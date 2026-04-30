package eml

import (
	"runtime"
	"sync"
)

var (
	hasSSE4    bool
	hasAVX2    bool
	hasAVX512  bool
	hasNeon    bool
	hasNeonDot bool
	simdOnce   sync.Once
)

func initSIMD() {
	simdOnce.Do(func() {
		detectSIMD()
	})
}

func detectSIMD() {
	switch runtime.GOARCH {
	case "amd64":
		detectAMD64SIMD()
	case "arm64":
		detectARM64SIMD()
	default:
		hasSSE4 = false
		hasAVX2 = false
		hasAVX512 = false
		hasNeon = false
		hasNeonDot = false
	}
}

// detectAMD64SIMD and detectARM64SIMD are implemented in architecture-specific files.

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
	chunk := 4
	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = nativeExp(x[j]) - nativeLog(y[j])
		}
	}
}

func avx512Eml(x, y, result []float64) {
	n := len(x)
	chunk := 8
	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = nativeExp(x[j]) - nativeLog(y[j])
		}
	}
}

func neonEml(x, y, result []float64) {
	n := len(x)
	chunk := 4
	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = nativeExp(x[j]) - nativeLog(y[j])
		}
	}
}

func scalarEml(x, y, result []float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		result[i] = nativeExp(x[i]) - nativeLog(y[i])
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

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = nativeExp(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = nativeExp(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func LogSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			if x[i] > 0 {
				result[i] = nativeLog(x[i])
			} else {
				result[i] = nan()
			}
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if x[j] > 0 {
					result[j] = nativeLog(x[j])
				} else {
					result[j] = nan()
				}
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func SqrtSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = nativeSqrt(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = nativeSqrt(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func SinSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = nativeSin(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = nativeSin(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func CosSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = nativeCos(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = nativeCos(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func SinCosSIMD(x []float64) (sin, cos []float64) {
	n := len(x)
	if n == 0 {
		return x, x
	}

	sin = make([]float64, n)
	cos = make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			sin[i], cos[i] = nativeSincos(x[i])
		}
		return
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				sin[j], cos[j] = nativeSincos(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return
}

func TanSIMD(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = nativeTan(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = nativeTan(x[j])
			}
		}(i, end)
	}
	wg.Wait()
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

	if n < 4 {
		for i := 0; i < n; i++ {
			result[i] = a[i] + b[i]
		}
		return result
	}

	if hasAVX512 && runtime.GOARCH == "amd64" {
		simdLen := (n / 8) * 8
		addAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b[i]
		}
		return result
	}

	if hasAVX2 && runtime.GOARCH == "amd64" {
		simdLen := (n / 4) * 4
		addAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b[i]
		}
		return result
	}

	if hasNeon && runtime.GOARCH == "arm64" {
		simdLen := (n / 2) * 2
		addNEON(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b[i]
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = a[j] + b[j]
			}
		}(i, end)
	}
	wg.Wait()
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

	if n < 4 {
		for i := 0; i < n; i++ {
			result[i] = a[i] - b[i]
		}
		return result
	}

	if hasAVX512 && runtime.GOARCH == "amd64" {
		simdLen := (n / 8) * 8
		subAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] - b[i]
		}
		return result
	}

	if hasAVX2 && runtime.GOARCH == "amd64" {
		simdLen := (n / 4) * 4
		subAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] - b[i]
		}
		return result
	}

	if hasNeon && runtime.GOARCH == "arm64" {
		simdLen := (n / 2) * 2
		subNEON(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] - b[i]
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = a[j] - b[j]
			}
		}(i, end)
	}
	wg.Wait()
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

	if n < 4 {
		for i := 0; i < n; i++ {
			result[i] = a[i] * b[i]
		}
		return result
	}

	if hasAVX512 && runtime.GOARCH == "amd64" {
		simdLen := (n / 8) * 8
		mulAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b[i]
		}
		return result
	}

	if hasAVX2 && runtime.GOARCH == "amd64" {
		simdLen := (n / 4) * 4
		mulAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b[i]
		}
		return result
	}

	if hasNeon && runtime.GOARCH == "arm64" {
		simdLen := (n / 2) * 2
		mulNEON(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b[i]
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = a[j] * b[j]
			}
		}(i, end)
	}
	wg.Wait()
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

	if n < 4 {
		for i := 0; i < n; i++ {
			result[i] = a[i] / b[i]
		}
		return result
	}

	if hasAVX512 && runtime.GOARCH == "amd64" {
		simdLen := (n / 8) * 8
		divAVX512(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] / b[i]
		}
		return result
	}

	if hasAVX2 && runtime.GOARCH == "amd64" {
		simdLen := (n / 4) * 4
		divAVX2(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] / b[i]
		}
		return result
	}

	if hasNeon && runtime.GOARCH == "arm64" {
		simdLen := (n / 2) * 2
		divNEON(a[:simdLen], b[:simdLen], result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] / b[i]
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = a[j] / b[j]
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func AddScalarSIMD(a []float64, b float64) []float64 {
	n := len(a)
	if n == 0 {
		return a
	}
	result := make([]float64, n)

	if n < 4 {
		for i := 0; i < n; i++ {
			result[i] = a[i] + b
		}
		return result
	}

	if hasAVX512 && runtime.GOARCH == "amd64" {
		simdLen := (n / 8) * 8
		addScalarAVX512(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b
		}
		return result
	}

	if hasAVX2 && runtime.GOARCH == "amd64" {
		simdLen := (n / 4) * 4
		addScalarAVX2(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b
		}
		return result
	}

	if hasNeon && runtime.GOARCH == "arm64" {
		simdLen := (n / 2) * 2
		addScalarNEON(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] + b
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = a[j] + b
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func MulScalarSIMD(a []float64, b float64) []float64 {
	n := len(a)
	if n == 0 {
		return a
	}
	result := make([]float64, n)

	if n < 4 {
		for i := 0; i < n; i++ {
			result[i] = a[i] * b
		}
		return result
	}

	if hasAVX512 && runtime.GOARCH == "amd64" {
		simdLen := (n / 8) * 8
		mulScalarAVX512(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b
		}
		return result
	}

	if hasAVX2 && runtime.GOARCH == "amd64" {
		simdLen := (n / 4) * 4
		mulScalarAVX2(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b
		}
		return result
	}

	if hasNeon && runtime.GOARCH == "arm64" {
		simdLen := (n / 2) * 2
		mulScalarNEON(a[:simdLen], b, result[:simdLen])
		for i := simdLen; i < n; i++ {
			result[i] = a[i] * b
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = a[j] * b
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func SinhBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = Sinh(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = Sinh(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func CoshBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = Cosh(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = Cosh(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func TanhBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = Tanh(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = Tanh(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func AsinhBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = Asinh(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = Asinh(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func AcoshBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = Acosh(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = Acosh(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}

func AtanhBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	if n < 256 {
		for i := 0; i < n; i++ {
			result[i] = Atanh(x[i])
		}
		return result
	}

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (n + numWorkers - 1) / numWorkers
	if chunkSize > 4096 {
		chunkSize = 4096
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				result[j] = Atanh(x[j])
			}
		}(i, end)
	}
	wg.Wait()
	return result
}