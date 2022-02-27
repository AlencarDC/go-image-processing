package effects

import (
	"errors"
	"fpi/photochopp"
)

func GaussianKernel() [][]float32 {
	return [][]float32{
		{0.0625, 0.125, 0.0625},
		{0.125, 0.25, 0.125},
		{0.0625, 0.125, 0.0625},
	}
}

type GaussianBlur struct{}

func (gb *GaussianBlur) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply Gaussian Blur to a nil image")
	}

	convolve := Convolve{Kernel: GaussianKernel(), Clampping: false}
	convolve.Apply(img)

	return nil
}
