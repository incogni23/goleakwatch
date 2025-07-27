package main

import (
	"fmt"
	"time"

	"github.com/incogni23/goleakwatch"
)

func main() {
	err := goleakwatch.DefaultCheck(func() {
		go func() {
			time.Sleep(1 * time.Second)
		}()
	})
	if err != nil {
		fmt.Println("Leak detected:", err)
	} else {
		fmt.Println("No leaks detected.")
	}
}
