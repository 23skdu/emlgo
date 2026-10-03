package arithmetic

import (
	"math"

	"github.com/emlgo/eml/internal/constants"
	"github.com/emlgo/eml/internal/eml"
)

var (
	isNaN       = eml.IsNaN
	isInf       = eml.IsInf
	nan         = eml.NaN
	inf         = eml.Inf
	nativeLog   = eml.Log
	nativeLog1p = eml.Log1p
	nativeSqrt  = eml.Sqrt
	nativeHypot = eml.Hypot

	nativeFloor = eml.Floor
	nativeCeil  = eml.Ceil
	nativeTrunc = eml.Trunc
	nativeRound = eml.Round
	nativeExp   = eml.Exp
)

// parallelMap applies fn to each element of x, storing results in result.
// For large slices, work is distributed across the shared worker pool.
func parallelMap(x, result []float64, fn func(float64) float64) {
	n := len(x)
	eml.ForEachChunk(n, func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = fn(x[i])
		}
	})
}

// parallelMap2 applies fn to paired elements of a and b, storing results in result.
func parallelMap2(a, b, result []float64, fn func(float64, float64) float64) {
	n := len(a)
	eml.ForEachChunk(n, func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = fn(a[i], b[i])
		}
	})
}

func Add(x, y float64) float64 {
	return x + y
}

func Sub(x, y float64) float64 {
	return x - y
}

func Mul(x, y float64) float64 {
	return x * y
}

func Div(x, y float64) float64 {
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

func FmaBatch(a, b, c []float64) []float64 {
	return eml.FmaSIMD(a, b, c)
}

func Log2Batch(x []float64) []float64 {
	result := make([]float64, len(x))
	eml.Log2SIMDTo(x, result)
	return result
}

func Log10Batch(x []float64) []float64 {
	result := make([]float64, len(x))
	eml.Log10SIMDTo(x, result)
	return result
}

func TanBatch(x []float64) []float64 {
	result := make([]float64, len(x))
	eml.TanSIMDTo(x, result)
	return result
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
	if x == 0 && y < 0 {
		return inf(1)
	}
	if x < 0 && isInteger(y) {
		intY := int(y)
		if intY == math.MinInt {
			return nan()
		}
		if intY < 0 {
			intY = -intY
		}
		if intY%2 == 0 {
			return powStable(-x, y)
		}
		return -powStable(-x, y)
	}
	return powStable(x, y)
}

// powStable computes x^y with improved stability near x=1.
// For x near 1: x^y = exp(y * log(x)) ≈ exp(y * log1p(x-1))
func powStable(x, y float64) float64 {
	absX1 := x - 1
	if absX1 > -0.5 && absX1 < 0.5 {
		return nativeExp(y * nativeLog1p(absX1))
	}
	return nativeExp(y * nativeLog(x))
}

func PowInt(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n < 0 {
		if n == math.MinInt {
			return 1 / (PowInt(x, math.MaxInt) * x)
		}
		return 1 / PowInt(x, -n)
	}
	// Binary exponentiation (exponentiation by squaring)
	result := 1.0
	base := x
	for n > 0 {
		if n&1 == 1 {
			result *= base
		}
		base *= base
		n >>= 1
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
	return nativeLog(x) / constants.Ln2
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
	return y
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
	return eml.AbsScalar(x)
}

func Neg(x float64) float64 {
	return -x
}

func Inv(x float64) float64 {
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
		// For x near 1, use Log1p for better accuracy
		if x > 0.5 && x < 1.5 {
			return nativeLog1p(x - 1)
		}
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
	return eml.FmaScalar(x, y, z)
}

func GCD(a, b int64) int64 {
	if a < 0 {
		if a == math.MinInt64 {
			a = math.MaxInt64
		} else {
			a = -a
		}
	}
	if b < 0 {
		if b == math.MinInt64 {
			b = math.MaxInt64
		} else {
			b = -b
		}
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Note: GCD(MinInt64, x) uses MaxInt64 as the absolute value approximation.
// This is a known limitation documented in docs/nextsteps.md.

func LCM(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	gcd := GCD(a, b)
	quot := a / gcd
	// Check overflow in quot * b
	if (quot > 0 && b > 0 && quot > math.MaxInt64/b) ||
		(quot > 0 && b < 0 && quot > math.MinInt64/b) ||
		(quot < 0 && b > 0 && quot < math.MinInt64/b) ||
		(quot < 0 && b < 0 && quot < math.MaxInt64/b) {
		return 0
	}
	return quot * b
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
		if a == math.MinInt {
			return math.MaxInt
		}
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

func AbsBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Abs)
	return result
}

func NegBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Neg)
	return result
}

func InvBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Inv)
	return result
}

func FloorBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Floor)
	return result
}

func CeilBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Ceil)
	return result
}

func TruncBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Trunc)
	return result
}

func Log1pBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Log1p)
	return result
}

func Expm1Batch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Expm1)
	return result
}

func PowBatch(x []float64, y float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, func(v float64) float64 { return Pow(v, y) })
	return result
}

func CbrtBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Cbrt)
	return result
}

func HypotBatch(x, y []float64) []float64 {
	n := len(x)
	if n == 0 || len(y) == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap2(x, y, result, Hypot)
	return result
}

func MaxBatch(x []float64, y float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, func(v float64) float64 { return Max(v, y) })
	return result
}

func MinBatch(x []float64, y float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, func(v float64) float64 { return Min(v, y) })
	return result
}

func ExpM1(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	return eml.Expm1(x)
}

func AddBatch(x, y []float64) []float64 {
	return eml.AddSIMD(x, y)
}

func AddBatchTo(x, y, result []float64) {
	eml.AddSIMDTo(x, y, result)
}

func SubBatch(x, y []float64) []float64 {
	return eml.SubSIMD(x, y)
}

func SubBatchTo(x, y, result []float64) {
	eml.SubSIMDTo(x, y, result)
}

func MulBatch(x, y []float64) []float64 {
	return eml.MulSIMD(x, y)
}

func MulBatchTo(x, y, result []float64) {
	eml.MulSIMDTo(x, y, result)
}

func DivBatch(x, y []float64) []float64 {
	return eml.DivSIMD(x, y)
}

func DivBatchTo(x, y, result []float64) {
	eml.DivSIMDTo(x, y, result)
}

func AddScalarBatch(x []float64, y float64) []float64 {
	return eml.AddScalarSIMD(x, y)
}

func MulScalarBatch(x []float64, y float64) []float64 {
	return eml.MulScalarSIMD(x, y)
}

