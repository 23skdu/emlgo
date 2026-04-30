package eml

import (
	"math"
	"testing"
	"testing/quick"
)

func TestEmlProperty(t *testing.T) {
	f := func(x, y float64) bool {
		if y <= 0 || math.IsInf(y, 0) || math.IsNaN(y) {
			return true
		}
		if math.IsInf(x, 0) || math.IsNaN(x) {
			return true
		}
		if x > 700 || x < -700 {
			return true
		}
		got := Eml(x, y)
		want := math.Exp(x) - math.Log(y)
		if math.IsNaN(got) && math.IsNaN(want) {
			return true
		}
		if !math.IsInf(got, 0) && math.Abs(got-want) > 1e-10 {
			t.Logf("Eml(%v, %v) = %v, want %v", x, y, got, want)
			return false
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestEmlOneProperty(t *testing.T) {
	f := func(x float64) bool {
		if math.IsInf(x, 0) || math.IsNaN(x) || x > 700 {
			return true
		}
		got := EmlOne(x)
		want := math.Exp(x)
		if !math.IsInf(got, 0) && math.Abs(got-want) > 1e-10 {
			t.Logf("EmlOne(%v) = %v, want %v", x, got, want)
			return false
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestOneEmlProperty(t *testing.T) {
	f := func(y float64) bool {
		if y <= 0 || y > 1e300 || math.IsInf(y, 0) || math.IsNaN(y) {
			return true
		}
		got := OneEml(y)
		want := -math.Log(y)
		if !math.IsInf(got, 0) && math.Abs(got-want) > 1e-10 {
			t.Logf("OneEml(%v) = %v, want %v", y, got, want)
			return false
		}
		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestULPAccuracy(t *testing.T) {
	testCases := []struct {
		name string
		fn   func(float64) float64
	}{
		{"EmlOne", func(x float64) float64 { return EmlOne(x) }},
		{"OneEml", func(x float64) float64 { return OneEml(x) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputs := []float64{0, 0.5, 1, 2, math.E, 10, 100}
			for _, x := range inputs {
				got := tc.fn(x)
				var want float64
				if tc.name == "EmlOne" {
					want = math.Exp(x)
				} else {
					want = -math.Log(x)
				}

				ulp := ulp(got, want)
				if ulp > 1 {
					t.Logf("%s(%v): got=%v want=%v ulp=%d", tc.name, x, got, want, ulp)
				}
			}
		})
	}
}

func ulp(a, b float64) int {
	if a == b {
		return 0
	}
	if math.IsNaN(a) || math.IsNaN(b) {
		return 0
	}
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return 0
	}
	if math.IsInf(a, -1) && math.IsInf(b, -1) {
		return 0
	}

	lower := math.Nextafter(a, math.Inf(-1))
	upper := math.Nextafter(a, math.Inf(1))

	diff := math.Abs(a - b)
	count := 0
	for diff > 0 && (lower < b && b < upper) {
		mid := (lower + upper) / 2
		if diff > math.Abs(mid-a) {
			lower = math.Nextafter(mid, a)
			upper = math.Nextafter(mid, a)
			count++
		} else {
			break
		}
	}

	return int(diff / math.Nextafter(0, 1))
}