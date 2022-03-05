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
	"fyne.io/fyne/v2/theme"
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

	btnSaveModified := widget.NewButtonWithIcon("Save Image", theme.DocumentSaveIcon(), func() {
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
		l := &effects.Luminance{}
		mainScreen.applyEffect(l)
		histogram, _ := mainScreen.modifiedImage.Histogram().ForChannel(photochopp.GrayChannel)
		window := NewHistogramWindow("histogram", "Histogram", app, *histogram)
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
	ctnColorQuantization := container.NewVBox(lblNumberOfColors, sliderNumberOfColors, btnColorQuantization)

	lblBrightnessValue := widget.NewLabel("Brightness value: 0")
	sliderBrightnessValue := widget.NewSlider(-255, 255)
	sliderBrightnessValue.SetValue(0)
	sliderBrightnessValue.Step = 1
	sliderBrightnessValue.OnChanged = func(f float64) {
		lblBrightnessValue.SetText("Brightness value: " + strconv.Itoa(int(f)))
	}

	btnBrightness := widget.NewButton("Brightness", func() {
		value := int(sliderBrightnessValue.Value)
		b := &effects.Brightness{Value: value}
		mainScreen.applyEffect(b)
	})
	ctnBrightness := container.NewVBox(lblBrightnessValue, sliderBrightnessValue, btnBrightness)

	lblContrastValue := widget.NewLabel("Contrast value: 0")
	sliderContrastValue := widget.NewSlider(-255, 255)
	sliderContrastValue.SetValue(0)
	sliderContrastValue.Step = 1
	sliderContrastValue.OnChanged = func(f float64) {
		lblContrastValue.SetText("Contrast value: " + strconv.Itoa(int(f)))
	}

	btnContrast := widget.NewButton("Contrast", func() {
		value := int(sliderContrastValue.Value)
		c := &effects.Contrast{Value: value}
		mainScreen.applyEffect(c)
	})
	ctnContrast := container.NewVBox(lblContrastValue, sliderContrastValue, btnContrast)

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

		screen.SetSrcImage(*mainScreen.modifiedImage)
		screen.SetDstImage(*heImage)

		window.SetContent(screen.Content())
		window.Show()

		var histBefore, histAfter [256]int32
		if mainScreen.modifiedImage.IsGrayScale() {
			histBefore = mainScreen.modifiedImage.Histogram().B
			histAfter = heImage.Histogram().B
		} else {
			beforeImg := mainScreen.modifiedImage.Copy()
			afterImg := heImage.Copy()
			luminace := effects.Luminance{}
			luminace.Apply(beforeImg)
			luminace.Apply(afterImg)
			scalingFactor := 255.0 / float32((beforeImg.Height() * beforeImg.Width()))
			histBefore, _ = beforeImg.Histogram().CumulativeHistogram(photochopp.GrayChannel, scalingFactor)
			histAfter, _ = afterImg.Histogram().CumulativeHistogram(photochopp.GrayChannel, scalingFactor)
		}

		windowHistBefore := NewHistogramWindow("hist-before", "Histogram Before", app, histBefore)
		windowHistAfter := NewHistogramWindow("hist-after", "Histogram After", app, histAfter)
		windowHistBefore.Show()
		windowHistAfter.Show()
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

	btnCustomKernelFilter := widget.NewButton("Custom Filter", func() {
		window := app.NewWindow("custom-kernel", "Custom Kernel Filter", 300, 200)
		screen := NewKernelPreviewScreen(app, window)
		screen.SetCallback(func(kernel [][]float32, emboss, grayScale bool) {
			if grayScale {
				luminance := &effects.Luminance{}
				mainScreen.applyEffect(luminance)
			}

			filter := &effects.Convolve{Kernel: kernel, ShouldEmboss: emboss}
			mainScreen.applyEffect(filter)
		})

		window.SetContent(screen.Content())
		window.Show()
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

	btnRotateClockwise := widget.NewButton("Rotate 90° Clockwise", func() {
		filter := &effects.Rotation90Degree{RotateClockwise: true}
		mainScreen.applyEffect(filter)
	})

	btnRotateCounterClockwise := widget.NewButton("Rotate 90° Counterclockwise", func() {
		filter := &effects.Rotation90Degree{RotateClockwise: false}
		mainScreen.applyEffect(filter)
	})

	lblXFactor := widget.NewLabel("X:")
	entryXFactor := component.NewNumericalEntry()
	entryXFactor.Text = "1"
	lblYFactor := widget.NewLabel("Y:")
	entryYFactor := component.NewNumericalEntry()
	entryYFactor.Text = "1"

	btnZoomOut := widget.NewButton("Zoom Out", func() {
		xFactor, xErr := strconv.ParseInt(entryXFactor.Text, 10, 32)
		yFactor, yErr := strconv.ParseInt(entryYFactor.Text, 10, 32)
		if xErr != nil || yErr != nil {
			return
		}
		filter := &effects.ZoomOut{XFactor: int(xFactor), YFactor: int(yFactor)}
		mainScreen.applyEffect(filter)
	})
	ctnZoomOut := container.NewVBox(container.NewGridWithColumns(2, container.NewHBox(lblXFactor, entryXFactor), container.NewHBox(lblYFactor, entryYFactor)), btnZoomOut)

	btnZoomIn := widget.NewButton("Zoom In", func() {
		filter := &effects.ZoomIn{}
		mainScreen.applyEffect(filter)
	})

	dlgHistogramMatching := component.NewLoadImageDialog(window, func(path string) {
		img, err := photochopp.NewImageFromFilePath(path)
		if err != nil {
			log.Println("histogram-matching: unable to load image")
		}

		luminance := &effects.Luminance{}
		luminance.Apply(img)

		copyImage := mainScreen.modifiedImage.Copy()
		hm := &effects.HistogramMatching{Target: *img}
		hm.Apply(copyImage)

		width, height := float32(img.Width()), float32(img.Height())
		window := app.NewWindow("matching", "Histogram Matching", width, height)
		screen := NewPreviewScreen(app, window)
		screen.SetCallback(mainScreen.updateModifiedImage)

		screen.SetSrcImage(*img)
		screen.SetDstImage(*copyImage)

		window.SetContent(screen.Content())
		window.Show()
	})

	btnHistogramMatching := widget.NewButton("Histogram Matching", func() {
		dlgHistogramMatching.Show()
	})

	// BUTTON CONTAINERS
	lblTransform := widget.NewLabelWithStyle("Transform", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	ctnTransformButtons := container.NewVBox(lblTransform, btnVFlip, btnHFlip, btnRotateClockwise, btnRotateCounterClockwise, btnZoomIn, ctnZoomOut)
	lblAdjustments := widget.NewLabelWithStyle("Adjustments", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	ctnAdjustmentsButtons := container.NewVBox(lblAdjustments, btnShowHistogram, btnGrayScale, btnNegative, btnHistogramEqualization, btnHistogramMatching, ctnColorQuantization, ctnBrightness, ctnContrast)
	lblFilters := widget.NewLabelWithStyle("Filters", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	ctnFiltersButtons := container.NewVBox(lblFilters, btnGaussianBlur, btnLaplacianFilter, btnHighPassFilter, btnHorizontalPrewittFilter, btnVerticalPrewittFilter, btnHorizontalSobelFilter, btnVerticalSobelFilter, btnCustomKernelFilter)

	// MAIN CONTAINER
	pnlEffectButtons := container.NewVBox(ctnTransformButtons, ctnAdjustmentsButtons, ctnFiltersButtons)
	scrollButtons := container.NewVScroll(container.NewPadded(pnlEffectButtons))
	scrollButtons.SetMinSize(fyne.NewSize(0, 650))
	ctnButtons := container.NewVBox(scrollButtons, layout.NewSpacer(), btnSaveModified)

	mainScreen.originalImage = nil
	mainScreen.modifiedImage = nil
	mainScreen.cnvOriginalImage = nil
	mainScreen.cnvModifiedImage = nil
	mainScreen.pnlOriginalImage = pnlOriginalImage
	mainScreen.pnlModifiedImage = pnlModifiedImage

	mainScreen.ctnMain = container.NewBorder(nil, nil, ctnButtons, nil, container.NewGridWithColumns(2, pnlOriginalImage, pnlModifiedImage))

	return mainScreen
}
