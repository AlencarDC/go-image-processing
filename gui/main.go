package main

import (
	"fmt"
	"fpi/photochopp"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type MainScreen struct {
	originalImage    *photochopp.Image
	modifiedImage    *photochopp.Image
	cnvOriginalImage *canvas.Image
	cnvModifiedImage *canvas.Image
	pnlOriginalImage *fyne.Container
	pnlModifiedImage *fyne.Container
	ctnMain          *fyne.Container
}

func (ms *MainScreen) loadImage(path string) {
	var originalErr, modifiedErr error
	ms.originalImage, originalErr = photochopp.NewImageFromFilePath(path)
	ms.modifiedImage, modifiedErr = photochopp.NewImageFromFilePath(path)

	if originalErr != nil || modifiedErr != nil {
		fmt.Println("dialog: error on load image")
	}

	ms.cnvOriginalImage = canvas.NewImageFromImage(*ms.originalImage.Image())
	ms.cnvOriginalImage.FillMode = canvas.ImageFillContain
	ms.cnvOriginalImage.ScaleMode = canvas.ImageScaleFastest
	ms.cnvOriginalImage.SetMinSize(fyne.Size{Width: float32(ms.originalImage.Width()), Height: float32(ms.originalImage.Height())})

	ms.cnvModifiedImage = canvas.NewImageFromImage(*ms.modifiedImage.Image())
	ms.cnvModifiedImage.FillMode = canvas.ImageFillContain
	ms.cnvModifiedImage.ScaleMode = canvas.ImageScaleFastest
	ms.cnvModifiedImage.SetMinSize(fyne.Size{Width: float32(ms.modifiedImage.Width()), Height: float32(ms.modifiedImage.Height())})

	if len(ms.pnlOriginalImage.Objects) > 1 {
		ms.pnlOriginalImage.Objects = ms.pnlOriginalImage.Objects[0:1]
		ms.pnlModifiedImage.Objects = ms.pnlModifiedImage.Objects[0:1]
	}

	ms.pnlOriginalImage.Add(container.NewCenter(ms.cnvOriginalImage))
	ms.pnlModifiedImage.Add(container.NewCenter(ms.cnvModifiedImage))
}

func (ms *MainScreen) applyEffect(effect photochopp.Effect) {
	if ms.modifiedImage == nil {
		log.Println("apply-effect: can not apply effect to a nil image")
		return
	}

	err := effect.Apply(ms.modifiedImage)

	if err != nil {
		log.Println("apply-effect: error during effect processing")
	}

	ms.cnvModifiedImage.Image = ms.modifiedImage.ImageFromRGBA()
	ms.cnvModifiedImage.Refresh()
}

func NewMainScreen(window *fyne.Window) *MainScreen {
	mainScreen := new(MainScreen)

	pnlOriginalImage := container.NewMax()
	pnlModifiedImage := container.NewMax()

	fileLoadDialog := dialog.NewFileOpen(func(fileURI fyne.URIReadCloser, err error) {
		if fileURI == nil || err != nil {
			return
		}
		mainScreen.loadImage(fileURI.URI().Path())
	}, *window)
	fileLoadDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png"}))

	btnFileLoad := widget.NewButton("Load image", func() {
		fileLoadDialog.Show()
	})
	pnlOriginalImage.Add(btnFileLoad)

	// BUTTONS
	btnVFlip := widget.NewButton("Vertical Flip", func() {
		vf := &photochopp.VerticalFlip{}
		mainScreen.applyEffect(vf)
	})

	btnHFlip := widget.NewButton("Horizontal Flip", func() {
		hf := &photochopp.HorizontalFlip{}
		mainScreen.applyEffect(hf)
	})

	btnGrayScale := widget.NewButton("Gray Scale (Luminance)", func() {
		l := &photochopp.Luminance{}
		mainScreen.applyEffect(l)
	})

	pnlEffectButtons := container.New(layout.NewVBoxLayout(), btnHFlip, btnVFlip, btnGrayScale)

	mainScreen.originalImage = nil
	mainScreen.modifiedImage = nil
	mainScreen.cnvOriginalImage = nil
	mainScreen.cnvModifiedImage = nil
	mainScreen.pnlOriginalImage = pnlOriginalImage
	mainScreen.pnlModifiedImage = pnlModifiedImage

	mainScreen.ctnMain = container.NewBorder(nil, nil, pnlEffectButtons, nil, container.NewGridWithColumns(2, pnlOriginalImage, pnlModifiedImage))

	return mainScreen
}

func main() {

	application := app.New()
	window := application.NewWindow("Photochopp v1.0")
	window.Resize(fyne.NewSize(1280, 720))

	mainScreen := NewMainScreen(&window)

	window.SetContent(mainScreen.ctnMain)
	window.ShowAndRun()
}
