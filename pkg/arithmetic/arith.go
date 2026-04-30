package arithmetic

import (
	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/logexp"
)

var (
	isNaN    = eml.IsNaN
	isInf    = eml.IsInf
	nan      = eml.NaN
	inf      = eml.Inf
	nativeLog   = eml.Log
	nativeSqrt  = eml.Sqrt
	nativeHypot = eml.Hypot
	nativeAbs   = eml.Abs
	nativeFloor = eml.Floor
	nativeCeil  = eml.Ceil
	nativeTrunc = eml.Trunc
	nativeRound = eml.Round
	nativeExp   = eml.Exp
	nativePow   = eml.Pow
)

func Add(x, y float64) float64 {
	if isNaN(x) || isNaN(y) {
		return nan()
	}
	return x + y
}

func Sub(x, y float64) float64 {
	if isNaN(x) || isNaN(y) {
		return nan()
	}
	return x - y
}

func Mul(x, y float64) float64 {
	if isNaN(x) || isNaN(y) {
		return nan()
	}
	if x == 0 || y == 0 {
		return 0
	}
	return x * y
}

func Div(x, y float64) float64 {
	if isNaN(x) || isNaN(y) {
		return nan()
	}
	if y == 0 {
		if x > 0 {
			return inf(1)
		} else if x < 0 {
			return inf(-1)
		}
		return nan()
	}
	return x / y
}

func Mod(x, y float64) float64 {
	if y == 0 || isNaN(x) || isNaN(y) {
		return nan()
	}
	return eml.Mod(x, y)
}

func Remainder(x, y float64) float64 {
	if y == 0 || isNaN(x) || isNaN(y) {
		return nan()
	}
	return eml.Remainder(x, y)
}

func Pow(x, y float64) float64 {
	if isNaN(x) || isNaN(y) {
		return nan()
	}
	if x == 0 && y > 0 {
		return 0
	}
	if x == 0 && y == 0 {
		return 1
	}
	if y == 0 {
		return 1
	}
	if x < 0 && !isInteger(y) {
		return nan()
	}
	if x == 0 {
		if y > 0 {
			return 0
		}
		if y < 0 {
			return inf(1)
		}
		return 1
	}
	if x < 0 && isInteger(y) {
		intY := int(y)
		if intY%2 == 0 {
			return logexp.Exp(y * logexp.Log(-x))
		}
		return -logexp.Exp(y * logexp.Log(-x))
	}
	return logexp.Exp(y * logexp.Log(x))
}

func PowInt(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n < 0 {
		return 1 / PowInt(x, -n)
	}
	result := float64(1)
	for i := 0; i < n; i++ {
		result *= x
	}
	return result
}

func LogBase(x, base float64) float64 {
	if x <= 0 || base <= 0 || base == 1 || isNaN(x) || isNaN(base) {
		return nan()
	}
	return nativeLog(x) / nativeLog(base)
}

func LogBase2(x float64) float64 {
	if x <= 0 || isNaN(x) {
		return nan()
	}
	return nativeLog(x) / 0.693147180559945309417232121458
}

func LogBase10(x float64) float64 {
	if x <= 0 || isNaN(x) {
		return nan()
	}
	return eml.Log10(x)
}

func Sqrt(x float64) float64 {
	return nativeSqrt(x)
}

func Cbrt(x float64) float64 {
	return eml.Cbrt(x)
}

func Hypot(x, y float64) float64 {
	return nativeHypot(x, y)
}

func Max(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x > y {
		return x
	}
	if y > x {
		return y
	}
	if x == 0 && y == 0 {
		return 0
	}
	if nativeAbs(x) > 0 {
		return x
	}
	return y
}

func Min(x, y float64) float64 {
	if isNaN(x) {
		return y
	}
	if isNaN(y) {
		return x
	}
	if x < y {
		return x
	}
	if y < x {
		return y
	}
	if x == 0 && y == 0 {
		return 0
	}
	if nativeAbs(x) > 0 {
		return y
	}
	return x
}

func Floor(x float64) float64 {
	return nativeFloor(x)
}

func Ceil(x float64) float64 {
	return nativeCeil(x)
}

func Trunc(x float64) float64 {
	return nativeTrunc(x)
}

func Round(x float64) float64 {
	return nativeRound(x)
}

func Abs(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	return nativeAbs(x)
}

func Neg(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x == 0 {
		return eml.Copysign(0, -1)
	}
	return -x
}

func Inv(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x == 0 {
		return inf(1)
	}
	return 1 / x
}

func Square(x float64) float64 {
	return x * x
}

func Cube(x float64) float64 {
	return x * x * x
}

func Exp(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if isInf(x, 1) {
		return inf(1)
	}
	if isInf(x, -1) {
		return 0
	}
	return nativeExp(x)
}

func Log(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x > 0 {
		return nativeLog(x)
	}
	if x == 0 {
		return inf(-1)
	}
	return nan()
}

func Log1p(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	if x > -1 {
		return eml.Log1p(x)
	}
	return nan()
}

func Expm1(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	return eml.Expm1(x)
}

func FMA(x, y, z float64) float64 {
	return x*y + z
}

func GCD(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	gcd := GCD(a, b)
	if a > 9223372036854775807/gcd || a < -9223372036854775807/gcd {
		return 0
	}
	return a / gcd * b
}

func isInteger(x float64) bool {
	_, frac := eml.Modf(x)
	return frac == 0
}

func IntAdd(a, b int) int { return a + b }
func IntSub(a, b int) int { return a - b }
func IntMul(a, b int) int { return a * b }

func IntDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

func IntMod(a, b int) int {
	if b == 0 {
		return 0
	}
	return a % b
}

func IntAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func UintAdd(a, b uint) uint { return a + b }
func UintSub(a, b uint) uint { return a - b }
func UintMul(a, b uint) uint { return a * b }

func UintDiv(a, b uint) uint {
	if b == 0 {
		return 0
	}
	return a / b
}

func UintMod(a, b uint) uint {
	if b == 0 {
		return 0
	}
	return a % b
}

func UintMax(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func UintMin(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func SqrtBatch(x []float64) []float64 {
	return eml.SqrtSIMD(x)
}

func ExpM1(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	return eml.Expm1(x)
}