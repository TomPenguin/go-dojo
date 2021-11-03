package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
)

func openImage(path string) (*image.Image, error) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("open failure: ", err)
		return nil, err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

type Filename struct {
	name string
	ext  string
}

func getFileName(path string) Filename {
	pos := strings.LastIndex(path, ".")
	return Filename{path[:pos], path[pos+1:]}
}

func jpgToPng(image *image.Image, filename Filename) error {
	fso, err := os.Create("./" + filename.name + ".png")
	if err != nil {
		panic(err)
	}
	defer fso.Close()

	return png.Encode(fso, *image)
}

func convert(path string) error {
	filename := getFileName(path)
	if filename.ext != "jpg" && filename.ext != "jpeg" {
		return nil
	}

	image, err := openImage(path)
	if err != nil {
		fmt.Println("failure opening image: ", err)
		return nil
	}

	return jpgToPng(image, filename)
}

func main() {
	flag.Parse()
	target := flag.Arg(0)

	fmt.Println(getFileName(target))
	convert(target)
}
