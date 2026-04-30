package trig

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/logexp"
)

func Sin(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	expIx := cexp(complex(0, x))
	return imag(expIx)
}

func Cos(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	expIx := cexp(complex(0, x))
	return real(expIx)
}

func cexp(z complex128) complex128 {
	a := real(z)
	b := imag(z)
	expA := logexp.Exp(a)
	return complex(expA*math.Cos(b), expA*math.Sin(b))
}

func Tan(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	sinx := Sin(x)
	cosx := Cos(x)
	if cosx == 0 {
		return math.Inf(1)
	}
	return sinx / cosx
}

func Cot(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	sinx := Sin(x)
	cosx := Cos(x)
	if sinx == 0 {
		return math.NaN()
	}
	return cosx / sinx
}

func Sec(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	cosx := Cos(x)
	if cosx == 0 {
		return math.Inf(1)
	}
	return 1 / cosx
}

func Csc(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	sinx := Sin(x)
	if sinx == 0 {
		return math.Inf(1)
	}
	return 1 / sinx
}

func Asin(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x < -1 || x > 1 {
		return math.NaN()
	}
	if x == 1 {
		return math.Pi / 2
	}
	if x == -1 {
		return -math.Pi / 2
	}
	if x == 0 {
		return 0
	}
	sqrt1mx2 := arithmetic.Sqrt(1 - x*x)
	if x > 0 {
		return math.Atan(x / sqrt1mx2)
	}
	return -math.Atan(-x / sqrt1mx2)
}

func Acos(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x < -1 || x > 1 {
		return math.NaN()
	}
	if x == 0 {
		return math.Pi / 2
	}
	if x == 1 {
		return 0
	}
	if x == -1 {
		return math.Pi
	}
	return math.Pi/2 - Asin(x)
}

func Atan(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	if x == 1 {
		return math.Pi / 4
	}
	if x == -1 {
		return -math.Pi / 4
	}
	if math.IsInf(x, 1) {
		return math.Pi / 2
	}
	if math.IsInf(x, -1) {
		return -math.Pi / 2
	}
	return math.Atan(x)
}

func Atan2(y, x float64) float64 {
	if math.IsNaN(y) || math.IsNaN(x) {
		return math.NaN()
	}
	if x > 0 {
		return Atan(y / x)
	}
	if x < 0 && y >= 0 {
		return Atan(y/x) + math.Pi
	}
	if x < 0 && y < 0 {
		return Atan(y/x) - math.Pi
	}
	if x == 0 && y > 0 {
		return math.Pi / 2
	}
	if x == 0 && y < 0 {
		return -math.Pi / 2
	}
	if x == 0 && y == 0 {
		return 0
	}
	return math.NaN()
}

func Acot(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	return math.Pi/2 - Atan(x)
}

func Asec(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x >= -1 && x <= 1 {
		return math.NaN()
	}
	return Acos(1 / x)
}

func Acsc(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x >= -1 && x <= 1 {
		return math.NaN()
	}
	return Asin(1 / x)
}

func Sinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex - emx) / 2
}

func Cosh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return math.Abs(x)
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex + emx) / 2
}

func Tanh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 1) {
		return 1
	}
	if math.IsInf(x, -1) {
		return -1
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex - emx) / (ex + emx)
}

func Coth(x float64) float64 {
	if math.IsNaN(x) || x == 0 {
		return math.NaN()
	}
	if math.IsInf(x, 1) {
		return 1
	}
	if math.IsInf(x, -1) {
		return -1
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex + emx) / (ex - emx)
}

func Sech(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return 0
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return 2 / (ex + emx)
}

func Csch(x float64) float64 {
	if math.IsNaN(x) || x == 0 {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return 0
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return 2 / (ex - emx)
}

func Asinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	if x == 0 {
		return 0
	}
	return logexp.Log(x + math.Sqrt(x*x+1))
}

func Acosh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x < 1 {
		return math.NaN()
	}
	if x == 1 {
		return 0
	}
	if math.IsInf(x, 0) {
		return x
	}
	return logexp.Log(x + math.Sqrt(x-1)*math.Sqrt(x+1))
}

func Atanh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x <= -1 || x >= 1 {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	return logexp.Log((1+x)/(1-x)) / 2
}

func Acoth(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x >= -1 && x <= 1 {
		return math.NaN()
	}
	return Atanh(1 / x)
}

func Asech(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x <= 0 || x > 1 {
		return math.NaN()
	}
	if x == 1 {
		return 0
	}
	return Acosh(1 / x)
}

func Acsch(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x == 0 {
		return math.Inf(1)
	}
	return Asinh(1 / x)
}

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func SinCos(x float64) (sin, cos float64) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN(), math.NaN()
	}
	sin = Sin(x)
	cos = Cos(x)
	return
}

func SinhCosh(x float64) (sinh, cosh float64) {
	if math.IsNaN(x) {
		return math.NaN(), math.NaN()
	}
	if math.IsInf(x, 0) {
		return x, math.Abs(x)
	}
	sinh = Sinh(x)
	cosh = Cosh(x)
	return
}

func SinBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	chunk := 4
	if eml.HasAVX512() {
		chunk = 8
	} else if eml.HasAVX2() || eml.HasNeon() {
		chunk = 4
	}

	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = Sin(x[j])
		}
	}
	return result
}

func CosBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	chunk := 4
	if eml.HasAVX512() {
		chunk = 8
	} else if eml.HasAVX2() || eml.HasNeon() {
		chunk = 4
	}

	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = Cos(x[j])
		}
	}
	return result
}

func SinCosBatch(x []float64) (sin, cos []float64) {
	n := len(x)
	if n == 0 {
		return x, x
	}
	sin = make([]float64, n)
	cos = make([]float64, n)

	chunk := 4
	if eml.HasAVX512() {
		chunk = 8
	} else if eml.HasAVX2() || eml.HasNeon() {
		chunk = 4
	}

	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			s, c := SinCos(x[j])
			sin[j] = s
			cos[j] = c
		}
	}
	return
}

func TanBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)

	chunk := 4
	if eml.HasAVX512() {
		chunk = 8
	} else if eml.HasAVX2() || eml.HasNeon() {
		chunk = 4
	}

	for i := 0; i < n; i += chunk {
		end := i + chunk
		if end > n {
			end = n
		}
		for j := i; j < end; j++ {
			result[j] = Tan(x[j])
		}
	}
	return result
}