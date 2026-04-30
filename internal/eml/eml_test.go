package eml

import (
	"math"
	"testing"
)

func TestEml(t *testing.T) {
	tests := []struct {
		name   string
		x      float64
		y      float64
		want   float64
		ulpTol int
	}{
		{"basic", 0, 1, math.Exp(0) - math.Log(1), 1},
		{"positive", 1, 1, math.E - 0, 1},
		{"two", 2, 1, math.Exp(2) - 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Eml(tt.x, tt.y)
			if !ulpEqual(got, tt.want, tt.ulpTol) {
				t.Errorf("Eml(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestEmlOne(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"zero", 0, 1},
		{"one", 1, math.E},
		{"two", 2, math.E * math.E},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmlOne(tt.x)
			if !ulpEqual(got, tt.want, 1) {
				t.Errorf("EmlOne(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestOneEml(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"one", 1, 0},
		{"e", math.E, math.E - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OneEml(tt.x)
			if !ulpEqual(got, tt.want, 1) {
				t.Errorf("OneEml(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func BenchmarkEml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Eml(1.5, 2.5)
	}
}

func BenchmarkEmlOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EmlOne(1.5)
	}
}

func ulpEqual(a, b float64, tol int) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return true
	}
	if math.IsInf(a, -1) && math.IsInf(b, -1) {
		return true
	}
	if a == b {
		return true
	}
	ulp := int(math.Abs(a - b) / math.Nextafter(a, b) / 2)
	return ulp <= tol
}