package constants

import (
	"math"

	"github.com/emlgo/eml/internal/eml"
)

const (
	One     float64 = 1.0
	E       float64 = math.E
	Pi      float64 = math.Pi
	NegOne  float64 = -1.0
	Two     float64 = 2.0
	Half    float64 = 0.5
	Sqrt2   float64 = 1.4142135623730951
	Sqrt3   float64 = 1.7320508075688772
	Ln2     float64 = 0.6931471805599453
	Ln10    float64 = 2.302585092994046
	SqrtPi  float64 = 1.7724538509055159
	Phi     float64 = 1.618033988749895
)

var I = complex(0, 1)

const (
	ComplexOne   complex128 = 1
	ComplexI     complex128 = complex(0, 1)
	ComplexNegI  complex128 = complex(0, -1)
)

func GenerateE() float64 {
	return E
}

func GeneratePi() float64 {
	return Pi
}

func GenerateI() complex128 {
	return I
}

func ExpOne() float64 {
	return E
}

func LogOne() float64 {
	return 0
}

func GenerateLn2() float64 {
	return Ln2
}

func GenerateLn10() float64 {
	return Ln10
}

func GenerateSqrt2() float64 {
	return Sqrt2
}

func GenerateSqrt3() float64 {
	return Sqrt3
}

func GenerateSqrtPi() float64 {
	return SqrtPi
}

func GeneratePhi() float64 {
	return Phi
}

func GenerateTwo() float64 {
	return Two
}

func GenerateHalf() float64 {
	return Half
}

func ComplexExp(z complex128) complex128 {
	return eml.EmlComplexOne(z)
}

func ComplexLog(z complex128) complex128 {
	if z == 0 {
		return complex(math.Inf(-1), 0)
	}
	arg := math.Atan2(imag(z), real(z))
	return complex(math.Log(math.Sqrt(real(z)*real(z)+imag(z)*imag(z))), arg)
}

func ComplexSin(z complex128) complex128 {
	if z == 0 {
		return 0
	}
	iz := complex(-imag(z), real(z))
	e := eml.EmlComplexOne(iz)
	neiz := eml.EmlComplexOne(-iz)
	return (e - neiz) / (2 * complex(0, 1))
}

func ComplexCos(z complex128) complex128 {
	conj := complex(-imag(z), real(z))
	exp_z := eml.EmlComplexOne(z)
	exp_conj := eml.EmlComplexOne(conj)
	return (exp_z + exp_conj) / 2
}