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
			fmt.Printf("Cannot fetch from url: %s \n", url)
			continue
		}

		b, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fail to fetch")
			continue
		}
		fmt.Printf("%s\n", b)
	}
}
