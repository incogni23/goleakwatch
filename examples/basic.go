package main

import (
	"fmt"
	"time"

	"github.com/incogni23/goleakwatch"
	"github.com/incogni23/goleakwatch/internal/errors"
	"github.com/incogni23/goleakwatch/internal/logger"
)

func main() {
	fmt.Println("🚀 goleakwatch Examples")
	fmt.Println("=======================")

	// Set up custom logger
	customLogger := logger.NewDefaultLogger(logger.INFO, nil)
	goleakwatch.SetLogger(customLogger)

	// Basic leak detection example
	fmt.Println("\n📋 Basic Leak Detection:")
	err := goleakwatch.DefaultCheck(func() {
		go func() {
			time.Sleep(1 * time.Second)
		}()
	})
	if err != nil {
		fmt.Println("Leak detected:", err)
		// Check if it's our custom error type
		if errors.IsLeakError(err) {
			if leakErr, ok := errors.GetLeakError(err); ok {
				fmt.Printf("Custom error details: %s\n", leakErr.Summary())
				fmt.Printf("Is significant leak: %v\n", leakErr.IsSignificant(2.0))
			}
		}
	} else {
		fmt.Println("No leaks detected.")
	}

	// Custom configuration example
	fmt.Println("\n⚙️ Custom Configuration:")
	err = goleakwatch.Check(func() {
		go func() {
			select {} // This will leak
		}()
	}, &goleakwatch.Config{
		Threshold:    0,
		Wait:         200 * time.Millisecond,
		EnableTrace:  false,
		FunctionName: "leakyFunction",
	})
	if err != nil {
		fmt.Println("Custom config leak detected:", err)
	}

	// Snapshot comparison example
	fmt.Println("\n📸 Snapshot Comparison:")
	err = goleakwatch.SnapshotCheck(func() {
		go func() {
			time.Sleep(50 * time.Millisecond)
		}()
	}, &goleakwatch.Config{
		Threshold:    1,
		Wait:         100 * time.Millisecond,
		EnableTrace:  false,
		FunctionName: "snapshotTest",
	})
	if err != nil {
		fmt.Println("Snapshot check error:", err)
	} else {
		fmt.Println("Snapshot check completed successfully.")
	}

	// Run benchmark examples
	runBenchmarkExamples()
}
