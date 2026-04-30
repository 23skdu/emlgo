package arithmetic

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
	"github.com/emlgo/eml/pkg/logexp"
)

//go:inline
func Add(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}
	return x + y
}

//go:inline
func Sub(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}
	return x - y
}

//go:inline
func Mul(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}
	if x == 0 || y == 0 {
		return 0
	}
	return x * y
}

//go:inline
func Div(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}
	if y == 0 {
		if x > 0 {
			return math.Inf(1)
		} else if x < 0 {
			return math.Inf(-1)
		}
		return math.NaN()
	}
	return x / y
}

func Mod(x, y float64) float64 {
	if y == 0 || math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}
	return math.Mod(x, y)
}

func Remainder(x, y float64) float64 {
	if y == 0 || math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}
	return math.Remainder(x, y)
}

func Pow(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
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
		return math.NaN()
	}
	if x == 0 {
		if y > 0 {
			return 0
		}
		if y < 0 {
			return math.Inf(1)
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
	if x <= 0 || base <= 0 || base == 1 || math.IsNaN(x) || math.IsNaN(base) {
		return math.NaN()
	}
	return logexp.Log(x) / logexp.Log(base)
}

func LogBase2(x float64) float64 {
	if x <= 0 || math.IsNaN(x) {
		return math.NaN()
	}
	return logexp.Log(x) / logexp.Log(2)
}

func LogBase10(x float64) float64 {
	if x <= 0 || math.IsNaN(x) {
		return math.NaN()
	}
	return logexp.Log(x) / logexp.Log(10)
}

//go:inline
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

func Cbrt(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	return math.Cbrt(x)
}

func Hypot(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}
	if math.IsInf(x, 0) || math.IsInf(y, 0) {
		return math.Inf(1)
	}
	return math.Hypot(x, y)
}

func Max(x, y float64) float64 {
	if math.IsNaN(x) {
		return y
	}
	if math.IsNaN(y) {
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
	if math.Signbit(x) {
		return x
	}
	return y
}

func Min(x, y float64) float64 {
	if math.IsNaN(x) {
		return y
	}
	if math.IsNaN(y) {
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
	if math.Signbit(x) {
		return y
	}
	return x
}

func Floor(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}
	return math.Floor(x)
}

func Ceil(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}
	return math.Ceil(x)
}

func Trunc(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}
	return math.Trunc(x)
}

func Round(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}
	return math.Round(x)
}

func Abs(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	return math.Abs(x)
}

func Neg(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x == 0 {
		return math.Copysign(0, -1)
	}
	return -x
}

func Inv(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x == 0 {
		return math.Inf(1)
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
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 1) {
		return math.Inf(1)
	}
	if math.IsInf(x, -1) {
		return 0
	}
	return logexp.Exp(x)
}

func Log(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x > 0 {
		return logexp.Log(x)
	}
	if x == 0 {
		return math.Inf(-1)
	}
	return math.NaN()
}

func Log1p(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x > -1 {
		return math.Log1p(x)
	}
	return math.NaN()
}

func ExpM1(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	return math.Expm1(x)
}

func FMA(x, y, z float64) float64 {
	return math.FMA(x, y, z)
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
	if a > math.MaxInt64/gcd || a < math.MinInt64/gcd {
		return 0
	}
	return a / gcd * b
}

func isInteger(x float64) bool {
	_, frac := math.Modf(x)
	return frac == 0
}

// Integer-specific operations (optimized for int type without float conversion)

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

// Unsigned integer-specific operations (optimized for uint type)

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