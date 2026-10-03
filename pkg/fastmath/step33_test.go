package fastmath

import (
	"math"
	"testing"

	"github.com/emlgo/eml/internal/eml"
)

func TestStep33_FastLogMinimaxSubnormal(t *testing.T) {
	var worstFastErr float64
	var worstMathErr float64

	for k := -1074; k <= -1022; k++ {
		x := math.Ldexp(1.0, k)
		wantExact := float64(k) * math.Ln2
		oracleVal := eml.ExactLogOracle(x)
		fastVal := FastLogMinimax(x)
		mathVal := math.Log(x)

		// Oracle matches exact mathematical value
		if math.Abs(oracleVal-wantExact) > 1e-12 {
			t.Errorf("oracle(2^%d) = %v, want %v", k, oracleVal, wantExact)
		}

		// FastLogMinimax has 0 error on power-of-two subnormals
		fastErr := math.Abs(fastVal - wantExact)
		if fastErr > worstFastErr {
			worstFastErr = fastErr
		}
		if fastErr > 1e-12 {
			t.Errorf("FastLogMinimax(2^%d) = %v, want %v (err %.3g)", k, fastVal, wantExact, fastErr)
		}

		mathErr := math.Abs(mathVal - wantExact)
		if mathErr > worstMathErr {
			worstMathErr = mathErr
		}
	}

	t.Logf("FastLogMinimax subnormal sweep [-1074, -1022]: worst err = %.3g, math.Log worst err = %.2f nats",
		worstFastErr, worstMathErr)

	if worstFastErr > 1e-12 {
		t.Fatalf("FastLogMinimax subnormal worst error %.3g exceeds 1e-12 budget", worstFastErr)
	}
	if worstMathErr < 30.0 {
		t.Logf("Notice: math.Log error is %.2f", worstMathErr)
	}
}
