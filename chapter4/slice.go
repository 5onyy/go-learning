package main

import "fmt"

func main() {
	// A simple way to rotate a slice left by n elements is to apply the reverse function three times, first to the leading n elements, then to the remaining elements, and finally to the whole slice.
	// (To rotate to the right, make the third call first.)
	s := []int{0, 1, 2, 3, 4, 5}
	// reverse(s[:2])
	// reverse(s[2:])
	// reverse(s)
	fmt.Println(cap(s), len(s))

	// append will allocate new memory if the len exceed capacity --> original array wont see the result

	// 	Updating the slice variable is required not just when calling append,but for any function that may change the length or capacity of a slice or make it refer to a different underlying array

	// s := make([]int, 0, 10)
	// s = append(s, 0, 1, 2, 3, 4, 5)
	fmt.Println(myAppend(s, 5))
	fmt.Println(cap(s), len(s))
	fmt.Println(s[:cap(s)])

	fmt.Println(remove(s, 3))
	fmt.Println(cap(s), len(s))
	fmt.Println(s[:cap(s)])
}

func myAppend(sl []int, v int) []int {
	return append(sl, v)
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func remove(sl []int, i int) []int {
	copy(sl[i:], sl[i+1:])
	return sl[:]
}
