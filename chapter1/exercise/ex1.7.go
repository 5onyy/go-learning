package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		response, err := http.Get(url)

		if err != nil {
			fmt.Printf("Cannot fetch from url: %v \n", err)
			continue
		}

		b, err := io.Copy(os.Stdout, response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fail to fetch")
			continue
		}
		fmt.Printf("%s\n", b)
	}
}
