// Package bigmath provides arbitrary-precision EML operations using math/big.
// This backend is used for symbolic verification, identity checking, and
// deep tree reduction where float64 precision is insufficient.
package bigmath

import (
	"math"
	"math/big"
)

// Prec defines the default precision in bits for big.Float operations.
const Prec = 256

var (
	one  = new(big.Float).SetPrec(Prec).SetInt64(1)
	two  = new(big.Float).SetPrec(Prec).SetInt64(2)
	half = new(big.Float).SetPrec(Prec).SetFloat64(0.5)
)

// Eml computes exp(x) - log(y) at arbitrary precision.
func Eml(x, y *big.Float) *big.Float {
	e := Exp(x)
	l := Log(y)
	return new(big.Float).Sub(e, l)
}

// Exp computes exp(x) at arbitrary precision using range reduction + Taylor series.
func Exp(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}

	if x.Sign() == 0 {
		return new(big.Float).SetPrec(prec).SetInt64(1)
	}

	if x.Sign() < 0 {
		threshold := new(big.Float).SetPrec(prec).SetFloat64(-750)
		if x.Cmp(threshold) < 0 {
			return new(big.Float).SetPrec(prec)
		}
	}

	absX := new(big.Float).Abs(x)
	wprec := prec + 64

	// Range reduction: exp(x) = exp(x/2^k)^(2^k)
	// Choose k so |x/2^k| < 0.5 for fast Taylor convergence
	k := 0
	{
		var f64abs float64
		if absX.IsInf() {
			f64abs = math.MaxFloat64
		} else {
			f64abs, _ = absX.Float64()
		}
		if f64abs > 0.5 {
			k = int(math.Ceil(math.Log2(f64abs))) + 2
		} else {
			k = 0
		}
	}

	shifted := new(big.Float).SetPrec(wprec)
	if k > 0 {
		divisor := new(big.Float).SetPrec(wprec).SetMantExp(
			new(big.Float).SetPrec(wprec).SetInt64(1), k)
		shifted.Quo(x, divisor)
	} else {
		shifted.Copy(x)
	}

	result := expTaylor(shifted, wprec)

	for i := 0; i < k; i++ {
		result.Mul(result, result)
	}

	return result.SetPrec(prec)
}

func expTaylor(x *big.Float, prec uint) *big.Float {
	result := new(big.Float).SetPrec(prec).SetInt64(1)
	term := new(big.Float).SetPrec(prec).SetInt64(1)

	threshold := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), -int(prec))

	maxIter := int(prec) + 100
	for i := int64(1); i < int64(maxIter); i++ {
		denom := new(big.Float).SetPrec(prec).SetInt64(i)
		term.Quo(term, denom)
		term.Mul(term, x)
		result.Add(result, term)
		if new(big.Float).Abs(term).Cmp(threshold) < 0 {
			break
		}
	}
	return result
}

// Log computes log(x) (natural logarithm) at arbitrary precision.
func Log(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}

	if x.Sign() <= 0 {
		panic("bigmath.Log: x must be positive")
	}

	if x.Cmp(one) == 0 {
		return new(big.Float).SetPrec(prec)
	}

	wprec := prec + 64

	// Write x = m * 2^e where 0.5 <= |m| < 1.0
	// ln(x) = ln(m) + e * ln(2)
	m2 := new(big.Float).SetPrec(wprec)
	exp2 := x.MantExp(m2)

	// Transform: t = (m2 - 1) / (m2 + 1)
	// ln(m2) = 2 * (t + t^3/3 + t^5/5 + ...)
	oneW := new(big.Float).SetPrec(wprec).SetInt64(1)
	m2Plus1 := new(big.Float).SetPrec(wprec).Add(m2, oneW)
	m2Minus1 := new(big.Float).SetPrec(wprec).Sub(m2, oneW)
	t := new(big.Float).SetPrec(wprec).Quo(m2Minus1, m2Plus1)
	t2 := new(big.Float).SetPrec(wprec).Mul(t, t)

	term := new(big.Float).SetPrec(wprec).Copy(t)
	sum := new(big.Float).SetPrec(wprec).Copy(t)
	threshold := new(big.Float).SetPrec(wprec).SetMantExp(
		new(big.Float).SetPrec(wprec).SetInt64(1), -int(wprec))

	maxIter := int(wprec) + 100
	for n := int64(3); n < int64(maxIter); n += 2 {
		term.Mul(term, t2)
		part := new(big.Float).SetPrec(wprec).Quo(term, new(big.Float).SetPrec(wprec).SetInt64(n))
		sum.Add(sum, part)
		if new(big.Float).Abs(part).Cmp(threshold) < 0 {
			break
		}
	}

	lnM := new(big.Float).SetPrec(wprec).Mul(new(big.Float).SetPrec(wprec).SetInt64(2), sum)

	if exp2 != 0 {
		ln2 := ln2Const(wprec)
		eTimesLn2 := new(big.Float).SetPrec(wprec).Mul(
			new(big.Float).SetPrec(wprec).SetInt64(int64(exp2)), ln2)
		lnM.Add(lnM, eTimesLn2)
	}

	return lnM.SetPrec(prec)
}

func ln2Const(prec uint) *big.Float {
	// ln(2) = 2 * sum_{n=0}^{inf} 1/((2n+1) * 3^(2n+1))
	// Compute 1/3 at full precision
	three := new(big.Float).SetPrec(prec).SetInt64(3)
	oneP := new(big.Float).SetPrec(prec).SetInt64(1)
	t := new(big.Float).SetPrec(prec).Quo(oneP, three)
	t2 := new(big.Float).SetPrec(prec).Mul(t, t)
	term := new(big.Float).SetPrec(prec).Copy(t)
	sum := new(big.Float).SetPrec(prec).Copy(t)
	threshold := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), -int(prec))

	maxIter := int(prec) + 100
	for n := int64(3); n < int64(maxIter); n += 2 {
		term.Mul(term, t2)
		part := new(big.Float).SetPrec(prec).Quo(term, new(big.Float).SetPrec(prec).SetInt64(n))
		sum.Add(sum, part)
		if new(big.Float).Abs(part).Cmp(threshold) < 0 {
			break
		}
	}
	return new(big.Float).SetPrec(prec).Mul(new(big.Float).SetPrec(prec).SetInt64(2), sum)
}

// Sin computes sin(x) at arbitrary precision.
func Sin(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}

	if x.Sign() == 0 {
		return new(big.Float).SetPrec(prec)
	}

	wprec := prec + 64
	reduced, negate := reduceTrig(x, wprec)
	result := sinTaylor(reduced, wprec)
	if negate {
		result.Neg(result)
	}
	return result.SetPrec(prec)
}

// Cos computes cos(x) at arbitrary precision.
func Cos(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}

	wprec := prec + 64

	// Reduce to [0, 2*pi)
	pi := piConst(wprec)
	twoPi := new(big.Float).SetPrec(wprec).Mul(new(big.Float).SetPrec(wprec).SetInt64(2), pi)
	xCopy := new(big.Float).SetPrec(wprec).Copy(x)
	q := new(big.Float).SetPrec(wprec)
	q.Quo(xCopy, twoPi)
	qInt, _ := q.Int(nil)
	qFloat := new(big.Float).SetPrec(wprec).SetInt(qInt)
	reduced := new(big.Float).SetPrec(wprec).Sub(xCopy, new(big.Float).SetPrec(wprec).Mul(qFloat, twoPi))
	if reduced.Sign() < 0 {
		reduced.Add(reduced, twoPi)
	}

	halfPi := new(big.Float).SetPrec(wprec).Mul(half, pi)

	if reduced.Cmp(halfPi) > 0 && reduced.Cmp(pi) <= 0 {
		reflected := new(big.Float).SetPrec(wprec).Sub(pi, reduced)
		return new(big.Float).Neg(cosTaylor(reflected, wprec)).SetPrec(prec)
	}
	threeHalfPi := new(big.Float).SetPrec(wprec).Mul(new(big.Float).SetPrec(wprec).SetFloat64(1.5), pi)
	if reduced.Cmp(pi) > 0 && reduced.Cmp(threeHalfPi) <= 0 {
		reflected := new(big.Float).SetPrec(wprec).Sub(reduced, pi)
		return new(big.Float).Neg(cosTaylor(reflected, wprec)).SetPrec(prec)
	}
	if reduced.Cmp(threeHalfPi) > 0 {
		reflected := new(big.Float).SetPrec(wprec).Sub(twoPi, reduced)
		return cosTaylor(reflected, wprec).SetPrec(prec)
	}
	return cosTaylor(reduced, wprec).SetPrec(prec)
}

func sinTaylor(x *big.Float, prec uint) *big.Float {
	x2 := new(big.Float).SetPrec(prec).Mul(x, x)
	negX2 := new(big.Float).SetPrec(prec).Neg(x2)
	term := new(big.Float).SetPrec(prec).Copy(x)
	sum := new(big.Float).SetPrec(prec).Copy(x)
	threshold := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), -int(prec))

	maxIter := int(prec) + 100
	for n := int64(1); n < int64(maxIter); n++ {
		k := 2 * n
		denom := new(big.Float).SetPrec(prec).SetInt64(k * (k + 1))
		term.Mul(term, negX2)
		term.Quo(term, denom)
		sum.Add(sum, term)
		if new(big.Float).Abs(term).Cmp(threshold) < 0 {
			break
		}
	}
	return sum
}

func cosTaylor(x *big.Float, prec uint) *big.Float {
	x2 := new(big.Float).SetPrec(prec).Mul(x, x)
	negX2 := new(big.Float).SetPrec(prec).Neg(x2)
	term := new(big.Float).SetPrec(prec).SetInt64(1)
	sum := new(big.Float).SetPrec(prec).SetInt64(1)
	threshold := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), -int(prec))

	maxIter := int(prec) + 100
	for n := int64(1); n < int64(maxIter); n++ {
		k := 2 * n
		denom := new(big.Float).SetPrec(prec).SetInt64(k * (k - 1))
		term.Mul(term, negX2)
		term.Quo(term, denom)
		sum.Add(sum, term)
		if new(big.Float).Abs(term).Cmp(threshold) < 0 {
			break
		}
	}
	return sum
}

func piConst(prec uint) *big.Float {
	return machinPi(prec)
}

func machinPi(prec uint) *big.Float {
	// pi/4 = 4*arctan(1/5) - arctan(1/239)
	// Compute 1/5 and 1/239 at full precision
	five := new(big.Float).SetPrec(prec).SetInt64(5)
	twoThirtyNine := new(big.Float).SetPrec(prec).SetInt64(239)
	oneP := new(big.Float).SetPrec(prec).SetInt64(1)

	a1 := arctanSeries(new(big.Float).SetPrec(prec).Quo(oneP, five), prec)
	a1.Mul(a1, new(big.Float).SetPrec(prec).SetInt64(4))

	a2 := arctanSeries(new(big.Float).SetPrec(prec).Quo(oneP, twoThirtyNine), prec)

	piOver4 := new(big.Float).SetPrec(prec).Sub(a1, a2)
	return piOver4.Mul(piOver4, new(big.Float).SetPrec(prec).SetInt64(4))
}

func arctanSeries(x *big.Float, prec uint) *big.Float {
	x2 := new(big.Float).SetPrec(prec).Mul(x, x)
	negX2 := new(big.Float).SetPrec(prec).Neg(x2)
	term := new(big.Float).SetPrec(prec).Copy(x)
	sum := new(big.Float).SetPrec(prec).Copy(x)
	threshold := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), -int(prec))

	maxIter := int(prec) + 100
	for n := int64(1); n < int64(maxIter); n++ {
		k := 2*n + 1
		term.Mul(term, negX2)
		part := new(big.Float).SetPrec(prec).Quo(term, new(big.Float).SetPrec(prec).SetInt64(k))
		sum.Add(sum, part)
		if new(big.Float).Abs(part).Cmp(threshold) < 0 {
			break
		}
	}
	return sum
}

func reduceTrig(x *big.Float, prec uint) (*big.Float, bool) {
	pi := piConst(prec)
	twoPi := new(big.Float).SetPrec(prec).Mul(new(big.Float).SetPrec(prec).SetInt64(2), pi)

	xCopy := new(big.Float).SetPrec(prec).Copy(x)
	q := new(big.Float).SetPrec(prec)
	q.Quo(xCopy, twoPi)
	qInt, _ := q.Int(nil)
	qFloat := new(big.Float).SetPrec(prec).SetInt(qInt)
	reduced := new(big.Float).SetPrec(prec).Sub(xCopy, new(big.Float).SetPrec(prec).Mul(qFloat, twoPi))

	negPi := new(big.Float).SetPrec(prec).Neg(pi)
	for reduced.Cmp(pi) > 0 {
		reduced.Sub(reduced, twoPi)
	}
	for reduced.Cmp(negPi) < 0 {
		reduced.Add(reduced, twoPi)
	}

	negate := false
	halfPi := new(big.Float).SetPrec(prec).Mul(half, pi)
	negHalfPi := new(big.Float).SetPrec(prec).Neg(halfPi)

	if reduced.Cmp(halfPi) > 0 {
		reduced.Sub(pi, reduced)
	} else if reduced.Cmp(negHalfPi) < 0 {
		reduced.Add(reduced, pi)
		negate = true
	}

	return reduced, negate
}

// Sqrt computes sqrt(x) at arbitrary precision.
func Sqrt(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}
	result := new(big.Float).SetPrec(prec)
	return result.Sqrt(x)
}

// Float64 converts a big.Float to float64.
func Float64(x *big.Float) float64 {
	f, _ := x.Float64()
	return f
}

// NewFloat creates a new big.Float from a float64 at default precision.
func NewFloat(x float64) *big.Float {
	return new(big.Float).SetPrec(Prec).SetFloat64(x)
}

// NewFloatFromString creates a new big.Float from a string at default precision.
func NewFloatFromString(s string) *big.Float {
	f, _, _ := new(big.Float).Parse(s, 10)
	return f.SetPrec(Prec)
}
