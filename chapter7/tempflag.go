package main

import (
	"flag"
	"fmt"
	"go-learning/chapter7/tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "The temperature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}
