package effects

import (
	"errors"
	"fpi/photochopp"
)

type VerticalFlip struct{}

func (e *VerticalFlip) Apply(img *photochopp.Image) (err error) {
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
