package photochopp

import (
	"image"
	"image/draw"
	"os"
)

type ColorChannel uint8

const (
	RedChannel   ColorChannel = 0
	GreenChannel ColorChannel = 1
	BlueChannel  ColorChannel = 2
)

type Image struct {
	img  *image.Image
	rgba *image.RGBA
}

func NewImageFromFilePath(path string) (*Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	return &Image{img: &img, rgba: rgba}, nil
}

func NewImage(rgba *image.RGBA) *Image {
	img := new(Image)
	img.rgba = rgba

	return img
}

func (img *Image) Image() *image.Image {
	return img.img
}

func (img *Image) RGBA() *image.RGBA {
	return img.rgba
}

func (img *Image) Width() int {
	return img.rgba.Bounds().Max.X
}

func (img *Image) Height() int {
	return img.rgba.Bounds().Max.Y
}

func (img *Image) ImageFromRGBA() image.Image {
	return img.rgba.SubImage(img.rgba.Bounds())
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
