package effects

import (
	"errors"
	"fpi/photochopp"
)

type GrayScaleHistogramEqualization struct{}

func (he *GrayScaleHistogramEqualization) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply histogram equalization to a nil image")
	}

	scalingFactor := 255.0 / float32((img.Height() * img.Width()))

	var cumulativeHistogram [256]int32
	if img.IsGrayScale() {
		cumulativeHistogram, err = img.Histogram().CumulativeHistogram(photochopp.GrayChannel, scalingFactor)
	} else {
		copyImg := img.Copy()
		luminance := Luminance{}
		luminance.Apply(copyImg)

		cumulativeHistogram, err = copyImg.Histogram().CumulativeHistogram(photochopp.GrayChannel, scalingFactor)
	}

	if err != nil {
		return err
	}

	for x := 0; x < img.Width(); x += 1 {
		for y := 0; y < img.Height(); y += 1 {
			pixel := img.Pixel(x, y)
			pixel[0] = uint8(cumulativeHistogram[pixel[0]])
			pixel[1] = uint8(cumulativeHistogram[pixel[1]])
			pixel[2] = uint8(cumulativeHistogram[pixel[2]])
		}
	}

	return nil
}
