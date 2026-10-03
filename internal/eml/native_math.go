package eml

import (
	"math"
)

// Bitwise helpers for IEEE 754 floating point
func f64bits(f float64) uint64     { return math.Float64bits(f) }
func f64frombits(b uint64) float64 { return math.Float64frombits(b) }

func signbit(x float64) bool {
	return f64bits(x)&(1<<63) != 0
}

func nan() float64 {
	return math.NaN()
}

func inf(sign int) float64 {
	return math.Inf(sign)
}

func isNaN(f float64) bool {
	return math.IsNaN(f)
}

func isInf(f float64, sign int) bool {
	return math.IsInf(f, sign)
}

func nativeSqrt(x float64) float64 {
	return math.Sqrt(x)
}

func nativeExp(x float64) float64 {
	return math.Exp(x)
}

func nativeLog(x float64) float64 {
	if x > 0 && x < 0x1p-1022 {
		return math.Log(x*0x1p54) - 54*math.Ln2
	}
	return math.Log(x)
}

// ExactLogOracle evaluates ln(x) for x > 0 by scaling subnormals into the normal
// range [2^-1020, 2^-968) before invoking the logarithm, eliminating toolchain subnormal truncation error.
func ExactLogOracle(x float64) float64 {
	if x <= 0 {
		if x == 0 {
			return Inf(-1)
		}
		return NaN()
	}
	if IsInf(x, 1) {
		return Inf(1)
	}
	if IsNaN(x) {
		return NaN()
	}
	if x < 0x1p-1022 {
		return math.Log(x*0x1p54) - 54*math.Ln2
	}
	return math.Log(x)
}

// logScalar is the scalar logarithm used by the batch dispatchers. It applies
// the same domain rule as pkg/logexp.Log -- NaN for negative argument and -Inf for zero --
// so that Log and LogBatch cannot disagree.
//
// The unguarded nativeLog is still used for the EML operator itself and for the
// fused LogDiv/LogSub kernels, where math.Log semantics are what we want.
func logScalar(x float64) float64 {
	if x < 0 {
		return NaN()
	}
	if x == 0 {
		return Inf(-1)
	}
	return nativeLog(x)
}

func nativeSin(x float64) float64 {
	return math.Sin(x)
}

func nativeCos(x float64) float64 {
	return math.Cos(x)
}

func nativeLog2(x float64) float64 {
	if x > 0 && x < 0x1p-1022 {
		return math.Log2(x*0x1p54) - 54
	}
	return math.Log2(x)
}

func nativeLog10(x float64) float64 {
	if x > 0 && x < 0x1p-1022 {
		return math.Log10(x*0x1p54) - 54*0.3010299956639811952137388947244930267681898814621085
	}
	return math.Log10(x)
}

func nativeSincos(x float64) (sin, cos float64) {
	return math.Sincos(x)
}

func nativeExpm1(x float64) float64 {
	return math.Expm1(x)
}

func nativeTan(x float64) float64 {
	return math.Tan(x)
}

func nativeAtan(x float64) float64 {
	return math.Atan(x)
}

func nativeAtan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

func nativeAsin(x float64) float64 {
	return math.Asin(x)
}

func nativeAcos(x float64) float64 {
	return math.Acos(x)
}

func nativeAsinh(x float64) float64 {
	return math.Asinh(x)
}

func nativeAcosh(x float64) float64 {
	return math.Acosh(x)
}

func nativeAtanh(x float64) float64 {
	return math.Atanh(x)
}

func nativeLog1p(x float64) float64 {
	return math.Log1p(x)
}

func nativePow(x, y float64) float64 {
	return math.Pow(x, y)
}

func nativeInv(x float64) float64 {
	return 1 / x
}

func nativeNeg(x float64) float64 {
	return -x
}

func nativeAbs(x float64) float64 {
	return math.Abs(x)
}

func nativeMod(x, y float64) float64 {
	return math.Mod(x, y)
}

func nativeRemainder(x, y float64) float64 {
	return math.Remainder(x, y)
}

func nativeHypot(x, y float64) float64 {
	return math.Hypot(x, y)
}

func nativeCbrt(x float64) float64 {
	return math.Cbrt(x)
}

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
	return y
}

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
	return y
}

func copysign(x, y float64) float64 {
	return math.Copysign(x, y)
}

func floor(x float64) float64 {
	return math.Floor(x)
}

func ceil(x float64) float64 {
	return math.Ceil(x)
}

func trunc(x float64) float64 {
	return math.Trunc(x)
}

func round(x float64) float64 {
	return math.Round(x)
}
