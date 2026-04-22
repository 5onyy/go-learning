package main

import (
	"fmt"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

// main demonstrates the usage of ByteCounter, a custom type that implements
// the io.Writer interface to count the number of bytes written to it.
//
// First, it writes "Hello" directly using the Write method and prints the count.
//
// Then, it resets the counter and uses fmt.Fprintf to write a formatted string
// through the ByteCounter. Note that fmt.Fprintf requires an io.Writer, which
// is defined as:
//
//	type Writer interface {
//	    Write(p []byte) (n int, err error)
//	}
//
// We pass &c (a pointer to c) instead of c because the Write method has a
// pointer receiver (*ByteCounter). In Go, a method with a pointer receiver
// is only in the method set of the pointer type (*ByteCounter), not the
// value type (ByteCounter). Therefore, only *ByteCounter satisfies the
// io.Writer interface, and we must pass &c to fmt.Fprintf. Passing c
// directly would result in a compile-time error:
//
//	"ByteCounter does not implement io.Writer (Write method has pointer receiver)"
func main() {
	var c ByteCounter
	c.Write([]byte("Hello"))
	fmt.Printf("c: %v\n", c)

	c = 0
	var name = "Dolly"
	// You pass &c because:
	// The Write method is defined on *ByteCounter (pointer receiver)
	// Only *ByteCounter implements io.Writer
	// &c has type *ByteCounter
	fmt.Fprintf(&c, "Hello, %s", name)
	fmt.Printf("c: %v\n", c)
}
