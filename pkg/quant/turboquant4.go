package quant

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"
)

// TQ4Lookup is the precomputed (cos, sin) table for 4-bit (16 bins) quantization.
var TQ4Lookup = func() []float32 {
	table := make([]float32, 16*2)
	for i := 0; i < 16; i++ {
		theta := (float64(i)/15.0)*2*math.Pi - math.Pi
		s, c := math.Sincos(theta)
		table[2*i] = float32(c)
		table[2*i+1] = float32(s)
	}
	return table
}()

var tqScratchPool = sync.Pool{
	New: func() any {
		return &tqScratch{
			recon:   make([]float32, 4096),
			indices: make([]byte, 4096),
		}
	},
}

type tqScratch struct {
	recon   []float32
	indices []byte
}

// Unpack4Bit unpacks 4-bit nibbles from packed into separate byte indices in dst.
// #nosec G602 - bounds verified by loop condition
func Unpack4Bit(packed []byte, dst []byte, count int) {
	numFullBytes := count / 2
	i := 0
	j := 0
	for ; i+7 < numFullBytes && j+15 < len(dst); i += 8 {
		p := packed[i : i+8]
		d := dst[j : j+16]
		// #nosec G602 - slice bounds checked by loop condition (i+7 < numFullBytes && j+15 < len(dst))
		b0, b1, b2, b3 := p[0], p[1], p[2], p[3]
		b4, b5, b6, b7 := p[4], p[5], p[6], p[7]

		// #nosec G602 - destination bounds checked by loop condition
		d[0] = b0 & 0x0F
		d[1] = b0 >> 4
		d[2] = b1 & 0x0F
		d[3] = b1 >> 4
		d[4] = b2 & 0x0F
		d[5] = b2 >> 4
		d[6] = b3 & 0x0F
		d[7] = b3 >> 4
		d[8] = b4 & 0x0F
		d[9] = b4 >> 4
		d[10] = b5 & 0x0F
		d[11] = b5 >> 4
		d[12] = b6 & 0x0F
		d[13] = b6 >> 4
		d[14] = b7 & 0x0F
		d[15] = b7 >> 4
		j += 16
	}
	for ; i < numFullBytes; i++ {
		b := packed[i]
		dst[j] = b & 0x0F
		dst[j+1] = b >> 4
		j += 2
	}
	if count%2 != 0 {
		dst[j] = packed[numFullBytes] & 0x0F
	}
}

// Pack4Bit packs separate 4-bit nibble indices (0-15) from src into contiguous bytes in dst.
func Pack4Bit(src []byte, dst []byte) {
	numFullPairs := len(src) / 2
	j := 0
	for i := 0; i < numFullPairs && j < len(dst); i++ {
		dst[j] = (src[2*i] & 0x0F) | ((src[2*i+1] & 0x0F) << 4)
		j++
	}
	if len(src)%2 != 0 && j < len(dst) {
		dst[j] = src[len(src)-1] & 0x0F
	}
}

// EncodeTurboQuant4 compresses a high-dimensional float32 vector into TurboQuant4 format.
// It computes recursive polar coordinate angles, quantizes them to 4-bit bins, and extracts QJL 1-bit residual signs.
func EncodeTurboQuant4(vec []float32, pow2 int) ([]byte, error) {
	if pow2 <= 1 || (pow2&(pow2-1)) != 0 {
		return nil, fmt.Errorf("quant: pow2 must be a power of two > 1, got %d", pow2)
	}
	if len(vec) == 0 || len(vec) > pow2 {
		return nil, fmt.Errorf("quant: vector dimension %d exceeds target pow2 %d", len(vec), pow2)
	}

	// Pad input to pow2
	padded := make([]float32, pow2)
	copy(padded, vec)

	angleCount := pow2 - 1
	qIndices := make([]byte, angleCount)

	// Recursive tree polar decomposition
	currentRadii := padded
	angleOffset := 0

	for currentLen := pow2; currentLen > 1; currentLen /= 2 {
		halfLen := currentLen / 2
		nextRadii := make([]float32, halfLen)
		for i := 0; i < halfLen; i++ {
			x := currentRadii[2*i]
			y := currentRadii[2*i+1]
			r := float32(math.Sqrt(float64(x*x + y*y)))
			theta := math.Atan2(float64(y), float64(x))

			bin := int(math.Round(((theta + math.Pi) / (2 * math.Pi)) * 15.0))
			if bin < 0 {
				bin = 0
			} else if bin > 15 {
				bin = 15
			}
			qIndices[angleOffset+i] = byte(bin)
			nextRadii[i] = r
		}
		angleOffset += halfLen
		currentRadii = nextRadii
	}

	radius := currentRadii[0]

	// Reconstruct to compute QJL residuals
	recon := make([]float32, pow2)
	ReconstructTQ4(radius, qIndices, TQ4Lookup, pow2, recon)

	// Compute QJL 1-bit residual mask
	qjlBytesCount := (pow2 + 7) / 8
	qjlBits := make([]byte, qjlBytesCount)
	for i := 0; i < len(vec); i++ {
		if vec[i] >= recon[i] {
			qjlBits[i/8] |= byte(1 << uint(i%8))
		}
	}

	// Pack 4-bit angles
	angleBytes := (angleCount*4 + 7) / 8
	packedAngles := make([]byte, angleBytes)
	Pack4Bit(qIndices, packedAngles)

	// Assemble final buffer: 4 bytes radius + packedAngles + qjlBits
	totalSize := 4 + angleBytes + qjlBytesCount
	out := make([]byte, totalSize)
	binary.LittleEndian.PutUint32(out[0:4], math.Float32bits(radius))
	copy(out[4:4+angleBytes], packedAngles)
	copy(out[4+angleBytes:], qjlBits)

	return out, nil
}

// PolarTransformBatch performs recursive 2D polar transformation on pairs of coordinates.
func PolarTransformBatch(src, dstRadii, dstAngles []float32) {
	n := len(src)
	halfN := n / 2
	for i := 0; i < halfN; i++ {
		x := src[2*i]
		y := src[2*i+1]
		dstRadii[i] = float32(math.Sqrt(float64(x*x + y*y)))
		dstAngles[i] = float32(math.Atan2(float64(y), float64(x)))
	}
}

// ReconstructTQ4 reconstructs vector coordinates from radius, angle indices, and lookup table.
func ReconstructTQ4(radius float32, qIndices []byte, lookupTable []float32, pow2 int, recon []float32) {
	recon[0] = radius
	angleCount := pow2 - 1
	currentLevelSize := 1
	angleOffset := angleCount
	for currentLevelSize < pow2 {
		angleOffset -= currentLevelSize
		for i := currentLevelSize - 1; i >= 0; i-- {
			r := recon[i]
			q := int(qIndices[angleOffset+i])
			c := lookupTable[2*q]
			s := lookupTable[2*q+1]
			recon[2*i] = r * c
			recon[2*i+1] = r * s
		}
		currentLevelSize *= 2
	}
}

// TurboQuant4Distance computes the L2 distance between a query vector and a TurboQuant4 encoded vector.
// Avoids allocations through a pooled scratch buffer.
func TurboQuant4Distance(query []float32, tqData []byte, dim int, pow2 int) (float32, error) {
	if len(tqData) < 4 {
		return 0, fmt.Errorf("quant: invalid tqData length %d", len(tqData))
	}
	buf := tqScratchPool.Get().(*tqScratch)
	defer tqScratchPool.Put(buf)

	return TurboQuant4DistanceScratch(query, tqData, dim, pow2, buf.recon, buf.indices)
}

// TurboQuant4DistanceScratch computes TurboQuant4 distance using caller-provided scratch buffers.
// Zero heap allocations.
func TurboQuant4DistanceScratch(query []float32, tqData []byte, dim int, pow2 int, recon []float32, qIndices []byte) (float32, error) {
	radius := math.Float32frombits(binary.LittleEndian.Uint32(tqData[0:4]))

	angleCount := pow2 - 1
	angleBytes := (angleCount*4 + 7) / 8
	if len(tqData) < 4+angleBytes {
		return 0, fmt.Errorf("quant: corrupted tqData (len=%d, expected>=%d)", len(tqData), 4+angleBytes)
	}

	packedAngles := tqData[4 : 4+angleBytes]
	qjlBits := tqData[4+angleBytes:]

	if cap(qIndices) < angleCount {
		qIndices = make([]byte, angleCount)
	}
	qIndices = qIndices[:angleCount]

	Unpack4Bit(packedAngles, qIndices, angleCount)

	if cap(recon) < pow2 {
		recon = make([]float32, pow2)
	}
	recon = recon[:pow2]

	ReconstructTQ4(radius, qIndices, TQ4Lookup, pow2, recon)

	// QJL residual correction
	correction := radius / float32(math.Sqrt(float64(pow2))) * 0.1
	limit := len(recon)
	if len(qjlBits)*8 < limit {
		limit = len(qjlBits) * 8
	}
	for i := 0; i < limit; i++ {
		if (qjlBits[i/8]>>uint(i%8))&1 != 0 {
			recon[i] += correction
		} else {
			recon[i] -= correction
		}
	}

	// L2 distance computation with 4-way loop unrolling for maximum throughput
	actualDim := dim
	if len(query) < actualDim {
		actualDim = len(query)
	}
	if len(recon) < actualDim {
		actualDim = len(recon)
	}

	var sum0, sum1, sum2, sum3 float32
	i := 0
	for ; i+7 < actualDim; i += 8 {
		d0 := query[i] - recon[i]
		d1 := query[i+1] - recon[i+1]
		d2 := query[i+2] - recon[i+2]
		d3 := query[i+3] - recon[i+3]
		d4 := query[i+4] - recon[i+4]
		d5 := query[i+5] - recon[i+5]
		d6 := query[i+6] - recon[i+6]
		d7 := query[i+7] - recon[i+7]

		sum0 += d0 * d0
		sum1 += d1 * d1
		sum2 += d2 * d2
		sum3 += d3 * d3
		sum0 += d4 * d4
		sum1 += d5 * d5
		sum2 += d6 * d6
		sum3 += d7 * d7
	}
	sum := sum0 + sum1 + sum2 + sum3
	for ; i < actualDim; i++ {
		d := query[i] - recon[i]
		sum += d * d
	}

	return float32(math.Sqrt(float64(sum))), nil
}
