package effects

import (
	"errors"
	"fpi/photochopp"
	"math"
)

type HistogramMatching struct {
	Target photochopp.Image
}

func (hm *HistogramMatching) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply histogram matching to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	targetAlpha := 255.0 / float32((hm.Target.Height() * hm.Target.Width()))
	targetCumulativeHistogram, err := hm.Target.Histogram().CumulativeHistogram(photochopp.BlueChannel, targetAlpha)
	if err != nil {
		return err
	}
	sourceAlpha := 255.0 / float32((img.Height() * img.Width()))
	sourceCumulativeHistogram, err := img.Histogram().CumulativeHistogram(photochopp.BlueChannel, sourceAlpha)
	if err != nil {
		return err
	}

	var HM [256]int32
	for i := 0; i < 256; i += 1 {
		histogramDiff := arrayAbsDiff(targetCumulativeHistogram[:], sourceCumulativeHistogram[i])
		indexOfMinValues := arrayIndexOfMinValues(histogramDiff)

		indexDiff := arrayAbsDiff(indexOfMinValues, int32(i))
		minIndex := arrayIndexOfMinValue(indexDiff)
		closestColor := indexOfMinValues[minIndex]

		HM[i] = closestColor
	}

	for x := 0; x < img.Width(); x += 1 {
		for y := 0; y < img.Height(); y += 1 {
			pixel := img.Pixel(x, y)
			pixel[0] = uint8(HM[pixel[0]])
			pixel[1] = uint8(HM[pixel[1]])
			pixel[2] = uint8(HM[pixel[2]])
		}
	}
	return nil
}

func arrayAbsDiff(array []int32, subValue int32) []int32 {
	var absDiffArray []int32
	for j := range array {
		value := int32(math.Abs(float64(subValue - array[j])))
		absDiffArray = append(absDiffArray, value)
	}

	return absDiffArray
}

func arrayIndexOfMinValues(array []int32) []int32 {
	var minIndexes []int32
	minValue := array[0]
	for j := range array {
		if array[j] == minValue {
			minIndexes = append(minIndexes, int32(j))
		} else if array[j] < minValue {
			minIndexes = nil
			minIndexes = append(minIndexes, int32(j))
			minValue = array[j]
		}
	}

	return minIndexes
}

func arrayIndexOfMinValue(array []int32) int32 {
	var index int32
	minValue := array[0]
	for j := range array {
		if array[j] < minValue {
			index = int32(j)
			minValue = array[j]
		}
	}

	return index
}
