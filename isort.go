package main

import _ "image/jpeg"
import "image/png"
import "image/color"
import "image"
import "path/filepath"
import "strings"
import "sort"
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

func writeImg(fileName string, rgba *image.RGBA) {
	out, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, rgba)
	if err != nil {
		panic(err)
	}
}

func sortRows(img image.Image) (*image.RGBA) {
	matrix := make([][]Pixel, img.Bounds().Max.Y)
	for j, row := range matrix {
		row = make([]Pixel, img.Bounds().Max.X)
		for i, pixel := range row {
			pixel.data = img.At(i, j)
			pixel.Luma()
			row[i] = pixel
		}
		matrix[j] = row
	}

	for _, row := range matrix {
		sort.Sort(ByLuma(row))
	}

	rgba := image.NewRGBA(img.Bounds())

	for j, row := range matrix {
		for i, pixel := range row {
			rgba.Set(i, j, pixel.data)
		}
	}

	return rgba
}

func sortCols(img image.Image) (*image.RGBA) {
	matrix := make([][]Pixel, img.Bounds().Max.X)
	for j, row := range matrix {
		row = make([]Pixel, img.Bounds().Max.Y)
		for i, pixel := range row {
			pixel.data = img.At(j, i)
			pixel.Luma()
			row[i] = pixel;
		}
		matrix[j] = row
	}

	for _, row := range matrix {
		sort.Sort(ByLuma(row))
	}

	rgba := image.NewRGBA(img.Bounds())

	for j, row := range matrix {
		for i, pixel := range row {
			rgba.Set(j, i, pixel.data)
		}
	}

	return rgba
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	src := openImg(os.Args[1])

	rgba := sortCols(src)
	outCols := strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1])) + "_cs.png"
	writeImg(outCols, rgba)

	rgba = sortRows(src)
	outRows := strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1])) + "_rs.png"
	writeImg(outRows, rgba)
}
