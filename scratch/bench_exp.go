package main

import (
	"fmt"
	"math"
	"time"
)

func FmaScalar(a, b, c float64) float64 {
	return a*b + c
}

func fastExp(x float64) float64 {
	const (
		Log2E = 1.44269504088896340736 // 1/ln(2)
		Ln2Hi = 0.69314718036912381649
		Ln2Lo = 1.9082149292705877000e-10
	)
	k := math.Round(x * Log2E)
	r := x - k*Ln2Hi - k*Ln2Lo

	p := 0.008333333333333333
	p = FmaScalar(p, r, 0.041666666666666664)
	p = FmaScalar(p, r, 0.16666666666666666)
	p = FmaScalar(p, r, 0.5)
	p = FmaScalar(p, r, 1.0)
	p = FmaScalar(p, r, 1.0)

	return math.Ldexp(p, int(k))
}

func main() {
	x := 1.23
	iterations := 10000000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = math.Exp(x)
	}
	fmt.Printf("math.Exp: %v\n", time.Since(start))

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = fastExp(x)
	}
	fmt.Printf("fastExp:  %v\n", time.Since(start))
	
	fmt.Printf("Diff: %v\n", math.Exp(x) - fastExp(x))
}
