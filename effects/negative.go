package effects

import (
	"errors"
	"fpi/photochopp"
)

type Negative struct{}

func (b *Negative) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply negative to a nil image")
	}

	width, height := img.Width(), img.Height()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.Pixel(x, y)
			pixel[0] = uint8(255 - pixel[0])
			pixel[1] = uint8(255 - pixel[1])
			pixel[2] = uint8(255 - pixel[2])
		}
	}

	return nil
}
