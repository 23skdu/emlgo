package eml

import (
	"sync"
	"testing"
)

// TestStopWorkerPoolIdempotent verifies that StopWorkerPool can be called
// multiple times without panicking (double-close protection via sync.Once).
//
// Note: the worker pool is package-level state initialized once in init().
// This test calls StopWorkerPool on the live pool. After the pool is stopped,
// parallelization falls back to sequential processing (channel send would block
// on a closed channel, so parallel paths should not be exercised after Stop).
func TestStopWorkerPoolIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping pool shutdown test in short mode")
	}
	// Calling StopWorkerPool many times sequentially must not panic.
	for i := 0; i < 10; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("StopWorkerPool panicked on call %d: %v", i, r)
				}
			}()
			StopWorkerPool()
		}()
	}
}

// TestStopWorkerPoolConcurrent verifies that concurrent calls to StopWorkerPool
// do not cause a double-close panic. Uses sync.Once internally.
func TestStopWorkerPoolConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping pool shutdown test in short mode")
	}

	const goroutines = 10
	panicked := make(chan interface{}, goroutines)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicked <- r
				}
			}()
			StopWorkerPool()
		}()
	}
	wg.Wait()
	close(panicked)

	for p := range panicked {
		t.Errorf("StopWorkerPool panicked: %v", p)
	}
}
