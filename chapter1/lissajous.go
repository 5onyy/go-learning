package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var cycles int = 5

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if param := r.URL.Query(); param.Has("cycles") {
			val, err := strconv.Atoi(param.Get("cycles"))
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid cycles parameter: %q - %s", param.Get("cycles"), err.Error()), http.StatusBadRequest)
				return
			}
			cycles = val
		}
		lissajous(w)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(w http.ResponseWriter) {
	fmt.Fprintf(w, "cycles = %d", cycles)
}
