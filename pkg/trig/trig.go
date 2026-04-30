package trig

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/logexp"
)

//go:inline
func Sin(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	return math.Sin(x)
}

//go:inline
func Cos(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	return math.Cos(x)
}

//go:inline
func Tan(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	return math.Tan(x)
}

//go:inline
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

//go:inline
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

//go:inline
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

//go:inline
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

//go:inline
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

//go:inline
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

//go:inline
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
	return eml.SinSIMD(x)
}

func CosBatch(x []float64) []float64 {
	return eml.CosSIMD(x)
}

func SinCosBatch(x []float64) (sin, cos []float64) {
	return eml.SinCosSIMD(x)
}

func TanBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	sin, cos := SinCosBatch(x)
	for i := 0; i < n; i++ {
		if cos[i] != 0 {
			result[i] = sin[i] / cos[i]
		} else {
			result[i] = math.Inf(1)
		}
	}
	return result
}