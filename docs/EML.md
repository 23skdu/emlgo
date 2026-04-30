# EML: All Elementary Functions from a Single Operator

## Overview

This document summarizes the research from "All elementary functions from a single operator" (arXiv:2603.21852v2) by Andrzej Odrzywołek.

## Key Discovery

A single binary operator called **EML** (Exp-Minus-Log) can reconstruct all elementary functions:

```
eml(x, y) = exp(x) - ln(y)
```

Together with the constant `1`, this operator forms a complete basis for computing:
- Arithmetic operations (+, -, *, /, ^)
- Trigonometric functions (sin, cos, tan, arcsin, arccos, arctan)
- Hyperbolic functions (sinh, cosh, tanh, arsinh, arcosh, artanh)
- Logarithms and exponentials
- Roots and powers
- Constants (e, π, i)

## The Grammar

Every EML expression becomes a binary tree of identical nodes:

```
S → 1 | eml(S, S)
```

This is isomorphic to Catalan structures and full binary trees.

## Example Expressions

| Function | EML Expression |
|----------|-----------------|
| e^x | eml(x, 1) |
| ln(x) | eml(1, eml(eml(1, x), 1)) |
| -x | eml(eml(1, eml(x, 1)), eml(1, 1)) |
| 1/x | eml(eml(1, x), eml(x, 1)) |
| x + y | eml(1, eml(eml(y, 1), eml(1, x))) |

## Significance

- Analogous to NAND gate for Boolean logic
- Enables uniform circuit representation of mathematical expressions
- Useful for symbolic regression and gradient-based optimization
- Can be implemented in hardware as a single instruction

## Related Operators

The paper also discovered related operators:
- **EDL**: edl(x, y) = exp(x) / ln(y) (requires constant e)
- **-EML**: -eml(y, x) = ln(x) - exp(y) (requires constant -∞)

## Applications

1. **EML Compiler**: Converts standard formulas to EML form
2. **Analog Computing**: Build circuits from uniform EML nodes
3. **Symbolic Regression**: Trainable EML trees for discovering formulas from data
4. **Single-Instruction Computers**: RPN calculator with single EML operation