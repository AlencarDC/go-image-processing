package photochopp

import (
	"image"
)

type ColorChannel uint8

const (
	RedChannel   ColorChannel = 0
	GreenChannel ColorChannel = 1
	BlueChannel  ColorChannel = 2
)

type Image struct {
	rgba *image.RGBA
}

func NewImage(rgba *image.RGBA) *Image {
	img := new(Image)
	img.rgba = rgba

	return img
}

func (img *Image) Width() int {
	return img.rgba.Bounds().Max.X
}

func (img *Image) Height() int {
	return img.rgba.Bounds().Max.Y
}

func (img *Image) Pixel(x, y int) []uint8 {
	begin := img.rgba.PixOffset(x, y)
	end := img.rgba.PixOffset(x, y) + 4

	return img.rgba.Pix[begin:end]
}

func (img *Image) PixelLine(y int) []uint8 {
	begin := img.rgba.PixOffset(0, y)
	end := img.rgba.PixOffset(img.Width(), y) - 1

	return img.rgba.Pix[begin:end]
}

func (img *Image) Histogram() *Histogram {
	histogram, _ := NewHistogram(img)

	return histogram
}
