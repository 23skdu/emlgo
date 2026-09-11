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
		t.Fatalf("TurboQuant4Distance returned unexpected error: %v", err)
	}
	if dist <= 0 || math.IsNaN(float64(dist)) {
		t.Errorf("TurboQuant4Distance invalid distance: %v", dist)
	}
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
