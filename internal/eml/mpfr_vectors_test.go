package eml

import (
	"math"
	"strconv"
	"testing"

	"github.com/emlgo/eml/internal/accuracyvectors"
)

func TestMPFRBatchGoldenVectors(t *testing.T) {
	rows, err := accuracyvectors.Rows()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if len(row) != 4 || (row[0] != "exp64" && row[0] != "log64") {
			continue
		}
		inputBits, err := strconv.ParseUint(row[1], 16, 64)
		if err != nil {
			t.Fatalf("parse input bits %q: %v", row[1], err)
		}
		wantBits, err := strconv.ParseUint(row[2], 16, 64)
		if err != nil {
			t.Fatalf("parse reference bits %q: %v", row[2], err)
		}
		t.Run(row[0]+"/"+row[1], func(t *testing.T) {
			x, want := math.Float64frombits(inputBits), math.Float64frombits(wantBits)
			var got float64
			var ulpBudget uint64
			switch row[0] {
			case "exp64":
				got = ExpSIMD([]float64{x})[0]
				ulpBudget = budgetExpULP
			case "log64":
				got = LogSIMD([]float64{x})[0]
				ulpBudget = budgetLogULP
			}
			if math.IsNaN(want) || math.IsInf(want, 0) {
				if math.IsNaN(got) != math.IsNaN(want) || math.IsInf(got, 0) != math.IsInf(want, 0) {
					t.Fatalf("%s(%g) = %g, want %g", row[0], x, got, want)
				}
				return
			}
			if distance := ulpDistance(got, want); distance > ulpBudget {
				t.Fatalf("%s(%g) = %.17g, want %.17g (%d ULP > %d)",
					row[0], x, got, want, distance, ulpBudget)
			}
		})
	}
}
