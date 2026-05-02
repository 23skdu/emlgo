package eml

import (
	"math"
	"unsafe"
)

// math constants removed to satisfy lint

//go:inline
func nan() float64 {
	return f64frombits(0x7FFFFFFFFFFFFFFF)
}

//go:inline
func inf(sign int) float64 {
	if sign > 0 {
		return f64frombits(0x7FF0000000000000)
	}
	return f64frombits(0xFFF0000000000000)
}

//go:inline
func isNaN(f float64) bool {
	return f != f
}

//go:inline
func isInf(f float64, sign int) bool {
	if sign > 0 {
		return f == inf(1)
	}
	return f == inf(-1)
}


//go:inline
func signbit(f float64) bool {
	return f64bits(f)>>63 == 1
}

//go:inline
func copysign(x, y float64) float64 {
	if (x > 0) == (y > 0) {
		return x
	}
	return -x
}

//go:inline
func f64bits(f float64) uint64 {
	// #nosec G103 - high performance float/bits conversion
	return *(*uint64)(unsafe.Pointer(&f))
}

//go:inline
func f64frombits(b uint64) float64 {
	// #nosec G103 - high performance float/bits conversion
	return *(*float64)(unsafe.Pointer(&b))
}

func nativeSqrt(x float64) float64 {
	return math.Sqrt(x)
}

//go:inline
func nativeExp(x float64) float64 {
	return expImpl(x)
}

func expImpl(x float64) float64 {
	return math.Exp(x)
}

func nativeLog(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x <= 0 {
		if x == 0 {
			return inf(-1)
		}
		return nan()
	}
	return math.Log(x)
}

//go:inline
func nativeSin(x float64) float64 {
	return math.Sin(x)
}

//go:inline
func nativeCos(x float64) float64 {
	return math.Cos(x)
}

//go:inline
func nativeSincos(x float64) (sin, cos float64) {
	return math.Sincos(x)
}


//go:inline
func floor(x float64) float64 {
	return math.Floor(x)
}

//go:inline
func ceil(x float64) float64 {
	return math.Ceil(x)
}

//go:inline
func trunc(x float64) float64 {
	return math.Trunc(x)
}

//go:inline
func round(x float64) float64 {
	return math.Round(x)
}

func nativeExpm1(x float64) float64 {
	return math.Expm1(x)
}

//go:inline
func nativeTan(x float64) float64 {
	s, c := nativeSincos(x)
	if c == 0 {
		return nan()
	}
	return s / c
}

//go:inline
func nativeAtan(x float64) float64 {
	return math.Atan(x)
}

//go:inline
func nativeAtan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

//go:inline
func nativeAsin(x float64) float64 {
	return math.Asin(x)
}

//go:inline
func nativeAcos(x float64) float64 {
	return math.Acos(x)
}

//go:inline
func nativeAsinh(x float64) float64 {
	return math.Asinh(x)
}

//go:inline
func nativeAcosh(x float64) float64 {
	return math.Acosh(x)
}

//go:inline
func nativeAtanh(x float64) float64 {
	return math.Atanh(x)
}

//go:inline
func nativeLog1p(x float64) float64 {
	return math.Log1p(x)
}

func nativePow(x, y float64) float64 {
	if x == 1 || y == 0 {
		return 1
	}
	if x == 0 {
		if y > 0 {
			return 0
		}
		return inf(1)
	}
	if x < 0 && !isInteger(y) {
		return nan()
	}
	if x < 0 {
		if int(y)%2 == 0 {
			return nativeExp(y * nativeLog(-x))
		}
		return -nativeExp(y * nativeLog(-x))
	}
	// x^y = exp(y * ln(x))
	return nativeExp(y * nativeLog(x))
}

func isInteger(x float64) bool {
	return x == floor(x)
}

//go:inline
func nativeInv(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x == 0 {
		return inf(1)
	}
	return 1 / x
}

//go:inline
func nativeNeg(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x == 0 {
		return copysign(0, -1)
	}
	return -x
}

//go:inline
func nativeAbs(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x < 0 {
		return -x
	}
	return x
}

//go:inline
func nativeMod(x, y float64) float64 {
	if isNaN(x) || isNaN(y) || y == 0 {
		return nan()
	}
	if y < 0 {
		y = -y
	}
	if x < 0 {
		return -nativeMod(-x, y)
	}
	for x >= y {
		x -= y
	}
	return x
}

//go:inline
func nativeRemainder(x, y float64) float64 {
	if isNaN(x) || isNaN(y) || y == 0 {
		return nan()
	}
	if y < 0 {
		y = -y
	}
	if x < 0 {
		return -nativeMod(-x, y)
	}
	return nativeMod(x, y)
}

//go:inline
func nativeHypot(x, y float64) float64 {
	if isNaN(x) || isNaN(y) {
		return nan()
	}
	if isInf(x, 0) || isInf(y, 0) {
		return inf(1)
	}
	if x == 0 {
		return nativeAbs(y)
	}
	if y == 0 {
		return nativeAbs(x)
	}

	if x < y {
		x, y = y, x
	}

	y = y / x
	return nativeAbs(x * nativeSqrt(1+y*y))
}

func nativeCbrt(x float64) float64 {
	if isNaN(x) || isInf(x, 0) || x == 0 {
		return x
	}

	sign := signbit(x)
	if sign {
		x = -x
	}

	exp := int(f64bits(x)>>52) - 1023
	frac := float64(f64bits(x)&((1<<52)-1)) / float64(1<<52)
	frac = nativePow(frac, 1.0/3.0)
	exp = exp / 3

	result := f64frombits(((uint64(exp+1023) << 52) | (uint64(frac*float64(1<<52)) & ((1<<52)-1))))

	if sign {
		return -result
	}
	return result
}

//go:inline
func nativeMax(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x > y {
		return x
	}
	if x < y {
		return y
	}
	if x == 0 && y == 0 {
		if signbit(x) {
			return y
		}
		return x
	}
	if x > y {
		return x
	}
	return y
}

//go:inline
func nativeMin(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x < y {
		return x
	}
	if x > y {
		return y
	}
	if x == 0 && y == 0 {
		if signbit(x) {
			return x
		}
		return y
	}
	if x < y {
		return x
	}
	return y
}

func nativeLog10(x float64) float64 {
	if isNaN(x) {
		return x
	}
	if x <= 0 {
		if x == 0 {
			return inf(-1)
		}
		return nan()
	}
	return nativeLog(x) / ln10()
}

//go:inline
func ln10() float64 {
	return 2.3025850929940456840179914546843642
}