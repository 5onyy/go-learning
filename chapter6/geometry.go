package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

// Traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Method of Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)

}

// Receiver name does not neet to be the same
func (q Point) DistanceManhattan(p Point) int {
	return int(math.Abs(q.X-p.X) + math.Abs(q.Y-p.Y))
}

// If we want to modify states of the struct or to avoid copying
// In a realistic program, convention dictates that if any method of Point has a pointer receiver, then all methods of Point should have a pointer receiver, even ones that don’t strictly need it.
// In this code, we will break this rule for leanrning purpose
func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

type Path []Point

// We cannot attach method to primitive type, instead, create a name for it
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

func main() {
	p := Point{1, 2}
	q := Point{4, 6}
	fmt.Println(q.Distance(p))
	fmt.Println(q.DistanceManhattan(p))

	perim := Path{
		{1, 1},
		{5, 1},
		{5, 4},
		{1, 1},
	}
	fmt.Println(perim.Distance())

	r := &Point{1, 2}
	r.ScaleBy(2)
	fmt.Println(*r)

	// Implicit call
	p.ScaleBy(2) // &p.ScaleBy(2)
	fmt.Printf("p: %v\n", p)

	// Implicit convert if the receiver parameter needs a *Point
	pptr := &p // (*p).ScaleBy(2)
	fmt.Printf("pptr.Distance(q): %v\n", pptr.Distance(q))

}
