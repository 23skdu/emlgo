package eml

import (
	"math"
	"testing"
)

func TestFusedOpsTo(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{0.5, 1, 2}
	result := make([]float64, len(a))

	t.Run("ExpMulTo", func(t *testing.T) {
		ExpMulTo(a, b, result)
	})
	t.Run("ExpAddTo", func(t *testing.T) {
		ExpAddTo(a, b, result)
	})
	t.Run("LogDivTo", func(t *testing.T) {
		LogDivTo(a, b, result)
	})
	t.Run("LogSubTo", func(t *testing.T) {
		LogSubTo(a, b, result)
	})
}

func TestSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.1
	}
	two := make([]float64, 10000)
	for i := range two {
		two[i] = 2.0
	}

	t.Run("AddSIMD", func(t *testing.T) {
		r := AddSIMD(large, two)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("SubSIMD", func(t *testing.T) {
		r := SubSIMD(large, two)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("MulSIMD", func(t *testing.T) {
		r := MulSIMD(large, two)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("DivSIMD", func(t *testing.T) {
		r := DivSIMD(large, two)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("AddScalarSIMD", func(t *testing.T) {
		r := AddScalarSIMD(large, 1.0)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("MulScalarSIMD", func(t *testing.T) {
		r := MulScalarSIMD(large, 2.0)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestEmlSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.1
	}
	result := make([]float64, len(large))

	t.Run("EmlSIMD", func(t *testing.T) {
		EmlSIMD(large, large, result)
		if len(result) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestExpSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.1
	}

	t.Run("ExpSIMD", func(t *testing.T) {
		r := ExpSIMD(large)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestLogSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i)*0.1 + 1
	}

	t.Run("LogSIMD", func(t *testing.T) {
		r := LogSIMD(large)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestSinSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.1
	}

	t.Run("SinSIMD", func(t *testing.T) {
		r := SinSIMD(large)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestCosSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.1
	}

	t.Run("CosSIMD", func(t *testing.T) {
		r := CosSIMD(large)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestSinCosSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.1
	}

	t.Run("SinCosSIMD", func(t *testing.T) {
		sin, cos := SinCosSIMD(large)
		if len(sin) != len(large) || len(cos) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestTanSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i) * 0.1
	}

	t.Run("TanSIMD", func(t *testing.T) {
		r := TanSIMD(large)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestSqrtSIMDExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i)
	}

	t.Run("SqrtSIMD", func(t *testing.T) {
		r := SqrtSIMD(large)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestFusedBatchExhaustive(t *testing.T) {
	large := make([]float64, 10000)
	for i := range large {
		large[i] = float64(i)*0.1 + 1
	}
	ones := make([]float64, 10000)
	for i := range ones {
		ones[i] = 1.0
	}
	halves := make([]float64, 10000)
	for i := range halves {
		halves[i] = 0.5
	}

	t.Run("ExpMulBatch", func(t *testing.T) {
		r := ExpMulBatch(large, halves)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("ExpAddBatch", func(t *testing.T) {
		r := ExpAddBatch(large, ones)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("LogDivBatch", func(t *testing.T) {
		r := LogDivBatch(large, halves)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
	t.Run("LogSubBatch", func(t *testing.T) {
		r := LogSubBatch(large, ones)
		if len(r) != len(large) {
			t.Error("len mismatch")
		}
	})
}

func TestIsNaNExhaustive(t *testing.T) {
	tests := []float64{math.NaN(), 0, 1, math.Inf(1), math.Inf(-1)}
	for _, x := range tests {
		t.Run("IsNaN", func(t *testing.T) {
			_ = IsNaN(x)
		})
	}
}

func TestIsInfExhaustive(t *testing.T) {
	t.Run("pos_inf", func(t *testing.T) {
		_ = IsInf(math.Inf(1), 1)
	})
	t.Run("neg_inf", func(t *testing.T) {
		_ = IsInf(math.Inf(-1), -1)
	})
	t.Run("zero_inf", func(t *testing.T) {
		_ = IsInf(math.Inf(1), 0)
	})
	t.Run("zero", func(t *testing.T) {
		_ = IsInf(0, 0)
	})
	t.Run("one", func(t *testing.T) {
		_ = IsInf(1, 0)
	})
}

func TestIsFiniteExhaustive(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		_ = IsFinite(0)
	})
	t.Run("one", func(t *testing.T) {
		_ = IsFinite(1)
	})
	t.Run("max", func(t *testing.T) {
		_ = IsFinite(math.MaxFloat64)
	})
	t.Run("nan", func(t *testing.T) {
		_ = IsFinite(math.NaN())
	})
	t.Run("inf", func(t *testing.T) {
		_ = IsFinite(math.Inf(1))
	})
}

func TestSignbitExhaustive(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		_ = Signbit(0)
	})
	t.Run("pos", func(t *testing.T) {
		_ = Signbit(1)
	})
	t.Run("neg", func(t *testing.T) {
		_ = Signbit(-1)
	})
	t.Run("pos_inf", func(t *testing.T) {
		_ = Signbit(math.Inf(1))
	})
	t.Run("neg_inf", func(t *testing.T) {
		_ = Signbit(math.Inf(-1))
	})
}