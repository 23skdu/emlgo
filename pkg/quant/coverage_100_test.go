package quant

import (
	"testing"
)

func TestTurboQuant4EdgeCases(t *testing.T) {
	// 1. clampAngleBin bounds
	if clampAngleBin(-5) != 0 {
		t.Errorf("clampAngleBin(-5) should be 0")
	}
	if clampAngleBin(20) != 15 {
		t.Errorf("clampAngleBin(20) should be 15")
	}

	// 2. Scratch pool returning unexpected type
	oldGet := getTQScratch
	getTQScratch = func() any { return "unexpected_object" }
	_, err := TurboQuant4Distance([]float32{1, 2, 3, 4}, make([]byte, 10), 4, 4)
	getTQScratch = oldGet
	if err == nil {
		t.Errorf("expected error when scratch pool contains unexpected type")
	}

	// 3. Short query and short recon in TurboQuant4DistanceScratch
	vec := []float32{1, 2, 3, 4}
	enc, err := EncodeTurboQuant4(vec, 4)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	// Short query (len(query) < dim)
	shortQuery := []float32{1, 2}
	dist, err := TurboQuant4DistanceScratch(shortQuery, enc, 4, 4, nil, nil)
	if err != nil {
		t.Errorf("unexpected error on short query: %v", err)
	}
	_ = dist

	// Short recon (dim > pow2)
	query := []float32{1, 2, 3, 4, 5, 6, 7, 8}
	dist2, err := TurboQuant4DistanceScratch(query, enc, 8, 4, nil, nil)
	if err != nil {
		t.Errorf("unexpected error on dim > pow2: %v", err)
	}
	_ = dist2

	// Short qjlBits (len(qjlBits)*8 < limit)
	// For pow2=4, angleCount=3, angleBytes=(3*4+7)/8=2, header=4. Total prefix=6.
	if len(enc) >= 6 {
		truncatedTQData := enc[:6]
		_, _ = TurboQuant4DistanceScratch(shortQuery, truncatedTQData, 4, 4, nil, nil)
	}
}
