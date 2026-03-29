package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()

	// Shared channel all goroutines send their results through
	ch := make(chan string)

	// Launch one goroutine per URL — all fetches run concurrently
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	// Collect exactly one result per goroutine; blocks until each sends
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}

	// Total elapsed time ≈ slowest fetch, not the sum of all fetches
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// fetch makes an HTTP GET request and sends a result summary to ch.
// ch is send-only (chan<- string) to enforce that fetch never reads from it.
func fetch(url string, ch chan<- string) {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		// Send the error as the result so main can print it
		ch <- fmt.Sprintf("while getting %s: %v", url, err)
		return
	}
	// Always close the body to free the underlying TCP connection
	defer resp.Body.Close()

	// Drain the body into /dev/null — we only care about byte count and timing
	nbytes, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()

	// Send formatted result: elapsed time, bytes received, URL
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
