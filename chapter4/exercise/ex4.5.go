package main

import "fmt"

func main() {
	var a []string
	for {
		var val string
		if _, err := fmt.Scanf("%s", &val); err != nil {
			break
		}
		a = append(a, val)
	}
	z := eliminate(a)
	fmt.Print(z)
}

func eliminate(a []string) []string {
	n := len(a)
	if n == 0 {
		return []string{}
	}
	var z []string
	z = append(z, a[0])
	for i := 1; i < n; i++ {
		if a[i] != a[i-1] {
			z = append(z, a[i])
		}
	}
	return z
}
