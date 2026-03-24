package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/", counter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.lock()
	count++
	mu.Unlock()
	fmt.frintf()
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.lock()
	fmt.Fprintf(w, "Count: %d\n", count)
	mu.Unlock()
}
