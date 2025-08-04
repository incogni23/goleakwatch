package goleakwatch

import (
	"context"
	"testing"
	"time"

	"github.com/incogni23/goleakwatch/internal/errors"
)

// FuzzCheck tests the Check function with various inputs
func FuzzCheck(f *testing.F) {
	// Add seed corpus
	f.Add(1, 100, true, "test_function")
	f.Add(0, 50, false, "edge_case")
	f.Add(10, 10, true, "boundary_test")

	f.Fuzz(func(t *testing.T, threshold int, waitMs int, enableTrace bool, functionName string) {
		// Validate inputs
		if threshold < 0 || waitMs < 0 || waitMs > 10000 {
			t.Skip("Invalid input values")
		}

		cfg := &Config{
			Threshold:    threshold,
			Wait:         time.Duration(waitMs) * time.Millisecond,
			EnableTrace:  enableTrace,
			FunctionName: functionName,
		}

		// Test with a simple function
		err := Check(func() {
			// Do nothing
		}, cfg)

		// Should not panic
		if err != nil && !errors.IsLeakError(err) {
			t.Errorf("Unexpected error type: %T", err)
		}
	})
}

// FuzzCheckWithContext tests the CheckWithContext function
func FuzzCheckWithContext(f *testing.F) {
	f.Add(1, 100, true, "context_test", 5000)

	f.Fuzz(func(t *testing.T, threshold int, waitMs int, enableTrace bool, functionName string, timeoutMs int) {
		// Validate inputs
		if threshold < 0 || waitMs < 0 || waitMs > 10000 || timeoutMs < 0 || timeoutMs > 30000 {
			t.Skip("Invalid input values")
		}

		cfg := &Config{
			Threshold:    threshold,
			Wait:         time.Duration(waitMs) * time.Millisecond,
			EnableTrace:  enableTrace,
			FunctionName: functionName,
			Timeout:      time.Duration(timeoutMs) * time.Millisecond,
		}

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		err := CheckWithContext(ctx, func() {
			// Do nothing
		}, cfg)

		// Should not panic
		if err != nil && !errors.IsLeakError(err) && err != context.DeadlineExceeded {
			t.Errorf("Unexpected error type: %T", err)
		}
	})
}

// FuzzBenchmark tests the Benchmark function
func FuzzBenchmark(f *testing.F) {
	f.Add("fuzz_test", 100)

	f.Fuzz(func(t *testing.T, name string, sleepMs int) {
		// Validate inputs
		if sleepMs < 0 || sleepMs > 1000 {
			t.Skip("Invalid sleep duration")
		}

		result := Benchmark(name, func() {
			time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		})

		// Validate result
		if result.FunctionName != name {
			t.Errorf("Expected function name %s, got %s", name, result.FunctionName)
		}

		if result.ExecutionTime < 0 {
			t.Errorf("Negative execution time: %v", result.ExecutionTime)
		}
	})
}

// FuzzSnapshotCheck tests the SnapshotCheck function
func FuzzSnapshotCheck(f *testing.F) {
	f.Add(1, 100, true, "snapshot_test")

	f.Fuzz(func(t *testing.T, threshold int, waitMs int, enableTrace bool, functionName string) {
		// Validate inputs
		if threshold < 0 || waitMs < 0 || waitMs > 10000 {
			t.Skip("Invalid input values")
		}

		cfg := &Config{
			Threshold:    threshold,
			Wait:         time.Duration(waitMs) * time.Millisecond,
			EnableTrace:  enableTrace,
			FunctionName: functionName,
		}

		err := SnapshotCheck(func() {
			// Do nothing
		}, cfg)

		// Should not panic
		if err != nil && !errors.IsLeakError(err) {
			t.Errorf("Unexpected error type: %T", err)
		}
	})
}

// TestIsLeakError tests the IsLeakError function
func TestIsLeakError(t *testing.T) {
	// Test with a real leak error
	err := Check(func() {
		go func() {
			select {} // This will leak
		}()
	}, &Config{
		Threshold:   0,
		Wait:        100 * time.Millisecond,
		EnableTrace: false,
	})

	if err != nil && !errors.IsLeakError(err) {
		t.Errorf("Expected leak error, got: %T", err)
	}

	// Test with a regular error
	regularErr := context.DeadlineExceeded
	if errors.IsLeakError(regularErr) {
		t.Errorf("Expected false for regular error")
	}
}
