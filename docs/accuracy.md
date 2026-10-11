# Numerical Accuracy and Semantics

This document describes the current numerical behavior and the assurance
boundaries of `emlgo`. It is not a claim that every operation is correctly
rounded or certified. Numerical results can differ across Go versions,
architectures, and execution backends; correctness budgets apply only where a
specific test or API contract says they do.

## Backend contract

| API family | Implementation/reference | Accuracy expectation | Domain and special values |
| :--- | :--- | :--- | :--- |
| `pkg/arithmetic`, `pkg/trig`, `pkg/hyper`, `pkg/logexp` scalar wrappers | Primarily Go `math`; see individual package source | Usually delegates to the corresponding Go operation. This does not promise correctly rounded results across every platform. | Follow Go `math` semantics unless the individual wrapper documents additional behavior. |
| `internal/eml` real batch/SIMD | Architecture dispatch, scalar fallbacks, and batch kernels | Batch-kernel tests define operation-specific mixed abs/relative and ULP budgets; those budgets are not interchangeable. | Test NaN/Inf classes, zeros, subnormals, vector tails, and each enabled backend against the operation's reference. |
| `pkg/fastmath` | Range-reduced polynomial and bit-manipulation approximations | Approximate by design. Current sweep budgets: `FastExpBitCast` 5% relative; `FastExpMinimax` `1.5e-6` relative on normal outputs and 2 ULP on nonzero subnormals; `FastExpF32` `2e-4` relative on normal outputs and 20 ULP on nonzero subnormals; `FastLogMinimax` `1e-7` absolute; `FastLogChebyshev` `5e-4` absolute; `FastLogF32` `1e-5` absolute. These are sampled regression-test budgets, not certified global bounds. | Each function's domain guards and overflow/underflow behavior are implemented in its source. FastMath functions must not be substituted for standard-accuracy APIs without accepting approximation error. |
| JIT and bytecode evaluation | Generated native code or interpreter | Compare against the AST/scalar reference with expression-specific budgets; no general correctly-rounded guarantee. | Parser, domain, and exceptional-value behavior must be tested per operation and target. JIT availability and executable-memory behavior are platform dependent. |
| Complex APIs | Go `math/cmplx` and batch wrappers | Follows Go complex math behavior unless a batch-specific tolerance is documented. | Branch cuts and exceptional values are significant; tests should include both sides of each documented cut. |
| `internal/eml/bigmath` | Project-owned `math/big.Float` algorithms | Precision is finite and algorithm-dependent; successful evaluation is not independent validation of another emlgo implementation. | Function-specific domain rules apply. Do not use as the only reference for validating itself. |

## Comparing floating-point results

- First compare semantic classes: NaN, positive/negative infinity, finite, and
  signed zero when the API promises to preserve it.
- For finite values away from zero, report ULP distance and relative error.
- Near zero or at cancellation points, use an explicit absolute-error budget;
  relative error is ill-conditioned there. A mixed budget has the form
  `abs(got-want) <= absTol + relTol*abs(want)`.
- Treat subnormal inputs and outputs explicitly. A normal-range-only reference
  does not establish subnormal accuracy.
- For fused operations, use the fused operation as the reference when fused
  rounding is part of the API contract; comparing to a separately rounded
  multiply then add tests a different operation.

## Independent reference vectors

`internal/accuracyvectors/mpfr-vectors.csv` contains representative MPFR 4.2.2
vectors for binary64/binary32 exponentials and logarithms. Inputs are exact
binary values; MPFR uses round-to-nearest and evaluates at the target result
precision (53 or 24 bits). Tests consume these vectors without MPFR installed.
To regenerate them with MPFR 4.2.2 development headers and libraries:

```bash
cc -Wall -Wextra -Werror scripts/generate_mpfr_vectors.c -lmpfr -lgmp \
  -o /tmp/generate_mpfr_vectors
/tmp/generate_mpfr_vectors internal/accuracyvectors/mpfr-vectors.csv
```

The generator rejects MPFR versions other than 4.2.2 so the provenance header
cannot silently change. The current fixture is deliberately a focused reference
set, not exhaustive coverage of all functions or inputs.

The FastMath regression sweep uses a fixed PRNG seed, dense binary64 exponential
inputs, reduction-boundary neighbors, exponent/mantissa-stratified logarithm
inputs, random positive binary32 logarithm bit patterns, and binary32
exponential reduction boundaries. It logs the worst observed finite error and
its input when run with `-v`; the MPFR fixtures independently anchor selected
normal and subnormal cases.

## Current validation limits

MPFR vectors provide an independent high-precision reference for the included
exp/log points. The Go `math` package remains the practical reference for many
other tests, but it is not an independent arbitrary-precision oracle. Existing
dense SIMD checks cover selected functions and input ranges; they do not certify
all inputs, all architectures, or all possible backends. FastMath sweeps and
golden vectors likewise do not prove global worst-case bounds. See
[nextsteps.md](nextsteps.md) Steps 41–50 for the remaining work to expand
reference coverage, domain testing, cross-platform execution, and
reproducibility.

Published accuracy numbers must identify the operation, precision, input
domain, backend, reference, error metric, sample method, and Go/architecture
environment. Do not describe a sampled maximum as a mathematical bound.
