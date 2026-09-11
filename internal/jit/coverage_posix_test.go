//go:build !windows && amd64 && (!js || !wasm)

package jit

import (
	"runtime"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

// -- AllocateExecutableMemory error paths via RLIMIT_AS --

func TestCompileMemoryAllocationError(t *testing.T) {
	runtime.LockOSThread()
	var oldRlim unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_AS, &oldRlim); err != nil {
		runtime.UnlockOSThread()
		t.Skip("cannot get RLIMIT_AS:", err)
	}
	zero := unix.Rlimit{Cur: 0, Max: oldRlim.Max}
	if err := unix.Setrlimit(unix.RLIMIT_AS, &zero); err != nil {
		runtime.UnlockOSThread()
		t.Skip("cannot set RLIMIT_AS:", err)
	}
	_, err := NewCompiler().Compile("x")
	if err == nil || !strings.Contains(err.Error(), "memory allocation") {
		_ = unix.Setrlimit(unix.RLIMIT_AS, &oldRlim) // #nosec G104 - best effort cleanup in test
		runtime.UnlockOSThread()
		t.Fatalf("expected memory allocation error, got: %v", err)
	}
	_ = unix.Setrlimit(unix.RLIMIT_AS, &oldRlim) // #nosec G104 - best effort cleanup in test
	runtime.UnlockOSThread()
}
