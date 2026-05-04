package trig

import (
	"math"
	"testing"
)

func TestTrigEdgeCases(t *testing.T) {
	t.Run("TanFast_ZeroCos", func(t *testing.T) {
		TanFast(math.Pi / 2)
	})
}
