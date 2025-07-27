package goleakwatch

import (
	"context"
	"testing"
	"time"
)

func BenchmarkCheckNoLeak(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check(func() {
			// No goroutines created
		}, Config{
			Threshold:   1,
			Wait:        10 * time.Millisecond,
			EnableTrace: false,
		})
	}
}

func BenchmarkCheckWithGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check(func() {
			go func() {
				time.Sleep(1 * time.Millisecond)
			}()
		}, Config{
			Threshold:   1,
			Wait:        10 * time.Millisecond,
			EnableTrace: false,
		})
	}
}

func BenchmarkDefaultCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check(func() {
			// No goroutines created
		}, Config{
			Threshold:   1,
			Wait:        10 * time.Millisecond,
			EnableTrace: false,
		})
	}
}

func BenchmarkWithContext(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		CheckWithContext(ctx, func() {
			// No goroutines created
		}, Config{
			Threshold:   1,
			Wait:        10 * time.Millisecond,
			EnableTrace: false,
		})
	}
}

func BenchmarkBenchmarkFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Benchmark("test", func() {
			// No goroutines created
		})
	}
}
