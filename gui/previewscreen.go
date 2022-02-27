package main

import (
	"fpi/photochopp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type PreviewScreen struct {
	window        fyne.Window
	image         photochopp.Image
	ctnMain       *fyne.Container
	applyCallback func(photochopp.Image)
}

func (ps *PreviewScreen) Content() fyne.CanvasObject {
	return ps.ctnMain
}

func (ps *PreviewScreen) SetCallback(callback func(photochopp.Image)) {
	ps.applyCallback = callback
}

func (ps *PreviewScreen) SetImage(img photochopp.Image) {
	ps.image = img

	cnvImage := canvas.NewImageFromImage(img.ImageFromRGBA())
	cnvImage.FillMode = canvas.ImageFillContain
	cnvImage.ScaleMode = canvas.ImageScaleFastest
	cnvImage.SetMinSize(fyne.Size{Width: float32(img.Width()), Height: float32(img.Height())})

	ps.ctnMain.Objects[0] = container.NewCenter(cnvImage)
}

func NewPreviewScreen(app App, window fyne.Window) *PreviewScreen {
	previewScreen := new(PreviewScreen)
	previewScreen.window = window

	btnApply := widget.NewButton("Apply", func() {
		if previewScreen.applyCallback != nil {
			previewScreen.applyCallback(previewScreen.image)
			previewScreen.window.Close()
		}
	})
	ctnButton := container.NewHBox(layout.NewSpacer(), btnApply)

	previewScreen.ctnMain = container.NewVBox(container.NewCenter(), ctnButton)

	return previewScreen
}
