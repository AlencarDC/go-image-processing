package effects

import "fpi/photochopp"

type Effect interface {
	Apply(img *photochopp.Image) (err error)
}
