package eml

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/emlgo/eml/internal/accuracyvectors"
)

func ulpDistance32(a, b float32) uint32 {
	x, y := orderedFloat32Bits(math.Float32bits(a)), orderedFloat32Bits(math.Float32bits(b))
	if x > y {
		return x - y
	}
	return y - x
}

func orderedFloat32Bits(bits uint32) uint32 {
	if bits&(uint32(1)<<31) != 0 {
		return ^bits
	}
	return bits | (uint32(1) << 31)
}

func TestMPFRBatchGoldenVectors(t *testing.T) {
	rows, err := accuracyvectors.Rows()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if len(row) != 4 {
			t.Fatalf("invalid MPFR vector row: %v", row)
		}
		op := row[0]
		t.Run(op+"/"+row[1], func(t *testing.T) {
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
				switch op {
				case "eml64":
					x := math.Float64frombits(xBits)
					y := math.Float64frombits(yBits)
					want := math.Float64frombits(wantBits)
					got := Eml(x, y)
					if distance := ulpDistance(got, want); distance > 4 {
						t.Fatalf("eml64(%g, %g) = %.17g, want %.17g (%d ULP)", x, y, got, want, distance)
					}
				case "pow64":
					x := math.Float64frombits(xBits)
					y := math.Float64frombits(yBits)
					want := math.Float64frombits(wantBits)
					got := math.Pow(x, y)
					if distance := ulpDistance(got, want); distance > 2 {
						t.Fatalf("pow64(%g, %g) = %.17g, want %.17g (%d ULP)", x, y, got, want, distance)
					}
				case "eml32":
					x := math.Float32frombits(uint32(xBits))       // #nosec G115
					y := math.Float32frombits(uint32(yBits))       // #nosec G115
					want := math.Float32frombits(uint32(wantBits)) // #nosec G115
					got := ExpSIMDF32([]float32{x})[0] - LogSIMDF32([]float32{y})[0]
					if distance := ulpDistance32(got, want); distance > 2 {
						t.Fatalf("eml32(%g, %g) = %.9g, want %.9g (%d ULP)", x, y, got, want, distance)
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

			if strings.HasSuffix(op, "32") {
				x := math.Float32frombits(uint32(inputBits))   // #nosec G115
				want := math.Float32frombits(uint32(wantBits)) // #nosec G115
				var got float32
				switch op {
				case "exp32":
					got = ExpSIMDF32([]float32{x})[0]
				case "log32":
					got = LogSIMDF32([]float32{x})[0]
				case "sin32":
					got = SinSIMDF32([]float32{x})[0]
				case "cos32":
					got = CosSIMDF32([]float32{x})[0]
				case "sqrt32":
					got = SqrtSIMDF32([]float32{x})[0]
				}
				if distance := ulpDistance32(got, want); distance > 2 {
					t.Fatalf("%s(%g) = %.9g, want %.9g (%d ULP)", op, x, got, want, distance)
				}
				return
			}

			x, want := math.Float64frombits(inputBits), math.Float64frombits(wantBits)
			var got float64
			var ulpBudget uint64 = 4
			switch op {
			case "exp64":
				got = ExpSIMD([]float64{x})[0]
				ulpBudget = budgetExpULP
			case "log64":
				got = LogSIMD([]float64{x})[0]
				ulpBudget = budgetLogULP
			case "expm1_64":
				got = Expm1(x)
			case "log1p_64":
				got = Log1p(x)
			case "sin64":
				got = SinSIMD([]float64{x})[0]
				ulpBudget = budgetSinULP
			case "cos64":
				got = CosSIMD([]float64{x})[0]
				ulpBudget = budgetCosULP
			case "tan64":
				got = TanSIMD([]float64{x})[0]
				ulpBudget = budgetTanULP
			case "sinh64":
				got = SinhBatch([]float64{x})[0]
			case "cosh64":
				got = CoshBatch([]float64{x})[0]
			case "tanh64":
				got = TanhBatch([]float64{x})[0]
			case "asinh64":
				got = AsinhBatch([]float64{x})[0]
			case "acosh64":
				got = AcoshBatch([]float64{x})[0]
			case "atanh64":
				got = AtanhBatch([]float64{x})[0]
			case "sqrt64":
				got = SqrtSIMD([]float64{x})[0]
			case "cbrt64":
				got = Cbrt(x)
			}
			if distance := ulpDistance(got, want); distance > ulpBudget {
				t.Fatalf("%s(%g) = %.17g, want %.17g (%d ULP > %d)",
					op, x, got, want, distance, ulpBudget)
			}
		})
	}
}
