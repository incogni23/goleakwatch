package utils

import (
	"context"
	"sync"
	"time"
)

// Timer provides utility functions for timing operations
type Timer struct {
	start time.Time
}

// NewTimer creates a new timer
func NewTimer() *Timer {
	return &Timer{
		start: time.Now(),
	}
}

// Elapsed returns the elapsed time since the timer was created
func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

// Reset resets the timer
func (t *Timer) Reset() {
	t.start = time.Now()
}

// WithTimeout executes a function with a timeout
func WithTimeout(ctx context.Context, timeout time.Duration, fn func() error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Retry executes a function with retry logic
func Retry(maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < maxAttempts {
			time.Sleep(delay)
		}
	}

	return lastErr
}

// ExponentialBackoff executes a function with exponential backoff
func ExponentialBackoff(maxAttempts int, initialDelay time.Duration, fn func() error) error {
	var lastErr error
	delay := initialDelay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < maxAttempts {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}

	return lastErr
}

// Debounce creates a debounced function that delays execution
func Debounce(delay time.Duration, fn func()) func() {
	var timer *time.Timer

	return func() {
		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(delay, fn)
	}
}

// Throttle creates a throttled function that limits execution frequency
func Throttle(interval time.Duration, fn func()) func() {
	var lastCall time.Time
	var mu sync.Mutex

	return func() {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		if now.Sub(lastCall) >= interval {
			fn()
			lastCall = now
		}
	}
}
