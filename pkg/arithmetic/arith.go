package arithmetic

import (
	"math"

	"github.com/23skdu/emlgo/internal/constants"
	"github.com/23skdu/emlgo/internal/eml"
)

func isNaN(x float64) bool             { return eml.IsNaN(x) }
func isInf(x float64, s int) bool      { return eml.IsInf(x, s) }
func nan() float64                     { return eml.NaN() }
func inf(s int) float64                { return eml.Inf(s) }
func nativeLog(x float64) float64      { return eml.Log(x) }
func nativeLog1p(x float64) float64    { return eml.Log1p(x) }
func nativeSqrt(x float64) float64     { return eml.Sqrt(x) }
func nativeHypot(x, y float64) float64 { return eml.Hypot(x, y) }
func nativeFloor(x float64) float64    { return eml.Floor(x) }
func nativeCeil(x float64) float64     { return eml.Ceil(x) }
func nativeTrunc(x float64) float64    { return eml.Trunc(x) }
func nativeRound(x float64) float64    { return eml.Round(x) }
func nativeExp(x float64) float64      { return eml.Exp(x) }

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

// Add returns the sum of x and y.
func Add(x, y float64) float64 {
	return x + y
}

// Sub returns the difference x - y.
func Sub(x, y float64) float64 {
	return x - y
}

// Mul returns the product of x and y.
func Mul(x, y float64) float64 {
	return x * y
}

// Div returns the quotient x / y.
func Div(x, y float64) float64 {
	return x / y
}

// Mod returns the floating-point remainder of x/y: the result has the sign of x
// and a magnitude less than |y|. It returns NaN if y is zero or either operand is NaN.
func Mod(x, y float64) float64 {
	if y == 0 || isNaN(x) || isNaN(y) {
		return nan()
	}
	return eml.Mod(x, y)
}

// Remainder returns the IEEE 754 floating-point remainder of x/y.
// It returns NaN if y is zero or either operand is NaN.
func Remainder(x, y float64) float64 {
	if y == 0 || isNaN(x) || isNaN(y) {
		return nan()
	}
	return eml.Remainder(x, y)
}

// FmaBatch returns a newly allocated slice containing a[i]*b[i] + c[i] for each index.
// It panics if the lengths differ; for empty inputs a is returned unchanged.
func FmaBatch(a, b, c []float64) []float64 {
	return eml.FmaSIMD(a, b, c)
}

// Log2Batch returns a newly allocated slice containing the base-2 logarithm of each element of x.
func Log2Batch(x []float64) []float64 {
	result := make([]float64, len(x))
	eml.Log2SIMDTo(x, result)
	return result
}

// Log10Batch returns a newly allocated slice containing the base-10 logarithm of each element of x.
func Log10Batch(x []float64) []float64 {
	result := make([]float64, len(x))
	eml.Log10SIMDTo(x, result)
	return result
}

// TanBatch returns a newly allocated slice containing the tangent of each element of x.
func TanBatch(x []float64) []float64 {
	result := make([]float64, len(x))
	eml.TanSIMDTo(x, result)
	return result
}

// Pow returns x raised to the power y. NaN operands yield NaN, 0^0 and x^0 yield 1,
// 0 raised to a negative power yields +Inf, and a negative base with a non-integer
// exponent yields NaN. A log1p-based path keeps accuracy near x = 1.
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
	if y == 0 || x == 1 || (x == -1 && isInf(y, 0)) {
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
	if absX1 > -0.125 && absX1 < 0.125 && math.Abs(y) <= 64 {
		return nativeExp(y * nativeLog1p(absX1))
	}
	return eml.Pow(x, y)
}

// PowInt raises x to the integer power n using binary exponentiation.
// A negative n yields the reciprocal; n == math.MinInt is handled without overflow.
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

// LogBase returns the logarithm of x in the given base, computed as log(x)/log(base).
// It returns NaN if x <= 0, base <= 0, base == 1, or either operand is NaN.
func LogBase(x, base float64) float64 {
	if x <= 0 || base <= 0 || base == 1 || isNaN(x) || isNaN(base) {
		return nan()
	}
	return nativeLog(x) / nativeLog(base)
}

// LogBase2 returns the base-2 logarithm of x.
// It returns NaN if x <= 0 or x is NaN.
func LogBase2(x float64) float64 {
	if x <= 0 || isNaN(x) {
		return nan()
	}
	return nativeLog(x) / constants.Ln2
}

// LogBase10 returns the base-10 logarithm of x.
// It returns NaN if x <= 0 or x is NaN.
func LogBase10(x float64) float64 {
	if x <= 0 || isNaN(x) {
		return nan()
	}
	return eml.Log10(x)
}

// Sqrt returns the square root of x.
func Sqrt(x float64) float64 {
	return nativeSqrt(x)
}

// Cbrt returns the cube root of x.
func Cbrt(x float64) float64 {
	return eml.Cbrt(x)
}

// Hypot returns Sqrt(x*x + y*y), taking care to avoid unnecessary overflow and underflow.
func Hypot(x, y float64) float64 {
	return nativeHypot(x, y)
}

// Max returns the larger of x and y. A NaN operand is ignored and the other
// operand is returned; the result is NaN only if both operands are NaN.
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

// Min returns the smaller of x and y. A NaN operand is ignored and the other
// operand is returned; the result is NaN only if both operands are NaN.
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

// Floor returns the largest integer value not greater than x.
func Floor(x float64) float64 {
	return nativeFloor(x)
}

// Ceil returns the smallest integer value not less than x.
func Ceil(x float64) float64 {
	return nativeCeil(x)
}

// Trunc returns the integer part of x, truncated toward zero.
func Trunc(x float64) float64 {
	return nativeTrunc(x)
}

// Round returns the nearest integer, rounding half away from zero.
func Round(x float64) float64 {
	return nativeRound(x)
}

// Abs returns the absolute value of x.
func Abs(x float64) float64 {
	return eml.AbsScalar(x)
}

// Neg returns -x.
func Neg(x float64) float64 {
	return -x
}

// Inv returns the reciprocal of x (1 / x).
func Inv(x float64) float64 {
	return 1 / x
}

// Square returns x * x.
func Square(x float64) float64 {
	return x * x
}

// Cube returns x * x * x.
func Cube(x float64) float64 {
	return x * x * x
}

// Exp returns e raised to the power x. NaN yields NaN, +Inf yields +Inf and -Inf yields 0.
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

// Log returns the natural logarithm of x. It returns NaN for x < 0 or NaN,
// and -Inf for x == 0.
func Log(x float64) float64 {
	if isNaN(x) || x < 0 {
		return nan()
	}
	if x == 0 {
		return inf(-1)
	}
	return nativeLog(x)
}

// Log1p returns log(1 + x), accurate for small x. It returns NaN for x < -1 or
// NaN, and -Inf for x == -1.
func Log1p(x float64) float64 {
	if isNaN(x) || x < -1 {
		return nan()
	}
	if x == -1 {
		return inf(-1)
	}
	return eml.Log1p(x)
}

// Expm1 returns e**x - 1, accurate for small x. NaN inputs yield NaN.
func Expm1(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	return eml.Expm1(x)
}

// FMA returns x*y + z, using a fused multiply-add instruction when available.
func FMA(x, y, z float64) float64 {
	return eml.FmaScalar(x, y, z)
}

func absUint64(x int64) uint64 {
	if x < 0 {
		return -uint64(x) // #nosec G115 -- intentional two's complement magnitude for MinInt64
	}
	return uint64(x) // #nosec G115 -- x >= 0
}

// GCD returns the non-negative greatest common divisor of a and b, computed with
// the Euclidean algorithm on their magnitudes; it saturates at math.MaxInt64.
func GCD(a, b int64) int64 {
	ua := absUint64(a)
	ub := absUint64(b)
	for ub != 0 {
		ua, ub = ub, ua%ub
	}
	if ua > uint64(math.MaxInt64) {
		return math.MaxInt64
	}
	return int64(ua) // #nosec G115 -- ua <= math.MaxInt64
}

// LCM returns a*b / gcd(a, b), with the sign of the product a*b.
// It returns 0 if either argument is zero or the result would overflow int64.
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

// IntAdd returns a + b.
func IntAdd(a, b int) int { return a + b }

// IntSub returns a - b.
func IntSub(a, b int) int { return a - b }

// IntMul returns a * b.
func IntMul(a, b int) int { return a * b }

// IntDiv returns a / b, or 0 if b is zero.
func IntDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// IntMod returns a % b, or 0 if b is zero.
func IntMod(a, b int) int {
	if b == 0 {
		return 0
	}
	return a % b
}

// IntAbs returns the absolute value of a; math.MinInt saturates to math.MaxInt.
func IntAbs(a int) int {
	if a < 0 {
		if a == math.MinInt {
			return math.MaxInt
		}
		return -a
	}
	return a
}

// IntMax returns the larger of a and b.
func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// IntMin returns the smaller of a and b.
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// UintAdd returns a + b.
func UintAdd(a, b uint) uint { return a + b }

// UintSub returns a - b.
func UintSub(a, b uint) uint { return a - b }

// UintMul returns a * b.
func UintMul(a, b uint) uint { return a * b }

// UintDiv returns a / b, or 0 if b is zero.
func UintDiv(a, b uint) uint {
	if b == 0 {
		return 0
	}
	return a / b
}

// UintMod returns a % b, or 0 if b is zero.
func UintMod(a, b uint) uint {
	if b == 0 {
		return 0
	}
	return a % b
}

// UintMax returns the larger of a and b.
func UintMax(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

// UintMin returns the smaller of a and b.
func UintMin(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

// SqrtBatch returns a newly allocated slice containing the square root of each element of x.
func SqrtBatch(x []float64) []float64 {
	return eml.SqrtSIMD(x)
}

// AbsBatch returns a newly allocated slice containing the absolute value of each
// element of x; an empty x is returned unchanged.
func AbsBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.AbsSIMDTo(x, result)
	return result
}

// NegBatch returns a newly allocated slice containing the negation of each element
// of x; an empty x is returned unchanged.
func NegBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.NegSIMDTo(x, result)
	return result
}

// InvBatch returns a newly allocated slice containing the reciprocal of each
// element of x; an empty x is returned unchanged.
func InvBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.InvSIMDTo(x, result)
	return result
}

// FloorBatch returns a newly allocated slice containing the floor of each element
// of x; an empty x is returned unchanged.
func FloorBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.FloorSIMDTo(x, result)
	return result
}

// CeilBatch returns a newly allocated slice containing the ceiling of each element
// of x; an empty x is returned unchanged.
func CeilBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.CeilSIMDTo(x, result)
	return result
}

// TruncBatch returns a newly allocated slice containing the truncation toward zero
// of each element of x; an empty x is returned unchanged.
func TruncBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.TruncSIMDTo(x, result)
	return result
}

// Log1pBatch returns a newly allocated slice containing Log1p of each element of x;
// an empty x is returned unchanged.
func Log1pBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Log1p)
	return result
}

// Expm1Batch returns a newly allocated slice containing Expm1 of each element of x;
// an empty x is returned unchanged.
func Expm1Batch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Expm1)
	return result
}

// PowBatch returns a newly allocated slice containing Pow(v, y) for each element v
// of x; an empty x is returned unchanged.
func PowBatch(x []float64, y float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, func(v float64) float64 { return Pow(v, y) })
	return result
}

// CbrtBatch returns a newly allocated slice containing the cube root of each element
// of x; an empty x is returned unchanged.
func CbrtBatch(x []float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap(x, result, Cbrt)
	return result
}

// HypotBatch returns a newly allocated slice containing Hypot(x[i], y[i]) for each
// index of x; if x or y is empty, x itself is returned unchanged. y must have at
// least len(x) elements.
func HypotBatch(x, y []float64) []float64 {
	n := len(x)
	if n == 0 || len(y) == 0 {
		return x
	}
	result := make([]float64, n)
	parallelMap2(x, y, result, Hypot)
	return result
}

// MaxBatch returns a newly allocated slice containing Max(x[i], y) for each element
// of x; an empty x is returned unchanged.
func MaxBatch(x []float64, y float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.ForEachChunk(n, func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = Max(x[i], y)
		}
	})
	return result
}

// MinBatch returns a newly allocated slice containing Min(x[i], y) for each element
// of x; an empty x is returned unchanged.
func MinBatch(x []float64, y float64) []float64 {
	n := len(x)
	if n == 0 {
		return x
	}
	result := make([]float64, n)
	eml.ForEachChunk(n, func(start, end int) {
		for i := start; i < end; i++ {
			result[i] = Min(x[i], y)
		}
	})
	return result
}

// ExpM1 returns e**x - 1, accurate for small x. NaN inputs yield NaN.
func ExpM1(x float64) float64 {
	if isNaN(x) {
		return nan()
	}
	return eml.Expm1(x)
}

// AddBatch returns a newly allocated slice containing x[i] + y[i]; the result does
// not alias x or y. It panics if the lengths differ.
func AddBatch(x, y []float64) []float64 {
	return eml.AddSIMD(x, y)
}

// AddBatchTo computes x[i] + y[i] and stores the results in result. It panics if
// the lengths differ; result may alias x or y with the same offset.
func AddBatchTo(x, y, result []float64) {
	eml.AddSIMDTo(x, y, result)
}

// SubBatch returns a newly allocated slice containing x[i] - y[i]; the result does
// not alias x or y. It panics if the lengths differ.
func SubBatch(x, y []float64) []float64 {
	return eml.SubSIMD(x, y)
}

// SubBatchTo computes x[i] - y[i] and stores the results in result. It panics if
// the lengths differ; result may alias x or y with the same offset.
func SubBatchTo(x, y, result []float64) {
	eml.SubSIMDTo(x, y, result)
}

// MulBatch returns a newly allocated slice containing x[i] * y[i]; the result does
// not alias x or y. It panics if the lengths differ.
func MulBatch(x, y []float64) []float64 {
	return eml.MulSIMD(x, y)
}

// MulBatchTo computes x[i] * y[i] and stores the results in result. It panics if
// the lengths differ; result may alias x or y with the same offset.
func MulBatchTo(x, y, result []float64) {
	eml.MulSIMDTo(x, y, result)
}

// DivBatch returns a newly allocated slice containing x[i] / y[i]; the result does
// not alias x or y. It panics if the lengths differ.
func DivBatch(x, y []float64) []float64 {
	return eml.DivSIMD(x, y)
}

// DivBatchTo computes x[i] / y[i] and stores the results in result. It panics if
// the lengths differ; result may alias x or y with the same offset.
func DivBatchTo(x, y, result []float64) {
	eml.DivSIMDTo(x, y, result)
}

// AddScalarBatch returns a newly allocated slice containing x[i] + y for each element of x.
func AddScalarBatch(x []float64, y float64) []float64 {
	return eml.AddScalarSIMD(x, y)
}

// MulScalarBatch returns a newly allocated slice containing x[i] * y for each element of x.
func MulScalarBatch(x []float64, y float64) []float64 {
	return eml.MulScalarSIMD(x, y)
}
