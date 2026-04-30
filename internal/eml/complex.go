package eml

import "math/cmplx"

func EmlComplex(x, y complex128) complex128 {
	return cmplx.Exp(x) - cmplx.Log(y)
}

func EmlComplexOne(x complex128) complex128 {
	return EmlComplex(x, 1)
}

func OneEmlComplex(x complex128) complex128 {
	return EmlComplex(1, x)
}