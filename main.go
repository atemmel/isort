package main

import _ "image/jpeg"
import _ "image/png"
import "image/color"
import "image"
import "fmt"
import "os"

func luma(in color.Color) (float32) {
	r, g, b, _ := in.RGBA()
	return float32(r) * 0.2126 + float32(g) * 0.7152 + float32(b) * 0.0722
}

func openImg(fileName string) (image.Image) {
	infile, err := os.Open(fileName)
	if(err != nil) {
		panic(err)
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if(err != nil) {
		panic(err)
	}

	return src
}

func main() {
	if(len(os.Args) < 2) {
		return
	}

	src := openImg(os.Args[1])
	dim := src.Bounds().Max

	fmt.Printf("%d x %d", dim.X, dim.Y)
}
