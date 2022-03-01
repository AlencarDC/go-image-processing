package main

import (
	"fpi/photochopp"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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

	gradient := canvas.NewHorizontalGradient(color.Black, color.White)
	gradientImage := canvas.NewImageFromImage(gradient.Generate(255*2, 8))
	gradientImage.FillMode = canvas.ImageFillOriginal
	gradientImage.ScaleMode = canvas.ImageScaleFastest
	gradientImage.SetMinSize(fyne.Size{Width: 255 * 2, Height: 8})

	lblBegin := widget.NewLabel("0")
	lblEnd := widget.NewLabel("255")
	ctnLabels := container.NewHBox(lblBegin, layout.NewSpacer(), lblEnd)

	hs.ctnMain = container.NewPadded(container.NewVBox(histogramBox, layout.NewSpacer(), gradientImage, ctnLabels))
}

func NewHistogramScreen(app App, window fyne.Window) *HistogramScreen {
	histogramScreen := new(HistogramScreen)
	histogramScreen.window = window

	histogramScreen.ctnMain = container.NewCenter()

	return histogramScreen
}
