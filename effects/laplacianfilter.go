package effects

import (
	"errors"
	"fpi/photochopp"
)

func LaplacianKernel() [][]float32 {
	return [][]float32{
		{0.0, -1.0, 0.0},
		{-1.0, 4.0, -1.0},
		{0.0, -1.0, 0.0},
	}
}

type LaplacianFilter struct{}

func (filter *LaplacianFilter) Apply(img *photochopp.Image) (err error) {
	if img == nil {
		return errors.New("effect: cannot apply Laplacian Filter to a nil image")
	}

	luminance := Luminance{}
	luminance.Apply(img)

	convolve := Convolve{Kernel: LaplacianKernel(), ShouldEmboss: true}
	convolve.Apply(img)

	return nil
}
