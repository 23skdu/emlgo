package eml

import "math/cmplx"

// Complex returns Exp(x) - Log(y) for complex128.
func Complex(x, y complex128) complex128 {
	return cmplx.Exp(x) - cmplx.Log(y)
}

// ComplexOne returns Complex(x, 1).
func ComplexOne(x complex128) complex128 {
	return Complex(x, 1)
}

// OneComplex returns Complex(1, x).
func OneComplex(x complex128) complex128 {
	return Complex(1, x)
}