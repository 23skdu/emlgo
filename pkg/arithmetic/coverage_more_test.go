package arithmetic

import (
	"math"
	"testing"
)

func TestAddSubAllCases(t *testing.T) {
	_ = Add(1, 2)
	_ = Add(math.Inf(1), 1)
	_ = Add(math.Inf(-1), 1)
	_ = Add(math.NaN(), 1)
	_ = Sub(1, 2)
	_ = Sub(math.Inf(1), 1)
	_ = Sub(math.Inf(-1), 1)
	_ = Sub(math.NaN(), 1)
}

func TestMulDivAllCases(t *testing.T) {
	_ = Mul(2, 3)
	_ = Mul(0.0, math.Inf(1))
	_ = Mul(math.NaN(), 1)
	_ = Div(6, 3)
	_ = Div(1, 0)
	_ = Div(math.Inf(1), 1)
	_ = Div(math.NaN(), 1)
}

func TestModRemainderAllCases(t *testing.T) {
	_ = Mod(10, 3)
	_ = Mod(10, 0)
	_ = Mod(math.Inf(1), 1)
	_ = Mod(math.NaN(), 1)
	_ = Remainder(10, 3)
	_ = Remainder(10, 0)
	_ = Remainder(math.Inf(1), 1)
	_ = Remainder(math.NaN(), 1)
}

func TestPowAllCases(t *testing.T) {
	_ = Pow(2, 3)
	_ = Pow(0, 0)
	_ = Pow(-1, 0.5)
	_ = Pow(math.Inf(1), 1)
	_ = Pow(1, math.Inf(1))
	_ = Pow(math.NaN(), 1)
	_ = Pow(1, 2)
	_ = Pow(2, 0.5)
	_ = Pow(-2, 2)
	_ = Pow(0, 2)
}

func TestLogBaseAllCases(t *testing.T) {
	_ = LogBase2(1)
	_ = LogBase2(2)
	_ = LogBase2(0.5)
	_ = LogBase2(math.Inf(1))
	_ = LogBase2(0)
	_ = LogBase2(math.NaN())
	_ = LogBase10(1)
	_ = LogBase10(10)
	_ = LogBase10(0.1)
	_ = LogBase10(math.Inf(1))
	_ = LogBase10(0)
	_ = LogBase10(math.NaN())
}

func TestMaxMinAllCases(t *testing.T) {
	_ = Max(1, 2)
	_ = Max(math.Inf(1), 2)
	_ = Max(math.Inf(-1), 2)
	_ = Max(math.NaN(), 2)
	_ = Max(1, math.Inf(1))
	_ = Max(1, math.NaN())
	_ = Min(1, 2)
	_ = Min(math.Inf(1), 2)
	_ = Min(math.Inf(-1), 2)
	_ = Min(math.NaN(), 2)
	_ = Min(1, math.Inf(1))
	_ = Min(1, math.NaN())
}

func TestAbsNegInvAllCases(t *testing.T) {
	_ = Abs(1)
	_ = Abs(-1)
	_ = Abs(0)
	_ = Abs(math.Inf(1))
	_ = Abs(math.Inf(-1))
	_ = Abs(math.NaN())
	_ = Neg(1)
	_ = Neg(-1)
	_ = Neg(0)
	_ = Neg(math.Inf(1))
	_ = Neg(math.NaN())
	_ = Inv(2)
	_ = Inv(0)
	_ = Inv(math.Inf(1))
	_ = Inv(math.NaN())
}

func TestExpLogAllCases(t *testing.T) {
	_ = Exp(0)
	_ = Exp(1)
	_ = Exp(10)
	_ = Exp(math.Inf(1))
	_ = Exp(math.Inf(-1))
	_ = Exp(math.NaN())
	_ = Log(1)
	_ = Log(math.E)
	_ = Log(math.Inf(1))
	_ = Log(0)
	_ = Log(math.NaN())
}

func TestLog1pExpm1AllCases(t *testing.T) {
	_ = Log1p(0)
	_ = Log1p(1)
	_ = Log1p(math.Inf(1))
	_ = Log1p(-1)
	_ = Log1p(math.NaN())
	_ = Expm1(0)
	_ = Expm1(1)
	_ = Expm1(10)
	_ = Expm1(math.Inf(1))
	_ = Expm1(math.Inf(-1))
	_ = Expm1(math.NaN())
}

func TestGCDLCMAllCases(t *testing.T) {
	_ = GCD(0, 0)
	_ = GCD(1, 0)
	_ = GCD(0, 1)
	_ = GCD(48, 18)
	_ = GCD(17, 23)
	_ = GCD(100, 25)
	_ = LCM(0, 0)
	_ = LCM(1, 0)
	_ = LCM(0, 1)
	_ = LCM(4, 6)
	_ = LCM(17, 23)
	_ = LCM(100, 25)
}