package effects

import (
	"errors"
	"fpi/photochopp"
	"math"
)

type Convolve struct {
	ShouldEmboss bool
	Kernel       [][]float32
}

func (c *Convolve) Apply(img *photochopp.Image) (err error) {
	if !c.isValidKernel() {
		return errors.New("convolve: provided kernel is invalid")
	}

	width, height := img.Width(), img.Height()
	kernelHeight := len(c.Kernel)
	kernelWidth := len(c.Kernel[0])
	imgCopy := img.Copy()

	for x := kernelWidth / 2; x < width-(kernelWidth/2); x += 1 {
		for y := kernelHeight / 2; y < height-(kernelHeight/2); y += 1 {
			convolvedPixel, err := c.convolvePixel(imgCopy, x, y)
			if err != nil {
				continue
			}

			if c.ShouldEmboss {
				convolvedPixel = c.embossPixel(convolvedPixel)
			}

			newPixel := c.clampPixel(convolvedPixel)

			pixel := img.Pixel(x, y)
			pixel[0] = newPixel[0]
			pixel[1] = newPixel[1]
			pixel[2] = newPixel[2]
		}
	}

	return nil
}

func (c *Convolve) convolvePixel(img *photochopp.Image, x, y int) ([3]float32, error) {
	width, height := img.Width(), img.Height()
	kernelHeight := len(c.Kernel)
	kernelWidth := len(c.Kernel[0])

	if !c.isValidPixel(x, y, width, height) {
		return [3]float32{}, errors.New("convolve: provided position is invalid")
	}

	var convolvedPixel [3]float32
	for i := -kernelHeight / 2; i <= kernelHeight/2; i += 1 {
		for j := -kernelWidth / 2; j <= kernelWidth/2; j += 1 {
			pixel := img.Pixel(x+j, y+i)
			kernelValue := c.Kernel[kernelHeight/2-i][kernelWidth/2-j]

			convolvedPixel[0] += float32(pixel[0]) * kernelValue
			convolvedPixel[1] += float32(pixel[1]) * kernelValue
			convolvedPixel[2] += float32(pixel[2]) * kernelValue
		}
	}
	return convolvedPixel, nil
}

func (c *Convolve) embossPixel(pixel [3]float32) [3]float32 {
	return [3]float32{pixel[0] + 127, pixel[1] + 127, pixel[2] + 127}
}

func (c *Convolve) clampPixel(pixel [3]float32) [3]uint8 {
	return [3]uint8{
		uint8(math.Max(0, math.Min(float64(pixel[0]), 255))),
		uint8(math.Max(0, math.Min(float64(pixel[1]), 255))),
		uint8(math.Max(0, math.Min(float64(pixel[2]), 255))),
	}
}

func (c *Convolve) isValidKernel() bool {
	kernelHeight := len(c.Kernel)
	kernelWidth := len(c.Kernel[0])

	return kernelHeight%2 == 1 && kernelWidth%2 == 1
}

func (c *Convolve) isValidPixel(x, y, width, height int) bool {
	kernelHeight := len(c.Kernel)
	kernelWidth := len(c.Kernel[0])

	isOkVertically := y-kernelHeight/2 >= 0 && y+kernelHeight/2 < height
	isOkHorizontally := x-kernelWidth/2 >= 0 && x+kernelWidth/2 < width

	return isOkVertically && isOkHorizontally
}
