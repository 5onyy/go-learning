package main

import (
	"fmt"
	"sort"
)

func isPalindrome(s sort.Interface) bool {
	for i, j := 0, s.Len()-1; i < j; i, j = i+1, j-1 {
		if s.Less(i, j) || s.Less(j, i) {
			return false
		}
	}
	return true
}

func main() {
	a := []int{1, 2, 3, 3, 2, 2}
	fmt.Printf("isPalindrome(a): %v\n", isPalindrome(sort.IntSlice(a)))
}
