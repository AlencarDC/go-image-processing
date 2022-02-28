package effects

import (
	"errors"
	"fpi/photochopp"
)

func HighPassKernel() [][]float32 {
	return [][]float32{
		{-1.0, -1.0, -1.0},
		{-1.0, 8.0, -1.0},
		{-1.0, -1.0, -1.0},
	}
}

type HighPassFilter struct{}

func (filter *HighPassFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply High Pass Filter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: HighPassKernel(), ShouldEmboss: false}
	convolve.Apply(img)

	return nil
}
