package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
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

func createEmptyPngFile(target string) *os.File {
	pos := strings.LastIndex(target, ".")
	fmt.Println(target[pos:] + ".png")
	fso, err := os.Create(target[pos:] + ".png")

	if err != nil {
		panic(err)
	}
	defer fso.Close()

	return fso
}

func getDestination(startingPoint string) func(filename string) string {
	abs, err := filepath.Abs(startingPoint)
	if err != nil {
		panic(err)
	}
	destinationDir := filepath.Join(abs, "..", "converted")

	return func(filename string) string {
		rel, err := filepath.Rel(startingPoint, filename)
		if err != nil {
			panic(err)
		}
		return filepath.Join(destinationDir, rel)
	}
}

func isJpeg(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".jpg" || ext == ".jpeg"
}

func toPng(image *image.Image, target string) error {
	return png.Encode(createEmptyPngFile(target), *image)
}

func convert(curPath string, getDestination func(path string) string) error {
	if !isJpeg(curPath) {
		return nil
	}

	fmt.Println("ss")

	image, err := openImage(curPath)
	if err != nil {
		fmt.Println("failure opening image: ", err)
		return nil
	}

	return toPng(image, getDestination(curPath))
}

func walkAndConvert(path string) {
	destination := getDestination(path)

	filepath.WalkDir(path, func(curPath string, _ fs.DirEntry, _ error) error {
		return convert(curPath, destination)
	})
}

func available(target string) bool {
	fileInfo, err := os.Stat(target)

	return err != nil && fileInfo.IsDir()
}

func main() {
	flag.Parse()
	target := flag.Arg(0)

	// if !available(target) {
	// 	fmt.Println("error:", target, "is not available.")
	// 	return
	// }

	walkAndConvert(target)
}
