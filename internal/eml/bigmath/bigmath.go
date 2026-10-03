// Package bigmath provides arbitrary-precision EML operations using math/big.
// This backend is used for symbolic verification, identity checking, and
// deep tree reduction where float64 precision is insufficient.
package bigmath

import (
	"fmt"
	"math/big"
	"sync"
)

// seriesMaxIterations computes a safe upper bound on series iterations for a given precision in bits.
func seriesMaxIterations(prec uint) int {
	return int(prec) + 64
}

// Prec defines the default precision in bits for big.Float operations.
const Prec = 256

var (
	one      = new(big.Float).SetPrec(Prec).SetInt64(1)
	half     = new(big.Float).SetPrec(Prec).SetFloat64(0.5)
	piCache  sync.Map
	ln2Cache sync.Map
)

// Eml computes exp(x) - log(y) at arbitrary precision.
func Eml(x, y *big.Float) *big.Float {
	e := Exp(x)
	l := Log(y)
	return new(big.Float).Sub(e, l)
}

// Exp computes e^x at arbitrary precision.
func Exp(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}

	if x.Sign() == 0 {
		return new(big.Float).SetPrec(prec).SetInt64(1)
	}

	neg := x.Sign() < 0
	absX := new(big.Float).Abs(x)
	wprec := prec + 64

	// Range reduction: exp(x) = exp(x/2^k)^(2^k)
	// Choose k so |x/2^k| < 0.5 for fast Taylor convergence without converting to float64
	k := 0
	expMant := absX.MantExp(nil)
	if expMant > -1 {
		k = expMant + 2
	}

	shifted := new(big.Float).SetPrec(wprec)
	if k > 0 {
		divisor := new(big.Float).SetPrec(wprec).SetMantExp(
			new(big.Float).SetPrec(wprec).SetInt64(1), k)
		shifted.Quo(absX, divisor)
	} else {
		shifted.Copy(absX)
	}

	result := expTaylor(shifted, wprec)

	for i := 0; i < k; i++ {
		result.Mul(result, result)
	}

	if neg {
		oneP := new(big.Float).SetPrec(wprec).SetInt64(1)
		result.Quo(oneP, result)
	}

	return result.SetPrec(prec)
}

func expTaylor(x *big.Float, prec uint) *big.Float {
	result := new(big.Float).SetPrec(prec).SetInt64(1)
	term := new(big.Float).SetPrec(prec).SetInt64(1)

	threshold := new(big.Float).SetPrec(prec).SetMantExp(
		new(big.Float).SetPrec(prec).SetInt64(1), -int(prec))

	denom := new(big.Float).SetPrec(prec)
	absTerm := new(big.Float).SetPrec(prec)
	maxIter := seriesMaxIterations(prec)
	for i := int64(1); i < int64(maxIter); i++ {
		denom.SetInt64(i)
		term.Quo(term, denom)
		term.Mul(term, x)
		result.Add(result, term)
		if absTerm.Abs(term).Cmp(threshold) < 0 {
			break
		}
	}
	return result
}

// Log computes log(x) (natural logarithm) at arbitrary precision.
// Domain: x > 0.
// If x is zero, Log returns -Inf (matching math.Log).
// If x is negative, Log returns the NaN-equivalent sentinel 0 (matching Asin/Acos,
// as math/big.Float does not represent NaN values).
func Log(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}

	if x.Sign() == 0 {
		return new(big.Float).SetPrec(prec).SetInf(true)
	}
	if x.Sign() < 0 {
		return new(big.Float).SetPrec(prec)
	}

	if x.Cmp(one) == 0 {
		return new(big.Float).SetPrec(prec)
	}

	wprec := prec + 64

	// Write x = m * 2^e where 0.5 <= |m| < 1.0
	// ln(x) = ln(m) + e * ln(2)
	m2 := new(big.Float).SetPrec(wprec)
	exp2 := x.MantExp(m2)

	// If m2 < sqrt(0.5) ~ 0.7071, shift m2 = 2*m2 and exp2--
	// so m2 in [sqrt(0.5), sqrt(2)], keeping |t| <= 0.1716 for faster convergence.
	sqrtHalf := new(big.Float).SetPrec(wprec).SetFloat64(0.7071067811865475244)
	if m2.Cmp(sqrtHalf) < 0 {
		twoW := new(big.Float).SetPrec(wprec).SetInt64(2)
		m2.Mul(m2, twoW)
		exp2--
	}

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

	part := new(big.Float).SetPrec(wprec)
	nBig := new(big.Float).SetPrec(wprec)
	absPart := new(big.Float).SetPrec(wprec)
	maxIter := seriesMaxIterations(wprec)
	for n := int64(3); n < int64(maxIter); n += 2 {
		term.Mul(term, t2)
		nBig.SetInt64(n)
		part.Quo(term, nBig)
		sum.Add(sum, part)
		if absPart.Abs(part).Cmp(threshold) < 0 {
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
	if v, ok := ln2Cache.Load(prec); ok {
		if bf, ok := v.(*big.Float); ok {
			return new(big.Float).SetPrec(prec).Set(bf)
		}
	}
	res := computeLn2Const(prec)
	ln2Cache.Store(prec, new(big.Float).SetPrec(prec).Set(res))
	return res
}

func computeLn2Const(prec uint) *big.Float {
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

	part := new(big.Float).SetPrec(prec)
	nBig := new(big.Float).SetPrec(prec)
	absPart := new(big.Float).SetPrec(prec)
	maxIter := seriesMaxIterations(prec)
	for n := int64(3); n < int64(maxIter); n += 2 {
		term.Mul(term, t2)
		nBig.SetInt64(n)
		part.Quo(term, nBig)
		sum.Add(sum, part)
		if absPart.Abs(part).Cmp(threshold) < 0 {
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

	denom := new(big.Float).SetPrec(prec)
	absTerm := new(big.Float).SetPrec(prec)
	maxIter := seriesMaxIterations(prec)
	for n := int64(1); n < int64(maxIter); n++ {
		k := 2 * n
		denom.SetInt64(k * (k + 1))
		term.Mul(term, negX2)
		term.Quo(term, denom)
		sum.Add(sum, term)
		if absTerm.Abs(term).Cmp(threshold) < 0 {
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

	denom := new(big.Float).SetPrec(prec)
	absTerm := new(big.Float).SetPrec(prec)
	maxIter := seriesMaxIterations(prec)
	for n := int64(1); n < int64(maxIter); n++ {
		k := 2 * n
		denom.SetInt64(k * (k - 1))
		term.Mul(term, negX2)
		term.Quo(term, denom)
		sum.Add(sum, term)
		if absTerm.Abs(term).Cmp(threshold) < 0 {
			break
		}
	}
	return sum
}

func piConst(prec uint) *big.Float {
	if v, ok := piCache.Load(prec); ok {
		if bf, ok := v.(*big.Float); ok {
			return new(big.Float).SetPrec(prec).Set(bf)
		}
	}
	res := machinPi(prec)
	piCache.Store(prec, new(big.Float).SetPrec(prec).Set(res))
	return res
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

	part := new(big.Float).SetPrec(prec)
	kBig := new(big.Float).SetPrec(prec)
	absPart := new(big.Float).SetPrec(prec)
	maxIter := seriesMaxIterations(prec)
	for n := int64(1); n < int64(maxIter); n++ {
		k := 2*n + 1
		term.Mul(term, negX2)
		kBig.SetInt64(k)
		part.Quo(term, kBig)
		sum.Add(sum, part)
		if absPart.Abs(part).Cmp(threshold) < 0 {
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

// NewFloatFromStringChecked parses a string representation of a floating-point number
// using base 0 (auto-detecting decimal or 0b/0o/0x prefixes) at default precision (Prec).
// It returns an error if the string is empty, contains invalid characters, or has trailing garbage.
func NewFloatFromStringChecked(s string) (*big.Float, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("bigmath: empty input string")
	}
	z := new(big.Float).SetPrec(Prec)
	f, _, err := z.Parse(s, 0)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// NewFloatFromString creates a new big.Float from a string at default precision.
// Deprecated: Use NewFloatFromStringChecked to handle parsing errors safely.
// If parsing fails, NewFloatFromString returns nil without panicking.
func NewFloatFromString(s string) *big.Float {
	f, err := NewFloatFromStringChecked(s)
	if err != nil {
		return nil
	}
	return f
}

// NewFloatFromInt creates a new big.Float from an int64 at default precision.
func NewFloatFromInt(n int64) *big.Float {
	return new(big.Float).SetPrec(Prec).SetInt64(n)
}

var cosFn = Cos

// Tan computes tan(x) = Sin(x)/Cos(x) at arbitrary precision.
// Panics if cos(x) is zero (i.e., x is an odd multiple of π/2).
func Tan(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}
	s := Sin(new(big.Float).SetPrec(prec).Copy(x))
	c := cosFn(new(big.Float).SetPrec(prec).Copy(x))
	if c.Sign() == 0 {
		// cos(x) = 0: return ±Inf analogous to math.Tan
		inf := new(big.Float).SetPrec(prec).SetInf(s.Sign() >= 0)
		return inf
	}
	return new(big.Float).SetPrec(prec).Quo(s, c)
}

// Atan computes arctan(x) at arbitrary precision.
// Uses range reduction so that arctanSeries is only called on |x| < 0.5.
func Atan(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}
	wprec := prec + 128

	sign := x.Sign()
	if sign == 0 {
		return new(big.Float).SetPrec(prec).SetInt64(0)
	}

	// Work on |x|.
	ax := new(big.Float).SetPrec(wprec).Abs(x)
	wone := new(big.Float).SetPrec(wprec).SetInt64(1)

	// Reduce large arguments: atan(x) = π/2 - atan(1/x) for |x| > 1.
	useRecip := false
	if ax.Cmp(wone) > 0 {
		useRecip = true
		ax.Quo(wone, ax)
	}

	// Half-angle reduction: atan(x) = 2*atan(x / (1 + sqrt(1+x²)))
	// Repeat until |x| < 0.5 for fast series convergence.
	halvings := 0
	whalf := new(big.Float).SetPrec(wprec).SetFloat64(0.5)
	x2 := new(big.Float).SetPrec(wprec)
	inner := new(big.Float).SetPrec(wprec)
	sq := new(big.Float).SetPrec(wprec)
	denom := new(big.Float).SetPrec(wprec)
	for ax.Cmp(whalf) >= 0 {
		x2.Mul(ax, ax)
		inner.Add(wone, x2)
		sq.Sqrt(inner)
		denom.Add(wone, sq)
		ax.Quo(ax, denom)
		halvings++
	}

	result := arctanSeries(ax, wprec)

	// Undo halvings.
	twoW := new(big.Float).SetPrec(wprec).SetInt64(1)
	two := new(big.Float).SetPrec(wprec).SetInt64(2)
	for i := 0; i < halvings; i++ {
		twoW.Mul(twoW, two)
	}
	result.Mul(result, twoW)

	// Undo reciprocal: atan(x) = π/2 - result.
	if useRecip {
		pi := piConst(wprec)
		piOver2 := new(big.Float).SetPrec(wprec).Quo(pi, new(big.Float).SetPrec(wprec).SetInt64(2))
		result.Sub(piOver2, result)
	}

	result.SetPrec(prec)
	if sign < 0 {
		result.Neg(result)
	}
	return result
}

// Asin computes arcsin(x) at arbitrary precision.
// Domain: x ∈ [-1, 1]. Returns NaN-equivalent (0) for out-of-range inputs.
func Asin(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}
	wprec := prec + 64

	oneP := new(big.Float).SetPrec(wprec).SetInt64(1)
	negOneP := new(big.Float).SetPrec(wprec).SetInt64(-1)

	if x.Cmp(oneP) > 0 || x.Cmp(negOneP) < 0 {
		return new(big.Float).SetPrec(prec) // 0 for out-of-range
	}
	if x.Cmp(oneP) == 0 {
		pi := piConst(prec)
		return new(big.Float).SetPrec(prec).Quo(pi, new(big.Float).SetPrec(prec).SetInt64(2))
	}
	if x.Cmp(negOneP) == 0 {
		pi := piConst(prec)
		r := new(big.Float).SetPrec(prec).Quo(pi, new(big.Float).SetPrec(prec).SetInt64(2))
		return r.Neg(r)
	}

	// asin(x) = atan(x / sqrt(1 - x²))
	xw := new(big.Float).SetPrec(wprec).Copy(x)
	x2 := new(big.Float).SetPrec(wprec).Mul(xw, xw)
	wone := new(big.Float).SetPrec(wprec).SetInt64(1)
	denom := new(big.Float).SetPrec(wprec).Sqrt(new(big.Float).SetPrec(wprec).Sub(wone, x2))
	arg := new(big.Float).SetPrec(wprec).Quo(xw, denom)
	return Atan(arg.SetPrec(prec))
}

// Acos computes arccos(x) at arbitrary precision.
// Domain: x ∈ [-1, 1].
func Acos(x *big.Float) *big.Float {
	prec := x.Prec()
	if prec == 0 {
		prec = Prec
	}
	wprec := prec + 64
	// acos(x) = π/2 - asin(x)
	pi := piConst(wprec)
	piOver2 := new(big.Float).SetPrec(wprec).Quo(pi, new(big.Float).SetPrec(wprec).SetInt64(2))
	asinX := Asin(new(big.Float).SetPrec(wprec).Copy(x))
	return new(big.Float).SetPrec(prec).Sub(piOver2, asinX)
}
