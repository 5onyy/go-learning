package main

import "fmt"

const MAX_N = 500

func main() {
	var n int
	fmt.Scanf("%d", &n)
	var a [MAX_N]int
	for i := 0; i < n; i++ {
		fmt.Scanf("%d", &a[i])
	}
	reverse(&a, n)
	fmt.Println(a[:n])
}

func reverse(a *[MAX_N]int, n int) {
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
