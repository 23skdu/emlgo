//go:build !arm64 || !linux
// +build !arm64 !linux

package eml

func detectSVE() { _ = hasSVE }

func addSVE(a, b, result []float64) { _, _, _ = a, b, result }
