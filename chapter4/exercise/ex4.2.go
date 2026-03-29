package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

type HASH_TYPE int

const (
	SHA_256 HASH_TYPE = iota
	SHA_384
	SHA_512
)

func main() {
	var data string
	fmt.Scanln(&data)
	calculateHash(data, SHA_512)
}

func calculateHash(data string, hashType ...HASH_TYPE) {
	_hashType := SHA_256
	if len(hashType) > 0 {
		_hashType = hashType[0]
	}
	switch {
	case _hashType == SHA_512:
		fmt.Printf("%x\n", sha512.Sum512([]byte(data)))
	case _hashType == SHA_384:
		fmt.Printf("%x\n", sha512.Sum384([]byte(data)))
	default:
		fmt.Printf("%x\n", sha256.Sum256([]byte(data)))
	}
}
