package effects

import (
	"errors"
	"fpi/photochopp"
)

type Rotation90Degree struct {
	RotateClockwise bool
}

func (r *Rotation90Degree) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply Rotation90Degree to a nil image")
	}

	imageCopy := img.Copy()
	ogWidth, ogHeight := imageCopy.Width(), imageCopy.Height()
	newWidth, newHeigth := ogHeight, ogWidth
	img.SetHeight(newHeigth)
	img.SetWidth(newWidth)

	for x := 0; x < newWidth; x++ {
		for y := 0; y < newHeigth; y++ {
			dstPixel := img.Pixel(x, y)

			var srcPixel []uint8
			if r.RotateClockwise {
				srcPixel = imageCopy.Pixel(y, ogHeight-1-x)
			} else {
				srcPixel = imageCopy.Pixel(ogWidth-1-y, x)
			}

			copy(dstPixel, srcPixel)
		}
	}
	return nil
}
