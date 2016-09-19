package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math"
	"os"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/BurntSushi/graphics-go/graphics/interp"
	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	frame := "/Users/indigolain/Downloads/upside_down1.jpg"
	file, err := os.Open(frame)
	if err != nil {
		panic(err)
	}

	//x, err := exif.Decode(file)
	//println(x)
	//if err != nil {
	//  panic(err)
	//}
	//lat, lng, _ := x.LatLong()
	//println(lat, lng)

	fixed, err := Process(file)
	if err != nil {
		panic(err)
	}

	out, err := os.Create("/Users/indigolain/Downloads/output.jpg")
	err = jpeg.Encode(out, fixed, nil)
	if err != nil {
		panic(err)
	}
}

var affines map[int]graphics.Affine = map[int]graphics.Affine{
	1: graphics.I,
	2: graphics.I.Scale(-1, 1),
	3: graphics.I.Scale(-1, -1),
	4: graphics.I.Scale(1, -1),
	5: graphics.I.Rotate(toRadian(90)).Scale(-1, 1),
	6: graphics.I.Rotate(toRadian(90)),
	7: graphics.I.Rotate(toRadian(-90)).Scale(-1, 1),
	8: graphics.I.Rotate(toRadian(-90)),
}

func Process(r io.Reader) (d image.Image, err error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	s, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return
	}
	o, err := ReadOrientation(bytes.NewReader(b))
	if err != nil {
		return s, nil
	}
	println(o)
	d = ApplyOrientation(s, o)
	return
}

func ReadOrientation(r io.Reader) (o int, err error) {
	e, err := exif.Decode(r)
	if err != nil {
		return
	}
	tag, err := e.Get(exif.Orientation)
	if err != nil {
		return
	}
	o, err = tag.Int(0)
	if err != nil {
		return
	}
	return
}

func ApplyOrientation(s image.Image, o int) (d draw.Image) {
	bounds := s.Bounds()
	if o >= 5 && o <= 8 {
		bounds = rotateRect(bounds)
	}
	d = image.NewRGBA64(bounds)
	affine := affines[o]
	affine.TransformCenter(d, s, interp.Bilinear)
	return
}

func toRadian(d float64) float64 {
	return math.Pi * d / 180
}

func rotateRect(r image.Rectangle) image.Rectangle {
	s := r.Size()
	return image.Rectangle{r.Min, image.Point{s.Y, s.X}}
}

func ShowImage(m image.Image) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, m, nil)
	if err != nil {
		panic(err)
	}
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println("IMAGE:" + enc)
}
