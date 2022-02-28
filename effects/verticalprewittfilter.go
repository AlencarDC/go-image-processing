package effects

import (
	"errors"
	"fpi/photochopp"
)

func VerticalPrewittKernel() [][]float32 {
	return [][]float32{
		{-1.0, -1.0, -1.0},
		{0.0, 0.0, 0.0},
		{1.0, 1.0, 1.0},
	}
}

type VerticalPrewittFilter struct{}

func (filter *VerticalPrewittFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply VerticalPrewittFilter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: VerticalPrewittKernel(), ShouldEmboss: true}
	convolve.Apply(img)

	return nil
}
