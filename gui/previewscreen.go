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
	srcImage      photochopp.Image
	dstImage      photochopp.Image
	cntImages     *fyne.Container
	ctnMain       *fyne.Container
	applyCallback func(photochopp.Image)
}

func (ps *PreviewScreen) Content() fyne.CanvasObject {
	return ps.ctnMain
}

func (ps *PreviewScreen) SetCallback(callback func(photochopp.Image)) {
	ps.applyCallback = callback
}

func (ps *PreviewScreen) SetSrcImage(img photochopp.Image) {
	ps.srcImage = img
	ps.update()
}

func (ps *PreviewScreen) SetDstImage(img photochopp.Image) {
	ps.dstImage = img
	ps.update()
}

func (ps *PreviewScreen) update() {
	ps.cntImages = container.NewGridWithColumns(1)
	if !ps.srcImage.IsEmpty() && !ps.dstImage.IsEmpty() {
		ps.cntImages = container.NewGridWithColumns(2)
	}

	if !ps.srcImage.IsEmpty() {
		cnvImage := ps.createCanvasImage(ps.srcImage)
		ps.cntImages.Add(container.NewCenter(cnvImage))
	}
	if !ps.dstImage.IsEmpty() {
		cnvImage := ps.createCanvasImage(ps.dstImage)
		ps.cntImages.Add(container.NewCenter(cnvImage))
	}
	ps.ctnMain.Objects[0] = ps.cntImages
}

func (ps *PreviewScreen) createCanvasImage(img photochopp.Image) *canvas.Image {
	size := fyne.Size{Width: float32(img.Width()), Height: float32(img.Height())}
	cnvImage := canvas.NewImageFromImage(img.ImageFromRGBA())
	cnvImage.FillMode = canvas.ImageFillOriginal
	cnvImage.ScaleMode = canvas.ImageScaleFastest
	cnvImage.SetMinSize(size)
	cnvImage.Resize(size)

	return cnvImage
}

func NewPreviewScreen(app App, window fyne.Window) *PreviewScreen {
	previewScreen := new(PreviewScreen)
	previewScreen.window = window

	btnApply := widget.NewButton("Apply", func() {
		if previewScreen.applyCallback != nil {
			previewScreen.applyCallback(previewScreen.dstImage)
			previewScreen.window.Close()
		}
	})
	ctnButton := container.NewHBox(layout.NewSpacer(), btnApply)
	ctnImages := container.NewGridWithColumns(2)

	previewScreen.ctnMain = container.NewVBox(ctnImages, ctnButton)

	return previewScreen
}
