package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

const (
	IMAGE_DIRECTORY_INPUT  = "./input"
	IMAGE_DIRECTORY_OUTPUT = "./output"
	JPEG_SUFFIX            = ".jpeg"
	IMAGE_WIDTH            = 100
	IMAGE_HEIGHT           = 80
)

func main() {
	itemsOfInputDirectory, err := ioutil.ReadDir(IMAGE_DIRECTORY_INPUT)
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range itemsOfInputDirectory {
		if !item.IsDir() && strings.HasSuffix(item.Name(), JPEG_SUFFIX) {
			if err := grayscaleAndResizeAndSave(item); err != nil {
				fmt.Printf("item with name '%s' failed", item.Name())
				log.Fatalln(err)
			}
		}
	}
}

func grayscaleAndResizeAndSave(item fs.FileInfo) error {
	inputFilePath := path.Join(IMAGE_DIRECTORY_INPUT, item.Name())
	img, err := imgio.Open(inputFilePath)
	if err != nil {
		return err
	}

	grayed := effect.Grayscale(img)
	resized := transform.Resize(grayed, IMAGE_WIDTH, IMAGE_HEIGHT, transform.CatmullRom)

	outputFilePath := path.Join(IMAGE_DIRECTORY_OUTPUT, item.Name()[:len(item.Name())-len(JPEG_SUFFIX)]+".png")
	if err := imgio.Save(outputFilePath, resized, imgio.PNGEncoder()); err != nil {
		return err
	}
	return nil
}
