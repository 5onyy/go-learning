package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Movie struct {
	Title  string
	Year   int  `json:"released"`        // alternative json name for field Year
	Color  bool `json:"color,omitempty"` // omit if 0 value / false
	Actors []string
}

var movies = []Movie{
	{
		Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"},
	},
	{
		Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"},
	},
	{
		Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"},
	},
}

func main() {
	data, err := json.Marshal(movies)
	data, err = json.MarshalIndent(movies, "", "\t")

	var decoded []Movie
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		log.Fatalf("JSON Unmarshaling failed: %s", err)
	}
	fmt.Printf("Decoded data: %+v\n", decoded)
	fmt.Printf("Original JSON:\n%s\n", data)
}
