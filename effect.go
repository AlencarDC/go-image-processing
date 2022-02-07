package photochopp

import (
	"errors"
	"fmt"
	"image"
)

type Effect interface {
	Apply(img image.RGBA) (err error)
}

type VerticalFlip struct{}

func (e *VerticalFlip) Apply(img *Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot flip vertically a nil image")
	}

	width, height := img.Width(), img.Height()

	lineBuff := make([]uint8, 4*width)
	for y := 0; y < height/2; y++ {
		upperLine := img.PixelLine(y)
		lowerLine := img.PixelLine(height - y - 1)

		copy(lineBuff, upperLine)
		copy(upperLine, lowerLine)
		copy(lowerLine, lineBuff)
	}

	return nil
}

type HorizontalFlip struct{}

func (e *HorizontalFlip) Apply(img *Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot flip horizontally a nil image")
	}

	width, height := img.Width(), img.Height()

	pixelBuff := make([]uint8, 4)
	for x := 0; x < width/2; x++ {
		for y := 0; y < height; y++ {
			leftPixel := img.Pixel(x, y)
			rightPixel := img.Pixel(width-x-1, y)

			copy(pixelBuff, leftPixel)
			copy(leftPixel, rightPixel)
			copy(rightPixel, pixelBuff)
		}
	}

	return nil
}

type Luminance struct{}

func (e *Luminance) Apply(img *Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply luminance to a nil image")
	}

	width, height := img.Width(), img.Height()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.Pixel(x, y)
			r := &pixel[0]
			g := &pixel[+1]
			b := &pixel[+2]

			luminance := 0.299*float64(*r) + 0.587*float64(*g) + 0.114*float64(*b)

			*r = uint8(luminance)
			*g = uint8(luminance)
			*b = uint8(luminance)
		}
	}

	return nil
}

type ColorQuantization struct {
	NumberOfDesiredColors int
}

func (cq *ColorQuantization) binSizes(histogram *Histogram) []int {
	channels := [...]ColorChannel{RedChannel, GreenChannel, BlueChannel}
	var binSizes []int
	for _, channel := range channels {
		lowest, _ := histogram.LowestColor(channel)
		highest, _ := histogram.HighestColor(channel)

		intensitySize := int(highest) - int(lowest) + 1

		binSize := intensitySize / cq.NumberOfDesiredColors
		fmt.Println(lowest, highest, intensitySize, binSize)
		binSizes = append(binSizes, binSize)
	}

	return binSizes
}

func (cq *ColorQuantization) newColor(currentColor uint8, binSize int) uint8 {
	newBinBegin := (int(currentColor) / binSize) * binSize

	return uint8(newBinBegin + binSize/2)
}

func (cq *ColorQuantization) Apply(img *Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply color quantization to a nil image")
	}

	histogram := img.Histogram()
	binSizes := cq.binSizes(histogram)

	width, height := img.Width(), img.Height()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.Pixel(x, y)
			pixel[0] = cq.newColor(pixel[0], binSizes[0])
			pixel[1] = cq.newColor(pixel[1], binSizes[1])
			pixel[2] = cq.newColor(pixel[2], binSizes[2])
		}
	}

	return nil
}
