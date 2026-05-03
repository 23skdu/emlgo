package eml

const (
	Pi          = 3.141592653589793238462643383279
	PiOver2     = 1.570796326794896619231321691640
	PiOver4     = 0.785398163397448309615660845820
	E           = 2.718281828459045235360287471353
	Sqrt2       = 1.414213562373095048801688724210
	Ln2         = 0.693147180559945309417232121458
	Ln10        = 2.302585092994045684017991454684
	MaxFloat64  = 1.797693134862315708145274237317043544980e+308
	SmallestNonZero = 4.940656458412465441765687928082213877e-324
)

func IsNaN(f float64) bool {
	return f != f
}

func IsInf(f float64, sign int) bool {
	if sign == 0 {
		return f == inf(1) || f == inf(-1)
	}
	if sign > 0 {
		return f == inf(1)
	}
	return f == inf(-1)
}

func IsFinite(f float64) bool {
	return f <= MaxFloat64 && f >= -MaxFloat64
}

func Signbit(f float64) bool {
	return f64bits(f)>>63 == 1
}

func NaN() float64 {
	return nan()
}

func Inf(sign int) float64 {
	return inf(sign)
}

//go:inline
func Floor(x float64) float64 {
	return floor(x)
}

//go:inline
func Ceil(x float64) float64 {
	return ceil(x)
}

//go:inline
func Trunc(x float64) float64 {
	return trunc(x)
}

//go:inline
func Round(x float64) float64 {
	return round(x)
}

//go:inline
func Abs(x float64) float64 {
	return nativeAbs(x)
}

//go:inline
func Neg(x float64) float64 {
	return nativeNeg(x)
}

//go:inline
func Inv(x float64) float64 {
	return nativeInv(x)
}

func Exp(x float64) float64 {
	return nativeExp(x)
}

func Log(x float64) float64 {
	return nativeLog(x)
}

func Log1p(x float64) float64 {
	return nativeLog1p(x)
}

func Expm1(x float64) float64 {
	return nativeExpm1(x)
}

func Sqrt(x float64) float64 {
	return nativeSqrt(x)
}

func Cbrt(x float64) float64 {
	return nativeCbrt(x)
}

func Pow(x, y float64) float64 {
	return nativePow(x, y)
}

func PowInt(x float64, n int) float64 {
	result := float64(1)
	for i := 0; i < n; i++ {
		result *= x
	}
	return result
}

func Sin(x float64) float64 {
	return nativeSin(x)
}

func Cos(x float64) float64 {
	return nativeCos(x)
}

func Tan(x float64) float64 {
	return nativeTan(x)
}

func Sincos(x float64) (sin, cos float64) {
	return nativeSincos(x)
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

func Sinh(x float64) float64 {
	if IsNaN(x) {
		return x
	}
	if IsInf(x, 0) {
		return x
	}
	if x > 709.78 || x < -709.78 {
		if x > 0 {
			return inf(1)
		}
		return inf(-1)
	}
	ex := nativeExp(x)
	emx := nativeExp(-x)
	return (ex - emx) / 2
}

func Cosh(x float64) float64 {
	if IsNaN(x) {
		return x
	}
	if IsInf(x, 0) {
		return nativeAbs(x)
	}
	if x > 709.78 || x < -709.78 {
		return inf(1)
	}
	ex := nativeExp(x)
	emx := nativeExp(-x)
	return (ex + emx) / 2
}

func Tanh(x float64) float64 {
	if IsNaN(x) {
		return x
	}
	if x > 709.78 {
		return 1
	}
	if x < -709.78 {
		return -1
	}
	if IsInf(x, 1) {
		return 1
	}
	if IsInf(x, -1) {
		return -1
	}
	ex := nativeExp(x)
	emx := nativeExp(-x)
	sum := ex + emx
	if IsInf(sum, 1) {
		if x > 0 {
			return 1
		}
		return -1
	}
	return (ex - emx) / sum
}

func Asinh(x float64) float64 {
	return nativeAsinh(x)
}

func Acosh(x float64) float64 {
	return nativeAcosh(x)
}

func Atanh(x float64) float64 {
	return nativeAtanh(x)
}

func Hypot(x, y float64) float64 {
	return nativeHypot(x, y)
}

func Max(x, y float64) float64 {
	return nativeMax(x, y)
}

func Min(x, y float64) float64 {
	return nativeMin(x, y)
}

func Mod(x, y float64) float64 {
	return nativeMod(x, y)
}

func Remainder(x, y float64) float64 {
	return nativeRemainder(x, y)
}

func Log10(x float64) float64 {
	return nativeLog10(x)
}

func Copysign(x, y float64) float64 {
	if (x > 0) == (y > 0) {
		return x
	}
	return -x
}

func Modf(x float64) (intPart, fracPart float64) {
	if isNaN(x) || isInf(x, 0) {
		return nan(), nan()
	}
	if x == 0 {
		return x, x
	}
	intPart = floor(x)
	fracPart = x - intPart
	return intPart, fracPart
}