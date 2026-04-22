package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, lengthInDuration("3m38s")},
	{"Go", "Moby", "Moby", 1992, lengthInDuration("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, lengthInDuration("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, lengthInDuration("4m24s")},
}

func lengthInDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

type SortByArtist []*Track // Just a trick. It is just the same slice but the comparison properties is define to compare by artist

func (a SortByArtist) Len() int           { return len(a) }
func (a SortByArtist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByArtist) Less(i, j int) bool { return a[i].Album < a[j].Artist }

// letting us define a new sort order by writing only the comparison function
// Incidentally, the concrete types that implement sort.Interface are not always slices;custom Sort is a struct type.

type customSort struct {
	t    []*Track
	less func(x, y *Track) bool // Consider less between 2 element in t
}

func (a customSort) Len() int           { return len(a.t) }
func (a customSort) Swap(i, j int)      { a.t[i], a.t[j] = a.t[j], a.t[i] } // We control what data to be compare
func (a customSort) Less(i, j int) bool { return a.less(a.t[i], a.t[j]) }   // We control what data to move, how data is moved

func main() {
	sort.Sort(SortByArtist(tracks)) // Convert to a type that satisfies sort.Interface
	printTracks(tracks)

	fmt.Printf("\n--------------------------------------------------------\n\n")
	sort.Sort(sort.Reverse(SortByArtist(tracks))) // Sort in reverse order
	printTracks(tracks)

	fmt.Printf("\n--------------------------------------------------------\n\n")
	// Sort using customSort
	sort.Sort(customSort{
		tracks,
		func(x, y *Track) bool {
			if x.Title != y.Title {
				return x.Title < y.Title
			}
			if x.Year != x.Year {
				return x.Year < y.Year
			}
			if x.Length != y.Length {
				return x.Length < y.Length
			}
			return false
		}})
	printTracks(tracks)
}
