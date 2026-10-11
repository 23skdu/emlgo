package fastmath

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/emlgo/eml/internal/accuracyvectors"
)

func TestMPFRGoldenVectors(t *testing.T) {
	rows, err := accuracyvectors.Rows()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if len(row) != 4 {
			t.Fatalf("invalid MPFR vector row: %v", row)
		}
		inputBits, err := strconv.ParseUint(row[1], 16, 64)
		if err != nil {
			t.Fatalf("parse input bits %q: %v", row[1], err)
		}
		wantBits, err := strconv.ParseUint(row[2], 16, 64)
		if err != nil {
			t.Fatalf("parse reference bits %q: %v", row[2], err)
		}
		t.Run(fmt.Sprintf("%s/%s", row[0], row[1]), func(t *testing.T) {
			switch row[0] {
			case "exp64":
				x := math.Float64frombits(inputBits)
				got, want := FastExpMinimax(x), math.Float64frombits(wantBits)
				if want == 0 {
					if got != 0 {
						t.Errorf("FastExpMinimax(%g) = %g, want 0", x, got)
					}
				} else if want < math.Float64frombits(0x0010000000000000) {
					if ulpDistance64(got, want) > 2 {
						t.Errorf("FastExpMinimax(%g) = %g, want %g (more than 2 ULP)", x, got, want)
					}
				} else if math.Abs(got-want) > 1.5e-6*math.Abs(want) {
					t.Errorf("FastExpMinimax(%g) = %g, want %g", x, got, want)
				}
			case "log64":
				x := math.Float64frombits(inputBits)
				got, want := FastLogMinimax(x), math.Float64frombits(wantBits)
				if math.Abs(got-want) > 1e-7 {
					t.Errorf("FastLogMinimax(%g) = %.17g, want %.17g", x, got, want)
				}
			case "exp32":
				x := math.Float32frombits(uint32(inputBits))
				got, want := FastExpF32(x), math.Float32frombits(uint32(wantBits))
				if want == 0 {
					if got != 0 {
						t.Errorf("FastExpF32(%g) = %g, want 0", x, got)
					}
				} else if want < math.Float32frombits(0x00800000) {
					if ulpDistance32(got, want) > 20 {
						t.Errorf("FastExpF32(%g) = %g, want %g (more than 20 ULP)", x, got, want)
					}
				} else if math.Abs(float64(got-want)) > 2e-4*math.Abs(float64(want)) {
					t.Errorf("FastExpF32(%g) = %g, want %g", x, got, want)
				}
			case "log32":
				x := math.Float32frombits(uint32(inputBits))
				got, want := FastLogF32(x), math.Float32frombits(uint32(wantBits))
				if math.Abs(float64(got-want)) > 1e-5 {
					t.Errorf("FastLogF32(%g) = %.9g, want %.9g", x, got, want)
				}
			default:
				t.Fatalf("unknown MPFR operation %q", row[0])
			}
		})
	}
}

func ulpDistance64(a, b float64) uint64 {
	x, y := orderedFloat64(math.Float64bits(a)), orderedFloat64(math.Float64bits(b))
	if x > y {
		return x - y
	}
	return y - x
}

func orderedFloat64(bits uint64) uint64 {
	if bits&(uint64(1)<<63) != 0 {
		return ^bits
	}
	return bits | (uint64(1) << 63)
}

func ulpDistance32(a, b float32) uint32 {
	x, y := orderedFloat32(math.Float32bits(a)), orderedFloat32(math.Float32bits(b))
	if x > y {
		return x - y
	}
	return y - x
}

func orderedFloat32(bits uint32) uint32 {
	if bits&(uint32(1)<<31) != 0 {
		return ^bits
	}
	return bits | (uint32(1) << 31)
}
