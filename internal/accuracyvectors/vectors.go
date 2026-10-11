package accuracyvectors

import _ "embed"

// MPFR contains deterministic reference vectors for accuracy regression tests.
//
//go:embed mpfr-vectors.csv
var MPFR string
