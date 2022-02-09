package main

import (
	"bytes"
	"fmt"
	"fpi/photochopp"
	"image"
	"image/jpeg"
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
	originalImage *canvas.Image
	modifiedImage *canvas.Image
	mainContainer *fyne.Container
}

func NewMainScreen() *MainScreen {
	return nil
}

func main() {

	application := app.New()
	window := application.NewWindow("Container")
	window.Resize(fyne.NewSize(1280, 720))

	// begin := canvas.NewText("(begin)", color.White)
	// end := canvas.NewText("(end)", color.White)

	rightPanel2 := container.NewMax()
	rightPanel := container.NewMax()

	var originalImage, modifiedImage *photochopp.Image
	var cnvOriginalImage, cnvModifiedImage *canvas.Image
	fileLoadDialog := dialog.NewFileOpen(func(fileURI fyne.URIReadCloser, err error) {
		if fileURI == nil || err != nil {
			return
		}
		fmt.Println("Arquivo selecionado", fileURI.URI())
		var originalErr, modifiedErr error
		originalImage, originalErr = photochopp.NewImageFromFilePath(fileURI.URI().Path())
		modifiedImage, modifiedErr = photochopp.NewImageFromFilePath(fileURI.URI().Path())
		if originalErr != nil || modifiedErr != nil {
			fmt.Println("dialog: error on load image")
		}

		cnvOriginalImage = canvas.NewImageFromImage(*originalImage.Image())
		cnvOriginalImage.FillMode = canvas.ImageFillContain
		cnvOriginalImage.ScaleMode = canvas.ImageScaleFastest
		cnvOriginalImage.SetMinSize(fyne.Size{Width: float32(originalImage.Width()), Height: float32(originalImage.Height())})

		cnvModifiedImage = canvas.NewImageFromImage(*modifiedImage.Image())
		cnvModifiedImage.FillMode = canvas.ImageFillContain
		cnvModifiedImage.ScaleMode = canvas.ImageScaleFastest
		cnvModifiedImage.SetMinSize(fyne.Size{Width: float32(modifiedImage.Width()), Height: float32(modifiedImage.Height())})

		if len(rightPanel.Objects) > 1 {
			rightPanel.Objects = rightPanel.Objects[0:1]
			rightPanel2.Objects = rightPanel2.Objects[0:1]
		}

		rightPanel.Add(container.NewCenter(cnvOriginalImage))
		rightPanel2.Add(container.NewCenter(cnvModifiedImage))
	}, window)
	fileLoadDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png"}))

	btnFileLoad := widget.NewButton("Load image", func() {
		fileLoadDialog.Show()
	})
	rightPanel.Add(btnFileLoad)

	// < BUTTONS -----------------------------------------
	btnHFlip := widget.NewButton("Horizontal Flip", func() {
		log.Println("Flip Horizontally")
		if modifiedImage == nil {
			log.Println("horizontal-flip: can not flip a nil image")
			return
		}
		hf := &photochopp.HorizontalFlip{}
		hf.Apply(modifiedImage)

		buf := new(bytes.Buffer)
		rgba := modifiedImage.RGBA()
		err := jpeg.Encode(buf, rgba, nil)
		if err != nil {
			log.Println("convert: error to convert image to buffer")
			return
		}
		newModifiedImage, _, err := image.Decode(buf)
		if err != nil {
			log.Println("convert: error to convert buffer to image")
			return
		}

		cnvModifiedImage.Image = newModifiedImage
		cnvModifiedImage.Refresh()
	})

	btnVFlip := widget.NewButton("Vertical Flip", func() {
		log.Println("Flip Vertically Start")
		if modifiedImage == nil {
			log.Println("vertical-flip: can not flip a nil image")
			return
		}
		vf := &photochopp.VerticalFlip{}
		vf.Apply(modifiedImage)

		buf := new(bytes.Buffer)
		rgba := modifiedImage.RGBA()
		err := jpeg.Encode(buf, rgba, nil)
		if err != nil {
			log.Println("convert: error to convert image to buffer")
			return
		}
		newModifiedImage, _, err := image.Decode(buf)
		if err != nil {
			log.Println("convert: error to convert buffer to image")
			return

		}

		cnvModifiedImage.Image = newModifiedImage
		cnvModifiedImage.Refresh()
		log.Println("Flip Vertically End")
	})

	leftPanel := container.New(layout.NewVBoxLayout(), btnHFlip, btnVFlip)
	// > BUTTONS -----------------------------------------

	window.SetContent(container.NewBorder(nil, nil, leftPanel, nil, container.NewGridWithColumns(2, rightPanel, rightPanel2)))
	window.ShowAndRun()
}
