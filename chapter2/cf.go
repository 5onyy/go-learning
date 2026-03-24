package main

import (
	"fmt"
	"go-learning/chapter2/popcount"
	"go-learning/chapter2/tempconv"
)

func main() {
	fmt.Printf("Boiling point: %g°F\n", tempconv.CToF(tempconv.BoilingC))
	fmt.Printf("Freezing point: %g°F\n", tempconv.CToF(tempconv.FreezingC))

	var a uint64 = 15

	fmt.Printf("%d", popcount.PopCount2(a))
}
