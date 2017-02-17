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
	cc := flag.Bool("c", false, "use a celtic cross layout")
	flag.Parse()
	images := make([]image.Image, 0, 3)
	for _, arg := range flag.Args() {
		images = append(images, readimg(arg))
	}
	var m image.RGBA
	if *cc {
		m = ccspread(images, glue)
	} else {
		m = linearspread(images, glue)
	}
	writeimg(m, fmt.Sprintf("%d.%s", int32(time.Now().Unix()), "png"))
}

func ccspread(images []image.Image, glue *int) image.RGBA {
	for len(images) < 10 {
		images = append(images, getunknown())
	}
	cardwidth, cardheight := getcwch(images)
	g := *glue
	canvas := *image.NewRGBA(image.Rect(0, 0, (3*cardwidth)+cardheight+(5*g), (5*g)+(4*cardheight)))
	cssize := (cardheight - cardwidth) / 2
	canchx := (2 * g) + cardwidth + cssize
	canchy := g + (cardheight / 2)
	incy := cardheight + g
	/* central area of the spread */
	copycard(images[4], canvas, canchx, canchy)
	copycard(images[0], canvas, canchx, canchy+incy)
	copycard(images[2], canvas, canchx, canchy+(2*incy))
	copycard(images[3], canvas, g, canchy+incy)
	copycard(rotatecard(images[1]), canvas, canchx-cssize, canchy+incy+cssize)
	copycard(images[5], canvas, canchx+cardwidth+cssize+g, canchy+incy)
	/* rod to the right of the spread */
	ranchx := canchx + (2 * cardwidth) + cssize + (3 * g)
	copycard(images[9], canvas, ranchx, g)
	copycard(images[8], canvas, ranchx, g+incy)
	copycard(images[7], canvas, ranchx, g+(2*incy))
	copycard(images[6], canvas, ranchx, g+(3*incy))
	return canvas
}

func rotatecard(img image.Image) image.Image {
	retval := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dy(), img.Bounds().Dx()))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			retval.Set(y, x, img.At(x, y))
		}
	}
	return retval
}

func copycard(img image.Image, canvas image.RGBA, anchorx, anchory int) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			canvas.Set(anchorx+x, anchory+y, img.At(x, y))
		}
	}
}

func linearspread(images []image.Image, glue *int) image.RGBA {
	canvas, cw, _ := gencanvas(images, *glue)
	for i, img := range images {
		copycard(img, canvas, ((i+1)**glue)+(i*cw), *glue)
	}
	return canvas
}

func getcwch(images []image.Image) (int, int) {
	cardwidth, cardheight := 0, 0
	for _, img := range images {
		if img.Bounds().Dx() > cardwidth {
			cardwidth = img.Bounds().Dx()
		}
		if img.Bounds().Dy() > cardheight {
			cardheight = img.Bounds().Dy()
		}
	}
	return cardwidth, cardheight
}

func gencanvas(images []image.Image, glue int) (image.RGBA, int, int) {
	cardwidth, cardheight := getcwch(images)
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
