package effects

import (
	"errors"
	"fpi/photochopp"
)

type Luminance struct{}

func (e *Luminance) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply luminance to a nil image")
	}

	width, height := img.Width(), img.Height()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.Pixel(x, y)
			r := &pixel[0]
			g := &pixel[1]
			b := &pixel[2]

			luminance := 0.299*float64(*r) + 0.587*float64(*g) + 0.114*float64(*b)

			*r = uint8(luminance)
			*g = uint8(luminance)
			*b = uint8(luminance)
		}
	}

	return nil
}
