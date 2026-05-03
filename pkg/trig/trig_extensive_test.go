package trig

import (
	"math"
	"testing"
)

func TestSinComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pi", math.Pi, 0},
		{"pi/2", math.Pi / 2, 1},
		{"pi/4", math.Pi / 4, math.Sin(math.Pi / 4)},
		{"neg", -math.Pi, 0},
		{"nan", math.NaN(), math.NaN()},
		{"inf", math.Inf(1), math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sin(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Sin(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Sin(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCosComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"pi", math.Pi, -1},
		{"pi/2", math.Pi / 2, 0},
		{"pi/4", math.Pi / 4, math.Cos(math.Pi / 4)},
		{"neg", -math.Pi, -1},
		{"nan", math.NaN(), math.NaN()},
		{"inf", math.Inf(1), math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cos(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Cos(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Cos(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestTanComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pi/4", math.Pi / 4, 1},
		{"pi/3", math.Pi / 3, math.Sqrt(3)},
		{"neg", -math.Pi / 4, -1},
		{"nan", math.NaN(), math.NaN()},
		{"inf", math.Inf(1), math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tan(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Tan(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Tan(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCotComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pi/4", math.Pi / 4, 1},
		{"pi/2", math.Pi / 2, 0},
		{"zero", 0, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cot(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Cot(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Cot(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestSecComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"pi/3", math.Pi / 3, 2},
		{"pi/2", math.Pi / 2, 1.6331239353195392e+16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sec(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Sec(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Sec(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCscComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"pi/2", math.Pi / 2, 1},
		{"pi/6", math.Pi / 6, 2},
		{"zero", 0, math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Csc(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Csc(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Csc(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAsinComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.Pi / 2},
		{"neg_one", -1, -math.Pi / 2},
		{"half", 0.5, math.Pi / 6},
		{"invalid", 2, math.NaN()},
		{"nan", math.NaN(), math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asin(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Asin(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Asin(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAcosComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, math.Pi / 2},
		{"one", 1, 0},
		{"neg_one", -1, math.Pi},
		{"invalid", 2, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acos(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Acos(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Acos(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAtanComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.Pi / 4},
		{"neg_one", -1, -math.Pi / 4},
		{"inf", math.Inf(1), math.Pi / 2},
		{"neg_inf", math.Inf(-1), -math.Pi / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atan(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Atan(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAtan2Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		y, x     float64
		expected float64
	}{
		{"first", 1, 1, math.Pi / 4},
		{"neg_y", -1, 1, -math.Pi / 4},
		{"zero_x", 1, 0, math.Pi / 2},
		{"zero_y", 0, 1, 0},
		{"zero_zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atan2(tt.y, tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Atan2(%v, %v) = %v, want %v", tt.y, tt.x, got, tt.expected)
			}
		})
	}
}

func TestAcotComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, math.Pi / 4},
		{"zero", 0, math.Pi / 2},
		{"neg", -1, 3 * math.Pi / 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acot(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Acot(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAsecComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"two", 2, math.Acos(0.5)},
		{"half", 0.5, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asec(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Asec(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Asec(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAcscComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"two", 2, math.Asin(0.5)},
		{"half", 0.5, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acsc(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Acsc(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Acsc(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestSinhComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.Sinh(1)},
		{"neg", -1, math.Sinh(-1)},
		{"inf", math.Inf(1), math.Inf(1)},
		{"neg_inf", math.Inf(-1), math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sinh(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Sinh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCoshComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"one", 1, math.Cosh(1)},
		{"neg", -1, math.Cosh(-1)},
		{"inf", math.Inf(1), math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cosh(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Cosh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestTanhComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.Tanh(1)},
		{"neg", -1, math.Tanh(-1)},
		{"large_pos", 100, 1},
		{"large_neg", -100, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tanh(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Tanh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCothComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 1 / math.Tanh(1)},
		{"zero", 0, math.NaN()},
		{"large", 100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Coth(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Coth(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Coth(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestSechComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"one", 1, 1 / math.Cosh(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sech(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Sech(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestCschComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 1 / math.Sinh(1)},
		{"zero", 0, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Csch(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Csch(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Csch(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAsinhComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"one", 1, math.Asinh(1)},
		{"neg", -1, math.Asinh(-1)},
		{"inf", math.Inf(1), math.Inf(1)},
		{"neg_inf", math.Inf(-1), math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asinh(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Asinh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAcoshComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"two", 2, math.Acosh(2)},
		{"half", 0.5, math.NaN()},
		{"inf", math.Inf(1), math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acosh(tt.x)
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(got) {
					t.Errorf("Acosh(%v) = %v, want NaN", tt.x, got)
				}
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Acosh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAtanhComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"half", 0.5, math.Atanh(0.5)},
		{"neg_half", -0.5, math.Atanh(-0.5)},
		{"one", 1, math.NaN()},
		{"boundary", 0.999999, math.Atanh(0.999999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atanh(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Atanh(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Atanh(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAcothComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"two", 2, math.Atanh(0.5)},
		{"half", 0.5, math.NaN()},
		{"neg_half", -0.5, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acoth(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Acoth(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Acoth(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAsechComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, 0},
		{"half", 0.5, math.Acosh(2)},
		{"zero", 0, math.NaN()},
		{"neg", -1, math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asech(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("Asech(%v) = %v, want NaN", tt.x, got)
			} else if !close(got, tt.expected, 1e-10) {
				t.Errorf("Asech(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestAcschComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"one", 1, math.Asinh(1)},
		{"zero", 0, math.Inf(1)},
		{"inf", math.Inf(1), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acsch(tt.x)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("Acsch(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func TestDegToRadComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		deg      float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pi", 180, math.Pi},
		{"pi/2", 90, math.Pi / 2},
		{"half_pi", 45, math.Pi / 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DegToRad(tt.deg)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("DegToRad(%v) = %v, want %v", tt.deg, got, tt.expected)
			}
		})
	}
}

func TestRadToDegComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		rad      float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pi", math.Pi, 180},
		{"pi/2", math.Pi / 2, 90},
		{"pi/4", math.Pi / 4, 45},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RadToDeg(tt.rad)
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("RadToDeg(%v) = %v, want %v", tt.rad, got, tt.expected)
			}
		})
	}
}

func TestSinCosComprehensive(t *testing.T) {
	sin, cos := SinCos(0)
	if sin != 0 || cos != 1 {
		t.Errorf("SinCos(0) = (%v, %v), want (0, 1)", sin, cos)
	}

	sin, cos = SinCos(math.Pi / 2)
	if !close(sin, 1, 1e-10) || !close(cos, 0, 1e-10) {
		t.Errorf("SinCos(pi/2) = (%v, %v), want (1, 0)", sin, cos)
	}
}

func TestSinhCoshComprehensive(t *testing.T) {
	sinh, cosh := SinhCosh(0)
	if sinh != 0 || cosh != 1 {
		t.Errorf("SinhCosh(0) = (%v, %v), want (0, 1)", sinh, cosh)
	}

	sinh, cosh = SinhCosh(1)
	if !close(sinh, math.Sinh(1), 1e-10) || !close(cosh, math.Cosh(1), 1e-10) {
		t.Errorf("SinhCosh(1) = (%v, %v)", sinh, cosh)
	}
}

func close(a, b, tol float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return true
	}
	if math.IsInf(a, -1) && math.IsInf(b, -1) {
		return true
	}
	return math.Abs(a-b) <= tol
}

func TestBatchFunctionsComprehensive(t *testing.T) {
	testData := []float64{0, math.Pi / 4, math.Pi / 2, math.Pi, math.Pi * 1.5}
	testDataSmall := []float64{0, math.Pi / 4}

	t.Run("SinBatch", func(t *testing.T) {
		result := SinBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("SinBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("SinBatch_empty", func(t *testing.T) {
		result := SinBatch([]float64{})
		if result != nil && len(result) != 0 {
			t.Errorf("SinBatch empty = %v", result)
		}
	})
	t.Run("SinBatch_small", func(t *testing.T) {
		result := SinBatch(testDataSmall)
		if len(result) != len(testDataSmall) {
			t.Errorf("SinBatch small length = %v, want %v", len(result), len(testDataSmall))
		}
	})
	t.Run("CosBatch", func(t *testing.T) {
		result := CosBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("CosBatch length = %v, want %v", len(result), len(testData))
		}
	})
	t.Run("CosBatch_empty", func(t *testing.T) {
		result := CosBatch([]float64{})
		if result != nil && len(result) != 0 {
			t.Errorf("CosBatch empty = %v", result)
		}
	})
	t.Run("SinCosBatch", func(t *testing.T) {
		sin, cos := SinCosBatch(testData)
		if len(sin) != len(testData) || len(cos) != len(testData) {
			t.Errorf("SinCosBatch length mismatch")
		}
	})
	t.Run("SinCosBatch_empty", func(t *testing.T) {
		sin, cos := SinCosBatch([]float64{})
		if (sin != nil && len(sin) != 0) || (cos != nil && len(cos) != 0) {
			t.Errorf("SinCosBatch empty = %v, %v", sin, cos)
		}
	})
	t.Run("TanBatch", func(t *testing.T) {
		result := TanBatch(testData)
		if len(result) != len(testData) {
			t.Errorf("TanBatch length = %v, want %v", len(result), len(testData))
		}
	})
}

func TestFastFunctionsComprehensive(t *testing.T) {
	testsSin := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pi", math.Pi, 0},
		{"pi/2", math.Pi / 2, 1},
		{"nan", math.NaN(), math.NaN()},
		{"inf", math.Inf(1), math.NaN()},
	}
	for _, tt := range testsSin {
		t.Run("SinFast/"+tt.name, func(t *testing.T) {
			got := SinFast(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("SinFast(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("SinFast(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}

	testsCos := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 1},
		{"pi", math.Pi, -1},
		{"pi/2", math.Pi / 2, 0},
		{"nan", math.NaN(), math.NaN()},
		{"inf", math.Inf(1), math.NaN()},
	}
	for _, tt := range testsCos {
		t.Run("CosFast/"+tt.name, func(t *testing.T) {
			got := CosFast(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("CosFast(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("CosFast(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}

	testsTan := []struct {
		name     string
		x        float64
		expected float64
	}{
		{"zero", 0, 0},
		{"pi/4", math.Pi / 4, 1},
		{"nan", math.NaN(), math.NaN()},
	}
	for _, tt := range testsTan {
		t.Run("TanFast/"+tt.name, func(t *testing.T) {
			got := TanFast(tt.x)
			if math.IsNaN(tt.expected) && !math.IsNaN(got) {
				t.Errorf("TanFast(%v) = %v, want NaN", tt.x, got)
			}
			if math.IsNaN(tt.expected) && math.IsNaN(got) {
				return
			}
			if !close(got, tt.expected, 1e-10) {
				t.Errorf("TanFast(%v) = %v, want %v", tt.x, got, tt.expected)
			}
		})
	}
}

func BenchmarkSin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sin(1.5)
	}
}

func BenchmarkCos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Cos(1.5)
	}
}

func BenchmarkTan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tan(1.5)
	}
}

func BenchmarkSinh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sinh(1.5)
	}
}

func BenchmarkCosh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Cosh(1.5)
	}
}

func BenchmarkTanh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tanh(1.5)
	}
}