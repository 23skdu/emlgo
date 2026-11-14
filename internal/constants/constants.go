package constants

import (
	"math"

	"github.com/23skdu/emlgo/internal/eml"
)

const (
	// One is the float64 constant 1.
	One float64 = 1.0
	// E is Euler's number (math.E).
	E float64 = math.E
	// Pi is π to float64 precision (math.Pi).
	Pi float64 = math.Pi
	// NegOne is the float64 constant -1.
	NegOne float64 = -1.0
	// Two is the float64 constant 2.
	Two float64 = 2.0
	// Half is the float64 constant 0.5.
	Half float64 = 0.5
	// Sqrt2 is √2 to float64 precision.
	Sqrt2 float64 = 1.4142135623730951
	// Sqrt3 is √3 to float64 precision.
	Sqrt3 float64 = 1.7320508075688772
	// Ln2 is the natural logarithm of 2.
	Ln2 float64 = 0.693147180559945309417232121458
	// Ln10 is the natural logarithm of 10.
	Ln10 float64 = 2.302585092994046
	// SqrtPi is √π to float64 precision.
	SqrtPi float64 = 1.7724538509055159
	// Phi is the golden ratio (1+√5)/2.
	Phi float64 = 1.618033988749895
	// ExpOverflow is the largest argument that math.Exp can evaluate without overflowing to +Inf.
	ExpOverflow float64 = 709.782712893384
)

// I is the imaginary unit, 0+1i, as a complex128.
var I = complex(0, 1)

const (
	// ComplexOne is the complex128 constant 1.
	ComplexOne complex128 = 1
	// ComplexI is the imaginary unit as a complex128.
	ComplexI complex128 = complex(0, 1)
	// ComplexNegI is the negative imaginary unit -i as a complex128.
	ComplexNegI complex128 = complex(0, -1)
)

// GenerateE computes e from the EML identity eml.Eml(1, 1) = exp(1) - log(1) = e - 0.
func GenerateE() float64 {
	// e = eml(1, 1) = exp(1) - log(1) = e - 0
	return eml.Eml(One, One)
}

// GeneratePi computes π as 4*atan(1), evaluating atan through complex
// logarithms obtained from eml.Complex.
func GeneratePi() float64 {
	// pi = 4 * atan(1)
	// atan(x) = i/2 * (log(1-ix) - log(1+ix))
	// Using eml.EmlComplex(0, z) = 1 - log(z)
	// log(z) = 1 - eml.EmlComplex(0, z)
	i := complex(0, 1)
	log1minusI := 1 - eml.Complex(0, 1-i)
	log1plusI := 1 - eml.Complex(0, 1+i)
	atan1 := (i / 2) * (log1minusI - log1plusI)
	return 4 * real(atan1)
}

// GenerateI returns the imaginary unit as a complex128.
func GenerateI() complex128 {
	return ComplexI
}

// ExpOne returns e, the value of exp(1).
func ExpOne() float64 {
	return E
}

// LogOne returns 0, the natural logarithm of 1.
func LogOne() float64 {
	return 0
}

// GenerateLn2 returns Ln2, the natural logarithm of 2, as a precomputed constant.
func GenerateLn2() float64 {
	return Ln2
}

// GenerateLn10 returns Ln10, the natural logarithm of 10, as a precomputed constant.
func GenerateLn10() float64 {
	return Ln10
}

// GenerateSqrt2 returns Sqrt2, √2, as a precomputed constant.
func GenerateSqrt2() float64 {
	return Sqrt2
}

// GenerateSqrt3 returns Sqrt3, √3, as a precomputed constant.
func GenerateSqrt3() float64 {
	return Sqrt3
}

// GenerateSqrtPi returns SqrtPi, √π, as a precomputed constant.
func GenerateSqrtPi() float64 {
	return SqrtPi
}

// GeneratePhi returns Phi, the golden ratio, as a precomputed constant.
func GeneratePhi() float64 {
	return Phi
}

// GenerateTwo returns Two, the float64 constant 2.
func GenerateTwo() float64 {
	return Two
}

// GenerateHalf returns Half, the float64 constant 0.5.
func GenerateHalf() float64 {
	return Half
}

// ComplexExp returns the complex exponential of z, computed as
// eml.ComplexOne(z) = exp(z) - log(1).
func ComplexExp(z complex128) complex128 {
	return eml.ComplexOne(z)
}

// ComplexLog returns the principal complex logarithm of z, or (-Inf, 0) when z is zero.
func ComplexLog(z complex128) complex128 {
	if z == 0 {
		return complex(math.Inf(-1), 0)
	}
	arg := math.Atan2(imag(z), real(z))
	return complex(math.Log(math.Sqrt(real(z)*real(z)+imag(z)*imag(z))), arg)
}

// ComplexSin returns the complex sine of z, computed as
// (exp(iz) - exp(-iz)) / 2i with the exponentials from eml.ComplexOne.
func ComplexSin(z complex128) complex128 {
	if z == 0 {
		return 0
	}
	iz := complex(-imag(z), real(z))
	e := eml.ComplexOne(iz)
	neiz := eml.ComplexOne(-iz)
	return (e - neiz) / (2 * complex(0, 1))
}

// ComplexCos returns the complex cosine of z, computed as
// (exp(iz) + exp(-iz)) / 2 with the exponentials from eml.ComplexOne.
func ComplexCos(z complex128) complex128 {
	iz := complex(-imag(z), real(z))
	niz := complex(imag(z), -real(z))
	return (eml.ComplexOne(iz) + eml.ComplexOne(niz)) / 2
}
