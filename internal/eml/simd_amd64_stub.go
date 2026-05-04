//go:build !amd64
// +build !amd64

package eml

import "runtime"

func addAVX2(a, b, result []float64)           { _, _, _ = a, b, result }
func subAVX2(a, b, result []float64)           { _, _, _ = a, b, result }
func mulAVX2(a, b, result []float64)           { _, _, _ = a, b, result }
func divAVX2(a, b, result []float64)           { _, _, _ = a, b, result }
func addScalarAVX2(a []float64, b float64, result []float64) { _, _, _ = a, b, result }
func mulScalarAVX2(a []float64, b float64, result []float64) { _, _, _ = a, b, result }

func addAVX512(a, b, result []float64)           { _, _, _ = a, b, result }
func subAVX512(a, b, result []float64)           { _, _, _ = a, b, result }
func mulAVX512(a, b, result []float64)           { _, _, _ = a, b, result }
func divAVX512(a, b, result []float64)           { _, _, _ = a, b, result }
func addScalarAVX512(a []float64, b float64, result []float64) { _, _, _ = a, b, result }
func mulScalarAVX512(a []float64, b float64, result []float64) { _, _, _ = a, b, result }

func sqrtAVX2(a, result []float64)   { _, _ = a, result }
func sqrtAVX512(a, result []float64) { _, _ = a, result }

func fmaAVX2(a, b, c, result []float64)   { _, _, _, _ = a, b, c, result }
func fmaAVX512(a, b, c, result []float64) { _, _, _, _ = a, b, c, result }

func detectAMD64SIMD() { _ = runtime.GOARCH }


