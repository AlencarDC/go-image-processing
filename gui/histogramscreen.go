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

	hLine := canvas.NewLine(color.White)
	hLine.StrokeWidth = 2
	hLine.Position1 = fyne.Position{X: 0, Y: MAX_HEIGHT + 4}
	hLine.Position2 = fyne.Position{X: 255 * 2, Y: MAX_HEIGHT + 4}
	histogramBox.Add(hLine)

	lBegin := canvas.NewText("0", color.White)
	histogramBox.Add(lBegin)
	lBegin.Move(fyne.Position{X: 0, Y: MAX_HEIGHT + 8})

	lEnd := canvas.NewText("255", color.White)
	histogramBox.Add(lEnd)
	lEnd.Move(fyne.Position{255 * 2, MAX_HEIGHT + 8})

	hs.ctnMain = container.NewMax(histogramBox)
}

func NewHistogramScreen(app App, window fyne.Window) *HistogramScreen {
	histogramScreen := new(HistogramScreen)
	histogramScreen.window = window

	histogramScreen.ctnMain = container.NewCenter()

	return histogramScreen
}
