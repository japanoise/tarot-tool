package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

func main() {
	glue := flag.Int("g", 30, "pixels of spacing glue between edges")
	flag.Parse()
	images := make([]image.Image, 0, 3)
	for _, arg := range flag.Args() {
		images = append(images, readimg(arg))
	}
	canvas, cw, _ := gencanvas(images, *glue)
	for i, img := range images {
		anchor := ((i + 1) * *glue) + (i * cw)
		for x := 0; x < img.Bounds().Dx(); x++ {
			for y := 0; y < img.Bounds().Dy(); y++ {
				canvas.Set(anchor+x, *glue+y, img.At(x, y))
				//log.Print(canvas.At(anchor + x, *glue + y))
			}
		}
	}
	writeimg(canvas, fmt.Sprintf("%d.%s", int32(time.Now().Unix()), "png"))
}

func gencanvas(images []image.Image, glue int) (image.RGBA, int, int) {
	cardwidth, cardheight := 0, 0
	for _, img := range images {
		if img.Bounds().Dx() > cardwidth {
			cardwidth = img.Bounds().Dx()
		}
		if img.Bounds().Dy() > cardheight {
			cardheight = img.Bounds().Dy()
		}
	}
	l := len(images)
	return *image.NewRGBA(image.Rect(0, 0, ((l+1)*glue)+(l*cardwidth), (2*glue)+cardheight)), cardwidth, cardheight
}

func readimg(filename string) image.Image {
	reader, err1 := os.Open(filename)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer reader.Close()
	m, _, err2 := image.Decode(reader)
	if err2 != nil {
		log.Fatal(err2)
	}
	return m
}

func writeimg(m image.RGBA, filename string) {
	writer, err := os.Create(filename)
	defer writer.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(writer, &m)
	if err != nil {
		log.Fatal(err)
	}
}
