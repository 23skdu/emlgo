package arithmetic

import (
	"math"

	"github.com/emlgo/eml/pkg/logexp"
)

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
	if y == 0 {
		return math.Inf(1)
	}
	return x / y
}

func Pow(x, y float64) float64 {
	if x <= 0 {
		return math.NaN()
	}
	return math.Pow(x, y)
}

func LogBase(x, base float64) float64 {
	if x <= 0 || base <= 0 || base == 1 {
		return math.NaN()
	}
	return math.Log(x) / math.Log(base)
}

func Sqrt(x float64) float64 {
	if x < 0 {
		return math.NaN()
	}
	return math.Sqrt(x)
}

func Neg(x float64) float64 {
	return -x
}

func Inv(x float64) float64 {
	if x == 0 {
		return math.Inf(1)
	}
	return 1 / x
}

func Square(x float64) float64 {
	return x * x
}

func Exp(x float64) float64 {
	return logexp.Exp(x)
}

func Log(x float64) float64 {
	return logexp.Log(x)
}