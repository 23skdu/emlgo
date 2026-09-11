package eml

import "math/cmplx"

// Complex returns Exp(x) - Log(y) for complex128.
func Complex(x, y complex128) complex128 {
	return cmplx.Exp(x) - cmplx.Log(y)
}

// Complex64 returns Exp(x) - Log(y) for complex64.
func Complex64(x, y complex64) complex64 {
	return complex64(cmplx.Exp(complex128(x)) - cmplx.Log(complex128(y)))
}

// ComplexOne returns Complex(x, 1).
func ComplexOne(x complex128) complex128 {
	return Complex(x, 1)
}

// OneComplex returns Complex(1, x).
func OneComplex(x complex128) complex128 {
	return Complex(1, x)
}
