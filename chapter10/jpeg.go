package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Need to import png, because it register the .png format (import -> init -> main) so that the image recognizes the format to decode, otherwise the format cannot be recognized
	"io"
	"log"
	"os"
)

func main() {
	if err := toJPEG(os.Stdin, os.Stdout); err != nil {
		log.Printf("Cannot convert png to JPEG, error %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, imgType, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Input format= %s\n", imgType)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}
