package main

import "fmt"

func main() {
	var a [3]int
	fmt.Println(a)

	// if an ellipsis ‘‘...’’ appears inp lace of the length, the array length is determined by the number of initializers
	q := [...]int{1, 2, 3, 4, 5}
	fmt.Println(q)

	// defines an array r with 100 elements, all zero exceptfor the last, which has value −1
	r := [...]int{15: -1}
	fmt.Println(r)

	c1 := [2]int{1, 2}
	c2 := [...]int{1, 2}
	c3 := [2]int{1, 3}
	fmt.Println(c1 == c2, c2 == c3, c1 == c3)
	d := [3]int{1, 2}
	fmt.Println(d)
	// fmt.Println(c1 == d)  Compiles error cannot compare [2]int == [3]int
}
