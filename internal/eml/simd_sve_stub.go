//go:build !arm64 || !linux
// +build !arm64 !linux

package eml

func detectSVE() { _ = hasSVE }
