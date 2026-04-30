package trig

import (
	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/logexp"
)

var (
	pi          = eml.Pi
	piOver2     = eml.PiOver2
	inf  = eml.Inf
	nan  = eml.NaN
	isNaN = eml.IsNaN
	isInf = eml.IsInf
	nativeSin  = eml.Sin
	nativeCos  = eml.Cos
	nativeTan  = eml.Tan
	nativeAsin = eml.Asin
	nativeAcos = eml.Acos
	nativeAtan = eml.Atan
	nativeAtan2 = eml.Atan2
	nativeSqrt = eml.Sqrt
	nativeAbs  = eml.Abs
	nativeSinh = eml.Sinh
	nativeCosh = eml.Cosh
	nativeTanh = eml.Tanh
)

func Sin(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	if x == 0 {
		return 0
	}
	return nativeSin(x)
}

func Cos(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return nativeCos(x)
}

func Tan(x float64) float64 {
	return nativeTan(x)
}

func Cot(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	sinx := Sin(x)
	cosx := Cos(x)
	if sinx == 0 {
		return nan()
	}
	return cosx / sinx
}

func Sec(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	cosx := Cos(x)
	if cosx == 0 {
		return inf(1)
	}
	return 1 / cosx
}

func Csc(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	sinx := Sin(x)
	if sinx == 0 {
		return inf(1)
	}
	return 1 / sinx
}

func Asin(x float64) float64 {
	return nativeAsin(x)
}

func Acos(x float64) float64 {
	return nativeAcos(x)
}

func Atan(x float64) float64 {
	return nativeAtan(x)
}

func Atan2(y, x float64) float64 {
	return nativeAtan2(y, x)
}

func Acot(x float64) float64 {
	if isNaN(x) || isInf(x, 0) {
		return nan()
	}
	return piOver2 - Atan(x)
}

func Asec(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x >= -1 && x <= 1 {
		return nan()
	}
	return Acos(1 / x)
}

func Acsc(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x >= -1 && x <= 1 {
		return nan()
	}
	return Asin(1 / x)
}

func Sinh(x float64) float64 {
	return nativeSinh(x)
}

func Cosh(x float64) float64 {
	return nativeCosh(x)
}

func Tanh(x float64) float64 {
	return nativeTanh(x)
}

func Coth(x float64) float64 {
	if isNaN(x) || x == 0 {
		return nan()
	}
	if isInf(x, 1) {
		return 1
	}
	if isInf(x, -1) {
		return -1
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex + emx) / (ex - emx)
}

func Sech(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if isInf(x, 0) {
		return 0
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return 2 / (ex + emx)
}

func Csch(x float64) float64 {
	if isNaN(x) || x == 0 {
		return nan()
	}
	if isInf(x, 0) {
		return 0
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return 2 / (ex - emx)
}

func Asinh(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if isInf(x, 0) {
		return x
	}
	if x == 0 {
		return 0
	}
	return logexp.Log(x + nativeSqrt(x*x+1))
}

func Acosh(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x < 1 {
		return nan()
	}
	if x == 1 {
		return 0
	}
	if isInf(x, 0) {
		return x
	}
	return logexp.Log(x + nativeSqrt(x-1)*nativeSqrt(x+1))
}

func Atanh(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x <= -1 || x >= 1 {
		return nan()
	}
	if x == 0 {
		return 0
	}
	return logexp.Log((1+x)/(1-x)) / 2
}

func Acoth(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x >= -1 && x <= 1 {
		return nan()
	}
	return Atanh(1 / x)
}

func Asech(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x <= 0 || x > 1 {
		return nan()
	}
	if x == 1 {
		return 0
	}
	return Acosh(1 / x)
}

func Acsch(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x == 0 {
		return inf(1)
	}
	return Asinh(1 / x)
}

func DegToRad(deg float64) float64 {
	return deg * pi / 180
}

func RadToDeg(rad float64) float64 {
	return rad * 180 / pi
}

func SinCos(x float64) (sin, cos float64) {
	if isNaN(x) || isInf(x, 0) {
		return nan(), nan()
	}
	sin = Sin(x)
	cos = Cos(x)
	return
}

func SinhCosh(x float64) (sinh, cosh float64) {
	if isNaN(x) {
		return nan(), nan()
	}
	if isInf(x, 0) {
		return x, nativeAbs(x)
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
	sin, cos := SinCosBatch(x)
	return eml.DivSIMD(sin, cos)
}