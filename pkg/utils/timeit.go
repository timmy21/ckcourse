package utils

import (
	"fmt"
	"time"
)

func Elapsed(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
