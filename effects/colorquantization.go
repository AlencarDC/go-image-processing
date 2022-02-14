package effects

import (
	"errors"
	"fpi/photochopp"
)

type ColorQuantization struct {
	NumberOfDesiredColors int
}

func (cq *ColorQuantization) binSizes(histogram *photochopp.Histogram) []int {
	channels := [...]photochopp.ColorChannel{photochopp.RedChannel, photochopp.GreenChannel, photochopp.BlueChannel}
	var binSizes []int
	for _, channel := range channels {
		lowest, _ := histogram.LowestColor(channel)
		highest, _ := histogram.HighestColor(channel)

		intensitySize := int(highest) - int(lowest) + 1

		binSize := intensitySize / cq.NumberOfDesiredColors
		binSizes = append(binSizes, binSize)
	}

	return binSizes
}

func (cq *ColorQuantization) newColor(currentColor uint8, binSize int) uint8 {
	newBinBegin := (int(currentColor) / binSize) * binSize

	return uint8(newBinBegin + binSize/2)
}

func (cq *ColorQuantization) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply color quantization to a nil image")
	}

	histogram := img.Histogram()
	binSizes := cq.binSizes(histogram)

	width, height := img.Width(), img.Height()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.Pixel(x, y)
			pixel[0] = cq.newColor(pixel[0], binSizes[0])
			pixel[1] = cq.newColor(pixel[1], binSizes[1])
			pixel[2] = cq.newColor(pixel[2], binSizes[2])
		}
	}

	return nil
}
