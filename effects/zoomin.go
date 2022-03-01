package effects

import (
	"errors"
	"fpi/photochopp"
	"image"
)

type ZoomIn struct{}

func (z *ZoomIn) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply ZoomIn to a nil image")
	}

	const XFactor int = 2
	const YFactor int = 2

	ogWidth, ogHeight := img.Width(), img.Height()
	newWidth, newHeight := img.Width()*XFactor-1, img.Height()*YFactor-1

	newRGBA := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	newImage := photochopp.NewImage(newRGBA)

	for x := 0; x < ogWidth; x += 1 {
		for y := 0; y < ogHeight; y += 1 {
			srcPixel := img.Pixel(x, y)
			dstPixel := newImage.Pixel(x*XFactor, y*YFactor)

			copy(dstPixel, srcPixel)
		}
	}

	for x := 1; x < newWidth; x += 2 {
		for y := 0; y < newHeight; y += 2 {
			leftPixel := newImage.Pixel(x-1, y)
			rightPixel := newImage.Pixel(x+1, y)

			middlePixel := newImage.Pixel(x, y)
			copy(middlePixel, z.interpolatedPixel(leftPixel, rightPixel))
		}
	}

	for x := 0; x < newWidth; x += 1 {
		for y := 1; y < newHeight; y += 2 {
			upPixel := newImage.Pixel(x, y-1)
			downPixel := newImage.Pixel(x, y+1)

			middlePixel := newImage.Pixel(x, y)
			copy(middlePixel, z.interpolatedPixel(upPixel, downPixel))
		}
	}

	img.SetRGBA(newImage.RGBA())

	return nil
}

func (z *ZoomIn) interpolatedPixel(pixel1, pixel2 []uint8) []uint8 {
	sum := []int{
		int(pixel1[0]) + int(pixel2[0]),
		int(pixel1[1]) + int(pixel2[1]),
		int(pixel1[2]) + int(pixel2[2]),
	}

	return []uint8{
		uint8(sum[0] / 2),
		uint8(sum[1] / 2),
		uint8(sum[2] / 2),
		255,
	}
}
