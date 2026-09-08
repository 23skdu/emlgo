package bigmath

import (
	"math"
	"math/big"
	"testing"
)

func TestEml_OneOne(t *testing.T) {
	// Eml(1, 1) = exp(1) - log(1) = e - 0 = e
	x := NewFloat(1)
	y := NewFloat(1)
	result := Eml(x, y)

	expected := math.E
	got := Float64(result)
	if math.Abs(got-expected) > 1e-10 {
		t.Errorf("Eml(1,1) = %v, want %v (diff %e)", got, expected, got-expected)
	}
}

func TestExp(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 1},
		{1, math.E},
		{2, math.Exp(2)},
		{-1, math.Exp(-1)},
		{10, math.Exp(10)},
		{-10, math.Exp(-10)},
	}
	for _, tt := range tests {
		result := Exp(NewFloat(tt.input))
		got := Float64(result)
		relErr := math.Abs(got-tt.expected) / math.Abs(tt.expected)
		if relErr > 1e-10 {
			t.Errorf("Exp(%v) = %v, want %v (relErr %e)", tt.input, got, tt.expected, relErr)
		}
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 0},
		{math.E, 1},
		{2, math.Log(2)},
		{0.5, math.Log(0.5)},
		{100, math.Log(100)},
	}
	for _, tt := range tests {
		result := Log(NewFloat(tt.input))
		got := Float64(result)
		relErr := math.Abs(got-tt.expected) / math.Max(math.Abs(tt.expected), 1e-15)
		if relErr > 1e-10 {
			t.Errorf("Log(%v) = %v, want %v (relErr %e)", tt.input, got, tt.expected, relErr)
		}
	}
}

func TestSin(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{math.Pi / 6, 0.5},
		{math.Pi / 4, math.Sqrt(2) / 2},
		{math.Pi / 2, 1},
		{math.Pi, 0},
		{-math.Pi / 2, -1},
		{3 * math.Pi / 2, -1},
	}
	for _, tt := range tests {
		result := Sin(NewFloat(tt.input))
		got := Float64(result)
		if math.Abs(got-tt.expected) > 1e-10 {
			t.Errorf("Sin(%v) = %v, want %v (diff %e)", tt.input, got, tt.expected, got-tt.expected)
		}
	}
}

func TestCos(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 1},
		{math.Pi / 3, 0.5},
		{math.Pi / 4, math.Sqrt(2) / 2},
		{math.Pi / 2, 0},
		{math.Pi, -1},
		{-math.Pi, -1},
	}
	for _, tt := range tests {
		result := Cos(NewFloat(tt.input))
		got := Float64(result)
		if math.Abs(got-tt.expected) > 1e-10 {
			t.Errorf("Cos(%v) = %v, want %v (diff %e)", tt.input, got, tt.expected, got-tt.expected)
		}
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{1, 1},
		{4, 2},
		{2, math.Sqrt(2)},
		{0.25, 0.5},
		{1e10, 1e5},
		{1e-10, 1e-5},
	}
	for _, tt := range tests {
		result := Sqrt(NewFloat(tt.input))
		got := Float64(result)
		relErr := math.Abs(got-tt.expected) / math.Max(math.Abs(tt.expected), 1e-15)
		if relErr > 1e-10 {
			t.Errorf("Sqrt(%v) = %v, want %v (relErr %e)", tt.input, got, tt.expected, relErr)
		}
	}
}

func TestTrigIdentity_SinCosSquared(t *testing.T) {
	// sin²(z) + cos²(z) = 1 at high precision
	highPrec := uint(512)
	testPoints := []float64{0.1, 0.5, 1.0, 1.5, 2.0, 3.0, math.Pi / 4, math.Pi / 3}

	for _, pt := range testPoints {
		bx := new(big.Float).SetPrec(highPrec).SetFloat64(pt)
		s := Sin(bx)
		c := Cos(bx)

		s2 := new(big.Float).SetPrec(highPrec).Mul(s, s)
		c2 := new(big.Float).SetPrec(highPrec).Mul(c, c)
		sum := new(big.Float).SetPrec(highPrec).Add(s2, c2)

		diff := new(big.Float).SetPrec(highPrec).Sub(sum, one)
		diff.Abs(diff)

		threshold := new(big.Float).SetPrec(highPrec).SetFloat64(1e-50)
		if diff.Cmp(threshold) > 0 {
			t.Errorf("sin²(%v) + cos²(%v) = %v, want 1 (diff %v)", pt, pt, Float64(sum), Float64(diff))
		}
	}
}

func TestEulerIdentity(t *testing.T) {
	// e^(i*pi) + 1 = 0
	// Using real/imag parts:
	// Real part: cos(pi) = -1, so e^(i*pi) = -1
	// Imag part: sin(pi) = 0
	prec := uint(512)
	bigPi := piConst(prec)

	cosResult := Cos(bigPi)
	sinResult := Sin(bigPi)

	// cos(pi) ≈ -1
	cosPlus1 := new(big.Float).Add(cosResult, one)
	cosPlus1.Abs(cosPlus1)
	threshold := new(big.Float).SetPrec(prec).SetFloat64(1e-50)
	if cosPlus1.Cmp(threshold) > 0 {
		t.Errorf("cos(pi) = %v, want -1 (diff from -1: %v)", Float64(cosResult), Float64(cosPlus1))
	}

	// sin(pi) ≈ 0
	sinAbs := new(big.Float).Abs(sinResult)
	if sinAbs.Cmp(threshold) > 0 {
		t.Errorf("sin(pi) = %v, want 0 (diff: %v)", Float64(sinResult), Float64(sinAbs))
	}
}

func TestEmlVsFloat64(t *testing.T) {
	// Compare big.Float Eml result with float64 result
	prec := uint(512)
	testCases := []struct {
		x, y float64
	}{
		{1, 1},
		{2, 3},
		{0.5, 2.0},
		{-1, 5},
	}

	for _, tc := range testCases {
		bx := new(big.Float).SetPrec(prec).SetFloat64(tc.x)
		by := new(big.Float).SetPrec(prec).SetFloat64(tc.y)

		bigResult := Eml(bx, by)
		floatResult := math.Exp(tc.x) - math.Log(tc.y)

		// Relative error
		bigF64 := Float64(bigResult)
		relErr := math.Abs(bigF64-floatResult) / math.Max(math.Abs(floatResult), 1e-15)
		if relErr > 1e-10 {
			t.Errorf("Eml(%v, %v): big=%v, float64=%v, relErr=%e",
				tc.x, tc.y, bigF64, floatResult, relErr)
		}
	}
}

func TestNewFloatFromString(t *testing.T) {
	f := NewFloatFromString("3.14159265358979323846")
	got := Float64(f)
	expected := math.Pi
	if math.Abs(got-expected) > 1e-15 {
		t.Errorf("NewFloatFromString(pi) = %v, want %v", got, expected)
	}
}
