package bytecode

import (
	"math"

	"github.com/emlgo/eml/pkg/fastmath"
)

// Eval evaluates the linear bytecode program against the given input variables.
// If scratch is provided and has capacity >= p.MaxStackDepth, it is reused to achieve
// 0 heap allocations. If scratch is nil or insufficient, a local 32-element stack array
// on the goroutine stack is utilized.
func (p *Program) Eval(vars []float64, scratch []float64) float64 {
	if p == nil || len(p.Ops) == 0 {
		return 0
	}

	var localStack [32]float64
	var stack []float64
	if scratch != nil && len(scratch) >= p.MaxStackDepth {
		stack = scratch
	} else if p.MaxStackDepth <= 32 {
		stack = localStack[:p.MaxStackDepth]
	} else {
		stack = make([]float64, p.MaxStackDepth)
	}

	sp := 0
	cIdx := 0
	vIdx := 0

	for _, op := range p.Ops {
		switch op {
		case OpConst:
			stack[sp] = p.Consts[cIdx]
			sp++
			cIdx++
		case OpVar:
			idx := p.VarIndices[vIdx]
			if int(idx) < len(vars) {
				stack[sp] = vars[idx]
			} else {
				stack[sp] = 0
			}
			sp++
			vIdx++
		case OpEML:
			sp--
			y := stack[sp]
			x := stack[sp-1]
			stack[sp-1] = fastmath.FastEml(x, y)
		case OpAdd:
			sp--
			stack[sp-1] += stack[sp]
		case OpSub:
			sp--
			stack[sp-1] -= stack[sp]
		case OpMul:
			sp--
			stack[sp-1] *= stack[sp]
		case OpDiv:
			sp--
			stack[sp-1] /= stack[sp]
		case OpPow:
			sp--
			stack[sp-1] = math.Pow(stack[sp-1], stack[sp])
		case OpNeg:
			stack[sp-1] = -stack[sp-1]
		case OpInv:
			stack[sp-1] = 1.0 / stack[sp-1]
		case OpSqrt:
			stack[sp-1] = math.Sqrt(stack[sp-1])
		case OpExp:
			stack[sp-1] = fastmath.FastExp(stack[sp-1])
		case OpLog:
			stack[sp-1] = fastmath.FastLog(stack[sp-1])
		}
	}

	if sp > 0 {
		return stack[0]
	}
	return 0
}

// EvalRegularized evaluates the bytecode with smooth domain regularization:
// ln_eps(y) = 0.5 * ln(y^2 + eps^2) and clamped exp.
// This prevents NaN cascading during genetic programming and numeric fitting.
func (p *Program) EvalRegularized(vars []float64, eps float64, scratch []float64) float64 {
	if p == nil || len(p.Ops) == 0 {
		return 0
	}
	if eps <= 0 {
		eps = fastmath.DefaultEpsilon
	}

	var localStack [32]float64
	var stack []float64
	if scratch != nil && len(scratch) >= p.MaxStackDepth {
		stack = scratch
	} else if p.MaxStackDepth <= 32 {
		stack = localStack[:p.MaxStackDepth]
	} else {
		stack = make([]float64, p.MaxStackDepth)
	}

	sp := 0
	cIdx := 0
	vIdx := 0

	for _, op := range p.Ops {
		switch op {
		case OpConst:
			stack[sp] = p.Consts[cIdx]
			sp++
			cIdx++
		case OpVar:
			idx := p.VarIndices[vIdx]
			if int(idx) < len(vars) {
				stack[sp] = vars[idx]
			} else {
				stack[sp] = 0
			}
			sp++
			vIdx++
		case OpEML:
			sp--
			y := stack[sp]
			x := stack[sp-1]
			stack[sp-1] = fastmath.FastEmlRegularized(x, y, eps)
		case OpAdd:
			sp--
			stack[sp-1] += stack[sp]
		case OpSub:
			sp--
			stack[sp-1] -= stack[sp]
		case OpMul:
			sp--
			stack[sp-1] *= stack[sp]
		case OpDiv:
			sp--
			denom := stack[sp]
			// Regularized division to prevent division by zero
			stack[sp-1] = stack[sp-1] * denom / (denom*denom + eps*eps)
		case OpPow:
			sp--
			stack[sp-1] = math.Pow(stack[sp-1], stack[sp])
		case OpNeg:
			stack[sp-1] = -stack[sp-1]
		case OpInv:
			denom := stack[sp-1]
			stack[sp-1] = denom / (denom*denom + eps*eps)
		case OpSqrt:
			v := stack[sp-1]
			if v < 0 {
				v = -v
			}
			stack[sp-1] = math.Sqrt(v)
		case OpExp:
			stack[sp-1] = fastmath.ExpClamped(stack[sp-1])
		case OpLog:
			stack[sp-1] = fastmath.LnRegularized(stack[sp-1], eps)
		}
	}

	if sp > 0 {
		return stack[0]
	}
	return 0
}
