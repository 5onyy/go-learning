package main

import "fmt"

func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}

func main() {
	f := squares()
	g := squares()
	fmt.Println(f()) // 1
	fmt.Println(f()) // 4
	fmt.Println(f()) // 9
	fmt.Println(f()) // 16
	fmt.Println(f()) // 25
	fmt.Println(g()) // 1

	// The squares example demonstrates that function values are not just code but can have state.
	//  Here again we see an example where the lifet imeofavar iable isnot deter mined byits scope: the variable x exists after squares has returned within main,even though x is hidden inside f
}
