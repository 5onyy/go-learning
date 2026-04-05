package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {

}

func title(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	ct := response.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}

	_, err = html.Parse(response.Body)
	if err != nil {
		return fmt.Errorf("Parsing %s as HTML failed: %v", url, err)
	}
	return nil
}
