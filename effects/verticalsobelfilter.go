package effects

import (
	"errors"
	"fpi/photochopp"
)

func VerticalSobelKernel() [][]float32 {
	return [][]float32{
		{-1.0, -2.0, -1.0},
		{0.0, 0.0, 0.0},
		{1.0, 2.0, 1.0},
	}
}

type VerticalSobelFilter struct{}

func (filter *VerticalSobelFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply VerticalSobelFilter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: VerticalSobelKernel(), ShouldEmboss: true}
	convolve.Apply(img)

	return nil
}
