package effects

import (
	"errors"
	"fpi/photochopp"
)

func HorizontalPrewittKernel() [][]float32 {
	return [][]float32{
		{-1.0, 0.0, 1.0},
		{-1.0, 0.0, 1.0},
		{-1.0, 0.0, 1.0},
	}
}

type HorizontalPrewittFilter struct{}

func (filter *HorizontalPrewittFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply HorizontalPrewittFilter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: HorizontalPrewittKernel(), ShouldEmboss: true}
	convolve.Apply(img)

	return nil
}
