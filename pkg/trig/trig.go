package trig

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
)

func Sin(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	return math.Sin(x)
}

func Cos(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	return math.Cos(x)
}

func Tan(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	cosx := math.Cos(x)
	if cosx == 0 {
		return math.Tan(x)
	}
	return math.Tan(x)
}

func Cot(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	tanx := math.Tan(x)
	if tanx == 0 {
		return math.NaN()
	}
	return 1 / tanx
}

func Sec(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	cosx := math.Cos(x)
	if cosx == 0 {
		return math.Inf(1)
	}
	return 1 / cosx
}

func Csc(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	sinx := math.Sin(x)
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
	return math.Asin(x)
}

func Acos(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x < -1 || x > 1 {
		return math.NaN()
	}
	return math.Acos(x)
}

func Atan(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	return math.Atan(x)
}

func Atan2(y, x float64) float64 {
	if math.IsNaN(y) || math.IsNaN(x) {
		return math.NaN()
	}
	return math.Atan2(y, x)
}

func Acot(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}
	if x == 0 {
		return math.Pi / 2
	}
	return math.Pi/2 - math.Atan(x)
}

func Asec(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x >= -1 && x <= 1 {
		return math.NaN()
	}
	return math.Acos(1 / x)
}

func Acsc(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x >= -1 && x <= 1 {
		return math.NaN()
	}
	return math.Asin(1 / x)
}

func Sinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	ex := eml.EmlOne(x)
	emx := eml.EmlOne(-x)
	return (ex - emx) / 2
}

func Cosh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return math.Abs(x)
	}
	ex := eml.EmlOne(x)
	emx := eml.EmlOne(-x)
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
	ex := eml.EmlOne(x)
	emx := eml.EmlOne(-x)
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
	ex := eml.EmlOne(x)
	emx := eml.EmlOne(-x)
	return (ex + emx) / (ex - emx)
}

func Sech(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return 0
	}
	ex := eml.EmlOne(x)
	emx := eml.EmlOne(-x)
	return 2 / (ex + emx)
}

func Csch(x float64) float64 {
	if math.IsNaN(x) || x == 0 {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return 0
	}
	ex := eml.EmlOne(x)
	emx := eml.EmlOne(-x)
	return 2 / (ex - emx)
}

func Asinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	return math.Asinh(x)
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
	return math.Acosh(x)
}

func Atanh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x <= -1 || x >= 1 {
		return math.NaN()
	}
	return math.Atanh(x)
}

func Acoth(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x >= -1 && x <= 1 {
		return math.NaN()
	}
	return math.Atanh(1 / x)
}

func Asech(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x <= 0 || x > 1 {
		return math.NaN()
	}
	return math.Acosh(1 / x)
}

func Acsch(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x == 0 {
		return math.Inf(1)
	}
	return math.Asinh(1 / x)
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
	return math.Sin(x), math.Cos(x)
}

func SinhCosh(x float64) (sinh, cosh float64) {
	if math.IsNaN(x) {
		return math.NaN(), math.NaN()
	}
	if math.IsInf(x, 0) {
		return x, math.Abs(x)
	}
	ex := eml.EmlOne(x)
	emx := eml.EmlOne(-x)
	sinh = (ex - emx) / 2
	cosh = (ex + emx) / 2
	return
}