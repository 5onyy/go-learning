package main

import (
	"fmt"
	"go-learning/chapter4/github"
	"log"
	"os"
)

func main() {
	res, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Number of issues: %d", res.TotalCount)
	for _, item := range res.Items {
		fmt.Printf("#%d \t User: %s \t Title: %s\n",
			item.Number, item.User.Login, item.Titttle)
	}
}
