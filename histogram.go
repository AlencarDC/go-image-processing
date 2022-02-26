package photochopp

import "errors"

type Histogram struct {
	R [256]int32
	G [256]int32
	B [256]int32
}

func (h *Histogram) ForChannel(channel ColorChannel) (*[256]int32, error) {
	switch channel {
	case RedChannel:
		return &h.R, nil
	case GreenChannel:
		return &h.G, nil
	case BlueChannel:
		return &h.B, nil
	default:
		return nil, errors.New("histogram: invalid color channel")
	}
}

func (h *Histogram) HighestColor(channel ColorChannel) (uint8, error) {
	hChannel, err := h.ForChannel(channel)
	if err != nil {
		return 0, err
	}

	color := 256 - 1
	for color >= 0 && hChannel[color] == 0 {
		color -= 1
	}

	return uint8(color), nil
}

func (h *Histogram) LowestColor(channel ColorChannel) (uint8, error) {
	hChannel, err := h.ForChannel(channel)
	if err != nil {
		return 0, err
	}

	color := 0
	for color < 256 && hChannel[color] == 0 {
		color += 1
	}

	return uint8(color), nil
}

func (h *Histogram) MaxPixelsCount(channel ColorChannel) (int32, error) {
	hChannel, err := h.ForChannel(channel)
	if err != nil {
		return 0, err
	}

	var maxPixelsCount int32
	for i, pixelsCount := range hChannel {
		if i == 0 || pixelsCount > maxPixelsCount {
			maxPixelsCount = pixelsCount
		}
	}

	return maxPixelsCount, nil
}

func NewHistogram(img *Image) (*Histogram, error) {
	if img == nil {
		return nil, errors.New("histogram: cannot create from nil image")
	}

	width, height := img.Width(), img.Height()

	histogram := new(Histogram)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := img.Pixel(x, y)
			histogram.R[pixel[0]] += 1
			histogram.G[pixel[1]] += 1
			histogram.B[pixel[2]] += 1
		}
	}

	return histogram, nil
}
