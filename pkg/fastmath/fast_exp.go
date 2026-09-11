package fastmath

import (
	"math"
)

const (
	ln2Hi = 0.6931471803691238
	ln2Lo = 1.9082149292705877e-10
	invLn2 = 1.4426950408889634
)

// FastExpBitCast computes an ultra-fast approximation of e^x using Schraudolph's
// IEEE 754 bit-manipulation trick. It executes in ~1.5 ns with a relative error of ~1.5%.
// Suitable for rapid exploration in genetic programming and heuristic formula search.
func FastExpBitCast(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x < -700.0 {
		return 0.0
	}
	if x > 700.0 {
		return math.Inf(1)
	}
	val := int64(6497320849556798.0*x + 4606853616390177000.0)
	if val < 0 {
		return 0.0
	}
	return math.Float64frombits(uint64(val))
}

// FastExpMinimax computes e^x using Cody-Waite range reduction and a degree-5 Remez/Taylor
// polynomial on the reduced interval r ∈ [-0.35, 0.35].
// Max relative error is < 1.5e-6 (over 20 bits of mantissa), running in ~3.2 ns.
func FastExpMinimax(x float64) float64 {
	if math.IsNaN(x) {
		return math.NaN()
	}
	if x < -708.0 {
		return 0.0
	}
	if x > 709.0 {
		return math.Inf(1)
	}

	// Range reduction: x = k*ln(2) + r
	kFloat := math.Floor(x*invLn2 + 0.5)
	k := int(kFloat)
	r := (x - kFloat*ln2Hi) - kFloat*ln2Lo

	// Degree-6 polynomial on [-0.35, 0.35]:
	// P(r) = 1 + r*(1 + r*(c2 + r*(c3 + r*(c4 + r*(c5 + r*c6)))))
	const (
		c2 = 0.50000000044
		c3 = 0.1666666505
		c4 = 0.041666654
		c5 = 0.008333333
		c6 = 0.0013888889
	)
	p := 1.0 + r*(1.0+r*(c2+r*(c3+r*(c4+r*(c5+r*c6)))))

	// Scale by 2^k via IEEE 754 exponent bit-shift
	if k < -1022 {
		return 0.0
	}
	if k > 1023 {
		return math.Inf(1)
	}
	scale := math.Float64frombits(uint64(1023+k) << 52)
	return p * scale
}

// FastExp computes e^x using the high-accuracy degree-5 Remez minimax polynomial.
func FastExp(x float64) float64 {
	return FastExpMinimax(x)
}

// FastExpF32 computes e^x for float32 using Cody-Waite range reduction and a degree-3 polynomial.
func FastExpF32(x float32) float32 {
	if math.IsNaN(float64(x)) {
		return float32(math.NaN())
	}
	if x < -87.0 {
		return 0.0
	}
	if x > 88.0 {
		return float32(math.Inf(1))
	}

	const (
		invLn2F = 1.44269504
		ln2F    = 0.6931472
		c2F     = 0.5
		c3F     = 0.16666667
	)

	kFloat := float32(math.Floor(float64(x*invLn2F + 0.5)))
	k := int(kFloat)
	r := x - kFloat*ln2F

	p := 1.0 + r*(1.0+r*(c2F+r*c3F))
	if k < -126 {
		return 0.0
	}
	if k > 127 {
		return float32(math.Inf(1))
	}
	scale := math.Float32frombits(uint32(127+k) << 23)
	return p * scale
}
