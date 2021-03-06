package photochopp

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

type ColorChannel uint8

const (
	GrayChannel  ColorChannel = 0
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

func (img *Image) SetRGBA(rgba *image.RGBA) {
	img.rgba = rgba
}

func (img *Image) IsEmpty() bool {
	return img.rgba == nil
}

func (img *Image) Width() int {
	return img.rgba.Bounds().Max.X
}

func (img *Image) SetWidth(width int) {
	img.rgba.Rect.Max.X = width
	img.rgba.Stride = width * 4
}

func (img *Image) Height() int {
	return img.rgba.Bounds().Max.Y
}

func (img *Image) SetHeight(height int) {
	img.rgba.Rect.Max.Y = height
}

func (img *Image) IsValidPosition(x, y int) bool {
	return x >= 0 && x < img.Width() && y >= 0 && y < img.Height()
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

func (img *Image) Copy() *Image {
	buffer := new(bytes.Buffer)
	err := jpeg.Encode(buffer, img.rgba, nil)
	if err != nil {
		return nil
	}

	imgCopy, _, err := image.Decode(buffer)
	if err != nil {
		return nil
	}

	bounds := imgCopy.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, imgCopy, bounds.Min, draw.Src)

	return &Image{img: &imgCopy, rgba: rgba}
}

func (img *Image) IsGrayScale() bool {
	width, heigth := img.Width(), img.Height()

	for i := 0; i < width; i += 1 {
		for j := 0; j < heigth; j += 1 {
			pixel := img.Pixel(i, j)
			if pixel[0] != pixel[1] || pixel[1] != pixel[2] {
				return false
			}
		}
	}

	return true
}
