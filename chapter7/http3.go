package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

var db database = database{"shoes": 50, "socks": 5}

func listHandler(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s %s:\n", item, price)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/list", listHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
