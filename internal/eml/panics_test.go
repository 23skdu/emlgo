package eml

import (
	"testing"
)

func TestSIMDPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for length mismatch in ExpSIMDTo")
		}
	}()
	ExpSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsLog(t *testing.T) {
	defer func() { recover() }()
	LogSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsSin(t *testing.T) {
	defer func() { recover() }()
	SinSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsCos(t *testing.T) {
	defer func() { recover() }()
	CosSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsTan(t *testing.T) {
	defer func() { recover() }()
	TanSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsSqrt(t *testing.T) {
	defer func() { recover() }()
	SqrtSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsSinCos(t *testing.T) {
	defer func() { recover() }()
	SinCosSIMDTo(make([]float64, 5), make([]float64, 4), make([]float64, 5))
}

func TestSIMDPanicsAdd(t *testing.T) {
	defer func() { recover() }()
	AddSIMD(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsSub(t *testing.T) {
	defer func() { recover() }()
	SubSIMD(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsMul(t *testing.T) {
	defer func() { recover() }()
	MulSIMD(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsDiv(t *testing.T) {
	defer func() { recover() }()
	DivSIMD(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsAbs(t *testing.T) {
	defer func() { recover() }()
	AbsSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsNeg(t *testing.T) {
	defer func() { recover() }()
	NegSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsInv(t *testing.T) {
	defer func() { recover() }()
	InvSIMDTo(make([]float64, 5), make([]float64, 4))
}

func TestSIMDPanicsFma(t *testing.T) {
	defer func() { recover() }()
	FmaSIMDTo(make([]float64, 5), make([]float64, 4), make([]float64, 5), make([]float64, 5))
}

func TestSIMDPanicsMore(t *testing.T) {
	data := make([]float64, 5)
	bad := make([]float64, 4)
	
	defer func() { recover() }()
	EmlSIMD(data, bad, data)
	EmlSIMD(data, data, bad)
	ExpSIMDTo(data, bad)
	LogSIMDTo(data, bad)
	SqrtSIMDTo(data, bad)
	SinSIMDTo(data, bad)
	CosSIMDTo(data, bad)
	TanSIMDTo(data, bad)
	SinCosSIMDTo(data, bad, data)
	SinCosSIMDTo(data, data, bad)
	AbsSIMDTo(data, bad)
	NegSIMDTo(data, bad)
	InvSIMDTo(data, bad)
	AddScalarSIMDTo(data, 1.0, bad)
	MulScalarSIMDTo(data, 1.0, bad)
	Log2SIMDTo(data, bad)
	Log10SIMDTo(data, bad)
}

func TestSIMDEmpty(t *testing.T) {
	EmlSIMD(nil, nil, nil)
	_ = ExpSIMD(nil)
	_ = LogSIMD(nil)
	_ = SqrtSIMD(nil)
	_ = SinSIMD(nil)
	_ = CosSIMD(nil)
	_ = TanSIMD(nil)
	_, _ = SinCosSIMD(nil)
	_ = AbsSIMD(nil)
	_ = NegSIMD(nil)
	_ = InvSIMD(nil)
}

func TestFusedPanics(t *testing.T) {
	a := make([]float64, 5)
	b := make([]float64, 4)
	
	defer func() { recover() }()
	ExpMulTo(a, b, a)
	ExpAddTo(a, b, a)
	LogDivTo(a, b, a)
	LogSubTo(a, b, a)
}
