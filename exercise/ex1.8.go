package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		var clean_url string = url
		if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
			clean_url = "https://" + url
		}

		response, err := http.Get(clean_url)

		if err != nil {
			continue
		}

		io.Copy(os.Stdout, response.Body)
		response.Body.Close()
	}
}
