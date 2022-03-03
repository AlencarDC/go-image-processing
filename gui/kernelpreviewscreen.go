package main

import (
	"fmt"
	"fpi/photochopp/gui/component"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type KernelPreviewScreen struct {
	window               fyne.Window
	shouldEmboss         bool
	shouldApplyGrayScale bool
	inputs               [3][3]*component.NumericalEntry
	ctnInputs            [3]*fyne.Container
	ctnMain              *fyne.Container
	applyCallback        func([][]float32, bool, bool)
}

func (ps *KernelPreviewScreen) Content() fyne.CanvasObject {
	return ps.ctnMain
}

func (ps *KernelPreviewScreen) SetCallback(callback func([][]float32, bool, bool)) {
	ps.applyCallback = callback
}

func (ps *KernelPreviewScreen) Kernel() ([][]float32, error) {
	kernel := make([][]float32, len(ps.inputs))
	for i := range ps.inputs {
		kernel[i] = make([]float32, len(ps.inputs[i]))
		for j := range ps.inputs[i] {
			value, err := strconv.ParseFloat(ps.inputs[i][j].Text, 32)
			if err != nil {
				return nil, err
			}
			kernel[i][j] = float32(value)
		}
	}

	return kernel, nil
}

func NewKernelPreviewScreen(app App, window fyne.Window) *KernelPreviewScreen {
	screen := new(KernelPreviewScreen)
	screen.window = window

	btnApply := widget.NewButton("Convolve", func() {
		if screen.applyCallback != nil {
			kernel, err := screen.Kernel()
			if err != nil {
				return
			}

			screen.applyCallback(kernel, screen.shouldEmboss, screen.shouldApplyGrayScale)
			screen.window.Close()
		}
	})
	ctnButton := container.NewHBox(layout.NewSpacer(), btnApply)

	for i := 0; i < 3; i += 1 {
		screen.ctnInputs[i] = container.NewGridWithColumns(3)
		for j := 0; j < 3; j += 1 {
			screen.inputs[i][j] = component.NewNumericalEntry()
			screen.inputs[i][j].Text = fmt.Sprintf("%f", 0.0)
			screen.ctnInputs[i].Add(screen.inputs[i][j])
		}
	}
	ctnKernel := container.NewGridWithRows(3, screen.ctnInputs[0], screen.ctnInputs[1], screen.ctnInputs[2])

	checkEmboss := widget.NewCheck("Emboss result", func(checked bool) {
		screen.shouldEmboss = checked
	})
	checkGrayScale := widget.NewCheck("Apply Gray Scale", func(checked bool) {
		screen.shouldApplyGrayScale = checked
	})

	screen.ctnMain = container.NewVBox(ctnKernel, checkEmboss, checkGrayScale, layout.NewSpacer(), ctnButton)

	return screen
}
