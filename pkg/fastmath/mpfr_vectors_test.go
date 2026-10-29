package fastmath

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/emlgo/eml/internal/accuracyvectors"
	"github.com/emlgo/eml/pkg/arithmetic"
	"github.com/emlgo/eml/pkg/hyper"
	"github.com/emlgo/eml/pkg/trig"
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
		t.Run(fmt.Sprintf("%s/%s", row[0], row[1]), func(t *testing.T) {
			if strings.Contains(row[1], ":") {
				parts := strings.SplitN(row[1], ":", 2)
				xBits, err := strconv.ParseUint(parts[0], 16, 64)
				if err != nil {
					t.Fatalf("parse x bits %q: %v", parts[0], err)
				}
				yBits, err := strconv.ParseUint(parts[1], 16, 64)
				if err != nil {
					t.Fatalf("parse y bits %q: %v", parts[1], err)
				}
				wantBits, err := strconv.ParseUint(row[2], 16, 64)
				if err != nil {
					t.Fatalf("parse reference bits %q: %v", row[2], err)
				}
				switch row[0] {
				case "eml64":
					x := math.Float64frombits(xBits)
					y := math.Float64frombits(yBits)
					got, want := FastEml(x, y), math.Float64frombits(wantBits)
					if math.Abs(got-want) > 2e-6*math.Max(1.0, math.Abs(want)) {
						t.Errorf("FastEml(%g, %g) = %.17g, want %.17g", x, y, got, want)
					}
				case "pow64":
					x := math.Float64frombits(xBits)
					y := math.Float64frombits(yBits)
					got, want := arithmetic.Pow(x, y), math.Float64frombits(wantBits)
					if ulpDistance64(got, want) > 4 {
						t.Errorf("Pow(%g, %g) = %.17g, want %.17g", x, y, got, want)
					}
				case "eml32":
					x := math.Float32frombits(uint32(xBits))                              // #nosec G115
					y := math.Float32frombits(uint32(yBits))                              // #nosec G115
					got, want := FastEmlF32(x, y), math.Float32frombits(uint32(wantBits)) // #nosec G115
					if math.Abs(float64(got-want)) > 5e-4*math.Max(1.0, math.Abs(float64(want))) {
						t.Errorf("FastEmlF32(%g, %g) = %.9g, want %.9g", x, y, got, want)
					}
				}
				return
			}

			inputBits, err := strconv.ParseUint(row[1], 16, 64)
			if err != nil {
				t.Fatalf("parse input bits %q: %v", row[1], err)
			}
			wantBits, err := strconv.ParseUint(row[2], 16, 64)
			if err != nil {
				t.Fatalf("parse reference bits %q: %v", row[2], err)
			}
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
			case "sin32":
				x := math.Float32frombits(uint32(inputBits))
				got, want := SinF32(x), math.Float32frombits(uint32(wantBits))
				if ulpDistance32(got, want) > 2 {
					t.Errorf("SinF32(%g) = %.9g, want %.9g", x, got, want)
				}
			case "cos32":
				x := math.Float32frombits(uint32(inputBits))
				got, want := CosF32(x), math.Float32frombits(uint32(wantBits))
				if ulpDistance32(got, want) > 2 {
					t.Errorf("CosF32(%g) = %.9g, want %.9g", x, got, want)
				}
			case "sqrt32":
				x := math.Float32frombits(uint32(inputBits))
				got, want := SqrtF32(x), math.Float32frombits(uint32(wantBits))
				if ulpDistance32(got, want) > 2 {
					t.Errorf("SqrtF32(%g) = %.9g, want %.9g", x, got, want)
				}
			default:
				x := math.Float64frombits(inputBits)
				want := math.Float64frombits(wantBits)
				var got float64
				switch row[0] {
				case "expm1_64":
					got = arithmetic.Expm1(x)
				case "log1p_64":
					got = arithmetic.Log1p(x)
				case "sin64":
					got = trig.Sin(x)
				case "cos64":
					got = trig.Cos(x)
				case "tan64":
					got = trig.Tan(x)
				case "sinh64":
					got = hyper.Sinh(x)
				case "cosh64":
					got = hyper.Cosh(x)
				case "tanh64":
					got = hyper.Tanh(x)
				case "asinh64":
					got = hyper.Asinh(x)
				case "acosh64":
					got = hyper.Acosh(x)
				case "atanh64":
					got = hyper.Atanh(x)
				case "sqrt64":
					got = arithmetic.Sqrt(x)
				case "cbrt64":
					got = arithmetic.Cbrt(x)
				default:
					t.Fatalf("unknown MPFR operation %q", row[0])
				}
				if ulpDistance64(got, want) > 8 {
					t.Errorf("%s(%g) = %.17g, want %.17g (%d ULP)", row[0], x, got, want, ulpDistance64(got, want))
				}
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
