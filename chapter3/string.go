package main

import "fmt"

func main() {
	s := "abc"
	b := []byte(s)
	b[2] = 'z'
	fmt.Println(b)
	fmt.Println(s)
}
