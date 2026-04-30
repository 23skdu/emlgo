package hyper

import (
	"math"

	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/logexp"
)

func Sinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	if x > 709.78 || x < -709.78 {
		if x > 0 {
			return math.Inf(1)
		}
		return math.Inf(-1)
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex - emx) / 2
}

func Cosh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return math.Abs(x)
	}
	if x > 709.78 || x < -709.78 {
		return math.Inf(1)
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	return (ex + emx) / 2
}

func Tanh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x > 709.78 {
		return 1
	}
	if x < -709.78 {
		return -1
	}
	if math.IsInf(x, 1) {
		return 1
	}
	if math.IsInf(x, -1) {
		return -1
	}
	ex := logexp.Exp(x)
	emx := logexp.Exp(-x)
	sum := ex + emx
	if math.IsInf(sum, 1) {
		if x > 0 {
			return 1
		}
		return -1
	}
	return (ex - emx) / sum
}

func Asinh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if math.IsInf(x, 0) {
		return x
	}
	if x == 0 {
		return 0
	}
	absX := math.Abs(x)
	if absX > 1e150 {
		// For extremely large x, use: asinh(x) ≈ ln(2x) + ln(1+1/(2x)) ≈ ln(2x) + 1/(2x)
		// But 2*x overflows, so use: ln(2) + ln(x) = ln(2) + ln(10)*log10(x)
		log10x := math.Log10(absX)
		approx := math.Ln2 + math.Ln10*log10x
		if x > 0 {
			return approx
		}
		return -approx
	}
	if absX > math.MaxFloat64/2 {
		if x > 0 {
			// For large x: asinh(x) ≈ ln(2x)
			// Use log10 to avoid overflow
			log10x := math.Log10(2 * absX)
			return math.Ln10 * log10x
		}
		log10x := math.Log10(2 * absX)
		return -math.Ln10 * log10x
	}
	term := arithmetic.Sqrt(x*x + 1)
	if math.IsInf(term, 1) {
		// Use log10 approximation
		log10x := math.Log10(2 * absX)
		if x > 0 {
			return math.Ln10 * log10x
		}
		return -math.Ln10 * log10x
	}
	return logexp.Log(x + term)
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
	if math.IsInf(x, 0) {
		return x
	}
	if x > math.MaxFloat64/2 {
		return logexp.Log(2*x) - 0.6931471805599453
	}
	return logexp.Log(x + arithmetic.Sqrt(x-1)*arithmetic.Sqrt(x+1))
}

func Atanh(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x <= -1 || x >= 1 {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}
	if math.Abs(x) > 0.9999999999999999 {
		if x > 0 {
			return math.Inf(1)
		}
		return math.Inf(-1)
	}
	return logexp.Log((1+x)/(1-x)) / 2
}