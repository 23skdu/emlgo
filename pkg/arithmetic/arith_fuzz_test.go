package arithmetic

import (
	"testing"
)

func FuzzAbs(f *testing.F) {
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(0.0)
	f.Fuzz(func(t *testing.T, a float64) {
		Abs(a)
	})
}

func FuzzAdd(f *testing.F) {
	f.Add(1.0, 2.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		Add(a, b)
	})
}

func FuzzMod(f *testing.F) {
	f.Add(10.0, 3.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		Mod(a, b)
	})
}

func FuzzRemainder(f *testing.F) {
	f.Add(10.0, 3.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		Remainder(a, b)
	})
}

func FuzzPow(f *testing.F) {
	f.Add(2.0, 3.0)
	f.Fuzz(func(t *testing.T, a, b float64) {
		Pow(a, b)
	})
}

func FuzzLogBase(f *testing.F) {
	f.Add(8.0)
	f.Fuzz(func(t *testing.T, a float64) {
		LogBase2(a)
		LogBase10(a)
	})
}
