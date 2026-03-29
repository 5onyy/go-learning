package main

import (
	"fmt"
)

func main() {
	var a []int
	for {
		var val int
		if _, err := fmt.Scanf("%d", &val); err != nil {
			break
		}
		a = append(a, val)
	}
	a = rotateLeft(a, 3)
	fmt.Print(a)
}

func rotateLeft(a []int, rot int) []int {
	for i := rot - 1; i >= 0; i-- {
		a = append(a, a[i])
	}
	return a[rot:]
}

// 0 1 2 3 4 5
