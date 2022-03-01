package main

import (
	"bytes"
	"fmt"
	"fpi/photochopp"
	"fpi/photochopp/effects"
	"fpi/photochopp/gui/component"
	"image/jpeg"
	"log"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

func (ms *MainScreen) Content() fyne.CanvasObject {
	return ms.ctnMain
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

func (ms *MainScreen) applyEffect(effect effects.Effect) {
	if ms.modifiedImage == nil {
		log.Println("apply-effect: can not apply effect to a nil image")
		return
	}

	err := effect.Apply(ms.modifiedImage)

	if err != nil {
		log.Println("apply-effect: error during effect processing", err)
		return
	}

	ms.updateModifiedImage(*ms.modifiedImage)
}

func (ms *MainScreen) updateModifiedImage(img photochopp.Image) {
	ms.cnvModifiedImage.Image = img.ImageFromRGBA()
	ms.cnvModifiedImage.SetMinSize(fyne.Size{Width: float32(ms.modifiedImage.Width()), Height: float32(ms.modifiedImage.Height())})
	ms.cnvModifiedImage.Refresh()
}

func (ms *MainScreen) saveModifiedImage() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, ms.modifiedImage.ImageFromRGBA(), nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewMainScreen(app App, window fyne.Window) *MainScreen {
	mainScreen := new(MainScreen)

	pnlOriginalImage := container.NewMax()
	pnlModifiedImage := container.NewMax()

	// SAVE IMAGE BUTTON
	dlgSaveImage := component.NewSaveImageDialog(window, mainScreen.saveModifiedImage)

	btnSaveModified := widget.NewButton("Save Image", func() {
		dlgSaveImage.Show()
	})

	// IMAGE LOAD BUTTON
	dlgImageLoad := component.NewLoadImageDialog(window, func(path string) {
		mainScreen.loadImage(path)

		_, filename := filepath.Split(path)
		dlgSaveImage.SetFileName("modified_" + filename)
	})

	btnFileLoad := widget.NewButton("Load image", func() {
		dlgImageLoad.Show()
	})
	pnlOriginalImage.Add(btnFileLoad)

	// HISTOGRAM BUTTON
	btnShowHistogram := widget.NewButton("Show Histogram", func() {
		window := app.NewWindow("histogram", "Histogram", 600, 300)
		screen := NewHistogramScreen(app, window)

		l := &effects.Luminance{}
		mainScreen.applyEffect(l)

		screen.PlotHistogram(*mainScreen.modifiedImage.Histogram())

		window.SetContent(screen.Content())
		window.Show()
	})

	// EFFECT BUTTONS
	btnVFlip := widget.NewButton("Vertical Flip", func() {
		vf := &effects.VerticalFlip{}
		mainScreen.applyEffect(vf)
	})

	btnHFlip := widget.NewButton("Horizontal Flip", func() {
		hf := &effects.HorizontalFlip{}
		mainScreen.applyEffect(hf)
	})

	btnGrayScale := widget.NewButton("Gray Scale (Luminance)", func() {
		l := &effects.Luminance{}
		mainScreen.applyEffect(l)
	})

	lblNumberOfColors := widget.NewLabel("Number of colors: 255")
	sliderNumberOfColors := widget.NewSlider(1, 255)
	sliderNumberOfColors.SetValue(255)
	sliderNumberOfColors.Step = 1
	sliderNumberOfColors.OnChanged = func(f float64) {
		lblNumberOfColors.SetText("Number of colors: " + strconv.Itoa(int(f)))
	}

	btnColorQuantization := widget.NewButton("Color Quantization", func() {
		nColors := int(sliderNumberOfColors.Value)
		cq := &effects.ColorQuantization{NumberOfDesiredColors: nColors}
		mainScreen.applyEffect(cq)
	})

	lblBrightnessValue := widget.NewLabel("Value: 0")
	sliderBrightnessValue := widget.NewSlider(-255, 255)
	sliderBrightnessValue.SetValue(0)
	sliderBrightnessValue.Step = 1
	sliderBrightnessValue.OnChanged = func(f float64) {
		lblBrightnessValue.SetText("Value: " + strconv.Itoa(int(f)))
	}

	btnBrightness := widget.NewButton("Brightness", func() {
		value := int(sliderBrightnessValue.Value)
		b := &effects.Brightness{Value: value}
		mainScreen.applyEffect(b)
	})

	lblContrastValue := widget.NewLabel("Value: 0")
	sliderContrastValue := widget.NewSlider(-255, 255)
	sliderContrastValue.SetValue(0)
	sliderContrastValue.Step = 1
	sliderContrastValue.OnChanged = func(f float64) {
		lblContrastValue.SetText("Value: " + strconv.Itoa(int(f)))
	}

	btnContrast := widget.NewButton("Contrast", func() {
		value := int(sliderContrastValue.Value)
		c := &effects.Contrast{Value: value}
		mainScreen.applyEffect(c)
	})

	btnNegative := widget.NewButton("Negative", func() {
		n := &effects.Negative{}
		mainScreen.applyEffect(n)
	})

	btnHistogramEqualization := widget.NewButton("Histogram Equalization", func() {
		heImage := mainScreen.modifiedImage.Copy()
		he := &effects.GrayScaleHistogramEqualization{}
		he.Apply(heImage)

		width, height := float32(heImage.Width()), float32(heImage.Height())

		window := app.NewWindow("equalization", "Histogram Equalization", width, height)
		screen := NewPreviewScreen(app, window)
		screen.SetCallback(mainScreen.updateModifiedImage)

		screen.SetImage(*heImage)

		window.SetContent(screen.Content())
		window.Show()
	})

	btnGaussianBlur := widget.NewButton("Gaussian Blur", func() {
		gb := &effects.GaussianBlur{}
		mainScreen.applyEffect(gb)
	})

	btnLaplacianFilter := widget.NewButton("Laplacian Filter", func() {
		filter := &effects.LaplacianFilter{}
		mainScreen.applyEffect(filter)
	})

	btnHighPassFilter := widget.NewButton("High Pass Filter", func() {
		filter := &effects.HighPassFilter{}
		mainScreen.applyEffect(filter)
	})

	btnHorizontalPrewittFilter := widget.NewButton("Prewitt Hx Filter", func() {
		filter := &effects.HorizontalPrewittFilter{}
		mainScreen.applyEffect(filter)
	})

	btnVerticalPrewittFilter := widget.NewButton("Prewitt Hy Filter", func() {
		filter := &effects.VerticalPrewittFilter{}
		mainScreen.applyEffect(filter)
	})

	btnHorizontalSobelFilter := widget.NewButton("Sobel Hx Filter", func() {
		filter := &effects.HorizontalSobelFilter{}
		mainScreen.applyEffect(filter)
	})

	btnVerticalSobelFilter := widget.NewButton("Sobel Hy Filter", func() {
		filter := &effects.VerticalSobelFilter{}
		mainScreen.applyEffect(filter)
	})

	btnRotateClockwiseFilter := widget.NewButton("Rotate 90° Clockwise", func() {
		filter := &effects.Rotation90Degree{RotateClockwise: true}
		mainScreen.applyEffect(filter)
	})

	btnRotateCounterClockwiseFilter := widget.NewButton("Rotate 90° Counterclockwise", func() {
		filter := &effects.Rotation90Degree{RotateClockwise: false}
		mainScreen.applyEffect(filter)
	})

	btnZoomOut := widget.NewButton("Zoom Out", func() {
		filter := &effects.ZoomOut{XFactor: 2, YFactor: 2}
		mainScreen.applyEffect(filter)
	})

	// MAIN CONTAINER
	pnlEffectButtons := container.New(layout.NewVBoxLayout(), btnHFlip, btnVFlip, btnGrayScale, lblNumberOfColors, sliderNumberOfColors, btnColorQuantization, btnShowHistogram, lblBrightnessValue, sliderBrightnessValue, btnBrightness, lblContrastValue, sliderContrastValue, btnContrast, btnNegative, btnHistogramEqualization, btnGaussianBlur, btnLaplacianFilter, btnHighPassFilter, btnHorizontalPrewittFilter, btnVerticalPrewittFilter, btnHorizontalSobelFilter, btnVerticalSobelFilter, btnRotateClockwiseFilter, btnRotateCounterClockwiseFilter, btnZoomOut, layout.NewSpacer(), btnSaveModified)

	mainScreen.originalImage = nil
	mainScreen.modifiedImage = nil
	mainScreen.cnvOriginalImage = nil
	mainScreen.cnvModifiedImage = nil
	mainScreen.pnlOriginalImage = pnlOriginalImage
	mainScreen.pnlModifiedImage = pnlModifiedImage

	mainScreen.ctnMain = container.NewBorder(nil, nil, pnlEffectButtons, nil, container.NewGridWithColumns(2, pnlOriginalImage, pnlModifiedImage))

	return mainScreen
}
