package bytecode

import "math"

// unaryFuncs maps a function name to its float64 implementation for every
// unary elementary function the VM supports.
//
// This is the single source of truth shared by the tokenizer, the code
// generator, the optimizer and the program printer. An earlier version had two
// independent tables -- isFunc() accepted sin and cos, while emitFunc() only
// knew eml/exp/log/sqrt -- so `sin(x)` was lexed as a call and then rejected at
// emit time with a confusing "unsupported function" error.
//
// The set matches internal/jit/functab.go exactly, so an expression that
// compiles for the JIT also compiles for the VM.
var unaryFuncs = map[string]mathFunc{
	"sin":   math.Sin,
	"cos":   math.Cos,
	"tan":   math.Tan,
	"exp":   math.Exp,
	"log":   math.Log,
	"ln":    math.Log,
	"sqrt":  math.Sqrt,
	"asin":  math.Asin,
	"acos":  math.Acos,
	"atan":  math.Atan,
	"abs":   math.Abs,
	"cbrt":  math.Cbrt,
	"log2":  math.Log2,
	"log10": math.Log10,
	"ceil":  math.Ceil,
	"floor": math.Floor,
	"trunc": math.Trunc,
	"round": math.Round,
	"sinh":  math.Sinh,
	"cosh":  math.Cosh,
	"tanh":  math.Tanh,
	"asinh": math.Asinh,
	"acosh": math.Acosh,
	"atanh": math.Atanh,
	"erf":   math.Erf,
	"gamma": math.Gamma,
}

// mathFunc is a unary float64 function.
type mathFunc func(float64) float64

// emlName is the only non-unary function the VM supports. It pops two operands,
// so it is handled separately from unaryFuncs everywhere.
const emlName = "eml"

// isFunc reports whether name denotes a function call. It covers both the unary
// elementary functions and the two-argument eml operator.
func isFunc(name string) bool {
	if name == emlName {
		return true
	}
	_, ok := unaryFuncs[name]
	return ok
}

// unaryOpcode returns the opcode implementing name.
func unaryOpcode(name string) (OpCode, bool) {
	op, ok := unaryCodes[name]
	return op, ok
}

// funcName returns the canonical source name for a unary opcode, used when
// printing programs.
func funcName(op OpCode) (string, bool) {
	name, ok := opcodeNames[op]
	return name, ok
}

// unaryCodes and opcodeNames are the bidirectional index over unaryFuncs. They
// are built once from it so that the three views cannot disagree.
var (
	unaryCodes   map[string]OpCode
	opcodeNames  map[OpCode]string
	canonicalNms map[string]string
)

func init() {
	unaryCodes = make(map[string]OpCode, len(unaryFuncs))
	opcodeNames = make(map[OpCode]string, len(unaryFuncs))
	canonicalNms = make(map[string]string, len(unaryFuncs))
	for name, op := range map[string]OpCode{
		"sin":   OpSin,
		"cos":   OpCos,
		"tan":   OpTan,
		"exp":   OpExp,
		"log":   OpLog,
		"ln":    OpLog,
		"sqrt":  OpSqrt,
		"asin":  OpAsin,
		"acos":  OpAcos,
		"atan":  OpAtan,
		"abs":   OpAbs,
		"cbrt":  OpCbrt,
		"log2":  OpLog2,
		"log10": OpLog10,
		"ceil":  OpCeil,
		"floor": OpFloor,
		"trunc": OpTrunc,
		"round": OpRound,
		"sinh":  OpSinh,
		"cosh":  OpCosh,
		"tanh":  OpTanh,
		"asinh": OpAsinh,
		"acosh": OpAcosh,
		"atanh": OpAtanh,
		"erf":   OpErf,
		"gamma": OpGamma,
	} {
		unaryCodes[name] = op
		if _, seen := opcodeNames[op]; !seen {
			opcodeNames[op] = name
			canonicalNms[name] = name
		}
		// Every alias must resolve through unaryFuncs too.
		if _, ok := unaryFuncs[name]; !ok {
			panic("bytecode: opcode table references unknown function " + name)
		}
	}
	if len(unaryCodes) != len(unaryFuncs) {
		panic("bytecode: unaryFuncs and unaryCodes have drifted apart")
	}

	opUnaryFns = make(map[OpCode]mathFunc, len(unaryFuncs))
	for name, op := range unaryCodes {
		opUnaryFns[op] = unaryFuncs[name]
		opUnaryTable[op] = unaryFuncs[name]
	}
}

// funcOpFirst and funcOpLast bound the contiguous block of function opcodes
// declared in program.go. Opcode space is allocated deliberately so that the
// whole block can be walked with a simple range.
const (
	funcOpFirst = OpSqrt
	funcOpLast  = OpGamma
)

// opUnaryTable is the direct flat-indexed dispatch table for unary opcodes.
// Opcode dispatch is a single array index lookup with no hash map overhead.
var opUnaryTable [256]mathFunc

// opUnaryFns is retained for backwards compatibility.
var opUnaryFns map[OpCode]mathFunc
