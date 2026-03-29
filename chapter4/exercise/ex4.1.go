package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	var c1, c2 string
	fmt.Scanln(&c1)
	fmt.Scanln(&c2)

	c1_hash := sha256.Sum256([]byte(c1))
	c2_hash := sha256.Sum256([]byte(c2))

	fmt.Println(countDifferentBits(c1_hash, c2_hash))
}

func countDifferentBits(a [32]byte, b [32]byte) int {
	var ans int = 0
	for i, _ := range a {
		b1, b2 := a[i], b[i]

		for b1 > 0 {
			if b1&1 != b2&1 {
				ans++
			}
			b1 = b1 & (b1 - 1)
			b2 = b2 & (b2 - 1)
		}
	}
	return ans
}
