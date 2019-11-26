package main

import _ "image/jpeg"
import "image/png"
import "image/color"
import "image"
import "path/filepath"
import "strings"
import "sort"
import "fmt"
import "os"

type Pixel struct {
	data       color.Color
	storedLuma float32
}

func (l *Pixel) Luma() {
	r, g, b, _ := l.data.RGBA()
	l.storedLuma = float32(r)*0.2126 + float32(g)*0.7152 + float32(b)*0.0722
}

type ByLuma []Pixel

func (l ByLuma) Len() int {
	return len(l)
}

func (l ByLuma) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l ByLuma) Less(i, j int) bool {
	return l[i].storedLuma < l[j].storedLuma
}

func openImg(fileName string) image.Image {
	infile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if err != nil {
		panic(err)
	}

	return src
}

func writeImg(fileName string, img image.Image) {
	out, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		panic(err)
	}
}

func sortRows(img image.Image) {
	matrix := make([][]Pixel, img.Bounds().Max.Y)
	for j := 0; j < img.Bounds().Max.Y; j++ {
		row := make([]Pixel, img.Bounds().Max.X)
		for i, p := range row {
			p.data = img.At(i, j)
			p.Luma()
		}
		matrix[j] = row
	}

	for _, row := range matrix {
		sort.Sort(ByLuma(row))
	}
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	src := openImg(os.Args[1])
	dim := src.Bounds().Max
	fmt.Printf("%d x %d\n", dim.X, dim.Y)

	fmt.Printf("Sorting pixels...\n")
	fmt.Printf("First pixels: ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", src.At(i, 0))
	}
	sortRows(src)

	out := strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1])) + "_sorted.png"
	fmt.Printf("\nWriting file data to %s...\n", out)
	writeImg(out, src)
}
