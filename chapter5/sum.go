package main

import "fmt"

func main() {
	fmt.Println(sum())
	fmt.Println(sum(3))
	fmt.Println(sum(1, 2, 3, 4))

	values := []int{1, 2, 3, 4}
	fmt.Println(sum(values...)) // Same as above, show how to invoke a variadic function when arguments are already in a slice
}

func sum(vals ...int) int { // Variadic function
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}
