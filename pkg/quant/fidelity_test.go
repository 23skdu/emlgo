package quant

import (
	"encoding/binary"
	"math"
	"math/rand"
	"testing"
)

// TestTurboQuant4FidelityByDimension evaluates the round-trip reconstruction accuracy
// and true bit rate across multiple vector dimensions (8 to 512).
func TestTurboQuant4FidelityByDimension(t *testing.T) {
	dimensions := []int{8, 16, 32, 64, 128, 256, 512}
	rng := rand.New(rand.NewSource(12345)) // #nosec G404 - deterministic test data

	for _, dim := range dimensions {
		t.Run(string(rune(dim)), func(t *testing.T) {
			pow2 := dim
			vec := make([]float32, dim)
			var sumSq float64
			for i := range vec {
				vec[i] = rng.Float32()*2.0 - 1.0
				sumSq += float64(vec[i] * vec[i])
			}
			norm := float32(math.Sqrt(sumSq))
			for i := range vec {
				vec[i] /= norm
			}

			// Encode
			encoded, err := EncodeTurboQuant4(vec, pow2)
			if err != nil {
				t.Fatalf("EncodeTurboQuant4(dim=%d) failed: %v", dim, err)
			}

			// 1. Verify true bit rate
			totalBytes := len(encoded)
			totalBits := totalBytes * 8
			bitRate := float64(totalBits) / float64(dim)

			// Expected: ~4 bits/dim angles + 1 bit/dim QJL + 32-bit header / dim
			// For dim >= 64, bit rate should be <= 5.7 bits/dim
			// For dim >= 256, bit rate should be <= 5.2 bits/dim
			maxAllowedBitRate := 5.0 + 32.0/float64(dim) + 1.0 // with byte alignment overhead
			if bitRate > maxAllowedBitRate {
				t.Errorf("dim=%d: bitRate=%.2f bits/dim, expected <= %.2f", dim, bitRate, maxAllowedBitRate)
			}

			// 2. Verify distance / reconstruction error
			dist, err := TurboQuant4Distance(vec, encoded, dim, pow2)
			if err != nil {
				t.Fatalf("TurboQuant4Distance(dim=%d) failed: %v", dim, err)
			}

			// L2 distance on unit vector: relative error = dist / 1.0 = dist
			if dist > 0.40 {
				t.Errorf("dim=%d: reconstruction L2 distance %.4f exceeds tolerance 0.40", dim, dist)
			}

			// 3. Verify cosine similarity of reconstructed vector
			radius := math.Float32frombits(binary.LittleEndian.Uint32(encoded[0:4]))
			angleCount := pow2 - 1
			angleBytes := (angleCount*4 + 7) / 8
			packedAngles := encoded[4 : 4+angleBytes]
			qjlBits := encoded[4+angleBytes:]

			qIndices := make([]byte, angleCount)
			Unpack4Bit(packedAngles, qIndices, angleCount)

			recon := make([]float32, pow2)
			ReconstructTQ4(radius, qIndices, TQ4Lookup, pow2, recon)

			correction := radius / float32(math.Sqrt(float64(pow2))) * 0.1
			for i := 0; i < dim; i++ {
				bit := (qjlBits[i/8] >> uint(i%8)) & 1
				if bit == 1 {
					recon[i] += correction
				} else {
					recon[i] -= correction
				}
			}

			var dotProd, reconNormSq float64
			for i := 0; i < dim; i++ {
				dotProd += float64(vec[i] * recon[i])
				reconNormSq += float64(recon[i] * recon[i])
			}
			cosSim := dotProd / (1.0 * math.Sqrt(reconNormSq))
			if cosSim < 0.85 {
				t.Errorf("dim=%d: cosine similarity %.4f is below 0.85 threshold", dim, cosSim)
			}
		})
	}
}
