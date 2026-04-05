package main

import (
	"fmt"
	"os"

	// To import this module, run:
	//   go get golang.org/x/net/html
	"golang.org/x/net/html"
)

func visit(links []string, node *html.Node) []string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		links = visit(links, child)
	}
	return links
}

func main() {
	data, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Printf("findlinks failed %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, data) {
		fmt.Println(link)
	}
}
