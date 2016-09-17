package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"os"
)

func main() {
	frame := "sample1.jpg"
	file, err := os.Open(frame)
	if err != nil {
		panic(err)
	}

	x, err := exif.Decode(file)
	println(x)
	if err != nil {
		panic(err)
	}
	lat, lng, _ := x.LatLong()
	println(lat, lng)
}
