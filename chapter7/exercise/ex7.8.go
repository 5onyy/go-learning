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

// Multi-tier sort, solution 1
// ---------------------------------------------------------------------------
type multiTierSort struct {
	t    []*Track
	less []func(a, b *Track) bool //order by priority, highest first
}

func (a multiTierSort) Len() int      { return len(a.t) }
func (a multiTierSort) Swap(i, j int) { a.t[i], a.t[j] = a.t[j], a.t[i] }
func (a multiTierSort) Less(i, j int) bool {
	t_i, t_j := a.t[i], a.t[j]
	// We loop through all less function and stop at the function that is highest priority and can determine the order of elements
	// In simple terms, it is similar to writing a comparison function, but this case it is dynamic
	for _, lessFunc := range a.less {
		if lessFunc(t_i, t_j) == true { // if t_i can be placed before t_j
			return true
		}
		if lessFunc(t_j, t_i) == true {
			return false
		}
	}
	return false
}

var byYear = func(a, b *Track) bool { return a.Year < b.Year }

var byLength = func(a, b *Track) bool { return a.Length < b.Length }

var byArtist = func(a, b *Track) bool { return a.Artist < b.Artist }

var byAlbum = func(a, b *Track) bool { return a.Album < b.Album }

var byTitle = func(a, b *Track) bool { return a.Title < b.Title }

func addSortKey(sortKeys []func(a, b *Track) bool, key func(a, b *Track) bool) []func(a, b *Track) bool {
	return append([]func(a, b *Track) bool{key}, sortKeys...)
	// sortKeys... unpacks (expands) a slice into individual elements.
}

func demoSolution1() {
	var sortKeys []func(a, b *Track) bool
	sortKeys = addSortKey(sortKeys, byTitle)
	sortKeys = addSortKey(sortKeys, byLength)
	sortKeys = addSortKey(sortKeys, byYear)

	// Priority: Year -> Length -> Tittle

	sort.Sort(multiTierSort{
		tracks,
		sortKeys,
	})
	printTracks(tracks)
}

// --------------------------------------------------------------

// Sort.Stable, Solution 2

func demoSolution2() {
	// Priority: Year -> Length -> Tittle

	sort.SliceStable(
		tracks,
		func(i, j int) bool {
			return tracks[i].Title < tracks[j].Title
		})

	sort.SliceStable(
		tracks,
		func(i, j int) bool {
			return tracks[i].Length < tracks[j].Length
		})

	sort.SliceStable(
		tracks,
		func(i, j int) bool {
			return tracks[i].Year < tracks[j].Year
		})

	printTracks(tracks)

}

func main() {
	demoSolution1()
	fmt.Printf("\n--------------------------------------------------------------\n\n")
	demoSolution2()
}
