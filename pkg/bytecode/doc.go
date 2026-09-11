// Package bytecode provides a data-oriented, zero-allocation virtual machine for mathematical
// expressions and symbolic search using Structure of Arrays (SoA) layout and Reverse Polish Notation (RPN).
//
// The core design eliminates pointer chasing, ensures prefetcher-friendly sequential memory access,
// reduces heap allocations during evaluation to zero, and amortizes expression interpretation overhead
// across large batch datasets via vectorized execution.
package bytecode
