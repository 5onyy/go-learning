package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"slices"
)

type BitSet struct {
	n    int
	data []uint64
}

func (bs *BitSet) Has(x uint64) bool {
	bucketNum, bit := x/64, x%64
	return bucketNum < uint64(len(bs.data)) && bs.data[bucketNum]&(1<<bit) != 0
}

func (bs *BitSet) Add(x uint64) {
	if bs.Has(x) {
		return
	}
	bs.n++
	bucketNum, bit := x/64, x%64
	for bucketNum >= uint64(len(bs.data)) {
		bs.data = append(bs.data, 0)
	}
	bs.data[bucketNum] |= (1 << bit)
}

func (bs *BitSet) Union(other *BitSet) {
	for i, item := range other.data {
		if i >= len(bs.data) {
			bs.data = append(bs.data, item)
			bs.n += bits.OnesCount64(item)
		} else {
			bs.n -= bits.OnesCount64(bs.data[i])
			bs.data[i] |= item
			bs.n += bits.OnesCount64(bs.data[i])
		}
	}
}

func (bs *BitSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, item := range bs.data {
		if item == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if (item>>j)&1 == 0 {
				continue
			} else {
				if buf.Len() > 1 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (bs *BitSet) Len() int {
	return bs.n
}

func (bs *BitSet) Remove(x uint64) {
	if !bs.Has(x) {
		return
	}
	bs.n--
	bucket, bit := x/64, x%64
	bs.data[bucket] &= ^(1 << bit)
}

func (bs *BitSet) Clear() {
	clear(bs.data)
	bs.data = slices.Clip(bs.data) // Release unsed capacity
	bs.n = 0
}

func (bs *BitSet) Copy() *BitSet {
	if bs == nil {
		return nil
	}
	return &BitSet{
		data: slices.Clone(bs.data),
		n:    bs.n,
	}
}

func main() {
	var x BitSet
	x.Add(5)
	x.Add(100)
	x.Add(64)
	fmt.Printf("x: %v\n", x.String())

	var y BitSet
	y.Add(2)

	x.Union(&y)
	fmt.Printf("x: %v\n", x.String())
	fmt.Printf("x.Len(): %v\n", x.Len())

	x.Remove(5)
	fmt.Printf("x: %v\n", x.String())
	fmt.Printf("x.Len(): %v\n", x.Len())

	// x.Clear()
	// fmt.Printf("x: %v\n", x.String())
	// fmt.Printf("x.Len(): %v\n", x.Len())

	clone := x.Copy()
	fmt.Printf("clone.String(): %v\n", clone.String())
}
