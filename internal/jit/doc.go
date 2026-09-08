// Package jit provides a JIT compiler for mathematical expressions targeting amd64.
//
// The compiler parses expressions into an AST, generates x86-64 machine code with SSE2
// instructions, and emits executable code via mmap/mprotect. Supports arithmetic operators,
// 16 math functions (sin, cos, exp, log, sqrt, tan, asin, acos, atan, abs, cbrt, log2,
// log10, ceil, floor, trunc), and power expressions with both integer and non-integer exponents.
//
// The JIT compiler is only available on amd64 platforms. Other architectures return an
// error from Compile.
package jit
