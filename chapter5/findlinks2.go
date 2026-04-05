package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		log.Printf("Find links for url: %s\n", url)
		links, err := findlLinks(url)
		if err != nil {
			log.Printf("Find links failed %s \n", err)
		}
		for _, link := range links {
			fmt.Println(link)
		}
	}
}

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

func findlLinks(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Cannot fetch from url: %s\n", url)
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Fetch failed with status: %s - %s \n", response.StatusCode, response.Status)
		return nil, err
	}
	doc, err := html.Parse(response.Body)
	response.Body.Close()
	if err != nil {
		log.Printf("Parse response Body failed with error: %s \n", err)
	}
	return visit(nil, doc), nil
}
