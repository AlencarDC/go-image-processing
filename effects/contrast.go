package effects

import (
	"errors"
	"fpi/photochopp"
	"math"
)

type Contrast struct {
	Value int
}

func (c *Contrast) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply contrast to a nil image")
	}

	value := int(math.Max(1, math.Min(float64(c.Value), 255)))

	width, height := img.Width(), img.Height()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.Pixel(x, y)
			pixel[0] = uint8(math.Max(0, math.Min(float64(int(pixel[0])*value), 255)))
			pixel[1] = uint8(math.Max(0, math.Min(float64(int(pixel[1])*value), 255)))
			pixel[2] = uint8(math.Max(0, math.Min(float64(int(pixel[2])*value), 255)))
		}
	}

	return nil
}
