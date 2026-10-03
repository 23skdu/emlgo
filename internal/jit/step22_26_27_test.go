package jit

import (
	"math"
	"testing"
)

func TestParseScientificNotationAndWhitespace(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"1e-5", 1e-5},
		{"2.5e+3", 2500.0},
		{"1.23E4", 12300.0},
		{"1e0", 1.0},
		{"1e2 + 5", 105.0},
		{"1 + \t 2", 3.0},
		{"1 +\n2", 3.0},
		{"1 +\r\n 2", 3.0},
		{"\t\n 5.5e-1 \t\r\n", 0.55},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			node, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tc.input, err)
			}
			got := Eval(node, 0.0)
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Eval(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseWithVarsValidation(t *testing.T) {
	// Should succeed when variables match declared set
	node, err := ParseWithVars("x + y", []string{"x", "y"})
	if err != nil {
		t.Fatalf("ParseWithVars failed for declared variables: %v", err)
	}
	if node == nil {
		t.Fatal("node is nil")
	}

	// Should fail when undeclared variable is used
	_, err = ParseWithVars("x + z", []string{"x", "y"})
	if err == nil {
		t.Fatal("expected error for undeclared variable 'z', got nil")
	}
}

func TestDecompilePrecedence(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want string
	}{
		{"mul over add", "(x + 1) * 2", "(x + 1.0) * 2.0"},
		{"add over mul", "x + 1 * 2", "x + 1.0 * 2.0"},
		{"div right mul", "x / (y * z)", "x / (y * z)"},
		{"sub right sub", "x - (y - z)", "x - (y - z)"},
		{"pow over add", "(x + 1)^2", "(x + 1.0)^2.0"},
		{"pow right neg", "x^(-2)", "x^(-2.0)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node, err := Parse(tc.expr)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tc.expr, err)
			}
			emlNode := Canonicalize(node)
			decomp := Decompile(emlNode)

			// The decompiled string must evaluate to the same value as the original
			reparsed, err := Parse(decomp)
			if err != nil {
				t.Fatalf("Failed to re-parse decompiled string %q: %v", decomp, err)
			}
			valOrig := EvalVars(node, map[string]float64{"x": 3.0, "y": 4.0, "z": 2.0})
			valRe := EvalVars(reparsed, map[string]float64{"x": 3.0, "y": 4.0, "z": 2.0})
			if math.Abs(valOrig-valRe) > 1e-9 {
				t.Fatalf("Value mismatch for %q -> %q: orig=%v, reparsed=%v",
					tc.expr, decomp, valOrig, valRe)
			}
		})
	}
}

func TestDecompileLaTeXNoSpuriousParens(t *testing.T) {
	node, err := Parse("sin(-x)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	emlNode := Canonicalize(node)
	latex := DecompileLaTeX(emlNode)
	// Must not be wrapped in spurious outer parens: want \sin(-x), not (\sin(-x))
	if latex == "(\\sin(-x))" {
		t.Errorf("spurious outer parens in DecompileLaTeX: %q", latex)
	}
}

func TestJITNegativeBasePowerAndUnaryMinus(t *testing.T) {
	tests := []struct {
		name string
		expr string
		x    float64
		want float64
	}{
		{"negative base odd power", "x^3", -2.0, -8.0},
		{"negative base even power", "x^4", -2.0, 16.0},
		{"negative base non-constant power", "x^(1 + 2)", -2.0, -8.0},
		{"zero to zero", "x^0", 0.0, 1.0},
		{"unary minus positive", "-x", 5.5, -5.5},
		{"unary minus negative", "-x", -4.2, 4.2},
		{"double unary minus", "-(-x)", -3.14, -3.14},
	}

	compiler := NewCompiler()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fn, err := compiler.Compile(tc.expr)
			if err != nil {
				t.Fatalf("Compile(%q) error: %v", tc.expr, err)
			}
			got := fn(tc.x)
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Compile(%q)(%v) = %v, want %v", tc.expr, tc.x, got, tc.want)
			}
		})
	}
}

func FuzzScientificFloatRoundTrip(f *testing.F) {
	f.Add(1e-5)
	f.Add(2.5e+3)
	f.Add(1.2345e-8)
	f.Add(-7.89e12)
	f.Add(0.0)

	f.Fuzz(func(t *testing.T, val float64) {
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return
		}
		node := Number{Value: val}
		emlNode := Canonicalize(node)
		decomp := Decompile(emlNode)
		reparsed, err := Parse(decomp)
		if err != nil {
			t.Fatalf("failed to parse decompiled %q for %v: %v", decomp, val, err)
		}
		reVal := Eval(reparsed, 0)
		diff := math.Abs(val - reVal)
		denom := math.Max(math.Abs(val), math.Abs(reVal))
		if denom > 1e-9 {
			diff /= denom
		}
		if diff > 1e-7 {
			t.Fatalf("round-trip mismatch: val=%v, decomp=%q, got=%v", val, decomp, reVal)
		}
	})
}
