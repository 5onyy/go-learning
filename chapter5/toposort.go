package main

import (
	"fmt"
	"sort"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":  {"discrete math"},
	"databases":        {"data structures"},
	"discrete math":    {"intro to programming"},
	"formal languages": {"discrete math"},
	"networks":         {"operating systems"},
	"operating systems": {
		"data structures",
		"computer organization",
	},
	"programming languages": {
		"data structures",
		"computer organization",
	},
}

func topoSort(prereqs map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)

	// When ananonymous function requires recursion, as in this example, we must first de clare a variable,and then assign the anonymous function to that var able
	visitAll = func(items []string) {
		for _, key := range items {
			if !seen[key] {
				seen[key] = true
				visitAll(prereqs[key])
				order = append(order, key)
			}
		}
	}

	var keys []string
	for key, val := range prereqs {
		fmt.Println(key, ":", val)
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	fmt.Println(keys)
	return order
}

func main() {
	for i, item := range topoSort(prereqs) {
		fmt.Printf("%d: %s\n", i+1, item)
	}
}
