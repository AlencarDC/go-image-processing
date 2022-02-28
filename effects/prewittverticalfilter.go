package effects

import (
	"errors"
	"fpi/photochopp"
)

func PrewittVerticalKernel() [][]float32 {
	return [][]float32{
		{-1.0, -1.0, -1.0},
		{0.0, 0.0, 0.0},
		{1.0, 1.0, 1.0},
	}
}

type PrewittVerticalFilter struct{}

func (filter *PrewittVerticalFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply PrewittVerticalFilter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: PrewittVerticalKernel(), ShouldEmboss: true}
	convolve.Apply(img)

	return nil
}
