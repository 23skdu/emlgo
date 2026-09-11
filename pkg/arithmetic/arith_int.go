package arithmetic

import (
	"math"
)

// DotProductInt8 computes the dot product of two int8 slices.
// Optimized with 4-way unrolling to maximize CPU instruction-level parallelism.
func DotProductInt8(a, b []int8) int32 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var sum0, sum1, sum2, sum3 int32
	i := 0
	for ; i+15 < n; i += 16 {
		sum0 += int32(a[i])*int32(b[i]) + int32(a[i+1])*int32(b[i+1]) + int32(a[i+2])*int32(b[i+2]) + int32(a[i+3])*int32(b[i+3])
		sum1 += int32(a[i+4])*int32(b[i+4]) + int32(a[i+5])*int32(b[i+5]) + int32(a[i+6])*int32(b[i+6]) + int32(a[i+7])*int32(b[i+7])
		sum2 += int32(a[i+8])*int32(b[i+8]) + int32(a[i+9])*int32(b[i+9]) + int32(a[i+10])*int32(b[i+10]) + int32(a[i+11])*int32(b[i+11])
		sum3 += int32(a[i+12])*int32(b[i+12]) + int32(a[i+13])*int32(b[i+13]) + int32(a[i+14])*int32(b[i+14]) + int32(a[i+15])*int32(b[i+15])
	}
	sum := sum0 + sum1 + sum2 + sum3
	for ; i < n; i++ {
		sum += int32(a[i]) * int32(b[i])
	}
	return sum
}

// DotProductUint8 computes the dot product of two uint8 slices.
func DotProductUint8(a, b []uint8) uint32 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var sum0, sum1, sum2, sum3 uint32
	i := 0
	for ; i+15 < n; i += 16 {
		sum0 += uint32(a[i])*uint32(b[i]) + uint32(a[i+1])*uint32(b[i+1]) + uint32(a[i+2])*uint32(b[i+2]) + uint32(a[i+3])*uint32(b[i+3])
		sum1 += uint32(a[i+4])*uint32(b[i+4]) + uint32(a[i+5])*uint32(b[i+5]) + uint32(a[i+6])*uint32(b[i+6]) + uint32(a[i+7])*uint32(b[i+7])
		sum2 += uint32(a[i+8])*uint32(b[i+8]) + uint32(a[i+9])*uint32(b[i+9]) + uint32(a[i+10])*uint32(b[i+10]) + uint32(a[i+11])*uint32(b[i+11])
		sum3 += uint32(a[i+12])*uint32(b[i+12]) + uint32(a[i+13])*uint32(b[i+13]) + uint32(a[i+14])*uint32(b[i+14]) + uint32(a[i+15])*uint32(b[i+15])
	}
	sum := sum0 + sum1 + sum2 + sum3
	for ; i < n; i++ {
		sum += uint32(a[i]) * uint32(b[i])
	}
	return sum
}

// L2SquaredInt8 computes the squared Euclidean distance between two int8 slices.
func L2SquaredInt8(a, b []int8) int32 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var sum0, sum1, sum2, sum3 int32
	i := 0
	for ; i+7 < n; i += 8 {
		d0 := int32(a[i]) - int32(b[i])
		d1 := int32(a[i+1]) - int32(b[i+1])
		d2 := int32(a[i+2]) - int32(b[i+2])
		d3 := int32(a[i+3]) - int32(b[i+3])
		d4 := int32(a[i+4]) - int32(b[i+4])
		d5 := int32(a[i+5]) - int32(b[i+5])
		d6 := int32(a[i+6]) - int32(b[i+6])
		d7 := int32(a[i+7]) - int32(b[i+7])
		sum0 += d0*d0 + d1*d1
		sum1 += d2*d2 + d3*d3
		sum2 += d4*d4 + d5*d5
		sum3 += d6*d6 + d7*d7
	}
	sum := sum0 + sum1 + sum2 + sum3
	for ; i < n; i++ {
		d := int32(a[i]) - int32(b[i])
		sum += d * d
	}
	return sum
}

// EuclideanDistanceInt8 computes the Euclidean (L2) distance between two int8 slices.
func EuclideanDistanceInt8(a, b []int8) float32 {
	return float32(math.Sqrt(float64(L2SquaredInt8(a, b))))
}

// CosineDistanceInt8 computes the cosine distance (1 - cos(theta)) between two int8 slices.
// Fuses dot product and norm accumulation into a single pass and uses intrinsic Sqrt to eliminate call overhead.
func CosineDistanceInt8(a, b []int8) float32 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n == 0 {
		return 1.0
	}
	var dot, normA, normB int64
	i := 0
	for ; i+7 < n; i += 8 {
		aChunk := a[i : i+8]
		bChunk := b[i : i+8]
		// #nosec G602 - bounds verified by loop condition (i+7 < n)
		a0, a1, a2, a3 := int64(aChunk[0]), int64(aChunk[1]), int64(aChunk[2]), int64(aChunk[3])
		a4, a5, a6, a7 := int64(aChunk[4]), int64(aChunk[5]), int64(aChunk[6]), int64(aChunk[7])
		b0, b1, b2, b3 := int64(bChunk[0]), int64(bChunk[1]), int64(bChunk[2]), int64(bChunk[3])
		b4, b5, b6, b7 := int64(bChunk[4]), int64(bChunk[5]), int64(bChunk[6]), int64(bChunk[7])

		dot += a0*b0 + a1*b1 + a2*b2 + a3*b3 + a4*b4 + a5*b5 + a6*b6 + a7*b7
		normA += a0*a0 + a1*a1 + a2*a2 + a3*a3 + a4*a4 + a5*a5 + a6*a6 + a7*a7
		normB += b0*b0 + b1*b1 + b2*b2 + b3*b3 + b4*b4 + b5*b5 + b6*b6 + b7*b7
	}
	for ; i < n; i++ {
		ai := int64(a[i])
		bi := int64(b[i])
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}
	if normA <= 0 || normB <= 0 {
		return 1.0
	}
	// math.Sqrt compiles directly to SQRTSD intrinsic
	similarity := float64(dot) / (math.Sqrt(float64(normA)) * math.Sqrt(float64(normB)))
	if similarity > 1.0 {
		similarity = 1.0
	} else if similarity < -1.0 {
		similarity = -1.0
	}
	return float32(1.0 - similarity)
}

// CosineDistanceUint8 computes the cosine distance (1 - cos(theta)) between two uint8 slices.
func CosineDistanceUint8(a, b []uint8) float32 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n == 0 {
		return 1.0
	}
	var dot, normA, normB uint64
	i := 0
	for ; i+7 < n; i += 8 {
		a0, a1, a2, a3 := uint64(a[i]), uint64(a[i+1]), uint64(a[i+2]), uint64(a[i+3])
		a4, a5, a6, a7 := uint64(a[i+4]), uint64(a[i+5]), uint64(a[i+6]), uint64(a[i+7])
		b0, b1, b2, b3 := uint64(b[i]), uint64(b[i+1]), uint64(b[i+2]), uint64(b[i+3])
		b4, b5, b6, b7 := uint64(b[i+4]), uint64(b[i+5]), uint64(b[i+6]), uint64(b[i+7])

		dot += a0*b0 + a1*b1 + a2*b2 + a3*b3 + a4*b4 + a5*b5 + a6*b6 + a7*b7
		normA += a0*a0 + a1*a1 + a2*a2 + a3*a3 + a4*a4 + a5*a5 + a6*a6 + a7*a7
		normB += b0*b0 + b1*b1 + b2*b2 + b3*b3 + b4*b4 + b5*b5 + b6*b6 + b7*b7
	}
	for ; i < n; i++ {
		ai := uint64(a[i])
		bi := uint64(b[i])
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}
	if normA == 0 || normB == 0 {
		return 1.0
	}
	similarity := float64(dot) / (math.Sqrt(float64(normA)) * math.Sqrt(float64(normB)))
	if similarity > 1.0 {
		similarity = 1.0
	} else if similarity < 0.0 {
		similarity = 0.0
	}
	return float32(1.0 - similarity)
}

// AddBatchInt8 adds two int8 slices element-wise with saturation.
func AddBatchInt8(a, b []int8) []int8 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	res := make([]int8, n)
	for i := 0; i < n; i++ {
		v := int32(a[i]) + int32(b[i])
		if v > math.MaxInt8 {
			v = math.MaxInt8
		} else if v < math.MinInt8 {
			v = math.MinInt8
		}
		res[i] = int8(v)
	}
	return res
}

// SubBatchInt8 subtracts two int8 slices element-wise with saturation.
func SubBatchInt8(a, b []int8) []int8 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	res := make([]int8, n)
	for i := 0; i < n; i++ {
		v := int32(a[i]) - int32(b[i])
		if v > math.MaxInt8 {
			v = math.MaxInt8
		} else if v < math.MinInt8 {
			v = math.MinInt8
		}
		res[i] = int8(v)
	}
	return res
}

// MulBatchInt8 multiplies two int8 slices element-wise with saturation.
func MulBatchInt8(a, b []int8) []int8 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	res := make([]int8, n)
	for i := 0; i < n; i++ {
		v := int32(a[i]) * int32(b[i])
		if v > math.MaxInt8 {
			v = math.MaxInt8
		} else if v < math.MinInt8 {
			v = math.MinInt8
		}
		res[i] = int8(v)
	}
	return res
}
