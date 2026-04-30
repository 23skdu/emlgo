package constants

import (
	"math"
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      float64
		expected float64
		tol      float64
	}{
		{"One", One, 1.0, 0},
		{"E", E, math.E, 0},
		{"Pi", Pi, math.Pi, 0},
		{"NegOne", NegOne, -1.0, 0},
		{"Two", Two, 2.0, 0},
		{"Half", Half, 0.5, 0},
		{"Sqrt2", Sqrt2, math.Sqrt2, 0},
		{"Sqrt3", Sqrt3, math.Sqrt(3), 1e-10},
		{"Ln2", Ln2, math.Ln2, 0},
		{"Ln10", Ln10, math.Ln10, 0},
		{"SqrtPi", SqrtPi, math.Sqrt(math.Pi), 1e-10},
		{"Phi", Phi, 1.618033988749895, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if math.Abs(tt.got-tt.expected) > tt.tol {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestComplexConstants(t *testing.T) {
	if ComplexOne != 1 {
		t.Error("ComplexOne should be 1")
	}
	if ComplexI != complex(0, 1) {
		t.Error("ComplexI should be i")
	}
	if ComplexNegI != complex(0, -1) {
		t.Error("ComplexNegI should be -i")
	}
}

func TestGenerateE(t *testing.T) {
	got := GenerateE()
	want := math.E
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GenerateE() = %v, want %v", got, want)
	}
}

func TestGeneratePi(t *testing.T) {
	got := GeneratePi()
	want := math.Pi
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GeneratePi() = %v, want %v", got, want)
	}
}

func TestGenerateI(t *testing.T) {
	got := GenerateI()
	want := complex(0, 1)
	if got != want {
		t.Errorf("GenerateI() = %v, want %v", got, want)
	}
}

func TestExpOne(t *testing.T) {
	got := ExpOne()
	want := math.E
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("ExpOne() = %v, want %v", got, want)
	}
}

func TestLogOne(t *testing.T) {
	got := LogOne()
	if got != 0 {
		t.Errorf("LogOne() = %v, want 0", got)
	}
}

func TestGenerateLn2(t *testing.T) {
	got := GenerateLn2()
	want := math.Ln2
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GenerateLn2() = %v, want %v", got, want)
	}
}

func TestGenerateLn10(t *testing.T) {
	got := GenerateLn10()
	want := math.Ln10
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GenerateLn10() = %v, want %v", got, want)
	}
}

func TestGenerateSqrt2(t *testing.T) {
	got := GenerateSqrt2()
	want := math.Sqrt2
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GenerateSqrt2() = %v, want %v", got, want)
	}
}

func TestGenerateSqrt3(t *testing.T) {
	got := GenerateSqrt3()
	want := math.Sqrt(3)
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GenerateSqrt3() = %v, want %v", got, want)
	}
}

func TestGenerateSqrtPi(t *testing.T) {
	got := GenerateSqrtPi()
	want := math.Sqrt(math.Pi)
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GenerateSqrtPi() = %v, want %v", got, want)
	}
}

func TestGeneratePhi(t *testing.T) {
	got := GeneratePhi()
	want := 1.618033988749895
	if math.Abs(got-want) > 1e-10 {
		t.Errorf("GeneratePhi() = %v, want %v", got, want)
	}
}

func TestComplexExp(t *testing.T) {
	z := complex(1, 0)
	got := ComplexExp(z)
	want := complex(math.E, 0)
	if math.Abs(real(got)-real(want)) > 1e-10 || math.Abs(imag(got)-imag(want)) > 1e-10 {
		t.Errorf("ComplexExp(%v) = %v, want %v", z, got, want)
	}
}

func TestComplexLog(t *testing.T) {
	z := complex(math.E, 0)
	got := ComplexLog(z)
	want := complex(1, 0)
	if math.Abs(real(got)-real(want)) > 1e-10 || math.Abs(imag(got)-imag(want)) > 1e-10 {
		t.Errorf("ComplexLog(%v) = %v, want %v", z, got, want)
	}
}

func TestComplexSin(t *testing.T) {
	z := complex(0, 0)
	got := ComplexSin(z)
	want := complex(0, 0)
	if math.Abs(real(got)-real(want)) > 1e-10 || math.Abs(imag(got)-imag(want)) > 1e-10 {
		t.Errorf("ComplexSin(%v) = %v, want %v", z, got, want)
	}
}

func TestComplexCos(t *testing.T) {
	z := complex(0, 0)
	got := ComplexCos(z)
	want := complex(1, 0)
	if math.Abs(real(got)-real(want)) > 1e-10 || math.Abs(imag(got)-imag(want)) > 1e-10 {
		t.Errorf("ComplexCos(%v) = %v, want %v", z, got, want)
	}
}