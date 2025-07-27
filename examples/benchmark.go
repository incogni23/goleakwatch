package main

import (
	"fmt"
	"time"

	"github.com/incogni23/goleakwatch"
)

func runBenchmarkExamples() {
	fmt.Println("🔍 Benchmarking Examples")
	fmt.Println("========================")

	// Example 1: Benchmark a function with no leaks
	fmt.Println("\n1. Benchmarking function with no leaks:")
	result1 := goleakwatch.Benchmark("noLeakFunction", func() {
		// This function doesn't create any goroutines
		time.Sleep(1 * time.Millisecond)
	})
	fmt.Println(result1)
	fmt.Printf("Is leaking: %t\n", result1.IsLeaking())

	// Example 2: Benchmark a function that creates temporary goroutines
	fmt.Println("\n2. Benchmarking function with temporary goroutines:")
	result2 := goleakwatch.Benchmark("tempGoroutines", func() {
		for i := 0; i < 3; i++ {
			go func() {
				time.Sleep(10 * time.Millisecond)
			}()
		}
		// Wait for goroutines to complete
		time.Sleep(20 * time.Millisecond)
	})
	fmt.Println(result2)
	fmt.Printf("Is leaking: %t\n", result2.IsLeaking())

	// Example 3: Benchmark a function with a leak
	fmt.Println("\n3. Benchmarking function with a leak:")
	result3 := goleakwatch.Benchmark("leakingFunction", func() {
		go func() {
			select {} // This goroutine will never exit
		}()
	})
	fmt.Println(result3)
	fmt.Printf("Is leaking: %t\n", result3.IsLeaking())

	// Example 4: Compare multiple functions
	fmt.Println("\n4. Comparing multiple functions:")

	functions := []struct {
		name string
		fn   func()
	}{
		{
			name: "fastFunction",
			fn: func() {
				// Fast, no goroutines
			},
		},
		{
			name: "slowFunction",
			fn: func() {
				time.Sleep(5 * time.Millisecond)
			},
		},
		{
			name: "goroutineFunction",
			fn: func() {
				go func() {
					time.Sleep(1 * time.Millisecond)
				}()
				time.Sleep(2 * time.Millisecond)
			},
		},
	}

	for _, f := range functions {
		result := goleakwatch.Benchmark(f.name, f.fn)
		fmt.Printf("  %s: %s\n", f.name, result)
	}
}
