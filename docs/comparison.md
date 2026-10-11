# emlgo and Go's `math` Package

## Performance claims

Older comparisons in this document used inconsistent benchmark setups and did
not retain enough raw output and environment metadata to reproduce the reported
ratios. Those performance claims have been withdrawn; they must not be treated
as current results or release claims. The repository does not claim that emlgo
is generally faster than Go's `math` package.

Use the reproducible procedure in [`performance.md`](performance.md) to measure
a specific workload on the target machine. The continuous-integration
regression gate compares selected batch benchmarks between revisions on the
same hosted runner; it does not compare emlgo against `math` or certify
cross-platform performance.

## API and numerical behavior

The public scalar wrappers generally delegate to Go's `math` package, with
individual exceptions and semantics documented in their package references.
`pkg/fastmath` intentionally includes approximate algorithms. Batch and fused
operations use their own implementations and are validated against
operation-specific tolerances; no single accuracy guarantee applies to every
backend.

See [`FUNCTIONS.md`](FUNCTIONS.md) for the API inventory and
[`accuracy.md`](accuracy.md) for numerical contracts, test budgets, and
assurance limits.
