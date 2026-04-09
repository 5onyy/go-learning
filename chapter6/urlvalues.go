package main

import (
	"fmt"
	"net/url"
)

func main() {
	m := url.Values{"lang": {"en"}} // direct construction
	m.Add("item", "1")
	m.Add("item", "2")

	fmt.Printf("m.Get(\"lang\"): %v\n", m.Get("lang"))
	fmt.Printf("m.Get(\"q\"): %v\n", m.Get("q"))
	fmt.Printf("m.Get(\"item\"): %v\n", m.Get("item")) // First value
	fmt.Printf("m[\"item\"]: %v\n", m["item"])         // Direct map access

	m = nil
	fmt.Printf("m[\"item\"]: %v\n", m["item"])
	m.Add("item", "3") // Assignment to entry in nil map

}
