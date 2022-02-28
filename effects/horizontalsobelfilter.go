package effects

import (
	"errors"
	"fpi/photochopp"
)

func HorizontalSobelKernel() [][]float32 {
	return [][]float32{
		{-1.0, 0.0, 1.0},
		{-2.0, 0.0, 2.0},
		{-1.0, 0.0, 1.0},
	}
}

type HorizontalSobelFilter struct{}

func (filter *HorizontalSobelFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply HorizontalSobelFilter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: HorizontalSobelKernel(), ShouldEmboss: true}
	convolve.Apply(img)

	return nil
}
