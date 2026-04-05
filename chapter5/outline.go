package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Println("Parse error %v: ", err)
		os.Exit(1)
	}
	outline(nil, doc)
}

func outline(stack []string, node *html.Node) {
	if node.Type == html.ElementNode {
		stack = append(stack, node.Data)
		fmt.Println(stack)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		outline(stack, child)
	}
}
