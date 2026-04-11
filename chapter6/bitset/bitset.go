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

func (bs *BitSet) Len() int {
	return bs.n
}

func (bs *BitSet) Has(x uint64) bool {
	bucketNum, bit := x/64, x%64
	return bucketNum < uint64(len(bs.data)) && bs.data[bucketNum]&(1<<bit) != 0
}

func (bs *BitSet) Add(x uint64) {

	bucketNum, bit := x/64, x%64
	for bucketNum >= uint64(len(bs.data)) {
		bs.data = append(bs.data, 0)
	}
	if bs.data[bucketNum]&(1<<bit) == 0 {
		bs.data[bucketNum] |= (1 << bit)
		bs.n++
	}

}

func (bs *BitSet) AddAll(v ...int) {
	if len(v) == 0 {
		return
	}
	for _, val := range v {
		bs.Add(uint64(val))
	}
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
	bs.data = slices.Clip(bs.data) // Release unused capacity
	bs.n = 0
}

func (bs *BitSet) Union(other *BitSet) {
	if other == nil {
		return
	}
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

func (bs *BitSet) Intersect(other *BitSet) {
	if bs == nil || other == nil {
		return
	}
	for i, item := range bs.data {
		if i >= len(other.data) {
			bs.n -= bits.OnesCount64(bs.data[i])
			bs.data[i] = 0
		} else {
			bs.n -= bits.OnesCount64(item)
			bs.data[i] &= other.data[i]
			bs.n += bits.OnesCount64(bs.data[i])
		}
	}
}

func (bs *BitSet) Difference(other *BitSet) {
	if bs == nil || other == nil {
		return
	}
	for i, item := range bs.data {
		if i >= len(other.data) {
			break
		}
		bs.n -= bits.OnesCount64(item)
		bs.data[i] &^= other.data[i] // A AND NOT(A AND B) = A AND NOT B
		bs.n += bits.OnesCount64(bs.data[i])
	}
}

func (bs *BitSet) SymetricDifference(other *BitSet) {
	if bs == nil || other == nil {
		return
	}
	for i, item := range other.data {
		if i >= len(bs.data) {
			bs.data = append(bs.data, item)
			bs.n += bits.OnesCount64(item)
		} else {
			bs.n -= bits.OnesCount64(bs.data[i])
			bs.data[i] ^= item // XOR
			bs.n += bits.OnesCount64(bs.data[i])
		}
	}
}

func (bs *BitSet) Elems() []uint64 {
	if bs == nil {
		return nil
	}
	elems := make([]uint64, 0, bs.n)
	for i, item := range bs.data {
		if item == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if (item>>j)&1 != 0 {
				elems = append(elems, uint64(64*i+j))
			}
		}
	}
	return elems
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
			}
			if buf.Len() > 1 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", 64*i+j)
		}
	}
	buf.WriteByte('}')
	return buf.String()
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
