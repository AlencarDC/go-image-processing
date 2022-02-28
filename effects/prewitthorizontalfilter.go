package effects

import (
	"errors"
	"fpi/photochopp"
)

func PrewittHorizontalKernel() [][]float32 {
	return [][]float32{
		{-1.0, 0.0, 1.0},
		{-1.0, 0.0, 1.0},
		{-1.0, 0.0, 1.0},
	}
}

type PrewittHorizontalFilter struct{}

func (filter *PrewittHorizontalFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply PrewittHorizontalFilter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: PrewittHorizontalKernel(), ShouldEmboss: true}
	convolve.Apply(img)

	return nil
}
