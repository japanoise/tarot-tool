package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

func main() {
	reader, err1 := os.Open("unknown.jpg")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer reader.Close()
	m, _, err2 := image.Decode(reader)
	if err2 != nil {
		log.Fatal(err2)
	}
	writeimg(m, "img.png")
}

func writeimg(m image.Image, filename string) {
	writer, err := os.Create(filename)
	defer writer.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(writer, m)
	if err != nil {
		log.Fatal(err)
	}
}
