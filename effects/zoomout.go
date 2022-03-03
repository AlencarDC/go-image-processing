package effects

import (
	"errors"
	"fpi/photochopp"
	"image"
)

type ZoomOut struct {
	XFactor int
	YFactor int
}

func (z *ZoomOut) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply ZoomOut to a nil image")
	}

	if z.XFactor < 1 || z.YFactor < 1 {
		return errors.New("effect: zoom out factors need to be greater than or equal to 1")
	}

	newWidth, newHeight := img.Width()/z.XFactor, img.Height()/z.YFactor

	newRGBA := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	newImage := photochopp.NewImage(newRGBA)

	for x := 0; x < newWidth; x += 1 {
		for y := 0; y < newHeight; y += 1 {
			rect := image.Rect(x*z.XFactor, y*z.YFactor, (x+1)*z.XFactor, (y+1)*z.YFactor)

			avgPixel := z.calculateAvgPixel(img, rect)

			pixel := newImage.Pixel(x, y)
			copy(pixel, avgPixel)
		}
	}

	img.SetRGBA(newImage.RGBA())

	return nil
}

func (z *ZoomOut) calculateAvgPixel(img *photochopp.Image, rect image.Rectangle) []uint8 {
	sum := []int{0, 0, 0}
	for i := rect.Min.X; i < rect.Max.X; i += 1 {
		for j := rect.Min.Y; j < rect.Max.Y; j += 1 {
			if img.IsValidPosition(i, j) {
				pixel := img.Pixel(i, j)
				sum[0] += int(pixel[0])
				sum[1] += int(pixel[1])
				sum[2] += int(pixel[2])
			}
		}
	}

	return []uint8{
		uint8(sum[0] / (z.XFactor * z.YFactor)),
		uint8(sum[1] / (z.XFactor * z.YFactor)),
		uint8(sum[2] / (z.XFactor * z.YFactor)),
		255,
	}
}
