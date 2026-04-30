package main

import (
	"fmt"
	"math"
	"time"
)

func FmaScalar(a, b, c float64) float64 {
	return a*b + c
}

func fastSin(x float64) float64 {
	// Simple range reduction for [0, pi/2]
	// This is NOT robust for large x, but good for local benchmarking
	
	x2 := x * x
	p := 0.008333333333333333 // 1/120
	p = FmaScalar(p, x2, -0.16666666666666666) // -1/6
	p = FmaScalar(p, x2, 1.0)
	return x * p
}

func main() {
	x := 0.5
	iterations := 10000000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = math.Sin(x)
	}
	fmt.Printf("math.Sin: %v\n", time.Since(start))

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = fastSin(x)
	}
	fmt.Printf("fastSin:  %v\n", time.Since(start))
	
	fmt.Printf("Diff: %v\n", math.Sin(x) - fastSin(x))
}
