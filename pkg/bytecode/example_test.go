package bytecode_test

import (
	"fmt"
	"math"

	"github.com/emlgo/eml/pkg/bytecode"
)

func ExampleCompileExpr() {
	prog, err := bytecode.CompileExpr("x^2 + 1")
	if err != nil {
		panic(err)
	}
	result := prog.Eval([]float64{3.0}, nil)
	fmt.Printf("(x^2 + 1)(3) = %.0f\n", result)
	// Output: (x^2 + 1)(3) = 10
}

func ExampleCompileExpr_multipleVariables() {
	prog, err := bytecode.CompileExpr("x + y")
	if err != nil {
		panic(err)
	}
	result := prog.Eval([]float64{3.0, 4.0}, nil)
	fmt.Printf("(x + y)(3, 4) = %.0f\n", result)
	// Output: (x + y)(3, 4) = 7
}

func ExampleCompileExpr_emlOperator() {
	prog, err := bytecode.CompileExpr("eml(x, y)")
	if err != nil {
		panic(err)
	}
	// eml(1, e) = exp(1) - ln(e) = e - 1
	result := prog.Eval([]float64{1.0, math.E}, nil)
	fmt.Printf("eml(1, e) = %.6f\n", result)
	// Output: eml(1, e) = 1.718282
}

func ExampleProgram_Eval_zeroAlloc() {
	prog, err := bytecode.CompileExpr("x^2 + y^2")
	if err != nil {
		panic(err)
	}
	scratch := make([]float64, prog.MaxStackDepth)
	for i := 0.0; i < 5; i++ {
		result := prog.Eval([]float64{i, i + 1}, scratch)
		fmt.Printf("(%.0f² + %.0f²) = %.0f\n", i, i+1, result)
	}
	// Output:
	// (0² + 1²) = 1
	// (1² + 2²) = 5
	// (2² + 3²) = 13
	// (3² + 4²) = 25
	// (4² + 5²) = 41
}

func ExampleCompileExpr_expLog() {
	prog, err := bytecode.CompileExpr("exp(x) - log(y)")
	if err != nil {
		panic(err)
	}
	result := prog.Eval([]float64{1.0, math.E}, nil)
	fmt.Printf("(exp(1) - log(e)) = %.6f\n", result)
	// Output: (exp(1) - log(e)) = 1.718282
}

func ExampleCompileExpr_sqrt() {
	prog, err := bytecode.CompileExpr("sqrt(x)")
	if err != nil {
		panic(err)
	}
	result := prog.Eval([]float64{144.0}, nil)
	fmt.Printf("sqrt(144) = %.0f\n", result)
	// Output: sqrt(144) = 12
}
