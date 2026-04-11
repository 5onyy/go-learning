package main

import (
	"fmt"
	"image/color"
	"math"
	"sync"
)

type Point struct{ X, Y float64 }

func (p *Point) Distance(q Point) float64 {
	return math.Hypot(p.X-q.X, p.Y-q.Y)
}

func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

type ColoredPoint struct {
	Point // Embeded, not inheritance, ColoredPoint is not a Point
	Color color.RGBA
}

// Methods can be declared only on named types (like Point) and pointers to them(*Point), but thanks to embedding, it’s possible and sometimes useful for unnamed struct types to have methods too.

var cache = struct { // This defines a struct without giving it a name.
	sync.Mutex
	mapping map[string]string
}{ // mmediate initialization
	mapping: make(map[string]string),
}

// Use when we didn't need to use this struct again

func lookUp(key string) string {
	cache.Lock()
	v := cache.mapping[key]
	cache.Unlock()
	return v
}

func main() {
	var cp ColoredPoint
	cp.X = 1
	fmt.Printf("cp.X: %v\n", cp.X)
	cp.Point.Y = 2
	fmt.Printf("cp.Y: %v\n", cp.Y)

	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	var p = ColoredPoint{Point{1, 4}, red}
	var q = ColoredPoint{Point{5, 4}, blue}

	fmt.Printf("p.Distance(q.Point): %v\n", p.Distance(q.Point)) // Must call distance to q.Point, because the Distance method only accept Point type
	p.ScaleBy(2)
	q.ScaleBy(2)
	fmt.Printf("p.Point: %v\n", p.Point)
	fmt.Printf("q.Point: %v\n", q.Point)

	ppoint := Point{1, 2}
	qpoint := Point{4, 6}

	distancefromP := ppoint.Distance // method value <- a function that binds a method Point.Distance to receiver distancefromP

	fmt.Printf("distancefromP(qpoint): %v\n", distancefromP(qpoint)) // can be invoke without receiver value

}
