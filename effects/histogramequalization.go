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
	histogram := img.Histogram().B
	cumulativeHistogram := he.computeCumulativeHistogram(scalingFactor, histogram)

	for x := 0; x < img.Width(); x += 1 {
		for y := 0; y < img.Height(); y += 1 {
			pixel := img.Pixel(x, y)
			pixel[0] = uint8(cumulativeHistogram[pixel[0]]) // only valid because R G B has the same value in this case
			pixel[1] = uint8(cumulativeHistogram[pixel[1]])
			pixel[2] = uint8(cumulativeHistogram[pixel[2]])
		}
	}

	return nil
}

func (he *GrayScaleHistogramEqualization) renormalizedValue(scalingFactor, histogramValue float32) int32 {
	return int32(scalingFactor * histogramValue)
}

func (he *GrayScaleHistogramEqualization) computeCumulativeHistogram(scalingFactor float32, histogram [256]int32) [256]int32 {
	var cumulativeHistogram [256]int32
	cumulativeHistogram[0] = he.renormalizedValue(scalingFactor, float32(histogram[0]))

	for i := 1; i < len(histogram); i += 1 {
		renormalizedValue := he.renormalizedValue(scalingFactor, float32(histogram[i]))
		cumulativeHistogram[i] = cumulativeHistogram[i-1] + renormalizedValue
	}

	return cumulativeHistogram
}
