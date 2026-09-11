package quant

import (
	"encoding/binary"
	"math"
	"math/rand"
	"testing"
)

func TestUnpack4Bit(t *testing.T) {
	packed := []byte{0x21, 0x43, 0x65}
	dst := make([]byte, 6)
	Unpack4Bit(packed, dst, 6)

	expected := []byte{1, 2, 3, 4, 5, 6}
	for i := range expected {
		if dst[i] != expected[i] {
			t.Errorf("Unpack4Bit[%d] = %d, want %d", i, dst[i], expected[i])
		}
	}

	// Odd count
	dstOdd := make([]byte, 5)
	Unpack4Bit(packed, dstOdd, 5)
	for i := 0; i < 5; i++ {
		if dstOdd[i] != expected[i] {
			t.Errorf("Unpack4Bit odd[%d] = %d, want %d", i, dstOdd[i], expected[i])
		}
	}
}

func TestPolarTransformBatch(t *testing.T) {
	src := []float32{3.0, 4.0, 1.0, 1.0}
	dstRadii := make([]float32, 2)
	dstAngles := make([]float32, 2)

	PolarTransformBatch(src, dstRadii, dstAngles)

	if math.Abs(float64(dstRadii[0]-5.0)) > 1e-5 {
		t.Errorf("Radius[0] = %f, want 5.0", dstRadii[0])
	}
	expectedTheta0 := float32(math.Atan2(4.0, 3.0))
	if math.Abs(float64(dstAngles[0]-expectedTheta0)) > 1e-5 {
		t.Errorf("Angle[0] = %f, want %f", dstAngles[0], expectedTheta0)
	}
}

func TestTurboQuant4Distance(t *testing.T) {
	dim := 128
	pow2 := 128
	angleCount := pow2 - 1
	angleBytes := (angleCount*4 + 7) / 8
	qjlBytes := (pow2 + 7) / 8

	tqData := make([]byte, 4+angleBytes+qjlBytes)
	binary.LittleEndian.PutUint32(tqData[0:4], math.Float32bits(1.0))

	// Fill with valid data
	for i := 4; i < len(tqData); i++ {
		tqData[i] = byte(i % 256)
	}

	query := make([]float32, dim)
	for i := range query {
		query[i] = 0.5
	}

	dist, err := TurboQuant4Distance(query, tqData, dim, pow2)
	if err != nil {
		t.Fatalf("TurboQuant4Distance failed: %v", err)
	}
	if dist <= 0 {
		t.Errorf("Expected positive distance, got %f", dist)
	}
}

func TestPackUnpackRoundtrip(t *testing.T) {
	orig := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 7}
	packed := make([]byte, (len(orig)+1)/2)
	Pack4Bit(orig, packed)

	unpacked := make([]byte, len(orig))
	Unpack4Bit(packed, unpacked, len(orig))

	for i := range orig {
		if unpacked[i] != orig[i] {
			t.Errorf("Mismatch at %d: unpacked=%d, orig=%d", i, unpacked[i], orig[i])
		}
	}
}

func TestEncodeTurboQuant4(t *testing.T) {
	dim := 16
	pow2 := 16
	vec := make([]float32, dim)
	for i := range vec {
		vec[i] = float32(i + 1)
	}

	encoded, err := EncodeTurboQuant4(vec, pow2)
	if err != nil {
		t.Fatalf("EncodeTurboQuant4 failed: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatal("EncodeTurboQuant4 returned empty slice")
	}

	// Distance from vector to its own encoded form should be very small
	dist, err := TurboQuant4Distance(vec, encoded, dim, pow2)
	if err != nil {
		t.Fatalf("TurboQuant4Distance on encoded failed: %v", err)
	}
	if dist < 0 || math.IsNaN(float64(dist)) {
		t.Errorf("Invalid distance: %v", dist)
	}

	// Error cases
	_, err = EncodeTurboQuant4(vec, 15) // not pow2
	if err == nil {
		t.Errorf("Expected error for non-pow2")
	}
	_, err = EncodeTurboQuant4(nil, 16)
	if err == nil {
		t.Errorf("Expected error for empty vec")
	}
	_, err = EncodeTurboQuant4(make([]float32, 32), 16)
	if err == nil {
		t.Errorf("Expected error for oversized vec")
	}

	// Distance error cases
	_, err = TurboQuant4Distance(vec, []byte{1, 2}, dim, pow2)
	if err == nil {
		t.Errorf("Expected error for short tqData")
	}
	_, err = TurboQuant4DistanceScratch(vec, []byte{1, 2, 3, 4}, dim, pow2, nil, nil)
	if err == nil {
		t.Errorf("Expected error for corrupted tqData")
	}
	// Test scratch reallocation
	_, _ = TurboQuant4DistanceScratch(vec, encoded, dim, pow2, make([]float32, 1), make([]byte, 1))
}

func BenchmarkTurboQuant4Distance(b *testing.B) {
	dim := 128
	pow2 := 128
	angleCount := pow2 - 1
	angleBytes := (angleCount*4 + 7) / 8
	qjlBytes := (pow2 + 7) / 8

	tqData := make([]byte, 4+angleBytes+qjlBytes)
	binary.LittleEndian.PutUint32(tqData[0:4], math.Float32bits(1.0))

	rng := rand.New(rand.NewSource(42)) // #nosec G404 - deterministic benchmark test data
	for i := 4; i < len(tqData); i++ {
		tqData[i] = byte(rng.Intn(256)) // #nosec G115 - value in [0, 255] fits in byte
	}

	query := make([]float32, dim)
	for i := range query {
		query[i] = rng.Float32()
	}

	for b.Loop() {
		_, _ = TurboQuant4Distance(query, tqData, dim, pow2)
	}
}
