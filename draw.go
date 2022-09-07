package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func fetchAssets() (assets map[string]image.Image) {
	openImage := func(name string) image.Image {
		f, e := os.Open(name + ".png")
		check(e)
		img, e := png.Decode(f)
		check(e)
		return img
	}

}