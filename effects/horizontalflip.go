package effects

import (
	"errors"
	"fpi/photochopp"
)

type HorizontalFlip struct{}

func (e *HorizontalFlip) Apply(img *photochopp.Image) (err error) {
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
