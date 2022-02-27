package main

import (
	"fpi/photochopp"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type HistogramScreen struct {
	window  fyne.Window
	ctnMain *fyne.Container
}

func (hs *HistogramScreen) Content() fyne.CanvasObject {
	return hs.ctnMain
}

func (hs *HistogramScreen) PlotHistogram(histogram photochopp.Histogram) {
	MAX_HEIGHT := float32(250.0)
	getLineHeight := func(pixelsCount, maxPixelsCount int32) float32 {
		return MAX_HEIGHT * (float32(pixelsCount) / float32(maxPixelsCount))
	}

	maxPixelsCount, _ := histogram.MaxPixelsCount(photochopp.BlueChannel)

	histogramBox := container.NewWithoutLayout()
	for i, pixelsCount := range histogram.B {
		lineHeight := getLineHeight(pixelsCount, maxPixelsCount)

		line := canvas.NewLine(color.White)
		line.StrokeWidth = 2
		line.Position1 = fyne.Position{X: float32(i) * 2, Y: MAX_HEIGHT - float32(lineHeight)}
		line.Position2 = fyne.Position{X: float32(i) * 2, Y: MAX_HEIGHT}

		histogramBox.Add(line)
	}

	hs.ctnMain = container.NewMax(histogramBox)
}

func NewHistogramScreen(app App, window fyne.Window) *HistogramScreen {
	histogramScreen := new(HistogramScreen)
	histogramScreen.window = window

	histogramScreen.ctnMain = container.NewCenter()

	return histogramScreen
}
